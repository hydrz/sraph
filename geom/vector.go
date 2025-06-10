package geom

import "math"

type Vector2[T Scalar] = Point[T]

func NewVector2[T Scalar](x, y T) Vector2[T] {
	return Vector2[T](NewPoint(x, y))
}

type Vector3[T Scalar] interface {
	// X returns the x component of the vector.
	X() T
	// Y returns the y component of the vector.
	Y() T
	// Z returns the z component of the vector.
	Z() T

	// Add adds another vector to this vector.
	Add(other Vector3[T]) Vector3[T]
	// Sub subtracts another vector from this vector.
	Sub(other Vector3[T]) Vector3[T]
	// Mul multiplies this vector by another vector component-wise.
	Mul(other Vector3[T]) Vector3[T]
	// Div divides this vector by another vector component-wise.
	Div(other Vector3[T]) Vector3[T]
	// Equal checks if this vector is equal to another vector.
	Equal(other Vector3[T]) bool

	// Scale scales this vector by a scalar.
	Scale(scalar T) Vector3[T]

	// Normalize returns the normalized (unit) vector.
	// If the vector's length is zero, returns a zero vector.
	// Normalization scales the vector to have a length of 1, while maintaining its direction.
	Normalize() Vector3[T]

	// Length returns the magnitude (Euclidean norm) of the vector.
	// It is defined as: length = sqrt(x^2 + y^2 + z^2).
	Length() T

	// Abs returns the component-wise absolute value of the vector.
	Abs() Vector3[T]
	// Floor returns the component-wise floor.
	Floor() Vector3[T]
	// Ceil returns the component-wise ceil.
	Ceil() Vector3[T]
	// Round returns the component-wise round.
	Round() Vector3[T]
	// Min returns the component-wise minimum with another vector.
	Min(other ...Vector3[T]) Vector3[T]
	// Max returns the component-wise maximum with another vector.
	Max(other ...Vector3[T]) Vector3[T]

	// Dot returns the dot product of this vector and another.
	Dot(other Vector3[T]) T
	// Cross returns the cross product of this vector and another.
	Cross(other Vector3[T]) Vector3[T]

	// Lerp linearly interpolates between this vector and another by t.
	Lerp(other Vector3[T], t T) Vector3[T]
	// Combine makes a linear combination of two vectors.
	Combine(b Vector3[T], bScale T) Vector3[T]

	// String returns a string representation of the vector. like "Vector3(1, 2, 3)".
	String() string
}

func NewVector3[T Scalar](x, y, z T) Vector3[T] {
	return vector3[T]{x: x, y: y, z: z}
}

type Vector4[T Scalar] interface {
	// X returns the x component of the vector.
	X() T
	// Y returns the y component of the vector.
	Y() T
	// Z returns the z component of the vector.
	Z() T
	// W returns the w component of the vector.
	W() T
	// XY returns the first two components as a Vector2.
	XY() Vector2[T]

	// Add adds another vector to this vector.
	Add(other Vector4[T]) Vector4[T]
	// Sub subtracts another vector from this vector.
	Sub(other Vector4[T]) Vector4[T]
	// Mul multiplies this vector by another vector component-wise.
	Mul(other Vector4[T]) Vector4[T]
	// Div divides this vector by another vector component-wise.
	Div(other Vector4[T]) Vector4[T]
	// Equal checks if this vector is equal to another vector.
	Equal(other Vector4[T]) bool

	// Scale scales this vector by a scalar.
	Scale(scalar T) Vector4[T]

	// Normalize returns the normalized (unit) vector.
	// If the vector's length is zero, returns a zero vector.
	// Normalization scales the vector to have a length of 1, while maintaining its direction.
	Normalize() Vector4[T]

	// Length returns the magnitude (Euclidean norm) of the vector.
	// It is defined as: length = sqrt(x^2 + y^2 + z^2).
	Length() T

	// Abs returns the component-wise absolute value of the vector.
	Abs() Vector4[T]
	// Floor returns the component-wise floor.
	Floor() Vector4[T]
	// Ceil returns the component-wise ceil.
	Ceil() Vector4[T]
	// Round returns the component-wise round.
	Round() Vector4[T]
	// Min returns the component-wise minimum with another vector.
	Min(other ...Vector4[T]) Vector4[T]
	// Max returns the component-wise maximum with another vector.
	Max(other ...Vector4[T]) Vector4[T]

	// Dot returns the dot product of this vector and another.
	Dot(other Vector4[T]) T
	// Cross returns the cross product of this vector and another.
	Cross(other Vector4[T]) Vector4[T]

	// Lerp linearly interpolates between this vector and another by t.
	Lerp(other Vector4[T], t T) Vector4[T]
	// Combine makes a linear combination of two vectors.
	Combine(b Vector4[T], bScale T) Vector4[T]

	// String returns a string representation of the vector. like "Vector4(1, 2, 3, 4)".
	String() string

	// IsFinite returns true if all components are finite.
	IsFinite() bool
}

