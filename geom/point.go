package geom

import (
	"image"
	"math"
)

type Quad[T TScalar] [4]Point[T]

// Transform applies a transformation matrix to the quad.
func (q Quad[T]) Transform(m Matrix[T]) Quad[T] {
	return Quad[T]{
		q[0].Transform(m),
		q[1].Transform(m),
		q[2].Transform(m),
		q[3].Transform(m),
	}
}

// Point represents a 2D point or vector in Cartesian coordinates.
type Point[T TScalar] struct {
	X T
	Y T
}

// NewPoint constructs a Point with the given coordinates.
func NewPoint[T TScalar](x, y T) Point[T] {
	return Point[T]{X: x, Y: y}
}

// NewPointPolar constructs a Point from polar coordinates (radius r, angle θ in radians).
func NewPointPolar[T TScalar](r T, θ Radians) Point[T] {
	s, c := math.Sincos(ToFloat64(θ))
	return Point[T]{X: r * T(c), Y: r * T(s)}
}

// NewPointGo converts an image.Point to a Point.
func NewPointGo[T TScalar](p image.Point) Point[T] {
	return Point[T]{X: T(p.X), Y: T(p.Y)}
}

// Add returns the vector sum of p and o.
func (p Point[T]) Add(o Point[T]) Point[T] { return Point[T]{p.X + o.X, p.Y + o.Y} }

// Sub returns the vector difference of p and o.
func (p Point[T]) Sub(o Point[T]) Point[T] { return Point[T]{p.X - o.X, p.Y - o.Y} }

// Mul returns the element-wise product of p and o.
func (p Point[T]) Mul(o Point[T]) Point[T] { return Point[T]{p.X * o.X, p.Y * o.Y} }

// MulSize returns the element-wise product of p and a Size.
func (p Point[T]) MulSize(o Size[T]) Point[T] { return Point[T]{p.X * o.Width, p.Y * o.Height} }

// Div returns the element-wise division of p by o.
func (p Point[T]) Div(o Point[T]) Point[T] { return Point[T]{p.X / o.X, p.Y / o.Y} }

// DivSize returns the element-wise division of p by a Size.
func (p Point[T]) DivSize(o Size[T]) Point[T] { return Point[T]{p.X / o.Width, p.Y / o.Height} }

// Neg returns the negation of p.
func (p Point[T]) Neg() Point[T] { return Point[T]{-p.X, -p.Y} }

// Min returns a point with the minimum value for each coordinate from p and o.
func (p Point[T]) Min(o Point[T]) Point[T] {
	return Point[T]{min(p.X, o.X), min(p.Y, o.Y)}
}

// Max returns a point with the maximum value for each coordinate from p and o.
func (p Point[T]) Max(o Point[T]) Point[T] {
	return Point[T]{max(p.X, o.X), max(p.Y, o.Y)}
}

// Abs returns a point with the absolute value of each coordinate.
func (p Point[T]) Abs() Point[T] {
	x := Cond(p.X < 0, -p.X, p.X)
	y := Cond(p.Y < 0, -p.Y, p.Y)
	return Point[T]{X: x, Y: y}
}

// Conj returns the conjugate of p, flipping the sign of Y.
func (p Point[T]) Conj() Point[T] {
	return Point[T]{X: p.X, Y: -p.Y}
}

// Floor returns a point with math.Floor applied to each coordinate.
func (p Point[T]) Floor() Point[T] {
	return Point[T]{
		X: T(math.Floor(ToFloat64(p.X))),
		Y: T(math.Floor(ToFloat64(p.Y))),
	}
}

// Ceil returns a point with math.Ceil applied to each coordinate.
func (p Point[T]) Ceil() Point[T] {
	return Point[T]{
		X: T(math.Ceil(ToFloat64(p.X))),
		Y: T(math.Ceil(ToFloat64(p.Y))),
	}
}

// Round returns a point with math.Round applied to each coordinate.
func (p Point[T]) Round() Point[T] {
	return Point[T]{
		X: T(math.Round(ToFloat64(p.X))),
		Y: T(math.Round(ToFloat64(p.Y))),
	}
}

// Equal reports whether p and o are nearly equal, using tolerance for floating-point types.
func (p Point[T]) Equal(o Point[T]) bool {
	return NearlyEqual(p.X, o.X) && NearlyEqual(p.Y, o.Y)
}

// IsFinite reports whether both coordinates of p are finite.
func (p Point[T]) IsFinite() bool {
	return IsFinite(p.X) && IsFinite(p.Y)
}

// IsZero reports whether both coordinates of p are nearly zero.
func (p Point[T]) IsZero() bool {
	return NearlyEqual(p.X, 0) && NearlyEqual(p.Y, 0)
}

// Translate returns a new point offset by the given vector.
func (p Point[T]) Translate(vector Point[T]) Point[T] {
	return Point[T]{X: p.X + vector.X, Y: p.Y + vector.Y}
}

