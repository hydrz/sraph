package geom

// RSTransform represents a 2D transformation consisting of rotation, uniform scale, and translation.
type RSTransform struct {
	ScaledCos  Scalar // Cosine of the rotation angle, multiplied by the uniform scale factor.
	ScaledSin  Scalar // Sine of the rotation angle, multiplied by the uniform scale factor.
	TranslateX Scalar // Translation offset along the X axis.
	TranslateY Scalar // Translation offset along the Y axis.
}

// NewRSTransform returns an RSTransform with the given translation origin, uniform scale, and rotation (in radians).
func NewRSTransform(origin Point, scale Scalar, radians Radians) RSTransform {
	cos, sin := NewMatrix().CosSin(radians)
	return RSTransform{
		ScaledCos:  cos * scale,
		ScaledSin:  sin * scale,
		TranslateX: origin.X(),
		TranslateY: origin.Y(),
	}
}

// TX returns the translation offset along the X axis.
func (r RSTransform) TX() Scalar {
	return r.TranslateX
}

// TY returns the translation offset along the Y axis.
func (r RSTransform) TY() Scalar {
	return r.TranslateY
}

// SX returns the scale factor along the X axis.
// If ScaledCos is zero, returns zero.
func (r RSTransform) SX() Scalar {
	if NearlyEqual(r.ScaledCos, 0) {
		return 0
	}
	return r.ScaledCos / Abs(r.ScaledCos)
}

// SY returns the scale factor along the Y axis.
// If ScaledCos is zero, returns zero.
func (r RSTransform) SY() Scalar {
	if NearlyEqual(r.ScaledCos, 0) {
		return 0
	}
	return r.ScaledSin / Abs(r.ScaledCos)
}

// IsAxisAligned reports whether the transformed quad will be axis-aligned.
func (r RSTransform) IsAxisAligned() bool {
	return NearlyEqual(r.ScaledCos, 0) || NearlyEqual(r.ScaledSin, 0)
}

// Matrix returns a 4x4 matrix representation of the RSTransform.
func (r RSTransform) Matrix() Matrix {
	return Matrix{
		r.ScaledCos, r.ScaledSin, 0, 0,
		-r.ScaledSin, r.ScaledCos, 0, 0,
		0, 0, 1, 0,
		r.TranslateX, r.TranslateY, 0, 1,
	}
}

// Quad returns the four corners of the transformed rectangle with the given width and height.
// The order is UpperLeft, UpperRight, LowerLeft, LowerRight.
func (r RSTransform) Quad(width, height Scalar) Quad {
	origin := NewPoint(r.TranslateX, r.TranslateY)
	dx := NewPoint(r.ScaledCos*width, r.ScaledSin*width)
	dy := NewPoint(-r.ScaledSin*height, r.ScaledCos*height)
	return Quad{
		origin,
		origin.Add(dx),
		origin.Add(dy),
		origin.Add(dx).Add(dy),
	}
}

// QuadSize returns the four corners of the transformed rectangle with the given size.
func (r RSTransform) QuadSize(size Size) Quad {
	return r.Quad(size.Width(), size.Height())
}

// Bounds returns the bounding rectangle of the transformed quad for the given width and height.
func (r RSTransform) Bounds(width, height Scalar) Rect {
	quad := r.Quad(width, height)
	return BoundingRect(quad[0], quad[1], quad[2], quad[3])
}

// BoundsSize returns the bounding rectangle of the transformed quad for the given size.
func (r RSTransform) BoundsSize(size Size) Rect {
	return r.Bounds(size.Width(), size.Height())
}
