package geom

import "math"

// Vector2 defines the interface for a 2D vector.
// Provides methods for vector arithmetic, normalization, and geometric queries.
type Vector2 = Point

// Vector3 defines the interface for a 3D vector.
// Provides methods for vector arithmetic, normalization, and geometric queries.
type Vector3 interface {
	X() Scalar
	Y() Scalar
	Z() Scalar
	Add(other Vector3) Vector3
	Sub(other Vector3) Vector3
	Mul(other Vector3) Vector3
	Div(other Vector3) Vector3
	Equal(other Vector3) bool
	Scale(scalar Scalar) Vector3
	Normalize() Vector3
	Length() Scalar
	Abs() Vector3
	Floor() Vector3
	Ceil() Vector3
	Round() Vector3
	Dot(other Vector3) Scalar
	Cross(other Vector3) Vector3
	Lerp(other Vector3, t Scalar) Vector3
	IsZero() bool
	Transform(m Matrix) Vector3
	TransformDirection(m Matrix) Vector3
	String() string
	Combine(other Vector3, factor Scalar) Vector3
}

// Vector4 defines the interface for a 4D vector.
type Vector4 interface {
	X() Scalar
	Y() Scalar
	Z() Scalar
	W() Scalar
	Add(other Vector4) Vector4
	Sub(other Vector4) Vector4
	Mul(other Vector4) Vector4
	Div(other Vector4) Vector4
	Equal(other Vector4) bool
	Scale(scalar Scalar) Vector4
	Normalize() Vector4
	Length() Scalar
	Abs() Vector4
	Floor() Vector4
	Ceil() Vector4
	Round() Vector4
	Dot(other Vector4) Scalar
	Lerp(other Vector4, t Scalar) Vector4
	IsZero() bool
	IsFinite() bool
	Transform(m Matrix) Vector4
	TransformDirection(m Matrix) Vector4
	String() string
}

func NewVector2[T Number](x, y T) Vector2 {
	return NewPoint(x, y)
}

// NewVector3 creates a new 3D vector.
func NewVector3[T Number](x, y, z T) Vector3 {
	return vector3[T]{x: x, y: y, z: z}
}

// NewVector4 creates a new 4D vector.
func NewVector4[T Number](x, y, z, w T) Vector4 {
	return vector4[T]{x: x, y: y, z: z, w: w}
}

type vector3[T Number] struct {
	x, y, z T
}

// X returns the X component of the vector.
func (v vector3[T]) X() Scalar {
	return Scalar(v.x)
}

// Y returns the Y component of the vector.
func (v vector3[T]) Y() Scalar {
	return Scalar(v.y)
}

// Z returns the Z component of the vector.
func (v vector3[T]) Z() Scalar {
	return Scalar(v.z)
}

// Add adds another vector to this vector.
func (v vector3[T]) Add(o Vector3) Vector3 {
	return vector3[T]{x: v.x + T(o.X()), y: v.y + T(o.Y()), z: v.z + T(o.Z())}
}

// Sub subtracts another vector from this vector.
func (v vector3[T]) Sub(other Vector3) Vector3 {
	return vector3[T]{x: v.x - T(other.X()), y: v.y - T(other.Y()), z: v.z - T(other.Z())}
}

// Mul multiplies this vector by another vector component-wise.
func (v vector3[T]) Mul(other Vector3) Vector3 {
	return vector3[T]{x: v.x * T(other.X()), y: v.y * T(other.Y()), z: v.z * T(other.Z())}
}

// Div divides this vector by another vector component-wise.
func (v vector3[T]) Div(other Vector3) Vector3 {
	return vector3[T]{x: v.x / T(other.X()), y: v.y / T(other.Y()), z: v.z / T(other.Z())}
}

// Equal checks if this vector is Eq to another vector.
func (v vector3[T]) Equal(o Vector3) bool {
	return NearlyEqual(v.x, T(o.X())) && NearlyEqual(v.y, T(o.Y())) && NearlyEqual(v.z, T(o.Z()))
}

// Scale scales this vector by a scalar.
func (v vector3[T]) Scale(o Scalar) Vector3 {
	return vector3[T]{x: v.x * T(o), y: v.y * T(o), z: v.z * T(o)}
}

