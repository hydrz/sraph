package render

import (
	"github.com/opensraph/sraph/gpu"
)

// Capabilities describes the rendering capabilities of the current GPU context.
// This includes information about supported features, limits, and formats.
type Capabilities interface {
	// SupportsComputeShaders returns true if compute shaders are supported
	SupportsComputeShaders() bool

	// SupportsFramebufferFetch returns true if framebuffer fetch is supported
	SupportsFramebufferFetch() bool

	// SupportsImplicitResolvingMSAA returns true if implicit MSAA resolving is supported
	SupportsImplicitResolvingMSAA() bool

	// SupportsSSBO returns true if Shader Storage Buffer Objects are supported
	SupportsSSBO() bool

	// SupportsTextureToTextureBlits returns true if texture-to-texture blits are supported
	SupportsTextureToTextureBlits() bool

	// SupportsDecalSamplerAddressMode returns true if decal sampler address mode is supported
	SupportsDecalSamplerAddressMode() bool

	// GetMaxTextureSize returns the maximum texture size supported
	GetMaxTextureSize() uint32

	// GetMaxBufferLength returns the maximum buffer size supported
	GetMaxBufferLength() uint64

	// GetDefaultColorFormat returns the default color format
	GetDefaultColorFormat() gpu.TextureFormat

	// GetDefaultStencilFormat returns the default stencil format
	GetDefaultStencilFormat() gpu.TextureFormat

	// GetDefaultDepthStencilFormat returns the default depth-stencil format
	GetDefaultDepthStencilFormat() gpu.TextureFormat

	// IsTextureFormatSupported returns true if the specified format is supported for textures
	IsTextureFormatSupported(format gpu.TextureFormat) bool

	// GetSupportsReadFromResolve returns true if reading from resolve attachments is supported
	GetSupportsReadFromResolve() bool

	// GetSupportsReadFromOnscreenTexture returns true if reading from onscreen textures is supported
	GetSupportsReadFromOnscreenTexture() bool
}

// PixelFormat represents different pixel formats and their capabilities
type PixelFormat struct {
	// Format is the GPU texture format
	Format gpu.TextureFormat

	// Type indicates whether this is a color, depth, or stencil format
	Type PixelFormatType

	// ComponentType indicates the data type of components
	ComponentType PixelFormatComponentType

	// ComponentCount is the number of components per pixel
	ComponentCount uint32

	// BytesPerPixel is the size in bytes per pixel
	BytesPerPixel uint32

	// IsSupported indicates if this format is supported by the current GPU
	IsSupported bool
}

// PixelFormatType defines the type of pixel format
type PixelFormatType int

const (
	PixelFormatTypeColor PixelFormatType = iota
	PixelFormatTypeDepth
	PixelFormatTypeStencil
	PixelFormatTypeDepthStencil
)

// PixelFormatComponentType defines the component data type
type PixelFormatComponentType int

const (
	PixelFormatComponentTypeUNorm PixelFormatComponentType = iota
	PixelFormatComponentTypeSNorm
	PixelFormatComponentTypeUInt
	PixelFormatComponentTypeSInt
	PixelFormatComponentTypeFloat
)

// DefaultCapabilities provides a default implementation of Capabilities
type DefaultCapabilities struct {
	// GPU context
	context gpu.Context

	// Cached capability flags
	supportsComputeShaders          bool
	supportsFramebufferFetch        bool
	supportsImplicitResolvingMSAA   bool
	supportsSSBO                    bool
	supportsTextureToTextureBlits   bool
	supportsDecalSamplerAddressMode bool
	supportsReadFromResolve         bool
	supportsReadFromOnscreenTexture bool

	// Limits
	maxTextureSize  uint32
	maxBufferLength uint64

	// Default formats
	defaultColorFormat        gpu.TextureFormat
	defaultStencilFormat      gpu.TextureFormat
	defaultDepthStencilFormat gpu.TextureFormat

	// Supported formats
	supportedFormats map[gpu.TextureFormat]bool
}

// NewCapabilities creates a new capabilities object by querying the GPU context
func NewCapabilities(context gpu.Context) Capabilities {
	caps := &DefaultCapabilities{
		context:          context,
		supportedFormats: make(map[gpu.TextureFormat]bool),
	}

	// Query capabilities from the GPU context
	caps.queryCapabilities()

	return caps
}

