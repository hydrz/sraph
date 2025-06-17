package x11

import (
	"fmt"
	"image"
	"log"
	"strings"
	"sync"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"github.com/opensraph/sraph/gio"
)

var _ gio.Window = (*x11Window)(nil)

// x11Window represents an X11 window implementation
type x11Window struct {
	gio.BaseWindow

	xDriver  *x11Driver
	xWindow  xproto.Window
	XContext xproto.Gcontext
	xEvent   chan xgb.Event

	// Clipboard support
	clipboardOwner    bool
	clipboardData     *gio.DataTransfer
	selectionProperty xproto.Atom

	// XDND drag-and-drop support
	xdndVersion      uint32
	xdndSourceWindow xproto.Window
	xdndDragContext  *gio.DataTransfer

	// Touch and pointer tracking (simplified without XI2)
	activeTouches  map[int]*gio.Touch // Maps touch ID to Touch
	lastPointerPos image.Point        // Last pointer position for motion detection

	mu sync.RWMutex
}

func newX11Window(driver *x11Driver, bw gio.BaseWindow) (*x11Window, error) {
	xw := &x11Window{
		xDriver:           driver,
		BaseWindow:        bw,
		clipboardData:     gio.NewDataTransfer(),
		selectionProperty: 0,
		xdndVersion:       5,
		activeTouches:     make(map[int]*gio.Touch),
		lastPointerPos:    image.Point{},
	}

	if err := xw.init(); err != nil {
		return nil, fmt.Errorf("x11driver: x11Window.init failed: %v", err)
	}

	// Start event handling
	go xw.loop()

	return xw, nil
}

func (xw *x11Window) init() error {
	conn := xw.xDriver.xc

	xWindow, err := xproto.NewWindowId(conn)
	if err != nil {
		return fmt.Errorf("x11driver: xproto.NewWindowId failed: %v", err)
	}

	xg, err := xproto.NewGcontextId(conn)
	if err != nil {
		return fmt.Errorf("x11driver: xproto.NewGcontextId failed: %v", err)
	}

	xw.mu.Lock()
	xw.xWindow = xWindow
	xw.XContext = xg
	xw.xEvent = make(chan xgb.Event)
	defer xw.mu.Unlock()

	width, height := 1024, 768
	opts := xw.Attr()
	if opts.Width > 0 {
		width = opts.Width
	}
	if opts.Height > 0 {
		height = opts.Height
	}

	xproto.CreateWindow(conn, xw.xDriver.xsi.RootDepth, xWindow, xw.xDriver.xsi.Root,
		0, 0, uint16(width), uint16(height), 0,
		xproto.WindowClassInputOutput, xw.xDriver.xsi.RootVisual,
		xproto.CwEventMask,
		[]uint32{0 |
			xproto.EventMaskKeyPress |
			xproto.EventMaskKeyRelease |
			xproto.EventMaskButtonPress |
			xproto.EventMaskButtonRelease |
			xproto.EventMaskPointerMotion |
			xproto.EventMaskExposure |
			xproto.EventMaskStructureNotify |
			xproto.EventMaskFocusChange,
		},
	)
	xw.xDriver.setProperty(xWindow, xw.xDriver.atomWMProtocols, xw.xDriver.atomWMDeleteWindow, xw.xDriver.atomWMTakeFocus)

	title := []byte(opts.Title)
	xproto.ChangeProperty(xw.xDriver.xc, xproto.PropModeReplace, xWindow, xw.xDriver.atomNETWMName, xw.xDriver.atomUTF8String, 8, uint32(len(title)), title)

	xproto.CreateGC(xw.xDriver.xc, xg, xproto.Drawable(xWindow), 0, nil)
	xproto.MapWindow(xw.xDriver.xc, xWindow)

	// Create a property atom for clipboard operations
	selectionProperty, err := xw.xDriver.internAtom("SRAPH_SELECTION")
	if err != nil {
		return fmt.Errorf("x11driver: failed to create selection property atom: %v", err)
	}
	xw.selectionProperty = selectionProperty

	// Enable XDND for drag-and-drop
	xdndAware := []uint32{xw.xdndVersion}
	xdndAwareBytes := make([]byte, 4)
	xdndAwareBytes[0] = uint8(xdndAware[0] >> 0)
	xdndAwareBytes[1] = uint8(xdndAware[0] >> 8)
	xdndAwareBytes[2] = uint8(xdndAware[0] >> 16)
	xdndAwareBytes[3] = uint8(xdndAware[0] >> 24)

	xdndAwareAtom, err := xw.xDriver.internAtom("XdndAware")
	if err != nil {
		return fmt.Errorf("x11driver: failed to intern XdndAware atom: %v", err)
	}

	xproto.ChangeProperty(xw.xDriver.xc, xproto.PropModeReplace, xWindow,
		xdndAwareAtom, xproto.AtomAtom, 32, 1, xdndAwareBytes)

	// Note: XI2 setup removed as it's not available
	log.Printf("x11driver: using basic X11 events for touch/pointer simulation")

	return nil
}

