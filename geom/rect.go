package geom

// Rect is a generic interface for axis-aligned rectangles.
// It provides methods for rectangle geometry, queries, and transformations.
// Rectangles are defined by four axis-aligned edges or by origin and size.
// An empty rectangle has width or height <= 0.
// All methods are immutable and return new values.
type Rect[T Scalar] interface {
	// === Basic Properties ===

	// Left returns the left edge.
	Left() T
	// Right returns the right edge.
	Right() T
	// Top returns the top edge.
	Top() T
	// Bottom returns the bottom edge.
	Bottom() T

	// X returns the x coordinate of the origin (left edge).
	X() T
	// Y returns the y coordinate of the origin (top edge).
	Y() T
	// Width returns the width (right - left).
	Width() T
	// Height returns the height (bottom - top).
	Height() T

	// === Geometric Queries ===

	// LeftTop returns the top-left corner.
	LeftTop() Point[T]
	// RightTop returns the top-right corner.
	RightTop() Point[T]
	// LeftBottom returns the bottom-left corner.
	LeftBottom() Point[T]
	// RightBottom returns the bottom-right corner.
	RightBottom() Point[T]

	// LTRB returns (left, top, right, bottom).
	LTRB() (T, T, T, T)
	// XYWH returns (x, y, width, height).
	XYWH() (T, T, T, T)

	// Origin returns the origin point (left, top).
	Origin() Point[T]
	// Size returns the size (width, height).
	Size() Size[T]
	// Area returns the area (width * height).
	Area() T
	// Center returns the center point.
	Center() Point[T]
	// Positive returns a rectangle with positive width and height.
	Positive() Rect[T]

	// Points returns the four corners: left-top, right-top, left-bottom, right-bottom.
	Points() [4]Point[T]

	// === State Checks ===

	// Equal reports whether two rectangles are equal.
	Equal(other Rect[T]) bool
	// IsEmpty reports whether the rectangle is empty (width or height <= 0).
	IsEmpty() bool
	// IsFinite reports whether all edges are finite (float types only).
	IsFinite() bool
	// IsSquare reports whether the rectangle is a square (width == height).
	IsSquare() bool
	// IsMaximum reports whether the rectangle covers all finite coordinates.
	IsMaximum() bool

	// === Spatial Relationships ===

	// ContainsExclusive reports whether the rectangle contains a point (excluding edges).
	ContainsExclusive(p Point[T]) bool
	// Contains reports whether the rectangle contains a point (including edges).
	Contains(p Point[T]) bool
	// ContainsRect reports whether the rectangle contains another rectangle.
	ContainsRect(r Rect[T]) bool
	// Intersects reports whether two rectangles intersect.
	Intersects(r Rect[T]) bool

	// === Set Operations ===

	// Intersection returns the intersection of two rectangles.
	Intersection(r Rect[T]) (Rect[T], bool)
	// IntersectOrEmpty returns the intersection or an empty rectangle.
	IntersectOrEmpty(r Rect[T]) Rect[T]
	// Union returns the union of two rectangles.
	Union(r Rect[T]) Rect[T]
	// Cutout subtracts another rectangle from this rectangle.
	Cutout(r Rect[T]) (Rect[T], bool)

	// === Transformations ===

	// Scale scales the rectangle by a scalar.
	Scale(scalar T) Rect[T]
	// ScaleXY scales the rectangle by sx and sy.
	ScaleXY(sx, sy T) Rect[T]
	// ScalePoint scales the rectangle to a given point.
	ScalePoint(p Point[T]) Rect[T]
	// ScaleSize scales the rectangle to a given size.
	ScaleSize(s Size[T]) Rect[T]
	// Translate moves the rectangle by a vector.
	Translate(vector Vector2[T]) Rect[T]
	// TranslateXY moves the rectangle by x and y.
	TranslateXY(x, y T) Rect[T]
	// Expand expands the rectangle by the given amount.
	Expand(amount T) Rect[T]
	// ExpandLTRB expands the rectangle by left, top, right, and bottom.
	ExpandLTRB(left, top, right, bottom T) Rect[T]
	// ExpandHV expands the rectangle by horizontal and vertical amounts.
	ExpandHV(horizontal, vertical T) Rect[T]
	// ExpandPoint expands the rectangle to include a point.
	ExpandPoint(p Point[T]) Rect[T]
	// ExpandSize expands the rectangle by a size.
	ExpandSize(s Size[T]) Rect[T]
	// Project projects the source rectangle into the coordinate space of this rectangle.
	Project(source Rect[T]) Rect[T]

	// Round rounds the rectangle's edges to the nearest integer.
	Round() Rect[Int32]
	// RoundOut rounds the edges outward to the nearest integer.
	RoundOut() Rect[Int32]
	// RoundIn rounds the edges inward (down) to the nearest integer.
	RoundIn() Rect[Int32]

	// === Matrix Transformations ===

	// Transform applies a matrix to the four corners.
	Transform(transform Matrix[T]) [4]Point[T]
	// TransformBounds applies an affine transformation and returns the bounding box.
	TransformBounds(transform Matrix[T]) Rect[T]
	// TransformClipBounds applies a perspective transformation, clips to bounds, and returns the bounding box.
	TransformClipBounds(transform Matrix[T], bounds Rect[T]) Rect[T]
	// NormalizingTransform returns a matrix that normalizes the rectangle to [0,1].
	NormalizingTransform() Matrix[T]

	// String returns a string representation, e.g. (LeftTop => RightBottom).
	String() string
}

