// Package sraph provides a UI toolkit for building beautiful, natively compiled applications.
//
// Sraph is inspired by Flutter and provides a reactive programming model for building
// user interfaces. It consists of several layers:
//
// 1. Application Layer: app package provides the main application framework
// 2. UI Layer: widget package provides declarative UI components
// 3. Element Layer: element package manages widget instantiation and lifecycle
// 4. Content Layer: display, entity, font packages handle rendering content
// 5. Geometry Layer: geom package provides mathematical primitives
// 6. Tessellation Layer: tess package handles path tessellation
// 7. Rendering Layer: render package manages the render object tree
// 8. Graphics Backend Layer: gpu, shader packages provide hardware abstraction
// 9. Platform Layer: gio package handles platform integration
//
// Example usage:
//
//	package main
//
//	import (
//		"context"
//		"github.com/opensraph/sraph"
//		"github.com/opensraph/sraph/widget"
//	)
//
//	func main() {
//		app := sraph.NewApp()
//		rootWidget := widget.NewText("Hello, World!", nil)
//		app.Run(context.Background(), rootWidget)
//	}
package sraph

import (
	"context"

	"github.com/opensraph/sraph/framework"
	"github.com/opensraph/sraph/widget"
)

// App represents a Sraph application.
type App = framework.Application

// Config represents application configuration.
type Config = framework.FrameworkConfig

// NewApp creates a new Sraph application.
func NewApp() *App {
	return framework.NewApplication()
}

// NewConfig creates a new application configuration with default values.
func NewConfig() *Config {
	return framework.DefaultFrameworkConfig()
}

// RunApp is a convenience function to create and run a Sraph application.
func RunApp(ctx context.Context, rootWidget widget.Widget, config *Config) error {
	return framework.RunApp(ctx, rootWidget, config)
}

// Initialize initializes the Sraph framework.
func Initialize(config *Config) error {
	return framework.Initialize(config)
}

// Shutdown shuts down the Sraph framework.
func Shutdown() error {
	return framework.Shutdown()
}
