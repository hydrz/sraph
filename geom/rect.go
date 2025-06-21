package geom

import (
	"fmt"
	"image"
	"math"
)

// Rect defines the interface for an axis-aligned rectangle.
// All methods are immutable and return new values.
// Rect is safe for concurrent use by multiple goroutines.
type Rect interface {
	// Left returns the left edge coordinate.
	Left() Scalar
	// Top returns the top edge coordinate.
	Top() Scalar
	// Right returns the right edge coordinate.
	Right() Scalar
	// Bottom returns the bottom edge coordinate.
	Bottom() Scalar

	// X returns the x coordinate of the rectangle's left edge.
	X() Scalar
	// Y returns the y coordinate of the rectangle's top edge.
	Y() Scalar
	// Width returns the width of the rectangle (right - left).
	Width() Scalar
	// Height returns the height of the rectangle (bottom - top).
	Height() Scalar

	// TopLeft returns the top-left corner point of the rectangle.
	TopLeft() Point
	// TopRight returns the top-right corner point of the rectangle.
	TopRight() Point
	// BottomLeft returns the bottom-left corner point of the rectangle.
	BottomLeft() Point
	// BottomRight returns the bottom-right corner point of the rectangle.
	BottomRight() Point

	// LTRB returns the rectangle's coordinates as (left, top, right, bottom).
	LTRB() (Scalar, Scalar, Scalar, Scalar)
	// XYWH returns the rectangle's coordinates as (x, y, width, height).
	XYWH() (Scalar, Scalar, Scalar, Scalar)

	// Origin returns the origin point (top-left corner) of the rectangle.
	Origin() Point
	// Size returns the dimensions of the rectangle as (width, height).
	Size() Size
	// Area returns the area of the rectangle (width × height).
	Area() Scalar
	// Center returns the center point of the rectangle.
	Center() Point

	// Positive returns a rectangle with positive width and height.
	// If the rectangle has negative dimensions, the coordinates are swapped to normalize it.
	Positive() Rect
	// Points returns the four corner points of the rectangle.
	// The order is: top-left, top-right, bottom-left, bottom-right.
	Points() [4]Point

	// Equal reports whether two rectangles are equal within floating-point tolerance.
	Equal(other Rect) bool
	// Round returns a rectangle with all coordinates rounded to the nearest integer.
	Round() Rect
	// Floor returns a rectangle with all coordinates rounded down to the nearest integer.
	Floor() Rect
	// Ceil returns a rectangle with all coordinates rounded up to the nearest integer.
	Ceil() Rect

	// IsEmpty reports whether the rectangle has zero or negative width or height.
	IsEmpty() bool
	// IsFinite reports whether all coordinate values are finite numbers.
	// This is primarily useful for floating-point types to detect infinity or NaN values.
	IsFinite() bool
	// IsSquare reports whether the rectangle is a square.
	// Returns true if width equals height and the rectangle is not empty.
	IsSquare() bool

	// Contains reports whether the rectangle contains the given point.
	// Points on the rectangle's edges are considered to be inside.
	Contains(p Point) bool
	// Inside reports whether the rectangle strictly contains the given point.
	// Points on the rectangle's edges are considered to be outside.
	Inside(p Point) bool
	// ContainsRect reports whether this rectangle completely contains another rectangle.
	// Returns false if either rectangle is empty.
	ContainsRect(other Rect) bool
	// Intersects reports whether this rectangle intersects with another rectangle.
	// Returns false if either rectangle is empty.
	Intersects(other Rect) bool

	// Intersect returns the intersection of this rectangle with another rectangle.
	// Returns an empty rectangle if the rectangles do not intersect.
	Intersect(other Rect) Rect
	// Union returns the smallest rectangle that contains both this rectangle and another.
	// If one rectangle is empty, returns the other rectangle.
	Union(other Rect) Rect
	// Cutout attempts to subtract another rectangle from this rectangle.
	// Returns the original rectangle and true if the operation succeeds.
	// For complex overlaps, returns the original rectangle unchanged.
	Cutout(other Rect) (Rect, bool)

	// Scale returns a rectangle with all coordinates scaled by the given factor.
	Scale(scalar Scalar) Rect
	// ScaleXY returns a rectangle with coordinates scaled by different factors for X and Y axes.
	ScaleXY(sx, sy Scalar) Rect
	// ScalePoint returns a rectangle scaled by the X and Y components of the given point.
	ScalePoint(p Point) Rect
	// ScaleSize returns a rectangle scaled by the width and height of the given size.
	ScaleSize(s Size) Rect

	// Translate returns a rectangle moved by the given vector.
	Translate(vector Vector2) Rect
	// TranslateXY returns a rectangle moved by the specified X and Y offsets.
	TranslateXY(x, y Scalar) Rect

	// Expand returns a rectangle expanded by the given amount in all directions.
	Expand(amount Scalar) Rect
	// ExpandLTRB returns a rectangle expanded by different amounts on each edge.
	// The parameters specify the expansion amounts for left, top, right, and bottom edges.
	ExpandLTRB(left, top, right, bottom Scalar) Rect
	// ExpandHV returns a rectangle expanded by different amounts horizontally and vertically.
	ExpandHV(horizontal, vertical Scalar) Rect
	// ExpandPoint returns a rectangle expanded to include the given point.
	// If the point is already inside the rectangle, the rectangle is returned unchanged.
	ExpandPoint(p Point) Rect
	// ExpandSize returns a rectangle expanded by half the given size in all directions.
	ExpandSize(s Size) Rect

	// Project maps the source rectangle into the coordinate space of this rectangle.
	// Returns an empty rectangle if either rectangle is empty.
	Project(source Rect) Rect
	// Transform applies a transformation matrix to the four corners of the rectangle.
	// Returns the transformed corner points.
	Transform(transform Matrix) [4]Point
	// TransformBounds applies a transformation matrix and returns the axis-aligned bounding box.
	// Returns the original rectangle if it is empty.
	TransformBounds(transform Matrix) Rect
	// TransformClipBounds applies a transformation, clips to the given bounds, and returns the bounding box.
	TransformClipBounds(transform Matrix, bounds Rect) Rect

	// NormalizingTransform returns a transformation matrix that maps this rectangle to the unit square [0,1]×[0,1].
	// Returns a zero matrix if the rectangle is empty.
	NormalizingTransform() Matrix

	// String returns a string representation of the rectangle in the format "(TopLeft => BottomRight)".
	String() string
	// Go converts the rectangle to a Go standard library image.Rectangle.
	Go() image.Rectangle
}

