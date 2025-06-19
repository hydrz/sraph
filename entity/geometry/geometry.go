// Package geometry provides geometric content types for entity rendering.
//
// This package defines geometric primitives that can be rendered as entities,
// including paths, rectangles, circles, and complex shapes.
package geometry

import (
	"github.com/opensraph/sraph/geom"
)

// Geometry represents a geometric shape that can be rendered.
type Geometry interface {
	// GetVertices returns the vertex data for this geometry.
	GetVertices() ([]Vertex, error)

	// GetIndices returns the index data for triangulation.
	GetIndices() ([]uint32, error)

	// GetBounds returns the bounding box of this geometry.
	GetBounds() geom.Rect[geom.F32]

	// IsConvex returns true if this geometry is convex.
	IsConvex() bool

	// GetType returns the type of geometry.
	GetType() GeometryType
}

// GeometryType identifies the type of geometry.
type GeometryType int

const (
	GeometryTypeRect GeometryType = iota
	GeometryTypeRoundRect
	GeometryTypeCircle
	GeometryTypeEllipse
	GeometryTypePath
	GeometryTypePolygon
	GeometryTypeText
)

// Vertex represents a single vertex in geometry.
type Vertex struct {
	// Position in 2D space
	Position geom.Point[geom.F32]

	// Texture coordinates
	TexCoord geom.Point[geom.F32]

	// Vertex color
	Color geom.Color

	// Normal vector (for 3D effects)
	Normal geom.Vector3[geom.F32]
}

// GeometryBuilder provides a fluent interface for building complex geometries.
type GeometryBuilder struct {
	vertices []Vertex
	indices  []uint32
	bounds   geom.Rect[geom.F32]
}

// NewGeometryBuilder creates a new geometry builder.
func NewGeometryBuilder() *GeometryBuilder {
	return &GeometryBuilder{
		vertices: make([]Vertex, 0),
		indices:  make([]uint32, 0),
	}
}

// AddVertex adds a vertex to the geometry.
func (gb *GeometryBuilder) AddVertex(vertex Vertex) *GeometryBuilder {
	gb.vertices = append(gb.vertices, vertex)
	gb.updateBounds(vertex.Position)
	return gb
}

// AddTriangle adds a triangle using vertex indices.
func (gb *GeometryBuilder) AddTriangle(i0, i1, i2 uint32) *GeometryBuilder {
	gb.indices = append(gb.indices, i0, i1, i2)
	return gb
}

// Build creates a custom geometry from the builder.
func (gb *GeometryBuilder) Build() Geometry {
	return &CustomGeometry{
		vertices: gb.vertices,
		indices:  gb.indices,
		bounds:   gb.bounds,
	}
}

func (gb *GeometryBuilder) updateBounds(point geom.Point[geom.F32]) {
	if len(gb.vertices) == 1 {
		// First vertex, initialize bounds
		gb.bounds = geom.NewRectXYWH(point.X, point.Y, 0, 0)
	} else {
		// Expand bounds to include new point
		gb.bounds = gb.bounds.ExpandToInclude(point)
	}
}

// CustomGeometry represents user-defined geometry.
type CustomGeometry struct {
	vertices []Vertex
	indices  []uint32
	bounds   geom.Rect[geom.F32]
}

// GetVertices implements Geometry.
func (cg *CustomGeometry) GetVertices() ([]Vertex, error) {
	return cg.vertices, nil
}

// GetIndices implements Geometry.
func (cg *CustomGeometry) GetIndices() ([]uint32, error) {
	return cg.indices, nil
}

// GetBounds implements Geometry.
func (cg *CustomGeometry) GetBounds() geom.Rect[geom.F32] {
	return cg.bounds
}

// IsConvex implements Geometry.
func (cg *CustomGeometry) IsConvex() bool {
	// TODO: Implement convexity test for custom geometry
	return false // Conservative default
}

// GetType implements Geometry.
func (cg *CustomGeometry) GetType() GeometryType {
	return GeometryTypePolygon // Default for custom geometry
}
