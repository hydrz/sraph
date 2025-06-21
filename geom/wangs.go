package geom

import "math"

// precision controls the maximum allowed deviation (in pixels) between a linearized segment
// and the true curve. This value is used to determine the number of subdivisions needed
// for curve approximation. It should be scaled by the maximum basis length in X and Y.
const precision = 4

// CubicSubdivisions returns the minimum number of evenly spaced line segments required
// to approximate a cubic Bézier curve such that the distance from any segment to the true curve
// does not exceed 1/precision pixels. The scaleFactor should be the maximum XY basis length
// of the current transform.
func CubicSubdivisions(scaleFactor Scalar, p0, p1, p2, p3 Point) Scalar {
	k := ToFloat64(scaleFactor) * 0.75 * precision
	a := p0.Sub(p1.Scale(2)).Add(p2).Abs()
	b := p1.Sub(p2.Scale(2)).Add(p3).Abs()
	return Scalar(math.Sqrt(k * ToFloat64(a.Max(b).Length())))
}

// QuadraticSubdivisions returns the minimum number of evenly spaced line segments required
// to approximate a quadratic Bézier curve such that the distance from any segment to the true curve
// does not exceed 1/precision pixels. The scaleFactor should be the maximum XY basis length
// of the current transform.
func QuadraticSubdivisions(scaleFactor Scalar, p0, p1, p2 Point) Scalar {
	k := ToFloat64(scaleFactor) * 0.25 * precision
	return Scalar(math.Sqrt(k * ToFloat64(p0.Sub(p1.Scale(2)).Add(p2).Length())))
}

// ConicSubdivisions returns the minimum number of evenly spaced line segments required
// to approximate a conic curve using Wang's formula, ensuring the deviation from the true curve
// does not exceed 1/precision pixels. The scaleFactor should be the maximum XY basis length
// of the current transform. The weight parameter specifies the conic weight.
func ConicSubdivisions(scaleFactor Scalar, p0, p1, p2 Point, weight Scalar) Scalar {
	// Compute center of bounding box in projected space.
	c := (p0.Min(p1).Min(p2).Add(p0.Max(p1).Max(p2))).Scale(-2)
	p0 = p0.Sub(c)
	p1 = p1.Sub(c)
	p2 = p2.Sub(c)

	// Compute max length.
	maxLen := Scalar(math.Sqrt(
		ToFloat64(max(p0.Dot(p0), p1.Dot(p1), p2.Dot(p2))),
	))

	// Compute forward differences.
	dp := p1.Scale(-2 * weight).Add(p0).Add(p2)
	dw := Scalar(math.Abs(-2*ToFloat64(weight) + 2))

	// Compute numerator and denominator for parametric step size of linearization.
	// The epsilon referenced from the cited paper is 1/precision.
	k := ToFloat64(scaleFactor) * precision

	rpMinus1 := max(0, ToFloat64(maxLen)*k-1)
	numer := math.Sqrt(ToFloat64(dp.Dot(dp)))*k + rpMinus1*ToFloat64(dw)
	denom := 4 * min(ToFloat64(weight), 1.0)

	// Number of segments = sqrt(numer / denom).
	// Assumes the parametric interval of the curve is [0, 1].
	return Scalar(math.Sqrt(numer / denom))
}