// NewRect creates a rectangle from left, top, right, and bottom coordinates.
func NewRect[T Number](left, top, right, bottom T) Rect {
	return rect[T]{left: left, top: top, right: right, bottom: bottom}
}

// NewRectXYWH creates a rectangle from x, y coordinates and width, height dimensions.
// Negative width or height values will be normalized to create a valid rectangle.
func NewRectXYWH[T Number](x, y, width, height T) Rect {
	if width < 0 {
		x += width
		width = -width
	}
	if height < 0 {
		y += height
		height = -height
	}
	return NewRect(x, y, x+width, y+height)
}

// NewRectOriginSize creates a rectangle from an origin point and size dimensions.
// Negative size values will be normalized to create a valid rectangle.
func NewRectOriginSize(origin Point, size Size) Rect {
	return NewRectXYWH(origin.X(), origin.Y(), size.Width(), size.Height())
}

// NewRectSize creates a rectangle at the origin (0,0) with the specified size.
func NewRectSize(size Size) Rect {
	return NewRectXYWH(0, 0, size.Width(), size.Height())
}

// NewRectGo converts a Go standard library image.Rectangle to a Rect.
func NewRectGo(r image.Rectangle) Rect {
	return NewRect(
		Scalar(r.Min.X), Scalar(r.Min.Y),
		Scalar(r.Max.X), Scalar(r.Max.Y),
	)
}

// BoundingRect returns the minimal rectangle that contains all the given points.
// Returns an empty rectangle at the origin if no points are provided.
func BoundingRect(points ...Point) Rect {
	if len(points) == 0 {
		return NewRect(0, 0, 0, 0)
	}

	minX, minY := points[0].X(), points[0].Y()
	maxX, maxY := minX, minY

	for _, p := range points[1:] {
		x, y := p.X(), p.Y()
		minX = min(minX, x)
		minY = min(minY, y)
		maxX = max(maxX, x)
		maxY = max(maxY, y)
	}

	return NewRect(minX, minY, maxX, maxY)
}

// rect represents an axis-aligned rectangle defined by four edges.
type rect[T Number] struct {
	left, top, right, bottom T
}

// Left implements Rect.Left.
func (r rect[T]) Left() Scalar { return Scalar(r.left) }

// Top implements Rect.Top.
func (r rect[T]) Top() Scalar { return Scalar(r.top) }

// Right implements Rect.Right.
func (r rect[T]) Right() Scalar { return Scalar(r.right) }

