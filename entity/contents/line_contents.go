package contents

// LineContents represents line/path rendering contents
type LineContents interface {
	Contents

	// SetGeometry sets the line geometry
	SetGeometry(geometry LineGeometry)

	// GetGeometry returns the line geometry
	GetGeometry() LineGeometry

	// SetStrokeWidth sets the stroke width
	SetStrokeWidth(width float32)

	// GetStrokeWidth returns the stroke width
	GetStrokeWidth() float32

	// SetStrokeCap sets the stroke cap style
	SetStrokeCap(cap StrokeCap)

	// GetStrokeCap returns the stroke cap style
	GetStrokeCap() StrokeCap

	// SetStrokeJoin sets the stroke join style
	SetStrokeJoin(join StrokeJoin)

	// GetStrokeJoin returns the stroke join style
	GetStrokeJoin() StrokeJoin

	// SetMiterLimit sets the miter limit for joins
	SetMiterLimit(limit float32)

	// GetMiterLimit returns the miter limit
	GetMiterLimit() float32

	// SetColor sets the line color
	SetColor(color Color)

	// GetColor returns the line color
	GetColor() Color
}

// StrokeCap defines the style of line caps
type StrokeCap int

const (
	StrokeCapButt StrokeCap = iota
	StrokeCapRound
	StrokeCapSquare
)

// StrokeJoin defines the style of line joins
type StrokeJoin int

const (
	StrokeJoinMiter StrokeJoin = iota
	StrokeJoinRound
	StrokeJoinBevel
)

// LineGeometry represents geometry for lines and paths
type LineGeometry interface {
	Geometry

	// GetVertices returns the line vertices
	GetVertices() []Point

	// GetIndices returns the line indices
	GetIndices() []uint16

	// IsClosed returns whether the line is closed (forms a loop)
	IsClosed() bool
}

// LineContentsImpl implements LineContents
type LineContentsImpl struct {
	*ContentsImpl
	geometry    LineGeometry
	strokeWidth float32
	strokeCap   StrokeCap
	strokeJoin  StrokeJoin
	miterLimit  float32
	color       Color
}

// NewLineContents creates new line contents
func NewLineContents() LineContents {
	return &LineContentsImpl{
		ContentsImpl: &ContentsImpl{
			isValid: true,
		},
		strokeWidth: 1.0,
		strokeCap:   StrokeCapButt,
		strokeJoin:  StrokeJoinMiter,
		miterLimit:  4.0,
		color:       Color{R: 0, G: 0, B: 0, A: 1},
	}
}

// SetGeometry sets the line geometry
func (l *LineContentsImpl) SetGeometry(geometry LineGeometry) {
	l.geometry = geometry
}

// GetGeometry returns the line geometry
func (l *LineContentsImpl) GetGeometry() LineGeometry {
	return l.geometry
}

// SetStrokeWidth sets the stroke width
func (l *LineContentsImpl) SetStrokeWidth(width float32) {
	l.strokeWidth = width
}

// GetStrokeWidth returns the stroke width
func (l *LineContentsImpl) GetStrokeWidth() float32 {
	return l.strokeWidth
}

// SetStrokeCap sets the stroke cap style
func (l *LineContentsImpl) SetStrokeCap(cap StrokeCap) {
	l.strokeCap = cap
}

// GetStrokeCap returns the stroke cap style
func (l *LineContentsImpl) GetStrokeCap() StrokeCap {
	return l.strokeCap
}

// SetStrokeJoin sets the stroke join style
func (l *LineContentsImpl) SetStrokeJoin(join StrokeJoin) {
	l.strokeJoin = join
}

// GetStrokeJoin returns the stroke join style
func (l *LineContentsImpl) GetStrokeJoin() StrokeJoin {
	return l.strokeJoin
}

// SetMiterLimit sets the miter limit for joins
func (l *LineContentsImpl) SetMiterLimit(limit float32) {
	l.miterLimit = limit
}

// GetMiterLimit returns the miter limit
func (l *LineContentsImpl) GetMiterLimit() float32 {
	return l.miterLimit
}

// SetColor sets the line color
func (l *LineContentsImpl) SetColor(color Color) {
	l.color = color
}

// GetColor returns the line color
func (l *LineContentsImpl) GetColor() Color {
	return l.color
}

// Render renders the line contents
func (l *LineContentsImpl) Render(context ContentContext, entity Entity, pass RenderPass) bool {
	if l.geometry == nil {
		return false
	}

	// TODO: Implement line rendering with stroke styles
	return false
}

// GetCoverage returns the coverage bounds of the line
func (l *LineContentsImpl) GetCoverage(entity Entity) Rect {
	if l.geometry == nil {
		return Rect{}
	}

	// Get base coverage from geometry and expand by stroke width
	coverage := l.geometry.GetCoverage(entity.GetTransform())
	expansion := l.strokeWidth / 2.0

	return Rect{
		X:      coverage.X - expansion,
		Y:      coverage.Y - expansion,
		Width:  coverage.Width + 2*expansion,
		Height: coverage.Height + 2*expansion,
	}
}

// Clone creates a copy of the line contents
func (l *LineContentsImpl) Clone() Contents {
	clone := &LineContentsImpl{
		ContentsImpl: l.ContentsImpl.Clone().(*ContentsImpl),
		strokeWidth:  l.strokeWidth,
		strokeCap:    l.strokeCap,
		strokeJoin:   l.strokeJoin,
		miterLimit:   l.miterLimit,
		color:        l.color,
	}

	// Clone geometry if it exists
	if l.geometry != nil {
		// TODO: Implement geometry cloning
		clone.geometry = l.geometry
	}

	return clone
}

// GetStrokeExpansion returns how much the stroke expands beyond the base geometry
func (l *LineContentsImpl) GetStrokeExpansion() float32 {
	expansion := l.strokeWidth / 2.0

	// Add extra expansion for round caps and joins
	if l.strokeCap == StrokeCapRound || l.strokeJoin == StrokeJoinRound {
		expansion += l.strokeWidth / 2.0
	}

	// Add extra expansion for miter joins
	if l.strokeJoin == StrokeJoinMiter {
		expansion += l.strokeWidth * l.miterLimit / 2.0
	}

	return expansion
}
