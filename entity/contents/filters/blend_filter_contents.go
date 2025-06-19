package filters

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// BlendFilterContents provides blend filter rendering capabilities
type BlendFilterContents struct {
	BaseFilterContents
	blendMode render.BlendMode
	sourceA   FilterInput
	sourceB   FilterInput
	opacity   float32
}

// NewBlendFilterContents creates a new blend filter contents
func NewBlendFilterContents(blendMode render.BlendMode, sourceA, sourceB FilterInput) *BlendFilterContents {
	return &BlendFilterContents{
		BaseFilterContents: NewBaseFilterContents(),
		blendMode:          blendMode,
		sourceA:            sourceA,
		sourceB:            sourceB,
		opacity:            1.0,
	}
}

// SetBlendMode sets the blend mode
func (b *BlendFilterContents) SetBlendMode(mode render.BlendMode) {
	b.blendMode = mode
}

// GetBlendMode returns the current blend mode
func (b *BlendFilterContents) GetBlendMode() render.BlendMode {
	return b.blendMode
}

// SetSourceA sets the first input source
func (b *BlendFilterContents) SetSourceA(source FilterInput) {
	b.sourceA = source
}

// GetSourceA returns the first input source
func (b *BlendFilterContents) GetSourceA() FilterInput {
	return b.sourceA
}

// SetSourceB sets the second input source
func (b *BlendFilterContents) SetSourceB(source FilterInput) {
	b.sourceB = source
}

// GetSourceB returns the second input source
func (b *BlendFilterContents) GetSourceB() FilterInput {
	return b.sourceB
}

// SetOpacity sets the opacity
func (b *BlendFilterContents) SetOpacity(opacity float32) {
	b.opacity = opacity
}

// GetOpacity returns the current opacity
func (b *BlendFilterContents) GetOpacity() float32 {
	return b.opacity
}

// Render implements the FilterContents interface
func (b *BlendFilterContents) Render(context *FilterContext, entity *Entity, pass *RenderPass) bool {
	// TODO: Implement blend filter rendering
	return false
}

// GetBounds returns the bounds of this filter
func (b *BlendFilterContents) GetBounds() geom.Rect {
	boundsA := b.sourceA.GetBounds()
	boundsB := b.sourceB.GetBounds()
	return geom.UnionRect(boundsA, boundsB)
}

// Clone creates a copy of this filter
func (b *BlendFilterContents) Clone() FilterContents {
	return &BlendFilterContents{
		BaseFilterContents: b.BaseFilterContents.Clone().(BaseFilterContents),
		blendMode:          b.blendMode,
		sourceA:            b.sourceA,
		sourceB:            b.sourceB,
		opacity:            b.opacity,
	}
}

// GetCoverage returns the coverage area for this filter
func (b *BlendFilterContents) GetCoverage(transform geom.Matrix) geom.Rect {
	return b.GetBounds()
}