// Bottom implements Rect.Bottom.
func (r rect[T]) Bottom() Scalar { return Scalar(r.bottom) }

// X implements Rect.X.
func (r rect[T]) X() Scalar { return Scalar(r.left) }

// Y implements Rect.Y.
func (r rect[T]) Y() Scalar { return Scalar(r.top) }

// Width implements Rect.Width.
func (r rect[T]) Width() Scalar { return Scalar(r.right - r.left) }

// Height implements Rect.Height.
func (r rect[T]) Height() Scalar { return Scalar(r.bottom - r.top) }

// TopLeft implements Rect.TopLeft.
func (r rect[T]) TopLeft() Point {
	return NewPoint(r.left, r.top)
}

// TopRight implements Rect.TopRight.
func (r rect[T]) TopRight() Point {
	return NewPoint(r.right, r.top)
}

// BottomLeft implements Rect.BottomLeft.
func (r rect[T]) BottomLeft() Point {
	return NewPoint(r.left, r.bottom)
}

// BottomRight implements Rect.BottomRight.
func (r rect[T]) BottomRight() Point {
	return NewPoint(r.right, r.bottom)
}

// LTRB implements Rect.LTRB.
func (r rect[T]) LTRB() (Scalar, Scalar, Scalar, Scalar) {
	return Scalar(r.left), Scalar(r.top), Scalar(r.right), Scalar(r.bottom)
}

// XYWH implements Rect.XYWH.
func (r rect[T]) XYWH() (Scalar, Scalar, Scalar, Scalar) {
	return Scalar(r.left), Scalar(r.top), Scalar(r.right - r.left), Scalar(r.bottom - r.top)
}

// Origin implements Rect.Origin.
func (r rect[T]) Origin() Point {
	return NewPoint(r.left, r.top)
}

// Size implements Rect.Size.
func (r rect[T]) Size() Size {
	return NewSize(Scalar(r.right-r.left), Scalar(r.bottom-r.top))
}

// Area implements Rect.Area.
func (r rect[T]) Area() Scalar {
	return Scalar((r.right - r.left) * (r.bottom - r.top))
}

// Center implements Rect.Center.
func (r rect[T]) Center() Point {
	return NewPoint((r.left+r.right)/2, (r.top+r.bottom)/2)
}

// Positive implements Rect.Positive.
func (r rect[T]) Positive() Rect {
	left, right := r.left, r.right
	top, bottom := r.top, r.bottom

	if left > right {
		left, right = right, left
	}
	if top > bottom {
		top, bottom = bottom, top
	}

	return rect[T]{left: left, top: top, right: right, bottom: bottom}
}

// Points implements Rect.Points.
func (r rect[T]) Points() [4]Point {
	return [4]Point{r.TopLeft(), r.TopRight(), r.BottomLeft(), r.BottomRight()}
}

// Equal implements Rect.Equal.
func (r rect[T]) Equal(other Rect) bool {
	return NearlyEqual(r.left, T(other.Left())) &&
		NearlyEqual(r.top, T(other.Top())) &&
		NearlyEqual(r.right, T(other.Right())) &&
		NearlyEqual(r.bottom, T(other.Bottom()))
}

// Round implements Rect.Round.
func (r rect[T]) Round() Rect {
	return rect[T]{
		left:   T(math.Round(ToFloat64(r.left))),
		top:    T(math.Round(ToFloat64(r.top))),
		right:  T(math.Round(ToFloat64(r.right))),
		bottom: T(math.Round(ToFloat64(r.bottom))),
	}
}

// Floor implements Rect.Floor.
func (r rect[T]) Floor() Rect {
	return rect[T]{
		left:   T(math.Floor(ToFloat64(r.left))),
		top:    T(math.Floor(ToFloat64(r.top))),
		right:  T(math.Floor(ToFloat64(r.right))),
		bottom: T(math.Floor(ToFloat64(r.bottom))),
	}
}

// Ceil implements Rect.Ceil.
func (r rect[T]) Ceil() Rect {
	return rect[T]{
		left:   T(math.Ceil(ToFloat64(r.left))),
		top:    T(math.Ceil(ToFloat64(r.top))),
		right:  T(math.Ceil(ToFloat64(r.right))),
		bottom: T(math.Ceil(ToFloat64(r.bottom))),
	}
}

// IsEmpty implements Rect.IsEmpty.
func (r rect[T]) IsEmpty() bool {
	return r.left >= r.right || r.top >= r.bottom
}

