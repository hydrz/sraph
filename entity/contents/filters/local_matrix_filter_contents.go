package filters

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// LocalMatrixFilterContents represents a filter that applies local matrix transformations
type LocalMatrixFilterContents struct {
	FilterContentsBase
	LocalMatrix geom.Matrix
	InputFilter FilterContents
}

// NewLocalMatrixFilterContents creates a new local matrix filter
func NewLocalMatrixFilterContents(matrix geom.Matrix, input FilterContents) *LocalMatrixFilterContents {
	return &LocalMatrixFilterContents{
		LocalMatrix: matrix,
		InputFilter: input,
	}
}

// Render renders the local matrix filter
func (f *LocalMatrixFilterContents) Render(renderer render.ContentRenderer, entity render.Entity) bool {
	// TODO: Implement local matrix transformation rendering
	// This involves applying a local coordinate transformation
	return false
}

// GetCoverage returns the coverage area of the filter
func (f *LocalMatrixFilterContents) GetCoverage(entity render.Entity) geom.Rect {
	// TODO: Transform the input coverage by the local matrix
	if f.InputFilter != nil {
		inputCoverage := f.InputFilter.GetCoverage(entity)
		return f.LocalMatrix.TransformRect(inputCoverage)
	}
	return entity.GetBounds()
}

// Clone creates a copy of the filter
func (f *LocalMatrixFilterContents) Clone() FilterContents {
	var clonedInput FilterContents
	if f.InputFilter != nil {
		clonedInput = f.InputFilter.Clone()
	}
	return &LocalMatrixFilterContents{
		LocalMatrix: f.LocalMatrix,
		InputFilter: clonedInput,
	}
}

// SetLocalMatrix sets the local transformation matrix
func (f *LocalMatrixFilterContents) SetLocalMatrix(matrix geom.Matrix) {
	f.LocalMatrix = matrix
}

// GetLocalMatrix returns the local transformation matrix
func (f *LocalMatrixFilterContents) GetLocalMatrix() geom.Matrix {
	return f.LocalMatrix
}

// SetInputFilter sets the input filter to be transformed
func (f *LocalMatrixFilterContents) SetInputFilter(input FilterContents) {
	f.InputFilter = input
}

// GetInputFilter returns the input filter
func (f *LocalMatrixFilterContents) GetInputFilter() FilterContents {
	return f.InputFilter
}
