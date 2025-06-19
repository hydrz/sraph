package render

import (
	"fmt"

	"github.com/opensraph/sraph/display"
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/gpu"
)

// RenderTarget represents a render target configuration.
type RenderTarget struct {
	colorAttachments []ColorAttachment
	depthAttachment  *DepthAttachment
	size             geom.Size[geom.F32]
	sampleCount      uint32
}

// ColorAttachment describes a color render target attachment.
type ColorAttachment struct {
	Texture     gpu.Texture
	View        gpu.TextureView
	LoadOp      gpu.LoadOp
	StoreOp     gpu.StoreOp
	ClearValue  gpu.Color
	ResolveView gpu.TextureView
}

// DepthAttachment describes a depth/stencil render target attachment.
type DepthAttachment struct {
	Texture           gpu.Texture
	View              gpu.TextureView
	DepthLoadOp       gpu.LoadOp
	DepthStoreOp      gpu.StoreOp
	DepthClearValue   float32
	StencilLoadOp     gpu.LoadOp
	StencilStoreOp    gpu.StoreOp
	StencilClearValue uint32
}

// NewRenderTarget creates a new render target.
func NewRenderTarget(desc RenderTargetDescriptor) (*RenderTarget, error) {
	if len(desc.ColorAttachments) == 0 {
		return nil, fmt.Errorf("render target must have at least one color attachment")
	}

	colorAttachments := make([]ColorAttachment, len(desc.ColorAttachments))
	copy(colorAttachments, desc.ColorAttachments)

	rt := &RenderTarget{
		colorAttachments: colorAttachments,
		size:             desc.Size,
		sampleCount:      desc.SampleCount,
	}

	if desc.DepthAttachment != nil {
		rt.depthAttachment = desc.DepthAttachment
	}

	return rt, nil
}

// RenderTargetDescriptor describes how to create a render target.
type RenderTargetDescriptor struct {
	ColorAttachments []ColorAttachment
	DepthAttachment  *DepthAttachment
	Size             geom.Size[geom.F32]
	SampleCount      uint32
}

// GetColorAttachments returns the color attachments.
func (rt *RenderTarget) GetColorAttachments() []ColorAttachment {
	return rt.colorAttachments
}

// GetDepthAttachment returns the depth attachment.
func (rt *RenderTarget) GetDepthAttachment() *DepthAttachment {
	return rt.depthAttachment
}

// GetSize returns the render target size.
func (rt *RenderTarget) GetSize() geom.Size[geom.F32] {
	return rt.size
}

// GetSampleCount returns the sample count for multisampling.
func (rt *RenderTarget) GetSampleCount() uint32 {
	return rt.sampleCount
}

// RenderPass encapsulates a GPU render pass with command recording.
type RenderPass struct {
	context        Context
	renderTarget   *RenderTarget
	commandEncoder gpu.CommandEncoder
	renderPass     gpu.RenderPassEncoder
	isActive       bool
}

// NewRenderPass creates a new render pass.
func NewRenderPass(context Context, renderTarget *RenderTarget, label string) (*RenderPass, error) {
	commandEncoder := context.GetDevice().CreateCommandEncoder(gpu.CommandEncoderDescriptor{
		Label: label,
	})

	// Build render pass descriptor
	colorAttachments := make([]gpu.RenderPassColorAttachment, len(renderTarget.colorAttachments))
	for i, attachment := range renderTarget.colorAttachments {
		colorAttachments[i] = gpu.RenderPassColorAttachment{
			View:       attachment.View,
			LoadOp:     attachment.LoadOp,
			StoreOp:    attachment.StoreOp,
			ClearValue: attachment.ClearValue,
		}
		if attachment.ResolveView != nil {
			colorAttachments[i].ResolveTarget = attachment.ResolveView
		}
	}

	passDesc := gpu.RenderPassDescriptor{
		Label:            label,
		ColorAttachments: colorAttachments,
	}

	// Add depth attachment if present
	if renderTarget.depthAttachment != nil {
		passDesc.DepthStencilAttachment = gpu.RenderPassDepthStencilAttachment{
			View:              renderTarget.depthAttachment.View,
			DepthLoadOp:       renderTarget.depthAttachment.DepthLoadOp,
			DepthStoreOp:      renderTarget.depthAttachment.DepthStoreOp,
			DepthClearValue:   renderTarget.depthAttachment.DepthClearValue,
			StencilLoadOp:     renderTarget.depthAttachment.StencilLoadOp,
			StencilStoreOp:    renderTarget.depthAttachment.StencilStoreOp,
			StencilClearValue: renderTarget.depthAttachment.StencilClearValue,
		}
	}

	renderPassEncoder := commandEncoder.BeginRenderPass(passDesc)

	return &RenderPass{
		context:        context,
		renderTarget:   renderTarget,
		commandEncoder: commandEncoder,
		renderPass:     renderPassEncoder,
		isActive:       true,
	}, nil
}

