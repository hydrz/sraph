package vulkan

import (
	"fmt"
	"sync"
)

// VulkanInstance manages Vulkan instance
type VulkanInstance struct {
	mu                sync.RWMutex
	handle            uintptr // VkInstance handle
	physicalDevices   []VulkanPhysicalDevice
	enabledLayers     []string
	enabledExtensions []string
	debugMessenger    uintptr // VkDebugUtilsMessengerEXT handle
	destroyed         bool
}

// VulkanPhysicalDevice represents a Vulkan physical device
type VulkanPhysicalDevice struct {
	Handle           uintptr // VkPhysicalDevice handle
	Properties       VulkanDeviceProperties
	Features         VulkanDeviceFeatures
	QueueFamilyProps []VulkanQueueFamilyProperties
	MemoryProperties VulkanPhysicalDeviceMemoryProperties
	SurfaceSupport   bool
}

// VulkanDeviceProperties contains device properties
type VulkanDeviceProperties struct {
	DeviceName     string
	DeviceType     uint32
	VendorID       uint32
	DeviceID       uint32
	DriverVersion  uint32
	APIVersion     uint32
	MaxTextureSize uint32
	MaxBufferSize  uint64
}

// VulkanDeviceFeatures contains device features
type VulkanDeviceFeatures struct {
	GeometryShader       bool
	TessellationShader   bool
	MultiViewport        bool
	SamplerAnisotropy    bool
	TextureCompressionBC bool
	DepthClamp           bool
}

// VulkanQueueFamilyProperties contains queue family properties
type VulkanQueueFamilyProperties struct {
	QueueFlags       uint32
	QueueCount       uint32
	TimestampBits    uint32
	SupportsGraphics bool
	SupportsCompute  bool
	SupportsTransfer bool
	SupportsSparse   bool
	SupportsPresent  bool
}

// VulkanPhysicalDeviceMemoryProperties contains memory properties
type VulkanPhysicalDeviceMemoryProperties struct {
	MemoryTypeCount uint32
	MemoryHeapCount uint32
	MemoryTypes     []VulkanMemoryType
	MemoryHeaps     []VulkanMemoryHeap
}

// VulkanMemoryType represents a memory type
type VulkanMemoryType struct {
	PropertyFlags uint32
	HeapIndex     uint32
}

// VulkanMemoryHeap represents a memory heap
type VulkanMemoryHeap struct {
	Size  uint64
	Flags uint32
}

// NewVulkanInstance creates a new Vulkan instance
func NewVulkanInstance(enableDebug bool) (*VulkanInstance, error) {
	instance := &VulkanInstance{
		enabledLayers:     []string{},
		enabledExtensions: []string{},
	}

	if enableDebug {
		instance.enabledLayers = append(instance.enabledLayers, "VK_LAYER_KHRONOS_validation")
		instance.enabledExtensions = append(instance.enabledExtensions, "VK_EXT_debug_utils")
	}

	// Add required extensions
	instance.enabledExtensions = append(instance.enabledExtensions,
		"VK_KHR_surface",
		"VK_KHR_win32_surface", // Platform specific, should be conditional
	)

	if err := instance.createInstance(); err != nil {
		return nil, fmt.Errorf("failed to create Vulkan instance: %v", err)
	}

	if err := instance.enumeratePhysicalDevices(); err != nil {
		instance.Destroy()
		return nil, fmt.Errorf("failed to enumerate physical devices: %v", err)
	}

	return instance, nil
}

