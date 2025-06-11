package geom

import (
	"math"
)

// Size represents a 2D size (width and height) in graphics programming.
// It is commonly used for describing the dimensions of rectangles, images, viewports, and other graphical objects.
// The interface provides a set of arithmetic and utility operations for manipulating and querying size values.
type Size[T Scalar] interface {
	// Width returns the width component.
	// In graphics, width is used to describe the horizontal extent of a shape, image, or viewport.
	Width() T
	// Height returns the height component.
	// In graphics, height is used to describe the vertical extent of a shape, image, or viewport.
	Height() T

	// Add returns the element-wise sum of this size and another.
	// Useful for combining the dimensions of two graphical objects.
	Add(other Size[T]) Size[T]
	// Sub returns the element-wise difference of this size and another.
	// Useful for calculating the remaining space or difference between two objects.
	Sub(other Size[T]) Size[T]
	// Mul returns the element-wise product of this size and another.
	// Can be used for scaling each dimension by another size, e.g., for proportional resizing.
	Mul(other Size[T]) Size[T]
	// Div returns the element-wise division of this size by another.
	// Useful for computing relative scaling factors or normalizing dimensions.
	Div(other Size[T]) Size[T]
	// Neg returns the negated size.
	// Rare in graphics, but can be used for certain mathematical operations.
	Neg() Size[T]

	// Scale scales the width and height by the same factor.
	// Commonly used for uniform scaling, such as resizing an image while maintaining aspect ratio.
	Scale(scale T) Size[T]
	// ScaleWH scales the width and height by the given factors.
	// Allows non-uniform scaling, e.g., stretching or shrinking only one dimension.
	ScaleWH(width, height T) Size[T]

	// Equal returns true if this size equals another.
	// Used for comparison in layout, collision, or rendering logic.
	Equal(other Size[T]) bool

	// Min returns the size with the minimum width and height among all.
	// Useful for bounding box calculations and fitting objects within constraints.
	Min(o ...Size[T]) Size[T]
	// Max returns the size with the maximum width and height among all.
	// Useful for determining the largest extent needed for rendering or layout.
	Max(o ...Size[T]) Size[T]
	// MinDimension returns the minimum of width and height.
	// Used to determine the limiting dimension, e.g., for fitting a square inside a rectangle.
	MinDimension() T
	// MaxDimension returns the maximum of width and height.
	// Used to determine the dominant dimension, e.g., for scaling or aspect ratio calculations.
	MaxDimension() T
	// Area returns the area (width * height).
	// Fundamental in graphics for pixel count, memory allocation, or hit-testing.
	Area() T
	// Abs returns the size with absolute width and height.
	// Ensures dimensions are non-negative, which is important for rendering and layout.
	Abs() Size[T]
	// Floor returns the size with width and height floored.
	// Useful for aligning to pixel boundaries or integer grid systems.
	Floor() Size[T]
	// Ceil returns the size with width and height ceiled.
	// Useful for ensuring enough space is allocated, avoiding clipping.
	Ceil() Size[T]
	// Round returns the size with width and height rounded.
	// Used for snapping to the nearest pixel or grid unit.
	Round() Size[T]
	// IsZero returns true if both width and height are zero.
	// Used to detect degenerate or empty objects.
	IsZero() bool
	// IsFinite returns true if both width and height are finite.
	// Important for validating geometry before rendering or computation.
	IsFinite() bool
	// IsInfinite returns true if either width or height is infinite.
	// Can be used to represent unbounded or unconstrained objects.
	IsInfinite() bool
	// IsSquare returns true if width equals height.
	// Useful for aspect ratio checks, e.g., for icons or tiles.
	IsSquare() bool
	// MipCount returns the mipmap count for the size.
	// Used in texture mapping to determine the number of mipmap levels for an image.
	MipCount() int

	// String returns a string representation of the size. like "(width, height)".
	String() string
}

func NewSize[T Scalar](width, height T) Size[T] {
	return size[T]{width, height}
}

func NewSizeInfinite[T Scalar]() Size[T] {
	maxValue := Max[T]()
	return size[T]{width: maxValue, height: maxValue}
}

// size is a generic type that represents a size with width and height.
type size[T Scalar] struct {
	width, height T
}