// NewRect returns a rectangle from left, top, right, bottom.
func NewRect[T Scalar](left, top, right, bottom T) Rect[T] {
	return &rect[T]{left, top, right, bottom}
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
	return &rect[T]{x, y, x + width, y + height}
}

// NewRectOriginSize returns a rectangle from an origin point and a size.
// Negative size will flip the rectangle accordingly.
func NewRectOriginSize[T Scalar](origin Point[T], size Size[T]) Rect[T] {
	if size.Width() < 0 {
		origin = NewPoint(origin.X()+size.Width(), origin.Y())
	}
	if size.Height() < 0 {
		origin = NewPoint(origin.X(), origin.Y()+size.Height())
	}
	return &rect[T]{origin.X(), origin.Y(), origin.X() + size.Width(), origin.Y() + size.Height()}
}

// NewRectSize returns a rectangle at (0,0) with the given size.
func NewRectSize[T Scalar](size Size[T]) Rect[T] {
	return &rect[T]{0, 0, size.Width(), size.Height()}
}

// NewRectMax returns a rectangle covering all finite coordinates.
func NewRectMax[T Scalar]() Rect[T] {
	maxVal := Max[T]()
	return NewRect(-maxVal, -maxVal, maxVal, maxVal)
}

// BoundingRect returns the minimal bounding rectangle for a set of points.
// If no points are given, returns an empty rectangle at (0,0).
func BoundingRect[T Scalar](points ...Point[T]) Rect[T] {
	if len(points) == 0 {
		return NewRect[T](0, 0, 0, 0)
	}

	left := points[0].X()
	top := points[0].Y()
	right := left
	bottom := top

	for _, p := range points {
		if p.X() < left {
			left = p.X()
		}
		if p.X() > right {
			right = p.X()
		}
		if p.Y() < top {
			top = p.Y()
		}
		if p.Y() > bottom {
			bottom = p.Y()
		}
	}

	return NewRect(left, top, right, bottom)
}

type rect[T Scalar] struct {
	left, top, right, bottom T
}

// === Basic Properties ===

// Left implements Rect.
func (r *rect[T]) Left() T {
	return r.left
}

// Right implements Rect.
func (r *rect[T]) Right() T {
	return r.right
}

// Top implements Rect.
func (r *rect[T]) Top() T {
	return r.top
}

// Bottom implements Rect.
func (r *rect[T]) Bottom() T {
	return r.bottom
}

// X implements Rect.
func (r *rect[T]) X() T {
	return r.left
}

// Y implements Rect.
func (r *rect[T]) Y() T {
	return r.top
}

