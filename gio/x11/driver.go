package x11

import (
	"fmt"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/render"
	"github.com/jezek/xgb/xproto"
	"github.com/opensraph/sraph/gio"
	"github.com/opensraph/sraph/gio/internal/x11key"
)

var _ gio.Driver = (*x11Driver)(nil)

func init() {
	gio.RegisterDriver(gio.DriverTypeX11, newX11Driver)
}

type x11Driver struct {
	xc      *xgb.Conn
	xsi     *xproto.ScreenInfo
	keysyms x11key.KeysymTable

	atomNETWMName               xproto.Atom
	atomUTF8String              xproto.Atom
	atomWMDeleteWindow          xproto.Atom
	atomWMProtocols             xproto.Atom
	atomWMTakeFocus             xproto.Atom
	atomWMChangeState           xproto.Atom
	atomNETWMState              xproto.Atom
	atomNETWMStateMaximizedHorz xproto.Atom
	atomNETWMStateMaximizedVert xproto.Atom
	atomNETWMStateAbove         xproto.Atom
	atomWMNormalHints           xproto.Atom
	atomWMSizeHints             xproto.Atom

	// Clipboard atoms
	atomClipboard xproto.Atom
	atomPrimary   xproto.Atom
	atomTargets   xproto.Atom
	atomMultiple  xproto.Atom
	atomText      xproto.Atom
	atomTextPlain xproto.Atom
	atomIncr      xproto.Atom

	// XDND atoms for drag-and-drop
	atomXdndEnter      xproto.Atom
	atomXdndPosition   xproto.Atom
	atomXdndStatus     xproto.Atom
	atomXdndLeave      xproto.Atom
	atomXdndDrop       xproto.Atom
	atomXdndFinished   xproto.Atom
	atomXdndSelection  xproto.Atom
	atomXdndTypeList   xproto.Atom
	atomXdndActionCopy xproto.Atom
	atomXdndActionMove xproto.Atom
	atomXdndActionLink xproto.Atom
	atomTextUriList    xproto.Atom
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

	s := &x11Driver{
		xc:  xc,
		xsi: xproto.Setup(xc).DefaultScreen(xc),
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("x11driver: init failed: %v", err)
	}

	return s, nil
}

// CreateWindow implements gio.Driver.
func (xd *x11Driver) CreateWindow(o gio.NewWindowOptions) (gio.Window, error) {
	bw := gio.NewBaseWindow(o)
	return newX11Window(xd, bw)
}

// Type implements gio.Driver.
func (xd *x11Driver) Type() gio.DriverType {
	return gio.DriverTypeX11
}

func (xd *x11Driver) init() error {
	if err := xd.initAtoms(); err != nil {
		return err
	}

	if err := xd.initKeyboardMapping(); err != nil {
		return err
	}

	return nil
}

func (xd *x11Driver) initAtoms() (err error) {
	atoms := map[string]*xproto.Atom{
		"_NET_WM_NAME":                 &xd.atomNETWMName,
		"UTF8_STRING":                  &xd.atomUTF8String,
		"WM_DELETE_WINDOW":             &xd.atomWMDeleteWindow,
		"WM_PROTOCOLS":                 &xd.atomWMProtocols,
		"WM_TAKE_FOCUS":                &xd.atomWMTakeFocus,
		"WM_CHANGE_STATE":              &xd.atomWMChangeState,
		"_NET_WM_STATE":                &xd.atomNETWMState,
		"_NET_WM_STATE_MAXIMIZED_HORZ": &xd.atomNETWMStateMaximizedHorz,
		"_NET_WM_STATE_MAXIMIZED_VERT": &xd.atomNETWMStateMaximizedVert,
		"_NET_WM_STATE_ABOVE":          &xd.atomNETWMStateAbove,
		"WM_NORMAL_HINTS":              &xd.atomWMNormalHints,
		"WM_SIZE_HINTS":                &xd.atomWMSizeHints,

		// Clipboard atoms
		"CLIPBOARD":  &xd.atomClipboard,
		"PRIMARY":    &xd.atomPrimary,
		"TARGETS":    &xd.atomTargets,
		"MULTIPLE":   &xd.atomMultiple,
		"TEXT":       &xd.atomText,
		"text/plain": &xd.atomTextPlain,
		"INCR":       &xd.atomIncr,

		// XDND atoms
		"XdndEnter":      &xd.atomXdndEnter,
		"XdndPosition":   &xd.atomXdndPosition,
		"XdndStatus":     &xd.atomXdndStatus,
		"XdndLeave":      &xd.atomXdndLeave,
		"XdndDrop":       &xd.atomXdndDrop,
		"XdndFinished":   &xd.atomXdndFinished,
		"XdndSelection":  &xd.atomXdndSelection,
		"XdndTypeList":   &xd.atomXdndTypeList,
		"XdndActionCopy": &xd.atomXdndActionCopy,
		"XdndActionMove": &xd.atomXdndActionMove,
		"XdndActionLink": &xd.atomXdndActionLink,
		"text/uri-list":  &xd.atomTextUriList,
	}

	for name, atom := range atoms {
		*atom, err = xd.internAtom(name)
		if err != nil {
			return fmt.Errorf("failed to intern atom %s: %v", name, err)
		}
	}
	return nil
}

