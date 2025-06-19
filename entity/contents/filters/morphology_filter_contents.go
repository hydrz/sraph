package filters

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// MorphologyOperation defines the type of morphology operation
type MorphologyOperation int

const (
	MorphologyDilate MorphologyOperation = iota
	MorphologyErode
)

// MorphologyFilterContents represents a filter that applies morphological operations
type MorphologyFilterContents struct {
	FilterContentsBase
	Operation MorphologyOperation
	RadiusX   float32
	RadiusY   float32
}

// NewMorphologyFilterContents creates a new morphology filter
func NewMorphologyFilterContents(operation MorphologyOperation, radiusX, radiusY float32) *MorphologyFilterContents {
	return &MorphologyFilterContents{
		Operation: operation,
		RadiusX:   radiusX,
		RadiusY:   radiusY,
	}
}

// Render renders the morphology filter
func (f *MorphologyFilterContents) Render(renderer render.ContentRenderer, entity render.Entity) bool {
	// TODO: Implement morphological operations (dilate/erode)
	// This involves expanding or contracting shapes
	return false
}

// GetCoverage returns the coverage area of the filter
func (f *MorphologyFilterContents) GetCoverage(entity render.Entity) geom.Rect {
	bounds := entity.GetBounds()
	// Morphology operations can expand the coverage
	expansion := f.getMaxRadius()
	return geom.Rect{
		Origin: geom.Point{
			X: bounds.Origin.X - expansion,
			Y: bounds.Origin.Y - expansion,
		},
		Size: geom.Size{
			Width:  bounds.Size.Width + 2*expansion,
			Height: bounds.Size.Height + 2*expansion,
		},
	}
}

// Clone creates a copy of the filter
func (f *MorphologyFilterContents) Clone() FilterContents {
	return &MorphologyFilterContents{
		Operation: f.Operation,
		RadiusX:   f.RadiusX,
		RadiusY:   f.RadiusY,
	}
}

// SetOperation sets the morphology operation
func (f *MorphologyFilterContents) SetOperation(operation MorphologyOperation) {
	f.Operation = operation
}

// GetOperation returns the morphology operation
func (f *MorphologyFilterContents) GetOperation() MorphologyOperation {
	return f.Operation
}

// SetRadiusX sets the X radius for the operation
func (f *MorphologyFilterContents) SetRadiusX(radius float32) {
	f.RadiusX = radius
}

// GetRadiusX returns the X radius
func (f *MorphologyFilterContents) GetRadiusX() float32 {
	return f.RadiusX
}

// SetRadiusY sets the Y radius for the operation
func (f *MorphologyFilterContents) SetRadiusY(radius float32) {
	f.RadiusY = radius
}

// GetRadiusY returns the Y radius
func (f *MorphologyFilterContents) GetRadiusY() float32 {
	return f.RadiusY
}

// getMaxRadius returns the maximum radius for coverage calculation
func (f *MorphologyFilterContents) getMaxRadius() float32 {
	if f.RadiusX > f.RadiusY {
		return f.RadiusX
	}
	return f.RadiusY
}
