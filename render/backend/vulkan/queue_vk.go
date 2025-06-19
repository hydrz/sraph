package vulkan

import (
	"github.com/opensraph/sraph/render"
	"github.com/vulkan-go/vulkan"
)

// QueueVK represents a Vulkan command queue
type QueueVK struct {
	device           vulkan.Device
	queue            vulkan.Queue
	queueFamilyIndex uint32
	trackedObjects   *TrackedObjectsVK
}

// NewQueueVK creates a new Vulkan queue wrapper
func NewQueueVK(device vulkan.Device, queue vulkan.Queue, queueFamilyIndex uint32) *QueueVK {
	return &QueueVK{
		device:           device,
		queue:            queue,
		queueFamilyIndex: queueFamilyIndex,
		trackedObjects:   NewTrackedObjectsVK(),
	}
}

// Submit submits a command buffer to the queue
func (q *QueueVK) Submit(commands []render.CommandBuffer, fence render.Fence) render.CommandSubmissionResult {
	// TODO: Convert command buffers to Vulkan and submit
	return render.CommandSubmissionResult{}
}

// GetCapabilities returns the queue capabilities
func (q *QueueVK) GetCapabilities() render.QueueCapabilities {
	// TODO: Return Vulkan queue capabilities
	return render.QueueCapabilities{}
}

// CreateCommandBuffer creates a new command buffer for this queue
func (q *QueueVK) CreateCommandBuffer() render.CommandBuffer {
	// TODO: Create Vulkan command buffer
	return nil
}

// WaitIdle waits for the queue to become idle
func (q *QueueVK) WaitIdle() error {
	ret := vulkan.QueueWaitIdle(q.queue)
	if ret != vulkan.Success {
		return render.NewError("Failed to wait for queue idle")
	}
	return nil
}

// GetQueue returns the underlying Vulkan queue
func (q *QueueVK) GetQueue() vulkan.Queue {
	return q.queue
}

// GetQueueFamilyIndex returns the queue family index
func (q *QueueVK) GetQueueFamilyIndex() uint32 {
	return q.queueFamilyIndex
}

// GetTrackedObjects returns the tracked objects manager
func (q *QueueVK) GetTrackedObjects() *TrackedObjectsVK {
	return q.trackedObjects
}

// InsertDebugMarker inserts a debug marker into the queue
func (q *QueueVK) InsertDebugMarker(label string) {
	// TODO: Insert debug marker if debug utils extension is available
}

// BeginDebugRegion begins a debug region
func (q *QueueVK) BeginDebugRegion(label string) {
	// TODO: Begin debug region if debug utils extension is available
}

// EndDebugRegion ends a debug region
func (q *QueueVK) EndDebugRegion() {
	// TODO: End debug region if debug utils extension is available
}
