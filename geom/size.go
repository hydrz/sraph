package geom

import (
	"fmt"
	"math"
)

// Size defines the interface for a two-dimensional size.
// It provides methods for arithmetic operations, comparison, and property queries.
type Size interface {
	// Width returns the width component.
	Width() Scalar
	// Height returns the height component.
	Height() Scalar

	// Add returns the element-wise sum of this size and another.
	Add(other Size) Size
	// Sub returns the element-wise difference of this size and another.
	Sub(other Size) Size
	// Mul returns the element-wise product of this size and another.
	Mul(other Size) Size
	// Div returns the element-wise quotient of this size and another.
	Div(other Size) Size
	// Neg returns the negation of this size.
	Neg() Size
	// Equal reports whether this size and another are equal within floating-point tolerance.
	Equal(other Size) bool

	// Scale returns a size scaled by the given factor in both dimensions.
	Scale(scale Scalar) Size
	// ScaleWH returns a size scaled by width and height factors.
	ScaleWH(width, height Scalar) Size
	// MinDimension returns the smaller of width and height.
	MinDimension() Scalar
	// MaxDimension returns the larger of width and height.
	MaxDimension() Scalar
	// Area returns the area of the size (width * height).
	Area() Scalar
	// Abs returns a size with non-negative width and height.
	Abs() Size
	// Floor returns a size with width and height rounded down to the nearest integer.
	Floor() Size
	// Ceil returns a size with width and height rounded up to the nearest integer.
	Ceil() Size
	// Round returns a size with width and height rounded to the nearest integer.
	Round() Size

	// IsZero reports whether both width and height are zero.
	IsZero() bool
	// IsFinite reports whether both width and height are finite numbers.
	IsFinite() bool
	// IsInfinite reports whether either width or height is infinite.
	IsInfinite() bool
	// IsSquare reports whether width and height are equal within floating-point tolerance.
	IsSquare() bool

	// MipCount returns the number of mipmap levels for the size. The result is at least 1.
	MipCount() int

	// String returns the string representation of the size.
	String() string
}

func NewSize[T Number](width, height T) Size {
	return size[T]{
		width:  width,
		height: height,
	}
}

// size represents a two-dimensional size with width and height.
// size is safe for concurrent use by multiple goroutines.
type size[T Number] struct {
	width, height T
}

// Width returns the width component.
func (s size[T]) Width() Scalar {
	return Scalar(s.width)
}

// Height returns the height component.
func (s size[T]) Height() Scalar {
	return Scalar(s.height)
}

// Add returns the element-wise sum of this size and another.
func (s size[T]) Add(other Size) Size {
	o := other.(size[T])
	return size[T]{width: s.width + o.width, height: s.height + o.height}
}

// Sub returns the element-wise difference of this size and another.
func (s size[T]) Sub(other Size) Size {
	o := other.(size[T])
	return size[T]{width: s.width - o.width, height: s.height - o.height}
}

// Mul returns the element-wise product of this size and another.
func (s size[T]) Mul(other Size) Size {
	o := other.(size[T])
	return size[T]{width: s.width * o.width, height: s.height * o.height}
}

// Div returns the element-wise quotient of this size and another.
func (s size[T]) Div(other Size) Size {
	o := other.(size[T])
	return size[T]{width: s.width / o.width, height: s.height / o.height}
}

// Neg returns the negation of this size.
func (s size[T]) Neg() Size {
	return size[T]{width: -s.width, height: -s.height}
}

// Equal reports whether this size and another are equal within floating-point tolerance.
func (s size[T]) Equal(other Size) bool {
	o := other.(size[T])
	return NearlyEqual(s.width, o.width) && NearlyEqual(s.height, o.height)
}

// Scale returns a size scaled by the given factor in both dimensions.
func (s size[T]) Scale(scale Scalar) Size {
	return size[T]{width: s.width * T(scale), height: s.height * T(scale)}
}

// ScaleWH returns a size scaled by width and height factors.
func (s size[T]) ScaleWH(width Scalar, height Scalar) Size {
	return size[T]{width: s.width * T(width), height: s.height * T(height)}
}

// MinDimension returns the smaller of width and height.
func (s size[T]) MinDimension() Scalar {
	if s.width < s.height {
		return Scalar(s.width)
	}
	return Scalar(s.height)
}

// MaxDimension returns the larger of width and height.
func (s size[T]) MaxDimension() Scalar {
	if s.width > s.height {
		return Scalar(s.width)
	}
	return Scalar(s.height)
}

// Area returns the area of the size (width * height).
func (s size[T]) Area() Scalar {
	return Scalar(s.width) * Scalar(s.height)
}

// Abs returns a size with non-negative width and height.
func (s size[T]) Abs() Size {
	var zero T
	w := s.width
	h := s.height
	if w < zero {
		w = -w
	}
	if h < zero {
		h = -h
	}
	return size[T]{width: w, height: h}
}

// Floor returns a size with width and height rounded down to the nearest integer.
func (s size[T]) Floor() Size {
	return size[T]{
		width:  T(math.Floor(float64(s.width))),
		height: T(math.Floor(float64(s.height))),
	}
}

// Ceil returns a size with width and height rounded up to the nearest integer.
func (s size[T]) Ceil() Size {
	return size[T]{
		width:  T(math.Ceil(float64(s.width))),
		height: T(math.Ceil(float64(s.height))),
	}
}

// Round returns a size with width and height rounded to the nearest integer.
func (s size[T]) Round() Size {
	return size[T]{
		width:  T(math.Round(float64(s.width))),
		height: T(math.Round(float64(s.height))),
	}
}

// IsZero reports whether both width and height are zero.
func (s size[T]) IsZero() bool {
	var zero T
	return NearlyEqual(s.width, zero) && NearlyEqual(s.height, zero)
}

// IsFinite reports whether both width and height are finite numbers.
func (s size[T]) IsFinite() bool {
	return IsFinite(s.width) && IsFinite(s.height)
}

// IsInfinite reports whether either width or height is infinite.
func (s size[T]) IsInfinite() bool {
	return math.IsInf(ToFloat64(s.width), 0) || math.IsInf(ToFloat64(s.height), 0)
}

// IsSquare reports whether width and height are equal within floating-point tolerance.
func (s size[T]) IsSquare() bool {
	return NearlyEqual(s.width, s.height)
}

// MipCount returns the number of mipmap levels for the size. The result is at least 1.
func (s size[T]) MipCount() int {
	w := int(s.width)
	h := int(s.height)
	count := 0
	for w > 1 || h > 1 {
		if w > 1 {
			w = w >> 1
		}
		if h > 1 {
			h = h >> 1
		}
		count++
	}
	return count + 1
}

// String returns the string representation of the size.
func (s size[T]) String() string {
	return fmt.Sprintf("(%v, %v)", s.width, s.height)
}
