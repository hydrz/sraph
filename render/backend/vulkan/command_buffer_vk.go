package vulkan

import (
	"fmt"
	"unsafe"

	"github.com/vulkan-go/vulkan"
)

// CommandBufferVK represents a Vulkan command buffer wrapper
type CommandBufferVK struct {
	context       *ContextVK
	commandBuffer vulkan.CommandBuffer
	commandPool   vulkan.CommandPool
	level         vulkan.CommandBufferLevel
	isRecording   bool
	isOneTime     bool
}

// NewCommandBufferVK creates a new command buffer
func NewCommandBufferVK(context *ContextVK, commandPool vulkan.CommandPool, level vulkan.CommandBufferLevel) (*CommandBufferVK, error) {
	if context == nil || commandPool == vulkan.NullCommandPool {
		return nil, fmt.Errorf("invalid parameters")
	}

	allocInfo := vulkan.CommandBufferAllocateInfo{
		SType:              vulkan.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        commandPool,
		Level:              level,
		CommandBufferCount: 1,
	}

	var commandBuffer vulkan.CommandBuffer
	if result := vulkan.AllocateCommandBuffers(context.device, &allocInfo, &commandBuffer); result != vulkan.Success {
		return nil, fmt.Errorf("failed to allocate command buffer: %s", result)
	}

	return &CommandBufferVK{
		context:       context,
		commandBuffer: commandBuffer,
		commandPool:   commandPool,
		level:         level,
		isRecording:   false,
		isOneTime:     false,
	}, nil
}

// Begin begins recording commands
func (cb *CommandBufferVK) Begin(flags vulkan.CommandBufferUsageFlags) error {
	if cb.isRecording {
		return fmt.Errorf("command buffer is already recording")
	}

	beginInfo := vulkan.CommandBufferBeginInfo{
		SType: vulkan.StructureTypeCommandBufferBeginInfo,
		Flags: flags,
	}

	if result := vulkan.BeginCommandBuffer(cb.commandBuffer, &beginInfo); result != vulkan.Success {
		return fmt.Errorf("failed to begin command buffer: %s", result)
	}

	cb.isRecording = true
	cb.isOneTime = (flags & vulkan.CommandBufferUsageOneTimeSubmitBit) != 0

	return nil
}

// End ends recording commands
func (cb *CommandBufferVK) End() error {
	if !cb.isRecording {
		return fmt.Errorf("command buffer is not recording")
	}

	if result := vulkan.EndCommandBuffer(cb.commandBuffer); result != vulkan.Success {
		return fmt.Errorf("failed to end command buffer: %s", result)
	}

	cb.isRecording = false
	return nil
}

// Reset resets the command buffer
func (cb *CommandBufferVK) Reset(flags vulkan.CommandBufferResetFlags) error {
	if cb.isRecording {
		return fmt.Errorf("cannot reset command buffer while recording")
	}

	if result := vulkan.ResetCommandBuffer(cb.commandBuffer, flags); result != vulkan.Success {
		return fmt.Errorf("failed to reset command buffer: %s", result)
	}

	return nil
}

// BindPipeline binds a graphics or compute pipeline
func (cb *CommandBufferVK) BindPipeline(bindPoint vulkan.PipelineBindPoint, pipeline vulkan.Pipeline) {
	if cb.isRecording {
		vulkan.CmdBindPipeline(cb.commandBuffer, bindPoint, pipeline)
	}
}

// BindVertexBuffers binds vertex buffers
func (cb *CommandBufferVK) BindVertexBuffers(firstBinding uint32, buffers []vulkan.Buffer, offsets []vulkan.DeviceSize) {
	if cb.isRecording && len(buffers) == len(offsets) {
		vulkan.CmdBindVertexBuffers(cb.commandBuffer, firstBinding, uint32(len(buffers)), buffers, offsets)
	}
}

// BindIndexBuffer binds an index buffer
func (cb *CommandBufferVK) BindIndexBuffer(buffer vulkan.Buffer, offset vulkan.DeviceSize, indexType vulkan.IndexType) {
	if cb.isRecording {
		vulkan.CmdBindIndexBuffer(cb.commandBuffer, buffer, offset, indexType)
	}
}

