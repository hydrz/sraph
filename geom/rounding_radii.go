package geom

// RoundingRadii defines the radii for the four corners of a rectangle.
// It is commonly used for rounded rectangle rendering and hit-testing in graphics.
type RoundingRadii[T Number] struct {
	TopLeft     Size[T]
	TopRight    Size[T]
	BottomLeft  Size[T]
	BottomRight Size[T]
}

// NewRoundingRadii creates a RoundingRadii with all four corners set to the same radius.
func NewRoundingRadii[T Number](radius T) RoundingRadii[T] {
	sz := Size[T]{radius, radius}
	return RoundingRadii[T]{
		TopLeft:     sz,
		TopRight:    sz,
		BottomLeft:  sz,
		BottomRight: sz,
	}
}

// NewRoundingRadiiLTRB creates a RoundingRadii with specified radii for each corner.
func NewRoundingRadiiLTRB[T Number](left, top, right, bottom T) RoundingRadii[T] {
	return RoundingRadii[T]{
		TopLeft:     Size[T]{left, top},
		TopRight:    Size[T]{right, top},
		BottomLeft:  Size[T]{left, bottom},
		BottomRight: Size[T]{right, bottom},
	}
}

// NewRoundingRadiiFromSizes creates a RoundingRadii with all four corners set to the same Size.
func NewRoundingRadiiFromSizes[T Number](sz Size[T]) RoundingRadii[T] {
	return RoundingRadii[T]{
		TopLeft:     sz,
		TopRight:    sz,
		BottomLeft:  sz,
		BottomRight: sz,
	}
}

// IsEmpty returns true if all radii are zero.
func (r RoundingRadii[T]) IsEmpty() bool {
	return r.TopLeft.IsZero() &&
		r.TopRight.IsZero() &&
		r.BottomLeft.IsZero() &&
		r.BottomRight.IsZero()
}

// IsFinite returns true if all radii are finite.
func (r RoundingRadii[T]) IsFinite() bool {
	return r.TopLeft.IsFinite() &&
		r.TopRight.IsFinite() &&
		r.BottomLeft.IsFinite() &&
		r.BottomRight.IsFinite()
}

// IsUniform returns true if all four corners are Eq.
func (r RoundingRadii[T]) IsUniform() bool {
	return r.TopLeft.Eq(r.TopRight) &&
		r.TopLeft.Eq(r.BottomLeft) &&
		r.TopLeft.Eq(r.BottomRight)
}

// Scale scales all radii by the given scalar.
func (r RoundingRadii[T]) Scale(scalar T) RoundingRadii[T] {
	return RoundingRadii[T]{
		TopLeft:     r.TopLeft.Scale(scalar),
		TopRight:    r.TopRight.Scale(scalar),
		BottomLeft:  r.BottomLeft.Scale(scalar),
		BottomRight: r.BottomRight.Scale(scalar),
	}
}

// ScaleToFit scales the radii so that the sum of the radii on each edge does not exceed the bounds.
func (r RoundingRadii[T]) ScaleToFit(bounds Rect[T]) RoundingRadii[T] {
	width := bounds.Width()
	height := bounds.Height()

	sumTop := r.TopLeft.Width + r.TopRight.Width
	sumBottom := r.BottomLeft.Width + r.BottomRight.Width
	sumLeft := r.TopLeft.Height + r.BottomLeft.Height
	sumRight := r.TopRight.Height + r.BottomRight.Height

	scaleX := T(1)
	scaleY := T(1)
	if sumTop > width {
		scaleX = width / sumTop
	}
	if sumBottom > width {
		s := width / sumBottom
		if s < scaleX {
			scaleX = s
		}
	}
	if sumLeft > height {
		scaleY = height / sumLeft
	}
	if sumRight > height {
		s := height / sumRight
		if s < scaleY {
			scaleY = s
		}
	}
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}
	if scale >= 1 {
		return r
	}
	return r.Scale(scale)
}

// Eq returns true if all four corners are Eq.
func (r RoundingRadii[T]) Eq(other RoundingRadii[T]) bool {
	return r.TopLeft.Eq(other.TopLeft) &&
		r.TopRight.Eq(other.TopRight) &&
		r.BottomLeft.Eq(other.BottomLeft) &&
		r.BottomRight.Eq(other.BottomRight)
}

// String returns a string representation of the rounding radii.
func (r RoundingRadii[T]) String() string {
	return "RoundingRadii{" +
		"TopLeft:" + r.TopLeft.String() +
		", TopRight:" + r.TopRight.String() +
		", BottomLeft:" + r.BottomLeft.String() +
		", BottomRight:" + r.BottomRight.String() +
		"}"
}
