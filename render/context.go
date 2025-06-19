// Package render provides advanced rendering capabilities for Sraph,
// implementing a modern GPU-accelerated renderer inspired by Flutter's Impeller.
//
// This package contains the core rendering abstractions and implementations
// that bridge between high-level graphics operations and low-level GPU commands.
package render

import (
	"fmt"

	"github.com/opensraph/sraph/gpu"
)

// Context represents a rendering context that manages GPU resources and state.
// It provides the foundation for all rendering operations.
type Context interface {
	// GetDevice returns the underlying GPU device
	GetDevice() gpu.Device

	// GetQueue returns the command queue for GPU operations
	GetQueue() gpu.Queue

	// GetCapabilities returns the rendering capabilities
	GetCapabilities() *Capabilities

	// GetResourceAllocator returns the resource allocator
	GetResourceAllocator() ResourceAllocator

	// GetShaderLibrary returns the shader library
	GetShaderLibrary() ShaderLibrary

	// GetPipelineLibrary returns the pipeline library
	GetPipelineLibrary() PipelineLibrary

	// GetSamplerLibrary returns the sampler library
	GetSamplerLibrary() SamplerLibrary

	// Shutdown shuts down the context and releases resources
	Shutdown() error
}

// Capabilities describes the rendering capabilities of the current backend.
type Capabilities struct {
	// SupportsAdvancedBlends indicates if advanced blend modes are supported
	SupportsAdvancedBlends bool

	// MaxTextureSize is the maximum texture dimension supported
	MaxTextureSize uint32

	// SupportsCompute indicates if compute shaders are supported
	SupportsCompute bool

	// SupportsTimestamps indicates if GPU timestamps are supported
	SupportsTimestamps bool

	// SupportsMultisampling indicates if MSAA is supported
	SupportsMultisampling bool

	// MaxSampleCount is the maximum sample count for MSAA
	MaxSampleCount uint32

	// SupportsInstancedDrawing indicates if instanced drawing is supported
	SupportsInstancedDrawing bool

	// DefaultColorFormat is the preferred color format for render targets
	DefaultColorFormat gpu.TextureFormat

	// DefaultDepthStencilFormat is the preferred depth/stencil format
	DefaultDepthStencilFormat gpu.TextureFormat
}

// ResourceAllocator manages GPU resource allocation and pooling.
type ResourceAllocator interface {
	// CreateBuffer creates a new buffer with the given descriptor
	CreateBuffer(desc BufferDescriptor) (Buffer, error)

	// CreateTexture creates a new texture with the given descriptor
	CreateTexture(desc TextureDescriptor) (Texture, error)

	// GetHostBuffer returns a reusable host buffer for uploading data
	GetHostBuffer(size uint64) (Buffer, error)

	// GetTransientBuffer returns a transient buffer for temporary data
	GetTransientBuffer(size uint64) (Buffer, error)
}

// ShaderLibrary manages shader compilation and caching.
type ShaderLibrary interface {
	// GetShader retrieves a compiled shader by name
	GetShader(name string) (gpu.ShaderModule, error)

	// CompileShader compiles shader source code
	CompileShader(source string, stage ShaderStage) (gpu.ShaderModule, error)

	// GetShaderFunction retrieves a shader function by name and stage
	GetShaderFunction(name string, stage ShaderStage) (ShaderFunction, error)
}

// PipelineLibrary manages render and compute pipeline creation and caching.
type PipelineLibrary interface {
	// GetRenderPipeline retrieves a render pipeline by descriptor
	GetRenderPipeline(desc RenderPipelineDescriptor) (gpu.RenderPipeline, error)

	// GetComputePipeline retrieves a compute pipeline by descriptor
	GetComputePipeline(desc ComputePipelineDescriptor) (gpu.ComputePipeline, error)

	// CreateRenderPipeline creates a new render pipeline
	CreateRenderPipeline(desc RenderPipelineDescriptor) (gpu.RenderPipeline, error)

	// CreateComputePipeline creates a new compute pipeline
	CreateComputePipeline(desc ComputePipelineDescriptor) (gpu.ComputePipeline, error)
}

