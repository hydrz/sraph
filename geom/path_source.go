package geom

// FillType defines the fill rule for a path.
type FillType uint8

const (
	FillTypeUnknown FillType = iota
	// FillTypeNonZero means non-zero winding fill.
	FillTypeNonZero
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
type PathReceiver[T Scalar] interface {
	// MoveTo moves the current point to p2, starting a new subpath.
	// willBeClosed indicates if the subpath will be closed.
	MoveTo(p2 Point[T], willBeClosed bool)
	// LineTo draws a straight line from the current point to p2.
	LineTo(p2 Point[T])
	// QuadTo draws a quadratic Bézier curve to p2 with control point cp.
	QuadTo(cp, p2 Point[T])
	// ConicTo draws a rational quadratic Bézier curve to p2 with control point cp and weight.
	// Returns false if not supported (default implementation).
	ConicTo(cp, p2 Point[T], weight float64) bool
	// CubicTo draws a cubic Bézier curve to p2 with control points cp1 and cp2.
	CubicTo(cp1, cp2, p2 Point[T])
	// Close closes the current subpath.
	Close()
	// PathEnd is called at the end of path traversal for cleanup (optional).
	PathEnd()
}

// PathSource provides path data for traversal and property queries.
type PathSource[T Scalar] interface {
	// FillType returns the fill rule for the path (e.g., non-zero, even-odd).
	FillType() FillType
	// Bounds returns the bounding rectangle of the path.
	Bounds() Rect[T]
	// IsConvex returns true if the path is convex.
	IsConvex() bool
	// Dispatch sends all path segments to the given PathReceiver.
	Dispatch(receiver PathReceiver[T])
}

// NewRectPathSource creates a new PathSource for rectangles.
func NewRectPathSource[T Scalar](rect Rect[T]) PathSource[T] {
	return &rectPathSource[T]{rect: rect}
}

// rectPathSource is a PathSource for rectangles.
type rectPathSource[T Scalar] struct {
	rect Rect[T]
}

// FillType implements PathSource.
func (r *rectPathSource[T]) FillType() FillType {
	return FillTypeNonZero
}

// Bounds implements PathSource.
func (r *rectPathSource[T]) Bounds() Rect[T] {
	return r.rect
}

// IsConvex implements PathSource.
func (r *rectPathSource[T]) IsConvex() bool {
	return true
}

// Dispatch implements PathSource.
func (r *rectPathSource[T]) Dispatch(receiver PathReceiver[T]) {
	if r.rect.IsEmpty() {
		return
	}

	// Draw rectangle as four lines
	topLeft := r.rect.LeftTop()
	topRight := r.rect.RightTop()
	bottomRight := r.rect.RightBottom()
	bottomLeft := r.rect.LeftBottom()

	receiver.MoveTo(topLeft, true)
	receiver.LineTo(topRight)
	receiver.LineTo(bottomRight)
	receiver.LineTo(bottomLeft)
	receiver.Close()
	receiver.PathEnd()
}

// NewEllipsePathSource creates a new PathSource for ellipses.
func NewEllipsePathSource[T Scalar](bounds Rect[T]) PathSource[T] {
	return &ellipsePathSource[T]{bounds: bounds}
}

// ellipsePathSource is a PathSource for ellipses.
type ellipsePathSource[T Scalar] struct {
	bounds Rect[T]
}

// FillType implements PathSource.
func (e *ellipsePathSource[T]) FillType() FillType {
	return FillTypeNonZero
}

// Bounds implements PathSource.
func (e *ellipsePathSource[T]) Bounds() Rect[T] {
	return e.bounds
}

// IsConvex implements PathSource.
func (e *ellipsePathSource[T]) IsConvex() bool {
	return true
}

// Dispatch implements PathSource.
func (e *ellipsePathSource[T]) Dispatch(receiver PathReceiver[T]) {
	if e.bounds.IsEmpty() {
		return
	}

	// Simplified ellipse drawing using quadratic curves
	// In practice, this would use more sophisticated approximation
	center := e.bounds.Center()
	halfWidth := e.bounds.Width() / T(2)
	halfHeight := e.bounds.Height() / T(2)

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
