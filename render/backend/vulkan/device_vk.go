package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// LogicalDeviceVK wraps Vulkan logical device
// It represents a logical connection to a physical device
type LogicalDeviceVK struct {
	device             vulkan.Device
	physicalDevice     vulkan.PhysicalDevice
	queueFamilyIndices QueueFamilyIndices
	graphicsQueue      vulkan.Queue
	presentQueue       vulkan.Queue
	computeQueue       vulkan.Queue
	transferQueue      vulkan.Queue
	enabledExtensions  []string
	enabledFeatures    vulkan.PhysicalDeviceFeatures
}

// QueueFamilyIndices contains queue family indices
type QueueFamilyIndices struct {
	GraphicsFamily *uint32
	PresentFamily  *uint32
	ComputeFamily  *uint32
	TransferFamily *uint32
}

// IsComplete checks if all required queue families are available
func (indices *QueueFamilyIndices) IsComplete() bool {
	return indices.GraphicsFamily != nil && indices.PresentFamily != nil
}

// DeviceCreateInfo contains device creation parameters
type DeviceCreateInfo struct {
	PhysicalDevice    vulkan.PhysicalDevice
	EnabledExtensions []string
	EnabledFeatures   vulkan.PhysicalDeviceFeatures
	Surface           vulkan.Surface
}

// NewLogicalDeviceVK creates a new logical device
func NewLogicalDeviceVK(createInfo *DeviceCreateInfo) (*LogicalDeviceVK, error) {
	if createInfo == nil {
		return nil, fmt.Errorf("device create info cannot be nil")
	}

	// Find queue families
	queueFamilyIndices := findQueueFamilies(createInfo.PhysicalDevice, createInfo.Surface)
	if !queueFamilyIndices.IsComplete() {
		return nil, fmt.Errorf("required queue families not available")
	}

	// Create queue create infos
	queueCreateInfos := createQueueCreateInfos(queueFamilyIndices)

	// Device create info
	deviceCreateInfo := vulkan.DeviceCreateInfo{
		SType:                vulkan.StructureTypeDeviceCreateInfo,
		QueueCreateInfoCount: uint32(len(queueCreateInfos)),
		PQueueCreateInfos:    queueCreateInfos,
		PEnabledFeatures:     &createInfo.EnabledFeatures,
	}

	// Set enabled extensions
	if len(createInfo.EnabledExtensions) > 0 {
		deviceCreateInfo.EnabledExtensionCount = uint32(len(createInfo.EnabledExtensions))
		deviceCreateInfo.PpEnabledExtensionNames = createInfo.EnabledExtensions
	}

	// Create device
	var device vulkan.Device
	if result := vulkan.CreateDevice(createInfo.PhysicalDevice, &deviceCreateInfo, nil, &device); result != vulkan.Success {
		return nil, fmt.Errorf("failed to create logical device: %s", result)
	}

	// Load device-level functions
	vulkan.InitDevice(device)

	// Get queue handles
	var graphicsQueue, presentQueue, computeQueue, transferQueue vulkan.Queue

	if queueFamilyIndices.GraphicsFamily != nil {
		vulkan.GetDeviceQueue(device, *queueFamilyIndices.GraphicsFamily, 0, &graphicsQueue)
	}

	if queueFamilyIndices.PresentFamily != nil {
		vulkan.GetDeviceQueue(device, *queueFamilyIndices.PresentFamily, 0, &presentQueue)
	}

	if queueFamilyIndices.ComputeFamily != nil {
		vulkan.GetDeviceQueue(device, *queueFamilyIndices.ComputeFamily, 0, &computeQueue)
	}

	if queueFamilyIndices.TransferFamily != nil {
		vulkan.GetDeviceQueue(device, *queueFamilyIndices.TransferFamily, 0, &transferQueue)
	}

	return &LogicalDeviceVK{
		device:             device,
		physicalDevice:     createInfo.PhysicalDevice,
		queueFamilyIndices: queueFamilyIndices,
		graphicsQueue:      graphicsQueue,
		presentQueue:       presentQueue,
		computeQueue:       computeQueue,
		transferQueue:      transferQueue,
		enabledExtensions:  createInfo.EnabledExtensions,
		enabledFeatures:    createInfo.EnabledFeatures,
	}, nil
}

