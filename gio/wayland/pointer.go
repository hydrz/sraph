package wayland

import (
	"log/slog"

	"github.com/opensraph/sraph/gio"
	"github.com/rajveermalviya/go-wayland/wayland/client"
)

const (
	pointerEventEnter        = 1 << 0
	pointerEventLeave        = 1 << 1
	pointerEventMotion       = 1 << 2
	pointerEventButton       = 1 << 3
	pointerEventAxis         = 1 << 4
	pointerEventAxisSource   = 1 << 5
	pointerEventAxisStop     = 1 << 6
	pointerEventAxisDiscrete = 1 << 7
)

// Linux input event codes
const (
	BtnLeft   = 0x110
	BtnRight  = 0x111
	BtnMiddle = 0x112
	BtnSide   = 0x113
	BtnExtra  = 0x114
)

// Pointer event handlers for WaylandWindow
func (w *WaylandWindow) handlePointerEnter(x, y float64, serial uint32) {
	pointerEvent := gio.NewPointerEvent()
	pointerEvent.PointerType = gio.PointerTypeMouse
	pointerEvent.Position.X = int(x)
	pointerEvent.Position.Y = int(y)
	pointerEvent.KeyState = gio.KeyStatePressed

	w.publishPointerEvent(pointerEvent)
}

func (w *WaylandWindow) handlePointerLeave(serial uint32) {
	pointerEvent := gio.NewPointerEvent()
	pointerEvent.PointerType = gio.PointerTypeMouse
	pointerEvent.KeyState = gio.KeyStateReleased

	w.publishPointerEvent(pointerEvent)
}

func (w *WaylandWindow) handlePointerMotion(x, y float64, time uint32) {
	pointerEvent := gio.NewPointerEvent()
	pointerEvent.PointerType = gio.PointerTypeMouse
	pointerEvent.Position.X = int(x)
	pointerEvent.Position.Y = int(y)

	w.publishPointerEvent(pointerEvent)
}

func (w *WaylandWindow) handlePointerButton(x, y float64, button uint32, pressed bool, serial, time uint32) {
	pointerEvent := gio.NewPointerEvent()
	pointerEvent.PointerType = gio.PointerTypeMouse
	pointerEvent.Position.X = int(x)
	pointerEvent.Position.Y = int(y)
	pointerEvent.MouseButton = waylandButtonToGio(button)

	if pressed {
		pointerEvent.KeyState = gio.KeyStatePressed
	} else {
		pointerEvent.KeyState = gio.KeyStateReleased
	}

	w.publishPointerEvent(pointerEvent)
}

func (w *WaylandWindow) handlePointerAxis(x, y, deltaX, deltaY float64, discreteX, discreteY int32, time uint32) {
	wheelEvent := gio.NewWheelEvent()
	wheelEvent.Position.X = int(x)
	wheelEvent.Position.Y = int(y)
	wheelEvent.DeltaX = deltaX
	wheelEvent.DeltaY = deltaY

	if w.IsClosed() {
		return
	}

	if err := w.Publish(wheelEvent); err != nil {
		slog.Error("failed to publish wheel event", "error", err)
	}
}

// publishPointerEvent publishes pointer event with error handling
func (w *WaylandWindow) publishPointerEvent(event *gio.PointerEvent) {
	if w.IsClosed() {
		return
	}

	if err := w.Publish(event); err != nil {
		slog.Error("failed to publish pointer event", "error", err)
	}
}

// Driver pointer event handlers
func (d *WaylandDriver) handlePointerEnter(e client.PointerEnterEvent) {
	d.updateSerial(e.Serial)

	d.eventState.pointerEvent.eventMask |= pointerEventEnter
	d.eventState.pointerEvent.serial = e.Serial
	d.eventState.pointerEvent.surfaceX = e.SurfaceX
	d.eventState.pointerEvent.surfaceY = e.SurfaceY

	// Find window for this surface
	if window := d.findWindowBySurface(e.Surface); window != nil {
		d.focusedWindow = window

		// Set cursor if available
		if d.cursorTheme != nil && d.pointer != nil {
			d.setCursor(e.Serial, "left_ptr")
		}
	}
}