func NewVector4[T Scalar](x, y, z, w T) Vector4[T] {
	return vector4[T]{x: x, y: y, z: z, w: w}
}

type vector3[T Scalar] struct {
	x, y, z T
}

// X implements Vector3.
func (v vector3[T]) X() T {
	return v.x
}

// Y implements Vector3.
func (v vector3[T]) Y() T {
	return v.y
}

// Z implements Vector3.
func (v vector3[T]) Z() T {
	return v.z
}

// Add implements Vector3.
func (v vector3[T]) Add(other Vector3[T]) Vector3[T] {
	return vector3[T]{
		x: v.x + other.X(),
		y: v.y + other.Y(),
		z: v.z + other.Z(),
	}
}

// Sub implements Vector3.
func (v vector3[T]) Sub(other Vector3[T]) Vector3[T] {
	return vector3[T]{
		x: v.x - other.X(),
		y: v.y - other.Y(),
		z: v.z - other.Z(),
	}
}

// Mul implements Vector3.
func (v vector3[T]) Mul(other Vector3[T]) Vector3[T] {
	return vector3[T]{
		x: v.x * other.X(),
		y: v.y * other.Y(),
		z: v.z * other.Z(),
	}
}

// Div implements Vector3.
func (v vector3[T]) Div(other Vector3[T]) Vector3[T] {
	return vector3[T]{
		x: v.x / other.X(),
		y: v.y / other.Y(),
		z: v.z / other.Z(),
	}
}

// Equal implements Vector3.
func (v vector3[T]) Equal(other Vector3[T]) bool {
	return Equal(v.x, other.X()) &&
		Equal(v.y, other.Y()) &&
		Equal(v.z, other.Z())
}

// Scale implements Vector3.
func (v vector3[T]) Scale(scalar T) Vector3[T] {
	return vector3[T]{x: v.x * scalar, y: v.y * scalar, z: v.z * scalar}
}

// Normalize implements Vector3.
func (v vector3[T]) Normalize() Vector3[T] {
	len := v.Length()
	var zero T
	if len == zero {
		return vector3[T]{x: zero, y: zero, z: zero}
	}
	return vector3[T]{x: v.x / len, y: v.y / len, z: v.z / len}
}

// Length implements Vector3.
func (v vector3[T]) Length() T {
	return T(math.Sqrt(float64(v.x*v.x + v.y*v.y + v.z*v.z)))
}

