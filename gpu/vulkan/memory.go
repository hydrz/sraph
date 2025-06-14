package vulkan

import (
	"fmt"
	"sync"
	"unsafe"
)

// VulkanMemoryAllocator manages Vulkan memory allocations
type VulkanMemoryAllocator struct {
	mu             sync.RWMutex
	device         *VulkanDevice
	allocations    map[uintptr]*VulkanMemoryAllocation
	heapUsage      []uint64
	totalAllocated uint64
}

// VulkanMemoryAllocation represents a memory allocation
type VulkanMemoryAllocation struct {
	Handle     uintptr // VkDeviceMemory handle
	Size       uint64
	MemoryType uint32
	Offset     uint64
	MappedPtr  unsafe.Pointer
	IsMapped   bool
	IsCoherent bool
}

// VulkanMemoryRequirements contains memory requirements
type VulkanMemoryRequirements struct {
	Size           uint64
	Alignment      uint64
	MemoryTypeBits uint32
}

// VulkanMemoryAllocateInfo contains allocation parameters
type VulkanMemoryAllocateInfo struct {
	Size            uint64
	MemoryTypeIndex uint32
	Requirements    VulkanMemoryRequirements
}

// NewVulkanMemoryAllocator creates a new memory allocator
func NewVulkanMemoryAllocator(device *VulkanDevice) *VulkanMemoryAllocator {
	allocator := &VulkanMemoryAllocator{
		device:      device,
		allocations: make(map[uintptr]*VulkanMemoryAllocation),
		heapUsage:   make([]uint64, device.physicalDevice.MemoryProperties.MemoryHeapCount),
	}
	return allocator
}

// AllocateMemory allocates device memory using real Vulkan API
func (vma *VulkanMemoryAllocator) AllocateMemory(allocInfo VulkanMemoryAllocateInfo) (*VulkanMemoryAllocation, error) {
	vma.mu.Lock()
	defer vma.mu.Unlock()

	if vma.device.IsDestroyed() {
		return nil, fmt.Errorf("device has been destroyed")
	}

	device := VkDevice(vma.device.handle)

	// Validate memory type index
	if allocInfo.MemoryTypeIndex >= vma.device.physicalDevice.MemoryProperties.MemoryTypeCount {
		return nil, fmt.Errorf("invalid memory type index: %d", allocInfo.MemoryTypeIndex)
	}

	// Check heap capacity
	memoryType := vma.device.physicalDevice.MemoryProperties.MemoryTypes[allocInfo.MemoryTypeIndex]
	heapIndex := memoryType.HeapIndex
	heap := vma.device.physicalDevice.MemoryProperties.MemoryHeaps[heapIndex]

	if vma.heapUsage[heapIndex]+allocInfo.Size > heap.Size {
		return nil, fmt.Errorf("insufficient memory in heap %d", heapIndex)
	}

	// Create Vulkan memory allocate info
	vkAllocInfo := VkMemoryAllocateInfo{
		SType:           VK_STRUCTURE_TYPE_MEMORY_ALLOCATE_INFO,
		PNext:           nil,
		AllocationSize:  allocInfo.Size,
		MemoryTypeIndex: allocInfo.MemoryTypeIndex,
	}

	var memory VkDeviceMemory
	result := vkAllocateMemory(device, &vkAllocInfo, nil, &memory)
	if result != VK_SUCCESS {
		return nil, fmt.Errorf("vkAllocateMemory failed with result %d", result)
	}

	allocation := &VulkanMemoryAllocation{
		Handle:     uintptr(memory),
		Size:       allocInfo.Size,
		MemoryType: allocInfo.MemoryTypeIndex,
		Offset:     0,
		MappedPtr:  nil,
		IsMapped:   false,
		IsCoherent: (memoryType.PropertyFlags & VK_MEMORY_PROPERTY_HOST_COHERENT_BIT) != 0,
	}

	vma.allocations[allocation.Handle] = allocation
	vma.heapUsage[heapIndex] += allocInfo.Size
	vma.totalAllocated += allocInfo.Size

	return allocation, nil
}

// FreeMemory frees device memory using real Vulkan API
func (vma *VulkanMemoryAllocator) FreeMemory(allocation *VulkanMemoryAllocation) error {
	vma.mu.Lock()
	defer vma.mu.Unlock()

	if allocation == nil {
		return fmt.Errorf("allocation is nil")
	}

	if _, exists := vma.allocations[allocation.Handle]; !exists {
		return fmt.Errorf("allocation not found")
	}

	device := VkDevice(vma.device.handle)

	// Unmap if mapped
	if allocation.IsMapped {
		vkUnmapMemory(device, VkDeviceMemory(allocation.Handle))
		allocation.IsMapped = false
		allocation.MappedPtr = nil
	}

	// Update heap usage
	memoryType := vma.device.physicalDevice.MemoryProperties.MemoryTypes[allocation.MemoryType]
	heapIndex := memoryType.HeapIndex
	vma.heapUsage[heapIndex] -= allocation.Size
	vma.totalAllocated -= allocation.Size

	// Free Vulkan memory
	vkFreeMemory(device, VkDeviceMemory(allocation.Handle), nil)
	delete(vma.allocations, allocation.Handle)

	return nil
}

