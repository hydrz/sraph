package entity

import (
	"github.com/opensraph/sraph/gpu"
	"github.com/opensraph/sraph/render"
)

// EntityPassTarget manages render targets for entity passes.
// It supports backdrop flipping and secondary color texture management.
type EntityPassTarget interface {
	// Flip flips the backdrop and returns a readable texture that can be
	// bound/sampled to restore the previous pass.
	Flip(renderer *ContentContext) gpu.Texture

	// GetRenderTarget returns the current render target
	GetRenderTarget() *render.RenderTarget

	// RemoveSecondary removes the cached secondary color texture
	RemoveSecondary()

	// IsValid returns true if the target is valid
	IsValid() bool
}

// DefaultEntityPassTarget provides a default implementation of EntityPassTarget
type DefaultEntityPassTarget struct {
	target                  *render.RenderTarget // Primary render target
	secondaryColorTexture   gpu.Texture          // Cached secondary color texture
	supportsReadFromResolve bool                 // Whether the target supports reading from resolve
	supportsImplicitMSAA    bool                 // Whether the target supports implicit MSAA
}

// NewEntityPassTarget creates a new entity pass target
func NewEntityPassTarget(
	renderTarget *render.RenderTarget,
	supportsReadFromResolve bool,
	supportsImplicitMSAA bool,
) EntityPassTarget {
	return &DefaultEntityPassTarget{
		target:                  renderTarget,
		secondaryColorTexture:   nil,
		supportsReadFromResolve: supportsReadFromResolve,
		supportsImplicitMSAA:    supportsImplicitMSAA,
	}
}

// Flip flips the backdrop and returns a readable texture
func (t *DefaultEntityPassTarget) Flip(renderer *ContentContext) gpu.Texture {
	if !t.IsValid() {
		return nil
	}

	// If we already have a secondary texture, return it
	if t.secondaryColorTexture != nil {
		return t.secondaryColorTexture
	}

	// Create a new secondary texture by copying from the current target
	if t.target != nil {
		colorTexture := t.target.GetColorTexture()
		if colorTexture != nil {
			// TODO: Implement texture copying/flipping logic
			// For now, just return the color texture directly
			t.secondaryColorTexture = colorTexture
			return t.secondaryColorTexture
		}
	}

	return nil
}

// GetRenderTarget returns the current render target
func (t *DefaultEntityPassTarget) GetRenderTarget() *render.RenderTarget {
	return t.target
}

// RemoveSecondary removes the cached secondary color texture
func (t *DefaultEntityPassTarget) RemoveSecondary() {
	t.secondaryColorTexture = nil
}

// IsValid returns true if the target is valid
func (t *DefaultEntityPassTarget) IsValid() bool {
	return t.target != nil && t.target.IsValid()
}

// SupportsReadFromResolve returns true if the target supports reading from resolve
func (t *DefaultEntityPassTarget) SupportsReadFromResolve() bool {
	return t.supportsReadFromResolve
}

// SupportsImplicitMSAA returns true if the target supports implicit MSAA
func (t *DefaultEntityPassTarget) SupportsImplicitMSAA() bool {
	return t.supportsImplicitMSAA
}