// queryCapabilities queries the GPU context for supported capabilities
func (c *DefaultCapabilities) queryCapabilities() {
	// TODO: Query actual capabilities from the GPU context
	// For now, set reasonable defaults

	c.supportsComputeShaders = true
	c.supportsFramebufferFetch = false
	c.supportsImplicitResolvingMSAA = true
	c.supportsSSBO = true
	c.supportsTextureToTextureBlits = true
	c.supportsDecalSamplerAddressMode = false
	c.supportsReadFromResolve = true
	c.supportsReadFromOnscreenTexture = true

	c.maxTextureSize = 16384
	c.maxBufferLength = 1024 * 1024 * 1024 // 1GB

	c.defaultColorFormat = gpu.TextureFormatRGBA8Unorm
	c.defaultStencilFormat = gpu.TextureFormatStencil8
	c.defaultDepthStencilFormat = gpu.TextureFormatDepth24PlusStencil8

	// Initialize supported formats
	c.initializeSupportedFormats()
}

// initializeSupportedFormats initializes the list of supported texture formats
func (c *DefaultCapabilities) initializeSupportedFormats() {
	// TODO: Query actual supported formats from GPU context
	// For now, assume common formats are supported

	supportedFormats := []gpu.TextureFormat{
		gpu.TextureFormatRGBA8Unorm,
		gpu.TextureFormatRGBA8UnormSrgb,
		gpu.TextureFormatBGRA8Unorm,
		gpu.TextureFormatBGRA8UnormSrgb,
		gpu.TextureFormatRGBA16Float,
		gpu.TextureFormatRGBA32Float,
		gpu.TextureFormatDepth32Float,
		gpu.TextureFormatDepth24Plus,
		gpu.TextureFormatDepth24PlusStencil8,
		gpu.TextureFormatStencil8,
		gpu.TextureFormatR8Unorm,
		gpu.TextureFormatRG8Unorm,
		gpu.TextureFormatR16Float,
		gpu.TextureFormatRG16Float,
	}

	for _, format := range supportedFormats {
		c.supportedFormats[format] = true
	}
}

// SupportsComputeShaders returns true if compute shaders are supported
func (c *DefaultCapabilities) SupportsComputeShaders() bool {
	return c.supportsComputeShaders
}

// SupportsFramebufferFetch returns true if framebuffer fetch is supported
func (c *DefaultCapabilities) SupportsFramebufferFetch() bool {
	return c.supportsFramebufferFetch
}

// SupportsImplicitResolvingMSAA returns true if implicit MSAA resolving is supported
func (c *DefaultCapabilities) SupportsImplicitResolvingMSAA() bool {
	return c.supportsImplicitResolvingMSAA
}

// SupportsSSBO returns true if Shader Storage Buffer Objects are supported
func (c *DefaultCapabilities) SupportsSSBO() bool {
	return c.supportsSSBO
}

// SupportsTextureToTextureBlits returns true if texture-to-texture blits are supported
func (c *DefaultCapabilities) SupportsTextureToTextureBlits() bool {
	return c.supportsTextureToTextureBlits
}

// SupportsDecalSamplerAddressMode returns true if decal sampler address mode is supported
func (c *DefaultCapabilities) SupportsDecalSamplerAddressMode() bool {
	return c.supportsDecalSamplerAddressMode
}

// GetMaxTextureSize returns the maximum texture size supported
func (c *DefaultCapabilities) GetMaxTextureSize() uint32 {
	return c.maxTextureSize
}

// GetMaxBufferLength returns the maximum buffer size supported
func (c *DefaultCapabilities) GetMaxBufferLength() uint64 {
	return c.maxBufferLength
}

// GetDefaultColorFormat returns the default color format
func (c *DefaultCapabilities) GetDefaultColorFormat() gpu.TextureFormat {
	return c.defaultColorFormat
}

// GetDefaultStencilFormat returns the default stencil format
func (c *DefaultCapabilities) GetDefaultStencilFormat() gpu.TextureFormat {
	return c.defaultStencilFormat
}

// GetDefaultDepthStencilFormat returns the default depth-stencil format
func (c *DefaultCapabilities) GetDefaultDepthStencilFormat() gpu.TextureFormat {
	return c.defaultDepthStencilFormat
}

// IsTextureFormatSupported returns true if the specified format is supported for textures
func (c *DefaultCapabilities) IsTextureFormatSupported(format gpu.TextureFormat) bool {
	return c.supportedFormats[format]
}

// GetSupportsReadFromResolve returns true if reading from resolve attachments is supported
func (c *DefaultCapabilities) GetSupportsReadFromResolve() bool {
	return c.supportsReadFromResolve
}

// GetSupportsReadFromOnscreenTexture returns true if reading from onscreen textures is supported
func (c *DefaultCapabilities) GetSupportsReadFromOnscreenTexture() bool {
	return c.supportsReadFromOnscreenTexture
}
