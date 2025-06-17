package x11

import (
	"image"
	"log"
	"time"

	"github.com/jezek/xgb/xproto"
	"github.com/opensraph/sraph/gio"
)

func (xw *x11Window) loop() {
	conn := xw.xDriver.xc

	for {
		ev, err := conn.WaitForEvent()
		if err != nil {
			log.Printf("x11driver: xproto.WaitForEvent: %v", err)
			continue
		}

		xw.handleEvent(ev)
	}
}

func (xw *x11Window) handleEvent(ev interface{}) {
	switch ev := ev.(type) {
	case xproto.DestroyNotifyEvent:
		if ev.Window == xw.xWindow {
			xw.handleDestroyNotify(ev)
		}

	case xproto.ClientMessageEvent:
		// Handle XDND drag-and-drop events
		if xw.isXDNDClientMessage(ev) {
			xw.handleXDNDClientMessage(ev)
			return
		}
		if ev.Window == xw.xWindow {
			xw.handleClientMessage(ev)
		}

	case xproto.ConfigureNotifyEvent:
		if ev.Window == xw.xWindow {
			xw.handleConfigureNotify(ev)
		}

	case xproto.ExposeEvent:
		if ev.Window == xw.xWindow {
			xw.handleExpose(ev)
		}

	case xproto.FocusInEvent:
		if ev.Event == xw.xWindow {
			xw.handleFocusIn(ev)
		}

	case xproto.FocusOutEvent:
		if ev.Event == xw.xWindow {
			xw.handleFocusOut(ev)
		}

	case xproto.KeyPressEvent:
		if ev.Event == xw.xWindow {
			xw.handleKeyPress(ev)
		}

	case xproto.KeyReleaseEvent:
		if ev.Event == xw.xWindow {
			xw.handleKeyRelease(ev)
		}

	case xproto.ButtonPressEvent:
		if ev.Event == xw.xWindow {
			xw.handleButtonPress(ev)
		}

	case xproto.ButtonReleaseEvent:
		if ev.Event == xw.xWindow {
			xw.handleButtonRelease(ev)
		}

	case xproto.MotionNotifyEvent:
		if ev.Event == xw.xWindow {
			xw.handleMotionNotify(ev)
		}

	case xproto.MapNotifyEvent:
		if ev.Window == xw.xWindow {
			xw.handleMapNotify(ev)
		}

	case xproto.UnmapNotifyEvent:
		if ev.Window == xw.xWindow {
			xw.handleUnmapNotify(ev)
		}

	case xproto.SelectionNotifyEvent:
		if ev.Requestor == xw.xWindow {
			xw.handleSelectionNotify(ev)
		}
	case xproto.SelectionRequestEvent:
		if ev.Owner == xw.xWindow {
			xw.handleSelectionRequest(ev)
		}

	case xproto.SelectionClearEvent:
		if ev.Owner == xw.xWindow {
			xw.handleSelectionClear(ev)
		}

	case xproto.PropertyNotifyEvent:
		if ev.Window == xw.xWindow {
			xw.handlePropertyNotify(ev)
		}

	case xproto.GeGenericEvent:
		// Handle XI2 generic events
		if ev.SequenceId() == uint16(xw.xDriver.xiOpcode) {
			xw.handleXI2GenericEvent(ev)
		}

	default:
		log.Printf("x11driver: unhandled event type %T", ev)
	}
}

func (xw *x11Window) handleDestroyNotify(ev xproto.DestroyNotifyEvent) {
	attr := xw.Attr()
	attr.State |= gio.WindowStateClosed
	xw.BaseWindow.SetAttr(attr)
}

func (xw *x11Window) handleClientMessage(ev xproto.ClientMessageEvent) {
	// Handle WM_DELETE_WINDOW protocol
	if ev.Type == xw.xDriver.atomWMProtocols &&
		len(ev.Data.Data32) > 0 &&
		xproto.Atom(ev.Data.Data32[0]) == xw.xDriver.atomWMDeleteWindow {
		attr := xw.Attr()
		attr.State |= gio.WindowStateClosed
		xw.BaseWindow.SetAttr(attr)
	}
}

