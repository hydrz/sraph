package gio

// Touch represents a single touch point
type Touch struct {
	Identifier    int     // Unique identifier for this touch
	ScreenX       float64 // X coordinate relative to the screen
	ScreenY       float64 // Y coordinate relative to the screen
	ClientX       float64 // X coordinate relative to the client area
	ClientY       float64 // Y coordinate relative to the client area
	RadiusX       float64 // X radius of the ellipse that most closely describes the touch area
	RadiusY       float64 // Y radius of the ellipse that most closely describes the touch area
	RotationAngle float64 // Rotation angle of the ellipse in degrees
	Force         float64 // Pressure/force applied (0.0 to 1.0)
}

// TouchList represents a list of Touch objects
type TouchList struct {
	Touches []*Touch
}

// Length returns the number of touches in the list
func (tl TouchList) Length() int {
	return len(tl.Touches)
}

// Item returns the touch at the specified index
func (tl TouchList) Item(index int) *Touch {
	if index < 0 || index >= len(tl.Touches) {
		return nil
	}
	return tl.Touches[index]
}