func (xd *x11Driver) internAtom(name string) (xproto.Atom, error) {
	r, err := xproto.InternAtom(xd.xc, false, uint16(len(name)), name).Reply()
	if err != nil {
		return 0, fmt.Errorf("x11driver: xproto.InternAtom failed: %v", err)
	}
	if r == nil {
		return 0, fmt.Errorf("x11driver: xproto.InternAtom failed")
	}
	return r.Atom, nil
}

func (xd *x11Driver) initKeyboardMapping() error {
	const (
		keyLo            = 8
		keyHi            = 255
		maxSymsPerKey    = 6
		xkNumLock        = 0xff7f
		xkModeSwitch     = 0xff7e
		xkISOLevel3Shift = 0xfe03
	)
	// Get keyboard mapping.
	km, err := xproto.GetKeyboardMapping(xd.xc, keyLo, keyHi-keyLo+1).Reply()
	if err != nil {
		return err
	}
	n := int(km.KeysymsPerKeycode)
	if n < 2 {
		return fmt.Errorf("x11driver: too few keysyms per keycode: %d", n)
	}
	// Fill keysyms table.
	for i := keyLo; i <= keyHi; i++ {
		base := (i - keyLo) * n
		for j := 0; j < maxSymsPerKey; j++ {
			if j < n {
				xd.keysyms.Table[i][j] = uint32(km.Keysyms[base+j])
			} else {
				xd.keysyms.Table[i][j] = 0
			}
		}
	}

	// Get modifier mapping.
	mm, err := xproto.GetModifierMapping(xd.xc).Reply()
	if err != nil {
		return err
	}
	xd.keysyms.NumLockMod = 0
	xd.keysyms.ModeSwitchMod = 0
	xd.keysyms.ISOLevel3ShiftMod = 0

	keycodesPerMod := int(mm.KeycodesPerModifier)
	for modifier := 0; modifier < 8; modifier++ {
		for i := 0; i < keycodesPerMod; i++ {
			keycode := mm.Keycodes[modifier*keycodesPerMod+i]
			ks := xd.keysyms.Table[keycode][0]
			switch ks {
			case xkNumLock:
				xd.keysyms.NumLockMod = 1 << uint(modifier)
			case xkModeSwitch:
				xd.keysyms.ModeSwitchMod = 1 << uint(modifier)
			case xkISOLevel3Shift:
				xd.keysyms.ISOLevel3ShiftMod = 1 << uint(modifier)
			}
		}
	}
	return nil
}

func (xd *x11Driver) setProperty(xw xproto.Window, prop xproto.Atom, values ...xproto.Atom) {
	b := make([]byte, len(values)*4)
	for i, v := range values {
		b[4*i+0] = uint8(v >> 0)
		b[4*i+1] = uint8(v >> 8)
		b[4*i+2] = uint8(v >> 16)
		b[4*i+3] = uint8(v >> 24)
	}
	xproto.ChangeProperty(xd.xc, xproto.PropModeReplace, xw, prop, xproto.AtomAtom, 32, uint32(len(values)), b)
}

// translateKeyCode converts X11 keycode to gio KeyCode using the keysym table
func (xd *x11Driver) translateKeyCode(keycode xproto.Keycode, state uint16) (rune, gio.KeyCode) {
	return xd.keysyms.Lookup(keycode, state)
}

// translateModifiers converts X11 modifier state to gio ModifierKey
func (xd *x11Driver) translateModifiers(state uint16) gio.ModifierKey {
	return x11key.KeyModifiers(state)
}

// translateMouseButton converts X11 button to gio MouseButton
func (xd *x11Driver) translateMouseButton(button xproto.Button) gio.MouseButton {
	switch button {
	case 1:
		return gio.MouseButtonLeft
	case 2:
		return gio.MouseButtonMiddle
	case 3:
		return gio.MouseButtonRight
	case 8:
		return gio.MouseButtonBack
	case 9:
		return gio.MouseButtonForward
	default:
		return gio.MouseButtonUnknown
	}
}

// translateWheelDelta converts X11 wheel button to scroll delta
func (xd *x11Driver) translateWheelDelta(button xproto.Button) (float64, float64) {
	switch button {
	case 4: // Scroll up
		return 0, -1
	case 5: // Scroll down
		return 0, 1
	case 6: // Scroll left
		return -1, 0
	case 7: // Scroll right
		return 1, 0
	default:
		return 0, 0
	}
}
