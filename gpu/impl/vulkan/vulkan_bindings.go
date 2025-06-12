package vulkan

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

// Vulkan constants
const (
	VK_SUCCESS                     = 0
	VK_INCOMPLETE                  = 5
	VK_ERROR_OUT_OF_HOST_MEMORY    = -1
	VK_ERROR_OUT_OF_DEVICE_MEMORY  = -2
	VK_ERROR_INITIALIZATION_FAILED = -3
	VK_ERROR_DEVICE_LOST           = -4
	VK_ERROR_MEMORY_MAP_FAILED     = -5
	VK_ERROR_LAYER_NOT_PRESENT     = -6
	VK_ERROR_EXTENSION_NOT_PRESENT = -7
	VK_ERROR_FEATURE_NOT_PRESENT   = -8
	VK_ERROR_INCOMPATIBLE_DRIVER   = -9

	VK_PHYSICAL_DEVICE_TYPE_OTHER          = 0
	VK_PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU = 1
	VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU   = 2
	VK_PHYSICAL_DEVICE_TYPE_VIRTUAL_GPU    = 3
	VK_PHYSICAL_DEVICE_TYPE_CPU            = 4

	VK_QUEUE_GRAPHICS_BIT       = 0x00000001
	VK_QUEUE_COMPUTE_BIT        = 0x00000002
	VK_QUEUE_TRANSFER_BIT       = 0x00000004
	VK_QUEUE_SPARSE_BINDING_BIT = 0x00000008

	VK_MEMORY_PROPERTY_DEVICE_LOCAL_BIT     = 0x00000001
	VK_MEMORY_PROPERTY_HOST_VISIBLE_BIT     = 0x00000002
	VK_MEMORY_PROPERTY_HOST_COHERENT_BIT    = 0x00000004
	VK_MEMORY_PROPERTY_HOST_CACHED_BIT      = 0x00000008
	VK_MEMORY_PROPERTY_LAZILY_ALLOCATED_BIT = 0x00000010

	VK_BUFFER_USAGE_TRANSFER_SRC_BIT         = 0x00000001
	VK_BUFFER_USAGE_TRANSFER_DST_BIT         = 0x00000002
	VK_BUFFER_USAGE_UNIFORM_TEXEL_BUFFER_BIT = 0x00000004
	VK_BUFFER_USAGE_STORAGE_TEXEL_BUFFER_BIT = 0x00000008
	VK_BUFFER_USAGE_UNIFORM_BUFFER_BIT       = 0x00000010
	VK_BUFFER_USAGE_STORAGE_BUFFER_BIT       = 0x00000020
	VK_BUFFER_USAGE_INDEX_BUFFER_BIT         = 0x00000040
	VK_BUFFER_USAGE_VERTEX_BUFFER_BIT        = 0x00000080
	VK_BUFFER_USAGE_INDIRECT_BUFFER_BIT      = 0x00000100

	VK_IMAGE_USAGE_TRANSFER_SRC_BIT             = 0x00000001
	VK_IMAGE_USAGE_TRANSFER_DST_BIT             = 0x00000002
	VK_IMAGE_USAGE_SAMPLED_BIT                  = 0x00000004
	VK_IMAGE_USAGE_STORAGE_BIT                  = 0x00000008
	VK_IMAGE_USAGE_COLOR_ATTACHMENT_BIT         = 0x00000010
	VK_IMAGE_USAGE_DEPTH_STENCIL_ATTACHMENT_BIT = 0x00000020

	VK_STRUCTURE_TYPE_APPLICATION_INFO         = 0
	VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO     = 1
	VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO = 2
	VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO       = 3
	VK_STRUCTURE_TYPE_BUFFER_CREATE_INFO       = 12
	VK_STRUCTURE_TYPE_IMAGE_CREATE_INFO        = 14
	VK_STRUCTURE_TYPE_MEMORY_ALLOCATE_INFO     = 16
)

// Vulkan handles
type VkInstance uintptr
type VkPhysicalDevice uintptr
type VkDevice uintptr
type VkQueue uintptr
type VkCommandPool uintptr
type VkCommandBuffer uintptr
type VkBuffer uintptr
type VkDeviceMemory uintptr
type VkImage uintptr
type VkImageView uintptr
type VkSampler uintptr

// Vulkan structures
type VkApplicationInfo struct {
	SType              uint32
	PNext              unsafe.Pointer
	PApplicationName   *byte
	ApplicationVersion uint32
	PEngineName        *byte
	EngineVersion      uint32
	ApiVersion         uint32
}

