package x11

import (
	"runtime"

	"github.com/opensraph/sraph/gio"
	"github.com/opensraph/sraph/gpu"
)

var _ gio.Driver = (*X11Driver)(nil)

func init() {
	runtime.LockOSThread()
	gio.RegisterDriver(gio.DriverTypeX11, func() gio.Driver {
		return &X11Driver{}
	})
}

type X11Driver struct{}

// CreateSurface implements gio.Driver.
func (x *X11Driver) CreateSurface() (gpu.Surface, error) {
	panic("unimplemented")
}

// CreateWindow implements gio.Driver.
func (x *X11Driver) CreateWindow(options gio.NewWindowOptions) (gio.Window, error) {
	panic("unimplemented")
}

// Type implements gio.Driver.
func (x *X11Driver) Type() gio.DriverType {
	return gio.DriverTypeX11
}
