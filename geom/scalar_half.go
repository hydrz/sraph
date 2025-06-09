package geom

import (
	"math"
	"strconv"
	"unsafe"
)

var _ Scalar = ScalarHalf(0)

// Half represents a 16-bit floating point number for GPU optimization.
// This provides memory-efficient storage for graphics operations where
// full 32-bit precision is not required. Half-precision floats are commonly
// used in GPU shaders and vertex data to reduce memory bandwidth.
// All operations on Half do not modify the original, but return a new Half.
type ScalarHalf uint16

// Add returns the sum of the current scalar and another scalar.
func (s ScalarHalf) Add(o Scalar) Scalar {
	return ScalarHalf(s.ToFloat32() + o.(ScalarHalf).ToFloat32())
}

// Sub returns the difference between the current scalar and another scalar.
func (s ScalarHalf) Sub(o Scalar) Scalar {
	return ScalarHalf(s.ToFloat32() - o.(ScalarHalf).ToFloat32())
}

// Mul returns the product of the current scalar and another scalar.
func (s ScalarHalf) Mul(o Scalar) Scalar {
	return ScalarHalf(s.ToFloat32() * o.(ScalarHalf).ToFloat32())
}

// Div returns the quotient of the current scalar divided by another scalar.
func (s ScalarHalf) Div(o Scalar) Scalar {
	return ScalarHalf(s.ToFloat32() / o.(ScalarHalf).ToFloat32())
}

// Neg returns the negated value of the scalar.
func (s ScalarHalf) Neg() Scalar {
	return s ^ 0x8000 // Toggle the sign bit
}

// Equal returns true if the current scalar is equal to another scalar.
func (s ScalarHalf) Equal(o Scalar) bool {
	const epsilon = ScalarHalf(0x1400) // 0.001 in half-precision
	diff := s - o.(ScalarHalf)
	if diff < 0 {
		diff = -diff
	}
	return diff <= epsilon
}

// Less returns true if the current scalar is less than another scalar.
func (s ScalarHalf) Less(o Scalar) bool {
	return s.ToFloat32() < o.(ScalarHalf).ToFloat32()
}

// LessEqual returns true if the current scalar is less than or equal to another scalar.
func (s ScalarHalf) LessEqual(o Scalar) bool {
	return s.ToFloat32() <= o.(ScalarHalf).ToFloat32()
}

// Greater returns true if the current scalar is greater than another scalar.
func (s ScalarHalf) Greater(o Scalar) bool {
	return s.ToFloat32() > o.(ScalarHalf).ToFloat32()
}

// GreaterEqual returns true if the current scalar is greater than or equal to another scalar.
func (s ScalarHalf) GreaterEqual(o Scalar) bool {
	return s.ToFloat32() >= o.(ScalarHalf).ToFloat32()
}

// NotEqual returns true if the current scalar is not equal to another scalar.
func (s ScalarHalf) NotEqual(o Scalar) bool {
	return !s.Equal(o)
}

// Min returns the minimum value among the current scalar and the provided scalars.
func (s ScalarHalf) Min(o ...Scalar) Scalar {
	minVal := s
	for _, v := range o {
		if minVal.Greater(v) {
			minVal = v.(ScalarHalf)
		}
	}
	return minVal
}

// Max returns the maximum value among the current scalar and the provided scalars.
func (s ScalarHalf) Max(o ...Scalar) Scalar {
	maxVal := s
	for _, v := range o {
		if maxVal.Less(v) {
			maxVal = v.(ScalarHalf)
		}
	}
	return maxVal
}

// IsNan returns true if the scalar is NaN (not a number).
func (s ScalarHalf) IsNan() bool {
	return (s & 0x7FFF) > 0x7C00
}

// IsInf returns true if the scalar is infinite.
func (s ScalarHalf) IsInf() bool {
	return (s & 0x7FFF) == 0x7C00
}

// IsZero returns true if the scalar is zero.
func (s ScalarHalf) IsZero() bool {
	return s&0x7FFF == 0 // Ignore sign bit
}

// IsFinite returns true if the scalar is finite (not NaN or infinite).
func (s ScalarHalf) IsFinite() bool {
	return !s.IsNan() && !s.IsInf()
}

// Clamp returns the scalar clamped between the given minimum and maximum values.
func (s ScalarHalf) Clamp(min, max Scalar) Scalar {
	minHalf := min.(ScalarHalf)
	maxHalf := max.(ScalarHalf)

	if s < minHalf {
		return minHalf
	} else if s > maxHalf {
		return maxHalf
	}
	return s
}

// Abs returns the absolute value of the scalar.
func (s ScalarHalf) Abs() Scalar {
	return ScalarHalf(s & 0x7FFF) // Clear the sign bit
}

