package entity

import (
	"github.com/opensraph/sraph/gpu"
	"github.com/opensraph/sraph/render"
)

// InlinePassContext manages inline rendering passes within entity rendering.
// It provides context for rendering operations that need to be performed
// inline within larger rendering passes.
type InlinePassContext interface {
	// IsValid returns true if the context is valid
	IsValid() bool

	// IsActive returns true if the context is currently active
	IsActive() bool

	// GetTexture returns the current texture being rendered to
	GetTexture() gpu.Texture

	// EndPass ends the current pass and returns success status
	EndPass(isOnscreen bool) bool

	// GetPassTarget returns the associated pass target
	GetPassTarget() EntityPassTarget

	// GetPassCount returns the number of passes executed
	GetPassCount() uint32

	// GetRenderPass returns the current render pass
	GetRenderPass() *render.RenderPass
}

// DefaultInlinePassContext provides a default implementation of InlinePassContext
type DefaultInlinePassContext struct {
	renderer      *ContentContext       // Content renderer reference
	passTarget    EntityPassTarget      // Associated pass target
	commandBuffer *render.CommandBuffer // Command buffer for recording commands
	renderPass    *render.RenderPass    // Current render pass
	passCount     uint32                // Number of passes executed
	isActive      bool                  // Whether the context is active
}

// NewInlinePassContext creates a new inline pass context
func NewInlinePassContext(renderer *ContentContext, passTarget EntityPassTarget) InlinePassContext {
	return &DefaultInlinePassContext{
		renderer:      renderer,
		passTarget:    passTarget,
		commandBuffer: nil,
		renderPass:    nil,
		passCount:     0,
		isActive:      false,
	}
}

// IsValid returns true if the context is valid
func (c *DefaultInlinePassContext) IsValid() bool {
	return c.renderer != nil && c.passTarget != nil && c.passTarget.IsValid()
}

// IsActive returns true if the context is currently active
func (c *DefaultInlinePassContext) IsActive() bool {
	return c.isActive && c.renderPass != nil
}

// GetTexture returns the current texture being rendered to
func (c *DefaultInlinePassContext) GetTexture() gpu.Texture {
	if !c.IsValid() {
		return nil
	}

	renderTarget := c.passTarget.GetRenderTarget()
	if renderTarget != nil {
		return renderTarget.GetColorTexture()
	}

	return nil
}

// EndPass ends the current pass and returns success status
func (c *DefaultInlinePassContext) EndPass(isOnscreen bool) bool {
	if !c.IsActive() {
		return false
	}

	// End the current render pass
	if c.renderPass != nil {
		success := c.renderPass.End()
		if !success {
			return false
		}
		c.renderPass = nil
	}

	// Submit command buffer if not onscreen
	if !isOnscreen && c.commandBuffer != nil {
		success := c.commandBuffer.Submit()
		if !success {
			return false
		}
		c.commandBuffer = nil
	}

	c.isActive = false
	c.passCount++

	return true
}

// GetPassTarget returns the associated pass target
func (c *DefaultInlinePassContext) GetPassTarget() EntityPassTarget {
	return c.passTarget
}

// GetPassCount returns the number of passes executed
func (c *DefaultInlinePassContext) GetPassCount() uint32 {
	return c.passCount
}

// GetRenderPass returns the current render pass
func (c *DefaultInlinePassContext) GetRenderPass() *render.RenderPass {
	return c.renderPass
}

// BeginPass begins a new render pass (internal method)
func (c *DefaultInlinePassContext) beginPass() bool {
	if !c.IsValid() || c.IsActive() {
		return false
	}

	// Create command buffer if needed
	if c.commandBuffer == nil {
		context := c.renderer.GetContext()
		if context == nil {
			return false
		}
		c.commandBuffer = context.CreateCommandBuffer()
		if c.commandBuffer == nil {
			return false
		}
	}

	// Begin render pass
	renderTarget := c.passTarget.GetRenderTarget()
	if renderTarget == nil {
		return false
	}

	c.renderPass = c.commandBuffer.CreateRenderPass(renderTarget)
	if c.renderPass == nil {
		return false
	}

	c.isActive = true
	return true
}