func (xw *x11Window) handleConfigureNotify(ev xproto.ConfigureNotifyEvent) {
	attr := xw.Attr()
	changed := false

	// Update size if changed
	if attr.Width != int(ev.Width) || attr.Height != int(ev.Height) {
		attr.Width = int(ev.Width)
		attr.Height = int(ev.Height)
		changed = true
	}

	// Update position if changed
	if attr.Position.X != int(ev.X) || attr.Position.Y != int(ev.Y) {
		attr.Position.X = int(ev.X)
		attr.Position.Y = int(ev.Y)
		changed = true
	}

	if changed {
		xw.BaseWindow.SetAttr(attr)
	}
}

func (xw *x11Window) handleExpose(ev xproto.ExposeEvent) {
	// Publish paint/redraw event

}

func (xw *x11Window) handleFocusIn(ev xproto.FocusInEvent) {
	attr := xw.Attr()
	if !attr.State.Contains(gio.WindowStateFocused) {
		attr.State |= gio.WindowStateFocused
		xw.BaseWindow.SetAttr(attr)
	}
}

func (xw *x11Window) handleFocusOut(ev xproto.FocusOutEvent) {
	attr := xw.Attr()
	if attr.State.Contains(gio.WindowStateFocused) {
		attr.State &^= gio.WindowStateFocused
		xw.BaseWindow.SetAttr(attr)
	}
}

func (xw *x11Window) handleKeyPress(ev xproto.KeyPressEvent) {
	_, keyCode := xw.xDriver.translateKeyCode(ev.Detail, ev.State)
	modifierKey := xw.xDriver.translateModifiers(ev.State)

	// Create keyboard event - baseWindow handles all repeat detection internally
	keyEvent := gio.NewKeyboardEvent()
	keyEvent.Code = keyCode
	keyEvent.ModifierKey = modifierKey
	keyEvent.State = gio.KeyStatePressed

	xw.Publish(keyEvent)
}

func (xw *x11Window) handleKeyRelease(ev xproto.KeyReleaseEvent) {
	_, keyCode := xw.xDriver.translateKeyCode(ev.Detail, ev.State)
	modifierKey := xw.xDriver.translateModifiers(ev.State)

	// Create keyboard event - baseWindow handles all repeat detection internally
	keyEvent := gio.NewKeyboardEvent()
	keyEvent.Code = keyCode
	keyEvent.ModifierKey = modifierKey
	keyEvent.State = gio.KeyStateReleased

	xw.Publish(keyEvent)
}

func (xw *x11Window) handleButtonPress(ev xproto.ButtonPressEvent) {
	button := xw.xDriver.translateMouseButton(ev.Detail)
	modifierKey := xw.xDriver.translateModifiers(ev.State)
	position := image.Point{X: int(ev.EventX), Y: int(ev.EventY)}

	// Handle mouse wheel events (buttons 4-7)
	if ev.Detail >= 4 && ev.Detail <= 7 {
		xw.handleWheelEvent(ev)
		return
	}

	// Create standard mouse event
	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = button
	mouseEvent.ModifierKey = modifierKey
	mouseEvent.Position = position
	xw.Publish(mouseEvent)
}

// handleWheelEvent handles mouse wheel events (buttons 4-7)
func (xw *x11Window) handleWheelEvent(ev xproto.ButtonPressEvent) {
	deltaX, deltaY := xw.xDriver.translateWheelDelta(ev.Detail)
	wheelEvent := gio.NewWheelEvent()
	wheelEvent.DeltaX = deltaX
	wheelEvent.DeltaY = deltaY
	wheelEvent.Position = image.Point{X: int(ev.EventX), Y: int(ev.EventY)}
	wheelEvent.ModifierKey = xw.xDriver.translateModifiers(ev.State)
	xw.Publish(wheelEvent)
}

