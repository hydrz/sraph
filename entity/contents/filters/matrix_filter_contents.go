package filters

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// MatrixFilterContents represents a filter that applies matrix transformations
type MatrixFilterContents struct {
	FilterContentsBase
	Matrix      geom.Matrix
	InputFilter FilterContents
}

// NewMatrixFilterContents creates a new matrix filter
func NewMatrixFilterContents(matrix geom.Matrix, input FilterContents) *MatrixFilterContents {
	return &MatrixFilterContents{
		Matrix:      matrix,
		InputFilter: input,
	}
}

// Render renders the matrix filter
func (f *MatrixFilterContents) Render(renderer render.ContentRenderer, entity render.Entity) bool {
	// TODO: Implement matrix transformation rendering
	// This involves applying geometric transformations
	return false
}

// GetCoverage returns the coverage area of the filter
func (f *MatrixFilterContents) GetCoverage(entity render.Entity) geom.Rect {
	// TODO: Transform the input coverage by the matrix
	if f.InputFilter != nil {
		inputCoverage := f.InputFilter.GetCoverage(entity)
		return f.Matrix.TransformRect(inputCoverage)
	}
	return f.Matrix.TransformRect(entity.GetBounds())
}

// Clone creates a copy of the filter
func (f *MatrixFilterContents) Clone() FilterContents {
	var clonedInput FilterContents
	if f.InputFilter != nil {
		clonedInput = f.InputFilter.Clone()
	}
	return &MatrixFilterContents{
		Matrix:      f.Matrix,
		InputFilter: clonedInput,
	}
}

// SetMatrix sets the transformation matrix
func (f *MatrixFilterContents) SetMatrix(matrix geom.Matrix) {
	f.Matrix = matrix
}

// GetMatrix returns the transformation matrix
func (f *MatrixFilterContents) GetMatrix() geom.Matrix {
	return f.Matrix
}

// SetInputFilter sets the input filter to be transformed
func (f *MatrixFilterContents) SetInputFilter(input FilterContents) {
	f.InputFilter = input
}

// GetInputFilter returns the input filter
func (f *MatrixFilterContents) GetInputFilter() FilterContents {
	return f.InputFilter
}
