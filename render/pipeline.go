// Pipeline interface and implementation for the render package.
// This provides graphics and compute pipeline functionality.

package render

import (
	"errors"
	"fmt"
	"sync"

	"github.com/opensraph/sraph/gpu"
)

// Pipeline represents a graphics or compute pipeline.
type Pipeline interface {
	// IsValid returns true if the pipeline is valid and can be used.
	IsValid() bool

	// GetDescriptor returns the pipeline descriptor.
	GetDescriptor() PipelineDescriptor

	// SetLabel sets a debug label for the pipeline.
	SetLabel(label string)
}

// PipelineType represents the type of pipeline.
type PipelineType int

const (
	PipelineTypeUnknown PipelineType = iota
	PipelineTypeRender
	PipelineTypeCompute
)

// PipelineDescriptor is the base interface for pipeline descriptors.
type PipelineDescriptor interface {
	// GetType returns the pipeline type.
	GetType() PipelineType

	// GetLabel returns the pipeline label.
	GetLabel() string

	// Validate validates the pipeline descriptor.
	Validate() error
}

// RenderPipeline represents a graphics rendering pipeline.
type RenderPipeline interface {
	Pipeline

	// GetVertexDescriptor returns the vertex descriptor.
	GetVertexDescriptor() *VertexDescriptor

	// GetRenderTargetPixelFormat returns the render target pixel format.
	GetRenderTargetPixelFormat() PixelFormat

	// GetSampleCount returns the sample count for multisampling.
	GetSampleCount() uint32

	// HasDepthAttachment returns true if the pipeline has a depth attachment.
	HasDepthAttachment() bool

	// HasStencilAttachment returns true if the pipeline has a stencil attachment.
	HasStencilAttachment() bool
}

// RenderPipelineDescriptor describes a render pipeline.
type RenderPipelineDescriptor struct {
	Label                   string
	VertexShader            ShaderFunction
	FragmentShader          ShaderFunction
	VertexDescriptor        *VertexDescriptor
	ColorAttachmentCount    uint32
	ColorAttachmentFormats  []PixelFormat
	HasDepthAttachment      bool
	DepthAttachmentFormat   PixelFormat
	HasStencilAttachment    bool
	StencilAttachmentFormat PixelFormat
	SampleCount             uint32
	DepthStencilDescriptor  *DepthStencilDescriptor
	BlendDescriptor         *BlendDescriptor
	WindingOrder            WindingOrder
	CullMode                CullMode
	PolygonMode             PolygonMode
}

// GetType returns the pipeline type.
func (rpd *RenderPipelineDescriptor) GetType() PipelineType {
	return PipelineTypeRender
}

// GetLabel returns the pipeline label.
func (rpd *RenderPipelineDescriptor) GetLabel() string {
	return rpd.Label
}

// Validate validates the pipeline descriptor.
func (rpd *RenderPipelineDescriptor) Validate() error {
	if rpd.VertexShader == nil {
		return errors.New("vertex shader cannot be nil")
	}
	if rpd.FragmentShader == nil {
		return errors.New("fragment shader cannot be nil")
	}
	if rpd.ColorAttachmentCount == 0 {
		return errors.New("color attachment count must be greater than 0")
	}
	if uint32(len(rpd.ColorAttachmentFormats)) != rpd.ColorAttachmentCount {
		return errors.New("color attachment formats count must match color attachment count")
	}
	if rpd.SampleCount == 0 {
		rpd.SampleCount = 1
	}
	return nil
}

// ComputePipelineDescriptor describes a compute pipeline.
type ComputePipelineDescriptor struct {
	Label           string
	ComputeShader   ShaderFunction
	ThreadgroupSize ThreadgroupSize
}

// GetType returns the pipeline type.
func (cpd *ComputePipelineDescriptor) GetType() PipelineType {
	return PipelineTypeCompute
}

// GetLabel returns the pipeline label.
func (cpd *ComputePipelineDescriptor) GetLabel() string {
	return cpd.Label
}

// Validate validates the pipeline descriptor.
func (cpd *ComputePipelineDescriptor) Validate() error {
	if cpd.ComputeShader == nil {
		return errors.New("compute shader cannot be nil")
	}
	return nil
}

// ThreadgroupSize represents the size of a compute threadgroup.
type ThreadgroupSize struct {
	Width  uint32
	Height uint32
	Depth  uint32
}

// WindingOrder represents the winding order for front-facing primitives.
type WindingOrder int

const (
	WindingOrderClockwise WindingOrder = iota
	WindingOrderCounterClockwise
)

// CullMode represents the culling mode for back-face culling.
type CullMode int

