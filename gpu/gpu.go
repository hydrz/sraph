package gpu

import "github.com/opensraph/sraph/gpu/impl"

//go:generate go run gen.go

var gpu = &gpuImpl{}

func NewInstance(descriptor InstanceDescriptor) (Instance, error) {
	return gpu.CreateInstance(descriptor)
}

func GetInstanceFeatures(features SupportedInstanceFeatures) error {
	return gpu.GetInstanceFeatures(features)
}

func GetInstanceLimits(limits InstanceLimits) (Status, error) {
	return gpu.GetInstanceLimits(limits)
}

func HasInstanceFeature(feature InstanceFeatureName) (bool, error) {
	return gpu.HasInstanceFeature(feature)
}

var _ GPU = (*gpuImpl)(nil)

// gpuImpl implements the main GPU interface
type gpuImpl struct{}

// CreateInstance creates a new WebGPU instance
func (g *gpuImpl) CreateInstance(descriptor InstanceDescriptor) (Instance, error) {
	instance := impl.NewInstance(descriptor)
	return instance, nil
}

// GetInstanceFeatures retrieves supported instance features
func (g *gpuImpl) GetInstanceFeatures(features SupportedInstanceFeatures) error {
	// In a real implementation, this would populate the features struct
	// with supported instance features
	_ = features
	return nil
}

// GetInstanceLimits retrieves instance limits
func (g *gpuImpl) GetInstanceLimits(limits InstanceLimits) (Status, error) {
	// In a real implementation, this would populate the limits struct
	// with supported instance limits
	_ = limits
	return StatusSuccess, nil
}

// HasInstanceFeature checks if an instance feature is supported
func (g *gpuImpl) HasInstanceFeature(feature InstanceFeatureName) (bool, error) {
	// In a real implementation, this would check if the specific feature is supported
	// For now, return false for all features as a conservative default
	_ = feature
	return false, nil
}
