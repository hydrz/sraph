package geom

import (
	"image"
	"math"
)

type Quad[T TScalar] = [4]Point[T]

// Point represents a 2D point with X and Y coordinates.
type Point[T TScalar] struct {
	X T
	Y T
}

func NewPoint[T TScalar](x, y T) Point[T] {
	return Point[T]{X: x, Y: y}
}

func NewPointGo[T TScalar](p image.Point) Point[T] {
	return Point[T]{X: T(p.X), Y: T(p.Y)}
}

// Add returns the sum of this point and another.
func (p Point[T]) Add(o Point[T]) Point[T] {
	return Point[T]{X: p.X + o.X, Y: p.Y + o.Y}
}

// Sub returns the difference of this point and another.
func (p Point[T]) Sub(o Point[T]) Point[T] {
	return Point[T]{X: p.X - o.X, Y: p.Y - o.Y}
}

// Mul returns the element-wise product of this point and another.
func (p Point[T]) Mul(o Point[T]) Point[T] {
	return Point[T]{X: p.X * o.X, Y: p.Y * o.Y}
}

// MulSize returns the element-wise product of this point and a Size.
func (p Point[T]) MulSize(o Size[T]) Point[T] {
	return Point[T]{X: p.X * o.Width, Y: p.Y * o.Height}
}

// Div returns the element-wise division of this point by another.
func (p Point[T]) Div(o Point[T]) Point[T] {
	return Point[T]{X: p.X / o.X, Y: p.Y / o.Y}
}

// DivSize returns the element-wise division of this point by a Size.
func (p Point[T]) DivSize(o Size[T]) Point[T] {
	return Point[T]{X: p.X / o.Width, Y: p.Y / o.Height}
}

// Neg returns the negated point.
func (p Point[T]) Neg() Point[T] {
	return Point[T]{X: -p.X, Y: -p.Y}
}

// Min returns the point with minimum values from this point and another.
func (p Point[T]) Min(o Point[T]) Point[T] {
	return Point[T]{
		X: min(p.X, o.X),
		Y: min(p.Y, o.Y),
	}
}

// Max returns the point with maximum values from this point and another.
func (p Point[T]) Max(o Point[T]) Point[T] {
	return Point[T]{
		X: max(p.X, o.X),
		Y: max(p.Y, o.Y),
	}
}

// Abs returns the point with absolute values.
func (p Point[T]) Abs() Point[T] {
	var zero T
	x := p.X
	y := p.Y
	if x < zero {
		x = -x
	}
	if y < zero {
		y = -y
	}
	return Point[T]{X: x, Y: y}
}

// Floor returns the point with floor applied to each coordinate.
func (p Point[T]) Floor() Point[T] {
	return Point[T]{
		X: T(math.Floor(ToFloat64(p.X))),
		Y: T(math.Floor(ToFloat64(p.Y))),
	}
}

// Ceil returns the point with ceil applied to each coordinate.
func (p Point[T]) Ceil() Point[T] {
	return Point[T]{
		X: T(math.Ceil(ToFloat64(p.X))),
		Y: T(math.Ceil(ToFloat64(p.Y))),
	}
}

// Round returns the point with round applied to each coordinate.
func (p Point[T]) Round() Point[T] {
	return Point[T]{
		X: T(math.Round(ToFloat64(p.X))),
		Y: T(math.Round(ToFloat64(p.Y))),
	}
}

// Eq reports whether p and o are equal.
func (p Point[T]) Eq(o Point[T]) bool {
	return NearlyEqual(p.X, o.X) && NearlyEqual(p.Y, o.Y)
}

// IsFinite returns true if both coordinates are finite.
func (p Point[T]) IsFinite() bool {
	return IsFinite(p.X) && IsFinite(p.Y)
}

// IsZero returns true if both coordinates are zero.
func (p Point[T]) IsZero() bool {
	var zero T
	return p.X == zero && p.Y == zero
}

// Translate returns a new point translated by the given vector.
func (p Point[T]) Translate(vector Point[T]) Point[T] {
	return Point[T]{X: p.X + vector.X, Y: p.Y + vector.Y}
}

// Scale returns the point scaled by a scalar.
func (p Point[T]) Scale(scale T) Point[T] {
	return Point[T]{X: p.X * scale, Y: p.Y * scale}
}

// Rotate rotates this point around the origin by the given angle in radians.
func (p Point[T]) Rotate(angle Radians) Point[T] {
	cos := T(math.Cos(ToFloat64(angle)))
	sin := T(math.Sin(ToFloat64(angle)))
	return Point[T]{
		X: p.X*cos - p.Y*sin,
		Y: p.X*sin + p.Y*cos,
	}
}

// Normalize returns a normalized version of this point.
// If the point is zero, it returns a default unit vector (1, 0).
func (p Point[T]) Normalize() Point[T] {
	len := p.Length()
	var zero T
	if len == zero {
		return Point[T]{X: 1, Y: 0}
	}
	return Point[T]{X: p.X / len, Y: p.Y / len}
}

// Dot calculates the dot product of this point and another.
func (p Point[T]) Dot(o Point[T]) T {
	return p.X*o.X + p.Y*o.Y
}

// Cross calculates the cross product of this point and another.
func (p Point[T]) Cross(o Point[T]) T {
	return p.X*o.Y - p.Y*o.X
}

// AngleTo returns the angle in radians between this point and another.
func (p Point[T]) AngleTo(o Point[T]) Radians {
	return Radians(math.Atan2(ToFloat64(p.Cross(o)), ToFloat64(p.Dot(o))))
}

// Distance returns the Euclidean distance between this point and another.
func (p Point[T]) Distance(o Point[T]) T {
	dx := p.X - o.X
	dy := p.Y - o.Y
	return T(math.Sqrt(ToFloat64(dx*dx + dy*dy)))
}

// DistanceSquared returns the squared Euclidean distance between this point and another.
func (p Point[T]) DistanceSquared(o Point[T]) T {
	dx := p.X - o.X
	dy := p.Y - o.Y
	return dx*dx + dy*dy
}

// Length returns the distance from this point to the origin (0,0).
func (p Point[T]) Length() T {
	return T(math.Sqrt(ToFloat64(p.X*p.X + p.Y*p.Y)))
}

// LengthSquared returns the squared distance from this point to the origin (0,0).
func (p Point[T]) LengthSquared() T {
	return p.X*p.X + p.Y*p.Y
}

// Reflect calculates the reflection vector of this point (vector) about the specified axis (vector).
func (p Point[T]) Reflect(axis Point[T]) Point[T] {
	axis = axis.Normalize()
	dot := p.Dot(axis)
	return p.Sub(axis.Scale(dot).Scale(2))
}

// Lerp performs linear interpolation between this point and another.
func (p Point[T]) Lerp(o Point[T], t T) Point[T] {
	return Point[T]{
		X: p.X + (o.X-p.X)*t,
		Y: p.Y + (o.Y-p.Y)*t,
	}
}

// String returns a string representation of the point, like "(x, y)".
func (p Point[T]) String() string {
	return "(" + ToString(p.X) + ", " + ToString(p.Y) + ")"
}

// ToGo converts this point to an image.Point.
func (p Point[T]) ToGo() image.Point {
	return image.Point{X: int(ToFloat64(p.X)), Y: int(ToFloat64(p.Y))}
}
