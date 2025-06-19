package sraph

import (
	"context"
	"log/slog"
)

type AppOptions struct {
	Ctx    context.Context
	Logger Logger
}

func (o *AppOptions) apply(opts ...AppOptions) {
	for _, opt := range opts {
		if opt.Ctx != nil {
			o.Ctx = opt.Ctx
		}
		if opt.Logger != nil {
			o.Logger = opt.Logger
		}
	}
}

var defaultAppOptions = AppOptions{
	Ctx:    context.Background(),
	Logger: slog.Default(),
}

type app struct {
	opt AppOptions
}

func (a *app) NewWindow(opts ...WindowOptions) Window {
	return &window{
		ctx: a.opt.Ctx,
	}
}

func (a *app) Run() error {
	select {
	case <-a.opt.Ctx.Done():
		return a.opt.Ctx.Err()
	default:
		// Simulate running the app
		return nil
	}
}