// Width implements Rect.
func (r *rect[T]) Width() T {
	return r.right - r.left
}

// Height implements Rect.
func (r *rect[T]) Height() T {
	return r.bottom - r.top
}

// === Geometric Queries ===

// LeftTop implements Rect.
func (r *rect[T]) LeftTop() Point[T] {
	return NewPoint(r.left, r.top)
}

// RightTop implements Rect.
func (r *rect[T]) RightTop() Point[T] {
	return NewPoint(r.right, r.top)
}

// LeftBottom implements Rect.
func (r *rect[T]) LeftBottom() Point[T] {
	return NewPoint(r.left, r.bottom)
}

// RightBottom implements Rect.
func (r *rect[T]) RightBottom() Point[T] {
	return NewPoint(r.right, r.bottom)
}

// LTRB implements Rect.
func (r *rect[T]) LTRB() (T, T, T, T) {
	return r.left, r.top, r.right, r.bottom
}

// XYWH implements Rect.
func (r *rect[T]) XYWH() (T, T, T, T) {
	return r.left, r.top, r.Width(), r.Height()
}

// Origin implements Rect.
func (r *rect[T]) Origin() Point[T] {
	return NewPoint(r.left, r.top)
}

// Size implements Rect.
func (r *rect[T]) Size() Size[T] {
	return NewSize(r.Width(), r.Height())
}

// Area implements Rect.
func (r *rect[T]) Area() T {
	return r.Width() * r.Height()
}

// Center implements Rect.
func (r *rect[T]) Center() Point[T] {
	return NewPoint((r.left+r.right)/T(2), (r.top+r.bottom)/T(2))
}

// Positive implements Rect.
func (r *rect[T]) Positive() Rect[T] {
	left, right := r.left, r.right
	top, bottom := r.top, r.bottom

	if left > right {
		left, right = right, left
	}
	if top > bottom {
		top, bottom = bottom, top
	}

	return NewRect(left, top, right, bottom)
}

// Points implements Rect.
func (r *rect[T]) Points() [4]Point[T] {
	return [4]Point[T]{
		r.LeftTop(),
		r.RightTop(),
		r.LeftBottom(),
		r.RightBottom(),
	}
}

// === State Checks ===

// Equal implements Rect.
func (r *rect[T]) Equal(other Rect[T]) bool {
	return Equal(r.left, other.Left()) &&
		Equal(r.top, other.Top()) &&
		Equal(r.right, other.Right()) &&
		Equal(r.bottom, other.Bottom())
}

// IsEmpty implements Rect.
func (r *rect[T]) IsEmpty() bool {
	return r.Width() <= 0 || r.Height() <= 0
}

// IsFinite implements Rect.
func (r *rect[T]) IsFinite() bool {
	return IsFinite(r.left) && IsFinite(r.top) && IsFinite(r.right) && IsFinite(r.bottom)
}

// IsSquare implements Rect.
func (r *rect[T]) IsSquare() bool {
	return Equal(r.Width(), r.Height()) && !r.IsEmpty()
}

// IsMaximum implements Rect.
func (r *rect[T]) IsMaximum() bool {
	maxVal := Max[T]()
	return Equal(r.left, -maxVal) && Equal(r.top, -maxVal) &&
		Equal(r.right, maxVal) && Equal(r.bottom, maxVal)
}

// === Spatial Relationships ===

// Contains implements Rect.
func (r *rect[T]) Contains(p Point[T]) bool {
	if r.IsEmpty() {
		return false
	}

	return p.X() >= r.left && p.X() <= r.right &&
		p.Y() >= r.top && p.Y() <= r.bottom
}

// ContainsExclusive implements Rect.
func (r *rect[T]) ContainsExclusive(p Point[T]) bool {
	if r.IsEmpty() {
		return false
	}
	return p.X() > r.left && p.X() < r.right &&
		p.Y() > r.top && p.Y() < r.bottom
}

// ContainsRect implements Rect.
func (r *rect[T]) ContainsRect(other Rect[T]) bool {
	if r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return other.Left() >= r.left && other.Right() <= r.right &&
		other.Top() >= r.top && other.Bottom() <= r.bottom
}

