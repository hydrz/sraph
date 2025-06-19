// Command queue interface and implementation for the render package.
// This provides a way to manage and execute command buffers.

package render

import (
	"errors"
	"sync"
	"time"

	"github.com/opensraph/sraph/gpu"
)

// CommandQueue represents a queue for executing command buffers.
type CommandQueue interface {
	// IsValid returns true if the command queue is valid and can be used.
	IsValid() bool

	// SetLabel sets a debug label for the command queue.
	SetLabel(label string)

	// Submit submits a command buffer for execution.
	Submit(commandBuffer CommandBuffer) error

	// SubmitMultiple submits multiple command buffers for execution.
	SubmitMultiple(commandBuffers []CommandBuffer) error

	// CreateCommandBuffer creates a new command buffer.
	CreateCommandBuffer() (CommandBuffer, error)

	// WaitIdle waits until all submitted commands have completed.
	WaitIdle() error

	// GetCapabilities returns the queue capabilities.
	GetCapabilities() *QueueCapabilities
}

// QueueType represents the type of command queue.
type QueueType int

const (
	QueueTypeUnknown QueueType = iota
	QueueTypeGraphics
	QueueTypeCompute
	QueueTypeTransfer
)

// QueueCapabilities describes the capabilities of a command queue.
type QueueCapabilities struct {
	Type               QueueType
	SupportsGraphics   bool
	SupportsCompute    bool
	SupportsTransfer   bool
	MaxCommandBuffers  uint32
	TimestampValidBits uint32
	MinTimestampPeriod float64
}

// CommandQueueDescriptor describes properties for creating a command queue.
type CommandQueueDescriptor struct {
	Type  QueueType
	Label string
}

// DefaultCommandQueue is the default implementation of CommandQueue.
type DefaultCommandQueue struct {
	device       gpu.Device
	queue        gpu.Queue
	allocator    ResourceAllocator
	capabilities *QueueCapabilities
	label        string
	isValid      bool
	mutex        sync.RWMutex
}

// NewCommandQueue creates a new command queue with the given descriptor.
func NewCommandQueue(device gpu.Device, allocator ResourceAllocator, descriptor *CommandQueueDescriptor) (CommandQueue, error) {
	if device == nil {
		return nil, errors.New("device cannot be nil")
	}
	if allocator == nil {
		return nil, errors.New("allocator cannot be nil")
	}
	if descriptor == nil {
		descriptor = &CommandQueueDescriptor{
			Type: QueueTypeGraphics,
		}
	}

	queue := &DefaultCommandQueue{
		device:    device,
		allocator: allocator,
		label:     descriptor.Label,
		isValid:   true,
		capabilities: &QueueCapabilities{
			Type:               descriptor.Type,
			SupportsGraphics:   descriptor.Type == QueueTypeGraphics,
			SupportsCompute:    descriptor.Type == QueueTypeCompute || descriptor.Type == QueueTypeGraphics,
			SupportsTransfer:   true, // All queues support transfer
			MaxCommandBuffers:  1000,
			TimestampValidBits: 64,
			MinTimestampPeriod: 1.0,
		},
	}

	// Create the underlying GPU queue
	gpuQueue, err := device.GetQueue()
	if err != nil {
		return nil, err
	}
	queue.queue = gpuQueue

	if descriptor.Label != "" {
		queue.SetLabel(descriptor.Label)
	}

	return queue, nil
}

// IsValid returns true if the command queue is valid and can be used.
func (cq *DefaultCommandQueue) IsValid() bool {
	cq.mutex.RLock()
	defer cq.mutex.RUnlock()
	return cq.isValid && cq.queue != nil
}

// SetLabel sets a debug label for the command queue.
func (cq *DefaultCommandQueue) SetLabel(label string) {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()

	cq.label = label
	if cq.queue != nil {
		cq.queue.SetLabel(label)
	}
}

// Submit submits a command buffer for execution.
func (cq *DefaultCommandQueue) Submit(commandBuffer CommandBuffer) error {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()

	if !cq.isValid {
		return errors.New("command queue is not valid")
	}

	if commandBuffer == nil {
		return errors.New("command buffer cannot be nil")
	}

	if !commandBuffer.IsValid() {
		return errors.New("command buffer is not valid")
	}

	// Submit the command buffer for execution
	return commandBuffer.SubmitCommands()
}

