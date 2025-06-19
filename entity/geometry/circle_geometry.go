package geometry

import (
	"math"

	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// CircleGeometry represents circular geometry
type CircleGeometry struct {
	BaseGeometry
	center      geom.Point
	radius      float32
	segments    int
	strokeWidth float32 // -1 for filled
	transform   geom.Matrix
}

// NewCircleGeometry creates a new filled circle geometry
func NewCircleGeometry(center geom.Point, radius float32) *CircleGeometry {
	return &CircleGeometry{
		BaseGeometry: NewBaseGeometry(),
		center:       center,
		radius:       radius,
		segments:     32,   // Default segments
		strokeWidth:  -1.0, // Filled
		transform:    geom.NewIdentityMatrix(),
	}
}

// NewStrokedCircleGeometry creates a new stroked circle geometry
func NewStrokedCircleGeometry(center geom.Point, radius float32, strokeWidth float32) *CircleGeometry {
	return &CircleGeometry{
		BaseGeometry: NewBaseGeometry(),
		center:       center,
		radius:       radius,
		segments:     32,
		strokeWidth:  strokeWidth,
		transform:    geom.NewIdentityMatrix(),
	}
}

// SetCenter sets the circle center
func (c *CircleGeometry) SetCenter(center geom.Point) {
	c.center = center
}

// GetCenter returns the circle center
func (c *CircleGeometry) GetCenter() geom.Point {
	return c.center
}

// SetRadius sets the circle radius
func (c *CircleGeometry) SetRadius(radius float32) {
	c.radius = radius
}

// GetRadius returns the circle radius
func (c *CircleGeometry) GetRadius() float32 {
	return c.radius
}

// SetSegments sets the number of segments for tessellation
func (c *CircleGeometry) SetSegments(segments int) {
	if segments >= 3 {
		c.segments = segments
	}
}

// GetSegments returns the number of segments
func (c *CircleGeometry) GetSegments() int {
	return c.segments
}

// SetStrokeWidth sets the stroke width (-1 for filled)
func (c *CircleGeometry) SetStrokeWidth(width float32) {
	c.strokeWidth = width
}

// GetStrokeWidth returns the stroke width
func (c *CircleGeometry) GetStrokeWidth() float32 {
	return c.strokeWidth
}

// SetTransform sets the transformation matrix
func (c *CircleGeometry) SetTransform(transform geom.Matrix) {
	c.transform = transform
}

// GetTransform returns the transformation matrix
func (c *CircleGeometry) GetTransform() geom.Matrix {
	return c.transform
}

// IsStroked returns true if this is a stroked circle
func (c *CircleGeometry) IsStroked() bool {
	return c.strokeWidth >= 0
}

// GetBounds returns the bounds of this geometry
func (c *CircleGeometry) GetBounds() geom.Rect {
	effectiveRadius := c.radius
	if c.IsStroked() {
		effectiveRadius += c.strokeWidth / 2
	}

	bounds := geom.Rect{
		X:      c.center.X - effectiveRadius,
		Y:      c.center.Y - effectiveRadius,
		Width:  effectiveRadius * 2,
		Height: effectiveRadius * 2,
	}

	return c.transform.TransformRect(bounds)
}

// GetVertexCount returns the number of vertices
func (c *CircleGeometry) GetVertexCount() int {
	if c.IsStroked() {
		return c.segments * 2 // Inner and outer vertices
	}
	return c.segments + 1 // Center + perimeter vertices
}

// GetIndexCount returns the number of indices
func (c *CircleGeometry) GetIndexCount() int {
	if c.IsStroked() {
		return c.segments * 6 // Two triangles per segment
	}
	return c.segments * 3 // One triangle per segment from center
}

// ComputeVertices generates vertex data for rendering
func (c *CircleGeometry) ComputeVertices() []render.Vertex {
	if c.IsStroked() {
		return c.computeStrokedVertices()
	}
	return c.computeFilledVertices()
}

// computeFilledVertices generates vertices for filled circle
func (c *CircleGeometry) computeFilledVertices() []render.Vertex {
	vertices := make([]render.Vertex, c.GetVertexCount())

	// Center vertex
	vertices[0] = render.Vertex{
		Position: c.center,
		TexCoord: geom.Point{X: 0.5, Y: 0.5},
	}

	// Perimeter vertices
	for i := 0; i < c.segments; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(c.segments)
		cos := float32(math.Cos(angle))
		sin := float32(math.Sin(angle))

		x := c.center.X + cos*c.radius
		y := c.center.Y + sin*c.radius

		vertices[i+1] = render.Vertex{
			Position: geom.Point{X: x, Y: y},
			TexCoord: geom.Point{
				X: 0.5 + cos*0.5,
				Y: 0.5 + sin*0.5,
			},
		}
	}

	return vertices
}

