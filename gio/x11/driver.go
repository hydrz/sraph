package x11

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/render"
	"github.com/jezek/xgb/shm"
	"github.com/jezek/xgb/xproto"
	"github.com/opensraph/sraph/gio"
	"github.com/opensraph/sraph/gio/internal/x11key"
)

var _ gio.Driver = (*x11Driver)(nil)

func init() {
	runtime.LockOSThread()
	gio.RegisterDriver(gio.DriverTypeX11, newX11Driver)
}

type x11Driver struct {
	xc      *xgb.Conn
	xsi     *xproto.ScreenInfo
	keysyms x11key.KeysymTable

	atomNETWMName      xproto.Atom
	atomUTF8String     xproto.Atom
	atomWMDeleteWindow xproto.Atom
	atomWMProtocols    xproto.Atom
	atomWMTakeFocus    xproto.Atom

	pixelsPerPt  float32
	pictformat24 render.Pictformat
	pictformat32 render.Pictformat

	// window32 and its related X11 resources is an unmapped window so that we
	// have a depth-32 window to create depth-32 pixmaps from, i.e. pixmaps
	// with an alpha channel. The root window isn't guaranteed to be depth-32.
	gcontext32 xproto.Gcontext
	window32   xproto.Window

	// opaqueP is a fully opaque, solid fill picture.
	opaqueP render.Picture
	useShm  bool

	uniformMu sync.Mutex
	uniformC  render.Color
	uniformP  render.Picture

	mu              sync.Mutex
	windows         map[xproto.Window]*x11Window
	nPendingUploads int
	completionKeys  []uint16
}

func newX11Driver() (driver gio.Driver, retError error) {
	xc, err := xgb.NewConn()
	if err != nil {
		return nil, fmt.Errorf("x11driver: xgb.NewConn failed: %v", err)
	}
	defer func() {
		if retError != nil {
			xc.Close()
		}
	}()

	if err := render.Init(xc); err != nil {
		return nil, fmt.Errorf("x11driver: render.Init failed: %v", err)
	}
	useShm := true
	if err := shm.Init(xc); err != nil {
		useShm = false
	}

	s := &x11Driver{
		xc:      xc,
		xsi:     xproto.Setup(xc).DefaultScreen(xc),
		windows: map[xproto.Window]*x11Window{},
		useShm:  useShm,
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("x11driver: init failed: %v", err)
	}

	return s, nil
}

// CreateWindow implements gio.Driver.
func (x *x11Driver) CreateWindow(o ...gio.NewWindowOptions) (gio.Window, error) {
	panic("unimplemented")
}

// Type implements gio.Driver.
func (x *x11Driver) Type() gio.DriverType {
	panic("unimplemented")
}

func (x *x11Driver) init() error {
	if err := x.initAtoms(); err != nil {
		return err
	}

	if err := x.initKeyboardMapping(); err != nil {
		return err
	}

	go x.run()

	return nil
}

func (s *x11Driver) initAtoms() (err error) {
	s.atomNETWMName, err = s.internAtom("_NET_WM_NAME")
	if err != nil {
		return err
	}
	s.atomUTF8String, err = s.internAtom("UTF8_STRING")
	if err != nil {
		return err
	}
	s.atomWMDeleteWindow, err = s.internAtom("WM_DELETE_WINDOW")
	if err != nil {
		return err
	}
	s.atomWMProtocols, err = s.internAtom("WM_PROTOCOLS")
	if err != nil {
		return err
	}
	s.atomWMTakeFocus, err = s.internAtom("WM_TAKE_FOCUS")
	if err != nil {
		return err
	}
	return nil
}

func (s *x11Driver) internAtom(name string) (xproto.Atom, error) {
	r, err := xproto.InternAtom(s.xc, false, uint16(len(name)), name).Reply()
	if err != nil {
		return 0, fmt.Errorf("x11driver: xproto.InternAtom failed: %v", err)
	}
	if r == nil {
		return 0, fmt.Errorf("x11driver: xproto.InternAtom failed")
	}
	return r.Atom, nil
}

