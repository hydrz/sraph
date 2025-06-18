package dl

import (
	"math"

	"github.com/opensraph/sraph/geom"
)

// PathBuilder provides a way to construct complex paths for drawing operations.
// It implements a subset of the SVG path specification and provides methods
// for building paths incrementally.
type PathBuilder struct {
	// segments contains the path segments
	segments []PathSegment
	// currentPoint tracks the current pen position
	currentPoint geom.Point[Scalar]
	// firstPoint tracks the first point of the current subpath for close operations
	firstPoint geom.Point[Scalar]
	// hasCurrentPoint indicates if currentPoint is valid
	hasCurrentPoint bool
	// bounds tracks the bounding rectangle of the path
	bounds geom.Rect[Scalar]
	// hasBounds indicates if bounds have been set
	hasBounds bool
}

// PathSegment represents a single segment in a path.
type PathSegment interface {
	// Type returns the type of this segment
	Type() PathSegmentType
	// EndPoint returns the end point of this segment
	EndPoint() geom.Point[Scalar]
	// Bounds returns the bounding rectangle of this segment
	Bounds() geom.Rect[Scalar]
}

// PathSegmentType defines the different types of path segments.
type PathSegmentType uint8

const (
	// PathSegmentTypeMoveTo moves the pen to a new position.
	PathSegmentTypeMoveTo PathSegmentType = iota
	// PathSegmentTypeLineTo draws a line to a new position.
	PathSegmentTypeLineTo
	// PathSegmentTypeQuadTo draws a quadratic bezier curve.
	PathSegmentTypeQuadTo
	// PathSegmentTypeCubicTo draws a cubic bezier curve.
	PathSegmentTypeCubicTo
	// PathSegmentTypeClose closes the current subpath.
	PathSegmentTypeClose
)

// MoveToSegment represents a move-to operation.
type MoveToSegment struct {
	Point geom.Point[Scalar]
}

// Type implements PathSegment.Type for MoveToSegment.
func (s *MoveToSegment) Type() PathSegmentType {
	return PathSegmentTypeMoveTo
}

// EndPoint implements PathSegment.EndPoint for MoveToSegment.
func (s *MoveToSegment) EndPoint() geom.Point[Scalar] {
	return s.Point
}

// Bounds implements PathSegment.Bounds for MoveToSegment.
func (s *MoveToSegment) Bounds() geom.Rect[Scalar] {
	return geom.NewRect(s.Point.X, s.Point.Y, s.Point.X, s.Point.Y)
}

// LineToSegment represents a line-to operation.
type LineToSegment struct {
	Point geom.Point[Scalar]
}

// Type implements PathSegment.Type for LineToSegment.
func (s *LineToSegment) Type() PathSegmentType {
	return PathSegmentTypeLineTo
}

// EndPoint implements PathSegment.EndPoint for LineToSegment.
func (s *LineToSegment) EndPoint() geom.Point[Scalar] {
	return s.Point
}

// Bounds implements PathSegment.Bounds for LineToSegment.
func (s *LineToSegment) Bounds() geom.Rect[Scalar] {
	return geom.NewRect(s.Point.X, s.Point.Y, s.Point.X, s.Point.Y)
}

// QuadToSegment represents a quadratic bezier curve operation.
type QuadToSegment struct {
	Control geom.Point[Scalar]
	Point   geom.Point[Scalar]
}

// Type implements PathSegment.Type for QuadToSegment.
func (s *QuadToSegment) Type() PathSegmentType {
	return PathSegmentTypeQuadTo
}

// EndPoint implements PathSegment.EndPoint for QuadToSegment.
func (s *QuadToSegment) EndPoint() geom.Point[Scalar] {
	return s.Point
}

// Bounds implements PathSegment.Bounds for QuadToSegment.
func (s *QuadToSegment) Bounds() geom.Rect[Scalar] {
	// For simplicity, return bounds of control points
	// A more accurate implementation would calculate the actual curve bounds
	minX := s.Control.X
	if s.Point.X < minX {
		minX = s.Point.X
	}
	maxX := s.Control.X
	if s.Point.X > maxX {
		maxX = s.Point.X
	}
	minY := s.Control.Y
	if s.Point.Y < minY {
		minY = s.Point.Y
	}
	maxY := s.Control.Y
	if s.Point.Y > maxY {
		maxY = s.Point.Y
	}
	return geom.NewRect(minX, minY, maxX, maxY)
}

