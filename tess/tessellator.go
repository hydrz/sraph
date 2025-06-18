package tess

import (
	"math"

	"github.com/opensraph/sraph/geom"
)

// PointArenaSize is the size of the point arena buffer stored on the tessellator.
const PointArenaSize = 4096

// Result represents the result of tessellation operations.
type Result uint8

const (
	// ResultSuccess indicates successful tessellation.
	ResultSuccess Result = iota
	// ResultInputError indicates an input error.
	ResultInputError
	// ResultTessellationError indicates a tessellation error.
	ResultTessellationError
)

// TessellatedVertexProc is a callback function for a VertexGenerator to deliver
// the vertices it computes as Point objects.
type TessellatedVertexProc[T geom.Scalar] func(p geom.Point[T])

// VertexGenerator is an object which produces a list of vertices as Points that
// tessellate a previously provided shape and delivers the vertices
// through a TessellatedVertexProc callback.
//
// The object can also provide advance information on how many
// vertices it will generate.
type VertexGenerator[T geom.Scalar] interface {
	// GetTriangleType returns the PrimitiveType that describes the relationship
	// among the list of vertices produced by the GenerateVertices method.
	//
	// Most generators will deliver TriangleStrip triangles
	GetTriangleType() PrimitiveType

	// GetVertexCount returns the number of vertices that the generator plans to
	// produce, if known.
	//
	// This value is advisory only and can be used to reserve space
	// where the vertices will be placed, but the count may be an
	// estimate.
	GetVertexCount() int

	// GenerateVertices generates the vertices and delivers them in the necessary
	// order (as required by the PrimitiveType) to the given callback function.
	GenerateVertices(proc TessellatedVertexProc[T])
}

// PrimitiveType defines the type of primitive used for rendering.
type PrimitiveType uint8

const (
	// PrimitiveTypeTriangleStrip represents triangle strip primitives.
	PrimitiveTypeTriangleStrip PrimitiveType = iota
	// PrimitiveTypeTriangleFan represents triangle fan primitives.
	PrimitiveTypeTriangleFan
)

// Trigs is essentially just a slice of Trig objects, but supports storing a
// reference to either a cached slice or a locally generated slice.
// The constructor will fill the slice with quarter circular samples
// for the indicated number of equal divisions if the slice is new.
//
// A given instance of Trigs will always contain at least 2 entries
// which is the minimum number of samples to traverse a quarter circle
// in a single step. The first sample will always be (0, 1) and the last
// sample will always be (1, 0).
type Trigs[T geom.Scalar] struct {
	trigs []geom.Trig[T]
}

// NewTrigs creates a new Trigs with the specified number of divisions.
func NewTrigs[T geom.Scalar](divisions int) *Trigs[T] {
	if divisions < 1 {
		divisions = 1
	}

	t := &Trigs[T]{
		trigs: make([]geom.Trig[T], divisions+1),
	}
	t.init(divisions)
	return t
}

// NewTrigsForPixelRadius creates a new Trigs for the given pixel radius.
func NewTrigsForPixelRadius[T geom.Scalar](pixelRadius T) *Trigs[T] {
	// Calculate the number of divisions based on pixel radius
	// This follows the same logic as the C++ implementation
	divisions := int(math.Ceil(pixelRadius.Float64() * math.Pi / 2.0))
	if divisions < 1 {
		divisions = 1
	}
	return NewTrigs[T](divisions)
}

// Size returns the number of trig values.
func (t *Trigs[T]) Size() int {
	return len(t.trigs)
}

// Get returns the trig value at the specified index.
func (t *Trigs[T]) Get(index int) geom.Trig[T] {
	return t.trigs[index]
}

// Slice returns the underlying slice of trig values.
func (t *Trigs[T]) Slice() []geom.Trig[T] {
	return t.trigs
}

// init fills the slice with the indicated number of equal divisions of
// trigonometric values.
func (t *Trigs[T]) init(divisions int) {
	// Fill with quarter circle samples
	for i := 0; i <= divisions; i++ {
		angle := geom.Radians(float64(i) * math.Pi / 2.0 / float64(divisions))
		t.trigs[i] = geom.NewTrig[T](angle)
	}
}

