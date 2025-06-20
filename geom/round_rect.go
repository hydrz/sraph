package geom

// RoundRect represents a rectangle with rounded corners defined by corner radii.
// The zero value is a rectangle at the origin with zero radii.
// RoundRect is safe for concurrent use by multiple goroutines.
type RoundRect[T TScalar] struct {
	Rect[T]
	Radii RoundingRadii[T]
}

// NewRoundRect creates a RoundRect with the specified rectangle and corner radii.
func NewRoundRect[T TScalar](rect Rect[T], radii RoundingRadii[T]) RoundRect[T] {
	return RoundRect[T]{
		Rect:  rect,
		Radii: radii,
	}
}

// NewRoundRectOval creates a RoundRect that forms a perfect oval.
// The corner radii are set to half the width and height of the rectangle.
func NewRoundRectOval[T TScalar](rect Rect[T]) RoundRect[T] {
	size := rect.Size()
	halfSize := Size[T]{Width: size.Width / T(2), Height: size.Height / T(2)}
	return NewRoundRect(
		rect,
		NewRoundingRadiiFromSizes(halfSize),
	)
}

// NewRoundRectRadius creates a RoundRect with uniform corner radius on all corners.
func NewRoundRectRadius[T TScalar](rect Rect[T], radius T) RoundRect[T] {
	return NewRoundRect(
		rect,
		NewRoundingRadii(radius),
	)
}

// NewRoundRectXY creates a RoundRect with uniform elliptical corner radii.
// All corners use the same xRadius for width and yRadius for height.
func NewRoundRectXY[T TScalar](rect Rect[T], xRadius, yRadius T) RoundRect[T] {
	return NewRoundRect(
		rect,
		NewRoundingRadiiFromSizes(Size[T]{Width: xRadius, Height: yRadius}),
	)
}

// NewRoundRectLTRB creates a RoundRect with individual corner radii.
// The parameters specify the radius for left-top, top-right, right-bottom, and bottom-left corners respectively.
func NewRoundRectLTRB[T TScalar](rect Rect[T], left, top, right, bottom T) RoundRect[T] {
	return NewRoundRect(
		rect,
		NewRoundingRadiiLTRB(left, top, right, bottom),
	)
}

// Bounds returns the bounding rectangle of the RoundRect.
func (r RoundRect[T]) Bounds() Rect[T] {
	return r.Rect
}

// Radius returns the corner radii configuration for all four corners.
func (r RoundRect[T]) Radius() RoundingRadii[T] {
	return r.Radii
}

// IsRect reports whether the RoundRect is a regular rectangle.
// Returns true if all corner radii are zero and the rectangle is not empty.
func (r RoundRect[T]) IsRect() bool {
	return !r.Rect.IsEmpty() && r.Radii.IsEmpty()
}

// IsOval reports whether the RoundRect forms a perfect oval.
// Returns true if all corner radii are uniform and equal to half the width and height.
func (r RoundRect[T]) IsOval() bool {
	return !r.Bounds().IsEmpty() && r.Radii.IsUniform() &&
		NearlyEqual(r.Radii.TopLeft.Width, r.Rect.Width()/T(2)) &&
		NearlyEqual(r.Radii.TopLeft.Height, r.Rect.Height()/T(2))
}

// Dispatch generates the path data for the RoundRect and sends it to the receiver.
// The path is constructed using line segments and conic curves for rounded corners.
// If includeEnd is true, PathEnd is called when path generation is complete.
func (r RoundRect[T]) Dispatch(receiver PathReceiver[T], includeEnd bool) {
	left := r.Rect.Left
	right := r.Rect.Right
	bottom := r.Rect.Bottom
	top := r.Rect.Top

	receiver.MoveTo(Point[T]{X: left + r.Radii.TopLeft.Width, Y: top}, true)
	receiver.LineTo(Point[T]{X: right - r.Radii.TopRight.Width, Y: top})

	receiver.ConicTo(
		Point[T]{X: right, Y: top},
		Point[T]{X: right, Y: top + r.Radii.TopRight.Height},
		Sqrt2Over2,
	)

	receiver.LineTo(Point[T]{X: right, Y: bottom - r.Radii.BottomRight.Height})

	receiver.ConicTo(
		Point[T]{X: right, Y: bottom},
		Point[T]{X: right - r.Radii.BottomRight.Width, Y: bottom},
		Sqrt2Over2,
	)

	receiver.LineTo(Point[T]{X: left + r.Radii.BottomLeft.Width, Y: bottom})
	receiver.ConicTo(
		Point[T]{X: left, Y: bottom},
		Point[T]{X: left, Y: bottom - r.Radii.BottomLeft.Height},
		Sqrt2Over2,
	)

	receiver.LineTo(Point[T]{X: left, Y: top + r.Radii.TopLeft.Height})
	receiver.ConicTo(
		Point[T]{X: left, Y: top},
		Point[T]{X: left + r.Radii.TopLeft.Width, Y: top},
		Sqrt2Over2,
	)

	receiver.Close()

	if includeEnd {
		receiver.PathEnd()
	}
}

// IsFinite reports whether both the rectangle bounds and all corner radii are finite.
// This method is primarily useful for floating-point types.
func (r RoundRect[T]) IsFinite() bool {
	return r.Rect.IsFinite() && r.Radii.IsFinite()
}

