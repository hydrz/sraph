package geom

import (
	"image"
	"math"
)

// Rect represents an axis-aligned rectangle defined by four edges.
// The zero value is a rectangle at the origin with zero width and height.
// All methods are immutable and return new values.
// Rect is safe for concurrent use by multiple goroutines.
type Rect[T TScalar] struct {
	Left   T
	Top    T
	Right  T
	Bottom T
}

// NewRect creates a rectangle from left, top, right, and bottom coordinates.
func NewRect[T TScalar](left, top, right, bottom T) Rect[T] {
	return Rect[T]{Left: left, Top: top, Right: right, Bottom: bottom}
}

// NewRectXYWH creates a rectangle from x, y coordinates and width, height dimensions.
// Negative width or height values will be normalized to create a valid rectangle.
func NewRectXYWH[T TScalar](x, y, width, height T) Rect[T] {
	if width < 0 {
		x += width
		width = -width
	}
	if height < 0 {
		y += height
		height = -height
	}
	return Rect[T]{Left: x, Top: y, Right: x + width, Bottom: y + height}
}

// NewRectOriginSize creates a rectangle from an origin point and size dimensions.
// Negative size values will be normalized to create a valid rectangle.
func NewRectOriginSize[T TScalar](origin Point[T], size Size[T]) Rect[T] {
	x, y := origin.X, origin.Y
	w, h := size.Width, size.Height
	if w < 0 {
		x += w
		w = -w
	}
	if h < 0 {
		y += h
		h = -h
	}
	return Rect[T]{Left: x, Top: y, Right: x + w, Bottom: y + h}
}

// NewRectSize creates a rectangle at the origin (0,0) with the specified size.
func NewRectSize[T TScalar](size Size[T]) Rect[T] {
	return Rect[T]{Left: 0, Top: 0, Right: size.Width, Bottom: size.Height}
}

// NewRectGo converts a Go standard library image.Rectangle to a Rect.
func NewRectGo[T TScalar](r image.Rectangle) Rect[T] {
	return Rect[T]{
		Left:   T(r.Min.X),
		Top:    T(r.Min.Y),
		Right:  T(r.Max.X),
		Bottom: T(r.Max.Y),
	}
}

// BoundingRect returns the minimal rectangle that contains all the given points.
// Returns an empty rectangle at the origin if no points are provided.
func BoundingRect[T TScalar](points ...Point[T]) Rect[T] {
	if len(points) == 0 {
		return NewRect[T](0, 0, 0, 0)
	}
	left := points[0].X
	top := points[0].Y
	right := left
	bottom := top
	for _, p := range points {
		if p.X < left {
			left = p.X
		}
		if p.X > right {
			right = p.X
		}
		if p.Y < top {
			top = p.Y
		}
		if p.Y > bottom {
			bottom = p.Y
		}
	}
	return NewRect(left, top, right, bottom)
}

// X returns the x coordinate of the rectangle's left edge.
func (r Rect[T]) X() T { return r.Left }

// Y returns the y coordinate of the rectangle's top edge.
func (r Rect[T]) Y() T { return r.Top }

// Width returns the width of the rectangle (right - left).
func (r Rect[T]) Width() T { return r.Right - r.Left }

// Height returns the height of the rectangle (bottom - top).
func (r Rect[T]) Height() T { return r.Bottom - r.Top }

// TopLeft returns the top-left corner point of the rectangle.
func (r Rect[T]) TopLeft() Point[T]     { return Point[T]{r.Left, r.Top} }
func (r Rect[T]) TopRight() Point[T]    { return Point[T]{r.Right, r.Top} }
func (r Rect[T]) BottomLeft() Point[T]  { return Point[T]{r.Left, r.Bottom} }
func (r Rect[T]) BottomRight() Point[T] { return Point[T]{r.Right, r.Bottom} }

// LTRB returns the rectangle's coordinates as (left, top, right, bottom).
func (r Rect[T]) LTRB() (T, T, T, T) { return r.Left, r.Top, r.Right, r.Bottom }

// XYWH returns the rectangle's coordinates as (x, y, width, height).
func (r Rect[T]) XYWH() (T, T, T, T) { return r.Left, r.Top, r.Width(), r.Height() }

// Origin returns the origin point (top-left corner) of the rectangle.
func (r Rect[T]) Origin() Point[T] { return Point[T]{r.Left, r.Top} }