// SamplerLibrary manages sampler state objects and caching.
type SamplerLibrary interface {
	// GetSampler retrieves a sampler by descriptor
	GetSampler(desc SamplerDescriptor) (gpu.Sampler, error)

	// CreateSampler creates a new sampler
	CreateSampler(desc SamplerDescriptor) (gpu.Sampler, error)
}

// Buffer represents a GPU buffer resource.
type Buffer interface {
	gpu.Buffer

	// GetDescriptor returns the buffer descriptor
	GetDescriptor() BufferDescriptor

	// SetLabel sets a debug label for the buffer
	SetLabel(label string)
}

// Texture represents a GPU texture resource.
type Texture interface {
	gpu.Texture

	// GetDescriptor returns the texture descriptor
	GetDescriptor() TextureDescriptor

	// SetLabel sets a debug label for the texture
	SetLabel(label string)
}

// ShaderFunction represents a compiled shader function.
type ShaderFunction interface {
	// GetName returns the function name
	GetName() string

	// GetStage returns the shader stage
	GetStage() ShaderStage

	// GetModule returns the underlying shader module
	GetModule() gpu.ShaderModule
}

// ShaderStage represents the stage of a shader in the graphics pipeline.
type ShaderStage uint32

const (
	ShaderStageVertex   ShaderStage = iota // Vertex shader stage
	ShaderStageFragment                    // Fragment/pixel shader stage
	ShaderStageCompute                     // Compute shader stage
)

// BufferDescriptor describes how to create a buffer.
type BufferDescriptor struct {
	Label string
	Size  uint64
	Usage gpu.BufferUsage
}

// TextureDescriptor describes how to create a texture.
type TextureDescriptor struct {
	Label         string
	Size          gpu.Extent3D
	Format        gpu.TextureFormat
	Usage         gpu.TextureUsage
	SampleCount   uint32
	MipLevelCount uint32
	Dimension     gpu.TextureDimension
}

// SamplerDescriptor describes how to create a sampler.
type SamplerDescriptor struct {
	Label         string
	AddressModeU  gpu.AddressMode
	AddressModeV  gpu.AddressMode
	AddressModeW  gpu.AddressMode
	MagFilter     gpu.FilterMode
	MinFilter     gpu.FilterMode
	MipmapFilter  gpu.MipmapFilterMode
	LodMinClamp   float32
	LodMaxClamp   float32
	Compare       gpu.CompareFunction
	MaxAnisotropy uint16
}

// RenderPipelineDescriptor describes how to create a render pipeline.
type RenderPipelineDescriptor struct {
	Label        string
	Layout       gpu.PipelineLayout
	Vertex       VertexState
	Primitive    PrimitiveState
	DepthStencil *DepthStencilState
	Multisample  MultisampleState
	Fragment     *FragmentState
}

// ComputePipelineDescriptor describes how to create a compute pipeline.
type ComputePipelineDescriptor struct {
	Label   string
	Layout  gpu.PipelineLayout
	Compute ComputeState
}

// VertexState describes the vertex processing state.
type VertexState struct {
	Module     gpu.ShaderModule
	EntryPoint string
	Buffers    []VertexBufferLayout
}

// PrimitiveState describes primitive assembly and rasterization state.
type PrimitiveState struct {
	Topology         gpu.PrimitiveTopology
	StripIndexFormat gpu.IndexFormat
	FrontFace        gpu.FrontFace
	CullMode         gpu.CullMode
}

// DepthStencilState describes depth and stencil testing state.
type DepthStencilState struct {
	Format              gpu.TextureFormat
	DepthWriteEnabled   bool
	DepthCompare        gpu.CompareFunction
	StencilFront        StencilFaceState
	StencilBack         StencilFaceState
	StencilReadMask     uint32
	StencilWriteMask    uint32
	DepthBias           int32
	DepthBiasSlopeScale float32
	DepthBiasClamp      float32
}

// StencilFaceState describes stencil testing for a single face.
type StencilFaceState struct {
	Compare     gpu.CompareFunction
	FailOp      gpu.StencilOperation
	DepthFailOp gpu.StencilOperation
	PassOp      gpu.StencilOperation
}

// MultisampleState describes multisampling state.
type MultisampleState struct {
	Count                  uint32
	Mask                   uint32
	AlphaToCoverageEnabled bool
}

// FragmentState describes fragment processing state.
type FragmentState struct {
	Module     gpu.ShaderModule
	EntryPoint string
	Targets    []ColorTargetState
}