// Contains reports whether the point lies within the RoundRect.
// Points on the curved edges are considered inside.
// The method handles both straight edges and rounded corners correctly.
func (r RoundRect[T]) Contains(p Point[T]) bool {
	if !r.Rect.Contains(p) {
		return false
	}

	var (
		upperLeftDirection  = Point[T]{X: -1, Y: -1}
		upperRightDirection = Point[T]{X: 1, Y: -1}
		lowerLeftDirection  = Point[T]{X: -1, Y: 1}
		lowerRightDirection = Point[T]{X: 1, Y: 1}
	)

	if !r.cornerContains(p, r.Rect.TopLeft(), upperLeftDirection, r.Radii.TopLeft) ||
		!r.cornerContains(p, r.Rect.TopRight(), upperRightDirection, r.Radii.TopRight) ||
		!r.cornerContains(p, r.Rect.BottomLeft(), lowerLeftDirection, r.Radii.BottomLeft) ||
		!r.cornerContains(p, r.Rect.BottomRight(), lowerRightDirection, r.Radii.BottomRight) {
		return false
	}

	return true
}

// cornerContains reports whether the point lies within the specified rounded corner.
// The corner is defined by its position, the direction vector, and elliptical radii.
// Returns true for corners with zero radii (sharp corners).
func (r RoundRect[T]) cornerContains(p Point[T], corner Point[T], direction Point[T], radii Size[T]) bool {
	if radii.IsZero() {
		return true
	}

	// Compute the positive X,Y distance between the corner and the point, in the direction of the corner.
	cornerRelative := corner.Sub(p).Mul(direction)

	// The distance from the "center" of the corner's elliptical curve.
	rp := Point[T]{X: radii.Width, Y: radii.Height}

	quadrantRelative := rp.Sub(cornerRelative)
	if quadrantRelative.X <= 0 || quadrantRelative.Y <= 0 {
		// Not within the curved quadrant of this corner, so "inside" relative to this one corner.
		return true
	}

	// Dividing the quadrantRelative point by the radii gives a corresponding location within a unit circle.
	unitCirclePoint := Point[T]{
		X: quadrantRelative.X / rp.X,
		Y: quadrantRelative.Y / rp.Y,
	}
	return unitCirclePoint.LengthSquared() <= T(1)
}

// RoundRectPathSource implements PathSource for rendering a single RoundRect.
// The zero value is not valid; use NewRoundRectPathSource to create instances.
type RoundRectPathSource[T TScalar] struct {
	roundRect RoundRect[T]
}

// NewRoundRectPathSource creates a PathSource that renders the given RoundRect.
func NewRoundRectPathSource[T TScalar](rr RoundRect[T]) PathSource[T] {
	return RoundRectPathSource[T]{roundRect: rr}
}

// FillType returns the fill rule for the RoundRect path.
// RoundRect always uses non-zero winding rule.
func (r RoundRectPathSource[T]) FillType() FillType {
	return FillTypeNonZero
}

// IsConvex reports whether the RoundRect path is convex.
// RoundRect is always convex regardless of corner radii.
func (r RoundRectPathSource[T]) IsConvex() bool {
	return true
}

// Bounds returns the bounding rectangle of the RoundRect.
func (r RoundRectPathSource[T]) Bounds() Rect[T] {
	return r.roundRect.Bounds()
}

// Dispatch generates the complete path data for the RoundRect.
func (r RoundRectPathSource[T]) Dispatch(receiver PathReceiver[T]) {
	r.roundRect.Dispatch(receiver, true)
}

// DiffRoundRectPathSource implements PathSource for the difference between two RoundRects.
// This creates a "donut" or "frame" shape by subtracting the inner RoundRect from the outer one.
// The zero value is not valid; use NewDiffRoundRectPathSource to create instances.
type DiffRoundRectPathSource[T TScalar] struct {
	outter RoundRect[T]
	inner  RoundRect[T]
}

// NewDiffRoundRectPathSource creates a PathSource that renders the difference between two RoundRects.
// The result is the outer RoundRect with the inner RoundRect subtracted from it.
func NewDiffRoundRectPathSource[T TScalar](outter, inner RoundRect[T]) PathSource[T] {
	return &DiffRoundRectPathSource[T]{
		outter: outter,
		inner:  inner,
	}
}

// FillType returns the fill rule for the difference path.
// Difference paths use even-odd winding rule to handle the subtraction correctly.
func (r DiffRoundRectPathSource[T]) FillType() FillType {
	return FillTypeEvenOdd
}

// IsConvex reports whether the difference path is convex.
// Difference paths are never convex as they contain holes.
func (r DiffRoundRectPathSource[T]) IsConvex() bool {
	return false
}

// Bounds returns the bounding rectangle of the difference path.
// This is the same as the outer RoundRect's bounds.
func (r DiffRoundRectPathSource[T]) Bounds() Rect[T] {
	return r.outter.Bounds()
}

// Dispatch generates the path data for both the outer and inner RoundRects.
// The outer path is generated first, followed by the inner path for subtraction.
func (r DiffRoundRectPathSource[T]) Dispatch(receiver PathReceiver[T]) {
	r.outter.Dispatch(receiver, false)
	r.inner.Dispatch(receiver, true)
}
