package geom

import "math"

// Vector2 is an alias for a 2D vector, implemented as Point.
type Vector2[T Scalar] = Point[T]

// Vector3 represents a 3D vector with X, Y, Z components.
type Vector3[T Scalar] struct {
	X T
	Y T
	Z T
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
func (v Vector3[T]) Eq(other Vector3[T]) bool {
	return Eq(v.X, other.X) &&
		Eq(v.Y, other.Y) &&
		Eq(v.Z, other.Z)
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
		X: T(math.Floor(v.X.Float64())),
		Y: T(math.Floor(v.Y.Float64())),
		Z: T(math.Floor(v.Z.Float64())),
	}
}

// Ceil returns the component-wise ceil.
func (v Vector3[T]) Ceil() Vector3[T] {
	return Vector3[T]{
		X: T(math.Ceil(v.X.Float64())),
		Y: T(math.Ceil(v.Y.Float64())),
		Z: T(math.Ceil(v.Z.Float64())),
	}
}

// Round returns the component-wise round.
func (v Vector3[T]) Round() Vector3[T] {
	return Vector3[T]{
		X: T(math.Round(v.X.Float64())),
		Y: T(math.Round(v.Y.Float64())),
		Z: T(math.Round(v.Z.Float64())),
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

// String returns a string representation of the vector, like "(1, 2, 3)".
func (v Vector3[T]) String() string {
	return "(" + v.X.String() + ", " + v.Y.String() + ", " + v.Z.String() + ")"
}

// Vector4 represents a 4D vector with X, Y, Z, W components.
type Vector4[T Scalar] struct {
	X T
	Y T
	Z T
	W T
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
func (v Vector4[T]) Eq(other Vector4[T]) bool {
	return Eq(v.X, other.X) &&
		Eq(v.Y, other.Y) &&
		Eq(v.Z, other.Z) &&
		Eq(v.W, other.W)
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
		X: T(math.Floor(v.X.Float64())),
		Y: T(math.Floor(v.Y.Float64())),
		Z: T(math.Floor(v.Z.Float64())),
		W: T(math.Floor(v.W.Float64())),
	}
}

// Ceil returns the component-wise ceil.
func (v Vector4[T]) Ceil() Vector4[T] {
	return Vector4[T]{
		X: T(math.Ceil(v.X.Float64())),
		Y: T(math.Ceil(v.Y.Float64())),
		Z: T(math.Ceil(v.Z.Float64())),
		W: T(math.Ceil(v.W.Float64())),
	}
}

// Round returns the component-wise round.
func (v Vector4[T]) Round() Vector4[T] {
	return Vector4[T]{
		X: T(math.Round(v.X.Float64())),
		Y: T(math.Round(v.Y.Float64())),
		Z: T(math.Round(v.Z.Float64())),
		W: T(math.Round(v.W.Float64())),
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

// String returns a string representation of the vector, like "(1, 2, 3, 4)".
func (v Vector4[T]) String() string {
	return "(" + v.X.String() + ", " + v.Y.String() + ", " + v.Z.String() + ", " + v.W.String() + ")"
}

// IsFinite returns true if all components are finite.
func (v Vector4[T]) IsFinite() bool {
	return IsFinite(v.X) && IsFinite(v.Y) && IsFinite(v.Z) && IsFinite(v.W)
}
