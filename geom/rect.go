package geom

import (
	"fmt"
	"image"
	"math"
)

// Rect defines the interface for an axis-aligned rectangle.
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

	// TopLeft returns the top-left corner point.
	TopLeft() Point
	// TopRight returns the top-right corner point.
	TopRight() Point
	// BottomLeft returns the bottom-left corner point.
	BottomLeft() Point
	// BottomRight returns the bottom-right corner point.
	BottomRight() Point

	// LTRB returns the rectangle's coordinates as (left, top, right, bottom).
	LTRB() (Scalar, Scalar, Scalar, Scalar)
	// XYWH returns the rectangle's coordinates as (x, y, width, height).
	XYWH() (Scalar, Scalar, Scalar, Scalar)

	// Origin returns the origin point (top-left corner).
	Origin() Point
	// Size returns the dimensions as (width, height).
	Size() Size
	// Area returns the area (width × height).
	Area() Scalar
	// Center returns the center point.
	Center() Point

	// Positive returns a rectangle with positive dimensions, normalizing if needed.
	Positive() Rect
	// Points returns the four corner points in order: top-left, top-right, bottom-left, bottom-right.
	Points() Quad

	// Equal reports whether two rectangles are equal within floating-point tolerance.
	Equal(o Rect) bool
	// Round returns a rectangle with coordinates rounded to nearest integer.
	Round() Rect
	// Floor returns a rectangle with coordinates rounded down.
	Floor() Rect
	// Ceil returns a rectangle with coordinates rounded up.
	Ceil() Rect

	// IsEmpty reports whether the rectangle has zero or negative dimensions.
	IsEmpty() bool
	// IsFinite reports whether all coordinates are finite numbers.
	IsFinite() bool
	// IsSquare reports whether the rectangle is a square with positive dimensions.
	IsSquare() bool
	// IsMaximum reports whether this is the maximum possible rectangle.
	IsMaximum() bool

	// Contains reports whether the rectangle contains the point (including edges).
	Contains(p Point) bool
	// ContainsInclusive reports whether the rectangle contains the point (including all edges).
	ContainsInclusive(p Point) bool
	// Inside reports whether the rectangle strictly contains the point (excluding edges).
	Inside(p Point) bool
	// ContainsRect reports whether this rectangle completely contains another.
	ContainsRect(o Rect) bool
	// Intersects reports whether this rectangle intersects with another.
	Intersects(o Rect) bool

	// Intersect returns the intersection with another rectangle.
	Intersect(o Rect) Rect
	// Union returns the smallest rectangle containing both rectangles.
	Union(o Rect) Rect
	// Cutout attempts to subtract another rectangle from this one.
	Cutout(o Rect) (Rect, bool)

	// Scale returns a rectangle scaled by the given factor.
	Scale(scalar Scalar) Rect
	// ScaleXY returns a rectangle scaled by different factors for X and Y axes.
	ScaleXY(sx, sy Scalar) Rect
	// ScalePoint returns a rectangle scaled by the point's components.
	ScalePoint(p Point) Rect
	// ScaleSize returns a rectangle scaled by the size's dimensions.
	ScaleSize(s Size) Rect

	// Translate returns a rectangle moved by the given vector.
	Translate(vector Vector2) Rect
	// TranslateXY returns a rectangle moved by the specified offsets.
	TranslateXY(x, y Scalar) Rect

	// Expand returns a rectangle expanded by the amount in all directions.
	Expand(amount Scalar) Rect
	// ExpandLTRB returns a rectangle expanded by different amounts on each edge.
	ExpandLTRB(left, top, right, bottom Scalar) Rect
	// ExpandHV returns a rectangle expanded horizontally and vertically.
	ExpandHV(horizontal, vertical Scalar) Rect
	// ExpandPoint returns a rectangle expanded to include the point.
	ExpandPoint(p Point) Rect
	// ExpandSize returns a rectangle expanded by half the size in all directions.
	ExpandSize(s Size) Rect

	// Project maps the source rectangle into this rectangle's coordinate space.
	Project(source Rect) Rect
	// Transform applies a transformation matrix to the four corners.
	Transform(transform Matrix) Quad
	// TransformBounds applies a transformation and returns the axis-aligned bounding box.
	TransformBounds(transform Matrix) Rect
	// TransformClipBounds applies a transformation, clips to bounds, and returns the bounding box.
	TransformClipBounds(transform Matrix, bounds Rect) Rect

	// NormalizingTransform returns a matrix that maps this rectangle to the unit square [0,1]×[0,1].
	NormalizingTransform() Matrix

	// String returns a string representation.
	String() string
	// Go converts to a Go standard library image.Rectangle.
	Go() image.Rectangle
}

