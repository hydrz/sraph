// Package framework provides the core framework components for the Sraph UI toolkit.
//
// This package integrates all the major components of the Sraph framework including
// widgets, elements, render objects, and bindings to provide a cohesive API.
package framework

import (
	"context"

	"github.com/opensraph/sraph/app"
	"github.com/opensraph/sraph/binding"
	"github.com/opensraph/sraph/element"
	"github.com/opensraph/sraph/render"
	"github.com/opensraph/sraph/widget"
)

// Application represents the main application framework that coordinates
// all components of the Sraph UI toolkit.
type Application struct {
	binding      binding.WidgetsBinding
	app          app.App
	rootWidget   widget.Widget
	rootElement  element.Element
	rootRenderer render.RenderObject
}

// NewApplication creates a new Sraph application.
func NewApplication() *Application {
	return &Application{
		binding: binding.NewDefaultWidgetsBinding(),
		app:     app.NewApp(),
	}
}

// Run starts the application with the given root widget.
func (a *Application) Run(ctx context.Context, rootWidget widget.Widget) error {
	a.rootWidget = rootWidget

	// Initialize the widget tree
	if err := a.initializeWidgetTree(); err != nil {
		return err
	}

	// Start the application
	return a.app.Run(ctx, rootWidget)
}

// SetRootWidget sets the root widget for the application.
func (a *Application) SetRootWidget(rootWidget widget.Widget) {
	a.rootWidget = rootWidget
	// TODO: Rebuild widget tree
}

// GetRootWidget returns the current root widget.
func (a *Application) GetRootWidget() widget.Widget {
	return a.rootWidget
}

// Binding returns the widgets binding instance.
func (a *Application) Binding() binding.WidgetsBinding {
	return a.binding
}

// App returns the app instance.
func (a *Application) App() app.App {
	return a.app
}

// initializeWidgetTree initializes the widget tree and element tree.
func (a *Application) initializeWidgetTree() error {
	// TODO: Create element tree from widget tree
	// TODO: Create render object tree from element tree
	// TODO: Setup bindings and connections
	return nil
}

// FrameworkConfig contains configuration options for the framework.
type FrameworkConfig struct {
	// EnableDebugMode enables debug features and logging.
	EnableDebugMode bool

	// WindowTitle specifies the default window title.
	WindowTitle string

	// WindowWidth specifies the default window width.
	WindowWidth int

	// WindowHeight specifies the default window height.
	WindowHeight int

	// VSync enables vertical synchronization.
	VSync bool

	// MaxFPS limits the maximum frames per second.
	MaxFPS int
}

// DefaultFrameworkConfig returns the default framework configuration.
func DefaultFrameworkConfig() *FrameworkConfig {
	return &FrameworkConfig{
		EnableDebugMode: false,
		WindowTitle:     "Sraph Application",
		WindowWidth:     800,
		WindowHeight:    600,
		VSync:           true,
		MaxFPS:          60,
	}
}

// Initialize initializes the Sraph framework with the given configuration.
func Initialize(config *FrameworkConfig) error {
	if config == nil {
		config = DefaultFrameworkConfig()
	}

	// TODO: Initialize platform services
	// TODO: Initialize GPU backends
	// TODO: Initialize font system
	// TODO: Initialize shader system
	// TODO: Register default widgets and render objects

	return nil
}

// Shutdown shuts down the Sraph framework and cleans up resources.
func Shutdown() error {
	// TODO: Cleanup resources
	// TODO: Shutdown platform services
	// TODO: Shutdown GPU backends

	return nil
}

// RunApp is a convenience function to create and run a Sraph application.
func RunApp(ctx context.Context, rootWidget widget.Widget, config *FrameworkConfig) error {
	if err := Initialize(config); err != nil {
		return err
	}
	defer Shutdown()

	app := NewApplication()
	return app.Run(ctx, rootWidget)
}
