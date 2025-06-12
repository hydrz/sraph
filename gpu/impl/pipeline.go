package impl

import (
	"fmt"
	"sync"
	"sync/atomic"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ PipelineLayout = (*pipelineLayout)(nil)
var _ QuerySet = (*querySet)(nil)

// pipelineLayout implements the PipelineLayout interface
type pipelineLayout struct {
	mu               sync.RWMutex
	refCount         int32
	label            string
	bindGroupLayouts []BindGroupLayout
	destroyed        bool
}

// newPipelineLayout creates a new WebGPU pipeline layout
func newPipelineLayout(descriptor PipelineLayoutDescriptor) PipelineLayout {
	return &pipelineLayout{
		refCount:         1,
		label:            descriptor.Label,
		bindGroupLayouts: descriptor.BindGroupLayouts,
	}
}

// NewPipelineLayout creates a new WebGPU pipeline layout (public factory function)
func NewPipelineLayout(descriptor PipelineLayoutDescriptor) PipelineLayout {
	return newPipelineLayout(descriptor)
}

// SetLabel sets the pipeline layout label
func (pl *pipelineLayout) SetLabel(label string) error {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	if pl.destroyed {
		return fmt.Errorf("pipeline layout has been destroyed")
	}

	pl.label = label
	return nil
}

// AddRef increments the reference count
func (pl *pipelineLayout) AddRef() error {
	if atomic.LoadInt32(&pl.refCount) <= 0 {
		return fmt.Errorf("pipeline layout has been destroyed")
	}

	atomic.AddInt32(&pl.refCount, 1)
	return nil
}

// Release decrements the reference count and destroys if zero
func (pl *pipelineLayout) Release() error {
	newCount := atomic.AddInt32(&pl.refCount, -1)
	if newCount == 0 {
		pl.mu.Lock()
		pl.destroyed = true
		pl.mu.Unlock()
	} else if newCount < 0 {
		return fmt.Errorf("reference count cannot be negative")
	}
	return nil
}

// querySet implements the QuerySet interface
type querySet struct {
	mu        sync.RWMutex
	refCount  int32
	label     string
	qType     QueryType
	count     uint32
	destroyed bool
}

// NewQuerySet creates a new WebGPU query set
func NewQuerySet(descriptor QuerySetDescriptor) QuerySet {
	return &querySet{
		refCount: 1,
		label:    descriptor.Label,
		qType:    descriptor.Type,
		count:    descriptor.Count,
	}
}

// Destroy destroys the query set
func (qs *querySet) Destroy() error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	if qs.destroyed {
		return fmt.Errorf("query set has already been destroyed")
	}

	qs.destroyed = true
	return nil
}

// GetCount gets the query count
func (qs *querySet) GetCount() (uint32, error) {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	if qs.destroyed {
		return 0, fmt.Errorf("query set has been destroyed")
	}

	return qs.count, nil
}

// GetType gets the query type
func (qs *querySet) GetType() (QueryType, error) {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	if qs.destroyed {
		return 0, fmt.Errorf("query set has been destroyed")
	}

	return qs.qType, nil
}

// SetLabel sets the query set label
func (qs *querySet) SetLabel(label string) error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	if qs.destroyed {
		return fmt.Errorf("query set has been destroyed")
	}

	qs.label = label
	return nil
}

// AddRef increments the reference count
func (qs *querySet) AddRef() error {
	if atomic.LoadInt32(&qs.refCount) <= 0 {
		return fmt.Errorf("query set has been destroyed")
	}

	atomic.AddInt32(&qs.refCount, 1)
	return nil
}

// Release decrements the reference count and destroys if zero
func (qs *querySet) Release() error {
	newCount := atomic.AddInt32(&qs.refCount, -1)
	if newCount == 0 {
		qs.mu.Lock()
		qs.destroyed = true
		qs.mu.Unlock()
	} else if newCount < 0 {
		return fmt.Errorf("reference count cannot be negative")
	}
	return nil
}
