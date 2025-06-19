package contents

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// GradientType defines the type of gradient
type GradientType int

const (
	GradientTypeLinear GradientType = iota
	GradientTypeRadial
	GradientTypeConical
	GradientTypeSweep
)

// GradientStop represents a color stop in a gradient
type GradientStop struct {
	Position float32    // Position along the gradient (0.0 to 1.0)
	Color    geom.Color // Color at this position
}

// BaseGradientContents provides base functionality for gradient contents
type BaseGradientContents struct {
	*BaseContents
	gradientType GradientType
	stops        []GradientStop
	tileMode     TileMode
}

// TileMode defines how gradients are tiled beyond their bounds
type TileMode int

const (
	TileModeClamp TileMode = iota
	TileModeRepeat
	TileModeMirror
	TileModeDecal
)

// NewBaseGradientContents creates a new base gradient contents
func NewBaseGradientContents(gradientType GradientType) *BaseGradientContents {
	return &BaseGradientContents{
		BaseContents: NewBaseContents(),
		gradientType: gradientType,
		stops:        make([]GradientStop, 0),
		tileMode:     TileModeClamp,
	}
}

// AddStop adds a color stop to the gradient
func (c *BaseGradientContents) AddStop(position float32, color geom.Color) {
	c.stops = append(c.stops, GradientStop{
		Position: position,
		Color:    color,
	})
}

// GetStops returns the gradient stops
func (c *BaseGradientContents) GetStops() []GradientStop {
	return c.stops
}

// SetTileMode sets the tile mode for the gradient
func (c *BaseGradientContents) SetTileMode(mode TileMode) {
	c.tileMode = mode
}

// GetTileMode returns the tile mode
func (c *BaseGradientContents) GetTileMode() TileMode {
	return c.tileMode
}

// GetGradientType returns the gradient type
func (c *BaseGradientContents) GetGradientType() GradientType {
	return c.gradientType
}

// LinearGradientContents renders linear gradients
type LinearGradientContents struct {
	*BaseGradientContents
	startPoint geom.Point
	endPoint   geom.Point
}

// NewLinearGradientContents creates a new linear gradient contents
func NewLinearGradientContents() *LinearGradientContents {
	return &LinearGradientContents{
		BaseGradientContents: NewBaseGradientContents(GradientTypeLinear),
		startPoint:           geom.Point{X: 0, Y: 0},
		endPoint:             geom.Point{X: 1, Y: 0},
	}
}

// SetStartPoint sets the start point of the linear gradient
func (c *LinearGradientContents) SetStartPoint(point geom.Point) {
	c.startPoint = point
}

// GetStartPoint returns the start point
func (c *LinearGradientContents) GetStartPoint() geom.Point {
	return c.startPoint
}

// SetEndPoint sets the end point of the linear gradient
func (c *LinearGradientContents) SetEndPoint(point geom.Point) {
	c.endPoint = point
}

// GetEndPoint returns the end point
func (c *LinearGradientContents) GetEndPoint() geom.Point {
	return c.endPoint
}

// GetCoverage returns the coverage area of this content
func (c *LinearGradientContents) GetCoverage(transform geom.Matrix) *geom.Rect {
	// For linear gradients, coverage is typically infinite
	// Return a large rect for practical purposes
	maxFloat := float32(1e6)
	bounds := geom.Rect{
		Origin: geom.Point{X: -maxFloat, Y: -maxFloat},
		Size:   geom.Size{Width: 2 * maxFloat, Height: 2 * maxFloat},
	}
	transformedBounds := transform.TransformRect(bounds)
	return &transformedBounds
}