// ComputeState describes compute processing state.
type ComputeState struct {
	Module     gpu.ShaderModule
	EntryPoint string
}

// ColorTargetState describes color target blending and write state.
type ColorTargetState struct {
	Format    gpu.TextureFormat
	Blend     *BlendState
	WriteMask gpu.ColorWriteMask
}

// BlendState describes color and alpha blending state.
type BlendState struct {
	Color BlendComponent
	Alpha BlendComponent
}

// BlendComponent describes blending for a single component.
type BlendComponent struct {
	Operation gpu.BlendOperation
	SrcFactor gpu.BlendFactor
	DstFactor gpu.BlendFactor
}

// VertexBufferLayout describes the layout of a vertex buffer.
type VertexBufferLayout struct {
	ArrayStride uint64
	StepMode    gpu.VertexStepMode
	Attributes  []VertexAttribute
}

// VertexAttribute describes a single vertex attribute.
type VertexAttribute struct {
	Format         gpu.VertexFormat
	Offset         uint64
	ShaderLocation uint32
}

// ImpellerContext implements the Context interface using Impeller principles.
type ImpellerContext struct {
	device            gpu.Device
	queue             gpu.Queue
	capabilities      *Capabilities
	resourceAllocator ResourceAllocator
	shaderLibrary     ShaderLibrary
	pipelineLibrary   PipelineLibrary
	samplerLibrary    SamplerLibrary
	isShutdown        bool
}

// NewImpellerContext creates a new Impeller-style rendering context.
func NewImpellerContext(device gpu.Device) (*ImpellerContext, error) {
	if device == nil {
		return nil, fmt.Errorf("device cannot be nil")
	}

	queue := device.Queue()
	if queue == nil {
		return nil, fmt.Errorf("failed to get device queue")
	}

	// Initialize capabilities based on device features
	capabilities := &Capabilities{
		SupportsAdvancedBlends:    true, // TODO: Query from device
		MaxTextureSize:            8192, // TODO: Query from device
		SupportsCompute:           true, // TODO: Query from device
		SupportsTimestamps:        true, // TODO: Query from device
		SupportsMultisampling:     true, // TODO: Query from device
		MaxSampleCount:            4,    // TODO: Query from device
		SupportsInstancedDrawing:  true, // TODO: Query from device
		DefaultColorFormat:        gpu.TextureFormatBGRA8Unorm,
		DefaultDepthStencilFormat: gpu.TextureFormatDepth24PlusStencil8,
	}

	ctx := &ImpellerContext{
		device:       device,
		queue:        queue,
		capabilities: capabilities,
	}

	// Initialize sub-systems
	ctx.resourceAllocator = NewResourceAllocator(device)
	ctx.shaderLibrary = NewShaderLibrary(device)
	ctx.pipelineLibrary = NewPipelineLibrary(device, ctx.shaderLibrary)
	ctx.samplerLibrary = NewSamplerLibrary(device)

	return ctx, nil
}

// GetDevice implements Context.
func (c *ImpellerContext) GetDevice() gpu.Device {
	return c.device
}

// GetQueue implements Context.
func (c *ImpellerContext) GetQueue() gpu.Queue {
	return c.queue
}

// GetCapabilities implements Context.
func (c *ImpellerContext) GetCapabilities() *Capabilities {
	return c.capabilities
}

// GetResourceAllocator implements Context.
func (c *ImpellerContext) GetResourceAllocator() ResourceAllocator {
	return c.resourceAllocator
}

// GetShaderLibrary implements Context.
func (c *ImpellerContext) GetShaderLibrary() ShaderLibrary {
	return c.shaderLibrary
}

// GetPipelineLibrary implements Context.
func (c *ImpellerContext) GetPipelineLibrary() PipelineLibrary {
	return c.pipelineLibrary
}

// GetSamplerLibrary implements Context.
func (c *ImpellerContext) GetSamplerLibrary() SamplerLibrary {
	return c.samplerLibrary
}

// Shutdown implements Context.
func (c *ImpellerContext) Shutdown() error {
	if c.isShutdown {
		return nil
	}

	c.isShutdown = true

	// TODO: Properly shutdown sub-systems and release resources
	return nil
}
