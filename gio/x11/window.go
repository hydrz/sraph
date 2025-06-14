package x11

import (
	"sync"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/render"
	"github.com/jezek/xgb/xproto"
	"github.com/opensraph/sraph/gio"
)

var _ gio.Window = (*X11Window)(nil)

type X11Window struct {
	width, height int

	xw xproto.Window
	xg xproto.Gcontext
	xp render.Picture
	xe chan xgb.Event

	mu sync.Mutex
}

// Receive implements gio.Window.
func (x *X11Window) Receive() <-chan gio.Event {
	panic("unimplemented")
}

// Release implements gio.Window.
func (x *X11Window) Release() {
	panic("unimplemented")
}

// Send implements gio.Window.
func (x *X11Window) Send(ch chan<- gio.Event) {
	panic("unimplemented")
}