const (
	CullModeNone CullMode = iota
	CullModeFront
	CullModeBack
)

// PolygonMode represents the polygon rasterization mode.
type PolygonMode int

const (
	PolygonModeFill PolygonMode = iota
	PolygonModeLine
	PolygonModePoint
)

// DepthStencilDescriptor describes depth and stencil testing.
type DepthStencilDescriptor struct {
	DepthTestEnabled     bool
	DepthWriteEnabled    bool
	DepthCompareFunction CompareFunction
	StencilTestEnabled   bool
	StencilFrontFace     StencilDescriptor
	StencilBackFace      StencilDescriptor
	StencilReadMask      uint32
	StencilWriteMask     uint32
}

// StencilDescriptor describes stencil testing for a face.
type StencilDescriptor struct {
	StencilCompareFunction    CompareFunction
	StencilFailureOperation   StencilOperation
	DepthFailureOperation     StencilOperation
	DepthStencilPassOperation StencilOperation
}

// CompareFunction represents a comparison function.
type CompareFunction int

const (
	CompareFunctionNever CompareFunction = iota
	CompareFunctionLess
	CompareFunctionEqual
	CompareFunctionLessEqual
	CompareFunctionGreater
	CompareFunctionNotEqual
	CompareFunctionGreaterEqual
	CompareFunctionAlways
)

// StencilOperation represents a stencil operation.
type StencilOperation int

const (
	StencilOperationKeep StencilOperation = iota
	StencilOperationZero
	StencilOperationReplace
	StencilOperationIncrementClamp
	StencilOperationDecrementClamp
	StencilOperationInvert
	StencilOperationIncrementWrap
	StencilOperationDecrementWrap
)

// BlendDescriptor describes color blending.
type BlendDescriptor struct {
	BlendEnabled                bool
	SourceColorBlendFactor      BlendFactor
	DestinationColorBlendFactor BlendFactor
	ColorBlendOperation         BlendOperation
	SourceAlphaBlendFactor      BlendFactor
	DestinationAlphaBlendFactor BlendFactor
	AlphaBlendOperation         BlendOperation
	ColorWriteMask              ColorWriteMask
}

// BlendFactor represents a blend factor.
type BlendFactor int

const (
	BlendFactorZero BlendFactor = iota
	BlendFactorOne
	BlendFactorSourceColor
	BlendFactorOneMinusSourceColor
	BlendFactorDestinationColor
	BlendFactorOneMinusDestinationColor
	BlendFactorSourceAlpha
	BlendFactorOneMinusSourceAlpha
	BlendFactorDestinationAlpha
	BlendFactorOneMinusDestinationAlpha
	BlendFactorSourceAlphaSaturated
)

// BlendOperation represents a blend operation.
type BlendOperation int

const (
	BlendOperationAdd BlendOperation = iota
	BlendOperationSubtract
	BlendOperationReverseSubtract
	BlendOperationMin
	BlendOperationMax
)

// ColorWriteMask represents which color components to write.
type ColorWriteMask int

const (
	ColorWriteMaskNone  ColorWriteMask = 0
	ColorWriteMaskRed   ColorWriteMask = 1 << 0
	ColorWriteMaskGreen ColorWriteMask = 1 << 1
	ColorWriteMaskBlue  ColorWriteMask = 1 << 2
	ColorWriteMaskAlpha ColorWriteMask = 1 << 3
	ColorWriteMaskAll   ColorWriteMask = ColorWriteMaskRed | ColorWriteMaskGreen | ColorWriteMaskBlue | ColorWriteMaskAlpha
)

// DefaultRenderPipeline is the default implementation of RenderPipeline.
type DefaultRenderPipeline struct {
	device     gpu.Device
	pipeline   gpu.RenderPipeline
	descriptor *RenderPipelineDescriptor
	label      string
	isValid    bool
	mutex      sync.RWMutex
}

