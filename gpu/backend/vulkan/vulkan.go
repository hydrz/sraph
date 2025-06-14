package vulkan

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/opensraph/sraph/gpu/backend"
	"github.com/opensraph/sraph/gpu/wgpu"
)

// VulkanBackend implements the wgpu.Adapter and backend.Backend interfaces.
type VulkanBackend struct {
	instance         VkInstance
	physicalDevice   VkPhysicalDevice
	properties       VkPhysicalDeviceProperties
	features         VkPhysicalDeviceFeatures
	memoryProperties VkPhysicalDeviceMemoryProperties
	queueFamilies    []VkQueueFamilyProperties
	initialized      bool
}

// Ensure VulkanBackend implements required interfaces.
var _ backend.Backend = (*VulkanBackend)(nil)
var _ wgpu.Adapter = (*VulkanBackend)(nil)

func init() {
	runtime.LockOSThread()
	vulkanBackend := &VulkanBackend{}
	backend.RegisterBackend(vulkanBackend)
}

// Features implements wgpu.Adapter.
func (vb *VulkanBackend) Features() (*wgpu.SupportedFeatures, error) {
	if !vb.initialized {
		if err := vb.Init(); err != nil {
			return nil, fmt.Errorf("failed to initialize Vulkan backend: %v", err)
		}
	}

	var features []wgpu.FeatureName

	// Map Vulkan features to WebGPU features
	if vb.features.DepthClamp == VK_TRUE {
		features = append(features, wgpu.FeatureNameDepthClipControl)
	}
	if vb.features.PipelineStatisticsQuery == VK_TRUE {
		features = append(features, wgpu.FeatureNameTimestampQuery)
	}
	if vb.features.TextureCompressionBC == VK_TRUE {
		features = append(features, wgpu.FeatureNameTextureCompressionBC)
	}
	if vb.features.TextureCompressionETC2 == VK_TRUE {
		features = append(features, wgpu.FeatureNameTextureCompressionETC2)
	}
	if vb.features.SamplerAnisotropy == VK_TRUE {
		features = append(features, wgpu.FeatureNameSamplerAnisotropy)
	}

	return &wgpu.SupportedFeatures{
		Features: features,
	}, nil
}

// HasFeature implements wgpu.Adapter.
func (vb *VulkanBackend) HasFeature(feature wgpu.FeatureName) (bool, error) {
	features, err := vb.Features()
	if err != nil {
		return false, err
	}

	for _, f := range features.Features {
		if f == feature {
			return true, nil
		}
	}
	return false, nil
}

// Info implements wgpu.Adapter.
func (vb *VulkanBackend) Info() (*wgpu.AdapterInfo, error) {
	if !vb.initialized {
		if err := vb.Init(); err != nil {
			return nil, fmt.Errorf("failed to initialize Vulkan backend: %v", err)
		}
	}

	// Convert device name from C string
	deviceName := bytePtrToString(&vb.properties.DeviceName[0], 256)

	// Map Vulkan device type to WebGPU adapter type
	var adapterType wgpu.AdapterType
	switch vb.properties.DeviceType {
	case VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU:
		adapterType = wgpu.AdapterTypeDiscreteGPU
	case VK_PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU:
		adapterType = wgpu.AdapterTypeIntegratedGPU
	case VK_PHYSICAL_DEVICE_TYPE_CPU:
		adapterType = wgpu.AdapterTypeCPU
	default:
		adapterType = wgpu.AdapterTypeUnknown
	}

	return &wgpu.AdapterInfo{
		Vendor:       getVendorName(vb.properties.VendorID),
		Architecture: "Unknown", // Vulkan doesn't provide architecture info directly
		Device:       deviceName,
		Description:  fmt.Sprintf("Vulkan %s", deviceName),
		BackendType:  wgpu.BackendTypeVulkan,
		AdapterType:  adapterType,
		VendorID:     vb.properties.VendorID,
		DeviceID:     vb.properties.DeviceID,
	}, nil
}

