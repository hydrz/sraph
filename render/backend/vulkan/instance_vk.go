package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// InstanceVK wraps Vulkan instance
// It represents the Vulkan API instance and manages global state
type InstanceVK struct {
	instance          vulkan.Instance
	enabledLayers     []string
	enabledExtensions []string
	debugUtils        *DebugUtilsVK
	surface           vulkan.Surface
}

// InstanceCreateInfo contains instance creation parameters
type InstanceCreateInfo struct {
	ApplicationName    string
	ApplicationVersion uint32
	EngineName         string
	EngineVersion      uint32
	APIVersion         uint32
	EnabledLayers      []string
	EnabledExtensions  []string
	EnableValidation   bool
}

// NewInstanceVK creates a new Vulkan instance
func NewInstanceVK(createInfo *InstanceCreateInfo) (*InstanceVK, error) {
	if createInfo == nil {
		return nil, fmt.Errorf("instance create info cannot be nil")
	}

	// Initialize Vulkan
	if err := vulkan.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize Vulkan: %v", err)
	}

	// Application info
	appInfo := vulkan.ApplicationInfo{
		SType:              vulkan.StructureTypeApplicationInfo,
		PApplicationName:   createInfo.ApplicationName,
		ApplicationVersion: createInfo.ApplicationVersion,
		PEngineName:        createInfo.EngineName,
		EngineVersion:      createInfo.EngineVersion,
		ApiVersion:         createInfo.APIVersion,
	}

	// Instance create info
	instanceCreateInfo := vulkan.InstanceCreateInfo{
		SType:            vulkan.StructureTypeInstanceCreateInfo,
		PApplicationInfo: &appInfo,
	}

	// Set enabled layers
	if len(createInfo.EnabledLayers) > 0 {
		instanceCreateInfo.EnabledLayerCount = uint32(len(createInfo.EnabledLayers))
		instanceCreateInfo.PpEnabledLayerNames = createInfo.EnabledLayers
	}

	// Set enabled extensions
	if len(createInfo.EnabledExtensions) > 0 {
		instanceCreateInfo.EnabledExtensionCount = uint32(len(createInfo.EnabledExtensions))
		instanceCreateInfo.PpEnabledExtensionNames = createInfo.EnabledExtensions
	}

	// Create instance
	var instance vulkan.Instance
	if result := vulkan.CreateInstance(&instanceCreateInfo, nil, &instance); result != vulkan.Success {
		return nil, fmt.Errorf("failed to create Vulkan instance: %s", result)
	}

	// Load instance-level functions
	vulkan.InitInstance(instance)

	// Create debug utils if validation is enabled
	var debugUtils *DebugUtilsVK
	if createInfo.EnableValidation {
		debugUtils = NewDebugUtilsVK(instance)
		debugUtils.SetEnabled(true)
	}

	return &InstanceVK{
		instance:          instance,
		enabledLayers:     createInfo.EnabledLayers,
		enabledExtensions: createInfo.EnabledExtensions,
		debugUtils:        debugUtils,
	}, nil
}

// GetInstance returns the Vulkan instance handle
func (inst *InstanceVK) GetInstance() vulkan.Instance {
	return inst.instance
}

// GetEnabledLayers returns the enabled layers
func (inst *InstanceVK) GetEnabledLayers() []string {
	return inst.enabledLayers
}

// GetEnabledExtensions returns the enabled extensions
func (inst *InstanceVK) GetEnabledExtensions() []string {
	return inst.enabledExtensions
}

// GetDebugUtils returns the debug utils manager
func (inst *InstanceVK) GetDebugUtils() *DebugUtilsVK {
	return inst.debugUtils
}

// CreateSurface creates a surface for the given window
func (inst *InstanceVK) CreateSurface(windowHandle uintptr) (vulkan.Surface, error) {
	// TODO: Implement platform-specific surface creation
	// This requires platform-specific code for different windowing systems
	var surface vulkan.Surface
	return surface, fmt.Errorf("surface creation not implemented")
}

// GetSurface returns the current surface
func (inst *InstanceVK) GetSurface() vulkan.Surface {
	return inst.surface
}

// SetSurface sets the surface
func (inst *InstanceVK) SetSurface(surface vulkan.Surface) {
	inst.surface = surface
}

// EnumeratePhysicalDevices enumerates available physical devices
func (inst *InstanceVK) EnumeratePhysicalDevices() ([]vulkan.PhysicalDevice, error) {
	var deviceCount uint32
	if result := vulkan.EnumeratePhysicalDevices(inst.instance, &deviceCount, nil); result != vulkan.Success {
		return nil, fmt.Errorf("failed to enumerate physical devices: %s", result)
	}

	if deviceCount == 0 {
		return nil, fmt.Errorf("no Vulkan-compatible devices found")
	}

	devices := make([]vulkan.PhysicalDevice, deviceCount)
	if result := vulkan.EnumeratePhysicalDevices(inst.instance, &deviceCount, devices); result != vulkan.Success {
		return nil, fmt.Errorf("failed to get physical devices: %s", result)
	}

	return devices, nil
}

// GetPhysicalDeviceProperties gets properties for a physical device
func (inst *InstanceVK) GetPhysicalDeviceProperties(device vulkan.PhysicalDevice) vulkan.PhysicalDeviceProperties {
	var props vulkan.PhysicalDeviceProperties
	vulkan.GetPhysicalDeviceProperties(device, &props)
	return props
}

// GetPhysicalDeviceFeatures gets features for a physical device
func (inst *InstanceVK) GetPhysicalDeviceFeatures(device vulkan.PhysicalDevice) vulkan.PhysicalDeviceFeatures {
	var features vulkan.PhysicalDeviceFeatures
	vulkan.GetPhysicalDeviceFeatures(device, &features)
	return features
}