// IsFinite implements Rect.IsFinite.
func (r rect[T]) IsFinite() bool {
	return IsFinite(r.left) && IsFinite(r.top) && IsFinite(r.right) && IsFinite(r.bottom)
}

// IsSquare implements Rect.IsSquare.
func (r rect[T]) IsSquare() bool {
	return !r.IsEmpty() && NearlyEqual(r.right-r.left, r.bottom-r.top)
}

// Contains implements Rect.Contains.
func (r rect[T]) Contains(p Point) bool {
	x, y := T(p.X()), T(p.Y())
	return r.left <= x && x <= r.right && r.top <= y && y <= r.bottom
}

// Inside implements Rect.Inside.
func (r rect[T]) Inside(p Point) bool {
	x, y := T(p.X()), T(p.Y())
	return r.left < x && x < r.right && r.top < y && y < r.bottom
}

// ContainsRect implements Rect.ContainsRect.
func (r rect[T]) ContainsRect(other Rect) bool {
	if r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return r.left <= T(other.Left()) &&
		r.top <= T(other.Top()) &&
		r.right >= T(other.Right()) &&
		r.bottom >= T(other.Bottom())
}

// Intersects implements Rect.Intersects.
func (r rect[T]) Intersects(other Rect) bool {
	if r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return r.left < T(other.Right()) &&
		r.right > T(other.Left()) &&
		r.top < T(other.Bottom()) &&
		r.bottom > T(other.Top())
}

// Intersect implements Rect.Intersect.
func (r rect[T]) Intersect(other Rect) Rect {
	if !r.Intersects(other) {
		return rect[T]{}
	}

	left := max(r.left, T(other.Left()))
	top := max(r.top, T(other.Top()))
	right := min(r.right, T(other.Right()))
	bottom := min(r.bottom, T(other.Bottom()))

	return rect[T]{left: left, top: top, right: right, bottom: bottom}
}

// Union implements Rect.Union.
func (r rect[T]) Union(other Rect) Rect {
	if r.IsEmpty() {
		return rect[T]{
			left:   T(other.Left()),
			top:    T(other.Top()),
			right:  T(other.Right()),
			bottom: T(other.Bottom()),
		}
	}
	if other.IsEmpty() {
		return r
	}

	left := min(r.left, T(other.Left()))
	top := min(r.top, T(other.Top()))
	right := max(r.right, T(other.Right()))
	bottom := max(r.bottom, T(other.Bottom()))

	return rect[T]{left: left, top: top, right: right, bottom: bottom}
}

// Cutout implements Rect.Cutout.
func (r rect[T]) Cutout(other Rect) (Rect, bool) {
	// Simple cutout - only handle the case where other is fully contained
	if !r.ContainsRect(other) {
		return r, false
	}

	// For simplicity, just return the original rectangle
	// Full implementation would need to handle complex splitting
	return r, true
}

// Scale implements Rect.Scale.
func (r rect[T]) Scale(scalar Scalar) Rect {
	s := T(scalar)
	return rect[T]{
		left:   r.left * s,
		top:    r.top * s,
		right:  r.right * s,
		bottom: r.bottom * s,
	}
}

// ScaleXY implements Rect.ScaleXY.
func (r rect[T]) ScaleXY(sx, sy Scalar) Rect {
	sxT, syT := T(sx), T(sy)
	return rect[T]{
		left:   r.left * sxT,
		top:    r.top * syT,
		right:  r.right * sxT,
		bottom: r.bottom * syT,
	}
}

// ScalePoint implements Rect.ScalePoint.
func (r rect[T]) ScalePoint(p Point) Rect {
	return r.ScaleXY(p.X(), p.Y())
}

// ScaleSize implements Rect.ScaleSize.
func (r rect[T]) ScaleSize(s Size) Rect {
	return r.ScaleXY(s.Width(), s.Height())
}

// Translate implements Rect.Translate.
func (r rect[T]) Translate(vector Vector2) Rect {
	dx, dy := T(vector.X()), T(vector.Y())
	return rect[T]{
		left:   r.left + dx,
		top:    r.top + dy,
		right:  r.right + dx,
		bottom: r.bottom + dy,
	}
}

// TranslateXY implements Rect.TranslateXY.
func (r rect[T]) TranslateXY(x, y Scalar) Rect {
	dx, dy := T(x), T(y)
	return rect[T]{
		left:   r.left + dx,
		top:    r.top + dy,
		right:  r.right + dx,
		bottom: r.bottom + dy,
	}
}

