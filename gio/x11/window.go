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
	eb            gio.EventBus
	width, height int

	xw xproto.Window
	xg xproto.Gcontext
	xp render.Picture
	xe chan xgb.Event

	mu sync.Mutex
}

// Release implements gio.Window.
func (x *X11Window) Release() {
	panic("unimplemented")
}
