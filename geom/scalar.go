package geom

import (
	"fmt"
	"math"
	"strconv"
)

// Scalar is a generic interface for numeric types used in geometry calculations.
type Scalar interface {
	// no unsigned integers because they are not compatible with negative values
	// which are common in geometry calculations.
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
	Float64() float64
	String() string
}

// ScalarEq compares two scalar values of type T with a tolerance.
//
// 0.001 for Float32 and 0.0000001 for Float64.
func ScalarEq[T Scalar](a, b T) bool {
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

// Int is a go built-in integer type.
type Int int

// Float64 implements Scalar.
func (i Int) Float64() float64 {
	return float64(i)
}

// String implements Scalar.
func (i Int) String() string {
	return strconv.Itoa(int(i))
}

// I26_6 is a signed 26.6 fixed-point number.
// The integer part ranges from -33554432 to 33554431.
// The format is: [integer(26)][fraction(6)].
// For example, the number one-and-a-quarter is I26_6(1<<6 + 1<<4).
type I26_6 int32

// Float64 implements Scalar.
func (x I26_6) Float64() float64 {
	return float64(x) / (1 << 6) // Divide by 2^6 to convert to float64
}

// String implements Scalar.
func (x I26_6) String() string {
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

// Radians represents an angle in radians.
// It is a generic struct for type safety and clarity in geometry calculations.
type Radians F32

// Degrees converts radians to degrees.
func (r Radians) Degrees() Degrees {
	return Degrees(r * 180 / math.Pi)
}

// Float64 implements Scalar.
func (r Radians) Float64() float64 {
	return F32(r).Float64()
}

// String implements Scalar.
func (r Radians) String() string {
	return F32(r).String() + " rad"
}

// Degrees represents an angle in degrees.
type Degrees F32

// Radians converts degrees to radians.
func (d Degrees) Radians() Radians {
	return Radians(d * math.Pi / 180)
}

// Float64 implements Scalar.
func (d Degrees) Float64() float64 {
	return F32(d).Float64()
}

// String implements Scalar.
func (d Degrees) String() string {
	return F32(d).String() + "°"
}
