package geom

import (
	"cmp"
	"fmt"
	"math"
	"strconv"
)

const (
	// 0.001
	Epsilon32 = 1e-3
	// 0.000001
	Epsilon64 = 1e-6

	// sqrt(2) / 2 == 1/sqrt(2)
	Sqrt2Over2 = math.Sqrt2 / 2
)

// Number is a generic interface for numeric types used in geometry calculations.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

// Floater is an interface for types that can be converted to float64.
//
// package geom used Floater interface to convert different numeric types to float64.
// If you want to extend the geometry package to support new numeric types,
// you need to implement the Floater interface.
// For example, the [Int26_6] type, which is a fixed-point integer type.
type Floater interface {
	Float64() float64
}

// Int26_6 is a signed 26.6 fixed-point number.
// The integer part ranges from -33554432 to 33554431.
// The format is: [integer(26)][fraction(6)].
// For example, the number one-and-a-quarter is Int26_6(1<<6 + 1<<4).
type Int26_6 int32

// Float64 implements [Floater].
func (x Int26_6) Float64() float64 {
	return float64(x) / (1 << 6) // Divide by 2^6 to convert to float64
}

// String implements [fmt.Stringer].
func (x Int26_6) String() string {
	const shift, mask = 6, 1<<6 - 1
	if x >= 0 {
		return fmt.Sprintf("%d:%02d", int32(x>>shift), int32(x&mask))
	}
	x = -x
	if x >= 0 {
		return fmt.Sprintf("-%d:%02d", int32(x>>shift), int32(x&mask))
	}
	return "-33554432:00" // The minimum value is -(1<<25).
}

// NearlyEqual compares two number values of type T with a tolerance.
//
// 0.001 for Float32 and 0.0000001 for Float64.
func NearlyEqual[T Number](a, b T) bool {
	if a == b {
		return true
	}
	var tolerance = Epsilon32
	if _, ok := any(a).(float64); ok {
		tolerance = Epsilon64
	}
	return Abs(a-b) < T(tolerance)
}

// Abs returns the absolute value of the number.
func Abs[T Number](s T) T {
	if s < 0 {
		return -s
	}
	return s
}

// IsFinite checks if the number is finite.
func IsFinite[T Number](s T) bool {
	return !(math.IsNaN(ToFloat64(s)) || math.IsInf(ToFloat64(s), 0))
}

// Clamp clamps the value between min and max.
func Clamp[T cmp.Ordered](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ToFloat64 converts a number to float64.
func ToFloat64[T Number](s T) float64 {
	if f, ok := any(s).(Floater); ok {
		return f.Float64()
	}
	return float64(s)
}

// ToString converts a number to its string representation.
func ToString[T Number](s T) string {
	if f, ok := any(s).(fmt.Stringer); ok {
		return f.String()
	}
	return strconv.FormatFloat(ToFloat64(s), 'f', -1, 64)
}
