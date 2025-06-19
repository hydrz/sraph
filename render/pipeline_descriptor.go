package render

// PipelineDescriptor is a base interface for pipeline descriptors
type PipelineDescriptor interface {
	// GetLabel returns the debug label for the pipeline
	GetLabel() string

	// IsValid returns true if the descriptor is valid
	IsValid() bool
}

// RenderPipelineDescriptor describes the configuration for a render pipeline
type RenderPipelineDescriptor struct {
	Label                   string
	VertexFunction          ShaderFunction
	FragmentFunction        ShaderFunction
	VertexDescriptor        VertexDescriptor
	ColorAttachmentCount    int
	ColorAttachments        []ColorAttachmentDescriptor
	DepthAttachmentFormat   PixelFormat
	StencilAttachmentFormat PixelFormat
	SampleCount             int
	AlphaToCoverageEnabled  bool
	AlphaToOneEnabled       bool
	RasterizationEnabled    bool
	WindingOrder            WindingOrder
	CullMode                CullMode
	FillMode                FillMode
	DepthClipMode           DepthClipMode
	DepthTestEnabled        bool
	DepthWriteEnabled       bool
	DepthCompareFunction    CompareFunction
	StencilTestEnabled      bool
	StencilFrontFace        StencilDescriptor
	StencilBackFace         StencilDescriptor
}

// ColorAttachmentDescriptor describes a color attachment for a render pipeline
type ColorAttachmentDescriptor struct {
	Format                      PixelFormat
	BlendingEnabled             bool
	SourceRGBBlendFactor        BlendFactor
	DestinationRGBBlendFactor   BlendFactor
	RGBBlendOperation           BlendOperation
	SourceAlphaBlendFactor      BlendFactor
	DestinationAlphaBlendFactor BlendFactor
	AlphaBlendOperation         BlendOperation
	WriteMask                   ColorWriteMask
}

// StencilDescriptor describes stencil test configuration
type StencilDescriptor struct {
	StencilCompareFunction    CompareFunction
	StencilFailureOperation   StencilOperation
	DepthFailureOperation     StencilOperation
	DepthStencilPassOperation StencilOperation
	ReadMask                  uint32
	WriteMask                 uint32
}

// Enumeration types for render pipeline configuration
type WindingOrder int

const (
	WindingOrderClockwise WindingOrder = iota
	WindingOrderCounterClockwise
)

type CullMode int

const (
	CullModeNone CullMode = iota
	CullModeFront
	CullModeBack
)

type FillMode int

const (
	FillModeFill FillMode = iota
	FillModeWireframe
	FillModePoint
)

type DepthClipMode int

const (
	DepthClipModeClip DepthClipMode = iota
	DepthClipModeClamp
)

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
	BlendFactorBlendColor
	BlendFactorOneMinusBlendColor
	BlendFactorBlendAlpha
	BlendFactorOneMinusBlendAlpha
)

type BlendOperation int

const (
	BlendOperationAdd BlendOperation = iota
	BlendOperationSubtract
	BlendOperationReverseSubtract
	BlendOperationMin
	BlendOperationMax
)

type ColorWriteMask int

const (
	ColorWriteMaskNone  ColorWriteMask = 0
	ColorWriteMaskRed   ColorWriteMask = 1 << 0
	ColorWriteMaskGreen ColorWriteMask = 1 << 1
	ColorWriteMaskBlue  ColorWriteMask = 1 << 2
	ColorWriteMaskAlpha ColorWriteMask = 1 << 3
	ColorWriteMaskAll   ColorWriteMask = ColorWriteMaskRed | ColorWriteMaskGreen | ColorWriteMaskBlue | ColorWriteMaskAlpha
)

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

// GetLabel returns the debug label for the pipeline
func (d *RenderPipelineDescriptor) GetLabel() string {
	return d.Label
}

// IsValid returns true if the descriptor is valid
func (d *RenderPipelineDescriptor) IsValid() bool {
	return d.VertexFunction != nil && d.FragmentFunction != nil
}

// RenderPipeline represents a render pipeline
type RenderPipeline interface {
	Pipeline

	// GetRenderPipelineDescriptor returns the render pipeline descriptor
	GetRenderPipelineDescriptor() RenderPipelineDescriptor
}

// RenderPipelineImpl is the default implementation of RenderPipeline
type RenderPipelineImpl struct {
	*PipelineImpl
	descriptor RenderPipelineDescriptor
}

// NewRenderPipeline creates a new render pipeline with the specified descriptor
func NewRenderPipeline(descriptor RenderPipelineDescriptor) RenderPipeline {
	return &RenderPipelineImpl{
		PipelineImpl: &PipelineImpl{
			label:   descriptor.Label,
			isValid: descriptor.IsValid(),
		},
		descriptor: descriptor,
	}
}

// GetRenderPipelineDescriptor returns the render pipeline descriptor
func (p *RenderPipelineImpl) GetRenderPipelineDescriptor() RenderPipelineDescriptor {
	return p.descriptor
}
