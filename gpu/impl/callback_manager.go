package impl

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

// CallbackTask represents a callback task with its execution context
type CallbackTask struct {
	ID          uint64
	FutureID    uint64
	Mode        CallbackMode
	ExecuteFunc func()
	TaskCtx     context.Context
	CancelFunc  context.CancelFunc
	CreatedAt   time.Time
	IsCompleted int32 // atomic flag
	IsCancelled int32 // atomic flag
}

// CallbackTaskManager manages WebGPU callbacks and their execution
type CallbackTaskManager struct {
	mu         sync.RWMutex
	tasks      map[uint64]*CallbackTask // task ID -> task
	futureTask map[uint64]*CallbackTask // future ID -> task
	idCounter  int64                    // atomic counter for task IDs
	ctx        context.Context
	cancelFunc context.CancelFunc
	eventQueue chan *CallbackTask
	workerWG   sync.WaitGroup
	isShutdown int32 // atomic flag
}

// NewCallbackTaskManager creates a new callback task manager
func NewCallbackTaskManager() *CallbackTaskManager {
	ctx, cancel := context.WithCancel(context.Background())

	manager := &CallbackTaskManager{
		tasks:      make(map[uint64]*CallbackTask),
		futureTask: make(map[uint64]*CallbackTask),
		ctx:        ctx,
		cancelFunc: cancel,
		eventQueue: make(chan *CallbackTask, 1000), // buffered channel for event processing
	}

	// Start worker goroutine for processing events
	manager.workerWG.Add(1)
	go manager.processEventLoop()

	return manager
}

// RegisterCallback registers a callback task with the manager
func (m *CallbackTaskManager) RegisterCallback(futureID uint64, mode CallbackMode, executeFunc func()) *CallbackTask {
	if atomic.LoadInt32(&m.isShutdown) != 0 {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	taskID := uint64(atomic.AddInt64(&m.idCounter, 1))
	taskCtx, cancelFunc := context.WithCancel(m.ctx)

	task := &CallbackTask{
		ID:          taskID,
		FutureID:    futureID,
		Mode:        mode,
		ExecuteFunc: executeFunc,
		TaskCtx:     taskCtx,
		CancelFunc:  cancelFunc,
		CreatedAt:   time.Now(),
	}

	m.tasks[taskID] = task
	m.futureTask[futureID] = task

	return task
}

// CompleteCallback marks a callback as completed and schedules its execution
func (m *CallbackTaskManager) CompleteCallback(futureID uint64) bool {
	if atomic.LoadInt32(&m.isShutdown) != 0 {
		return false
	}

	m.mu.Lock()
	task, exists := m.futureTask[futureID]
	m.mu.Unlock()

	if !exists || atomic.LoadInt32(&task.IsCancelled) != 0 {
		return false
	}

	if !atomic.CompareAndSwapInt32(&task.IsCompleted, 0, 1) {
		return false // already completed
	}

	switch task.Mode {
	case CallbackModeWaitAnyOnly:
		// Only execute when explicitly waited for
		return true

	case CallbackModeAllowProcessEvents:
		// Queue for processing in ProcessEvents
		select {
		case m.eventQueue <- task:
			return true
		default:
			// Queue is full, execute synchronously as fallback
			m.executeCallback(task)
			return true
		}

	case CallbackModeAllowSpontaneous:
		// Execute immediately in background
		go m.executeCallback(task)
		return true

	default:
		return false
	}
}

// CancelCallback cancels a pending callback
func (m *CallbackTaskManager) CancelCallback(futureID uint64) bool {
	m.mu.Lock()
	task, exists := m.futureTask[futureID]
	m.mu.Unlock()

	if !exists {
		return false
	}

	if atomic.CompareAndSwapInt32(&task.IsCancelled, 0, 1) {
		task.CancelFunc()
		return true
	}

	return false
}

// ProcessEvents processes callbacks that are queued for event processing
func (m *CallbackTaskManager) ProcessEvents() error {
	if atomic.LoadInt32(&m.isShutdown) != 0 {
		return nil
	}

	// Process all queued events
	for {
		select {
		case task := <-m.eventQueue:
			if atomic.LoadInt32(&task.IsCancelled) == 0 && atomic.LoadInt32(&task.IsCompleted) != 0 {
				m.executeCallback(task)
			}
		default:
			return nil // no more events to process
		}
	}
}

// WaitForCallback waits for a specific callback to complete and executes it if in WaitAnyOnly mode
func (m *CallbackTaskManager) WaitForCallback(futureID uint64, timeoutNS uint64) bool {
	if atomic.LoadInt32(&m.isShutdown) != 0 {
		return false
	}

	m.mu.RLock()
	task, exists := m.futureTask[futureID]
	m.mu.RUnlock()

	if !exists {
		return false
	}

	// Check if already completed
	if atomic.LoadInt32(&task.IsCompleted) != 0 {
		if task.Mode == CallbackModeWaitAnyOnly && atomic.LoadInt32(&task.IsCancelled) == 0 {
			m.executeCallback(task)
		}
		return true
	}

	// Wait for completion with timeout
	var timeout <-chan time.Time
	if timeoutNS > 0 {
		timeout = time.After(time.Duration(timeoutNS) * time.Nanosecond)
	}

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return false // timeout

		case <-task.TaskCtx.Done():
			return false // cancelled

		case <-ticker.C:
			if atomic.LoadInt32(&task.IsCompleted) != 0 {
				if task.Mode == CallbackModeWaitAnyOnly && atomic.LoadInt32(&task.IsCancelled) == 0 {
					m.executeCallback(task)
				}
				return true
			}
		}
	}
}

