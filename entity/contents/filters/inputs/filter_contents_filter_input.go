package inputs

import (
	"github.com/opensraph/sraph/render"
)

// FilterContentsFilterInput represents an input from another filter for filter operations
type FilterContentsFilterInput struct {
	BaseFilterInput
	FilterContents interface{} // TODO: Define proper filter contents interface
	Texture        *render.Texture
}

// NewFilterContentsFilterInput creates a new filter contents filter input
func NewFilterContentsFilterInput(filterContents interface{}) *FilterContentsFilterInput {
	return &FilterContentsFilterInput{
		BaseFilterInput: BaseFilterInput{
			InputType: FilterInputTypeFilter,
		},
		FilterContents: filterContents,
	}
}

// GetTexture returns the texture representation of the filter contents
func (f *FilterContentsFilterInput) GetTexture() *render.Texture {
	if f.Texture != nil {
		return f.Texture
	}

	// TODO: Render filter contents to texture if not already cached
	return nil
}

// IsValid returns true if the filter contents input is valid
func (f *FilterContentsFilterInput) IsValid() bool {
	return f.FilterContents != nil
}

// Clone creates a copy of the filter contents input
func (f *FilterContentsFilterInput) Clone() FilterInput {
	return &FilterContentsFilterInput{
		BaseFilterInput: f.BaseFilterInput,
		FilterContents:  f.FilterContents,
		Texture:         f.Texture,
	}
}

// GetFilterContents returns the underlying filter contents
func (f *FilterContentsFilterInput) GetFilterContents() interface{} {
	return f.FilterContents
}

// SetFilterContents sets the underlying filter contents
func (f *FilterContentsFilterInput) SetFilterContents(filterContents interface{}) {
	f.FilterContents = filterContents
	// Invalidate cached texture when filter contents change
	f.Texture = nil
}

// CacheTexture sets the cached texture representation
func (f *FilterContentsFilterInput) CacheTexture(texture *render.Texture) {
	f.Texture = texture
}

// InvalidateCache clears the cached texture
func (f *FilterContentsFilterInput) InvalidateCache() {
	f.Texture = nil
}