func (s *x11Driver) initKeyboardMapping() error {
	const keyLo, keyHi = 8, 255
	km, err := xproto.GetKeyboardMapping(s.xc, keyLo, keyHi-keyLo+1).Reply()
	if err != nil {
		return err
	}
	n := int(km.KeysymsPerKeycode)
	if n < 2 {
		return fmt.Errorf("x11driver: too few keysyms per keycode: %d", n)
	}
	for i := keyLo; i <= keyHi; i++ {
		for j := 0; j < 6; j++ {
			if j < n {
				s.keysyms.Table[i][j] = uint32(km.Keysyms[(i-keyLo)*n+j])
			} else {
				s.keysyms.Table[i][j] = 0
			}
		}
	}

	// Figure out which modifier is the numlock modifier (see chapter 12.7 of the XLib Manual).
	mm, err := xproto.GetModifierMapping(s.xc).Reply()
	if err != nil {
		return err
	}
	s.keysyms.NumLockMod, s.keysyms.ModeSwitchMod, s.keysyms.ISOLevel3ShiftMod = 0, 0, 0
	numLockFound, modeSwitchFound, isoLevel3ShiftFound := false, false, false
modifierSearchLoop:
	for modifier := 0; modifier < 8; modifier++ {
		for i := 0; i < int(mm.KeycodesPerModifier); i++ {
			const (
				// XK_Num_Lock, XK_Mode_switch and XK_ISO_Level3_Shift from /usr/include/X11/keysymdef.h.
				xkNumLock        = 0xff7f
				xkModeSwitch     = 0xff7e
				xkISOLevel3Shift = 0xfe03
			)
			switch s.keysyms.Table[mm.Keycodes[modifier*int(mm.KeycodesPerModifier)+i]][0] {
			case xkNumLock:
				s.keysyms.NumLockMod = 1 << uint(modifier)
				numLockFound = true
			case xkModeSwitch:
				s.keysyms.ModeSwitchMod = 1 << uint(modifier)
				modeSwitchFound = true
			case xkISOLevel3Shift:
				s.keysyms.ISOLevel3ShiftMod = 1 << uint(modifier)
				isoLevel3ShiftFound = true
			}
			if numLockFound && modeSwitchFound && isoLevel3ShiftFound {
				break modifierSearchLoop
			}
		}
	}
	return nil
}

func (s *x11Driver) run() {
	keyboardChanged := false
	for {
		ev, err := s.xc.WaitForEvent()
		if err != nil {
			log.Printf("x11driver: xproto.WaitForEvent: %v", err)
			continue
		}

		noWindowFound := false
		switch ev := ev.(type) {
		case xproto.DestroyNotifyEvent:
			s.mu.Lock()
			delete(s.windows, ev.Window)
			s.mu.Unlock()

		case shm.CompletionEvent:
			s.mu.Lock()
			s.completionKeys = append(s.completionKeys, ev.Sequence)
			s.mu.Unlock()

		case xproto.ClientMessageEvent:
			if ev.Type != s.atomWMProtocols || ev.Format != 32 {
				break
			}
			switch xproto.Atom(ev.Data.Data32[0]) {
			case s.atomWMDeleteWindow:
				if w := s.findWindow(ev.Window); w != nil {

				} else {
					noWindowFound = true
				}
			case s.atomWMTakeFocus:
				xproto.SetInputFocus(s.xc, xproto.InputFocusParent, ev.Window, xproto.Timestamp(ev.Data.Data32[1]))
			}

		case xproto.ConfigureNotifyEvent:
			if w := s.findWindow(ev.Window); w != nil {
			} else {
				noWindowFound = true
			}

		case xproto.ExposeEvent:
			if w := s.findWindow(ev.Window); w != nil {
				// A non-zero Count means that there are more expose events
				// coming. For example, a non-rectangular exposure (e.g. from a
				// partially overlapped window) will result in multiple expose
				// events whose dirty rectangles combine to define the dirty
				// region. Go's paint events do not provide dirty regions, so
				// we only pass on the final X11 expose event.
				if ev.Count == 0 {
				}
			} else {
				noWindowFound = true
			}

		case xproto.FocusInEvent:
			if w := s.findWindow(ev.Event); w != nil {
			} else {
				noWindowFound = true
			}

		case xproto.FocusOutEvent:
			if w := s.findWindow(ev.Event); w != nil {
			} else {
				noWindowFound = true
			}

		case xproto.KeyPressEvent:
			if keyboardChanged {
				keyboardChanged = false
				s.initKeyboardMapping()
			}
			if w := s.findWindow(ev.Event); w != nil {
			} else {
				noWindowFound = true
			}

		case xproto.KeyReleaseEvent:
			if keyboardChanged {
				keyboardChanged = false
				s.initKeyboardMapping()
			}
			if w := s.findWindow(ev.Event); w != nil {
			} else {
				noWindowFound = true
			}

		case xproto.ButtonPressEvent:
			if w := s.findWindow(ev.Event); w != nil {

			} else {
				noWindowFound = true
			}

		case xproto.ButtonReleaseEvent:
			if w := s.findWindow(ev.Event); w != nil {

			} else {
				noWindowFound = true
			}

		case xproto.MotionNotifyEvent:
			if w := s.findWindow(ev.Event); w != nil {
			} else {
				noWindowFound = true
			}

		case xproto.MappingNotifyEvent:
			if ev.Request == xproto.MappingModifier || ev.Request == xproto.MappingKeyboard {
				keyboardChanged = true
			}
		}

		if noWindowFound {
			log.Printf("x11driver: no window found for event %T", ev)
		}
	}
}

func (s *x11Driver) findWindow(key xproto.Window) *x11Window {
	s.mu.Lock()
	w := s.windows[key]
	s.mu.Unlock()
	return w
}
