package geom

import (
	"fmt"
	"math"
)

// Quaternion represents a quaternion for 3D rotations.
// Quaternions are used in 3D graphics and robotics to represent rotations
// and orientations because they avoid gimbal lock and provide smooth interpolation.
type Quaternion[T Number] struct {
	X, Y, Z, W T
}

// NewQuaternionFromAxisAngle creates a new quaternion from an axis and angle.
func NewQuaternionFromAxisAngle[T Number](axis Vector3[T], angle Radians) Quaternion[T] {
	axis = axis.Normalize()
	halfAngle := ToFloat64(angle) / 2
	sinHalfAngle := T(math.Sin(halfAngle))
	cosHalfAngle := T(math.Cos(halfAngle))
	return Quaternion[T]{
		X: axis.X * sinHalfAngle,
		Y: axis.Y * sinHalfAngle,
		Z: axis.Z * sinHalfAngle,
		W: cosHalfAngle,
	}
}

// Add adds another quaternion to this one and returns the result.
func (q Quaternion[T]) Add(other Quaternion[T]) Quaternion[T] {
	return Quaternion[T]{
		X: q.X + other.X,
		Y: q.Y + other.Y,
		Z: q.Z + other.Z,
		W: q.W + other.W,
	}
}

// Sub subtracts another quaternion from this one and returns the result.
func (q Quaternion[T]) Sub(other Quaternion[T]) Quaternion[T] {
	return Quaternion[T]{
		X: q.X - other.X,
		Y: q.Y - other.Y,
		Z: q.Z - other.Z,
		W: q.W - other.W,
	}
}

// Mul multiplies this quaternion by another and returns the result.
// Hamilton product.
func (q Quaternion[T]) Mul(other Quaternion[T]) Quaternion[T] {
	x1, y1, z1, w1 := q.X, q.Y, q.Z, q.W
	x2, y2, z2, w2 := other.X, other.Y, other.Z, other.W
	return Quaternion[T]{
		X: w1*x2 + x1*w2 + y1*z2 - z1*y2,
		Y: w1*y2 - x1*z2 + y1*w2 + z1*x2,
		Z: w1*z2 + x1*y2 - y1*x2 + z1*w2,
		W: w1*w2 - x1*x2 - y1*y2 - z1*z2,
	}
}

// Div divides this quaternion by another and returns the result.
// q / r = q * r^-1
func (q Quaternion[T]) Div(other Quaternion[T]) Quaternion[T] {
	return q.Mul(other.Invert())
}

// Eq checks if this quaternion is Eq to another.
func (q Quaternion[T]) Eq(other Quaternion[T]) bool {
	return NearlyEqual(q.X, other.X) &&
		NearlyEqual(q.Y, other.Y) &&
		NearlyEqual(q.Z, other.Z) &&
		NearlyEqual(q.W, other.W)
}

// Scale scales the quaternion by a scalar value and returns the result.
func (q Quaternion[T]) Scale(scalar T) Quaternion[T] {
	return Quaternion[T]{
		X: q.X * scalar,
		Y: q.Y * scalar,
		Z: q.Z * scalar,
		W: q.W * scalar,
	}
}

// Neg negates the quaternion (inverts the sign of all components).
// Negation is useful for reversing the direction of a rotation.
func (q Quaternion[T]) Neg() Quaternion[T] {
	return Quaternion[T]{
		X: -q.X,
		Y: -q.Y,
		Z: -q.Z,
		W: -q.W,
	}
}

// Length calculates the length (magnitude) of the quaternion.
func (q Quaternion[T]) Length() T {
	return T(math.Sqrt(ToFloat64(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)))
}

// Dot calculates the dot product of this quaternion with another.
// The dot product is used to measure the similarity or angle between two quaternions.
// It is commonly used in interpolation and blending of rotations.
func (q Quaternion[T]) Dot(other Quaternion[T]) T {
	return q.X*other.X + q.Y*other.Y + q.Z*other.Z + q.W*other.W
}

// Normalize normalizes the quaternion to unit length and returns the result.
// Normalization is important in many applications, such as 3D rotations,
// to ensure the quaternion represents a valid rotation (unit quaternion).
func (q Quaternion[T]) Normalize() Quaternion[T] {
	len := q.Length()
	var zero T
	if len == zero {
		return Quaternion[T]{}
	}
	return Quaternion[T]{
		X: q.X / len,
		Y: q.Y / len,
		Z: q.Z / len,
		W: q.W / len,
	}
}

// Invert returns the inverse of this quaternion (negates the vector part, keeps the scalar part).
// The inverse is used to undo a rotation or compute relative rotations.
// The inverse of a quaternion q = (x, y, z, w) is given by q^-1 = (-x, -y, -z, w) / (x^2 + y^2 + z^2 + w^2).
func (q Quaternion[T]) Invert() Quaternion[T] {
	normSq := q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W
	var zero T
	if normSq == zero {
		return Quaternion[T]{}
	}
	return Quaternion[T]{
		X: -q.X / normSq,
		Y: -q.Y / normSq,
		Z: -q.Z / normSq,
		W: q.W / normSq,
	}
}

// Slerp performs spherical linear interpolation (SLERP) between
// this quaternion and another quaternion by a factor of time (0.0 to 1.0).
// Formula: q' = q * sin((1-t) * θ) / sin(θ) + other * sin(t * θ) / sin(θ)
func (q Quaternion[T]) Slerp(other Quaternion[T], time float64) Quaternion[T] {
	time = Clamp(time, 0.0, 1.0) // Ensure time is between 0 and 1
	cosine := ToFloat64(q.Dot(other))
	if NearlyEqual(T(cosine), 1.0) {
		// Spherical Interpolation
		sine := math.Sqrt(1.0 - cosine*cosine)
		angle := math.Atan2(sine, cosine)
		sineInverse := 1.0 / sine
		c0 := T(math.Sin((1.0-time)*angle) * sineInverse)
		c1 := T(math.Sin(time*angle) * sineInverse)
		return q.Scale(c0).Add(other.Scale(c1)).Normalize()
	} else {
		// Linear Interpolation
		return q.Scale(T(1.0) - T(time)).Add(other.Scale(T(time))).Normalize()
	}
}

// RotateVector3 rotates a 3D vector by this quaternion and returns the result.
// This is useful for applying the quaternion as a rotation to a vector in 3D space,
// such as transforming directions or points in graphics and simulation.
// Formula: v' = v + 2 * w * (qv × v) + 2 * (qv × (qv × v))
// where qv is the vector part of the quaternion (x, y, z) and w is the scalar part.
func (q Quaternion[T]) RotateVector3(vector Vector3[T]) Vector3[T] {
	// v' = q * v * q^-1
	vx, vy, vz := vector.X, vector.Y, vector.Z
	qx, qy, qz, qw := q.X, q.Y, q.Z, q.W

	// Quaternion-vector multiplication (optimized)
	// t = 2 * cross(q.xyz, v)
	tx := T(2) * (qy*vz - qz*vy)
	ty := T(2) * (qz*vx - qx*vz)
	tz := T(2) * (qx*vy - qy*vx)

	// v' = v + qw * t + cross(q.xyz, t)
	rx := vx + qw*tx + (qy*tz - qz*ty)
	ry := vy + qw*ty + (qz*tx - qx*tz)
	rz := vz + qw*tz + (qx*ty - qy*tx)

	return Vector3[T]{rx, ry, rz}
}

// String returns a string representation. like "(x, y, z, w)".
func (q Quaternion[T]) String() string {
	return fmt.Sprintf("(%v, %v, %v, %v)", q.X, q.Y, q.Z, q.W)
}
