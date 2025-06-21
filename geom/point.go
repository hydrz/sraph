package geom

import (
	"fmt"
	"image"
	"math"
)

// Quad defines the interface for a 4-point polygon (quadrilateral).
type Quad [4]Point

// Transform applies a transformation matrix to the quad.
func (q Quad) Transform(m Matrix) Quad {
	return Quad{
		q[0].Transform(m),
		q[1].Transform(m),
		q[2].Transform(m),
		q[3].Transform(m),
	}
}

// Point defines the interface for a 2D point or vector in Cartesian coordinates.
// It provides methods for arithmetic operations, geometric queries, and transformations.
type Point interface {
	// X returns the X coordinate.
	X() Scalar
	// Y returns the Y coordinate.
	Y() Scalar

	// Add returns the vector sum of this point and another.
	Add(other Point) Point
	// Sub returns the vector difference of this point and another.
	Sub(other Point) Point
	// Mul returns the element-wise product of this point and another.
	Mul(other Point) Point
	// MulSize returns the element-wise product of this point and a size.
	MulSize(size Size) Point
	// Div returns the element-wise division of this point by another.
	Div(other Point) Point
	// DivSize returns the element-wise division of this point by a size.
	DivSize(size Size) Point
	// Neg returns the negation of this point.
	Neg() Point
	// Min returns a point with the minimum value for each coordinate from this and another point.
	Min(other Point) Point
	// Max returns a point with the maximum value for each coordinate from this and another point.
	Max(other Point) Point
	// Abs returns a point with the absolute value of each coordinate.
	Abs() Point
	// Conj returns the conjugate of this point, flipping the sign of Y.
	Conj() Point
	// Floor returns a point with math.Floor applied to each coordinate.
	Floor() Point
	// Ceil returns a point with math.Ceil applied to each coordinate.
	Ceil() Point
	// Round returns a point with math.Round applied to each coordinate.
	Round() Point
	// Equal reports whether this point and another are nearly equal, using tolerance for floating-point types.
	Equal(other Point) bool
	// IsFinite reports whether both coordinates of this point are finite.
	IsFinite() bool
	// IsZero reports whether both coordinates of this point are nearly zero.
	IsZero() bool
	// Translate returns a new point offset by the given vector.
	Translate(vector Point) Point
	// Scale returns a new point scaled by the given scalar.
	Scale(scale Scalar) Point
	// Rotate returns a new point rotated around the origin by angle (in radians).
	Rotate(angle Radians) Point
	// Normalize returns a unit vector in the direction of this point. If zero, returns (1, 0).
	Normalize() Point
	// Dot returns the dot product of this point and another.
	Dot(other Point) Scalar
	// Cross returns the 2D cross product of this point and another.
	Cross(other Point) Scalar
	// Phase returns the angle (in radians) between this point and the positive X axis, in [-Pi, Pi].
	Phase() Radians
	// Polar returns the length and phase (angle in radians) of this point.
	Polar() (Scalar, Radians)
	// AngleTo returns the angle (in radians) between this point and another.
	AngleTo(other Point) Radians
	// Distance returns the Euclidean distance between this point and another.
	Distance(other Point) Scalar
	// DistanceSquared returns the squared Euclidean distance between this point and another.
	DistanceSquared(other Point) Scalar
	// Length returns the Euclidean length of this point, i.e., the distance from the origin.
	Length() Scalar
	// LengthSquared returns the squared Euclidean length of this point.
	LengthSquared() Scalar
	// Reflect returns the reflection of this point about the given axis vector.
	// The axis does not need to be normalized.
	Reflect(axis Point) Point
	// Lerp returns the linear interpolation between this point and another by parameter t in [0,1].
	Lerp(other Point, t Scalar) Point
	// Complex returns the complex128 representation of this point.
	Complex() complex128
	// Go converts this point to image.Point, truncating coordinates to int.
	Go() image.Point
	// Transform applies the given matrix transformation to this point.
	Transform(m Matrix) Point
	// TransformDirection applies the given matrix transformation to this point, treating it as a direction vector.
	TransformDirection(m Matrix) Point
	// TransformHomogenous transforms this 2D point to 3D homogeneous coordinates.
	TransformHomogenous(m Matrix) Vector3
	// String returns a string representation of this point in the form "(x, y)".
	String() string
}