type VkInstanceCreateInfo struct {
	SType                   uint32
	PNext                   unsafe.Pointer
	Flags                   uint32
	PApplicationInfo        *VkApplicationInfo
	EnabledLayerCount       uint32
	PpEnabledLayerNames     **byte
	EnabledExtensionCount   uint32
	PpEnabledExtensionNames **byte
}

type VkPhysicalDeviceProperties struct {
	ApiVersion        uint32
	DriverVersion     uint32
	VendorID          uint32
	DeviceID          uint32
	DeviceType        uint32
	DeviceName        [256]byte
	PipelineCacheUUID [16]byte
	Limits            VkPhysicalDeviceLimits
	SparseProperties  VkPhysicalDeviceSparseProperties
}

type VkPhysicalDeviceLimits struct {
	MaxImageDimension1D                             uint32
	MaxImageDimension2D                             uint32
	MaxImageDimension3D                             uint32
	MaxImageDimensionCube                           uint32
	MaxImageArrayLayers                             uint32
	MaxTexelBufferElements                          uint32
	MaxUniformBufferRange                           uint32
	MaxStorageBufferRange                           uint32
	MaxPushConstantsSize                            uint32
	MaxMemoryAllocationCount                        uint32
	MaxSamplerAllocationCount                       uint32
	BufferImageGranularity                          uint64
	SparseAddressSpaceSize                          uint64
	MaxBoundDescriptorSets                          uint32
	MaxPerStageDescriptorSamplers                   uint32
	MaxPerStageDescriptorUniformBuffers             uint32
	MaxPerStageDescriptorStorageBuffers             uint32
	MaxPerStageDescriptorSampledImages              uint32
	MaxPerStageDescriptorStorageImages              uint32
	MaxPerStageDescriptorInputAttachments           uint32
	MaxPerStageResources                            uint32
	MaxDescriptorSetSamplers                        uint32
	MaxDescriptorSetUniformBuffers                  uint32
	MaxDescriptorSetUniformBuffersDynamic           uint32
	MaxDescriptorSetStorageBuffers                  uint32
	MaxDescriptorSetStorageBuffersDynamic           uint32
	MaxDescriptorSetSampledImages                   uint32
	MaxDescriptorSetStorageImages                   uint32
	MaxDescriptorSetInputAttachments                uint32
	MaxVertexInputAttributes                        uint32
	MaxVertexInputBindings                          uint32
	MaxVertexInputAttributeOffset                   uint32
	MaxVertexInputBindingStride                     uint32
	MaxVertexOutputComponents                       uint32
	MaxTessellationGenerationLevel                  uint32
	MaxTessellationPatchSize                        uint32
	MaxTessellationControlPerVertexInputComponents  uint32
	MaxTessellationControlPerVertexOutputComponents uint32
	MaxTessellationControlPerPatchOutputComponents  uint32
	MaxTessellationControlTotalOutputComponents     uint32
	MaxTessellationEvaluationInputComponents        uint32
	MaxTessellationEvaluationOutputComponents       uint32
	MaxGeometryShaderInvocations                    uint32
	MaxGeometryInputComponents                      uint32
	MaxGeometryOutputComponents                     uint32
	MaxGeometryOutputVertices                       uint32
	MaxGeometryTotalOutputComponents                uint32
	MaxFragmentInputComponents                      uint32
	MaxFragmentOutputAttachments                    uint32
	MaxFragmentDualSrcAttachments                   uint32
	MaxFragmentCombinedOutputResources              uint32
	MaxComputeSharedMemorySize                      uint32
	MaxComputeWorkGroupCount                        [3]uint32
	MaxComputeWorkGroupInvocations                  uint32
	MaxComputeWorkGroupSize                         [3]uint32
	SubPixelPrecisionBits                           uint32
	SubTexelPrecisionBits                           uint32
	MipmapPrecisionBits                             uint32
	MaxDrawIndexedIndexValue                        uint32
	MaxDrawIndirectCount                            uint32
	MaxSamplerLodBias                               float32
	MaxSamplerAnisotropy                            float32
	MaxViewports                                    uint32
	MaxViewportDimensions                           [2]uint32
	ViewportBoundsRange                             [2]float32
	ViewportSubPixelBits                            uint32
	MinMemoryMapAlignment                           uintptr
	MinTexelBufferOffsetAlignment                   uint64
	MinUniformBufferOffsetAlignment                 uint64
	MinStorageBufferOffsetAlignment                 uint64
	MinTexelOffset                                  int32
	MaxTexelOffset                                  uint32
	MinTexelGatherOffset                            int32
	MaxTexelGatherOffset                            uint32
	MinInterpolationOffset                          float32
	MaxInterpolationOffset                          float32
	SubPixelInterpolationOffsetBits                 uint32
	MaxFramebufferWidth                             uint32
	MaxFramebufferHeight                            uint32
	MaxFramebufferLayers                            uint32
	FramebufferColorSampleCounts                    uint32
	FramebufferDepthSampleCounts                    uint32
	FramebufferStencilSampleCounts                  uint32
	FramebufferNoAttachmentsSampleCounts            uint32
	MaxColorAttachments                             uint32
	SampledImageColorSampleCounts                   uint32
	SampledImageIntegerSampleCounts                 uint32
	SampledImageDepthSampleCounts                   uint32
	SampledImageStencilSampleCounts                 uint32
	StorageImageSampleCounts                        uint32
	MaxSampleMaskWords                              uint32
	TimestampComputeAndGraphics                     uint32
	TimestampPeriod                                 float32
	MaxClipDistances                                uint32
	MaxCullDistances                                uint32
	MaxCombinedClipAndCullDistances                 uint32
	DiscreteQueuePriorities                         uint32
	PointSizeRange                                  [2]float32
	LineWidthRange                                  [2]float32
	PointSizeGranularity                            float32
	LineWidthGranularity                            float32
	StrictLines                                     uint32
	StandardSampleLocations                         uint32
	OptimalBufferCopyOffsetAlignment                uint64
	OptimalBufferCopyRowPitchAlignment              uint64
	NonCoherentAtomSize                             uint64
}

