package geometry

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// StrokePathGeometry represents geometry for stroked path rendering
type StrokePathGeometry struct {
	BaseGeometry
	path         geom.Path
	strokeParams StrokeParameters
	transform    geom.Matrix
	tessellated  bool
	vertices     []render.Vertex
	indices      []uint16
}

// StrokeParameters defines stroke rendering parameters
type StrokeParameters struct {
	Width      float32
	Cap        LineCap
	Join       LineJoin
	MiterLimit float32
}

// NewStrokePathGeometry creates a new stroke path geometry
func NewStrokePathGeometry(path geom.Path, params StrokeParameters) *StrokePathGeometry {
	return &StrokePathGeometry{
		BaseGeometry: NewBaseGeometry(),
		path:         path,
		strokeParams: params,
		transform:    geom.NewIdentityMatrix(),
		tessellated:  false,
	}
}

// SetPath sets the path to be stroked
func (s *StrokePathGeometry) SetPath(path geom.Path) {
	s.path = path
	s.tessellated = false
}

// GetPath returns the current path
func (s *StrokePathGeometry) GetPath() geom.Path {
	return s.path
}

// SetStrokeParameters sets the stroke parameters
func (s *StrokePathGeometry) SetStrokeParameters(params StrokeParameters) {
	s.strokeParams = params
	s.tessellated = false
}

// GetStrokeParameters returns the stroke parameters
func (s *StrokePathGeometry) GetStrokeParameters() StrokeParameters {
	return s.strokeParams
}

// SetTransform sets the transformation matrix
func (s *StrokePathGeometry) SetTransform(transform geom.Matrix) {
	s.transform = transform
}

// GetTransform returns the transformation matrix
func (s *StrokePathGeometry) GetTransform() geom.Matrix {
	return s.transform
}

// GetBounds returns the bounds of this geometry
func (s *StrokePathGeometry) GetBounds() geom.Rect {
	bounds := s.path.GetBounds()

	// Expand by stroke width
	halfWidth := s.strokeParams.Width / 2
	expandedBounds := geom.Rect{
		X:      bounds.X - halfWidth,
		Y:      bounds.Y - halfWidth,
		Width:  bounds.Width + s.strokeParams.Width,
		Height: bounds.Height + s.strokeParams.Width,
	}

	return s.transform.TransformRect(expandedBounds)
}

// Tessellate performs stroke tessellation
func (s *StrokePathGeometry) Tessellate() error {
	// TODO: Implement stroke tessellation
	// This would involve:
	// 1. Generating stroke outline from path
	// 2. Handling line caps and joins
	// 3. Creating triangulated mesh

	// Placeholder implementation
	s.vertices = []render.Vertex{}
	s.indices = []uint16{}
	s.tessellated = true

	return nil
}

// GetVertexCount returns the number of vertices
func (s *StrokePathGeometry) GetVertexCount() int {
	if !s.tessellated {
		s.Tessellate()
	}
	return len(s.vertices)
}

// GetIndexCount returns the number of indices
func (s *StrokePathGeometry) GetIndexCount() int {
	if !s.tessellated {
		s.Tessellate()
	}
	return len(s.indices)
}

// ComputeVertices generates vertex data for rendering
func (s *StrokePathGeometry) ComputeVertices() []render.Vertex {
	if !s.tessellated {
		s.Tessellate()
	}
	return s.vertices
}

// GetIndexBuffer generates index data
func (s *StrokePathGeometry) GetIndexBuffer() []uint16 {
	if !s.tessellated {
		s.Tessellate()
	}
	return s.indices
}

// Clone creates a copy of this geometry
func (s *StrokePathGeometry) Clone() Geometry {
	return &StrokePathGeometry{
		BaseGeometry: s.BaseGeometry.Clone().(BaseGeometry),
		path:         s.path,
		strokeParams: s.strokeParams,
		transform:    s.transform,
		tessellated:  false,
	}
}

// GetCoverage returns the coverage area for this geometry
func (s *StrokePathGeometry) GetCoverage(transform geom.Matrix) geom.Rect {
	return s.GetBounds()
}
