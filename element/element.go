// Package element provides the Element framework for managing widget lifecycle.
//
// Elements are the instantiations of Widgets at particular locations in the tree.
// Elements have a lifecycle and maintain state across builds.
package element

import (
	"github.com/opensraph/sraph/render"
)

// Element represents an instantiation of a Widget at a particular location in the tree.
// Elements have the following lifecycle:
// 1. The framework calls CreateElement to create an element.
// 2. The framework calls mount to add the element to the tree at a given location.
// 3. The framework calls build to ask the element to describe its children.
// 4. The framework may call update to change the widget that configures this element.
// 5. The framework may call deactivate and then activate to move the element to a new location.
// 6. The framework calls unmount to remove the element from the tree.
type Element interface {
	// Widget returns the widget that configures this element.
	Widget() interface{} // Widget interface

	// RenderObject returns the render object for this element, if any.
	RenderObject() render.RenderObject

	// Mount this element at the given location in the tree.
	Mount(parent Element, slot interface{})

	// Update this element's configuration with a new widget.
	Update(newWidget interface{}) // Widget interface

	// Unmount this element from the tree.
	Unmount()

	// Activate this element after it has been moved to a new location.
	Activate()

	// Deactivate this element before moving it to a new location.
	Deactivate()

	// MarkNeedsBuild marks this element as needing to be rebuilt.
	MarkNeedsBuild()

	// Rebuild this element if it has been marked as needing a rebuild.
	Rebuild()
}

// ComponentElement is an Element that uses a Widget as its configuration.
type ComponentElement interface {
	Element

	// Build asks the widget to describe its children given the current configuration.
	Build() interface{} // Widget interface
}

// RenderObjectElement is an Element that uses a RenderObjectWidget as its configuration.
type RenderObjectElement interface {
	Element

	// UpdateRenderObject updates the render object's configuration.
	UpdateRenderObject()
}

// StatelessElement is an Element that uses a StatelessWidget as its configuration.
type StatelessElement struct {
	// TODO: Implement StatelessElement
}

// StatefulElement is an Element that uses a StatefulWidget as its configuration.
type StatefulElement struct {
	// TODO: Implement StatefulElement
}

// LeafRenderObjectElement is a RenderObjectElement that has no children.
type LeafRenderObjectElement struct {
	// TODO: Implement LeafRenderObjectElement
}

// SingleChildRenderObjectElement is a RenderObjectElement that has exactly one child.
type SingleChildRenderObjectElement struct {
	// TODO: Implement SingleChildRenderObjectElement
}

// MultiChildRenderObjectElement is a RenderObjectElement that has multiple children.
type MultiChildRenderObjectElement struct {
	// TODO: Implement MultiChildRenderObjectElement
}
