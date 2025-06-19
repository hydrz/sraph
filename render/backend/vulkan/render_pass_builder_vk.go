package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// RenderPassBuilderVK builds Vulkan render passes
type RenderPassBuilderVK struct {
	context     *ContextVK
	attachments []AttachmentDescriptor
	extent      vulkan.Extent2D
	multisampling bool
	samples     vulkan.SampleCountFlagBits
}

// AttachmentDescriptor describes a render pass attachment
type AttachmentDescriptor struct {
	Format      vulkan.Format
	LoadOp      vulkan.AttachmentLoadOp
	StoreOp     vulkan.AttachmentStoreOp
	StencilLoadOp  vulkan.AttachmentLoadOp
	StencilStoreOp vulkan.AttachmentStoreOp
	InitialLayout  vulkan.ImageLayout
	FinalLayout    vulkan.ImageLayout
	Type        AttachmentType
	ClearValue  vulkan.ClearValue
}

// NewRenderPassBuilderVK creates a new render pass builder
func NewRenderPassBuilderVK(context *ContextVK) *RenderPassBuilderVK {
	return &RenderPassBuilderVK{
		context:     context,
		attachments: make([]AttachmentDescriptor, 0),
		samples:     vulkan.SampleCount1Bit,
	}
}

// SetExtent sets the render pass extent
func (builder *RenderPassBuilderVK) SetExtent(width, height uint32) *RenderPassBuilderVK {
	builder.extent = vulkan.Extent2D{
		Width:  width,
		Height: height,
	}
	return builder
}

// SetMultisampling enables multisampling with the specified sample count
func (builder *RenderPassBuilderVK) SetMultisampling(samples vulkan.SampleCountFlagBits) *RenderPassBuilderVK {
	builder.multisampling = true
	builder.samples = samples
	return builder
}

// AddColorAttachment adds a color attachment to the render pass
func (builder *RenderPassBuilderVK) AddColorAttachment(format vulkan.Format, loadOp vulkan.AttachmentLoadOp, storeOp vulkan.AttachmentStoreOp) *RenderPassBuilderVK {
	attachment := AttachmentDescriptor{
		Format:        format,
		LoadOp:        loadOp,
		StoreOp:       storeOp,
		StencilLoadOp: vulkan.AttachmentLoadOpDontCare,
		StencilStoreOp: vulkan.AttachmentStoreOpDontCare,
		InitialLayout: vulkan.ImageLayoutUndefined,
		FinalLayout:   vulkan.ImageLayoutPresentSrcKhr,
		Type:          AttachmentTypeColor,
	}

	// Set clear color to black
	attachment.ClearValue.SetColor([]float32{0.0, 0.0, 0.0, 1.0})

	builder.attachments = append(builder.attachments, attachment)
	return builder
}

// AddColorAttachmentWithClear adds a color attachment with custom clear color
func (builder *RenderPassBuilderVK) AddColorAttachmentWithClear(format vulkan.Format, loadOp vulkan.AttachmentLoadOp, storeOp vulkan.AttachmentStoreOp, clearColor [4]float32) *RenderPassBuilderVK {
	attachment := AttachmentDescriptor{
		Format:        format,
		LoadOp:        loadOp,
		StoreOp:       storeOp,
		StencilLoadOp: vulkan.AttachmentLoadOpDontCare,
		StencilStoreOp: vulkan.AttachmentStoreOpDontCare,
		InitialLayout: vulkan.ImageLayoutUndefined,
		FinalLayout:   vulkan.ImageLayoutPresentSrcKhr,
		Type:          AttachmentTypeColor,
	}

	attachment.ClearValue.SetColor(clearColor[:])

	builder.attachments = append(builder.attachments, attachment)
	return builder
}

// AddDepthAttachment adds a depth attachment to the render pass
func (builder *RenderPassBuilderVK) AddDepthAttachment(format vulkan.Format, loadOp vulkan.AttachmentLoadOp, storeOp vulkan.AttachmentStoreOp) *RenderPassBuilderVK {
	attachment := AttachmentDescriptor{
		Format:        format,
		LoadOp:        loadOp,
		StoreOp:       storeOp,
		StencilLoadOp: vulkan.AttachmentLoadOpDontCare,
		StencilStoreOp: vulkan.AttachmentStoreOpDontCare,
		InitialLayout: vulkan.ImageLayoutUndefined,
		FinalLayout:   vulkan.ImageLayoutDepthStencilAttachmentOptimal,
		Type:          AttachmentTypeDepth,
	}

	// Set clear depth to 1.0
	attachment.ClearValue.SetDepthStencil(1.0, 0)

	builder.attachments = append(builder.attachments, attachment)
	return builder
}

// AddDepthStencilAttachment adds a depth-stencil attachment to the render pass
func (builder *RenderPassBuilderVK) AddDepthStencilAttachment(format vulkan.Format, loadOp vulkan.AttachmentLoadOp, storeOp vulkan.AttachmentStoreOp, stencilLoadOp vulkan.AttachmentLoadOp, stencilStoreOp vulkan.AttachmentStoreOp) *RenderPassBuilderVK {
	attachment := AttachmentDescriptor{
		Format:        format,
		LoadOp:        loadOp,
		StoreOp:       storeOp,
		StencilLoadOp: stencilLoadOp,
		StencilStoreOp: stencilStoreOp,
		InitialLayout: vulkan.ImageLayoutUndefined,
		FinalLayout:   vulkan.ImageLayoutDepthStencilAttachmentOptimal,
		Type:          AttachmentTypeDepthStencil,
	}

	// Set clear depth to 1.0 and stencil to 0
	attachment.ClearValue.SetDepthStencil(1.0, 0)

	builder.attachments = append(builder.attachments, attachment)
	return builder
}

