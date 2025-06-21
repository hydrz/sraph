package geom

// RoundRect represents a rectangle with rounded corners defined by corner radii.
// The zero value is a rectangle at the origin with zero radii.
// RoundRect is safe for concurrent use by multiple goroutines.
type RoundRect struct {
	Rect
	Radii RoundingRadii
}

// NewRoundRect creates a RoundRect with the specified rectangle and corner radii.
func NewRoundRect(rect Rect, radii RoundingRadii) RoundRect {
	return RoundRect{
		Rect:  rect,
		Radii: radii,
	}
}

// NewRoundRectOval creates a RoundRect that forms a perfect oval.
// The corner radii are set to half the width and height of the rectangle.
func NewRoundRectOval(rect Rect) RoundRect {
	size := rect.Size()
	return NewRoundRect(
		rect,
		NewRoundingRadiiFromSizes(size.Scale(0.5)),
	)
}

// NewRoundRectRadius creates a RoundRect with uniform corner radius on all corners.
func NewRoundRectRadius(rect Rect, radius Scalar) RoundRect {
	return NewRoundRect(
		rect,
		NewRoundingRadii(radius),
	)
}

// NewRoundRectXY creates a RoundRect with uniform elliptical corner radii.
// All corners use the same xRadius for width and yRadius for height.
func NewRoundRectXY(rect Rect, xRadius, yRadius Scalar) RoundRect {
	return NewRoundRect(
		rect,
		NewRoundingRadiiFromSizes(NewSize(xRadius, yRadius)),
	)
}

// NewRoundRectLTRB creates a RoundRect with individual corner radii.
// The parameters specify the radius for left-top, top-right, right-bottom, and bottom-left corners respectively.
func NewRoundRectLTRB(rect Rect, left, top, right, bottom Scalar) RoundRect {
	return NewRoundRect(
		rect,
		NewRoundingRadiiLTRB(left, top, right, bottom),
	)
}

// Bounds returns the bounding rectangle of the RoundRect.
func (r RoundRect) Bounds() Rect {
	return r.Rect
}

// Radius returns the corner radii configuration for all four corners.
func (r RoundRect) Radius() RoundingRadii {
	return r.Radii
}

// IsRect reports whether the RoundRect is a regular rectangle.
// Returns true if all corner radii are zero and the rectangle is not empty.
func (r RoundRect) IsRect() bool {
	return !r.Rect.IsEmpty() && r.Radii.IsEmpty()
}

// IsOval reports whether the RoundRect forms a perfect oval.
// Returns true if all corner radii are uniform and equal to half the width and height.
func (r RoundRect) IsOval() bool {
	return !r.Bounds().IsEmpty() && r.Radii.IsUniform() &&
		r.Radii.TopLeft().Width() == r.Bounds().Width()/2 &&
		r.Radii.TopLeft().Height() == r.Bounds().Height()/2
}

// Dispatch generates the path data for the RoundRect and sends it to the receiver.
// The path is constructed using line segments and conic curves for rounded corners.
// If includeEnd is true, PathEnd is called when path generation is complete.
func (r RoundRect) Dispatch(receiver PathReceiver, includeEnd bool) {
	left := r.Rect.Left()
	right := r.Rect.Right()
	bottom := r.Rect.Bottom()
	top := r.Rect.Top()

	receiver.MoveTo(NewPoint(left+r.Radii.TopLeft().Width(), top), true)
	receiver.LineTo(NewPoint(right-r.Radii.TopRight().Width(), top))

	receiver.ConicTo(
		NewPoint(right, top),
		NewPoint(right, top+r.Radii.TopRight().Height()),
		Sqrt2Over2,
	)

	receiver.LineTo(NewPoint(right, bottom-r.Radii.BottomRight().Height()))

	receiver.ConicTo(
		NewPoint(right, bottom),
		NewPoint(right-r.Radii.BottomRight().Width(), bottom),
		Sqrt2Over2,
	)

	receiver.LineTo(NewPoint(left+r.Radii.BottomLeft().Width(), bottom))
	receiver.ConicTo(
		NewPoint(left, bottom),
		NewPoint(left, bottom-r.Radii.BottomLeft().Height()),
		Sqrt2Over2,
	)

	receiver.LineTo(NewPoint(left, top+r.Radii.TopLeft().Height()))
	receiver.ConicTo(
		NewPoint(left, top),
		NewPoint(left+r.Radii.TopLeft().Width(), top),
		Sqrt2Over2,
	)

	receiver.Close()

	if includeEnd {
		receiver.PathEnd()
	}
}

// IsFinite reports whether both the rectangle bounds and all corner radii are finite.
// This method is primarily useful for floating-point types.
func (r RoundRect) IsFinite() bool {
	return r.Rect.IsFinite() && r.Radii.IsFinite()
}

