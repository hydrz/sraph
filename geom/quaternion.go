package geom

import (
	"fmt"
	"math"
)

// Quaternion defines the interface for a quaternion representing 3D rotations.
// Provides methods for quaternion arithmetic, normalization, and vector rotation.
type Quaternion interface {
	// X returns the X component of the quaternion.
	X() Scalar
	// Y returns the Y component of the quaternion.
	Y() Scalar
	// Z returns the Z component of the quaternion.
	Z() Scalar
	// W returns the W component (scalar part) of the quaternion.
	W() Scalar

	// Add returns the component-wise sum of this quaternion and another.
	Add(other Quaternion) Quaternion
	// Sub returns the component-wise difference of this quaternion and another.
	Sub(other Quaternion) Quaternion
	// Mul returns the Hamilton product of this quaternion and another.
	// This corresponds to the composition of two 3D rotations.
	Mul(other Quaternion) Quaternion
	// Div returns the quotient of this quaternion divided by another.
	// This is equivalent to multiplying by the inverse of the other quaternion.
	Div(other Quaternion) Quaternion

	// Equal reports whether this quaternion and another are equal within floating-point tolerance.
	Equal(other Quaternion) bool
	// Scale returns this quaternion scaled by the given scalar value.
	Scale(scalar Scalar) Quaternion
	// Neg returns the negation of this quaternion (all components negated).
	Neg() Quaternion

	// Length returns the magnitude (norm) of the quaternion.
	// For unit quaternions representing rotations, this should be 1.
	Length() Scalar
	// Dot returns the dot product of this quaternion and another.
	Dot(other Quaternion) Scalar
	// Normalize returns a unit quaternion in the same direction as this quaternion.
	// If this quaternion has zero length, returns a zero quaternion.
	Normalize() Quaternion
	// Invert returns the inverse (conjugate divided by squared magnitude) of this quaternion.
	// For unit quaternions, this is equivalent to the conjugate.
	Invert() Quaternion

	// Slerp performs spherical linear interpolation between this quaternion and another.
	// The parameter t should be in the range [0,1], where t=0 returns this quaternion
	// and t=1 returns the other quaternion.
	Slerp(other Quaternion, t float64) Quaternion
	// RotateVector3 applies this quaternion rotation to a 3D vector and returns the rotated vector.
	RotateVector3(vector Vector3) Vector3

	// String returns a string representation of the quaternion in the form "(x, y, z, w)".
	String() string
}

// quaternion is the generic implementation of the Quaternion interface.
type quaternion[T Number] struct {
	x, y, z, w T
}

// NewQuaternion creates a new quaternion with the given components.
func NewQuaternion(x, y, z, w Scalar) Quaternion {
	return quaternion[Scalar]{x, y, z, w}
}

// NewQuaternionFromAxisAngle creates a new quaternion from an axis and angle.
func NewQuaternionFromAxisAngle(axis Vector3, angle Radians) Quaternion {
	axis = axis.Normalize()
	halfAngle := ToFloat64(angle) / 2
	sinHalfAngle := Scalar(math.Sin(halfAngle))
	cosHalfAngle := Scalar(math.Cos(halfAngle))
	return quaternion[Scalar]{
		x: axis.X() * sinHalfAngle,
		y: axis.Y() * sinHalfAngle,
		z: axis.Z() * sinHalfAngle,
		w: cosHalfAngle,
	}
}

// X implements Quaternion.X.
func (q quaternion[T]) X() Scalar { return Scalar(q.x) }

// Y implements Quaternion.Y.
func (q quaternion[T]) Y() Scalar { return Scalar(q.y) }

// Z implements Quaternion.Z.
func (q quaternion[T]) Z() Scalar { return Scalar(q.z) }

// W implements Quaternion.W.
func (q quaternion[T]) W() Scalar { return Scalar(q.w) }

// Add implements Quaternion.Add.
func (q quaternion[T]) Add(other Quaternion) Quaternion {
	return quaternion[T]{
		x: q.x + T(other.X()),
		y: q.y + T(other.Y()),
		z: q.z + T(other.Z()),
		w: q.w + T(other.W()),
	}
}

// Sub implements Quaternion.Sub.
func (q quaternion[T]) Sub(other Quaternion) Quaternion {
	return quaternion[T]{
		x: q.x - T(other.X()),
		y: q.y - T(other.Y()),
		z: q.z - T(other.Z()),
		w: q.w - T(other.W()),
	}
}

