package dl

import (
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
	// GetType returns the type of this segment
	GetType() PathSegmentType
	// GetEndPoint returns the end point of this segment
	GetEndPoint() geom.Point[Scalar]
	// GetBounds returns the bounding rectangle of this segment
	GetBounds() geom.Rect[Scalar]
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

// GetType implements PathSegment.GetType for MoveToSegment.
func (s *MoveToSegment) GetType() PathSegmentType {
	return PathSegmentTypeMoveTo
}

// GetEndPoint implements PathSegment.GetEndPoint for MoveToSegment.
func (s *MoveToSegment) GetEndPoint() geom.Point[Scalar] {
	return s.Point
}

// GetBounds implements PathSegment.GetBounds for MoveToSegment.
func (s *MoveToSegment) GetBounds() geom.Rect[Scalar] {
	return geom.NewRect(s.Point.X, s.Point.Y, s.Point.X, s.Point.Y)
}

// LineToSegment represents a line-to operation.
type LineToSegment struct {
	Point geom.Point[Scalar]
}

// GetType implements PathSegment.GetType for LineToSegment.
func (s *LineToSegment) GetType() PathSegmentType {
	return PathSegmentTypeLineTo
}

// GetEndPoint implements PathSegment.GetEndPoint for LineToSegment.
func (s *LineToSegment) GetEndPoint() geom.Point[Scalar] {
	return s.Point
}

// GetBounds implements PathSegment.GetBounds for LineToSegment.
func (s *LineToSegment) GetBounds() geom.Rect[Scalar] {
	return geom.NewRect(s.Point.X, s.Point.Y, s.Point.X, s.Point.Y)
}

// QuadToSegment represents a quadratic bezier curve operation.
type QuadToSegment struct {
	Control geom.Point[Scalar]
	Point   geom.Point[Scalar]
}

// GetType implements PathSegment.GetType for QuadToSegment.
func (s *QuadToSegment) GetType() PathSegmentType {
	return PathSegmentTypeQuadTo
}

// GetEndPoint implements PathSegment.GetEndPoint for QuadToSegment.
func (s *QuadToSegment) GetEndPoint() geom.Point[Scalar] {
	return s.Point
}

// GetBounds implements PathSegment.GetBounds for QuadToSegment.
func (s *QuadToSegment) GetBounds() geom.Rect[Scalar] {
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

// GetType implements PathSegment.GetType for CubicToSegment.
func (s *CubicToSegment) GetType() PathSegmentType {
	return PathSegmentTypeCubicTo
}

// GetEndPoint implements PathSegment.GetEndPoint for CubicToSegment.
func (s *CubicToSegment) GetEndPoint() geom.Point[Scalar] {
	return s.Point
}

// GetBounds implements PathSegment.GetBounds for CubicToSegment.
func (s *CubicToSegment) GetBounds() geom.Rect[Scalar] {
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

// GetType implements PathSegment.GetType for CloseSegment.
func (s *CloseSegment) GetType() PathSegmentType {
	return PathSegmentTypeClose
}

// GetEndPoint implements PathSegment.GetEndPoint for CloseSegment.
func (s *CloseSegment) GetEndPoint() geom.Point[Scalar] {
	// Close segments don't have a meaningful end point
	return geom.Point[Scalar]{}
}

// GetBounds implements PathSegment.GetBounds for CloseSegment.
func (s *CloseSegment) GetBounds() geom.Rect[Scalar] {
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
	pb.updateBounds(segment.GetBounds())
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
	pb.updateBounds(segment.GetBounds())
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

// GetBounds returns the bounding rectangle of the path.
func (pb *PathBuilder) GetBounds() geom.Rect[Scalar] {
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

// GetBounds returns the bounding rectangle of the path.
func (p *Path) GetBounds() geom.Rect[Scalar] {
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