// NewRect creates a rectangle from left, top, right, and bottom coordinates.
func NewRect[T Number](left, top, right, bottom T) Rect {
	return rect[T]{left: left, top: top, right: right, bottom: bottom}
}

// NewRectXYWH creates a rectangle from x, y coordinates and width, height dimensions.
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

// rect is the generic implementation of the Rect interface.
type rect[T Number] struct {
	left, top, right, bottom T
}

// Left implements Rect.Left.
func (r rect[T]) Left() Scalar { return ToScalar(r.left) }

// Top implements Rect.Top.
func (r rect[T]) Top() Scalar { return ToScalar(r.top) }

// Right implements Rect.Right.
func (r rect[T]) Right() Scalar { return ToScalar(r.right) }

// Bottom implements Rect.Bottom.
func (r rect[T]) Bottom() Scalar { return ToScalar(r.bottom) }

// X implements Rect.X.
func (r rect[T]) X() Scalar { return r.Left() }

// Y implements Rect.Y.
func (r rect[T]) Y() Scalar { return r.Top() }

// Width implements Rect.Width.
func (r rect[T]) Width() Scalar { return r.Right() - r.Left() }

// Height implements Rect.Height.
func (r rect[T]) Height() Scalar { return r.Bottom() - r.Top() }

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
	return r.Left(), r.Top(), r.Right(), r.Bottom()
}

// XYWH implements Rect.XYWH.
func (r rect[T]) XYWH() (Scalar, Scalar, Scalar, Scalar) {
	return r.X(), r.Y(), r.Width(), r.Height()
}

// Origin implements Rect.Origin.
func (r rect[T]) Origin() Point {
	return NewPoint(r.left, r.top)
}

// Size implements Rect.Size.
func (r rect[T]) Size() Size {
	return NewSize(r.Width(), r.Height())
}

// Area implements Rect.Area.
func (r rect[T]) Area() Scalar {
	return r.Width() * r.Height()
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

	return NewRect(left, top, right, bottom)
}

// Points implements Rect.Points.
func (r rect[T]) Points() Quad {
	return Quad{r.TopLeft(), r.TopRight(), r.BottomLeft(), r.BottomRight()}
}

// Equal implements Rect.Equal.
func (r rect[T]) Equal(o Rect) bool {
	return NearlyEqual(r.Left(), o.Left()) &&
		NearlyEqual(r.Top(), o.Top()) &&
		NearlyEqual(r.Right(), o.Right()) &&
		NearlyEqual(r.Bottom(), o.Bottom())
}

// Round implements Rect.Round.
func (r rect[T]) Round() Rect {
	return NewRect(
		T(math.Round(ToFloat64(r.left))),
		T(math.Round(ToFloat64(r.top))),
		T(math.Round(ToFloat64(r.right))),
		T(math.Round(ToFloat64(r.bottom))),
	)
}

// Floor implements Rect.Floor.
func (r rect[T]) Floor() Rect {
	return NewRect(
		T(math.Floor(ToFloat64(r.left))),
		T(math.Floor(ToFloat64(r.top))),
		T(math.Floor(ToFloat64(r.right))),
		T(math.Floor(ToFloat64(r.bottom))),
	)
}

