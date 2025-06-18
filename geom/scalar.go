package geom

import (
	"cmp"
	"math"
	"strconv"
)

// Scalar is a generic interface for numeric types used in geometry calculations.
type Scalar interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
	Float64() float64
	String() string
}

// Eq compares two scalar values of type T with a tolerance.
//
// 0.001 for Float32 and 0.0000001 for Float64.
func Eq[T Scalar](a, b T) bool {
	if a == b {
		return true
	}
	tolerance := Epsilon32
	if _, ok := any(a).(F64); ok {
		tolerance = Epsilon64
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff.Float64() < tolerance
}

// IsFinite checks if the scalar value is finite.
func IsFinite[T Scalar](s T) bool {
	if _, ok := any(s).(F64); ok {
		return !math.IsNaN(s.Float64()) && !math.IsInf(s.Float64(), 0)
	}
	if _, ok := any(s).(F32); ok {
		return !math.IsNaN(s.Float64()) && !math.IsInf(s.Float64(), 0)
	}
	return true // For integer types, we assume they are finite
}

func Abs[T Scalar](s T) T {
	if s < 0 {
		return -s
	}
	return s
}

// Clamp clamps the scalar value between min and max.
func Clamp[T cmp.Ordered](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// F32 is a 32-bit floating point scalar type.
type F32 float32

// Float64 implements Scalar.
func (f F32) Float64() float64 {
	return float64(f)
}

// String implements Scalar.
func (f F32) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

// Value returns the underlying float32 value.
func (f F32) Value() float32 {
	return float32(f)
}

// F64 is a 64-bit floating point scalar type.
type F64 float64

// Float64 implements Scalar.
func (f F64) Float64() float64 {
	return float64(f)
}

// String implements Scalar.
func (f F64) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

// I32 is a 32-bit integer scalar type.
type I32 int32

// Float64 implements Scalar.
func (i I32) Float64() float64 {
	return float64(i)
}

// String implements Scalar.
func (i I32) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// I64 is a 64-bit integer scalar type.
type I64 int64

// Float64 implements Scalar.
func (i I64) Float64() float64 {
	return float64(i)
}

// String implements Scalar.
func (i I64) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// I26_6 is a signed 26.6 fixed-point number.
// The integer part ranges from -33554432 to 33554431.
// The format is: [integer(26)][fraction(6)].
// For example, the number one-and-a-quarter is I26_6(1<<6 + 1<<4).
type I26_6 int32

// Float64 implements Scalar.
func (f I26_6) Float64() float64 {
	return float64(f) / (1 << 6) // Divide by 2^6 to convert to float64
}

// String implements Scalar.
func (f I26_6) String() string {
	return strconv.FormatFloat(f.Float64(), 'f', -1, 64)
}

// Radians represents an angle in radians.
// It is a generic struct for type safety and clarity in geometry calculations.
type Radians[T Scalar] struct {
	radians T
}

func NewRadians[T Scalar](radians T) Radians[T] {
	return Radians[T]{radians: radians}
}

// Degrees converts radians to degrees.
func (r Radians[T]) Degrees() Degrees[T] {
	return Degrees[T]{degrees: T(r.Float64() * 180 / math.Pi)}
}

// Float64 returns the float64 value of the radians.
func (r Radians[T]) Float64() float64 {
	return r.radians.Float64()
}

// String returns a string representation of the radians value.
func (r Radians[T]) String() string {
	return strconv.FormatFloat(r.radians.Float64(), 'f', -1, 64) + "rad"
}

// Degrees represents an angle in degrees.
type Degrees[T Scalar] struct {
	degrees T
}

func NewDegrees[T Scalar](degrees T) Degrees[T] {
	return Degrees[T]{degrees: degrees}
}

// Radians converts degrees to radians.
func (d Degrees[T]) Radians() Radians[T] {
	return Radians[T]{radians: T(d.Float64() * math.Pi / 180)}
}

// Float64 returns the float64 value of the degrees.
func (d Degrees[T]) Float64() float64 {
	return d.degrees.Float64()
}

// String returns a string representation of the degrees value.
func (d Degrees[T]) String() string {
	return strconv.FormatFloat(d.degrees.Float64(), 'f', -1, 64) + "°"
}
