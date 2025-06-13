package impl

import (
	"fmt"
	"sync"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ Sampler = (*sampler)(nil)

// sampler implements the Sampler interface
type sampler struct {
	mu            sync.RWMutex
	label         string
	addressModeU  AddressMode
	addressModeV  AddressMode
	addressModeW  AddressMode
	magFilter     FilterMode
	minFilter     FilterMode
	mipmapFilter  MipmapFilterMode
	lodMinClamp   float32
	lodMaxClamp   float32
	compare       CompareFunction
	maxAnisotropy uint16
	destroyed     bool
}

// NewSampler creates a new WebGPU sampler (public factory function)
func NewSampler(descriptor SamplerDescriptor) Sampler {
	return &sampler{
		label:         descriptor.Label,
		addressModeU:  descriptor.AddressModeU,
		addressModeV:  descriptor.AddressModeV,
		addressModeW:  descriptor.AddressModeW,
		magFilter:     descriptor.MagFilter,
		minFilter:     descriptor.MinFilter,
		mipmapFilter:  descriptor.MipmapFilter,
		lodMinClamp:   descriptor.LodMinClamp,
		lodMaxClamp:   descriptor.LodMaxClamp,
		compare:       descriptor.Compare,
		maxAnisotropy: descriptor.MaxAnisotropy,
	}
}

// SetLabel sets the sampler label
func (s *sampler) SetLabel(label string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.destroyed {
		return fmt.Errorf("sampler has been destroyed")
	}

	s.label = label
	return nil
}