// WindowID returns the X11 window ID.
func (xw *x11Window) WindowID() gio.WindowID {
	xw.mu.RLock()
	defer xw.mu.RUnlock()
	return gio.WindowID(xw.xWindow)
}

func (xw *x11Window) SetAttr(attr gio.WindowAttr) {
	xw.mu.Lock()
	defer xw.mu.Unlock()

	currentAttr := xw.Attr()
	if currentAttr.Equal(attr) {
		return
	}

	// Validate attributes first
	if err := xw.validateAttributes(attr); err != nil {
		return // Error already logged
	}

	// Handle window closure first - early return
	if attr.State.Contains(gio.WindowStateClosed) {
		xw.handleWindowClose()
		xw.BaseWindow.SetAttr(attr)
		return
	}

	// Update basic properties
	xw.updateTitle(attr, currentAttr)
	xw.updateGeometry(attr, currentAttr)

	// Handle state changes
	xw.updateWindowStates(attr, currentAttr)

	// Handle resizability
	xw.updateResizability(attr, currentAttr)

	// Update base attributes last
	xw.BaseWindow.SetAttr(attr)
}

// validateAttributes validates window attribute values
func (xw *x11Window) validateAttributes(attr gio.WindowAttr) error {
	const maxDimension = 32767 // X11 coordinate limit

	if attr.Width < 0 || attr.Height < 0 {
		return fmt.Errorf("x11driver: negative dimensions: width=%d, height=%d", attr.Width, attr.Height)
	}

	if attr.Width > maxDimension || attr.Height > maxDimension {
		return fmt.Errorf("x11driver: dimensions too large: width=%d, height=%d (limit=%d)",
			attr.Width, attr.Height, maxDimension)
	}

	if attr.Position.X < -maxDimension || attr.Position.X > maxDimension ||
		attr.Position.Y < -maxDimension || attr.Position.Y > maxDimension {
		return fmt.Errorf("x11driver: position out of range: x=%d, y=%d (range=±%d)",
			attr.Position.X, attr.Position.Y, maxDimension)
	}

	return nil
}

// handleWindowClose destroys the window
func (xw *x11Window) handleWindowClose() {
	if err := xproto.DestroyWindowChecked(xw.xDriver.xc, xw.xWindow).Check(); err != nil {
		// Window might already be destroyed, log but don't fail
		fmt.Printf("x11driver: failed to destroy window: %v\n", err)
	}
}

// updateTitle updates the window title if changed
func (xw *x11Window) updateTitle(attr, currentAttr gio.WindowAttr) {
	if attr.Title == "" || attr.Title == currentAttr.Title {
		return
	}

	titleBytes := []byte(attr.Title)
	xproto.ChangeProperty(xw.xDriver.xc, xproto.PropModeReplace, xw.xWindow,
		xw.xDriver.atomNETWMName, xw.xDriver.atomUTF8String, 8,
		uint32(len(titleBytes)), titleBytes)
}

// updateGeometry updates window size and position
func (xw *x11Window) updateGeometry(attr, currentAttr gio.WindowAttr) {
	sizeChanged := attr.Width != currentAttr.Width || attr.Height != currentAttr.Height
	posChanged := attr.Position.X != currentAttr.Position.X || attr.Position.Y != currentAttr.Position.Y

	if sizeChanged && attr.Width > 0 && attr.Height > 0 {
		err := xproto.ConfigureWindowChecked(xw.xDriver.xc, xw.xWindow,
			xproto.ConfigWindowWidth|xproto.ConfigWindowHeight,
			[]uint32{uint32(attr.Width), uint32(attr.Height)}).Check()
		if err != nil {
			fmt.Printf("x11driver: failed to resize window: %v\n", err)
		}
	}

	if posChanged {
		err := xproto.ConfigureWindowChecked(xw.xDriver.xc, xw.xWindow,
			xproto.ConfigWindowX|xproto.ConfigWindowY,
			[]uint32{uint32(attr.Position.X), uint32(attr.Position.Y)}).Check()
		if err != nil {
			fmt.Printf("x11driver: failed to move window: %v\n", err)
		}
	}
}

