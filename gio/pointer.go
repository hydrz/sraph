package gio

type PointerType uint8

// PointerType constants represent the type of pointer
const (
	PointerTypeUnknown PointerType = iota // Unknown pointer type
	PointerTypeMouse                      // Mouse pointer
	PointerTypePen                        // Pen pointer
	PointerTypeTouch                      // Touch pointer
)

// String returns the string representation of the PointerType
func (pt PointerType) String() string {
	switch pt {
	case PointerTypeMouse:
		return "Mouse"
	case PointerTypePen:
		return "Pen"
	case PointerTypeTouch:
		return "Touch"
	default:
		return "Unknown"
	}
}
