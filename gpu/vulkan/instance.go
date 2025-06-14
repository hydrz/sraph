package vulkan

import (
	"fmt"
	"sync"
	"unsafe"
)

// VulkanInstance manages Vulkan instance
type VulkanInstance struct {
	mu                sync.RWMutex
	handle            VkInstance
	physicalDevices   []VulkanPhysicalDevice
	enabledLayers     []string
	enabledExtensions []string
	debugMessenger    VkDebugUtilsMessengerEXT
	destroyed         bool
}

// VulkanPhysicalDevice represents a Vulkan physical device
type VulkanPhysicalDevice struct {
	handle              VkPhysicalDevice
	properties          VkPhysicalDeviceProperties
	features            VkPhysicalDeviceFeatures
	queueFamilies       []VkQueueFamilyProperties
	memoryProperties    VkPhysicalDeviceMemoryProperties
	supportedExtensions []string
}

// VulkanDeviceProperties contains device properties
type VulkanDeviceProperties struct {
	ApiVersion        uint32
	DriverVersion     uint32
	VendorID          uint32
	DeviceID          uint32
	DeviceType        uint32
	DeviceName        string
	PipelineCacheUUID [16]byte
}

// VulkanDeviceFeatures contains device features
type VulkanDeviceFeatures struct {
	RobustBufferAccess  bool
	FullDrawIndexUint32 bool
	ImageCubeArray      bool
	IndependentBlend    bool
	GeometryShader      bool
	TessellationShader  bool
	SampleRateShading   bool
	// Add more features as needed
}

// VulkanQueueFamilyProperties contains queue family properties
type VulkanQueueFamilyProperties struct {
	QueueFlags                  uint32
	QueueCount                  uint32
	TimestampValidBits          uint32
	MinImageTransferGranularity VkExtent3D
}

