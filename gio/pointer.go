package gio

type PointerType uint8

// PointerType constants represent the type of pointer
const (
	PointerTypeUnknown PointerType = iota // Unknown pointer type
	PointerTypeMouse                      // Mouse pointer
	PointerTypePen                        // Pen pointer
	PointerTypeTouch                      // Touch pointer
)

type PointerEvent struct {
	MouseEvent
	PointerID          int         // Unique identifier for the pointer
	Pressure           float64     // Pressure applied by the pointer (0.0 to 1.0)
	TangentialPressure float64     // Pressure applied tangentially (0.0 to 1.0)
	TiltX              float64     // Tilt angle in the X direction (degrees)
	TiltY              float64     // Tilt angle in the Y direction (degrees)
	Twist              float64     // Twist angle (degrees)
	Width              float64     // Width of the pointer contact area (CSS pixels)
	Height             float64     // Height of the pointer contact area (CSS pixels)
	PointerType        PointerType // Type of pointer (mouse, pen, touch)
	IsPrimary          bool        // Whether this is the primary pointer for the device
}
