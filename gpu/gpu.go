package gpu

//go:generate go run gen.go

var instance = NewInstance(InstanceDescriptor{})

func RequestAdapter(descriptor RequestAdapterOptions) (Adapter, error) {
	return instance.RequestAdapter(descriptor)
}

func GetPreferredCanvasFormat() TextureFormat {
	panic("GetPreferredCanvasFormat is not implemented in this mock GPU package")
}

func WGSLLanguageFeatures() (*SupportedWGSLLanguageFeatures, error) {
	return instance.WGSLLanguageFeatures()
}
