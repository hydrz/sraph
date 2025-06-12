package gpu

import (
	"fmt"
	"sync"
)

var _ CommandEncoder = (*commandEncoder)(nil)

// commandEncoder implements the CommandEncoder interface
type commandEncoder struct {
	mu        sync.RWMutex
	refCount  int32
	label     string
	commands  []interface{} // Store commands for later execution
	finished  bool
	destroyed bool
}

// BeginComputePass begins a compute pass
func (ce *commandEncoder) BeginComputePass(descriptor ComputePassDescriptor) (ComputePassEncoder, error) {
	ce.mu.RLock()
	defer ce.mu.RUnlock()

	if ce.destroyed {
		return nil, fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return nil, fmt.Errorf("command encoder has been finished")
	}

	encoder := &computePassEncoder{
		refCount:        1,
		label:           descriptor.Label,
		timestampWrites: descriptor.TimestampWrites,
	}

	return encoder, nil
}

// BeginRenderPass begins a render pass
func (ce *commandEncoder) BeginRenderPass(descriptor RenderPassDescriptor) (RenderPassEncoder, error) {
	ce.mu.RLock()
	defer ce.mu.RUnlock()

	if ce.destroyed {
		return nil, fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return nil, fmt.Errorf("command encoder has been finished")
	}

	encoder := &renderPassEncoder{
		refCount:               1,
		label:                  descriptor.Label,
		colorAttachments:       descriptor.ColorAttachments,
		depthStencilAttachment: descriptor.DepthStencilAttachment,
		occlusionQuerySet:      descriptor.OcclusionQuerySet,
		timestampWrites:        descriptor.TimestampWrites,
	}

	return encoder, nil
}

// ClearBuffer clears a buffer
func (ce *commandEncoder) ClearBuffer(buffer Buffer, offset uint64, size uint64) error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return fmt.Errorf("command encoder has been finished")
	}

	// In a real implementation, this would record a clear buffer command
	_ = buffer
	_ = offset
	_ = size

	return nil
}

// CopyBufferToBuffer copies data between buffers
func (ce *commandEncoder) CopyBufferToBuffer(source Buffer, sourceOffset uint64, destination Buffer, destinationOffset uint64, size uint64) error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return fmt.Errorf("command encoder has been finished")
	}

	// In a real implementation, this would record a buffer copy command
	_ = source
	_ = sourceOffset
	_ = destination
	_ = destinationOffset
	_ = size

	return nil
}

// CopyBufferToTexture copies data from buffer to texture
func (ce *commandEncoder) CopyBufferToTexture(source TexelCopyBufferInfo, destination TexelCopyTextureInfo, copySize Extent3D) error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return fmt.Errorf("command encoder has been finished")
	}

	// In a real implementation, this would record a buffer to texture copy command
	_ = source
	_ = destination
	_ = copySize

	return nil
}

// CopyTextureToBuffer copies data from texture to buffer
func (ce *commandEncoder) CopyTextureToBuffer(source TexelCopyTextureInfo, destination TexelCopyBufferInfo, copySize Extent3D) error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return fmt.Errorf("command encoder has been finished")
	}

	// In a real implementation, this would record a texture to buffer copy command
	_ = source
	_ = destination
	_ = copySize

	return nil
}

// CopyTextureToTexture copies data between textures
func (ce *commandEncoder) CopyTextureToTexture(source TexelCopyTextureInfo, destination TexelCopyTextureInfo, copySize Extent3D) error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return fmt.Errorf("command encoder has been finished")
	}

	// In a real implementation, this would record a texture to texture copy command
	_ = source
	_ = destination
	_ = copySize

	return nil
}

// Finish finishes recording commands and returns a command buffer
func (ce *commandEncoder) Finish(descriptor CommandBufferDescriptor) (CommandBuffer, error) {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return nil, fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return nil, fmt.Errorf("command encoder has already been finished")
	}

	ce.finished = true

	commandBuffer := &commandBuffer{
		refCount: 1,
		label:    descriptor.Label,
		commands: ce.commands,
	}

	return commandBuffer, nil
}

// InsertDebugMarker inserts a debug marker
func (ce *commandEncoder) InsertDebugMarker(markerLabel string) error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return fmt.Errorf("command encoder has been finished")
	}

	// In a real implementation, this would insert a debug marker
	_ = markerLabel

	return nil
}

// PopDebugGroup pops a debug group
func (ce *commandEncoder) PopDebugGroup() error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return fmt.Errorf("command encoder has been finished")
	}

	// In a real implementation, this would pop a debug group
	return nil
}

// PushDebugGroup pushes a debug group
func (ce *commandEncoder) PushDebugGroup(groupLabel string) error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return fmt.Errorf("command encoder has been finished")
	}

	// In a real implementation, this would push a debug group
	_ = groupLabel

	return nil
}

// ResolveQuerySet resolves a query set
func (ce *commandEncoder) ResolveQuerySet(querySet QuerySet, firstQuery uint32, queryCount uint32, destination Buffer, destinationOffset uint64) error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return fmt.Errorf("command encoder has been finished")
	}

	// In a real implementation, this would resolve a query set
	_ = querySet
	_ = firstQuery
	_ = queryCount
	_ = destination
	_ = destinationOffset

	return nil
}

// SetLabel sets the command encoder label
func (ce *commandEncoder) SetLabel(label string) error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	ce.label = label
	return nil
}

// WriteTimestamp writes a timestamp
func (ce *commandEncoder) WriteTimestamp(querySet QuerySet, queryIndex uint32) error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	if ce.finished {
		return fmt.Errorf("command encoder has been finished")
	}

	// In a real implementation, this would write a timestamp
	_ = querySet
	_ = queryIndex

	return nil
}

// AddRef increments the reference count
func (ce *commandEncoder) AddRef() error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	ce.refCount++
	return nil
}

// Release decrements the reference count and destroys if zero
func (ce *commandEncoder) Release() error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.destroyed {
		return fmt.Errorf("command encoder has been destroyed")
	}

	ce.refCount--
	if ce.refCount <= 0 {
		ce.destroyed = true
	}

	return nil
}

// commandBuffer implements the CommandBuffer interface
type commandBuffer struct {
	mu        sync.RWMutex
	refCount  int32
	label     string
	commands  []interface{}
	destroyed bool
}

// SetLabel sets the command buffer label
func (cb *commandBuffer) SetLabel(label string) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.destroyed {
		return fmt.Errorf("command buffer has been destroyed")
	}

	cb.label = label
	return nil
}

// AddRef increments the reference count
func (cb *commandBuffer) AddRef() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.destroyed {
		return fmt.Errorf("command buffer has been destroyed")
	}

	cb.refCount++
	return nil
}

// Release decrements the reference count and destroys if zero
func (cb *commandBuffer) Release() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.destroyed {
		return fmt.Errorf("command buffer has been destroyed")
	}

	cb.refCount--
	if cb.refCount <= 0 {
		cb.destroyed = true
	}

	return nil
}
