package impl

import . "github.com/opensraph/sraph/gpu/wgpu"

// Common shader stage combinations
const (
	ShaderStageVertexFragment ShaderStage = ShaderStageVertex | ShaderStageFragment
	ShaderStageAll            ShaderStage = ShaderStageVertex | ShaderStageFragment | ShaderStageCompute
)

// Common buffer usage combinations
const (
	BufferUsageVertexIndex    BufferUsage = BufferUsageVertex | BufferUsageIndex
	BufferUsageUniformStorage BufferUsage = BufferUsageUniform | BufferUsageStorage
	BufferUsageCopyAll        BufferUsage = BufferUsageCopySrc | BufferUsageCopyDst
)

// Common texture usage combinations
const (
	TextureUsageRenderTarget TextureUsage = TextureUsageRenderAttachment | TextureUsageTextureBinding
	TextureUsageStorageRead  TextureUsage = TextureUsageStorageBinding | TextureUsageTextureBinding
	TextureUsageCopyAll      TextureUsage = TextureUsageCopySrc | TextureUsageCopyDst
)

// Default values for common descriptors
var (
	DefaultSamplerDescriptor = SamplerDescriptor{
		Label:         "Default Sampler",
		AddressModeU:  AddressModeClampToEdge,
		AddressModeV:  AddressModeClampToEdge,
		AddressModeW:  AddressModeClampToEdge,
		MagFilter:     FilterModeLinear,
		MinFilter:     FilterModeLinear,
		MipmapFilter:  MipmapFilterModeLinear,
		LodMinClamp:   0.0,
		LodMaxClamp:   32.0,
		Compare:       CompareFunctionUndefined,
		MaxAnisotropy: 1,
	}

	DefaultBlendState = BlendState{
		Color: BlendComponent{
			Operation: BlendOperationAdd,
			SrcFactor: BlendFactorOne,
			DstFactor: BlendFactorZero,
		},
		Alpha: BlendComponent{
			Operation: BlendOperationAdd,
			SrcFactor: BlendFactorOne,
			DstFactor: BlendFactorZero,
		},
	}

	DefaultMultisampleState = MultisampleState{
		Count:                  1,
		Mask:                   0xFFFFFFFF,
		AlphaToCoverageEnabled: false,
	}

	DefaultPrimitiveState = PrimitiveState{
		Topology:         PrimitiveTopologyTriangleList,
		StripIndexFormat: IndexFormatUndefined,
		FrontFace:        FrontFaceCCW,
		CullMode:         CullModeNone,
		UnclippedDepth:   false,
	}
)

// Predefined color constants
var (
	ColorTransparent = Color{R: 0.0, G: 0.0, B: 0.0, A: 0.0}
	ColorBlack       = Color{R: 0.0, G: 0.0, B: 0.0, A: 1.0}
	ColorWhite       = Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0}
	ColorRed         = Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0}
	ColorGreen       = Color{R: 0.0, G: 1.0, B: 0.0, A: 1.0}
	ColorBlue        = Color{R: 0.0, G: 0.0, B: 1.0, A: 1.0}
)

// Common texture formats grouped by usage
var (
	UnormFormats = []TextureFormat{
		TextureFormatR8Unorm,
		TextureFormatRG8Unorm,
		TextureFormatRGBA8Unorm,
		TextureFormatBGRA8Unorm,
		TextureFormatRGB10A2Unorm,
	}

	FloatFormats = []TextureFormat{
		TextureFormatR16Float,
		TextureFormatRG16Float,
		TextureFormatRGBA16Float,
		TextureFormatR32Float,
		TextureFormatRG32Float,
		TextureFormatRGBA32Float,
	}

	DepthFormats = []TextureFormat{
		TextureFormatDepth16Unorm,
		TextureFormatDepth24Plus,
		TextureFormatDepth32Float,
		TextureFormatDepth24PlusStencil8,
		TextureFormatDepth32FloatStencil8,
	}

	CompressedFormats = []TextureFormat{
		TextureFormatBC1RGBAUnorm,
		TextureFormatBC2RGBAUnorm,
		TextureFormatBC3RGBAUnorm,
		TextureFormatBC4RUnorm,
		TextureFormatBC5RGUnorm,
		TextureFormatBC6HRGBUfloat,
		TextureFormatBC7RGBAUnorm,
	}
)
