package geom

// RoundingRadii defines the radii for the four corners of a rectangle.
type RoundingRadii struct {
	topLeft     Size
	topRight    Size
	bottomLeft  Size
	bottomRight Size
}

// NewRoundingRadii creates a RoundingRadii with all four corners set to the same radius.
func NewRoundingRadii(radius Scalar) RoundingRadii {
	sz := NewSize(radius, radius)
	return RoundingRadii{
		topLeft:     sz,
		topRight:    sz,
		bottomLeft:  sz,
		bottomRight: sz,
	}
}

// NewRoundingRadiiLTRB creates a RoundingRadii with specified radii for each corner.
func NewRoundingRadiiLTRB(left, top, right, bottom Scalar) RoundingRadii {
	return RoundingRadii{
		topLeft:     NewSize(left, top),
		topRight:    NewSize(right, top),
		bottomLeft:  NewSize(left, bottom),
		bottomRight: NewSize(right, bottom),
	}
}

// NewRoundingRadiiFromSizes creates a RoundingRadii with all four corners set to the same Size.
func NewRoundingRadiiFromSizes(sz Size) RoundingRadii {
	return RoundingRadii{
		topLeft:     sz,
		topRight:    sz,
		bottomLeft:  sz,
		bottomRight: sz,
	}
}

// TopLeft returns the top-left corner size.
func (r RoundingRadii) TopLeft() Size { return r.topLeft }

// TopRight returns the top-right corner size.
func (r RoundingRadii) TopRight() Size { return r.topRight }

// BottomLeft returns the bottom-left corner size.
func (r RoundingRadii) BottomLeft() Size { return r.bottomLeft }

// BottomRight returns the bottom-right corner size.
func (r RoundingRadii) BottomRight() Size { return r.bottomRight }

// IsEmpty reports whether all radii are zero.
func (r RoundingRadii) IsEmpty() bool {
	return r.topLeft.IsZero() &&
		r.topRight.IsZero() &&
		r.bottomLeft.IsZero() &&
		r.bottomRight.IsZero()
}

// IsFinite reports whether all radii are finite.
func (r RoundingRadii) IsFinite() bool {
	return r.topLeft.IsFinite() &&
		r.topRight.IsFinite() &&
		r.bottomLeft.IsFinite() &&
		r.bottomRight.IsFinite()
}

// IsUniform reports whether all four corners are equal.
func (r RoundingRadii) IsUniform() bool {
	return r.topLeft.Equal(r.topRight) &&
		r.topLeft.Equal(r.bottomLeft) &&
		r.topLeft.Equal(r.bottomRight)
}

// Scale scales all radii by the given scalar.
func (r RoundingRadii) Scale(scalar Scalar) RoundingRadii {
	return RoundingRadii{
		topLeft:     r.topLeft.Scale(scalar),
		topRight:    r.topRight.Scale(scalar),
		bottomLeft:  r.bottomLeft.Scale(scalar),
		bottomRight: r.bottomRight.Scale(scalar),
	}
}

// ScaleToFit scales the radii so that the sum of the radii on each edge does not exceed the bounds.
func (r RoundingRadii) ScaleToFit(bounds Rect) RoundingRadii {
	width := bounds.Width()
	height := bounds.Height()

	sumTop := r.topLeft.Width() + r.topRight.Width()
	sumBottom := r.bottomLeft.Width() + r.bottomRight.Width()
	sumLeft := r.topLeft.Height() + r.bottomLeft.Height()
	sumRight := r.topRight.Height() + r.bottomRight.Height()

	scaleX := Scalar(1)
	scaleY := Scalar(1)
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

// Equal reports whether all four corners are equal.
func (r RoundingRadii) Equal(other RoundingRadii) bool {
	return r.topLeft.Equal(other.TopLeft()) &&
		r.topRight.Equal(other.TopRight()) &&
		r.bottomLeft.Equal(other.BottomLeft()) &&
		r.bottomRight.Equal(other.BottomRight())
}

// String returns a string representation of the rounding radii.
func (r RoundingRadii) String() string {
	return "(" + r.topLeft.String() + ", " +
		r.topRight.String() + ", " +
		r.bottomLeft.String() + ", " +
		r.bottomRight.String() + ")"
}
