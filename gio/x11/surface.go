package x11

import (
	"github.com/opensraph/sraph/gpu"
)

var _ gpu.Surface = (*x11Surface)(nil)

type x11Surface struct {
	config *gpu.SurfaceConfiguration
}

// Capabilities implements gpu.Surface.
func (x *x11Surface) Capabilities(adapter gpu.Adapter) (*gpu.SurfaceCapabilities, error) {
	panic("unimplemented")
}

// Configure implements gpu.Surface.
func (x *x11Surface) Configure(config gpu.SurfaceConfiguration) {
	panic("unimplemented")
}

// CurrentTexture implements gpu.Surface.
func (x *x11Surface) CurrentTexture() *gpu.SurfaceTexture {
	panic("unimplemented")
}

// Present implements gpu.Surface.
func (x *x11Surface) Present() error {
	panic("unimplemented")
}

// SetLabel implements gpu.Surface.
func (x *x11Surface) SetLabel(label string) {
	panic("unimplemented")
}

// Unconfigure implements gpu.Surface.
func (x *x11Surface) Unconfigure() {
	panic("unimplemented")
}
