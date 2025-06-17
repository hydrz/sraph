package x11

import (
	"encoding/binary"
	"fmt"
	"image"
	"log"

	"github.com/jezek/xgb/xproto"
	"github.com/opensraph/sraph/gio"
)

const (
	// XI2 event types
	XI_TouchBegin    = 18
	XI_TouchUpdate   = 19
	XI_TouchEnd      = 20
	XI_Motion        = 6
	XI_ButtonPress   = 4
	XI_ButtonRelease = 5

	// XI2 device types
	XIMasterPointer = 3
	XISlavePointer  = 4

	// XI2 touch modes
	XIDirectTouch    = 1
	XIDependentTouch = 2
)

// XI2Event represents a generic XI2 event
type XI2Event struct {
	Type      uint8
	Length    uint16
	EventType uint16
	Deviceid  uint16
	Time      uint32
	Detail    uint32
	Root      xproto.Window
	Event     xproto.Window
	Child     xproto.Window
	RootX     int32
	RootY     int32
	EventX    int32
	EventY    int32
	Flags     uint32
	Valuators []float64
	Mods      struct {
		Base      uint32
		Latched   uint32
		Locked    uint32
		Effective uint32
	}
	Group struct {
		Base      uint8
		Latched   uint8
		Locked    uint8
		Effective uint8
	}
	SourceId uint16
}

// getAllTouches returns all currently active touches
func (xw *x11Window) getAllTouches() gio.TouchList {
	xw.mu.RLock()
	defer xw.mu.RUnlock()

	touches := make([]*gio.Touch, 0, len(xw.activeTouches))
	for _, touch := range xw.activeTouches {
		touches = append(touches, touch)
	}

	return gio.TouchList{Touches: touches}
}

// handleXI2GenericEvent handles XI2 generic events
func (xw *x11Window) handleXI2GenericEvent(ev xproto.GeGenericEvent) {
	// Parse the XI2 event data
	xi2Event, err := parseXI2Event(ev.Bytes())
	if err != nil {
		log.Printf("x11driver: failed to parse XI2 event: %v", err)
		return
	}

	// Only handle events for our window
	if xi2Event.Event != xw.xWindow {
		return
	}

	switch xi2Event.EventType {
	case XI_TouchBegin:
		xw.handleXI2TouchBegin(xi2Event)
	case XI_TouchUpdate:
		xw.handleXI2TouchUpdate(xi2Event)
	case XI_TouchEnd:
		xw.handleXI2TouchEnd(xi2Event)
	case XI_ButtonPress:
		xw.handleXI2ButtonPress(xi2Event)
	case XI_ButtonRelease:
		xw.handleXI2ButtonRelease(xi2Event)
	case XI_Motion:
		xw.handleXI2Motion(xi2Event)
	default:
		log.Printf("x11driver: unhandled XI2 event type %d", xi2Event.EventType)
	}
}

// handleXI2TouchBegin handles XI2 touch begin events
func (xw *x11Window) handleXI2TouchBegin(ev *XI2Event) {
	touch := &gio.Touch{
		Identifier:    int(ev.Detail),
		ScreenX:       float64(ev.RootX),
		ScreenY:       float64(ev.RootY),
		ClientX:       float64(ev.EventX),
		ClientY:       float64(ev.EventY),
		RadiusX:       5.0, // Default radius
		RadiusY:       5.0,
		Force:         1.0, // Default force
		RotationAngle: 0.0,
	}

	// Extract pressure from valuators if available
	if len(ev.Valuators) > 2 {
		touch.Force = ev.Valuators[2] // Pressure is often the 3rd valuator
	}

	// Store active touch
	xw.mu.Lock()
	xw.activeTouches[int(ev.Detail)] = touch
	xw.mu.Unlock()

	// Create and publish touch event
	touchEvent := gio.NewTouchEvent()
	touchEvent.ChangedTouches = gio.TouchList{Touches: []*gio.Touch{touch}}
	touchEvent.Touches = xw.getAllTouches()
	touchEvent.ModifierKey = xw.xDriver.translateModifiers(uint16(ev.Mods.Effective))

	xw.Publish(touchEvent)
}

