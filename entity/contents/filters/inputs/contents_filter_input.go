package inputs

import (
	"github.com/opensraph/sraph/render"
)

// ContentsFilterInput represents an input from contents for filter operations
type ContentsFilterInput struct {
	BaseFilterInput
	Contents render.Contents
	Texture  *render.Texture
}

// NewContentsFilterInput creates a new contents filter input
func NewContentsFilterInput(contents render.Contents) *ContentsFilterInput {
	return &ContentsFilterInput{
		BaseFilterInput: BaseFilterInput{
			InputType: FilterInputTypeContents,
		},
		Contents: contents,
	}
}

// GetTexture returns the texture representation of the contents
func (c *ContentsFilterInput) GetTexture() *render.Texture {
	if c.Texture != nil {
		return c.Texture
	}

	// TODO: Render contents to texture if not already cached
	return nil
}

// IsValid returns true if the contents input is valid
func (c *ContentsFilterInput) IsValid() bool {
	return c.Contents != nil
}

// Clone creates a copy of the contents input
func (c *ContentsFilterInput) Clone() FilterInput {
	return &ContentsFilterInput{
		BaseFilterInput: c.BaseFilterInput,
		Contents:        c.Contents,
		Texture:         c.Texture,
	}
}

// GetContents returns the underlying contents
func (c *ContentsFilterInput) GetContents() render.Contents {
	return c.Contents
}

// SetContents sets the underlying contents
func (c *ContentsFilterInput) SetContents(contents render.Contents) {
	c.Contents = contents
	// Invalidate cached texture when contents change
	c.Texture = nil
}

// CacheTexture sets the cached texture representation
func (c *ContentsFilterInput) CacheTexture(texture *render.Texture) {
	c.Texture = texture
}

// InvalidateCache clears the cached texture
func (c *ContentsFilterInput) InvalidateCache() {
	c.Texture = nil
}