// Mul implements Quaternion.Mul (Hamilton product).
func (q quaternion[T]) Mul(other Quaternion) Quaternion {
	x1, y1, z1, w1 := q.x, q.y, q.z, q.w
	x2, y2, z2, w2 := T(other.X()), T(other.Y()), T(other.Z()), T(other.W())
	return quaternion[T]{
		x: w1*x2 + x1*w2 + y1*z2 - z1*y2,
		y: w1*y2 - x1*z2 + y1*w2 + z1*x2,
		z: w1*z2 + x1*y2 - y1*x2 + z1*w2,
		w: w1*w2 - x1*x2 - y1*y2 - z1*z2,
	}
}

// Div implements Quaternion.Div.
func (q quaternion[T]) Div(other Quaternion) Quaternion {
	return q.Mul(other.Invert())
}

// Equal implements Quaternion.Equal.
func (q quaternion[T]) Equal(other Quaternion) bool {
	return NearlyEqual(q.x, T(other.X())) &&
		NearlyEqual(q.y, T(other.Y())) &&
		NearlyEqual(q.z, T(other.Z())) &&
		NearlyEqual(q.w, T(other.W()))
}

// Scale implements Quaternion.Scale.
func (q quaternion[T]) Scale(scalar Scalar) Quaternion {
	return quaternion[T]{
		x: q.x * T(scalar),
		y: q.y * T(scalar),
		z: q.z * T(scalar),
		w: q.w * T(scalar),
	}
}

// Neg implements Quaternion.Neg.
func (q quaternion[T]) Neg() Quaternion {
	return quaternion[T]{x: -q.x, y: -q.y, z: -q.z, w: -q.w}
}

// Length implements Quaternion.Length.
func (q quaternion[T]) Length() Scalar {
	return Scalar(math.Sqrt(float64(q.x*q.x + q.y*q.y + q.z*q.z + q.w*q.w)))
}

// Dot implements Quaternion.Dot.
func (q quaternion[T]) Dot(other Quaternion) Scalar {
	return Scalar(q.x*T(other.X()) + q.y*T(other.Y()) + q.z*T(other.Z()) + q.w*T(other.W()))
}

// Normalize implements Quaternion.Normalize.
func (q quaternion[T]) Normalize() Quaternion {
	len := q.Length()
	if len == 0 {
		return quaternion[T]{}
	}
	return quaternion[T]{
		x: q.x / T(len),
		y: q.y / T(len),
		z: q.z / T(len),
		w: q.w / T(len),
	}
}

// Invert implements Quaternion.Invert.
func (q quaternion[T]) Invert() Quaternion {
	normSq := q.x*q.x + q.y*q.y + q.z*q.z + q.w*q.w
	var zero T
	if normSq == zero {
		return quaternion[T]{}
	}
	return quaternion[T]{
		x: -q.x / normSq,
		y: -q.y / normSq,
		z: -q.z / normSq,
		w: q.w / normSq,
	}
}

// Slerp implements Quaternion.Slerp.
func (q quaternion[T]) Slerp(other Quaternion, t float64) Quaternion {
	cosine := ToFloat64(q.Dot(other))
	if NearlyEqual(T(cosine), 1.0) {
		sine := math.Sqrt(1.0 - cosine*cosine)
		angle := math.Atan2(sine, cosine)
		sineInverse := 1.0 / sine
		c0 := T(math.Sin((1.0-t)*angle) * sineInverse)
		c1 := T(math.Sin(t*angle) * sineInverse)
		return q.Scale(Scalar(c0)).Add(other.Scale(Scalar(c1))).Normalize()
	} else {
		return q.Scale(Scalar(1.0 - t)).Add(other.Scale(Scalar(t))).Normalize()
	}
}

// RotateVector3 implements Quaternion.RotateVector3.
func (q quaternion[T]) RotateVector3(vector Vector3) Vector3 {
	vx, vy, vz := T(vector.X()), T(vector.Y()), T(vector.Z())
	qx, qy, qz, qw := q.x, q.y, q.z, q.w
	tx := T(2) * (qy*vz - qz*vy)
	ty := T(2) * (qz*vx - qx*vz)
	tz := T(2) * (qx*vy - qy*vx)
	rx := vx + qw*tx + (qy*tz - qz*ty)
	ry := vy + qw*ty + (qz*tx - qx*tz)
	rz := vz + qw*tz + (qx*ty - qy*tx)
	return NewVector3(rx, ry, rz)
}

// String implements Quaternion.String.
func (q quaternion[T]) String() string {
	return fmt.Sprintf("(%v, %v, %v, %v)", q.x, q.y, q.z, q.w)
}