// Normalize returns the normalized (unit) vector.
// If the vector's length is zero, returns a zero vector.
// Normalization scales the vector to have a length of 1, while maintaining its direction.
func (v vector3[T]) Normalize() Vector3 {
	len := T(v.Length())
	if len == 0 {
		return vector3[T]{}
	}
	return vector3[T]{x: v.x / len, y: v.y / len, z: v.z / len}
}

// Length returns the magnitude (Euclidean norm) of the vector.
// It is defined as: length = sqrt(x^2 + y^2 + z^2).
func (v vector3[T]) Length() Scalar {
	return Scalar(math.Sqrt(float64(v.x*v.x + v.y*v.y + v.z*v.z)))
}

// Abs returns the component-wise absolute value of the vector.
func (v vector3[T]) Abs() Vector3 {
	var zero T
	x, y, z := v.x, v.y, v.z
	if x < zero {
		x = -x
	}
	if y < zero {
		y = -y
	}
	if z < zero {
		z = -z
	}
	return vector3[T]{x: x, y: y, z: z}
}

// Floor returns the component-wise floor.
func (v vector3[T]) Floor() Vector3 {
	return vector3[T]{
		x: T(math.Floor(ToFloat64(v.x))),
		y: T(math.Floor(ToFloat64(v.y))),
		z: T(math.Floor(ToFloat64(v.z))),
	}
}

// Ceil returns the component-wise ceil.
func (v vector3[T]) Ceil() Vector3 {
	return vector3[T]{
		x: T(math.Ceil(ToFloat64(v.x))),
		y: T(math.Ceil(ToFloat64(v.y))),
		z: T(math.Ceil(ToFloat64(v.z))),
	}
}

// Round returns the component-wise round.
func (v vector3[T]) Round() Vector3 {
	return vector3[T]{
		x: T(math.Round(ToFloat64(v.x))),
		y: T(math.Round(ToFloat64(v.y))),
		z: T(math.Round(ToFloat64(v.z))),
	}
}

// Dot returns the dot product of this vector and another.
func (v vector3[T]) Dot(other Vector3) Scalar {
	o := other.(vector3[T])
	return Scalar(v.x*o.x + v.y*o.y + v.z*o.z)
}

// Cross returns the cross product of this vector and another.
func (v vector3[T]) Cross(other Vector3) Vector3 {
	o := other.(vector3[T])
	return vector3[T]{
		x: v.y*o.z - v.z*o.y,
		y: v.z*o.x - v.x*o.z,
		z: v.x*o.y - v.y*o.x,
	}
}

// Lerp linearly interpolates between this vector and another by t.
func (v vector3[T]) Lerp(other Vector3, t Scalar) Vector3 {
	o := other.(vector3[T])
	return vector3[T]{
		x: v.x + (o.x-v.x)*T(t),
		y: v.y + (o.y-v.y)*T(t),
		z: v.z + (o.z-v.z)*T(t),
	}
}

// IsZero checks if all components are nearly zero.
func (v vector3[T]) IsZero() bool {
	return NearlyEqual(v.x, 0) && NearlyEqual(v.y, 0) && NearlyEqual(v.z, 0)
}

// Transform applies a transformation matrix to this vector.
func (v vector3[T]) Transform(m Matrix) Vector3 {
	w := v.x*T(m[3]) + v.y*T(m[7]) + v.z*T(m[11]) + T(m[15])
	r := vector3[T]{
		x: v.x*T(m[0]) + v.y*T(m[4]) + v.z*T(m[8]) + T(m[12]),
		y: v.x*T(m[1]) + v.y*T(m[5]) + v.z*T(m[9]) + T(m[13]),
		z: v.x*T(m[2]) + v.y*T(m[6]) + v.z*T(m[10]) + T(m[14]),
	}
	if w != 0 {
		w = 1 / w
	}
	return r.Scale(Scalar(w))
}

// TransformDirection applies a transformation matrix to this vector,
// treating it as a direction vector (ignoring translation).
func (v vector3[T]) TransformDirection(m Matrix) Vector3 {
	return vector3[T]{
		x: v.x*T(m[0]) + v.y*T(m[4]) + v.z*T(m[8]),
		y: v.x*T(m[1]) + v.y*T(m[5]) + v.z*T(m[9]),
		z: v.x*T(m[2]) + v.y*T(m[6]) + v.z*T(m[10]),
	}
}

