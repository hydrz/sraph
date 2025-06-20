package geom

import "strconv"

const kernelRadiusPerSigma = 1.73205080756887729352744634150587236694280525381038062805580697 // See https://oeis.org/A002194

// Sigma represents the standard deviation ("sigma") for Gaussian distributions in filter operations.
// Sigma is measured in the pixel grid of the filter input and determines the spread of the Gaussian.
type Sigma struct {
	sigma Scalar
}

// ToRadians returns the kernel radius in radians corresponding to the Sigma value for convolution filters.
// For Gaussian blur, the radius is linearly related to Sigma. Returns 0 if Sigma is not greater than 0.5.
func (s Sigma) ToRadians() Radians {
	if s.sigma > 0.5 {
		return Radians((s.sigma - 0.5) * kernelRadiusPerSigma)
	}
	return Radians(0.0)
}

// NewSigma returns a Sigma value corresponding to the given kernel radius in radians.
// If the radius is negative, the result is zero.
func NewSigma(radius Radians) Sigma {
	if radius < 0.0 {
		return Sigma{sigma: 0.0}
	}
	return Sigma{sigma: Scalar(radius)/kernelRadiusPerSigma + 0.5}
}

// String returns a string representation of the Sigma value.
func (s Sigma) String() string {
	return strconv.FormatFloat(s.sigma, 'f', -1, 64) + "σ"
}
