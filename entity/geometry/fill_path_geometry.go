package geometry

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// FillPathGeometry represents geometry for filled path rendering
type FillPathGeometry struct {
	BaseGeometry
	path        geom.Path
	fillRule    FillRule
	innerRect   *geom.Rect // Optional inner rectangle for optimization
	transform   geom.Matrix
	tessellated bool
	vertices    []render.Vertex
	indices     []uint16
}

// NewFillPathGeometry creates a new fill path geometry
func NewFillPathGeometry(path geom.Path, fillRule FillRule) *FillPathGeometry {
	return &FillPathGeometry{
		BaseGeometry: NewBaseGeometry(),
		path:         path,
		fillRule:     fillRule,
		innerRect:    nil,
		transform:    geom.NewIdentityMatrix(),
		tessellated:  false,
	}
}

// NewFillPathGeometryWithInnerRect creates a fill path geometry with inner rectangle optimization
func NewFillPathGeometryWithInnerRect(path geom.Path, fillRule FillRule, innerRect geom.Rect) *FillPathGeometry {
	return &FillPathGeometry{
		BaseGeometry: NewBaseGeometry(),
		path:         path,
		fillRule:     fillRule,
		innerRect:    &innerRect,
		transform:    geom.NewIdentityMatrix(),
		tessellated:  false,
	}
}

// SetPath sets the path to be filled
func (f *FillPathGeometry) SetPath(path geom.Path) {
	f.path = path
	f.tessellated = false // Mark as needing re-tessellation
}

// GetPath returns the current path
func (f *FillPathGeometry) GetPath() geom.Path {
	return f.path
}

// SetFillRule sets the fill rule
func (f *FillPathGeometry) SetFillRule(rule FillRule) {
	f.fillRule = rule
	f.tessellated = false
}

// GetFillRule returns the current fill rule
func (f *FillPathGeometry) GetFillRule() FillRule {
	return f.fillRule
}

// SetInnerRect sets the inner rectangle for optimization
func (f *FillPathGeometry) SetInnerRect(rect *geom.Rect) {
	f.innerRect = rect
}

// GetInnerRect returns the inner rectangle
func (f *FillPathGeometry) GetInnerRect() *geom.Rect {
	return f.innerRect
}

// SetTransform sets the transformation matrix
func (f *FillPathGeometry) SetTransform(transform geom.Matrix) {
	f.transform = transform
}

// GetTransform returns the transformation matrix
func (f *FillPathGeometry) GetTransform() geom.Matrix {
	return f.transform
}

// GetBounds returns the bounds of this geometry
func (f *FillPathGeometry) GetBounds() geom.Rect {
	bounds := f.path.GetBounds()
	return f.transform.TransformRect(bounds)
}

// IsTessellated returns true if the path has been tessellated
func (f *FillPathGeometry) IsTessellated() bool {
	return f.tessellated
}

// Tessellate performs path tessellation
func (f *FillPathGeometry) Tessellate() error {
	// TODO: Implement path tessellation
	// This would involve:
	// 1. Converting path commands to triangles
	// 2. Handling fill rule (non-zero vs even-odd)
	// 3. Optimizing with inner rectangle if available

	// Placeholder implementation
	f.vertices = []render.Vertex{}
	f.indices = []uint16{}
	f.tessellated = true

	return nil
}

// GetVertexCount returns the number of vertices
func (f *FillPathGeometry) GetVertexCount() int {
	if !f.tessellated {
		f.Tessellate()
	}
	return len(f.vertices)
}

// GetIndexCount returns the number of indices
func (f *FillPathGeometry) GetIndexCount() int {
	if !f.tessellated {
		f.Tessellate()
	}
	return len(f.indices)
}

// ComputeVertices generates vertex data for rendering
func (f *FillPathGeometry) ComputeVertices() []render.Vertex {
	if !f.tessellated {
		f.Tessellate()
	}
	return f.vertices
}

// GetIndexBuffer generates index data
func (f *FillPathGeometry) GetIndexBuffer() []uint16 {
	if !f.tessellated {
		f.Tessellate()
	}
	return f.indices
}

// Clone creates a copy of this geometry
func (f *FillPathGeometry) Clone() Geometry {
	clone := &FillPathGeometry{
		BaseGeometry: f.BaseGeometry.Clone().(BaseGeometry),
		path:         f.path,
		fillRule:     f.fillRule,
		transform:    f.transform,
		tessellated:  false, // Force re-tessellation
	}

	if f.innerRect != nil {
		innerRect := *f.innerRect
		clone.innerRect = &innerRect
	}

	return clone
}

// GetCoverage returns the coverage area for this geometry
func (f *FillPathGeometry) GetCoverage(transform geom.Matrix) geom.Rect {
	combinedTransform := transform.Multiply(f.transform)
	bounds := f.path.GetBounds()
	return combinedTransform.TransformRect(bounds)
}

// CanApplyMaskFilter returns whether mask filters can be applied
func (f *FillPathGeometry) CanApplyMaskFilter() bool {
	return true
}
