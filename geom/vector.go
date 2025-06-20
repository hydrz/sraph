package geom

import "math"

// Vector2 is an alias for a 2D vector, implemented as Point.
type Vector2[T TScalar] = Point[T]

func NewVector2[T TScalar](x, y T) Vector2[T] {
	return Vector2[T]{X: x, Y: y}
}

// Vector3 represents a 3D vector with X, Y, Z components.
type Vector3[T TScalar] struct {
	X T
	Y T
	Z T
}

func NewVector3[T TScalar](x, y, z T) Vector3[T] {
	return Vector3[T]{X: x, Y: y, Z: z}
}

// Add adds another vector to this vector.
func (v Vector3[T]) Add(other Vector3[T]) Vector3[T] {
	return Vector3[T]{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

// Sub subtracts another vector from this vector.
func (v Vector3[T]) Sub(other Vector3[T]) Vector3[T] {
	return Vector3[T]{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

// Mul multiplies this vector by another vector component-wise.
func (v Vector3[T]) Mul(other Vector3[T]) Vector3[T] {
	return Vector3[T]{
		X: v.X * other.X,
		Y: v.Y * other.Y,
		Z: v.Z * other.Z,
	}
}

// Div divides this vector by another vector component-wise.
func (v Vector3[T]) Div(other Vector3[T]) Vector3[T] {
	return Vector3[T]{
		X: v.X / other.X,
		Y: v.Y / other.Y,
		Z: v.Z / other.Z,
	}
}

// Eq checks if this vector is Eq to another vector.
func (v Vector3[T]) Equal(other Vector3[T]) bool {
	return NearlyEqual(v.X, other.X) &&
		NearlyEqual(v.Y, other.Y) &&
		NearlyEqual(v.Z, other.Z)
}

// Scale scales this vector by a scalar.
func (v Vector3[T]) Scale(scalar T) Vector3[T] {
	return Vector3[T]{X: v.X * scalar, Y: v.Y * scalar, Z: v.Z * scalar}
}

// Normalize returns the normalized (unit) vector.
// If the vector's length is zero, returns a zero vector.
// Normalization scales the vector to have a length of 1, while maintaining its direction.
func (v Vector3[T]) Normalize() Vector3[T] {
	len := v.Length()
	var zero T
	if len == zero {
		return Vector3[T]{X: zero, Y: zero, Z: zero}
	}
	return Vector3[T]{X: v.X / len, Y: v.Y / len, Z: v.Z / len}
}

// Length returns the magnitude (Euclidean norm) of the vector.
// It is defined as: length = sqrt(x^2 + y^2 + z^2).
func (v Vector3[T]) Length() T {
	return T(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Abs returns the component-wise absolute value of the vector.
func (v Vector3[T]) Abs() Vector3[T] {
	var zero T
	x, y, z := v.X, v.Y, v.Z
	if x < zero {
		x = -x
	}
	if y < zero {
		y = -y
	}
	if z < zero {
		z = -z
	}
	return Vector3[T]{X: x, Y: y, Z: z}
}

// Floor returns the component-wise floor.
func (v Vector3[T]) Floor() Vector3[T] {
	return Vector3[T]{
		X: T(math.Floor(ToFloat64(v.X))),
		Y: T(math.Floor(ToFloat64(v.Y))),
		Z: T(math.Floor(ToFloat64(v.Z))),
	}
}

// Ceil returns the component-wise ceil.
func (v Vector3[T]) Ceil() Vector3[T] {
	return Vector3[T]{
		X: T(math.Ceil(ToFloat64(v.X))),
		Y: T(math.Ceil(ToFloat64(v.Y))),
		Z: T(math.Ceil(ToFloat64(v.Z))),
	}
}

// Round returns the component-wise round.
func (v Vector3[T]) Round() Vector3[T] {
	return Vector3[T]{
		X: T(math.Round(ToFloat64(v.X))),
		Y: T(math.Round(ToFloat64(v.Y))),
		Z: T(math.Round(ToFloat64(v.Z))),
	}
}

// Dot returns the dot product of this vector and another.
func (v Vector3[T]) Dot(other Vector3[T]) T {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross returns the cross product of this vector and another.
func (v Vector3[T]) Cross(other Vector3[T]) Vector3[T] {
	return Vector3[T]{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// Lerp linearly interpolates between this vector and another by t.
func (v Vector3[T]) Lerp(other Vector3[T], t T) Vector3[T] {
	return Vector3[T]{
		X: v.X + (other.X-v.X)*t,
		Y: v.Y + (other.Y-v.Y)*t,
		Z: v.Z + (other.Z-v.Z)*t,
	}
}

// Combine makes a linear combination of two vectors.
func (v Vector3[T]) Combine(b Vector3[T], bScale T) Vector3[T] {
	return Vector3[T]{
		X: v.X + b.X*bScale,
		Y: v.Y + b.Y*bScale,
		Z: v.Z + b.Z*bScale,
	}
}

// IsZero checks if all components are nearly zero.
func (v Vector3[T]) IsZero() bool {
	return NearlyEqual(v.X, 0) && NearlyEqual(v.Y, 0) && NearlyEqual(v.Z, 0)
}

// Transform applies a transformation matrix to this vector.
func (v Vector3[T]) Transform(m Matrix[T]) Vector3[T] {
	w := v.X*m[3] + v.Y*m[7] + v.Z*m[11] + m[15]
	r := Vector3[T]{
		v.X*m[0] + v.Y*m[4] + v.Z*m[8] + m[12],
		v.X*m[1] + v.Y*m[5] + v.Z*m[9] + m[13],
		v.X*m[2] + v.Y*m[6] + v.Z*m[10] + m[14],
	}
	if w != 0 {
		w = 1 / w
	}
	return r.Scale(w)
}

// TransformDirection applies a transformation matrix to this vector,
// treating it as a direction vector (ignoring translation).
func (v Vector3[T]) TransformDirection(m Matrix[T]) Vector3[T] {
	return Vector3[T]{
		X: v.X*m[0] + v.Y*m[4] + v.Z*m[8],
		Y: v.X*m[1] + v.Y*m[5] + v.Z*m[9],
		Z: v.X*m[2] + v.Y*m[6] + v.Z*m[10],
	}
}

// String returns a string representation of the vector, like "(1, 2, 3)".
func (v Vector3[T]) String() string {
	return "(" + ToString(v.X) + ", " + ToString(v.Y) + ", " + ToString(v.Z) + ")"
}

// Vector4 represents a 4D vector with X, Y, Z, W components.
type Vector4[T TScalar] struct {
	X T
	Y T
	Z T
	W T
}

func NewVector4[T TScalar](x, y, z, w T) Vector4[T] {
	return Vector4[T]{X: x, Y: y, Z: z, W: w}
}

// XY returns the first two components as a Vector2.
func (v Vector4[T]) XY() Vector2[T] {
	return Vector2[T]{X: v.X, Y: v.Y}
}

// Add adds another vector to this vector.
func (v Vector4[T]) Add(other Vector4[T]) Vector4[T] {
	return Vector4[T]{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
		W: v.W + other.W,
	}
}

// Sub subtracts another vector from this vector.
func (v Vector4[T]) Sub(other Vector4[T]) Vector4[T] {
	return Vector4[T]{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
		W: v.W - other.W,
	}
}

// Mul multiplies this vector by another vector component-wise.
func (v Vector4[T]) Mul(other Vector4[T]) Vector4[T] {
	return Vector4[T]{
		X: v.X * other.X,
		Y: v.Y * other.Y,
		Z: v.Z * other.Z,
		W: v.W * other.W,
	}
}

// Div divides this vector by another vector component-wise.
func (v Vector4[T]) Div(other Vector4[T]) Vector4[T] {
	return Vector4[T]{
		X: v.X / other.X,
		Y: v.Y / other.Y,
		Z: v.Z / other.Z,
		W: v.W / other.W,
	}
}

// Eq checks if this vector is Eq to another vector.
func (v Vector4[T]) Equal(other Vector4[T]) bool {
	return NearlyEqual(v.X, other.X) &&
		NearlyEqual(v.Y, other.Y) &&
		NearlyEqual(v.Z, other.Z) &&
		NearlyEqual(v.W, other.W)
}

// Scale scales this vector by a scalar.
func (v Vector4[T]) Scale(scalar T) Vector4[T] {
	return Vector4[T]{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
		W: v.W * scalar,
	}
}

// Normalize returns the normalized (unit) vector.
// If the vector's length is zero, returns a zero vector.
// Normalization scales the vector to have a length of 1, while maintaining its direction.
func (v Vector4[T]) Normalize() Vector4[T] {
	len := v.Length()
	var zero T
	if len == zero {
		return Vector4[T]{X: zero, Y: zero, Z: zero, W: zero}
	}
	return Vector4[T]{
		X: v.X / len,
		Y: v.Y / len,
		Z: v.Z / len,
		W: v.W / len,
	}
}

// Length returns the magnitude (Euclidean norm) of the vector.
// It is defined as: length = sqrt(x^2 + y^2 + z^2 + w^2).
func (v Vector4[T]) Length() T {
	return T(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W)))
}

// Abs returns the component-wise absolute value of the vector.
func (v Vector4[T]) Abs() Vector4[T] {
	var zero T
	x, y, z, w := v.X, v.Y, v.Z, v.W
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
	return Vector4[T]{X: x, Y: y, Z: z, W: w}
}

// Floor returns the component-wise floor.
func (v Vector4[T]) Floor() Vector4[T] {
	return Vector4[T]{
		X: T(math.Floor(ToFloat64(v.X))),
		Y: T(math.Floor(ToFloat64(v.Y))),
		Z: T(math.Floor(ToFloat64(v.Z))),
		W: T(math.Floor(ToFloat64(v.W))),
	}
}

// Ceil returns the component-wise ceil.
func (v Vector4[T]) Ceil() Vector4[T] {
	return Vector4[T]{
		X: T(math.Ceil(ToFloat64(v.X))),
		Y: T(math.Ceil(ToFloat64(v.Y))),
		Z: T(math.Ceil(ToFloat64(v.Z))),
		W: T(math.Ceil(ToFloat64(v.W))),
	}
}

// Round returns the component-wise round.
func (v Vector4[T]) Round() Vector4[T] {
	return Vector4[T]{
		X: T(math.Round(ToFloat64(v.X))),
		Y: T(math.Round(ToFloat64(v.Y))),
		Z: T(math.Round(ToFloat64(v.Z))),
		W: T(math.Round(ToFloat64(v.W))),
	}
}

// Dot returns the dot product of this vector and another.
func (v Vector4[T]) Dot(other Vector4[T]) T {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z + v.W*other.W
}

// Cross returns the cross product of this vector and another.
// Cross product is not well-defined for 4D vectors, so this returns a zero vector.
func (v Vector4[T]) Cross(other Vector4[T]) Vector4[T] {
	var zero T
	return Vector4[T]{X: zero, Y: zero, Z: zero, W: zero}
}

// Lerp linearly interpolates between this vector and another by t.
func (v Vector4[T]) Lerp(other Vector4[T], t T) Vector4[T] {
	return Vector4[T]{
		X: v.X + (other.X-v.X)*t,
		Y: v.Y + (other.Y-v.Y)*t,
		Z: v.Z + (other.Z-v.Z)*t,
		W: v.W + (other.W-v.W)*t,
	}
}

// Combine makes a linear combination of two vectors.
func (v Vector4[T]) Combine(b Vector4[T], bScale T) Vector4[T] {
	return Vector4[T]{
		X: v.X + b.X*bScale,
		Y: v.Y + b.Y*bScale,
		Z: v.Z + b.Z*bScale,
		W: v.W + b.W*bScale,
	}
}

// IsZero checks if all components are nearly zero.
func (v Vector4[T]) IsZero() bool {
	return NearlyEqual(v.X, 0) && NearlyEqual(v.Y, 0) &&
		NearlyEqual(v.Z, 0) && NearlyEqual(v.W, 0)
}

// IsFinite returns true if all components are finite.
func (v Vector4[T]) IsFinite() bool {
	return IsFinite(v.X) && IsFinite(v.Y) && IsFinite(v.Z) && IsFinite(v.W)
}

// Transform applies a transformation matrix to this vector.
func (v Vector4[T]) Transform(m Matrix[T]) Vector4[T] {
	return Vector4[T]{
		v.X*m[0] + v.Y*m[4] + v.Z*m[8] + v.W*m[12],
		v.X*m[1] + v.Y*m[5] + v.Z*m[9] + v.W*m[13],
		v.X*m[2] + v.Y*m[6] + v.Z*m[10] + v.W*m[14],
		v.X*m[3] + v.Y*m[7] + v.Z*m[11] + v.W*m[15],
	}
}

// TransformDirection applies a transformation matrix to this vector,
// treating it as a direction vector (ignoring translation).
func (v Vector4[T]) TransformDirection(m Matrix[T]) Vector4[T] {
	return Vector4[T]{
		X: v.X*m[0] + v.Y*m[4] + v.Z*m[8],
		Y: v.X*m[1] + v.Y*m[5] + v.Z*m[9],
		Z: v.X*m[2] + v.Y*m[6] + v.Z*m[10],
		W: v.W,
	}
}

// String returns a string representation of the vector, like "(1, 2, 3, 4)".
func (v Vector4[T]) String() string {
	return "(" + ToString(v.X) + ", " + ToString(v.Y) + ", " + ToString(v.Z) + ", " + ToString(v.W) + ")"
}
