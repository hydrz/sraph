package geom

import (
	"fmt"
	"math"
	"time"
)

// Quaternion represents a quaternion for 3D rotations.
// Quaternions are used in 3D graphics and robotics to represent rotations
// and orientations because they avoid gimbal lock and provide smooth interpolation.
type Quaternion[T Scalar] interface {
	// X returns the x component of the quaternion.
	X() T
	// Y returns the y component of the quaternion.
	Y() T
	// Z returns the z component of the quaternion.
	Z() T
	// W returns the w component of the quaternion.
	W() T

	// Add adds another quaternion to this one and returns the result.
	Add(other Quaternion[T]) Quaternion[T]
	// Sub subtracts another quaternion from this one and returns the result.
	Sub(other Quaternion[T]) Quaternion[T]
	// Mul multiplies this quaternion by another and returns the result.
	Mul(other Quaternion[T]) Quaternion[T]
	// Div divides this quaternion by another and returns the result.
	Div(other Quaternion[T]) Quaternion[T]

	// Equal checks if this quaternion is equal to another.
	Equal(other Quaternion[T]) bool

	// Scale scales the quaternion by a scalar value and returns the result.
	Scale(scalar T) Quaternion[T]

	// Neg negates the quaternion (inverts the sign of all components).
	// Negation is useful for reversing the direction of a rotation.
	Neg() Quaternion[T]

	// Length calculates the length (magnitude) of the quaternion.
	Length() T

	// Dot calculates the dot product of this quaternion with another.
	// The dot product is used to measure the similarity or angle between two quaternions.
	// It is commonly used in interpolation and blending of rotations.
	Dot(other Quaternion[T]) T

	// Normalize normalizes the quaternion to unit length and returns the result.
	// Normalization is important in many applications, such as 3D rotations,
	// to ensure the quaternion represents a valid rotation (unit quaternion).
	Normalize() Quaternion[T]

	// Invert returns the inverse of this quaternion (negates the vector part, keeps the scalar part).
	// The inverse is used to undo a rotation or compute relative rotations.
	// The inverse of a quaternion q = (x, y, z, w) is given by q^-1 = (-x, -y, -z, w) / (x^2 + y^2 + z^2 + w^2).
	Invert() Quaternion[T]

	// Slerp performs spherical linear interpolation (SLERP) between
	// this quaternion and another quaternion by a factor of time (0.0 to 1.0).
	// SLERP is used for smooth interpolation of rotations, commonly in animation and orientation blending.
	// Formula: slerp = sin((1-t)*θ)/sin(θ) * q1 + sin(t*θ)/sin(θ) * q2
	Slerp(other Quaternion[T], time time.Time) Quaternion[T]

	// RotateVector3 rotates a 3D vector by this quaternion and returns the result.
	// This is useful for applying the quaternion as a rotation to a vector in 3D space,
	// such as transforming directions or points in graphics and simulation.
	// Formula: v' = v + 2 * w * (qv × v) + 2 * (qv × (qv × v))
	// where qv is the vector part of the quaternion (x, y, z) and w is the scalar part.
	RotateVector3(vector Vector3[T]) Vector3[T]

	// String returns a string representation. like "(x, y, z, w)".
	String() string
}

// NewQuaternion creates a new quaternion with the given components.
func NewQuaternion[T Scalar](x, y, z, w T) Quaternion[T] {
	return quaternion[T]{x: x, y: y, z: z, w: w}
}

// NewQuaternionFromAxisAngle creates a new quaternion from an axis and angle.
func NewQuaternionFromAxisAngle[T Scalar](axis Vector3[T], angle Radians) Quaternion[T] {
	// Normalize the axis vector
	axis = axis.Normalize()
	sinHalfAngle := T(math.Sin(float64(angle / 2)))
	cosHalfAngle := T(math.Cos(float64(angle / 2)))

	return &quaternion[T]{
		x: axis.X() * sinHalfAngle,
		y: axis.Y() * sinHalfAngle,
		z: axis.Z() * sinHalfAngle,
		w: cosHalfAngle,
	}
}

type quaternion[T Scalar] struct {
	x, y, z, w T
}

// X implements Quaternion.
func (q quaternion[T]) X() T {
	return q.x
}

// Y implements Quaternion.
func (q quaternion[T]) Y() T {
	return q.y
}

// Z implements Quaternion.
func (q quaternion[T]) Z() T {
	return q.z
}

// W implements Quaternion.
func (q quaternion[T]) W() T {
	return q.w
}

// Add implements Quaternion.
func (q quaternion[T]) Add(other Quaternion[T]) Quaternion[T] {
	return &quaternion[T]{
		x: q.x + other.X(),
		y: q.y + other.Y(),
		z: q.z + other.Z(),
		w: q.w + other.W(),
	}
}

// Sub implements Quaternion.
func (q quaternion[T]) Sub(other Quaternion[T]) Quaternion[T] {
	return &quaternion[T]{
		x: q.x - other.X(),
		y: q.y - other.Y(),
		z: q.z - other.Z(),
		w: q.w - other.W(),
	}
}