// GetDevice returns the Vulkan device handle
func (ld *LogicalDeviceVK) GetDevice() vulkan.Device {
	return ld.device
}

// GetPhysicalDevice returns the physical device
func (ld *LogicalDeviceVK) GetPhysicalDevice() vulkan.PhysicalDevice {
	return ld.physicalDevice
}

// GetQueueFamilyIndices returns the queue family indices
func (ld *LogicalDeviceVK) GetQueueFamilyIndices() QueueFamilyIndices {
	return ld.queueFamilyIndices
}

// GetGraphicsQueue returns the graphics queue
func (ld *LogicalDeviceVK) GetGraphicsQueue() vulkan.Queue {
	return ld.graphicsQueue
}

// GetPresentQueue returns the present queue
func (ld *LogicalDeviceVK) GetPresentQueue() vulkan.Queue {
	return ld.presentQueue
}

// GetComputeQueue returns the compute queue
func (ld *LogicalDeviceVK) GetComputeQueue() vulkan.Queue {
	return ld.computeQueue
}

// GetTransferQueue returns the transfer queue
func (ld *LogicalDeviceVK) GetTransferQueue() vulkan.Queue {
	return ld.transferQueue
}

// GetEnabledExtensions returns the enabled extensions
func (ld *LogicalDeviceVK) GetEnabledExtensions() []string {
	return ld.enabledExtensions
}

// GetEnabledFeatures returns the enabled features
func (ld *LogicalDeviceVK) GetEnabledFeatures() vulkan.PhysicalDeviceFeatures {
	return ld.enabledFeatures
}

// WaitIdle waits for the device to become idle
func (ld *LogicalDeviceVK) WaitIdle() error {
	if result := vulkan.DeviceWaitIdle(ld.device); result != vulkan.Success {
		return fmt.Errorf("failed to wait for device idle: %s", result)
	}
	return nil
}

// Destroy destroys the logical device
func (ld *LogicalDeviceVK) Destroy() {
	if ld.device != vulkan.NullDevice {
		vulkan.DestroyDevice(ld.device, nil)
		ld.device = vulkan.NullDevice
	}
}

// findQueueFamilies finds queue families for a physical device
func findQueueFamilies(physicalDevice vulkan.PhysicalDevice, surface vulkan.Surface) QueueFamilyIndices {
	var queueFamilyCount uint32
	vulkan.GetPhysicalDeviceQueueFamilyProperties(physicalDevice, &queueFamilyCount, nil)

	queueFamilies := make([]vulkan.QueueFamilyProperties, queueFamilyCount)
	vulkan.GetPhysicalDeviceQueueFamilyProperties(physicalDevice, &queueFamilyCount, queueFamilies)

	indices := QueueFamilyIndices{}

	for i, queueFamily := range queueFamilies {
		index := uint32(i)

		// Check for graphics support
		if queueFamily.QueueFlags&vulkan.QueueFlags(vulkan.QueueGraphicsBit) != 0 {
			indices.GraphicsFamily = &index
		}

		// Check for compute support
		if queueFamily.QueueFlags&vulkan.QueueFlags(vulkan.QueueComputeBit) != 0 {
			indices.ComputeFamily = &index
		}

		// Check for transfer support
		if queueFamily.QueueFlags&vulkan.QueueFlags(vulkan.QueueTransferBit) != 0 {
			indices.TransferFamily = &index
		}

		// Check for presentation support
		if surface != vulkan.NullSurface {
			var presentSupport vulkan.Bool32
			vulkan.GetPhysicalDeviceSurfaceSupport(physicalDevice, index, surface, &presentSupport)
			if presentSupport == vulkan.True {
				indices.PresentFamily = &index
			}
		}

		// Early exit if we have everything we need
		if indices.IsComplete() {
			break
		}
	}

	return indices
}

