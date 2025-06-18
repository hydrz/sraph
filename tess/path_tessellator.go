package tess

import (
	"math"

	"github.com/opensraph/sraph/geom"
)

// VertexWriter is an interface for generating a multi contour polyline as a triangle strip.
type VertexWriter[T geom.Scalar] interface {
	// Write writes a single point to the vertex buffer.
	Write(point geom.Point[T])
	// EndContour marks the end of a contour.
	EndContour()
}

// SegmentReceiver is an interface for receiving pruned path segments.
type SegmentReceiver[T geom.Scalar] interface {
	// BeginContour starts a new contour with the given origin point.
	// Every set of path segments will be surrounded by a Begin/EndContour pair with the same origin point.
	BeginContour(origin geom.Point[T], willBeClosed bool)

	// RecordLine records a line segment from p1 to p2.
	// Guaranteed to be non-degenerate except in the single case of stroking
	// where we have a MoveTo followed by any number of degenerate path segments.
	// p1 will always be the last recorded point.
	RecordLine(p1, p2 geom.Point[T])

	// RecordQuad records a quadratic Bézier curve.
	// Guaranteed to be non-degenerate (not a line).
	// p1 will always be the last recorded point.
	RecordQuad(p1, cp, p2 geom.Point[T])

	// RecordConic records a rational quadratic Bézier curve (conic section).
	// Guaranteed to be non-degenerate (not a quad or line).
	// p1 will always be the last recorded point.
	RecordConic(p1, cp, p2 geom.Point[T], weight T)

	// RecordCubic records a cubic Bézier curve.
	// Guaranteed to be trivially non-degenerate (not all 4 points the same).
	// p1 will always be the last recorded point.
	RecordCubic(p1, cp1, cp2, p2 geom.Point[T])

	// EndContour ends the current contour.
	// Every set of path segments will be surrounded by a Begin/EndContour pair with the same origin point.
	// The boolean indicates if the path was closed as the result of an explicit Close invocation.
	EndContour(origin geom.Point[T], withClose bool)
}

// Quad represents a quadratic Bézier curve.
type Quad[T geom.Scalar] struct {
	P1 geom.Point[T] // Start point
	CP geom.Point[T] // Control point
	P2 geom.Point[T] // End point
}

// Last returns the last point of the quad.
func (q Quad[T]) Last() geom.Point[T] {
	return q.P2
}

// Solve evaluates the quadratic curve at parameter t.
func (q Quad[T]) Solve(t T) geom.Point[T] {
	// B(t) = (1-t)²P1 + 2(1-t)t*CP + t²P2
	oneMinusT := 1 - t
	oneMinusTSq := oneMinusT * oneMinusT
	twoOneMinusT := 2 * oneMinusT * t
	tSq := t * t

	return geom.Point[T]{
		X: oneMinusTSq*q.P1.X + twoOneMinusT*q.CP.X + tSq*q.P2.X,
		Y: oneMinusTSq*q.P1.Y + twoOneMinusT*q.CP.Y + tSq*q.P2.Y,
	}
}

// SubdivisionCount returns the minimum number of subdivisions needed for the given scale factor.
func (q Quad[T]) SubdivisionCount(scale T) T {
	return geom.QuadraticSubdivisions(scale, q.P1, q.CP, q.P2)
}

// GetStartDirection returns the starting direction vector of the quad.
func (q Quad[T]) GetStartDirection() *geom.Point[T] {
	dir := q.CP.Sub(q.P1)
	if dir.IsZero() {
		dir = q.P2.Sub(q.P1)
	}
	if dir.IsZero() {
		return nil
	}
	normalized := dir.Normalize()
	return &normalized
}

// GetEndDirection returns the ending direction vector of the quad.
func (q Quad[T]) GetEndDirection() *geom.Point[T] {
	dir := q.P2.Sub(q.CP)
	if dir.IsZero() {
		dir = q.P2.Sub(q.P1)
	}
	if dir.IsZero() {
		return nil
	}
	normalized := dir.Normalize()
	return &normalized
}

