package gio

type MouseButton uint8

// MouseButton constants represent the mouse buttons
const (
	MouseButtonUnknown MouseButton = iota
	MouseButtonLeft                = 1 << iota // Left mouse button (alias for MouseButton1)
	MouseButtonRight                           // Right mouse button (alias for MouseButton2)
	MouseButtonMiddle                          // Middle mouse button (alias for MouseButton3)
	MouseButtonBack                            // Back button (usually on the side of the mouse)
	MouseButtonForward                         // Forward button (usually on the side of the mouse)

	MouseButton1 = MouseButtonLeft    // Alias for left mouse button
	MouseButton2 = MouseButtonRight   // Alias for right mouse button
	MouseButton3 = MouseButtonMiddle  // Alias for middle mouse button
	MouseButton4 = MouseButtonBack    // Alias for back button
	MouseButton5 = MouseButtonForward // Alias for forward button
	MouseButton6 = 1 << 5             // Additional button (if available)
	MouseButton7 = 1 << 6             // Additional button (if available)
	MouseButton8 = 1 << 7             // Additional button (if available)
)

func (mb MouseButton) String() string {
	switch mb {
	case MouseButtonUnknown:
		return "Unknown"
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
	default:
		return "MouseButton" + string(mb)
	}
}

func (mb MouseButton) Contains(button MouseButton) bool {
	return mb&button != 0
}

type DeltaMode uint8

// DeltaMode constants represent the unit of measurement for deltas
const (
	DeltaModeUnknown DeltaMode = iota // Unknown delta mode
	DeltaModePixels                   // Pixels
	DeltaModeLines                    // Lines
	DeltaModePages                    // Pages
)
