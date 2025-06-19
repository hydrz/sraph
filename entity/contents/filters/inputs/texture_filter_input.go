package inputs

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// TextureFilterInput represents a texture input for filter operations
type TextureFilterInput struct {
	BaseFilterInput
	Texture *render.Texture
}

// NewTextureFilterInput creates a new texture filter input
func NewTextureFilterInput(texture *render.Texture) *TextureFilterInput {
	input := &TextureFilterInput{
		BaseFilterInput: BaseFilterInput{
			InputType: FilterInputTypeTexture,
		},
		Texture: texture,
	}

	// Set bounds based on texture size if available
	if texture != nil {
		// TODO: Get actual texture dimensions
		input.Bounds = geom.Rect{
			Origin: geom.Point{X: 0, Y: 0},
			Size:   geom.Size{Width: 1, Height: 1}, // Placeholder
		}
	}

	return input
}

// GetTexture returns the underlying texture
func (t *TextureFilterInput) GetTexture() *render.Texture {
	return t.Texture
}

// IsValid returns true if the texture input is valid
func (t *TextureFilterInput) IsValid() bool {
	return t.Texture != nil
}

// Clone creates a copy of the texture input
func (t *TextureFilterInput) Clone() FilterInput {
	return &TextureFilterInput{
		BaseFilterInput: t.BaseFilterInput,
		Texture:         t.Texture,
	}
}

// SetTexture sets the underlying texture
func (t *TextureFilterInput) SetTexture(texture *render.Texture) {
	t.Texture = texture

	// Update bounds based on new texture
	if texture != nil {
		// TODO: Get actual texture dimensions and update bounds
	}
}

// GetTextureSize returns the size of the texture
func (t *TextureFilterInput) GetTextureSize() geom.Size {
	if t.Texture == nil {
		return geom.Size{Width: 0, Height: 0}
	}

	// TODO: Return actual texture size
	return geom.Size{Width: 1, Height: 1}
}

// GetTextureFormat returns the format of the texture
func (t *TextureFilterInput) GetTextureFormat() render.PixelFormat {
	if t.Texture == nil {
		return render.PixelFormatUnknown
	}

	// TODO: Return actual texture format
	return render.PixelFormatRGBA8
}
