package geom

// NewRoundRect creates a new RoundRect with the given rectangle and corner radii.
func NewRoundRect[T Scalar](rect Rect[T], radii RoundingRadii[T]) RoundRect[T] {
	return RoundRect[T]{
		Rect:  rect,
		radii: radii,
	}
}

// NewRoundRectOval creates a new RoundRect that is an oval, with radii Eq to half the width and height of the rectangle.
func NewRoundRectOval[T Scalar](rect Rect[T]) RoundRect[T] {
	return NewRoundRect(
		rect,
		NewRoundingRadiiFromSizes(rect.Size().Scale(1/2)),
	)
}

// NewRoundRectRadius creates a new RoundRect with the given rectangle and uniform corner radius.
func NewRoundRectRadius[T Scalar](rect Rect[T], radius T) RoundRect[T] {
	return NewRoundRect(
		rect,
		NewRoundingRadii(radius),
	)
}

func NewRoundRectXY[T Scalar](rect Rect[T], xRadius, yRadius T) RoundRect[T] {
	return NewRoundRect(
		rect,
		NewRoundingRadiiFromSizes(Size[T]{Width: xRadius, Height: yRadius}),
	)
}

func NewRoundRectLTRB[T Scalar](rect Rect[T], left, top, right, bottom T) RoundRect[T] {
	return NewRoundRect(
		rect,
		NewRoundingRadiiLTRB(left, top, right, bottom),
	)
}

// RoundRect represents a rectangle with rounded corners.
type RoundRect[T Scalar] struct {
	Rect[T]
	radii RoundingRadii[T]
}

// Bounds returns the bounding rectangle of the round rect.
func (r *RoundRect[T]) Bounds() Rect[T] {
	return r.Rect
}

// Radius returns the radii for all four corners.
func (r *RoundRect[T]) Radius() RoundingRadii[T] {
	return r.radii
}

// IsRect returns true if all corner radii are zero and the rectangle is not empty.
func (r *RoundRect[T]) IsRect() bool {
	return !r.Rect.IsEmpty() && r.radii.IsEmpty()
}

// IsOval returns true if all corner radii are Eq and Eq to half the width/height.
func (r *RoundRect[T]) IsOval() bool {
	return !r.Bounds().IsEmpty() && r.radii.IsUniform() &&
		NearlyEq(r.radii.TopLeft.Width, r.Rect.Width()/T(2)) &&
		NearlyEq(r.radii.TopLeft.Height, r.Rect.Height()/T(2))
}

// Dispatch sends the path data of the round rect to the given PathReceiver.
// If includeEnd is true, PathEnd will be called at the end.
func (r *RoundRect[T]) Dispatch(receiver PathReceiver[T], includeEnd bool) {
	left := r.Rect.Left
	right := r.Rect.Right
	bottom := r.Rect.Bottom
	top := r.Rect.Top

	receiver.MoveTo(Point[T]{X: left + r.radii.TopLeft.Width, Y: top}, true)
	receiver.LineTo(Point[T]{X: right - r.radii.TopRight.Width, Y: top})

	receiver.ConicTo(
		Point[T]{X: right, Y: top},
		Point[T]{X: right, Y: top + r.radii.TopRight.Height},
		Sqrt2Over2,
	)

	receiver.LineTo(Point[T]{X: right, Y: bottom - r.radii.BottomRight.Height})

	receiver.ConicTo(
		Point[T]{X: right, Y: bottom},
		Point[T]{X: right - r.radii.BottomRight.Width, Y: bottom},
		Sqrt2Over2,
	)

	receiver.LineTo(Point[T]{X: left + r.radii.BottomLeft.Width, Y: bottom})
	receiver.ConicTo(
		Point[T]{X: left, Y: bottom},
		Point[T]{X: left, Y: bottom - r.radii.BottomLeft.Height},
		Sqrt2Over2,
	)

	receiver.LineTo(Point[T]{X: left, Y: top + r.radii.TopLeft.Height})
	receiver.ConicTo(
		Point[T]{X: left, Y: top},
		Point[T]{X: left + r.radii.TopLeft.Width, Y: top},
		Sqrt2Over2,
	)

	receiver.Close()

	if includeEnd {
		receiver.PathEnd()
	}
}

// IsFinite checks if both the rectangle and its radii are finite.
func (r *RoundRect[T]) IsFinite() bool {
	return r.Rect.IsFinite() && r.radii.IsFinite()
}

// Contains checks if the point is contained within the rounded rectangle.
func (r *RoundRect[T]) Contains(p Point[T]) bool {
	if !r.Rect.Contains(p) {
		return false
	}

	var (
		upperLeftDirection  = Point[T]{X: -1, Y: -1}
		upperRightDirection = Point[T]{X: 1, Y: -1}
		lowerLeftDirection  = Point[T]{X: -1, Y: 1}
		lowerRightDirection = Point[T]{X: 1, Y: 1}
	)

	if !r.cornerContains(p, r.Rect.LeftTop(), upperLeftDirection, r.radii.TopLeft) ||
		!r.cornerContains(p, r.Rect.RightTop(), upperRightDirection, r.radii.TopRight) ||
		!r.cornerContains(p, r.Rect.LeftBottom(), lowerLeftDirection, r.radii.BottomLeft) ||
		!r.cornerContains(p, r.Rect.RightBottom(), lowerRightDirection, r.radii.BottomRight) {
		return false
	}

	return true
}

// cornerContains checks if the point p is contained within the rounded corner defined by corner, direction, and radii.
func (r *RoundRect[T]) cornerContains(p Point[T], corner Point[T], direction Point[T], radii Size[T]) bool {
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

// RoundRectPathSource implements PathSource for a single RoundRect.
type RoundRectPathSource[T Scalar] struct {
	roundRect RoundRect[T]
}

func NewRoundRectPathSource[T Scalar](rr RoundRect[T]) PathSource[T] {
	return &RoundRectPathSource[T]{roundRect: rr}
}

// FillType implements PathSource.
func (r *RoundRectPathSource[T]) FillType() FillType {
	return FillTypeNonZero
}

// IsConvex implements PathSource.
func (r *RoundRectPathSource[T]) IsConvex() bool {
	return true
}

// Bounds implements PathSource.
func (r *RoundRectPathSource[T]) Bounds() Rect[T] {
	return r.roundRect.Bounds()
}

// Dispatch implements PathSource.
func (r *RoundRectPathSource[T]) Dispatch(receiver PathReceiver[T]) {
	r.roundRect.Dispatch(receiver, true)
}

// DiffRoundRectPathSource implements PathSource for the difference between two RoundRects.
type DiffRoundRectPathSource[T Scalar] struct {
	outter RoundRect[T]
	inner  RoundRect[T]
}

func NewDiffRoundRectPathSource[T Scalar](outter, inner RoundRect[T]) PathSource[T] {
	return &DiffRoundRectPathSource[T]{
		outter: outter,
		inner:  inner,
	}
}

// FillType implements PathSource.
func (r *DiffRoundRectPathSource[T]) FillType() FillType {
	return FillTypeEvenOdd
}

// IsConvex implements PathSource.
func (r *DiffRoundRectPathSource[T]) IsConvex() bool {
	return false
}

// Bounds implements PathSource.
func (r *DiffRoundRectPathSource[T]) Bounds() Rect[T] {
	return r.outter.Bounds()
}

// Dispatch implements PathSource.
func (r *DiffRoundRectPathSource[T]) Dispatch(receiver PathReceiver[T]) {
	r.outter.Dispatch(receiver, false)
	r.inner.Dispatch(receiver, true)
}