type VkPhysicalDeviceSparseProperties struct {
	ResidencyStandard2DBlockShape            uint32
	ResidencyStandard2DMultisampleBlockShape uint32
	ResidencyStandard3DBlockShape            uint32
	ResidencyAlignedMipSize                  uint32
	ResidencyNonResidentStrict               uint32
}

type VkQueueFamilyProperties struct {
	QueueFlags                  uint32
	QueueCount                  uint32
	TimestampValidBits          uint32
	MinImageTransferGranularity VkExtent3D
}

type VkExtent3D struct {
	Width  uint32
	Height uint32
	Depth  uint32
}

type VkPhysicalDeviceMemoryProperties struct {
	MemoryTypeCount uint32
	MemoryTypes     [32]VkMemoryType
	MemoryHeapCount uint32
	MemoryHeaps     [16]VkMemoryHeap
}

type VkMemoryType struct {
	PropertyFlags uint32
	HeapIndex     uint32
}

type VkMemoryHeap struct {
	Size  uint64
	Flags uint32
}

type VkDeviceQueueCreateInfo struct {
	SType            uint32
	PNext            unsafe.Pointer
	Flags            uint32
	QueueFamilyIndex uint32
	QueueCount       uint32
	PQueuePriorities *float32
}

type VkDeviceCreateInfo struct {
	SType                   uint32
	PNext                   unsafe.Pointer
	Flags                   uint32
	QueueCreateInfoCount    uint32
	PQueueCreateInfos       *VkDeviceQueueCreateInfo
	EnabledLayerCount       uint32
	PpEnabledLayerNames     **byte
	EnabledExtensionCount   uint32
	PpEnabledExtensionNames **byte
	PEnabledFeatures        *VkPhysicalDeviceFeatures
}

