package gpu

import (
	"context"
)

var gpu = &gpuImpl{}

func CreateInstance(ctx context.Context, descriptor InstanceDescriptor) (Instance, error) {
	return gpu.CreateInstance(ctx, descriptor)
}

func GetInstanceFeatures(ctx context.Context, features SupportedInstanceFeatures) error {
	return gpu.GetInstanceFeatures(ctx, features)
}

func GetInstanceLimits(ctx context.Context, limits InstanceLimits) (Status, error) {
	return gpu.GetInstanceLimits(ctx, limits)
}

func HasInstanceFeature(ctx context.Context, feature InstanceFeatureName) (bool, error) {
	return gpu.HasInstanceFeature(ctx, feature)
}

var _ GPU = (*gpuImpl)(nil)

// gpuImpl implements the main GPU interface
type gpuImpl struct {
}

// CreateInstance creates a new WebGPU instance
func (g *gpuImpl) CreateInstance(ctx context.Context, descriptor InstanceDescriptor) (Instance, error) {
	instance := newInstance(descriptor)
	return instance, nil
}

// GetInstanceFeatures retrieves supported instance features
func (g *gpuImpl) GetInstanceFeatures(ctx context.Context, features SupportedInstanceFeatures) error {
	// In a real implementation, this would populate the features struct
	// with supported instance features
	_ = features
	return nil
}

// GetInstanceLimits retrieves instance limits
func (g *gpuImpl) GetInstanceLimits(ctx context.Context, limits InstanceLimits) (Status, error) {
	// In a real implementation, this would populate the limits struct
	// with supported instance limits
	_ = limits
	return StatusSuccess, nil
}

// HasInstanceFeature checks if an instance feature is supported
func (g *gpuImpl) HasInstanceFeature(ctx context.Context, feature InstanceFeatureName) (bool, error) {
	// In a real implementation, this would check if the specific feature is supported
	// For now, return false for all features as a conservative default
	_ = feature
	return false, nil
}
