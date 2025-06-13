package impl

import (
	"fmt"
	"sync"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ BindGroup = (*bindGroup)(nil)
var _ BindGroupLayout = (*bindGroupLayout)(nil)

// bindGroup implements the BindGroup interface
type bindGroup struct {
	mu        sync.RWMutex
	label     string
	layout    BindGroupLayout
	entries   []BindGroupEntry
	destroyed bool
}

// SetLabel sets the bind group label
func (bg *bindGroup) SetLabel(label string) error {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	if bg.destroyed {
		return fmt.Errorf("bind group has been destroyed")
	}

	bg.label = label
	return nil
}

// NewBindGroup creates a new WebGPU bind group (public factory function)
func NewBindGroup(descriptor BindGroupDescriptor) BindGroup {
	return &bindGroup{
		label:   descriptor.Label,
		layout:  descriptor.Layout,
		entries: descriptor.Entries,
	}
}

// bindGroupLayout implements the BindGroupLayout interface
type bindGroupLayout struct {
	mu        sync.RWMutex
	label     string
	entries   []BindGroupLayoutEntry
	destroyed bool
}

// SetLabel sets the bind group layout label
func (bgl *bindGroupLayout) SetLabel(label string) error {
	bgl.mu.Lock()
	defer bgl.mu.Unlock()

	if bgl.destroyed {
		return fmt.Errorf("bind group layout has been destroyed")
	}

	bgl.label = label
	return nil
}

// NewBindGroupLayout creates a new WebGPU bind group layout (public factory function)
func NewBindGroupLayout(descriptor BindGroupLayoutDescriptor) BindGroupLayout {
	return &bindGroupLayout{
		label:   descriptor.Label,
		entries: descriptor.Entries,
	}
}