// executeCallback safely executes a callback task
func (m *CallbackTaskManager) executeCallback(task *CallbackTask) {
	if task == nil || atomic.LoadInt32(&task.IsCancelled) != 0 {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			// Log panic but don't crash the application
			// In a real implementation, this would use proper logging
		}
	}()

	if task.ExecuteFunc != nil {
		task.ExecuteFunc()
	}

	// Clean up the task
	m.cleanupTask(task)
}

// cleanupTask removes a completed task from the manager
func (m *CallbackTaskManager) cleanupTask(task *CallbackTask) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.tasks, task.ID)
	delete(m.futureTask, task.FutureID)
}

// processEventLoop runs in background to handle spontaneous callback execution
func (m *CallbackTaskManager) processEventLoop() {
	defer m.workerWG.Done()

	for {
		select {
		case <-m.ctx.Done():
			return

		case task := <-m.eventQueue:
			if atomic.LoadInt32(&task.IsCancelled) == 0 && atomic.LoadInt32(&task.IsCompleted) != 0 {
				m.executeCallback(task)
			}
		}
	}
}

// GetPendingCallbackCount returns the number of pending callbacks
func (m *CallbackTaskManager) GetPendingCallbackCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.tasks)
}

// CancelAllCallbacks cancels all pending callbacks
func (m *CallbackTaskManager) CancelAllCallbacks() {
	m.mu.Lock()
	tasks := make([]*CallbackTask, 0, len(m.tasks))
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}
	m.mu.Unlock()

	for _, task := range tasks {
		atomic.StoreInt32(&task.IsCancelled, 1)
		task.CancelFunc()
	}
}

// Shutdown gracefully shuts down the callback manager
func (m *CallbackTaskManager) Shutdown() error {
	if !atomic.CompareAndSwapInt32(&m.isShutdown, 0, 1) {
		return nil // already shutdown
	}

	// Cancel all pending callbacks
	m.CancelAllCallbacks()

	// Close event queue and cancel context
	close(m.eventQueue)
	m.cancelFunc()

	// Wait for worker goroutine to finish
	m.workerWG.Wait()

	// Clear all tasks
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tasks = make(map[uint64]*CallbackTask)
	m.futureTask = make(map[uint64]*CallbackTask)

	return nil
}

