package contents

import (
	"github.com/opensraph/sraph/font"
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// TextContents renders text using glyph atlases
type TextContents struct {
	*BaseContents
	textFrame *font.TextFrame
	color     geom.Color
	position  geom.Point
	scale     float32
}

// NewTextContents creates a new text contents
func NewTextContents() *TextContents {
	return &TextContents{
		BaseContents: NewBaseContents(),
		textFrame:    nil,
		color:        geom.ColorBlack,
		position:     geom.Point{X: 0, Y: 0},
		scale:        1.0,
	}
}

// SetTextFrame sets the text frame to render
func (c *TextContents) SetTextFrame(frame *font.TextFrame) {
	c.textFrame = frame
}

// GetTextFrame returns the text frame
func (c *TextContents) GetTextFrame() *font.TextFrame {
	return c.textFrame
}

// SetColor sets the text color
func (c *TextContents) SetColor(color geom.Color) {
	c.color = color
}

// GetColor returns the text color
func (c *TextContents) GetColor() geom.Color {
	return c.color
}

// SetPosition sets the text position
func (c *TextContents) SetPosition(position geom.Point) {
	c.position = position
}

// GetPosition returns the text position
func (c *TextContents) GetPosition() geom.Point {
	return c.position
}

// SetScale sets the text scale
func (c *TextContents) SetScale(scale float32) {
	c.scale = scale
}

// GetScale returns the text scale
func (c *TextContents) GetScale() float32 {
	return c.scale
}

// GetCoverage returns the coverage area of this content
func (c *TextContents) GetCoverage(transform geom.Matrix) *geom.Rect {
	if c.textFrame == nil {
		return nil
	}

	bounds := c.textFrame.GetBounds()
	if bounds == nil {
		return nil
	}

	// Apply scale and position
	scaledBounds := geom.Rect{
		Origin: geom.Point{
			X: c.position.X + bounds.Origin.X*c.scale,
			Y: c.position.Y + bounds.Origin.Y*c.scale,
		},
		Size: geom.Size{
			Width:  bounds.Size.Width * c.scale,
			Height: bounds.Size.Height * c.scale,
		},
	}

	transformedBounds := transform.TransformRect(scaledBounds)
	return &transformedBounds
}

// Render renders this content using the provided context and render pass
func (c *TextContents) Render(
	context *ContentContext,
	pass *render.RenderPass,
	transform geom.Matrix,
	entity Entity,
) bool {
	if c.textFrame == nil || context == nil || pass == nil {
		return false
	}

	// Get the text shader
	shader := context.GetTextShader()
	if shader == nil {
		return false
	}

	// Create pipeline
	pipeline := context.GetPipelineLibrary().GetTextPipeline(c.GetBlendMode())
	if pipeline == nil {
		return false
	}

	// Generate glyph vertices
	vertices := c.generateGlyphVertices(transform)
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

	// Get glyph atlas texture
	atlasTexture := context.GetGlyphAtlas().GetTexture()
	if atlasTexture == nil {
		return false
	}

	// Get sampler for glyph atlas
	sampler := context.GetSamplerLibrary().GetLinearSampler()
	if sampler == nil {
		return false
	}

	// Bind resources
	pass.SetPipeline(pipeline)
	pass.SetVertexBuffer(vertexBuffer)
	pass.BindUniformBuffer(uniformBuffer, 0)
	pass.BindTexture(atlasTexture, 0)
	pass.BindSampler(sampler, 0)

	// Draw
	return pass.Draw(len(vertices))
}

// IsOpaque returns true if this content is fully opaque
func (c *TextContents) IsOpaque() bool {
	// Text is typically not opaque due to anti-aliasing
	return false
}

// Clone creates a copy of this content
func (c *TextContents) Clone() Contents {
	clone := NewTextContents()
	clone.SetTextFrame(c.textFrame) // Text frames are typically immutable
	clone.SetColor(c.color)
	clone.SetPosition(c.position)
	clone.SetScale(c.scale)
	clone.SetBlendMode(c.GetBlendMode())
	clone.SetOpacity(c.GetOpacity())
	return clone
}

// generateGlyphVertices generates vertex data for glyph rendering
func (c *TextContents) generateGlyphVertices(transform geom.Matrix) []GlyphVertex {
	if c.textFrame == nil {
		return nil
	}

	var vertices []GlyphVertex

	// Get glyph runs from text frame
	runs := c.textFrame.GetRuns()
	currentX := c.position.X
	currentY := c.position.Y

	for _, run := range runs {
		glyphs := run.GetGlyphs()
		font := run.GetFont()

		for _, glyph := range glyphs {
			// Get glyph metrics
			metrics := font.GetGlyphMetrics(glyph.GetID())
			if metrics == nil {
				continue
			}

			// Calculate glyph position and size
			glyphX := currentX + glyph.GetOffset().X*c.scale
			glyphY := currentY + glyph.GetOffset().Y*c.scale
			glyphWidth := metrics.GetBounds().Size.Width * c.scale
			glyphHeight := metrics.GetBounds().Size.Height * c.scale

			// Get texture coordinates from glyph atlas
			atlasRect := context.GetGlyphAtlas().GetGlyphRect(glyph.GetID())
			if atlasRect == nil {
				continue
			}

			// Create quad for this glyph
			vertices = append(vertices,
				GlyphVertex{
					Position: geom.Point{X: glyphX, Y: glyphY},
					TexCoord: atlasRect.Origin,
				},
				GlyphVertex{
					Position: geom.Point{X: glyphX + glyphWidth, Y: glyphY},
					TexCoord: geom.Point{X: atlasRect.Origin.X + atlasRect.Size.Width, Y: atlasRect.Origin.Y},
				},
				GlyphVertex{
					Position: geom.Point{X: glyphX + glyphWidth, Y: glyphY + glyphHeight},
					TexCoord: geom.Point{X: atlasRect.Origin.X + atlasRect.Size.Width, Y: atlasRect.Origin.Y + atlasRect.Size.Height},
				},
				GlyphVertex{
					Position: geom.Point{X: glyphX, Y: glyphY + glyphHeight},
					TexCoord: geom.Point{X: atlasRect.Origin.X, Y: atlasRect.Origin.Y + atlasRect.Size.Height},
				},
			)

			// Advance position
			currentX += metrics.GetAdvance() * c.scale
		}

		// Move to next line if needed
		if run.IsLineBreak() {
			currentX = c.position.X
			currentY += font.GetLineHeight() * c.scale
		}
	}

	return vertices
}

// createUniforms creates uniform data for rendering
func (c *TextContents) createUniforms(transform geom.Matrix, entity Entity) *TextUniforms {
	return &TextUniforms{
		Transform: transform,
		Color:     c.color,
		Opacity:   c.GetOpacity(),
	}
}

// GlyphVertex represents a vertex for glyph rendering
type GlyphVertex struct {
	Position geom.Point
	TexCoord geom.Point
}

// TextUniforms represents uniform data for text rendering
type TextUniforms struct {
	Transform geom.Matrix
	Color     geom.Color
	Opacity   float32
}
