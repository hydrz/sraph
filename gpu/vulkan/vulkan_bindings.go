package vulkan

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

// Vulkan constants
const (
	VK_TRUE  = 1
	VK_FALSE = 0

	// API Version
	VK_API_VERSION_1_0 = (1 << 22) | (0 << 12) | (0)
	VK_API_VERSION_1_1 = (1 << 22) | (1 << 12) | (0)
	VK_API_VERSION_1_2 = (1 << 22) | (2 << 12) | (0)
	VK_API_VERSION_1_3 = (1 << 22) | (3 << 12) | (0)

	// Result codes
	VK_SUCCESS                     = 0
	VK_NOT_READY                   = 1
	VK_TIMEOUT                     = 2
	VK_EVENT_SET                   = 3
	VK_EVENT_RESET                 = 4
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
	VK_ERROR_TOO_MANY_OBJECTS      = -10
	VK_ERROR_FORMAT_NOT_SUPPORTED  = -11

	// Physical device types
	VK_PHYSICAL_DEVICE_TYPE_OTHER          = 0
	VK_PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU = 1
	VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU   = 2
	VK_PHYSICAL_DEVICE_TYPE_VIRTUAL_GPU    = 3
	VK_PHYSICAL_DEVICE_TYPE_CPU            = 4

	// Queue flags
	VK_QUEUE_GRAPHICS_BIT       = 0x00000001
	VK_QUEUE_COMPUTE_BIT        = 0x00000002
	VK_QUEUE_TRANSFER_BIT       = 0x00000004
	VK_QUEUE_SPARSE_BINDING_BIT = 0x00000008

	// Memory property flags
	VK_MEMORY_PROPERTY_DEVICE_LOCAL_BIT     = 0x00000001
	VK_MEMORY_PROPERTY_HOST_VISIBLE_BIT     = 0x00000002
	VK_MEMORY_PROPERTY_HOST_COHERENT_BIT    = 0x00000004
	VK_MEMORY_PROPERTY_HOST_CACHED_BIT      = 0x00000008
	VK_MEMORY_PROPERTY_LAZILY_ALLOCATED_BIT = 0x00000010
	VK_MEMORY_PROPERTY_PROTECTED_BIT        = 0x00000020

	// Structure types
	VK_STRUCTURE_TYPE_APPLICATION_INFO         = 0
	VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO     = 1
	VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO = 2
	VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO       = 3
	VK_STRUCTURE_TYPE_SUBMIT_INFO              = 4
	VK_STRUCTURE_TYPE_MEMORY_ALLOCATE_INFO     = 5
	VK_STRUCTURE_TYPE_MAPPED_MEMORY_RANGE      = 6
	VK_STRUCTURE_TYPE_BIND_SPARSE_INFO         = 7
	VK_STRUCTURE_TYPE_FENCE_CREATE_INFO        = 8
	VK_STRUCTURE_TYPE_SEMAPHORE_CREATE_INFO    = 9
	VK_STRUCTURE_TYPE_EVENT_CREATE_INFO        = 10
	VK_STRUCTURE_TYPE_QUERY_POOL_CREATE_INFO   = 11
	VK_STRUCTURE_TYPE_BUFFER_CREATE_INFO       = 12
	VK_STRUCTURE_TYPE_BUFFER_VIEW_CREATE_INFO  = 13
	VK_STRUCTURE_TYPE_IMAGE_CREATE_INFO        = 14
	VK_STRUCTURE_TYPE_IMAGE_VIEW_CREATE_INFO   = 15

	// Default values
	VK_WHOLE_SIZE           = ^uint64(0)
	VK_ATTACHMENT_UNUSED    = ^uint32(0)
	VK_QUEUE_FAMILY_IGNORED = ^uint32(0)
	VK_SUBPASS_EXTERNAL     = ^uint32(0)
)

// Vulkan handles
type VkInstance uintptr
type VkPhysicalDevice uintptr
type VkDevice uintptr
type VkQueue uintptr
type VkSemaphore uintptr
type VkCommandBuffer uintptr
type VkFence uintptr
type VkDeviceMemory uintptr
type VkBuffer uintptr
type VkImage uintptr
type VkEvent uintptr
type VkQueryPool uintptr
type VkBufferView uintptr
type VkImageView uintptr
type VkShaderModule uintptr
type VkPipelineCache uintptr
type VkPipelineLayout uintptr
type VkRenderPass uintptr
type VkPipeline uintptr
type VkDescriptorSetLayout uintptr
type VkSampler uintptr
type VkDescriptorPool uintptr
type VkDescriptorSet uintptr
type VkFramebuffer uintptr
type VkCommandPool uintptr
type VkSurfaceKHR uintptr
type VkSwapchainKHR uintptr
type VkDebugUtilsMessengerEXT uintptr

