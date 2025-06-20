package geom

import (
	"math"
)

// Size represents a 2D size (width and height) in graphics programming.
// It is commonly used for describing the dimensions of rectangles, images, viewports, and other graphical objects.
// The struct provides a set of arithmetic and utility operations for manipulating and querying size values.
type Size[T Number] struct {
	Width, Height T
}

// Add returns the element-wise sum of this size and another.
// Useful for combining the dimensions of two graphical objects.
func (s Size[T]) Add(other Size[T]) Size[T] {
	return Size[T]{
		Width:  s.Width + other.Width,
		Height: s.Height + other.Height,
	}
}

// Sub returns the element-wise difference of this size and another.
// Useful for calculating the remaining space or difference between two objects.
func (s Size[T]) Sub(other Size[T]) Size[T] {
	return Size[T]{
		Width:  s.Width - other.Width,
		Height: s.Height - other.Height,
	}
}

// Mul returns the element-wise product of this size and another.
// Can be used for scaling each dimension by another size, e.g., for proportional resizing.
func (s Size[T]) Mul(other Size[T]) Size[T] {
	return Size[T]{
		Width:  s.Width * other.Width,
		Height: s.Height * other.Height,
	}
}

// Div returns the element-wise division of this size by another.
// Useful for computing relative scaling factors or normalizing dimensions.
func (s Size[T]) Div(other Size[T]) Size[T] {
	return Size[T]{
		Width:  s.Width / other.Width,
		Height: s.Height / other.Height,
	}
}

// Neg returns the negated size.
// Rare in graphics, but can be used for certain mathematical operations.
func (s Size[T]) Neg() Size[T] {
	return Size[T]{
		Width:  -s.Width,
		Height: -s.Height,
	}
}

// Scale scales the width and height by the same factor.
// Commonly used for uniform scaling, such as resizing an image while maintaining aspect ratio.
func (s Size[T]) Scale(scale T) Size[T] {
	return Size[T]{
		Width:  s.Width * scale,
		Height: s.Height * scale,
	}
}

// ScaleWH scales the width and height by the given factors.
// Allows non-uniform scaling, e.g., stretching or shrinking only one dimension.
func (s Size[T]) ScaleWH(width, height T) Size[T] {
	return Size[T]{
		Width:  s.Width * width,
		Height: s.Height * height,
	}
}

// Eq reports whether s and o are equal.
func (s Size[T]) Eq(o Size[T]) bool {
	return NearlyEqual(s.Width, o.Width) && NearlyEqual(s.Height, o.Height)
}

// Min returns the size with the minimum width and height among all.
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

// Max returns the size with the maximum width and height among all.
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

// MinDimension returns the minimum of width and height.
// Used to determine the limiting dimension, e.g., for fitting a square inside a rectangle.
func (s Size[T]) MinDimension() T {
	if s.Width < s.Height {
		return s.Width
	}
	return s.Height
}

// MaxDimension returns the maximum of width and height.
// Used to determine the dominant dimension, e.g., for scaling or aspect ratio calculations.
func (s Size[T]) MaxDimension() T {
	if s.Width > s.Height {
		return s.Width
	}
	return s.Height
}

// Area returns the area (width * height).
// Fundamental in graphics for pixel count, memory allocation, or hit-testing.
func (s Size[T]) Area() T {
	return s.Width * s.Height
}

// Abs returns the size with absolute width and height.
// Ensures dimensions are non-negative, which is important for rendering and layout.
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

// Floor returns the size with width and height floored.
// Useful for aligning to pixel boundaries or integer grid systems.
func (s Size[T]) Floor() Size[T] {
	return Size[T]{
		Width:  T(math.Floor(ToFloat64(s.Width))),
		Height: T(math.Floor(ToFloat64(s.Height))),
	}
}

// Ceil returns the size with width and height ceiled.
// Useful for ensuring enough space is allocated, avoiding clipping.
func (s Size[T]) Ceil() Size[T] {
	return Size[T]{
		Width:  T(math.Ceil(ToFloat64(s.Width))),
		Height: T(math.Ceil(ToFloat64(s.Height))),
	}
}

// Round returns the size with width and height rounded.
// Used for snapping to the nearest pixel or grid unit.
func (s Size[T]) Round() Size[T] {
	return Size[T]{
		Width:  T(math.Round(ToFloat64(s.Width))),
		Height: T(math.Round(ToFloat64(s.Height))),
	}
}

// IsZero returns true if both width and height are zero.
// Used to detect degenerate or empty objects.
func (s Size[T]) IsZero() bool {
	var zero T
	return NearlyEqual(s.Width, zero) && NearlyEqual(s.Height, zero)
}

// IsFinite returns true if both width and height are finite.
// Important for validating geometry before rendering or computation.
func (s Size[T]) IsFinite() bool {
	return IsFinite(s.Width) && IsFinite(s.Height)
}

// IsInfinite returns true if either width or height is infinite.
// Can be used to represent unbounded or unconstrained objects.
func (s Size[T]) IsInfinite() bool {
	return math.IsInf(ToFloat64(s.Width), 0) || math.IsInf(ToFloat64(s.Height), 0)
}

// IsSquare returns true if width and height are equal.
// Useful for aspect ratio checks, e.g., for icons or tiles.
func (s Size[T]) IsSquare() bool {
	return NearlyEqual(s.Width, s.Height)
}

// MipCount returns the mipmap count for the size.
// Used in texture mapping to determine the number of mipmap levels for an image.
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

// String returns a string representation of the size, like "(width, height)".
func (s Size[T]) String() string {
	return "(" + ToString(s.Width) + ", " + ToString(s.Height) + ")"
}
