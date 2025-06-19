package render

// Pool represents a resource pool for managing GPU resources efficiently
type Pool[T any] interface {
	// Get retrieves a resource from the pool or creates a new one
	Get() T
	
	// Put returns a resource to the pool for reuse
	Put(resource T)
	
	// Clear removes all resources from the pool
	Clear()
	
	// Size returns the current size of the pool
	Size() int
	
	// SetMaxSize sets the maximum size of the pool
	SetMaxSize(maxSize int)
	
	// GetMaxSize returns the maximum size of the pool
	GetMaxSize() int
}

// CreateFunc is a function type for creating new resources
type CreateFunc[T any] func() T

// ResetFunc is a function type for resetting resources before reuse
type ResetFunc[T any] func(resource T)

// PoolImpl is a generic implementation of Pool
type PoolImpl[T any] struct {
	resources []T
	createFn  CreateFunc[T]
	resetFn   ResetFunc[T]
	maxSize   int
}

// NewPool creates a new resource pool with the specified create and reset functions
func NewPool[T any](createFn CreateFunc[T], resetFn ResetFunc[T]) Pool[T] {
	return &PoolImpl[T]{
		resources: make([]T, 0),
		createFn:  createFn,
		resetFn:   resetFn,
		maxSize:   100, // Default max size
	}
}

// Get retrieves a resource from the pool or creates a new one
func (p *PoolImpl[T]) Get() T {
	if len(p.resources) > 0 {
		// Pop from the end for better performance
		resource := p.resources[len(p.resources)-1]
		p.resources = p.resources[:len(p.resources)-1]
		return resource
	}
	
	// Create new resource if pool is empty
	return p.createFn()
}

// Put returns a resource to the pool for reuse
func (p *PoolImpl[T]) Put(resource T) {
	if len(p.resources) >= p.maxSize {
		// Pool is full, discard the resource
		return
	}
	
	// Reset the resource if reset function is provided
	if p.resetFn != nil {
		p.resetFn(resource)
	}
	
	// Add to pool
	p.resources = append(p.resources, resource)
}

// Clear removes all resources from the pool
func (p *PoolImpl[T]) Clear() {
	p.resources = p.resources[:0] // Clear slice but keep capacity
}

// Size returns the current size of the pool
func (p *PoolImpl[T]) Size() int {
	return len(p.resources)
}

// SetMaxSize sets the maximum size of the pool
func (p *PoolImpl[T]) SetMaxSize(maxSize int) {
	p.maxSize = maxSize
	
	// Trim pool if it exceeds new max size
	if len(p.resources) > maxSize {
		p.resources = p.resources[:maxSize]
	}
}

// GetMaxSize returns the maximum size of the pool
func (p *PoolImpl[T]) GetMaxSize() int {
	return p.maxSize
}

// BufferPool is a specialized pool for buffers
type BufferPool interface {
	Pool[Buffer]
	
	// GetBuffer retrieves a buffer with the specified size and usage
	GetBuffer(size int, usage BufferUsage) Buffer
}

// BufferPoolImpl implements BufferPool
type BufferPoolImpl struct {
	*PoolImpl[Buffer]
	bufferSize  int
	bufferUsage BufferUsage
}

// NewBufferPool creates a new buffer pool
func NewBufferPool(bufferSize int, usage BufferUsage) BufferPool {
	createFn := func() Buffer {
		return NewBuffer(bufferSize, usage)
	}
	
	resetFn := func(buffer Buffer) {
		// TODO: Reset buffer state if needed
	}
	
	return &BufferPoolImpl{
		PoolImpl:    NewPool(createFn, resetFn).(*PoolImpl[Buffer]),
		bufferSize:  bufferSize,
		bufferUsage: usage,
	}
}

// GetBuffer retrieves a buffer with the specified size and usage
func (bp *BufferPoolImpl) GetBuffer(size int, usage BufferUsage) Buffer {
	if size == bp.bufferSize && usage == bp.bufferUsage {
		return bp.Get()
	}
	
	// Create new buffer with different specs
	return NewBuffer(size, usage)
}
