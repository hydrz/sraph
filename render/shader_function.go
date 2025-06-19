// Shader function interface and implementation for the render package.
// This provides shader function abstraction for pipelines.

package render

import (
	"errors"
	"fmt"
	"sync"

	"github.com/opensraph/sraph/gpu"
)

// ShaderFunction represents a compiled shader function.
type ShaderFunction interface {
	// IsValid returns true if the shader function is valid and can be used.
	IsValid() bool

	// GetType returns the shader type.
	GetType() ShaderType

	// GetEntryPoint returns the entry point function name.
	GetEntryPoint() string

	// GetSource returns the shader source code.
	GetSource() []byte

	// GetMetadata returns shader metadata.
	GetMetadata() *ShaderMetadata

	// SetLabel sets a debug label for the shader function.
	SetLabel(label string)
}

// ShaderType represents the type of shader.
type ShaderType int

const (
	ShaderTypeUnknown ShaderType = iota
	ShaderTypeVertex
	ShaderTypeFragment
	ShaderTypeCompute
	ShaderTypeGeometry
	ShaderTypeTessellationControl
	ShaderTypeTessellationEvaluation
)

// ShaderStage represents the shader stage.
type ShaderStage int

const (
	ShaderStageVertex ShaderStage = iota
	ShaderStageFragment
	ShaderStageCompute
)

// ShaderMetadata contains metadata about a shader function.
type ShaderMetadata struct {
	Name          string
	EntryPoint    string
	Type          ShaderType
	Stage         ShaderStage
	Uniforms      []UniformInfo
	Attributes    []AttributeInfo
	Bindings      []BindingInfo
	WorkgroupSize [3]uint32 // For compute shaders
}

// UniformInfo describes a uniform variable in the shader.
type UniformInfo struct {
	Name    string
	Type    string
	Size    uint32
	Offset  uint32
	Binding uint32
	Set     uint32
}

// AttributeInfo describes a vertex attribute in the shader.
type AttributeInfo struct {
	Name     string
	Type     string
	Location uint32
	Format   string
}

// BindingInfo describes a resource binding in the shader.
type BindingInfo struct {
	Name    string
	Type    BindingType
	Binding uint32
	Set     uint32
}

// BindingType represents the type of resource binding.
type BindingType int

const (
	BindingTypeUnknown BindingType = iota
	BindingTypeUniformBuffer
	BindingTypeStorageBuffer
	BindingTypeTexture
	BindingTypeSampler
	BindingTypeStorageTexture
)

// DefaultShaderFunction is the default implementation of ShaderFunction.
type DefaultShaderFunction struct {
	shaderModule gpu.ShaderModule
	metadata     *ShaderMetadata
	source       []byte
	label        string
	isValid      bool
	mutex        sync.RWMutex
}

// NewShaderFunction creates a new shader function with the given source code.
func NewShaderFunction(device gpu.Device, source []byte, metadata *ShaderMetadata) (ShaderFunction, error) {
	if device == nil {
		return nil, errors.New("device cannot be nil")
	}
	if len(source) == 0 {
		return nil, errors.New("shader source cannot be empty")
	}
	if metadata == nil {
		return nil, errors.New("shader metadata cannot be nil")
	}

	shaderFunc := &DefaultShaderFunction{
		source:   source,
		metadata: metadata,
		label:    metadata.Name,
		isValid:  false,
	}

	// Create the shader module
	moduleDesc := &gpu.ShaderModuleDescriptor{
		Label:  metadata.Name,
		Source: string(source),
	}

	module, err := device.CreateShaderModule(moduleDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to create shader module: %w", err)
	}

	shaderFunc.shaderModule = module
	shaderFunc.isValid = true

	return shaderFunc, nil
}

// IsValid returns true if the shader function is valid and can be used.
func (sf *DefaultShaderFunction) IsValid() bool {
	sf.mutex.RLock()
	defer sf.mutex.RUnlock()
	return sf.isValid && sf.shaderModule != nil
}

// GetType returns the shader type.
func (sf *DefaultShaderFunction) GetType() ShaderType {
	sf.mutex.RLock()
	defer sf.mutex.RUnlock()
	return sf.metadata.Type
}

// GetEntryPoint returns the entry point function name.
func (sf *DefaultShaderFunction) GetEntryPoint() string {
	sf.mutex.RLock()
	defer sf.mutex.RUnlock()
	return sf.metadata.EntryPoint
}

// GetSource returns the shader source code.
func (sf *DefaultShaderFunction) GetSource() []byte {
	sf.mutex.RLock()
	defer sf.mutex.RUnlock()
	// Return a copy to prevent modifications
	source := make([]byte, len(sf.source))
	copy(source, sf.source)
	return source
}

// GetMetadata returns shader metadata.
func (sf *DefaultShaderFunction) GetMetadata() *ShaderMetadata {
	sf.mutex.RLock()
	defer sf.mutex.RUnlock()
	return sf.metadata
}

// SetLabel sets a debug label for the shader function.
func (sf *DefaultShaderFunction) SetLabel(label string) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	sf.label = label
	if sf.shaderModule != nil {
		sf.shaderModule.SetLabel(label)
	}
}

