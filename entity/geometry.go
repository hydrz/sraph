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

// RectGeometry represents a rectangular geometry.
type RectGeometry struct {
	rect geom.Rect[geom.F32]
}

// NewRectGeometry creates a new rectangular geometry.
func NewRectGeometry(rect geom.Rect[geom.F32]) *RectGeometry {
	return &RectGeometry{rect: rect}
}

// GetVertices implements Geometry.
func (r *RectGeometry) GetVertices() ([]Vertex, error) {
	vertices := []Vertex{
		{Position: r.rect.LeftTop(), TexCoord: geom.Pt[geom.F32](0, 0)},
		{Position: r.rect.RightTop(), TexCoord: geom.Pt[geom.F32](1, 0)},
		{Position: r.rect.RightBottom(), TexCoord: geom.Pt[geom.F32](1, 1)},
		{Position: r.rect.LeftBottom(), TexCoord: geom.Pt[geom.F32](0, 1)},
	}
	return vertices, nil
}

// GetIndices implements Geometry.
func (r *RectGeometry) GetIndices() ([]uint32, error) {
	// Two triangles for a rectangle
	indices := []uint32{0, 1, 2, 0, 2, 3}
	return indices, nil
}

// GetBounds implements Geometry.
func (r *RectGeometry) GetBounds() geom.Rect[geom.F32] {
	return r.rect
}

// IsConvex implements Geometry.
func (r *RectGeometry) IsConvex() bool {
	return true
}

// GetType implements Geometry.
func (r *RectGeometry) GetType() GeometryType {
	return GeometryTypeRect
}

// CircleGeometry represents a circular geometry.
type CircleGeometry struct {
	center   geom.Point[geom.F32]
	radius   geom.F32
	segments int
}

// NewCircleGeometry creates a new circular geometry.
func NewCircleGeometry(center geom.Point[geom.F32], radius geom.F32, segments int) *CircleGeometry {
	if segments < 3 {
		segments = 32 // Default reasonable number of segments
	}
	return &CircleGeometry{
		center:   center,
		radius:   radius,
		segments: segments,
	}
}

// GetVertices implements Geometry.
func (c *CircleGeometry) GetVertices() ([]Vertex, error) {
	vertices := make([]Vertex, c.segments+1) // +1 for center vertex

	// Center vertex
	vertices[0] = Vertex{
		Position: c.center,
		TexCoord: geom.Pt[geom.F32](0.5, 0.5),
	}

	// Circle vertices
	angleStep := 2.0 * geom.Pi / geom.F32(c.segments)
	for i := 0; i < c.segments; i++ {
		angle := geom.F32(i) * angleStep
		x := c.center.X + c.radius*geom.F32(geom.Cos(float64(angle)))
		y := c.center.Y + c.radius*geom.F32(geom.Sin(float64(angle)))

		// Texture coordinates map from [-1,1] to [0,1]
		texX := (x-c.center.X)/c.radius*0.5 + 0.5
		texY := (y-c.center.Y)/c.radius*0.5 + 0.5

		vertices[i+1] = Vertex{
			Position: geom.Pt[geom.F32](x, y),
			TexCoord: geom.Pt[geom.F32](texX, texY),
		}
	}

	return vertices, nil
}

// GetIndices implements Geometry.
func (c *CircleGeometry) GetIndices() ([]uint32, error) {
	indices := make([]uint32, c.segments*3)

	for i := 0; i < c.segments; i++ {
		base := i * 3
		indices[base] = 0 // Center vertex
		indices[base+1] = uint32(i + 1)
		indices[base+2] = uint32((i+1)%c.segments + 1)
	}

	return indices, nil
}

// GetBounds implements Geometry.
func (c *CircleGeometry) GetBounds() geom.Rect[geom.F32] {
	return geom.NewRectLTRB(
		c.center.X-c.radius,
		c.center.Y-c.radius,
		c.center.X+c.radius,
		c.center.Y+c.radius,
	)
}

