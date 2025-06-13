package gpu

import "github.com/opensraph/sraph/gpu/impl"

//go:generate go run gen.go

var instance = impl.NewInstance(InstanceDescriptor{})

func RequestAdapter(descriptor RequestAdapterOptions) (Adapter, error) {
	return instance.RequestAdapter(descriptor)
}

func GetPreferredCanvasFormat() TextureFormat {
	panic("GetPreferredCanvasFormat is not implemented in this mock GPU package")
}

func WGSLLanguageFeatures() (*SupportedWGSLLanguageFeatures, error) {
	return instance.WGSLLanguageFeatures()
}