// Size returns the dimensions of the rectangle as (width, height).
func (r Rect[T]) Size() Size[T] { return Size[T]{r.Width(), r.Height()} }

// Area returns the area of the rectangle (width × height).
func (r Rect[T]) Area() T { return r.Width() * r.Height() }

// Center returns the center point of the rectangle.
func (r Rect[T]) Center() Point[T] {
	return Point[T]{(r.Left + r.Right) / T(2), (r.Top + r.Bottom) / T(2)}
}

// Positive returns a rectangle with positive width and height.
// If the rectangle has negative dimensions, the coordinates are swapped to normalize it.
func (r Rect[T]) Positive() Rect[T] {
	left, right := r.Left, r.Right
	top, bottom := r.Top, r.Bottom
	if left > right {
		left, right = right, left
	}
	if top > bottom {
		top, bottom = bottom, top
	}
	return NewRect(left, top, right, bottom)
}

// Points returns the four corner points of the rectangle.
// The order is: top-left, top-right, bottom-left, bottom-right.
func (r Rect[T]) Points() [4]Point[T] {
	return [4]Point[T]{
		r.TopLeft(),
		r.TopRight(),
		r.BottomLeft(),
		r.BottomRight(),
	}
}

// Equal reports whether two rectangles are equal within floating-point tolerance.
func (r Rect[T]) Equal(other Rect[T]) bool {
	return NearlyEqual(r.Left, other.Left) &&
		NearlyEqual(r.Top, other.Top) &&
		NearlyEqual(r.Right, other.Right) &&
		NearlyEqual(r.Bottom, other.Bottom)
}

// Round returns a rectangle with all coordinates rounded to the nearest integer.
func (r Rect[T]) Round() Rect[T] {
	return NewRect(
		T(math.Round(ToFloat64(r.Left))),
		T(math.Round(ToFloat64(r.Top))),
		T(math.Round(ToFloat64(r.Right))),
		T(math.Round(ToFloat64(r.Bottom))),
	)
}

// Floor returns a rectangle with all coordinates rounded down to the nearest integer.
func (r Rect[T]) Floor() Rect[T] {
	return NewRect(
		T(math.Floor(ToFloat64(r.Left))),
		T(math.Floor(ToFloat64(r.Top))),
		T(math.Floor(ToFloat64(r.Right))),
		T(math.Floor(ToFloat64(r.Bottom))),
	)
}

// Ceil returns a rectangle with all coordinates rounded up to the nearest integer.
func (r Rect[T]) Ceil() Rect[T] {
	return NewRect(
		T(math.Ceil(ToFloat64(r.Left))),
		T(math.Ceil(ToFloat64(r.Top))),
		T(math.Ceil(ToFloat64(r.Right))),
		T(math.Ceil(ToFloat64(r.Bottom))),
	)
}

// IsEmpty reports whether the rectangle has zero or negative width or height.
func (r Rect[T]) IsEmpty() bool {
	return r.Width() <= 0 || r.Height() <= 0
}

// IsFinite reports whether all coordinate values are finite numbers.
// This is primarily useful for floating-point types to detect infinity or NaN values.
func (r Rect[T]) IsFinite() bool {
	return IsFinite(r.Left) && IsFinite(r.Top) && IsFinite(r.Right) && IsFinite(r.Bottom)
}

// IsSquare reports whether the rectangle is a square.
// Returns true if width equals height and the rectangle is not empty.
func (r Rect[T]) IsSquare() bool {
	return NearlyEqual(r.Width(), r.Height()) && !r.IsEmpty()
}

// Contains reports whether the rectangle contains the given point.
// Points on the rectangle's edges are considered to be inside.
func (r Rect[T]) Contains(p Point[T]) bool {
	if r.IsEmpty() {
		return false
	}
	return p.X >= r.Left && p.X <= r.Right &&
		p.Y >= r.Top && p.Y <= r.Bottom
}

// Inside reports whether the rectangle strictly contains the given point.
// Points on the rectangle's edges are considered to be outside.
func (r Rect[T]) Inside(p Point[T]) bool {
	if r.IsEmpty() {
		return false
	}
	return p.X > r.Left && p.X < r.Right &&
		p.Y > r.Top && p.Y < r.Bottom
}

