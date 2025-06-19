package contents

import "math"

// LinearGradientContents represents a linear gradient
type LinearGradientContents interface {
	ColorSourceContents

	// SetStartPoint sets the start point of the gradient
	SetStartPoint(point Point)

	// GetStartPoint returns the start point
	GetStartPoint() Point

	// SetEndPoint sets the end point of the gradient
	SetEndPoint(point Point)

	// GetEndPoint returns the end point
	GetEndPoint() Point

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

// LinearGradientContentsImpl implements LinearGradientContents
type LinearGradientContentsImpl struct {
	*ColorSourceContentsImpl
	startPoint Point
	endPoint   Point
	colors     []Color
	stops      []float32
	tileMode   TileMode
}

// NewLinearGradientContents creates a new linear gradient
func NewLinearGradientContents() LinearGradientContents {
	return &LinearGradientContentsImpl{
		ColorSourceContentsImpl: NewColorSourceContents(),
		startPoint:              Point{X: 0, Y: 0},
		endPoint:                Point{X: 1, Y: 0},
		colors:                  make([]Color, 0),
		stops:                   make([]float32, 0),
		tileMode:                TileModeClamp,
	}
}

// SetStartPoint sets the start point of the gradient
func (l *LinearGradientContentsImpl) SetStartPoint(point Point) {
	l.startPoint = point
}

// GetStartPoint returns the start point
func (l *LinearGradientContentsImpl) GetStartPoint() Point {
	return l.startPoint
}

// SetEndPoint sets the end point of the gradient
func (l *LinearGradientContentsImpl) SetEndPoint(point Point) {
	l.endPoint = point
}

// GetEndPoint returns the end point
func (l *LinearGradientContentsImpl) GetEndPoint() Point {
	return l.endPoint
}

// SetColors sets the gradient colors and stops
func (l *LinearGradientContentsImpl) SetColors(colors []Color, stops []float32) {
	l.colors = make([]Color, len(colors))
	copy(l.colors, colors)

	if len(stops) > 0 {
		l.stops = make([]float32, len(stops))
		copy(l.stops, stops)
	} else {
		// Generate evenly spaced stops
		l.stops = make([]float32, len(colors))
		for i := range l.stops {
			l.stops[i] = float32(i) / float32(len(colors)-1)
		}
	}
}

// GetColors returns the gradient colors
func (l *LinearGradientContentsImpl) GetColors() []Color {
	return l.colors
}

// GetStops returns the gradient stops
func (l *LinearGradientContentsImpl) GetStops() []float32 {
	return l.stops
}

// SetTileMode sets how the gradient tiles
func (l *LinearGradientContentsImpl) SetTileMode(mode TileMode) {
	l.tileMode = mode
}

// GetTileMode returns the tile mode
func (l *LinearGradientContentsImpl) GetTileMode() TileMode {
	return l.tileMode
}

// Render renders the linear gradient
func (l *LinearGradientContentsImpl) Render(context ContentContext, entity Entity, pass RenderPass) bool {
	if len(l.colors) < 2 {
		return false
	}

	// TODO: Implement linear gradient rendering using shaders
	return false
}

// Clone creates a copy of the linear gradient contents
func (l *LinearGradientContentsImpl) Clone() Contents {
	clone := &LinearGradientContentsImpl{
		ColorSourceContentsImpl: l.ColorSourceContentsImpl,
		startPoint:              l.startPoint,
		endPoint:                l.endPoint,
		tileMode:                l.tileMode,
	}

	// Deep copy colors and stops
	clone.colors = make([]Color, len(l.colors))
	copy(clone.colors, l.colors)

	clone.stops = make([]float32, len(l.stops))
	copy(clone.stops, l.stops)

	return clone
}

// GetDirection returns the normalized direction vector of the gradient
func (l *LinearGradientContentsImpl) GetDirection() Point {
	dx := l.endPoint.X - l.startPoint.X
	dy := l.endPoint.Y - l.startPoint.Y
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if length == 0 {
		return Point{X: 1, Y: 0}
	}

	return Point{X: dx / length, Y: dy / length}
}

// GetLength returns the length of the gradient vector
func (l *LinearGradientContentsImpl) GetLength() float32 {
	dx := l.endPoint.X - l.startPoint.X
	dy := l.endPoint.Y - l.startPoint.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}
