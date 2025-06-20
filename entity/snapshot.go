package entity

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// Snapshot represents a captured rendering state that can be rendered as an entity.
// Based on Impeller's Snapshot class design.
type Snapshot interface {
	// GetTexture returns the texture containing the snapshot data
	GetTexture() render.Texture

	// GetTransform returns the transform applied to the snapshot
	GetTransform() geom.Matrix[geom.F32]

	// GetOpacity returns the opacity of the snapshot (0.0 to 1.0)
	GetOpacity() geom.F32

	// GetSize returns the size of the snapshot in pixels
	GetSize() geom.Size[geom.F32]

	// GetCoverage returns the coverage area in local coordinates
	GetCoverage() geom.Rect[geom.F32]

	// IsValid returns true if this snapshot is valid for rendering
	IsValid() bool

	// GetSamplerDescriptor returns the sampler configuration for the snapshot
	GetSamplerDescriptor() SamplerDescriptor
}

// SnapshotImpl is the default implementation of Snapshot
type SnapshotImpl struct {
	texture           render.Texture
	transform         geom.Matrix[geom.F32]
	opacity           geom.F32
	coverage          geom.Rect[geom.F32]
	samplerDescriptor SamplerDescriptor
}

// NewSnapshot creates a new snapshot with the given parameters
func NewSnapshot(
	texture render.Texture,
	transform geom.Matrix[geom.F32],
	opacity geom.F32,
	coverage geom.Rect[geom.F32],
) Snapshot {
	return &SnapshotImpl{
		texture:   texture,
		transform: transform,
		opacity:   opacity,
		coverage:  coverage,
		samplerDescriptor: SamplerDescriptor{
			MinFilter:    render.FilterModeLinear,
			MagFilter:    render.FilterModeLinear,
			AddressModeU: render.AddressModeClampToEdge,
			AddressModeV: render.AddressModeClampToEdge,
		},
	}
}

// GetTexture returns the texture containing the snapshot data
func (s *SnapshotImpl) GetTexture() render.Texture {
	return s.texture
}

// GetTransform returns the transform applied to the snapshot
func (s *SnapshotImpl) GetTransform() geom.Matrix[geom.F32] {
	return s.transform
}

// GetOpacity returns the opacity of the snapshot
func (s *SnapshotImpl) GetOpacity() geom.F32 {
	return s.opacity
}

// GetSize returns the size of the snapshot in pixels
func (s *SnapshotImpl) GetSize() geom.Size[geom.F32] {
	if s.texture == nil {
		return geom.Size[geom.F32]{}
	}

	desc := s.texture.GetDescriptor()
	return geom.Size[geom.F32]{
		Width:  geom.F32(desc.Size.Width),
		Height: geom.F32(desc.Size.Height),
	}
}

// GetCoverage returns the coverage area in local coordinates
func (s *SnapshotImpl) GetCoverage() geom.Rect[geom.F32] {
	return s.coverage
}

// IsValid returns true if this snapshot is valid for rendering
func (s *SnapshotImpl) IsValid() bool {
	return s.texture != nil && s.texture.IsValid()
}

// GetSamplerDescriptor returns the sampler configuration
func (s *SnapshotImpl) GetSamplerDescriptor() SamplerDescriptor {
	return s.samplerDescriptor
}

// SetSamplerDescriptor sets the sampler configuration
func (s *SnapshotImpl) SetSamplerDescriptor(desc SamplerDescriptor) {
	s.samplerDescriptor = desc
}

// SamplerDescriptor describes sampler configuration for snapshots
type SamplerDescriptor struct {
	MinFilter    render.FilterMode
	MagFilter    render.FilterMode
	AddressModeU render.AddressMode
	AddressModeV render.AddressMode
}

// SnapshotBuilder helps build snapshots with a fluent interface
type SnapshotBuilder struct {
	texture           render.Texture
	transform         geom.Matrix[geom.F32]
	opacity           geom.F32
	coverage          geom.Rect[geom.F32]
	samplerDescriptor SamplerDescriptor
}

// NewSnapshotBuilder creates a new snapshot builder
func NewSnapshotBuilder() *SnapshotBuilder {
	return &SnapshotBuilder{
		transform: geom.Matrix[geom.F32]{}, // Identity matrix
		opacity:   1.0,
		samplerDescriptor: SamplerDescriptor{
			MinFilter:    render.FilterModeLinear,
			MagFilter:    render.FilterModeLinear,
			AddressModeU: render.AddressModeClampToEdge,
			AddressModeV: render.AddressModeClampToEdge,
		},
	}
}

