package geom

// RoundingRadii defines the radii for the four corners of a rectangle.
// It is commonly used for rounded rectangle rendering and hit-testing in graphics.
type RoundingRadii[T Scalar] interface {
	// TopLeft returns the radii for the top-left corner.
	TopLeft() Size[T]
	// TopRight returns the radii for the top-right corner.
	TopRight() Size[T]
	// BottomLeft returns the radii for the bottom-left corner.
	BottomLeft() Size[T]
	// BottomRight returns the radii for the bottom-right corner.
	BottomRight() Size[T]

	// IsEmpty returns true if all corner radii are zero.
	IsEmpty() bool
	// IsFinite returns true if all corner radii are finite.
	IsFinite() bool
	// IsUniform returns true if all four corners have the same radii.
	IsUniform() bool

	// Scale scales all corner radii by the given factor.
	Scale(scale T) RoundingRadii[T]

	// ScaleToFit returns a scaled RoundingRadii.
	// So that the sum of two radii on each edge does not exceed the rect's width or height.
	ScaleToFit(bounds Rect[T]) RoundingRadii[T]

	// Equal compares two RoundingRadii for equality.
	Equal(other RoundingRadii[T]) bool

	// String returns a string representation of the rounding radii.
	String() string
}

// NewRoundingRadii creates a RoundingRadii with all four corners set to the same radius.
func NewRoundingRadii[T Scalar](radius T) RoundingRadii[T] {
	return roundingRadii[T]{
		leftTop:     NewSize(radius, radius),
		rightTop:    NewSize(radius, radius),
		leftBottom:  NewSize(radius, radius),
		rightBottom: NewSize(radius, radius),
	}
}

// NewRoundingRadiiLTRB creates a RoundingRadii with each corner specified individually.
func NewRoundingRadiiLTRB[T Scalar](left, top, right, bottom T) RoundingRadii[T] {
	return roundingRadii[T]{
		leftTop:     NewSize(left, top),
		rightTop:    NewSize(right, top),
		leftBottom:  NewSize(left, bottom),
		rightBottom: NewSize(right, bottom),
	}
}

// NewRoundingRadiiFromSizes creates a RoundingRadii with all four corners set to the same Size.
func NewRoundingRadiiFromSizes[T Scalar](radii Size[T]) RoundingRadii[T] {
	return roundingRadii[T]{
		leftTop:     radii,
		rightTop:    radii,
		leftBottom:  radii,
		rightBottom: radii,
	}
}

type roundingRadii[T Scalar] struct {
	leftTop     Size[T]
	rightTop    Size[T]
	leftBottom  Size[T]
	rightBottom Size[T]
}

// TopLeft implements RoundingRadii.
func (r roundingRadii[T]) TopLeft() Size[T] {
	return r.leftTop
}

// TopRight implements RoundingRadii.
func (r roundingRadii[T]) TopRight() Size[T] {
	return r.rightTop
}

// BottomLeft implements RoundingRadii.
func (r roundingRadii[T]) BottomLeft() Size[T] {
	return r.leftBottom
}

// BottomRight implements RoundingRadii.
func (r roundingRadii[T]) BottomRight() Size[T] {
	return r.rightBottom
}

// IsEmpty implements RoundingRadii.
func (r roundingRadii[T]) IsEmpty() bool {
	return r.leftTop.Equal(NewSize[T](0, 0)) &&
		r.rightTop.Equal(NewSize[T](0, 0)) &&
		r.leftBottom.Equal(NewSize[T](0, 0)) &&
		r.rightBottom.Equal(NewSize[T](0, 0))
}

// IsFinite implements RoundingRadii.
func (r roundingRadii[T]) IsFinite() bool {
	return r.leftTop.IsFinite() &&
		r.rightTop.IsFinite() &&
		r.leftBottom.IsFinite() &&
		r.rightBottom.IsFinite()
}

// IsUniform implements RoundingRadii.
func (r roundingRadii[T]) IsUniform() bool {
	return r.leftTop.Equal(r.rightTop) &&
		r.leftTop.Equal(r.leftBottom) &&
		r.leftTop.Equal(r.rightBottom)
}

// Scale implements RoundingRadii.
func (r roundingRadii[T]) Scale(scale T) RoundingRadii[T] {
	return roundingRadii[T]{
		leftTop:     r.leftTop.Scale(scale),
		rightTop:    r.rightTop.Scale(scale),
		leftBottom:  r.leftBottom.Scale(scale),
		rightBottom: r.rightBottom.Scale(scale),
	}
}

// ScaleToFit implements RoundingRadii.
func (r roundingRadii[T]) ScaleToFit(bounds Rect[T]) RoundingRadii[T] {
	width := bounds.Width()
	height := bounds.Height()

	// For each edge, the sum of the two adjacent corner radii must not exceed the edge length.
	// Compute scale factors for each edge.
	scaleX := T(1)
	scaleY := T(1)

	sumTop := r.leftTop.Width() + r.rightTop.Width()
	sumBottom := r.leftBottom.Width() + r.rightBottom.Width()
	sumLeft := r.leftTop.Height() + r.leftBottom.Height()
	sumRight := r.rightTop.Height() + r.rightBottom.Height()

	if sumTop > width && sumTop != 0 {
		scaleX = width / sumTop
	}
	if sumBottom > width && sumBottom != 0 {
		s := width / sumBottom
		if s < scaleX {
			scaleX = s
		}
	}
	if sumLeft > height && sumLeft != 0 {
		scaleY = height / sumLeft
	}
	if sumRight > height && sumRight != 0 {
		s := height / sumRight
		if s < scaleY {
			scaleY = s
		}
	}
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}
	if scale > 1 {
		scale = 1
	}
	return r.Scale(scale)
}

// Equal implements RoundingRadii.
func (r roundingRadii[T]) Equal(other RoundingRadii[T]) bool {
	return r.leftTop.Equal(other.TopLeft()) &&
		r.rightTop.Equal(other.TopRight()) &&
		r.leftBottom.Equal(other.BottomLeft()) &&
		r.rightBottom.Equal(other.BottomRight())
}

// String implements RoundingRadii.
func (r roundingRadii[T]) String() string {
	return "(" +
		"ul: " + r.leftTop.String() + ", " +
		"ur: " + r.rightTop.String() + ", " +
		"ll: " + r.leftBottom.String() + ", " +
		"lr: " + r.rightBottom.String() +
		")"
}