// Intersects implements Rect.
func (r *rect[T]) Intersects(other Rect[T]) bool {
	if r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return r.left < other.Right() && r.right > other.Left() &&
		r.top < other.Bottom() && r.bottom > other.Top()
}

// === Set Operations ===

// Intersection implements Rect.
func (r *rect[T]) Intersection(other Rect[T]) (Rect[T], bool) {
	if !r.Intersects(other) {
		return NewRect[T](0, 0, 0, 0), false
	}

	left := r.left
	if other.Left() > left {
		left = other.Left()
	}
	top := r.top
	if other.Top() > top {
		top = other.Top()
	}
	right := r.right
	if other.Right() < right {
		right = other.Right()
	}
	bottom := r.bottom
	if other.Bottom() < bottom {
		bottom = other.Bottom()
	}

	result := NewRect(left, top, right, bottom)
	return result, !result.IsEmpty()
}

// IntersectOrEmpty implements Rect.
func (r *rect[T]) IntersectOrEmpty(other Rect[T]) Rect[T] {
	result, ok := r.Intersection(other)
	if !ok {
		return NewRect[T](0, 0, 0, 0)
	}
	return result
}

// Union implements Rect.
func (r *rect[T]) Union(other Rect[T]) Rect[T] {
	if r.IsEmpty() {
		return other
	}
	if other.IsEmpty() {
		return r
	}

	left := r.left
	if other.Left() < left {
		left = other.Left()
	}
	top := r.top
	if other.Top() < top {
		top = other.Top()
	}
	right := r.right
	if other.Right() > right {
		right = other.Right()
	}
	bottom := r.bottom
	if other.Bottom() > bottom {
		bottom = other.Bottom()
	}

	return NewRect(left, top, right, bottom)
}

// Cutout implements Rect.
func (r *rect[T]) Cutout(other Rect[T]) (Rect[T], bool) {
	if !r.Intersects(other) {
		return r, true
	}

	// Simple case: if other completely contains this rect, result is empty
	if other.ContainsRect(r) {
		return NewRect[T](0, 0, 0, 0), false
	}

	// For simplicity, return the original rect if partial overlap
	// A full implementation would return multiple rects
	return r, true
}

// === Transformations ===

// Scale implements Rect.
func (r *rect[T]) Scale(scalar T) Rect[T] {
	return NewRect(r.left*scalar, r.top*scalar, r.right*scalar, r.bottom*scalar)
}

// ScaleXY implements Rect.
func (r *rect[T]) ScaleXY(sx, sy T) Rect[T] {
	return NewRect(r.left*sx, r.top*sy, r.right*sx, r.bottom*sy)
}

// ScalePoint implements Rect.
func (r *rect[T]) ScalePoint(p Point[T]) Rect[T] {
	return r.ScaleXY(p.X(), p.Y())
}

// ScaleSize implements Rect.
func (r *rect[T]) ScaleSize(s Size[T]) Rect[T] {
	return r.ScaleXY(s.Width(), s.Height())
}

// Translate implements Rect.
func (r *rect[T]) Translate(vector Vector2[T]) Rect[T] {
	return NewRect(r.left+vector.X(), r.top+vector.Y(), r.right+vector.X(), r.bottom+vector.Y())
}

// TranslateXY implements Rect.
func (r *rect[T]) TranslateXY(x, y T) Rect[T] {
	return NewRect(r.left+x, r.top+y, r.right+x, r.bottom+y)
}

// Expand implements Rect.
func (r *rect[T]) Expand(amount T) Rect[T] {
	return NewRect(r.left-amount, r.top-amount, r.right+amount, r.bottom+amount)
}

// ExpandLTRB implements Rect.
func (r *rect[T]) ExpandLTRB(left, top, right, bottom T) Rect[T] {
	return NewRect(r.left-left, r.top-top, r.right+right, r.bottom+bottom)
}

