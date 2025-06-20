// Package contents provides content types for entity rendering.
// This package defines the various types of content that can be rendered
// as part of entities, including solid colors, gradients, textures, and text.
package contents

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/gpu"
	"github.com/opensraph/sraph/render"
)

// Contents represents renderable content that can be attached to an entity.
// This is the base interface for all content types in the rendering system.
type Contents interface {
	// GetCoverage returns the coverage area of this content in local coordinates
	GetCoverage(transform geom.Matrix) *geom.Rect

	// Render renders this content using the provided context and render pass
	Render(context *ContentContext, pass *render.RenderPass, transform geom.Matrix, entity Entity) bool

	// GetBlendMode returns the blend mode for this content
	GetBlendMode() BlendMode

	// SetBlendMode sets the blend mode for this content
	SetBlendMode(mode BlendMode)

	// IsOpaque returns true if this content is fully opaque
	IsOpaque() bool

	// Clone creates a copy of this content
	Clone() Contents
}

// Entity represents an entity in the rendering system (forward declaration)
type Entity interface {
	GetTransform() geom.Matrix
	GetContents() Contents
	GetBlendMode() BlendMode
}

// ContentContext provides context for content rendering
// Based on Impeller's ContentContext which manages pipeline state and resources
type ContentContext struct {
	renderContext      render.Context
	typographerContext interface{} // TODO: Replace with proper typographer type

	// Pipeline and resource management
	pipelineCache map[ContentContextOptions]*render.Pipeline
	samplerCache  map[SamplerDescriptor]*render.Sampler

	// Host buffer for uploading uniform data
	hostBuffer *render.HostBuffer

	// Texture and atlas management
	lazyGlyphAtlas interface{} // TODO: Replace with glyph atlas type

	// Capabilities and features
	capabilities             *render.Capabilities
	isAdvancedBlendSupported bool
}

// ContentContextOptions defines pipeline state configuration
// Based on Impeller's ContentContextOptions for pipeline variants
type ContentContextOptions struct {
	// Stencil configuration
	StencilMode StencilMode

	// Primitive type being rendered
	PrimitiveType PrimitiveType

	// Sampling options for textures
	SampleCount uint32

	// Color attachment format
	ColorAttachmentFormat gpu.TextureFormat

	// Depth/stencil format
	DepthStencilFormat gpu.TextureFormat

	// Blend mode
	BlendMode BlendMode

	// Wire frame mode for debugging
	Wireframe bool
}

// StencilMode defines how stencil testing should be performed
type StencilMode uint8

const (
	// Turn off stencil test
	StencilModeIgnore StencilMode = iota

	// Stencil-then-cover operations
	StencilModeNonZeroFill
	StencilModeEvenOddFill
	StencilModeCoverCompare
	StencilModeCoverCompareInverted

	// Overdraw prevention for strokes
	StencilModeOverdrawPrevent
	StencilModeOverdrawRestore
)

// PrimitiveType defines the type of primitive being rendered
type PrimitiveType uint8

const (
	PrimitiveTypeTriangleList PrimitiveType = iota
	PrimitiveTypeTriangleStrip
	PrimitiveTypeLineList
	PrimitiveTypeLineStrip
	PrimitiveTypePointList
)

// SamplerDescriptor describes sampler state
type SamplerDescriptor struct {
	MinFilter    gpu.FilterMode
	MagFilter    gpu.FilterMode
	MipFilter    gpu.MipFilterMode
	AddressModeU gpu.AddressMode
	AddressModeV gpu.AddressMode
	AddressModeW gpu.AddressMode
}

// BlendMode represents different blending modes for content composition
type BlendMode int

const (
	BlendModeClear BlendMode = iota
	BlendModeSrc
	BlendModeDst
	BlendModeSrcOver
	BlendModeDstOver
	BlendModeSrcIn
	BlendModeDstIn
	BlendModeSrcOut
	BlendModeDstOut
	BlendModeSrcATop
	BlendModeDstATop
	BlendModeXor
	BlendModePlus
	BlendModeModulate
	BlendModeScreen
	BlendModeOverlay
	BlendModeDarken
	BlendModeLighten
	BlendModeColorDodge
	BlendModeColorBurn
	BlendModeHardLight
	BlendModeSoftLight
	BlendModeDifference
	BlendModeExclusion
	BlendModeMultiply
	BlendModeHue
	BlendModeSaturation
	BlendModeColor
	BlendModeLuminosity
)

// BaseContents provides a base implementation for Contents interface
type BaseContents struct {
	blendMode BlendMode
	opacity   float32
}

// NewBaseContents creates a new base contents instance
func NewBaseContents() *BaseContents {
	return &BaseContents{
		blendMode: BlendModeSrcOver,
		opacity:   1.0,
	}
}

// GetBlendMode returns the blend mode for this content
func (c *BaseContents) GetBlendMode() BlendMode {
	return c.blendMode
}

// SetBlendMode sets the blend mode for this content
func (c *BaseContents) SetBlendMode(mode BlendMode) {
	c.blendMode = mode
}

// GetOpacity returns the opacity value
func (c *BaseContents) GetOpacity() float32 {
	return c.opacity
}

// SetOpacity sets the opacity value
func (c *BaseContents) SetOpacity(opacity float32) {
	if opacity < 0 {
		opacity = 0
	} else if opacity > 1 {
		opacity = 1
	}
	c.opacity = opacity
}

// IsOpaque returns true if this content is fully opaque
func (c *BaseContents) IsOpaque() bool {
	return c.opacity >= 1.0 && c.blendMode == BlendModeSrcOver
}

// ContentRenderer provides methods for rendering different types of content
type ContentRenderer struct {
	context *ContentContext
}

// NewContentRenderer creates a new content renderer
func NewContentRenderer(context *ContentContext) *ContentRenderer {
	return &ContentRenderer{context: context}
}

// RenderContents renders the given contents
func (r *ContentRenderer) RenderContents(
	contents Contents,
	pass *render.RenderPass,
	transform geom.Matrix,
	entity Entity,
) bool {
	if contents == nil {
		return false
	}

	return contents.Render(r.context, pass, transform, entity)
}

// GetCoverageForContents computes coverage for the given contents
func (r *ContentRenderer) GetCoverageForContents(
	contents Contents,
	transform geom.Matrix,
) *geom.Rect {
	if contents == nil {
		return nil
	}

	return contents.GetCoverage(transform)
}
