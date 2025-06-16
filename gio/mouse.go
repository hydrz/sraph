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

func (mb MouseButton) Contains(button MouseButton) bool {
	return mb&button != 0
}

type MouseEvent struct {
	BaseEvent
	Button      MouseButton // The mouse button that was pressed
	ModifierKey ModifierKey // Modifier keys pressed during the event
	ClientX     int         // X coordinate relative to the viewport
	ClientY     int         // Y coordinate relative to the viewport
	MovementX   int         // Change in X since the last event
	MovementY   int         // Change in Y since the last event
	PageX       int         // X coordinate relative to the document
	PageY       int         // Y coordinate relative to the document
	ScreenX     int         // X coordinate relative to the screen
	ScreenY     int         // Y coordinate relative to the screen
}

type DeltaMode uint8

// DeltaMode constants represent the unit of measurement for deltas
const (
	DeltaModeUnknown DeltaMode = iota // Unknown delta mode
	DeltaModePixels                   // Pixels
	DeltaModeLines                    // Lines
	DeltaModePages                    // Pages
)

type WheelEvent struct {
	BaseEvent
	DeltaX    float64   // Horizontal scroll amount
	DeltaY    float64   // Vertical scroll amount
	DeltaZ    float64   // Depth scroll amount (if applicable)
	DeltaMode DeltaMode // Unit of measurement for deltas
}
