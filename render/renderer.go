package render

import (
	"fmt"
	"sync"

	"github.com/opensraph/sraph/entity"
)

// Renderer is the main interface for rendering entities to render targets.
// Based on Impeller's Renderer class design.
type Renderer interface {
	// Render renders entities to the specified render target
	Render(target RenderTarget, entities []entity.Entity) error

	// RenderWithCallback renders with a custom setup callback
	RenderWithCallback(target RenderTarget, callback RenderCallback) error

	// GetContext returns the rendering context
	GetContext() Context

	// SetDebugLabel sets a debug label for the renderer
	SetDebugLabel(label string)

	// GetStats returns rendering statistics
	GetStats() RendererStats
}

// RenderCallback is a function that performs custom rendering operations
type RenderCallback func(pass RenderPass) error

// RendererStats contains rendering performance statistics
type RendererStats struct {
	// Frame statistics
	FramesRendered    uint64
	TotalRenderTime   float64 // in milliseconds
	AverageRenderTime float64 // in milliseconds
	LastFrameTime     float64 // in milliseconds

	// Entity statistics
	EntitiesRendered  uint64
	DrawCallsIssued   uint64
	TrianglesRendered uint64
	VerticesProcessed uint64

	// Resource statistics
	TexturesUsed  uint32
	BuffersUsed   uint32
	PipelinesUsed uint32

	// Memory statistics
	GPUMemoryUsed uint64 // in bytes
	CPUMemoryUsed uint64 // in bytes

	// Performance metrics
	GPUUtilization float32 // 0.0 to 1.0
	FrameRate      float32 // frames per second
}

// RendererImpl is the main implementation of the Renderer interface
type RendererImpl struct {
	context    Context
	debugLabel string
	stats      RendererStats
	statsMutex sync.RWMutex

	// Resource pools
	commandBufferPool Pool[CommandBuffer]
	renderPassPool    Pool[RenderPass]

	// Pipeline cache
	pipelineCache map[string]Pipeline
	cacheMutex    sync.RWMutex

	// Render state
	currentPass   RenderPass
	currentTarget RenderTarget

	// Configuration
	enableStats bool
	enableDebug bool
}

// NewRenderer creates a new renderer with the given context
func NewRenderer(context Context) Renderer {
	return &RendererImpl{
		context:       context,
		pipelineCache: make(map[string]Pipeline),
		enableStats:   true,
		enableDebug:   false,

		// Initialize resource pools
		commandBufferPool: NewPool[CommandBuffer](
			func() CommandBuffer {
				return context.CreateCommandBuffer()
			},
			func(cb CommandBuffer) {
				cb.Reset()
			},
		),
		renderPassPool: NewPool[RenderPass](
			func() RenderPass {
				// Create a default render pass
				return nil // TODO: Implement render pass creation
			},
			func(rp RenderPass) {
				// Reset render pass
			},
		),
	}
}

// Render renders entities to the specified render target
func (r *RendererImpl) Render(target RenderTarget, entities []entity.Entity) error {
	if target == nil {
		return fmt.Errorf("render target cannot be nil")
	}

	// Update stats
	if r.enableStats {
		r.statsMutex.Lock()
		r.stats.FramesRendered++
		r.statsMutex.Unlock()
	}

	// Get command buffer from pool
	commandBuffer := r.commandBufferPool.Get()
	defer r.commandBufferPool.Put(commandBuffer)

	// Begin command buffer
	if err := commandBuffer.Begin(); err != nil {
		return fmt.Errorf("failed to begin command buffer: %w", err)
	}

	// Create render pass
	renderPass, err := r.createRenderPass(commandBuffer, target)
	if err != nil {
		return fmt.Errorf("failed to create render pass: %w", err)
	}

	// Render entities
	err = r.renderEntities(renderPass, entities)
	if err != nil {
		renderPass.End()
		return fmt.Errorf("failed to render entities: %w", err)
	}

	// End render pass
	renderPass.End()

	// End and submit command buffer
	if err := commandBuffer.End(); err != nil {
		return fmt.Errorf("failed to end command buffer: %w", err)
	}

	if err := r.context.GetCommandQueue().Submit(commandBuffer); err != nil {
		return fmt.Errorf("failed to submit command buffer: %w", err)
	}

	return nil
}

// RenderWithCallback renders with a custom setup callback
func (r *RendererImpl) RenderWithCallback(target RenderTarget, callback RenderCallback) error {
	if target == nil {
		return fmt.Errorf("render target cannot be nil")
	}
	if callback == nil {
		return fmt.Errorf("render callback cannot be nil")
	}

	// Get command buffer from pool
	commandBuffer := r.commandBufferPool.Get()
	defer r.commandBufferPool.Put(commandBuffer)

	// Begin command buffer
	if err := commandBuffer.Begin(); err != nil {
		return fmt.Errorf("failed to begin command buffer: %w", err)
	}

	// Create render pass
	renderPass, err := r.createRenderPass(commandBuffer, target)
	if err != nil {
		return fmt.Errorf("failed to create render pass: %w", err)
	}

	// Execute callback
	err = callback(renderPass)
	if err != nil {
		renderPass.End()
		return fmt.Errorf("render callback failed: %w", err)
	}

	// End render pass
	renderPass.End()

	// End and submit command buffer
	if err := commandBuffer.End(); err != nil {
		return fmt.Errorf("failed to end command buffer: %w", err)
	}

	if err := r.context.GetCommandQueue().Submit(commandBuffer); err != nil {
		return fmt.Errorf("failed to submit command buffer: %w", err)
	}

	return nil
}

