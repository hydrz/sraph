// Package main provides a shapes drawing example for the Sraph UI toolkit.
package main

import (
	"context"
	"log"

	"github.com/opensraph/sraph/app"
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/widget"
)

// ShapesWidget is a widget that displays various geometric shapes.
type ShapesWidget struct {
	key widget.Key
}

// Build implements Widget.
func (s *ShapesWidget) Build(context widget.BuildContext) widget.Widget {
	// TODO: Return a widget that draws various shapes (circles, rectangles, etc.)
	// This is a placeholder implementation
	return s
}

// Key implements Widget.
func (s *ShapesWidget) Key() widget.Key {
	return s.key
}

// SetKey implements Widget.
func (s *ShapesWidget) SetKey(key widget.Key) {
	s.key = key
}

// NewShapesWidget creates a new ShapesWidget.
func NewShapesWidget() widget.Widget {
	return &ShapesWidget{}
}

// Shape represents a drawable shape.
type Shape struct {
	Position geom.Point[geom.F32]
	Size     geom.Size[geom.F32]
	Color    geom.Color
	Type     ShapeType
}

// ShapeType represents the type of shape.
type ShapeType int

const (
	ShapeTypeRectangle ShapeType = iota
	ShapeTypeCircle
	ShapeTypeTriangle
	ShapeTypeLine
)

func main() {
	// Initialize the application
	application := app.NewApp()
	if application == nil {
		log.Fatal("Failed to create application")
	}

	// Create the root widget
	rootWidget := NewShapesWidget()

	// Run the application
	ctx := context.Background()
	if err := application.Run(ctx, rootWidget); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
