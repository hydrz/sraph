package impl

import (
	"fmt"
	"sync"
	"time"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ Adapter = (*adapter)(nil)

// adapter implements the Adapter interface
type adapter struct {
	mu          sync.RWMutex
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

// RequestDevice requests a WebGPU device with improved callback handling
func (a *adapter) RequestDevice(descriptor DeviceDescriptor, callback RequestDeviceCallbackInfo) Future {
	a.mu.RLock()
	defer a.mu.RUnlock()

	future := Future{
		Id: GenerateFutureId(),
	}

	if a.destroyed {
		// Use global callback registry for error
		GlobalCallbackRegistry().RequestDevice(future.Id, callback, RequestDeviceStatusError, nil, "adapter has been destroyed")
		GlobalCallbackManager().Complete(future.Id)
		return future
	}

	// Start async device creation
	go func() {
		// Simulate device creation process
		time.Sleep(time.Millisecond * 10) // Simulate work

		device := NewDevice(descriptor)

		// Register success callback
		GlobalCallbackRegistry().RequestDevice(future.Id, callback, RequestDeviceStatusSuccess, device, "")
		GlobalCallbackManager().Complete(future.Id)
	}()

	return future
}