// NewPoint constructs a Point with the given coordinates.
func NewPoint[T Number](x, y T) Point {
	return point[T]{x: x, y: y}
}

// NewPointPolar constructs a Point from polar coordinates (radius r, angle θ in radians).
func NewPointPolar[T Number](r T, θ Radians) Point {
	s, c := math.Sincos(float64(θ))
	return point[T]{x: r * T(c), y: r * T(s)}
}

// NewPointGo converts an image.Point to a Point.
func NewPointGo[T Number](p image.Point) Point {
	return point[T]{x: T(p.X), y: T(p.Y)}
}

// point is the generic implementation of the Point interface.
type point[T Number] struct {
	x T
	y T
}

// X implements Point.X.
func (p point[T]) X() Scalar { return Scalar(p.x) }

// Y implements Point.Y.
func (p point[T]) Y() Scalar { return Scalar(p.y) }

// Add implements Point.Add.
func (p point[T]) Add(o Point) Point {
	return NewPoint(p.X()+o.X(), p.Y()+o.Y())
}

// Sub implements Point.Sub.
func (p point[T]) Sub(o Point) Point {
	return NewPoint(p.X()-o.X(), p.Y()-o.Y())
}

// Mul implements Point.Mul.
func (p point[T]) Mul(o Point) Point {
	return NewPoint(p.X()*o.X(), p.Y()*o.Y())
}

// MulSize implements Point.MulSize.
func (p point[T]) MulSize(s Size) Point {
	return NewPoint(p.X()*s.Width(), p.Y()*s.Height())
}

// Div implements Point.Div.
func (p point[T]) Div(o Point) Point {
	return NewPoint(p.X()/o.X(), p.Y()/o.Y())
}

// DivSize implements Point.DivSize.
func (p point[T]) DivSize(s Size) Point {
	return NewPoint(p.X()/s.Width(), p.Y()/s.Height())
}

// Neg implements Point.Neg.
func (p point[T]) Neg() Point { return point[T]{x: -p.x, y: -p.y} }

// Min implements Point.Min.
func (p point[T]) Min(o Point) Point {
	return point[T]{x: min(p.x, T(o.X())), y: min(p.y, T(o.Y()))}
}

// Max implements Point.Max.
func (p point[T]) Max(o Point) Point {
	return point[T]{x: max(p.x, T(o.X())), y: max(p.y, T(o.Y()))}
}

// Abs implements Point.Abs.
func (p point[T]) Abs() Point {
	var zero T
	x := p.x
	y := p.y
	if x < zero {
		x = -x
	}
	if y < zero {
		y = -y
	}
	return point[T]{x: x, y: y}
}

// Conj implements Point.Conj.
func (p point[T]) Conj() Point { return point[T]{x: p.x, y: -p.y} }

// Floor implements Point.Floor.
func (p point[T]) Floor() Point {
	return point[T]{x: T(math.Floor(float64(p.x))), y: T(math.Floor(float64(p.y)))}
}

// Ceil implements Point.Ceil.
func (p point[T]) Ceil() Point {
	return point[T]{x: T(math.Ceil(float64(p.x))), y: T(math.Ceil(float64(p.y)))}
}

// Round implements Point.Round.
func (p point[T]) Round() Point {
	return point[T]{x: T(math.Round(float64(p.x))), y: T(math.Round(float64(p.y)))}
}

// Equal implements Point.Equal.
func (p point[T]) Equal(other Point) bool {
	return NearlyEqual(p.x, T(other.X())) && NearlyEqual(p.y, T(other.Y()))
}

// IsFinite implements Point.IsFinite.
func (p point[T]) IsFinite() bool {
	return IsFinite(p.x) && IsFinite(p.y)
}

// IsZero implements Point.IsZero.
func (p point[T]) IsZero() bool {
	return NearlyEqual(p.x, 0) && NearlyEqual(p.y, 0)
}

// Translate implements Point.Translate.
func (p point[T]) Translate(v Point) Point {
	return point[T]{x: p.x + T(v.X()), y: p.y + T(v.Y())}
}

// Scale implements Point.Scale.
func (p point[T]) Scale(scale Scalar) Point {
	return NewPoint(p.X()*scale, p.Y()*scale)
}

