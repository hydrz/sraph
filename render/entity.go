package render

import (
	"fmt"

	"github.com/opensraph/sraph/display"
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/gpu"
)

// Entity represents a renderable entity that can be processed by the renderer.
// This is inspired by Impeller's entity system.
type Entity interface {
	// GetTransform returns the transformation matrix for this entity.
	GetTransform() geom.Matrix[geom.F32]

	// GetBounds returns the bounds of this entity in local coordinates.
	GetBounds() geom.Rect[geom.F32]

	// GetBlendMode returns the blend mode for this entity.
	GetBlendMode() display.BlendMode

	// Render renders this entity using the given render context.
	Render(ctx *EntityRenderContext) error
}

// EntityRenderContext provides context for rendering entities.
type EntityRenderContext struct {
	context    Context
	renderPass *RenderPass
	transform  geom.Matrix[geom.F32]
	clipBounds geom.Rect[geom.F32]
	alpha      float32
}

// NewEntityRenderContext creates a new entity render context.
func NewEntityRenderContext(context Context, renderPass *RenderPass) *EntityRenderContext {
	return &EntityRenderContext{
		context:    context,
		renderPass: renderPass,
		transform:  geom.NewMatrix[geom.F32](),
		clipBounds: geom.Rect[geom.F32]{},
		alpha:      1.0,
	}
}

// GetContext returns the rendering context.
func (ctx *EntityRenderContext) GetContext() Context {
	return ctx.context
}

// GetRenderPass returns the current render pass.
func (ctx *EntityRenderContext) GetRenderPass() *RenderPass {
	return ctx.renderPass
}

// GetTransform returns the current transformation matrix.
func (ctx *EntityRenderContext) GetTransform() geom.Matrix[geom.F32] {
	return ctx.transform
}

// SetTransform sets the transformation matrix.
func (ctx *EntityRenderContext) SetTransform(transform geom.Matrix[geom.F32]) {
	ctx.transform = transform
}

// GetClipBounds returns the current clipping bounds.
func (ctx *EntityRenderContext) GetClipBounds() geom.Rect[geom.F32] {
	return ctx.clipBounds
}

// SetClipBounds sets the clipping bounds.
func (ctx *EntityRenderContext) SetClipBounds(bounds geom.Rect[geom.F32]) {
	ctx.clipBounds = bounds
}

// GetAlpha returns the current alpha value.
func (ctx *EntityRenderContext) GetAlpha() float32 {
	return ctx.alpha
}

// SetAlpha sets the alpha value.
func (ctx *EntityRenderContext) SetAlpha(alpha float32) {
	ctx.alpha = alpha
}

// PushTransform pushes a new transformation onto the transform stack.
func (ctx *EntityRenderContext) PushTransform(transform geom.Matrix[geom.F32]) {
	ctx.transform = ctx.transform.Mul(transform)
}

// EntityPass represents a rendering pass for entities.
type EntityPass struct {
	context      Context
	renderTarget *RenderTarget
	entities     []Entity
	delegate     EntityPassDelegate
}

// EntityPassDelegate handles entity-specific rendering operations.
type EntityPassDelegate interface {
	// CanApplyOpacityPeephole determines if opacity peephole optimization can be applied.
	CanApplyOpacityPeephole() bool

	// SetupRenderPass sets up the render pass for entity rendering.
	SetupRenderPass(renderPass *RenderPass) error

	// GetClearColor returns the clear color for the render pass.
	GetClearColor(size geom.Size[geom.F32]) display.Color
}

// NewEntityPass creates a new entity pass.
func NewEntityPass(context Context, renderTarget *RenderTarget) *EntityPass {
	return &EntityPass{
		context:      context,
		renderTarget: renderTarget,
		entities:     make([]Entity, 0),
	}
}

// AddEntity adds an entity to the pass.
func (ep *EntityPass) AddEntity(entity Entity) {
	ep.entities = append(ep.entities, entity)
}

// SetDelegate sets the entity pass delegate.
func (ep *EntityPass) SetDelegate(delegate EntityPassDelegate) {
	ep.delegate = delegate
}

// Render renders all entities in the pass.
func (ep *EntityPass) Render() error {
	if len(ep.entities) == 0 {
		return nil
	}

	// Create render pass
	renderPass, err := NewRenderPass(ep.context, ep.renderTarget, "Entity Pass")
	if err != nil {
		return fmt.Errorf("failed to create render pass: %w", err)
	}

	// Setup render pass using delegate
	if ep.delegate != nil {
		if err := ep.delegate.SetupRenderPass(renderPass); err != nil {
			return fmt.Errorf("failed to setup render pass: %w", err)
		}
	}

	// Create render context
	renderContext := NewEntityRenderContext(ep.context, renderPass)

	// Render all entities
	for _, entity := range ep.entities {
		if err := entity.Render(renderContext); err != nil {
			return fmt.Errorf("failed to render entity: %w", err)
		}
	}

	// End render pass and submit
	commandBuffer := renderPass.End()
	if commandBuffer != nil {
		ep.context.GetQueue().Submit([]gpu.CommandBuffer{commandBuffer})
	}

	return nil
}

// SolidColorEntity represents a simple solid color entity.
type SolidColorEntity struct {
	color     display.Color
	bounds    geom.Rect[geom.F32]
	transform geom.Matrix[geom.F32]
}

// NewSolidColorEntity creates a new solid color entity.
func NewSolidColorEntity(color display.Color, bounds geom.Rect[geom.F32]) *SolidColorEntity {
	return &SolidColorEntity{
		color:     color,
		bounds:    bounds,
		transform: geom.NewMatrix[geom.F32](),
	}
}

// GetTransform implements Entity.
func (sce *SolidColorEntity) GetTransform() geom.Matrix[geom.F32] {
	return sce.transform
}

// SetTransform sets the transformation matrix.
func (sce *SolidColorEntity) SetTransform(transform geom.Matrix[geom.F32]) {
	sce.transform = transform
}

// GetBounds implements Entity.
func (sce *SolidColorEntity) GetBounds() geom.Rect[geom.F32] {
	return sce.bounds
}

// GetBlendMode implements Entity.
func (sce *SolidColorEntity) GetBlendMode() display.BlendMode {
	return display.BlendModeSrcOver
}

// Render implements Entity.
func (sce *SolidColorEntity) Render(ctx *EntityRenderContext) error {
	// TODO: Implement solid color rendering using appropriate shaders
	// For now, this is a placeholder implementation
	return nil
}

// GetColor returns the entity's color.
func (sce *SolidColorEntity) GetColor() display.Color {
	return sce.color
}

// SetColor sets the entity's color.
func (sce *SolidColorEntity) SetColor(color display.Color) {
	sce.color = color
}
