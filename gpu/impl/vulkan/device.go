package vulkan

import (
	"fmt"
	"sync"
)

// VulkanDevice represents a Vulkan logical device
type VulkanDevice struct {
	mu                sync.RWMutex
	handle            uintptr // VkDevice handle
	physicalDevice    *VulkanPhysicalDevice
	instance          *VulkanInstance
	queues            map[uint32]*VulkanQueue
	commandPools      map[uint32]*VulkanCommandPool
	memoryAllocator   *VulkanMemoryAllocator
	enabledFeatures   VulkanDeviceFeatures
	enabledExtensions []string
	destroyed         bool
}

// VulkanQueue represents a Vulkan queue
type VulkanQueue struct {
	Handle      uintptr // VkQueue handle
	FamilyIndex uint32
	QueueIndex  uint32
	Properties  VulkanQueueFamilyProperties
}

// VulkanDeviceCreateInfo contains device creation parameters
type VulkanDeviceCreateInfo struct {
	PhysicalDevice     *VulkanPhysicalDevice
	Instance           *VulkanInstance
	RequiredFeatures   VulkanDeviceFeatures
	RequiredExtensions []string
}

// NewVulkanDevice creates a new Vulkan logical device
func NewVulkanDevice(createInfo VulkanDeviceCreateInfo) (*VulkanDevice, error) {
	device := &VulkanDevice{
		physicalDevice:    createInfo.PhysicalDevice,
		instance:          createInfo.Instance,
		queues:            make(map[uint32]*VulkanQueue),
		commandPools:      make(map[uint32]*VulkanCommandPool),
		enabledFeatures:   createInfo.RequiredFeatures,
		enabledExtensions: createInfo.RequiredExtensions,
	}

	if err := device.createLogicalDevice(); err != nil {
		return nil, fmt.Errorf("failed to create logical device: %v", err)
	}

	if err := device.createQueues(); err != nil {
		device.Destroy()
		return nil, fmt.Errorf("failed to create queues: %v", err)
	}

	if err := device.createCommandPools(); err != nil {
		device.Destroy()
		return nil, fmt.Errorf("failed to create command pools: %v", err)
	}

	// Initialize memory allocator
	device.memoryAllocator = NewVulkanMemoryAllocator(device)

	return device, nil
}

// createLogicalDevice creates the Vulkan logical device using real Vulkan API
func (vd *VulkanDevice) createLogicalDevice() error {
	physicalDevice := VkPhysicalDevice(vd.physicalDevice.Handle)

	// Create queue create infos for all available queue families
	queueCreateInfos := make([]VkDeviceQueueCreateInfo, len(vd.physicalDevice.QueueFamilyProps))
	queuePriorities := make([][]float32, len(queueCreateInfos))

	for i, queueFamily := range vd.physicalDevice.QueueFamilyProps {
		// Create priority array for this queue family
		queuePriorities[i] = make([]float32, queueFamily.QueueCount)
		for j := range queuePriorities[i] {
			queuePriorities[i][j] = 1.0
		}

		queueCreateInfos[i] = VkDeviceQueueCreateInfo{
			SType:            VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO,
			PNext:            nil,
			Flags:            0,
			QueueFamilyIndex: uint32(i),
			QueueCount:       queueFamily.QueueCount,
			PQueuePriorities: &queuePriorities[i][0],
		}
	}

	// Convert enabled features
	features := VkPhysicalDeviceFeatures{
		GeometryShader:       boolToUint32(vd.enabledFeatures.GeometryShader),
		TessellationShader:   boolToUint32(vd.enabledFeatures.TessellationShader),
		MultiViewport:        boolToUint32(vd.enabledFeatures.MultiViewport),
		SamplerAnisotropy:    boolToUint32(vd.enabledFeatures.SamplerAnisotropy),
		TextureCompressionBC: boolToUint32(vd.enabledFeatures.TextureCompressionBC),
		DepthClamp:           boolToUint32(vd.enabledFeatures.DepthClamp),
	}

	// Convert extension names
	extNames, extCStrs := toCStringArray(vd.enabledExtensions)
	defer func() {
		_ = extCStrs // Keep alive
	}()

	// Create device create info
	createInfo := VkDeviceCreateInfo{
		SType:                   VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO,
		PNext:                   nil,
		Flags:                   0,
		QueueCreateInfoCount:    uint32(len(queueCreateInfos)),
		PQueueCreateInfos:       &queueCreateInfos[0],
		EnabledLayerCount:       0,
		PpEnabledLayerNames:     nil,
		EnabledExtensionCount:   uint32(len(vd.enabledExtensions)),
		PpEnabledExtensionNames: extNames,
		PEnabledFeatures:        &features,
	}

	var device VkDevice
	result := vkCreateDevice(physicalDevice, &createInfo, nil, &device)
	if result != VK_SUCCESS {
		return fmt.Errorf("vkCreateDevice failed with result %d", result)
	}

	vd.handle = uintptr(device)
	return nil
}

