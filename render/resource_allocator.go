package render

import (
	"fmt"
	"sync"

	"github.com/opensraph/sraph/gpu"
)

// resourceAllocator implements ResourceAllocator interface.
type resourceAllocator struct {
	device       gpu.Device
	bufferPool   sync.Pool
	texturePool  sync.Pool
	hostBuffers  map[uint64]Buffer
	buffersMutex sync.RWMutex
}

// NewResourceAllocator creates a new resource allocator.
func NewResourceAllocator(device gpu.Device) ResourceAllocator {
	return &resourceAllocator{
		device:      device,
		hostBuffers: make(map[uint64]Buffer),
	}
}

// CreateBuffer implements ResourceAllocator.
func (ra *resourceAllocator) CreateBuffer(desc BufferDescriptor) (Buffer, error) {
	gpuBuffer := ra.device.CreateBuffer(gpu.BufferDescriptor{
		Label: desc.Label,
		Size:  desc.Size,
		Usage: desc.Usage,
	})

	if gpuBuffer == nil {
		return nil, fmt.Errorf("failed to create GPU buffer")
	}

	return &bufferImpl{
		Buffer:     gpuBuffer,
		descriptor: desc,
	}, nil
}

// CreateTexture implements ResourceAllocator.
func (ra *resourceAllocator) CreateTexture(desc TextureDescriptor) (Texture, error) {
	gpuTexture := ra.device.CreateTexture(gpu.TextureDescriptor{
		Label:         desc.Label,
		Size:          desc.Size,
		Format:        desc.Format,
		Usage:         desc.Usage,
		SampleCount:   desc.SampleCount,
		MipLevelCount: desc.MipLevelCount,
		Dimension:     desc.Dimension,
	})

	if gpuTexture == nil {
		return nil, fmt.Errorf("failed to create GPU texture")
	}

	return &textureImpl{
		Texture:    gpuTexture,
		descriptor: desc,
	}, nil
}

// GetHostBuffer implements ResourceAllocator.
func (ra *resourceAllocator) GetHostBuffer(size uint64) (Buffer, error) {
	ra.buffersMutex.Lock()
	defer ra.buffersMutex.Unlock()

	if buffer, exists := ra.hostBuffers[size]; exists {
		return buffer, nil
	}

	desc := BufferDescriptor{
		Label: "Host Buffer",
		Size:  size,
		Usage: gpu.BufferUsageMapWrite | gpu.BufferUsageCopySrc,
	}

	buffer, err := ra.CreateBuffer(desc)
	if err != nil {
		return nil, err
	}

	ra.hostBuffers[size] = buffer
	return buffer, nil
}

// GetTransientBuffer implements ResourceAllocator.
func (ra *resourceAllocator) GetTransientBuffer(size uint64) (Buffer, error) {
	desc := BufferDescriptor{
		Label: "Transient Buffer",
		Size:  size,
		Usage: gpu.BufferUsageVertex | gpu.BufferUsageIndex | gpu.BufferUsageUniform,
	}

	return ra.CreateBuffer(desc)
}

// bufferImpl implements Buffer interface.
type bufferImpl struct {
	gpu.Buffer
	descriptor BufferDescriptor
}

// GetDescriptor implements Buffer.
func (b *bufferImpl) GetDescriptor() BufferDescriptor {
	return b.descriptor
}

// SetLabel implements Buffer.
func (b *bufferImpl) SetLabel(label string) {
	b.descriptor.Label = label
	b.Buffer.SetLabel(label)
}

// textureImpl implements Texture interface.
type textureImpl struct {
	gpu.Texture
	descriptor TextureDescriptor
}

// GetDescriptor implements Texture.
func (t *textureImpl) GetDescriptor() TextureDescriptor {
	return t.descriptor
}

// SetLabel implements Texture.
func (t *textureImpl) SetLabel(label string) {
	t.descriptor.Label = label
	t.Texture.SetLabel(label)
}
