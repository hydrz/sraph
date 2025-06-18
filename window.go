package sraph

import (
	"context"

	"github.com/opensraph/sraph/geom"
)

type Window interface {
	Run(e Element) error
}

type WindowOptions struct {
	Title string
	Size  geom.Size[geom.F32]
	Pos   geom.Point[geom.F32]
}

var defaultWindowOptions = WindowOptions{
	Title: "Sraph Window",
	Size:  geom.Size[geom.F32]{Width: 800, Height: 600},
	Pos:   geom.Point[geom.F32]{X: 100, Y: 100},
}

type window struct {
	ctx     context.Context
	options WindowOptions
}

func (w *window) Run(e Element) error {
	// Implement the logic to run the window with the given element
	// This is a placeholder for the actual implementation
	return nil
}
