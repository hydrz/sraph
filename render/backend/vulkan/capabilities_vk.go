package vulkan

import (
	"github.com/opensraph/sraph/gpu"
	"github.com/opensraph/sraph/render"
	"github.com/vulkan-go/vulkan"
)

// CapabilitiesVK represents Vulkan device capabilities
type CapabilitiesVK struct {
	device                       vulkan.Device
	physicalDevice               vulkan.PhysicalDevice
	deviceProperties             vulkan.PhysicalDeviceProperties
	deviceFeatures               vulkan.PhysicalDeviceFeatures
	memoryProperties             vulkan.PhysicalDeviceMemoryProperties
	queueFamilyProperties        []vulkan.QueueFamilyProperties
	supportedExtensions          []vulkan.ExtensionProperties
	supportedLayers              []vulkan.LayerProperties
	maxTextureSize               uint32
	maxRenderTargetSize          uint32
	maxVertexAttributes          uint32
	maxUniformBufferRange        uint32
	maxStorageBufferRange        uint32
	maxPushConstantSize          uint32
	maxDescriptorSets            uint32
	supportsGeometryShaders      bool
	supportsTessellation         bool
	supportsComputeShaders       bool
	supportsMultisampling        bool
	supportsAnisotropicFiltering bool
}

// NewCapabilitiesVK creates a new Vulkan capabilities object
func NewCapabilitiesVK(device vulkan.Device, physicalDevice vulkan.PhysicalDevice) *CapabilitiesVK {
	caps := &CapabilitiesVK{
		device:         device,
		physicalDevice: physicalDevice,
	}

	caps.queryCapabilities()
	return caps
}

// SupportsBackendType checks if the backend type is supported
func (c *CapabilitiesVK) SupportsBackendType(backendType gpu.BackendType) bool {
	return backendType == gpu.BackendTypeVulkan
}

// SupportsTextureFormat checks if a texture format is supported
func (c *CapabilitiesVK) SupportsTextureFormat(format render.PixelFormat) bool {
	vkFormat := c.convertPixelFormatToVulkan(format)
	if vkFormat == vulkan.FormatUndefined {
		return false
	}

	var formatProperties vulkan.FormatProperties
	vulkan.GetPhysicalDeviceFormatProperties(c.physicalDevice, vkFormat, &formatProperties)

	return formatProperties.OptimalTilingFeatures != 0
}

// GetMaxTextureSize returns the maximum texture size
func (c *CapabilitiesVK) GetMaxTextureSize() uint32 {
	return c.maxTextureSize
}

// GetMaxRenderTargetSize returns the maximum render target size
func (c *CapabilitiesVK) GetMaxRenderTargetSize() uint32 {
	return c.maxRenderTargetSize
}

// GetMaxVertexAttributes returns the maximum number of vertex attributes
func (c *CapabilitiesVK) GetMaxVertexAttributes() uint32 {
	return c.maxVertexAttributes
}

// GetMaxUniformBufferRange returns the maximum uniform buffer range
func (c *CapabilitiesVK) GetMaxUniformBufferRange() uint32 {
	return c.maxUniformBufferRange
}

// GetMaxStorageBufferRange returns the maximum storage buffer range
func (c *CapabilitiesVK) GetMaxStorageBufferRange() uint32 {
	return c.maxStorageBufferRange
}

// GetMaxPushConstantSize returns the maximum push constant size
func (c *CapabilitiesVK) GetMaxPushConstantSize() uint32 {
	return c.maxPushConstantSize
}

// SupportsGeometryShaders returns whether geometry shaders are supported
func (c *CapabilitiesVK) SupportsGeometryShaders() bool {
	return c.supportsGeometryShaders
}

// SupportsTessellation returns whether tessellation is supported
func (c *CapabilitiesVK) SupportsTessellation() bool {
	return c.supportsTessellation
}

// SupportsComputeShaders returns whether compute shaders are supported
func (c *CapabilitiesVK) SupportsComputeShaders() bool {
	return c.supportsComputeShaders
}

// SupportsMultisampling returns whether multisampling is supported
func (c *CapabilitiesVK) SupportsMultisampling() bool {
	return c.supportsMultisampling
}

// SupportsAnisotropicFiltering returns whether anisotropic filtering is supported
func (c *CapabilitiesVK) SupportsAnisotropicFiltering() bool {
	return c.supportsAnisotropicFiltering
}

