package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// SharedObjectVK manages shared Vulkan objects and resources
type SharedObjectVK struct {
	context     *ContextVK
	buffers     map[string]*BufferVK
	textures    map[string]*TextureVK
	descriptorSets map[string]vulkan.DescriptorSet
	samplers    map[string]*SamplerVK
	refCounts   map[string]int
}

// BufferVK represents a Vulkan buffer
type BufferVK struct {
	buffer       vulkan.Buffer
	deviceMemory vulkan.DeviceMemory
	size         vulkan.DeviceSize
	usage        vulkan.BufferUsageFlags
	properties   vulkan.MemoryPropertyFlags
}

// NewSharedObjectVK creates a new shared object manager
func NewSharedObjectVK(context *ContextVK) *SharedObjectVK {
	return &SharedObjectVK{
		context:        context,
		buffers:        make(map[string]*BufferVK),
		textures:       make(map[string]*TextureVK),
		descriptorSets: make(map[string]vulkan.DescriptorSet),
		samplers:       make(map[string]*SamplerVK),
		refCounts:      make(map[string]int),
	}
}

// GetBuffer retrieves or creates a shared buffer
func (so *SharedObjectVK) GetBuffer(name string, size vulkan.DeviceSize, usage vulkan.BufferUsageFlags, properties vulkan.MemoryPropertyFlags) (*BufferVK, error) {
	if buffer, exists := so.buffers[name]; exists {
		so.refCounts[name]++
		return buffer, nil
	}

	buffer, err := so.createBuffer(size, usage, properties)
	if err != nil {
		return nil, fmt.Errorf("failed to create buffer %s: %w", name, err)
	}

	so.buffers[name] = buffer
	so.refCounts[name] = 1

	return buffer, nil
}

// GetTexture retrieves or creates a shared texture
func (so *SharedObjectVK) GetTexture(name string, descriptor *TextureDescriptor) (*TextureVK, error) {
	if texture, exists := so.textures[name]; exists {
		so.refCounts[name]++
		return texture, nil
	}

	texture, err := NewTextureVK(so.context, descriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to create texture %s: %w", name, err)
	}

	so.textures[name] = texture
	so.refCounts[name] = 1

	return texture, nil
}

// GetSampler retrieves or creates a shared sampler
func (so *SharedObjectVK) GetSampler(name string, descriptor *SamplerDescriptor) (*SamplerVK, error) {
	if sampler, exists := so.samplers[name]; exists {
		so.refCounts[name]++
		return sampler, nil
	}

	sampler, err := NewSamplerVK(so.context, descriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to create sampler %s: %w", name, err)
	}

	so.samplers[name] = sampler
	so.refCounts[name] = 1

	return sampler, nil
}

// ReleaseBuffer releases a reference to a shared buffer
func (so *SharedObjectVK) ReleaseBuffer(name string) {
	if count, exists := so.refCounts[name]; exists {
		so.refCounts[name] = count - 1
		if so.refCounts[name] <= 0 {
			if buffer, exists := so.buffers[name]; exists {
				buffer.Destroy(so.context)
				delete(so.buffers, name)
			}
			delete(so.refCounts, name)
		}
	}
}

// ReleaseTexture releases a reference to a shared texture
func (so *SharedObjectVK) ReleaseTexture(name string) {
	if count, exists := so.refCounts[name]; exists {
		so.refCounts[name] = count - 1
		if so.refCounts[name] <= 0 {
			if texture, exists := so.textures[name]; exists {
				texture.Destroy()
				delete(so.textures, name)
			}
			delete(so.refCounts, name)
		}
	}
}

// ReleaseSampler releases a reference to a shared sampler
func (so *SharedObjectVK) ReleaseSampler(name string) {
	if count, exists := so.refCounts[name]; exists {
		so.refCounts[name] = count - 1
		if so.refCounts[name] <= 0 {
			if sampler, exists := so.samplers[name]; exists {
				sampler.Destroy()
				delete(so.samplers, name)
			}
			delete(so.refCounts, name)
		}
	}
}