// Ceil implements Rect.Ceil.
func (r rect[T]) Ceil() Rect {
	return NewRect(
		T(math.Ceil(ToFloat64(r.left))),
		T(math.Ceil(ToFloat64(r.top))),
		T(math.Ceil(ToFloat64(r.right))),
		T(math.Ceil(ToFloat64(r.bottom))),
	)
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

// IsMaximum implements Rect.IsMaximum.
func (r rect[T]) IsMaximum() bool {
	// Check if this rectangle represents the maximum possible rectangle
	// by comparing against very large but representable values
	maxVal := Scalar(1e30) // Use a large but safe value for both float32 and float64
	return r.Left() <= -maxVal && r.Top() <= -maxVal &&
		r.Right() >= maxVal && r.Bottom() >= maxVal
}

// Contains implements Rect.Contains.
func (r rect[T]) Contains(p Point) bool {
	if r.IsEmpty() {
		return false
	}
	x, y := p.X(), p.Y()
	return r.Left() <= x && x < r.Right() && r.Top() <= y && y < r.Bottom()
}

// ContainsInclusive implements Rect.ContainsInclusive.
func (r rect[T]) ContainsInclusive(p Point) bool {
	if r.IsEmpty() {
		return false
	}
	x, y := p.X(), p.Y()
	return r.Left() <= x && x <= r.Right() && r.Top() <= y && y <= r.Bottom()
}

// Inside implements Rect.Inside.
func (r rect[T]) Inside(p Point) bool {
	x, y := p.X(), p.Y()
	return r.Left() < x && x < r.Right() && r.Top() < y && y < r.Bottom()
}

// ContainsRect implements Rect.ContainsRect.
func (r rect[T]) ContainsRect(o Rect) bool {
	if r.IsEmpty() {
		return false
	}
	if o.IsEmpty() {
		return true // An empty rectangle is contained within any non-empty rectangle
	}
	return r.Left() <= o.Left() &&
		r.Top() <= o.Top() &&
		r.Right() >= o.Right() &&
		r.Bottom() >= o.Bottom()
}

// Intersects implements Rect.Intersects.
func (r rect[T]) Intersects(o Rect) bool {
	if r.IsEmpty() || o.IsEmpty() {
		return false
	}
	return r.Left() < o.Right() &&
		r.Right() > o.Left() &&
		r.Top() < o.Bottom() &&
		r.Bottom() > o.Top()
}

// Intersect implements Rect.Intersect.
func (r rect[T]) Intersect(o Rect) Rect {
	if !r.Intersects(o) {
		return NewRect[T](0, 0, 0, 0) // Return empty rectangle
	}

	left := max(r.Left(), o.Left())
	top := max(r.Top(), o.Top())
	right := min(r.Right(), o.Right())
	bottom := min(r.Bottom(), o.Bottom())

	return NewRect(left, top, right, bottom)
}

// Union implements Rect.Union.
func (r rect[T]) Union(o Rect) Rect {
	if r.IsEmpty() {
		return NewRect(o.Left(), o.Top(), o.Right(), o.Bottom())
	}
	if o.IsEmpty() {
		return r
	}

	left := min(r.Left(), o.Left())
	top := min(r.Top(), o.Top())
	right := max(r.Right(), o.Right())
	bottom := max(r.Bottom(), o.Bottom())

	return NewRect(left, top, right, bottom)
}

// Cutout implements Rect.Cutout.
func (r rect[T]) Cutout(o Rect) (Rect, bool) {
	if r.IsEmpty() {
		// Empty rectangles cannot be cut out.
		return NewRect[T](0, 0, 0, 0), false
	}
	if o.IsEmpty() {
		// Cutting by empty rectangle returns original.
		return r, true
	}

	aLeft, aTop, aRight, aBottom := r.LTRB()
	bLeft, bTop, bRight, bBottom := o.LTRB()

	// Check if cutout completely contains the source rectangle
	if bLeft <= aLeft && bRight >= aRight && bTop <= aTop && bBottom >= aBottom {
		// Full cutout, return empty.
		return NewRect[T](0, 0, 0, 0), false
	}

	// Check for vertical cuts (cutout spans full width)
	if bLeft <= aLeft && bRight >= aRight {
		if bTop <= aTop && bBottom > aTop {
			// Cuts off the top.
			return NewRect(aLeft, bBottom, aRight, aBottom), true
		}
		if bBottom >= aBottom && bTop < aBottom {
			// Cuts off the bottom.
			return NewRect(aLeft, aTop, aRight, bTop), true
		}
	}

	// Check for horizontal cuts (cutout spans full height)
	if bTop <= aTop && bBottom >= aBottom {
		if bLeft <= aLeft && bRight > aLeft {
			// Cuts off the left.
			return NewRect(bRight, aTop, aRight, aBottom), true
		}
		if bRight >= aRight && bLeft < aRight {
			// Cuts off the right.
			return NewRect(aLeft, aTop, bLeft, aBottom), true
		}
	}

	// No valid cutout, return original rectangle.
	return r, true
}

// Scale implements Rect.Scale.
func (r rect[T]) Scale(scalar Scalar) Rect {
	return NewRect(
		r.Left()*scalar,
		r.Top()*scalar,
		r.Right()*scalar,
		r.Bottom()*scalar,
	)
}

// ScaleXY implements Rect.ScaleXY.
func (r rect[T]) ScaleXY(sx, sy Scalar) Rect {
	return NewRect(
		r.Left()*sx,
		r.Top()*sy,
		r.Right()*sx,
		r.Bottom()*sy,
	)
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
	return NewRect(
		r.Left()+vector.X(),
		r.Top()+vector.Y(),
		r.Right()+vector.X(),
		r.Bottom()+vector.Y(),
	)
}

// TranslateXY implements Rect.TranslateXY.
func (r rect[T]) TranslateXY(x, y Scalar) Rect {
	return NewRect(
		r.Left()+x,
		r.Top()+y,
		r.Right()+x,
		r.Bottom()+y,
	)
}

// Expand implements Rect.Expand.
func (r rect[T]) Expand(amount Scalar) Rect {
	return NewRect(
		r.Left()-amount,
		r.Top()-amount,
		r.Right()+amount,
		r.Bottom()+amount,
	)
}

// ExpandLTRB implements Rect.ExpandLTRB.
func (r rect[T]) ExpandLTRB(left, top, right, bottom Scalar) Rect {
	return NewRect(
		r.Left()-left,
		r.Top()-top,
		r.Right()+right,
		r.Bottom()+bottom,
	)
}

// ExpandHV implements Rect.ExpandHV.
func (r rect[T]) ExpandHV(horizontal, vertical Scalar) Rect {
	return NewRect(
		r.Left()-horizontal,
		r.Top()-vertical,
		r.Right()+horizontal,
		r.Bottom()+vertical,
	)
}

// ExpandPoint implements Rect.ExpandPoint.
func (r rect[T]) ExpandPoint(p Point) Rect {
	if r.IsEmpty() {
		x, y := p.X(), p.Y()
		return NewRect(x, y, x, y)
	}

	return NewRect(
		min(r.Left(), p.X()),
		min(r.Top(), p.Y()),
		max(r.Right(), p.X()),
		max(r.Bottom(), p.Y()),
	)
}

// ExpandSize implements Rect.ExpandSize.
func (r rect[T]) ExpandSize(s Size) Rect {
	halfW, halfH := s.Width()/2, s.Height()/2
	return NewRect(
		r.Left()-halfW,
		r.Top()-halfH,
		r.Right()+halfW,
		r.Bottom()+halfH,
	)
}

// Project implements Rect.Project.
func (r rect[T]) Project(source Rect) Rect {
	if r.IsEmpty() {
		return NewRect[T](0, 0, 0, 0)
	}
	// Transform source rectangle: shift by -r.origin, then scale by 1/r.dimensions
	shifted := NewRect(
		source.Left()-r.Left(),
		source.Top()-r.Top(),
		source.Right()-r.Left(),
		source.Bottom()-r.Top(),
	)

	scaleX := 1.0 / r.Width()
	scaleY := 1.0 / r.Height()

	return shifted.ScaleXY(scaleX, scaleY)
}

// Transform implements Rect.Transform.
func (r rect[T]) Transform(transform Matrix) Quad {
	return Quad{
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
