package geometry

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// VerticesGeometry represents custom vertex geometry
type VerticesGeometry struct {
	BaseGeometry
	vertices  []render.Vertex
	indices   []uint16
	mode      PrimitiveMode
	transform geom.Matrix
}

// PrimitiveMode defines how vertices are interpreted
type PrimitiveMode int

const (
	// PrimitiveModeTriangles triangle list
	PrimitiveModeTriangles PrimitiveMode = iota
	// PrimitiveModeTriangleStrip triangle strip
	PrimitiveModeTriangleStrip
	// PrimitiveModeTriangleFan triangle fan
	PrimitiveModeTriangleFan
	// PrimitiveModeLines line list
	PrimitiveModeLines
	// PrimitiveModeLineStrip line strip
	PrimitiveModeLineStrip
	// PrimitiveModePoints point list
	PrimitiveModePoints
)

// NewVerticesGeometry creates a new vertices geometry
func NewVerticesGeometry(vertices []render.Vertex, indices []uint16, mode PrimitiveMode) *VerticesGeometry {
	return &VerticesGeometry{
		BaseGeometry: NewBaseGeometry(),
		vertices:     vertices,
		indices:      indices,
		mode:         mode,
		transform:    geom.NewIdentityMatrix(),
	}
}

// SetVertices sets the vertex data
func (v *VerticesGeometry) SetVertices(vertices []render.Vertex) {
	v.vertices = vertices
}

// GetVertices returns the vertex data
func (v *VerticesGeometry) GetVertices() []render.Vertex {
	return v.vertices
}

// SetIndices sets the index data
func (v *VerticesGeometry) SetIndices(indices []uint16) {
	v.indices = indices
}

// GetIndices returns the index data
func (v *VerticesGeometry) GetIndices() []uint16 {
	return v.indices
}

// SetPrimitiveMode sets the primitive mode
func (v *VerticesGeometry) SetPrimitiveMode(mode PrimitiveMode) {
	v.mode = mode
}

// GetPrimitiveMode returns the primitive mode
func (v *VerticesGeometry) GetPrimitiveMode() PrimitiveMode {
	return v.mode
}

// SetTransform sets the transformation matrix
func (v *VerticesGeometry) SetTransform(transform geom.Matrix) {
	v.transform = transform
}

// GetTransform returns the current transformation matrix
func (v *VerticesGeometry) GetTransform() geom.Matrix {
	return v.transform
}

// GetBounds returns the bounds of this geometry
func (v *VerticesGeometry) GetBounds() geom.Rect {
	if len(v.vertices) == 0 {
		return geom.Rect{}
	}

	// Find min/max coordinates
	minX, maxX := v.vertices[0].Position.X, v.vertices[0].Position.X
	minY, maxY := v.vertices[0].Position.Y, v.vertices[0].Position.Y

	for _, vertex := range v.vertices[1:] {
		if vertex.Position.X < minX {
			minX = vertex.Position.X
		}
		if vertex.Position.X > maxX {
			maxX = vertex.Position.X
		}
		if vertex.Position.Y < minY {
			minY = vertex.Position.Y
		}
		if vertex.Position.Y > maxY {
			maxY = vertex.Position.Y
		}
	}

	bounds := geom.Rect{
		X:      minX,
		Y:      minY,
		Width:  maxX - minX,
		Height: maxY - minY,
	}

	return v.transform.TransformRect(bounds)
}

// GetVertexCount returns the number of vertices
func (v *VerticesGeometry) GetVertexCount() int {
	return len(v.vertices)
}

// GetIndexCount returns the number of indices
func (v *VerticesGeometry) GetIndexCount() int {
	return len(v.indices)
}

// ComputeVertices returns the vertex data
func (v *VerticesGeometry) ComputeVertices() []render.Vertex {
	return v.vertices
}

// GetIndexBuffer returns the index data
func (v *VerticesGeometry) GetIndexBuffer() []uint16 {
	return v.indices
}

// Clone creates a copy of this geometry
func (v *VerticesGeometry) Clone() Geometry {
	// Deep copy vertices and indices
	verticesCopy := make([]render.Vertex, len(v.vertices))
	copy(verticesCopy, v.vertices)

	indicesCopy := make([]uint16, len(v.indices))
	copy(indicesCopy, v.indices)

	return &VerticesGeometry{
		BaseGeometry: v.BaseGeometry.Clone().(BaseGeometry),
		vertices:     verticesCopy,
		indices:      indicesCopy,
		mode:         v.mode,
		transform:    v.transform,
	}
}
