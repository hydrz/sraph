package geom

import (
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