func (xw *x11Window) handleButtonRelease(ev xproto.ButtonReleaseEvent) {
	// Skip wheel events for release
	if ev.Detail >= 4 && ev.Detail <= 7 {
		return
	}

	button := xw.xDriver.translateMouseButton(ev.Detail)
	modifierKey := xw.xDriver.translateModifiers(ev.State)
	position := image.Point{X: int(ev.EventX), Y: int(ev.EventY)}

	// Create standard mouse event
	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = button
	mouseEvent.ModifierKey = modifierKey
	mouseEvent.Position = position
	xw.Publish(mouseEvent)
}

func (xw *x11Window) handleMotionNotify(ev xproto.MotionNotifyEvent) {
	modifierKey := xw.xDriver.translateModifiers(ev.State)
	position := image.Point{X: int(ev.EventX), Y: int(ev.EventY)}

	// Create standard mouse move event
	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = gio.MouseButtonUnknown
	mouseEvent.ModifierKey = modifierKey
	mouseEvent.Position = position
	xw.Publish(mouseEvent)

	xw.lastPointerPos = position
}

func (xw *x11Window) handleMapNotify(ev xproto.MapNotifyEvent) {
	attr := xw.Attr()
	if !attr.State.Contains(gio.WindowStateVisible) {
		attr.State |= gio.WindowStateVisible
		xw.BaseWindow.SetAttr(attr)
	}
}

func (xw *x11Window) handleUnmapNotify(ev xproto.UnmapNotifyEvent) {
	attr := xw.Attr()
	if attr.State.Contains(gio.WindowStateVisible) {
		attr.State &^= gio.WindowStateVisible
		xw.BaseWindow.SetAttr(attr)
	}
}

// handleSelectionNotify handles clipboard paste (SelectionNotify)
func (xw *x11Window) handleSelectionNotify(ev xproto.SelectionNotifyEvent) {
	clipboardEvent := gio.NewClipboardEvent()

	if ev.Property == 0 {
		log.Printf("x11driver: SelectionNotify received but property is 0 (no clipboard data)")
		xw.Publish(clipboardEvent)
		return
	}

	reply, err := xproto.GetProperty(xw.xDriver.xc, false, xw.xWindow, ev.Property,
		xproto.GetPropertyTypeAny, 0, (1<<24)-1).Reply()
	if err != nil {
		log.Printf("x11driver: failed to get clipboard property: %v", err)
		xw.Publish(clipboardEvent)
		return
	}

	if reply == nil || len(reply.Value) == 0 {
		log.Printf("x11driver: clipboard property is empty")
		xw.Publish(clipboardEvent)
		return
	}

	var text string
	var mimeType string

	switch reply.Type {
	case xw.xDriver.atomUTF8String:
		text = string(reply.Value)
		mimeType = "text/plain"
	case xproto.AtomString:
		text = string(reply.Value)
		mimeType = "text/plain"
	case xw.xDriver.atomTextUriList:
		text = string(reply.Value)
		mimeType = "text/uri-list"
		// Convert URI list to files
		files := xw.convertXdndUriListToFiles(text)
		for _, filePath := range files {
			if err := clipboardEvent.Data.AddFile(filePath); err != nil {
				log.Printf("x11driver: failed to add file to clipboard: %v", err)
			}
		}
	default:
		text = string(reply.Value)
		mimeType = "application/octet-stream"
	}

	if err := clipboardEvent.Data.SetData(mimeType, text); err != nil {
		log.Printf("x11driver: failed to set clipboard data: %v", err)
	}

	// Clean up the property
	xproto.DeleteProperty(xw.xDriver.xc, xw.xWindow, ev.Property)

	xw.Publish(clipboardEvent)
}

