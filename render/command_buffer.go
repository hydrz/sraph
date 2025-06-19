// Command buffer interface and implementation for the render package.
// This provides a way to record and submit rendering commands.

package render

import (
	"errors"
	"fmt"
	"sync"

	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/gpu"
)

// CommandBuffer represents a buffer for recording rendering commands.
type CommandBuffer interface {
	// IsValid returns true if the command buffer is valid and can be used.
	IsValid() bool

	// SetLabel sets a debug label for the command buffer.
	SetLabel(label string)

	// CreateRenderPass creates a render pass with the given descriptor.
	CreateRenderPass(descriptor *RenderPassDescriptor) (RenderPass, error)

	// CreateBlitPass creates a blit pass.
	CreateBlitPass() (BlitPass, error)

	// CreateComputePass creates a compute pass.
	CreateComputePass() (ComputePass, error)

	// SubmitCommands submits all recorded commands for execution.
	SubmitCommands() error

	// WaitUntilCompleted waits until all commands have been completed.
	WaitUntilCompleted() error

	// GetResourceAllocator returns the resource allocator for this command buffer.
	GetResourceAllocator() ResourceAllocator
}

// CommandBufferType represents the type of command buffer.
type CommandBufferType int

const (
	CommandBufferTypeUnknown CommandBufferType = iota
	CommandBufferTypePrimary
	CommandBufferTypeSecondary
)

// CommandBufferDescriptor describes properties for creating a command buffer.
type CommandBufferDescriptor struct {
	Type  CommandBufferType
	Label string
}

// DefaultCommandBuffer is the default implementation of CommandBuffer.
type DefaultCommandBuffer struct {
	device        gpu.Device
	commandBuffer *gpu.CommandBuffer
	allocator     ResourceAllocator
	label         string
	isValid       bool
	mutex         sync.RWMutex
}

// NewCommandBuffer creates a new command buffer with the given descriptor.
func NewCommandBuffer(device gpu.Device, allocator ResourceAllocator, descriptor *CommandBufferDescriptor) (CommandBuffer, error) {
	if device == nil {
		return nil, errors.New("device cannot be nil")
	}
	if allocator == nil {
		return nil, errors.New("allocator cannot be nil")
	}
	if descriptor == nil {
		descriptor = &CommandBufferDescriptor{
			Type: CommandBufferTypePrimary,
		}
	}

	cmdBuffer := &DefaultCommandBuffer{
		device:    device,
		allocator: allocator,
		label:     descriptor.Label,
		isValid:   true,
	}

	// Create the underlying GPU command buffer
	gpuCmdBuffer, err := device.CreateCommandBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to create GPU command buffer: %w", err)
	}

	cmdBuffer.commandBuffer = gpuCmdBuffer

	if descriptor.Label != "" {
		cmdBuffer.SetLabel(descriptor.Label)
	}

	return cmdBuffer, nil
}

// IsValid returns true if the command buffer is valid and can be used.
func (cb *DefaultCommandBuffer) IsValid() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.isValid && cb.commandBuffer != nil
}

// SetLabel sets a debug label for the command buffer.
func (cb *DefaultCommandBuffer) SetLabel(label string) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.label = label
	if cb.commandBuffer != nil {
		cb.commandBuffer.SetLabel(label)
	}
}

// CreateRenderPass creates a render pass with the given descriptor.
func (cb *DefaultCommandBuffer) CreateRenderPass(descriptor *RenderPassDescriptor) (RenderPass, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if !cb.isValid {
		return nil, errors.New("command buffer is not valid")
	}

	if descriptor == nil {
		return nil, errors.New("render pass descriptor cannot be nil")
	}

	return NewRenderPass(cb.commandBuffer, descriptor)
}

// CreateBlitPass creates a blit pass.
func (cb *DefaultCommandBuffer) CreateBlitPass() (BlitPass, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if !cb.isValid {
		return nil, errors.New("command buffer is not valid")
	}

	return NewBlitPass(cb.commandBuffer)
}

