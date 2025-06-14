package impl

import (
	"fmt"
	"sync"

	. "github.com/opensraph/sraph/gpu/wgpu"
)

var _ PipelineLayout = (*pipelineLayout)(nil)
var _ QuerySet = (*querySet)(nil)

// pipelineLayout implements the PipelineLayout interface
type pipelineLayout struct {
	mu               sync.RWMutex
	label            string
	bindGroupLayouts []BindGroupLayout
	destroyed        bool
}

// newPipelineLayout creates a new WebGPU pipeline layout
func newPipelineLayout(descriptor PipelineLayoutDescriptor) PipelineLayout {
	return &pipelineLayout{
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

// querySet implements the QuerySet interface
type querySet struct {
	mu        sync.RWMutex
	label     string
	qType     QueryType
	count     uint32
	destroyed bool
}

// NewQuerySet creates a new WebGPU query set
func NewQuerySet(descriptor QuerySetDescriptor) QuerySet {
	return &querySet{
		label: descriptor.Label,
		qType: descriptor.Type,
		count: descriptor.Count,
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

// Count gets the query count
func (qs *querySet) Count() (uint32, error) {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	if qs.destroyed {
		return 0, fmt.Errorf("query set has been destroyed")
	}

	return qs.count, nil
}

// Type gets the query type
func (qs *querySet) Type() (QueryType, error) {
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