// String returns a string representation of the vector, like "(1, 2, 3)".
func (v vector3[T]) String() string {
	return "(" + ToString(v.x) + ", " + ToString(v.y) + ", " + ToString(v.z) + ")"
}

// Combine combines this vector with another vector by adding the first vector
// and multiplying the second vector by a scalar factor.
func (v vector3[T]) Combine(other Vector3, factor Scalar) Vector3 {
	return vector3[T]{
		x: v.x + T(factor)*T(other.X()),
		y: v.y + T(factor)*T(other.Y()),
		z: v.z + T(factor)*T(other.Z()),
	}
}

type vector4[T Number] struct {
	x, y, z, w T
}

// X returns the X component of the vector.
func (v vector4[T]) X() Scalar {
	return Scalar(v.x)
}

// Y returns the Y component of the vector.
func (v vector4[T]) Y() Scalar {
	return Scalar(v.y)
}

// Z returns the Z component of the vector.
func (v vector4[T]) Z() Scalar {
	return Scalar(v.z)
}

// W returns the W component of the vector.
func (v vector4[T]) W() Scalar {
	return Scalar(v.w)
}

// Add adds another vector to this vector.
func (v vector4[T]) Add(other Vector4) Vector4 {
	return vector4[T]{
		x: v.x + T(other.X()),
		y: v.y + T(other.Y()),
		z: v.z + T(other.Z()),
		w: v.w + T(other.W()),
	}
}

// Sub subtracts another vector from this vector.
func (v vector4[T]) Sub(other Vector4) Vector4 {
	return vector4[T]{
		x: v.x - T(other.X()),
		y: v.y - T(other.Y()),
		z: v.z - T(other.Z()),
		w: v.w - T(other.W()),
	}
}

// Mul multiplies this vector by another vector component-wise.
func (v vector4[T]) Mul(other Vector4) Vector4 {
	return vector4[T]{
		x: v.x * T(other.X()),
		y: v.y * T(other.Y()),
		z: v.z * T(other.Z()),
		w: v.w * T(other.W()),
	}
}

// Div divides this vector by another vector component-wise.
func (v vector4[T]) Div(other Vector4) Vector4 {
	return vector4[T]{
		x: v.x / T(other.X()),
		y: v.y / T(other.Y()),
		z: v.z / T(other.Z()),
		w: v.w / T(other.W()),
	}
}

// Equal checks if this vector is Eq to another vector.
func (v vector4[T]) Equal(other Vector4) bool {
	return NearlyEqual(v.x, T(other.X())) &&
		NearlyEqual(v.y, T(other.Y())) &&
		NearlyEqual(v.z, T(other.Z())) &&
		NearlyEqual(v.w, T(other.W()))
}

// Scale scales this vector by a scalar.
func (v vector4[T]) Scale(scalar Scalar) Vector4 {
	return vector4[T]{
		x: v.x * T(scalar),
		y: v.y * T(scalar),
		z: v.z * T(scalar),
		w: v.w * T(scalar),
	}
}

// Normalize returns the normalized (unit) vector.
// If the vector's length is zero, returns a zero vector.
// Normalization scales the vector to have a length of 1, while maintaining its direction.
func (v vector4[T]) Normalize() Vector4 {
	len := T(v.Length())
	if len == 0 {
		return vector4[T]{}
	}
	return vector4[T]{x: v.x / len, y: v.y / len, z: v.z / len, w: v.w / len}
}

// Length returns the magnitude (Euclidean norm) of the vector.
// It is defined as: length = sqrt(x^2 + y^2 + z^2 + w^2).
func (v vector4[T]) Length() Scalar {
	return Scalar(math.Sqrt(float64(v.x*v.x + v.y*v.y + v.z*v.z + v.w*v.w)))
}

// Abs returns the component-wise absolute value of the vector.
func (v vector4[T]) Abs() Vector4 {
	var zero T
	x, y, z, w := v.x, v.y, v.z, v.w
	if x < zero {
		x = -x
	}
	if y < zero {
		y = -y
	}
	if z < zero {
		z = -z
	}
	if w < zero {
		w = -w
	}
	return vector4[T]{x: x, y: y, z: z, w: w}
}

