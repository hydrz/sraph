package render

import (
	"github.com/opensraph/sraph/geom"
)

// Surface represents a render surface that can be presented to the screen.
// It encapsulates a render target and provides presentation capabilities.
type Surface interface {
	// GetSize returns the size of the surface in pixels
	GetSize() geom.Size

	// IsValid returns true if the surface is valid and ready for rendering
	IsValid() bool

	// GetRenderTarget returns the render target associated with this surface
	GetRenderTarget() RenderTarget

	// Present presents the rendered content to the screen
	Present() error
}

// SurfaceImpl is the default implementation of Surface
type SurfaceImpl struct {
	renderTarget RenderTarget
	size         geom.Size
	isValid      bool
}

// NewSurface creates a new surface with the specified render target
func NewSurface(target RenderTarget) Surface {
	return &SurfaceImpl{
		renderTarget: target,
		size:         target.GetSize(),
		isValid:      target != nil && target.IsValid(),
	}
}

// GetSize returns the size of the surface in pixels
func (s *SurfaceImpl) GetSize() geom.Size {
	return s.size
}

// IsValid returns true if the surface is valid and ready for rendering
func (s *SurfaceImpl) IsValid() bool {
	return s.isValid && s.renderTarget != nil && s.renderTarget.IsValid()
}

// GetRenderTarget returns the render target associated with this surface
func (s *SurfaceImpl) GetRenderTarget() RenderTarget {
	return s.renderTarget
}

// Present presents the rendered content to the screen
func (s *SurfaceImpl) Present() error {
	// TODO: Implement surface presentation logic
	return nil
}
