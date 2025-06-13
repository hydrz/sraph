package impl

import (
	"fmt"
	"sync"

	. "github.com/opensraph/sraph/gpu/wgpu"
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

// CompilationInfo implements ShaderModule.CompilationInfo.
// Returns nil as a stub, since async callback is not part of the interface.
func (sm *shaderModule) CompilationInfo() error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.destroyed {
		return fmt.Errorf("shader module has been destroyed")
	}

	// TODO: Implement actual compilation info retrieval if needed.
	return nil
}

// SetLabel implements ShaderModule.SetLabel.
func (sm *shaderModule) SetLabel(label string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.destroyed {
		return fmt.Errorf("shader module has been destroyed")
	}

	sm.label = label
	return nil
}