// Vulkan structures
type VkApplicationInfo struct {
	SType              uint32
	PNext              uintptr
	PApplicationName   *byte
	ApplicationVersion uint32
	PEngineName        *byte
	EngineVersion      uint32
	ApiVersion         uint32
}

type VkInstanceCreateInfo struct {
	SType                   uint32
	PNext                   uintptr
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
	PNext            uintptr
	Flags            uint32
	QueueFamilyIndex uint32
	QueueCount       uint32
	PQueuePriorities *float32
}

type VkDeviceCreateInfo struct {
	SType                   uint32
	PNext                   uintptr
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

// Vulkan function pointers
var (
	vulkanLib uintptr

	// Instance functions
	vkCreateInstance                         func(pCreateInfo *VkInstanceCreateInfo, pAllocator uintptr, pInstance *VkInstance) int32
	vkDestroyInstance                        func(instance VkInstance, pAllocator uintptr)
	vkEnumeratePhysicalDevices               func(instance VkInstance, pPhysicalDeviceCount *uint32, pPhysicalDevices *VkPhysicalDevice) int32
	vkGetPhysicalDeviceProperties            func(physicalDevice VkPhysicalDevice, pProperties *VkPhysicalDeviceProperties)
	vkGetPhysicalDeviceFeatures              func(physicalDevice VkPhysicalDevice, pFeatures *VkPhysicalDeviceFeatures)
	vkGetPhysicalDeviceQueueFamilyProperties func(physicalDevice VkPhysicalDevice, pQueueFamilyPropertyCount *uint32, pQueueFamilyProperties *VkQueueFamilyProperties)
	vkGetPhysicalDeviceMemoryProperties      func(physicalDevice VkPhysicalDevice, pMemoryProperties *VkPhysicalDeviceMemoryProperties)

	// Device functions
	vkCreateDevice   func(physicalDevice VkPhysicalDevice, pCreateInfo *VkDeviceCreateInfo, pAllocator uintptr, pDevice *VkDevice) int32
	vkDestroyDevice  func(device VkDevice, pAllocator uintptr)
	vkGetDeviceQueue func(device VkDevice, queueFamilyIndex uint32, queueIndex uint32, pQueue *VkQueue)

	// Memory functions
	vkAllocateMemory func(device VkDevice, pAllocateInfo uintptr, pAllocator uintptr, pMemory *VkDeviceMemory) int32
	vkFreeMemory     func(device VkDevice, memory VkDeviceMemory, pAllocator uintptr)
	vkMapMemory      func(device VkDevice, memory VkDeviceMemory, offset uint64, size uint64, flags uint32, ppData *unsafe.Pointer) int32
	vkUnmapMemory    func(device VkDevice, memory VkDeviceMemory)
)

// LoadVulkanLibrary loads the Vulkan library and function pointers
func LoadVulkanLibrary() error {
	var err error

	// Determine library name based on platform
	libName := getVulkanLibraryName()

	vulkanLib, err = purego.Dlopen(libName, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load Vulkan library %s: %v", libName, err)
	}

	return loadVulkanFunctions()
}

// loadVulkanFunctions loads all required Vulkan function pointers
func loadVulkanFunctions() error {
	var err error

	// Load global functions directly from library
	vkCreateInstanceAddr, err := purego.Dlsym(vulkanLib, "vkCreateInstance")
	if err != nil {
		return fmt.Errorf("failed to load vkCreateInstance: %v", err)
	}
	purego.RegisterFunc(&vkCreateInstance, vkCreateInstanceAddr)

	// Note: Other functions will be loaded after instance creation using vkGetInstanceProcAddr
	// This is a simplified version - in practice, you'd load functions dynamically

	return nil
}

// UnloadVulkanLibrary unloads the Vulkan library
func UnloadVulkanLibrary() error {
	if vulkanLib != 0 {
		err := purego.Dlclose(vulkanLib)
		vulkanLib = 0
		return err
	}
	return nil
}

// Platform detection functions
func getVulkanLibraryName() string {
	switch runtime.GOOS {
	case "windows":
		return "vulkan-1.dll"
	case "darwin":
		return "libvulkan.1.dylib"
	default:
		return "libvulkan.so.1"
	}
}

// Helper function to convert Go string to C string
func cString(s string) *byte {
	if s == "" {
		return nil
	}
	b := make([]byte, len(s)+1)
	copy(b, s)
	return &b[0]
}

// Helper function to convert Go string slice to C string array
func cStringArray(strs []string) (**byte, []uintptr) {
	if len(strs) == 0 {
		return nil, nil
	}

	ptrs := make([]uintptr, len(strs))
	for i, s := range strs {
		ptrs[i] = uintptr(unsafe.Pointer(cString(s)))
	}

	return (**byte)(unsafe.Pointer(&ptrs[0])), ptrs
}