// SetTexture sets the texture for the snapshot
func (sb *SnapshotBuilder) SetTexture(texture render.Texture) *SnapshotBuilder {
	sb.texture = texture

	// Auto-calculate coverage from texture size if not set
	if sb.coverage.IsEmpty() && texture != nil {
		desc := texture.GetDescriptor()
		sb.coverage = geom.Rect[geom.F32]{
			X:      0,
			Y:      0,
			Width:  geom.F32(desc.Size.Width),
			Height: geom.F32(desc.Size.Height),
		}
	}

	return sb
}

// SetTransform sets the transform for the snapshot
func (sb *SnapshotBuilder) SetTransform(transform geom.Matrix[geom.F32]) *SnapshotBuilder {
	sb.transform = transform
	return sb
}

// SetOpacity sets the opacity for the snapshot
func (sb *SnapshotBuilder) SetOpacity(opacity geom.F32) *SnapshotBuilder {
	sb.opacity = opacity
	return sb
}

// SetCoverage sets the coverage area for the snapshot
func (sb *SnapshotBuilder) SetCoverage(coverage geom.Rect[geom.F32]) *SnapshotBuilder {
	sb.coverage = coverage
	return sb
}

// SetSamplerDescriptor sets the sampler configuration
func (sb *SnapshotBuilder) SetSamplerDescriptor(desc SamplerDescriptor) *SnapshotBuilder {
	sb.samplerDescriptor = desc
	return sb
}

// SetLinearSampling configures linear sampling
func (sb *SnapshotBuilder) SetLinearSampling() *SnapshotBuilder {
	sb.samplerDescriptor.MinFilter = render.FilterModeLinear
	sb.samplerDescriptor.MagFilter = render.FilterModeLinear
	return sb
}

// SetNearestSampling configures nearest neighbor sampling
func (sb *SnapshotBuilder) SetNearestSampling() *SnapshotBuilder {
	sb.samplerDescriptor.MinFilter = render.FilterModeNearest
	sb.samplerDescriptor.MagFilter = render.FilterModeNearest
	return sb
}

// Build creates the snapshot
func (sb *SnapshotBuilder) Build() Snapshot {
	return &SnapshotImpl{
		texture:           sb.texture,
		transform:         sb.transform,
		opacity:           sb.opacity,
		coverage:          sb.coverage,
		samplerDescriptor: sb.samplerDescriptor,
	}
}

// CreateSnapshotFromTexture is a convenience function to create a snapshot from a texture
func CreateSnapshotFromTexture(texture render.Texture) Snapshot {
	return NewSnapshotBuilder().
		SetTexture(texture).
		Build()
}

// CreateSnapshotWithTransform creates a snapshot with a transform
func CreateSnapshotWithTransform(texture render.Texture, transform geom.Matrix[geom.F32]) Snapshot {
	return NewSnapshotBuilder().
		SetTexture(texture).
		SetTransform(transform).
		Build()
}

// CreateSnapshotWithOpacity creates a snapshot with opacity
func CreateSnapshotWithOpacity(texture render.Texture, opacity geom.F32) Snapshot {
	return NewSnapshotBuilder().
		SetTexture(texture).
		SetOpacity(opacity).
		Build()
}

// SnapshotRenderer provides utilities for rendering snapshots
type SnapshotRenderer struct {
	context render.Context
}

// NewSnapshotRenderer creates a new snapshot renderer
func NewSnapshotRenderer(context render.Context) *SnapshotRenderer {
	return &SnapshotRenderer{
		context: context,
	}
}

// RenderSnapshot renders a snapshot to a render pass
func (sr *SnapshotRenderer) RenderSnapshot(
	snapshot Snapshot,
	renderPass render.RenderPass,
	transform geom.Matrix[geom.F32],
) error {
	if !snapshot.IsValid() {
		return nil // Skip invalid snapshots
	}

	// TODO: Implement snapshot rendering
	// This would involve:
	// 1. Creating a quad geometry
	// 2. Setting up texture sampling
	// 3. Applying transforms and opacity
	// 4. Issuing draw commands

	return nil
}

// CreateSnapshotTexture creates a texture suitable for use in snapshots
func (sr *SnapshotRenderer) CreateSnapshotTexture(size geom.Size[geom.F32]) (render.Texture, error) {
	desc := render.TextureDescriptor{
		Size: render.Size{
			Width:  uint32(size.Width),
			Height: uint32(size.Height),
			Depth:  1,
		},
		Format:      render.TextureFormatRGBA8Unorm,
		Usage:       render.TextureUsageRenderTarget | render.TextureUsageShaderRead,
		SampleCount: 1,
	}

	return sr.context.GetResourceAllocator().CreateTexture(desc)
}
