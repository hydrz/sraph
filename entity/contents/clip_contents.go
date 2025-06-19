package contents

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// ClipOperation defines the type of clipping operation
type ClipOperation int

const (
	ClipOperationDifference ClipOperation = iota
	ClipOperationIntersect
)

// ClipContents performs clipping operations on subsequent rendering
type ClipContents struct {
	*BaseContents
	path          *geom.Path
	clipOperation ClipOperation
	antiAliased   bool
}

// NewClipContents creates a new clip contents
func NewClipContents() *ClipContents {
	return &ClipContents{
		BaseContents:  NewBaseContents(),
		path:          nil,
		clipOperation: ClipOperationIntersect,
		antiAliased:   true,
	}
}

// SetPath sets the clipping path
func (c *ClipContents) SetPath(path *geom.Path) {
	c.path = path
}

// GetPath returns the clipping path
func (c *ClipContents) GetPath() *geom.Path {
	return c.path
}

// SetClipOperation sets the clip operation type
func (c *ClipContents) SetClipOperation(operation ClipOperation) {
	c.clipOperation = operation
}

// GetClipOperation returns the clip operation type
func (c *ClipContents) GetClipOperation() ClipOperation {
	return c.clipOperation
}

// SetAntiAliased sets whether the clip should be anti-aliased
func (c *ClipContents) SetAntiAliased(antiAliased bool) {
	c.antiAliased = antiAliased
}

// IsAntiAliased returns whether the clip is anti-aliased
func (c *ClipContents) IsAntiAliased() bool {
	return c.antiAliased
}

// GetCoverage returns the coverage area of this content
func (c *ClipContents) GetCoverage(transform geom.Matrix) *geom.Rect {
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
func (c *ClipContents) Render(
	context *ContentContext,
	pass *render.RenderPass,
	transform geom.Matrix,
	entity Entity,
) bool {
	if c.path == nil || context == nil || pass == nil {
		return false
	}

	// Clip contents typically modify the stencil buffer rather than rendering directly
	return c.renderToStencil(context, pass, transform, entity)
}

// renderToStencil renders the clip to the stencil buffer
func (c *ClipContents) renderToStencil(
	context *ContentContext,
	pass *render.RenderPass,
	transform geom.Matrix,
	entity Entity,
) bool {

	// Get the stencil shader
	shader := context.GetStencilShader()
	if shader == nil {
		return false
	}

	// Create stencil pipeline
	pipeline := context.GetPipelineLibrary().GetStencilPipeline(c.clipOperation)
	if pipeline == nil {
		return false
	}

	// Generate vertices for the clip path
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

	// Configure stencil state
	stencilState := c.createStencilState()
	pass.SetStencilState(stencilState)

	// Bind resources
	pass.SetPipeline(pipeline)
	pass.SetVertexBuffer(vertexBuffer)
	pass.BindUniformBuffer(uniformBuffer, 0)

	// Draw to stencil buffer
	return pass.Draw(len(vertices))
}

// IsOpaque returns true if this content is fully opaque
func (c *ClipContents) IsOpaque() bool {
	// Clip contents don't contribute to color, so they're not opaque
	return false
}

// Clone creates a copy of this content
func (c *ClipContents) Clone() Contents {
	clone := NewClipContents()
	clone.SetClipOperation(c.clipOperation)
	clone.SetAntiAliased(c.antiAliased)
	clone.SetBlendMode(c.GetBlendMode())
	clone.SetOpacity(c.GetOpacity())
	if c.path != nil {
		clone.SetPath(c.path.Copy())
	}
	return clone
}

// generateVertices generates vertex data for the clip path
func (c *ClipContents) generateVertices(transform geom.Matrix) []ClipVertex {
	if c.path == nil {
		return nil
	}

	// TODO: Implement proper path tessellation for clipping
	// This would typically involve:
	// 1. Tessellating the path into triangles
	// 2. Applying the transform
	// 3. Generating appropriate vertices for stencil rendering

	// For now, return empty vertices as a placeholder
	return []ClipVertex{}
}

// createUniforms creates uniform data for clip rendering
func (c *ClipContents) createUniforms(transform geom.Matrix, entity Entity) *ClipUniforms {
	return &ClipUniforms{
		Transform: transform,
		ClipOp:    int32(c.clipOperation),
	}
}

// createStencilState creates the stencil state for clipping
func (c *ClipContents) createStencilState() *render.StencilState {
	switch c.clipOperation {
	case ClipOperationIntersect:
		return &render.StencilState{
			StencilTest:   true,
			StencilFunc:   render.StencilFuncAlways,
			StencilRef:    1,
			StencilMask:   0xFF,
			StencilPassOp: render.StencilOpIncrement,
			StencilFailOp: render.StencilOpKeep,
			DepthFailOp:   render.StencilOpKeep,
		}
	case ClipOperationDifference:
		return &render.StencilState{
			StencilTest:   true,
			StencilFunc:   render.StencilFuncAlways,
			StencilRef:    1,
			StencilMask:   0xFF,
			StencilPassOp: render.StencilOpDecrement,
			StencilFailOp: render.StencilOpKeep,
			DepthFailOp:   render.StencilOpKeep,
		}
	default:
		return &render.StencilState{
			StencilTest: false,
		}
	}
}

// GetClipBounds returns the bounds of the clip area
func (c *ClipContents) GetClipBounds(transform geom.Matrix) *geom.Rect {
	return c.GetCoverage(transform)
}

// ClipVertex represents a vertex for clip rendering
type ClipVertex struct {
	Position geom.Point
}

// ClipUniforms represents uniform data for clip rendering
type ClipUniforms struct {
	Transform geom.Matrix
	ClipOp    int32
}
