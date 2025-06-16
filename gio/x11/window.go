package x11

import (
	"sync"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/render"
	"github.com/jezek/xgb/xproto"
	"github.com/opensraph/sraph/gio"
	"github.com/opensraph/sraph/gpu"
)

var _ gio.PlatformWindow = (*x11Window)(nil)

type x11Window struct {
	attr     gio.WindowAttr
	eventBus *gio.EventBus

	driver        *x11Driver
	xw            xproto.Window
	xg            xproto.Gcontext
	xp            render.Picture
	xevents       chan xgb.Event
	width, height int

	mu       sync.Mutex
	released bool
}

// Close implements gio.PlatformWindow.
func (x *x11Window) Close() error {
	panic("unimplemented")
}

// Event implements gio.PlatformWindow.
func (x *x11Window) Event() *gio.EventBus {
	panic("unimplemented")
}

// SetAttr implements gio.PlatformWindow.
func (x *x11Window) SetAttr(attr gio.WindowAttr) error {
	panic("unimplemented")
}

// Surface implements gio.PlatformWindow.
func (x *x11Window) Surface() (gpu.Surface, error) {
	panic("unimplemented")
}

func newX11Window(d *x11Driver, o gio.WindowAttr) (*x11Window, error) {

}
