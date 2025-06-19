// Package main provides an animation example for the Sraph UI toolkit.
package main

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/opensraph/sraph/app"
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/widget"
)

// AnimationController controls animation timing and interpolation.
type AnimationController struct {
	duration  time.Duration
	startTime time.Time
	running   bool
	value     float32
}

// NewAnimationController creates a new animation controller.
func NewAnimationController(duration time.Duration) *AnimationController {
	return &AnimationController{
		duration: duration,
		value:    0.0,
	}
}

// Start starts the animation.
func (c *AnimationController) Start() {
	c.startTime = time.Now()
	c.running = true
}

// Stop stops the animation.
func (c *AnimationController) Stop() {
	c.running = false
}

// Update updates the animation value based on elapsed time.
func (c *AnimationController) Update() {
	if !c.running {
		return
	}

	elapsed := time.Since(c.startTime)
	if elapsed >= c.duration {
		c.value = 1.0
		c.running = false
		return
	}

	// Linear interpolation
	c.value = float32(elapsed.Seconds() / c.duration.Seconds())
}

// Value returns the current animation value (0.0 to 1.0).
func (c *AnimationController) Value() float32 {
	return c.value
}

// IsRunning returns true if the animation is currently running.
func (c *AnimationController) IsRunning() bool {
	return c.running
}

// AnimatedWidget is a widget that displays an animated shape.
type AnimatedWidget struct {
	key        widget.Key
	controller *AnimationController
}

// Build implements Widget.
func (a *AnimatedWidget) Build(context widget.BuildContext) widget.Widget {
	// Update animation
	a.controller.Update()

	// TODO: Return a widget that draws an animated shape
	// The animation value can be used to interpolate position, size, color, etc.
	// This is a placeholder implementation
	return a
}

// Key implements Widget.
func (a *AnimatedWidget) Key() widget.Key {
	return a.key
}

// SetKey implements Widget.
func (a *AnimatedWidget) SetKey(key widget.Key) {
	a.key = key
}

// GetAnimationValue returns the current animation value with easing applied.
func (a *AnimatedWidget) GetAnimationValue() float32 {
	t := a.controller.Value()
	// Apply ease-in-out cubic easing
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - float32(math.Pow(float64(-2*t+2), 3))/2
}

// GetPosition returns the interpolated position based on animation value.
func (a *AnimatedWidget) GetPosition() geom.Point[geom.F32] {
	value := a.GetAnimationValue()
	startX := geom.F32(50.0)
	endX := geom.F32(200.0)
	startY := geom.F32(50.0)
	endY := geom.F32(150.0)

	x := startX + (endX-startX)*geom.F32(value)
	y := startY + (endY-startY)*geom.F32(value)

	return geom.Point[geom.F32]{X: x, Y: y}
}

// GetColor returns the interpolated color based on animation value.
func (a *AnimatedWidget) GetColor() geom.Color {
	value := a.GetAnimationValue()
	startColor := geom.NewColor(1.0, 0.0, 0.0, 1.0) // Red
	endColor := geom.NewColor(0.0, 0.0, 1.0, 1.0)   // Blue

	r := startColor.R + (endColor.R-startColor.R)*geom.F32(value)
	g := startColor.G + (endColor.G-startColor.G)*geom.F32(value)
	b := startColor.B + (endColor.B-startColor.B)*geom.F32(value)
	alpha := startColor.A + (endColor.A-startColor.A)*geom.F32(value)

	return geom.Color{R: r, G: g, B: b, A: alpha}
}

// NewAnimatedWidget creates a new AnimatedWidget.
func NewAnimatedWidget() *AnimatedWidget {
	controller := NewAnimationController(2 * time.Second)
	controller.Start()

	return &AnimatedWidget{
		controller: controller,
	}
}

// AnimationDemoWidget demonstrates various animation capabilities.
type AnimationDemoWidget struct {
	key widget.Key
}

// Build implements Widget.
func (d *AnimationDemoWidget) Build(context widget.BuildContext) widget.Widget {
	// TODO: Create a layout with multiple animated widgets
	// This is a placeholder implementation
	return d
}

// Key implements Widget.
func (d *AnimationDemoWidget) Key() widget.Key {
	return d.key
}

// SetKey implements Widget.
func (d *AnimationDemoWidget) SetKey(key widget.Key) {
	d.key = key
}

// NewAnimationDemoWidget creates a new AnimationDemoWidget.
func NewAnimationDemoWidget() widget.Widget {
	return &AnimationDemoWidget{}
}

func main() {
	// Initialize the application
	application := app.NewApp()
	if application == nil {
		log.Fatal("Failed to create application")
	}

	// Create the root widget
	rootWidget := NewAnimationDemoWidget()

	// Run the application
	ctx := context.Background()
	if err := application.Run(ctx, rootWidget); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
