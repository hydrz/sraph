package contents

import (
	"math"
)

// SweepGradientContents represents a sweep (angular) gradient
type SweepGradientContents interface {
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

// SweepGradientContentsImpl implements SweepGradientContents
type SweepGradientContentsImpl struct {
	*ColorSourceContentsImpl
	center     Point
	startAngle float32
	endAngle   float32
	colors     []Color
	stops      []float32
	tileMode   TileMode
}

// NewSweepGradientContents creates a new sweep gradient
func NewSweepGradientContents() SweepGradientContents {
	return &SweepGradientContentsImpl{
		ColorSourceContentsImpl: NewColorSourceContents(),
		center:                  Point{X: 0.5, Y: 0.5},
		startAngle:              0,
		endAngle:                2 * math.Pi,
		colors:                  make([]Color, 0),
		stops:                   make([]float32, 0),
		tileMode:                TileModeClamp,
	}
}

// SetCenter sets the center point of the gradient
func (s *SweepGradientContentsImpl) SetCenter(center Point) {
	s.center = center
}

// GetCenter returns the center point
func (s *SweepGradientContentsImpl) GetCenter() Point {
	return s.center
}

// SetStartAngle sets the start angle in radians
func (s *SweepGradientContentsImpl) SetStartAngle(angle float32) {
	s.startAngle = angle
}

// GetStartAngle returns the start angle
func (s *SweepGradientContentsImpl) GetStartAngle() float32 {
	return s.startAngle
}

// SetEndAngle sets the end angle in radians
func (s *SweepGradientContentsImpl) SetEndAngle(angle float32) {
	s.endAngle = angle
}

// GetEndAngle returns the end angle
func (s *SweepGradientContentsImpl) GetEndAngle() float32 {
	return s.endAngle
}

// SetColors sets the gradient colors and stops
func (s *SweepGradientContentsImpl) SetColors(colors []Color, stops []float32) {
	s.colors = make([]Color, len(colors))
	copy(s.colors, colors)

	if len(stops) > 0 {
		s.stops = make([]float32, len(stops))
		copy(s.stops, stops)
	} else {
		// Generate evenly spaced stops
		s.stops = make([]float32, len(colors))
		for i := range s.stops {
			s.stops[i] = float32(i) / float32(len(colors)-1)
		}
	}
}

// GetColors returns the gradient colors
func (s *SweepGradientContentsImpl) GetColors() []Color {
	return s.colors
}

// GetStops returns the gradient stops
func (s *SweepGradientContentsImpl) GetStops() []float32 {
	return s.stops
}

// SetTileMode sets how the gradient tiles
func (s *SweepGradientContentsImpl) SetTileMode(mode TileMode) {
	s.tileMode = mode
}

// GetTileMode returns the tile mode
func (s *SweepGradientContentsImpl) GetTileMode() TileMode {
	return s.tileMode
}

// Render renders the sweep gradient
func (s *SweepGradientContentsImpl) Render(context ContentContext, entity Entity, pass RenderPass) bool {
	if len(s.colors) < 2 {
		return false
	}

	// TODO: Implement sweep gradient rendering using shaders
	return false
}

// Clone creates a copy of the sweep gradient contents
func (s *SweepGradientContentsImpl) Clone() Contents {
	clone := &SweepGradientContentsImpl{
		ColorSourceContentsImpl: s.ColorSourceContentsImpl,
		center:                  s.center,
		startAngle:              s.startAngle,
		endAngle:                s.endAngle,
		tileMode:                s.tileMode,
	}

	// Deep copy colors and stops
	clone.colors = make([]Color, len(s.colors))
	copy(clone.colors, s.colors)

	clone.stops = make([]float32, len(s.stops))
	copy(clone.stops, s.stops)

	return clone
}

// GetAngleFromCenter calculates the angle from center to a point
func (s *SweepGradientContentsImpl) GetAngleFromCenter(point Point) float32 {
	dx := point.X - s.center.X
	dy := point.Y - s.center.Y
	return float32(math.Atan2(float64(dy), float64(dx)))
}

// GetColorAtAngle returns the interpolated color at the given angle
func (s *SweepGradientContentsImpl) GetColorAtAngle(angle float32) Color {
	if len(s.colors) == 0 {
		return Color{R: 0, G: 0, B: 0, A: 1}
	}

	if len(s.colors) == 1 {
		return s.ApplyColorModifiers(s.colors[0])
	}

	// Normalize angle to [0, 1] range
	normalizedAngle := (angle - s.startAngle) / (s.endAngle - s.startAngle)

	// Apply tile mode
	switch s.tileMode {
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
	for i := 0; i < len(s.stops)-1; i++ {
		if normalizedAngle >= s.stops[i] && normalizedAngle <= s.stops[i+1] {
			// Interpolate between colors
			t := (normalizedAngle - s.stops[i]) / (s.stops[i+1] - s.stops[i])
			color1 := s.colors[i]
			color2 := s.colors[i+1]

			interpolated := Color{
				R: color1.R + t*(color2.R-color1.R),
				G: color1.G + t*(color2.G-color1.G),
				B: color1.B + t*(color2.B-color1.B),
				A: color1.A + t*(color2.A-color1.A),
			}

			return s.ApplyColorModifiers(interpolated)
		}
	}

	// Return last color if no segment found
	return s.ApplyColorModifiers(s.colors[len(s.colors)-1])
}