// Width implements Size.
func (s size[T]) Width() T {
	return s.width
}

// Height implements Size.
func (s size[T]) Height() T {
	return s.height
}

// Add implements Size.
func (s size[T]) Add(other Size[T]) Size[T] {
	return size[T]{
		width:  s.width + other.Width(),
		height: s.height + other.Height(),
	}
}

// Sub implements Size.
func (s size[T]) Sub(other Size[T]) Size[T] {
	return size[T]{
		width:  s.width - other.Width(),
		height: s.height - other.Height(),
	}
}

// Mul implements Size.
func (s size[T]) Mul(other Size[T]) Size[T] {
	return size[T]{
		width:  s.width * other.Width(),
		height: s.height * other.Height(),
	}
}

// Div implements Size.
func (s size[T]) Div(other Size[T]) Size[T] {
	return size[T]{
		width:  s.width / other.Width(),
		height: s.height / other.Height(),
	}
}

// Neg implements Size.
func (s size[T]) Neg() Size[T] {
	return size[T]{
		width:  -s.width,
		height: -s.height,
	}
}

// Scale implements Size.
func (s size[T]) Scale(scale T) Size[T] {
	return size[T]{
		width:  s.width * scale,
		height: s.height * scale,
	}
}

// ScaleWH implements Size.
func (s size[T]) ScaleWH(width, height T) Size[T] {
	return size[T]{
		width:  s.width * width,
		height: s.height * height,
	}
}

// Equal implements Size.
func (s size[T]) Equal(other Size[T]) bool {
	return Equal(s.width, other.Width()) && Equal(s.height, other.Height())
}

// Min implements Size.
func (s size[T]) Min(o ...Size[T]) Size[T] {
	minW, minH := s.width, s.height
	for _, other := range o {
		if other.Width() < minW {
			minW = other.Width()
		}
		if other.Height() < minH {
			minH = other.Height()
		}
	}
	return size[T]{width: minW, height: minH}
}

// Max implements Size.
func (s size[T]) Max(o ...Size[T]) Size[T] {
	maxW, maxH := s.width, s.height
	for _, other := range o {
		if other.Width() > maxW {
			maxW = other.Width()
		}
		if other.Height() > maxH {
			maxH = other.Height()
		}
	}
	return size[T]{width: maxW, height: maxH}
}

// MinDimension implements Size.
func (s size[T]) MinDimension() T {
	if s.width < s.height {
		return s.width
	}
	return s.height
}

// MaxDimension implements Size.
func (s size[T]) MaxDimension() T {
	if s.width > s.height {
		return s.width
	}
	return s.height
}

// Area implements Size.
func (s size[T]) Area() T {
	return s.width * s.height
}

// Abs implements Size.
func (s size[T]) Abs() Size[T] {
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

// Floor implements Size.
func (s size[T]) Floor() Size[T] {
	return size[T]{
		width:  T(math.Floor(s.width.Float64())),
		height: T(math.Floor(s.height.Float64())),
	}
}

// Ceil implements Size.
func (s size[T]) Ceil() Size[T] {
	return size[T]{
		width:  T(math.Ceil(s.width.Float64())),
		height: T(math.Ceil(s.height.Float64())),
	}
}

// Round implements Size.
func (s size[T]) Round() Size[T] {
	return size[T]{
		width:  T(math.Round(s.width.Float64())),
		height: T(math.Round(s.height.Float64())),
	}
}

// IsZero implements Size.
func (s size[T]) IsZero() bool {
	var zero T
	return Equal(s.width, zero) && Equal(s.height, zero)
}

// IsFinite implements Size.
func (s size[T]) IsFinite() bool {
	return IsFinite(s.width) && IsFinite(s.height)
}

// IsInfinite implements Size.
func (s size[T]) IsInfinite() bool {
	return math.IsInf(s.width.Float64(), 0) || math.IsInf(s.height.Float64(), 0)
}

// IsSquare implements Size.
func (s size[T]) IsSquare() bool {
	return Equal(s.width, s.height)
}

// MipCount implements Size.
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

// String implements Size.
func (s size[T]) String() string {
	return "(" + s.width.String() + ", " + s.height.String() + ")"
}
