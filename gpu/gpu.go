package gpu

//go:generate go run ./gen -yaml ./gen/webgpu.yml -go wgpu.go

var instance Instance

func RequestAdapter(descriptor RequestAdapterOptions) (Adapter, error) {
	return instance.RequestAdapter(descriptor)
}

func GetPreferredCanvasFormat() TextureFormat {
	panic("GetPreferredCanvasFormat is not implemented in this mock GPU package")
}

func WGSLLanguageFeatures() (*SupportedWGSLLanguageFeatures, error) {
	return instance.WGSLLanguageFeatures()
}