// createInstance creates the Vulkan instance using real Vulkan API
func (vi *VulkanInstance) createInstance() error {
	// Create application info
	appInfo := VkApplicationInfo{
		SType:              VK_STRUCTURE_TYPE_APPLICATION_INFO,
		PNext:              nil,
		PApplicationName:   toCString("WebGPU Application"),
		ApplicationVersion: 1,
		PEngineName:        toCString("WebGPU Engine"),
		EngineVersion:      1,
		ApiVersion:         0x00401000, // Vulkan 1.1
	}

	// Convert layer names to C strings
	layerNames, layerCStrs := toCStringArray(vi.enabledLayers)
	defer func() {
		// Keep cstrs alive
		_ = layerCStrs
	}()

	// Convert extension names to C strings
	extNames, extCStrs := toCStringArray(vi.enabledExtensions)
	defer func() {
		// Keep cstrs alive
		_ = extCStrs
	}()

	// Create instance info
	createInfo := VkInstanceCreateInfo{
		SType:                   VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO,
		PNext:                   nil,
		Flags:                   0,
		PApplicationInfo:        &appInfo,
		EnabledLayerCount:       uint32(len(vi.enabledLayers)),
		PpEnabledLayerNames:     layerNames,
		EnabledExtensionCount:   uint32(len(vi.enabledExtensions)),
		PpEnabledExtensionNames: extNames,
	}

	var instance VkInstance
	result := vkCreateInstance(&createInfo, nil, &instance)
	if result != VK_SUCCESS {
		return fmt.Errorf("vkCreateInstance failed with result %d", result)
	}

	vi.handle = uintptr(instance)
	return nil
}

// enumeratePhysicalDevices enumerates available physical devices using real Vulkan API
func (vi *VulkanInstance) enumeratePhysicalDevices() error {
	instance := VkInstance(vi.handle)

	// Get device count
	var deviceCount uint32
	result := vkEnumeratePhysicalDevices(instance, &deviceCount, nil)
	if result != VK_SUCCESS {
		return fmt.Errorf("vkEnumeratePhysicalDevices failed with result %d", result)
	}

	if deviceCount == 0 {
		return fmt.Errorf("no Vulkan physical devices found")
	}

	// Get devices
	devices := make([]VkPhysicalDevice, deviceCount)
	result = vkEnumeratePhysicalDevices(instance, &deviceCount, &devices[0])
	if result != VK_SUCCESS {
		return fmt.Errorf("vkEnumeratePhysicalDevices failed with result %d", result)
	}

	// Convert to internal format
	vi.physicalDevices = make([]VulkanPhysicalDevice, deviceCount)
	for i, device := range devices {
		vulkanDevice, err := vi.convertPhysicalDevice(device)
		if err != nil {
			return fmt.Errorf("failed to convert physical device %d: %v", i, err)
		}
		vi.physicalDevices[i] = vulkanDevice
	}

	return nil
}

// convertPhysicalDevice converts VkPhysicalDevice to VulkanPhysicalDevice
func (vi *VulkanInstance) convertPhysicalDevice(device VkPhysicalDevice) (VulkanPhysicalDevice, error) {
	var properties VkPhysicalDeviceProperties
	vkGetPhysicalDeviceProperties(device, &properties)

	var features VkPhysicalDeviceFeatures
	vkGetPhysicalDeviceFeatures(device, &features)

	var memoryProperties VkPhysicalDeviceMemoryProperties
	vkGetPhysicalDeviceMemoryProperties(device, &memoryProperties)

	// Get queue family properties
	var queueFamilyCount uint32
	vkGetPhysicalDeviceQueueFamilyProperties(device, &queueFamilyCount, nil)

	queueFamilies := make([]VkQueueFamilyProperties, queueFamilyCount)
	if queueFamilyCount > 0 {
		vkGetPhysicalDeviceQueueFamilyProperties(device, &queueFamilyCount, &queueFamilies[0])
	}

	// Convert properties
	vulkanDevice := VulkanPhysicalDevice{
		Handle: uintptr(device),
		Properties: VulkanDeviceProperties{
			DeviceName:     cStringToGo(properties.DeviceName[:]),
			DeviceType:     properties.DeviceType,
			VendorID:       properties.VendorID,
			DeviceID:       properties.DeviceID,
			DriverVersion:  properties.DriverVersion,
			APIVersion:     properties.ApiVersion,
			MaxTextureSize: properties.Limits.MaxImageDimension2D,
			MaxBufferSize:  uint64(properties.Limits.MaxStorageBufferRange),
		},
		Features: VulkanDeviceFeatures{
			GeometryShader:       features.GeometryShader != 0,
			TessellationShader:   features.TessellationShader != 0,
			MultiViewport:        features.MultiViewport != 0,
			SamplerAnisotropy:    features.SamplerAnisotropy != 0,
			TextureCompressionBC: features.TextureCompressionBC != 0,
			DepthClamp:           features.DepthClamp != 0,
		},
		QueueFamilyProps: convertQueueFamilyProperties(queueFamilies),
		MemoryProperties: convertMemoryProperties(memoryProperties),
		SurfaceSupport:   true, // Assume surface support for now
	}

	return vulkanDevice, nil
}