// createQueues creates device queues using real Vulkan API
func (vd *VulkanDevice) createQueues() error {
	device := VkDevice(vd.handle)

	for familyIndex, familyProps := range vd.physicalDevice.QueueFamilyProps {
		queueCount := familyProps.QueueCount
		if queueCount > 4 {
			queueCount = 4 // Limit to 4 queues per family
		}

		for queueIndex := uint32(0); queueIndex < queueCount; queueIndex++ {
			var queue VkQueue
			vkGetDeviceQueue(device, uint32(familyIndex), queueIndex, &queue)

			vulkanQueue := &VulkanQueue{
				Handle:      uintptr(queue),
				FamilyIndex: uint32(familyIndex),
				QueueIndex:  queueIndex,
				Properties:  familyProps,
			}

			vd.queues[uint32(familyIndex)*16+queueIndex] = vulkanQueue
		}
	}

	return nil
}

// createCommandPools creates command pools for each queue family
func (vd *VulkanDevice) createCommandPools() error {
	for familyIndex := range vd.physicalDevice.QueueFamilyProps {
		pool, err := NewVulkanCommandPool(vd, uint32(familyIndex))
		if err != nil {
			return fmt.Errorf("failed to create command pool for family %d: %v", familyIndex, err)
		}
		vd.commandPools[uint32(familyIndex)] = pool
	}
	return nil
}

// GetQueue returns a queue for the specified family and index
func (vd *VulkanDevice) GetQueue(familyIndex, queueIndex uint32) *VulkanQueue {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.queues[familyIndex*16+queueIndex]
}

// GetGraphicsQueue returns the first available graphics queue
func (vd *VulkanDevice) GetGraphicsQueue() *VulkanQueue {
	vd.mu.RLock()
	defer vd.mu.RUnlock()

	for _, queue := range vd.queues {
		if queue.Properties.SupportsGraphics {
			return queue
		}
	}
	return nil
}

// GetComputeQueue returns the first available compute queue
func (vd *VulkanDevice) GetComputeQueue() *VulkanQueue {
	vd.mu.RLock()
	defer vd.mu.RUnlock()

	for _, queue := range vd.queues {
		if queue.Properties.SupportsCompute {
			return queue
		}
	}
	return nil
}

// GetTransferQueue returns the first available transfer queue
func (vd *VulkanDevice) GetTransferQueue() *VulkanQueue {
	vd.mu.RLock()
	defer vd.mu.RUnlock()

	for _, queue := range vd.queues {
		if queue.Properties.SupportsTransfer {
			return queue
		}
	}
	return nil
}

// GetCommandPool returns a command pool for the specified queue family
func (vd *VulkanDevice) GetCommandPool(familyIndex uint32) *VulkanCommandPool {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.commandPools[familyIndex]
}

// GetMemoryAllocator returns the memory allocator
func (vd *VulkanDevice) GetMemoryAllocator() *VulkanMemoryAllocator {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.memoryAllocator
}

// WaitIdle waits for the device to become idle using real Vulkan API
func (vd *VulkanDevice) WaitIdle() error {
	vd.mu.RLock()
	defer vd.mu.RUnlock()

	if vd.destroyed {
		return fmt.Errorf("device has been destroyed")
	}

	result := vkDeviceWaitIdle(VkDevice(vd.handle))
	if result != VK_SUCCESS {
		return fmt.Errorf("vkDeviceWaitIdle failed with result %d", result)
	}

	return nil
}

// GetPhysicalDevice returns the physical device
func (vd *VulkanDevice) GetPhysicalDevice() *VulkanPhysicalDevice {
	return vd.physicalDevice
}

// GetInstance returns the Vulkan instance
func (vd *VulkanDevice) GetInstance() *VulkanInstance {
	return vd.instance
}

// GetHandle returns the Vulkan device handle
func (vd *VulkanDevice) GetHandle() uintptr {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.handle
}

// IsDestroyed checks if the device is destroyed
func (vd *VulkanDevice) IsDestroyed() bool {
	vd.mu.RLock()
	defer vd.mu.RUnlock()
	return vd.destroyed
}

// Destroy destroys the Vulkan device using real Vulkan API
func (vd *VulkanDevice) Destroy() error {
	vd.mu.Lock()
	defer vd.mu.Unlock()

	if vd.destroyed {
		return nil
	}

	// Wait for device to become idle
	if vd.handle != 0 {
		vkDeviceWaitIdle(VkDevice(vd.handle))
	}

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
		vkDestroyDevice(VkDevice(vd.handle), nil)
		vd.handle = 0
	}

	vd.queues = nil
	vd.destroyed = true

	return nil
}

// Helper function to convert bool to uint32
func boolToUint32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

// HasFeature checks if a device feature is enabled
func (vd *VulkanDevice) HasFeature(feature string) bool {
	switch feature {
	case "geometry_shader":
		return vd.enabledFeatures.GeometryShader
	case "tessellation_shader":
		return vd.enabledFeatures.TessellationShader
	case "multi_viewport":
		return vd.enabledFeatures.MultiViewport
	case "sampler_anisotropy":
		return vd.enabledFeatures.SamplerAnisotropy
	case "texture_compression_bc":
		return vd.enabledFeatures.TextureCompressionBC
	case "depth_clamp":
		return vd.enabledFeatures.DepthClamp
	default:
		return false
	}
}

// HasExtension checks if a device extension is enabled
func (vd *VulkanDevice) HasExtension(extension string) bool {
	for _, ext := range vd.enabledExtensions {
		if ext == extension {
			return true
		}
	}
	return false
}