// BindDescriptorSets binds descriptor sets
func (cb *CommandBufferVK) BindDescriptorSets(bindPoint vulkan.PipelineBindPoint, layout vulkan.PipelineLayout, firstSet uint32, descriptorSets []vulkan.DescriptorSet, dynamicOffsets []uint32) {
	if cb.isRecording {
		vulkan.CmdBindDescriptorSets(cb.commandBuffer, bindPoint, layout, firstSet, uint32(len(descriptorSets)), descriptorSets, uint32(len(dynamicOffsets)), dynamicOffsets)
	}
}

// Draw draws primitives
func (cb *CommandBufferVK) Draw(vertexCount, instanceCount, firstVertex, firstInstance uint32) {
	if cb.isRecording {
		vulkan.CmdDraw(cb.commandBuffer, vertexCount, instanceCount, firstVertex, firstInstance)
	}
}

// DrawIndexed draws indexed primitives
func (cb *CommandBufferVK) DrawIndexed(indexCount, instanceCount, firstIndex uint32, vertexOffset int32, firstInstance uint32) {
	if cb.isRecording {
		vulkan.CmdDrawIndexed(cb.commandBuffer, indexCount, instanceCount, firstIndex, vertexOffset, firstInstance)
	}
}

// Dispatch dispatches compute work groups
func (cb *CommandBufferVK) Dispatch(groupCountX, groupCountY, groupCountZ uint32) {
	if cb.isRecording {
		vulkan.CmdDispatch(cb.commandBuffer, groupCountX, groupCountY, groupCountZ)
	}
}

// CopyBuffer copies data between buffers
func (cb *CommandBufferVK) CopyBuffer(srcBuffer, dstBuffer vulkan.Buffer, regions []vulkan.BufferCopy) {
	if cb.isRecording {
		vulkan.CmdCopyBuffer(cb.commandBuffer, srcBuffer, dstBuffer, uint32(len(regions)), regions)
	}
}

// CopyBufferToImage copies data from buffer to image
func (cb *CommandBufferVK) CopyBufferToImage(srcBuffer vulkan.Buffer, dstImage vulkan.Image, dstImageLayout vulkan.ImageLayout, regions []vulkan.BufferImageCopy) {
	if cb.isRecording {
		vulkan.CmdCopyBufferToImage(cb.commandBuffer, srcBuffer, dstImage, dstImageLayout, uint32(len(regions)), regions)
	}
}

// PipelineBarrier inserts a pipeline barrier
func (cb *CommandBufferVK) PipelineBarrier(srcStageMask, dstStageMask vulkan.PipelineStageFlags, dependencyFlags vulkan.DependencyFlags, memoryBarriers []vulkan.MemoryBarrier, bufferMemoryBarriers []vulkan.BufferMemoryBarrier, imageMemoryBarriers []vulkan.ImageMemoryBarrier) {
	if cb.isRecording {
		vulkan.CmdPipelineBarrier(cb.commandBuffer, srcStageMask, dstStageMask, dependencyFlags,
			uint32(len(memoryBarriers)), memoryBarriers,
			uint32(len(bufferMemoryBarriers)), bufferMemoryBarriers,
			uint32(len(imageMemoryBarriers)), imageMemoryBarriers)
	}
}

// SetViewport sets the viewport
func (cb *CommandBufferVK) SetViewport(firstViewport uint32, viewports []vulkan.Viewport) {
	if cb.isRecording {
		vulkan.CmdSetViewport(cb.commandBuffer, firstViewport, uint32(len(viewports)), viewports)
	}
}

// SetScissor sets the scissor rectangles
func (cb *CommandBufferVK) SetScissor(firstScissor uint32, scissors []vulkan.Rect2D) {
	if cb.isRecording {
		vulkan.CmdSetScissor(cb.commandBuffer, firstScissor, uint32(len(scissors)), scissors)
	}
}

