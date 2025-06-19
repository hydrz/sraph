package render

import (
	"fmt"
	"sync"

	"github.com/opensraph/sraph/gpu"
)

// samplerLibrary implements SamplerLibrary interface.
type samplerLibrary struct {
	device     gpu.Device
	samplers   map[string]gpu.Sampler
	cacheMutex sync.RWMutex
}

// NewSamplerLibrary creates a new sampler library.
func NewSamplerLibrary(device gpu.Device) SamplerLibrary {
	return &samplerLibrary{
		device:   device,
		samplers: make(map[string]gpu.Sampler),
	}
}

// GetSampler implements SamplerLibrary.
func (sl *samplerLibrary) GetSampler(desc SamplerDescriptor) (gpu.Sampler, error) {
	key := generateSamplerKey(desc)

	sl.cacheMutex.RLock()
	if sampler, exists := sl.samplers[key]; exists {
		sl.cacheMutex.RUnlock()
		return sampler, nil
	}
	sl.cacheMutex.RUnlock()

	// Create new sampler if not cached
	return sl.CreateSampler(desc)
}

// CreateSampler implements SamplerLibrary.
func (sl *samplerLibrary) CreateSampler(desc SamplerDescriptor) (gpu.Sampler, error) {
	gpuDesc := gpu.SamplerDescriptor{
		Label:         desc.Label,
		AddressModeU:  desc.AddressModeU,
		AddressModeV:  desc.AddressModeV,
		AddressModeW:  desc.AddressModeW,
		MagFilter:     desc.MagFilter,
		MinFilter:     desc.MinFilter,
		MipmapFilter:  desc.MipmapFilter,
		LodMinClamp:   desc.LodMinClamp,
		LodMaxClamp:   desc.LodMaxClamp,
		Compare:       desc.Compare,
		MaxAnisotropy: desc.MaxAnisotropy,
	}

	sampler := sl.device.CreateSampler(gpuDesc)
	if sampler == nil {
		return nil, fmt.Errorf("failed to create sampler")
	}

	// Cache the sampler
	key := generateSamplerKey(desc)
	sl.cacheMutex.Lock()
	sl.samplers[key] = sampler
	sl.cacheMutex.Unlock()

	return sampler, nil
}

// generateSamplerKey creates a unique key for sampler caching.
func generateSamplerKey(desc SamplerDescriptor) string {
	// TODO: Implement proper hash key generation based on descriptor content
	return fmt.Sprintf("sampler_%s_%d_%d_%d_%d_%d_%d_%f_%f_%d_%d",
		desc.Label,
		desc.AddressModeU,
		desc.AddressModeV,
		desc.AddressModeW,
		desc.MagFilter,
		desc.MinFilter,
		desc.MipmapFilter,
		desc.LodMinClamp,
		desc.LodMaxClamp,
		desc.Compare,
		desc.MaxAnisotropy,
	)
}
