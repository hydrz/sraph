package impl

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ Instance = (*instance)(nil)

// instance implements the Instance interface
type instance struct {
	mu               sync.RWMutex
	refCount         int32
	features         []InstanceFeatureName
	limits           InstanceLimits
	destroyed        bool
	callbackManager  *CallbackManager
	callbackRegistry *CallbackRegistry
}

// NewInstance creates a new WebGPU instance
func NewInstance(descriptor InstanceDescriptor) *instance {
	callbackManager := NewCallbackManager()
	instance := &instance{
		refCount:         1,
		features:         descriptor.RequiredFeatures,
		limits:           descriptor.RequiredLimits,
		callbackManager:  callbackManager,
		callbackRegistry: newCallbackRegistry(callbackManager),
	}
	return instance
}

// CreateSurface creates a new surface for rendering
func (i *instance) CreateSurface(descriptor SurfaceDescriptor) (Surface, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return nil, fmt.Errorf("instance has been destroyed")
	}

	surface := newSurface(descriptor.Label)
	return surface, nil
}

// GetWGSLLanguageFeatures retrieves supported WGSL language features
func (i *instance) GetWGSLLanguageFeatures(features SupportedWGSLLanguageFeatures) (Status, error) {
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

// GetInstanceFeatures retrieves supported instance features
func (i *instance) GetInstanceFeatures(features SupportedInstanceFeatures) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return fmt.Errorf("instance has been destroyed")
	}

	features.Features = append(features.Features, i.features...)
	return nil
}

// GetInstanceLimits retrieves instance limits
func (i *instance) GetInstanceLimits(limits InstanceLimits) (Status, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return StatusError, fmt.Errorf("instance has been destroyed")
	}

	limits = i.limits
	return StatusSuccess, nil
}

// HasInstanceFeature checks if an instance feature is supported (public method)
func (i *instance) HasInstanceFeature(feature InstanceFeatureName) (bool, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return false, fmt.Errorf("instance has been destroyed")
	}

	return i.hasInstanceFeature(feature), nil
}

// HasWGSLLanguageFeature checks if a WGSL language feature is supported
func (i *instance) HasWGSLLanguageFeature(feature WGSLLanguageFeatureName) (bool, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return false, fmt.Errorf("instance has been destroyed")
	}

	// In a real implementation, this would check against supported WGSL features
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

// ProcessEvents processes pending events and callbacks with improved synchronization
func (i *instance) ProcessEvents() error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return fmt.Errorf("instance has been destroyed")
	}

	// Process callbacks that are ready for event processing
	i.callbackManager.ProcessEvents()
	return nil
}

// RequestAdapter requests a WebGPU adapter with improved callback handling
func (i *instance) RequestAdapter(options RequestAdapterOptions, callback RequestAdapterCallbackInfo) Future {
	i.mu.RLock()
	defer i.mu.RUnlock()

	future := Future{
		Id: GenerateFutureId(),
	}

	if i.destroyed {
		// Register error callback
		i.callbackRegistry.RequestAdapter(future.Id, callback, RequestAdapterStatusError, nil, "instance has been destroyed")
		i.callbackManager.Complete(future.Id)
		return future
	}

	// Start async adapter creation
	go func() {
		// Simulate adapter creation process
		time.Sleep(time.Millisecond * 10) // Simulate work

		// Create adapter based on options
		var backendType BackendType = BackendTypeVulkan
		if options.BackendType != BackendTypeUndefined {
			backendType = options.BackendType
		}

		adapter := NewAdapter(backendType, AdapterTypeDiscreteGPU)

		// Register success callback
		i.callbackRegistry.RequestAdapter(future.Id, callback, RequestAdapterStatusSuccess, adapter, "")
		i.callbackManager.Complete(future.Id)
	}()

	return future
}