// Abs implements Vector3.
func (v vector3[T]) Abs() Vector3[T] {
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

// Floor implements Vector3.
func (v vector3[T]) Floor() Vector3[T] {
	return vector3[T]{
		x: T(math.Floor(v.x.Float64())),
		y: T(math.Floor(v.y.Float64())),
		z: T(math.Floor(v.z.Float64())),
	}
}

// Ceil implements Vector3.
func (v vector3[T]) Ceil() Vector3[T] {
	return vector3[T]{
		x: T(math.Ceil(v.x.Float64())),
		y: T(math.Ceil(v.y.Float64())),
		z: T(math.Ceil(v.z.Float64())),
	}
}

// Round implements Vector3.
func (v vector3[T]) Round() Vector3[T] {
	return vector3[T]{
		x: T(math.Round(v.x.Float64())),
		y: T(math.Round(v.y.Float64())),
		z: T(math.Round(v.z.Float64())),
	}
}

// Min implements Vector3.
func (v vector3[T]) Min(other ...Vector3[T]) Vector3[T] {
	minX, minY, minZ := v.x, v.y, v.z
	for _, o := range other {
		if o.X() < minX {
			minX = o.X()
		}
		if o.Y() < minY {
			minY = o.Y()
		}
		if o.Z() < minZ {
			minZ = o.Z()
		}
	}
	return vector3[T]{x: minX, y: minY, z: minZ}
}

// Max implements Vector3.
func (v vector3[T]) Max(other ...Vector3[T]) Vector3[T] {
	maxX, maxY, maxZ := v.x, v.y, v.z
	for _, o := range other {
		if o.X() > maxX {
			maxX = o.X()
		}
		if o.Y() > maxY {
			maxY = o.Y()
		}
		if o.Z() > maxZ {
			maxZ = o.Z()
		}
	}
	return vector3[T]{x: maxX, y: maxY, z: maxZ}
}

// Dot implements Vector3.
func (v vector3[T]) Dot(other Vector3[T]) T {
	return v.x*other.X() + v.y*other.Y() + v.z*other.Z()
}

// Cross implements Vector3.
func (v vector3[T]) Cross(other Vector3[T]) Vector3[T] {
	return vector3[T]{
		x: v.y*other.Z() - v.z*other.Y(),
		y: v.z*other.X() - v.x*other.Z(),
		z: v.x*other.Y() - v.y*other.X(),
	}
}

// Lerp implements Vector3.
func (v vector3[T]) Lerp(other Vector3[T], t T) Vector3[T] {
	return vector3[T]{
		x: v.x + (other.X()-v.x)*t,
		y: v.y + (other.Y()-v.y)*t,
		z: v.z + (other.Z()-v.z)*t,
	}
}

// Combine implements Vector3.
func (v vector3[T]) Combine(b Vector3[T], bScale T) Vector3[T] {
	return vector3[T]{
		x: v.x + b.X()*bScale,
		y: v.y + b.Y()*bScale,
		z: v.z + b.Z()*bScale,
	}
}

// String implements Vector3.
func (v vector3[T]) String() string {
	return "Vector3(" + v.x.String() + ", " + v.y.String() + ", " + v.z.String() + ")"
}

type vector4[T Scalar] struct {
	x, y, z, w T
}

// X implements Vector4.
func (v vector4[T]) X() T {
	return v.x
}

// Y implements Vector4.
func (v vector4[T]) Y() T {
	return v.y
}

// Z implements Vector4.
func (v vector4[T]) Z() T {
	return v.z
}

// W implements Vector4.
func (v vector4[T]) W() T {
	return v.w
}

// XY implements Vector4.
func (v vector4[T]) XY() Vector2[T] {
	return NewVector2(v.x, v.y)
}

// Add implements Vector4.
func (v vector4[T]) Add(other Vector4[T]) Vector4[T] {
	return vector4[T]{
		x: v.x + other.X(),
		y: v.y + other.Y(),
		z: v.z + other.Z(),
		w: v.w + other.W(),
	}
}

// Sub implements Vector4.
func (v vector4[T]) Sub(other Vector4[T]) Vector4[T] {
	return vector4[T]{
		x: v.x - other.X(),
		y: v.y - other.Y(),
		z: v.z - other.Z(),
		w: v.w - other.W(),
	}
}

// Mul implements Vector4.
func (v vector4[T]) Mul(other Vector4[T]) Vector4[T] {
	return vector4[T]{
		x: v.x * other.X(),
		y: v.y * other.Y(),
		z: v.z * other.Z(),
		w: v.w * other.W(),
	}
}

// Div implements Vector4.
func (v vector4[T]) Div(other Vector4[T]) Vector4[T] {
	return vector4[T]{
		x: v.x / other.X(),
		y: v.y / other.Y(),
		z: v.z / other.Z(),
		w: v.w / other.W(),
	}
}

// Equal implements Vector4.
func (v vector4[T]) Equal(other Vector4[T]) bool {
	return Equal(v.x, other.X()) &&
		Equal(v.y, other.Y()) &&
		Equal(v.z, other.Z()) &&
		Equal(v.w, other.W())
}

// Scale implements Vector4.
func (v vector4[T]) Scale(scalar T) Vector4[T] {
	return vector4[T]{
		x: v.x * scalar,
		y: v.y * scalar,
		z: v.z * scalar,
		w: v.w * scalar,
	}
}

// Normalize implements Vector4.
func (v vector4[T]) Normalize() Vector4[T] {
	len := v.Length()
	var zero T
	if len == zero {
		return vector4[T]{x: zero, y: zero, z: zero, w: zero}
	}
	return vector4[T]{
		x: v.x / len,
		y: v.y / len,
		z: v.z / len,
		w: v.w / len,
	}
}

// Length implements Vector4.
func (v vector4[T]) Length() T {
	return T(math.Sqrt(float64(v.x*v.x + v.y*v.y + v.z*v.z + v.w*v.w)))
}

// Abs implements Vector4.
func (v vector4[T]) Abs() Vector4[T] {
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

// Floor implements Vector4.
func (v vector4[T]) Floor() Vector4[T] {
	return vector4[T]{
		x: T(math.Floor(v.x.Float64())),
		y: T(math.Floor(v.y.Float64())),
		z: T(math.Floor(v.z.Float64())),
		w: T(math.Floor(v.w.Float64())),
	}
}

// Ceil implements Vector4.
func (v vector4[T]) Ceil() Vector4[T] {
	return vector4[T]{
		x: T(math.Ceil(v.x.Float64())),
		y: T(math.Ceil(v.y.Float64())),
		z: T(math.Ceil(v.z.Float64())),
		w: T(math.Ceil(v.w.Float64())),
	}
}

// Round implements Vector4.
func (v vector4[T]) Round() Vector4[T] {
	return vector4[T]{
		x: T(math.Round(v.x.Float64())),
		y: T(math.Round(v.y.Float64())),
		z: T(math.Round(v.z.Float64())),
		w: T(math.Round(v.w.Float64())),
	}
}

// Min implements Vector4.
func (v vector4[T]) Min(other ...Vector4[T]) Vector4[T] {
	minX, minY, minZ, minW := v.x, v.y, v.z, v.w
	for _, o := range other {
		if o.X() < minX {
			minX = o.X()
		}
		if o.Y() < minY {
			minY = o.Y()
		}
		if o.Z() < minZ {
			minZ = o.Z()
		}
		if o.W() < minW {
			minW = o.W()
		}
	}
	return vector4[T]{x: minX, y: minY, z: minZ, w: minW}
}

// Max implements Vector4.
func (v vector4[T]) Max(other ...Vector4[T]) Vector4[T] {
	maxX, maxY, maxZ, maxW := v.x, v.y, v.z, v.w
	for _, o := range other {
		if o.X() > maxX {
			maxX = o.X()
		}
		if o.Y() > maxY {
			maxY = o.Y()
		}
		if o.Z() > maxZ {
			maxZ = o.Z()
		}
		if o.W() > maxW {
			maxW = o.W()
		}
	}
	return vector4[T]{x: maxX, y: maxY, z: maxZ, w: maxW}
}

// Dot implements Vector4.
func (v vector4[T]) Dot(other Vector4[T]) T {
	return v.x*other.X() + v.y*other.Y() + v.z*other.Z() + v.w*other.W()
}

// Cross implements Vector4.
func (v vector4[T]) Cross(other Vector4[T]) Vector4[T] {
	// Cross product is not well-defined for 4D vectors.
	// Here we return a zero vector.
	var zero T
	return vector4[T]{x: zero, y: zero, z: zero, w: zero}
}

// Lerp implements Vector4.
func (v vector4[T]) Lerp(other Vector4[T], t T) Vector4[T] {
	return vector4[T]{
		x: v.x + (other.X()-v.x)*t,
		y: v.y + (other.Y()-v.y)*t,
		z: v.z + (other.Z()-v.z)*t,
		w: v.w + (other.W()-v.w)*t,
	}
}

// Combine implements Vector4.
func (v vector4[T]) Combine(b Vector4[T], bScale T) Vector4[T] {
	return vector4[T]{
		x: v.x + b.X()*bScale,
		y: v.y + b.Y()*bScale,
		z: v.z + b.Z()*bScale,
		w: v.w + b.W()*bScale,
	}
}

// String implements Vector4.
func (v vector4[T]) String() string {
	return "Vector4(" + v.x.String() + ", " + v.y.String() + ", " + v.z.String() + ", " + v.w.String() + ")"
}

// IsFinite implements Vector4.
func (v vector4[T]) IsFinite() bool {
	return IsFinite(v.x) && IsFinite(v.y) && IsFinite(v.z) && IsFinite(v.w)
}
