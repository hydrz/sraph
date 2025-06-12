package impl

import (
	"fmt"
	"sync"
	"sync/atomic"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ Instance = (*instance)(nil)

// instance implements the Instance interface
type instance struct {
	mu              sync.RWMutex
	refCount        int32
	features        []InstanceFeatureName
	limits          InstanceLimits
	destroyed       bool
	callbackManager *CallbackTaskManager
	callbackHelper  *CallbackHelper
}

// NewInstance creates a new WebGPU instance
func NewInstance(descriptor InstanceDescriptor) *instance {
	callbackManager := NewCallbackTaskManager()
	callbackHelper := NewCallbackHelper(callbackManager)

	instance := &instance{
		refCount:        1,
		features:        descriptor.RequiredFeatures,
		limits:          descriptor.RequiredLimits,
		callbackManager: callbackManager,
		callbackHelper:  callbackHelper,
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

	surface := &surface{
		refCount: 1,
		label:    descriptor.Label,
	}

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

// HasWGSLLanguageFeature checks if a WGSL language feature is supported
func (i *instance) HasWGSLLanguageFeature(feature WGSLLanguageFeatureName) (bool, error) {
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
func (i *instance) ProcessEvents() error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return fmt.Errorf("instance has been destroyed")
	}

	// Process callbacks that are queued for event processing
	return i.callbackManager.ProcessEvents()
}

// RequestAdapter requests a WebGPU adapter
func (i *instance) RequestAdapter(options RequestAdapterOptions, callback RequestAdapterCallbackInfo) Future {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Create a future for async adapter request
	future := Future{
		Id: generateFutureId(),
	}

	if i.destroyed {
		// Register callback with error status
		i.callbackHelper.RegisterRequestAdapterCallback(
			future.Id,
			callback,
			RequestAdapterStatusError,
			nil,
			"instance has been destroyed",
		)
		// Complete the callback immediately
		i.callbackManager.CompleteCallback(future.Id)
		return future
	}

	// Register the callback for async completion
	i.callbackHelper.RegisterRequestAdapterCallback(
		future.Id,
		callback,
		RequestAdapterStatusSuccess,
		nil, // Will be set when adapter is created
		"",
	)

	// In a real implementation, this would initiate async adapter creation
	// and call the callback when complete
	go func() {
		adapter := &adapter{
			refCount:    1,
			backendType: options.BackendType,
			adapterType: AdapterTypeDiscreteGPU, // Default to discrete GPU
			features:    getDefaultAdapterFeatures(),
			limits:      getDefaultAdapterLimits(),
		}

		// Update the callback with the actual adapter
		i.callbackHelper.RegisterRequestAdapterCallback(
			future.Id,
			callback,
			RequestAdapterStatusSuccess,
			adapter,
			"",
		)

		// Mark the callback as completed
		i.callbackManager.CompleteCallback(future.Id)
	}()

	return future
}

// WaitAny waits for any of the given futures to complete
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

	// Wait for the specific future to complete
	if i.callbackManager.WaitForCallback(futures.Future.Id, timeoutNS) {
		futures.Completed = true
		return WaitStatusSuccess, nil
	}

	// Check if it timed out or errored
	if timeoutNS > 0 {
		return WaitStatusTimedOut, nil
	}

	return WaitStatusError, fmt.Errorf("failed to wait for future")
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
		// Shutdown the callback manager
		if i.callbackManager != nil {
			i.callbackManager.Shutdown()
		}
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

// Thread-safe future ID generation
var futureIdCounter int64

func generateFutureId() uint64 {
	return uint64(atomic.AddInt64(&futureIdCounter, 1))
}