// updateWindowStates handles window state changes
func (xw *x11Window) updateWindowStates(attr, currentAttr gio.WindowAttr) {
	conn := xw.xDriver.xc

	// Handle maximized state
	maximizedChanged := attr.State.Contains(gio.WindowStateMaximized) != currentAttr.State.Contains(gio.WindowStateMaximized)
	if maximizedChanged {
		action := uint32(0) // _NET_WM_STATE_REMOVE
		if attr.State.Contains(gio.WindowStateMaximized) {
			action = 1 // _NET_WM_STATE_ADD
		}

		xw.sendNetWMStateMessage(action, xw.xDriver.atomNETWMStateMaximizedHorz, xw.xDriver.atomNETWMStateMaximizedVert)
	}

	// Handle floating (always on top) state
	floatingChanged := attr.State.Contains(gio.WindowStateFloating) != currentAttr.State.Contains(gio.WindowStateFloating)
	if floatingChanged {
		action := uint32(0) // _NET_WM_STATE_REMOVE
		if attr.State.Contains(gio.WindowStateFloating) {
			action = 1 // _NET_WM_STATE_ADD
		}

		xw.sendNetWMStateMessage(action, xw.xDriver.atomNETWMStateAbove, 0)
	}

	// Handle focus state
	focusChanged := attr.State.Contains(gio.WindowStateFocused) != currentAttr.State.Contains(gio.WindowStateFocused)
	if focusChanged && attr.State.Contains(gio.WindowStateFocused) {
		xproto.SetInputFocus(conn, xproto.InputFocusPointerRoot, xw.xWindow, xproto.TimeCurrentTime)
	}

	// Handle visible state
	visibleChanged := attr.State.Contains(gio.WindowStateVisible) != currentAttr.State.Contains(gio.WindowStateVisible)
	if visibleChanged {
		if attr.State.Contains(gio.WindowStateVisible) {
			xproto.MapWindow(conn, xw.xWindow)
		} else {
			xproto.UnmapWindow(conn, xw.xWindow)
		}
	}
}

// updateResizability handles window resizability changes
func (xw *x11Window) updateResizability(attr, currentAttr gio.WindowAttr) {
	resizableChanged := attr.State.Contains(gio.WindowStateResizable) != currentAttr.State.Contains(gio.WindowStateResizable)
	if !resizableChanged && attr.Width == currentAttr.Width && attr.Height == currentAttr.Height {
		return
	}

	width := attr.Width
	height := attr.Height
	if width <= 0 {
		width = 1024
	}
	if height <= 0 {
		height = 768
	}

	var minWidth, minHeight, maxWidth, maxHeight uint32
	if attr.State.Contains(gio.WindowStateResizable) {
		minWidth, minHeight = 100, 100     // Reasonable minimum
		maxWidth, maxHeight = 65535, 65535 // Maximum possible	} else {
		// Fixed size window
		minWidth, minHeight = uint32(width), uint32(height)
		maxWidth, maxHeight = uint32(width), uint32(height)
	}

	// Set WM_NORMAL_HINTS to control resizability
	hints := make([]uint32, 18)
	hints[0] = 1<<4 | 1<<5 | 1<<6 // PMinSize | PMaxSize | PResizeInc
	hints[5] = minWidth
	hints[6] = minHeight
	hints[7] = maxWidth
	hints[8] = maxHeight
	hints[9] = 1  // width_inc
	hints[10] = 1 // height_inc

	hintsBytes := make([]byte, len(hints)*4)
	for i, h := range hints {
		hintsBytes[4*i+0] = uint8(h >> 0)
		hintsBytes[4*i+1] = uint8(h >> 8)
		hintsBytes[4*i+2] = uint8(h >> 16)
		hintsBytes[4*i+3] = uint8(h >> 24)
	}

	xproto.ChangeProperty(xw.xDriver.xc, xproto.PropModeReplace, xw.xWindow,
		xw.xDriver.atomWMNormalHints, xw.xDriver.atomWMSizeHints, 32,
		uint32(len(hints)), hintsBytes)
}