// ContainsRect reports whether this rectangle completely contains another rectangle.
// Returns false if either rectangle is empty.
func (r Rect[T]) ContainsRect(other Rect[T]) bool {
	if r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return other.Left >= r.Left && other.Right <= r.Right &&
		other.Top >= r.Top && other.Bottom <= r.Bottom
}

// Intersects reports whether this rectangle intersects with another rectangle.
// Returns false if either rectangle is empty.
func (r Rect[T]) Intersects(other Rect[T]) bool {
	if r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return r.Left < other.Right && r.Right > other.Left &&
		r.Top < other.Bottom && r.Bottom > other.Top
}

// Intersect returns the intersection of this rectangle with another rectangle.
// Returns an empty rectangle if the rectangles do not intersect.
func (r Rect[T]) Intersect(other Rect[T]) Rect[T] {
	if !r.Intersects(other) {
		return Rect[T]{}
	}
	left := r.Left
	if other.Left > left {
		left = other.Left
	}
	top := r.Top
	if other.Top > top {
		top = other.Top
	}
	right := r.Right
	if other.Right < right {
		right = other.Right
	}
	bottom := r.Bottom
	if other.Bottom < bottom {
		bottom = other.Bottom
	}
	result := NewRect(left, top, right, bottom)

	if result.IsEmpty() {
		return Rect[T]{}
	}

	return result
}

// Union returns the smallest rectangle that contains both this rectangle and another.
// If one rectangle is empty, returns the other rectangle.
func (r Rect[T]) Union(other Rect[T]) Rect[T] {
	if r.IsEmpty() {
		return other
	}
	if other.IsEmpty() {
		return r
	}
	left := r.Left
	if other.Left < left {
		left = other.Left
	}
	top := r.Top
	if other.Top < top {
		top = other.Top
	}
	right := r.Right
	if other.Right > right {
		right = other.Right
	}
	bottom := r.Bottom
	if other.Bottom > bottom {
		bottom = other.Bottom
	}
	return NewRect(left, top, right, bottom)
}

// Cutout attempts to subtract another rectangle from this rectangle.
// Returns the original rectangle and true if the operation succeeds.
// For complex overlaps, returns the original rectangle unchanged.
func (r Rect[T]) Cutout(other Rect[T]) (Rect[T], bool) {
	if !r.Intersects(other) {
		return r, true
	}
	if other.ContainsRect(r) {
		return NewRect[T](0, 0, 0, 0), false
	}
	return r, true
}

// Scale returns a rectangle with all coordinates scaled by the given factor.
func (r Rect[T]) Scale(scalar T) Rect[T] {
	return NewRect(r.Left*scalar, r.Top*scalar, r.Right*scalar, r.Bottom*scalar)
}

// ScaleXY returns a rectangle with coordinates scaled by different factors for X and Y axes.
func (r Rect[T]) ScaleXY(sx, sy T) Rect[T] {
	return NewRect(r.Left*sx, r.Top*sy, r.Right*sx, r.Bottom*sy)
}

// ScalePoint returns a rectangle scaled by the X and Y components of the given point.
func (r Rect[T]) ScalePoint(p Point[T]) Rect[T] {
	return r.ScaleXY(p.X, p.Y)
}

// ScaleSize returns a rectangle scaled by the width and height of the given size.
func (r Rect[T]) ScaleSize(s Size[T]) Rect[T] {
	return r.ScaleXY(s.Width, s.Height)
}

// Translate returns a rectangle moved by the given vector.
func (r Rect[T]) Translate(vector Vector2[T]) Rect[T] {
	return NewRect(r.Left+vector.X, r.Top+vector.Y, r.Right+vector.X, r.Bottom+vector.Y)
}

// TranslateXY returns a rectangle moved by the specified X and Y offsets.
func (r Rect[T]) TranslateXY(x, y T) Rect[T] {
	return NewRect(r.Left+x, r.Top+y, r.Right+x, r.Bottom+y)
}

// Expand returns a rectangle expanded by the given amount in all directions.
func (r Rect[T]) Expand(amount T) Rect[T] {
	return NewRect(r.Left-amount, r.Top-amount, r.Right+amount, r.Bottom+amount)
}

// ExpandLTRB returns a rectangle expanded by different amounts on each edge.
// The parameters specify the expansion amounts for left, top, right, and bottom edges.
func (r Rect[T]) ExpandLTRB(left, top, right, bottom T) Rect[T] {
	return NewRect(r.Left-left, r.Top-top, r.Right+right, r.Bottom+bottom)
}