// ArcIteration describes the iteration through a set of angle vectors
// in a Trigs structure to render the points along an arc.
type ArcIteration[T geom.Scalar] struct {
	// Start is the true begin angle of the arc, expressed as unit direction vector.
	Start geom.Point[T]
	// End is the true end angle of the arc, expressed as unit direction vector.
	End geom.Point[T]
	// Quadrants is the variable number of quadrants that have to be iterated.
	Quadrants []Quadrant[T]
}

// Quadrant represents the axis to multiply by each Trig value and the half-open [start, end)
// range of the Trig vector over which to compute.
type Quadrant[T geom.Scalar] struct {
	Axis       geom.Point[T]
	StartIndex int
	EndIndex   int
}

// GetPointCount returns the number of points in this quadrant.
func (q *Quadrant[T]) GetPointCount() int {
	if q.StartIndex >= q.EndIndex {
		return 0
	}
	return q.EndIndex - q.StartIndex
}

// GetPointCount returns the total number of points in the arc iteration.
func (a *ArcIteration[T]) GetPointCount() int {
	count := 2 // start and end points
	for _, quadrant := range a.Quadrants {
		count += quadrant.GetPointCount()
	}
	return count
}

// EllipticalVertexGeneratorData holds the data needed for elliptical vertex generation.
type EllipticalVertexGeneratorData[T geom.Scalar] struct {
	// ReferenceCenters - Circles and Ellipses only use one of these points.
	// RoundCapLines use both as the endpoints of the unexpanded line.
	// A round rect can specify its interior rectangle by using the
	// 2 points as opposing corners.
	ReferenceCenters [2]geom.Point[T]
	// Radii - Circular shapes have the same value in radii.Width and radii.Height
	Radii geom.Size[T]
	// HalfWidth is only used in cases where the generator will be
	// generating 2 different outlines, such as StrokedCircle
	HalfWidth T
}

// EllipticalVertexGenerator is the VertexGenerator implementation common to all shapes
// that are based on a polygonal representation of an ellipse.
type EllipticalVertexGenerator[T geom.Scalar] struct {
	trigs           *Trigs[T]
	data            EllipticalVertexGeneratorData[T]
	verticesPerTrig int
	generator       func(*Trigs[T], EllipticalVertexGeneratorData[T], TessellatedVertexProc[T])
}

// GetTriangleType implements VertexGenerator.
func (e *EllipticalVertexGenerator[T]) GetTriangleType() PrimitiveType {
	return PrimitiveTypeTriangleStrip
}

// GetVertexCount implements VertexGenerator.
func (e *EllipticalVertexGenerator[T]) GetVertexCount() int {
	return e.trigs.Size() * e.verticesPerTrig
}

// GenerateVertices implements VertexGenerator.
func (e *EllipticalVertexGenerator[T]) GenerateVertices(proc TessellatedVertexProc[T]) {
	e.generator(e.trigs, e.data, proc)
}

// ArcVertexGenerator is the VertexGenerator implementation common to all shapes
// that are based on a polygonal representation of an arc.
type ArcVertexGenerator[T geom.Scalar] struct {
	iteration            ArcIteration[T]
	trigs                *Trigs[T]
	ovalBounds           geom.Rect[T]
	useCenter            bool
	halfWidth            T
	cap                  geom.StrokeCap
	supportsTriangleFans bool
}

// GetTriangleType implements VertexGenerator.
func (a *ArcVertexGenerator[T]) GetTriangleType() PrimitiveType {
	if a.supportsTriangleFans && a.useCenter {
		return PrimitiveTypeTriangleFan
	}
	return PrimitiveTypeTriangleStrip
}

// GetVertexCount implements VertexGenerator.
func (a *ArcVertexGenerator[T]) GetVertexCount() int {
	count := a.iteration.GetPointCount()
	if a.halfWidth > 0 {
		count *= 2 // For stroked arcs, we need inner and outer vertices
	}
	if a.useCenter {
		count++ // Add center point for filled arcs
	}
	return count
}

