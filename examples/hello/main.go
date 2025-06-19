// Package main provides a simple "Hello, World!" example for the Sraph UI toolkit.
package main

import (
	"context"
	"log"

	"github.com/opensraph/sraph/app"
	"github.com/opensraph/sraph/widget"
)

// HelloWidget is a simple stateless widget that displays "Hello, World!".
type HelloWidget struct {
	key widget.Key
}

// Build implements Widget.
func (h *HelloWidget) Build(context widget.BuildContext) widget.Widget {
	// TODO: Return a text widget displaying "Hello, World!"
	// This is a placeholder implementation
	return h
}

// Key implements Widget.
func (h *HelloWidget) Key() widget.Key {
	return h.key
}

// SetKey implements Widget.
func (h *HelloWidget) SetKey(key widget.Key) {
	h.key = key
}

// NewHelloWidget creates a new HelloWidget.
func NewHelloWidget() widget.Widget {
	return &HelloWidget{}
}

func main() {
	// Initialize the application
	application := app.NewApp()
	if application == nil {
		log.Fatal("Failed to create application")
	}

	// Create the root widget
	rootWidget := NewHelloWidget()

	// Run the application
	ctx := context.Background()
	if err := application.Run(ctx, rootWidget); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
