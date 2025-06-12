package vulkan

import (
	"fmt"
	"sync"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

// VulkanBackend implements WebGPU backend using Vulkan
type VulkanBackend struct {
	mu        sync.RWMutex
	instance  *VulkanInstance
	devices   map[*VulkanPhysicalDevice]*VulkanDevice
	enabled   bool
	destroyed bool
}

// VulkanBackendConfig contains backend configuration
type VulkanBackendConfig struct {
	EnableDebug        bool
	EnableValidation   bool
	RequiredExtensions []string
	ApplicationName    string
	ApplicationVersion uint32
	EngineName         string
	EngineVersion      uint32
}

var (
	globalVulkanBackend *VulkanBackend
	vulkanBackendOnce   sync.Once
)

// GetVulkanBackend returns the global Vulkan backend instance
func GetVulkanBackend() *VulkanBackend {
	vulkanBackendOnce.Do(func() {
		globalVulkanBackend = &VulkanBackend{
			devices: make(map[*VulkanPhysicalDevice]*VulkanDevice),
		}
	})
	return globalVulkanBackend
}

// Initialize initializes the Vulkan backend with real Vulkan library loading
func (vb *VulkanBackend) Initialize(config VulkanBackendConfig) error {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	if vb.enabled {
		return fmt.Errorf("Vulkan backend is already initialized")
	}

	if vb.destroyed {
		return fmt.Errorf("Vulkan backend has been destroyed")
	}

	// Load Vulkan library
	if err := LoadVulkanLibrary(); err != nil {
		return fmt.Errorf("failed to load Vulkan library: %v", err)
	}

	// Create Vulkan instance
	instance, err := NewVulkanInstance(config.EnableDebug)
	if err != nil {
		UnloadVulkanLibrary()
		return fmt.Errorf("failed to create Vulkan instance: %v", err)
	}

	vb.instance = instance
	vb.enabled = true

	return nil
}

// Shutdown shuts down the Vulkan backend and unloads library
func (vb *VulkanBackend) Shutdown() error {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	if vb.destroyed {
		return nil
	}

	// Destroy all devices first
	for _, device := range vb.devices {
		if device != nil {
			device.Destroy()
		}
	}
	vb.devices = make(map[*VulkanPhysicalDevice]*VulkanDevice)

	// Destroy instance
	if vb.instance != nil {
		vb.instance.Destroy()
		vb.instance = nil
	}

	// Unload Vulkan library
	if err := UnloadVulkanLibrary(); err != nil {
		// Log warning but don't fail shutdown
		fmt.Printf("Warning: failed to unload Vulkan library: %v\n", err)
	}

	vb.enabled = false
	vb.destroyed = true

	return nil
}

// IsEnabled checks if the backend is enabled
func (vb *VulkanBackend) IsEnabled() bool {
	vb.mu.RLock()
	defer vb.mu.RUnlock()
	return vb.enabled && !vb.destroyed
}

// GetInstance returns the Vulkan instance
func (vb *VulkanBackend) GetInstance() *VulkanInstance {
	vb.mu.RLock()
	defer vb.mu.RUnlock()
	return vb.instance
}

// EnumerateAdapters enumerates available adapters
func (vb *VulkanBackend) EnumerateAdapters() ([]Adapter, error) {
	vb.mu.RLock()
	defer vb.mu.RUnlock()

	if !vb.enabled {
		return nil, fmt.Errorf("Vulkan backend is not initialized")
	}

	if vb.destroyed {
		return nil, fmt.Errorf("Vulkan backend has been destroyed")
	}

	if vb.instance == nil {
		return nil, fmt.Errorf("Vulkan instance is nil")
	}

	physicalDevices := vb.instance.GetPhysicalDevices()
	if len(physicalDevices) == 0 {
		return nil, fmt.Errorf("no Vulkan physical devices found")
	}

	var adapters []Adapter
	for i := range physicalDevices {
		adapter := NewVulkanAdapter(vb, &physicalDevices[i], i)
		adapters = append(adapters, adapter)
	}

	return adapters, nil
}

// CreateDevice creates a Vulkan device
func (vb *VulkanBackend) CreateDevice(physicalDevice *VulkanPhysicalDevice, descriptor DeviceDescriptor) (*VulkanDevice, error) {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	if !vb.enabled {
		return nil, fmt.Errorf("Vulkan backend is not initialized")
	}

	if vb.destroyed {
		return nil, fmt.Errorf("Vulkan backend has been destroyed")
	}

	if physicalDevice == nil {
		return nil, fmt.Errorf("physical device cannot be nil")
	}

	// Check if device already exists
	if existingDevice, exists := vb.devices[physicalDevice]; exists && existingDevice != nil {
		return existingDevice, nil
	}

	// Convert WebGPU features to Vulkan features
	vulkanFeatures := VulkanDeviceFeatures{
		GeometryShader:       vb.hasFeature(descriptor.RequiredFeatures, FeatureNameDepthClipControl),
		TessellationShader:   false, // Not directly mapped
		MultiViewport:        true,  // Usually available
		SamplerAnisotropy:    true,  // Usually available
		TextureCompressionBC: vb.hasFeature(descriptor.RequiredFeatures, FeatureNameTextureCompressionBC),
		DepthClamp:           vb.hasFeature(descriptor.RequiredFeatures, FeatureNameDepthClipControl),
	}

	// Create device
	createInfo := VulkanDeviceCreateInfo{
		PhysicalDevice:     physicalDevice,
		Instance:           vb.instance,
		RequiredFeatures:   vulkanFeatures,
		RequiredExtensions: []string{"VK_KHR_swapchain"},
	}

	device, err := NewVulkanDevice(createInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vulkan device: %v", err)
	}

	vb.devices[physicalDevice] = device
	return device, nil
}

// hasFeature checks if a WebGPU feature is requested
func (vb *VulkanBackend) hasFeature(features []FeatureName, feature FeatureName) bool {
	for _, f := range features {
		if f == feature {
			return true
		}
	}
	return false
}

// GetDevice returns an existing device for a physical device
func (vb *VulkanBackend) GetDevice(physicalDevice *VulkanPhysicalDevice) *VulkanDevice {
	vb.mu.RLock()
	defer vb.mu.RUnlock()

	if physicalDevice == nil {
		return nil
	}

	return vb.devices[physicalDevice]
}

// RemoveDevice removes a device from tracking
func (vb *VulkanBackend) RemoveDevice(physicalDevice *VulkanPhysicalDevice) {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	if physicalDevice != nil {
		delete(vb.devices, physicalDevice)
	}
}

// GetDeviceCount returns the number of created devices
func (vb *VulkanBackend) GetDeviceCount() int {
	vb.mu.RLock()
	defer vb.mu.RUnlock()
	return len(vb.devices)
}

// GetSupportedFeatures returns supported WebGPU features
func (vb *VulkanBackend) GetSupportedFeatures() []FeatureName {
	return []FeatureName{
		FeatureNameDepthClipControl,
		FeatureNameDepth32FloatStencil8,
		FeatureNameTimestampQuery,
		FeatureNameTextureCompressionBC,
		FeatureNameTextureCompressionBCSliced3D,
		FeatureNameIndirectFirstInstance,
		FeatureNameRG11B10UfloatRenderable,
		FeatureNameBGRA8UnormStorage,
		FeatureNameFloat32Filterable,
		FeatureNameFloat32Blendable,
	}
}

// GetSupportedLimits returns supported WebGPU limits
func (vb *VulkanBackend) GetSupportedLimits() Limits {
	// In a real implementation, this would query Vulkan limits
	return Limits{
		MaxTextureDimension1D:                     16384,
		MaxTextureDimension2D:                     16384,
		MaxTextureDimension3D:                     2048,
		MaxTextureArrayLayers:                     2048,
		MaxBindGroups:                             8,
		MaxBindGroupsPlusVertexBuffers:            32,
		MaxBindingsPerBindGroup:                   1000,
		MaxDynamicUniformBuffersPerPipelineLayout: 8,
		MaxDynamicStorageBuffersPerPipelineLayout: 4,
		MaxSampledTexturesPerShaderStage:          32,
		MaxSamplersPerShaderStage:                 16,
		MaxStorageBuffersPerShaderStage:           8,
		MaxStorageTexturesPerShaderStage:          8,
		MaxUniformBuffersPerShaderStage:           15,
		MaxUniformBufferBindingSize:               65536,
		MaxStorageBufferBindingSize:               1 << 30, // 1GB
		MinUniformBufferOffsetAlignment:           256,
		MinStorageBufferOffsetAlignment:           256,
		MaxVertexBuffers:                          16,
		MaxBufferSize:                             1 << 30, // 1GB
		MaxVertexAttributes:                       32,
		MaxVertexBufferArrayStride:                2048,
		MaxInterStageShaderVariables:              32,
		MaxColorAttachments:                       8,
		MaxColorAttachmentBytesPerSample:          32,
		MaxComputeWorkgroupStorageSize:            32768,
		MaxComputeInvocationsPerWorkgroup:         1024,
		MaxComputeWorkgroupSizeX:                  1024,
		MaxComputeWorkgroupSizeY:                  1024,
		MaxComputeWorkgroupSizeZ:                  64,
		MaxComputeWorkgroupsPerDimension:          65535,
		MaxImmediateSize:                          16777216,
	}
}

// IsDestroyed checks if the backend is destroyed
func (vb *VulkanBackend) IsDestroyed() bool {
	vb.mu.RLock()
	defer vb.mu.RUnlock()
	return vb.destroyed
}

// CreateSurface creates a Vulkan surface for the given platform data
func (vb *VulkanBackend) CreateSurface(platformData interface{}) (*VulkanSurface, error) {
	vb.mu.RLock()
	defer vb.mu.RUnlock()

	if !vb.enabled {
		return nil, fmt.Errorf("Vulkan backend is not initialized")
	}

	if vb.destroyed {
		return nil, fmt.Errorf("Vulkan backend has been destroyed")
	}

	if vb.instance == nil {
		return nil, fmt.Errorf("Vulkan instance is nil")
	}

	// In a real implementation, platformData would contain window handle
	return NewVulkanSurface(vb.instance, nil)
}