// CreateComputePass creates a compute pass.
func (cb *DefaultCommandBuffer) CreateComputePass() (ComputePass, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if !cb.isValid {
		return nil, errors.New("command buffer is not valid")
	}

	return NewComputePass(cb.commandBuffer)
}

// SubmitCommands submits all recorded commands for execution.
func (cb *DefaultCommandBuffer) SubmitCommands() error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if !cb.isValid {
		return errors.New("command buffer is not valid")
	}

	if cb.commandBuffer == nil {
		return errors.New("underlying command buffer is nil")
	}

	return cb.commandBuffer.Submit()
}

// WaitUntilCompleted waits until all commands have been completed.
func (cb *DefaultCommandBuffer) WaitUntilCompleted() error {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	if !cb.isValid {
		return errors.New("command buffer is not valid")
	}

	if cb.commandBuffer == nil {
		return errors.New("underlying command buffer is nil")
	}

	return cb.commandBuffer.Wait()
}

// GetResourceAllocator returns the resource allocator for this command buffer.
func (cb *DefaultCommandBuffer) GetResourceAllocator() ResourceAllocator {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.allocator
}

// ComputePass interface for compute operations.
type ComputePass interface {
	// SetLabel sets a debug label for the compute pass.
	SetLabel(label string)

	// SetComputePipeline sets the compute pipeline for this pass.
	SetComputePipeline(pipeline ComputePipeline)

	// SetCommandBuffer sets the command buffer for this pass.
	SetCommandBuffer(buffer Buffer, offset uint64, binding uint32)

	// DispatchThreadgroups dispatches compute threadgroups.
	DispatchThreadgroups(threadgroupsPerGrid, threadsPerThreadgroup geom.Vector3) error

	// EncodeCommands encodes the compute commands.
	EncodeCommands() error
}

// ComputePipeline represents a compute pipeline.
type ComputePipeline interface {
	// IsValid returns true if the pipeline is valid.
	IsValid() bool

	// GetDescriptor returns the pipeline descriptor.
	GetDescriptor() *ComputePipelineDescriptor
}

// ComputePipelineDescriptor describes a compute pipeline.
type ComputePipelineDescriptor struct {
	Label           string
	ComputeShader   ShaderFunction
	ThreadGroupSize geom.Vector3
}

// DefaultComputePass is the default implementation of ComputePass.
type DefaultComputePass struct {
	commandBuffer *gpu.CommandBuffer
	pipeline      ComputePipeline
	label         string
	isValid       bool
	mutex         sync.RWMutex
}

// NewComputePass creates a new compute pass.
func NewComputePass(commandBuffer *gpu.CommandBuffer) (ComputePass, error) {
	if commandBuffer == nil {
		return nil, errors.New("command buffer cannot be nil")
	}

	return &DefaultComputePass{
		commandBuffer: commandBuffer,
		isValid:       true,
	}, nil
}

// SetLabel sets a debug label for the compute pass.
func (cp *DefaultComputePass) SetLabel(label string) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	cp.label = label
}

// SetComputePipeline sets the compute pipeline for this pass.
func (cp *DefaultComputePass) SetComputePipeline(pipeline ComputePipeline) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	cp.pipeline = pipeline
}

// SetCommandBuffer sets the command buffer for this pass.
func (cp *DefaultComputePass) SetCommandBuffer(buffer Buffer, offset uint64, binding uint32) {
	// Implementation would bind the buffer to the compute shader
}

// DispatchThreadgroups dispatches compute threadgroups.
func (cp *DefaultComputePass) DispatchThreadgroups(threadgroupsPerGrid, threadsPerThreadgroup geom.Vector3) error {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	if !cp.isValid {
		return errors.New("compute pass is not valid")
	}

	// Implementation would dispatch the compute work
	return nil
}

// EncodeCommands encodes the compute commands.
func (cp *DefaultComputePass) EncodeCommands() error {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	if !cp.isValid {
		return errors.New("compute pass is not valid")
	}

	// Implementation would encode the compute commands
	return nil
}
