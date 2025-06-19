// Vertex descriptor interface and implementation for the render package.
// This describes the layout and format of vertex data.

package render

import (
	"errors"
	"fmt"

	"github.com/opensraph/sraph/gpu"
)

// VertexDescriptor describes the layout of vertex data.
type VertexDescriptor struct {
	Attributes []VertexAttribute
	Layouts    []VertexBufferLayout
}

// VertexSemantic represents the semantic meaning of a vertex attribute.
type VertexSemantic int

const (
	VertexSemanticUnknown VertexSemantic = iota
	VertexSemanticPosition
	VertexSemanticNormal
	VertexSemanticTangent
	VertexSemanticBitangent
	VertexSemanticColor
	VertexSemanticTexCoord
	VertexSemanticJoints
	VertexSemanticWeights
	VertexSemanticCustom
)

// String returns the string representation of the vertex semantic.
func (vs VertexSemantic) String() string {
	switch vs {
	case VertexSemanticPosition:
		return "Position"
	case VertexSemanticNormal:
		return "Normal"
	case VertexSemanticTangent:
		return "Tangent"
	case VertexSemanticBitangent:
		return "Bitangent"
	case VertexSemanticColor:
		return "Color"
	case VertexSemanticTexCoord:
		return "TexCoord"
	case VertexSemanticJoints:
		return "Joints"
	case VertexSemanticWeights:
		return "Weights"
	case VertexSemanticCustom:
		return "Custom"
	default:
		return "Unknown"
	}
}

// Validate validates the vertex descriptor.
func (vd *VertexDescriptor) Validate() error {
	if len(vd.Attributes) == 0 {
		return errors.New("vertex descriptor must have at least one attribute")
	}

	if len(vd.Layouts) == 0 {
		return errors.New("vertex descriptor must have at least one layout")
	}

	// Check for duplicate attribute locations
	locationMap := make(map[uint32]bool)
	for _, attr := range vd.Attributes {
		if locationMap[attr.ShaderLocation] {
			return fmt.Errorf("duplicate attribute location: %d", attr.ShaderLocation)
		}
		locationMap[attr.ShaderLocation] = true
	}

	// Validate layouts
	for _, layout := range vd.Layouts {
		if layout.ArrayStride == 0 {
			return errors.New("vertex buffer layout has zero stride")
		}
	}

	return nil
}

// GetAttribute returns the attribute at the given shader location.
func (vd *VertexDescriptor) GetAttribute(location uint32) (*VertexAttribute, error) {
	for i := range vd.Attributes {
		if vd.Attributes[i].ShaderLocation == location {
			return &vd.Attributes[i], nil
		}
	}
	return nil, fmt.Errorf("no attribute found at location %d", location)
}

// GetLayout returns the layout at the given index.
func (vd *VertexDescriptor) GetLayout(index int) (*VertexBufferLayout, error) {
	if index < 0 || index >= len(vd.Layouts) {
		return nil, fmt.Errorf("layout index %d out of range", index)
	}
	return &vd.Layouts[index], nil
}

// GetTotalStride returns the total stride of all vertex buffers.
func (vd *VertexDescriptor) GetTotalStride() uint64 {
	var totalStride uint64
	for _, layout := range vd.Layouts {
		totalStride += layout.ArrayStride
	}
	return totalStride
}

// NewVertexDescriptor creates a new vertex descriptor.
func NewVertexDescriptor() *VertexDescriptor {
	return &VertexDescriptor{
		Attributes: make([]VertexAttribute, 0),
		Layouts:    make([]VertexBufferLayout, 0),
	}
}

// AddAttribute adds a vertex attribute to the descriptor.
func (vd *VertexDescriptor) AddAttribute(format gpu.VertexFormat, offset uint64, shaderLocation uint32) {
	attr := VertexAttribute{
		Format:         format,
		Offset:         offset,
		ShaderLocation: shaderLocation,
	}
	vd.Attributes = append(vd.Attributes, attr)
}

// AddLayout adds a vertex buffer layout to the descriptor.
func (vd *VertexDescriptor) AddLayout(arrayStride uint64, stepMode gpu.VertexStepMode, attributes []VertexAttribute) {
	layout := VertexBufferLayout{
		ArrayStride: arrayStride,
		StepMode:    stepMode,
		Attributes:  attributes,
	}
	vd.Layouts = append(vd.Layouts, layout)
}

// VertexDescriptorBuilder provides a builder pattern for creating vertex descriptors.
type VertexDescriptorBuilder struct {
	descriptor *VertexDescriptor
}

// NewVertexDescriptorBuilder creates a new vertex descriptor builder.
func NewVertexDescriptorBuilder() *VertexDescriptorBuilder {
	return &VertexDescriptorBuilder{
		descriptor: NewVertexDescriptor(),
	}
}

// WithAttribute adds an attribute to the vertex descriptor.
func (vdb *VertexDescriptorBuilder) WithAttribute(format gpu.VertexFormat, offset uint64, shaderLocation uint32) *VertexDescriptorBuilder {
	vdb.descriptor.AddAttribute(format, offset, shaderLocation)
	return vdb
}

// WithLayout adds a layout to the vertex descriptor.
func (vdb *VertexDescriptorBuilder) WithLayout(arrayStride uint64, stepMode gpu.VertexStepMode, attributes []VertexAttribute) *VertexDescriptorBuilder {
	vdb.descriptor.AddLayout(arrayStride, stepMode, attributes)
	return vdb
}

// Build builds and returns the vertex descriptor.
func (vdb *VertexDescriptorBuilder) Build() (*VertexDescriptor, error) {
	if err := vdb.descriptor.Validate(); err != nil {
		return nil, err
	}
	return vdb.descriptor, nil
}