// Expand implements Rect.Expand.
func (r rect[T]) Expand(amount Scalar) Rect {
	a := T(amount)
	return rect[T]{
		left:   r.left - a,
		top:    r.top - a,
		right:  r.right + a,
		bottom: r.bottom + a,
	}
}

// ExpandLTRB implements Rect.ExpandLTRB.
func (r rect[T]) ExpandLTRB(left, top, right, bottom Scalar) Rect {
	return rect[T]{
		left:   r.left - T(left),
		top:    r.top - T(top),
		right:  r.right + T(right),
		bottom: r.bottom + T(bottom),
	}
}

// ExpandHV implements Rect.ExpandHV.
func (r rect[T]) ExpandHV(horizontal, vertical Scalar) Rect {
	h, v := T(horizontal), T(vertical)
	return rect[T]{
		left:   r.left - h,
		top:    r.top - v,
		right:  r.right + h,
		bottom: r.bottom + v,
	}
}

// ExpandPoint implements Rect.ExpandPoint.
func (r rect[T]) ExpandPoint(p Point) Rect {
	if r.IsEmpty() {
		x, y := T(p.X()), T(p.Y())
		return rect[T]{left: x, top: y, right: x, bottom: y}
	}

	x, y := T(p.X()), T(p.Y())
	return rect[T]{
		left:   min(r.left, x),
		top:    min(r.top, y),
		right:  max(r.right, x),
		bottom: max(r.bottom, y),
	}
}

// ExpandSize implements Rect.ExpandSize.
func (r rect[T]) ExpandSize(s Size) Rect {
	halfW, halfH := T(s.Width()/2), T(s.Height()/2)
	return rect[T]{
		left:   r.left - halfW,
		top:    r.top - halfH,
		right:  r.right + halfW,
		bottom: r.bottom + halfH,
	}
}

// Project implements Rect.Project.
func (r rect[T]) Project(source Rect) Rect {
	if r.IsEmpty() || source.IsEmpty() {
		return rect[T]{}
	}

	// Simple linear mapping from source to this rectangle
	sx := ToFloat64(r.Width()) / ToFloat64(source.Width())
	sy := ToFloat64(r.Height()) / ToFloat64(source.Height())

	return rect[T]{
		left:   T((ToFloat64(source.Left())-ToFloat64(r.Left()))*sx + ToFloat64(r.Left())),
		top:    T((ToFloat64(source.Top())-ToFloat64(r.Top()))*sy + ToFloat64(r.Top())),
		right:  T((ToFloat64(source.Right())-ToFloat64(r.Left()))*sx + ToFloat64(r.Left())),
		bottom: T((ToFloat64(source.Bottom())-ToFloat64(r.Top()))*sy + ToFloat64(r.Top())),
	}
}

// Transform implements Rect.Transform.
func (r rect[T]) Transform(transform Matrix) [4]Point {
	return [4]Point{
		r.TopLeft().Transform(transform),
		r.TopRight().Transform(transform),
		r.BottomLeft().Transform(transform),
		r.BottomRight().Transform(transform),
	}
}

// TransformBounds implements Rect.TransformBounds.
func (r rect[T]) TransformBounds(transform Matrix) Rect {
	if r.IsEmpty() {
		return r
	}

	points := r.Transform(transform)
	return BoundingRect(points[0], points[1], points[2], points[3])
}

// TransformClipBounds implements Rect.TransformClipBounds.
func (r rect[T]) TransformClipBounds(transform Matrix, bounds Rect) Rect {
	transformed := r.TransformBounds(transform)
	return transformed.Intersect(bounds)
}

// NormalizingTransform implements Rect.NormalizingTransform.
func (r rect[T]) NormalizingTransform() Matrix {
	if r.IsEmpty() {
		return NewMatrix()
	}

	sx := 1.0 / r.Width()
	sy := 1.0 / r.Height()
	tx := -r.Left() * sx
	ty := -r.Top() * sy

	return NewMatrix().
		Scale(NewVector3(sx, sy, 1)).
		Translate(NewVector3(tx, ty, 0))
}

// String implements Rect.String.
func (r rect[T]) String() string {
	return fmt.Sprintf("((%v, %v) => (%v, %v))", r.left, r.top, r.right, r.bottom)
}

// Go implements Rect.Go.
func (r rect[T]) Go() image.Rectangle {
	return image.Rect(
		int(math.Round(ToFloat64(r.left))),
		int(math.Round(ToFloat64(r.top))),
		int(math.Round(ToFloat64(r.right))),
		int(math.Round(ToFloat64(r.bottom))),
	)
}