// Conic represents a rational quadratic Bézier curve (conic section).
type Conic[T geom.Scalar] struct {
	P1     geom.Point[T] // Start point
	CP     geom.Point[T] // Control point
	P2     geom.Point[T] // End point
	Weight T             // Weight factor
}

// Last returns the last point of the conic.
func (c Conic[T]) Last() geom.Point[T] {
	return c.P2
}

// Solve evaluates the conic curve at parameter t.
func (c Conic[T]) Solve(t T) geom.Point[T] {
	// Rational Bézier curve formula
	oneMinusT := 1 - t
	w0 := oneMinusT * oneMinusT
	w1 := 2 * oneMinusT * t * c.Weight
	w2 := t * t

	denom := w0 + w1 + w2
	if denom == 0 {
		return c.P1 // Fallback
	}

	return geom.Point[T]{
		X: (w0*c.P1.X + w1*c.CP.X + w2*c.P2.X) / denom,
		Y: (w0*c.P1.Y + w1*c.CP.Y + w2*c.P2.Y) / denom,
	}
}

// SubdivisionCount returns the minimum number of subdivisions needed for the given scale factor.
func (c Conic[T]) SubdivisionCount(scale T) T {
	return geom.ConicSubdivisions(scale, c.P1, c.CP, c.P2, c.Weight)
}

// GetStartDirection returns the starting direction vector of the conic.
func (c Conic[T]) GetStartDirection() *geom.Point[T] {
	dir := c.CP.Sub(c.P1)
	if dir.IsZero() {
		dir = c.P2.Sub(c.P1)
	}
	if dir.IsZero() {
		return nil
	}
	normalized := dir.Normalize()
	return &normalized
}

// Cubic represents a cubic Bézier curve.
type Cubic[T geom.Scalar] struct {
	P1  geom.Point[T] // Start point
	CP1 geom.Point[T] // First control point
	CP2 geom.Point[T] // Second control point
	P2  geom.Point[T] // End point
}

// Last returns the last point of the cubic.
func (c Cubic[T]) Last() geom.Point[T] {
	return c.P2
}

// Solve evaluates the cubic curve at parameter t.
func (c Cubic[T]) Solve(t T) geom.Point[T] {
	// B(t) = (1-t)³P1 + 3(1-t)²t*CP1 + 3(1-t)t²*CP2 + t³P2
	oneMinusT := 1 - t
	oneMinusTSq := oneMinusT * oneMinusT
	oneMinusTCub := oneMinusTSq * oneMinusT
	tSq := t * t
	tCub := tSq * t

	coeff1 := oneMinusTCub
	coeff2 := 3 * oneMinusTSq * t
	coeff3 := 3 * oneMinusT * tSq
	coeff4 := tCub

	return geom.Point[T]{
		X: coeff1*c.P1.X + coeff2*c.CP1.X + coeff3*c.CP2.X + coeff4*c.P2.X,
		Y: coeff1*c.P1.Y + coeff2*c.CP1.Y + coeff3*c.CP2.Y + coeff4*c.P2.Y,
	}
}

// SubdivisionCount returns the minimum number of subdivisions needed for the given scale factor.
func (c Cubic[T]) SubdivisionCount(scale T) T {
	return geom.CubicSubdivisions(scale, c.P1, c.CP1, c.CP2, c.P2)
}

// GetStartDirection returns the starting direction vector of the cubic.
func (c Cubic[T]) GetStartDirection() *geom.Point[T] {
	dir := c.CP1.Sub(c.P1)
	if dir.IsZero() {
		dir = c.CP2.Sub(c.P1)
		if dir.IsZero() {
			dir = c.P2.Sub(c.P1)
		}
	}
	if dir.IsZero() {
		return nil
	}
	normalized := dir.Normalize()
	return &normalized
}

