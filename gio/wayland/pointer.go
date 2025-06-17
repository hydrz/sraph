package wayland

import (
	"log"

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
	if w.IsClosed() {
		return
	}

	pointerEvent := gio.NewPointerEvent()
	pointerEvent.PointerType = gio.PointerTypeMouse
	pointerEvent.Position.X = int(x)
	pointerEvent.Position.Y = int(y)
	pointerEvent.KeyState = gio.KeyStatePressed

	if err := w.Publish(pointerEvent); err != nil {
		log.Printf("wayland: failed to publish pointer enter event: %v", err)
	}
}

func (w *WaylandWindow) handlePointerLeave(serial uint32) {
	if w.IsClosed() {
		return
	}

	pointerEvent := gio.NewPointerEvent()
	pointerEvent.PointerType = gio.PointerTypeMouse
	pointerEvent.KeyState = gio.KeyStateReleased

	if err := w.Publish(pointerEvent); err != nil {
		log.Printf("wayland: failed to publish pointer leave event: %v", err)
	}
}

func (w *WaylandWindow) handlePointerMotion(x, y float64, time uint32) {
	if w.IsClosed() {
		return
	}

	pointerEvent := gio.NewPointerEvent()
	pointerEvent.PointerType = gio.PointerTypeMouse
	pointerEvent.Position.X = int(x)
	pointerEvent.Position.Y = int(y)

	if err := w.Publish(pointerEvent); err != nil {
		log.Printf("wayland: failed to publish pointer motion event: %v", err)
	}
}

func (w *WaylandWindow) handlePointerButton(x, y float64, button uint32, pressed bool, serial, time uint32) {
	if w.IsClosed() {
		return
	}

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

	if err := w.Publish(pointerEvent); err != nil {
		log.Printf("wayland: failed to publish pointer button event: %v", err)
	}
}

func (w *WaylandWindow) handlePointerAxis(x, y, deltaX, deltaY float64, discreteX, discreteY int32, time uint32) {
	if w.IsClosed() {
		return
	}

	wheelEvent := gio.NewWheelEvent()
	wheelEvent.Position.X = int(x)
	wheelEvent.Position.Y = int(y)
	wheelEvent.DeltaX = deltaX
	wheelEvent.DeltaY = deltaY

	if err := w.Publish(wheelEvent); err != nil {
		log.Printf("wayland: failed to publish wheel event: %v", err)
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
		d.pointer.SetCursor(e.Serial, nil, 0, 0)
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
	pe := d.eventState.pointerEvent

	if d.focusedWindow == nil {
		// Reset event state
		d.eventState.pointerEvent = pointerEvent{
			surfaceX: pe.surfaceX,
			surfaceY: pe.surfaceY,
		}
		return
	}

	// Handle enter event
	if (pe.eventMask & pointerEventEnter) != 0 {
		d.focusedWindow.handlePointerEnter(pe.surfaceX, pe.surfaceY, pe.serial)
	}

	// Handle leave event
	if (pe.eventMask & pointerEventLeave) != 0 {
		d.focusedWindow.handlePointerLeave(pe.serial)
	}

	// Handle motion event
	if (pe.eventMask & pointerEventMotion) != 0 {
		d.focusedWindow.handlePointerMotion(pe.surfaceX, pe.surfaceY, pe.time)
	}

	// Handle button event
	if (pe.eventMask & pointerEventButton) != 0 {
		pressed := pe.state == uint32(client.PointerButtonStatePressed)
		d.focusedWindow.handlePointerButton(pe.surfaceX, pe.surfaceY, pe.button, pressed, pe.serial, pe.time)
	}

	// Handle axis events (scrolling)
	const axisEvents = pointerEventAxis | pointerEventAxisSource | pointerEventAxisStop | pointerEventAxisDiscrete
	if (pe.eventMask & axisEvents) != 0 {
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

	// Reset event state, keeping surface location
	d.eventState.pointerEvent = pointerEvent{
		surfaceX: pe.surfaceX,
		surfaceY: pe.surfaceY,
	}
}

// attachPointer sets up the pointer interface
func (d *WaylandDriver) attachPointer() {
	pointer, err := d.seat.GetPointer()
	if err != nil {
		log.Printf("wayland: failed to get pointer: %v", err)
		return
	}
	d.pointer = pointer

	// Set up pointer event handlers - following examples pattern
	d.pointer.SetEnterHandler(d.handlePointerEnter)
	d.pointer.SetLeaveHandler(d.handlePointerLeave)
	d.pointer.SetMotionHandler(d.handlePointerMotion)
	d.pointer.SetButtonHandler(d.handlePointerButton)
	d.pointer.SetAxisHandler(d.handlePointerAxis)
	d.pointer.SetAxisSourceHandler(d.handlePointerAxisSource)
	d.pointer.SetAxisStopHandler(d.handlePointerAxisStop)
	d.pointer.SetAxisDiscreteHandler(d.handlePointerAxisDiscrete)
	d.pointer.SetFrameHandler(d.handlePointerFrame)

	log.Printf("wayland: pointer interface registered")
}

// releasePointer releases the pointer interface
func (d *WaylandDriver) releasePointer() {
	if d.pointer != nil && d.seatVersion >= 3 {
		if err := d.pointer.Release(); err != nil {
			log.Printf("wayland: failed to release pointer: %v", err)
		}
	}
	d.pointer = nil
	log.Printf("wayland: pointer interface released")
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
