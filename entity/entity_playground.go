package entity

import (
	"github.com/opensraph/sraph/font"
	"github.com/opensraph/sraph/gpu"
	"github.com/opensraph/sraph/render"
)

// EntityPlayground provides a testing and development environment for entity rendering.
// It extends the base playground functionality with entity-specific features.
type EntityPlayground interface {
	// SetTypographerContext sets the typographer context for text rendering
	SetTypographerContext(typographerContext font.TypographerContext)

	// GetContentContext returns the content context for rendering
	GetContentContext() *ContentContext

	// OpenPlaygroundHere opens a playground with a single entity
	OpenPlaygroundHere(entity Entity) bool

	// OpenPlaygroundHereWithCallback opens a playground with a callback function
	OpenPlaygroundHereWithCallback(callback EntityPlaygroundCallback) bool

	// IsPlaygroundEnabled returns true if playground functionality is enabled
	IsPlaygroundEnabled() bool
}

// EntityPlaygroundCallback is a function type for playground rendering callbacks
type EntityPlaygroundCallback func(context *ContentContext, pass *render.RenderPass) bool

// PlaygroundSwitches holds configuration switches for playground functionality
type PlaygroundSwitches struct {
	EnablePlayground bool // Whether playground is enabled
	EnableDebugInfo  bool // Whether to show debug information
	EnableWireframe  bool // Whether to enable wireframe rendering
}

// DefaultEntityPlayground provides a default implementation of EntityPlayground
type DefaultEntityPlayground struct {
	typographerContext font.TypographerContext // Typographer context for text rendering
	gpuContext         gpu.Context             // GPU context
	switches           *PlaygroundSwitches     // Configuration switches
	contentContext     *ContentContext         // Cached content context
}

// NewEntityPlayground creates a new entity playground
func NewEntityPlayground(gpuContext gpu.Context) EntityPlayground {
	// Create default typographer context
	typographerContext := font.NewSkiaTypographerContext()

	return &DefaultEntityPlayground{
		typographerContext: typographerContext,
		gpuContext:         gpuContext,
		switches: &PlaygroundSwitches{
			EnablePlayground: true,
			EnableDebugInfo:  false,
			EnableWireframe:  false,
		},
		contentContext: nil,
	}
}

// SetTypographerContext sets the typographer context for text rendering
func (p *DefaultEntityPlayground) SetTypographerContext(typographerContext font.TypographerContext) {
	p.typographerContext = typographerContext
	// Invalidate cached content context
	p.contentContext = nil
}

// GetContentContext returns the content context for rendering
func (p *DefaultEntityPlayground) GetContentContext() *ContentContext {
	if p.contentContext == nil {
		p.contentContext = NewContentContext(p.gpuContext, p.typographerContext)
	}
	return p.contentContext
}

// OpenPlaygroundHere opens a playground with a single entity
func (p *DefaultEntityPlayground) OpenPlaygroundHere(entity Entity) bool {
	if !p.IsPlaygroundEnabled() {
		return false
	}

	// Create a callback that renders the entity
	callback := func(context *ContentContext, pass *render.RenderPass) bool {
		if entity == nil {
			return false
		}

		// Render the entity
		return entity.Render(context, pass)
	}

	return p.OpenPlaygroundHereWithCallback(callback)
}

// OpenPlaygroundHereWithCallback opens a playground with a callback function
func (p *DefaultEntityPlayground) OpenPlaygroundHereWithCallback(callback EntityPlaygroundCallback) bool {
	if !p.IsPlaygroundEnabled() || callback == nil {
		return false
	}

	// TODO: Implement actual playground rendering loop
	// This would typically:
	// 1. Create a render target
	// 2. Begin a render pass
	// 3. Call the callback function
	// 4. End the render pass and present
	// 5. Handle input and UI

	contentContext := p.GetContentContext()
	if contentContext == nil {
		return false
	}

	// Create a simple render target for testing
	renderTarget := p.createTestRenderTarget()
	if renderTarget == nil {
		return false
	}

	// Create command buffer
	commandBuffer := p.gpuContext.CreateCommandBuffer()
	if commandBuffer == nil {
		return false
	}

	// Create render pass
	renderPass := commandBuffer.CreateRenderPass(renderTarget)
	if renderPass == nil {
		return false
	}

	// Call the callback
	success := callback(contentContext, renderPass)

	// End the pass
	renderPass.End()

	// Submit the command buffer
	commandBuffer.Submit()

	return success
}

// IsPlaygroundEnabled returns true if playground functionality is enabled
func (p *DefaultEntityPlayground) IsPlaygroundEnabled() bool {
	return p.switches != nil && p.switches.EnablePlayground
}

// SetPlaygroundEnabled enables or disables playground functionality
func (p *DefaultEntityPlayground) SetPlaygroundEnabled(enabled bool) {
	if p.switches == nil {
		p.switches = &PlaygroundSwitches{}
	}
	p.switches.EnablePlayground = enabled
}

// SetDebugInfoEnabled enables or disables debug information display
func (p *DefaultEntityPlayground) SetDebugInfoEnabled(enabled bool) {
	if p.switches == nil {
		p.switches = &PlaygroundSwitches{}
	}
	p.switches.EnableDebugInfo = enabled
}

// SetWireframeEnabled enables or disables wireframe rendering
func (p *DefaultEntityPlayground) SetWireframeEnabled(enabled bool) {
	if p.switches == nil {
		p.switches = &PlaygroundSwitches{}
	}
	p.switches.EnableWireframe = enabled
}

// GetSwitches returns the current playground switches
func (p *DefaultEntityPlayground) GetSwitches() *PlaygroundSwitches {
	return p.switches
}

// createTestRenderTarget creates a simple render target for testing
func (p *DefaultEntityPlayground) createTestRenderTarget() *render.RenderTarget {
	// TODO: Implement actual render target creation
	// This is a placeholder implementation

	// Create a default sized render target
	size := gpu.Size{Width: 800, Height: 600}

	// Create color texture
	colorTexture := p.createColorTexture(size)
	if colorTexture == nil {
		return nil
	}

	// Create render target
	renderTarget := render.NewRenderTarget()
	if renderTarget == nil {
		return nil
	}

	renderTarget.SetColorTexture(colorTexture)

	return renderTarget
}

// createColorTexture creates a color texture for testing
func (p *DefaultEntityPlayground) createColorTexture(size gpu.Size) gpu.Texture {
	// TODO: Implement actual texture creation
	// This is a placeholder implementation
	return nil
}
