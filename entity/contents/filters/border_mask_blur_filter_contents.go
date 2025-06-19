package filters

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// BorderMaskBlurFilterContents represents a filter that creates a blur effect with border masking
type BorderMaskBlurFilterContents struct {
	*FilterContents
	BlurRadius   float32
	BorderRadius float32
	MaskTexture  *render.Texture
	BlurStyle    BlurStyle
}

// BlurStyle defines different blur styles for border mask blur
type BlurStyle int

const (
	BlurStyleNormal BlurStyle = iota
	BlurStyleSolid
	BlurStyleOuter
	BlurStyleInner
)

// NewBorderMaskBlurFilterContents creates a new border mask blur filter
func NewBorderMaskBlurFilterContents(blurRadius, borderRadius float32) *BorderMaskBlurFilterContents {
	return &BorderMaskBlurFilterContents{
		FilterContents: NewFilterContents(),
		BlurRadius:     blurRadius,
		BorderRadius:   borderRadius,
		BlurStyle:      BlurStyleNormal,
	}
}

// WithMaskTexture sets the mask texture for the blur
func (f *BorderMaskBlurFilterContents) WithMaskTexture(texture *render.Texture) *BorderMaskBlurFilterContents {
	f.MaskTexture = texture
	return f
}

// WithBlurStyle sets the blur style
func (f *BorderMaskBlurFilterContents) WithBlurStyle(style BlurStyle) *BorderMaskBlurFilterContents {
	f.BlurStyle = style
	return f
}

// Render renders the border mask blur filter
func (f *BorderMaskBlurFilterContents) Render(renderer render.ContentRenderer, entity render.Entity) bool {
	// TODO: Implement border mask blur rendering
	// This involves creating a blur effect with proper border masking
	return false
}

// GetCoverage returns the coverage area of the filter
func (f *BorderMaskBlurFilterContents) GetCoverage(entity render.Entity) geom.Rect {
	// TODO: Calculate coverage including blur expansion and border effects
	return geom.Rect{}
}

// Clone creates a copy of the filter
func (f *BorderMaskBlurFilterContents) Clone() FilterContents {
	return &BorderMaskBlurFilterContents{
		FilterContents: f.FilterContents.Clone().(*FilterContents),
		BlurRadius:     f.BlurRadius,
		BorderRadius:   f.BorderRadius,
		MaskTexture:    f.MaskTexture,
		BlurStyle:      f.BlurStyle,
	}
}

// SetBlurRadius sets the blur radius
func (f *BorderMaskBlurFilterContents) SetBlurRadius(radius float32) {
	f.BlurRadius = radius
}

// GetBlurRadius returns the blur radius
func (f *BorderMaskBlurFilterContents) GetBlurRadius() float32 {
	return f.BlurRadius
}

// SetBorderRadius sets the border radius
func (f *BorderMaskBlurFilterContents) SetBorderRadius(radius float32) {
	f.BorderRadius = radius
}

// GetBorderRadius returns the border radius
func (f *BorderMaskBlurFilterContents) GetBorderRadius() float32 {
	return f.BorderRadius
}