type VkPhysicalDeviceFeatures struct {
	RobustBufferAccess                      uint32
	FullDrawIndexUint32                     uint32
	ImageCubeArray                          uint32
	IndependentBlend                        uint32
	GeometryShader                          uint32
	TessellationShader                      uint32
	SampleRateShading                       uint32
	DualSrcBlend                            uint32
	LogicOp                                 uint32
	MultiDrawIndirect                       uint32
	DrawIndirectFirstInstance               uint32
	DepthClamp                              uint32
	DepthBiasClamp                          uint32
	FillModeNonSolid                        uint32
	DepthBounds                             uint32
	WideLines                               uint32
	LargePoints                             uint32
	AlphaToOne                              uint32
	MultiViewport                           uint32
	SamplerAnisotropy                       uint32
	TextureCompressionETC2                  uint32
	TextureCompressionASTC_LDR              uint32
	TextureCompressionBC                    uint32
	OcclusionQueryPrecise                   uint32
	PipelineStatisticsQuery                 uint32
	VertexPipelineStoresAndAtomics          uint32
	FragmentStoresAndAtomics                uint32
	ShaderTessellationAndGeometryPointSize  uint32
	ShaderImageGatherExtended               uint32
	ShaderStorageImageExtendedFormats       uint32
	ShaderStorageImageMultisample           uint32
	ShaderStorageImageReadWithoutFormat     uint32
	ShaderStorageImageWriteWithoutFormat    uint32
	ShaderUniformBufferArrayDynamicIndexing uint32
	ShaderSampledImageArrayDynamicIndexing  uint32
	ShaderStorageBufferArrayDynamicIndexing uint32
	ShaderStorageImageArrayDynamicIndexing  uint32
	ShaderClipDistance                      uint32
	ShaderCullDistance                      uint32
	ShaderFloat64                           uint32
	ShaderInt64                             uint32
	ShaderInt16                             uint32
	ShaderResourceResidency                 uint32
	ShaderResourceMinLod                    uint32
	SparseBinding                           uint32
	SparseResidencyBuffer                   uint32
	SparseResidencyImage2D                  uint32
	SparseResidencyImage3D                  uint32
	SparseResidency2Samples                 uint32
	SparseResidency4Samples                 uint32
	SparseResidency8Samples                 uint32
	SparseResidency16Samples                uint32
	SparseResidencyAliased                  uint32
	VariableMultisampleRate                 uint32
	InheritedQueries                        uint32
}

type VkBufferCreateInfo struct {
	SType                 uint32
	PNext                 unsafe.Pointer
	Flags                 uint32
	Size                  uint64
	Usage                 uint32
	SharingMode           uint32
	QueueFamilyIndexCount uint32
	PQueueFamilyIndices   *uint32
}

type VkMemoryRequirements struct {
	Size           uint64
	Alignment      uint64
	MemoryTypeBits uint32
}

type VkMemoryAllocateInfo struct {
	SType           uint32
	PNext           unsafe.Pointer
	AllocationSize  uint64
	MemoryTypeIndex uint32
}

// Vulkan function pointers
var (
	vulkanLib                                uintptr
	vkCreateInstance                         func(pCreateInfo *VkInstanceCreateInfo, pAllocator unsafe.Pointer, pInstance *VkInstance) int32
	vkDestroyInstance                        func(instance VkInstance, pAllocator unsafe.Pointer)
	vkEnumeratePhysicalDevices               func(instance VkInstance, pPhysicalDeviceCount *uint32, pPhysicalDevices *VkPhysicalDevice) int32
	vkGetPhysicalDeviceProperties            func(physicalDevice VkPhysicalDevice, pProperties *VkPhysicalDeviceProperties)
	vkGetPhysicalDeviceFeatures              func(physicalDevice VkPhysicalDevice, pFeatures *VkPhysicalDeviceFeatures)
	vkGetPhysicalDeviceQueueFamilyProperties func(physicalDevice VkPhysicalDevice, pQueueFamilyPropertyCount *uint32, pQueueFamilyProperties *VkQueueFamilyProperties)
	vkGetPhysicalDeviceMemoryProperties      func(physicalDevice VkPhysicalDevice, pMemoryProperties *VkPhysicalDeviceMemoryProperties)
	vkCreateDevice                           func(physicalDevice VkPhysicalDevice, pCreateInfo *VkDeviceCreateInfo, pAllocator unsafe.Pointer, pDevice *VkDevice) int32
	vkDestroyDevice                          func(device VkDevice, pAllocator unsafe.Pointer)
	vkGetDeviceQueue                         func(device VkDevice, queueFamilyIndex uint32, queueIndex uint32, pQueue *VkQueue)
	vkDeviceWaitIdle                         func(device VkDevice) int32
	vkCreateBuffer                           func(device VkDevice, pCreateInfo *VkBufferCreateInfo, pAllocator unsafe.Pointer, pBuffer *VkBuffer) int32
	vkDestroyBuffer                          func(device VkDevice, buffer VkBuffer, pAllocator unsafe.Pointer)
	vkGetBufferMemoryRequirements            func(device VkDevice, buffer VkBuffer, pMemoryRequirements *VkMemoryRequirements)
	vkAllocateMemory                         func(device VkDevice, pAllocateInfo *VkMemoryAllocateInfo, pAllocator unsafe.Pointer, pMemory *VkDeviceMemory) int32
	vkFreeMemory                             func(device VkDevice, memory VkDeviceMemory, pAllocator unsafe.Pointer)
	vkBindBufferMemory                       func(device VkDevice, buffer VkBuffer, memory VkDeviceMemory, memoryOffset uint64) int32
	vkMapMemory                              func(device VkDevice, memory VkDeviceMemory, offset uint64, size uint64, flags uint32, ppData *unsafe.Pointer) int32
	vkUnmapMemory                            func(device VkDevice, memory VkDeviceMemory)
)

