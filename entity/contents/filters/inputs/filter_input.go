package inputs

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// FilterInput represents an input source for filter operations
type FilterInput interface {
	// GetTexture returns the texture representation of this input
	GetTexture() *render.Texture

	// GetBounds returns the bounds of this input
	GetBounds() geom.Rect

	// IsValid returns true if this input is valid and can be used
	IsValid() bool

	// Clone creates a copy of this input
	Clone() FilterInput
}

// FilterInputType defines the type of filter input
type FilterInputType int

const (
	FilterInputTypeContents FilterInputType = iota
	FilterInputTypeTexture
	FilterInputTypeFilter
	FilterInputTypePlaceholder
)

// BaseFilterInput provides common functionality for filter inputs
type BaseFilterInput struct {
	InputType FilterInputType
	Bounds    geom.Rect
}

// GetBounds returns the bounds of the input
func (b *BaseFilterInput) GetBounds() geom.Rect {
	return b.Bounds
}

// SetBounds sets the bounds of the input
func (b *BaseFilterInput) SetBounds(bounds geom.Rect) {
	b.Bounds = bounds
}

// GetInputType returns the type of this input
func (b *BaseFilterInput) GetInputType() FilterInputType {
	return b.InputType
}