// CubicToSegment represents a cubic bezier curve operation.
type CubicToSegment struct {
	Control1 geom.Point[Scalar]
	Control2 geom.Point[Scalar]
	Point    geom.Point[Scalar]
}

// Type implements PathSegment.Type for CubicToSegment.
func (s *CubicToSegment) Type() PathSegmentType {
	return PathSegmentTypeCubicTo
}

// EndPoint implements PathSegment.EndPoint for CubicToSegment.
func (s *CubicToSegment) EndPoint() geom.Point[Scalar] {
	return s.Point
}

// Bounds implements PathSegment.Bounds for CubicToSegment.
func (s *CubicToSegment) Bounds() geom.Rect[Scalar] {
	// For simplicity, return bounds of control points
	// A more accurate implementation would calculate the actual curve bounds
	minX := s.Control1.X
	if s.Control2.X < minX {
		minX = s.Control2.X
	}
	if s.Point.X < minX {
		minX = s.Point.X
	}

	maxX := s.Control1.X
	if s.Control2.X > maxX {
		maxX = s.Control2.X
	}
	if s.Point.X > maxX {
		maxX = s.Point.X
	}

	minY := s.Control1.Y
	if s.Control2.Y < minY {
		minY = s.Control2.Y
	}
	if s.Point.Y < minY {
		minY = s.Point.Y
	}

	maxY := s.Control1.Y
	if s.Control2.Y > maxY {
		maxY = s.Control2.Y
	}
	if s.Point.Y > maxY {
		maxY = s.Point.Y
	}

	return geom.NewRect(minX, minY, maxX, maxY)
}

// CloseSegment represents a close path operation.
type CloseSegment struct{}

// Type implements PathSegment.Type for CloseSegment.
func (s *CloseSegment) Type() PathSegmentType {
	return PathSegmentTypeClose
}

// EndPoint implements PathSegment.EndPoint for CloseSegment.
func (s *CloseSegment) EndPoint() geom.Point[Scalar] {
	// Close segments don't have a meaningful end point
	return geom.Point[Scalar]{}
}

// Bounds implements PathSegment.Bounds for CloseSegment.
func (s *CloseSegment) Bounds() geom.Rect[Scalar] {
	// Close segments don't contribute to bounds
	return geom.Rect[Scalar]{}
}

// NewPathBuilder creates a new PathBuilder.
func NewPathBuilder() *PathBuilder {
	return &PathBuilder{
		segments: make([]PathSegment, 0),
	}
}

// MoveTo moves the pen to the specified point.
func (pb *PathBuilder) MoveTo(x, y Scalar) *PathBuilder {
	point := geom.Point[Scalar]{X: x, Y: y}
	pb.segments = append(pb.segments, &MoveToSegment{Point: point})
	pb.currentPoint = point
	pb.firstPoint = point
	pb.hasCurrentPoint = true
	pb.updateBounds(geom.NewRect(x, y, x, y))
	return pb
}

// LineTo draws a line from the current point to the specified point.
func (pb *PathBuilder) LineTo(x, y Scalar) *PathBuilder {
	point := geom.Point[Scalar]{X: x, Y: y}
	pb.segments = append(pb.segments, &LineToSegment{Point: point})
	pb.currentPoint = point
	pb.hasCurrentPoint = true
	pb.updateBounds(geom.NewRect(x, y, x, y))
	return pb
}

// QuadTo draws a quadratic bezier curve from the current point.
func (pb *PathBuilder) QuadTo(cx, cy, x, y Scalar) *PathBuilder {
	control := geom.Point[Scalar]{X: cx, Y: cy}
	point := geom.Point[Scalar]{X: x, Y: y}
	segment := &QuadToSegment{Control: control, Point: point}
	pb.segments = append(pb.segments, segment)
	pb.currentPoint = point
	pb.hasCurrentPoint = true
	pb.updateBounds(segment.Bounds())
	return pb
}

// CubicTo draws a cubic bezier curve from the current point.
func (pb *PathBuilder) CubicTo(c1x, c1y, c2x, c2y, x, y Scalar) *PathBuilder {
	control1 := geom.Point[Scalar]{X: c1x, Y: c1y}
	control2 := geom.Point[Scalar]{X: c2x, Y: c2y}
	point := geom.Point[Scalar]{X: x, Y: y}
	segment := &CubicToSegment{Control1: control1, Control2: control2, Point: point}
	pb.segments = append(pb.segments, segment)
	pb.currentPoint = point
	pb.hasCurrentPoint = true
	pb.updateBounds(segment.Bounds())
	return pb
}