func (d *WaylandDriver) handlePointerLeave(e client.PointerLeaveEvent) {
	d.eventState.pointerEvent.eventMask |= pointerEventLeave
	d.eventState.pointerEvent.serial = e.Serial

	// Clear cursor
	if d.pointer != nil {
		if err := d.pointer.SetCursor(e.Serial, nil, 0, 0); err != nil {
			slog.Error("failed to clear cursor", "error", err)
		}
	}
}

func (d *WaylandDriver) handlePointerMotion(e client.PointerMotionEvent) {
	d.eventState.pointerEvent.eventMask |= pointerEventMotion
	d.eventState.pointerEvent.time = e.Time
	d.eventState.pointerEvent.surfaceX = e.SurfaceX
	d.eventState.pointerEvent.surfaceY = e.SurfaceY
}

func (d *WaylandDriver) handlePointerButton(e client.PointerButtonEvent) {
	d.updateSerial(e.Serial)

	d.eventState.pointerEvent.eventMask |= pointerEventButton
	d.eventState.pointerEvent.serial = e.Serial
	d.eventState.pointerEvent.time = e.Time
	d.eventState.pointerEvent.button = e.Button
	d.eventState.pointerEvent.state = e.State
}

func (d *WaylandDriver) handlePointerAxis(e client.PointerAxisEvent) {
	d.eventState.pointerEvent.eventMask |= pointerEventAxis
	d.eventState.pointerEvent.time = e.Time
	d.eventState.pointerEvent.axes[e.Axis].valid = true
	d.eventState.pointerEvent.axes[e.Axis].value = e.Value
}

func (d *WaylandDriver) handlePointerAxisSource(e client.PointerAxisSourceEvent) {
	d.eventState.pointerEvent.eventMask |= pointerEventAxisSource
	d.eventState.pointerEvent.axisSource = e.AxisSource
}

func (d *WaylandDriver) handlePointerAxisStop(e client.PointerAxisStopEvent) {
	d.eventState.pointerEvent.eventMask |= pointerEventAxisStop
	d.eventState.pointerEvent.time = e.Time
	d.eventState.pointerEvent.axes[e.Axis].valid = true
}

func (d *WaylandDriver) handlePointerAxisDiscrete(e client.PointerAxisDiscreteEvent) {
	d.eventState.pointerEvent.eventMask |= pointerEventAxisDiscrete
	d.eventState.pointerEvent.axes[e.Axis].valid = true
	d.eventState.pointerEvent.axes[e.Axis].discrete = e.Discrete
}

func (d *WaylandDriver) handlePointerFrame(e client.PointerFrameEvent) {
	d.processPointerEvents()
}

// processPointerEvents processes batched pointer events
func (d *WaylandDriver) processPointerEvents() {
	pe := d.eventState.pointerEvent

	if d.focusedWindow == nil {
		d.resetPointerEvent(&pe)
		return
	}

	// Process events in order
	eventHandlers := []struct {
		mask    int
		handler func()
	}{
		{pointerEventEnter, func() {
			d.focusedWindow.handlePointerEnter(pe.surfaceX, pe.surfaceY, pe.serial)
		}},
		{pointerEventLeave, func() {
			d.focusedWindow.handlePointerLeave(pe.serial)
		}},
		{pointerEventMotion, func() {
			d.focusedWindow.handlePointerMotion(pe.surfaceX, pe.surfaceY, pe.time)
		}},
		{pointerEventButton, func() {
			pressed := pe.state == uint32(client.PointerButtonStatePressed)
			d.focusedWindow.handlePointerButton(pe.surfaceX, pe.surfaceY, pe.button, pressed, pe.serial, pe.time)
		}},
	}

	for _, handler := range eventHandlers {
		if (pe.eventMask & handler.mask) != 0 {
			handler.handler()
		}
	}

	// Handle axis events
	d.processAxisEvents(&pe)

	d.resetPointerEvent(&pe)
}

