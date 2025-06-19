package display

// SamplingOptions defines options for sampling textures and images.
type SamplingOptions struct {
	filterMode    FilterMode
	imageSampling ImageSampling
	useCubic      bool
	cubicB        float32
	cubicC        float32
}

// NewSamplingOptionsNearest creates sampling options for nearest neighbor filtering.
func NewSamplingOptionsNearest() SamplingOptions {
	return SamplingOptions{
		filterMode:    FilterModeNearest,
		imageSampling: ImageSamplingNearestNeighbor,
		useCubic:      false,
	}
}

// NewSamplingOptionsLinear creates sampling options for linear filtering.
func NewSamplingOptionsLinear() SamplingOptions {
	return SamplingOptions{
		filterMode:    FilterModeLinear,
		imageSampling: ImageSamplingLinear,
		useCubic:      false,
	}
}

// NewSamplingOptionsMipmap creates sampling options for mipmap linear filtering.
func NewSamplingOptionsMipmap() SamplingOptions {
	return SamplingOptions{
		filterMode:    FilterModeLinear,
		imageSampling: ImageSamplingMipmapLinear,
		useCubic:      false,
	}
}

// NewSamplingOptionsCubic creates sampling options for cubic filtering.
func NewSamplingOptionsCubic(b, c float32) SamplingOptions {
	return SamplingOptions{
		filterMode:    FilterModeLinear,
		imageSampling: ImageSamplingCubic,
		useCubic:      true,
		cubicB:        b,
		cubicC:        c,
	}
}

// NewSamplingOptionsCatmullRom creates sampling options for Catmull-Rom cubic filtering.
func NewSamplingOptionsCatmullRom() SamplingOptions {
	return NewSamplingOptionsCubic(0.0, 0.5)
}

// NewSamplingOptionsMitchell creates sampling options for Mitchell cubic filtering.
func NewSamplingOptionsMitchell() SamplingOptions {
	return NewSamplingOptionsCubic(1.0/3.0, 1.0/3.0)
}

// NewSamplingOptionsAniso creates sampling options for anisotropic filtering.
func NewSamplingOptionsAniso(maxAnisotropy float32) SamplingOptions {
	return SamplingOptions{
		filterMode:    FilterModeLinear,
		imageSampling: ImageSamplingMipmapLinear,
		useCubic:      false,
		// TODO: Add anisotropy support
	}
}

// FilterMode returns the filter mode.
func (s SamplingOptions) FilterMode() FilterMode {
	return s.filterMode
}

// ImageSampling returns the image sampling mode.
func (s SamplingOptions) ImageSampling() ImageSampling {
	return s.imageSampling
}

// UseCubic returns true if cubic filtering is enabled.
func (s SamplingOptions) UseCubic() bool {
	return s.useCubic
}

// CubicB returns the cubic B parameter.
func (s SamplingOptions) CubicB() float32 {
	return s.cubicB
}

// CubicC returns the cubic C parameter.
func (s SamplingOptions) CubicC() float32 {
	return s.cubicC
}

// IsEqual returns true if two sampling options are equal.
func (s SamplingOptions) IsEqual(other SamplingOptions) bool {
	return s.filterMode == other.filterMode &&
		s.imageSampling == other.imageSampling &&
		s.useCubic == other.useCubic &&
		s.cubicB == other.cubicB &&
		s.cubicC == other.cubicC
}

// Predefined sampling options
var (
	// SamplingOptionsNearest represents nearest neighbor sampling.
	SamplingOptionsNearest = NewSamplingOptionsNearest()
	// SamplingOptionsLinear represents linear sampling.
	SamplingOptionsLinear = NewSamplingOptionsLinear()
	// SamplingOptionsMipmap represents mipmap linear sampling.
	SamplingOptionsMipmap = NewSamplingOptionsMipmap()
	// SamplingOptionsCatmullRom represents Catmull-Rom cubic sampling.
	SamplingOptionsCatmullRom = NewSamplingOptionsCatmullRom()
	// SamplingOptionsMitchell represents Mitchell cubic sampling.
	SamplingOptionsMitchell = NewSamplingOptionsMitchell()
)
