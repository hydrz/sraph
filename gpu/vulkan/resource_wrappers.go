package vulkan

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/opensraph/sraph/gpu/impl"
	. "github.com/opensraph/sraph/gpu/wgpu"
)

// VulkanBufferWrapper wraps Vulkan buffer to implement WebGPU Buffer interface
type VulkanBufferWrapper struct {
	vulkanBuffer interface{} // Should be *vulkan.VulkanBuffer
	baseBuffer   Buffer
}

// NewVulkanBufferWrapper creates a wrapper around Vulkan buffer
func NewVulkanBufferWrapper(vulkanBuffer interface{}, descriptor BufferDescriptor) Buffer {
	baseBuffer := impl.NewBuffer(descriptor)

	return &VulkanBufferWrapper{
		vulkanBuffer: vulkanBuffer,
		baseBuffer:   baseBuffer,
	}
}

// Delegate Buffer interface methods to base buffer with Vulkan-specific handling
func (vbw *VulkanBufferWrapper) Destroy() error {
	// First destroy Vulkan buffer if it has a Destroy method
	if destroyer, ok := vbw.vulkanBuffer.(interface{ Destroy() error }); ok {
		if err := destroyer.Destroy(); err != nil {
			return fmt.Errorf("failed to destroy Vulkan buffer: %v", err)
		}
	}

	// Then destroy base buffer
	return vbw.baseBuffer.Destroy()
}

func (vbw *VulkanBufferWrapper) GetConstMappedRange(offset uintptr, size uintptr) (unsafe.Pointer, error) {
	return vbw.baseBuffer.GetConstMappedRange(offset, size)
}

func (vbw *VulkanBufferWrapper) GetMappedRange(offset uintptr, size uintptr) (unsafe.Pointer, error) {
	return vbw.baseBuffer.GetMappedRange(offset, size)
}

func (vbw *VulkanBufferWrapper) GetMapState() (BufferMapState, error) {
	return vbw.baseBuffer.GetMapState()
}

func (vbw *VulkanBufferWrapper) GetSize() (uint64, error) {
	return vbw.baseBuffer.GetSize()
}

func (vbw *VulkanBufferWrapper) GetUsage() (BufferUsage, error) {
	return vbw.baseBuffer.GetUsage()
}

func (vbw *VulkanBufferWrapper) MapAsync(mode MapMode, offset uintptr, size uintptr, callback BufferMapCallbackInfo) Future {
	return vbw.baseBuffer.MapAsync(mode, offset, size, callback)
}

func (vbw *VulkanBufferWrapper) ReadMappedRange(offset uintptr, data unsafe.Pointer, size uintptr) (Status, error) {
	return vbw.baseBuffer.ReadMappedRange(offset, data, size)
}

func (vbw *VulkanBufferWrapper) SetLabel(label string) error {
	return vbw.baseBuffer.SetLabel(label)
}

func (vbw *VulkanBufferWrapper) Unmap() error {
	return vbw.baseBuffer.Unmap()
}

func (vbw *VulkanBufferWrapper) WriteMappedRange(offset uintptr, data unsafe.Pointer, size uintptr) (Status, error) {
	return vbw.baseBuffer.WriteMappedRange(offset, data, size)
}

// VulkanTextureWrapper wraps Vulkan texture to implement WebGPU Texture interface
type VulkanTextureWrapper struct {
	vulkanTexture interface{} // Should be *vulkan.VulkanTexture
	baseTexture   Texture
}

// NewVulkanTextureWrapper creates a wrapper around Vulkan texture
func NewVulkanTextureWrapper(vulkanTexture interface{}, descriptor TextureDescriptor) Texture {
	baseTexture := impl.NewTexture(descriptor)

	return &VulkanTextureWrapper{
		vulkanTexture: vulkanTexture,
		baseTexture:   baseTexture,
	}
}

// Delegate Texture interface methods to base texture
func (vtw *VulkanTextureWrapper) CreateView(descriptor TextureViewDescriptor) (TextureView, error) {
	return vtw.baseTexture.CreateView(descriptor)
}

func (vtw *VulkanTextureWrapper) Destroy() error {
	// First destroy Vulkan texture if it has a Destroy method
	if destroyer, ok := vtw.vulkanTexture.(interface{ Destroy() error }); ok {
		if err := destroyer.Destroy(); err != nil {
			return fmt.Errorf("failed to destroy Vulkan texture: %v", err)
		}
	}

	return vtw.baseTexture.Destroy()
}

func (vtw *VulkanTextureWrapper) GetDepthOrArrayLayers() (uint32, error) {
	return vtw.baseTexture.GetDepthOrArrayLayers()
}