// Contains reports whether the point lies within the RoundRect.
// Points on the curved edges are considered inside.
// The method handles both straight edges and rounded corners correctly.
func (r RoundRect) Contains(p Point) bool {
	if !r.Rect.Contains(p) {
		return false
	}

	var (
		upperLeftDirection  = NewPoint(-1, -1)
		upperRightDirection = NewPoint(1, -1)
		lowerLeftDirection  = NewPoint(-1, 1)
		lowerRightDirection = NewPoint(1, 1)
	)

	if !r.cornerContains(p, r.Rect.TopLeft(), upperLeftDirection, r.Radii.TopLeft()) ||
		!r.cornerContains(p, r.Rect.TopRight(), upperRightDirection, r.Radii.TopRight()) ||
		!r.cornerContains(p, r.Rect.BottomLeft(), lowerLeftDirection, r.Radii.BottomLeft()) ||
		!r.cornerContains(p, r.Rect.BottomRight(), lowerRightDirection, r.Radii.BottomRight()) {
		return false
	}

	return true
}

// cornerContains reports whether the point lies within the specified rounded corner.
// The corner is defined by its position, the direction vector, and elliptical radii.
// Returns true for corners with zero radii (sharp corners).
func (r RoundRect) cornerContains(p Point, corner Point, direction Point, radii Size) bool {
	if radii.IsZero() {
		return true
	}

	// Compute the positive X,Y distance between the corner and the point, in the direction of the corner.
	cornerRelative := corner.Sub(p).Mul(direction)

	// The distance from the "center" of the corner's elliptical curve.
	rp := NewPoint(radii.Width(), radii.Height())

	quadrantRelative := rp.Sub(cornerRelative)
	if quadrantRelative.X() <= 0 || quadrantRelative.Y() <= 0 {
		// Not within the curved quadrant of this corner, so "inside" relative to this one corner.
		return true
	}

	// Dividing the quadrantRelative point by the radii gives a corresponding location within a unit circle.
	unitCirclePoint := NewPoint(
		quadrantRelative.X()/rp.X(),
		quadrantRelative.Y()/rp.Y(),
	)
	return unitCirclePoint.LengthSquared() <= 1
}

// RoundRectPathSource implements PathSource for rendering a single RoundRect.
// The zero value is not valid; use NewRoundRectPathSource to create instances.
type RoundRectPathSource struct {
	roundRect RoundRect
}

// NewRoundRectPathSource creates a PathSource that renders the given RoundRect.
func NewRoundRectPathSource(rr RoundRect) PathSource {
	return RoundRectPathSource{roundRect: rr}
}

// FillType returns the fill rule for the RoundRect path.
// RoundRect always uses non-zero winding rule.
func (r RoundRectPathSource) FillType() FillType {
	return FillTypeNonZero
}

// IsConvex reports whether the RoundRect path is convex.
// RoundRect is always convex regardless of corner radii.
func (r RoundRectPathSource) IsConvex() bool {
	return true
}

// Bounds returns the bounding rectangle of the RoundRect.
func (r RoundRectPathSource) Bounds() Rect {
	return r.roundRect.Bounds()
}

// Dispatch generates the complete path data for the RoundRect.
func (r RoundRectPathSource) Dispatch(receiver PathReceiver) {
	r.roundRect.Dispatch(receiver, true)
}

// DiffRoundRectPathSource implements PathSource for the difference between two RoundRects.
// This creates a "donut" or "frame" shape by subtracting the inner RoundRect from the outer one.
// The zero value is not valid; use NewDiffRoundRectPathSource to create instances.
type DiffRoundRectPathSource struct {
	outter RoundRect
	inner  RoundRect
}

// NewDiffRoundRectPathSource creates a PathSource that renders the difference between two RoundRects.
// The result is the outer RoundRect with the inner RoundRect subtracted from it.
func NewDiffRoundRectPathSource(outter, inner RoundRect) PathSource {
	return &DiffRoundRectPathSource{
		outter: outter,
		inner:  inner,
	}
}

// FillType returns the fill rule for the difference path.
// Difference paths use even-odd winding rule to handle the subtraction correctly.
func (r DiffRoundRectPathSource) FillType() FillType {
	return FillTypeEvenOdd
}

// IsConvex reports whether the difference path is convex.
// Difference paths are never convex as they contain holes.
func (r DiffRoundRectPathSource) IsConvex() bool {
	return false
}

// Bounds returns the bounding rectangle of the difference path.
// This is the same as the outer RoundRect's bounds.
func (r DiffRoundRectPathSource) Bounds() Rect {
	return r.outter.Bounds()
}

// Dispatch generates the path data for both the outer and inner RoundRects.
// The outer path is generated first, followed by the inner path for subtraction.
func (r DiffRoundRectPathSource) Dispatch(receiver PathReceiver) {
	r.outter.Dispatch(receiver, false)
	r.inner.Dispatch(receiver, true)
}