// String returns the string representation of the shader type.
func (st ShaderType) String() string {
	switch st {
	case ShaderTypeVertex:
		return "Vertex"
	case ShaderTypeFragment:
		return "Fragment"
	case ShaderTypeCompute:
		return "Compute"
	case ShaderTypeGeometry:
		return "Geometry"
	case ShaderTypeTessellationControl:
		return "TessellationControl"
	case ShaderTypeTessellationEvaluation:
		return "TessellationEvaluation"
	default:
		return "Unknown"
	}
}

// String returns the string representation of the shader stage.
func (ss ShaderStage) String() string {
	switch ss {
	case ShaderStageVertex:
		return "Vertex"
	case ShaderStageFragment:
		return "Fragment"
	case ShaderStageCompute:
		return "Compute"
	default:
		return "Unknown"
	}
}

// String returns the string representation of the binding type.
func (bt BindingType) String() string {
	switch bt {
	case BindingTypeUniformBuffer:
		return "UniformBuffer"
	case BindingTypeStorageBuffer:
		return "StorageBuffer"
	case BindingTypeTexture:
		return "Texture"
	case BindingTypeSampler:
		return "Sampler"
	case BindingTypeStorageTexture:
		return "StorageTexture"
	default:
		return "Unknown"
	}
}

// ShaderFunctionBuilder provides a builder pattern for creating shader functions.
type ShaderFunctionBuilder struct {
	device   gpu.Device
	source   []byte
	metadata *ShaderMetadata
}

// NewShaderFunctionBuilder creates a new shader function builder.
func NewShaderFunctionBuilder(device gpu.Device) *ShaderFunctionBuilder {
	return &ShaderFunctionBuilder{
		device: device,
		metadata: &ShaderMetadata{
			Uniforms:   make([]UniformInfo, 0),
			Attributes: make([]AttributeInfo, 0),
			Bindings:   make([]BindingInfo, 0),
		},
	}
}

// WithSource sets the shader source code.
func (sfb *ShaderFunctionBuilder) WithSource(source []byte) *ShaderFunctionBuilder {
	sfb.source = source
	return sfb
}

// WithName sets the shader name.
func (sfb *ShaderFunctionBuilder) WithName(name string) *ShaderFunctionBuilder {
	sfb.metadata.Name = name
	return sfb
}

// WithEntryPoint sets the shader entry point.
func (sfb *ShaderFunctionBuilder) WithEntryPoint(entryPoint string) *ShaderFunctionBuilder {
	sfb.metadata.EntryPoint = entryPoint
	return sfb
}

// WithType sets the shader type.
func (sfb *ShaderFunctionBuilder) WithType(shaderType ShaderType) *ShaderFunctionBuilder {
	sfb.metadata.Type = shaderType
	return sfb
}

// WithStage sets the shader stage.
func (sfb *ShaderFunctionBuilder) WithStage(stage ShaderStage) *ShaderFunctionBuilder {
	sfb.metadata.Stage = stage
	return sfb
}

// WithUniform adds a uniform to the shader metadata.
func (sfb *ShaderFunctionBuilder) WithUniform(name, uniformType string, size, offset, binding, set uint32) *ShaderFunctionBuilder {
	uniform := UniformInfo{
		Name:    name,
		Type:    uniformType,
		Size:    size,
		Offset:  offset,
		Binding: binding,
		Set:     set,
	}
	sfb.metadata.Uniforms = append(sfb.metadata.Uniforms, uniform)
	return sfb
}

// WithAttribute adds an attribute to the shader metadata.
func (sfb *ShaderFunctionBuilder) WithAttribute(name, attrType string, location uint32, format string) *ShaderFunctionBuilder {
	attribute := AttributeInfo{
		Name:     name,
		Type:     attrType,
		Location: location,
		Format:   format,
	}
	sfb.metadata.Attributes = append(sfb.metadata.Attributes, attribute)
	return sfb
}

// WithBinding adds a binding to the shader metadata.
func (sfb *ShaderFunctionBuilder) WithBinding(name string, bindingType BindingType, binding, set uint32) *ShaderFunctionBuilder {
	bindingInfo := BindingInfo{
		Name:    name,
		Type:    bindingType,
		Binding: binding,
		Set:     set,
	}
	sfb.metadata.Bindings = append(sfb.metadata.Bindings, bindingInfo)
	return sfb
}

// WithWorkgroupSize sets the workgroup size for compute shaders.
func (sfb *ShaderFunctionBuilder) WithWorkgroupSize(x, y, z uint32) *ShaderFunctionBuilder {
	sfb.metadata.WorkgroupSize = [3]uint32{x, y, z}
	return sfb
}

// Build builds and returns the shader function.
func (sfb *ShaderFunctionBuilder) Build() (ShaderFunction, error) {
	if sfb.device == nil {
		return nil, errors.New("device is required")
	}
	if len(sfb.source) == 0 {
		return nil, errors.New("shader source is required")
	}
	if sfb.metadata.Name == "" {
		return nil, errors.New("shader name is required")
	}
	if sfb.metadata.EntryPoint == "" {
		sfb.metadata.EntryPoint = "main"
	}

	return NewShaderFunction(sfb.device, sfb.source, sfb.metadata)
}
