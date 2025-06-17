package gio

type PointerType uint8

// PointerType constants represent the type of pointer
const (
	PointerTypeUnknown PointerType = 1 << iota // Unknown pointer type
	PointerTypeMouse                           // Mouse pointer
	PointerTypePen                             // Pen pointer
	PointerTypeTouch                           // Touch pointer

	PointerTypeAll = PointerTypeMouse | PointerTypePen | PointerTypeTouch // All pointer types
)

func (pt PointerType) Contains(pointer PointerType) bool {
	return pt&pointer != 0
}

// String returns the string representation of the PointerType
func (pt PointerType) String() string {
	switch pt {
	case PointerTypeMouse:
		return "Mouse"
	case PointerTypePen:
		return "Pen"
	case PointerTypeTouch:
		return "Touch"
	case PointerTypeAll:
		return "All"
	default:
		return "Unknown"
	}
}

type MouseButton uint8

// MouseButton constants represent the mouse buttons
const (
	MouseButtonUnknown MouseButton = 1 << iota
	MouseButtonLeft                // Left mouse button (alias for MouseButton1)
	MouseButtonRight               // Right mouse button (alias for MouseButton2)
	MouseButtonMiddle              // Middle mouse button (alias for MouseButton3)
	MouseButtonBack                // Back button (usually on the side of the mouse)
	MouseButtonForward             // Forward button (usually on the side of the mouse)

	MouseButtonX1 // Additional button (if available)
	MouseButtonX2 // Additional button (if available)
)

func (mb MouseButton) Contains(button MouseButton) bool {
	return mb&button != 0
}

func (mb MouseButton) String() string {
	switch mb {
	case MouseButtonLeft:
		return "Left"
	case MouseButtonRight:
		return "Right"
	case MouseButtonMiddle:
		return "Middle"
	case MouseButtonBack:
		return "Back"
	case MouseButtonForward:
		return "Forward"
	case MouseButtonX1:
		return "X1"
	case MouseButtonX2:
		return "X2"
	default:
		return "Unknown"
	}
}
