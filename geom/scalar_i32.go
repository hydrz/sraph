package geom

import (
	"math"
	"strconv"
)

var _ Scalar = ScalarI32(0)

type ScalarI32 int32

// Add returns the sum of the current scalar and another scalar.
func (s ScalarI32) Add(o Scalar) Scalar {
	return s + o.(ScalarI32)
}

// Sub returns the difference between the current scalar and another scalar.
func (s ScalarI32) Sub(o Scalar) Scalar {
	return s - o.(ScalarI32)
}

// Mul returns the product of the current scalar and another scalar.
func (s ScalarI32) Mul(o Scalar) Scalar {
	return s * o.(ScalarI32)
}

// Div returns the quotient of the current scalar divided by another scalar.
func (s ScalarI32) Div(o Scalar) Scalar {
	return s / o.(ScalarI32)
}

// Neg returns the negated value of the scalar.
func (s ScalarI32) Neg() Scalar {
	return -s
}

// Equal returns true if the current scalar is equal to another scalar.
func (s ScalarI32) Equal(o Scalar) bool {
	return s == o.(ScalarI32)
}

// Less returns true if the current scalar is less than another scalar.
func (s ScalarI32) Less(o Scalar) bool {
	return s < o.(ScalarI32)
}

// LessEqual returns true if the current scalar is less than or equal to another scalar.
func (s ScalarI32) LessEqual(o Scalar) bool {
	return s <= o.(ScalarI32)
}

// Greater returns true if the current scalar is greater than another scalar.
func (s ScalarI32) Greater(o Scalar) bool {
	return s > o.(ScalarI32)
}

// GreaterEqual returns true if the current scalar is greater than or equal to another scalar.
func (s ScalarI32) GreaterEqual(o Scalar) bool {
	return s >= o.(ScalarI32)
}

// NotEqual returns true if the current scalar is not equal to another scalar.
func (s ScalarI32) NotEqual(o Scalar) bool {
	return s != o.(ScalarI32)
}

// Min returns the minimum value among the current scalar and the provided scalars.
func (s ScalarI32) Min(o ...Scalar) Scalar {
	min := s
	for _, v := range o {
		f := v.(ScalarI32)
		if f < min {
			min = f
		}
	}
	return min
}

// Max returns the maximum value among the current scalar and the provided scalars.
func (s ScalarI32) Max(o ...Scalar) Scalar {
	max := s
	for _, v := range o {
		f := v.(ScalarI32)
		if f > max {
			max = f
		}
	}
	return max
}

// IsNan returns true if the scalar is NaN (not a number).
func (s ScalarI32) IsNan() bool {
	return false
}

// IsInf returns true if the scalar is infinite.
func (s ScalarI32) IsInf() bool {
	return false
}

// IsZero returns true if the scalar is zero.
func (s ScalarI32) IsZero() bool {
	return s == 0
}

// IsFinite returns true if the scalar is a finite number (not NaN or infinite).
func (s ScalarI32) IsFinite() bool {
	return true
}

// Clamp returns the scalar clamped between the given minimum and maximum values.
func (s ScalarI32) Clamp(min, max Scalar) Scalar {
	f := s
	fmin := min.(ScalarI32)
	fmax := max.(ScalarI32)
	if f < fmin {
		return fmin
	}
	if f > fmax {
		return fmax
	}
	return f
}

// Abs returns the absolute value of the scalar.
func (s ScalarI32) Abs() Scalar {
	if s < 0 {
		return -s
	}
	return s
}

// Pow returns the current scalar raised to the power of another scalar.
func (s ScalarI32) Pow(o Scalar) Scalar {
	return ScalarI32(int32(math.Pow(float64(s), float64(o.(ScalarI32)))))
}

// Sqrt returns the square root of the scalar.
func (s ScalarI32) Sqrt() Scalar {
	return ScalarI32(int32(math.Sqrt(float64(s))))
}

// Floor returns the greatest integer value less than or equal to the scalar.
func (s ScalarI32) Floor() Scalar {
	return s // Already integer
}

// Ceil returns the smallest integer value greater than or equal to the scalar.
func (s ScalarI32) Ceil() Scalar {
	return s // Already integer
}

// Round returns the nearest integer to the scalar.
func (s ScalarI32) Round() Scalar {
	return s // Already integer
}

// Sin returns the sine of the scalar (in radians).
func (s ScalarI32) Sin() Scalar {
	return ScalarI32(int32(math.Sin(float64(s))))
}

// Cos returns the cosine of the scalar (in radians).
func (s ScalarI32) Cos() Scalar {
	return ScalarI32(int32(math.Cos(float64(s))))
}

// Tan returns the tangent of the scalar (in radians).
func (s ScalarI32) Tan() Scalar {
	return ScalarI32(int32(math.Tan(float64(s))))
}

// Acos returns the arccosine of the scalar (in radians).
func (s ScalarI32) Acos() Scalar {
	return ScalarI32(int32(math.Acos(float64(s))))
}

// Atan returns the arctangent of the scalar (in radians).
func (s ScalarI32) Atan() Scalar {
	return ScalarI32(int32(math.Atan(float64(s))))
}

// Atan2 returns the arctangent of the quotient of the current scalar and another scalar.
func (s ScalarI32) Atan2(o Scalar) Scalar {
	return ScalarI32(int32(math.Atan2(float64(s), float64(o.(ScalarI32)))))
}

// String returns the string representation of the scalar.
func (s ScalarI32) String() string {
	return strconv.Itoa(int(s))
}
