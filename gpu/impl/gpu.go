package impl

import (
	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ GPU = (*gpu)(nil)

// gpu implements the GPU interface
type gpu struct{}

// NewGPU creates a new WebGPU GPU instance
func NewGPU() GPU {
	return &gpu{}
}

// CreateInstance creates a new WebGPU instance
func (g *gpu) CreateInstance(descriptor InstanceDescriptor) (Instance, error) {
	return NewInstance(descriptor), nil
}

// GetInstanceFeatures retrieves supported instance features
func (g *gpu) GetInstanceFeatures(features SupportedInstanceFeatures) error {
	// Add default instance features
	features.Features = []InstanceFeatureName{
		InstanceFeatureNameTimedWaitAnyEnable,
		InstanceFeatureNameShaderSourceSPIRV,
		InstanceFeatureNameMultipleDevicesPerAdapter,
	}
	return nil
}

// GetInstanceLimits retrieves instance limits
func (g *gpu) GetInstanceLimits(limits InstanceLimits) (Status, error) {
	limits.TimedWaitAnyMaxCount = 64
	return StatusSuccess, nil
}

// HasInstanceFeature checks if an instance feature is supported
func (g *gpu) HasInstanceFeature(feature InstanceFeatureName) (bool, error) {
	supportedFeatures := []InstanceFeatureName{
		InstanceFeatureNameTimedWaitAnyEnable,
		InstanceFeatureNameShaderSourceSPIRV,
		InstanceFeatureNameMultipleDevicesPerAdapter,
	}

	for _, f := range supportedFeatures {
		if f == feature {
			return true, nil
		}
	}
	return false, nil
}
