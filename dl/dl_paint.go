package dl

import (
	"github.com/opensraph/sraph/geom"
)

// Paint defines how geometry is filled or stroked.
type Paint struct {
	color       Color
	blendMode   BlendMode
	strokeCap   geom.StrokeCap
	strokeJoin  geom.StrokeJoin
	width       float32
	miterLimit  float32
	isStroke    bool
	isAntiAlias bool
	style       PaintStyle
	// TODO: Add support for shaders, image filters, color filters, etc.
}

// NewPaint creates a new Paint with default values.
func NewPaint() Paint {
	return Paint{
		color:       ColorBlack,
		blendMode:   BlendModeSrcOver,
		strokeCap:   geom.StrokeCapButt,
		strokeJoin:  geom.StrokeJoinMiter,
		width:       1.0,
		miterLimit:  4.0,
		isStroke:    false,
		isAntiAlias: true,
	}
}

// Color returns the paint color.
func (p Paint) Color() Color {
	return p.color
}

// SetColor sets the paint color and returns the updated paint.
func (p Paint) SetColor(color Color) Paint {
	p.color = color
	return p
}

// BlendMode returns the blend mode.
func (p Paint) BlendMode() BlendMode {
	return p.blendMode
}

// SetBlendMode sets the blend mode and returns the updated paint.
func (p Paint) SetBlendMode(mode BlendMode) Paint {
	p.blendMode = mode
	return p
}

// StrokeCap returns the stroke cap style.
func (p Paint) StrokeCap() geom.StrokeCap {
	return p.strokeCap
}

// SetStrokeCap sets the stroke cap style and returns the updated paint.
func (p Paint) SetStrokeCap(cap geom.StrokeCap) Paint {
	p.strokeCap = cap
	return p
}

// StrokeJoin returns the stroke join style.
func (p Paint) StrokeJoin() geom.StrokeJoin {
	return p.strokeJoin
}

// SetStrokeJoin sets the stroke join style and returns the updated paint.
func (p Paint) SetStrokeJoin(join geom.StrokeJoin) Paint {
	p.strokeJoin = join
	return p
}

// StrokeWidth returns the stroke width.
func (p Paint) StrokeWidth() float32 {
	return p.width
}

// SetStrokeWidth sets the stroke width and returns the updated paint.
func (p Paint) SetStrokeWidth(width float32) Paint {
	p.width = width
	return p
}

// StrokeMiterLimit returns the stroke miter limit.
func (p Paint) StrokeMiterLimit() float32 {
	return p.miterLimit
}

// SetStrokeMiterLimit sets the stroke miter limit and returns the updated paint.
func (p Paint) SetStrokeMiterLimit(limit float32) Paint {
	p.miterLimit = limit
	return p
}

// IsStroke returns true if the paint is set to stroke mode.
func (p Paint) IsStroke() bool {
	return p.isStroke
}

// SetStroke sets the paint to stroke mode and returns the updated paint.
func (p Paint) SetStroke(stroke bool) Paint {
	p.isStroke = stroke
	return p
}

// IsAntiAlias returns true if anti-aliasing is enabled.
func (p Paint) IsAntiAlias() bool {
	return p.isAntiAlias
}

// SetAntiAlias sets anti-aliasing and returns the updated paint.
func (p Paint) SetAntiAlias(antiAlias bool) Paint {
	p.isAntiAlias = antiAlias
	return p
}

// Style returns the paint style (fill or stroke).
func (p Paint) Style() PaintStyle {
	return p.style
}

// SetStyle sets the paint style and returns the updated paint.
func (p Paint) SetStyle(style PaintStyle) Paint {
	p.style = style
	return p
}

// GetStrokeStyle returns the stroke style for use with geometry operations.
func (p Paint) GetStrokeStyle() geom.StrokeStyle[geom.F32] {
	return geom.StrokeStyle[geom.F32]{
		Width:      geom.F32(p.width),
		Cap:        p.strokeCap,
		Join:       p.strokeJoin,
		MiterLimit: geom.F32(p.miterLimit),
	}
}