// LoadVulkanLibrary loads the Vulkan library and function pointers
func LoadVulkanLibrary() error {
	var libName string
	switch {
	case isWindows():
		libName = "vulkan-1.dll"
	case isMacOS():
		libName = "libvulkan.1.dylib"
	default:
		libName = "libvulkan.so.1"
	}

	lib, err := purego.Dlopen(libName, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load Vulkan library %s: %v", libName, err)
	}
	vulkanLib = lib

	// Load function pointers
	if err := loadVulkanFunctions(); err != nil {
		return fmt.Errorf("failed to load Vulkan functions: %v", err)
	}

	return nil
}

// loadVulkanFunctions loads all required Vulkan function pointers
func loadVulkanFunctions() error {
	funcMap := map[string]interface{}{
		"vkCreateInstance":                         &vkCreateInstance,
		"vkDestroyInstance":                        &vkDestroyInstance,
		"vkEnumeratePhysicalDevices":               &vkEnumeratePhysicalDevices,
		"vkGetPhysicalDeviceProperties":            &vkGetPhysicalDeviceProperties,
		"vkGetPhysicalDeviceFeatures":              &vkGetPhysicalDeviceFeatures,
		"vkGetPhysicalDeviceQueueFamilyProperties": &vkGetPhysicalDeviceQueueFamilyProperties,
		"vkGetPhysicalDeviceMemoryProperties":      &vkGetPhysicalDeviceMemoryProperties,
		"vkCreateDevice":                           &vkCreateDevice,
		"vkDestroyDevice":                          &vkDestroyDevice,
		"vkGetDeviceQueue":                         &vkGetDeviceQueue,
		"vkDeviceWaitIdle":                         &vkDeviceWaitIdle,
		"vkCreateBuffer":                           &vkCreateBuffer,
		"vkDestroyBuffer":                          &vkDestroyBuffer,
		"vkGetBufferMemoryRequirements":            &vkGetBufferMemoryRequirements,
		"vkAllocateMemory":                         &vkAllocateMemory,
		"vkFreeMemory":                             &vkFreeMemory,
		"vkBindBufferMemory":                       &vkBindBufferMemory,
		"vkMapMemory":                              &vkMapMemory,
		"vkUnmapMemory":                            &vkUnmapMemory,
	}

	for name, funcPtr := range funcMap {
		sym, err := purego.Dlsym(vulkanLib, name)
		if err != nil {
			return fmt.Errorf("failed to load function %s: %v", name, err)
		}
		purego.RegisterFunc(funcPtr, sym)
	}

	return nil
}

// UnloadVulkanLibrary unloads the Vulkan library
func UnloadVulkanLibrary() error {
	if vulkanLib != 0 {
		return purego.Dlclose(vulkanLib)
	}
	return nil
}

// Platform detection functions
func isWindows() bool {
	// Implementation specific to Go runtime
	return false // Simplified for example
}

func isMacOS() bool {
	// Implementation specific to Go runtime
	return false // Simplified for example
}

// Helper function to convert Go string to C string
func toCString(s string) *byte {
	if s == "" {
		return nil
	}
	b := make([]byte, len(s)+1)
	copy(b, s)
	return &b[0]
}

// Helper function to convert Go string slice to C string array
func toCStringArray(strs []string) (**byte, [][]byte) {
	if len(strs) == 0 {
		return nil, nil
	}

	cstrs := make([][]byte, len(strs))
	ptrs := make([]*byte, len(strs))

	for i, s := range strs {
		cstrs[i] = make([]byte, len(s)+1)
		copy(cstrs[i], s)
		ptrs[i] = &cstrs[i][0]
	}

	return &ptrs[0], cstrs
}
