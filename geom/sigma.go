package geom

import "strconv"

const kernelRadiusPerSigma = 1.73205080756887729352744634150587236694280525381038062805580697 // See https://oeis.org/A002194

// Sigma represents the standard deviation ("sigma") for Gaussian distributions in filter operations.
// It is measured in terms of the local space pixel grid of the filter input.
// Sigma determines how wide the Gaussian distribution stretches.
type Sigma[T Scalar] struct {
	sigma T
}

// ToRadians converts the Sigma value to a kernel radius in radians for convolution filters.
// For Gaussian blur kernels, the radius has a linear relationship with Sigma.
// Returns 0 if Sigma is not greater than 0.5.
func (s Sigma[T]) ToRadians() Radians[T] {
	if s.sigma.Float64() > 0.5 {
		return NewRadians[T](T((s.sigma.Float64() - 0.5) * kernelRadiusPerSigma))
	}
	return NewRadians[T](0.0)
}

// NewSigma creates a Sigma value from a given kernel radius in radians.
// If the radius is negative, returns 0.
func NewSigma[T Scalar](radius Radians[T]) Sigma[T] {
	if radius.Float64() < 0.0 {
		return Sigma[T]{sigma: T(0.0)}
	}

	return Sigma[T]{sigma: T(radius.Float64()/kernelRadiusPerSigma + 0.5)}
}

// Float64 implements Scalar for Sigma.
func (s Sigma[T]) Float64() float64 {
	return s.sigma.Float64()
}

// String returns a string representation of the Sigma value.
func (s Sigma[T]) String() string {
	return strconv.FormatFloat(s.sigma.Float64(), 'f', -1, 64) + "σ"
}
