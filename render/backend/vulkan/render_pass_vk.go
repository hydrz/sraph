package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// RenderPassVK represents a Vulkan render pass implementation
type RenderPassVK struct {
	context      *ContextVK
	renderPass   vulkan.RenderPass
	framebuffer  vulkan.Framebuffer
	attachments  []AttachmentVK
	extent       vulkan.Extent2D
	clearValues  []vulkan.ClearValue
	isActive     bool
	commandBuffer vulkan.CommandBuffer
}

// AttachmentVK represents a render pass attachment
type AttachmentVK struct {
	Texture    *TextureVK
	LoadOp     vulkan.AttachmentLoadOp
	StoreOp    vulkan.AttachmentStoreOp
	ClearValue vulkan.ClearValue
	Type       AttachmentType
}

// AttachmentType defines the type of attachment
type AttachmentType int

const (
	AttachmentTypeColor AttachmentType = iota
	AttachmentTypeDepth
	AttachmentTypeStencil
	AttachmentTypeDepthStencil
)

// NewRenderPassVK creates a new Vulkan render pass
func NewRenderPassVK(context *ContextVK, descriptor *RenderPassDescriptor) (*RenderPassVK, error) {
	if context == nil || descriptor == nil {
		return nil, fmt.Errorf("invalid parameters")
	}

	rp := &RenderPassVK{
		context:     context,
		attachments: make([]AttachmentVK, len(descriptor.Attachments)),
		extent:      descriptor.Extent,
		clearValues: make([]vulkan.ClearValue, len(descriptor.Attachments)),
	}

	// Copy attachments
	copy(rp.attachments, descriptor.Attachments)

	// Prepare clear values
	for i, attachment := range rp.attachments {
		rp.clearValues[i] = attachment.ClearValue
	}

	// Create render pass
	if err := rp.createRenderPass(); err != nil {
		return nil, fmt.Errorf("failed to create render pass: %w", err)
	}

	// Create framebuffer
	if err := rp.createFramebuffer(); err != nil {
		rp.Destroy()
		return nil, fmt.Errorf("failed to create framebuffer: %w", err)
	}

	return rp, nil
}

// createRenderPass creates the Vulkan render pass object
func (rp *RenderPassVK) createRenderPass() error {
	attachmentDescriptions := make([]vulkan.AttachmentDescription, len(rp.attachments))
	colorAttachmentRefs := make([]vulkan.AttachmentReference, 0)
	var depthAttachmentRef *vulkan.AttachmentReference

	// Create attachment descriptions and references
	for i, attachment := range rp.attachments {
		attachmentDescriptions[i] = vulkan.AttachmentDescription{
			Format:         attachment.Texture.GetFormat(),
			Samples:        vulkan.SampleCount1Bit,
			LoadOp:         attachment.LoadOp,
			StoreOp:        attachment.StoreOp,
			StencilLoadOp:  vulkan.AttachmentLoadOpDontCare,
			StencilStoreOp: vulkan.AttachmentStoreOpDontCare,
			InitialLayout:  vulkan.ImageLayoutUndefined,
			FinalLayout:    vulkan.ImageLayoutPresentSrcKhr,
		}

		attachmentRef := vulkan.AttachmentReference{
			Attachment: uint32(i),
			Layout:     vulkan.ImageLayoutColorAttachmentOptimal,
		}

		switch attachment.Type {
		case AttachmentTypeColor:
			attachmentDescriptions[i].FinalLayout = vulkan.ImageLayoutColorAttachmentOptimal
			colorAttachmentRefs = append(colorAttachmentRefs, attachmentRef)
		case AttachmentTypeDepth, AttachmentTypeDepthStencil:
			attachmentDescriptions[i].FinalLayout = vulkan.ImageLayoutDepthStencilAttachmentOptimal
			attachmentRef.Layout = vulkan.ImageLayoutDepthStencilAttachmentOptimal
			depthAttachmentRef = &attachmentRef
		}
	}

	// Create subpass description
	subpass := vulkan.SubpassDescription{
		PipelineBindPoint:    vulkan.PipelineBindPointGraphics,
		ColorAttachmentCount: uint32(len(colorAttachmentRefs)),
		PColorAttachments:    colorAttachmentRefs,
		PDepthStencilAttachment: depthAttachmentRef,
	}

	// Create subpass dependency
	dependency := vulkan.SubpassDependency{
		SrcSubpass:    vulkan.SubpassExternal,
		DstSubpass:    0,
		SrcStageMask:  vulkan.PipelineStageColorAttachmentOutputBit,
		SrcAccessMask: 0,
		DstStageMask:  vulkan.PipelineStageColorAttachmentOutputBit,
		DstAccessMask: vulkan.AccessColorAttachmentWriteBit,
	}

	// Create render pass
	renderPassInfo := vulkan.RenderPassCreateInfo{
		SType:           vulkan.StructureTypeRenderPassCreateInfo,
		AttachmentCount: uint32(len(attachmentDescriptions)),
		PAttachments:    attachmentDescriptions,
		SubpassCount:    1,
		PSubpasses:      []vulkan.SubpassDescription{subpass},
		DependencyCount: 1,
		PDependencies:   []vulkan.SubpassDependency{dependency},
	}

	var renderPass vulkan.RenderPass
	if result := vulkan.CreateRenderPass(rp.context.device, &renderPassInfo, nil, &renderPass); result != vulkan.Success {
		return fmt.Errorf("failed to create render pass: %s", result)
	}

	rp.renderPass = renderPass
	return nil
}

