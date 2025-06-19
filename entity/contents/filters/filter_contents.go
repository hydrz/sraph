// Package filters provides content filters for entity rendering.
// This package contains various filters that can be applied to content,
// including blur filters, color filters, and image filters.
package filters

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// FilterContents represents a filter that can be applied to content
type FilterContents interface {
	// GetCoverage returns the coverage area affected by this filter
	GetCoverage(input geom.Rect) *geom.Rect

	// Render applies the filter to the input and renders to the output
	Render(context FilterContext, input FilterInput, output FilterOutput) bool

	// IsValid returns true if the filter is valid
	IsValid() bool

	// GetRequiredMipCount returns the number of mip levels required
	GetRequiredMipCount() int
}

// FilterContext provides context for filter operations
type FilterContext interface {
	// GetRenderContext returns the render context
	GetRenderContext() *render.Context

	// CreateTexture creates a new texture for filter operations
	CreateTexture(size geom.Size, format TextureFormat) Texture

	// CreateRenderTarget creates a render target for filter output
	CreateRenderTarget(texture Texture) *render.RenderTarget
}

// FilterInput represents input to a filter operation
type FilterInput interface {
	// GetTexture returns the input texture
	GetTexture() Texture

	// GetCoverage returns the coverage area of the input
	GetCoverage() geom.Rect

	// GetTransform returns the transform applied to the input
	GetTransform() geom.Matrix
}

// FilterOutput represents output from a filter operation
type FilterOutput interface {
	// GetTexture returns the output texture
	GetTexture() Texture

	// GetRenderTarget returns the render target for output
	GetRenderTarget() *render.RenderTarget
}

// Texture represents a texture used in filter operations
type Texture interface {
	// GetSize returns the texture size
	GetSize() geom.Size

	// GetFormat returns the texture format
	GetFormat() TextureFormat

	// IsValid returns true if the texture is valid
	IsValid() bool
}

// TextureFormat represents texture pixel formats
type TextureFormat int

const (
	TextureFormatRGBA8 TextureFormat = iota
	TextureFormatRGBA16F
	TextureFormatRGBA32F
	TextureFormatR8
	TextureFormatR16F
	TextureFormatR32F
)

// BaseFilterContents provides base implementation for filters
type BaseFilterContents struct {
	requiredMipCount int
	isValid          bool
}

// NewBaseFilterContents creates a new base filter contents
func NewBaseFilterContents() *BaseFilterContents {
	return &BaseFilterContents{
		requiredMipCount: 1,
		isValid:          true,
	}
}

// GetRequiredMipCount returns the number of mip levels required
func (f *BaseFilterContents) GetRequiredMipCount() int {
	return f.requiredMipCount
}

// SetRequiredMipCount sets the number of mip levels required
func (f *BaseFilterContents) SetRequiredMipCount(count int) {
	f.requiredMipCount = count
}

// IsValid returns true if the filter is valid
func (f *BaseFilterContents) IsValid() bool {
	return f.isValid
}

// SetValid sets the validity of the filter
func (f *BaseFilterContents) SetValid(valid bool) {
	f.isValid = valid
}

// BlurFilterContents applies blur effects to content
type BlurFilterContents struct {
	*BaseFilterContents
	blurRadius float32
	blurStyle  BlurStyle
}

// BlurStyle defines the style of blur
type BlurStyle int

const (
	BlurStyleNormal BlurStyle = iota
	BlurStyleSolid
	BlurStyleOuter
	BlurStyleInner
)

// NewBlurFilterContents creates a new blur filter
func NewBlurFilterContents(radius float32) *BlurFilterContents {
	return &BlurFilterContents{
		BaseFilterContents: NewBaseFilterContents(),
		blurRadius:         radius,
		blurStyle:          BlurStyleNormal,
	}
}

// SetBlurRadius sets the blur radius
func (f *BlurFilterContents) SetBlurRadius(radius float32) {
	f.blurRadius = radius
}

// GetBlurRadius returns the blur radius
func (f *BlurFilterContents) GetBlurRadius() float32 {
	return f.blurRadius
}

// SetBlurStyle sets the blur style
func (f *BlurFilterContents) SetBlurStyle(style BlurStyle) {
	f.blurStyle = style
}

// GetBlurStyle returns the blur style
func (f *BlurFilterContents) GetBlurStyle() BlurStyle {
	return f.blurStyle
}

// GetCoverage returns the coverage area affected by this filter
func (f *BlurFilterContents) GetCoverage(input geom.Rect) *geom.Rect {
	// Blur expands the coverage by the blur radius
	expansion := f.blurRadius * 2
	expanded := geom.Rect{
		Origin: geom.Point{
			X: input.Origin.X - expansion,
			Y: input.Origin.Y - expansion,
		},
		Size: geom.Size{
			Width:  input.Size.Width + 2*expansion,
			Height: input.Size.Height + 2*expansion,
		},
	}
	return &expanded
}

// Render applies the blur filter
func (f *BlurFilterContents) Render(context FilterContext, input FilterInput, output FilterOutput) bool {
	// TODO: Implement blur filter rendering
	// This would typically involve:
	// 1. Creating intermediate textures for separable blur
	// 2. Performing horizontal blur pass
	// 3. Performing vertical blur pass
	// 4. Compositing the result
	return false
}

// ColorFilterContents applies color transformations to content
type ColorFilterContents struct {
	*BaseFilterContents
	colorMatrix geom.Matrix4
}

// NewColorFilterContents creates a new color filter
func NewColorFilterContents(matrix geom.Matrix4) *ColorFilterContents {
	return &ColorFilterContents{
		BaseFilterContents: NewBaseFilterContents(),
		colorMatrix:        matrix,
	}
}

// SetColorMatrix sets the color transformation matrix
func (f *ColorFilterContents) SetColorMatrix(matrix geom.Matrix4) {
	f.colorMatrix = matrix
}

// GetColorMatrix returns the color transformation matrix
func (f *ColorFilterContents) GetColorMatrix() geom.Matrix4 {
	return f.colorMatrix
}

// GetCoverage returns the coverage area affected by this filter
func (f *ColorFilterContents) GetCoverage(input geom.Rect) *geom.Rect {
	// Color filters don't change coverage
	return &input
}

// Render applies the color filter
func (f *ColorFilterContents) Render(context FilterContext, input FilterInput, output FilterOutput) bool {
	// TODO: Implement color filter rendering
	// This would involve applying the color matrix transformation
	return false
}
