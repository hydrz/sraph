package geometry

import (
	"github.com/hydrz/sraph/geom"
	"github.com/hydrz/sraph/render"
)

// PathGeometry represents path-based geometry
type PathGeometry struct {
	BaseGeometry
	path      geom.Path
	fillRule  FillRule
	transform geom.Matrix
}

// FillRule defines how path filling is determined
type FillRule int

const (
	// FillRuleNonZero non-zero winding rule
	FillRuleNonZero FillRule = iota
	// FillRuleEvenOdd even-odd rule
	FillRuleEvenOdd
)

// NewPathGeometry creates a new path geometry
func NewPathGeometry(path geom.Path) *PathGeometry {
	return &PathGeometry{
		BaseGeometry: NewBaseGeometry(),
		path:         path,
		fillRule:     FillRuleNonZero,
		transform:    geom.NewIdentityMatrix(),
	}
}

// SetPath sets the path
func (p *PathGeometry) SetPath(path geom.Path) {
	p.path = path
}

// GetPath returns the current path
func (p *PathGeometry) GetPath() geom.Path {
	return p.path
}

// SetFillRule sets the fill rule
func (p *PathGeometry) SetFillRule(rule FillRule) {
	p.fillRule = rule
}

// GetFillRule returns the current fill rule
func (p *PathGeometry) GetFillRule() FillRule {
	return p.fillRule
}

// SetTransform sets the transformation matrix
func (p *PathGeometry) SetTransform(transform geom.Matrix) {
	p.transform = transform
}

// GetTransform returns the current transformation matrix
func (p *PathGeometry) GetTransform() geom.Matrix {
	return p.transform
}

// GetBounds returns the bounds of this geometry
func (p *PathGeometry) GetBounds() geom.Rect {
	bounds := p.path.GetBounds()
	return p.transform.TransformRect(bounds)
}

// GetVertexCount returns the number of vertices (estimated)
func (p *PathGeometry) GetVertexCount() int {
	// TODO: Calculate based on path complexity
	return 0
}

// GetIndexCount returns the number of indices (estimated)
func (p *PathGeometry) GetIndexCount() int {
	// TODO: Calculate based on triangulation
	return 0
}

// ComputeVertices generates vertex data for rendering
func (p *PathGeometry) ComputeVertices() []render.Vertex {
	// TODO: Implement path tessellation
	// This would involve:
	// 1. Tessellating the path into triangles
	// 2. Handling curves and complex shapes
	// 3. Applying fill rule
	return []render.Vertex{}
}

// GetIndexBuffer generates index data
func (p *PathGeometry) GetIndexBuffer() []uint16 {
	// TODO: Implement path triangulation indices
	return []uint16{}
}

// Clone creates a copy of this geometry
func (p *PathGeometry) Clone() Geometry {
	return &PathGeometry{
		BaseGeometry: p.BaseGeometry.Clone().(BaseGeometry),
		path:         p.path,
		fillRule:     p.fillRule,
		transform:    p.transform,
	}
}

// StrokePathGeometry represents stroked path geometry
type StrokePathGeometry struct {
	BaseGeometry
	path       geom.Path
	width      float32
	cap        LineCap
	join       LineJoin
	miterLimit float32
	transform  geom.Matrix
}

// LineJoin defines how line segments are joined
type LineJoin int

const (
	// LineJoinMiter miter join
	LineJoinMiter LineJoin = iota
	// LineJoinRound round join
	LineJoinRound
	// LineJoinBevel bevel join
	LineJoinBevel
)

// NewStrokePathGeometry creates a new stroked path geometry
func NewStrokePathGeometry(path geom.Path, width float32) *StrokePathGeometry {
	return &StrokePathGeometry{
		BaseGeometry: NewBaseGeometry(),
		path:         path,
		width:        width,
		cap:          LineCapButt,
		join:         LineJoinMiter,
		miterLimit:   4.0,
		transform:    geom.NewIdentityMatrix(),
	}
}

// SetPath sets the path
func (s *StrokePathGeometry) SetPath(path geom.Path) {
	s.path = path
}

// GetPath returns the current path
func (s *StrokePathGeometry) GetPath() geom.Path {
	return s.path
}

// SetWidth sets the stroke width
func (s *StrokePathGeometry) SetWidth(width float32) {
	s.width = width
}

// GetWidth returns the stroke width
func (s *StrokePathGeometry) GetWidth() float32 {
	return s.width
}

// SetCap sets the line cap style
func (s *StrokePathGeometry) SetCap(cap LineCap) {
	s.cap = cap
}

// GetCap returns the line cap style
func (s *StrokePathGeometry) GetCap() LineCap {
	return s.cap
}

// SetJoin sets the line join style
func (s *StrokePathGeometry) SetJoin(join LineJoin) {
	s.join = join
}

// GetJoin returns the line join style
func (s *StrokePathGeometry) GetJoin() LineJoin {
	return s.join
}

// SetMiterLimit sets the miter limit
func (s *StrokePathGeometry) SetMiterLimit(limit float32) {
	s.miterLimit = limit
}

// GetMiterLimit returns the miter limit
func (s *StrokePathGeometry) GetMiterLimit() float32 {
	return s.miterLimit
}

// SetTransform sets the transformation matrix
func (s *StrokePathGeometry) SetTransform(transform geom.Matrix) {
	s.transform = transform
}

// GetTransform returns the current transformation matrix
func (s *StrokePathGeometry) GetTransform() geom.Matrix {
	return s.transform
}

// GetBounds returns the bounds of this geometry
func (s *StrokePathGeometry) GetBounds() geom.Rect {
	bounds := s.path.GetBounds()

	// Expand by stroke width
	halfWidth := s.width / 2
	expandedBounds := geom.Rect{
		X:      bounds.X - halfWidth,
		Y:      bounds.Y - halfWidth,
		Width:  bounds.Width + s.width,
		Height: bounds.Height + s.width,
	}

	return s.transform.TransformRect(expandedBounds)
}

// GetVertexCount returns the number of vertices (estimated)
func (s *StrokePathGeometry) GetVertexCount() int {
	// TODO: Calculate based on path complexity and stroke parameters
	return 0
}

// GetIndexCount returns the number of indices (estimated)
func (s *StrokePathGeometry) GetIndexCount() int {
	// TODO: Calculate based on stroke tessellation
	return 0
}

// ComputeVertices generates vertex data for rendering
func (s *StrokePathGeometry) ComputeVertices() []render.Vertex {
	// TODO: Implement path stroking
	// This would involve:
	// 1. Generating stroke outline
	// 2. Handling line caps and joins
	// 3. Tessellating the stroke
	return []render.Vertex{}
}

// GetIndexBuffer generates index data
func (s *StrokePathGeometry) GetIndexBuffer() []uint16 {
	// TODO: Implement stroke triangulation indices
	return []uint16{}
}

// Clone creates a copy of this geometry
func (s *StrokePathGeometry) Clone() Geometry {
	return &StrokePathGeometry{
		BaseGeometry: s.BaseGeometry.Clone().(BaseGeometry),
		path:         s.path,
		width:        s.width,
		cap:          s.cap,
		join:         s.join,
		miterLimit:   s.miterLimit,
		transform:    s.transform,
	}
}

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