// ExpandHV returns a rectangle expanded by different amounts horizontally and vertically.
func (r Rect[T]) ExpandHV(horizontal, vertical T) Rect[T] {
	return NewRect(r.Left-horizontal, r.Top-vertical, r.Right+horizontal, r.Bottom+vertical)
}

// ExpandPoint returns a rectangle expanded to include the given point.
// If the point is already inside the rectangle, the rectangle is returned unchanged.
func (r Rect[T]) ExpandPoint(p Point[T]) Rect[T] {
	left := r.Left
	if p.X < left {
		left = p.X
	}
	top := r.Top
	if p.Y < top {
		top = p.Y
	}
	right := r.Right
	if p.X > right {
		right = p.X
	}
	bottom := r.Bottom
	if p.Y > bottom {
		bottom = p.Y
	}
	return NewRect(left, top, right, bottom)
}

// ExpandSize returns a rectangle expanded by half the given size in all directions.
func (r Rect[T]) ExpandSize(s Size[T]) Rect[T] {
	return r.ExpandHV(s.Width/T(2), s.Height/T(2))
}

// Project maps the source rectangle into the coordinate space of this rectangle.
// Returns an empty rectangle if either rectangle is empty.
func (r Rect[T]) Project(source Rect[T]) Rect[T] {
	if r.IsEmpty() || source.IsEmpty() {
		return NewRect[T](0, 0, 0, 0)
	}
	scaleX := r.Width() / source.Width()
	scaleY := r.Height() / source.Height()
	return NewRect(
		r.Left+scaleX*(source.Left-source.Left),
		r.Top+scaleY*(source.Top-source.Top),
		r.Left+scaleX*(source.Right-source.Left),
		r.Top+scaleY*(source.Bottom-source.Top),
	)
}

// Transform applies a transformation matrix to the four corners of the rectangle.
// Returns the transformed corner points.
func (r Rect[T]) Transform(transform Matrix[T]) [4]Point[T] {
	corners := r.Points()
	return [4]Point[T]{
		transform.transformPoint(corners[0]),
		transform.transformPoint(corners[1]),
		transform.transformPoint(corners[2]),
		transform.transformPoint(corners[3]),
	}
}

// TransformBounds applies a transformation matrix and returns the axis-aligned bounding box.
// Returns the original rectangle if it is empty.
func (r Rect[T]) TransformBounds(transform Matrix[T]) Rect[T] {
	if r.IsEmpty() {
		return r
	}
	transformed := r.Transform(transform)
	left := transformed[0].X
	top := transformed[0].Y
	right := left
	bottom := top
	for i := 1; i < 4; i++ {
		x, y := transformed[i].X, transformed[i].Y
		if x < left {
			left = x
		}
		if x > right {
			right = x
		}
		if y < top {
			top = y
		}
		if y > bottom {
			bottom = y
		}
	}
	return NewRect(left, top, right, bottom)
}

// TransformClipBounds applies a transformation, clips to the given bounds, and returns the bounding box.
func (r Rect[T]) TransformClipBounds(transform Matrix[T], bounds Rect[T]) Rect[T] {
	transformed := r.TransformBounds(transform)
	return transformed.Intersect(bounds)
}

// NormalizingTransform returns a transformation matrix that maps this rectangle to the unit square [0,1]×[0,1].
// Returns a zero matrix if the rectangle is empty.
func (r Rect[T]) NormalizingTransform() Matrix[T] {
	if r.IsEmpty() {
		return Matrix[T]{}
	}
	scaleX := 1 / r.Width()
	scaleY := 1 / r.Height()
	matrix := Matrix[T]{}
	matrix = matrix.Scale(Vector2[T]{scaleX, scaleY})
	matrix = matrix.Translate(Vector2[T]{-r.Left, -r.Top})
	return matrix
}

// String returns a string representation of the rectangle in the format "(TopLeft => BottomRight)".
func (r Rect[T]) String() string {
	return "(" + r.TopLeft().String() + " => " + r.BottomRight().String() + ")"
}

// Go converts the rectangle to a Go standard library image.Rectangle.
func (r Rect[T]) Go() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: int(r.Left), Y: int(r.Top)},
		Max: image.Point{X: int(r.Right), Y: int(r.Bottom)},
	}
}