// sendNetWMStateMessage sends a _NET_WM_STATE client message
func (xw *x11Window) sendNetWMStateMessage(action uint32, prop1, prop2 xproto.Atom) {
	data := []uint32{action, uint32(prop1), uint32(prop2), 0, 0}
	ev := xproto.ClientMessageEvent{
		Format: 32,
		Window: xw.xWindow,
		Type:   xw.xDriver.atomNETWMState,
		Data:   xproto.ClientMessageDataUnionData32New(data),
	}

	xproto.SendEvent(xw.xDriver.xc, false, xw.xDriver.xsi.Root,
		xproto.EventMaskSubstructureNotify|xproto.EventMaskSubstructureRedirect,
		string(ev.Bytes()))
}

// SetClipboardData sets data to the system clipboard
func (xw *x11Window) SetClipboardData(data *gio.DataTransfer) error {
	xw.mu.Lock()
	defer xw.mu.Unlock()

	// Store clipboard data
	xw.clipboardData = data

	// Take ownership of the clipboard selection
	xproto.SetSelectionOwner(xw.xDriver.xc, xw.xWindow, xw.xDriver.atomClipboard, xproto.TimeCurrentTime)

	// Verify we got ownership
	reply, err := xproto.GetSelectionOwner(xw.xDriver.xc, xw.xDriver.atomClipboard).Reply()
	if err != nil {
		return fmt.Errorf("x11driver: failed to verify clipboard ownership: %v", err)
	}

	xw.clipboardOwner = (reply.Owner == xw.xWindow)
	return nil
}

// GetClipboardData requests data from the system clipboard
func (xw *x11Window) GetClipboardData() error {
	// Request clipboard content by converting the selection
	xproto.ConvertSelection(xw.xDriver.xc, xw.xWindow, xw.xDriver.atomClipboard,
		xw.xDriver.atomUTF8String, xw.selectionProperty, xproto.TimeCurrentTime)
	return nil
}

// sendXdndStatus sends XDND status message
func (xw *x11Window) sendXdndStatus(sourceWindow xproto.Window, accept bool, action xproto.Atom) {
	flags := uint32(0)
	if accept {
		flags |= 1 // Accept drop
	}

	data := []uint32{uint32(xw.xWindow), flags, 0, 0, uint32(action)}
	ev := xproto.ClientMessageEvent{
		Format: 32,
		Window: sourceWindow,
		Type:   xw.xDriver.atomXdndStatus,
		Data:   xproto.ClientMessageDataUnionData32New(data),
	}

	xproto.SendEvent(xw.xDriver.xc, false, sourceWindow, 0, string(ev.Bytes()))
}

// sendXdndFinished sends XDND finished message
func (xw *x11Window) sendXdndFinished(sourceWindow xproto.Window, action xproto.Atom) {
	data := []uint32{uint32(xw.xWindow), 1, uint32(action), 0, 0} // 1 = accept
	ev := xproto.ClientMessageEvent{
		Format: 32,
		Window: sourceWindow,
		Type:   xw.xDriver.atomXdndFinished,
		Data:   xproto.ClientMessageDataUnionData32New(data),
	}

	xproto.SendEvent(xw.xDriver.xc, false, sourceWindow, 0, string(ev.Bytes()))
}

// parseXdndTypeList parses XDND type list from a property
func (xw *x11Window) parseXdndTypeList(property xproto.Atom) []string {
	reply, err := xproto.GetProperty(xw.xDriver.xc, false, xw.xdndSourceWindow,
		property, xproto.AtomAtom, 0, 1024).Reply()
	if err != nil || reply == nil {
		return nil
	}

	var types []string
	for i := 0; i < len(reply.Value); i += 4 {
		if i+3 < len(reply.Value) {
			atom := xproto.Atom(uint32(reply.Value[i]) |
				uint32(reply.Value[i+1])<<8 |
				uint32(reply.Value[i+2])<<16 |
				uint32(reply.Value[i+3])<<24)

			nameReply, err := xproto.GetAtomName(xw.xDriver.xc, atom).Reply()
			if err == nil && nameReply != nil {
				types = append(types, nameReply.Name)
			}
		}
	}

	return types
}

// convertXdndUriListToFiles converts text/uri-list to file paths
func (xw *x11Window) convertXdndUriListToFiles(uriList string) []string {
	var files []string
	lines := strings.Split(uriList, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "file://") {
			// Remove file:// prefix and decode URI
			filePath := strings.TrimPrefix(line, "file://")
			files = append(files, filePath)
		}
	}

	return files
}
