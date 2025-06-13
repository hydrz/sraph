package impl

import (
	"fmt"
	"sync"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ ShaderModule = (*shaderModule)(nil)

// shaderModule implements the ShaderModule interface
type shaderModule struct {
	mu        sync.RWMutex
	label     string
	code      string
	destroyed bool
}

// newShaderModule creates a new WebGPU shader module
func newShaderModule(descriptor ShaderModuleDescriptor) ShaderModule {
	return &shaderModule{
		label: descriptor.Label,
	}
}

// NewShaderModule creates a new WebGPU shader module (public factory function)
func NewShaderModule(descriptor ShaderModuleDescriptor) ShaderModule {
	return newShaderModule(descriptor)
}

// GetCompilationInfo gets compilation information
func (sm *shaderModule) GetCompilationInfo(callback CompilationInfoCallbackInfo) Future {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	future := Future{
		Id: GenerateFutureId(),
	}

	if sm.destroyed {
		// Register error callback
		GlobalCallbackRegistry().CompilationInfo(future.Id, callback, CompilationInfoRequestStatusCallbackCancelled, CompilationInfo{})
		GlobalCallbackManager().Complete(future.Id)
		return future
	}

	// Simulate async compilation
	go func() {
		// Mock compilation - in a real implementation this would compile WGSL/SPIR-V
		compilationInfo := CompilationInfo{
			Messages: []CompilationMessage{},
		}

		// Register success callback
		GlobalCallbackRegistry().CompilationInfo(future.Id, callback, CompilationInfoRequestStatusSuccess, compilationInfo)
		GlobalCallbackManager().Complete(future.Id)
	}()

	return future
}

// SetLabel sets the shader module label
func (sm *shaderModule) SetLabel(label string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.destroyed {
		return fmt.Errorf("shader module has been destroyed")
	}

	sm.label = label
	return nil
}
