package geom

// RSTransform describes a 2D transformation composed of rotation, scale, and translation.
// Commonly used for efficient sprite transformations and compatible with Flutter's RSTransform class.
// The transformation is represented as:
//
//	x' = ScaledCos * x - ScaledSin * y + TranslateX
//	y' = ScaledSin * x + ScaledCos * y + TranslateY
type RSTransform[T TScalar] struct {
	ScaledCos  T // Cosine of the rotation angle, multiplied by the uniform scale factor.
	ScaledSin  T // Sine of the rotation angle, multiplied by the uniform scale factor.
	TranslateX T // Translation offset along the X axis.
	TranslateY T // Translation offset along the Y axis.
}

// NewRSTransform constructs an RSTransform from the given origin, scale, and radians.
func NewRSTransform[T TScalar](origin Point[T], scale T, radians Radians) RSTransform[T] {
	cos, sin := NewMatrix[T]().CosSin(radians)
	return RSTransform[T]{
		ScaledCos:  cos * scale,
		ScaledSin:  sin * scale,
		TranslateX: origin.X,
		TranslateY: origin.Y,
	}
}

// TX returns the translation along the X axis.
func (r RSTransform[T]) TX() T {
	return r.TranslateX
}

// TY returns the translation along the Y axis.
func (r RSTransform[T]) TY() T {
	return r.TranslateY
}

// SX returns the scale factor along the X axis.
func (r RSTransform[T]) SX() T {
	var zero T
	if NearlyEqual(r.ScaledCos, zero) {
		return zero
	}
	return r.ScaledCos / Abs(r.ScaledCos)
}

// SY returns the scale factor along the Y axis.
func (r RSTransform[T]) SY() T {
	var zero T
	if NearlyEqual(r.ScaledCos, zero) {
		return zero
	}
	return r.ScaledSin / Abs(r.ScaledCos)
}

// IsAxisAligned returns true if the resulting transformed quad will be axis-aligned.
func (r RSTransform[T]) IsAxisAligned() bool {
	var zero T
	return NearlyEqual(r.ScaledCos, zero) || NearlyEqual(r.ScaledSin, zero)
}

// Matrix returns the 4x4 matrix representing this RSTransform.
func (r RSTransform[T]) Matrix() Matrix[T] {
	return Matrix[T]{
		r.ScaledCos, r.ScaledSin, 0, 0,
		-r.ScaledSin, r.ScaledCos, 0, 0,
		0, 0, 1, 0,
		r.TranslateX, r.TranslateY, 0, 1,
	}
}

// Quad returns the 4 corner points of the transformed quad for a sub-image of the given width and height.
// The order is UpperLeft, UpperRight, LowerLeft, LowerRight.
func (r RSTransform[T]) Quad(width, height T) Quad[T] {
	origin := Point[T]{r.TranslateX, r.TranslateY}
	dx := Point[T]{r.ScaledCos * width, r.ScaledSin * width}
	dy := Point[T]{-r.ScaledSin * height, r.ScaledCos * height}
	return Quad[T]{
		origin,
		origin.Add(dx),
		origin.Add(dy),
		origin.Add(dx).Add(dy),
	}
}

// QuadSize returns the 4 corner points for a sub-image of the given size.
func (r RSTransform[T]) QuadSize(size Size[T]) Quad[T] {
	return r.Quad(size.Width, size.Height)
}

// Bounds returns the bounding rectangle of the transformed quad for the given width and height.
func (r RSTransform[T]) Bounds(width, height T) Rect[T] {
	quad := r.Quad(width, height)
	return BoundingRect(quad[0], quad[1], quad[2], quad[3])
}

// BoundsSize returns the bounding rectangle for the given size.
func (r RSTransform[T]) BoundsSize(size Size[T]) Rect[T] {
	return r.Bounds(size.Width, size.Height)
}
