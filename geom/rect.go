package geom

import "image"

// Rect represents an axis-aligned rectangle defined by four edges or by origin and size.
// All methods are immutable and return new values.
type Rect[T Scalar] struct {
	Left   T
	Top    T
	Right  T
	Bottom T
}

// NewRect returns a rectangle from left, top, right, bottom.
func NewRect[T Scalar](left, top, right, bottom T) Rect[T] {
	return Rect[T]{Left: left, Top: top, Right: right, Bottom: bottom}
}

// NewRectXYWH returns a rectangle from x, y, width, height.
// Negative width/height will flip the rectangle accordingly.
func NewRectXYWH[T Scalar](x, y, width, height T) Rect[T] {
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

// NewRectOriginSize returns a rectangle from an origin point and a size.
// Negative size will flip the rectangle accordingly.
func NewRectOriginSize[T Scalar](origin Point[T], size Size[T]) Rect[T] {
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

// NewRectSize returns a rectangle at (0,0) with the given size.
func NewRectSize[T Scalar](size Size[T]) Rect[T] {
	return Rect[T]{Left: 0, Top: 0, Right: size.Width, Bottom: size.Height}
}

// NewRectFromGo converts a Go image.Rectangle to a Rect.
// The Go rectangle is inclusive on the left and top, exclusive on the right and bottom.
func NewRectFromGo[T Scalar](r image.Rectangle) Rect[T] {
	return Rect[T]{
		Left:   T(r.Min.X),
		Top:    T(r.Min.Y),
		Right:  T(r.Max.X),
		Bottom: T(r.Max.Y),
	}
}

// BoundingRect returns the minimal bounding rectangle for a set of points.
// If no points are given, returns an empty rectangle at (0,0).
func BoundingRect[T Scalar](points ...Point[T]) Rect[T] {
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

// X returns the x coordinate of the origin (left edge).
func (r Rect[T]) X() T {
	return r.Left
}

// Y returns the y coordinate of the origin (top edge).
func (r Rect[T]) Y() T {
	return r.Top
}

// Width returns the width (right - left).
func (r Rect[T]) Width() T {
	return r.Right - r.Left
}

// Height returns the height (bottom - top).
func (r Rect[T]) Height() T {
	return r.Bottom - r.Top
}

// LeftTop returns the top-left corner.
func (r Rect[T]) LeftTop() Point[T] {
	return Point[T]{r.Left, r.Top}
}

// RightTop returns the top-right corner.
func (r Rect[T]) RightTop() Point[T] {
	return Point[T]{r.Right, r.Top}
}

// LeftBottom returns the bottom-left corner.
func (r Rect[T]) LeftBottom() Point[T] {
	return Point[T]{r.Left, r.Bottom}
}

// RightBottom returns the bottom-right corner.
func (r Rect[T]) RightBottom() Point[T] {
	return Point[T]{r.Right, r.Bottom}
}

// LTRB returns (left, top, right, bottom).
func (r Rect[T]) LTRB() (T, T, T, T) {
	return r.Left, r.Top, r.Right, r.Bottom
}

// XYWH returns (x, y, width, height).
func (r Rect[T]) XYWH() (T, T, T, T) {
	return r.Left, r.Top, r.Width(), r.Height()
}

// Origin returns the origin point (left, top).
func (r Rect[T]) Origin() Point[T] {
	return Point[T]{r.Left, r.Top}
}

// Size returns the size (width, height).
func (r Rect[T]) Size() Size[T] {
	return Size[T]{r.Width(), r.Height()}
}

// Area returns the area (width * height).
func (r Rect[T]) Area() T {
	return r.Width() * r.Height()
}

// Center returns the center point.
func (r Rect[T]) Center() Point[T] {
	return Point[T]{(r.Left + r.Right) / T(2), (r.Top + r.Bottom) / T(2)}
}

// Positive returns a rectangle with positive width and height.
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

// Points returns the four corners: left-top, right-top, left-bottom, right-bottom.
func (r Rect[T]) Points() [4]Point[T] {
	return [4]Point[T]{
		r.LeftTop(),
		r.RightTop(),
		r.LeftBottom(),
		r.RightBottom(),
	}
}

// Eq reports whether two rectangles are Eq.
func (r Rect[T]) Eq(other Rect[T]) bool {
	return ScalarEq(r.Left, other.Left) &&
		ScalarEq(r.Top, other.Top) &&
		ScalarEq(r.Right, other.Right) &&
		ScalarEq(r.Bottom, other.Bottom)
}

// IsEmpty reports whether the rectangle is empty (width or height <= 0).
func (r Rect[T]) IsEmpty() bool {
	return r.Width() <= 0 || r.Height() <= 0
}

// IsFinite reports whether all edges are finite (float types only).
func (r Rect[T]) IsFinite() bool {
	return IsFinite(r.Left) && IsFinite(r.Top) && IsFinite(r.Right) && IsFinite(r.Bottom)
}

// IsSquare reports whether the rectangle is a square (width == height).
func (r Rect[T]) IsSquare() bool {
	return ScalarEq(r.Width(), r.Height()) && !r.IsEmpty()
}

// ContainsExclusive reports whether the rectangle contains a point (excluding edges).
func (r Rect[T]) ContainsExclusive(p Point[T]) bool {
	if r.IsEmpty() {
		return false
	}
	return p.X > r.Left && p.X < r.Right &&
		p.Y > r.Top && p.Y < r.Bottom
}

// Contains reports whether the rectangle contains a point (including edges).
func (r Rect[T]) Contains(p Point[T]) bool {
	if r.IsEmpty() {
		return false
	}
	return p.X >= r.Left && p.X <= r.Right &&
		p.Y >= r.Top && p.Y <= r.Bottom
}

// ContainsRect reports whether the rectangle contains another rectangle.
func (r Rect[T]) ContainsRect(other Rect[T]) bool {
	if r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return other.Left >= r.Left && other.Right <= r.Right &&
		other.Top >= r.Top && other.Bottom <= r.Bottom
}

// Intersects reports whether two rectangles intersect.
func (r Rect[T]) Intersects(other Rect[T]) bool {
	if r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return r.Left < other.Right && r.Right > other.Left &&
		r.Top < other.Bottom && r.Bottom > other.Top
}

// Intersect returns the intersection of two rectangles.
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

// Union returns the union of two rectangles.
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

// Cutout subtracts another rectangle from this rectangle.
// For simplicity, returns the original rect if partial overlap.
func (r Rect[T]) Cutout(other Rect[T]) (Rect[T], bool) {
	if !r.Intersects(other) {
		return r, true
	}
	if other.ContainsRect(r) {
		return NewRect[T](0, 0, 0, 0), false
	}
	return r, true
}

// Scale scales the rectangle by a scalar.
func (r Rect[T]) Scale(scalar T) Rect[T] {
	return NewRect(r.Left*scalar, r.Top*scalar, r.Right*scalar, r.Bottom*scalar)
}

// ScaleXY scales the rectangle by sx and sy.
func (r Rect[T]) ScaleXY(sx, sy T) Rect[T] {
	return NewRect(r.Left*sx, r.Top*sy, r.Right*sx, r.Bottom*sy)
}

// ScalePoint scales the rectangle to a given point.
func (r Rect[T]) ScalePoint(p Point[T]) Rect[T] {
	return r.ScaleXY(p.X, p.Y)
}

// ScaleSize scales the rectangle to a given size.
func (r Rect[T]) ScaleSize(s Size[T]) Rect[T] {
	return r.ScaleXY(s.Width, s.Height)
}

// Translate moves the rectangle by a vector.
func (r Rect[T]) Translate(vector Vector2[T]) Rect[T] {
	return NewRect(r.Left+vector.X, r.Top+vector.Y, r.Right+vector.X, r.Bottom+vector.Y)
}

// TranslateXY moves the rectangle by x and y.
func (r Rect[T]) TranslateXY(x, y T) Rect[T] {
	return NewRect(r.Left+x, r.Top+y, r.Right+x, r.Bottom+y)
}

// Expand expands the rectangle by the given amount.
func (r Rect[T]) Expand(amount T) Rect[T] {
	return NewRect(r.Left-amount, r.Top-amount, r.Right+amount, r.Bottom+amount)
}

// ExpandLTRB expands the rectangle by left, top, right, and bottom.
func (r Rect[T]) ExpandLTRB(left, top, right, bottom T) Rect[T] {
	return NewRect(r.Left-left, r.Top-top, r.Right+right, r.Bottom+bottom)
}

// ExpandHV expands the rectangle by horizontal and vertical amounts.
func (r Rect[T]) ExpandHV(horizontal, vertical T) Rect[T] {
	return NewRect(r.Left-horizontal, r.Top-vertical, r.Right+horizontal, r.Bottom+vertical)
}

// ExpandPoint expands the rectangle to include a point.
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

// ExpandSize expands the rectangle by a size.
func (r Rect[T]) ExpandSize(s Size[T]) Rect[T] {
	return r.ExpandHV(s.Width/T(2), s.Height/T(2))
}

// Project projects the source rectangle into the coordinate space of this rectangle.
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

// Round rounds the rectangle's edges to the nearest integer.
func (r Rect[T]) Round() Rect[I32] {
	return NewRect(
		I32(r.Left.Float64()+0.5),
		I32(r.Top.Float64()+0.5),
		I32(r.Right.Float64()+0.5),
		I32(r.Bottom.Float64()+0.5),
	)
}

// RoundOut rounds the edges outward to the nearest integer.
func (r Rect[T]) RoundOut() Rect[I32] {
	return NewRect(
		I32(r.Left.Float64()),
		I32(r.Top.Float64()),
		I32(r.Right.Float64()+0.999),
		I32(r.Bottom.Float64()+0.999),
	)
}

// RoundIn rounds the edges inward (down) to the nearest integer.
func (r Rect[T]) RoundIn() Rect[I32] {
	return NewRect(
		I32(r.Left.Float64()+0.999),
		I32(r.Top.Float64()+0.999),
		I32(r.Right.Float64()),
		I32(r.Bottom.Float64()),
	)
}

// Transform applies a matrix to the four corners.
func (r Rect[T]) Transform(transform Matrix[T]) [4]Point[T] {
	corners := r.Points()
	return [4]Point[T]{
		transform.transformPoint(corners[0]),
		transform.transformPoint(corners[1]),
		transform.transformPoint(corners[2]),
		transform.transformPoint(corners[3]),
	}
}

// TransformBounds applies an affine transformation and returns the bounding box.
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

// TransformClipBounds applies a perspective transformation, clips to bounds, and returns the bounding box.
func (r Rect[T]) TransformClipBounds(transform Matrix[T], bounds Rect[T]) Rect[T] {
	transformed := r.TransformBounds(transform)
	return transformed.Intersect(bounds)
}

// NormalizingTransform returns a matrix that normalizes the rectangle to [0,1].
func (r Rect[T]) NormalizingTransform() Matrix[T] {
	if r.IsEmpty() {
		return Matrix[T]{}
	}
	scaleX := T(1) / r.Width()
	scaleY := T(1) / r.Height()
	matrix := Matrix[T]{}
	matrix = matrix.Scale(Vector2[T]{scaleX, scaleY})
	matrix = matrix.Translate(Vector2[T]{-r.Left, -r.Top})
	return matrix
}

// String returns a string representation, e.g. (LeftTop => RightBottom).
func (r Rect[T]) String() string {
	return "(" + r.LeftTop().String() + " => " + r.RightBottom().String() + ")"
}

// ToGo converts the rectangle to a Go image.Rectangle.
func (r Rect[T]) ToGo() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: int(r.Left), Y: int(r.Top)},
		Max: image.Point{X: int(r.Right), Y: int(r.Bottom)},
	}
}
