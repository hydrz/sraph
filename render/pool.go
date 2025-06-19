package render

import (
	"sync"
	"time"
)

// Pool represents a resource pool for managing GPU resources efficiently
// Based on Impeller's Pool template class design
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

	// Trim removes unused resources beyond the trim size
	Trim(trimSize int)

	// GetHitRate returns the cache hit rate (0.0 to 1.0)
	GetHitRate() float64

	// GetStats returns pool statistics
	GetStats() PoolStats
}

// PoolStats contains statistics about pool usage
type PoolStats struct {
	// Current number of resources in pool
	Size int

	// Maximum pool size
	MaxSize int

	// Total number of Get() calls
	TotalGets int

	// Number of cache hits (reused resources)
	CacheHits int

	// Number of cache misses (new resources created)
	CacheMisses int

	// Cache hit rate (0.0 to 1.0)
	HitRate float64

	// Total resources created
	TotalCreated int

	// Total resources destroyed
	TotalDestroyed int
}

// CreateFunc is a function type for creating new resources
type CreateFunc[T any] func() T

// ResetFunc is a function type for resetting resources before reuse
type ResetFunc[T any] func(resource T)

// DestroyFunc is a function type for destroying resources
type DestroyFunc[T any] func(resource T)

// ValidateFunc is a function type for validating if a resource is still usable
type ValidateFunc[T any] func(resource T) bool

// PoolEntry wraps a resource with metadata
type PoolEntry[T any] struct {
	Resource T
	LastUsed time.Time
	UseCount int
	IsValid  bool
}

// PoolImpl is a generic implementation of Pool based on Impeller's design
type PoolImpl[T any] struct {
	// Thread safety
	mutex sync.RWMutex

	// Resource storage
	resources []PoolEntry[T]

	// Configuration
	maxSize  int
	trimSize int

	// Function handlers
	createFn   CreateFunc[T]
	resetFn    ResetFunc[T]
	destroyFn  DestroyFunc[T]
	validateFn ValidateFunc[T]

	// Statistics
	stats PoolStats

	// Resource lifetime management
	maxIdleTime time.Duration
}

// NewPool creates a new resource pool with the specified functions
func NewPool[T any](createFn CreateFunc[T], resetFn ResetFunc[T]) Pool[T] {
	return &PoolImpl[T]{
		resources:   make([]PoolEntry[T], 0),
		createFn:    createFn,
		resetFn:     resetFn,
		maxSize:     100, // Default max size
		trimSize:    50,  // Default trim size
		maxIdleTime: 5 * time.Minute,
		stats: PoolStats{
			MaxSize: 100,
		},
	}
}

// NewPoolWithOptions creates a new pool with custom options
func NewPoolWithOptions[T any](
	createFn CreateFunc[T],
	resetFn ResetFunc[T],
	destroyFn DestroyFunc[T],
	validateFn ValidateFunc[T],
	maxSize int,
	maxIdleTime time.Duration,
) Pool[T] {
	return &PoolImpl[T]{
		resources:   make([]PoolEntry[T], 0, maxSize),
		createFn:    createFn,
		resetFn:     resetFn,
		destroyFn:   destroyFn,
		validateFn:  validateFn,
		maxSize:     maxSize,
		trimSize:    maxSize / 2,
		maxIdleTime: maxIdleTime,
		stats: PoolStats{
			MaxSize: maxSize,
		},
	}
}

// Get retrieves a resource from the pool or creates a new one
func (p *PoolImpl[T]) Get() T {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.stats.TotalGets++

	// Try to find a valid resource in the pool
	for i := len(p.resources) - 1; i >= 0; i-- {
		entry := &p.resources[i]

		// Check if resource is still valid
		if p.validateFn != nil && !p.validateFn(entry.Resource) {
			// Remove invalid resource
			p.removeResourceAt(i)
			continue
		}

		// Check if resource has expired
		if p.maxIdleTime > 0 && time.Since(entry.LastUsed) > p.maxIdleTime {
			p.removeResourceAt(i)
			continue
		}

		// Found a valid resource
		resource := entry.Resource
		entry.UseCount++
		entry.LastUsed = time.Now()

		// Move to end (most recently used)
		p.resources = append(p.resources[:i], p.resources[i+1:]...)

		// Reset resource before reuse
		if p.resetFn != nil {
			p.resetFn(resource)
		}

		p.stats.CacheHits++
		p.updateHitRate()
		return resource
	}

	// No valid resource found, create a new one
	resource := p.createFn()
	p.stats.CacheMisses++
	p.stats.TotalCreated++
	p.updateHitRate()

	return resource
}

