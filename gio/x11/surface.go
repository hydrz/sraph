package x11

import "github.com/opensraph/sraph/gpu/wgpu"

var _ wgpu.Surface = &X11Surface{}

type X11Surface struct{}

// Capabilities implements wgpu.Surface.
func (x *X11Surface) Capabilities(adapter wgpu.Adapter) (*wgpu.SurfaceCapabilities, error) {
	panic("unimplemented")
}

// Configure implements wgpu.Surface.
func (x *X11Surface) Configure(config wgpu.SurfaceConfiguration) error {
	panic("unimplemented")
}

// CurrentTexture implements wgpu.Surface.
func (x *X11Surface) CurrentTexture() (*wgpu.SurfaceTexture, error) {
	panic("unimplemented")
}

// Present implements wgpu.Surface.
func (x *X11Surface) Present() error {
	panic("unimplemented")
}

// SetLabel implements wgpu.Surface.
func (x *X11Surface) SetLabel(label string) error {
	panic("unimplemented")
}

// Unconfigure implements wgpu.Surface.
func (x *X11Surface) Unconfigure() error {
	panic("unimplemented")
}
