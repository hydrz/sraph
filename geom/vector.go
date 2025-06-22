package geom

import "math"

// Vector2 defines the interface for a 2D vector.
// Provides methods for vector arithmetic, normalization, and geometric queries.
type Vector2 = Point

// Vector3 defines the interface for a 3D vector.
// Provides methods for vector arithmetic, normalization, and geometric queries.
type Vector3 interface {
	// X returns the X component.
	X() Scalar
	// Y returns the Y component.
	Y() Scalar
	// Z returns the Z component.
	Z() Scalar
	// Add returns the vector sum of this vector and another.
	Add(other Vector3) Vector3
	// Sub returns the vector difference of this vector and another.
	Sub(other Vector3) Vector3
	// Mul returns the element-wise product of this vector and another.
	Mul(other Vector3) Vector3
	// Div returns the element-wise division of this vector by another.
	Div(other Vector3) Vector3
	// Equal reports whether this vector and another are nearly equal.
	Equal(other Vector3) bool
	// Scale returns a new vector scaled by the given scalar.
	Scale(scalar Scalar) Vector3
	// Normalize returns a unit vector in the direction of this vector. If zero, returns zero vector.
	Normalize() Vector3
	// Length returns the Euclidean length of this vector.
	Length() Scalar
	// Abs returns a vector with the absolute value of each component.
	Abs() Vector3
	// Floor returns a vector with math.Floor applied to each component.
	Floor() Vector3
	// Ceil returns a vector with math.Ceil applied to each component.
	Ceil() Vector3
	// Round returns a vector with math.Round applied to each component.
	Round() Vector3
	// Dot returns the dot product of this vector and another.
	Dot(other Vector3) Scalar
	// Cross returns the cross product of this vector and another.
	Cross(other Vector3) Vector3
	// Lerp returns the linear interpolation between this vector and another by parameter t in [0,1].
	Lerp(other Vector3, t Scalar) Vector3
	// IsZero reports whether all components of this vector are nearly zero.
	IsZero() bool
	// Transform applies the given matrix transformation to this vector.
	Transform(m Matrix) Vector3
	// TransformDirection applies the given matrix transformation to this vector, treating it as a direction vector.
	TransformDirection(m Matrix) Vector3
	// String returns a string representation of this vector.
	String() string
	// Combine combines this vector with another by adding the first and multiplying the second by a factor.
	Combine(other Vector3, factor Scalar) Vector3
}