// GetContext returns the rendering context
func (r *RendererImpl) GetContext() Context {
	return r.context
}

// SetDebugLabel sets a debug label for the renderer
func (r *RendererImpl) SetDebugLabel(label string) {
	r.debugLabel = label
}

// GetStats returns rendering statistics
func (r *RendererImpl) GetStats() RendererStats {
	r.statsMutex.RLock()
	defer r.statsMutex.RUnlock()
	return r.stats
}

// Private helper methods

func (r *RendererImpl) createRenderPass(commandBuffer CommandBuffer, target RenderTarget) (RenderPass, error) {
	// TODO: Implement render pass creation based on target
	// This should create a render pass with appropriate color/depth attachments
	return nil, fmt.Errorf("render pass creation not implemented")
}

func (r *RendererImpl) renderEntities(renderPass RenderPass, entities []entity.Entity) error {
	for _, ent := range entities {
		if err := r.renderEntity(renderPass, ent); err != nil {
			return fmt.Errorf("failed to render entity: %w", err)
		}
	}
	return nil
}

func (r *RendererImpl) renderEntity(renderPass RenderPass, ent entity.Entity) error {
	// Get entity content
	content := ent.GetContents()
	if content == nil {
		return nil // Skip entities without content
	}

	// Apply entity transform
	transform := ent.GetTransform()
	renderPass.SetTransform(transform)

	// Set blend mode
	blendMode := ent.GetBlendMode()
	renderPass.SetBlendMode(blendMode)

	// Render the content
	return content.Render(renderPass, ent)
}

func (r *RendererImpl) getPipeline(key string) Pipeline {
	r.cacheMutex.RLock()
	pipeline, exists := r.pipelineCache[key]
	r.cacheMutex.RUnlock()

	if !exists {
		// TODO: Create pipeline based on key
		r.cacheMutex.Lock()
		// Double-check after acquiring write lock
		if pipeline, exists = r.pipelineCache[key]; !exists {
			// Create new pipeline
			pipeline = r.createPipeline(key)
			r.pipelineCache[key] = pipeline
		}
		r.cacheMutex.Unlock()
	}

	return pipeline
}

func (r *RendererImpl) createPipeline(key string) Pipeline {
	// TODO: Implement pipeline creation based on key
	return nil
}

// UpdateStats updates rendering statistics
func (r *RendererImpl) UpdateStats(deltaTime float64, drawCalls uint64, triangles uint64, vertices uint64) {
	if !r.enableStats {
		return
	}

	r.statsMutex.Lock()
	defer r.statsMutex.Unlock()

	r.stats.LastFrameTime = deltaTime
	r.stats.TotalRenderTime += deltaTime
	r.stats.DrawCallsIssued += drawCalls
	r.stats.TrianglesRendered += triangles
	r.stats.VerticesProcessed += vertices

	// Calculate average render time
	if r.stats.FramesRendered > 0 {
		r.stats.AverageRenderTime = r.stats.TotalRenderTime / float64(r.stats.FramesRendered)
	}

	// Calculate frame rate
	if deltaTime > 0 {
		r.stats.FrameRate = float32(1000.0 / deltaTime) // Convert ms to FPS
	}
}

// SetStatsEnabled enables or disables statistics collection
func (r *RendererImpl) SetStatsEnabled(enabled bool) {
	r.enableStats = enabled
}

// SetDebugEnabled enables or disables debug features
func (r *RendererImpl) SetDebugEnabled(enabled bool) {
	r.enableDebug = enabled
}

// ClearStats resets all statistics
func (r *RendererImpl) ClearStats() {
	r.statsMutex.Lock()
	defer r.statsMutex.Unlock()
	r.stats = RendererStats{}
}

// Cleanup cleans up renderer resources
func (r *RendererImpl) Cleanup() {
	// Clear pipeline cache
	r.cacheMutex.Lock()
	for _, pipeline := range r.pipelineCache {
		if pipeline != nil {
			// TODO: Destroy pipeline
		}
	}
	r.pipelineCache = make(map[string]Pipeline)
	r.cacheMutex.Unlock()

	// Clear resource pools
	r.commandBufferPool.Clear()
	r.renderPassPool.Clear()
}

// RenderEntityDirect provides direct entity rendering without pools
func (r *RendererImpl) RenderEntityDirect(renderPass RenderPass, entity entity.Entity) error {
	return r.renderEntity(renderPass, entity)
}

// BeginFrame marks the beginning of a frame for statistics
func (r *RendererImpl) BeginFrame() {
	// TODO: Implement frame timing
}

// EndFrame marks the end of a frame for statistics
func (r *RendererImpl) EndFrame() {
	// TODO: Implement frame timing and statistics update
}