// Scale returns a new point scaled by the given scalar.
func (p Point[T]) Scale(scale T) Point[T] {
	return Point[T]{X: p.X * scale, Y: p.Y * scale}
}

// Rotate returns a new point rotated around the origin by angle (in radians).
func (p Point[T]) Rotate(angle Radians) Point[T] {
	cos := T(math.Cos(ToFloat64(angle)))
	sin := T(math.Sin(ToFloat64(angle)))
	return Point[T]{
		X: p.X*cos - p.Y*sin,
		Y: p.X*sin + p.Y*cos,
	}
}

// Normalize returns a unit vector in the direction of p.
// If p is zero, returns (1, 0).
func (p Point[T]) Normalize() Point[T] {
	len := p.Length()
	var zero T
	if len == zero {
		return Point[T]{X: 1, Y: 0}
	}
	return Point[T]{X: p.X / len, Y: p.Y / len}
}

// Dot returns the dot product of p and o.
func (p Point[T]) Dot(o Point[T]) T {
	return p.X*o.X + p.Y*o.Y
}

// Cross returns the 2D cross product of p and o.
func (p Point[T]) Cross(o Point[T]) T {
	return p.X*o.Y - p.Y*o.X
}

// Phase returns the angle (in radians) between p and the positive X axis, in [-Pi, Pi].
func (p Point[T]) Phase() (θ Radians) {
	return Radians(math.Atan2(ToFloat64(p.Y), ToFloat64(p.X)))
}

// Polar returns the length and phase (angle in radians) of p.
func (p Point[T]) Polar() (r T, θ Radians) {
	return p.Length(), p.Phase()
}

// AngleTo returns the angle (in radians) between p and o.
func (p Point[T]) AngleTo(o Point[T]) Radians {
	return Radians(math.Atan2(ToFloat64(p.Cross(o)), ToFloat64(p.Dot(o))))
}

// Distance returns the Euclidean distance between p and o.
func (p Point[T]) Distance(o Point[T]) T {
	return T(math.Hypot(ToFloat64(p.X-o.X), ToFloat64(p.Y-o.Y)))
}

// DistanceSquared returns the squared Euclidean distance between p and o.
func (p Point[T]) DistanceSquared(o Point[T]) T {
	dx := p.X - o.X
	dy := p.Y - o.Y
	return dx*dx + dy*dy
}

// Length returns the Euclidean length of p, i.e., the distance from the origin.
func (p Point[T]) Length() T {
	return p.Distance(Point[T]{})
}

// LengthSquared returns the squared Euclidean length of p.
func (p Point[T]) LengthSquared() T {
	return p.X*p.X + p.Y*p.Y
}

// Reflect returns the reflection of p about the given axis vector.
func (p Point[T]) Reflect(axis Point[T]) Point[T] {
	axis = axis.Normalize()
	dot := p.Dot(axis)
	return p.Sub(axis.Scale(dot).Scale(2))
}

// Lerp returns the linear interpolation between p and o by parameter t in [0,1].
func (p Point[T]) Lerp(o Point[T], t T) Point[T] {
	return Point[T]{
		X: p.X + (o.X-p.X)*t,
		Y: p.Y + (o.Y-p.Y)*t,
	}
}

// Complex returns the complex128 representation of p.
func (p Point[T]) Complex() complex128 {
	return complex(ToFloat64(p.X), ToFloat64(p.Y))
}

// Go converts p to image.Point, truncating coordinates to int.
func (p Point[T]) Go() image.Point {
	return image.Point{X: int(ToFloat64(p.X)), Y: int(ToFloat64(p.Y))}
}

// Transform applies the given matrix transformation to the point p.
func (p Point[T]) Transform(m Matrix[T]) Point[T] {
	w := p.X*m[3] + p.Y*m[7] + m[15]
	r := Point[T]{
		X: p.X*m[0] + p.Y*m[4] + m[12],
		Y: p.X*m[1] + p.Y*m[5] + m[13],
	}
	if w != 0 {
		w = 1 / w
	}
	return Point[T]{r.X * w, r.Y * w}
}

// TransformDirection applies the given matrix transformation to the point p,
// treating it as a direction vector (ignoring translation).
func (p Point[T]) TransformDirection(m Matrix[T]) Point[T] {
	return Point[T]{
		X: p.X*m[0] + p.Y*m[4],
		Y: p.X*m[1] + p.Y*m[5],
	}
}

// TransformHomogenous transforms a 2D point to 3D homogeneous coordinates.
// It applies the transformation matrix to the point and returns a Vector3
func (p Point[T]) TransformHomogenous(m Matrix[T]) Vector3[T] {
	return Vector3[T]{
		X: p.X*m[0] + p.Y*m[4] + m[12],
		Y: p.X*m[1] + p.Y*m[5] + m[13],
		Z: p.X*m[3] + p.Y*m[7] + m[15],
	}
}

// String returns a string representation of p in the form "(x, y)".
func (p Point[T]) String() string {
	return "(" + ToString(p.X) + ", " + ToString(p.Y) + ")"
}
