package geom

import (
	"math"
	"strconv"
)

// Scalar is a generic interface for numeric types used in geometry calculations.
// It supports conversion to float64 and string representation.
type Scalar interface {
	~int32 | ~float32 | ~float64
	ToFloat64() float64
	String() string
}

// NewScalar creates a new scalar value of type T.
func NewScalar[T Scalar](value T) T {
	return value
}

// MaxScalar returns the maximum value for the given scalar type T.
func MaxScalar[T Scalar]() T {
	var zero T
	switch any(zero).(type) {
	case I32:
		return any(I32(math.MaxInt32)).(T)
	case F32:
		return any(F32(math.MaxFloat32)).(T)
	case F64:
		return any(F64(math.MaxFloat64)).(T)
	case FixedI26_6:
		return any(FixedI26_6(math.MaxInt32)).(T)
	default:
		// fallback for unknown types
		return zero
	}
}

// ScalarEqualF32 can avoid type assertion for F32 comparisons.
func ScalarEqualF32(a, b F32) bool {
	diff := a.ToFloat64() - b.ToFloat64()
	return math.Abs(diff) < Epsilon32
}

// ScalarEqual compares two scalar values of type T with a tolerance.
//
// 0.001 for F32 and 0.0000001 for F64.
func ScalarEqual[T Scalar](a, b T) bool {
	tolerance := Epsilon32

	if _, ok := any(a).(F64); ok {
		tolerance = Epsilon64
	} else if _, ok := any(b).(F64); ok {
		tolerance = Epsilon64
	}

	diff := a.ToFloat64() - b.ToFloat64()
	return math.Abs(diff) < tolerance
}

// ScalarIsFinite checks if the scalar value is finite.
func ScalarIsFinite[T Scalar](s T) bool {
	if _, ok := any(s).(F64); ok {
		return !math.IsNaN(s.ToFloat64()) && !math.IsInf(s.ToFloat64(), 0)
	}
	if _, ok := any(s).(F32); ok {
		return !math.IsNaN(float64(s.ToFloat64())) && !math.IsInf(float64(s.ToFloat64()), 0)
	}
	return true // For integer types, we assume they are finite
}

// I32 is a 32-bit integer scalar type.
type I32 int32

// ToFloat64 implements Scalar.
func (i I32) ToFloat64() float64 {
	return float64(i)
}

// String implements Scalar.
func (i I32) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// F32 is a 32-bit floating point scalar type.
type F32 float32

// ToFloat64 implements Scalar.
func (f F32) ToFloat64() float64 {
	return float64(f)
}

// String implements Scalar.
func (f F32) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

// Unwarp returns the underlying float32 value.
func (f F32) Unwarp() any {
	return float32(f)
}

// F64 is a 64-bit floating point scalar type.
type F64 float64

// ToFloat64 implements Scalar.
func (f F64) ToFloat64() float64 {
	return float64(f)
}

// String implements Scalar.
func (f F64) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

// FixedI26_6 is a signed 26.6 fixed-point number.
// The integer part ranges from -33554432 to 33554431.
// The format is: [integer(26)][fraction(6)].
// For example, the number one-and-a-quarter is FixedI26_6(1<<6 + 1<<4).
type FixedI26_6 int32

// ToFloat64 implements Scalar.
func (f FixedI26_6) ToFloat64() float64 {
	return float64(f) / (1 << 6) // Divide by 2^6 to convert to float64
}

// String implements Scalar.
func (f FixedI26_6) String() string {
	return strconv.FormatFloat(f.ToFloat64(), 'f', -1, 64)
}
