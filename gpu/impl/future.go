package impl

import (
	"context"
	"fmt"
	"sync"
	"time"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

// FutureManager manages Future objects and their completion
type FutureManager struct {
	mu          sync.RWMutex
	futures     map[uint64]*FutureState
	nextID      uint64
	callbackMgr *CallbackManager
}

// FutureState represents the state of a future
type FutureState struct {
	ID          uint64
	Completed   bool
	CompletedAt time.Time
	Error       error
	Result      interface{}
	ctx         context.Context
	cancel      context.CancelFunc
	done        chan struct{}
	once        sync.Once
}

// NewFutureManager creates a new future manager
func NewFutureManager(callbackMgr *CallbackManager) *FutureManager {
	return &FutureManager{
		futures:     make(map[uint64]*FutureState),
		callbackMgr: callbackMgr,
	}
}

// CreateFuture creates a new future
func (fm *FutureManager) CreateFuture() Future {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.nextID++
	futureID := fm.nextID

	ctx, cancel := context.WithCancel(context.Background())
	state := &FutureState{
		ID:     futureID,
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}

	fm.futures[futureID] = state

	return Future{Id: futureID}
}

// CompleteFuture marks a future as completed
func (fm *FutureManager) CompleteFuture(futureID uint64, result interface{}, err error) bool {
	fm.mu.Lock()
	state, exists := fm.futures[futureID]
	fm.mu.Unlock()

	if !exists {
		return false
	}

	state.once.Do(func() {
		state.Completed = true
		state.CompletedAt = time.Now()
		state.Result = result
		state.Error = err
		close(state.done)
		state.cancel()
	})

	// Notify callback manager
	return fm.callbackMgr.Complete(futureID)
}

// WaitForFuture waits for a future to complete
func (fm *FutureManager) WaitForFuture(futureID uint64, timeout time.Duration) (bool, interface{}, error) {
	fm.mu.RLock()
	state, exists := fm.futures[futureID]
	fm.mu.RUnlock()

	if !exists {
		return false, nil, fmt.Errorf("future %d not found", futureID)
	}

	ctx := state.ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	select {
	case <-state.done:
		return true, state.Result, state.Error
	case <-ctx.Done():
		return false, nil, ctx.Err()
	}
}

// IsFutureCompleted checks if a future is completed
func (fm *FutureManager) IsFutureCompleted(futureID uint64) bool {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	if state, exists := fm.futures[futureID]; exists {
		return state.Completed
	}
	return false
}

// CleanupFuture removes a completed future
func (fm *FutureManager) CleanupFuture(futureID uint64) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if state, exists := fm.futures[futureID]; exists {
		state.cancel()
		delete(fm.futures, futureID)
	}
}

// CleanupExpiredFutures removes futures that completed more than the specified duration ago
func (fm *FutureManager) CleanupExpiredFutures(maxAge time.Duration) int {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	now := time.Now()
	cleaned := 0

	for id, state := range fm.futures {
		if state.Completed && now.Sub(state.CompletedAt) > maxAge {
			state.cancel()
			delete(fm.futures, id)
			cleaned++
		}
	}

	return cleaned
}

// GetFutureCount returns the number of tracked futures
func (fm *FutureManager) GetFutureCount() (total, completed int) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	total = len(fm.futures)
	for _, state := range fm.futures {
		if state.Completed {
			completed++
		}
	}
	return
}

// Shutdown shuts down the future manager
func (fm *FutureManager) Shutdown() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	for _, state := range fm.futures {
		state.cancel()
	}
	fm.futures = make(map[uint64]*FutureState)
}

// Global future manager
var globalFutureManager *FutureManager

// InitializeFutureManager initializes the global future manager
func InitializeFutureManager() {
	globalFutureManager = NewFutureManager(globalCallbackManager)
}

// GetGlobalFutureManager returns the global future manager
func GetGlobalFutureManager() *FutureManager {
	return globalFutureManager
}

// GenerateFutureId creates a new future and returns its ID
func GenerateFutureId() uint64 {
	if globalFutureManager == nil {
		InitializeFutureManager()
	}
	future := globalFutureManager.CreateFuture()
	return future.Id
}
