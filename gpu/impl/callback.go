package impl

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

// CallbackTask represents a pending callback
type CallbackTask struct {
	FutureID uint64
	Mode     CallbackMode
	Execute  func()
	Done     chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc
}

// CallbackManager manages WebGPU callbacks with improved synchronization
type CallbackManager struct {
	mu            sync.RWMutex
	tasks         map[uint64]*CallbackTask
	completedIDs  map[uint64]bool
	shutdown      int32
	wg            sync.WaitGroup
	eventCond     *sync.Cond
	pendingEvents int32
}

// NewCallbackManager creates a new callback manager
func NewCallbackManager() *CallbackManager {
	cm := &CallbackManager{
		tasks:        make(map[uint64]*CallbackTask),
		completedIDs: make(map[uint64]bool),
	}
	cm.eventCond = sync.NewCond(&cm.mu)
	return cm
}

// Register registers a callback for execution with improved synchronization
func (m *CallbackManager) Register(futureID uint64, mode CallbackMode, fn func()) {
	if atomic.LoadInt32(&m.shutdown) != 0 {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	m.mu.Lock()
	defer m.mu.Unlock()

	task := &CallbackTask{
		FutureID: futureID,
		Mode:     mode,
		Execute:  fn,
		Done:     make(chan struct{}),
		ctx:      ctx,
		cancel:   cancel,
	}

	m.tasks[futureID] = task
	atomic.AddInt32(&m.pendingEvents, 1)
	m.eventCond.Signal()
}

// Complete marks a callback as ready and executes it based on its mode
func (m *CallbackManager) Complete(futureID uint64) bool {
	if atomic.LoadInt32(&m.shutdown) != 0 {
		return false
	}

	m.mu.Lock()
	task, exists := m.tasks[futureID]
	if !exists {
		m.mu.Unlock()
		return false
	}

	// Mark as completed
	m.completedIDs[futureID] = true
	m.mu.Unlock()

	// Handle different callback modes with proper synchronization
	switch task.Mode {
	case CallbackModeAllowSpontaneous:
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			m.executeSafe(task)
		}()
	case CallbackModeAllowProcessEvents:
		// Mark as ready for ProcessEvents
		close(task.Done)
		m.mu.Lock()
		atomic.AddInt32(&m.pendingEvents, 1)
		m.eventCond.Signal()
		m.mu.Unlock()
	case CallbackModeWaitAnyOnly:
		// Only execute when explicitly waited for
		close(task.Done)
	}

	return true
}

// ProcessEvents processes callbacks that are ready for event processing
func (m *CallbackManager) ProcessEvents() {
	if atomic.LoadInt32(&m.shutdown) != 0 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var readyTasks []*CallbackTask
	for futureID, task := range m.tasks {
		if task.Mode == CallbackModeAllowProcessEvents && m.completedIDs[futureID] {
			select {
			case <-task.Done:
				readyTasks = append(readyTasks, task)
			default:
			}
		}
	}

	// Execute ready tasks without holding the main lock
	for _, task := range readyTasks {
		m.mu.Unlock()
		m.executeSafe(task)
		m.mu.Lock()
	}
}

// Wait waits for a specific callback to complete with timeout support
func (m *CallbackManager) Wait(futureID uint64, timeout time.Duration) bool {
	if atomic.LoadInt32(&m.shutdown) != 0 {
		return false
	}

	m.mu.RLock()
	task, exists := m.tasks[futureID]
	m.mu.RUnlock()

	if !exists {
		return false
	}

	ctx := task.ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	select {
	case <-task.Done:
		if task.Mode == CallbackModeWaitAnyOnly {
			m.executeSafe(task)
		}
		return true
	case <-ctx.Done():
		return false
	}
}