// IsConvex implements Geometry.
func (c *CircleGeometry) IsConvex() bool {
	return true
}

// GetType implements Geometry.
func (c *CircleGeometry) GetType() GeometryType {
	return GeometryTypeCircle
}

// PathGeometry represents a path-based geometry.
type PathGeometry struct {
	path        geom.PathSource[geom.F32]
	tessellated []Vertex
	indices     []uint32
	bounds      geom.Rect[geom.F32]
	needsUpdate bool
}

// NewPathGeometry creates a new path geometry.
func NewPathGeometry(path geom.PathSource[geom.F32]) *PathGeometry {
	return &PathGeometry{
		path:        path,
		bounds:      path.Bounds(),
		needsUpdate: true,
	}
}

// GetVertices implements Geometry.
func (p *PathGeometry) GetVertices() ([]Vertex, error) {
	if p.needsUpdate {
		if err := p.tessellate(); err != nil {
			return nil, err
		}
	}
	return p.tessellated, nil
}

// GetIndices implements Geometry.
func (p *PathGeometry) GetIndices() ([]uint32, error) {
	if p.needsUpdate {
		if err := p.tessellate(); err != nil {
			return nil, err
		}
	}
	return p.indices, nil
}

// GetBounds implements Geometry.
func (p *PathGeometry) GetBounds() geom.Rect[geom.F32] {
	return p.bounds
}

// IsConvex implements Geometry.
func (p *PathGeometry) IsConvex() bool {
	return p.path.IsConvex()
}

// GetType implements Geometry.
func (p *PathGeometry) GetType() GeometryType {
	return GeometryTypePath
}

// tessellate converts the path to triangle vertices.
func (p *PathGeometry) tessellate() error {
	// TODO: Implement path tessellation
	// This would use the tess package to convert the path to triangles
	// For now, create a simple placeholder

	bounds := p.path.Bounds()
	p.tessellated = []Vertex{
		{Position: bounds.LeftTop()},
		{Position: bounds.RightTop()},
		{Position: bounds.RightBottom()},
		{Position: bounds.LeftBottom()},
	}

	p.indices = []uint32{0, 1, 2, 0, 2, 3}
	p.needsUpdate = false

	return nil
}

// RoundRectGeometry represents a rounded rectangle geometry.
type RoundRectGeometry struct {
	roundRect geom.RoundRect[geom.F32]
	segments  int
}

// NewRoundRectGeometry creates a new rounded rectangle geometry.
func NewRoundRectGeometry(roundRect geom.RoundRect[geom.F32], segments int) *RoundRectGeometry {
	if segments < 4 {
		segments = 8 // Default segments per corner
	}
	return &RoundRectGeometry{
		roundRect: roundRect,
		segments:  segments,
	}
}

// GetVertices implements Geometry.
func (rr *RoundRectGeometry) GetVertices() ([]Vertex, error) {
	// TODO: Implement rounded rectangle tessellation
	// For now, approximate with regular rectangle
	rect := rr.roundRect.Rect()
	rectGeom := NewRectGeometry(rect)
	return rectGeom.GetVertices()
}

// GetIndices implements Geometry.
func (rr *RoundRectGeometry) GetIndices() ([]uint32, error) {
	// TODO: Implement rounded rectangle tessellation
	// For now, approximate with regular rectangle
	rect := rr.roundRect.Rect()
	rectGeom := NewRectGeometry(rect)
	return rectGeom.GetIndices()
}

// GetBounds implements Geometry.
func (rr *RoundRectGeometry) GetBounds() geom.Rect[geom.F32] {
	return rr.roundRect.Rect()
}

// IsConvex implements Geometry.
func (rr *RoundRectGeometry) IsConvex() bool {
	return true
}

// GetType implements Geometry.
func (rr *RoundRectGeometry) GetType() GeometryType {
	return GeometryTypeRoundRect
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
