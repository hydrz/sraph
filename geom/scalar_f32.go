package geom

import (
	"math"
	"strconv"
)

var _ Scalar = ScalarF32(0)

type ScalarF32 float32

// Add returns the sum of the current scalar and another scalar.
func (s ScalarF32) Add(o Scalar) Scalar {
	return s + o.(ScalarF32)
}

// Sub returns the difference between the current scalar and another scalar.
func (s ScalarF32) Sub(o Scalar) Scalar {
	return s - o.(ScalarF32)
}

// Mul returns the product of the current scalar and another scalar.
func (s ScalarF32) Mul(o Scalar) Scalar {
	return s * o.(ScalarF32)
}

// Div returns the quotient of the current scalar divided by another scalar.
func (s ScalarF32) Div(o Scalar) Scalar {
	return s / o.(ScalarF32)
}

// Neg returns the negated value of the scalar.
func (s ScalarF32) Neg() Scalar {
	return -s
}

// Equal returns true if the current scalar is equal to another scalar.
func (s ScalarF32) Equal(o Scalar) bool {
	const epsilon = ScalarF32(1e-3) // 0.001

	diff := s.Sub(o)
	if diff.Less(ScalarF32(0)) {
		diff = diff.Neg()
	}

	return diff.LessEqual(epsilon)
}

// Less returns true if the current scalar is less than another scalar.
func (s ScalarF32) Less(o Scalar) bool {
	return s < o.(ScalarF32)
}

// LessEqual returns true if the current scalar is less than or equal to another scalar.
func (s ScalarF32) LessEqual(o Scalar) bool {
	return s <= o.(ScalarF32)
}

// Greater returns true if the current scalar is greater than another scalar.
func (s ScalarF32) Greater(o Scalar) bool {
	return s > o.(ScalarF32)
}

// GreaterEqual returns true if the current scalar is greater than or equal to another scalar.
func (s ScalarF32) GreaterEqual(o Scalar) bool {
	return s >= o.(ScalarF32)
}

// NotEqual returns true if the current scalar is not equal to another scalar.
func (s ScalarF32) NotEqual(o Scalar) bool {
	return !s.Equal(o)
}

// Min returns the minimum value among the current scalar and the provided scalars.
func (s ScalarF32) Min(o ...Scalar) Scalar {
	min := s
	for _, v := range o {
		f := v.(ScalarF32)
		if f < min {
			min = f
		}
	}
	return min
}

// Max returns the maximum value among the current scalar and the provided scalars.
func (s ScalarF32) Max(o ...Scalar) Scalar {
	max := s
	for _, v := range o {
		f := v.(ScalarF32)
		if f > max {
			max = f
		}
	}
	return max
}

// IsNan returns true if the scalar is NaN (not a number).
func (s ScalarF32) IsNan() bool {
	return math.IsNaN(float64(s))
}

// IsInf returns true if the scalar is infinite.
func (s ScalarF32) IsInf() bool {
	return math.IsInf(float64(s), 0)
}

// IsZero returns true if the scalar is zero.
func (s ScalarF32) IsZero() bool {
	return s == 0
}

// IsFinite returns true if the scalar is finite (not NaN or infinite).
func (s ScalarF32) IsFinite() bool {
	return !s.IsNan() && !s.IsInf()
}

// Clamp returns the scalar clamped between the given minimum and maximum values.
func (s ScalarF32) Clamp(min, max Scalar) Scalar {
	f := s
	fmin := min.(ScalarF32)
	fmax := max.(ScalarF32)
	if f < fmin {
		return fmin
	}
	if f > fmax {
		return fmax
	}
	return s
}

// Abs returns the absolute value of the scalar.
func (s ScalarF32) Abs() Scalar {
	return ScalarF32(math.Abs(float64(s)))
}

// Pow returns the current scalar raised to the power of another scalar.
func (s ScalarF32) Pow(o Scalar) Scalar {
	return ScalarF32(math.Pow(float64(s), float64(o.(ScalarF32))))
}

// Sqrt returns the square root of the scalar.
func (s ScalarF32) Sqrt() Scalar {
	return ScalarF32(math.Sqrt(float64(s)))
}

// Floor returns the greatest integer value less than or equal to the scalar.
func (s ScalarF32) Floor() Scalar {
	return ScalarF32(math.Floor(float64(s)))
}

// Ceil returns the smallest integer value greater than or equal to the scalar.
func (s ScalarF32) Ceil() Scalar {
	return ScalarF32(math.Ceil(float64(s)))
}

// Round returns the nearest integer to the scalar.
func (s ScalarF32) Round() Scalar {
	return ScalarF32(math.Round(float64(s)))
}

// Sin returns the sine of the scalar (in radians).
func (s ScalarF32) Sin() Scalar {
	return ScalarF32(math.Sin(float64(s)))
}

// Cos returns the cosine of the scalar (in radians).
func (s ScalarF32) Cos() Scalar {
	return ScalarF32(math.Cos(float64(s)))
}

// Tan returns the tangent of the scalar (in radians).
func (s ScalarF32) Tan() Scalar {
	return ScalarF32(math.Tan(float64(s)))
}

// Acos returns the arccosine of the scalar (in radians).
func (s ScalarF32) Acos() Scalar {
	return ScalarF32(math.Acos(float64(s)))
}

// Atan returns the arctangent of the scalar (in radians).
func (s ScalarF32) Atan() Scalar {
	return ScalarF32(math.Atan(float64(s)))
}

// Atan2 returns the arctangent of the quotient of the current scalar and another scalar.
func (s ScalarF32) Atan2(o Scalar) Scalar {
	return ScalarF32(math.Atan2(float64(s), float64(o.(ScalarF32))))
}

// String returns the string representation of the scalar.
func (s ScalarF32) String() string {
	return strconv.FormatFloat(float64(s), 'f', -1, 32)
}