// SubmitMultiple submits multiple command buffers for execution.
func (cq *DefaultCommandQueue) SubmitMultiple(commandBuffers []CommandBuffer) error {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()

	if !cq.isValid {
		return errors.New("command queue is not valid")
	}

	if len(commandBuffers) == 0 {
		return errors.New("command buffers list cannot be empty")
	}

	// Submit all command buffers
	for i, cmdBuffer := range commandBuffers {
		if cmdBuffer == nil {
			return errors.New("command buffer at index %d is nil")
		}
		if !cmdBuffer.IsValid() {
			return errors.New("command buffer at index %d is not valid")
		}

		if err := cmdBuffer.SubmitCommands(); err != nil {
			return err
		}
	}

	return nil
}

// CreateCommandBuffer creates a new command buffer.
func (cq *DefaultCommandQueue) CreateCommandBuffer() (CommandBuffer, error) {
	cq.mutex.RLock()
	defer cq.mutex.RUnlock()

	if !cq.isValid {
		return nil, errors.New("command queue is not valid")
	}

	descriptor := &CommandBufferDescriptor{
		Type: CommandBufferTypePrimary,
	}

	return NewCommandBuffer(cq.device, cq.allocator, descriptor)
}

// WaitIdle waits until all submitted commands have completed.
func (cq *DefaultCommandQueue) WaitIdle() error {
	cq.mutex.RLock()
	defer cq.mutex.RUnlock()

	if !cq.isValid {
		return errors.New("command queue is not valid")
	}

	if cq.queue == nil {
		return errors.New("underlying queue is nil")
	}

	// Wait for the queue to become idle
	return cq.queue.WaitIdle()
}

// GetCapabilities returns the queue capabilities.
func (cq *DefaultCommandQueue) GetCapabilities() *QueueCapabilities {
	cq.mutex.RLock()
	defer cq.mutex.RUnlock()
	return cq.capabilities
}

// CommandExecutor provides high-level command execution utilities.
type CommandExecutor struct {
	queue     CommandQueue
	allocator ResourceAllocator
}

// NewCommandExecutor creates a new command executor.
func NewCommandExecutor(queue CommandQueue, allocator ResourceAllocator) *CommandExecutor {
	return &CommandExecutor{
		queue:     queue,
		allocator: allocator,
	}
}

// ExecuteCommands executes a function that records commands into a command buffer.
func (ce *CommandExecutor) ExecuteCommands(recordFunc func(CommandBuffer) error) error {
	if ce.queue == nil {
		return errors.New("command queue is nil")
	}

	// Create a command buffer
	cmdBuffer, err := ce.queue.CreateCommandBuffer()
	if err != nil {
		return err
	}

	// Record commands
	if err := recordFunc(cmdBuffer); err != nil {
		return err
	}

	// Submit the command buffer
	if err := ce.queue.Submit(cmdBuffer); err != nil {
		return err
	}

	// Wait for completion
	return cmdBuffer.WaitUntilCompleted()
}

// ExecuteCommandsAsync executes commands asynchronously.
func (ce *CommandExecutor) ExecuteCommandsAsync(recordFunc func(CommandBuffer) error) error {
	if ce.queue == nil {
		return errors.New("command queue is nil")
	}

	// Create a command buffer
	cmdBuffer, err := ce.queue.CreateCommandBuffer()
	if err != nil {
		return err
	}

	// Record commands
	if err := recordFunc(cmdBuffer); err != nil {
		return err
	}

	// Submit the command buffer (don't wait)
	return ce.queue.Submit(cmdBuffer)
}

// SubmissionInfo holds information about a command submission.
type SubmissionInfo struct {
	CommandBuffers []CommandBuffer
	SubmissionTime time.Time
	CompletionTime *time.Time
	Status         SubmissionStatus
}

// SubmissionStatus represents the status of a command submission.
type SubmissionStatus int

const (
	SubmissionStatusPending SubmissionStatus = iota
	SubmissionStatusExecuting
	SubmissionStatusCompleted
	SubmissionStatusFailed
)

// String returns the string representation of the submission status.
func (s SubmissionStatus) String() string {
	switch s {
	case SubmissionStatusPending:
		return "Pending"
	case SubmissionStatusExecuting:
		return "Executing"
	case SubmissionStatusCompleted:
		return "Completed"
	case SubmissionStatusFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}
