package contents

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// SolidRRectLikeBlurContents provides blur rendering for rounded rectangle-like shapes
type SolidRRectLikeBlurContents struct {
	BaseContents
	color         geom.Color
	path          geom.Path
	blurRadius    float32
	blurStyle     BlurStyle
	cornerRadius  float32
	cornerMask    CornerMask
	blendMode     render.BlendMode
	shadowOffset  geom.Point
	shadowOpacity float32
}

// CornerMask defines which corners should be rounded
type CornerMask int

const (
	// CornerMaskNone no corners rounded
	CornerMaskNone CornerMask = 0
	// CornerMaskTopLeft round top-left corner
	CornerMaskTopLeft CornerMask = 1 << 0
	// CornerMaskTopRight round top-right corner
	CornerMaskTopRight CornerMask = 1 << 1
	// CornerMaskBottomLeft round bottom-left corner
	CornerMaskBottomLeft CornerMask = 1 << 2
	// CornerMaskBottomRight round bottom-right corner
	CornerMaskBottomRight CornerMask = 1 << 3
	// CornerMaskAll round all corners
	CornerMaskAll CornerMask = CornerMaskTopLeft | CornerMaskTopRight | CornerMaskBottomLeft | CornerMaskBottomRight
)

// NewSolidRRectLikeBlurContents creates a new solid rounded rectangle-like blur contents
func NewSolidRRectLikeBlurContents(color geom.Color, path geom.Path, blurRadius float32) *SolidRRectLikeBlurContents {
	return &SolidRRectLikeBlurContents{
		BaseContents:  NewBaseContents(),
		color:         color,
		path:          path,
		blurRadius:    blurRadius,
		blurStyle:     BlurStyleNormal,
		cornerRadius:  0,
		cornerMask:    CornerMaskAll,
		blendMode:     render.BlendModeSourceOver,
		shadowOffset:  geom.Point{X: 0, Y: 0},
		shadowOpacity: 1.0,
	}
}

// SetColor sets the solid color
func (s *SolidRRectLikeBlurContents) SetColor(color geom.Color) {
	s.color = color
}

// GetColor returns the current color
func (s *SolidRRectLikeBlurContents) GetColor() geom.Color {
	return s.color
}

// SetPath sets the path for the shape
func (s *SolidRRectLikeBlurContents) SetPath(path geom.Path) {
	s.path = path
}

// GetPath returns the current path
func (s *SolidRRectLikeBlurContents) GetPath() geom.Path {
	return s.path
}

// SetBlurRadius sets the blur radius
func (s *SolidRRectLikeBlurContents) SetBlurRadius(radius float32) {
	s.blurRadius = radius
}

// GetBlurRadius returns the current blur radius
func (s *SolidRRectLikeBlurContents) GetBlurRadius() float32 {
	return s.blurRadius
}

// SetBlurStyle sets the blur style
func (s *SolidRRectLikeBlurContents) SetBlurStyle(style BlurStyle) {
	s.blurStyle = style
}

// GetBlurStyle returns the current blur style
func (s *SolidRRectLikeBlurContents) GetBlurStyle() BlurStyle {
	return s.blurStyle
}

// SetCornerRadius sets the corner radius
func (s *SolidRRectLikeBlurContents) SetCornerRadius(radius float32) {
	s.cornerRadius = radius
}

// GetCornerRadius returns the current corner radius
func (s *SolidRRectLikeBlurContents) GetCornerRadius() float32 {
	return s.cornerRadius
}

// SetCornerMask sets which corners should be rounded
func (s *SolidRRectLikeBlurContents) SetCornerMask(mask CornerMask) {
	s.cornerMask = mask
}

// GetCornerMask returns the current corner mask
func (s *SolidRRectLikeBlurContents) GetCornerMask() CornerMask {
	return s.cornerMask
}

// SetBlendMode sets the blend mode
func (s *SolidRRectLikeBlurContents) SetBlendMode(mode render.BlendMode) {
	s.blendMode = mode
}

// GetBlendMode returns the current blend mode
func (s *SolidRRectLikeBlurContents) GetBlendMode() render.BlendMode {
	return s.blendMode
}

