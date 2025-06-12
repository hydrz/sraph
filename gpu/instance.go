package gpu

import (
	"context"
	"fmt"
	"sync"
)

var _ Instance = (*instance)(nil)

// instance implements the Instance interface
type instance struct {
	mu        sync.RWMutex
	refCount  int32
	features  []InstanceFeatureName
	limits    InstanceLimits
	destroyed bool
}

// newInstance creates a new WebGPU instance
func newInstance(descriptor InstanceDescriptor) *instance {
	instance := &instance{
		refCount: 1,
		features: descriptor.RequiredFeatures,
		limits:   descriptor.RequiredLimits,
	}
	return instance
}

// CreateSurface creates a new surface for rendering
func (i *instance) CreateSurface(ctx context.Context, descriptor SurfaceDescriptor) (Surface, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return nil, fmt.Errorf("instance has been destroyed")
	}

	surface := &surface{
		refCount: 1,
		label:    descriptor.Label,
	}

	return surface, nil
}

// GetFeatures retrieves supported features
func (i *instance) GetFeatures(ctx context.Context, features SupportedInstanceFeatures) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return fmt.Errorf("instance has been destroyed")
	}

	features.Features = append(features.Features, i.features...)
	return nil
}

// GetLimits retrieves instance limits
func (i *instance) GetLimits(ctx context.Context, limits InstanceLimits) (Status, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return StatusError, fmt.Errorf("instance has been destroyed")
	}

	limits.TimedWaitAnyMaxCount = i.limits.TimedWaitAnyMaxCount
	return StatusSuccess, nil
}

// GetWGSLLanguageFeatures retrieves supported WGSL language features
func (i *instance) GetWGSLLanguageFeatures(ctx context.Context, features SupportedWGSLLanguageFeatures) (Status, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return StatusError, fmt.Errorf("instance has been destroyed")
	}

	// Add default WGSL language features
	features.Features = []WGSLLanguageFeatureName{
		WGSLLanguageFeatureNameReadonlyAndReadwriteStorageTextures,
		WGSLLanguageFeatureNamePacked4x8IntegerDotProduct,
	}

	return StatusSuccess, nil
}

// HasFeature checks if an instance feature is supported
func (i *instance) HasFeature(ctx context.Context, feature InstanceFeatureName) (bool, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return false, fmt.Errorf("instance has been destroyed")
	}

	for _, f := range i.features {
		if f == feature {
			return true, nil
		}
	}
	return false, nil
}

// HasWGSLLanguageFeature checks if a WGSL language feature is supported
func (i *instance) HasWGSLLanguageFeature(ctx context.Context, feature WGSLLanguageFeatureName) (bool, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return false, fmt.Errorf("instance has been destroyed")
	}

	// Check supported WGSL language features
	supportedFeatures := []WGSLLanguageFeatureName{
		WGSLLanguageFeatureNameReadonlyAndReadwriteStorageTextures,
		WGSLLanguageFeatureNamePacked4x8IntegerDotProduct,
	}

	for _, f := range supportedFeatures {
		if f == feature {
			return true, nil
		}
	}
	return false, nil
}

// ProcessEvents processes pending events
func (i *instance) ProcessEvents(ctx context.Context) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return fmt.Errorf("instance has been destroyed")
	}

	// Process any pending asynchronous operations
	// In a real implementation, this would handle callbacks and async operations
	return nil
}

// RequestAdapter requests a WebGPU adapter
func (i *instance) RequestAdapter(ctx context.Context, options RequestAdapterOptions) Future {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Create a future for async adapter request
	future := Future{
		Id: generateFutureId(),
	}

	// In a real implementation, this would initiate async adapter creation
	// and call the callback when complete
	go func() {
		adapter := &adapter{
			refCount:    1,
			backendType: options.BackendType,
			adapterType: AdapterTypeDiscreteGPU, // Default to discrete GPU
		}

		// Simulate async completion
		// callback would be called here with the adapter
		_ = adapter
	}()

	return future
}

// WaitAny waits for any of the given futures to complete
func (i *instance) WaitAny(ctx context.Context, futureCount uintptr, futures FutureWaitInfo, timeoutNS uint64) (WaitStatus, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return WaitStatusError, fmt.Errorf("instance has been destroyed")
	}

	// In a real implementation, this would wait for futures to complete
	// For now, return success immediately
	return WaitStatusSuccess, nil
}

// AddRef increments the reference count
func (i *instance) AddRef(ctx context.Context) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.destroyed {
		return fmt.Errorf("instance has been destroyed")
	}

	i.refCount++
	return nil
}

// Release decrements the reference count and destroys if zero
func (i *instance) Release(ctx context.Context) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.destroyed {
		return fmt.Errorf("instance has been destroyed")
	}

	i.refCount--
	if i.refCount <= 0 {
		i.destroyed = true
	}

	return nil
}

var futureIdCounter uint64 = 1

func generateFutureId() uint64 {
	// In a real implementation, this would be thread-safe
	futureIdCounter++
	return futureIdCounter
}
