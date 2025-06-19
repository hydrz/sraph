package vulkan

import (
	"sync"

	"github.com/opensraph/sraph/render"
	"github.com/vulkan-go/vulkan"
)

// SamplerLibraryVK manages Vulkan samplers
type SamplerLibraryVK struct {
	device   vulkan.Device
	mutex    sync.RWMutex
	samplers map[string]*SamplerVK
}

// NewSamplerLibraryVK creates a new Vulkan sampler library
func NewSamplerLibraryVK(device vulkan.Device) *SamplerLibraryVK {
	return &SamplerLibraryVK{
		device:   device,
		samplers: make(map[string]*SamplerVK),
	}
}

// GetSampler gets or creates a sampler with the given descriptor
func (s *SamplerLibraryVK) GetSampler(descriptor render.SamplerDescriptor) render.Sampler {
	key := s.generateSamplerKey(descriptor)

	s.mutex.RLock()
	if sampler, exists := s.samplers[key]; exists {
		s.mutex.RUnlock()
		return sampler
	}
	s.mutex.RUnlock()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Double-check after acquiring write lock
	if sampler, exists := s.samplers[key]; exists {
		return sampler
	}

	// Create new sampler
	sampler, err := s.createSampler(descriptor)
	if err != nil {
		return nil
	}

	s.samplers[key] = sampler
	return sampler
}

// CreateSampler creates a new Vulkan sampler
func (s *SamplerLibraryVK) CreateSampler(descriptor render.SamplerDescriptor) (*SamplerVK, error) {
	return s.createSampler(descriptor)
}

// GetSamplerCount returns the number of cached samplers
func (s *SamplerLibraryVK) GetSamplerCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.samplers)
}

// ClearCache clears all cached samplers
func (s *SamplerLibraryVK) ClearCache() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, sampler := range s.samplers {
		sampler.Cleanup()
	}

	s.samplers = make(map[string]*SamplerVK)
}

// Cleanup cleans up the sampler library
func (s *SamplerLibraryVK) Cleanup() {
	s.ClearCache()
}

// createSampler creates a new Vulkan sampler from descriptor
func (s *SamplerLibraryVK) createSampler(descriptor render.SamplerDescriptor) (*SamplerVK, error) {
	createInfo := vulkan.SamplerCreateInfo{
		SType:                   vulkan.StructureTypeSamplerCreateInfo,
		MagFilter:               s.convertFilterMode(descriptor.MagFilter),
		MinFilter:               s.convertFilterMode(descriptor.MinFilter),
		MipmapMode:              s.convertMipmapMode(descriptor.MipFilter),
		AddressModeU:            s.convertAddressMode(descriptor.AddressModeU),
		AddressModeV:            s.convertAddressMode(descriptor.AddressModeV),
		AddressModeW:            s.convertAddressMode(descriptor.AddressModeW),
		MipLodBias:              0.0,
		AnisotropyEnable:        vulkan.Bool32(0),
		MaxAnisotropy:           1.0,
		CompareEnable:           vulkan.Bool32(0),
		CompareOp:               vulkan.CompareOpAlways,
		MinLod:                  0.0,
		MaxLod:                  vulkan.LodClampNone,
		BorderColor:             vulkan.BorderColorIntOpaqueBlack,
		UnnormalizedCoordinates: vulkan.Bool32(0),
	}

	// Handle anisotropic filtering
	if descriptor.MaxAnisotropy > 1.0 {
		createInfo.AnisotropyEnable = vulkan.Bool32(1)
		createInfo.MaxAnisotropy = descriptor.MaxAnisotropy
	}

	var sampler vulkan.Sampler
	ret := vulkan.CreateSampler(s.device, &createInfo, nil, &sampler)
	if ret != vulkan.Success {
		return nil, render.NewError("Failed to create Vulkan sampler")
	}

	return NewSamplerVK(s.device, sampler, descriptor), nil
}

// generateSamplerKey generates a unique key for sampler caching
func (s *SamplerLibraryVK) generateSamplerKey(descriptor render.SamplerDescriptor) string {
	// TODO: Generate a proper hash key based on descriptor properties
	return "sampler_key"
}

// convertFilterMode converts render filter mode to Vulkan filter
func (s *SamplerLibraryVK) convertFilterMode(filter render.FilterMode) vulkan.Filter {
	switch filter {
	case render.FilterModeNearest:
		return vulkan.FilterNearest
	case render.FilterModeLinear:
		return vulkan.FilterLinear
	default:
		return vulkan.FilterLinear
	}
}

// convertMipmapMode converts render filter mode to Vulkan mipmap mode
func (s *SamplerLibraryVK) convertMipmapMode(filter render.FilterMode) vulkan.SamplerMipmapMode {
	switch filter {
	case render.FilterModeNearest:
		return vulkan.SamplerMipmapModeNearest
	case render.FilterModeLinear:
		return vulkan.SamplerMipmapModeLinear
	default:
		return vulkan.SamplerMipmapModeLinear
	}
}

// convertAddressMode converts render address mode to Vulkan address mode
func (s *SamplerLibraryVK) convertAddressMode(mode render.AddressMode) vulkan.SamplerAddressMode {
	switch mode {
	case render.AddressModeClampToEdge:
		return vulkan.SamplerAddressModeClampToEdge
	case render.AddressModeRepeat:
		return vulkan.SamplerAddressModeRepeat
	case render.AddressModeMirrorRepeat:
		return vulkan.SamplerAddressModeMirroredRepeat
	case render.AddressModeClampToBorder:
		return vulkan.SamplerAddressModeClampToBorder
	default:
		return vulkan.SamplerAddressModeClampToEdge
	}
}