// Mul implements Quaternion.
func (q quaternion[T]) Mul(other Quaternion[T]) Quaternion[T] {
	// Hamilton product
	x1, y1, z1, w1 := q.x, q.y, q.z, q.w
	x2, y2, z2, w2 := other.X(), other.Y(), other.Z(), other.W()
	return &quaternion[T]{
		x: w1*x2 + x1*w2 + y1*z2 - z1*y2,
		y: w1*y2 - x1*z2 + y1*w2 + z1*x2,
		z: w1*z2 + x1*y2 - y1*x2 + z1*w2,
		w: w1*w2 - x1*x2 - y1*y2 - z1*z2,
	}
}

// Div implements Quaternion.
func (q quaternion[T]) Div(other Quaternion[T]) Quaternion[T] {
	// q / r = q * r^-1
	return q.Mul(other.Invert())
}

// Equal implements Quaternion.
func (q quaternion[T]) Equal(other Quaternion[T]) bool {
	return Equal(q.x, other.X()) &&
		Equal(q.y, other.Y()) &&
		Equal(q.z, other.Z()) &&
		Equal(q.w, other.W())
}

// Scale implements Quaternion.
func (q quaternion[T]) Scale(scalar T) Quaternion[T] {
	return &quaternion[T]{
		x: q.x * scalar,
		y: q.y * scalar,
		z: q.z * scalar,
		w: q.w * scalar,
	}
}

// Neg implements Quaternion.
func (q quaternion[T]) Neg() Quaternion[T] {
	return &quaternion[T]{
		x: -q.x,
		y: -q.y,
		z: -q.z,
		w: -q.w,
	}
}

// Length implements Quaternion.
func (q quaternion[T]) Length() T {
	return T(math.Sqrt(float64(q.x*q.x + q.y*q.y + q.z*q.z + q.w*q.w)))
}

// Dot implements Quaternion.
func (q quaternion[T]) Dot(other Quaternion[T]) T {
	return q.x*other.X() + q.y*other.Y() + q.z*other.Z() + q.w*other.W()
}

// Normalize implements Quaternion.
func (q quaternion[T]) Normalize() Quaternion[T] {
	len := q.Length()
	var zero T
	if len == zero {
		return &quaternion[T]{x: zero, y: zero, z: zero, w: zero}
	}
	return &quaternion[T]{
		x: q.x / len,
		y: q.y / len,
		z: q.z / len,
		w: q.w / len,
	}
}

// Invert implements Quaternion.
func (q quaternion[T]) Invert() Quaternion[T] {
	normSq := q.x*q.x + q.y*q.y + q.z*q.z + q.w*q.w
	var zero T
	if normSq == zero {
		return &quaternion[T]{x: zero, y: zero, z: zero, w: zero}
	}
	return &quaternion[T]{
		x: -q.x / normSq,
		y: -q.y / normSq,
		z: -q.z / normSq,
		w: q.w / normSq,
	}
}

// Slerp implements Quaternion.
func (q quaternion[T]) Slerp(other Quaternion[T], t time.Time) Quaternion[T] {
	// slerp = sin((1-t)*θ)/sin(θ) * q1 + sin(t*θ)/sin(θ) * q2
	frac := T(float64(t.UnixNano()%1e9) / 1e9)
	dot := q.Dot(other)
	// Clamp dot to [-1, 1]
	dot = Clamp(dot, T(-1), T(1))
	one := T(1)
	zero := T(0)
	if dot.Float64() > 1-Epsilon64 {
		// Linear interpolation for very close quaternions
		return (&quaternion[T]{
			x: q.x + (other.X()-q.x)*frac,
			y: q.y + (other.Y()-q.y)*frac,
			z: q.z + (other.Z()-q.z)*frac,
			w: q.w + (other.W()-q.w)*frac,
		}).Normalize()
	}
	theta := T(math.Acos(dot.Float64()))
	sinTheta := T(math.Sin(theta.Float64()))
	if Equal(sinTheta, zero) {
		return q
	}

	a := T(math.Sin(((one - frac) * theta).Float64())) / sinTheta
	b := T(math.Sin((frac * theta).Float64())) / sinTheta
	return &quaternion[T]{
		x: q.x*a + other.X()*b,
		y: q.y*a + other.Y()*b,
		z: q.z*a + other.Z()*b,
		w: q.w*a + other.W()*b,
	}
}

// RotateVector3 implements Quaternion.
func (q quaternion[T]) RotateVector3(vector Vector3[T]) Vector3[T] {
	// v' = q * v * q^-1
	vx, vy, vz := vector.X(), vector.Y(), vector.Z()
	qx, qy, qz, qw := q.x, q.y, q.z, q.w

	// Quaternion-vector multiplication (optimized)
	// t = 2 * cross(q.xyz, v)
	tx := T(2) * (qy*vz - qz*vy)
	ty := T(2) * (qz*vx - qx*vz)
	tz := T(2) * (qx*vy - qy*vx)

	// v' = v + qw * t + cross(q.xyz, t)
	rx := vx + qw*tx + (qy*tz - qz*ty)
	ry := vy + qw*ty + (qz*tx - qx*tz)
	rz := vz + qw*tz + (qx*ty - qy*tx)

	return NewVector3(rx, ry, rz)
}

// String implements Quaternion.
func (q quaternion[T]) String() string {
	return fmt.Sprintf("(%v, %v, %v, %v)", q.x, q.y, q.z, q.w)
}