// WaitAny waits for any of the given futures with improved synchronization
func (m *CallbackManager) WaitAny(futures []uint64, timeout time.Duration) (uint64, bool) {
	if atomic.LoadInt32(&m.shutdown) != 0 || len(futures) == 0 {
		return 0, false
	}

	// Create a context for timeout handling
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// Channel to receive completion notifications
	done := make(chan uint64, len(futures))

	// Set up goroutines to wait for each future
	var wg sync.WaitGroup
	for _, futureID := range futures {
		wg.Add(1)
		go func(id uint64) {
			defer wg.Done()
			if m.Wait(id, 0) { // No timeout for individual waits
				select {
				case done <- id:
				default:
				}
			}
		}(futureID)
	}

	// Goroutine to close done channel when all futures are processed
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for first completion or timeout
	select {
	case futureID := <-done:
		return futureID, true
	case <-ctx.Done():
		return 0, false
	}
}

// WaitForEvents waits for events to be available for processing
func (m *CallbackManager) WaitForEvents(timeout time.Duration) bool {
	if atomic.LoadInt32(&m.shutdown) != 0 {
		return false
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// If events are already pending, return immediately
	if atomic.LoadInt32(&m.pendingEvents) > 0 {
		return true
	}

	// Wait for events with timeout
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		done := make(chan bool, 1)
		go func() {
			m.eventCond.Wait()
			done <- true
		}()

		m.mu.Unlock()
		select {
		case <-done:
			m.mu.Lock()
			return atomic.LoadInt32(&m.pendingEvents) > 0
		case <-timer.C:
			m.mu.Lock()
			return false
		}
	} else {
		m.eventCond.Wait()
		return atomic.LoadInt32(&m.pendingEvents) > 0
	}
}

// executeSafe safely executes a callback task
func (m *CallbackManager) executeSafe(task *CallbackTask) {
	defer func() {
		if r := recover(); r != nil {
			// Log error in real implementation
		}
		m.cleanup(task.FutureID)
		atomic.AddInt32(&m.pendingEvents, -1)
	}()

	if task.Execute != nil {
		task.Execute()
	}
}

// cleanup removes a completed task
func (m *CallbackManager) cleanup(futureID uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if task, exists := m.tasks[futureID]; exists {
		task.cancel()
		delete(m.tasks, futureID)
		delete(m.completedIDs, futureID)
	}
}

// Shutdown gracefully shuts down the callback manager
func (m *CallbackManager) Shutdown() {
	if !atomic.CompareAndSwapInt32(&m.shutdown, 0, 1) {
		return
	}

	m.mu.Lock()
	// Cancel all pending tasks
	for _, task := range m.tasks {
		task.cancel()
		select {
		case <-task.Done:
		default:
			close(task.Done)
		}
	}

	// Clear all data
	m.tasks = make(map[uint64]*CallbackTask)
	m.completedIDs = make(map[uint64]bool)
	atomic.StoreInt32(&m.pendingEvents, 0)

	// Wake up any waiting goroutines
	m.eventCond.Broadcast()
	m.mu.Unlock()

	// Wait for all spawned goroutines to complete
	m.wg.Wait()
}

// GetPendingEventCount returns the number of pending events
func (m *CallbackManager) GetPendingEventCount() int32 {
	return atomic.LoadInt32(&m.pendingEvents)
}

// IsShutdown checks if the callback manager is shut down
func (m *CallbackManager) IsShutdown() bool {
	return atomic.LoadInt32(&m.shutdown) != 0
}

// CallbackRegistry provides helper functions for registering callbacks
type CallbackRegistry struct {
	manager *CallbackManager
}

// newCallbackRegistry creates a new callback registry
func newCallbackRegistry(manager *CallbackManager) *CallbackRegistry {
	return &CallbackRegistry{manager: manager}
}

// RequestAdapter registers a RequestAdapter callback
func (r *CallbackRegistry) RequestAdapter(futureID uint64, callback RequestAdapterCallbackInfo, status RequestAdapterStatus, adapter Adapter, message string) {
	r.manager.Register(futureID, callback.Mode, func() {
		if callback.Callback != nil {
			callback.Callback(status, adapter, message)
		}
	})
}

