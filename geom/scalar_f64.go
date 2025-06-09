package geom

import (
	"math"
	"strconv"
)

var _ Scalar = ScalarF64(0)

type ScalarF64 float64

// Add returns the sum of the current scalar and another scalar.
func (s ScalarF64) Add(o Scalar) Scalar {
	return s + o.(ScalarF64)
}

// Sub returns the difference between the current scalar and another scalar.
func (s ScalarF64) Sub(o Scalar) Scalar {
	return s - o.(ScalarF64)
}

// Mul returns the product of the current scalar and another scalar.
func (s ScalarF64) Mul(o Scalar) Scalar {
	return s * o.(ScalarF64)
}

// Div returns the quotient of the current scalar divided by another scalar.
func (s ScalarF64) Div(o Scalar) Scalar {
	return s / o.(ScalarF64)
}

// Neg returns the negated value of the scalar.
func (s ScalarF64) Neg() Scalar {
	return -s
}

// Equal returns true if the current scalar is equal to another scalar.
func (s ScalarF64) Equal(o Scalar) bool {
	const epsilon = ScalarF64(1e+6) // 0.000001
	diff := s.Sub(o)
	if diff.Less(ScalarF64(0)) {
		diff = diff.Neg()
	}
	return diff.LessEqual(epsilon)
}

// Less returns true if the current scalar is less than another scalar.
func (s ScalarF64) Less(o Scalar) bool {
	return s < o.(ScalarF64)
}

// LessEqual returns true if the current scalar is less than or equal to another scalar.
func (s ScalarF64) LessEqual(o Scalar) bool {
	return s <= o.(ScalarF64)
}

// Greater returns true if the current scalar is greater than another scalar.
func (s ScalarF64) Greater(o Scalar) bool {
	return s > o.(ScalarF64)
}

// GreaterEqual returns true if the current scalar is greater than or equal to another scalar.
func (s ScalarF64) GreaterEqual(o Scalar) bool {
	return s >= o.(ScalarF64)
}

// NotEqual returns true if the current scalar is not equal to another scalar.
func (s ScalarF64) NotEqual(o Scalar) bool {
	return !(s.Equal(o))
}

// Min returns the minimum value among the current scalar and the provided scalars.
func (s ScalarF64) Min(o ...Scalar) Scalar {
	min := s
	for _, v := range o {
		f := v.(ScalarF64)
		if f < min {
			min = f
		}
	}
	return min
}

// Max returns the maximum value among the current scalar and the provided scalars.
func (s ScalarF64) Max(o ...Scalar) Scalar {
	max := s
	for _, v := range o {
		f := v.(ScalarF64)
		if f > max {
			max = f
		}
	}
	return max
}

// IsNan returns true if the scalar is NaN (not a number).
func (s ScalarF64) IsNan() bool {
	return math.IsNaN(float64(s))
}

// IsInf returns true if the scalar is infinite.
func (s ScalarF64) IsInf() bool {
	return math.IsInf(float64(s), 0)
}

// IsZero returns true if the scalar is zero.
func (s ScalarF64) IsZero() bool {
	return float64(s) == 0
}

// IsFinite returns true if the scalar is finite (not NaN or infinite).
func (s ScalarF64) IsFinite() bool {
	return !s.IsNan() && !s.IsInf()
}

// Clamp returns the scalar clamped between the given minimum and maximum values.
func (s ScalarF64) Clamp(min, max Scalar) Scalar {
	f := s
	fmin := min.(ScalarF64)
	fmax := max.(ScalarF64)
	if f < fmin {
		return fmin
	}
	if f > fmax {
		return fmax
	}
	return f
}

// Abs returns the absolute value of the scalar.
func (s ScalarF64) Abs() Scalar {
	return ScalarF64(math.Abs(float64(s)))
}

// Pow returns the current scalar raised to the power of another scalar.
func (s ScalarF64) Pow(o Scalar) Scalar {
	return ScalarF64(math.Pow(float64(s), float64(o.(ScalarF64))))
}

// Sqrt returns the square root of the scalar.
func (s ScalarF64) Sqrt() Scalar {
	return ScalarF64(math.Sqrt(float64(s)))
}

// Floor returns the greatest integer value less than or equal to the scalar.
func (s ScalarF64) Floor() Scalar {
	return ScalarF64(math.Floor(float64(s)))
}

// Ceil returns the smallest integer value greater than or equal to the scalar.
func (s ScalarF64) Ceil() Scalar {
	return ScalarF64(math.Ceil(float64(s)))
}

// Round returns the nearest integer to the scalar.
func (s ScalarF64) Round() Scalar {
	return ScalarF64(math.Round(float64(s)))
}

// Sin returns the sine of the scalar (in radians).
func (s ScalarF64) Sin() Scalar {
	return ScalarF64(math.Sin(float64(s)))
}

// Cos returns the cosine of the scalar (in radians).
func (s ScalarF64) Cos() Scalar {
	return ScalarF64(math.Cos(float64(s)))
}

// Tan returns the tangent of the scalar (in radians).
func (s ScalarF64) Tan() Scalar {
	return ScalarF64(math.Tan(float64(s)))
}

// Acos returns the arccosine of the scalar (in radians).
func (s ScalarF64) Acos() Scalar {
	return ScalarF64(math.Acos(float64(s)))
}

// Atan returns the arctangent of the scalar (in radians).
func (s ScalarF64) Atan() Scalar {
	return ScalarF64(math.Atan(float64(s)))
}

// Atan2 returns the arctangent of the quotient of the current scalar and another scalar.
func (s ScalarF64) Atan2(o Scalar) Scalar {
	return ScalarF64(math.Atan2(float64(s), float64(o.(ScalarF64))))
}

// String returns the string representation of the scalar.
func (s ScalarF64) String() string {
	return strconv.FormatFloat(float64(s), 'f', -1, 32)
}
