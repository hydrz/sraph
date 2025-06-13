package impl

import (
	"fmt"
	"sync"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ Instance = (*instance)(nil)

// instance implements the Instance interface
type instance struct {
	mu        sync.RWMutex
	features  []InstanceFeatureName
	limits    InstanceLimits
	destroyed bool
}

// NewInstance creates a new WebGPU instance
func NewInstance(descriptor InstanceDescriptor) Instance {
	return &instance{
		features:  descriptor.RequiredFeatures,
		limits:    descriptor.RequiredLimits,
		destroyed: false,
	}
}

// CreateSurface implements Instance.CreateSurface.
func (i *instance) CreateSurface(descriptor SurfaceDescriptor) (Surface, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return nil, fmt.Errorf("instance has been destroyed")
	}

	return newSurface(descriptor.Label), nil
}

// WGSLLanguageFeatures implements Instance.WGSLLanguageFeatures.
func (i *instance) WGSLLanguageFeatures() (*SupportedWGSLLanguageFeatures, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return nil, fmt.Errorf("instance has been destroyed")
	}

	return &SupportedWGSLLanguageFeatures{
		Features: []WGSLLanguageFeatureName{
			WGSLLanguageFeatureNameReadonlyAndReadwriteStorageTextures,
			WGSLLanguageFeatureNamePacked4x8IntegerDotProduct,
		},
	}, nil
}

// HasWGSLLanguageFeature implements Instance.HasWGSLLanguageFeature.
func (i *instance) HasWGSLLanguageFeature(feature WGSLLanguageFeatureName) (bool, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return false, fmt.Errorf("instance has been destroyed")
	}

	supported := []WGSLLanguageFeatureName{
		WGSLLanguageFeatureNameReadonlyAndReadwriteStorageTextures,
		WGSLLanguageFeatureNamePacked4x8IntegerDotProduct,
	}
	for _, f := range supported {
		if f == feature {
			return true, nil
		}
	}
	return false, nil
}

// RequestAdapter implements Instance.RequestAdapter.
func (i *instance) RequestAdapter(options RequestAdapterOptions) (Adapter, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.destroyed {
		return nil, fmt.Errorf("instance has been destroyed")
	}

	backendType := BackendTypeVulkan
	if options.BackendType != BackendTypeUndefined {
		backendType = options.BackendType
	}
	return NewAdapter(backendType, AdapterTypeDiscreteGPU), nil
}