// GenerateVertices implements VertexGenerator.
func (a *ArcVertexGenerator[T]) GenerateVertices(proc TessellatedVertexProc[T]) {
	if a.halfWidth > 0 {
		generateStrokedArc(a.trigs, &a.iteration, a.ovalBounds, a.halfWidth, a.cap, proc)
	} else if a.supportsTriangleFans {
		generateFilledArcFan(a.trigs, &a.iteration, a.ovalBounds, a.useCenter, proc)
	} else {
		generateFilledArcStrip(a.trigs, &a.iteration, a.ovalBounds, a.useCenter, proc)
	}
}

// Tessellator is a utility that generates triangles of the specified fill type
// given a polyline. This happens on the CPU.
//
// Also contains functionality for optimized generation of circles and ellipses.
//
// This object is not thread safe, and its methods must not be called from multiple threads.
type Tessellator[T geom.Scalar] struct {
	// Data for various Circle/EllipseGenerator classes, cached per
	// Tessellator instance which is usually the foreground life of an app
	// if not longer.
	precomputedTrigs [300]*Trigs[T]

	// Used for polyline generation
	pointBuffer []geom.Point[T]
	indexBuffer []uint16

	// Used for stroke path generation
	strokePoints []geom.Point[T]
}

// CircleTolerance is the pixel tolerance used by the algorithm to determine how
// many divisions to create for a circle.
//
// No point on the polygon of vertices should deviate from the
// true circle by more than this tolerance.
const CircleTolerance = 0.1

// NewTessellator creates a new Tessellator instance.
func NewTessellator[T geom.Scalar]() *Tessellator[T] {
	return &Tessellator[T]{
		pointBuffer:  make([]geom.Point[T], 0, PointArenaSize),
		indexBuffer:  make([]uint16, 0),
		strokePoints: make([]geom.Point[T], 0, PointArenaSize),
	}
}

// GetTrigsForDeviceRadius returns a vector of Trig (cos, sin pairs) structs for a 90 degree
// circle quadrant of the specified pixel radius.
func (t *Tessellator[T]) GetTrigsForDeviceRadius(pixelRadius T) *Trigs[T] {
	divisions := int(math.Ceil(pixelRadius.Float64() * math.Pi / 2.0 / CircleTolerance))
	if divisions < 1 {
		divisions = 1
	}
	if divisions >= len(t.precomputedTrigs) {
		return NewTrigs[T](divisions)
	}

	// Use cached trig values if available
	if t.precomputedTrigs[divisions] == nil {
		t.precomputedTrigs[divisions] = NewTrigs[T](divisions)
	}
	return t.precomputedTrigs[divisions]
}

// GetStrokePointCache retrieves a pre-allocated arena of PointArenaSize points.
func (t *Tessellator[T]) GetStrokePointCache() []geom.Point[T] {
	t.strokePoints = t.strokePoints[:0] // Reset slice but keep capacity
	return t.strokePoints
}

// FilledCircle creates a VertexGenerator that can produce vertices for
// a filled circle of the given radius around the given center.
func (t *Tessellator[T]) FilledCircle(center geom.Point[T], radius T, pixelRadius T) *EllipticalVertexGenerator[T] {
	trigs := t.GetTrigsForDeviceRadius(pixelRadius)

	data := EllipticalVertexGeneratorData[T]{
		ReferenceCenters: [2]geom.Point[T]{center, center},
		Radii:            geom.Size[T]{Width: radius, Height: radius},
		HalfWidth:        0,
	}

	return &EllipticalVertexGenerator[T]{
		trigs:           trigs,
		data:            data,
		verticesPerTrig: 1,
		generator:       generateFilledCircle[T],
	}
}