// computeStrokedVertices generates vertices for stroked circle
func (c *CircleGeometry) computeStrokedVertices() []render.Vertex {
	vertices := make([]render.Vertex, c.GetVertexCount())

	innerRadius := c.radius - c.strokeWidth/2
	outerRadius := c.radius + c.strokeWidth/2

	// Ensure inner radius is not negative
	if innerRadius < 0 {
		innerRadius = 0
	}

	for i := 0; i < c.segments; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(c.segments)
		cos := float32(math.Cos(angle))
		sin := float32(math.Sin(angle))

		// Inner vertex
		innerX := c.center.X + cos*innerRadius
		innerY := c.center.Y + sin*innerRadius
		vertices[i*2] = render.Vertex{
			Position: geom.Point{X: innerX, Y: innerY},
			TexCoord: geom.Point{X: 0, Y: float32(i) / float32(c.segments)},
		}

		// Outer vertex
		outerX := c.center.X + cos*outerRadius
		outerY := c.center.Y + sin*outerRadius
		vertices[i*2+1] = render.Vertex{
			Position: geom.Point{X: outerX, Y: outerY},
			TexCoord: geom.Point{X: 1, Y: float32(i) / float32(c.segments)},
		}
	}

	return vertices
}

// GetIndexBuffer generates index data
func (c *CircleGeometry) GetIndexBuffer() []uint16 {
	if c.IsStroked() {
		return c.getStrokedIndexBuffer()
	}
	return c.getFilledIndexBuffer()
}

// getFilledIndexBuffer generates indices for filled circle
func (c *CircleGeometry) getFilledIndexBuffer() []uint16 {
	indices := make([]uint16, c.GetIndexCount())

	for i := 0; i < c.segments; i++ {
		base := i * 3
		indices[base] = 0                              // Center
		indices[base+1] = uint16(i + 1)                // Current vertex
		indices[base+2] = uint16((i+1)%c.segments + 1) // Next vertex
	}

	return indices
}

// getStrokedIndexBuffer generates indices for stroked circle
func (c *CircleGeometry) getStrokedIndexBuffer() []uint16 {
	indices := make([]uint16, c.GetIndexCount())

	for i := 0; i < c.segments; i++ {
		base := i * 6
		current := uint16(i * 2)
		next := uint16(((i + 1) % c.segments) * 2)

		// First triangle
		indices[base] = current
		indices[base+1] = current + 1
		indices[base+2] = next

		// Second triangle
		indices[base+3] = next
		indices[base+4] = current + 1
		indices[base+5] = next + 1
	}

	return indices
}

// Clone creates a copy of this geometry
func (c *CircleGeometry) Clone() Geometry {
	return &CircleGeometry{
		BaseGeometry: c.BaseGeometry.Clone().(BaseGeometry),
		center:       c.center,
		radius:       c.radius,
		segments:     c.segments,
		strokeWidth:  c.strokeWidth,
		transform:    c.transform,
	}
}

// GetCoverage returns the coverage area for this geometry
func (c *CircleGeometry) GetCoverage(transform geom.Matrix) geom.Rect {
	combinedTransform := transform.Multiply(c.transform)
	bounds := c.GetBounds()
	return combinedTransform.TransformRect(bounds)
}