// createFramebuffer creates the framebuffer
func (rp *RenderPassVK) createFramebuffer() error {
	imageViews := make([]vulkan.ImageView, len(rp.attachments))
	for i, attachment := range rp.attachments {
		imageViews[i] = attachment.Texture.GetImageView()
	}

	framebufferInfo := vulkan.FramebufferCreateInfo{
		SType:           vulkan.StructureTypeFramebufferCreateInfo,
		RenderPass:      rp.renderPass,
		AttachmentCount: uint32(len(imageViews)),
		PAttachments:    imageViews,
		Width:           rp.extent.Width,
		Height:          rp.extent.Height,
		Layers:          1,
	}

	var framebuffer vulkan.Framebuffer
	if result := vulkan.CreateFramebuffer(rp.context.device, &framebufferInfo, nil, &framebuffer); result != vulkan.Success {
		return fmt.Errorf("failed to create framebuffer: %s", result)
	}

	rp.framebuffer = framebuffer
	return nil
}

// Begin begins the render pass
func (rp *RenderPassVK) Begin(commandBuffer vulkan.CommandBuffer) error {
	if rp.isActive {
		return fmt.Errorf("render pass is already active")
	}

	rp.commandBuffer = commandBuffer

	renderPassBeginInfo := vulkan.RenderPassBeginInfo{
		SType:       vulkan.StructureTypeRenderPassBeginInfo,
		RenderPass:  rp.renderPass,
		Framebuffer: rp.framebuffer,
		RenderArea: vulkan.Rect2D{
			Offset: vulkan.Offset2D{X: 0, Y: 0},
			Extent: rp.extent,
		},
		ClearValueCount: uint32(len(rp.clearValues)),
		PClearValues:    rp.clearValues,
	}

	vulkan.CmdBeginRenderPass(commandBuffer, &renderPassBeginInfo, vulkan.SubpassContentsInline)
	rp.isActive = true

	return nil
}

// End ends the render pass
func (rp *RenderPassVK) End() error {
	if !rp.isActive {
		return fmt.Errorf("render pass is not active")
	}

	vulkan.CmdEndRenderPass(rp.commandBuffer)
	rp.isActive = false
	rp.commandBuffer = vulkan.NullCommandBuffer

	return nil
}

// SetViewport sets the viewport for the render pass
func (rp *RenderPassVK) SetViewport(viewport vulkan.Viewport) {
	if rp.isActive && rp.commandBuffer != vulkan.NullCommandBuffer {
		vulkan.CmdSetViewport(rp.commandBuffer, 0, 1, []vulkan.Viewport{viewport})
	}
}

// SetScissor sets the scissor rectangle for the render pass
func (rp *RenderPassVK) SetScissor(scissor vulkan.Rect2D) {
	if rp.isActive && rp.commandBuffer != vulkan.NullCommandBuffer {
		vulkan.CmdSetScissor(rp.commandBuffer, 0, 1, []vulkan.Rect2D{scissor})
	}
}

// GetRenderPass returns the Vulkan render pass handle
func (rp *RenderPassVK) GetRenderPass() vulkan.RenderPass {
	return rp.renderPass
}

// GetFramebuffer returns the Vulkan framebuffer handle
func (rp *RenderPassVK) GetFramebuffer() vulkan.Framebuffer {
	return rp.framebuffer
}

// GetExtent returns the render pass extent
func (rp *RenderPassVK) GetExtent() vulkan.Extent2D {
	return rp.extent
}

// IsActive returns true if the render pass is currently active
func (rp *RenderPassVK) IsActive() bool {
	return rp.isActive
}

// Destroy destroys the render pass and frees resources
func (rp *RenderPassVK) Destroy() {
	if rp.context == nil {
		return
	}

	if rp.framebuffer != vulkan.NullFramebuffer {
		vulkan.DestroyFramebuffer(rp.context.device, rp.framebuffer, nil)
		rp.framebuffer = vulkan.NullFramebuffer
	}

	if rp.renderPass != vulkan.NullRenderPass {
		vulkan.DestroyRenderPass(rp.context.device, rp.renderPass, nil)
		rp.renderPass = vulkan.NullRenderPass
	}
}

// RenderPassDescriptor describes render pass creation parameters
type RenderPassDescriptor struct {
	Attachments []AttachmentVK
	Extent      vulkan.Extent2D
}
