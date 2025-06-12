package impl

import (
	"fmt"
	"sync"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ Adapter = (*adapter)(nil)

// adapter implements the Adapter interface
type adapter struct {
	mu          sync.RWMutex
	refCount    int32
	backendType BackendType
	adapterType AdapterType
	features    []FeatureName
	limits      Limits
	info        AdapterInfo
	destroyed   bool
}

// NewAdapter creates a new WebGPU adapter
func NewAdapter(backendType BackendType, adapterType AdapterType) Adapter {
	adapter := &adapter{
		refCount:    1,
		backendType: backendType,
		adapterType: adapterType,
		features: []FeatureName{
			FeatureNameDepthClipControl,
			FeatureNameTimestampQuery,
			FeatureNameTextureCompressionBC,
		},
		limits: getDefaultLimits(),
		info: AdapterInfo{
			Vendor:       "WebGPU Implementation",
			Architecture: "Unknown",
			Device:       "Generic Device",
			Description:  "WebGPU Adapter",
			BackendType:  backendType,
			AdapterType:  adapterType,
			VendorID:     0,
			DeviceID:     0,
		},
	}
	return adapter
}

// GetFeatures retrieves supported features
func (a *adapter) GetFeatures(features SupportedFeatures) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.destroyed {
		return fmt.Errorf("adapter has been destroyed")
	}

	features.Features = append(features.Features, a.features...)
	return nil
}

// GetInfo retrieves adapter information
func (a *adapter) GetInfo(info AdapterInfo) (Status, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.destroyed {
		return StatusError, fmt.Errorf("adapter has been destroyed")
	}

	info.Vendor = a.info.Vendor
	info.Architecture = a.info.Architecture
	info.Device = a.info.Device
	info.Description = a.info.Description
	info.BackendType = a.info.BackendType
	info.AdapterType = a.info.AdapterType
	info.VendorID = a.info.VendorID
	info.DeviceID = a.info.DeviceID

	return StatusSuccess, nil
}

// GetLimits retrieves adapter limits
func (a *adapter) GetLimits(limits Limits) (Status, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.destroyed {
		return StatusError, fmt.Errorf("adapter has been destroyed")
	}

	limits = a.limits
	return StatusSuccess, nil
}

// HasFeature checks if a feature is supported
func (a *adapter) HasFeature(feature FeatureName) (bool, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.destroyed {
		return false, fmt.Errorf("adapter has been destroyed")
	}

	for _, f := range a.features {
		if f == feature {
			return true, nil
		}
	}
	return false, nil
}

// RequestDevice requests a WebGPU device
func (a *adapter) RequestDevice(descriptor DeviceDescriptor, callback RequestDeviceCallbackInfo) Future {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Create a future for async device request
	future := Future{
		Id: generateFutureId(),
	}

	if a.destroyed {
		// Call callback with error if adapter is destroyed
		go func() {
			if callback.Callback != nil {
				callback.Callback(RequestDeviceStatusError, nil, "adapter has been destroyed")
			}
		}()
		return future
	}

	// In a real implementation, this would initiate async device creation
	go func() {
		device := &device{
			refCount:  1,
			label:     descriptor.Label,
			features:  descriptor.RequiredFeatures,
			limits:    descriptor.RequiredLimits,
			queue:     newQueue(descriptor.DefaultQueue),
			destroyed: false,
		}

		// Call callback with success
		if callback.Callback != nil {
			callback.Callback(RequestDeviceStatusSuccess, device, "")
		}
	}()

	return future
}

// AddRef increments the reference count
func (a *adapter) AddRef() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.destroyed {
		return fmt.Errorf("adapter has been destroyed")
	}

	a.refCount++
	return nil
}

// Release decrements the reference count and destroys if zero
func (a *adapter) Release() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.destroyed {
		return fmt.Errorf("adapter has been destroyed")
	}

	a.refCount--
	if a.refCount <= 0 {
		a.destroyed = true
	}

	return nil
}

// getDefaultLimits returns default WebGPU limits
func getDefaultLimits() Limits {
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
		MaxImmediateSize:                          1024,
	}
}