// Close closes the current subpath by drawing a line back to the first point.
func (pb *PathBuilder) Close() *PathBuilder {
	pb.segments = append(pb.segments, &CloseSegment{})
	if pb.hasCurrentPoint {
		pb.currentPoint = pb.firstPoint
	}
	return pb
}

// AddRect adds a rectangle to the path.
func (pb *PathBuilder) AddRect(rect geom.Rect[Scalar]) *PathBuilder {
	pb.MoveTo(rect.Left, rect.Top)
	pb.LineTo(rect.Right, rect.Top)
	pb.LineTo(rect.Right, rect.Bottom)
	pb.LineTo(rect.Left, rect.Bottom)
	pb.Close()
	return pb
}

// AddCircle adds a circle to the path using cubic bezier approximations.
func (pb *PathBuilder) AddCircle(center geom.Point[Scalar], radius Scalar) *PathBuilder {
	// Use the standard 4-segment cubic bezier approximation of a circle
	// The control point distance is approximately 0.552 * radius
	c := radius * 0.552284749831

	x, y := center.X, center.Y

	pb.MoveTo(x+radius, y)
	pb.CubicTo(x+radius, y+c, x+c, y+radius, x, y+radius)
	pb.CubicTo(x-c, y+radius, x-radius, y+c, x-radius, y)
	pb.CubicTo(x-radius, y-c, x-c, y-radius, x, y-radius)
	pb.CubicTo(x+c, y-radius, x+radius, y-c, x+radius, y)
	pb.Close()

	return pb
}

// AddEllipse adds an ellipse to the path.
func (pb *PathBuilder) AddEllipse(bounds geom.Rect[Scalar]) *PathBuilder {
	cx := (bounds.Left + bounds.Right) / 2
	cy := (bounds.Top + bounds.Bottom) / 2
	rx := (bounds.Right - bounds.Left) / 2
	ry := (bounds.Bottom - bounds.Top) / 2

	// Control point distances for ellipse
	cx_offset := rx * 0.552284749831
	cy_offset := ry * 0.552284749831

	pb.MoveTo(cx+rx, cy)
	pb.CubicTo(cx+rx, cy+cy_offset, cx+cx_offset, cy+ry, cx, cy+ry)
	pb.CubicTo(cx-cx_offset, cy+ry, cx-rx, cy+cy_offset, cx-rx, cy)
	pb.CubicTo(cx-rx, cy-cy_offset, cx-cx_offset, cy-ry, cx, cy-ry)
	pb.CubicTo(cx+cx_offset, cy-ry, cx+rx, cy-cy_offset, cx+rx, cy)
	pb.Close()

	return pb
}

// AddRoundRect adds a rounded rectangle to the path.
func (pb *PathBuilder) AddRoundRect(rrect geom.RoundRect[Scalar]) *PathBuilder {
	rect := rrect.Rect
	// TODO: Get actual radii from RoundRect - for now use a default
	radius := Scalar(10) // This should come from rrect.GetRadii()

	pb.MoveTo(rect.Left+radius, rect.Top)
	pb.LineTo(rect.Right-radius, rect.Top)
	pb.QuadTo(rect.Right, rect.Top, rect.Right, rect.Top+radius)
	pb.LineTo(rect.Right, rect.Bottom-radius)
	pb.QuadTo(rect.Right, rect.Bottom, rect.Right-radius, rect.Bottom)
	pb.LineTo(rect.Left+radius, rect.Bottom)
	pb.QuadTo(rect.Left, rect.Bottom, rect.Left, rect.Bottom-radius)
	pb.LineTo(rect.Left, rect.Top+radius)
	pb.QuadTo(rect.Left, rect.Top, rect.Left+radius, rect.Top)
	pb.Close()

	return pb
}

// AddSuperellipse adds a superellipse to the path.
// A superellipse is a generalization of an ellipse with adjustable curvature.
func (pb *PathBuilder) AddSuperellipse(bounds geom.Rect[Scalar], n float32) *PathBuilder {
	// TODO: Implement superellipse path generation
	// This would create a superellipse with the given exponent n
	return pb
}

