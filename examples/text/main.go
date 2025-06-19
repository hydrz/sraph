// Package main provides a text rendering example for the Sraph UI toolkit.
package main

import (
	"context"
	"log"

	"github.com/opensraph/sraph/app"
	"github.com/opensraph/sraph/font"
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/widget"
)

// TextWidget is a widget that displays formatted text.
type TextWidget struct {
	key   widget.Key
	text  string
	style *font.TextStyle
}

// Build implements Widget.
func (t *TextWidget) Build(context widget.BuildContext) widget.Widget {
	// TODO: Return a widget that renders the text with the given style
	// This is a placeholder implementation
	return t
}

// Key implements Widget.
func (t *TextWidget) Key() widget.Key {
	return t.key
}

// SetKey implements Widget.
func (t *TextWidget) SetKey(key widget.Key) {
	t.key = key
}

// SetText sets the text to display.
func (t *TextWidget) SetText(text string) {
	t.text = text
}

// SetStyle sets the text style.
func (t *TextWidget) SetStyle(style *font.TextStyle) {
	t.style = style
}

// Text returns the current text.
func (t *TextWidget) Text() string {
	return t.text
}

// Style returns the current text style.
func (t *TextWidget) Style() *font.TextStyle {
	return t.style
}

// NewTextWidget creates a new TextWidget.
func NewTextWidget(text string, style *font.TextStyle) *TextWidget {
	return &TextWidget{
		text:  text,
		style: style,
	}
}

// TextDemoWidget demonstrates various text rendering capabilities.
type TextDemoWidget struct {
	key widget.Key
}

// Build implements Widget.
func (d *TextDemoWidget) Build(context widget.BuildContext) widget.Widget {
	// TODO: Create a layout with multiple text widgets showing different styles
	// This is a placeholder implementation
	return d
}

// Key implements Widget.
func (d *TextDemoWidget) Key() widget.Key {
	return d.key
}

// SetKey implements Widget.
func (d *TextDemoWidget) SetKey(key widget.Key) {
	d.key = key
}

// NewTextDemoWidget creates a new TextDemoWidget.
func NewTextDemoWidget() widget.Widget {
	return &TextDemoWidget{}
}

func main() {
	// Initialize the application
	application := app.NewApp()
	if application == nil {
		log.Fatal("Failed to create application")
	}

	// Create some sample text styles
	normalStyle := &font.TextStyle{
		FontSize:   16.0,
		FontWeight: font.WeightNormal,
		FontStyle:  font.StyleNormal,
		Color:      geom.NewColor(0.0, 0.0, 0.0, 1.0),
	}

	titleStyle := &font.TextStyle{
		FontSize:   24.0,
		FontWeight: font.WeightBold,
		FontStyle:  font.StyleNormal,
		Color:      geom.NewColor(0.2, 0.2, 0.8, 1.0),
	}

	_ = normalStyle
	_ = titleStyle

	// Create the root widget
	rootWidget := NewTextDemoWidget()

	// Run the application
	ctx := context.Background()
	if err := application.Run(ctx, rootWidget); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