// Vector4 defines the interface for a 4D vector.
type Vector4 interface {
	// X returns the X component.
	X() Scalar
	// Y returns the Y component.
	Y() Scalar
	// Z returns the Z component.
	Z() Scalar
	// W returns the W component.
	W() Scalar
	// Add returns the vector sum of this vector and another.
	Add(other Vector4) Vector4
	// Sub returns the vector difference of this vector and another.
	Sub(other Vector4) Vector4
	// Mul returns the element-wise product of this vector and another.
	Mul(other Vector4) Vector4
	// Div returns the element-wise division of this vector by another.
	Div(other Vector4) Vector4
	// Equal reports whether this vector and another are nearly equal.
	Equal(other Vector4) bool
	// Scale returns a new vector scaled by the given scalar.
	Scale(scalar Scalar) Vector4
	// Normalize returns a unit vector in the direction of this vector. If zero, returns zero vector.
	Normalize() Vector4
	// Length returns the Euclidean length of this vector.
	Length() Scalar
	// Abs returns a vector with the absolute value of each component.
	Abs() Vector4
	// Floor returns a vector with math.Floor applied to each component.
	Floor() Vector4
	// Ceil returns a vector with math.Ceil applied to each component.
	Ceil() Vector4
	// Round returns a vector with math.Round applied to each component.
	Round() Vector4
	// Dot returns the dot product of this vector and another.
	Dot(other Vector4) Scalar
	// Lerp returns the linear interpolation between this vector and another by parameter t in [0,1].
	Lerp(other Vector4, t Scalar) Vector4
	// IsZero reports whether all components of this vector are nearly zero.
	IsZero() bool
	// IsFinite reports whether all components of this vector are finite.
	IsFinite() bool
	// Transform applies the given matrix transformation to this vector.
	Transform(m Matrix) Vector4
	// TransformDirection applies the given matrix transformation to this vector, treating it as a direction vector.
	TransformDirection(m Matrix) Vector4
	// String returns a string representation of this vector.
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

// X implements Vector3.X.
func (v vector3[T]) X() Scalar {
	return ToScalar(v.x)
}

// Y implements Vector3.Y.
func (v vector3[T]) Y() Scalar {
	return ToScalar(v.y)
}

// Z implements Vector3.Z.
func (v vector3[T]) Z() Scalar {
	return ToScalar(v.z)
}

// Add implements Vector3.Add.
func (v vector3[T]) Add(o Vector3) Vector3 {
	return NewVector3(v.X()+o.X(), v.Y()+o.Y(), v.Z()+o.Z())
}

// Sub implements Vector3.Sub.
func (v vector3[T]) Sub(other Vector3) Vector3 {
	return NewVector3(v.X()-other.X(), v.Y()-other.Y(), v.Z()-other.Z())
}

// Mul implements Vector3.Mul.
func (v vector3[T]) Mul(other Vector3) Vector3 {
	return NewVector3(v.X()*other.X(), v.Y()*other.Y(), v.Z()*other.Z())
}

// Div implements Vector3.Div.
func (v vector3[T]) Div(other Vector3) Vector3 {
	return NewVector3(v.X()/other.X(), v.Y()/other.Y(), v.Z()/other.Z())
}

// Equal implements Vector3.Equal.
func (v vector3[T]) Equal(o Vector3) bool {
	return NearlyEqual(v.X(), o.X()) && NearlyEqual(v.Y(), o.Y()) && NearlyEqual(v.Z(), o.Z())
}

// Scale implements Vector3.Scale.
func (v vector3[T]) Scale(scalar Scalar) Vector3 {
	return NewVector3(v.X()*scalar, v.Y()*scalar, v.Z()*scalar)
}

// Normalize implements Vector3.Normalize.
func (v vector3[T]) Normalize() Vector3 {
	len := v.Length()
	if len == 0 {
		return NewVector3[T](0, 0, 0)
	}
	return NewVector3(v.X()/len, v.Y()/len, v.Z()/len)
}

// Length implements Vector3.Length.
func (v vector3[T]) Length() Scalar {
	return ToScalar(math.Sqrt(ToFloat64(v.x*v.x + v.y*v.y + v.z*v.z)))
}

// Abs implements Vector3.Abs.
func (v vector3[T]) Abs() Vector3 {
	return NewVector3(Abs(v.x), Abs(v.y), Abs(v.z))
}

// Floor implements Vector3.Floor.
func (v vector3[T]) Floor() Vector3 {
	return vector3[T]{
		x: T(math.Floor(ToFloat64(v.x))),
		y: T(math.Floor(ToFloat64(v.y))),
		z: T(math.Floor(ToFloat64(v.z))),
	}
}

// Ceil implements Vector3.Ceil.
func (v vector3[T]) Ceil() Vector3 {
	return vector3[T]{
		x: T(math.Ceil(ToFloat64(v.x))),
		y: T(math.Ceil(ToFloat64(v.y))),
		z: T(math.Ceil(ToFloat64(v.z))),
	}
}

// Round implements Vector3.Round.
func (v vector3[T]) Round() Vector3 {
	return vector3[T]{
		x: T(math.Round(ToFloat64(v.x))),
		y: T(math.Round(ToFloat64(v.y))),
		z: T(math.Round(ToFloat64(v.z))),
	}
}

// Dot implements Vector3.Dot.
func (v vector3[T]) Dot(o Vector3) Scalar {
	return v.X()*o.X() + v.Y()*o.Y() + v.Z()*o.Z()
}

// Cross implements Vector3.Cross.
func (v vector3[T]) Cross(other Vector3) Vector3 {
	return NewVector3(
		v.Y()*other.Z()-v.Z()*other.Y(),
		v.Z()*other.X()-v.X()*other.Z(),
		v.X()*other.Y()-v.Y()*other.X(),
	)
}

// Lerp implements Vector3.Lerp.
func (v vector3[T]) Lerp(other Vector3, t Scalar) Vector3 {
	return NewVector3(
		v.X()+(other.X()-v.X())*t,
		v.Y()+(other.Y()-v.Y())*t,
		v.Z()+(other.Z()-v.Z())*t,
	)
}

// IsZero implements Vector3.IsZero.
func (v vector3[T]) IsZero() bool {
	return NearlyEqual(v.x, 0) && NearlyEqual(v.y, 0) && NearlyEqual(v.z, 0)
}

// Transform implements Vector3.Transform.
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

// TransformDirection implements Vector3.TransformDirection.
func (v vector3[T]) TransformDirection(m Matrix) Vector3 {
	return NewVector3(
		v.X()*m[0]+v.Y()*m[4]+v.Z()*m[8],
		v.X()*m[1]+v.Y()*m[5]+v.Z()*m[9],
		v.X()*m[2]+v.Y()*m[6]+v.Z()*m[10],
	)
}

// String implements Vector3.String.
func (v vector3[T]) String() string {
	return "(" +
		ToString(v.x) + ", " +
		ToString(v.y) + ", " +
		ToString(v.z) + ")"
}

// Combine implements Vector3.Combine.
func (v vector3[T]) Combine(other Vector3, factor Scalar) Vector3 {
	return NewVector3(
		v.X()+other.X()*factor,
		v.Y()+other.Y()*factor,
		v.Z()+other.Z()*factor,
	)
}

type vector4[T Number] struct {
	x, y, z, w T
}

// X implements Vector4.X.
func (v vector4[T]) X() Scalar {
	return ToScalar(v.x)
}

// Y implements Vector4.Y.
func (v vector4[T]) Y() Scalar {
	return ToScalar(v.y)
}

// Z implements Vector4.Z.
func (v vector4[T]) Z() Scalar {
	return ToScalar(v.z)
}

// W implements Vector4.W.
func (v vector4[T]) W() Scalar {
	return ToScalar(v.w)
}

// Add implements Vector4.Add.
func (v vector4[T]) Add(other Vector4) Vector4 {
	return NewVector4(v.X()+other.X(), v.Y()+other.Y(), v.Z()+other.Z(), v.W()+other.W())
}

// Sub implements Vector4.Sub.
func (v vector4[T]) Sub(other Vector4) Vector4 {
	return NewVector4(v.X()-other.X(), v.Y()-other.Y(), v.Z()-other.Z(), v.W()-other.W())
}

// Mul implements Vector4.Mul.
func (v vector4[T]) Mul(other Vector4) Vector4 {
	return NewVector4(v.X()*other.X(), v.Y()*other.Y(), v.Z()*other.Z(), v.W()*other.W())
}

// Div implements Vector4.Div.
func (v vector4[T]) Div(other Vector4) Vector4 {
	return NewVector4(v.X()/other.X(), v.Y()/other.Y(), v.Z()/other.Z(), v.W()/other.W())
}

// Equal implements Vector4.Equal.
func (v vector4[T]) Equal(other Vector4) bool {
	return NearlyEqual(v.X(), other.X()) &&
		NearlyEqual(v.Y(), other.Y()) &&
		NearlyEqual(v.Z(), other.Z()) &&
		NearlyEqual(v.W(), other.W())
}

// Scale implements Vector4.Scale.
func (v vector4[T]) Scale(scalar Scalar) Vector4 {
	return NewVector4(v.X()*scalar, v.Y()*scalar, v.Z()*scalar, v.W()*scalar)
}

// Normalize implements Vector4.Normalize.
func (v vector4[T]) Normalize() Vector4 {
	len := v.Length()
	if len == 0 {
		return NewVector4[T](0, 0, 0, 0)
	}
	return NewVector4(v.X()/len, v.Y()/len, v.Z()/len, v.W()/len)
}

// Length implements Vector4.Length.
func (v vector4[T]) Length() Scalar {
	return ToScalar(math.Sqrt(ToFloat64(v.x*v.x + v.y*v.y + v.z*v.z + v.w*v.w)))
}

// Abs implements Vector4.Abs.
func (v vector4[T]) Abs() Vector4 {
	return NewVector4(Abs(v.x), Abs(v.y), Abs(v.z), Abs(v.w))
}

// Floor implements Vector4.Floor.
func (v vector4[T]) Floor() Vector4 {
	return vector4[T]{
		x: T(math.Floor(ToFloat64(v.x))),
		y: T(math.Floor(ToFloat64(v.y))),
		z: T(math.Floor(ToFloat64(v.z))),
		w: T(math.Floor(ToFloat64(v.w))),
	}
}

// Ceil implements Vector4.Ceil.
func (v vector4[T]) Ceil() Vector4 {
	return vector4[T]{
		x: T(math.Ceil(ToFloat64(v.x))),
		y: T(math.Ceil(ToFloat64(v.y))),
		z: T(math.Ceil(ToFloat64(v.z))),
		w: T(math.Ceil(ToFloat64(v.w))),
	}
}

// Round implements Vector4.Round.
func (v vector4[T]) Round() Vector4 {
	return vector4[T]{
		x: T(math.Round(ToFloat64(v.x))),
		y: T(math.Round(ToFloat64(v.y))),
		z: T(math.Round(ToFloat64(v.z))),
		w: T(math.Round(ToFloat64(v.w))),
	}
}

// Dot implements Vector4.Dot.
func (v vector4[T]) Dot(o Vector4) Scalar {
	return v.X()*o.X() + v.Y()*o.Y() + v.Z()*o.Z() + v.W()*o.W()
}

// Lerp implements Vector4.Lerp.
func (v vector4[T]) Lerp(o Vector4, t Scalar) Vector4 {
	return NewVector4(
		v.X()+(o.X()-v.X())*t,
		v.Y()+(o.Y()-v.Y())*t,
		v.Z()+(o.Z()-v.Z())*t,
		v.W()+(o.W()-v.W())*t,
	)
}

// IsZero implements Vector4.IsZero.
func (v vector4[T]) IsZero() bool {
	return NearlyEqual(v.x, 0) && NearlyEqual(v.y, 0) &&
		NearlyEqual(v.z, 0) && NearlyEqual(v.w, 0)
}

// IsFinite implements Vector4.IsFinite.
func (v vector4[T]) IsFinite() bool {
	return IsFinite(v.x) && IsFinite(v.y) && IsFinite(v.z) && IsFinite(v.w)
}

// Transform implements Vector4.Transform.
func (v vector4[T]) Transform(m Matrix) Vector4 {
	return NewVector4(
		v.X()*m[0]+v.Y()*m[4]+v.Z()*m[8]+v.W()*m[12],
		v.X()*m[1]+v.Y()*m[5]+v.Z()*m[9]+v.W()*m[13],
		v.X()*m[2]+v.Y()*m[6]+v.Z()*m[10]+v.W()*m[14],
		v.X()*m[3]+v.Y()*m[7]+v.Z()*m[11]+v.W()*m[15],
	)
}

// TransformDirection implements Vector4.TransformDirection.
func (v vector4[T]) TransformDirection(m Matrix) Vector4 {
	return NewVector4(
		v.X()*m[0]+v.Y()*m[4]+v.Z()*m[8],
		v.X()*m[1]+v.Y()*m[5]+v.Z()*m[9],
		v.X()*m[2]+v.Y()*m[6]+v.Z()*m[10],
		v.W(), // Direction vectors do not change w
	)
}

// String implements Vector4.String.
func (v vector4[T]) String() string {
	return "(" +
		ToString(v.x) + ", " +
		ToString(v.y) + ", " +
		ToString(v.z) + ", " +
		ToString(v.w) + ")"
}
