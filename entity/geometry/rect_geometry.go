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
