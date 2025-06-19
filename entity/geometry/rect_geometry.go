package geometry

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// RectGeometry represents rectangular geometry
type RectGeometry struct {
	BaseGeometry
	rect      geom.Rect
	transform geom.Matrix
}

// NewRectGeometry creates a new rectangular geometry
func NewRectGeometry(rect geom.Rect) *RectGeometry {
	return &RectGeometry{
		BaseGeometry: NewBaseGeometry(),
		rect:         rect,
		transform:    geom.NewIdentityMatrix(),
	}
}

// SetRect sets the rectangle
func (r *RectGeometry) SetRect(rect geom.Rect) {
	r.rect = rect
}

// GetRect returns the current rectangle
func (r *RectGeometry) GetRect() geom.Rect {
	return r.rect
}

// SetTransform sets the transformation matrix
func (r *RectGeometry) SetTransform(transform geom.Matrix) {
	r.transform = transform
}

// GetTransform returns the current transformation matrix
func (r *RectGeometry) GetTransform() geom.Matrix {
	return r.transform
}

// GetBounds returns the bounds of this geometry
func (r *RectGeometry) GetBounds() geom.Rect {
	return r.transform.TransformRect(r.rect)
}

// GetVertexCount returns the number of vertices
func (r *RectGeometry) GetVertexCount() int {
	return 4 // Rectangle has 4 vertices
}

// GetIndexCount returns the number of indices
func (r *RectGeometry) GetIndexCount() int {
	return 6 // Two triangles = 6 indices
}

// GetPositionBuffer generates position vertex data
func (r *RectGeometry) GetPositionBuffer() []float32 {
	bounds := r.GetBounds()
	return []float32{
		bounds.X, bounds.Y, // Bottom-left
		bounds.X + bounds.Width, bounds.Y, // Bottom-right
		bounds.X + bounds.Width, bounds.Y + bounds.Height, // Top-right
		bounds.X, bounds.Y + bounds.Height, // Top-left
	}
}

// GetTexCoordBuffer generates texture coordinate data
func (r *RectGeometry) GetTexCoordBuffer() []float32 {
	return []float32{
		0.0, 0.0, // Bottom-left
		1.0, 0.0, // Bottom-right
		1.0, 1.0, // Top-right
		0.0, 1.0, // Top-left
	}
}

// GetIndexBuffer generates index data
func (r *RectGeometry) GetIndexBuffer() []uint16 {
	return []uint16{
		0, 1, 2, // First triangle
		2, 3, 0, // Second triangle
	}
}

// ComputeVertices generates vertex data for rendering
func (r *RectGeometry) ComputeVertices() []render.Vertex {
	positions := r.GetPositionBuffer()
	texCoords := r.GetTexCoordBuffer()

	vertices := make([]render.Vertex, 4)
	for i := 0; i < 4; i++ {
		vertices[i] = render.Vertex{
			Position: geom.Point{
				X: positions[i*2],
				Y: positions[i*2+1],
			},
			TexCoord: geom.Point{
				X: texCoords[i*2],
				Y: texCoords[i*2+1],
			},
		}
	}

	return vertices
}

// Clone creates a copy of this geometry
func (r *RectGeometry) Clone() Geometry {
	return &RectGeometry{
		BaseGeometry: r.BaseGeometry.Clone().(BaseGeometry),
		rect:         r.rect,
		transform:    r.transform,
	}
}

// CircleGeometry represents circular geometry
type CircleGeometry struct {
	BaseGeometry
	center    geom.Point
	radius    float32
	segments  int
	transform geom.Matrix
}

// NewCircleGeometry creates a new circular geometry
func NewCircleGeometry(center geom.Point, radius float32, segments int) *CircleGeometry {
	if segments < 3 {
		segments = 32 // Default segments
	}

	return &CircleGeometry{
		BaseGeometry: NewBaseGeometry(),
		center:       center,
		radius:       radius,
		segments:     segments,
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

// SetSegments sets the number of segments
func (c *CircleGeometry) SetSegments(segments int) {
	if segments >= 3 {
		c.segments = segments
	}
}

// GetSegments returns the number of segments
func (c *CircleGeometry) GetSegments() int {
	return c.segments
}

// SetTransform sets the transformation matrix
func (c *CircleGeometry) SetTransform(transform geom.Matrix) {
	c.transform = transform
}

// GetTransform returns the current transformation matrix
func (c *CircleGeometry) GetTransform() geom.Matrix {
	return c.transform
}

// GetBounds returns the bounds of this geometry
func (c *CircleGeometry) GetBounds() geom.Rect {
	rect := geom.Rect{
		X:      c.center.X - c.radius,
		Y:      c.center.Y - c.radius,
		Width:  c.radius * 2,
		Height: c.radius * 2,
	}
	return c.transform.TransformRect(rect)
}

// GetVertexCount returns the number of vertices
func (c *CircleGeometry) GetVertexCount() int {
	return c.segments + 1 // Center vertex + perimeter vertices
}

// GetIndexCount returns the number of indices
func (c *CircleGeometry) GetIndexCount() int {
	return c.segments * 3 // Each segment forms a triangle with center
}

// ComputeVertices generates vertex data for rendering
func (c *CircleGeometry) ComputeVertices() []render.Vertex {
	vertices := make([]render.Vertex, c.GetVertexCount())

	// Center vertex
	vertices[0] = render.Vertex{
		Position: c.center,
		TexCoord: geom.Point{X: 0.5, Y: 0.5},
	}

	// Perimeter vertices
	for i := 0; i < c.segments; i++ {
		angle := float32(i) * 2.0 * geom.Pi / float32(c.segments)
		cos := geom.Cos(angle)
		sin := geom.Sin(angle)

		vertices[i+1] = render.Vertex{
			Position: geom.Point{
				X: c.center.X + cos*c.radius,
				Y: c.center.Y + sin*c.radius,
			},
			TexCoord: geom.Point{
				X: 0.5 + cos*0.5,
				Y: 0.5 + sin*0.5,
			},
		}
	}

	return vertices
}

// GetIndexBuffer generates index data
func (c *CircleGeometry) GetIndexBuffer() []uint16 {
	indices := make([]uint16, c.GetIndexCount())

	for i := 0; i < c.segments; i++ {
		base := i * 3
		indices[base] = 0                              // Center
		indices[base+1] = uint16(i + 1)                // Current perimeter vertex
		indices[base+2] = uint16((i+1)%c.segments + 1) // Next perimeter vertex
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
		transform:    c.transform,
	}
}