// createQueueCreateInfos creates queue create info structures
func createQueueCreateInfos(indices QueueFamilyIndices) []vulkan.DeviceQueueCreateInfo {
	queueCreateInfos := make([]vulkan.DeviceQueueCreateInfo, 0)
	uniqueQueueFamilies := make(map[uint32]bool)

	// Collect unique queue family indices
	if indices.GraphicsFamily != nil {
		uniqueQueueFamilies[*indices.GraphicsFamily] = true
	}
	if indices.PresentFamily != nil {
		uniqueQueueFamilies[*indices.PresentFamily] = true
	}
	if indices.ComputeFamily != nil {
		uniqueQueueFamilies[*indices.ComputeFamily] = true
	}
	if indices.TransferFamily != nil {
		uniqueQueueFamilies[*indices.TransferFamily] = true
	}

	// Create queue create infos
	queuePriority := float32(1.0)
	for queueFamily := range uniqueQueueFamilies {
		queueCreateInfo := vulkan.DeviceQueueCreateInfo{
			SType:            vulkan.StructureTypeDeviceQueueCreateInfo,
			QueueFamilyIndex: queueFamily,
			QueueCount:       1,
			PQueuePriorities: []float32{queuePriority},
		}
		queueCreateInfos = append(queueCreateInfos, queueCreateInfo)
	}

	return queueCreateInfos
}

// GetRequiredDeviceExtensions returns the required device extensions
func GetRequiredDeviceExtensions() []string {
	return []string{
		vulkan.KhrSwapchainExtensionName,
	}
}

// DeviceManagerVK manages multiple logical devices
// It provides centralized device management
type DeviceManagerVK struct {
	instance      *InstanceVK
	devices       []*LogicalDeviceVK
	primaryDevice *LogicalDeviceVK
}

// NewDeviceManagerVK creates a new device manager
func NewDeviceManagerVK(instance *InstanceVK) *DeviceManagerVK {
	return &DeviceManagerVK{
		instance: instance,
		devices:  make([]*LogicalDeviceVK, 0),
	}
}

// CreateDevice creates a new logical device
func (dm *DeviceManagerVK) CreateDevice(createInfo *DeviceCreateInfo) (*LogicalDeviceVK, error) {
	device, err := NewLogicalDeviceVK(createInfo)
	if err != nil {
		return nil, err
	}

	dm.devices = append(dm.devices, device)

	// Set as primary device if it's the first one
	if dm.primaryDevice == nil {
		dm.primaryDevice = device
	}

	return device, nil
}

// CreateDefaultDevice creates a device with default settings
func (dm *DeviceManagerVK) CreateDefaultDevice() (*LogicalDeviceVK, error) {
	// Find best physical device
	physicalDevice, err := dm.instance.FindBestPhysicalDevice()
	if err != nil {
		return nil, err
	}

	// Get required extensions
	requiredExtensions := GetRequiredDeviceExtensions()

	// Get available features
	availableFeatures := dm.instance.GetPhysicalDeviceFeatures(physicalDevice)

	// Create device with default settings
	createInfo := &DeviceCreateInfo{
		PhysicalDevice:    physicalDevice,
		EnabledExtensions: requiredExtensions,
		EnabledFeatures:   availableFeatures,
		Surface:           dm.instance.GetSurface(),
	}

	return dm.CreateDevice(createInfo)
}

// GetPrimaryDevice returns the primary device
func (dm *DeviceManagerVK) GetPrimaryDevice() *LogicalDeviceVK {
	return dm.primaryDevice
}

// GetDeviceCount returns the number of devices
func (dm *DeviceManagerVK) GetDeviceCount() int {
	return len(dm.devices)
}

// WaitAllDevicesIdle waits for all devices to become idle
func (dm *DeviceManagerVK) WaitAllDevicesIdle() error {
	for _, device := range dm.devices {
		if err := device.WaitIdle(); err != nil {
			return err
		}
	}
	return nil
}

// Destroy destroys all devices
func (dm *DeviceManagerVK) Destroy() {
	for _, device := range dm.devices {
		device.Destroy()
	}
	dm.devices = dm.devices[:0]
	dm.primaryDevice = nil
}