// PushConstants pushes constants to the pipeline
func (cb *CommandBufferVK) PushConstants(layout vulkan.PipelineLayout, stageFlags vulkan.ShaderStageFlags, offset uint32, data []byte) {
	if cb.isRecording && len(data) > 0 {
		vulkan.CmdPushConstants(cb.commandBuffer, layout, stageFlags, offset, uint32(len(data)), unsafe.Pointer(&data[0]))
	}
}

// BeginRenderPass begins a render pass
func (cb *CommandBufferVK) BeginRenderPass(renderPassBegin *vulkan.RenderPassBeginInfo, contents vulkan.SubpassContents) {
	if cb.isRecording {
		vulkan.CmdBeginRenderPass(cb.commandBuffer, renderPassBegin, contents)
	}
}

// EndRenderPass ends the current render pass
func (cb *CommandBufferVK) EndRenderPass() {
	if cb.isRecording {
		vulkan.CmdEndRenderPass(cb.commandBuffer)
	}
}

// NextSubpass advances to the next subpass
func (cb *CommandBufferVK) NextSubpass(contents vulkan.SubpassContents) {
	if cb.isRecording {
		vulkan.CmdNextSubpass(cb.commandBuffer, contents)
	}
}

// BlitImage blits an image region
func (cb *CommandBufferVK) BlitImage(srcImage vulkan.Image, srcImageLayout vulkan.ImageLayout, dstImage vulkan.Image, dstImageLayout vulkan.ImageLayout, regions []vulkan.ImageBlit, filter vulkan.Filter) {
	if cb.isRecording {
		vulkan.CmdBlitImage(cb.commandBuffer, srcImage, srcImageLayout, dstImage, dstImageLayout, uint32(len(regions)), regions, filter)
	}
}

// ClearColorImage clears a color image
func (cb *CommandBufferVK) ClearColorImage(image vulkan.Image, imageLayout vulkan.ImageLayout, color *vulkan.ClearColorValue, ranges []vulkan.ImageSubresourceRange) {
	if cb.isRecording {
		vulkan.CmdClearColorImage(cb.commandBuffer, image, imageLayout, color, uint32(len(ranges)), ranges)
	}
}

// ClearDepthStencilImage clears a depth/stencil image
func (cb *CommandBufferVK) ClearDepthStencilImage(image vulkan.Image, imageLayout vulkan.ImageLayout, depthStencil *vulkan.ClearDepthStencilValue, ranges []vulkan.ImageSubresourceRange) {
	if cb.isRecording {
		vulkan.CmdClearDepthStencilImage(cb.commandBuffer, image, imageLayout, depthStencil, uint32(len(ranges)), ranges)
	}
}

// WriteTimestamp writes a timestamp
func (cb *CommandBufferVK) WriteTimestamp(pipelineStage vulkan.PipelineStageFlags, queryPool vulkan.QueryPool, query uint32) {
	if cb.isRecording {
		vulkan.CmdWriteTimestamp(cb.commandBuffer, pipelineStage, queryPool, query)
	}
}

// GetCommandBuffer returns the Vulkan command buffer handle
func (cb *CommandBufferVK) GetCommandBuffer() vulkan.CommandBuffer {
	return cb.commandBuffer
}

// IsRecording returns true if the command buffer is currently recording
func (cb *CommandBufferVK) IsRecording() bool {
	return cb.isRecording
}

// IsOneTime returns true if this is a one-time submit command buffer
func (cb *CommandBufferVK) IsOneTime() bool {
	return cb.isOneTime
}

// GetLevel returns the command buffer level
func (cb *CommandBufferVK) GetLevel() vulkan.CommandBufferLevel {
	return cb.level
}

// Destroy destroys the command buffer (returns it to the pool)
func (cb *CommandBufferVK) Destroy() {
	if cb.context != nil && cb.commandBuffer != vulkan.NullCommandBuffer {
		vulkan.FreeCommandBuffers(cb.context.device, cb.commandPool, 1, []vulkan.CommandBuffer{cb.commandBuffer})
		cb.commandBuffer = vulkan.NullCommandBuffer
	}
}
