package entity

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// EntityPass represents a single rendering pass for entities
// Based on Impeller's EntityPass class design
type EntityPass interface {
	// GetElementCount returns the number of elements in this pass
	GetElementCount() int

	// GetElements returns all elements in this pass
	GetElements() []Element

	// AddElement adds an element to this pass
	AddElement(element Element)

	// Render renders this pass to the given render target
	Render(context *EntityPassContext, target render.RenderTarget) error

	// GetClearColor returns the clear color for this pass
	GetClearColor() *geom.Color

	// SetClearColor sets the clear color for this pass
	SetClearColor(color *geom.Color)

	// GetSize returns the size of this pass
	GetSize() geom.ISize

	// SetSize sets the size of this pass
	SetSize(size geom.ISize)
}

// EntityPassContext provides context for entity pass rendering
type EntityPassContext struct {
	RenderContext  render.Context
	PassTarget     render.RenderTarget
	LazyGlyphAtlas interface{} // TODO: Replace with proper glyph atlas
}

// Element represents a single element in an entity pass
// This can be either an Entity or a nested EntityPass
type Element interface {
	// Render renders this element
	Render(context *EntityPassContext, pass render.RenderPass) error

	// GetCoverage returns the coverage area of this element
	GetCoverage(transform geom.Matrix[geom.F32]) *geom.Rect[geom.F32]

	// IsValid returns true if this element is valid for rendering
	IsValid() bool
}

// EntityPassImpl is the default implementation of EntityPass
type EntityPassImpl struct {
	elements         []Element
	clearColor       *geom.Color
	size             geom.ISize
	renderToOnscreen bool
	backdrop         *EntityPassImpl
}

// NewEntityPass creates a new entity pass
func NewEntityPass() EntityPass {
	return &EntityPassImpl{
		elements: make([]Element, 0),
		size:     geom.ISize{Width: 0, Height: 0},
	}
}

// GetElementCount returns the number of elements in this pass
func (ep *EntityPassImpl) GetElementCount() int {
	return len(ep.elements)
}

// GetElements returns all elements in this pass
func (ep *EntityPassImpl) GetElements() []Element {
	return ep.elements
}

// AddElement adds an element to this pass
func (ep *EntityPassImpl) AddElement(element Element) {
	if element != nil {
		ep.elements = append(ep.elements, element)
	}
}

// Render renders this pass to the given render target
func (ep *EntityPassImpl) Render(context *EntityPassContext, target render.RenderTarget) error {
	// Create render pass
	renderPass, err := context.RenderContext.CreateRenderPass(target)
	if err != nil {
		return err
	}
	defer renderPass.End()

	// Clear if needed
	if ep.clearColor != nil {
		renderPass.SetClearColor(*ep.clearColor)
	}

	// Render all elements
	for _, element := range ep.elements {
		if element.IsValid() {
			if err := element.Render(context, renderPass); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetClearColor returns the clear color for this pass
func (ep *EntityPassImpl) GetClearColor() *geom.Color {
	return ep.clearColor
}

// SetClearColor sets the clear color for this pass
func (ep *EntityPassImpl) SetClearColor(color *geom.Color) {
	ep.clearColor = color
}

// GetSize returns the size of this pass
func (ep *EntityPassImpl) GetSize() geom.ISize {
	return ep.size
}

// SetSize sets the size of this pass
func (ep *EntityPassImpl) SetSize(size geom.ISize) {
	ep.size = size
}

// EntityElement wraps an Entity to implement the Element interface
type EntityElement struct {
	entity Entity
}

// NewEntityElement creates a new entity element
func NewEntityElement(entity Entity) Element {
	return &EntityElement{entity: entity}
}

// Render renders the entity element
func (ee *EntityElement) Render(context *EntityPassContext, pass render.RenderPass) error {
	if ee.entity == nil {
		return nil
	}

	// Get entity properties
	transform := ee.entity.GetTransform()
	contents := ee.entity.GetContents()
	blendMode := ee.entity.GetBlendMode()

	// Apply transform
	pass.SetTransform(transform)

	// Set blend mode
	pass.SetBlendMode(blendMode)

	// Render contents
	if contents != nil {
		return contents.Render(context, pass, ee.entity)
	}

	return nil
}

// GetCoverage returns the coverage area of the entity
func (ee *EntityElement) GetCoverage(transform geom.Matrix[geom.F32]) *geom.Rect[geom.F32] {
	if ee.entity == nil {
		return nil
	}

	contents := ee.entity.GetContents()
	if contents == nil {
		return nil
	}

	// Combine entity transform with passed transform
	combinedTransform := transform.Multiply(ee.entity.GetTransform())
	return contents.GetCoverage(combinedTransform)
}

// IsValid returns true if the entity element is valid
func (ee *EntityElement) IsValid() bool {
	return ee.entity != nil && ee.entity.GetContents() != nil
}

// PassElement wraps a nested EntityPass to implement the Element interface
type PassElement struct {
	pass EntityPass
}

// NewPassElement creates a new pass element
func NewPassElement(pass EntityPass) Element {
	return &PassElement{pass: pass}
}

// Render renders the nested pass
func (pe *PassElement) Render(context *EntityPassContext, pass render.RenderPass) error {
	if pe.pass == nil {
		return nil
	}

	// Create a sub-target for the nested pass
	// TODO: Implement proper render target management
	return pe.pass.Render(context, context.PassTarget)
}

// GetCoverage returns the coverage area of the nested pass
func (pe *PassElement) GetCoverage(transform geom.Matrix[geom.F32]) *geom.Rect[geom.F32] {
	if pe.pass == nil {
		return nil
	}

	// For nested passes, return the full size as coverage
	size := pe.pass.GetSize()
	return &geom.Rect[geom.F32]{
		Origin: geom.Point[geom.F32]{X: 0, Y: 0},
		Size:   geom.Size[geom.F32]{Width: geom.F32(size.Width), Height: geom.F32(size.Height)},
	}
}

// IsValid returns true if the nested pass is valid
func (pe *PassElement) IsValid() bool {
	return pe.pass != nil && pe.pass.GetElementCount() > 0
}

// EntityPassBuilder helps build entity passes with a fluent interface
type EntityPassBuilder struct {
	pass *EntityPassImpl
}

// NewEntityPassBuilder creates a new entity pass builder
func NewEntityPassBuilder() *EntityPassBuilder {
	return &EntityPassBuilder{
		pass: &EntityPassImpl{
			elements: make([]Element, 0),
		},
	}
}

// AddEntity adds an entity to the pass
func (epb *EntityPassBuilder) AddEntity(entity Entity) *EntityPassBuilder {
	epb.pass.AddElement(NewEntityElement(entity))
	return epb
}

// AddPass adds a nested pass to the pass
func (epb *EntityPassBuilder) AddPass(pass EntityPass) *EntityPassBuilder {
	epb.pass.AddElement(NewPassElement(pass))
	return epb
}

// SetClearColor sets the clear color
func (epb *EntityPassBuilder) SetClearColor(color geom.Color) *EntityPassBuilder {
	epb.pass.SetClearColor(&color)
	return epb
}

// SetSize sets the pass size
func (epb *EntityPassBuilder) SetSize(size geom.ISize) *EntityPassBuilder {
	epb.pass.SetSize(size)
	return epb
}

// Build returns the constructed entity pass
func (epb *EntityPassBuilder) Build() EntityPass {
	return epb.pass
}