// processAxisEvents processes scroll wheel events
func (d *WaylandDriver) processAxisEvents(pe *pointerEvent) {
	const axisEvents = pointerEventAxis | pointerEventAxisSource | pointerEventAxisStop | pointerEventAxisDiscrete
	if (pe.eventMask & axisEvents) == 0 {
		return
	}

	for axis := 0; axis < 2; axis++ {
		if !pe.axes[axis].valid {
			continue
		}

		deltaX, deltaY := 0.0, 0.0
		discreteX, discreteY := int32(0), int32(0)

		if axis == int(client.PointerAxisHorizontalScroll) {
			deltaX = pe.axes[axis].value
			discreteX = pe.axes[axis].discrete
		} else if axis == int(client.PointerAxisVerticalScroll) {
			deltaY = pe.axes[axis].value
			discreteY = pe.axes[axis].discrete
		}

		d.focusedWindow.handlePointerAxis(pe.surfaceX, pe.surfaceY, deltaX, deltaY, discreteX, discreteY, pe.time)
	}
}

// resetPointerEvent resets event state
func (d *WaylandDriver) resetPointerEvent(pe *pointerEvent) {
	d.eventState.pointerEvent = pointerEvent{
		surfaceX: pe.surfaceX,
		surfaceY: pe.surfaceY,
	}
}

// attachPointer sets up the pointer interface with better error handling
func (d *WaylandDriver) attachPointer() {
	pointer, err := d.seat.GetPointer()
	if err != nil {
		slog.Error("failed to get pointer", "error", err)
		return
	}
	d.pointer = pointer

	// Set up pointer event handlers
	d.pointer.SetEnterHandler(d.handlePointerEnter)
	d.pointer.SetLeaveHandler(d.handlePointerLeave)
	d.pointer.SetMotionHandler(d.handlePointerMotion)
	d.pointer.SetButtonHandler(d.handlePointerButton)
	d.pointer.SetAxisHandler(d.handlePointerAxis)
	d.pointer.SetAxisSourceHandler(d.handlePointerAxisSource)
	d.pointer.SetAxisStopHandler(d.handlePointerAxisStop)
	d.pointer.SetAxisDiscreteHandler(d.handlePointerAxisDiscrete)
	d.pointer.SetFrameHandler(d.handlePointerFrame)

	slog.Info("pointer interface registered")
}

// releasePointer releases the pointer interface
func (d *WaylandDriver) releasePointer() {
	if d.pointer != nil && d.seatVersion >= 3 {
		if err := d.pointer.Release(); err != nil {
			slog.Error("failed to release pointer", "error", err)
		}
	}
	d.pointer = nil
	slog.Info("pointer interface released")
}

// waylandButtonToGio converts Wayland button to gio MouseButton
func waylandButtonToGio(button uint32) gio.MouseButton {
	switch button {
	case BtnLeft:
		return gio.MouseButtonLeft
	case BtnRight:
		return gio.MouseButtonRight
	case BtnMiddle:
		return gio.MouseButtonMiddle
	case BtnSide:
		return gio.MouseButtonBack
	case BtnExtra:
		return gio.MouseButtonForward
	default:
		return gio.MouseButtonUnknown
	}
}

// setCursor sets the cursor for the pointer
func (d *WaylandDriver) setCursor(serial uint32, name string) {
	if d.cursorTheme == nil || d.pointer == nil {
		return
	}

	cursor := d.cursorTheme.GetCursor(name)
	if cursor == nil {
		return
	}

	image := cursor.Images[0]

	surface, err := d.compositor.CreateSurface()
	if err != nil {
		slog.Error("failed to create cursor surface", "error", err)
		return
	}

	buffer, err := image.GetBuffer()
	if err != nil {
		slog.Error("failed to get cursor buffer", "error", err)
		return
	}

	if buffer != nil {
		if err := surface.Attach(buffer, 0, 0); err != nil {
			slog.Error("failed to attach cursor buffer", "error", err)
			return
		}
		if err := surface.Damage(0, 0, int32(image.Width), int32(image.Height)); err != nil {
			slog.Error("failed to damage cursor surface", "error", err)
		}
		if err := surface.Commit(); err != nil {
			slog.Error("failed to commit cursor surface", "error", err)
			return
		}

		hotspotX := int32(image.HotspotX)
		hotspotY := int32(image.HotspotY)
		if err := d.pointer.SetCursor(serial, surface, hotspotX, hotspotY); err != nil {
			slog.Error("failed to set cursor", "error", err)
		}
	}
}
