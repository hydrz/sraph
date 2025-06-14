package gpu

import (
	"sync"

	. "github.com/opensraph/sraph/gpu/wgpu"
)

var _ Instance = (*instanceImpl)(nil)

// instanceImpl implements the Instance interface
type instanceImpl struct {
	mu       sync.RWMutex
	features []InstanceFeatureName
	limits   InstanceLimits
}

// NewInstance creates a new WebGPU instance
func NewInstance(descriptor InstanceDescriptor) Instance {
	return &instanceImpl{
		features: descriptor.RequiredFeatures,
		limits:   descriptor.RequiredLimits,
	}
}

// RequestAdapter implements Instance.RequestAdapter.
func (i *instanceImpl) RequestAdapter(options RequestAdapterOptions) (Adapter, error) {

}

// CreateSurface implements Instance.CreateSurface.
func (i *instanceImpl) CreateSurface(descriptor SurfaceDescriptor) (Surface, error) {

}

// WGSLLanguageFeatures implements Instance.WGSLLanguageFeatures.
func (i *instanceImpl) WGSLLanguageFeatures() (*SupportedWGSLLanguageFeatures, error) {

}

// HasWGSLLanguageFeature implements Instance.HasWGSLLanguageFeature.
func (i *instanceImpl) HasWGSLLanguageFeature(feature WGSLLanguageFeatureName) (bool, error) {

}

func (i *instanceImpl) EnumerateAdapters(options RequestAdapterOptions) ([]Adapter, error) {

}