// SetShadowOffset sets the shadow offset
func (s *SolidRRectLikeBlurContents) SetShadowOffset(offset geom.Point) {
	s.shadowOffset = offset
}

// GetShadowOffset returns the current shadow offset
func (s *SolidRRectLikeBlurContents) GetShadowOffset() geom.Point {
	return s.shadowOffset
}

// SetShadowOpacity sets the shadow opacity
func (s *SolidRRectLikeBlurContents) SetShadowOpacity(opacity float32) {
	s.shadowOpacity = opacity
}

// GetShadowOpacity returns the current shadow opacity
func (s *SolidRRectLikeBlurContents) GetShadowOpacity() float32 {
	return s.shadowOpacity
}

// Render implements the Contents interface
func (s *SolidRRectLikeBlurContents) Render(context *ContentContext, entity *Entity, pass *RenderPass) bool {
	// TODO: Implement solid rounded rectangle-like blur rendering
	// This would involve:
	// 1. Creating a blur kernel based on radius and style
	// 2. Rendering the shape with selective corner rounding
	// 3. Applying shadow offset and opacity
	// 4. Applying the blur effect
	// 5. Compositing with the specified blend mode
	return false
}

// GetBounds returns the bounds of this contents including blur and shadow
func (s *SolidRRectLikeBlurContents) GetBounds() geom.Rect {
	bounds := s.path.GetBounds()

	// Expand bounds by blur radius
	expansion := s.blurRadius * 2

	// Account for shadow offset
	shadowExpansion := geom.Point{
		X: s.shadowOffset.X,
		Y: s.shadowOffset.Y,
	}

	minX := bounds.X - expansion
	minY := bounds.Y - expansion
	maxX := bounds.X + bounds.Width + expansion
	maxY := bounds.Y + bounds.Height + expansion

	// Adjust for shadow
	if shadowExpansion.X < 0 {
		minX += shadowExpansion.X
	} else {
		maxX += shadowExpansion.X
	}

	if shadowExpansion.Y < 0 {
		minY += shadowExpansion.Y
	} else {
		maxY += shadowExpansion.Y
	}

	return geom.Rect{
		X:      minX,
		Y:      minY,
		Width:  maxX - minX,
		Height: maxY - minY,
	}
}

// Clone creates a copy of this contents
func (s *SolidRRectLikeBlurContents) Clone() Contents {
	return &SolidRRectLikeBlurContents{
		BaseContents:  s.BaseContents.Clone().(BaseContents),
		color:         s.color,
		path:          s.path,
		blurRadius:    s.blurRadius,
		blurStyle:     s.blurStyle,
		cornerRadius:  s.cornerRadius,
		cornerMask:    s.cornerMask,
		blendMode:     s.blendMode,
		shadowOffset:  s.shadowOffset,
		shadowOpacity: s.shadowOpacity,
	}
}

// GetCoverage returns the coverage area for this contents
func (s *SolidRRectLikeBlurContents) GetCoverage(transform geom.Matrix) geom.Rect {
	return s.GetBounds()
}

// SolidRRectLikeBlurFactory creates solid rounded rectangle-like blur contents instances
type SolidRRectLikeBlurFactory struct{}

// CreateSolidRRectLikeBlur creates a new solid rounded rectangle-like blur contents
func (f *SolidRRectLikeBlurFactory) CreateSolidRRectLikeBlur(color geom.Color, path geom.Path, blurRadius float32) Contents {
	return NewSolidRRectLikeBlurContents(color, path, blurRadius)
}

// CreateSolidRRectLikeBlurWithCorners creates a new solid rounded rectangle-like blur with corner settings
func (f *SolidRRectLikeBlurFactory) CreateSolidRRectLikeBlurWithCorners(color geom.Color, path geom.Path, blurRadius float32, cornerRadius float32, cornerMask CornerMask) Contents {
	contents := NewSolidRRectLikeBlurContents(color, path, blurRadius)
	contents.SetCornerRadius(cornerRadius)
	contents.SetCornerMask(cornerMask)
	return contents
}
