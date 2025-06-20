package geom

import (
	"fmt"
	"math"
	"strconv"
)

// Scalar is a generic interface for numeric types used in geometry calculations.
type Scalar = float32

// Radians represents an angle in radians.
// It is a generic struct for type safety and clarity in geometry calculations.
type Radians Scalar

// Degrees converts radians to degrees.
func (r Radians) Degrees() Degrees {
	return Degrees(r * 180 / math.Pi)
}

// String implements Scalar.
func (r Radians) String() string {
	return strconv.FormatFloat(float64(r), 'f', -1, 64) + " rad"
}

// Degrees represents an angle in degrees.
type Degrees Scalar

// Radians converts degrees to radians.
func (d Degrees) Radians() Radians {
	return Radians(d * math.Pi / 180)
}

// String implements [fmt.Stringer].
func (d Degrees) String() string {
	return strconv.FormatFloat(float64(d), 'f', -1, 64) + "°"
}

// I26_6 is a signed 26.6 fixed-point number.
// The integer part ranges from -33554432 to 33554431.
// The format is: [integer(26)][fraction(6)].
// For example, the number one-and-a-quarter is I26_6(1<<6 + 1<<4).
type I26_6 int32

// Float64 implements [Floater].
func (x I26_6) Float64() float64 {
	return float64(x) / (1 << 6) // Divide by 2^6 to convert to float64
}

// String implements [fmt.Stringer].
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