// Pow returns the current scalar raised to the power of another scalar.
func (s ScalarHalf) Pow(o Scalar) Scalar {
	f := math.Pow(float64(s.ToFloat32()), float64(o.(ScalarHalf).ToFloat32()))
	return float32ToHalf(float32(f))
}

// Sqrt returns the square root of the scalar.
func (s ScalarHalf) Sqrt() Scalar {
	f := math.Sqrt(float64(s.ToFloat32()))
	return float32ToHalf(float32(f))
}

// Floor returns the greatest integer value less than or equal to the scalar.
func (s ScalarHalf) Floor() Scalar {
	f := math.Floor(float64(s.ToFloat32()))
	return float32ToHalf(float32(f))
}

// Ceil returns the smallest integer value greater than or equal to the scalar.
func (s ScalarHalf) Ceil() Scalar {
	f := math.Ceil(float64(s.ToFloat32()))
	return float32ToHalf(float32(f))
}

// Round returns the nearest integer to the scalar.
func (s ScalarHalf) Round() Scalar {
	f := math.Round(float64(s.ToFloat32()))
	return float32ToHalf(float32(f))
}

// Sin returns the sine of the scalar (in radians).
func (s ScalarHalf) Sin() Scalar {
	f := math.Sin(float64(s.ToFloat32()))
	return float32ToHalf(float32(f))
}

// Cos returns the cosine of the scalar (in radians).
func (s ScalarHalf) Cos() Scalar {
	f := math.Cos(float64(s.ToFloat32()))
	return float32ToHalf(float32(f))
}

// Tan returns the tangent of the scalar (in radians).
func (s ScalarHalf) Tan() Scalar {
	f := math.Tan(float64(s.ToFloat32()))
	return float32ToHalf(float32(f))
}

// Acos returns the arccosine of the scalar (in radians).
func (s ScalarHalf) Acos() Scalar {
	f := math.Acos(float64(s.ToFloat32()))
	return float32ToHalf(float32(f))
}

// Atan returns the arctangent of the scalar (in radians).
func (s ScalarHalf) Atan() Scalar {
	f := math.Atan(float64(s.ToFloat32()))
	return float32ToHalf(float32(f))
}

// Atan2 returns the arctangent of the quotient of the current scalar and another scalar.
func (s ScalarHalf) Atan2(o Scalar) Scalar {
	f := math.Atan2(float64(s.ToFloat32()), float64(o.(ScalarHalf).ToFloat32()))
	return float32ToHalf(float32(f))
}

// String returns the string representation of the scalar.
func (s ScalarHalf) String() string {
	return strconv.FormatFloat(float64(s.ToFloat32()), 'f', -1, 16)
}

func (s ScalarHalf) ToFloat32() float32 {
	bits := uint16(s)
	sign := (bits >> 15) & 0x0001
	exp := (bits >> 10) & 0x001F
	mant := bits & 0x03FF

	var f uint32
	if exp == 0 {
		if mant == 0 {
			f = uint32(sign) << 31
		} else {
			// Denormalized
			expVal := -14
			mantVal := float32(mant) / 1024.0
			val := math.Ldexp(float64(mantVal), expVal)
			if sign != 0 {
				val = -val
			}
			return float32(val)
		}
	} else if exp == 31 {
		f = (uint32(sign) << 31) | 0x7F800000 | (uint32(mant) << 13)
	} else {
		f = (uint32(sign) << 31) | ((uint32(exp) + 112) << 23) | (uint32(mant) << 13)
	}
	return *(*float32)(unsafe.Pointer(&f))
}

// Helper: convert float32 to ScalarHalf.
func float32ToHalf(f float32) ScalarHalf {
	// Use unsafe to get the bit representation
	bits := *(*uint32)(unsafe.Pointer(&f))

	sign := (bits >> 16) & 0x8000
	expSigned := int(((bits >> 23) & 0xFF)) - 127 + 15
	mantissa := bits & 0x7FFFFF

	// Handle special cases
	if expSigned <= 0 {
		// Underflow to zero or denormal
		if expSigned < -10 {
			return ScalarHalf(sign)
		}
		// Denormal
		mantissa = (mantissa | 0x800000) >> uint32(1-expSigned)
		return ScalarHalf(sign | (mantissa >> 13))
	}

	exp := uint32(expSigned)

	if exp >= 31 {
		// Overflow to infinity or NaN
		if exp == 128 && mantissa != 0 {
			// NaN
			return ScalarHalf(sign | 0x7C00 | (mantissa >> 13))
		}
		// Infinity
		return ScalarHalf(sign | 0x7C00)
	}

	// Normal case
	return ScalarHalf(sign | (exp << 10) | (mantissa >> 13))
}