// handleSelectionRequest handles clipboard copy (SelectionRequest)
func (xw *x11Window) handleSelectionRequest(ev xproto.SelectionRequestEvent) {
	xw.mu.RLock()
	hasData := xw.clipboardOwner && xw.clipboardData != nil
	xw.mu.RUnlock()

	if !hasData {
		// Send failure response
		xw.sendSelectionNotify(ev, 0)
		return
	}

	// Handle TARGETS request
	if ev.Target == xw.xDriver.atomTargets {
		xw.handleTargetsRequest(ev)
		return
	}

	// Handle text requests
	var data string
	var found bool

	switch ev.Target {
	case xw.xDriver.atomUTF8String, xw.xDriver.atomTextPlain, xproto.AtomString:
		data, found = xw.clipboardData.GetText()
		if !found {
			if dataInterface, exists := xw.clipboardData.GetData("text/plain"); exists {
				if str, ok := dataInterface.(string); ok {
					data = str
					found = true
				}
			}
		}
	}

	if found {
		// Set the property with the data
		xproto.ChangeProperty(xw.xDriver.xc, xproto.PropModeReplace, ev.Requestor,
			ev.Property, ev.Target, 8, uint32(len(data)), []byte(data))
		xw.sendSelectionNotify(ev, ev.Property)
	} else {
		xw.sendSelectionNotify(ev, 0)
	}
}

// handleSelectionClear handles losing clipboard ownership
func (xw *x11Window) handleSelectionClear(ev xproto.SelectionClearEvent) {
	xw.mu.Lock()
	xw.clipboardOwner = false
	xw.mu.Unlock()

	log.Printf("x11driver: lost clipboard ownership")
}

// handlePropertyNotify handles property change notifications
func (xw *x11Window) handlePropertyNotify(ev xproto.PropertyNotifyEvent) {
	// Handle INCR (incremental) clipboard transfers if needed
	log.Printf("x11driver: property notify: atom=%d, state=%d", ev.Atom, ev.State)
}

// handleTargetsRequest handles TARGETS selection request
func (xw *x11Window) handleTargetsRequest(ev xproto.SelectionRequestEvent) {
	targets := []xproto.Atom{
		xw.xDriver.atomTargets,
		xw.xDriver.atomUTF8String,
		xproto.AtomString,
		xw.xDriver.atomTextPlain,
	}

	targetBytes := make([]byte, len(targets)*4)
	for i, target := range targets {
		targetBytes[4*i+0] = uint8(target >> 0)
		targetBytes[4*i+1] = uint8(target >> 8)
		targetBytes[4*i+2] = uint8(target >> 16)
		targetBytes[4*i+3] = uint8(target >> 24)
	}

	xproto.ChangeProperty(xw.xDriver.xc, xproto.PropModeReplace, ev.Requestor,
		ev.Property, xproto.AtomAtom, 32, uint32(len(targets)), targetBytes)

	xw.sendSelectionNotify(ev, ev.Property)
}

// sendSelectionNotify sends SelectionNotify response
func (xw *x11Window) sendSelectionNotify(req xproto.SelectionRequestEvent, property xproto.Atom) {
	ev := xproto.SelectionNotifyEvent{
		Time:      req.Time,
		Requestor: req.Requestor,
		Selection: req.Selection,
		Target:    req.Target,
		Property:  property,
	}

	xproto.SendEvent(xw.xDriver.xc, false, req.Requestor, 0, string(ev.Bytes()))
}

// isXDNDClientMessage checks if a ClientMessage is an XDND drag-and-drop event
func (xw *x11Window) isXDNDClientMessage(ev xproto.ClientMessageEvent) bool {
	return ev.Type == xw.xDriver.atomXdndEnter ||
		ev.Type == xw.xDriver.atomXdndPosition ||
		ev.Type == xw.xDriver.atomXdndLeave ||
		ev.Type == xw.xDriver.atomXdndDrop
}

