package contents

import (
	"math"
)

// RadialGradientContents represents a radial gradient
type RadialGradientContents interface {
	ColorSourceContents

	// SetCenter sets the center point of the gradient
	SetCenter(center Point)

	// GetCenter returns the center point
	GetCenter() Point

	// SetRadius sets the radius of the gradient
	SetRadius(radius float32)

	// GetRadius returns the radius
	GetRadius() float32

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

// RadialGradientContentsImpl implements RadialGradientContents
type RadialGradientContentsImpl struct {
	*ColorSourceContentsImpl
	center   Point
	radius   float32
	colors   []Color
	stops    []float32
	tileMode TileMode
}

// NewRadialGradientContents creates a new radial gradient
func NewRadialGradientContents() RadialGradientContents {
	return &RadialGradientContentsImpl{
		ColorSourceContentsImpl: NewColorSourceContents(),
		center:                  Point{X: 0.5, Y: 0.5},
		radius:                  0.5,
		colors:                  make([]Color, 0),
		stops:                   make([]float32, 0),
		tileMode:                TileModeClamp,
	}
}

// SetCenter sets the center point of the gradient
func (r *RadialGradientContentsImpl) SetCenter(center Point) {
	r.center = center
}

// GetCenter returns the center point
func (r *RadialGradientContentsImpl) GetCenter() Point {
	return r.center
}

// SetRadius sets the radius of the gradient
func (r *RadialGradientContentsImpl) SetRadius(radius float32) {
	r.radius = radius
}

// GetRadius returns the radius
func (r *RadialGradientContentsImpl) GetRadius() float32 {
	return r.radius
}

// SetColors sets the gradient colors and stops
func (r *RadialGradientContentsImpl) SetColors(colors []Color, stops []float32) {
	r.colors = make([]Color, len(colors))
	copy(r.colors, colors)

	if len(stops) > 0 {
		r.stops = make([]float32, len(stops))
		copy(r.stops, stops)
	} else {
		// Generate evenly spaced stops
		r.stops = make([]float32, len(colors))
		for i := range r.stops {
			r.stops[i] = float32(i) / float32(len(colors)-1)
		}
	}
}

// GetColors returns the gradient colors
func (r *RadialGradientContentsImpl) GetColors() []Color {
	return r.colors
}

// GetStops returns the gradient stops
func (r *RadialGradientContentsImpl) GetStops() []float32 {
	return r.stops
}

// SetTileMode sets how the gradient tiles
func (r *RadialGradientContentsImpl) SetTileMode(mode TileMode) {
	r.tileMode = mode
}

// GetTileMode returns the tile mode
func (r *RadialGradientContentsImpl) GetTileMode() TileMode {
	return r.tileMode
}

// Render renders the radial gradient
func (r *RadialGradientContentsImpl) Render(context ContentContext, entity Entity, pass RenderPass) bool {
	if len(r.colors) < 2 {
		return false
	}

	// TODO: Implement radial gradient rendering using shaders
	return false
}

// Clone creates a copy of the radial gradient contents
func (r *RadialGradientContentsImpl) Clone() Contents {
	clone := &RadialGradientContentsImpl{
		ColorSourceContentsImpl: r.ColorSourceContentsImpl,
		center:                  r.center,
		radius:                  r.radius,
		tileMode:                r.tileMode,
	}

	// Deep copy colors and stops
	clone.colors = make([]Color, len(r.colors))
	copy(clone.colors, r.colors)

	clone.stops = make([]float32, len(r.stops))
	copy(clone.stops, r.stops)

	return clone
}

// GetColorAt returns the interpolated color at the given distance from center
func (r *RadialGradientContentsImpl) GetColorAt(distance float32) Color {
	if len(r.colors) == 0 {
		return Color{R: 0, G: 0, B: 0, A: 1}
	}

	if len(r.colors) == 1 {
		return r.ApplyColorModifiers(r.colors[0])
	}

	// Normalize distance to [0, 1] range
	normalizedDistance := distance / r.radius

	// Apply tile mode
	switch r.tileMode {
	case TileModeClamp:
		if normalizedDistance < 0 {
			normalizedDistance = 0
		} else if normalizedDistance > 1 {
			normalizedDistance = 1
		}
	case TileModeRepeat:
		normalizedDistance = normalizedDistance - float32(math.Floor(float64(normalizedDistance)))
	case TileModeMirror:
		normalizedDistance = normalizedDistance - float32(math.Floor(float64(normalizedDistance)))
		if int(math.Floor(float64(normalizedDistance*2)))%2 == 1 {
			normalizedDistance = 1 - normalizedDistance
		}
	case TileModeDecal:
		if normalizedDistance < 0 || normalizedDistance > 1 {
			return Color{R: 0, G: 0, B: 0, A: 0}
		}
	}

	// Find the color segment
	for i := 0; i < len(r.stops)-1; i++ {
		if normalizedDistance >= r.stops[i] && normalizedDistance <= r.stops[i+1] {
			// Interpolate between colors
			t := (normalizedDistance - r.stops[i]) / (r.stops[i+1] - r.stops[i])
			color1 := r.colors[i]
			color2 := r.colors[i+1]

			interpolated := Color{
				R: color1.R + t*(color2.R-color1.R),
				G: color1.G + t*(color2.G-color1.G),
				B: color1.B + t*(color2.B-color1.B),
				A: color1.A + t*(color2.A-color1.A),
			}

			return r.ApplyColorModifiers(interpolated)
		}
	}

	// Return last color if no segment found
	return r.ApplyColorModifiers(r.colors[len(r.colors)-1])
}

// DistanceFromCenter calculates the distance from a point to the gradient center
func (r *RadialGradientContentsImpl) DistanceFromCenter(point Point) float32 {
	dx := point.X - r.center.X
	dy := point.Y - r.center.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}