// ExpandHV implements Rect.
func (r *rect[T]) ExpandHV(horizontal, vertical T) Rect[T] {
	return NewRect(r.left-horizontal, r.top-vertical, r.right+horizontal, r.bottom+vertical)
}

// ExpandPoint implements Rect.
func (r *rect[T]) ExpandPoint(p Point[T]) Rect[T] {
	left := r.left
	if p.X() < left {
		left = p.X()
	}
	top := r.top
	if p.Y() < top {
		top = p.Y()
	}
	right := r.right
	if p.X() > right {
		right = p.X()
	}
	bottom := r.bottom
	if p.Y() > bottom {
		bottom = p.Y()
	}

	return NewRect(left, top, right, bottom)
}

// ExpandSize implements Rect.
func (r *rect[T]) ExpandSize(s Size[T]) Rect[T] {
	return r.ExpandHV(s.Width()/T(2), s.Height()/T(2))
}

// Project implements Rect.
func (r *rect[T]) Project(source Rect[T]) Rect[T] {
	if r.IsEmpty() || source.IsEmpty() {
		return NewRect[T](0, 0, 0, 0)
	}

	scaleX := r.Width() / source.Width()
	scaleY := r.Height() / source.Height()

	return NewRect(
		r.left+scaleX*(source.Left()-source.Left()),
		r.top+scaleY*(source.Top()-source.Top()),
		r.left+scaleX*(source.Right()-source.Left()),
		r.top+scaleY*(source.Bottom()-source.Top()),
	)
}

// Round implements Rect.
func (r *rect[T]) Round() Rect[Int32] {
	return NewRect(
		Int32(r.left.Float64()+0.5),
		Int32(r.top.Float64()+0.5),
		Int32(r.right.Float64()+0.5),
		Int32(r.bottom.Float64()+0.5),
	)
}

// RoundOut implements Rect.
func (r *rect[T]) RoundOut() Rect[Int32] {
	return NewRect(
		Int32(r.left.Float64()),
		Int32(r.top.Float64()),
		Int32(r.right.Float64()+0.999),
		Int32(r.bottom.Float64()+0.999),
	)
}

// RoundIn implements Rect.
func (r *rect[T]) RoundIn() Rect[Int32] {
	return NewRect(
		Int32(r.left.Float64()+0.999),
		Int32(r.top.Float64()+0.999),
		Int32(r.right.Float64()),
		Int32(r.bottom.Float64()),
	)
}

// === Matrix Transformations ===

// Transform implements Rect.
func (r *rect[T]) Transform(transform Matrix[T]) [4]Point[T] {
	corners := r.Points()
	return [4]Point[T]{
		transform.TransformPoint(corners[0]),
		transform.TransformPoint(corners[1]),
		transform.TransformPoint(corners[2]),
		transform.TransformPoint(corners[3]),
	}
}

// TransformBounds implements Rect.
func (r *rect[T]) TransformBounds(transform Matrix[T]) Rect[T] {
	if r.IsEmpty() {
		return r
	}

	transformed := r.Transform(transform)

	left := transformed[0].X()
	top := transformed[0].Y()
	right := left
	bottom := top

	for i := 1; i < 4; i++ {
		x, y := transformed[i].X(), transformed[i].Y()
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

// TransformClipBounds implements Rect.
func (r *rect[T]) TransformClipBounds(transform Matrix[T], bounds Rect[T]) Rect[T] {
	transformed := r.TransformBounds(transform)
	return transformed.IntersectOrEmpty(bounds)
}

// NormalizingTransform implements Rect.
func (r *rect[T]) NormalizingTransform() Matrix[T] {
	if r.IsEmpty() {
		return NewMatrix[T]()
	}

	scaleX := T(1) / r.Width()
	scaleY := T(1) / r.Height()

	matrix := NewMatrix[T]()
	matrix = matrix.Scale2D(NewVector2(scaleX, scaleY))
	matrix = matrix.Translate2D(NewVector2(-r.left, -r.top))

	return matrix
}

// String implements Rect.
func (r *rect[T]) String() string {
	return "(" + r.LeftTop().String() + " => " + r.RightBottom().String() + ")"
}
