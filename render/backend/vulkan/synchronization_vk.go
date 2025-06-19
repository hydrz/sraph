package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// FenceVK wraps Vulkan fence
// It provides CPU-GPU synchronization primitive
type FenceVK struct {
	device   vulkan.Device
	fence    vulkan.Fence
	signaled bool
}

// NewFenceVK creates a new Vulkan fence
func NewFenceVK(device vulkan.Device, signaled bool) (*FenceVK, error) {
	var flags vulkan.FenceCreateFlags
	if signaled {
		flags = vulkan.FenceCreateFlags(vulkan.FenceCreateSignaledBit)
	}

	createInfo := vulkan.FenceCreateInfo{
		SType: vulkan.StructureTypeFenceCreateInfo,
		Flags: flags,
	}

	var fence vulkan.Fence
	if result := vulkan.CreateFence(device, &createInfo, nil, &fence); result != vulkan.Success {
		return nil, fmt.Errorf("failed to create fence: %s", result)
	}

	return &FenceVK{
		device:   device,
		fence:    fence,
		signaled: signaled,
	}, nil
}

// GetFence returns the Vulkan fence handle
func (f *FenceVK) GetFence() vulkan.Fence {
	return f.fence
}

// Wait waits for the fence to be signaled
func (f *FenceVK) Wait(timeout uint64) error {
	if result := vulkan.WaitForFences(f.device, 1, []vulkan.Fence{f.fence}, vulkan.True, timeout); result != vulkan.Success {
		return fmt.Errorf("failed to wait for fence: %s", result)
	}
	f.signaled = true
	return nil
}

// Reset resets the fence to unsignaled state
func (f *FenceVK) Reset() error {
	if result := vulkan.ResetFences(f.device, 1, []vulkan.Fence{f.fence}); result != vulkan.Success {
		return fmt.Errorf("failed to reset fence: %s", result)
	}
	f.signaled = false
	return nil
}

// IsSignaled checks if the fence is signaled
func (f *FenceVK) IsSignaled() bool {
	status := vulkan.GetFenceStatus(f.device, f.fence)
	signaled := status == vulkan.Success
	f.signaled = signaled
	return signaled
}

// Destroy destroys the fence
func (f *FenceVK) Destroy() {
	if f.fence != vulkan.NullFence {
		vulkan.DestroyFence(f.device, f.fence, nil)
		f.fence = vulkan.NullFence
	}
}

// SemaphoreVK wraps Vulkan semaphore
// It provides GPU-GPU synchronization primitive
type SemaphoreVK struct {
	device    vulkan.Device
	semaphore vulkan.Semaphore
}

// NewSemaphoreVK creates a new Vulkan semaphore
func NewSemaphoreVK(device vulkan.Device) (*SemaphoreVK, error) {
	createInfo := vulkan.SemaphoreCreateInfo{
		SType: vulkan.StructureTypeSemaphoreCreateInfo,
	}

	var semaphore vulkan.Semaphore
	if result := vulkan.CreateSemaphore(device, &createInfo, nil, &semaphore); result != vulkan.Success {
		return nil, fmt.Errorf("failed to create semaphore: %s", result)
	}

	return &SemaphoreVK{
		device:    device,
		semaphore: semaphore,
	}, nil
}

// GetSemaphore returns the Vulkan semaphore handle
func (s *SemaphoreVK) GetSemaphore() vulkan.Semaphore {
	return s.semaphore
}

// Destroy destroys the semaphore
func (s *SemaphoreVK) Destroy() {
	if s.semaphore != vulkan.NullSemaphore {
		vulkan.DestroySemaphore(s.device, s.semaphore, nil)
		s.semaphore = vulkan.NullSemaphore
	}
}

// EventVK wraps Vulkan event
// It provides fine-grained synchronization primitive
type EventVK struct {
	device vulkan.Device
	event  vulkan.Event
}

// NewEventVK creates a new Vulkan event
func NewEventVK(device vulkan.Device) (*EventVK, error) {
	createInfo := vulkan.EventCreateInfo{
		SType: vulkan.StructureTypeEventCreateInfo,
	}

	var event vulkan.Event
	if result := vulkan.CreateEvent(device, &createInfo, nil, &event); result != vulkan.Success {
		return nil, fmt.Errorf("failed to create event: %s", result)
	}

	return &EventVK{
		device: device,
		event:  event,
	}, nil
}