// GetDeviceProperties returns the Vulkan device properties
func (c *CapabilitiesVK) GetDeviceProperties() vulkan.PhysicalDeviceProperties {
	return c.deviceProperties
}

// GetDeviceFeatures returns the Vulkan device features
func (c *CapabilitiesVK) GetDeviceFeatures() vulkan.PhysicalDeviceFeatures {
	return c.deviceFeatures
}

// GetMemoryProperties returns the Vulkan memory properties
func (c *CapabilitiesVK) GetMemoryProperties() vulkan.PhysicalDeviceMemoryProperties {
	return c.memoryProperties
}

// GetQueueFamilyProperties returns the queue family properties
func (c *CapabilitiesVK) GetQueueFamilyProperties() []vulkan.QueueFamilyProperties {
	return c.queueFamilyProperties
}

// queryCapabilities queries and caches device capabilities
func (c *CapabilitiesVK) queryCapabilities() {
	// Get device properties
	vulkan.GetPhysicalDeviceProperties(c.physicalDevice, &c.deviceProperties)

	// Get device features
	vulkan.GetPhysicalDeviceFeatures(c.physicalDevice, &c.deviceFeatures)

	// Get memory properties
	vulkan.GetPhysicalDeviceMemoryProperties(c.physicalDevice, &c.memoryProperties)

	// Query queue family properties
	var queueFamilyCount uint32
	vulkan.GetPhysicalDeviceQueueFamilyProperties(c.physicalDevice, &queueFamilyCount, nil)
	c.queueFamilyProperties = make([]vulkan.QueueFamilyProperties, queueFamilyCount)
	vulkan.GetPhysicalDeviceQueueFamilyProperties(c.physicalDevice, &queueFamilyCount, c.queueFamilyProperties)

	// Extract capabilities from properties
	c.maxTextureSize = c.deviceProperties.Limits.MaxImageDimension2D
	c.maxRenderTargetSize = c.deviceProperties.Limits.MaxFramebufferWidth
	c.maxVertexAttributes = c.deviceProperties.Limits.MaxVertexInputAttributes
	c.maxUniformBufferRange = c.deviceProperties.Limits.MaxUniformBufferRange
	c.maxStorageBufferRange = c.deviceProperties.Limits.MaxStorageBufferRange
	c.maxPushConstantSize = c.deviceProperties.Limits.MaxPushConstantsSize

	// Check feature support
	c.supportsGeometryShaders = c.deviceFeatures.GeometryShader == vulkan.True
	c.supportsTessellation = c.deviceFeatures.TessellationShader == vulkan.True
	c.supportsComputeShaders = true // Vulkan always supports compute
	c.supportsMultisampling = true  // Vulkan always supports multisampling
	c.supportsAnisotropicFiltering = c.deviceFeatures.SamplerAnisotropy == vulkan.True
}

// convertPixelFormatToVulkan converts render pixel format to Vulkan format
func (c *CapabilitiesVK) convertPixelFormatToVulkan(format render.PixelFormat) vulkan.Format {
	switch format {
	case render.PixelFormatR8:
		return vulkan.FormatR8Unorm
	case render.PixelFormatRG8:
		return vulkan.FormatR8g8Unorm
	case render.PixelFormatRGBA8:
		return vulkan.FormatR8g8b8a8Unorm
	case render.PixelFormatBGRA8:
		return vulkan.FormatB8g8r8a8Unorm
	case render.PixelFormatR16Float:
		return vulkan.FormatR16Sfloat
	case render.PixelFormatRG16Float:
		return vulkan.FormatR16g16Sfloat
	case render.PixelFormatRGBA16Float:
		return vulkan.FormatR16g16b16a16Sfloat
	case render.PixelFormatR32Float:
		return vulkan.FormatR32Sfloat
	case render.PixelFormatRG32Float:
		return vulkan.FormatR32g32Sfloat
	case render.PixelFormatRGBA32Float:
		return vulkan.FormatR32g32b32a32Sfloat
	case render.PixelFormatDepth16:
		return vulkan.FormatD16Unorm
	case render.PixelFormatDepth24Stencil8:
		return vulkan.FormatD24UnormS8Uint
	case render.PixelFormatDepth32Float:
		return vulkan.FormatD32Sfloat
	default:
		return vulkan.FormatUndefined
	}
}
