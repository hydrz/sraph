package vulkan

import (
	"github.com/opensraph/sraph/render"
	"github.com/vulkan-go/vulkan"
)

// SamplerVK represents a Vulkan sampler
type SamplerVK struct {
	device     vulkan.Device
	sampler    vulkan.Sampler
	descriptor render.SamplerDescriptor
}

// NewSamplerVK creates a new Vulkan sampler wrapper
func NewSamplerVK(device vulkan.Device, sampler vulkan.Sampler, descriptor render.SamplerDescriptor) *SamplerVK {
	return &SamplerVK{
		device:     device,
		sampler:    sampler,
		descriptor: descriptor,
	}
}

// GetDescriptor returns the sampler descriptor
func (s *SamplerVK) GetDescriptor() render.SamplerDescriptor {
	return s.descriptor
}

// GetSampler returns the underlying Vulkan sampler
func (s *SamplerVK) GetSampler() vulkan.Sampler {
	return s.sampler
}

// GetMagFilter returns the magnification filter
func (s *SamplerVK) GetMagFilter() render.FilterMode {
	return s.descriptor.MagFilter
}

// GetMinFilter returns the minification filter
func (s *SamplerVK) GetMinFilter() render.FilterMode {
	return s.descriptor.MinFilter
}

// GetMipFilter returns the mipmap filter
func (s *SamplerVK) GetMipFilter() render.FilterMode {
	return s.descriptor.MipFilter
}

// GetAddressModeU returns the U address mode
func (s *SamplerVK) GetAddressModeU() render.AddressMode {
	return s.descriptor.AddressModeU
}

// GetAddressModeV returns the V address mode
func (s *SamplerVK) GetAddressModeV() render.AddressMode {
	return s.descriptor.AddressModeV
}

// GetAddressModeW returns the W address mode
func (s *SamplerVK) GetAddressModeW() render.AddressMode {
	return s.descriptor.AddressModeW
}

// GetMaxAnisotropy returns the maximum anisotropy level
func (s *SamplerVK) GetMaxAnisotropy() float32 {
	return s.descriptor.MaxAnisotropy
}

// IsAnisotropyEnabled returns whether anisotropic filtering is enabled
func (s *SamplerVK) IsAnisotropyEnabled() bool {
	return s.descriptor.MaxAnisotropy > 1.0
}

// Cleanup cleans up the sampler resources
func (s *SamplerVK) Cleanup() {
	if s.sampler != vulkan.NullHandle {
		vulkan.DestroySampler(s.device, s.sampler, nil)
		s.sampler = vulkan.NullHandle
	}
}

// IsValid returns whether the sampler is valid
func (s *SamplerVK) IsValid() bool {
	return s.sampler != vulkan.NullHandle
}

// GetHashCode returns a hash code for the sampler
func (s *SamplerVK) GetHashCode() uint64 {
	// TODO: Generate proper hash based on descriptor
	return 0
}

// Equals compares this sampler with another
func (s *SamplerVK) Equals(other *SamplerVK) bool {
	if other == nil {
		return false
	}

	// Compare descriptors
	return s.descriptor.MagFilter == other.descriptor.MagFilter &&
		s.descriptor.MinFilter == other.descriptor.MinFilter &&
		s.descriptor.MipFilter == other.descriptor.MipFilter &&
		s.descriptor.AddressModeU == other.descriptor.AddressModeU &&
		s.descriptor.AddressModeV == other.descriptor.AddressModeV &&
		s.descriptor.AddressModeW == other.descriptor.AddressModeW &&
		s.descriptor.MaxAnisotropy == other.descriptor.MaxAnisotropy
}

// Clone creates a copy of the sampler (shares underlying Vulkan sampler)
func (s *SamplerVK) Clone() *SamplerVK {
	return &SamplerVK{
		device:     s.device,
		sampler:    s.sampler, // Shared reference
		descriptor: s.descriptor,
	}
}
