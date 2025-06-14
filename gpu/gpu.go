package gpu

import (
	_ "github.com/goccy/go-yaml"
	_ "github.com/santhosh-tekuri/jsonschema/v5"
)

//go:generate go run gen.go -go wgpu.go

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
