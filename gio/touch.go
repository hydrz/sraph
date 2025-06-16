package gio

// Touch represents a single touch point on the screen
type Touch struct {
	Identifier int // Unique identifier for the touch point
	ScreenX    int // X coordinate relative to the screen
	ScreenY    int // Y coordinate relative to the screen
	ClientX    int // X coordinate relative to the viewport
	ClientY    int // Y coordinate relative to the viewport
	PageX      int // X coordinate relative to the document
	PageY      int // Y coordinate relative to the document
}

type TouchList []Touch

type TouchEvent struct {
	BaseEvent
	ChangedTouches TouchList   // List of touches that changed since the last event
	ModifierKey    ModifierKey // Modifier keys pressed during the event
	TargetTouches  TouchList   // List of touches currently on the target
	Touches        TouchList   // List of all touches on the screen
}