// IsShutdown returns whether the manager has been shut down
func (m *CallbackTaskManager) IsShutdown() bool {
	return atomic.LoadInt32(&m.isShutdown) != 0
}

// Additional helper methods for managing callback states
func (m *CallbackTaskManager) GetCallbackStatus(futureID uint64) (bool, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, exists := m.futureTask[futureID]
	if !exists {
		return false, false // not found, not completed
	}

	isCompleted := atomic.LoadInt32(&task.IsCompleted) != 0
	isCancelled := atomic.LoadInt32(&task.IsCancelled) != 0

	return true, isCompleted && !isCancelled // found, success status
}

// GetCallbackStats returns statistics about callback management
func (m *CallbackTaskManager) GetCallbackStats() (pending, completed, cancelled int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, task := range m.tasks {
		isCompleted := atomic.LoadInt32(&task.IsCompleted) != 0
		isCancelled := atomic.LoadInt32(&task.IsCancelled) != 0

		if isCancelled {
			cancelled++
		} else if isCompleted {
			completed++
		} else {
			pending++
		}
	}

	return pending, completed, cancelled
}

// WaitForAnyCallback waits for any callback to complete within the timeout
func (m *CallbackTaskManager) WaitForAnyCallback(timeoutNS uint64) (uint64, bool) {
	if atomic.LoadInt32(&m.isShutdown) != 0 {
		return 0, false
	}

	// Get a snapshot of current tasks
	m.mu.RLock()
	tasks := make([]*CallbackTask, 0, len(m.futureTask))
	for _, task := range m.futureTask {
		tasks = append(tasks, task)
	}
	m.mu.RUnlock()

	if len(tasks) == 0 {
		return 0, false
	}

	// Check if any are already completed
	for _, task := range tasks {
		if atomic.LoadInt32(&task.IsCompleted) != 0 && atomic.LoadInt32(&task.IsCancelled) == 0 {
			if task.Mode == CallbackModeWaitAnyOnly {
				m.executeCallback(task)
			}
			return task.FutureID, true
		}
	}

	// Wait for any to complete with timeout
	var timeout <-chan time.Time
	if timeoutNS > 0 {
		timeout = time.After(time.Duration(timeoutNS) * time.Nanosecond)
	}

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return 0, false // timeout

		case <-m.ctx.Done():
			return 0, false // shutdown

		case <-ticker.C:
			// Check all tasks again
			for _, task := range tasks {
				if atomic.LoadInt32(&task.IsCompleted) != 0 && atomic.LoadInt32(&task.IsCancelled) == 0 {
					if task.Mode == CallbackModeWaitAnyOnly {
						m.executeCallback(task)
					}
					return task.FutureID, true
				}
			}
		}
	}
}

// CallbackHelper provides helper functions for common callback patterns
type CallbackHelper struct {
	manager *CallbackTaskManager
}

// NewCallbackHelper creates a new callback helper
func NewCallbackHelper(manager *CallbackTaskManager) *CallbackHelper {
	return &CallbackHelper{
		manager: manager,
	}
}

// RegisterRequestAdapterCallback registers a RequestAdapter callback
func (h *CallbackHelper) RegisterRequestAdapterCallback(futureID uint64, callback RequestAdapterCallbackInfo,
	status RequestAdapterStatus, adapter Adapter, message string) *CallbackTask {

	executeFunc := func() {
		if callback.Callback != nil {
			callback.Callback(status, adapter, message)
		}
	}

	return h.manager.RegisterCallback(futureID, callback.Mode, executeFunc)
}

// RegisterRequestDeviceCallback registers a RequestDevice callback
func (h *CallbackHelper) RegisterRequestDeviceCallback(futureID uint64, callback RequestDeviceCallbackInfo,
	status RequestDeviceStatus, device Device, message string) *CallbackTask {

	executeFunc := func() {
		if callback.Callback != nil {
			callback.Callback(status, device, message)
		}
	}

	return h.manager.RegisterCallback(futureID, callback.Mode, executeFunc)
}

