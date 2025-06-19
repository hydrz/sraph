package contents

import (
	"github.com/opensraph/sraph/geom"
)

// SolidRRectBlurContents provides solid rounded rectangle blur rendering
type SolidRRectBlurContents struct {
	BaseContents
	color       geom.Color
	path        geom.Path
	blurRadius  float32
	blurStyle   BlurStyle
	cornerRadii [4]float32 // top-left, top-right, bottom-right, bottom-left
}

// BlurStyle defines how blur is applied
type BlurStyle int

const (
	// BlurStyleNormal applies blur normally
	BlurStyleNormal BlurStyle = iota
	// BlurStyleSolid applies solid blur
	BlurStyleSolid
	// BlurStyleOuter applies outer blur only
	BlurStyleOuter
	// BlurStyleInner applies inner blur only
	BlurStyleInner
)

// NewSolidRRectBlurContents creates a new solid rounded rectangle blur contents
func NewSolidRRectBlurContents(color geom.Color, path geom.Path, blurRadius float32) *SolidRRectBlurContents {
	return &SolidRRectBlurContents{
		BaseContents: NewBaseContents(),
		color:        color,
		path:         path,
		blurRadius:   blurRadius,
		blurStyle:    BlurStyleNormal,
		cornerRadii:  [4]float32{0, 0, 0, 0},
	}
}

// SetColor sets the solid color
func (s *SolidRRectBlurContents) SetColor(color geom.Color) {
	s.color = color
}

// GetColor returns the current color
func (s *SolidRRectBlurContents) GetColor() geom.Color {
	return s.color
}

// SetPath sets the path for the rounded rectangle
func (s *SolidRRectBlurContents) SetPath(path geom.Path) {
	s.path = path
}

// GetPath returns the current path
func (s *SolidRRectBlurContents) GetPath() geom.Path {
	return s.path
}

// SetBlurRadius sets the blur radius
func (s *SolidRRectBlurContents) SetBlurRadius(radius float32) {
	s.blurRadius = radius
}

// GetBlurRadius returns the current blur radius
func (s *SolidRRectBlurContents) GetBlurRadius() float32 {
	return s.blurRadius
}

// SetBlurStyle sets the blur style
func (s *SolidRRectBlurContents) SetBlurStyle(style BlurStyle) {
	s.blurStyle = style
}

// GetBlurStyle returns the current blur style
func (s *SolidRRectBlurContents) GetBlurStyle() BlurStyle {
	return s.blurStyle
}

// SetCornerRadii sets the corner radii for the rounded rectangle
func (s *SolidRRectBlurContents) SetCornerRadii(topLeft, topRight, bottomRight, bottomLeft float32) {
	s.cornerRadii = [4]float32{topLeft, topRight, bottomRight, bottomLeft}
}

// GetCornerRadii returns the corner radii
func (s *SolidRRectBlurContents) GetCornerRadii() [4]float32 {
	return s.cornerRadii
}

// Render implements the Contents interface
func (s *SolidRRectBlurContents) Render(context *ContentContext, entity *Entity, pass *RenderPass) bool {
	// TODO: Implement solid rounded rectangle blur rendering
	// This would involve:
	// 1. Creating a blur kernel based on radius and style
	// 2. Rendering the rounded rectangle with corner radii
	// 3. Applying the blur effect
	// 4. Compositing with the specified blend mode
	return false
}

// GetBounds returns the bounds of this contents including blur
func (s *SolidRRectBlurContents) GetBounds() geom.Rect {
	bounds := s.path.GetBounds()
	// Expand bounds by blur radius
	expansion := s.blurRadius * 2
	return geom.Rect{
		X:      bounds.X - expansion,
		Y:      bounds.Y - expansion,
		Width:  bounds.Width + expansion*2,
		Height: bounds.Height + expansion*2,
	}
}

// Clone creates a copy of this contents
func (s *SolidRRectBlurContents) Clone() Contents {
	return &SolidRRectBlurContents{
		BaseContents: s.BaseContents.Clone().(BaseContents),
		color:        s.color,
		path:         s.path,
		blurRadius:   s.blurRadius,
		blurStyle:    s.blurStyle,
		cornerRadii:  s.cornerRadii,
	}
}

// GetCoverage returns the coverage area for this contents
func (s *SolidRRectBlurContents) GetCoverage(transform geom.Matrix) geom.Rect {
	return s.GetBounds()
}

// SolidRRectBlurFactory creates solid rounded rectangle blur contents instances
type SolidRRectBlurFactory struct{}

// CreateSolidRRectBlur creates a new solid rounded rectangle blur contents
func (f *SolidRRectBlurFactory) CreateSolidRRectBlur(color geom.Color, path geom.Path, blurRadius float32) Contents {
	return NewSolidRRectBlurContents(color, path, blurRadius)
}

// CreateSolidRRectBlurWithStyle creates a new solid rounded rectangle blur with specified style
func (f *SolidRRectBlurFactory) CreateSolidRRectBlurWithStyle(color geom.Color, path geom.Path, blurRadius float32, style BlurStyle) Contents {
	contents := NewSolidRRectBlurContents(color, path, blurRadius)
	contents.SetBlurStyle(style)
	return contents
}
