package contents

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/gpu"
	"github.com/opensraph/sraph/render"
)

// TextureContents renders textured fills
type TextureContents struct {
	*BaseContents
	texture         gpu.Texture
	samplingMode    SamplingMode
	sourceRect      *geom.Rect
	destinationRect *geom.Rect
}

// SamplingMode defines how textures are sampled
type SamplingMode int

const (
	SamplingModeNearest SamplingMode = iota
	SamplingModeLinear
)

// NewTextureContents creates a new texture contents
func NewTextureContents() *TextureContents {
	return &TextureContents{
		BaseContents:    NewBaseContents(),
		texture:         nil,
		samplingMode:    SamplingModeLinear,
		sourceRect:      nil,
		destinationRect: nil,
	}
}

// SetTexture sets the texture to render
func (c *TextureContents) SetTexture(texture gpu.Texture) {
	c.texture = texture
}

// GetTexture returns the texture
func (c *TextureContents) GetTexture() gpu.Texture {
	return c.texture
}

// SetSamplingMode sets the texture sampling mode
func (c *TextureContents) SetSamplingMode(mode SamplingMode) {
	c.samplingMode = mode
}

// GetSamplingMode returns the texture sampling mode
func (c *TextureContents) GetSamplingMode() SamplingMode {
	return c.samplingMode
}

// SetSourceRect sets the source rectangle within the texture
func (c *TextureContents) SetSourceRect(rect *geom.Rect) {
	c.sourceRect = rect
}

// GetSourceRect returns the source rectangle
func (c *TextureContents) GetSourceRect() *geom.Rect {
	return c.sourceRect
}

// SetDestinationRect sets the destination rectangle for rendering
func (c *TextureContents) SetDestinationRect(rect *geom.Rect) {
	c.destinationRect = rect
}

// GetDestinationRect returns the destination rectangle
func (c *TextureContents) GetDestinationRect() *geom.Rect {
	return c.destinationRect
}

// GetCoverage returns the coverage area of this content
func (c *TextureContents) GetCoverage(transform geom.Matrix) *geom.Rect {
	var bounds geom.Rect

	if c.destinationRect != nil {
		bounds = *c.destinationRect
	} else if c.texture != nil {
		// Use texture size as bounds
		size := c.texture.GetSize()
		bounds = geom.Rect{
			Origin: geom.Point{X: 0, Y: 0},
			Size:   geom.Size{Width: float32(size.Width), Height: float32(size.Height)},
		}
	} else {
		return nil
	}

	transformedBounds := transform.TransformRect(bounds)
	return &transformedBounds
}

// Render renders this content using the provided context and render pass
func (c *TextureContents) Render(
	context *ContentContext,
	pass *render.RenderPass,
	transform geom.Matrix,
	entity Entity,
) bool {
	if c.texture == nil || context == nil || pass == nil {
		return false
	}

	// Get the texture shader
	shader := context.GetTextureShader()
	if shader == nil {
		return false
	}

	// Create pipeline
	pipeline := context.GetPipelineLibrary().GetTexturePipeline(c.GetBlendMode())
	if pipeline == nil {
		return false
	}

	// Set up vertex data
	vertices := c.generateVertices(transform)
	if len(vertices) == 0 {
		return false
	}

	// Create vertex buffer
	vertexBuffer := context.GetHostBuffer().Emplace(vertices)
	if vertexBuffer == nil {
		return false
	}

	// Set up uniforms
	uniforms := c.createUniforms(transform, entity)
	uniformBuffer := context.GetHostBuffer().Emplace(uniforms)
	if uniformBuffer == nil {
		return false
	}

	// Get sampler
	sampler := context.GetSamplerLibrary().GetSampler(c.convertSamplingMode())
	if sampler == nil {
		return false
	}

	// Bind resources
	pass.SetPipeline(pipeline)
	pass.SetVertexBuffer(vertexBuffer)
	pass.BindUniformBuffer(uniformBuffer, 0)
	pass.BindTexture(c.texture, 0)
	pass.BindSampler(sampler, 0)

	// Draw
	return pass.Draw(len(vertices))
}