// GetEndDirection returns the ending direction vector of the cubic.
func (c Cubic[T]) GetEndDirection() *geom.Point[T] {
	dir := c.P2.Sub(c.CP2)
	if dir.IsZero() {
		dir = c.P2.Sub(c.CP1)
		if dir.IsZero() {
			dir = c.P2.Sub(c.P1)
		}
	}
	if dir.IsZero() {
		return nil
	}
	normalized := dir.Normalize()
	return &normalized
}

// PathToFilledSegments converts a path source to filled segments.
func PathToFilledSegments[T geom.Scalar](source geom.PathSource[T], receiver SegmentReceiver[T]) {
	pruner := &pathPruner[T]{
		receiver:   receiver,
		isStroking: false,
	}
	source.Dispatch(pruner)
}

// PathToStrokedSegments converts a path source to stroked segments.
func PathToStrokedSegments[T geom.Scalar](source geom.PathSource[T], receiver SegmentReceiver[T]) {
	pruner := &pathPruner[T]{
		receiver:   receiver,
		isStroking: true,
	}
	source.Dispatch(pruner)
}

// CountFillStorage counts the storage requirements for filled path vertices.
func CountFillStorage[T geom.Scalar](source geom.PathSource[T], scale T) (pointCount, contourCount int) {
	counter := &storageCounter[T]{scale: scale}
	PathToFilledSegments(source, counter)
	return counter.pointCount, counter.contourCount
}

// PathToFilledVertices converts a path source to filled vertices.
func PathToFilledVertices[T geom.Scalar](source geom.PathSource[T], writer VertexWriter[T], scale T) {
	fillWriter := &pathFillWriter[T]{
		writer: writer,
		scale:  scale,
	}
	PathToFilledSegments(source, fillWriter)
}

// pathPruner is a utility path receiver that prunes empty contours and degenerate path segments.
//
// Some simplifications and guarantees that it implements:
//   - remove duplicate MoveTo operations
//   - ensure Begin/EndContour on every sub-path
//   - ensure a single degenerate line for empty stroked sub-paths
//   - ensure line back to origin for filled sub-paths
//   - trivial Quad to Line
//   - trivial Conic to Quad
//   - trivial Conic to Line
type pathPruner[T geom.Scalar] struct {
	receiver            SegmentReceiver[T]
	isStroking          bool
	contourHasSegments  bool
	contourHasPoints    bool
	contourWillBeClosed bool
	contourOrigin       geom.Point[T]
	currentPoint        geom.Point[T]
}

// MoveTo implements geom.PathReceiver.
func (p *pathPruner[T]) MoveTo(p2 geom.Point[T], willBeClosed bool) {
	if p.isStroking {
		if p.contourHasSegments && !p.contourHasPoints {
			// If we had actual path segments, but none of them went anywhere
			// (i.e. they never generated any points) then we have to record a
			// 0-length line so that stroker can draw "cap boxes"
			p.receiver.RecordLine(p.contourOrigin, p.contourOrigin)
		}
	} else { // !isStroking
		if !p.currentPoint.Eq(p.contourOrigin) {
			// We help fill operations out by manually connecting back to the
			// contour origin - basically all fill operations implicitly close
			// their contours. If the current point is not at the contour
			// origin then we must have encountered both segments and points.
			p.receiver.RecordLine(p.currentPoint, p.contourOrigin)
		}
	}
	if p.contourHasSegments {
		// contourHasSegments implies we have called BeginContour at some
		// point in time, so we need to end it as we've "moved on".
		p.receiver.EndContour(p.contourOrigin, false)
	}
	p.contourOrigin = p2
	p.currentPoint = p2
	p.contourHasSegments = false
	p.contourHasPoints = false
	p.contourWillBeClosed = willBeClosed
	// We will not record a BeginContour for this potential new contour
	// until we get an actual segment within the contour.
	// See segmentEncountered()
}

// LineTo implements geom.PathReceiver.
func (p *pathPruner[T]) LineTo(p2 geom.Point[T]) {
	p.segmentEncountered()
	if !p2.Eq(p.currentPoint) {
		p.receiver.RecordLine(p.currentPoint, p2)
		p.currentPoint = p2
		p.contourHasPoints = true
	}
}

