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

// Features implements Adapter.Features.
func (a *adapter) Features() (*SupportedFeatures, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.destroyed {
		return nil, fmt.Errorf("adapter has been destroyed")
	}

	return &SupportedFeatures{Features: append([]FeatureName{}, a.features...)}, nil
}

// Info implements Adapter.Info.
func (a *adapter) Info() (*AdapterInfo, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.destroyed {
		return nil, fmt.Errorf("adapter has been destroyed")
	}

	info := a.info
	return &info, nil
}

// Limits implements Adapter.Limits.
func (a *adapter) Limits() (*Limits, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.destroyed {
		return nil, fmt.Errorf("adapter has been destroyed")
	}

	limits := a.limits
	return &limits, nil
}

// HasFeature implements Adapter.HasFeature.
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

// RequestDevice implements Adapter.RequestDevice.
func (a *adapter) RequestDevice(descriptor DeviceDescriptor) (Device, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.destroyed {
		return nil, fmt.Errorf("adapter has been destroyed")
	}

	return NewDevice(descriptor), nil
}
