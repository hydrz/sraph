package geom

import (
	"math"
	"strconv"
)

// Scalar is a generic interface for numeric types used in geometry calculations.
// It supports conversion to float64 and string representation.
type Scalar interface {
	~int32 | ~float32 | ~float64
	Float64() float64
	String() string
}

// New creates a new scalar value of type T.
func New[T Scalar](value T) T {
	return value
}

// Max returns the maximum value for the given scalar type T.
func Max[T Scalar]() T {
	var zero T
	switch any(zero).(type) {
	case Int32:
		return any(Int32(math.MaxInt32)).(T)
	case Float32:
		return any(Float32(math.MaxFloat32)).(T)
	case Float64:
		return any(Float64(math.MaxFloat64)).(T)
	case Fixed26_6:
		return any(Fixed26_6(math.MaxInt32)).(T)
	default:
		// fallback for unknown types
		return zero
	}
}

// EqualFloat32 can avoid type assertion for Float32 comparisons.
func EqualFloat32(a, b Float32) bool {
	diff := a - b
	return math.Abs(float64(diff)) < Epsilon32
}

// EqualFloat64 can avoid type assertion for Float64 comparisons.
func EqualFloat64(a, b Float64) bool {
	diff := a - b
	return math.Abs(float64(diff)) < Epsilon32
}

// Equal compares two scalar values of type T with a tolerance.
//
// 0.001 for Float32 and 0.0000001 for Float64.
func Equal[T Scalar](a, b T) bool {
	tolerance := Epsilon32

	if _, ok := any(a).(Float64); ok {
		tolerance = Epsilon64
	} else if _, ok := any(b).(Float64); ok {
		tolerance = Epsilon64
	}

	diff := a.Float64() - b.Float64()
	return math.Abs(diff) < tolerance
}

// IsFinite checks if the scalar value is finite.
func IsFinite[T Scalar](s T) bool {
	if _, ok := any(s).(Float64); ok {
		return !math.IsNaN(s.Float64()) && !math.IsInf(s.Float64(), 0)
	}
	if _, ok := any(s).(Float32); ok {
		return !math.IsNaN(float64(s.Float64())) && !math.IsInf(float64(s.Float64()), 0)
	}
	return true // For integer types, we assume they are finite
}

// Clamp clamps the scalar value between min and max.
func Clamp[T Scalar](s, min, max T) T {
	if s.Float64() < min.Float64() {
		return min
	}
	if s.Float64() > max.Float64() {
		return max
	}
	return s
}

// Int32 is a 32-bit integer scalar type.
type Int32 int32

// Float64 implements Scalar.
func (i Int32) Float64() float64 {
	return float64(i)
}

// String implements Scalar.
func (i Int32) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// Float32 is a 32-bit floating point scalar type.
type Float32 float32

// Float64 implements Scalar.
func (f Float32) Float64() float64 {
	return float64(f)
}

// String implements Scalar.
func (f Float32) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

// Value returns the underlying float32 value.
func (f Float32) Value() float32 {
	return float32(f)
}

// Float64 is a 64-bit floating point scalar type.
type Float64 float64

// Float64 implements Scalar.
func (f Float64) Float64() float64 {
	return float64(f)
}

// String implements Scalar.
func (f Float64) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

// Fixed26_6 is a signed 26.6 fixed-point number.
// The integer part ranges from -33554432 to 33554431.
// The format is: [integer(26)][fraction(6)].
// For example, the number one-and-a-quarter is Fixed26_6(1<<6 + 1<<4).
type Fixed26_6 int32

// Float64 implements Scalar.
func (f Fixed26_6) Float64() float64 {
	return float64(f) / (1 << 6) // Divide by 2^6 to convert to float64
}

// String implements Scalar.
func (f Fixed26_6) String() string {
	return strconv.FormatFloat(f.Float64(), 'f', -1, 64)
}

// Radians represents an angle in radians.
// It is a float64 type, but defined separately to clarify its purpose in geometry calculations.
// Radians is used to represent angles in radians, which is common in geometry and trigonometry.
// The conversion between radians and degrees can be done using the formulas:
// Degrees = Radians * (180 / π)
type Radians float64

func (r Radians) Degrees() Degrees {
	return Degrees(r * (180 / math.Pi))
}

func (r Radians) String() string {
	return strconv.FormatFloat(float64(r), 'f', -1, 64) + "rad"
}

// Degrees represents an angle in degrees.
// It is a float64 type, but defined separately to clarify its purpose in geometry calculations.
// Degrees is used to represent angles in degrees, which is often more intuitive for users.
// The conversion between radians and degrees can be done using the formulas:
// Radians = Degrees * (π / 180)
type Degrees float64

func (d Degrees) Radians() Radians {
	return Radians(d * (math.Pi / 180))
}

func (d Degrees) String() string {
	return strconv.FormatFloat(float64(d), 'f', -1, 64) + "°"
}