// handleXI2TouchUpdate handles XI2 touch update events
func (xw *x11Window) handleXI2TouchUpdate(ev *XI2Event) {
	xw.mu.Lock()
	touch, exists := xw.activeTouches[int(ev.Detail)]
	if !exists {
		xw.mu.Unlock()
		return
	}

	// Update touch position and properties
	touch.ScreenX = float64(ev.RootX)
	touch.ScreenY = float64(ev.RootY)
	touch.ClientX = float64(ev.EventX)
	touch.ClientY = float64(ev.EventY)

	// Update pressure from valuators if available
	if len(ev.Valuators) > 2 {
		touch.Force = ev.Valuators[2]
	}
	xw.mu.Unlock()

	// Create and publish touch event
	touchEvent := gio.NewTouchEvent()
	touchEvent.ChangedTouches = gio.TouchList{Touches: []*gio.Touch{touch}}
	touchEvent.Touches = xw.getAllTouches()
	touchEvent.ModifierKey = xw.xDriver.translateModifiers(uint16(ev.Mods.Effective))

	xw.Publish(touchEvent)
}

// handleXI2TouchEnd handles XI2 touch end events
func (xw *x11Window) handleXI2TouchEnd(ev *XI2Event) {
	xw.mu.Lock()
	touch, exists := xw.activeTouches[int(ev.Detail)]
	if !exists {
		xw.mu.Unlock()
		return
	}

	// Update final position
	touch.ScreenX = float64(ev.RootX)
	touch.ScreenY = float64(ev.RootY)
	touch.ClientX = float64(ev.EventX)
	touch.ClientY = float64(ev.EventY)

	// Remove from active touches
	delete(xw.activeTouches, int(ev.Detail))
	xw.mu.Unlock()

	// Create and publish touch event
	touchEvent := gio.NewTouchEvent()
	touchEvent.ChangedTouches = gio.TouchList{Touches: []*gio.Touch{touch}}
	touchEvent.Touches = xw.getAllTouches()
	touchEvent.ModifierKey = xw.xDriver.translateModifiers(uint16(ev.Mods.Effective))

	xw.Publish(touchEvent)
}

// handleXI2ButtonPress handles XI2 button press events for enhanced pointer support
func (xw *x11Window) handleXI2ButtonPress(ev *XI2Event) {
	// Check if this is from a stylus or enhanced pointer device
	isStylus := xw.isXI2PointerDevice(ev.Deviceid, ev.SourceId)

	if isStylus {
		// Create pointer event for stylus
		pointerEvent := gio.NewPointerEvent()
		pointerEvent.PointerID = int(ev.Deviceid)
		pointerEvent.Button = xw.xDriver.translateMouseButton(xproto.Button(ev.Detail))
		pointerEvent.ModifierKey = xw.xDriver.translateModifiers(uint16(ev.Mods.Effective))
		pointerEvent.Position = image.Point{X: int(ev.EventX), Y: int(ev.EventY)}

		// Extract pressure and tilt information from valuators
		xw.extractXI2PointerData(pointerEvent, ev.Valuators)
		pointerEvent.PointerType = gio.PointerTypePen
		pointerEvent.IsPrimary = ev.Deviceid == ev.SourceId

		xw.Publish(pointerEvent)
	} else {
		// Fall back to regular mouse event handling
		xw.handleButtonPress(xproto.ButtonPressEvent{
			Detail: xproto.Button(ev.Detail),
			EventX: int16(ev.EventX),
			EventY: int16(ev.EventY),
			State:  uint16(ev.Mods.Effective),
			Event:  ev.Event,
		})
	}
}

// handleXI2ButtonRelease handles XI2 button release events
func (xw *x11Window) handleXI2ButtonRelease(ev *XI2Event) {
	isStylus := xw.isXI2PointerDevice(ev.Deviceid, ev.SourceId)

	if isStylus {
		pointerEvent := gio.NewPointerEvent()
		pointerEvent.PointerID = int(ev.Deviceid)
		pointerEvent.Button = xw.xDriver.translateMouseButton(xproto.Button(ev.Detail))
		pointerEvent.ModifierKey = xw.xDriver.translateModifiers(uint16(ev.Mods.Effective))
		pointerEvent.Position = image.Point{X: int(ev.EventX), Y: int(ev.EventY)}

		xw.extractXI2PointerData(pointerEvent, ev.Valuators)
		pointerEvent.PointerType = gio.PointerTypePen
		pointerEvent.IsPrimary = ev.Deviceid == ev.SourceId

		xw.Publish(pointerEvent)
	} else {
		xw.handleButtonRelease(xproto.ButtonReleaseEvent{
			Detail: xproto.Button(ev.Detail),
			EventX: int16(ev.EventX),
			EventY: int16(ev.EventY),
			State:  uint16(ev.Mods.Effective),
			Event:  ev.Event,
		})
	}
}