// Limits implements wgpu.Adapter.
func (vb *VulkanBackend) Limits() (*wgpu.Limits, error) {
	if !vb.initialized {
		if err := vb.Init(); err != nil {
			return nil, fmt.Errorf("failed to initialize Vulkan backend: %v", err)
		}
	}

	limits := &vb.properties.Limits

	return &wgpu.Limits{
		MaxTextureDimension1D:                     limits.MaxImageDimension1D,
		MaxTextureDimension2D:                     limits.MaxImageDimension2D,
		MaxTextureDimension3D:                     limits.MaxImageDimension3D,
		MaxTextureArrayLayers:                     limits.MaxImageArrayLayers,
		MaxBindGroups:                             limits.MaxBoundDescriptorSets,
		MaxBindGroupsPlusVertexBuffers:            limits.MaxBoundDescriptorSets + 8, // Conservative estimate
		MaxBindingsPerBindGroup:                   1000,                              // Conservative default
		MaxDynamicUniformBuffersPerPipelineLayout: limits.MaxDescriptorSetUniformBuffersDynamic,
		MaxDynamicStorageBuffersPerPipelineLayout: limits.MaxDescriptorSetStorageBuffersDynamic,
		MaxSampledTexturesPerShaderStage:          limits.MaxPerStageDescriptorSampledImages,
		MaxSamplersPerShaderStage:                 limits.MaxPerStageDescriptorSamplers,
		MaxStorageBuffersPerShaderStage:           limits.MaxPerStageDescriptorStorageBuffers,
		MaxStorageTexturesPerShaderStage:          limits.MaxPerStageDescriptorStorageImages,
		MaxUniformBuffersPerShaderStage:           limits.MaxPerStageDescriptorUniformBuffers,
		MaxUniformBufferBindingSize:               uint64(limits.MaxUniformBufferRange),
		MaxStorageBufferBindingSize:               uint64(limits.MaxStorageBufferRange),
		MinUniformBufferOffsetAlignment:           uint32(limits.MinUniformBufferOffsetAlignment),
		MinStorageBufferOffsetAlignment:           uint32(limits.MinStorageBufferOffsetAlignment),
		MaxVertexBuffers:                          limits.MaxVertexInputBindings,
		MaxBufferSize:                             ^uint64(0), // Use max value
		MaxVertexAttributes:                       limits.MaxVertexInputAttributes,
		MaxVertexBufferArrayStride:                limits.MaxVertexInputBindingStride,
		MaxInterStageShaderVariables:              limits.MaxFragmentInputComponents / 4, // Estimate
		MaxColorAttachments:                       limits.MaxColorAttachments,
		MaxColorAttachmentBytesPerSample:          32, // Conservative default
		MaxComputeWorkgroupStorageSize:            limits.MaxComputeSharedMemorySize,
		MaxComputeInvocationsPerWorkgroup:         limits.MaxComputeWorkGroupInvocations,
		MaxComputeWorkgroupSizeX:                  limits.MaxComputeWorkGroupSize[0],
		MaxComputeWorkgroupSizeY:                  limits.MaxComputeWorkGroupSize[1],
		MaxComputeWorkgroupSizeZ:                  limits.MaxComputeWorkGroupSize[2],
		MaxComputeWorkgroupsPerDimension:          limits.MaxComputeWorkGroupCount[0],
		MaxImmediateSize:                          16777216, // 16MB default
	}, nil
}

