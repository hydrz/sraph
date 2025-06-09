package geom

import (
	"fmt"
	"math"
)

var _ Scalar = ScalarInt26_6(0)

// ScalarInt26_6 is a signed 26.6 fixed-point number.
//
// The integer part ranges from -33554432 to 33554431, inclusive. The
// fractional part has 6 bits of precision.
//
// For example, the number one-and-a-quarter is ScalarInt26_6(1<<6 + 1<<4).
type ScalarInt26_6 int32

// Add returns the sum of the current scalar and another scalar.
func (s ScalarInt26_6) Add(o Scalar) Scalar {
	return s + o.(ScalarInt26_6)
}

// Sub returns the difference between the current scalar and another scalar.
func (s ScalarInt26_6) Sub(o Scalar) Scalar {
	return s - o.(ScalarInt26_6)
}

// Mul returns the product of the current scalar and another scalar.
func (s ScalarInt26_6) Mul(o Scalar) Scalar {
	return ScalarInt26_6((int64(s) * int64(o.(ScalarInt26_6))) >> 6)
}

// Div returns the quotient of the current scalar divided by another scalar.
func (s ScalarInt26_6) Div(o Scalar) Scalar {
	return ScalarInt26_6((int64(s) << 6) / int64(o.(ScalarInt26_6)))
}

// Neg returns the negated value of the scalar.
func (s ScalarInt26_6) Neg() Scalar {
	return -s
}

// Equal returns true if the current scalar is equal to another scalar.
func (s ScalarInt26_6) Equal(o Scalar) bool {
	return s == o.(ScalarInt26_6)
}

// Less returns true if the current scalar is less than another scalar.
func (s ScalarInt26_6) Less(o Scalar) bool {
	return s < o.(ScalarInt26_6)
}

// LessEqual returns true if the current scalar is less than or equal to another scalar.
func (s ScalarInt26_6) LessEqual(o Scalar) bool {
	return s <= o.(ScalarInt26_6)
}

// Greater returns true if the current scalar is greater than another scalar.
func (s ScalarInt26_6) Greater(o Scalar) bool {
	return s > o.(ScalarInt26_6)
}

// GreaterEqual returns true if the current scalar is greater than or equal to another scalar.
func (s ScalarInt26_6) GreaterEqual(o Scalar) bool {
	return s >= o.(ScalarInt26_6)
}

// NotEqual returns true if the current scalar is not equal to another scalar.
func (s ScalarInt26_6) NotEqual(o Scalar) bool {
	return s != o.(ScalarInt26_6)
}

// Min returns the minimum value among the current scalar and the provided scalars.
func (s ScalarInt26_6) Min(o ...Scalar) Scalar {
	min := s
	for _, v := range o {
		f := v.(ScalarInt26_6)
		if f < min {
			min = f
		}
	}
	return min
}

// Max returns the maximum value among the current scalar and the provided scalars.
func (s ScalarInt26_6) Max(o ...Scalar) Scalar {
	max := s
	for _, v := range o {
		f := v.(ScalarInt26_6)
		if f > max {
			max = f
		}
	}
	return max
}

// IsNan returns true if the scalar is NaN (not a number).
func (s ScalarInt26_6) IsNan() bool {
	return false
}

// IsInf returns true if the scalar is infinite.
func (s ScalarInt26_6) IsInf() bool {
	return false
}

// IsZero returns true if the scalar is zero.
func (s ScalarInt26_6) IsZero() bool {
	return s == 0
}

// IsFinite returns true if the scalar is finite (not NaN or infinite).
func (s ScalarInt26_6) IsFinite() bool {
	// Since ScalarInt26_6 is a fixed-point integer type, it is always finite.
	return true
}

// Clamp returns the scalar clamped between the given minimum and maximum values.
func (s ScalarInt26_6) Clamp(min, max Scalar) Scalar {
	f := s
	fmin := min.(ScalarInt26_6)
	fmax := max.(ScalarInt26_6)
	if f < fmin {
		return fmin
	}
	if f > fmax {
		return fmax
	}
	return f
}

// Abs returns the absolute value of the scalar.
func (s ScalarInt26_6) Abs() Scalar {
	if s < 0 {
		return -s
	}
	return s
}

// Pow returns the current scalar raised to the power of another scalar.
func (s ScalarInt26_6) Pow(o Scalar) Scalar {
	base := float64(s) / 64.0
	exp := float64(o.(ScalarInt26_6)) / 64.0
	return ScalarInt26_6(int32(math.Round(math.Pow(base, exp) * 64)))
}

// Sqrt returns the square root of the scalar.
func (s ScalarInt26_6) Sqrt() Scalar {
	val := float64(s) / 64.0
	return ScalarInt26_6(int32(math.Round(math.Sqrt(val) * 64)))
}

// Floor returns the greatest integer value less than or equal to the scalar.
func (s ScalarInt26_6) Floor() Scalar {
	return ScalarInt26_6((s >> 6) << 6)
}

// Ceil returns the smallest integer value greater than or equal to the scalar.
func (s ScalarInt26_6) Ceil() Scalar {
	if s&0x3F == 0 {
		return s
	}
	return ScalarInt26_6(((s >> 6) + 1) << 6)
}

// Round returns the nearest integer to the scalar.
func (s ScalarInt26_6) Round() Scalar {
	frac := s & 0x3F
	if frac >= 32 {
		return ScalarInt26_6(((s >> 6) + 1) << 6)
	}
	return ScalarInt26_6((s >> 6) << 6)
}

// Sin returns the sine of the scalar (in radians).
func (s ScalarInt26_6) Sin() Scalar {
	val := float64(s) / 64.0
	return ScalarInt26_6(int32(math.Round(math.Sin(val) * 64)))
}

// Cos returns the cosine of the scalar (in radians).
func (s ScalarInt26_6) Cos() Scalar {
	val := float64(s) / 64.0
	return ScalarInt26_6(int32(math.Round(math.Cos(val) * 64)))
}

// Tan returns the tangent of the scalar (in radians).
func (s ScalarInt26_6) Tan() Scalar {
	val := float64(s) / 64.0
	return ScalarInt26_6(int32(math.Round(math.Tan(val) * 64)))
}

// Acos returns the arccosine of the scalar (in radians).
func (s ScalarInt26_6) Acos() Scalar {
	val := float64(s) / 64.0
	return ScalarInt26_6(int32(math.Round(math.Acos(val) * 64)))
}

// Atan returns the arctangent of the scalar (in radians).
func (s ScalarInt26_6) Atan() Scalar {
	val := float64(s) / 64.0
	return ScalarInt26_6(int32(math.Round(math.Atan(val) * 64)))
}

// Atan2 returns the arctangent of the quotient of the current scalar and another scalar.
func (s ScalarInt26_6) Atan2(o Scalar) Scalar {
	a := float64(s) / 64.0
	b := float64(o.(ScalarInt26_6)) / 64.0
	return ScalarInt26_6(int32(math.Round(math.Atan2(a, b) * 64)))
}

// String returns the string representation of the scalar.
func (s ScalarInt26_6) String() string {
	const shift, mask = 6, 1<<6 - 1
	if s >= 0 {
		return fmt.Sprintf("%d:%02d", int32(s>>shift), int32(s&mask))
	}
	ns := -s
	if ns >= 0 {
		return fmt.Sprintf("-%d:%02d", int32(ns>>shift), int32(ns&mask))
	}
	return "-33554432:00" // The minimum value is -(1<<25).
}