// AddPolygon adds a regular polygon to the path.
func (pb *PathBuilder) AddPolygon(center geom.Point[Scalar], radius Scalar, sides int, rotation Scalar) *PathBuilder {
	if sides < 3 {
		return pb
	}

	angleStep := 2.0 * 3.14159265359 / float64(sides)
	startAngle := float64(rotation)

	// Calculate first point
	firstAngle := startAngle
	firstX := center.X + radius*Scalar(cos(firstAngle))
	firstY := center.Y + radius*Scalar(sin(firstAngle))
	pb.MoveTo(firstX, firstY)

	// Add lines to other vertices
	for i := 1; i < sides; i++ {
		angle := startAngle + float64(i)*angleStep
		x := center.X + radius*Scalar(cos(angle))
		y := center.Y + radius*Scalar(sin(angle))
		pb.LineTo(x, y)
	}

	pb.Close()
	return pb
}

// AddStar adds a star shape to the path.
func (pb *PathBuilder) AddStar(center geom.Point[Scalar], outerRadius, innerRadius Scalar, points int, rotation Scalar) *PathBuilder {
	if points < 3 {
		return pb
	}

	angleStep := 3.14159265359 / float64(points)
	startAngle := float64(rotation)

	// Start at first outer point
	firstAngle := startAngle
	firstX := center.X + outerRadius*Scalar(cos(firstAngle))
	firstY := center.Y + outerRadius*Scalar(sin(firstAngle))
	pb.MoveTo(firstX, firstY)

	// Alternate between outer and inner points
	for i := 0; i < points*2; i++ {
		radius := outerRadius
		if i%2 == 1 {
			radius = innerRadius
		}
		angle := startAngle + float64(i+1)*angleStep
		x := center.X + radius*Scalar(cos(angle))
		y := center.Y + radius*Scalar(sin(angle))
		pb.LineTo(x, y)
	}

	pb.Close()
	return pb
}

// AddArrow adds an arrow shape to the path.
func (pb *PathBuilder) AddArrow(start, end geom.Point[Scalar], headLength, headWidth Scalar) *PathBuilder {
	// Calculate arrow direction
	dx := end.X - start.X
	dy := end.Y - start.Y
	length := Scalar(sqrt(float64(dx*dx + dy*dy)))

	if length == 0 {
		return pb
	}

	// Normalize direction
	dirX := dx / length
	dirY := dy / length

	// Calculate perpendicular
	perpX := -dirY
	perpY := dirX

	// Arrow shaft
	pb.MoveTo(start.X, start.Y)
	shaftEndX := end.X - dirX*headLength
	shaftEndY := end.Y - dirY*headLength
	pb.LineTo(shaftEndX, shaftEndY)

	// Arrow head
	headBaseX := shaftEndX + perpX*headWidth/2
	headBaseY := shaftEndY + perpY*headWidth/2
	pb.LineTo(headBaseX, headBaseY)
	pb.LineTo(end.X, end.Y)

	headBaseX2 := shaftEndX - perpX*headWidth/2
	headBaseY2 := shaftEndY - perpY*headWidth/2
	pb.LineTo(headBaseX2, headBaseY2)
	pb.LineTo(shaftEndX, shaftEndY)

	return pb
}

// AddDashedLine adds a dashed line to the path.
func (pb *PathBuilder) AddDashedLine(start, end geom.Point[Scalar], dashPattern []Scalar, phase Scalar) *PathBuilder {
	// TODO: Implement dashed line generation
	// This would create a dashed line following the dash pattern
	pb.MoveTo(start.X, start.Y)
	pb.LineTo(end.X, end.Y)
	return pb
}

// AddWavyLine adds a wavy line to the path.
func (pb *PathBuilder) AddWavyLine(start, end geom.Point[Scalar], amplitude, frequency Scalar) *PathBuilder {
	// Calculate line direction and length
	dx := end.X - start.X
	dy := end.Y - start.Y
	length := Scalar(sqrt(float64(dx*dx + dy*dy)))

	if length == 0 {
		return pb
	}

	// Normalize direction and calculate perpendicular
	dirX := dx / length
	dirY := dy / length
	perpX := -dirY
	perpY := dirX

	pb.MoveTo(start.X, start.Y)

	// Create wavy line with multiple segments
	segments := int(length * frequency / 10) // Approximate number of segments
	if segments < 2 {
		segments = 2
	}

	for i := 1; i <= segments; i++ {
		t := Scalar(i) / Scalar(segments)
		baseX := start.X + t*dx
		baseY := start.Y + t*dy

		// Calculate wave offset
		wavePhase := t * frequency * 2 * 3.14159265359
		waveOffset := amplitude * Scalar(sin(float64(wavePhase)))

		// Apply wave offset perpendicular to line direction
		x := baseX + perpX*waveOffset
		y := baseY + perpY*waveOffset

		pb.LineTo(x, y)
	}

	return pb
}