// createBuffer creates a new Vulkan buffer
func (so *SharedObjectVK) createBuffer(size vulkan.DeviceSize, usage vulkan.BufferUsageFlags, properties vulkan.MemoryPropertyFlags) (*BufferVK, error) {
	bufferInfo := vulkan.BufferCreateInfo{
		SType:       vulkan.StructureTypeBufferCreateInfo,
		Size:        size,
		Usage:       usage,
		SharingMode: vulkan.SharingModeExclusive,
	}

	var buffer vulkan.Buffer
	if result := vulkan.CreateBuffer(so.context.device, &bufferInfo, nil, &buffer); result != vulkan.Success {
		return nil, fmt.Errorf("failed to create buffer: %s", result)
	}

	// Get memory requirements
	var memRequirements vulkan.MemoryRequirements
	vulkan.GetBufferMemoryRequirements(so.context.device, buffer, &memRequirements)
	memRequirements.Deref()

	// Find suitable memory type (simplified)
	memTypeIndex := uint32(0) // TODO: Implement proper memory type selection

	// Allocate memory
	allocInfo := vulkan.MemoryAllocateInfo{
		SType:           vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  memRequirements.Size,
		MemoryTypeIndex: memTypeIndex,
	}

	var deviceMemory vulkan.DeviceMemory
	if result := vulkan.AllocateMemory(so.context.device, &allocInfo, nil, &deviceMemory); result != vulkan.Success {
		vulkan.DestroyBuffer(so.context.device, buffer, nil)
		return nil, fmt.Errorf("failed to allocate memory: %s", result)
	}

	// Bind memory to buffer
	if result := vulkan.BindBufferMemory(so.context.device, buffer, deviceMemory, 0); result != vulkan.Success {
		vulkan.FreeMemory(so.context.device, deviceMemory, nil)
		vulkan.DestroyBuffer(so.context.device, buffer, nil)
		return nil, fmt.Errorf("failed to bind memory: %s", result)
	}

	return &BufferVK{
		buffer:       buffer,
		deviceMemory: deviceMemory,
		size:         size,
		usage:        usage,
		properties:   properties,
	}, nil
}

// GetBufferCount returns the number of shared buffers
func (so *SharedObjectVK) GetBufferCount() int {
	return len(so.buffers)
}

// GetTextureCount returns the number of shared textures
func (so *SharedObjectVK) GetTextureCount() int {
	return len(so.textures)
}

// GetSamplerCount returns the number of shared samplers
func (so *SharedObjectVK) GetSamplerCount() int {
	return len(so.samplers)
}

// Clear releases all shared objects
func (so *SharedObjectVK) Clear() {
	// Destroy all buffers
	for name, buffer := range so.buffers {
		buffer.Destroy(so.context)
		delete(so.buffers, name)
	}

	// Destroy all textures
	for name, texture := range so.textures {
		texture.Destroy()
		delete(so.textures, name)
	}

	// Destroy all samplers
	for name, sampler := range so.samplers {
		sampler.Destroy()
		delete(so.samplers, name)
	}

	// Clear reference counts
	for name := range so.refCounts {
		delete(so.refCounts, name)
	}
}

// Destroy destroys the shared object manager and all resources
func (so *SharedObjectVK) Destroy() {
	so.Clear()
}

// BufferVK methods

// GetBuffer returns the Vulkan buffer handle
func (b *BufferVK) GetBuffer() vulkan.Buffer {
	return b.buffer
}

// GetDeviceMemory returns the device memory handle
func (b *BufferVK) GetDeviceMemory() vulkan.DeviceMemory {
	return b.deviceMemory
}

// GetSize returns the buffer size
func (b *BufferVK) GetSize() vulkan.DeviceSize {
	return b.size
}

// GetUsage returns the buffer usage flags
func (b *BufferVK) GetUsage() vulkan.BufferUsageFlags {
	return b.usage
}

// Map maps the buffer memory for CPU access
func (b *BufferVK) Map(context *ContextVK, offset vulkan.DeviceSize, size vulkan.DeviceSize) (unsafe.Pointer, error) {
	var data unsafe.Pointer
	if result := vulkan.MapMemory(context.device, b.deviceMemory, offset, size, 0, &data); result != vulkan.Success {
		return nil, fmt.Errorf("failed to map memory: %s", result)
	}
	return data, nil
}

// Unmap unmaps the buffer memory
func (b *BufferVK) Unmap(context *ContextVK) {
	vulkan.UnmapMemory(context.device, b.deviceMemory)
}

// Destroy destroys the buffer and frees memory
func (b *BufferVK) Destroy(context *ContextVK) {
	if b.buffer != vulkan.NullBuffer {
		vulkan.DestroyBuffer(context.device, b.buffer, nil)
		b.buffer = vulkan.NullBuffer
	}

	if b.deviceMemory != vulkan.NullDeviceMemory {
		vulkan.FreeMemory(context.device, b.deviceMemory, nil)
		b.deviceMemory = vulkan.NullDeviceMemory
	}
}

// SamplerDescriptor describes sampler creation parameters
type SamplerDescriptor struct {
	MagFilter    vulkan.Filter
	MinFilter    vulkan.Filter
	AddressModeU vulkan.SamplerAddressMode
	AddressModeV vulkan.SamplerAddressMode
	AddressModeW vulkan.SamplerAddressMode
	Anisotropy   bool
	MaxAnisotropy float32
	CompareEnable bool
	CompareOp    vulkan.CompareOp
	MinLod       float32
	MaxLod       float32
}