// handleXI2Motion handles XI2 motion events
func (xw *x11Window) handleXI2Motion(ev *XI2Event) {
	isStylus := xw.isXI2PointerDevice(ev.Deviceid, ev.SourceId)

	if isStylus {
		pointerEvent := gio.NewPointerEvent()
		pointerEvent.PointerID = int(ev.Deviceid)
		pointerEvent.Button = gio.MouseButtonUnknown // No button pressed for motion
		pointerEvent.ModifierKey = xw.xDriver.translateModifiers(uint16(ev.Mods.Effective))
		pointerEvent.Position = image.Point{X: int(ev.EventX), Y: int(ev.EventY)}

		xw.extractXI2PointerData(pointerEvent, ev.Valuators)
		pointerEvent.PointerType = gio.PointerTypePen
		pointerEvent.IsPrimary = ev.Deviceid == ev.SourceId

		xw.Publish(pointerEvent)
	} else {
		xw.handleMotionNotify(xproto.MotionNotifyEvent{
			EventX: int16(ev.EventX),
			EventY: int16(ev.EventY),
			State:  uint16(ev.Mods.Effective),
			Event:  ev.Event,
		})
	}
}

// isXI2PointerDevice checks if a device ID corresponds to a stylus or enhanced pointer
func (xw *x11Window) isXI2PointerDevice(deviceID, sourceID uint16) bool {
	// Simple heuristics without XI2 device info:
	// - Different device and source IDs often indicate stylus
	// - Device IDs > 10 are often stylus devices (heuristic)
	return deviceID != sourceID || deviceID > 10
}

// extractXI2PointerData extracts pressure, tilt and other data from XI2 valuators
func (xw *x11Window) extractXI2PointerData(pointerEvent *gio.PointerEvent, valuators []float64) {
	// Default values
	pointerEvent.Pressure = 1.0
	pointerEvent.TiltX = 0.0
	pointerEvent.TiltY = 0.0
	pointerEvent.Twist = 0.0
	pointerEvent.Width = 1.0
	pointerEvent.Height = 1.0

	// Extract valuator data
	// Common mapping: [X, Y, Pressure, TiltX, TiltY, Twist]
	if len(valuators) > 2 {
		pressure := valuators[2]
		if pressure >= 0 && pressure <= 1 {
			pointerEvent.Pressure = pressure
		}
	}

	if len(valuators) > 3 {
		pointerEvent.TiltX = valuators[3]
	}

	if len(valuators) > 4 {
		pointerEvent.TiltY = valuators[4]
	}

	if len(valuators) > 5 {
		pointerEvent.Twist = valuators[5]
	}
}

// parseXI2Event parses a raw XI2 event from bytes
func parseXI2Event(data []byte) (*XI2Event, error) {
	if len(data) < 32 {
		return nil, fmt.Errorf("XI2 event data too short: %d bytes", len(data))
	}

	event := &XI2Event{}

	// Parse basic event structure
	event.Type = data[0]
	event.Length = binary.LittleEndian.Uint16(data[2:4])
	event.EventType = binary.LittleEndian.Uint16(data[4:6])
	event.Deviceid = binary.LittleEndian.Uint16(data[6:8])
	event.Time = binary.LittleEndian.Uint32(data[8:12])
	event.Detail = binary.LittleEndian.Uint32(data[12:16])
	event.Root = xproto.Window(binary.LittleEndian.Uint32(data[16:20]))
	event.Event = xproto.Window(binary.LittleEndian.Uint32(data[20:24]))
	event.Child = xproto.Window(binary.LittleEndian.Uint32(data[24:28]))

	if len(data) >= 48 {
		event.RootX = int32(binary.LittleEndian.Uint32(data[28:32]))
		event.RootY = int32(binary.LittleEndian.Uint32(data[32:36]))
		event.EventX = int32(binary.LittleEndian.Uint32(data[36:40]))
		event.EventY = int32(binary.LittleEndian.Uint32(data[40:44]))
		event.Flags = binary.LittleEndian.Uint32(data[44:48])
	}

	// Parse modifier state if available
	if len(data) >= 64 {
		event.Mods.Base = binary.LittleEndian.Uint32(data[48:52])
		event.Mods.Latched = binary.LittleEndian.Uint32(data[52:56])
		event.Mods.Locked = binary.LittleEndian.Uint32(data[56:60])
		event.Mods.Effective = binary.LittleEndian.Uint32(data[60:64])
	}

	// Parse group state if available
	if len(data) >= 68 {
		event.Group.Base = data[64]
		event.Group.Latched = data[65]
		event.Group.Locked = data[66]
		event.Group.Effective = data[67]
	}

	// Parse source ID if available
	if len(data) >= 70 {
		event.SourceId = binary.LittleEndian.Uint16(data[68:70])
	}

	return event, nil
}