// RequestDevice implements wgpu.Adapter.
func (vb *VulkanBackend) RequestDevice(descriptor wgpu.DeviceDescriptor) (wgpu.Device, error) {
	if !vb.initialized {
		if err := vb.Init(); err != nil {
			return nil, fmt.Errorf("failed to initialize Vulkan backend: %v", err)
		}
	}

	// Find graphics queue family
	graphicsQueueFamily := uint32(0)
	found := false
	for i, queueFamily := range vb.queueFamilies {
		if queueFamily.QueueFlags&VK_QUEUE_GRAPHICS_BIT != 0 {
			graphicsQueueFamily = uint32(i)
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("no graphics queue family found")
	}

	// Create device queue create info
	queuePriority := float32(1.0)
	queueCreateInfo := VkDeviceQueueCreateInfo{
		SType:            VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO,
		QueueFamilyIndex: graphicsQueueFamily,
		QueueCount:       1,
		PQueuePriorities: &queuePriority,
	}

	// Create device
	deviceCreateInfo := VkDeviceCreateInfo{
		SType:                VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO,
		QueueCreateInfoCount: 1,
		PQueueCreateInfos:    &queueCreateInfo,
		PEnabledFeatures:     &vb.features,
	}

	var device VkDevice
	result := vkCreateDevice(vb.physicalDevice, &deviceCreateInfo, 0, &device)
	if result != VK_SUCCESS {
		return nil, fmt.Errorf("failed to create Vulkan device: %d", result)
	}

	// TODO: Create and return a proper VulkanDevice implementation
	return nil, fmt.Errorf("device creation not fully implemented yet")
}

// Type implements backend.Backend.
func (vb *VulkanBackend) Type() wgpu.BackendType {
	return wgpu.BackendTypeVulkan
}

// Init implements backend.Backend.
func (vb *VulkanBackend) Init() error {
	if vb.initialized {
		return nil
	}

	// Load Vulkan library
	if err := LoadVulkanLibrary(); err != nil {
		return fmt.Errorf("failed to load Vulkan library: %v", err)
	}

	// Create Vulkan instance
	if err := vb.createInstance(); err != nil {
		return fmt.Errorf("failed to create Vulkan instance: %v", err)
	}

	// Select physical device
	if err := vb.selectPhysicalDevice(); err != nil {
		return fmt.Errorf("failed to select physical device: %v", err)
	}

	// Query device properties and features
	vb.queryDeviceInfo()

	vb.initialized = true
	return nil
}

// createInstance creates a Vulkan instance
func (vb *VulkanBackend) createInstance() error {
	appInfo := VkApplicationInfo{
		SType:              VK_STRUCTURE_TYPE_APPLICATION_INFO,
		PApplicationName:   cString("Sraph Application"),
		ApplicationVersion: VK_API_VERSION_1_0,
		PEngineName:        cString("Sraph Engine"),
		EngineVersion:      VK_API_VERSION_1_0,
		ApiVersion:         VK_API_VERSION_1_0,
	}

	createInfo := VkInstanceCreateInfo{
		SType:            VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO,
		PApplicationInfo: &appInfo,
	}

	result := vkCreateInstance(&createInfo, 0, &vb.instance)
	if result != VK_SUCCESS {
		return fmt.Errorf("failed to create Vulkan instance: %d", result)
	}

	return nil
}

// selectPhysicalDevice selects the best available physical device
func (vb *VulkanBackend) selectPhysicalDevice() error {
	var deviceCount uint32
	result := vkEnumeratePhysicalDevices(vb.instance, &deviceCount, nil)
	if result != VK_SUCCESS || deviceCount == 0 {
		return fmt.Errorf("no Vulkan physical devices found")
	}

	devices := make([]VkPhysicalDevice, deviceCount)
	result = vkEnumeratePhysicalDevices(vb.instance, &deviceCount, &devices[0])
	if result != VK_SUCCESS {
		return fmt.Errorf("failed to enumerate physical devices")
	}

	// Select the first discrete GPU, or fallback to the first device
	vb.physicalDevice = devices[0]
	for _, device := range devices {
		var props VkPhysicalDeviceProperties
		vkGetPhysicalDeviceProperties(device, &props)
		if props.DeviceType == VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU {
			vb.physicalDevice = device
			break
		}
	}

	return nil
}

// queryDeviceInfo queries device properties, features, and queue families
func (vb *VulkanBackend) queryDeviceInfo() {
	// Get device properties
	vkGetPhysicalDeviceProperties(vb.physicalDevice, &vb.properties)

	// Get device features
	vkGetPhysicalDeviceFeatures(vb.physicalDevice, &vb.features)

	// Get memory properties
	vkGetPhysicalDeviceMemoryProperties(vb.physicalDevice, &vb.memoryProperties)

	// Get queue family properties
	var queueFamilyCount uint32
	vkGetPhysicalDeviceQueueFamilyProperties(vb.physicalDevice, &queueFamilyCount, nil)
	if queueFamilyCount > 0 {
		vb.queueFamilies = make([]VkQueueFamilyProperties, queueFamilyCount)
		vkGetPhysicalDeviceQueueFamilyProperties(vb.physicalDevice, &queueFamilyCount, &vb.queueFamilies[0])
	}
}

// Helper functions

// bytePtrToString converts a byte pointer to a Go string
func bytePtrToString(ptr *byte, maxLen int) string {
	if ptr == nil {
		return ""
	}

	bytes := (*[256]byte)(unsafe.Pointer(ptr))
	for i := 0; i < maxLen; i++ {
		if bytes[i] == 0 {
			return string(bytes[:i])
		}
	}
	return string(bytes[:maxLen])
}

// getVendorName maps vendor ID to vendor name
func getVendorName(vendorID uint32) string {
	switch vendorID {
	case 0x1002:
		return "AMD"
	case 0x1010:
		return "ImgTec"
	case 0x10DE:
		return "NVIDIA"
	case 0x13B5:
		return "ARM"
	case 0x5143:
		return "Qualcomm"
	case 0x8086:
		return "Intel"
	default:
		return "Unknown"
	}
}

// Cleanup function to be called when the backend is no longer needed
func (vb *VulkanBackend) Release() error {
	if vb.instance != 0 {
		vkDestroyInstance(vb.instance, 0)
		vb.instance = 0
	}
	UnloadVulkanLibrary()
	vb.initialized = false
	return nil
}