// Helper functions for conversion
func cStringToGo(data []byte) string {
	// Find null terminator
	end := 0
	for i, b := range data {
		if b == 0 {
			end = i
			break
		}
	}
	return string(data[:end])
}

func convertQueueFamilyProperties(props []VkQueueFamilyProperties) []VulkanQueueFamilyProperties {
	result := make([]VulkanQueueFamilyProperties, len(props))
	for i, prop := range props {
		result[i] = VulkanQueueFamilyProperties{
			QueueFlags:       prop.QueueFlags,
			QueueCount:       prop.QueueCount,
			TimestampBits:    prop.TimestampValidBits,
			SupportsGraphics: (prop.QueueFlags & VK_QUEUE_GRAPHICS_BIT) != 0,
			SupportsCompute:  (prop.QueueFlags & VK_QUEUE_COMPUTE_BIT) != 0,
			SupportsTransfer: (prop.QueueFlags & VK_QUEUE_TRANSFER_BIT) != 0,
			SupportsSparse:   (prop.QueueFlags & VK_QUEUE_SPARSE_BINDING_BIT) != 0,
			SupportsPresent:  true, // Check surface support separately
		}
	}
	return result
}

func convertMemoryProperties(props VkPhysicalDeviceMemoryProperties) VulkanPhysicalDeviceMemoryProperties {
	memTypes := make([]VulkanMemoryType, props.MemoryTypeCount)
	for i := uint32(0); i < props.MemoryTypeCount; i++ {
		memTypes[i] = VulkanMemoryType{
			PropertyFlags: props.MemoryTypes[i].PropertyFlags,
			HeapIndex:     props.MemoryTypes[i].HeapIndex,
		}
	}

	memHeaps := make([]VulkanMemoryHeap, props.MemoryHeapCount)
	for i := uint32(0); i < props.MemoryHeapCount; i++ {
		memHeaps[i] = VulkanMemoryHeap{
			Size:  props.MemoryHeaps[i].Size,
			Flags: props.MemoryHeaps[i].Flags,
		}
	}

	return VulkanPhysicalDeviceMemoryProperties{
		MemoryTypeCount: props.MemoryTypeCount,
		MemoryHeapCount: props.MemoryHeapCount,
		MemoryTypes:     memTypes,
		MemoryHeaps:     memHeaps,
	}
}

// Destroy destroys the Vulkan instance using real Vulkan API
func (vi *VulkanInstance) Destroy() error {
	vi.mu.Lock()
	defer vi.mu.Unlock()

	if vi.destroyed {
		return nil
	}

	if vi.handle != 0 {
		vkDestroyInstance(VkInstance(vi.handle), nil)
		vi.handle = 0
	}

	vi.physicalDevices = nil
	vi.destroyed = true

	return nil
}

// IsDestroyed checks if the instance is destroyed
func (vi *VulkanInstance) IsDestroyed() bool {
	vi.mu.RLock()
	defer vi.mu.RUnlock()
	return vi.destroyed
}

// GetHandle returns the Vulkan instance handle
func (vi *VulkanInstance) GetHandle() uintptr {
	vi.mu.RLock()
	defer vi.mu.RUnlock()
	return vi.handle
}