// WaitAny waits for any of the given futures to complete with improved implementation
func (i *instance) WaitAny(futureCount uintptr, futures FutureWaitInfo, timeoutNS uint64) (WaitStatus, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return WaitStatusError, fmt.Errorf("instance has been destroyed")
	}

	if futureCount == 0 {
		return WaitStatusError, fmt.Errorf("no futures to wait for")
	}

	// Check if timed wait is supported and timeout is specified
	if timeoutNS > 0 && !i.hasInstanceFeature(InstanceFeatureNameTimedWaitAnyEnable) {
		return WaitStatusError, fmt.Errorf("timed wait is not enabled")
	}

	timeout := time.Duration(timeoutNS) * time.Nanosecond

	// Use the improved callback manager for waiting
	if futures.Future.Id > 0 {
		success := i.callbackManager.Wait(futures.Future.Id, timeout)
		if success {
			futures.Completed = true
			return WaitStatusSuccess, nil
		} else {
			return WaitStatusTimedOut, nil
		}
	}

	return WaitStatusError, fmt.Errorf("invalid future")
}

// AddRef increments the reference count
func (i *instance) AddRef() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.destroyed {
		return fmt.Errorf("instance has been destroyed")
	}

	atomic.AddInt32(&i.refCount, 1)
	return nil
}

// Release decrements the reference count and destroys if zero
func (i *instance) Release() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.destroyed {
		return fmt.Errorf("instance has been destroyed")
	}

	if atomic.AddInt32(&i.refCount, -1) <= 0 {
		i.destroyed = true
		// Shutdown callback manager when instance is destroyed
		i.callbackManager.Shutdown()
	}

	return nil
}

// hasInstanceFeature checks if an instance feature is supported internally
func (i *instance) hasInstanceFeature(feature InstanceFeatureName) bool {
	for _, f := range i.features {
		if f == feature {
			return true
		}
	}
	return false
}

// Helper functions for default configurations
func getDefaultAdapterFeatures() []FeatureName {
	return []FeatureName{
		FeatureNameDepthClipControl,
		FeatureNameTimestampQuery,
		FeatureNameTextureCompressionBC,
	}
}

func getDefaultAdapterLimits() Limits {
	return Limits{
		MaxTextureDimension1D:                     8192,
		MaxTextureDimension2D:                     8192,
		MaxTextureDimension3D:                     2048,
		MaxTextureArrayLayers:                     256,
		MaxBindGroups:                             4,
		MaxBindGroupsPlusVertexBuffers:            24,
		MaxBindingsPerBindGroup:                   1000,
		MaxDynamicUniformBuffersPerPipelineLayout: 8,
		MaxDynamicStorageBuffersPerPipelineLayout: 4,
		MaxSampledTexturesPerShaderStage:          16,
		MaxSamplersPerShaderStage:                 16,
		MaxStorageBuffersPerShaderStage:           8,
		MaxStorageTexturesPerShaderStage:          4,
		MaxUniformBuffersPerShaderStage:           12,
		MaxUniformBufferBindingSize:               65536,
		MaxStorageBufferBindingSize:               134217728,
		MinUniformBufferOffsetAlignment:           256,
		MinStorageBufferOffsetAlignment:           256,
		MaxVertexBuffers:                          8,
		MaxBufferSize:                             268435456,
		MaxVertexAttributes:                       16,
		MaxVertexBufferArrayStride:                2048,
		MaxInterStageShaderVariables:              16,
		MaxColorAttachments:                       8,
		MaxColorAttachmentBytesPerSample:          32,
		MaxComputeWorkgroupStorageSize:            16384,
		MaxComputeInvocationsPerWorkgroup:         256,
		MaxComputeWorkgroupSizeX:                  256,
		MaxComputeWorkgroupSizeY:                  256,
		MaxComputeWorkgroupSizeZ:                  64,
		MaxComputeWorkgroupsPerDimension:          65535,
		MaxImmediateSize:                          16777216,
	}
}
