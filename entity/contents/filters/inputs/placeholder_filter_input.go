package inputs

import (
	"github.com/hydrz/sraph/geom"
	"github.com/hydrz/sraph/render"
)

// PlaceholderFilterInput represents a placeholder input for filter operations
// Used when an input is expected but not yet available
type PlaceholderFilterInput struct {
	BaseFilterInput
	PlaceholderColor geom.Color
	PlaceholderSize  geom.Size
}

// NewPlaceholderFilterInput creates a new placeholder filter input
func NewPlaceholderFilterInput(size geom.Size, color geom.Color) *PlaceholderFilterInput {
	return &PlaceholderFilterInput{
		BaseFilterInput: BaseFilterInput{
			InputType: FilterInputTypePlaceholder,
			Bounds: geom.Rect{
				Origin: geom.Point{X: 0, Y: 0},
				Size:   size,
			},
		},
		PlaceholderColor: color,
		PlaceholderSize:  size,
	}
}

// GetTexture returns a solid color texture with the placeholder color
func (p *PlaceholderFilterInput) GetTexture() *render.Texture {
	// TODO: Generate a solid color texture with the placeholder color and size
	return nil
}

// IsValid returns true if the placeholder input is valid
func (p *PlaceholderFilterInput) IsValid() bool {
	return p.PlaceholderSize.Width > 0 && p.PlaceholderSize.Height > 0
}

// Clone creates a copy of the placeholder input
func (p *PlaceholderFilterInput) Clone() FilterInput {
	return &PlaceholderFilterInput{
		BaseFilterInput:  p.BaseFilterInput,
		PlaceholderColor: p.PlaceholderColor,
		PlaceholderSize:  p.PlaceholderSize,
	}
}

// GetPlaceholderColor returns the placeholder color
func (p *PlaceholderFilterInput) GetPlaceholderColor() geom.Color {
	return p.PlaceholderColor
}

// SetPlaceholderColor sets the placeholder color
func (p *PlaceholderFilterInput) SetPlaceholderColor(color geom.Color) {
	p.PlaceholderColor = color
}

// GetPlaceholderSize returns the placeholder size
func (p *PlaceholderFilterInput) GetPlaceholderSize() geom.Size {
	return p.PlaceholderSize
}

// SetPlaceholderSize sets the placeholder size
func (p *PlaceholderFilterInput) SetPlaceholderSize(size geom.Size) {
	p.PlaceholderSize = size
	p.Bounds.Size = size
}
