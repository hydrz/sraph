package render

// PipelineBuilder provides a fluent interface for constructing rendering pipelines
// Used to build both graphics and compute pipelines with validation
type PipelineBuilder interface {
	// SetVertexShader sets the vertex shader for the pipeline
	SetVertexShader(shader ShaderFunction) PipelineBuilder

	// SetFragmentShader sets the fragment shader for the pipeline
	SetFragmentShader(shader ShaderFunction) PipelineBuilder

	// SetVertexDescriptor sets the vertex input descriptor
	SetVertexDescriptor(descriptor VertexDescriptor) PipelineBuilder

	// SetBlendMode sets the blending configuration
	SetBlendMode(mode BlendMode) PipelineBuilder

	// SetColorFormat sets the color attachment format
	SetColorFormat(format PixelFormat) PipelineBuilder

	// SetDepthFormat sets the depth attachment format
	SetDepthFormat(format PixelFormat) PipelineBuilder

	// SetStencilFormat sets the stencil attachment format
	SetStencilFormat(format PixelFormat) PipelineBuilder

	// SetPrimitiveType sets the primitive topology
	SetPrimitiveType(primitiveType PrimitiveType) PipelineBuilder

	// SetCullMode sets the face culling mode
	SetCullMode(mode CullMode) PipelineBuilder

	// SetWindingOrder sets the triangle winding order
	SetWindingOrder(order WindingOrder) PipelineBuilder

	// SetDepthTest enables or disables depth testing
	SetDepthTest(enabled bool) PipelineBuilder

	// SetDepthWrite enables or disables depth writing
	SetDepthWrite(enabled bool) PipelineBuilder

	// SetStencilTest enables or disables stencil testing
	SetStencilTest(enabled bool) PipelineBuilder

	// Build creates the final pipeline
	Build() (Pipeline, error)

	// IsValid returns true if the current configuration is valid
	IsValid() bool
}

// ComputePipelineBuilder provides a fluent interface for constructing compute pipelines
type ComputePipelineBuilder interface {
	// SetComputeShader sets the compute shader for the pipeline
	SetComputeShader(shader ShaderFunction) ComputePipelineBuilder

	// SetThreadgroupSize sets the threadgroup size
	SetThreadgroupSize(x, y, z uint32) ComputePipelineBuilder

	// Build creates the final compute pipeline
	Build() (ComputePipeline, error)

	// IsValid returns true if the current configuration is valid
	IsValid() bool
}

// GraphicsPipelineBuilder implements PipelineBuilder for graphics pipelines
type GraphicsPipelineBuilder struct {
	vertexShader       ShaderFunction
	fragmentShader     ShaderFunction
	vertexDescriptor   VertexDescriptor
	blendMode          BlendMode
	colorFormat        PixelFormat
	depthFormat        PixelFormat
	stencilFormat      PixelFormat
	primitiveType      PrimitiveType
	cullMode           CullMode
	windingOrder       WindingOrder
	depthTestEnabled   bool
	depthWriteEnabled  bool
	stencilTestEnabled bool
}

// NewGraphicsPipelineBuilder creates a new graphics pipeline builder
func NewGraphicsPipelineBuilder() *GraphicsPipelineBuilder {
	return &GraphicsPipelineBuilder{
		blendMode:          BlendModeSourceOver,
		colorFormat:        PixelFormatRGBA8,
		primitiveType:      PrimitiveTypeTriangle,
		cullMode:           CullModeNone,
		windingOrder:       WindingOrderClockwise,
		depthTestEnabled:   false,
		depthWriteEnabled:  false,
		stencilTestEnabled: false,
	}
}

// SetVertexShader implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetVertexShader(shader ShaderFunction) PipelineBuilder {
	b.vertexShader = shader
	return b
}

// SetFragmentShader implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetFragmentShader(shader ShaderFunction) PipelineBuilder {
	b.fragmentShader = shader
	return b
}

// SetVertexDescriptor implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetVertexDescriptor(descriptor VertexDescriptor) PipelineBuilder {
	b.vertexDescriptor = descriptor
	return b
}

// SetBlendMode implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetBlendMode(mode BlendMode) PipelineBuilder {
	b.blendMode = mode
	return b
}

// SetColorFormat implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetColorFormat(format PixelFormat) PipelineBuilder {
	b.colorFormat = format
	return b
}

// SetDepthFormat implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetDepthFormat(format PixelFormat) PipelineBuilder {
	b.depthFormat = format
	return b
}

// SetStencilFormat implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetStencilFormat(format PixelFormat) PipelineBuilder {
	b.stencilFormat = format
	return b
}

// SetPrimitiveType implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetPrimitiveType(primitiveType PrimitiveType) PipelineBuilder {
	b.primitiveType = primitiveType
	return b
}

// SetCullMode implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetCullMode(mode CullMode) PipelineBuilder {
	b.cullMode = mode
	return b
}

// SetWindingOrder implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetWindingOrder(order WindingOrder) PipelineBuilder {
	b.windingOrder = order
	return b
}

// SetDepthTest implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetDepthTest(enabled bool) PipelineBuilder {
	b.depthTestEnabled = enabled
	return b
}

// SetDepthWrite implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetDepthWrite(enabled bool) PipelineBuilder {
	b.depthWriteEnabled = enabled
	return b
}

// SetStencilTest implements PipelineBuilder
func (b *GraphicsPipelineBuilder) SetStencilTest(enabled bool) PipelineBuilder {
	b.stencilTestEnabled = enabled
	return b
}

// IsValid implements PipelineBuilder
func (b *GraphicsPipelineBuilder) IsValid() bool {
	return b.vertexShader != nil && b.fragmentShader != nil
}

// Build implements PipelineBuilder
func (b *GraphicsPipelineBuilder) Build() (Pipeline, error) {
	if !b.IsValid() {
		return nil, ErrInvalidPipelineConfiguration
	}

	// TODO: Implement actual pipeline creation
	// This would involve creating the appropriate backend-specific pipeline
	return nil, ErrNotImplemented
}

// DefaultComputePipelineBuilder implements ComputePipelineBuilder
type DefaultComputePipelineBuilder struct {
	computeShader   ShaderFunction
	threadgroupSize [3]uint32
}

// NewComputePipelineBuilder creates a new compute pipeline builder
func NewComputePipelineBuilder() *DefaultComputePipelineBuilder {
	return &DefaultComputePipelineBuilder{
		threadgroupSize: [3]uint32{1, 1, 1},
	}
}

// SetComputeShader implements ComputePipelineBuilder
func (b *DefaultComputePipelineBuilder) SetComputeShader(shader ShaderFunction) ComputePipelineBuilder {
	b.computeShader = shader
	return b
}

// SetThreadgroupSize implements ComputePipelineBuilder
func (b *DefaultComputePipelineBuilder) SetThreadgroupSize(x, y, z uint32) ComputePipelineBuilder {
	b.threadgroupSize = [3]uint32{x, y, z}
	return b
}

// IsValid implements ComputePipelineBuilder
func (b *DefaultComputePipelineBuilder) IsValid() bool {
	return b.computeShader != nil
}

// Build implements ComputePipelineBuilder
func (b *DefaultComputePipelineBuilder) Build() (ComputePipeline, error) {
	if !b.IsValid() {
		return nil, ErrInvalidPipelineConfiguration
	}

	// TODO: Implement actual compute pipeline creation
	return nil, ErrNotImplemented
}
