package gpu

import (
	"context"
	"fmt"
	"sync"
)

var _ BindGroup = (*bindGroup)(nil)
var _ BindGroupLayout = (*bindGroupLayout)(nil)

// bindGroup implements the BindGroup interface
type bindGroup struct {
	mu        sync.RWMutex
	refCount  int32
	label     string
	layout    BindGroupLayout
	entries   []BindGroupEntry
	destroyed bool
}

// SetLabel sets the bind group label
func (bg *bindGroup) SetLabel(ctx context.Context, label string) error {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	if bg.destroyed {
		return fmt.Errorf("bind group has been destroyed")
	}

	bg.label = label
	return nil
}

// AddRef increments the reference count
func (bg *bindGroup) AddRef(ctx context.Context) error {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	if bg.destroyed {
		return fmt.Errorf("bind group has been destroyed")
	}

	bg.refCount++
	return nil
}

// Release decrements the reference count and destroys if zero
func (bg *bindGroup) Release(ctx context.Context) error {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	if bg.destroyed {
		return fmt.Errorf("bind group has been destroyed")
	}

	bg.refCount--
	if bg.refCount <= 0 {
		bg.destroyed = true
	}

	return nil
}

// bindGroupLayout implements the BindGroupLayout interface
type bindGroupLayout struct {
	mu        sync.RWMutex
	refCount  int32
	label     string
	entries   []BindGroupLayoutEntry
	destroyed bool
}

// SetLabel sets the bind group layout label
func (bgl *bindGroupLayout) SetLabel(ctx context.Context, label string) error {
	bgl.mu.Lock()
	defer bgl.mu.Unlock()

	if bgl.destroyed {
		return fmt.Errorf("bind group layout has been destroyed")
	}

	bgl.label = label
	return nil
}

// AddRef increments the reference count
func (bgl *bindGroupLayout) AddRef(ctx context.Context) error {
	bgl.mu.Lock()
	defer bgl.mu.Unlock()

	if bgl.destroyed {
		return fmt.Errorf("bind group layout has been destroyed")
	}

	bgl.refCount++
	return nil
}

// Release decrements the reference count and destroys if zero
func (bgl *bindGroupLayout) Release(ctx context.Context) error {
	bgl.mu.Lock()
	defer bgl.mu.Unlock()

	if bgl.destroyed {
		return fmt.Errorf("bind group layout has been destroyed")
	}

	bgl.refCount--
	if bgl.refCount <= 0 {
		bgl.destroyed = true
	}

	return nil
}