// Render renders this content using the provided context and render pass
func (c *LinearGradientContents) Render(
	context *ContentContext,
	pass *render.RenderPass,
	transform geom.Matrix,
	entity Entity,
) bool {
	if len(c.stops) == 0 || context == nil || pass == nil {
		return false
	}

	// Get the linear gradient shader
	shader := context.GetLinearGradientShader()
	if shader == nil {
		return false
	}

	// Create pipeline
	pipeline := context.GetPipelineLibrary().GetLinearGradientPipeline(c.GetBlendMode())
	if pipeline == nil {
		return false
	}

	// Set up vertex data (typically a full-screen quad)
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

// Clone creates a copy of this content
func (c *LinearGradientContents) Clone() Contents {
	clone := NewLinearGradientContents()
	clone.SetStartPoint(c.startPoint)
	clone.SetEndPoint(c.endPoint)
	clone.SetBlendMode(c.GetBlendMode())
	clone.SetOpacity(c.GetOpacity())
	clone.SetTileMode(c.GetTileMode())

	// Copy stops
	for _, stop := range c.stops {
		clone.AddStop(stop.Position, stop.Color)
	}

	return clone
}

// generateVertices generates vertex data for gradient rendering
func (c *LinearGradientContents) generateVertices(transform geom.Matrix) []GradientVertex {
	// Create a full-screen quad for gradient rendering
	// TODO: Optimize this to only cover the actual gradient area
	return []GradientVertex{
		{Position: geom.Point{X: -1, Y: -1}},
		{Position: geom.Point{X: 1, Y: -1}},
		{Position: geom.Point{X: 1, Y: 1}},
		{Position: geom.Point{X: -1, Y: 1}},
	}
}

// createUniforms creates uniform data for rendering
func (c *LinearGradientContents) createUniforms(transform geom.Matrix, entity Entity) *LinearGradientUniforms {
	return &LinearGradientUniforms{
		Transform:  transform,
		StartPoint: c.startPoint,
		EndPoint:   c.endPoint,
		Opacity:    c.GetOpacity(),
		TileMode:   int32(c.tileMode),
		StopCount:  int32(len(c.stops)),
		Stops:      c.convertStops(),
	}
}

// convertStops converts gradient stops to shader format
func (c *LinearGradientContents) convertStops() [16]GradientStopUniform {
	var stops [16]GradientStopUniform

	for i, stop := range c.stops {
		if i >= 16 {
			break // Limit to 16 stops for shader compatibility
		}
		stops[i] = GradientStopUniform{
			Position: stop.Position,
			Color:    stop.Color,
		}
	}

	return stops
}

// RadialGradientContents renders radial gradients
type RadialGradientContents struct {
	*BaseGradientContents
	center geom.Point
	radius float32
}

// NewRadialGradientContents creates a new radial gradient contents
func NewRadialGradientContents() *RadialGradientContents {
	return &RadialGradientContents{
		BaseGradientContents: NewBaseGradientContents(GradientTypeRadial),
		center:               geom.Point{X: 0.5, Y: 0.5},
		radius:               0.5,
	}
}

// SetCenter sets the center point of the radial gradient
func (c *RadialGradientContents) SetCenter(center geom.Point) {
	c.center = center
}

// GetCenter returns the center point
func (c *RadialGradientContents) GetCenter() geom.Point {
	return c.center
}

// SetRadius sets the radius of the radial gradient
func (c *RadialGradientContents) SetRadius(radius float32) {
	c.radius = radius
}

// GetRadius returns the radius
func (c *RadialGradientContents) GetRadius() float32 {
	return c.radius
}

// GetCoverage returns the coverage area of this content
func (c *RadialGradientContents) GetCoverage(transform geom.Matrix) *geom.Rect {
	// For radial gradients, coverage is typically the circle bounds
	bounds := geom.Rect{
		Origin: geom.Point{X: c.center.X - c.radius, Y: c.center.Y - c.radius},
		Size:   geom.Size{Width: 2 * c.radius, Height: 2 * c.radius},
	}
	transformedBounds := transform.TransformRect(bounds)
	return &transformedBounds
}

// Render renders this content using the provided context and render pass
func (c *RadialGradientContents) Render(
	context *ContentContext,
	pass *render.RenderPass,
	transform geom.Matrix,
	entity Entity,
) bool {
	// Similar to linear gradient rendering but with radial shader
	// TODO: Implement radial gradient rendering
	return false
}

// Clone creates a copy of this content
func (c *RadialGradientContents) Clone() Contents {
	clone := NewRadialGradientContents()
	clone.SetCenter(c.center)
	clone.SetRadius(c.radius)
	clone.SetBlendMode(c.GetBlendMode())
	clone.SetOpacity(c.GetOpacity())
	clone.SetTileMode(c.GetTileMode())

	// Copy stops
	for _, stop := range c.stops {
		clone.AddStop(stop.Position, stop.Color)
	}

	return clone
}

// GradientVertex represents a vertex for gradient rendering
type GradientVertex struct {
	Position geom.Point
}

// LinearGradientUniforms represents uniform data for linear gradient rendering
type LinearGradientUniforms struct {
	Transform  geom.Matrix
	StartPoint geom.Point
	EndPoint   geom.Point
	Opacity    float32
	TileMode   int32
	StopCount  int32
	Stops      [16]GradientStopUniform
}

// GradientStopUniform represents a gradient stop in uniform buffer format
type GradientStopUniform struct {
	Position float32
	Color    geom.Color
}
