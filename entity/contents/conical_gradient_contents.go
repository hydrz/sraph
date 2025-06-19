package contents

import (
	"math"
)

// ConicalGradientContents represents a conical (angular) gradient
type ConicalGradientContents interface {
	ColorSourceContents

	// SetCenter sets the center point of the gradient
	SetCenter(center Point)

	// GetCenter returns the center point
	GetCenter() Point

	// SetStartAngle sets the start angle in radians
	SetStartAngle(angle float32)

	// GetStartAngle returns the start angle
	GetStartAngle() float32

	// SetEndAngle sets the end angle in radians
	SetEndAngle(angle float32)

	// GetEndAngle returns the end angle
	GetEndAngle() float32

	// SetColors sets the gradient colors and stops
	SetColors(colors []Color, stops []float32)

	// GetColors returns the gradient colors
	GetColors() []Color

	// GetStops returns the gradient stops
	GetStops() []float32

	// SetTileMode sets how the gradient tiles
	SetTileMode(mode TileMode)

	// GetTileMode returns the tile mode
	GetTileMode() TileMode
}

// ConicalGradientContentsImpl implements ConicalGradientContents
type ConicalGradientContentsImpl struct {
	*ColorSourceContentsImpl
	center     Point
	startAngle float32
	endAngle   float32
	colors     []Color
	stops      []float32
	tileMode   TileMode
}

// NewConicalGradientContents creates a new conical gradient
func NewConicalGradientContents() ConicalGradientContents {
	return &ConicalGradientContentsImpl{
		ColorSourceContentsImpl: NewColorSourceContents(),
		center:                  Point{X: 0, Y: 0},
		startAngle:              0,
		endAngle:                2 * math.Pi,
		colors:                  make([]Color, 0),
		stops:                   make([]float32, 0),
		tileMode:                TileModeClamp,
	}
}

// SetCenter sets the center point of the gradient
func (c *ConicalGradientContentsImpl) SetCenter(center Point) {
	c.center = center
}

// GetCenter returns the center point
func (c *ConicalGradientContentsImpl) GetCenter() Point {
	return c.center
}

// SetStartAngle sets the start angle in radians
func (c *ConicalGradientContentsImpl) SetStartAngle(angle float32) {
	c.startAngle = angle
}

// GetStartAngle returns the start angle
func (c *ConicalGradientContentsImpl) GetStartAngle() float32 {
	return c.startAngle
}

// SetEndAngle sets the end angle in radians
func (c *ConicalGradientContentsImpl) SetEndAngle(angle float32) {
	c.endAngle = angle
}

// GetEndAngle returns the end angle
func (c *ConicalGradientContentsImpl) GetEndAngle() float32 {
	return c.endAngle
}

// SetColors sets the gradient colors and stops
func (c *ConicalGradientContentsImpl) SetColors(colors []Color, stops []float32) {
	c.colors = make([]Color, len(colors))
	copy(c.colors, colors)

	if len(stops) > 0 {
		c.stops = make([]float32, len(stops))
		copy(c.stops, stops)
	} else {
		// Generate evenly spaced stops
		c.stops = make([]float32, len(colors))
		for i := range c.stops {
			c.stops[i] = float32(i) / float32(len(colors)-1)
		}
	}
}

// GetColors returns the gradient colors
func (c *ConicalGradientContentsImpl) GetColors() []Color {
	return c.colors
}

// GetStops returns the gradient stops
func (c *ConicalGradientContentsImpl) GetStops() []float32 {
	return c.stops
}

// SetTileMode sets how the gradient tiles
func (c *ConicalGradientContentsImpl) SetTileMode(mode TileMode) {
	c.tileMode = mode
}

// GetTileMode returns the tile mode
func (c *ConicalGradientContentsImpl) GetTileMode() TileMode {
	return c.tileMode
}

// Render renders the conical gradient
func (c *ConicalGradientContentsImpl) Render(context ContentContext, entity Entity, pass RenderPass) bool {
	if len(c.colors) < 2 {
		return false
	}

	// TODO: Implement conical gradient rendering using shaders
	return false
}

// Clone creates a copy of the conical gradient contents
func (c *ConicalGradientContentsImpl) Clone() Contents {
	clone := &ConicalGradientContentsImpl{
		ColorSourceContentsImpl: c.ColorSourceContentsImpl,
		center:                  c.center,
		startAngle:              c.startAngle,
		endAngle:                c.endAngle,
		tileMode:                c.tileMode,
	}

	// Deep copy colors and stops
	clone.colors = make([]Color, len(c.colors))
	copy(clone.colors, c.colors)

	clone.stops = make([]float32, len(c.stops))
	copy(clone.stops, c.stops)

	return clone
}

// GetColorAt returns the interpolated color at the given angle
func (c *ConicalGradientContentsImpl) GetColorAt(angle float32) Color {
	if len(c.colors) == 0 {
		return Color{R: 0, G: 0, B: 0, A: 1}
	}

	if len(c.colors) == 1 {
		return c.ApplyColorModifiers(c.colors[0])
	}

	// Normalize angle to [0, 1] range
	normalizedAngle := (angle - c.startAngle) / (c.endAngle - c.startAngle)

	// Apply tile mode
	switch c.tileMode {
	case TileModeClamp:
		if normalizedAngle < 0 {
			normalizedAngle = 0
		} else if normalizedAngle > 1 {
			normalizedAngle = 1
		}
	case TileModeRepeat:
		normalizedAngle = normalizedAngle - float32(math.Floor(float64(normalizedAngle)))
	case TileModeMirror:
		normalizedAngle = normalizedAngle - float32(math.Floor(float64(normalizedAngle)))
		if int(math.Floor(float64(normalizedAngle*2)))%2 == 1 {
			normalizedAngle = 1 - normalizedAngle
		}
	case TileModeDecal:
		if normalizedAngle < 0 || normalizedAngle > 1 {
			return Color{R: 0, G: 0, B: 0, A: 0}
		}
	}

	// Find the color segment
	for i := 0; i < len(c.stops)-1; i++ {
		if normalizedAngle >= c.stops[i] && normalizedAngle <= c.stops[i+1] {
			// Interpolate between colors
			t := (normalizedAngle - c.stops[i]) / (c.stops[i+1] - c.stops[i])
			color1 := c.colors[i]
			color2 := c.colors[i+1]

			interpolated := Color{
				R: color1.R + t*(color2.R-color1.R),
				G: color1.G + t*(color2.G-color1.G),
				B: color1.B + t*(color2.B-color1.B),
				A: color1.A + t*(color2.A-color1.A),
			}

			return c.ApplyColorModifiers(interpolated)
		}
	}

	// Return last color if no segment found
	return c.ApplyColorModifiers(c.colors[len(c.colors)-1])
}