// QuadTo implements geom.PathReceiver.
func (p *pathPruner[T]) QuadTo(cp, p2 geom.Point[T]) {
	if cp.Eq(p.currentPoint) || p2.Eq(cp) {
		// If all 3 are the same, LineTo will handle that for us
		p.LineTo(p2)
	} else {
		p.segmentEncountered()
		p.receiver.RecordQuad(p.currentPoint, cp, p2)
		p.currentPoint = p2
		p.contourHasPoints = true
	}
}

// ConicTo implements geom.PathReceiver.
func (p *pathPruner[T]) ConicTo(cp, p2 geom.Point[T], weight float64) bool {
	if weight == 1.0 {
		p.QuadTo(cp, p2)
	} else if cp.Eq(p.currentPoint) || p2.Eq(cp) || weight == 0.0 {
		p.LineTo(p2)
	} else {
		p.segmentEncountered()
		p.receiver.RecordConic(p.currentPoint, cp, p2, T(weight))
		p.currentPoint = p2
		p.contourHasPoints = true
	}
	return true
}

// CubicTo implements geom.PathReceiver.
func (p *pathPruner[T]) CubicTo(cp1, cp2, p2 geom.Point[T]) {
	p.segmentEncountered()
	if !cp1.Eq(p.currentPoint) || !cp2.Eq(p.currentPoint) || !p2.Eq(p.currentPoint) {
		// We could check if 3 of the 4 points are equal and simplify to a
		// LineTo, but that quantity of compares is overkill for the unlikely
		// case that it will happen. Checking for simplifying to a QuadTo
		// would involve computing the intersection point of the control
		// polygon edges which is too expensive to be worth the benefit.
		p.receiver.RecordCubic(p.currentPoint, cp1, cp2, p2)
		p.currentPoint = p2
		p.contourHasPoints = true
	}
}

// Close implements geom.PathReceiver.
func (p *pathPruner[T]) Close() {
	// Even a {MoveTo(); Close();} sequence generates a "cap box" at the
	// contour origin location, so we always consider this an "encountered"
	// segment.
	p.segmentEncountered()
	if p.isStroking {
		if !p.contourHasPoints {
			p.receiver.RecordLine(p.currentPoint, p.contourOrigin)
			p.contourHasPoints = true
		}
	} else { // !isStroking
		if !p.currentPoint.Eq(p.contourOrigin) {
			p.receiver.RecordLine(p.currentPoint, p.contourOrigin)
		}
	}
	p.receiver.EndContour(p.contourOrigin, true)
	// The following mirrors the actions of MoveTo - we remain open to
	// recording a new contour from this origin point as if we had had
	// a MoveTo, but we perform no other processing that a MoveTo implies.
	p.currentPoint = p.contourOrigin
	p.contourHasSegments = false
	p.contourHasPoints = false
	// We will not record a BeginContour for this potential new contour
	// until we get an actual segment within the contour.
	// See segmentEncountered()
}

// PathEnd implements geom.PathReceiver.
func (p *pathPruner[T]) PathEnd() {
	if !p.isStroking && !p.currentPoint.Eq(p.contourOrigin) {
		p.receiver.RecordLine(p.currentPoint, p.contourOrigin)
	}
	if p.contourHasSegments {
		p.receiver.EndContour(p.contourOrigin, false)
	}
}

func (p *pathPruner[T]) segmentEncountered() {
	if !p.contourHasSegments {
		p.receiver.BeginContour(p.contourOrigin, p.contourWillBeClosed)
		p.contourHasSegments = true
	}
}

// storageCounter counts storage requirements for tessellation.
type storageCounter[T geom.Scalar] struct {
	scale        T
	pointCount   int
	contourCount int
}

// BeginContour implements SegmentReceiver.
func (s *storageCounter[T]) BeginContour(origin geom.Point[T], willBeClosed bool) {
	// This is a new contour
	s.contourCount++

	// This contour will have an implicit "from" point that will be
	// be delivered with the corresponding Segment methods below.
	s.pointCount++
}