// NewRenderPipeline creates a new render pipeline with the given descriptor.
func NewRenderPipeline(device gpu.Device, descriptor *RenderPipelineDescriptor) (RenderPipeline, error) {
	if device == nil {
		return nil, errors.New("device cannot be nil")
	}
	if descriptor == nil {
		return nil, errors.New("descriptor cannot be nil")
	}
	if err := descriptor.Validate(); err != nil {
		return nil, fmt.Errorf("invalid descriptor: %w", err)
	}

	pipeline := &DefaultRenderPipeline{
		device:     device,
		descriptor: descriptor,
		label:      descriptor.Label,
		isValid:    false,
	}

	// Create the underlying GPU render pipeline
	gpuPipeline, err := device.CreateRenderPipeline(&gpu.RenderPipelineDescriptor{
		Label: descriptor.Label,
		// Additional GPU-specific descriptor fields would be set here
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create GPU render pipeline: %w", err)
	}

	pipeline.pipeline = gpuPipeline
	pipeline.isValid = true

	return pipeline, nil
}

// IsValid returns true if the pipeline is valid and can be used.
func (rp *DefaultRenderPipeline) IsValid() bool {
	rp.mutex.RLock()
	defer rp.mutex.RUnlock()
	return rp.isValid && rp.pipeline != nil
}

// GetDescriptor returns the pipeline descriptor.
func (rp *DefaultRenderPipeline) GetDescriptor() PipelineDescriptor {
	rp.mutex.RLock()
	defer rp.mutex.RUnlock()
	return rp.descriptor
}

// SetLabel sets a debug label for the pipeline.
func (rp *DefaultRenderPipeline) SetLabel(label string) {
	rp.mutex.Lock()
	defer rp.mutex.Unlock()

	rp.label = label
	if rp.pipeline != nil {
		rp.pipeline.SetLabel(label)
	}
}

// GetVertexDescriptor returns the vertex descriptor.
func (rp *DefaultRenderPipeline) GetVertexDescriptor() *VertexDescriptor {
	rp.mutex.RLock()
	defer rp.mutex.RUnlock()
	return rp.descriptor.VertexDescriptor
}

// GetRenderTargetPixelFormat returns the render target pixel format.
func (rp *DefaultRenderPipeline) GetRenderTargetPixelFormat() PixelFormat {
	rp.mutex.RLock()
	defer rp.mutex.RUnlock()
	if len(rp.descriptor.ColorAttachmentFormats) > 0 {
		return rp.descriptor.ColorAttachmentFormats[0]
	}
	return PixelFormatUnknown
}

// GetSampleCount returns the sample count for multisampling.
func (rp *DefaultRenderPipeline) GetSampleCount() uint32 {
	rp.mutex.RLock()
	defer rp.mutex.RUnlock()
	return rp.descriptor.SampleCount
}

// HasDepthAttachment returns true if the pipeline has a depth attachment.
func (rp *DefaultRenderPipeline) HasDepthAttachment() bool {
	rp.mutex.RLock()
	defer rp.mutex.RUnlock()
	return rp.descriptor.HasDepthAttachment
}

// HasStencilAttachment returns true if the pipeline has a stencil attachment.
func (rp *DefaultRenderPipeline) HasStencilAttachment() bool {
	rp.mutex.RLock()
	defer rp.mutex.RUnlock()
	return rp.descriptor.HasStencilAttachment
}

// DefaultComputePipeline is the default implementation of ComputePipeline.
type DefaultComputePipeline struct {
	device     gpu.Device
	pipeline   gpu.ComputePipeline
	descriptor *ComputePipelineDescriptor
	label      string
	isValid    bool
	mutex      sync.RWMutex
}

// NewComputePipeline creates a new compute pipeline with the given descriptor.
func NewComputePipeline(device gpu.Device, descriptor *ComputePipelineDescriptor) (ComputePipeline, error) {
	if device == nil {
		return nil, errors.New("device cannot be nil")
	}
	if descriptor == nil {
		return nil, errors.New("descriptor cannot be nil")
	}
	if err := descriptor.Validate(); err != nil {
		return nil, fmt.Errorf("invalid descriptor: %w", err)
	}

	pipeline := &DefaultComputePipeline{
		device:     device,
		descriptor: descriptor,
		label:      descriptor.Label,
		isValid:    false,
	}

	// Create the underlying GPU compute pipeline
	gpuPipeline, err := device.CreateComputePipeline(&gpu.ComputePipelineDescriptor{
		Label: descriptor.Label,
		// Additional GPU-specific descriptor fields would be set here
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create GPU compute pipeline: %w", err)
	}

	pipeline.pipeline = gpuPipeline
	pipeline.isValid = true

	return pipeline, nil
}

// IsValid returns true if the pipeline is valid and can be used.
func (cp *DefaultComputePipeline) IsValid() bool {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()
	return cp.isValid && cp.pipeline != nil
}

// GetDescriptor returns the pipeline descriptor.
func (cp *DefaultComputePipeline) GetDescriptor() PipelineDescriptor {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()
	return cp.descriptor
}

// SetLabel sets a debug label for the pipeline.
func (cp *DefaultComputePipeline) SetLabel(label string) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	cp.label = label
	if cp.pipeline != nil {
		cp.pipeline.SetLabel(label)
	}
}