// Floor returns the component-wise floor.
func (v vector4[T]) Floor() Vector4 {
	return vector4[T]{
		x: T(math.Floor(ToFloat64(v.x))),
		y: T(math.Floor(ToFloat64(v.y))),
		z: T(math.Floor(ToFloat64(v.z))),
		w: T(math.Floor(ToFloat64(v.w))),
	}
}

// Ceil returns the component-wise ceil.
func (v vector4[T]) Ceil() Vector4 {
	return vector4[T]{
		x: T(math.Ceil(ToFloat64(v.x))),
		y: T(math.Ceil(ToFloat64(v.y))),
		z: T(math.Ceil(ToFloat64(v.z))),
		w: T(math.Ceil(ToFloat64(v.w))),
	}
}

// Round returns the component-wise round.
func (v vector4[T]) Round() Vector4 {
	return vector4[T]{
		x: T(math.Round(ToFloat64(v.x))),
		y: T(math.Round(ToFloat64(v.y))),
		z: T(math.Round(ToFloat64(v.z))),
		w: T(math.Round(ToFloat64(v.w))),
	}
}

// Dot returns the dot product of this vector and another.
func (v vector4[T]) Dot(other Vector4) Scalar {
	o := other.(vector4[T])
	return Scalar(v.x*o.x + v.y*o.y + v.z*o.z + v.w*o.w)
}

// Cross returns the cross product of this vector and another.
// Cross product is not well-defined for 4D vectors, so this returns a zero vector.
func (v vector4[T]) Cross(other Vector4) Vector4 {
	var zero T
	return vector4[T]{x: zero, y: zero, z: zero, w: zero}
}

// Lerp linearly interpolates between this vector and another by t.
func (v vector4[T]) Lerp(o Vector4, t Scalar) Vector4 {
	return vector4[T]{
		x: v.x + (T(o.X())-v.x)*T(t),
		y: v.y + (T(o.Y())-v.y)*T(t),
		z: v.z + (T(o.Z())-v.z)*T(t),
		w: v.w + (T(o.W())-v.w)*T(t),
	}
}

// IsZero checks if all components are nearly zero.
func (v vector4[T]) IsZero() bool {
	return NearlyEqual(v.x, 0) && NearlyEqual(v.y, 0) &&
		NearlyEqual(v.z, 0) && NearlyEqual(v.w, 0)
}

// IsFinite returns true if all components are finite.
func (v vector4[T]) IsFinite() bool {
	return IsFinite(v.x) && IsFinite(v.y) && IsFinite(v.z) && IsFinite(v.w)
}

// Transform applies a transformation matrix to this vector.
func (v vector4[T]) Transform(m Matrix) Vector4 {
	return vector4[T]{
		x: v.x*T(m[0]) + v.y*T(m[4]) + v.z*T(m[8]) + v.w*T(m[12]),
		y: v.x*T(m[1]) + v.y*T(m[5]) + v.z*T(m[9]) + v.w*T(m[13]),
		z: v.x*T(m[2]) + v.y*T(m[6]) + v.z*T(m[10]) + v.w*T(m[14]),
		w: v.x*T(m[3]) + v.y*T(m[7]) + v.z*T(m[11]) + v.w*T(m[15]),
	}
}

// TransformDirection applies a transformation matrix to this vector,
// treating it as a direction vector (ignoring translation).
func (v vector4[T]) TransformDirection(m Matrix) Vector4 {
	return vector4[T]{
		x: v.x*T(m[0]) + v.y*T(m[4]) + v.z*T(m[8]),
		y: v.x*T(m[1]) + v.y*T(m[5]) + v.z*T(m[9]),
		z: v.x*T(m[2]) + v.y*T(m[6]) + v.z*T(m[10]),
		w: v.w, // Direction vectors do not change w
	}
}

// String returns a string representation of the vector, like "(1, 2, 3, 4)".
func (v vector4[T]) String() string {
	return "(" + ToString(v.x) + ", " + ToString(v.y) + ", " + ToString(v.z) + ", " + ToString(v.w) + ")"
}
