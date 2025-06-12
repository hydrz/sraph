package gpu

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

var _ ShaderModule = (*shaderModule)(nil)

// shaderModule implements the ShaderModule interface
type shaderModule struct {
	mu        sync.RWMutex
	refCount  int32
	label     string
	code      string
	destroyed bool
}

// newShaderModule creates a new WebGPU shader module
func newShaderModule(descriptor ShaderModuleDescriptor) ShaderModule {
	return &shaderModule{
		refCount: 1,
		label:    descriptor.Label,
	}
}

// GetCompilationInfo gets compilation information
func (sm *shaderModule) GetCompilationInfo(ctx context.Context) Future {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	future := Future{
		Id: generateFutureId(),
	}

	// In a real implementation, this would compile the shader and return info asynchronously
	go func() {
		// Mock compilation - in a real implementation this would compile WGSL/SPIR-V
		_ = "compilation complete"
	}()

	return future
}

// SetLabel sets the shader module label
func (sm *shaderModule) SetLabel(ctx context.Context, label string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.destroyed {
		return fmt.Errorf("shader module has been destroyed")
	}

	sm.label = label
	return nil
}

// AddRef increments the reference count
func (sm *shaderModule) AddRef(ctx context.Context) error {
	if atomic.LoadInt32(&sm.refCount) <= 0 {
		return fmt.Errorf("shader module has been destroyed")
	}

	atomic.AddInt32(&sm.refCount, 1)
	return nil
}

// Release decrements the reference count and destroys if zero
func (sm *shaderModule) Release(ctx context.Context) error {
	newCount := atomic.AddInt32(&sm.refCount, -1)
	if newCount == 0 {
		sm.mu.Lock()
		sm.destroyed = true
		sm.mu.Unlock()
	} else if newCount < 0 {
		return fmt.Errorf("reference count cannot be negative")
	}
	return nil
}