func (vtw *VulkanTextureWrapper) GetDimension() (TextureDimension, error) {
	return vtw.baseTexture.GetDimension()
}

func (vtw *VulkanTextureWrapper) GetFormat() (TextureFormat, error) {
	return vtw.baseTexture.GetFormat()
}

func (vtw *VulkanTextureWrapper) GetHeight() (uint32, error) {
	return vtw.baseTexture.GetHeight()
}

func (vtw *VulkanTextureWrapper) GetMipLevelCount() (uint32, error) {
	return vtw.baseTexture.GetMipLevelCount()
}

func (vtw *VulkanTextureWrapper) GetSampleCount() (uint32, error) {
	return vtw.baseTexture.GetSampleCount()
}

func (vtw *VulkanTextureWrapper) GetUsage() (TextureUsage, error) {
	return vtw.baseTexture.GetUsage()
}

func (vtw *VulkanTextureWrapper) GetWidth() (uint32, error) {
	return vtw.baseTexture.GetWidth()
}

func (vtw *VulkanTextureWrapper) SetLabel(label string) error {
	return vtw.baseTexture.SetLabel(label)
}

// VulkanSamplerWrapper wraps Vulkan sampler to implement WebGPU Sampler interface
type VulkanSamplerWrapper struct {
	vulkanSampler interface{} // Should be *vulkan.VulkanSampler
	baseSampler   Sampler
}

// NewVulkanSamplerWrapper creates a wrapper around Vulkan sampler
func NewVulkanSamplerWrapper(vulkanSampler interface{}, descriptor SamplerDescriptor) Sampler {
	baseSampler := impl.NewSampler(descriptor)

	return &VulkanSamplerWrapper{
		vulkanSampler: vulkanSampler,
		baseSampler:   baseSampler,
	}
}

// Delegate Sampler interface methods to base sampler
func (vsw *VulkanSamplerWrapper) SetLabel(label string) error {
	return vsw.baseSampler.SetLabel(label)
}

// ResourceWrapper provides a generic interface for all resource wrappers
type ResourceWrapper interface {
	GetBackendType() BackendType
	GetNativeHandle() interface{}
	IsDestroyed() bool
}

// BackendResourceManager manages backend-specific resource wrappers
type BackendResourceManager struct {
	mu        sync.RWMutex
	resources map[uintptr]ResourceWrapper
	nextID    uintptr
}

var globalResourceManager = &BackendResourceManager{
	resources: make(map[uintptr]ResourceWrapper),
	nextID:    1,
}

// RegisterResource registers a resource wrapper for tracking
func RegisterResource(wrapper ResourceWrapper) uintptr {
	globalResourceManager.mu.Lock()
	defer globalResourceManager.mu.Unlock()

	id := globalResourceManager.nextID
	globalResourceManager.nextID++
	globalResourceManager.resources[id] = wrapper

	return id
}

// UnregisterResource removes a resource wrapper from tracking
func UnregisterResource(id uintptr) {
	globalResourceManager.mu.Lock()
	defer globalResourceManager.mu.Unlock()

	delete(globalResourceManager.resources, id)
}

// GetResourceCount returns the number of tracked resources
func GetResourceCount() int {
	globalResourceManager.mu.RLock()
	defer globalResourceManager.mu.RUnlock()

	return len(globalResourceManager.resources)
}

// GetResourcesByBackend returns resources filtered by backend type
func GetResourcesByBackend(backendType BackendType) []ResourceWrapper {
	globalResourceManager.mu.RLock()
	defer globalResourceManager.mu.RUnlock()

	var result []ResourceWrapper
	for _, wrapper := range globalResourceManager.resources {
		if wrapper.GetBackendType() == backendType {
			result = append(result, wrapper)
		}
	}

	return result
}

// CleanupDestroyedResources removes destroyed resources from tracking
func CleanupDestroyedResources() int {
	globalResourceManager.mu.Lock()
	defer globalResourceManager.mu.Unlock()

	cleaned := 0
	for id, wrapper := range globalResourceManager.resources {
		if wrapper.IsDestroyed() {
			delete(globalResourceManager.resources, id)
			cleaned++
		}
	}

	return cleaned
}

// GetResourceStats returns statistics about tracked resources
func GetResourceStats() map[BackendType]int {
	globalResourceManager.mu.RLock()
	defer globalResourceManager.mu.RUnlock()

	stats := make(map[BackendType]int)
	for _, wrapper := range globalResourceManager.resources {
		stats[wrapper.GetBackendType()]++
	}

	return stats
}