// Put returns a resource to the pool for reuse
func (p *PoolImpl[T]) Put(resource T) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Don't add if pool is at max capacity
	if len(p.resources) >= p.maxSize {
		if p.destroyFn != nil {
			p.destroyFn(resource)
			p.stats.TotalDestroyed++
		}
		return
	}

	// Validate resource before adding to pool
	if p.validateFn != nil && !p.validateFn(resource) {
		if p.destroyFn != nil {
			p.destroyFn(resource)
			p.stats.TotalDestroyed++
		}
		return
	}

	// Add resource to pool
	entry := PoolEntry[T]{
		Resource: resource,
		LastUsed: time.Now(),
		UseCount: 0,
		IsValid:  true,
	}

	p.resources = append(p.resources, entry)
	p.stats.Size = len(p.resources)
}

// Clear removes all resources from the pool
func (p *PoolImpl[T]) Clear() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.destroyFn != nil {
		for _, entry := range p.resources {
			p.destroyFn(entry.Resource)
			p.stats.TotalDestroyed++
		}
	}

	p.resources = p.resources[:0]
	p.stats.Size = 0
}

// Size returns the current size of the pool
func (p *PoolImpl[T]) Size() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return len(p.resources)
}

// SetMaxSize sets the maximum size of the pool
func (p *PoolImpl[T]) SetMaxSize(maxSize int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.maxSize = maxSize
	p.stats.MaxSize = maxSize

	// Trim if current size exceeds new max size
	if len(p.resources) > maxSize {
		p.trimToSize(maxSize)
	}
}

// GetMaxSize returns the maximum size of the pool
func (p *PoolImpl[T]) GetMaxSize() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.maxSize
}

// Trim removes unused resources beyond the trim size
func (p *PoolImpl[T]) Trim(trimSize int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.trimToSize(trimSize)
}

// GetHitRate returns the cache hit rate
func (p *PoolImpl[T]) GetHitRate() float64 {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.stats.HitRate
}

// GetStats returns pool statistics
func (p *PoolImpl[T]) GetStats() PoolStats {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	statsCopy := p.stats
	statsCopy.Size = len(p.resources)
	return statsCopy
}

// Internal helper methods

func (p *PoolImpl[T]) removeResourceAt(index int) {
	if p.destroyFn != nil {
		p.destroyFn(p.resources[index].Resource)
		p.stats.TotalDestroyed++
	}
	p.resources = append(p.resources[:index], p.resources[index+1:]...)
	p.stats.Size = len(p.resources)
}

func (p *PoolImpl[T]) trimToSize(size int) {
	if len(p.resources) <= size {
		return
	}

	// Remove oldest resources first
	toRemove := len(p.resources) - size
	for i := 0; i < toRemove; i++ {
		if p.destroyFn != nil {
			p.destroyFn(p.resources[i].Resource)
			p.stats.TotalDestroyed++
		}
	}

	p.resources = p.resources[toRemove:]
	p.stats.Size = len(p.resources)
}

func (p *PoolImpl[T]) updateHitRate() {
	if p.stats.TotalGets > 0 {
		p.stats.HitRate = float64(p.stats.CacheHits) / float64(p.stats.TotalGets)
	} else {
		p.stats.HitRate = 0.0
	}
}

// CleanupExpiredResources removes resources that have exceeded max idle time
func (p *PoolImpl[T]) CleanupExpiredResources() {
	if p.maxIdleTime <= 0 {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	now := time.Now()
	i := 0
	for i < len(p.resources) {
		if now.Sub(p.resources[i].LastUsed) > p.maxIdleTime {
			p.removeResourceAt(i)
		} else {
			i++
		}
	}
}

// SetDestroyFunc sets the destroy function for resources
func (p *PoolImpl[T]) SetDestroyFunc(destroyFn DestroyFunc[T]) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.destroyFn = destroyFn
}

// SetValidateFunc sets the validation function for resources
func (p *PoolImpl[T]) SetValidateFunc(validateFn ValidateFunc[T]) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.validateFn = validateFn
}

// SetMaxIdleTime sets the maximum idle time for resources
func (p *PoolImpl[T]) SetMaxIdleTime(maxIdleTime time.Duration) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.maxIdleTime = maxIdleTime
}
