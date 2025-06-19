package contents

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// SolidRSuperellipseBlurContents provides blur rendering for rounded superellipse shapes
type SolidRSuperellipseBlurContents struct {
	BaseContents
	color         geom.Color
	path          geom.Path
	blurRadius    float32
	blurStyle     BlurStyle
	superellipseN float32 // Superellipse parameter (n > 0)
	cornerRadius  float32
	blendMode     render.BlendMode
	shadowOffset  geom.Point
	shadowOpacity float32
	antialiasing  bool
}

// NewSolidRSuperellipseBlurContents creates a new solid rounded superellipse blur contents
func NewSolidRSuperellipseBlurContents(color geom.Color, path geom.Path, blurRadius float32) *SolidRSuperellipseBlurContents {
	return &SolidRSuperellipseBlurContents{
		BaseContents:  NewBaseContents(),
		color:         color,
		path:          path,
		blurRadius:    blurRadius,
		blurStyle:     BlurStyleNormal,
		superellipseN: 2.0, // Default to circular (n=2)
		cornerRadius:  0,
		blendMode:     render.BlendModeSourceOver,
		shadowOffset:  geom.Point{X: 0, Y: 0},
		shadowOpacity: 1.0,
		antialiasing:  true,
	}
}

// SetColor sets the solid color
func (s *SolidRSuperellipseBlurContents) SetColor(color geom.Color) {
	s.color = color
}

// GetColor returns the current color
func (s *SolidRSuperellipseBlurContents) GetColor() geom.Color {
	return s.color
}

// SetPath sets the path for the superellipse
func (s *SolidRSuperellipseBlurContents) SetPath(path geom.Path) {
	s.path = path
}

// GetPath returns the current path
func (s *SolidRSuperellipseBlurContents) GetPath() geom.Path {
	return s.path
}

// SetBlurRadius sets the blur radius
func (s *SolidRSuperellipseBlurContents) SetBlurRadius(radius float32) {
	s.blurRadius = radius
}

// GetBlurRadius returns the current blur radius
func (s *SolidRSuperellipseBlurContents) GetBlurRadius() float32 {
	return s.blurRadius
}

// SetBlurStyle sets the blur style
func (s *SolidRSuperellipseBlurContents) SetBlurStyle(style BlurStyle) {
	s.blurStyle = style
}

// GetBlurStyle returns the current blur style
func (s *SolidRSuperellipseBlurContents) GetBlurStyle() BlurStyle {
	return s.blurStyle
}

// SetSuperellipseN sets the superellipse parameter
// n = 1: diamond/rhombus shape
// n = 2: circle/ellipse (default)
// n > 2: increasingly rectangular with rounded corners
func (s *SolidRSuperellipseBlurContents) SetSuperellipseN(n float32) {
	if n > 0 {
		s.superellipseN = n
	}
}

// GetSuperellipseN returns the current superellipse parameter
func (s *SolidRSuperellipseBlurContents) GetSuperellipseN() float32 {
	return s.superellipseN
}

// SetCornerRadius sets the corner radius
func (s *SolidRSuperellipseBlurContents) SetCornerRadius(radius float32) {
	s.cornerRadius = radius
}

// GetCornerRadius returns the current corner radius
func (s *SolidRSuperellipseBlurContents) GetCornerRadius() float32 {
	return s.cornerRadius
}

// SetBlendMode sets the blend mode
func (s *SolidRSuperellipseBlurContents) SetBlendMode(mode render.BlendMode) {
	s.blendMode = mode
}

// GetBlendMode returns the current blend mode
func (s *SolidRSuperellipseBlurContents) GetBlendMode() render.BlendMode {
	return s.blendMode
}

// SetShadowOffset sets the shadow offset
func (s *SolidRSuperellipseBlurContents) SetShadowOffset(offset geom.Point) {
	s.shadowOffset = offset
}

// GetShadowOffset returns the current shadow offset
func (s *SolidRSuperellipseBlurContents) GetShadowOffset() geom.Point {
	return s.shadowOffset
}

// SetShadowOpacity sets the shadow opacity
func (s *SolidRSuperellipseBlurContents) SetShadowOpacity(opacity float32) {
	s.shadowOpacity = opacity
}

// GetShadowOpacity returns the current shadow opacity
func (s *SolidRSuperellipseBlurContents) GetShadowOpacity() float32 {
	return s.shadowOpacity
}

// SetAntialiasing enables or disables antialiasing
func (s *SolidRSuperellipseBlurContents) SetAntialiasing(enabled bool) {
	s.antialiasing = enabled
}

// GetAntialiasing returns whether antialiasing is enabled
func (s *SolidRSuperellipseBlurContents) GetAntialiasing() bool {
	return s.antialiasing
}

// Render implements the Contents interface
func (s *SolidRSuperellipseBlurContents) Render(context *ContentContext, entity *Entity, pass *RenderPass) bool {
	// TODO: Implement solid rounded superellipse blur rendering
	// This would involve:
	// 1. Generating superellipse shape based on n parameter
	// 2. Creating a blur kernel based on radius and style
	// 3. Applying corner radius if specified
	// 4. Applying shadow offset and opacity
	// 5. Applying the blur effect with antialiasing
	// 6. Compositing with the specified blend mode
	return false
}

// GetBounds returns the bounds of this contents including blur and shadow
func (s *SolidRSuperellipseBlurContents) GetBounds() geom.Rect {
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
func (s *SolidRSuperellipseBlurContents) Clone() Contents {
	return &SolidRSuperellipseBlurContents{
		BaseContents:  s.BaseContents.Clone().(BaseContents),
		color:         s.color,
		path:          s.path,
		blurRadius:    s.blurRadius,
		blurStyle:     s.blurStyle,
		superellipseN: s.superellipseN,
		cornerRadius:  s.cornerRadius,
		blendMode:     s.blendMode,
		shadowOffset:  s.shadowOffset,
		shadowOpacity: s.shadowOpacity,
		antialiasing:  s.antialiasing,
	}
}

// GetCoverage returns the coverage area for this contents
func (s *SolidRSuperellipseBlurContents) GetCoverage(transform geom.Matrix) geom.Rect {
	return s.GetBounds()
}

// SolidRSuperellipseBlurFactory creates solid rounded superellipse blur contents instances
type SolidRSuperellipseBlurFactory struct{}

// CreateSolidRSuperellipseBlur creates a new solid rounded superellipse blur contents
func (f *SolidRSuperellipseBlurFactory) CreateSolidRSuperellipseBlur(color geom.Color, path geom.Path, blurRadius float32) Contents {
	return NewSolidRSuperellipseBlurContents(color, path, blurRadius)
}

// CreateSolidRSuperellipseBlurWithN creates a new solid rounded superellipse blur with specified n parameter
func (f *SolidRSuperellipseBlurFactory) CreateSolidRSuperellipseBlurWithN(color geom.Color, path geom.Path, blurRadius float32, n float32) Contents {
	contents := NewSolidRSuperellipseBlurContents(color, path, blurRadius)
	contents.SetSuperellipseN(n)
	return contents
}

// CreateSolidRSuperellipseBlurWithCorners creates a new solid rounded superellipse blur with corner settings
func (f *SolidRSuperellipseBlurFactory) CreateSolidRSuperellipseBlurWithCorners(color geom.Color, path geom.Path, blurRadius float32, n float32, cornerRadius float32) Contents {
	contents := NewSolidRSuperellipseBlurContents(color, path, blurRadius)
	contents.SetSuperellipseN(n)
	contents.SetCornerRadius(cornerRadius)
	return contents
}
