package filters

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// LinearToSRGBFilterContents represents a filter that converts linear RGB to sRGB color space
type LinearToSRGBFilterContents struct {
	FilterContentsBase
	Gamma float32
}

// NewLinearToSRGBFilterContents creates a new linear to sRGB filter
func NewLinearToSRGBFilterContents() *LinearToSRGBFilterContents {
	return &LinearToSRGBFilterContents{
		Gamma: 2.2, // Standard sRGB gamma
	}
}

// WithGamma sets a custom gamma value
func (f *LinearToSRGBFilterContents) WithGamma(gamma float32) *LinearToSRGBFilterContents {
	f.Gamma = gamma
	return f
}

// Render renders the linear to sRGB conversion filter
func (f *LinearToSRGBFilterContents) Render(renderer render.ContentRenderer, entity render.Entity) bool {
	// TODO: Implement linear to sRGB color space conversion
	// This involves applying gamma correction
	return false
}

// GetCoverage returns the coverage area of the filter
func (f *LinearToSRGBFilterContents) GetCoverage(entity render.Entity) geom.Rect {
	// Color space conversion doesn't change geometry
	return entity.GetBounds()
}

// Clone creates a copy of the filter
func (f *LinearToSRGBFilterContents) Clone() FilterContents {
	return &LinearToSRGBFilterContents{
		Gamma: f.Gamma,
	}
}

// SRGBToLinearFilterContents represents a filter that converts sRGB to linear RGB color space
type SRGBToLinearFilterContents struct {
	FilterContentsBase
	Gamma float32
}

// NewSRGBToLinearFilterContents creates a new sRGB to linear filter
func NewSRGBToLinearFilterContents() *SRGBToLinearFilterContents {
	return &SRGBToLinearFilterContents{
		Gamma: 2.2, // Standard sRGB gamma
	}
}

// WithGamma sets a custom gamma value
func (f *SRGBToLinearFilterContents) WithGamma(gamma float32) *SRGBToLinearFilterContents {
	f.Gamma = gamma
	return f
}

// Render renders the sRGB to linear conversion filter
func (f *SRGBToLinearFilterContents) Render(renderer render.ContentRenderer, entity render.Entity) bool {
	// TODO: Implement sRGB to linear color space conversion
	// This involves applying inverse gamma correction
	return false
}

// GetCoverage returns the coverage area of the filter
func (f *SRGBToLinearFilterContents) GetCoverage(entity render.Entity) geom.Rect {
	// Color space conversion doesn't change geometry
	return entity.GetBounds()
}

// Clone creates a copy of the filter
func (f *SRGBToLinearFilterContents) Clone() FilterContents {
	return &SRGBToLinearFilterContents{
		Gamma: f.Gamma,
	}
}
