package geom

import (
	"cmp"
	"fmt"
	"math"
	"strconv"
)

const (
	// Epsilon32 is the tolerance value for float32 comparisons (0.001).
	Epsilon32 = 1e-3
	// Epsilon64 is the tolerance value for float64 comparisons (0.000001).
	Epsilon64 = 1e-6

	// Sqrt2Over2 represents sqrt(2)/2, commonly used in geometry calculations.
	Sqrt2Over2 = math.Sqrt2 / 2
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Floater represents types that can be converted to float64 for geometric calculations.
// This interface allows extending the geometry package to support custom numeric types
// such as fixed-point integers or other specialized number representations.
type Floater interface {
	Float64() float64
}

type Scalar float32

type Radians Scalar

func (r Radians) Degrees() Degrees {
	return Degrees(r) * (180.0 / math.Pi)
}

type Degrees Scalar

func (d Degrees) Radians() Radians {
	return Radians(d) * (math.Pi / 180.0)
}

// NearlyEqual reports whether two scalar values are approximately equal within a tolerance.
// Uses Epsilon32 (0.001) for float32 types and Epsilon64 (0.000001) for float64 types.
// This function handles floating-point precision issues in geometric calculations.
func NearlyEqual[A Number, B Number](a A, b B) bool {
	bb := A(b) // Convert b to the same type as a for comparison
	if a == bb {
		return true
	}
	var tolerance = Epsilon32
	if _, ok := any(a).(float64); ok {
		tolerance = Epsilon64
	}
	return ToFloat64(Abs(a-bb)) < tolerance
}

// Abs returns the absolute value of a scalar number.
func Abs[T Number](s T) T {
	if s < 0 {
		return -s
	}
	return s
}

// IsFinite reports whether the scalar value is finite (not NaN or infinite).
// This is useful for validating geometric calculations and preventing errors.
func IsFinite[T Number](s T) bool {
	return !(math.IsNaN(ToFloat64(s)) || math.IsInf(ToFloat64(s), 0))
}

// Clamp restricts a value to lie within the specified minimum and maximum bounds.
// Returns min if value < min, max if value > max, otherwise returns value unchanged.
func Clamp[T cmp.Ordered](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Cond returns ifTrue if condition is true, otherwise returns ifFalse.
// This function provides a concise way to select between two values based on a condition.
func Cond[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

// ToFloat64 converts a scalar value to float64.
// If the value implements Floater interface, uses its Float64 method for conversion.
// Otherwise, performs a direct type conversion to float64.
func ToFloat64[T Number](s T) float64 {
	if f, ok := any(s).(Floater); ok {
		return f.Float64()
	}
	return float64(s)
}

// ToString converts a scalar value to its string representation.
// Uses the String method if the value implements fmt.Stringer, otherwise formats as float64.
func ToString[T Number](s T) string {
	if f, ok := any(s).(fmt.Stringer); ok {
		return f.String()
	}
	return strconv.FormatFloat(ToFloat64(s), 'f', -1, 64)
}

// Int26_6 represents a signed 26.6 fixed-point number.
// The integer part uses 26 bits (range: -33554432 to 33554431) and the fractional part uses 6 bits.
// For example, 1.25 is represented as Int26_6(1<<6 + 1<<4) = Int26_6(80).
// This type is commonly used in graphics and typography for sub-pixel precision.
type Int26_6 int32

// Float64 converts the fixed-point number to a float64 value.
func (x Int26_6) Float64() float64 {
	return float64(x) / (1 << 6) // Divide by 2^6 to convert to float64
}

// String returns a string representation in "integer:fraction" format.
// The fraction part is displayed as a 2-digit value from 00 to 63.
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
