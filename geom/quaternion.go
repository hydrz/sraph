package geom

import (
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
	Add(o Quaternion) Quaternion
	// Sub returns the component-wise difference of this quaternion and another.
	Sub(o Quaternion) Quaternion
	// Mul returns the Hamilton product of this quaternion and another.
	// This corresponds to the composition of two 3D rotations.
	Mul(o Quaternion) Quaternion
	// Div returns the quotient of this quaternion divided by another.
	// This is equivalent to multiplying by the inverse of the other quaternion.
	Div(o Quaternion) Quaternion

	// Equal reports whether this quaternion and another are equal within floating-point tolerance.
	Equal(o Quaternion) bool
	// Scale returns this quaternion scaled by the given scalar value.
	Scale(scalar Scalar) Quaternion
	// Neg returns the negation of this quaternion (all components negated).
	Neg() Quaternion

	// Length returns the magnitude (norm) of the quaternion.
	// For unit quaternions representing rotations, this should be 1.
	Length() Scalar
	// Dot returns the dot product of this quaternion and another.
	Dot(o Quaternion) Scalar
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
func NewQuaternion[T Number](x, y, z, w T) Quaternion {
	return quaternion[T]{x, y, z, w}
}

// NewQuaternionFromAxisAngle creates a new quaternion from an axis and angle.
func NewQuaternionFromAxisAngle(axis Vector3, angle Radians) Quaternion {
	axis = axis.Normalize()
	halfAngle := ToFloat64(angle) / 2
	sinHalfAngle := Scalar(math.Sin(halfAngle))
	cosHalfAngle := Scalar(math.Cos(halfAngle))
	return NewQuaternion(
		axis.X()*sinHalfAngle,
		axis.Y()*sinHalfAngle,
		axis.Z()*sinHalfAngle,
		cosHalfAngle,
	)
}

// X implements Quaternion.X.
func (q quaternion[T]) X() Scalar { return ToScalar(q.x) }

// Y implements Quaternion.Y.
func (q quaternion[T]) Y() Scalar { return ToScalar(q.y) }

// Z implements Quaternion.Z.
func (q quaternion[T]) Z() Scalar { return ToScalar(q.z) }

// W implements Quaternion.W.
func (q quaternion[T]) W() Scalar { return ToScalar(q.w) }

// Add implements Quaternion.Add.
func (q quaternion[T]) Add(o Quaternion) Quaternion {
	return NewQuaternion(
		q.X()+o.X(),
		q.Y()+o.Y(),
		q.Z()+o.Z(),
		q.W()+o.W(),
	)
}

// Sub implements Quaternion.Sub.
func (q quaternion[T]) Sub(other Quaternion) Quaternion {
	return NewQuaternion(
		q.X()-other.X(),
		q.Y()-other.Y(),
		q.Z()-other.Z(),
		q.W()-other.W(),
	)
}

// Mul implements Quaternion.Mul (Hamilton product).
func (q quaternion[T]) Mul(o Quaternion) Quaternion {
	x1, y1, z1, w1 := q.X(), q.Y(), q.Z(), q.W()
	x2, y2, z2, w2 := o.X(), o.Y(), o.Z(), o.W()
	return NewQuaternion(
		w1*x2+x1*w2+y1*z2-z1*y2,
		w1*y2-x1*z2+y1*w2+z1*x2,
		w1*z2+x1*y2-y1*x2+z1*w2,
		w1*w2-x1*x2-y1*y2-z1*z2,
	)
}

// Div implements Quaternion.Div.
func (q quaternion[T]) Div(o Quaternion) Quaternion {
	return q.Mul(o.Invert())
}

// Equal implements Quaternion.Equal.
func (q quaternion[T]) Equal(other Quaternion) bool {
	return NearlyEqual(q.X(), other.X()) &&
		NearlyEqual(q.Y(), other.Y()) &&
		NearlyEqual(q.Z(), other.Z()) &&
		NearlyEqual(q.W(), other.W())
}

// Scale implements Quaternion.Scale.
func (q quaternion[T]) Scale(scalar Scalar) Quaternion {
	return NewQuaternion(
		q.X()*scalar,
		q.Y()*scalar,
		q.Z()*scalar,
		q.W()*scalar,
	)
}

// Neg implements Quaternion.Neg.
func (q quaternion[T]) Neg() Quaternion {
	return NewQuaternion(
		-q.X(),
		-q.Y(),
		-q.Z(),
		-q.W(),
	)
}

// Length implements Quaternion.Length.
func (q quaternion[T]) Length() Scalar {
	return ToScalar(math.Sqrt(ToFloat64(q.x*q.x + q.y*q.y + q.z*q.z + q.w*q.w)))
}

// Dot implements Quaternion.Dot.
func (q quaternion[T]) Dot(o Quaternion) Scalar {
	return q.X()*o.X() +
		q.Y()*o.Y() +
		q.Z()*o.Z() +
		q.W()*o.W()
}

// Normalize implements Quaternion.Normalize.
func (q quaternion[T]) Normalize() Quaternion {
	len := q.Length()
	if len == 0 {
		return quaternion[T]{}
	}
	return NewQuaternion(
		q.X()/len,
		q.Y()/len,
		q.Z()/len,
		q.W()/len,
	)
}

// Invert implements Quaternion.Invert.
func (q quaternion[T]) Invert() Quaternion {
	normSq := q.X()*q.X() + q.Y()*q.Y() + q.Z()*q.Z() + q.W()*q.W()
	if normSq == 0 {
		return quaternion[T]{}
	}
	return NewQuaternion(
		-q.X()/normSq,
		-q.Y()/normSq,
		-q.Z()/normSq,
		q.W()/normSq,
	)
}

// Slerp implements Quaternion.Slerp.
func (q quaternion[T]) Slerp(o Quaternion, t float64) Quaternion {
	cosine := ToFloat64(q.Dot(o))
	if !NearlyEqual(cosine, 1.0) {
		return q.Scale(Scalar(1.0 - t)).Add(o.Scale(Scalar(t))).Normalize()
	}

	sine := math.Sqrt(1.0 - cosine*cosine)
	angle := math.Atan2(sine, cosine)
	sineInverse := 1.0 / sine
	c0 := ToScalar(math.Sin((1.0-t)*angle) * sineInverse)
	c1 := ToScalar(math.Sin(t*angle) * sineInverse)
	return q.Scale(c0).Add(o.Scale(c1)).Normalize()
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
	return "(" + ToString(q.x) + ", " +
		ToString(q.y) + ", " +
		ToString(q.z) + ", " +
		ToString(q.w) + ")"
}
