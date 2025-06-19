package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// AllocatorVK provides memory allocation functionality
// It manages Vulkan device memory allocation and deallocation
type AllocatorVK struct {
	device         vulkan.Device
	physicalDevice vulkan.PhysicalDevice
	memoryProps    vulkan.PhysicalDeviceMemoryProperties
	allocations    []*AllocationVK
	totalAllocated vulkan.DeviceSize
}

// AllocationVK represents a memory allocation
type AllocationVK struct {
	memory     vulkan.DeviceMemory
	size       vulkan.DeviceSize
	memoryType uint32
	offset     vulkan.DeviceSize
	mapped     bool
	mappedPtr  uintptr
}

// NewAllocatorVK creates a new Vulkan allocator
func NewAllocatorVK(device vulkan.Device, physicalDevice vulkan.PhysicalDevice) *AllocatorVK {
	var memProps vulkan.PhysicalDeviceMemoryProperties
	vulkan.GetPhysicalDeviceMemoryProperties(physicalDevice, &memProps)

	return &AllocatorVK{
		device:         device,
		physicalDevice: physicalDevice,
		memoryProps:    memProps,
		allocations:    make([]*AllocationVK, 0),
		totalAllocated: 0,
	}
}

// AllocateMemory allocates device memory
func (a *AllocatorVK) AllocateMemory(size vulkan.DeviceSize, memoryTypeBits uint32, properties vulkan.MemoryPropertyFlags) (*AllocationVK, error) {
	memoryType, err := a.findMemoryType(memoryTypeBits, properties)
	if err != nil {
		return nil, err
	}

	allocInfo := vulkan.MemoryAllocateInfo{
		SType:           vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  size,
		MemoryTypeIndex: memoryType,
	}

	var memory vulkan.DeviceMemory
	if result := vulkan.AllocateMemory(a.device, &allocInfo, nil, &memory); result != vulkan.Success {
		return nil, fmt.Errorf("failed to allocate memory: %s", result)
	}

	allocation := &AllocationVK{
		memory:     memory,
		size:       size,
		memoryType: memoryType,
		offset:     0,
		mapped:     false,
	}

	a.allocations = append(a.allocations, allocation)
	a.totalAllocated += size

	return allocation, nil
}

// FreeMemory frees device memory
func (a *AllocatorVK) FreeMemory(allocation *AllocationVK) {
	if allocation == nil {
		return
	}

	if allocation.mapped {
		a.UnmapMemory(allocation)
	}

	if allocation.memory != vulkan.NullDeviceMemory {
		vulkan.FreeMemory(a.device, allocation.memory, nil)
		a.totalAllocated -= allocation.size
	}

	// Remove from allocations list
	for i, alloc := range a.allocations {
		if alloc == allocation {
			a.allocations = append(a.allocations[:i], a.allocations[i+1:]...)
			break
		}
	}
}

// MapMemory maps memory for CPU access
func (a *AllocatorVK) MapMemory(allocation *AllocationVK) (uintptr, error) {
	if allocation.mapped {
		return allocation.mappedPtr, nil
	}

	var data uintptr
	if result := vulkan.MapMemory(a.device, allocation.memory, allocation.offset, allocation.size, 0, &data); result != vulkan.Success {
		return 0, fmt.Errorf("failed to map memory: %s", result)
	}

	allocation.mapped = true
	allocation.mappedPtr = data

	return data, nil
}

// UnmapMemory unmaps memory
func (a *AllocatorVK) UnmapMemory(allocation *AllocationVK) {
	if allocation.mapped {
		vulkan.UnmapMemory(a.device, allocation.memory)
		allocation.mapped = false
		allocation.mappedPtr = 0
	}
}

// FlushMemory flushes mapped memory
func (a *AllocatorVK) FlushMemory(allocation *AllocationVK, offset, size vulkan.DeviceSize) error {
	mappedRange := vulkan.MappedMemoryRange{
		SType:  vulkan.StructureTypeMappedMemoryRange,
		Memory: allocation.memory,
		Offset: offset,
		Size:   size,
	}

	if result := vulkan.FlushMappedMemoryRanges(a.device, 1, []vulkan.MappedMemoryRange{mappedRange}); result != vulkan.Success {
		return fmt.Errorf("failed to flush memory: %s", result)
	}

	return nil
}

