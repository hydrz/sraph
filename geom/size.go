package geom

import (
	"math"
)

// Size represents a two-dimensional size with width and height.
// Size is safe for concurrent use by multiple goroutines.
type Size[T TScalar] struct {
	Width, Height T
}

// NewSize returns a Size with the given width and height.
func NewSize[T TScalar](width, height T) Size[T] {
	return Size[T]{
		Width:  width,
		Height: height,
	}
}

// Add returns the element-wise sum of s and other.
func (s Size[T]) Add(other Size[T]) Size[T] {
	return Size[T]{
		Width:  s.Width + other.Width,
		Height: s.Height + other.Height,
	}
}

// Sub returns the element-wise difference of s and other.
func (s Size[T]) Sub(other Size[T]) Size[T] {
	return Size[T]{
		Width:  s.Width - other.Width,
		Height: s.Height - other.Height,
	}
}

// Mul returns the element-wise product of s and other.
func (s Size[T]) Mul(other Size[T]) Size[T] {
	return Size[T]{
		Width:  s.Width * other.Width,
		Height: s.Height * other.Height,
	}
}

// Div returns the element-wise quotient of s and other.
func (s Size[T]) Div(other Size[T]) Size[T] {
	return Size[T]{
		Width:  s.Width / other.Width,
		Height: s.Height / other.Height,
	}
}

// Neg returns the negation of s.
func (s Size[T]) Neg() Size[T] {
	return Size[T]{
		Width:  -s.Width,
		Height: -s.Height,
	}
}

// Scale returns a Size scaled by the given factor in both dimensions.
func (s Size[T]) Scale(scale T) Size[T] {
	return Size[T]{
		Width:  s.Width * scale,
		Height: s.Height * scale,
	}
}

// ScaleWH returns a Size scaled by width and height factors.
func (s Size[T]) ScaleWH(width, height T) Size[T] {
	return Size[T]{
		Width:  s.Width * width,
		Height: s.Height * height,
	}
}

// Equal reports whether s and o have equal width and height within floating-point tolerance.
func (s Size[T]) Equal(o Size[T]) bool {
	return NearlyEqual(s.Width, o.Width) && NearlyEqual(s.Height, o.Height)
}

// Min returns a Size with the minimum width and height among s and all others.
func (s Size[T]) Min(o ...Size[T]) Size[T] {
	minW, minH := s.Width, s.Height
	for _, other := range o {
		if other.Width < minW {
			minW = other.Width
		}
		if other.Height < minH {
			minH = other.Height
		}
	}
	return Size[T]{Width: minW, Height: minH}
}

// Max returns a Size with the maximum width and height among s and all others.
func (s Size[T]) Max(o ...Size[T]) Size[T] {
	maxW, maxH := s.Width, s.Height
	for _, other := range o {
		if other.Width > maxW {
			maxW = other.Width
		}
		if other.Height > maxH {
			maxH = other.Height
		}
	}
	return Size[T]{Width: maxW, Height: maxH}
}

// MinDimension returns the smaller of width and height.
func (s Size[T]) MinDimension() T {
	if s.Width < s.Height {
		return s.Width
	}
	return s.Height
}

// MaxDimension returns the larger of width and height.
func (s Size[T]) MaxDimension() T {
	if s.Width > s.Height {
		return s.Width
	}
	return s.Height
}

// Area returns the area of the size (width * height).
func (s Size[T]) Area() T {
	return s.Width * s.Height
}

// Abs returns a Size with non-negative width and height.
func (s Size[T]) Abs() Size[T] {
	var zero T
	w := s.Width
	h := s.Height
	if w < zero {
		w = -w
	}
	if h < zero {
		h = -h
	}
	return Size[T]{Width: w, Height: h}
}

// Floor returns a Size with width and height rounded down to the nearest integer.
func (s Size[T]) Floor() Size[T] {
	return Size[T]{
		Width:  T(math.Floor(ToFloat64(s.Width))),
		Height: T(math.Floor(ToFloat64(s.Height))),
	}
}

// Ceil returns a Size with width and height rounded up to the nearest integer.
func (s Size[T]) Ceil() Size[T] {
	return Size[T]{
		Width:  T(math.Ceil(ToFloat64(s.Width))),
		Height: T(math.Ceil(ToFloat64(s.Height))),
	}
}

// Round returns a Size with width and height rounded to the nearest integer.
func (s Size[T]) Round() Size[T] {
	return Size[T]{
		Width:  T(math.Round(ToFloat64(s.Width))),
		Height: T(math.Round(ToFloat64(s.Height))),
	}
}

// IsZero reports whether both width and height are zero.
func (s Size[T]) IsZero() bool {
	var zero T
	return NearlyEqual(s.Width, zero) && NearlyEqual(s.Height, zero)
}

// IsFinite reports whether both width and height are finite numbers.
func (s Size[T]) IsFinite() bool {
	return IsFinite(s.Width) && IsFinite(s.Height)
}

// IsInfinite reports whether either width or height is infinite.
func (s Size[T]) IsInfinite() bool {
	return math.IsInf(ToFloat64(s.Width), 0) || math.IsInf(ToFloat64(s.Height), 0)
}

// IsSquare reports whether width and height are equal within floating-point tolerance.
func (s Size[T]) IsSquare() bool {
	return NearlyEqual(s.Width, s.Height)
}

// MipCount returns the number of mipmap levels for the size.
// The result is at least 1.
func (s Size[T]) MipCount() int {
	w := int(s.Width)
	h := int(s.Height)
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

// String returns the string representation of the size in the form "(width, height)".
func (s Size[T]) String() string {
	return "(" + ToString(s.Width) + ", " + ToString(s.Height) + ")"
}