// handleXDNDClientMessage handles XDND drag-and-drop events
func (xw *x11Window) handleXDNDClientMessage(ev xproto.ClientMessageEvent) {
	switch ev.Type {
	case xw.xDriver.atomXdndEnter:
		xw.handleXdndEnter(ev)
	case xw.xDriver.atomXdndPosition:
		xw.handleXdndPosition(ev)
	case xw.xDriver.atomXdndLeave:
		xw.handleXdndLeave(ev)
	case xw.xDriver.atomXdndDrop:
		xw.handleXdndDrop(ev)
	}
}

// handleXdndEnter handles XDND enter event
func (xw *x11Window) handleXdndEnter(ev xproto.ClientMessageEvent) {
	data := ev.Data.Data32
	xw.xdndSourceWindow = xproto.Window(data[0])
	xw.xdndVersion = (data[1] >> 24) & 0xFF

	dragEvent := gio.NewDragEvent()
	xw.xdndDragContext = dragEvent.Data

	// Parse data types
	if (data[1] & 1) != 0 {
		// More than 3 types, read from XdndTypeList property
		types := xw.parseXdndTypeList(xw.xDriver.atomXdndTypeList)
		for _, typeName := range types {
			dragEvent.Data.SetData(typeName, "")
		}
	} else {
		// Up to 3 types in the message
		for i := 2; i <= 4; i++ {
			if data[i] != 0 {
				atomName, err := xproto.GetAtomName(xw.xDriver.xc, xproto.Atom(data[i])).Reply()
				if err == nil && atomName != nil {
					dragEvent.Data.SetData(atomName.Name, "")
				}
			}
		}
	}

	log.Printf("x11driver: XDND enter from window %d, version %d", xw.xdndSourceWindow, xw.xdndVersion)
	xw.Publish(dragEvent)
}

// handleXdndPosition handles XDND position event
func (xw *x11Window) handleXdndPosition(ev xproto.ClientMessageEvent) {
	data := ev.Data.Data32

	// Extract position (data[2] contains x and y coordinates)
	x := int16(data[2] >> 16)
	y := int16(data[2] & 0xFFFF)

	dragEvent := gio.NewDragEvent()
	if xw.xdndDragContext != nil {
		dragEvent.Data = xw.xdndDragContext
	}

	log.Printf("x11driver: XDND position at (%d, %d)", x, y)

	// Send status response (accept the drop)
	xw.sendXdndStatus(xw.xdndSourceWindow, true, xw.xDriver.atomXdndActionCopy)

	xw.Publish(dragEvent)
}

// handleXdndLeave handles XDND leave event
func (xw *x11Window) handleXdndLeave(ev xproto.ClientMessageEvent) {
	dragEvent := gio.NewDragEvent()
	if xw.xdndDragContext != nil {
		dragEvent.Data = xw.xdndDragContext
	}

	log.Printf("x11driver: XDND leave")

	// Clean up drag context
	xw.xdndDragContext = nil
	xw.xdndSourceWindow = 0

	xw.Publish(dragEvent)
}

// handleXdndDrop handles XDND drop event
func (xw *x11Window) handleXdndDrop(ev xproto.ClientMessageEvent) {
	dragEvent := gio.NewDragEvent()
	if xw.xdndDragContext != nil {
		dragEvent.Data = xw.xdndDragContext
	}

	// Request the actual data
	xproto.ConvertSelection(xw.xDriver.xc, xw.xWindow, xw.xDriver.atomXdndSelection,
		xw.xDriver.atomTextUriList, xw.selectionProperty, xproto.TimeCurrentTime)

	log.Printf("x11driver: XDND drop")

	// Send finished message
	xw.sendXdndFinished(xw.xdndSourceWindow, xw.xDriver.atomXdndActionCopy)

	// Wait a bit for the selection data, then publish the event
	go func() {
		time.Sleep(100 * time.Millisecond)
		xw.Publish(dragEvent)

		// Clean up
		xw.xdndDragContext = nil
		xw.xdndSourceWindow = 0
	}()
}