// SetPipeline sets the render pipeline for this pass.
func (rp *RenderPass) SetPipeline(pipeline gpu.RenderPipeline) {
	if !rp.isActive {
		return
	}
	rp.renderPass.SetPipeline(pipeline)
}

// SetBindGroup sets a bind group for the render pass.
func (rp *RenderPass) SetBindGroup(groupIndex uint32, bindGroup gpu.BindGroup, dynamicOffsets []uint32) {
	if !rp.isActive {
		return
	}
	rp.renderPass.SetBindGroup(groupIndex, bindGroup, dynamicOffsets)
}

// SetVertexBuffer sets a vertex buffer for the render pass.
func (rp *RenderPass) SetVertexBuffer(slot uint32, buffer gpu.Buffer, offset uint64, size uint64) {
	if !rp.isActive {
		return
	}
	rp.renderPass.SetVertexBuffer(slot, buffer, offset, size)
}

// SetIndexBuffer sets an index buffer for the render pass.
func (rp *RenderPass) SetIndexBuffer(buffer gpu.Buffer, format gpu.IndexFormat, offset uint64, size uint64) {
	if !rp.isActive {
		return
	}
	rp.renderPass.SetIndexBuffer(buffer, format, offset, size)
}

// Draw issues a draw command.
func (rp *RenderPass) Draw(vertexCount uint32, instanceCount uint32, firstVertex uint32, firstInstance uint32) {
	if !rp.isActive {
		return
	}
	rp.renderPass.Draw(vertexCount, instanceCount, firstVertex, firstInstance)
}

// DrawIndexed issues an indexed draw command.
func (rp *RenderPass) DrawIndexed(indexCount uint32, instanceCount uint32, firstIndex uint32, baseVertex int32, firstInstance uint32) {
	if !rp.isActive {
		return
	}
	rp.renderPass.DrawIndexed(indexCount, instanceCount, firstIndex, baseVertex, firstInstance)
}

// SetViewport sets the viewport for the render pass.
func (rp *RenderPass) SetViewport(x float32, y float32, width float32, height float32, minDepth float32, maxDepth float32) {
	if !rp.isActive {
		return
	}
	rp.renderPass.SetViewport(x, y, width, height, minDepth, maxDepth)
}

// SetScissorRect sets the scissor rectangle for the render pass.
func (rp *RenderPass) SetScissorRect(x uint32, y uint32, width uint32, height uint32) {
	if !rp.isActive {
		return
	}
	rp.renderPass.SetScissorRect(x, y, width, height)
}

// InsertDebugMarker inserts a debug marker.
func (rp *RenderPass) InsertDebugMarker(markerLabel string) {
	if !rp.isActive {
		return
	}
	rp.renderPass.InsertDebugMarker(markerLabel)
}

// PushDebugGroup pushes a debug group.
func (rp *RenderPass) PushDebugGroup(groupLabel string) {
	if !rp.isActive {
		return
	}
	rp.renderPass.PushDebugGroup(groupLabel)
}

// PopDebugGroup pops a debug group.
func (rp *RenderPass) PopDebugGroup() {
	if !rp.isActive {
		return
	}
	rp.renderPass.PopDebugGroup()
}

// End ends the render pass and returns the command buffer.
func (rp *RenderPass) End() gpu.CommandBuffer {
	if !rp.isActive {
		return nil
	}

	rp.renderPass.End()
	rp.isActive = false

	return rp.commandEncoder.Finish(gpu.CommandBufferDescriptor{
		Label: "Render Pass Commands",
	})
}

// EntityRenderer handles rendering of display list entities.
type EntityRenderer struct {
	context Context
}

// NewEntityRenderer creates a new entity renderer.
func NewEntityRenderer(context Context) *EntityRenderer {
	return &EntityRenderer{
		context: context,
	}
}

// RenderDisplayList renders a display list to the given render target.
func (er *EntityRenderer) RenderDisplayList(displayList *display.DisplayList, renderTarget *RenderTarget) error {
	if displayList == nil {
		return fmt.Errorf("display list cannot be nil")
	}
	if renderTarget == nil {
		return fmt.Errorf("render target cannot be nil")
	}

	// Use content renderer to handle display list rendering
	contentRenderer := NewContentRenderer(er.context)
	return contentRenderer.RenderDisplayList(displayList, renderTarget)
}
