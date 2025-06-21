package geom

// FillType defines the fill rule for a path.
type FillType uint8

const (
	// FillTypeNonZero means non-zero winding fill.
	FillTypeNonZero = iota
	// FillTypeEvenOdd means even-odd fill.
	FillTypeEvenOdd
)

// Convexity defines the convexity property of a path.
type Convexity uint8

const (
	// ConvexityUnknown means convexity is unknown.
	ConvexityUnknown Convexity = iota
	// ConvexityConvex means the path is convex.
	ConvexityConvex
)

// PathReceiver receives path segments (lines, curves, etc.) during path traversal.
// Typically used by path iterators or renderers.
type PathReceiver interface {
	// MoveTo moves the current point to p2, starting a new subpath.
	// willBeClosed indicates if the subpath will be closed.
	MoveTo(p2 Point, willBeClosed bool)
	// LineTo draws a straight line from the current point to p2.
	LineTo(p2 Point)
	// QuadTo draws a quadratic Bézier curve to p2 with control point cp.
	QuadTo(cp, p2 Point)
	// ConicTo draws a rational quadratic Bézier curve to p2 with control point cp and weight.
	// Returns false if not supported (default implementation).
	ConicTo(cp, p2 Point, weight float64) bool
	// CubicTo draws a cubic Bézier curve to p2 with control points cp1 and cp2.
	CubicTo(cp1, cp2, p2 Point)
	// Close closes the current subpath.
	Close()
	// PathEnd is called at the end of path traversal for cleanup (optional).
	PathEnd()
}

// PathSource provides path data for traversal and property queries.
type PathSource interface {
	// FillType returns the fill rule for the path (e.g., non-zero, even-odd).
	FillType() FillType
	// Bounds returns the bounding rectangle of the path.
	Bounds() Rect
	// IsConvex returns true if the path is convex.
	IsConvex() bool
	// Dispatch sends all path segments to the given PathReceiver.
	Dispatch(receiver PathReceiver)
}

// NewRectPathSource creates a new PathSource for rectangles.
func NewRectPathSource(rect Rect) PathSource {
	return rectPathSource{rect: rect}
}

// rectPathSource is a PathSource for rectangles.
type rectPathSource struct {
	rect Rect
}

// FillType implements PathSource.
func (r rectPathSource) FillType() FillType {
	return FillTypeNonZero
}

// Bounds implements PathSource.
func (r rectPathSource) Bounds() Rect {
	return r.rect
}

// IsConvex implements PathSource.
func (r rectPathSource) IsConvex() bool {
	return true
}

// Dispatch implements PathSource.
func (r rectPathSource) Dispatch(receiver PathReceiver) {
	if r.rect.IsEmpty() {
		return
	}

	// Draw rectangle as four lines
	topLeft := r.rect.TopLeft()
	topRight := r.rect.TopRight()
	bottomRight := r.rect.BottomRight()
	bottomLeft := r.rect.BottomLeft()

	receiver.MoveTo(topLeft, true)
	receiver.LineTo(topRight)
	receiver.LineTo(bottomRight)
	receiver.LineTo(bottomLeft)
	receiver.Close()
	receiver.PathEnd()
}

// NewEllipsePathSource creates a new PathSource for ellipses.
func NewEllipsePathSource(bounds Rect) PathSource {
	return ellipsePathSource{bounds: bounds}
}

// ellipsePathSource is a PathSource for ellipses.
type ellipsePathSource struct {
	bounds Rect
}

// FillType implements PathSource.
func (e ellipsePathSource) FillType() FillType {
	return FillTypeNonZero
}

// Bounds implements PathSource.
func (e ellipsePathSource) Bounds() Rect {
	return e.bounds
}

// IsConvex implements PathSource.
func (e ellipsePathSource) IsConvex() bool {
	return true
}

// Dispatch implements PathSource.
func (e ellipsePathSource) Dispatch(receiver PathReceiver) {
	if e.bounds.IsEmpty() {
		return
	}

	// Simplified ellipse drawing using quadratic curves
	// In practice, this would use more sophisticated approximation
	center := e.bounds.Center()
	halfWidth := e.bounds.Width() / 2
	halfHeight := e.bounds.Height() / 2

	// Start at rightmost point
	start := NewPoint(center.X()+halfWidth, center.Y())
	receiver.MoveTo(start, true)

	// Approximate ellipse with 4 quadratic curves
	// This is a simplified implementation
	// Top-right quadrant
	cp1 := NewPoint(center.X()+halfWidth, center.Y()-halfHeight)
	p1 := NewPoint(center.X(), center.Y()-halfHeight)
	receiver.QuadTo(cp1, p1)

	// Top-left quadrant
	cp2 := NewPoint(center.X()-halfWidth, center.Y()-halfHeight)
	p2 := NewPoint(center.X()-halfWidth, center.Y())
	receiver.QuadTo(cp2, p2)

	// Bottom-left quadrant
	cp3 := NewPoint(center.X()-halfWidth, center.Y()+halfHeight)
	p3 := NewPoint(center.X(), center.Y()+halfHeight)
	receiver.QuadTo(cp3, p3)

	// Bottom-right quadrant
	cp4 := NewPoint(center.X()+halfWidth, center.Y()+halfHeight)
	receiver.QuadTo(cp4, start)

	receiver.Close()
	receiver.PathEnd()
}