// AddSpiral adds a spiral to the path.
func (pb *PathBuilder) AddSpiral(center geom.Point[Scalar], startRadius, endRadius Scalar, turns float64, clockwise bool) *PathBuilder {
	if turns <= 0 {
		return pb
	}

	steps := int(turns * 50) // Number of line segments per turn
	angleStep := 2 * 3.14159265359 * turns / float64(steps)
	radiusStep := (endRadius - startRadius) / Scalar(steps)

	if !clockwise {
		angleStep = -angleStep
	}

	// Start point
	x := center.X + startRadius
	y := center.Y
	pb.MoveTo(x, y)

	// Generate spiral points
	for i := 1; i <= steps; i++ {
		angle := float64(i) * angleStep
		radius := startRadius + Scalar(i)*radiusStep
		x := center.X + radius*Scalar(cos(angle))
		y := center.Y + radius*Scalar(sin(angle))
		pb.LineTo(x, y)
	}

	return pb
}

// Helper functions for mathematical operations
func cos(x float64) float64 {
	return math.Cos(x)
}

func sin(x float64) float64 {
	return math.Sin(x)
}

func sqrt(x float64) float64 {
	return math.Sqrt(x)
}

// Bounds returns the bounding rectangle of the path.
func (pb *PathBuilder) Bounds() geom.Rect[Scalar] {
	return pb.bounds
}

// GetSegments returns the path segments.
func (pb *PathBuilder) GetSegments() []PathSegment {
	return pb.segments
}

// IsEmpty returns true if the path has no segments.
func (pb *PathBuilder) IsEmpty() bool {
	return len(pb.segments) == 0
}

// Reset clears all path segments.
func (pb *PathBuilder) Reset() *PathBuilder {
	pb.segments = pb.segments[:0]
	pb.currentPoint = geom.Point[Scalar]{}
	pb.firstPoint = geom.Point[Scalar]{}
	pb.hasCurrentPoint = false
	pb.bounds = geom.Rect[Scalar]{}
	pb.hasBounds = false
	return pb
}

// updateBounds updates the path bounds with a new rectangle.
func (pb *PathBuilder) updateBounds(newBounds geom.Rect[Scalar]) {
	if !pb.hasBounds {
		pb.bounds = newBounds
		pb.hasBounds = true
	} else {
		pb.bounds = pb.bounds.Union(newBounds)
	}
}

// Build creates a Path from the current segments.
func (pb *PathBuilder) Build() *Path {
	// Create a copy of segments to ensure immutability
	segments := make([]PathSegment, len(pb.segments))
	copy(segments, pb.segments)

	return &Path{
		segments: segments,
		bounds:   pb.bounds,
	}
}

// Path represents an immutable sequence of path segments.
type Path struct {
	segments []PathSegment
	bounds   geom.Rect[Scalar]
}

// Bounds returns the bounding rectangle of the path.
func (p *Path) Bounds() geom.Rect[Scalar] {
	return p.bounds
}

// GetSegments returns the path segments.
func (p *Path) GetSegments() []PathSegment {
	return p.segments
}

// IsEmpty returns true if the path has no segments.
func (p *Path) IsEmpty() bool {
	return len(p.segments) == 0
}

// FillType implements PathSource.FillType.
func (p *Path) FillType() geom.FillType {
	// TODO: Store fill type in Path struct
	return geom.FillTypeNonZero
}

// IsConvex implements PathSource.IsConvex.
func (p *Path) IsConvex() bool {
	// TODO: Implement convexity analysis
	return false
}

// Dispatch implements PathSource.Dispatch.
func (p *Path) Dispatch(receiver geom.PathReceiver[Scalar]) {
	for _, segment := range p.segments {
		switch seg := segment.(type) {
		case *MoveToSegment:
			receiver.MoveTo(seg.Point, false) // TODO: Determine if path will be closed
		case *LineToSegment:
			receiver.LineTo(seg.Point)
		case *QuadToSegment:
			receiver.QuadTo(seg.Control, seg.Point)
		case *CubicToSegment:
			receiver.CubicTo(seg.Control1, seg.Control2, seg.Point)
		case *CloseSegment:
			receiver.Close()
		}
	}
}
