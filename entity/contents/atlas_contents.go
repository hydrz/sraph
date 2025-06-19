package contents

import (
	"github.com/opensraph/sraph/display"
)

// AtlasContents represents atlas-based rendering contents for batched texture rendering
type AtlasContents interface {
	Contents

	// AddTextureQuad adds a texture quad to the atlas
	AddTextureQuad(texture Texture, sourceRect, destRect display.Rect, color display.Color)

	// SetSamplerDescriptor sets the sampler descriptor for texture sampling
	SetSamplerDescriptor(descriptor SamplerDescriptor)

	// GetQuadCount returns the number of quads in the atlas
	GetQuadCount() int

	// Clear removes all quads from the atlas
	Clear()
}

// TextureQuad represents a single texture quad in an atlas
type TextureQuad struct {
	Texture    Texture
	SourceRect display.Rect
	DestRect   display.Rect
	Color      display.Color
}

// AtlasContentsImpl implements AtlasContents
type AtlasContentsImpl struct {
	*ContentsImpl
	quads            []TextureQuad
	samplerDesc      SamplerDescriptor
	blendMode        display.BlendMode
	hasColorFilter   bool
	colorFilterColor display.Color
}

// NewAtlasContents creates new atlas contents
func NewAtlasContents() AtlasContents {
	return &AtlasContentsImpl{
		ContentsImpl: &ContentsImpl{
			isValid: true,
		},
		quads:       make([]TextureQuad, 0),
		samplerDesc: SamplerDescriptor{},
		blendMode:   display.BlendModeNormal,
	}
}

// AddTextureQuad adds a texture quad to the atlas
func (a *AtlasContentsImpl) AddTextureQuad(texture Texture, sourceRect, destRect display.Rect, color display.Color) {
	quad := TextureQuad{
		Texture:    texture,
		SourceRect: sourceRect,
		DestRect:   destRect,
		Color:      color,
	}
	a.quads = append(a.quads, quad)
}

// SetSamplerDescriptor sets the sampler descriptor for texture sampling
func (a *AtlasContentsImpl) SetSamplerDescriptor(descriptor SamplerDescriptor) {
	a.samplerDesc = descriptor
}

// GetQuadCount returns the number of quads in the atlas
func (a *AtlasContentsImpl) GetQuadCount() int {
	return len(a.quads)
}

// Clear removes all quads from the atlas
func (a *AtlasContentsImpl) Clear() {
	a.quads = a.quads[:0]
}

// Render renders the atlas contents
func (a *AtlasContentsImpl) Render(context ContentContext, entity Entity, pass RenderPass) bool {
	if len(a.quads) == 0 {
		return true
	}

	// TODO: Implement atlas rendering with batched quads
	return false
}

// GetCoverage returns the coverage bounds of all quads
func (a *AtlasContentsImpl) GetCoverage(entity Entity) display.Rect {
	if len(a.quads) == 0 {
		return display.Rect{}
	}

	bounds := a.quads[0].DestRect
	for i := 1; i < len(a.quads); i++ {
		bounds = bounds.Union(a.quads[i].DestRect)
	}

	return bounds
}

// Clone creates a copy of the atlas contents
func (a *AtlasContentsImpl) Clone() Contents {
	clone := &AtlasContentsImpl{
		ContentsImpl:     a.ContentsImpl.Clone().(*ContentsImpl),
		samplerDesc:      a.samplerDesc,
		blendMode:        a.blendMode,
		hasColorFilter:   a.hasColorFilter,
		colorFilterColor: a.colorFilterColor,
	}

	// Deep copy quads
	clone.quads = make([]TextureQuad, len(a.quads))
	copy(clone.quads, a.quads)

	return clone
}

// SetBlendMode sets the blend mode for the atlas
func (a *AtlasContentsImpl) SetBlendMode(mode display.BlendMode) {
	a.blendMode = mode
}

// GetBlendMode returns the blend mode
func (a *AtlasContentsImpl) GetBlendMode() display.BlendMode {
	return a.blendMode
}

// SetColorFilter sets a color filter for all quads
func (a *AtlasContentsImpl) SetColorFilter(color display.Color) {
	a.hasColorFilter = true
	a.colorFilterColor = color
}

// ClearColorFilter removes the color filter
func (a *AtlasContentsImpl) ClearColorFilter() {
	a.hasColorFilter = false
	a.colorFilterColor = display.Color{}
}
