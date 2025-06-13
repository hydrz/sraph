package gpu

import "github.com/opensraph/sraph/gpu/impl"

//go:generate go run gen.go

var gpu = impl.NewGPU()

func CreateInstance(descriptor InstanceDescriptor) (Instance, error) {
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
