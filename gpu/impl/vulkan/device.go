package vulkan

import (
	"fmt"
	"sync"
	"unsafe"
)

// VulkanDevice represents a Vulkan logical device
type VulkanDevice struct {
	mu                  sync.RWMutex
	handle              VkDevice
	physicalDevice      *VulkanPhysicalDevice
	instance            *VulkanInstance
	graphicsQueue       VkQueue
	computeQueue        VkQueue
	transferQueue       VkQueue
	presentQueue        VkQueue
	graphicsQueueFamily uint32
	computeQueueFamily  uint32
	transferQueueFamily uint32
	presentQueueFamily  uint32
	memoryAllocator     *VulkanMemoryAllocator
	commandPools        map[uint32]*VulkanCommandPool
	enabledExtensions   []string
	features            VkPhysicalDeviceFeatures
	destroyed           bool
}

// VulkanQueue represents a Vulkan queue
type VulkanQueue struct {
	handle      VkQueue
	familyIndex uint32
	queueIndex  uint32
	flags       uint32
}

// NewVulkanDevice creates a new Vulkan device
func NewVulkanDevice(instance *VulkanInstance, physicalDevice *VulkanPhysicalDevice, extensions []string) (*VulkanDevice, error) {
	device := &VulkanDevice{
		instance:          instance,
		physicalDevice:    physicalDevice,
		enabledExtensions: extensions,
		commandPools:      make(map[uint32]*VulkanCommandPool),
	}

	// Create the logical device
	if err := device.createDevice(); err != nil {
		return nil, err
	}

	// Get device queues
	device.getDeviceQueues()

	// Create memory allocator
	device.memoryAllocator = NewVulkanMemoryAllocator(device)

	return device, nil
}

// createDevice creates the Vulkan logical device
func (vd *VulkanDevice) createDevice() error {
	vd.mu.Lock()
	defer vd.mu.Unlock()

	if vd.destroyed {
		return fmt.Errorf("device has been destroyed")
	}

	// Find queue families
	graphicsFamily, hasGraphics := vd.physicalDevice.GetGraphicsQueueFamilyIndex()
	computeFamily, hasCompute := vd.physicalDevice.GetComputeQueueFamilyIndex()

	if !hasGraphics {
		return fmt.Errorf("no graphics queue family found")
	}

	vd.graphicsQueueFamily = graphicsFamily
	if hasCompute {
		vd.computeQueueFamily = computeFamily
	} else {
		vd.computeQueueFamily = graphicsFamily // Fallback to graphics queue
	}

	// Create queue create infos
	queueFamilies := []uint32{graphicsFamily}
	if hasCompute && computeFamily != graphicsFamily {
		queueFamilies = append(queueFamilies, computeFamily)
	}

	queueCreateInfos := make([]VkDeviceQueueCreateInfo, len(queueFamilies))
	queuePriority := float32(1.0)

	for i, family := range queueFamilies {
		queueCreateInfos[i] = VkDeviceQueueCreateInfo{
			SType:            VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO,
			PNext:            0,
			Flags:            0,
			QueueFamilyIndex: family,
			QueueCount:       1,
			PQueuePriorities: &queuePriority,
		}
	}

	// Convert extension names
	var extensionNames **byte
	if len(vd.enabledExtensions) > 0 {
		extensionNames, _ = cStringArray(vd.enabledExtensions)
	}

	// Device create info
	createInfo := VkDeviceCreateInfo{
		SType:                   VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO,
		PNext:                   0,
		Flags:                   0,
		QueueCreateInfoCount:    uint32(len(queueCreateInfos)),
		PQueueCreateInfos:       &queueCreateInfos[0],
		EnabledLayerCount:       0,
		PpEnabledLayerNames:     nil,
		EnabledExtensionCount:   uint32(len(vd.enabledExtensions)),
		PpEnabledExtensionNames: extensionNames,
		PEnabledFeatures:        &vd.physicalDevice.features,
	}

	// Create device
	result := vkCreateDevice(vd.physicalDevice.handle, &createInfo, 0, &vd.handle)
	if result != VK_SUCCESS {
		return fmt.Errorf("failed to create Vulkan device: %d", result)
	}

	return nil
}

