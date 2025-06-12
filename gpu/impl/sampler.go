package impl

import (
	"fmt"
	"sync"
	"sync/atomic"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ Sampler = (*sampler)(nil)

// sampler implements the Sampler interface
type sampler struct {
	mu            sync.RWMutex
	refCount      int32
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

// newSampler creates a new WebGPU sampler
func newSampler(descriptor SamplerDescriptor) Sampler {
	return &sampler{
		refCount:      1,
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

// AddRef increments the reference count
func (s *sampler) AddRef() error {
	if atomic.LoadInt32(&s.refCount) <= 0 {
		return fmt.Errorf("sampler has been destroyed")
	}

	atomic.AddInt32(&s.refCount, 1)
	return nil
}

// Release decrements the reference count and destroys if zero
func (s *sampler) Release() error {
	newCount := atomic.AddInt32(&s.refCount, -1)
	if newCount == 0 {
		s.mu.Lock()
		s.destroyed = true
		s.mu.Unlock()
	} else if newCount < 0 {
		return fmt.Errorf("reference count cannot be negative")
	}
	return nil
}