// GetEvent returns the Vulkan event handle
func (e *EventVK) GetEvent() vulkan.Event {
	return e.event
}

// Set sets the event
func (e *EventVK) Set() error {
	if result := vulkan.SetEvent(e.device, e.event); result != vulkan.Success {
		return fmt.Errorf("failed to set event: %s", result)
	}
	return nil
}

// Reset resets the event
func (e *EventVK) Reset() error {
	if result := vulkan.ResetEvent(e.device, e.event); result != vulkan.Success {
		return fmt.Errorf("failed to reset event: %s", result)
	}
	return nil
}

// IsSet checks if the event is set
func (e *EventVK) IsSet() bool {
	status := vulkan.GetEventStatus(e.device, e.event)
	return status == vulkan.EventSet
}

// Destroy destroys the event
func (e *EventVK) Destroy() {
	if e.event != vulkan.NullEvent {
		vulkan.DestroyEvent(e.device, e.event, nil)
		e.event = vulkan.NullEvent
	}
}

// SynchronizationManagerVK manages synchronization objects
// It provides centralized management of fences, semaphores, and events
type SynchronizationManagerVK struct {
	device     vulkan.Device
	fences     []*FenceVK
	semaphores []*SemaphoreVK
	events     []*EventVK
}

// NewSynchronizationManagerVK creates a new synchronization manager
func NewSynchronizationManagerVK(device vulkan.Device) *SynchronizationManagerVK {
	return &SynchronizationManagerVK{
		device:     device,
		fences:     make([]*FenceVK, 0),
		semaphores: make([]*SemaphoreVK, 0),
		events:     make([]*EventVK, 0),
	}
}

// CreateFence creates a new fence
func (sm *SynchronizationManagerVK) CreateFence(signaled bool) (*FenceVK, error) {
	fence, err := NewFenceVK(sm.device, signaled)
	if err != nil {
		return nil, err
	}
	sm.fences = append(sm.fences, fence)
	return fence, nil
}

// CreateSemaphore creates a new semaphore
func (sm *SynchronizationManagerVK) CreateSemaphore() (*SemaphoreVK, error) {
	semaphore, err := NewSemaphoreVK(sm.device)
	if err != nil {
		return nil, err
	}
	sm.semaphores = append(sm.semaphores, semaphore)
	return semaphore, nil
}

// CreateEvent creates a new event
func (sm *SynchronizationManagerVK) CreateEvent() (*EventVK, error) {
	event, err := NewEventVK(sm.device)
	if err != nil {
		return nil, err
	}
	sm.events = append(sm.events, event)
	return event, nil
}

// WaitForFences waits for multiple fences
func (sm *SynchronizationManagerVK) WaitForFences(fences []*FenceVK, waitAll bool, timeout uint64) error {
	if len(fences) == 0 {
		return nil
	}

	vkFences := make([]vulkan.Fence, len(fences))
	for i, fence := range fences {
		vkFences[i] = fence.GetFence()
	}

	waitAllFlag := vulkan.False
	if waitAll {
		waitAllFlag = vulkan.True
	}

	if result := vulkan.WaitForFences(sm.device, uint32(len(vkFences)), vkFences, waitAllFlag, timeout); result != vulkan.Success {
		return fmt.Errorf("failed to wait for fences: %s", result)
	}

	for _, fence := range fences {
		fence.signaled = true
	}

	return nil
}

// ResetFences resets multiple fences
func (sm *SynchronizationManagerVK) ResetFences(fences []*FenceVK) error {
	if len(fences) == 0 {
		return nil
	}

	vkFences := make([]vulkan.Fence, len(fences))
	for i, fence := range fences {
		vkFences[i] = fence.GetFence()
	}

	if result := vulkan.ResetFences(sm.device, uint32(len(vkFences)), vkFences); result != vulkan.Success {
		return fmt.Errorf("failed to reset fences: %s", result)
	}

	for _, fence := range fences {
		fence.signaled = false
	}

	return nil
}

// Destroy destroys all synchronization objects
func (sm *SynchronizationManagerVK) Destroy() {
	for _, fence := range sm.fences {
		fence.Destroy()
	}
	sm.fences = sm.fences[:0]

	for _, semaphore := range sm.semaphores {
		semaphore.Destroy()
	}
	sm.semaphores = sm.semaphores[:0]

	for _, event := range sm.events {
		event.Destroy()
	}
	sm.events = sm.events[:0]
}
