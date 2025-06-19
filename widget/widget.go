// Package widget provides UI widgets and layout functionality for the Sraph UI toolkit.
//
// This package defines the core widget system including basic widgets like buttons,
// text, containers, and layout management.
package widget

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// Widget represents a user interface widget.
type Widget interface {
	// Build builds the widget tree and returns child widgets.
	Build(context BuildContext) Widget

	// Key returns the key that identifies this widget.
	Key() Key

	// SetKey sets the key for this widget.
	SetKey(key Key)
}

// StatelessWidget represents a widget that does not maintain state.
type StatelessWidget interface {
	Widget

	// Build builds the widget tree for a stateless widget.
	Build(context BuildContext) Widget
}

// StatefulWidget represents a widget that maintains state.
type StatefulWidget interface {
	Widget

	// CreateState creates the state for this widget.
	CreateState() State
}

// State represents the state of a StatefulWidget.
type State interface {
	// Build builds the widget tree using the current state.
	Build(context BuildContext) Widget

	// InitState initializes the state when the widget is first created.
	InitState()

	// Dispose disposes of the state when the widget is removed.
	Dispose()

	// SetState schedules a rebuild of the widget with updated state.
	SetState(updateFn func())

	// Widget returns the widget that owns this state.
	Widget() StatefulWidget
}

// BuildContext provides context information during widget building.
type BuildContext interface {
	// Widget returns the widget being built.
	Widget() Widget

	// Size returns the available size for the widget.
	Size() geom.Size[geom.F32]

	// FindAncestorWidget finds an ancestor widget of the given type.
	FindAncestorWidget(widgetType interface{}) Widget

	// FindRenderObject finds the render object for this context.
	FindRenderObject() render.RenderObject
}

// Key represents a unique identifier for a widget.
type Key interface {
	// String returns the string representation of the key.
	String() string
}

// ValueKey is a key that uses a value for identification.
type ValueKey struct {
	Value interface{}
}

// String implements Key.
func (k ValueKey) String() string {
	if k.Value == nil {
		return "ValueKey(nil)"
	}
	return "ValueKey(" + k.Value.(string) + ")"
}

// Container represents a widget that can contain other widgets.
type Container interface {
	Widget

	// Child returns the child widget.
	Child() Widget

	// SetChild sets the child widget.
	SetChild(child Widget)
}

// MultiChildWidget represents a widget that can contain multiple children.
type MultiChildWidget interface {
	Widget

	// Children returns the child widgets.
	Children() []Widget

	// SetChildren sets the child widgets.
	SetChildren(children []Widget)

	// AddChild adds a child widget.
	AddChild(child Widget)

	// RemoveChild removes a child widget.
	RemoveChild(child Widget)
}

// RenderObjectWidget represents a widget that has an associated render object.
type RenderObjectWidget interface {
	Widget

	// CreateRenderObject creates the render object for this widget.
	CreateRenderObject() render.RenderObject

	// UpdateRenderObject updates the render object with new properties.
	UpdateRenderObject(renderObject render.RenderObject)
}

// LeafRenderObjectWidget represents a render object widget with no children.
type LeafRenderObjectWidget interface {
	RenderObjectWidget
}

// SingleChildRenderObjectWidget represents a render object widget with one child.
type SingleChildRenderObjectWidget interface {
	RenderObjectWidget
	Container
}

// MultiChildRenderObjectWidget represents a render object widget with multiple children.
type MultiChildRenderObjectWidget interface {
	RenderObjectWidget
	MultiChildWidget
}

// NewValueKey creates a new ValueKey with the given value.
func NewValueKey(value interface{}) Key {
	return ValueKey{Value: value}
}