// VulkanPhysicalDeviceMemoryProperties contains memory properties
type VulkanPhysicalDeviceMemoryProperties struct {
	MemoryTypes []VulkanMemoryType
	MemoryHeaps []VulkanMemoryHeap
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
func NewVulkanInstance(appName, engineName string, enableValidation bool) (*VulkanInstance, error) {
	// Ensure Vulkan library is loaded
	if vulkanLib == 0 {
		if err := LoadVulkanLibrary(); err != nil {
			return nil, fmt.Errorf("failed to load Vulkan library: %v", err)
		}
	}

	instance := &VulkanInstance{
		enabledLayers:     make([]string, 0),
		enabledExtensions: make([]string, 0),
	}

	// Add validation layers if requested
	if enableValidation {
		instance.enabledLayers = append(instance.enabledLayers, "VK_LAYER_KHRONOS_validation")
		instance.enabledExtensions = append(instance.enabledExtensions, "VK_EXT_debug_utils")
	}

	// Create the instance
	if err := instance.createInstance(appName, engineName); err != nil {
		return nil, err
	}

	// Enumerate physical devices
	if err := instance.enumeratePhysicalDevices(); err != nil {
		instance.Destroy()
		return nil, err
	}

	return instance, nil
}

// createInstance creates the Vulkan instance using real Vulkan API
func (vi *VulkanInstance) createInstance(appName, engineName string) error {
	vi.mu.Lock()
	defer vi.mu.Unlock()

	if vi.destroyed {
		return fmt.Errorf("instance has been destroyed")
	}

	// Application info
	appInfo := VkApplicationInfo{
		SType:              VK_STRUCTURE_TYPE_APPLICATION_INFO,
		PNext:              0,
		PApplicationName:   cString(appName),
		ApplicationVersion: VK_API_VERSION_1_0,
		PEngineName:        cString(engineName),
		EngineVersion:      VK_API_VERSION_1_0,
		ApiVersion:         VK_API_VERSION_1_0,
	}
	// Convert layer names
	var layerNames **byte
	if len(vi.enabledLayers) > 0 {
		layerNames, _ = cStringArray(vi.enabledLayers)
	}

	// Convert extension names
	var extensionNames **byte
	if len(vi.enabledExtensions) > 0 {
		extensionNames, _ = cStringArray(vi.enabledExtensions)
	}

	// Instance create info
	createInfo := VkInstanceCreateInfo{
		SType:                   VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO,
		PNext:                   0,
		Flags:                   0,
		PApplicationInfo:        &appInfo,
		EnabledLayerCount:       uint32(len(vi.enabledLayers)),
		PpEnabledLayerNames:     layerNames,
		EnabledExtensionCount:   uint32(len(vi.enabledExtensions)),
		PpEnabledExtensionNames: extensionNames,
	}

	// Create instance
	result := vkCreateInstance(&createInfo, 0, &vi.handle)
	if result != VK_SUCCESS {
		return fmt.Errorf("failed to create Vulkan instance: %d", result)
	}

	return nil
}

// enumeratePhysicalDevices enumerates available physical devices using real Vulkan API
func (vi *VulkanInstance) enumeratePhysicalDevices() error {
	vi.mu.Lock()
	defer vi.mu.Unlock()

	if vi.destroyed {
		return fmt.Errorf("instance has been destroyed")
	}

	// First call to get count
	var deviceCount uint32
	result := vkEnumeratePhysicalDevices(vi.handle, &deviceCount, nil)
	if result != VK_SUCCESS {
		return fmt.Errorf("failed to enumerate physical devices: %d", result)
	}

	if deviceCount == 0 {
		return fmt.Errorf("no Vulkan-compatible devices found")
	}

	// Second call to get devices
	devices := make([]VkPhysicalDevice, deviceCount)
	result = vkEnumeratePhysicalDevices(vi.handle, &deviceCount, &devices[0])
	if result != VK_SUCCESS {
		return fmt.Errorf("failed to get physical devices: %d", result)
	}

	// Convert to VulkanPhysicalDevice
	vi.physicalDevices = make([]VulkanPhysicalDevice, len(devices))
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
	var vulkanDevice VulkanPhysicalDevice
	vulkanDevice.handle = device

	// Get device properties
	vkGetPhysicalDeviceProperties(device, &vulkanDevice.properties)

	// Get device features
	vkGetPhysicalDeviceFeatures(device, &vulkanDevice.features)

	// Get queue family properties
	var queueFamilyCount uint32
	vkGetPhysicalDeviceQueueFamilyProperties(device, &queueFamilyCount, nil)

	if queueFamilyCount > 0 {
		queueFamilies := make([]VkQueueFamilyProperties, queueFamilyCount)
		vkGetPhysicalDeviceQueueFamilyProperties(device, &queueFamilyCount, &queueFamilies[0])
		vulkanDevice.queueFamilies = queueFamilies
	}

	// Get memory properties
	vkGetPhysicalDeviceMemoryProperties(device, &vulkanDevice.memoryProperties)

	return vulkanDevice, nil
}

// GetPhysicalDevices returns all physical devices
func (vi *VulkanInstance) GetPhysicalDevices() []VulkanPhysicalDevice {
	vi.mu.RLock()
	defer vi.mu.RUnlock()

	devices := make([]VulkanPhysicalDevice, len(vi.physicalDevices))
	copy(devices, vi.physicalDevices)
	return devices
}

// GetPhysicalDevice returns a specific physical device by index
func (vi *VulkanInstance) GetPhysicalDevice(index int) (*VulkanPhysicalDevice, error) {
	vi.mu.RLock()
	defer vi.mu.RUnlock()

	if index < 0 || index >= len(vi.physicalDevices) {
		return nil, fmt.Errorf("invalid physical device index: %d", index)
	}

	return &vi.physicalDevices[index], nil
}

// Destroy destroys the Vulkan instance using real Vulkan API
func (vi *VulkanInstance) Destroy() error {
	vi.mu.Lock()
	defer vi.mu.Unlock()

	if vi.destroyed {
		return nil
	}

	// Destroy debug messenger if created
	if vi.debugMessenger != 0 {
		// vkDestroyDebugUtilsMessengerEXT would be called here
		vi.debugMessenger = 0
	}

	// Destroy instance
	if vi.handle != 0 {
		vkDestroyInstance(vi.handle, 0)
		vi.handle = 0
	}

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
func (vi *VulkanInstance) GetHandle() VkInstance {
	vi.mu.RLock()
	defer vi.mu.RUnlock()
	return vi.handle
}

// Helper functions for conversion
func (vd *VulkanPhysicalDevice) GetDeviceType() string {
	switch vd.properties.DeviceType {
	case VK_PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU:
		return "Integrated GPU"
	case VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU:
		return "Discrete GPU"
	case VK_PHYSICAL_DEVICE_TYPE_VIRTUAL_GPU:
		return "Virtual GPU"
	case VK_PHYSICAL_DEVICE_TYPE_CPU:
		return "CPU"
	default:
		return "Other"
	}
}

func (vd *VulkanPhysicalDevice) GetDeviceName() string {
	// Convert C string to Go string
	name := (*[256]byte)(unsafe.Pointer(&vd.properties.DeviceName[0]))
	for i, b := range name {
		if b == 0 {
			return string(name[:i])
		}
	}
	return string(name[:])
}

func (vd *VulkanPhysicalDevice) HasGraphicsQueue() bool {
	for _, qf := range vd.queueFamilies {
		if qf.QueueFlags&VK_QUEUE_GRAPHICS_BIT != 0 {
			return true
		}
	}
	return false
}

func (vd *VulkanPhysicalDevice) HasComputeQueue() bool {
	for _, qf := range vd.queueFamilies {
		if qf.QueueFlags&VK_QUEUE_COMPUTE_BIT != 0 {
			return true
		}
	}
	return false
}

func (vd *VulkanPhysicalDevice) GetGraphicsQueueFamilyIndex() (uint32, bool) {
	for i, qf := range vd.queueFamilies {
		if qf.QueueFlags&VK_QUEUE_GRAPHICS_BIT != 0 {
			return uint32(i), true
		}
	}
	return 0, false
}

func (vd *VulkanPhysicalDevice) GetComputeQueueFamilyIndex() (uint32, bool) {
	for i, qf := range vd.queueFamilies {
		if qf.QueueFlags&VK_QUEUE_COMPUTE_BIT != 0 {
			return uint32(i), true
		}
	}
	return 0, false
}