// IsOpaque returns true if this content is fully opaque
func (c *TextureContents) IsOpaque() bool {
	// Texture contents are considered opaque if the base is opaque and
	// the texture doesn't have an alpha channel (which we can't easily check)
	return c.BaseContents.IsOpaque()
}

// Clone creates a copy of this content
func (c *TextureContents) Clone() Contents {
	clone := NewTextureContents()
	clone.SetTexture(c.texture)
	clone.SetSamplingMode(c.samplingMode)
	clone.SetBlendMode(c.GetBlendMode())
	clone.SetOpacity(c.GetOpacity())
	if c.sourceRect != nil {
		sourceRect := *c.sourceRect
		clone.SetSourceRect(&sourceRect)
	}
	if c.destinationRect != nil {
		destRect := *c.destinationRect
		clone.SetDestinationRect(&destRect)
	}
	return clone
}

// generateVertices generates vertex data for texture rendering
func (c *TextureContents) generateVertices(transform geom.Matrix) []TextureVertex {
	var bounds geom.Rect
	var texCoords geom.Rect

	// Determine destination bounds
	if c.destinationRect != nil {
		bounds = *c.destinationRect
	} else if c.texture != nil {
		size := c.texture.GetSize()
		bounds = geom.Rect{
			Origin: geom.Point{X: 0, Y: 0},
			Size:   geom.Size{Width: float32(size.Width), Height: float32(size.Height)},
		}
	} else {
		return nil
	}

	// Determine texture coordinates
	if c.sourceRect != nil && c.texture != nil {
		textureSize := c.texture.GetSize()
		texCoords = geom.Rect{
			Origin: geom.Point{
				X: c.sourceRect.Origin.X / float32(textureSize.Width),
				Y: c.sourceRect.Origin.Y / float32(textureSize.Height),
			},
			Size: geom.Size{
				Width:  c.sourceRect.Size.Width / float32(textureSize.Width),
				Height: c.sourceRect.Size.Height / float32(textureSize.Height),
			},
		}
	} else {
		texCoords = geom.Rect{
			Origin: geom.Point{X: 0, Y: 0},
			Size:   geom.Size{Width: 1, Height: 1},
		}
	}

	// Create quad vertices
	return []TextureVertex{
		{Position: bounds.Origin, TexCoord: texCoords.Origin},
		{Position: geom.Point{X: bounds.Origin.X + bounds.Size.Width, Y: bounds.Origin.Y},
			TexCoord: geom.Point{X: texCoords.Origin.X + texCoords.Size.Width, Y: texCoords.Origin.Y}},
		{Position: geom.Point{X: bounds.Origin.X + bounds.Size.Width, Y: bounds.Origin.Y + bounds.Size.Height},
			TexCoord: geom.Point{X: texCoords.Origin.X + texCoords.Size.Width, Y: texCoords.Origin.Y + texCoords.Size.Height}},
		{Position: geom.Point{X: bounds.Origin.X, Y: bounds.Origin.Y + bounds.Size.Height},
			TexCoord: geom.Point{X: texCoords.Origin.X, Y: texCoords.Origin.Y + texCoords.Size.Height}},
	}
}

// createUniforms creates uniform data for rendering
func (c *TextureContents) createUniforms(transform geom.Matrix, entity Entity) *TextureUniforms {
	return &TextureUniforms{
		Transform: transform,
		Opacity:   c.GetOpacity(),
	}
}

// convertSamplingMode converts SamplingMode to render system sampling mode
func (c *TextureContents) convertSamplingMode() render.SamplingMode {
	switch c.samplingMode {
	case SamplingModeNearest:
		return render.SamplingModeNearest
	case SamplingModeLinear:
		return render.SamplingModeLinear
	default:
		return render.SamplingModeLinear
	}
}

// TextureVertex represents a vertex for texture rendering
type TextureVertex struct {
	Position geom.Point
	TexCoord geom.Point
}

// TextureUniforms represents uniform data for texture rendering
type TextureUniforms struct {
	Transform geom.Matrix
	Opacity   float32
}