// Build creates the render pass from the builder configuration
func (builder *RenderPassBuilderVK) Build() (*RenderPassVK, error) {
	if len(builder.attachments) == 0 {
		return nil, fmt.Errorf("no attachments specified")
	}

	if builder.extent.Width == 0 || builder.extent.Height == 0 {
		return nil, fmt.Errorf("invalid extent specified")
	}

	// Create attachment descriptions
	attachmentDescriptions := make([]vulkan.AttachmentDescription, len(builder.attachments))
	colorAttachmentRefs := make([]vulkan.AttachmentReference, 0)
	var depthAttachmentRef *vulkan.AttachmentReference

	for i, attachment := range builder.attachments {
		attachmentDescriptions[i] = vulkan.AttachmentDescription{
			Format:         attachment.Format,
			Samples:        builder.samples,
			LoadOp:         attachment.LoadOp,
			StoreOp:        attachment.StoreOp,
			StencilLoadOp:  attachment.StencilLoadOp,
			StencilStoreOp: attachment.StencilStoreOp,
			InitialLayout:  attachment.InitialLayout,
			FinalLayout:    attachment.FinalLayout,
		}

		attachmentRef := vulkan.AttachmentReference{
			Attachment: uint32(i),
		}

		switch attachment.Type {
		case AttachmentTypeColor:
			attachmentRef.Layout = vulkan.ImageLayoutColorAttachmentOptimal
			colorAttachmentRefs = append(colorAttachmentRefs, attachmentRef)
		case AttachmentTypeDepth, AttachmentTypeDepthStencil:
			attachmentRef.Layout = vulkan.ImageLayoutDepthStencilAttachmentOptimal
			depthAttachmentRef = &attachmentRef
		}
	}

	// Create subpass description
	subpass := vulkan.SubpassDescription{
		PipelineBindPoint:       vulkan.PipelineBindPointGraphics,
		ColorAttachmentCount:    uint32(len(colorAttachmentRefs)),
		PColorAttachments:       colorAttachmentRefs,
		PDepthStencilAttachment: depthAttachmentRef,
	}

	// Create subpass dependencies
	dependencies := []vulkan.SubpassDependency{
		{
			SrcSubpass:    vulkan.SubpassExternal,
			DstSubpass:    0,
			SrcStageMask:  vulkan.PipelineStageColorAttachmentOutputBit | vulkan.PipelineStageEarlyFragmentTestsBit,
			SrcAccessMask: 0,
			DstStageMask:  vulkan.PipelineStageColorAttachmentOutputBit | vulkan.PipelineStageEarlyFragmentTestsBit,
			DstAccessMask: vulkan.AccessColorAttachmentWriteBit | vulkan.AccessDepthStencilAttachmentWriteBit,
		},
	}

	// Create render pass
	renderPassInfo := vulkan.RenderPassCreateInfo{
		SType:           vulkan.StructureTypeRenderPassCreateInfo,
		AttachmentCount: uint32(len(attachmentDescriptions)),
		PAttachments:    attachmentDescriptions,
		SubpassCount:    1,
		PSubpasses:      []vulkan.SubpassDescription{subpass},
		DependencyCount: uint32(len(dependencies)),
		PDependencies:   dependencies,
	}

	var renderPass vulkan.RenderPass
	if result := vulkan.CreateRenderPass(builder.context.device, &renderPassInfo, nil, &renderPass); result != vulkan.Success {
		return nil, fmt.Errorf("failed to create render pass: %s", result)
	}

	// Create render pass object
	rp := &RenderPassVK{
		context:    builder.context,
		renderPass: renderPass,
		extent:     builder.extent,
	}

	return rp, nil
}

// Reset resets the builder to its initial state
func (builder *RenderPassBuilderVK) Reset() *RenderPassBuilderVK {
	builder.attachments = builder.attachments[:0]
	builder.extent = vulkan.Extent2D{}
	builder.multisampling = false
	builder.samples = vulkan.SampleCount1Bit
	return builder
}

// GetAttachmentCount returns the number of attachments
func (builder *RenderPassBuilderVK) GetAttachmentCount() int {
	return len(builder.attachments)
}

// GetColorAttachmentCount returns the number of color attachments
func (builder *RenderPassBuilderVK) GetColorAttachmentCount() int {
	count := 0
	for _, attachment := range builder.attachments {
		if attachment.Type == AttachmentTypeColor {
			count++
		}
	}
	return count
}

// HasDepthAttachment returns true if there is a depth attachment
func (builder *RenderPassBuilderVK) HasDepthAttachment() bool {
	for _, attachment := range builder.attachments {
		if attachment.Type == AttachmentTypeDepth || attachment.Type == AttachmentTypeDepthStencil {
			return true
		}
	}
	return false
}