// StrokedCircle creates a VertexGenerator that can produce vertices for
// a stroked circle of the given radius and half_width around the given center.
func (t *Tessellator[T]) StrokedCircle(center geom.Point[T], radius, halfWidth T, pixelRadius T) *EllipticalVertexGenerator[T] {
	trigs := t.GetTrigsForDeviceRadius(pixelRadius)

	data := EllipticalVertexGeneratorData[T]{
		ReferenceCenters: [2]geom.Point[T]{center, center},
		Radii:            geom.Size[T]{Width: radius, Height: radius},
		HalfWidth:        halfWidth,
	}

	return &EllipticalVertexGenerator[T]{
		trigs:           trigs,
		data:            data,
		verticesPerTrig: 2, // Inner and outer vertices
		generator:       generateStrokedCircle[T],
	}
}

// FilledEllipse creates a VertexGenerator that can produce vertices for
// a filled ellipse inscribed within the given bounds.
func (t *Tessellator[T]) FilledEllipse(bounds geom.Rect[T], pixelRadius T) *EllipticalVertexGenerator[T] {
	trigs := t.GetTrigsForDeviceRadius(pixelRadius)
	center := bounds.Center()
	radii := geom.Size[T]{
		Width:  bounds.Width() / 2,
		Height: bounds.Height() / 2,
	}

	data := EllipticalVertexGeneratorData[T]{
		ReferenceCenters: [2]geom.Point[T]{center, center},
		Radii:            radii,
		HalfWidth:        0,
	}

	return &EllipticalVertexGenerator[T]{
		trigs:           trigs,
		data:            data,
		verticesPerTrig: 1,
		generator:       generateFilledEllipse[T],
	}
}

// RoundCapLine creates a VertexGenerator that can produce vertices for
// a line with round end caps of the given radius.
func (t *Tessellator[T]) RoundCapLine(p0, p1 geom.Point[T], radius T, pixelRadius T) *EllipticalVertexGenerator[T] {
	trigs := t.GetTrigsForDeviceRadius(pixelRadius)

	data := EllipticalVertexGeneratorData[T]{
		ReferenceCenters: [2]geom.Point[T]{p0, p1},
		Radii:            geom.Size[T]{Width: radius, Height: radius},
		HalfWidth:        0,
	}

	return &EllipticalVertexGenerator[T]{
		trigs:           trigs,
		data:            data,
		verticesPerTrig: 2,
		generator:       generateRoundCapLine[T],
	}
}

// ComputeArcQuadrantIterations computes the quadrant iterations for an arc
// with the given start angle and sweep.
func ComputeArcQuadrantIterations[T geom.Scalar](trigCount int, start, sweep geom.Degrees) ArcIteration[T] {
	// Convert degrees to radians
	startRad := start.Radians()
	sweepRad := sweep.Radians()
	endRad := startRad + sweepRad

	// Create start and end unit vectors
	startVec := geom.Point[T]{
		X: T(math.Cos(startRad.Float64())),
		Y: T(math.Sin(startRad.Float64())),
	}
	endVec := geom.Point[T]{
		X: T(math.Cos(endRad.Float64())),
		Y: T(math.Sin(endRad.Float64())),
	}

	iteration := ArcIteration[T]{
		Start: startVec,
		End:   endVec,
	}

	// Compute quadrants based on the arc span
	// This is a simplified implementation - the full C++ version is more complex
	if math.Abs(sweepRad.Float64()) >= math.Pi/2 {
		// Arc spans multiple quadrants
		quadrant := Quadrant[T]{
			Axis:       geom.Point[T]{X: 1, Y: 0}, // X-axis
			StartIndex: 0,
			EndIndex:   trigCount,
		}
		iteration.Quadrants = append(iteration.Quadrants, quadrant)
	} else {
		// Arc is within a single quadrant
		quadrant := Quadrant[T]{
			Axis:       geom.Point[T]{X: 1, Y: 0}, // X-axis
			StartIndex: 0,
			EndIndex:   int(math.Abs(sweepRad.Float64()) * float64(trigCount) / (math.Pi / 2)),
		}
		if quadrant.EndIndex > trigCount {
			quadrant.EndIndex = trigCount
		}
		iteration.Quadrants = append(iteration.Quadrants, quadrant)
	}

	return iteration
}
