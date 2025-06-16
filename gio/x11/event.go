package x11

import (
	"image"
	"log"

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

	case xproto.MappingNotifyEvent:
		xw.handleMappingNotify(ev)

	case xproto.MapNotifyEvent:
		if ev.Window == xw.xWindow {
			xw.handleMapNotify(ev)
		}

	case xproto.UnmapNotifyEvent:
		if ev.Window == xw.xWindow {
			xw.handleUnmapNotify(ev)
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

	keyEvent := gio.NewKeyboardEvent()
	keyEvent.Code = keyCode
	keyEvent.ModifierKey = modifierKey
	keyEvent.Repeat = false // TODO: detect key repeat

	xw.Publish(keyEvent)
}

func (xw *x11Window) handleKeyRelease(ev xproto.KeyReleaseEvent) {
	_, keyCode := xw.xDriver.translateKeyCode(ev.Detail, ev.State)
	modifierKey := xw.xDriver.translateModifiers(ev.State)

	keyEvent := gio.NewKeyboardEvent()
	keyEvent.Code = keyCode
	keyEvent.ModifierKey = modifierKey
	keyEvent.Repeat = false

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

	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = button
	mouseEvent.ModifierKey = modifierKey
	mouseEvent.Position = position

	xw.Publish(mouseEvent)
}

func (xw *x11Window) handleButtonRelease(ev xproto.ButtonReleaseEvent) {
	// Skip wheel events for release
	if ev.Detail >= 4 && ev.Detail <= 7 {
		return
	}

	button := xw.xDriver.translateMouseButton(ev.Detail)
	modifierKey := xw.xDriver.translateModifiers(ev.State)
	position := image.Point{X: int(ev.EventX), Y: int(ev.EventY)}

	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = button
	mouseEvent.ModifierKey = modifierKey
	mouseEvent.Position = position

	xw.Publish(mouseEvent)
}

func (xw *x11Window) handleWheelEvent(ev xproto.ButtonPressEvent) {
	deltaX, deltaY := xw.xDriver.translateWheelDelta(ev.Detail)

	modifierKey := xw.xDriver.translateModifiers(ev.State)
	position := image.Point{X: int(ev.EventX), Y: int(ev.EventY)}

	wheelEvent := gio.NewWheelEvent()
	wheelEvent.DeltaX = deltaX
	wheelEvent.DeltaY = deltaY
	wheelEvent.ModifierKey = modifierKey
	wheelEvent.Position = position

	xw.Publish(wheelEvent)
}

func (xw *x11Window) handleMotionNotify(ev xproto.MotionNotifyEvent) {
	modifierKey := xw.xDriver.translateModifiers(ev.State)
	position := image.Point{X: int(ev.EventX), Y: int(ev.EventY)}

	// Create mouse move event (no button pressed)
	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = gio.MouseButtonUnknown
	mouseEvent.ModifierKey = modifierKey
	mouseEvent.Position = position

	xw.Publish(mouseEvent)
}

func (xw *x11Window) handleMappingNotify(ev xproto.MappingNotifyEvent) {
	// Refresh keyboard mapping when it changes
	// This is typically handled by the X11 driver
	log.Printf("x11driver: keyboard mapping changed")
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
