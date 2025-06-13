package vulkan

import (
	"fmt"
	"sync"
	"time"

	"github.com/opensraph/sraph/gpu/impl"
	. "github.com/opensraph/sraph/gpu/webgpu"
)

// VulkanAdapter implements the WebGPU Adapter interface using Vulkan
type VulkanAdapter struct {
	mu             sync.RWMutex
	backend        *VulkanBackend
	physicalDevice *VulkanPhysicalDevice
	adapterIndex   int
	baseAdapter    Adapter // Embed base adapter functionality
	destroyed      bool
}

// NewVulkanAdapter creates a new Vulkan adapter
func NewVulkanAdapter(backend *VulkanBackend, physicalDevice *VulkanPhysicalDevice, index int) *VulkanAdapter {
	// Create base adapter with Vulkan-specific info
	baseAdapter := impl.NewAdapter(BackendTypeVulkan, convertDeviceType(physicalDevice.Properties.DeviceType))

	return &VulkanAdapter{
		backend:        backend,
		physicalDevice: physicalDevice,
		adapterIndex:   index,
		baseAdapter:    baseAdapter,
	}
}

// GetFeatures retrieves supported features - override to add Vulkan-specific features
func (va *VulkanAdapter) GetFeatures(features SupportedFeatures) error {
	va.mu.RLock()
	defer va.mu.RUnlock()

	if va.destroyed {
		return fmt.Errorf("adapter has been destroyed")
	}

	// Get base features first
	if err := va.baseAdapter.GetFeatures(features); err != nil {
		return err
	}

	// Add Vulkan-specific features based on physical device capabilities
	vulkanFeatures := va.getVulkanSpecificFeatures()
	features.Features = append(features.Features, vulkanFeatures...)

	return nil
}

// GetInfo retrieves adapter information - override to provide Vulkan-specific info
func (va *VulkanAdapter) GetInfo(info AdapterInfo) (Status, error) {
	va.mu.RLock()
	defer va.mu.RUnlock()

	if va.destroyed {
		return StatusError, fmt.Errorf("adapter has been destroyed")
	}

	// Populate with Vulkan physical device info
	info.Vendor = va.physicalDevice.Properties.DeviceName
	info.Architecture = "Vulkan"
	info.Device = va.physicalDevice.Properties.DeviceName
	info.Description = fmt.Sprintf("Vulkan Device: %s", va.physicalDevice.Properties.DeviceName)
	info.BackendType = BackendTypeVulkan
	info.AdapterType = convertDeviceType(va.physicalDevice.Properties.DeviceType)
	info.VendorID = va.physicalDevice.Properties.VendorID
	info.DeviceID = va.physicalDevice.Properties.DeviceID

	return StatusSuccess, nil
}

// GetLimits retrieves adapter limits - override to provide Vulkan-specific limits
func (va *VulkanAdapter) GetLimits(limits Limits) (Status, error) {
	va.mu.RLock()
	defer va.mu.RUnlock()

	if va.destroyed {
		return StatusError, fmt.Errorf("adapter has been destroyed")
	}

	// Convert Vulkan limits to WebGPU limits
	limits = va.convertVulkanLimits()
	return StatusSuccess, nil
}

// HasFeature checks if a feature is supported - override for Vulkan-specific checks
func (va *VulkanAdapter) HasFeature(feature FeatureName) (bool, error) {
	va.mu.RLock()
	defer va.mu.RUnlock()

	if va.destroyed {
		return false, fmt.Errorf("adapter has been destroyed")
	}

	// Check Vulkan-specific features first
	if va.hasVulkanFeature(feature) {
		return true, nil
	}

	// Fall back to base adapter
	return va.baseAdapter.HasFeature(feature)
}

// RequestDevice requests a WebGPU device - override to create Vulkan device
func (va *VulkanAdapter) RequestDevice(descriptor DeviceDescriptor, callback RequestDeviceCallbackInfo) Future {
	va.mu.RLock()
	defer va.mu.RUnlock()

	future := Future{
		Id: impl.GenerateFutureId(), // Use impl package helper
	}

	if va.destroyed {
		impl.GlobalCallbackRegistry().RequestDevice(future.Id, callback, RequestDeviceStatusError, nil, "adapter has been destroyed")
		impl.GlobalCallbackManager().Complete(future.Id)
		return future
	}

	// Start async Vulkan device creation
	go func() {
		time.Sleep(time.Millisecond * 20) // Vulkan initialization takes longer

		device, err := va.backend.CreateDevice(va.physicalDevice, descriptor)
		if err != nil {
			impl.GlobalCallbackRegistry().RequestDevice(future.Id, callback, RequestDeviceStatusError, nil, err.Error())
		} else {
			// Wrap Vulkan device to implement WebGPU Device interface
			webgpuDevice := NewVulkanDeviceWrapper(device)
			impl.GlobalCallbackRegistry().RequestDevice(future.Id, callback, RequestDeviceStatusSuccess, webgpuDevice, "")
		}
		impl.GlobalCallbackManager().Complete(future.Id)
	}()

	return future
}

// Private helper methods for Vulkan-specific functionality

func (va *VulkanAdapter) getVulkanSpecificFeatures() []FeatureName {
	var features []FeatureName

	if va.physicalDevice.Features.GeometryShader {
		features = append(features, FeatureNameDepthClipControl)
	}
	if va.physicalDevice.Features.TextureCompressionBC {
		features = append(features, FeatureNameTextureCompressionBC)
	}
	if va.physicalDevice.Features.SamplerAnisotropy {
		features = append(features, FeatureNameFloat32Filterable)
	}
	if va.physicalDevice.Features.DepthClamp {
		features = append(features, FeatureNameDepth32FloatStencil8)
	}

	return features
}

func (va *VulkanAdapter) hasVulkanFeature(feature FeatureName) bool {
	switch feature {
	case FeatureNameDepthClipControl:
		return va.physicalDevice.Features.GeometryShader
	case FeatureNameTextureCompressionBC:
		return va.physicalDevice.Features.TextureCompressionBC
	case FeatureNameFloat32Filterable:
		return va.physicalDevice.Features.SamplerAnisotropy
	case FeatureNameDepth32FloatStencil8:
		return va.physicalDevice.Features.DepthClamp
	default:
		return false
	}
}

func (va *VulkanAdapter) convertVulkanLimits() Limits {
	// Convert Vulkan physical device limits to WebGPU limits
	return Limits{
		MaxTextureDimension1D:                     16384,
		MaxTextureDimension2D:                     va.physicalDevice.Properties.MaxTextureSize,
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
		MaxStorageBufferBindingSize:               uint64(va.physicalDevice.Properties.MaxBufferSize),
		MinUniformBufferOffsetAlignment:           256,
		MinStorageBufferOffsetAlignment:           256,
		MaxVertexBuffers:                          16,
		MaxBufferSize:                             va.physicalDevice.Properties.MaxBufferSize,
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

func convertDeviceType(vulkanType uint32) AdapterType {
	switch vulkanType {
	case 2: // VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
		return AdapterTypeDiscreteGPU
	case 1: // VK_PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
		return AdapterTypeIntegratedGPU
	case 4: // VK_PHYSICAL_DEVICE_TYPE_CPU
		return AdapterTypeCPU
	default:
		return AdapterTypeUnknown
	}
}
