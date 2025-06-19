package contents

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// SolidColorContents renders solid color fills
type SolidColorContents struct {
	*BaseContents
	color geom.Color
	path  *geom.Path
}

// NewSolidColorContents creates a new solid color contents
func NewSolidColorContents() *SolidColorContents {
	return &SolidColorContents{
		BaseContents: NewBaseContents(),
		color:        geom.ColorWhite,
		path:         nil,
	}
}

// SetColor sets the fill color
func (c *SolidColorContents) SetColor(color geom.Color) {
	c.color = color
}

// GetColor returns the fill color
func (c *SolidColorContents) GetColor() geom.Color {
	return c.color
}

// SetPath sets the path to fill
func (c *SolidColorContents) SetPath(path *geom.Path) {
	c.path = path
}

// GetPath returns the path to fill
func (c *SolidColorContents) GetPath() *geom.Path {
	return c.path
}

// GetCoverage returns the coverage area of this content
func (c *SolidColorContents) GetCoverage(transform geom.Matrix) *geom.Rect {
	if c.path == nil {
		return nil
	}

	bounds := c.path.GetBounds()
	if bounds == nil {
		return nil
	}

	transformedBounds := transform.TransformRect(*bounds)
	return &transformedBounds
}

// Render renders this content using the provided context and render pass
func (c *SolidColorContents) Render(
	context *ContentContext,
	pass *render.RenderPass,
	transform geom.Matrix,
	entity Entity,
) bool {
	if c.path == nil || context == nil || pass == nil {
		return false
	}

	// Get the solid color shader
	shader := context.GetSolidColorShader()
	if shader == nil {
		return false
	}

	// Create pipeline
	pipeline := context.GetPipelineLibrary().GetSolidColorPipeline(c.GetBlendMode())
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

	// Bind resources
	pass.SetPipeline(pipeline)
	pass.SetVertexBuffer(vertexBuffer)
	pass.BindUniformBuffer(uniformBuffer, 0)

	// Draw
	return pass.Draw(len(vertices))
}

// IsOpaque returns true if this content is fully opaque
func (c *SolidColorContents) IsOpaque() bool {
	return c.BaseContents.IsOpaque() && c.color.Alpha >= 1.0
}

// Clone creates a copy of this content
func (c *SolidColorContents) Clone() Contents {
	clone := NewSolidColorContents()
	clone.SetColor(c.color)
	clone.SetBlendMode(c.GetBlendMode())
	clone.SetOpacity(c.GetOpacity())
	if c.path != nil {
		clone.SetPath(c.path.Copy())
	}
	return clone
}

// generateVertices generates vertex data for the path
func (c *SolidColorContents) generateVertices(transform geom.Matrix) []SolidColorVertex {
	if c.path == nil {
		return nil
	}

	// TODO: Implement proper path tessellation
	// For now, return empty vertices
	return []SolidColorVertex{}
}

// createUniforms creates uniform data for rendering
func (c *SolidColorContents) createUniforms(transform geom.Matrix, entity Entity) *SolidColorUniforms {
	return &SolidColorUniforms{
		Color:     c.color,
		Transform: transform,
		Opacity:   c.GetOpacity(),
	}
}

// SolidColorVertex represents a vertex for solid color rendering
type SolidColorVertex struct {
	Position geom.Point
}

// SolidColorUniforms represents uniform data for solid color rendering
type SolidColorUniforms struct {
	Color     geom.Color
	Transform geom.Matrix
	Opacity   float32
}