// RecordLine implements SegmentReceiver.
func (s *storageCounter[T]) RecordLine(p1, p2 geom.Point[T]) {
	s.pointCount++
}

// RecordQuad implements SegmentReceiver.
func (s *storageCounter[T]) RecordQuad(p1, cp, p2 geom.Point[T]) {
	quad := Quad[T]{P1: p1, CP: cp, P2: p2}
	subdivisions := int(math.Ceil(quad.SubdivisionCount(s.scale).Float64()))
	s.pointCount += subdivisions
}

// RecordConic implements SegmentReceiver.
func (s *storageCounter[T]) RecordConic(p1, cp, p2 geom.Point[T], weight T) {
	conic := Conic[T]{P1: p1, CP: cp, P2: p2, Weight: weight}
	subdivisions := int(math.Ceil(conic.SubdivisionCount(s.scale).Float64()))
	s.pointCount += subdivisions
}

// RecordCubic implements SegmentReceiver.
func (s *storageCounter[T]) RecordCubic(p1, cp1, cp2, p2 geom.Point[T]) {
	cubic := Cubic[T]{P1: p1, CP1: cp1, CP2: cp2, P2: p2}
	subdivisions := int(math.Ceil(cubic.SubdivisionCount(s.scale).Float64()))
	s.pointCount += subdivisions
}

// EndContour implements SegmentReceiver.
func (s *storageCounter[T]) EndContour(origin geom.Point[T], withClose bool) {
	// If the close operation would have resulted in an additional line
	// segment then the pruner will call RecordLine independently.
	// We count contours in the BeginContour method
}

// pathFillWriter writes path segments as vertices for filling.
type pathFillWriter[T geom.Scalar] struct {
	writer VertexWriter[T]
	scale  T
}

// BeginContour implements SegmentReceiver.
func (p *pathFillWriter[T]) BeginContour(origin geom.Point[T], willBeClosed bool) {
	p.writer.Write(origin)
}

// RecordLine implements SegmentReceiver.
func (p *pathFillWriter[T]) RecordLine(p1, p2 geom.Point[T]) {
	p.writer.Write(p2)
}

// RecordQuad implements SegmentReceiver.
func (p *pathFillWriter[T]) RecordQuad(p1, cp, p2 geom.Point[T]) {
	quad := Quad[T]{P1: p1, CP: cp, P2: p2}
	count := math.Ceil(quad.SubdivisionCount(p.scale).Float64())
	for i := 1; i < int(count); i++ {
		t := T(float64(i) / count)
		point := quad.Solve(t)
		p.writer.Write(point)
	}
	p.writer.Write(p2)
}

// RecordConic implements SegmentReceiver.
func (p *pathFillWriter[T]) RecordConic(p1, cp, p2 geom.Point[T], weight T) {
	conic := Conic[T]{P1: p1, CP: cp, P2: p2, Weight: weight}
	count := math.Ceil(conic.SubdivisionCount(p.scale).Float64())
	for i := 1; i < int(count); i++ {
		t := T(float64(i) / count)
		point := conic.Solve(t)
		p.writer.Write(point)
	}
	p.writer.Write(p2)
}

// RecordCubic implements SegmentReceiver.
func (p *pathFillWriter[T]) RecordCubic(p1, cp1, cp2, p2 geom.Point[T]) {
	cubic := Cubic[T]{P1: p1, CP1: cp1, CP2: cp2, P2: p2}
	count := math.Ceil(cubic.SubdivisionCount(p.scale).Float64())
	for i := 1; i < int(count); i++ {
		t := T(float64(i) / count)
		point := cubic.Solve(t)
		p.writer.Write(point)
	}
	p.writer.Write(p2)
}

// EndContour implements SegmentReceiver.
func (p *pathFillWriter[T]) EndContour(origin geom.Point[T], withClose bool) {
	p.writer.EndContour()
}
