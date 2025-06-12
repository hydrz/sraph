package impl

import (
	"fmt"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

// ValidateBufferDescriptor validates a buffer descriptor
func ValidateBufferDescriptor(desc BufferDescriptor, limits Limits) error {
	if desc.Size == 0 {
		return fmt.Errorf("buffer size cannot be zero")
	}

	if desc.Size > limits.MaxBufferSize {
		return fmt.Errorf("buffer size %d exceeds limit %d", desc.Size, limits.MaxBufferSize)
	}

	// Validate usage flags
	if desc.Usage == BufferUsageNone {
		return fmt.Errorf("buffer must have at least one usage flag")
	}

	// Check for conflicting usage flags
	if (desc.Usage&BufferUsageMapRead) != 0 && (desc.Usage&BufferUsageMapWrite) != 0 {
		return fmt.Errorf("buffer cannot have both MapRead and MapWrite usage")
	}

	return nil
}

// ValidateTextureDescriptor validates a texture descriptor
func ValidateTextureDescriptor(desc TextureDescriptor, limits Limits) error {
	if desc.Size.Width == 0 || desc.Size.Height == 0 || desc.Size.DepthOrArrayLayers == 0 {
		return fmt.Errorf("texture dimensions cannot be zero")
	}

	// Check dimension limits
	switch desc.Dimension {
	case TextureDimension1D:
		if desc.Size.Width > limits.MaxTextureDimension1D {
			return fmt.Errorf("1D texture width %d exceeds limit %d", desc.Size.Width, limits.MaxTextureDimension1D)
		}
	case TextureDimension2D:
		if desc.Size.Width > limits.MaxTextureDimension2D || desc.Size.Height > limits.MaxTextureDimension2D {
			return fmt.Errorf("2D texture dimensions exceed limit %d", limits.MaxTextureDimension2D)
		}
	case TextureDimension3D:
		maxDim := limits.MaxTextureDimension3D
		if desc.Size.Width > maxDim || desc.Size.Height > maxDim || desc.Size.DepthOrArrayLayers > maxDim {
			return fmt.Errorf("3D texture dimensions exceed limit %d", maxDim)
		}
	}

	// Validate array layers
	if desc.Size.DepthOrArrayLayers > limits.MaxTextureArrayLayers {
		return fmt.Errorf("texture array layers %d exceed limit %d", desc.Size.DepthOrArrayLayers, limits.MaxTextureArrayLayers)
	}

	// Validate usage flags
	if desc.Usage == TextureUsageNone {
		return fmt.Errorf("texture must have at least one usage flag")
	}

	return nil
}

// ValidateBindGroupLayoutDescriptor validates a bind group layout descriptor
func ValidateBindGroupLayoutDescriptor(desc BindGroupLayoutDescriptor, limits Limits) error {
	if uint32(len(desc.Entries)) > limits.MaxBindingsPerBindGroup {
		return fmt.Errorf("bind group layout has %d entries, limit is %d", len(desc.Entries), limits.MaxBindingsPerBindGroup)
	}

	// Check for duplicate bindings
	bindings := make(map[uint32]bool)
	for _, entry := range desc.Entries {
		if bindings[entry.Binding] {
			return fmt.Errorf("duplicate binding %d in bind group layout", entry.Binding)
		}
		bindings[entry.Binding] = true
	}

	return nil
}

// ValidateRenderPipelineDescriptor validates a render pipeline descriptor
func ValidateRenderPipelineDescriptor(desc RenderPipelineDescriptor, limits Limits) error {
	// Validate vertex stage
	if desc.Vertex.Module == nil {
		return fmt.Errorf("vertex shader module is required")
	}

	if desc.Vertex.EntryPoint == "" {
		return fmt.Errorf("vertex shader entry point is required")
	}

	// Validate vertex buffers
	if uint32(len(desc.Vertex.Buffers)) > limits.MaxVertexBuffers {
		return fmt.Errorf("vertex buffers count %d exceeds limit %d", len(desc.Vertex.Buffers), limits.MaxVertexBuffers)
	}

	// Count total vertex attributes
	totalAttributes := 0
	for _, buffer := range desc.Vertex.Buffers {
		totalAttributes += len(buffer.Attributes)
	}

	if uint32(totalAttributes) > limits.MaxVertexAttributes {
		return fmt.Errorf("total vertex attributes %d exceed limit %d", totalAttributes, limits.MaxVertexAttributes)
	}

	// Validate fragment stage if present
	if desc.Fragment.Module != nil {
		if desc.Fragment.EntryPoint == "" {
			return fmt.Errorf("fragment shader entry point is required when module is specified")
		}

		if uint32(len(desc.Fragment.Targets)) > limits.MaxColorAttachments {
			return fmt.Errorf("color targets count %d exceeds limit %d", len(desc.Fragment.Targets), limits.MaxColorAttachments)
		}
	}

	return nil
}

// ValidateComputePipelineDescriptor validates a compute pipeline descriptor
func ValidateComputePipelineDescriptor(desc ComputePipelineDescriptor) error {
	if desc.Compute.Module == nil {
		return fmt.Errorf("compute shader module is required")
	}

	if desc.Compute.EntryPoint == "" {
		return fmt.Errorf("compute shader entry point is required")
	}

	return nil
}
