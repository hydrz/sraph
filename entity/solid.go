package entity

import (
	"github.com/opensraph/sraph/geom"
)

// SolidContents represents solid color content that can be rendered.
// This is inspired by Impeller's solid color rendering.
type SolidContents struct {
	color geom.Color
	path  geom.PathSource[geom.F32]
}

// NewSolidContents creates a new solid color contents.
func NewSolidContents(color geom.Color, path geom.PathSource[geom.F32]) *SolidContents {
	return &SolidContents{
		color: color,
		path:  path,
	}
}

// Render implements Contents.
func (s *SolidContents) Render(ctx *ContentContext, entity *Entity, pass RenderPass) bool {
	// TODO: Implement actual solid color rendering
	// This will use the GPU backend to render solid colors
	return true
}

// GetCoverage implements Contents.
func (s *SolidContents) GetCoverage(entity *Entity) *geom.Rect[geom.F32] {
	if s.path != nil {
		bounds := s.path.Bounds()
		return &bounds
	}
	return nil
}

// CanInheritOpacity implements Contents.
func (s *SolidContents) CanInheritOpacity(entity *Entity) bool {
	return true
}

// Clone implements Contents.
func (s *SolidContents) Clone() Contents {
	return &SolidContents{
		color: s.color,
		path:  s.path, // TODO: Clone path if needed
	}
}

// GetColor returns the solid color.
func (s *SolidContents) GetColor() geom.Color {
	return s.color
}

// SetColor sets the solid color.
func (s *SolidContents) SetColor(color geom.Color) {
	s.color = color
}

// GetPath returns the path to be filled with solid color.
func (s *SolidContents) GetPath() geom.PathSource[geom.F32] {
	return s.path
}

// SetPath sets the path to be filled with solid color.
func (s *SolidContents) SetPath(path geom.PathSource[geom.F32]) {
	s.path = path
}