// MapMemory maps device memory to host memory using real Vulkan API
func (vma *VulkanMemoryAllocator) MapMemory(allocation *VulkanMemoryAllocation, offset, size uint64) (unsafe.Pointer, error) {
	vma.mu.Lock()
	defer vma.mu.Unlock()

	if allocation == nil {
		return nil, fmt.Errorf("allocation is nil")
	}

	if allocation.IsMapped {
		return nil, fmt.Errorf("memory is already mapped")
	}

	if offset+size > allocation.Size {
		return nil, fmt.Errorf("map range exceeds allocation size")
	}

	// Check if memory type is host visible
	memoryType := vma.device.physicalDevice.MemoryProperties.MemoryTypes[allocation.MemoryType]
	if (memoryType.PropertyFlags & VK_MEMORY_PROPERTY_HOST_VISIBLE_BIT) == 0 {
		return nil, fmt.Errorf("memory type is not host visible")
	}

	device := VkDevice(vma.device.handle)
	var mappedPtr unsafe.Pointer

	result := vkMapMemory(device, VkDeviceMemory(allocation.Handle), offset, size, 0, &mappedPtr)
	if result != VK_SUCCESS {
		return nil, fmt.Errorf("vkMapMemory failed with result %d", result)
	}

	allocation.MappedPtr = mappedPtr
	allocation.IsMapped = true

	return mappedPtr, nil
}

// UnmapMemory unmaps device memory using real Vulkan API
func (vma *VulkanMemoryAllocator) UnmapMemory(allocation *VulkanMemoryAllocation) error {
	vma.mu.Lock()
	defer vma.mu.Unlock()

	return vma.unmapMemoryLocked(allocation)
}

// unmapMemoryLocked unmaps memory (caller must hold lock)
func (vma *VulkanMemoryAllocator) unmapMemoryLocked(allocation *VulkanMemoryAllocation) error {
	if allocation == nil {
		return fmt.Errorf("allocation is nil")
	}

	if !allocation.IsMapped {
		return fmt.Errorf("memory is not mapped")
	}

	device := VkDevice(vma.device.handle)
	vkUnmapMemory(device, VkDeviceMemory(allocation.Handle))

	allocation.MappedPtr = nil
	allocation.IsMapped = false

	return nil
}

// FlushMappedMemory flushes mapped memory ranges
func (vma *VulkanMemoryAllocator) FlushMappedMemory(allocation *VulkanMemoryAllocation, offset, size uint64) error {
	vma.mu.RLock()
	defer vma.mu.RUnlock()

	if allocation == nil {
		return fmt.Errorf("allocation is nil")
	}

	if !allocation.IsMapped {
		return fmt.Errorf("memory is not mapped")
	}

	// If memory is coherent, no need to flush
	if allocation.IsCoherent {
		return nil
	}

	// In a real implementation, this would call vkFlushMappedMemoryRanges
	return nil
}

// InvalidateMappedMemory invalidates mapped memory ranges
func (vma *VulkanMemoryAllocator) InvalidateMappedMemory(allocation *VulkanMemoryAllocation, offset, size uint64) error {
	vma.mu.RLock()
	defer vma.mu.RUnlock()

	if allocation == nil {
		return fmt.Errorf("allocation is nil")
	}

	if !allocation.IsMapped {
		return fmt.Errorf("memory is not mapped")
	}

	// If memory is coherent, no need to invalidate
	if allocation.IsCoherent {
		return nil
	}

	// In a real implementation, this would call vkInvalidateMappedMemoryRanges
	return nil
}

// FindMemoryType finds suitable memory type
func (vma *VulkanMemoryAllocator) FindMemoryType(typeFilter uint32, properties uint32) (uint32, error) {
	memProps := vma.device.physicalDevice.MemoryProperties

	for i := uint32(0); i < memProps.MemoryTypeCount; i++ {
		if (typeFilter&(1<<i)) != 0 && (memProps.MemoryTypes[i].PropertyFlags&properties) == properties {
			return i, nil
		}
	}

	return 0, fmt.Errorf("failed to find suitable memory type")
}

// GetMemoryUsage returns memory usage statistics
func (vma *VulkanMemoryAllocator) GetMemoryUsage() (totalAllocated uint64, heapUsage []uint64) {
	vma.mu.RLock()
	defer vma.mu.RUnlock()

	// Copy heap usage to avoid data races
	heapUsageCopy := make([]uint64, len(vma.heapUsage))
	copy(heapUsageCopy, vma.heapUsage)

	return vma.totalAllocated, heapUsageCopy
}

// Destroy destroys the memory allocator
func (vma *VulkanMemoryAllocator) Destroy() error {
	vma.mu.Lock()
	defer vma.mu.Unlock()

	// Free all remaining allocations
	for _, allocation := range vma.allocations {
		if allocation.IsMapped {
			vma.unmapMemoryLocked(allocation)
		}
		// In a real implementation, this would call vkFreeMemory
	}

	vma.allocations = nil
	vma.heapUsage = nil
	vma.totalAllocated = 0

	return nil
}