// getDeviceQueues gets device queues
func (vd *VulkanDevice) getDeviceQueues() {
	vkGetDeviceQueue(vd.handle, vd.graphicsQueueFamily, 0, &vd.graphicsQueue)
	vkGetDeviceQueue(vd.handle, vd.computeQueueFamily, 0, &vd.computeQueue)

	// Transfer queue is typically the same as graphics queue
	vd.transferQueue = vd.graphicsQueue
	vd.transferQueueFamily = vd.graphicsQueueFamily

	// Present queue is typically the same as graphics queue
	vd.presentQueue = vd.graphicsQueue
	vd.presentQueueFamily = vd.graphicsQueueFamily
}

// GetHandle returns the Vulkan device handle
func (vd *VulkanDevice) GetHandle() VkDevice {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.handle
}

// GetPhysicalDevice returns the physical device
func (vd *VulkanDevice) GetPhysicalDevice() *VulkanPhysicalDevice {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.physicalDevice
}

// GetInstance returns the Vulkan instance
func (vd *VulkanDevice) GetInstance() *VulkanInstance {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.instance
}

// GetGraphicsQueue returns the graphics queue
func (vd *VulkanDevice) GetGraphicsQueue() VkQueue {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.graphicsQueue
}

// GetComputeQueue returns the compute queue
func (vd *VulkanDevice) GetComputeQueue() VkQueue {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.computeQueue
}

// GetGraphicsQueueFamily returns the graphics queue family index
func (vd *VulkanDevice) GetGraphicsQueueFamily() uint32 {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.graphicsQueueFamily
}

// GetComputeQueueFamily returns the compute queue family index
func (vd *VulkanDevice) GetComputeQueueFamily() uint32 {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.computeQueueFamily
}

// GetMemoryAllocator returns the memory allocator
func (vd *VulkanDevice) GetMemoryAllocator() *VulkanMemoryAllocator {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.memoryAllocator
}

// CreateCommandPool creates a command pool for the specified queue family
func (vd *VulkanDevice) CreateCommandPool(queueFamilyIndex uint32) (*VulkanCommandPool, error) {
	vd.mu.Lock()
	defer vd.mu.Unlock()

	if vd.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	// Check if command pool already exists for this queue family
	if pool, exists := vd.commandPools[queueFamilyIndex]; exists {
		return pool, nil
	}

	pool := &VulkanCommandPool{
		device:           vd,
		queueFamilyIndex: queueFamilyIndex,
		commandBuffers:   make(map[uintptr]*VulkanCommandBuffer),
	}

	// Note: In a real implementation, you would create the VkCommandPool here
	// For now, we'll simulate it
	pool.handle = uintptr(unsafe.Pointer(pool)) // Placeholder

	vd.commandPools[queueFamilyIndex] = pool
	return pool, nil
}

// GetCommandPool gets or creates a command pool for the specified queue family
func (vd *VulkanDevice) GetCommandPool(queueFamilyIndex uint32) (*VulkanCommandPool, error) {
	vd.mu.RLock()
	if pool, exists := vd.commandPools[queueFamilyIndex]; exists {
		vd.mu.RUnlock()
		return pool, nil
	}
	vd.mu.RUnlock()

	// Create new command pool
	return vd.CreateCommandPool(queueFamilyIndex)
}

// WaitIdle waits for the device to become idle
func (vd *VulkanDevice) WaitIdle() error {
	vd.mu.RLock()
	defer vd.mu.RUnlock()

	if vd.destroyed {
		return fmt.Errorf("device has been destroyed")
	}

	// Note: In a real implementation, you would call vkDeviceWaitIdle
	// For now, we'll simulate it
	return nil
}

// Destroy destroys the Vulkan device
func (vd *VulkanDevice) Destroy() error {
	vd.mu.Lock()
	defer vd.mu.Unlock()

	if vd.destroyed {
		return nil
	}

	// Wait for device to be idle
	// vkDeviceWaitIdle(vd.handle)

	// Destroy command pools
	for _, pool := range vd.commandPools {
		pool.Destroy()
	}
	vd.commandPools = nil

	// Destroy memory allocator
	if vd.memoryAllocator != nil {
		vd.memoryAllocator.Destroy()
		vd.memoryAllocator = nil
	}

	// Destroy device
	if vd.handle != 0 {
		vkDestroyDevice(vd.handle, 0)
		vd.handle = 0
	}

	vd.destroyed = true
	return nil
}

// IsDestroyed checks if the device is destroyed
func (vd *VulkanDevice) IsDestroyed() bool {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.destroyed
}
