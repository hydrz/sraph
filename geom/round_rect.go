package geom

// RoundRect represents a rectangle with rounded corners.
// It extends Rect and provides additional methods for rounded rectangle geometry and path generation.
type RoundRect[T Scalar] interface {
	Rect[T]
	// Bounds returns the bounding rectangle of the round rect.
	Bounds() Rect[T]
	// Radius returns the radii for all four corners.
	Radius() RoundingRadii[T]

	// IsRect returns true if all corner radii are zero and the rectangle is not empty.
	IsRect() bool
	// IsOval returns true if all corner radii are equal and equal to half the width/height.
	IsOval() bool

	// Dispatch sends the path data of the round rect to the given PathReceiver.
	// If includeEnd is true, PathEnd will be called at the end.
	Dispatch(receiver PathReceiver[T], includeEnd bool)
}

// NewRoundRect creates a new RoundRect with the given rectangle and corner radii.
func NewRoundRect[T Scalar](rect Rect[T], radii RoundingRadii[T]) RoundRect[T] {
	return &roundRect[T]{
		Rect:  rect,
		radii: radii,
	}
}

// NewRoundRectOval creates a new RoundRect that is an oval, with radii equal to half the width and height of the rectangle.
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
		NewRoundingRadiiFromSizes(NewSize(xRadius, yRadius)),
	)
}

func NewRoundRectLTRB[T Scalar](rect Rect[T], left, top, right, bottom T) RoundRect[T] {
	return NewRoundRect(
		rect,
		NewRoundingRadiiLTRB(left, top, right, bottom),
	)
}

type roundRect[T Scalar] struct {
	Rect[T]
	radii RoundingRadii[T]
}

// Bounds returns the bounding rectangle of the round rect.
func (r *roundRect[T]) Bounds() Rect[T] {
	return r.Rect
}

// Radius returns the radii for all four corners.
func (r *roundRect[T]) Radius() RoundingRadii[T] {
	return r.radii
}

// IsRect returns true if all corner radii are zero and the rectangle is not empty.
func (r *roundRect[T]) IsRect() bool {
	return !r.Rect.IsEmpty() && r.radii.IsEmpty()
}

// IsOval returns true if all corner radii are equal and equal to half the width/height.
func (r *roundRect[T]) IsOval() bool {
	return !r.Bounds().IsEmpty() && r.radii.IsUniform() &&
		Equal(r.radii.TopLeft().Width(), r.Rect.Width()/2) &&
		Equal(r.radii.TopLeft().Height(), r.Rect.Height()/2)
}

// Dispatch sends the path data of the round rect to the given PathReceiver.
// If includeEnd is true, PathEnd will be called at the end.
func (r *roundRect[T]) Dispatch(receiver PathReceiver[T], includeEnd bool) {
	left := r.Rect.Left()
	right := r.Rect.Right()
	bottom := r.Rect.Bottom()
	top := r.Rect.Top()

	receiver.MoveTo(NewPoint(left+r.radii.TopLeft().Width(), r.Rect.Top()), true)
	receiver.LineTo(NewPoint(right-r.radii.TopRight().Width(), r.Rect.Top()))

	receiver.ConicTo(NewPoint(right, top),
		NewPoint(right, top+r.radii.TopRight().Height()),
		Sqrt2Over2)

	receiver.LineTo(NewPoint(right, bottom-r.radii.BottomRight().Height()))

	receiver.ConicTo(NewPoint(right, bottom),
		NewPoint(right-r.radii.BottomRight().Width(), bottom),
		Sqrt2Over2)

	receiver.LineTo(NewPoint(left+r.radii.BottomLeft().Width(), bottom))
	receiver.ConicTo(NewPoint(left, bottom),
		NewPoint(left, bottom-r.radii.BottomLeft().Height()),
		Sqrt2Over2)

	receiver.LineTo(NewPoint(left, top+r.radii.TopLeft().Height()))
	receiver.ConicTo(NewPoint(left, top),
		NewPoint(left+r.radii.TopLeft().Width(), top),
		Sqrt2Over2)

	receiver.Close()

	if includeEnd {
		receiver.PathEnd()
	}
}

// IsFinite overrides Rect's IsFinite method to check both the rectangle and its radii.
func (r *roundRect[T]) IsFinite() bool {
	return r.Rect.IsFinite() && r.radii.IsFinite()
}

// Contains overrides Rect's Contains method to check containment within the rounded rectangle.
func (r *roundRect[T]) Contains(p Point[T]) bool {
	if !r.Rect.Contains(p) {
		return false
	}

	var (
		roundRectUpperLeftDirection  = point[T]{-1, -1}
		roundRectUpperRightDirection = point[T]{1, -1}
		roundRectLowerLeftDirection  = point[T]{-1, 1}
		roundRectLowerRightDirection = point[T]{1, 1}
	)

	if (!r.cornerContains(p, r.Rect.LeftTop(), roundRectUpperLeftDirection, r.radii.TopLeft())) ||
		(!r.cornerContains(p, r.Rect.RightTop(), roundRectUpperRightDirection, r.radii.TopRight())) ||
		(!r.cornerContains(p, r.Rect.LeftBottom(), roundRectLowerLeftDirection, r.radii.BottomLeft())) ||
		(!r.cornerContains(p, r.Rect.RightBottom(), roundRectLowerRightDirection, r.radii.BottomRight())) {
		return false
	}

	return true
}

// cornerContains checks if the point p is contained within the rounded corner defined by corner, direction, and radii.
func (r *roundRect[T]) cornerContains(p Point[T], corner Point[T], direction point[T], radii Size[T]) bool {
	// This corner is not curved, therefore the containment is the same as
	// the previously checked bounds containment.
	if radii.IsZero() {
		return true
	}

	// The positive X,Y distance between the corner and the point.
	cornerRelative := corner.Sub(p).Mul(direction)

	// The distance from the "center" of the corner's elliptical curve.
	// If both numbers are positive then we need to do an elliptical distance
	// check to determine if it is inside the curve.
	// If either number is negative, then the point is outside this quadrant
	// and is governed by inclusion in the bounds and inclusion within other
	// corners of this round rect. In that case, we return true here to allow
	// further evaluation within other quadrants.
	rp := NewPoint(radii.Width(), radii.Height())

	quadrantRelative := rp.Sub(cornerRelative)
	if quadrantRelative.X() <= 0 || quadrantRelative.Y() <= 0 {
		// Not within the curved quadrant of this corner, therefore "inside"
		// relative to this one corner.
		return true
	}

	// Dividing the quadrantRelative point by the radii gives a corresponding
	// location within a unit circle which can be more easily tested for
	// containment. We can use x^2 + y^2 and compare it against the radius
	// squared (1.0) to avoid the sqrt.
	unitCirclePoint := quadrantRelative.Div(rp)
	return unitCirclePoint.LengthSquared() <= 1.0
}

// RoundRectPathSource implements PathSource for a single RoundRect.
// It emits the path for the round rect as a filled shape.
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
// It emits the path for the outer round rect minus the inner round rect, using even-odd fill.
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
