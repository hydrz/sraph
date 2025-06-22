package geom

import (
	"math"
)

// Size defines the interface for a two-dimensional size.
type Size interface {
	// Width returns the width component.
	Width() Scalar
	// Height returns the height component.
	Height() Scalar

	// Add returns the element-wise sum of this size and another.
	Add(o Size) Size
	// Sub returns the element-wise difference of this size and another.
	Sub(o Size) Size
	// Mul returns the element-wise product of this size and another.
	Mul(o Size) Size
	// Div returns the element-wise quotient of this size and another.
	Div(o Size) Size
	// Neg returns the negation of this size.
	Neg() Size
	// Equal reports whether this size and another are equal within floating-point tolerance.
	Equal(o Size) bool

	// Scale returns a size scaled by the given factor in both dimensions.
	Scale(scale Scalar) Size
	// ScaleWH returns a size scaled by width and height factors independently.
	ScaleWH(width, height Scalar) Size
	// MinDimension returns the smaller of width and height.
	MinDimension() Scalar
	// MaxDimension returns the larger of width and height.
	MaxDimension() Scalar
	// Area returns the area (width * height).
	Area() Scalar
	// Abs returns a size with non-negative dimensions.
	Abs() Size
	// Floor returns a size with dimensions rounded down.
	Floor() Size
	// Ceil returns a size with dimensions rounded up.
	Ceil() Size
	// Round returns a size with dimensions rounded to nearest integer.
	Round() Size

	// IsZero reports whether both dimensions are zero.
	IsZero() bool
	// IsFinite reports whether both dimensions are finite.
	IsFinite() bool
	// IsInfinite reports whether either dimension is infinite.
	IsInfinite() bool
	// IsSquare reports whether width equals height within tolerance.
	IsSquare() bool

	// MipCount returns the number of mipmap levels. Useful for texture operations.
	MipCount() int

	// String returns the string representation.
	String() string
}

// NewSize constructs a Size with the given dimensions.
func NewSize[T Number](width, height T) Size {
	return size[T]{
		width:  width,
		height: height,
	}
}

// size is the generic implementation of the Size interface.
type size[T Number] struct {
	width, height T
}

// Width implements Size.Width.
func (s size[T]) Width() Scalar {
	return ToScalar(s.width)
}

// Height implements Size.Height.
func (s size[T]) Height() Scalar {
	return ToScalar(s.height)
}

// Add implements Size.Add.
func (s size[T]) Add(o Size) Size {
	return NewSize(
		s.Width()+o.Width(),
		s.Height()+o.Height(),
	)
}

// Sub implements Size.Sub.
func (s size[T]) Sub(other Size) Size {
	return NewSize(
		s.Width()-other.Width(),
		s.Height()-other.Height(),
	)
}

// Mul implements Size.Mul.
func (s size[T]) Mul(other Size) Size {
	return NewSize(
		s.Width()*other.Width(),
		s.Height()*other.Height(),
	)
}

// Div implements Size.Div.
func (s size[T]) Div(other Size) Size {
	return NewSize(
		s.Width()/other.Width(),
		s.Height()/other.Height(),
	)
}

// Neg implements Size.Neg.
func (s size[T]) Neg() Size {
	return NewSize(
		-s.Width(),
		-s.Height(),
	)
}

// Equal implements Size.Equal.
func (s size[T]) Equal(o Size) bool {
	return NearlyEqual(s.Width(), o.Width()) &&
		NearlyEqual(s.Height(), o.Height())
}

// Scale implements Size.Scale.
func (s size[T]) Scale(scale Scalar) Size {
	return NewSize(
		s.Width()*scale,
		s.Height()*scale,
	)
}

// ScaleWH implements Size.ScaleWH.
func (s size[T]) ScaleWH(w Scalar, h Scalar) Size {
	return NewSize(
		s.Width()*w,
		s.Height()*h,
	)
}

// MinDimension implements Size.MinDimension.
func (s size[T]) MinDimension() Scalar {
	if s.width < s.height {
		return ToScalar(s.width)
	}
	return ToScalar(s.height)
}

// MaxDimension implements Size.MaxDimension.
func (s size[T]) MaxDimension() Scalar {
	if s.width > s.height {
		return ToScalar(s.width)
	}
	return ToScalar(s.height)
}

// Area implements Size.Area.
func (s size[T]) Area() Scalar {
	return ToScalar(s.width * s.height)
}

// Abs implements Size.Abs.
func (s size[T]) Abs() Size {
	return size[T]{
		width:  Abs(s.width),
		height: Abs(s.height),
	}
}

// Floor implements Size.Floor.
func (s size[T]) Floor() Size {
	return size[T]{
		width:  T(math.Floor(ToFloat64(s.width))),
		height: T(math.Floor(ToFloat64(s.height))),
	}
}

// Ceil implements Size.Ceil.
func (s size[T]) Ceil() Size {
	return size[T]{
		width:  T(math.Ceil(ToFloat64(s.width))),
		height: T(math.Ceil(ToFloat64(s.height))),
	}
}

// Round implements Size.Round.
func (s size[T]) Round() Size {
	return size[T]{
		width:  T(math.Round(ToFloat64(s.width))),
		height: T(math.Round(ToFloat64(s.height))),
	}
}

// IsZero implements Size.IsZero.
func (s size[T]) IsZero() bool {
	return NearlyEqual(s.width, 0) && NearlyEqual(s.height, 0)
}

// IsFinite implements Size.IsFinite.
func (s size[T]) IsFinite() bool {
	return IsFinite(s.width) && IsFinite(s.height)
}

// IsInfinite implements Size.IsInfinite.
func (s size[T]) IsInfinite() bool {
	return math.IsInf(ToFloat64(s.width), 0) || math.IsInf(ToFloat64(s.height), 0)
}

// IsSquare implements Size.IsSquare.
func (s size[T]) IsSquare() bool {
	return NearlyEqual(s.width, s.height)
}

// MipCount implements Size.MipCount.
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

// String implements Size.String.
func (s size[T]) String() string {
	return "(" + ToString(s.width) + ", " + ToString(s.height) + ")"
}