// InvalidateMemory invalidates mapped memory
func (a *AllocatorVK) InvalidateMemory(allocation *AllocationVK, offset, size vulkan.DeviceSize) error {
	mappedRange := vulkan.MappedMemoryRange{
		SType:  vulkan.StructureTypeMappedMemoryRange,
		Memory: allocation.memory,
		Offset: offset,
		Size:   size,
	}

	if result := vulkan.InvalidateMappedMemoryRanges(a.device, 1, []vulkan.MappedMemoryRange{mappedRange}); result != vulkan.Success {
		return fmt.Errorf("failed to invalidate memory: %s", result)
	}

	return nil
}

// GetMemoryRequirements gets memory requirements for a buffer
func (a *AllocatorVK) GetBufferMemoryRequirements(buffer vulkan.Buffer) vulkan.MemoryRequirements {
	var memReqs vulkan.MemoryRequirements
	vulkan.GetBufferMemoryRequirements(a.device, buffer, &memReqs)
	return memReqs
}

// GetImageMemoryRequirements gets memory requirements for an image
func (a *AllocatorVK) GetImageMemoryRequirements(image vulkan.Image) vulkan.MemoryRequirements {
	var memReqs vulkan.MemoryRequirements
	vulkan.GetImageMemoryRequirements(a.device, image, &memReqs)
	return memReqs
}

// BindBufferMemory binds buffer to memory
func (a *AllocatorVK) BindBufferMemory(buffer vulkan.Buffer, allocation *AllocationVK) error {
	if result := vulkan.BindBufferMemory(a.device, buffer, allocation.memory, allocation.offset); result != vulkan.Success {
		return fmt.Errorf("failed to bind buffer memory: %s", result)
	}
	return nil
}

// BindImageMemory binds image to memory
func (a *AllocatorVK) BindImageMemory(image vulkan.Image, allocation *AllocationVK) error {
	if result := vulkan.BindImageMemory(a.device, image, allocation.memory, allocation.offset); result != vulkan.Success {
		return fmt.Errorf("failed to bind image memory: %s", result)
	}
	return nil
}

// GetTotalAllocated returns the total allocated memory size
func (a *AllocatorVK) GetTotalAllocated() vulkan.DeviceSize {
	return a.totalAllocated
}

// GetAllocationCount returns the number of active allocations
func (a *AllocatorVK) GetAllocationCount() int {
	return len(a.allocations)
}

// findMemoryType finds a suitable memory type
func (a *AllocatorVK) findMemoryType(typeFilter uint32, properties vulkan.MemoryPropertyFlags) (uint32, error) {
	memProperties := a.memoryProps

	for i := uint32(0); i < memProperties.MemoryTypeCount; i++ {
		if (typeFilter&(1<<i)) != 0 && (memProperties.MemoryTypes[i].PropertyFlags&properties) == properties {
			return i, nil
		}
	}

	return 0, fmt.Errorf("failed to find suitable memory type")
}

// Destroy destroys the allocator and frees all allocations
func (a *AllocatorVK) Destroy() {
	// Free all remaining allocations
	for len(a.allocations) > 0 {
		a.FreeMemory(a.allocations[0])
	}
}

// GetMemory returns the Vulkan device memory handle
func (alloc *AllocationVK) GetMemory() vulkan.DeviceMemory {
	return alloc.memory
}

// GetSize returns the allocation size
func (alloc *AllocationVK) GetSize() vulkan.DeviceSize {
	return alloc.size
}

// GetOffset returns the allocation offset
func (alloc *AllocationVK) GetOffset() vulkan.DeviceSize {
	return alloc.offset
}

// GetMemoryType returns the memory type index
func (alloc *AllocationVK) GetMemoryType() uint32 {
	return alloc.memoryType
}

// IsMapped returns whether the allocation is mapped
func (alloc *AllocationVK) IsMapped() bool {
	return alloc.mapped
}

// GetMappedPtr returns the mapped pointer
func (alloc *AllocationVK) GetMappedPtr() uintptr {
	return alloc.mappedPtr
}