// RegisterBufferMapCallback registers a Buffer map callback
func (h *CallbackHelper) RegisterBufferMapCallback(futureID uint64, callback BufferMapCallbackInfo,
	status MapAsyncStatus, message string) *CallbackTask {

	executeFunc := func() {
		if callback.Callback != nil {
			callback.Callback(status, message)
		}
	}

	return h.manager.RegisterCallback(futureID, callback.Mode, executeFunc)
}

// RegisterCreateComputePipelineAsyncCallback registers a CreateComputePipelineAsync callback
func (h *CallbackHelper) RegisterCreateComputePipelineAsyncCallback(futureID uint64, callback CreateComputePipelineAsyncCallbackInfo,
	status CreatePipelineAsyncStatus, pipeline ComputePipeline, message string) *CallbackTask {

	executeFunc := func() {
		if callback.Callback != nil {
			callback.Callback(status, pipeline, message)
		}
	}

	return h.manager.RegisterCallback(futureID, callback.Mode, executeFunc)
}

// RegisterCreateRenderPipelineAsyncCallback registers a CreateRenderPipelineAsync callback
func (h *CallbackHelper) RegisterCreateRenderPipelineAsyncCallback(futureID uint64, callback CreateRenderPipelineAsyncCallbackInfo,
	status CreatePipelineAsyncStatus, pipeline RenderPipeline, message string) *CallbackTask {

	executeFunc := func() {
		if callback.Callback != nil {
			callback.Callback(status, pipeline, message)
		}
	}

	return h.manager.RegisterCallback(futureID, callback.Mode, executeFunc)
}

// RegisterDeviceLostCallback registers a DeviceLost callback
func (h *CallbackHelper) RegisterDeviceLostCallback(futureID uint64, callback DeviceLostCallbackInfo,
	device Device, reason DeviceLostReason, message string) *CallbackTask {

	executeFunc := func() {
		if callback.Callback != nil {
			callback.Callback(device, reason, message)
		}
	}

	return h.manager.RegisterCallback(futureID, callback.Mode, executeFunc)
}

// RegisterPopErrorScopeCallback registers a PopErrorScope callback
func (h *CallbackHelper) RegisterPopErrorScopeCallback(futureID uint64, callback PopErrorScopeCallbackInfo,
	status PopErrorScopeStatus, errorType ErrorType, message string) *CallbackTask {

	executeFunc := func() {
		if callback.Callback != nil {
			callback.Callback(status, errorType, message)
		}
	}

	return h.manager.RegisterCallback(futureID, callback.Mode, executeFunc)
}

// RegisterQueueWorkDoneCallback registers a QueueWorkDone callback
func (h *CallbackHelper) RegisterQueueWorkDoneCallback(futureID uint64, callback QueueWorkDoneCallbackInfo,
	status QueueWorkDoneStatus, message string) *CallbackTask {

	executeFunc := func() {
		if callback.Callback != nil {
			callback.Callback(status, message)
		}
	}

	return h.manager.RegisterCallback(futureID, callback.Mode, executeFunc)
}

// RegisterCompilationInfoCallback registers a CompilationInfo callback
func (h *CallbackHelper) RegisterCompilationInfoCallback(futureID uint64, callback CompilationInfoCallbackInfo,
	status CompilationInfoRequestStatus, compilationInfo CompilationInfo) *CallbackTask {

	executeFunc := func() {
		if callback.Callback != nil {
			callback.Callback(status, compilationInfo)
		}
	}

	return h.manager.RegisterCallback(futureID, callback.Mode, executeFunc)
}

// RegisterUncapturedErrorCallback registers an UncapturedError callback
// Note: UncapturedError callbacks don't use the future system, they're immediate
func (h *CallbackHelper) RegisterUncapturedErrorCallback(callback UncapturedErrorCallbackInfo,
	device Device, errorType ErrorType, message string) {

	// UncapturedError callbacks are executed immediately
	if callback.Callback != nil {
		callback.Callback(device, errorType, message)
	}
}
