package vulkan

import (
	"github.com/opensraph/sraph/render"
	"github.com/vulkan-go/vulkan"
)

// VertexDescriptorVK represents a Vulkan vertex input descriptor
type VertexDescriptorVK struct {
	bindingDescriptions   []vulkan.VertexInputBindingDescription
	attributeDescriptions []vulkan.VertexInputAttributeDescription
	inputState            vulkan.PipelineVertexInputStateCreateInfo
}

// NewVertexDescriptorVK creates a new Vulkan vertex descriptor
func NewVertexDescriptorVK() *VertexDescriptorVK {
	return &VertexDescriptorVK{
		bindingDescriptions:   make([]vulkan.VertexInputBindingDescription, 0),
		attributeDescriptions: make([]vulkan.VertexInputAttributeDescription, 0),
	}
}

// NewVertexDescriptorVKFromRender creates a Vulkan vertex descriptor from render descriptor
func NewVertexDescriptorVKFromRender(descriptor render.VertexDescriptor) *VertexDescriptorVK {
	vkDescriptor := NewVertexDescriptorVK()

	// TODO: Convert render vertex descriptor to Vulkan format
	// This would involve mapping render.VertexFormat to Vulkan formats
	// and setting up binding and attribute descriptions

	return vkDescriptor
}

// AddBinding adds a vertex input binding
func (v *VertexDescriptorVK) AddBinding(binding uint32, stride uint32, inputRate vulkan.VertexInputRate) {
	bindingDesc := vulkan.VertexInputBindingDescription{
		Binding:   binding,
		Stride:    stride,
		InputRate: inputRate,
	}
	v.bindingDescriptions = append(v.bindingDescriptions, bindingDesc)
}

// AddAttribute adds a vertex input attribute
func (v *VertexDescriptorVK) AddAttribute(location, binding, offset uint32, format vulkan.Format) {
	attributeDesc := vulkan.VertexInputAttributeDescription{
		Location: location,
		Binding:  binding,
		Format:   format,
		Offset:   offset,
	}
	v.attributeDescriptions = append(v.attributeDescriptions, attributeDesc)
}

// GetInputState returns the pipeline vertex input state create info
func (v *VertexDescriptorVK) GetInputState() vulkan.PipelineVertexInputStateCreateInfo {
	v.inputState = vulkan.PipelineVertexInputStateCreateInfo{
		SType:                           vulkan.StructureTypePipelineVertexInputStateCreateInfo,
		VertexBindingDescriptionCount:   uint32(len(v.bindingDescriptions)),
		PVertexBindingDescriptions:      v.bindingDescriptions,
		VertexAttributeDescriptionCount: uint32(len(v.attributeDescriptions)),
		PVertexAttributeDescriptions:    v.attributeDescriptions,
	}
	return v.inputState
}

// GetBindingDescriptions returns the vertex input binding descriptions
func (v *VertexDescriptorVK) GetBindingDescriptions() []vulkan.VertexInputBindingDescription {
	return v.bindingDescriptions
}

// GetAttributeDescriptions returns the vertex input attribute descriptions
func (v *VertexDescriptorVK) GetAttributeDescriptions() []vulkan.VertexInputAttributeDescription {
	return v.attributeDescriptions
}

// Clear clears all binding and attribute descriptions
func (v *VertexDescriptorVK) Clear() {
	v.bindingDescriptions = v.bindingDescriptions[:0]
	v.attributeDescriptions = v.attributeDescriptions[:0]
}

// GetBindingCount returns the number of vertex bindings
func (v *VertexDescriptorVK) GetBindingCount() uint32 {
	return uint32(len(v.bindingDescriptions))
}

// GetAttributeCount returns the number of vertex attributes
func (v *VertexDescriptorVK) GetAttributeCount() uint32 {
	return uint32(len(v.attributeDescriptions))
}

// ConvertVertexFormat converts a render vertex format to Vulkan format
func ConvertVertexFormat(format render.VertexFormat) vulkan.Format {
	switch format {
	case render.VertexFormatFloat:
		return vulkan.FormatR32Sfloat
	case render.VertexFormatFloat2:
		return vulkan.FormatR32g32Sfloat
	case render.VertexFormatFloat3:
		return vulkan.FormatR32g32b32Sfloat
	case render.VertexFormatFloat4:
		return vulkan.FormatR32g32b32a32Sfloat
	case render.VertexFormatUByte4:
		return vulkan.FormatR8g8b8a8Unorm
	case render.VertexFormatUShort2:
		return vulkan.FormatR16g16Uint
	case render.VertexFormatUShort4:
		return vulkan.FormatR16g16b16a16Uint
	default:
		return vulkan.FormatUndefined
	}
}

// ConvertVertexInputRate converts a render vertex step to Vulkan input rate
func ConvertVertexInputRate(step render.VertexStep) vulkan.VertexInputRate {
	switch step {
	case render.VertexStepVertex:
		return vulkan.VertexInputRateVertex
	case render.VertexStepInstance:
		return vulkan.VertexInputRateInstance
	default:
		return vulkan.VertexInputRateVertex
	}
}
