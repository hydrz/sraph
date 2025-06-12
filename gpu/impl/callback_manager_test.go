package impl

import (
	"sync"
	"testing"
	"time"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

func TestCallbackTaskManager_Basic(t *testing.T) {
	manager := NewCallbackTaskManager()
	defer manager.Shutdown()

	if manager.IsShutdown() {
		t.Fatal("manager should not be shutdown initially")
	}

	count := manager.GetPendingCallbackCount()
	if count != 0 {
		t.Errorf("expected 0 pending callbacks, got %d", count)
	}
}

func TestCallbackTaskManager_RegisterAndComplete(t *testing.T) {
	manager := NewCallbackTaskManager()
	defer manager.Shutdown()

	executed := false
	futureID := uint64(123)

	task := manager.RegisterCallback(futureID, CallbackModeAllowProcessEvents, func() {
		executed = true
	})

	if task == nil {
		t.Fatal("expected task to be registered")
	}

	if task.FutureID != futureID {
		t.Errorf("expected future ID %d, got %d", futureID, task.FutureID)
	}

	// Complete the callback
	success := manager.CompleteCallback(futureID)
	if !success {
		t.Error("expected callback completion to succeed")
	}

	// Process events to execute the callback
	err := manager.ProcessEvents()
	if err != nil {
		t.Errorf("ProcessEvents failed: %v", err)
	}

	// Give it a moment for async execution
	time.Sleep(10 * time.Millisecond)

	if !executed {
		t.Error("callback was not executed")
	}
}

func TestCallbackTaskManager_WaitForCallback(t *testing.T) {
	manager := NewCallbackTaskManager()
	defer manager.Shutdown()

	executed := false
	futureID := uint64(456)

	task := manager.RegisterCallback(futureID, CallbackModeWaitAnyOnly, func() {
		executed = true
	})

	if task == nil {
		t.Fatal("expected task to be registered")
	}

	// Complete callback in background
	go func() {
		time.Sleep(10 * time.Millisecond)
		manager.CompleteCallback(futureID)
	}()

	// Wait for completion
	success := manager.WaitForCallback(futureID, uint64(100*time.Millisecond))
	if !success {
		t.Error("expected wait to succeed")
	}

	if !executed {
		t.Error("callback was not executed during wait")
	}
}

func TestCallbackTaskManager_CancelCallback(t *testing.T) {
	manager := NewCallbackTaskManager()
	defer manager.Shutdown()

	executed := false
	futureID := uint64(789)

	task := manager.RegisterCallback(futureID, CallbackModeAllowProcessEvents, func() {
		executed = true
	})

	if task == nil {
		t.Fatal("expected task to be registered")
	}

	// Cancel the callback
	success := manager.CancelCallback(futureID)
	if !success {
		t.Error("expected callback cancellation to succeed")
	}

	// Try to complete the callback
	completed := manager.CompleteCallback(futureID)
	if completed {
		t.Error("expected completion to fail after cancellation")
	}

	// Process events
	manager.ProcessEvents()
	time.Sleep(10 * time.Millisecond)

	if executed {
		t.Error("callback should not have been executed after cancellation")
	}
}

func TestCallbackTaskManager_SpontaneousMode(t *testing.T) {
	manager := NewCallbackTaskManager()
	defer manager.Shutdown()

	var executed bool
	var mu sync.Mutex
	futureID := uint64(999)

	task := manager.RegisterCallback(futureID, CallbackModeAllowSpontaneous, func() {
		mu.Lock()
		executed = true
		mu.Unlock()
	})

	if task == nil {
		t.Fatal("expected task to be registered")
	}

	// Complete the callback
	success := manager.CompleteCallback(futureID)
	if !success {
		t.Error("expected callback completion to succeed")
	}

	// Wait for spontaneous execution
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	wasExecuted := executed
	mu.Unlock()

	if !wasExecuted {
		t.Error("callback was not executed spontaneously")
	}
}

func TestCallbackTaskManager_GetStats(t *testing.T) {
	manager := NewCallbackTaskManager()
	defer manager.Shutdown()

	// Register some callbacks
	futureID1 := uint64(1001)
	futureID2 := uint64(1002)
	futureID3 := uint64(1003)

	manager.RegisterCallback(futureID1, CallbackModeWaitAnyOnly, func() {})
	manager.RegisterCallback(futureID2, CallbackModeWaitAnyOnly, func() {})
	manager.RegisterCallback(futureID3, CallbackModeWaitAnyOnly, func() {})

	pending, completed, cancelled := manager.GetCallbackStats()
	if pending != 3 || completed != 0 || cancelled != 0 {
		t.Errorf("expected (3,0,0), got (%d,%d,%d)", pending, completed, cancelled)
	}

	// Complete one
	manager.CompleteCallback(futureID1)

	// Cancel one
	manager.CancelCallback(futureID2)

	pending, completed, cancelled = manager.GetCallbackStats()
	if pending != 1 || completed != 1 || cancelled != 1 {
		t.Errorf("expected (1,1,1), got (%d,%d,%d)", pending, completed, cancelled)
	}
}

func TestCallbackHelper_RequestAdapterCallback(t *testing.T) {
	manager := NewCallbackTaskManager()
	defer manager.Shutdown()

	helper := NewCallbackHelper(manager)

	var receivedStatus RequestAdapterStatus
	var receivedAdapter Adapter
	var receivedMessage string

	callback := RequestAdapterCallbackInfo{
		Mode: CallbackModeAllowProcessEvents,
		Callback: func(status RequestAdapterStatus, adapter Adapter, message string) {
			receivedStatus = status
			receivedAdapter = adapter
			receivedMessage = message
		},
	}

	futureID := uint64(2001)
	testAdapter := &adapter{refCount: 1}

	task := helper.RegisterRequestAdapterCallback(
		futureID,
		callback,
		RequestAdapterStatusSuccess,
		testAdapter,
		"test message",
	)

	if task == nil {
		t.Fatal("expected task to be registered")
	}

	// Complete and process
	manager.CompleteCallback(futureID)
	manager.ProcessEvents()
	time.Sleep(10 * time.Millisecond)

	if receivedStatus != RequestAdapterStatusSuccess {
		t.Errorf("expected status %d, got %d", RequestAdapterStatusSuccess, receivedStatus)
	}

	if receivedAdapter != testAdapter {
		t.Error("adapter mismatch")
	}

	if receivedMessage != "test message" {
		t.Errorf("expected 'test message', got '%s'", receivedMessage)
	}
}

func TestCallbackTaskManager_Shutdown(t *testing.T) {
	manager := NewCallbackTaskManager()

	// Register some callbacks
	manager.RegisterCallback(1, CallbackModeWaitAnyOnly, func() {})
	manager.RegisterCallback(2, CallbackModeWaitAnyOnly, func() {})

	count := manager.GetPendingCallbackCount()
	if count != 2 {
		t.Errorf("expected 2 pending callbacks, got %d", count)
	}

	// Shutdown
	err := manager.Shutdown()
	if err != nil {
		t.Errorf("shutdown failed: %v", err)
	}

	if !manager.IsShutdown() {
		t.Error("manager should be shutdown")
	}

	// Try to register after shutdown
	task := manager.RegisterCallback(3, CallbackModeWaitAnyOnly, func() {})
	if task != nil {
		t.Error("expected registration to fail after shutdown")
	}
}
