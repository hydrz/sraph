package impl

import (
	. "github.com/opensraph/sraph/gpu/wgpu"
)

var _ GPU = (*gpu)(nil)

// gpu implements the GPU interface
type gpu struct{}

// NewGPU creates a new WebGPU GPU instance
func NewGPU() GPU {
	return &gpu{}
}

// CreateInstance implements GPU.CreateInstance.
func (g *gpu) CreateInstance(descriptor InstanceDescriptor) (Instance, error) {
	return NewInstance(descriptor), nil
}

// InstanceFeatures implements GPU.InstanceFeatures.
func (g *gpu) InstanceFeatures() (*SupportedInstanceFeatures, error) {
	return &SupportedInstanceFeatures{
		Features: []InstanceFeatureName{
			InstanceFeatureNameTimedWaitAnyEnable,
			InstanceFeatureNameShaderSourceSPIRV,
			InstanceFeatureNameMultipleDevicesPerAdapter,
		},
	}, nil
}

// InstanceLimits implements GPU.InstanceLimits.
func (g *gpu) InstanceLimits() (*InstanceLimits, error) {
	return &InstanceLimits{
		TimedWaitAnyMaxCount: 64,
	}, nil
}

// HasInstanceFeature implements GPU.HasInstanceFeature.
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
