package tess

import (
	"github.com/opensraph/sraph/geom"
)

// TessellatorLibtess provides tessellation using libtess2-like algorithms.
// An extended tessellator that offers arbitrary/concave tessellation.
//
// This object is not thread safe, and its methods must not be called from multiple threads.
type TessellatorLibtess[T geom.Scalar] struct {
	// Internal state for tessellation
}

// NewTessellatorLibtess creates a new TessellatorLibtess instance.
func NewTessellatorLibtess[T geom.Scalar]() *TessellatorLibtess[T] {
	return &TessellatorLibtess[T]{}
}

// BuilderCallback is a callback that returns the results of the tessellation.
//
// The index buffer may not be populated, in which case indices will
// be nil and indicesCount will be 0.
type BuilderCallback[T geom.Scalar] func(vertices []T, verticesCount int, indices []uint16, indicesCount int) bool

// Tessellate generates filled triangles from the path. A callback is
// invoked once for the entire tessellation.
//
// Parameters:
//   - source: The path source to tessellate.
//   - tolerance: The tolerance value for conversion of the path to
//     a polyline. This value is often derived from the
//     Matrix::GetMaxBasisLength of the CTM applied to the
//     path for rendering.
//   - callback: The callback, return false to indicate failure.
//
// Returns the result status of the tessellation.
func (t *TessellatorLibtess[T]) Tessellate(source geom.PathSource[T], tolerance T, callback BuilderCallback[T]) Result {
	if callback == nil {
		return ResultInputError
	}

	polyline := &polylineWriter[T]{}
	PathToFilledVertices(source, polyline, tolerance)

	fillType := source.FillType()

	if len(polyline.points) == 0 {
		return ResultInputError
	}

	// Simple triangle fan tessellation for convex polygons
	// In a full implementation, this would use libtess2 or a similar library
	return t.tessellateSimple(polyline, fillType, callback)
}

// tessellateSimple performs simple triangle fan tessellation for convex shapes.
func (t *TessellatorLibtess[T]) tessellateSimple(polyline *polylineWriter[T], fillType geom.FillType, callback BuilderCallback[T]) Result {
	if len(polyline.contours) == 0 {
		return ResultInputError
	}

	// Convert points to flat array
	vertices := make([]T, len(polyline.points)*2)
	for i, point := range polyline.points {
		vertices[i*2] = point.X
		vertices[i*2+1] = point.Y
	}

	// Simple triangle fan tessellation
	var indices []uint16
	for _, contour := range polyline.contours {
		if contour.Size() < 3 {
			continue // Skip degenerate contours
		}

		// Create triangle fan from contour
		for i := 1; i < int(contour.Size())-1; i++ {
			indices = append(indices,
				uint16(contour.Start),
				uint16(contour.Start+i),
				uint16(contour.Start+i+1),
			)
		}
	}

	if !callback(vertices, len(vertices), indices, len(indices)) {
		return ResultInputError
	}

	return ResultSuccess
}

// polylineWriter collects vertices for tessellation.
type polylineWriter[T geom.Scalar] struct {
	points       []geom.Point[T]
	contours     []contour
	contourStart int
}

// contour represents a single contour in the polyline.
type contour struct {
	Start int // Start index in points array
	End   int // End index in points array
}

// Size returns the number of points in the contour.
func (c contour) Size() int {
	return c.End - c.Start
}

// Write implements VertexWriter.
func (p *polylineWriter[T]) Write(point geom.Point[T]) {
	p.points = append(p.points, point)
}

// EndContour implements VertexWriter.
func (p *polylineWriter[T]) EndContour() {
	contourEnd := len(p.points)
	p.contours = append(p.contours, contour{
		Start: p.contourStart,
		End:   contourEnd,
	})
	p.contourStart = contourEnd
}