// RequestDevice registers a RequestDevice callback
func (r *CallbackRegistry) RequestDevice(futureID uint64, callback RequestDeviceCallbackInfo, status RequestDeviceStatus, device Device, message string) {
	r.manager.Register(futureID, callback.Mode, func() {
		if callback.Callback != nil {
			callback.Callback(status, device, message)
		}
	})
}

// BufferMap registers a Buffer map callback
func (r *CallbackRegistry) BufferMap(futureID uint64, callback BufferMapCallbackInfo, status MapAsyncStatus, message string) {
	r.manager.Register(futureID, callback.Mode, func() {
		if callback.Callback != nil {
			callback.Callback(status, message)
		}
	})
}

// CreateComputePipelineAsync registers a CreateComputePipelineAsync callback
func (r *CallbackRegistry) CreateComputePipelineAsync(futureID uint64, callback CreateComputePipelineAsyncCallbackInfo, status CreatePipelineAsyncStatus, pipeline ComputePipeline, message string) {
	r.manager.Register(futureID, callback.Mode, func() {
		if callback.Callback != nil {
			callback.Callback(status, pipeline, message)
		}
	})
}

// CreateRenderPipelineAsync registers a CreateRenderPipelineAsync callback
func (r *CallbackRegistry) CreateRenderPipelineAsync(futureID uint64, callback CreateRenderPipelineAsyncCallbackInfo, status CreatePipelineAsyncStatus, pipeline RenderPipeline, message string) {
	r.manager.Register(futureID, callback.Mode, func() {
		if callback.Callback != nil {
			callback.Callback(status, pipeline, message)
		}
	})
}

// DeviceLost registers a DeviceLost callback
func (r *CallbackRegistry) DeviceLost(futureID uint64, callback DeviceLostCallbackInfo, device Device, reason DeviceLostReason, message string) {
	r.manager.Register(futureID, callback.Mode, func() {
		if callback.Callback != nil {
			callback.Callback(device, reason, message)
		}
	})
}

// PopErrorScope registers a PopErrorScope callback
func (r *CallbackRegistry) PopErrorScope(futureID uint64, callback PopErrorScopeCallbackInfo, status PopErrorScopeStatus, errorType ErrorType, message string) {
	r.manager.Register(futureID, callback.Mode, func() {
		if callback.Callback != nil {
			callback.Callback(status, errorType, message)
		}
	})
}

// QueueWorkDone registers a QueueWorkDone callback
func (r *CallbackRegistry) QueueWorkDone(futureID uint64, callback QueueWorkDoneCallbackInfo, status QueueWorkDoneStatus, message string) {
	r.manager.Register(futureID, callback.Mode, func() {
		if callback.Callback != nil {
			callback.Callback(status, message)
		}
	})
}

// CompilationInfo registers a CompilationInfo callback
func (r *CallbackRegistry) CompilationInfo(futureID uint64, callback CompilationInfoCallbackInfo, status CompilationInfoRequestStatus, compilationInfo CompilationInfo) {
	r.manager.Register(futureID, callback.Mode, func() {
		if callback.Callback != nil {
			callback.Callback(status, compilationInfo)
		}
	})
}

// UncapturedError executes an UncapturedError callback immediately (no future system)
func (r *CallbackRegistry) UncapturedError(callback UncapturedErrorCallbackInfo, device Device, errorType ErrorType, message string) {
	if callback.Callback != nil {
		callback.Callback(device, errorType, message)
	}
}

// Global callback manager instance
var globalCallbackManager = NewCallbackManager()
var globalCallbackRegistry = newCallbackRegistry(globalCallbackManager)

// GlobalCallbackManager returns the global callback manager
func GlobalCallbackManager() *CallbackManager {
	return globalCallbackManager
}

// GlobalCallbackRegistry returns the global callback registry
func GlobalCallbackRegistry() *CallbackRegistry {
	return globalCallbackRegistry
}
