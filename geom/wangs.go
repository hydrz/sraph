package geom

import "math"

// Precision controls how closely linearized segments must approximate the true curve.
// Don't allow linearized segments to be off by more than 1/4th of a pixel from
// the true curve. This value should be scaled by the max basis of the X and Y directions.
const precision = 4

// CubicSubdivisions returns the minimum number of evenly spaced (in parametric sense)
// line segments that the cubic curve must be subdivided into to guarantee all segments
// stay within a distance of "1/precision" pixels from the true curve.
//
// The scaleFactor should be the max basis XY of the current transform.
func CubicSubdivisions[T Scalar](scaleFactor T, p0, p1, p2, p3 Point[T]) T {
	k := scaleFactor.Float64() * 0.75 * precision
	a := p0.Sub(p1.Scale(2)).Add(p2).Abs()
	b := p1.Sub(p2.Scale(2)).Add(p3).Abs()
	return T(math.Sqrt(k * a.Max(b).Length().Float64()))
}

// QuadraticSubdivisions returns the minimum number of evenly spaced (in parametric sense)
// line segments that the quadratic curve must be subdivided into to guarantee all segments
// stay within a distance of "1/precision" pixels from the true curve.
//
// The scaleFactor should be the max basis XY of the current transform.
func QuadraticSubdivisions[T Scalar](scaleFactor T, p0, p1, p2 Point[T]) T {
	k := scaleFactor.Float64() * 0.25 * precision
	return T(math.Sqrt(k * p0.Sub(p1.Scale(2)).Add(p2).Length().Float64()))
}

// ConicSubdivisions returns Wang's formula specialized for a conic curve.
func ConicSubdivisions[T Scalar](scaleFactor T, p0, p1, p2 Point[T], weight T) T {
	// Compute center of bounding box in projected space
	c := (p0.Min(p1).Min(p2).Add(p0.Max(p1).Max(p2))).Scale(1 / 2)
	p0 = p0.Sub(c)
	p1 = p1.Sub(c)
	p2 = p2.Sub(c)

	// Compute max length
	maxLen := T(math.Sqrt(max(p0.Dot(p0).Float64(), p1.Dot(p1).Float64(), p2.Dot(p2).Float64())))

	// Compute forward differences
	dp := p1.Scale(-2 * weight).Add(p0).Add(p2)
	dw := T(math.Abs(-2*weight.Float64() + 2))

	// Compute numerator and denominator for parametric step size of
	// linearization. Here, the epsilon referenced from the cited paper
	// is 1/precision.
	k := scaleFactor.Float64() * precision

	rpMinus1 := max(0, maxLen.Float64()*k-1)
	numer := math.Sqrt(dp.Dot(dp).Float64())*k + rpMinus1*dw.Float64()
	denom := 4 * min(weight.Float64(), 1.0)

	// Number of segments = sqrt(numer / denom).
	// This assumes parametric interval of curve being linearized is
	//   [t0,t1] = [0, 1].
	// If not, the number of segments is (tmax - tmin) / sqrt(denom / numer).

	return T(math.Sqrt(numer / denom))
}