// Rotate implements Point.Rotate.
func (p point[T]) Rotate(angle Radians) Point {
	cos := T(math.Cos(float64(angle)))
	sin := T(math.Sin(float64(angle)))
	return point[T]{x: p.x*cos - p.y*sin, y: p.x*sin + p.y*cos}
}

// Normalize implements Point.Normalize.
func (p point[T]) Normalize() Point {
	len := p.Length()
	if len == 0 {
		return NewPoint(1, 0)
	}
	return NewPoint(p.X()/len, p.Y()/len)
}

// Dot implements Point.Dot.
func (p point[T]) Dot(o Point) Scalar {
	return p.X()*o.X() + p.Y()*o.Y()
}

// Cross implements Point.Cross.
func (p point[T]) Cross(o Point) Scalar {
	return p.X()*o.Y() - p.Y()*o.X()
}

// Phase implements Point.Phase.
func (p point[T]) Phase() Radians {
	return Radians(math.Atan2(float64(p.y), float64(p.x)))
}

// Polar implements Point.Polar.
func (p point[T]) Polar() (Scalar, Radians) {
	return p.Length(), p.Phase()
}

// AngleTo implements Point.AngleTo.
func (p point[T]) AngleTo(o Point) Radians {
	return Radians(
		math.Atan2(
			ToFloat64(p.x*T(o.Y()))-ToFloat64(p.y*T(o.X())),
			ToFloat64(p.x*T(o.X()))+ToFloat64(p.y*T(o.Y())),
		),
	)
}

// Distance implements Point.Distance.
func (p point[T]) Distance(o Point) Scalar {
	return Scalar(
		math.Sqrt(ToFloat64(p.DistanceSquared(o))),
	)
}

// DistanceSquared implements Point.DistanceSquared.
func (p point[T]) DistanceSquared(o Point) Scalar {
	dx := o.X() - p.X()
	dy := o.Y() - p.Y()
	return dx*dx + dy*dy
}

// Length implements Point.Length.
func (p point[T]) Length() Scalar {
	return Scalar(
		math.Sqrt(ToFloat64(p.LengthSquared())),
	)
}

// LengthSquared implements Point.LengthSquared.
func (p point[T]) LengthSquared() Scalar {
	return p.X()*p.X() + p.Y()*p.Y()
}

// Reflect implements Point.Reflect.
func (p point[T]) Reflect(axis Point) Point {
	return p.Sub(
		axis.Scale(p.Dot(axis)).Scale(2),
	)
}

// Lerp implements Point.Lerp.
func (p point[T]) Lerp(o Point, t Scalar) Point {
	return point[T]{
		x: p.x + (T(o.X())-p.x)*T(t),
		y: p.y + (T(o.Y())-p.y)*T(t),
	}
}

// Complex implements Point.Complex.
func (p point[T]) Complex() complex128 {
	return complex(float64(p.x), float64(p.y))
}

// Go implements Point.Go.
func (p point[T]) Go() image.Point {
	return image.Point{X: int(p.x), Y: int(p.y)}
}

// Transform implements Point.Transform.
func (p point[T]) Transform(m Matrix) Point {
	w := p.x*T(m[3]) + p.y*T(m[7]) + T(m[15])
	rx := p.x*T(m[0]) + p.y*T(m[4]) + T(m[12])
	ry := p.x*T(m[1]) + p.y*T(m[5]) + T(m[13])
	if w != 0 {
		w = 1 / w
	}
	return point[T]{x: rx * w, y: ry * w}
}

// TransformDirection implements Point.TransformDirection.
func (p point[T]) TransformDirection(m Matrix) Point {
	return point[T]{
		x: p.x*T(m[0]) + p.y*T(m[4]),
		y: p.x*T(m[1]) + p.y*T(m[5]),
	}
}

// TransformHomogenous implements Point.TransformHomogenous.
func (p point[T]) TransformHomogenous(m Matrix) Vector3 {
	return NewVector3(
		p.x*T(m[0])+p.y*T(m[4])+T(m[12]),
		p.x*T(m[1])+p.y*T(m[5])+T(m[13]),
		p.x*T(m[3])+p.y*T(m[7])+T(m[15]),
	)
}

// String implements Point.String.
func (p point[T]) String() string {
	return fmt.Sprintf("(%v, %v)", p.x, p.y)
}
