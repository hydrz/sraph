package impl

import (
	"fmt"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

// getDefaultLimits returns default WebGPU limits for devices
func getDefaultLimits() Limits {
	return Limits{
		MaxTextureDimension1D:                     8192,
		MaxTextureDimension2D:                     8192,
		MaxTextureDimension3D:                     2048,
		MaxTextureArrayLayers:                     256,
		MaxBindGroups:                             4,
		MaxBindGroupsPlusVertexBuffers:            24,
		MaxBindingsPerBindGroup:                   1000,
		MaxDynamicUniformBuffersPerPipelineLayout: 8,
		MaxDynamicStorageBuffersPerPipelineLayout: 4,
		MaxSampledTexturesPerShaderStage:          16,
		MaxSamplersPerShaderStage:                 16,
		MaxStorageBuffersPerShaderStage:           8,
		MaxStorageTexturesPerShaderStage:          4,
		MaxUniformBuffersPerShaderStage:           12,
		MaxUniformBufferBindingSize:               65536,
		MaxStorageBufferBindingSize:               134217728,
		MinUniformBufferOffsetAlignment:           256,
		MinStorageBufferOffsetAlignment:           256,
		MaxVertexBuffers:                          8,
		MaxBufferSize:                             268435456,
		MaxVertexAttributes:                       16,
		MaxVertexBufferArrayStride:                2048,
		MaxInterStageShaderVariables:              16,
		MaxColorAttachments:                       8,
		MaxColorAttachmentBytesPerSample:          32,
		MaxComputeWorkgroupStorageSize:            16384,
		MaxComputeInvocationsPerWorkgroup:         256,
		MaxComputeWorkgroupSizeX:                  256,
		MaxComputeWorkgroupSizeY:                  256,
		MaxComputeWorkgroupSizeZ:                  64,
		MaxComputeWorkgroupsPerDimension:          65535,
		MaxImmediateSize:                          16777216,
	}
}

// ValidateTextureSize validates texture dimensions
func ValidateTextureSize(dimension TextureDimension, size Extent3D, limits Limits) error {
	switch dimension {
	case TextureDimension1D:
		if size.Width > limits.MaxTextureDimension1D {
			return fmt.Errorf("1D texture width %d exceeds limit %d", size.Width, limits.MaxTextureDimension1D)
		}
	case TextureDimension2D:
		if size.Width > limits.MaxTextureDimension2D || size.Height > limits.MaxTextureDimension2D {
			return fmt.Errorf("2D texture dimensions exceed limit %d", limits.MaxTextureDimension2D)
		}
	case TextureDimension3D:
		maxDim := limits.MaxTextureDimension3D
		if size.Width > maxDim || size.Height > maxDim || size.DepthOrArrayLayers > maxDim {
			return fmt.Errorf("3D texture dimensions exceed limit %d", maxDim)
		}
	}
	return nil
}

// GetTextureFormatBlockSize returns the block size in bytes for a texture format
func GetTextureFormatBlockSize(format TextureFormat) uint32 {
	switch format {
	case TextureFormatR8Unorm, TextureFormatR8Snorm, TextureFormatR8Uint, TextureFormatR8Sint:
		return 1
	case TextureFormatR16Uint, TextureFormatR16Sint, TextureFormatR16Float:
		return 2
	case TextureFormatRG8Unorm, TextureFormatRG8Snorm, TextureFormatRG8Uint, TextureFormatRG8Sint:
		return 2
	case TextureFormatR32Float, TextureFormatR32Uint, TextureFormatR32Sint:
		return 4
	case TextureFormatRG16Uint, TextureFormatRG16Sint, TextureFormatRG16Float:
		return 4
	case TextureFormatRGBA8Unorm, TextureFormatRGBA8UnormSrgb, TextureFormatRGBA8Snorm, TextureFormatRGBA8Uint, TextureFormatRGBA8Sint:
		return 4
	case TextureFormatBGRA8Unorm, TextureFormatBGRA8UnormSrgb:
		return 4
	case TextureFormatRG32Float, TextureFormatRG32Uint, TextureFormatRG32Sint:
		return 8
	case TextureFormatRGBA16Uint, TextureFormatRGBA16Sint, TextureFormatRGBA16Float:
		return 8
	case TextureFormatRGBA32Float, TextureFormatRGBA32Uint, TextureFormatRGBA32Sint:
		return 16
	default:
		return 1 // Default fallback
	}
}

// IsDepthFormat checks if a texture format is a depth format
func IsDepthFormat(format TextureFormat) bool {
	switch format {
	case TextureFormatDepth16Unorm, TextureFormatDepth24Plus, TextureFormatDepth32Float, TextureFormatDepth24PlusStencil8, TextureFormatDepth32FloatStencil8:
		return true
	default:
		return false
	}
}

// IsStencilFormat checks if a texture format has stencil component
func IsStencilFormat(format TextureFormat) bool {
	switch format {
	case TextureFormatStencil8, TextureFormatDepth24PlusStencil8, TextureFormatDepth32FloatStencil8:
		return true
	default:
		return false
	}
}

// CalculateBufferAlignment calculates proper buffer alignment
func CalculateBufferAlignment(usage BufferUsage, limits Limits) uint32 {
	if (usage & BufferUsageUniform) != 0 {
		return limits.MinUniformBufferOffsetAlignment
	}
	if (usage & BufferUsageStorage) != 0 {
		return limits.MinStorageBufferOffsetAlignment
	}
	return 1 // Default alignment
}
