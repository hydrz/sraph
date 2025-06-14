package x11

import "github.com/opensraph/sraph/gpu"

var _ gpu.Surface = (*X11Surface)(nil)

type X11Surface struct{}

// Capabilities implements gpu.Surface.
func (x *X11Surface) Capabilities(adapter gpu.Adapter) (*gpu.SurfaceCapabilities, error) {
	panic("unimplemented")
}

// Configure implements gpu.Surface.
func (x *X11Surface) Configure(config gpu.SurfaceConfiguration) {
	panic("unimplemented")
}

// CurrentTexture implements gpu.Surface.
func (x *X11Surface) CurrentTexture() *gpu.SurfaceTexture {
	panic("unimplemented")
}

// Present implements gpu.Surface.
func (x *X11Surface) Present() error {
	panic("unimplemented")
}

// SetLabel implements gpu.Surface.
func (x *X11Surface) SetLabel(label string) {
	panic("unimplemented")
}

// Unconfigure implements gpu.Surface.
func (x *X11Surface) Unconfigure() {
	panic("unimplemented")
}