// GetPhysicalDeviceMemoryProperties gets memory properties for a physical device
func (inst *InstanceVK) GetPhysicalDeviceMemoryProperties(device vulkan.PhysicalDevice) vulkan.PhysicalDeviceMemoryProperties {
	var memProps vulkan.PhysicalDeviceMemoryProperties
	vulkan.GetPhysicalDeviceMemoryProperties(device, &memProps)
	return memProps
}

// FindBestPhysicalDevice finds the best physical device for rendering
func (inst *InstanceVK) FindBestPhysicalDevice() (vulkan.PhysicalDevice, error) {
	devices, err := inst.EnumeratePhysicalDevices()
	if err != nil {
		return vulkan.NullPhysicalDevice, err
	}

	// Score devices and pick the best one
	bestDevice := vulkan.NullPhysicalDevice
	bestScore := -1

	for _, device := range devices {
		score := inst.scorePhysicalDevice(device)
		if score > bestScore {
			bestScore = score
			bestDevice = device
		}
	}

	if bestDevice == vulkan.NullPhysicalDevice {
		return vulkan.NullPhysicalDevice, fmt.Errorf("no suitable physical device found")
	}

	return bestDevice, nil
}

// scorePhysicalDevice scores a physical device based on its capabilities
func (inst *InstanceVK) scorePhysicalDevice(device vulkan.PhysicalDevice) int {
	props := inst.GetPhysicalDeviceProperties(device)
	features := inst.GetPhysicalDeviceFeatures(device)

	score := 0

	// Prefer discrete GPUs
	if props.DeviceType == vulkan.PhysicalDeviceTypeDiscreteGpu {
		score += 1000
	}

	// Maximum possible size of textures affects graphics quality
	score += int(props.Limits.MaxImageDimension2D)

	// Application requires geometry shaders
	if features.GeometryShader == vulkan.False {
		return 0 // Disqualify devices without geometry shader
	}

	// Check for required extensions
	if !inst.checkDeviceExtensionSupport(device, []string{vulkan.KhrSwapchainExtensionName}) {
		return 0 // Disqualify devices without swapchain support
	}

	// Check swapchain support
	if inst.surface != vulkan.NullSurface {
		if !inst.checkSwapchainSupport(device) {
			return 0 // Disqualify devices without adequate swapchain support
		}
	}

	return score
}

// checkDeviceExtensionSupport checks if device supports required extensions
func (inst *InstanceVK) checkDeviceExtensionSupport(device vulkan.PhysicalDevice, requiredExtensions []string) bool {
	var extensionCount uint32
	vulkan.EnumerateDeviceExtensionProperties(device, "", &extensionCount, nil)

	if extensionCount == 0 {
		return len(requiredExtensions) == 0
	}

	availableExtensions := make([]vulkan.ExtensionProperties, extensionCount)
	vulkan.EnumerateDeviceExtensionProperties(device, "", &extensionCount, availableExtensions)

	for _, required := range requiredExtensions {
		found := false
		for _, available := range availableExtensions {
			if required == vulkan.ToString(available.ExtensionName[:]) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// checkSwapchainSupport checks if device has adequate swapchain support
func (inst *InstanceVK) checkSwapchainSupport(device vulkan.PhysicalDevice) bool {
	// Check surface capabilities
	var capabilities vulkan.SurfaceCapabilities
	if result := vulkan.GetPhysicalDeviceSurfaceCapabilities(device, inst.surface, &capabilities); result != vulkan.Success {
		return false
	}

	// Check surface formats
	var formatCount uint32
	vulkan.GetPhysicalDeviceSurfaceFormats(device, inst.surface, &formatCount, nil)

	// Check present modes
	var presentModeCount uint32
	vulkan.GetPhysicalDeviceSurfacePresentModes(device, inst.surface, &presentModeCount, nil)

	return formatCount > 0 && presentModeCount > 0
}

// Destroy destroys the Vulkan instance
func (inst *InstanceVK) Destroy() {
	if inst.debugUtils != nil {
		inst.debugUtils.Destroy()
	}

	if inst.surface != vulkan.NullSurface {
		vulkan.DestroySurface(inst.instance, inst.surface, nil)
		inst.surface = vulkan.NullSurface
	}

	if inst.instance != vulkan.NullInstance {
		vulkan.DestroyInstance(inst.instance, nil)
		inst.instance = vulkan.NullInstance
	}
}

// GetRequiredExtensions returns the required instance extensions
func GetRequiredInstanceExtensions() []string {
	// TODO: Get platform-specific extensions
	return []string{
		vulkan.KhrSurfaceExtensionName,
		// Add platform-specific surface extensions
	}
}

// CheckInstanceExtensionSupport checks if required extensions are supported
func CheckInstanceExtensionSupport(requiredExtensions []string) error {
	var extensionCount uint32
	if result := vulkan.EnumerateInstanceExtensionProperties("", &extensionCount, nil); result != vulkan.Success {
		return fmt.Errorf("failed to enumerate instance extensions: %s", result)
	}

	if extensionCount == 0 {
		if len(requiredExtensions) > 0 {
			return fmt.Errorf("no instance extensions available")
		}
		return nil
	}

	availableExtensions := make([]vulkan.ExtensionProperties, extensionCount)
	if result := vulkan.EnumerateInstanceExtensionProperties("", &extensionCount, availableExtensions); result != vulkan.Success {
		return fmt.Errorf("failed to get instance extensions: %s", result)
	}

	for _, required := range requiredExtensions {
		found := false
		for _, available := range availableExtensions {
			if required == vulkan.ToString(available.ExtensionName[:]) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("required extension not available: %s", required)
		}
	}

	return nil
}
