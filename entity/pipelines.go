package entity

import (
	"github.com/opensraph/sraph/render"
)

// PipelineType represents different types of rendering pipelines
type PipelineType int

const (
	// PipelineTypeSolid for solid color rendering
	PipelineTypeSolid PipelineType = iota
	// PipelineTypeGradient for gradient rendering
	PipelineTypeGradient
	// PipelineTypeTexture for texture rendering
	PipelineTypeTexture
	// PipelineTypeText for text rendering
	PipelineTypeText
	// PipelineTypeBlur for blur effects
	PipelineTypeBlur
	// PipelineTypeRuntimeEffect for runtime effects
	PipelineTypeRuntimeEffect
	// PipelineTypeFramebufferBlend for framebuffer blending
	PipelineTypeFramebufferBlend
)

// PipelineDescriptor describes a rendering pipeline configuration
type PipelineDescriptor struct {
	Type             PipelineType
	VertexShader     *render.ShaderFunction
	FragmentShader   *render.ShaderFunction
	VertexDescriptor *render.VertexDescriptor
	BlendMode        render.BlendMode
	Culling          render.CullMode
	DepthTest        bool
	DepthWrite       bool
	Wireframe        bool
	Samples          int
	Label            string
}

// PipelineManager manages rendering pipelines for entities
type PipelineManager struct {
	library   render.PipelineLibrary
	pipelines map[PipelineType]render.Pipeline
	context   render.Context
}

// NewPipelineManager creates a new pipeline manager
func NewPipelineManager(context render.Context, library render.PipelineLibrary) *PipelineManager {
	return &PipelineManager{
		library:   library,
		pipelines: make(map[PipelineType]render.Pipeline),
		context:   context,
	}
}

// GetPipeline returns a pipeline for the specified type
func (p *PipelineManager) GetPipeline(pipelineType PipelineType) render.Pipeline {
	if pipeline, exists := p.pipelines[pipelineType]; exists {
		return pipeline
	}

	// Create and cache pipeline
	descriptor := p.createPipelineDescriptor(pipelineType)
	pipeline := p.library.CreatePipeline(descriptor)
	p.pipelines[pipelineType] = pipeline

	return pipeline
}

// createPipelineDescriptor creates a pipeline descriptor for the given type
func (p *PipelineManager) createPipelineDescriptor(pipelineType PipelineType) *PipelineDescriptor {
	switch pipelineType {
	case PipelineTypeSolid:
		return p.createSolidPipelineDescriptor()
	case PipelineTypeGradient:
		return p.createGradientPipelineDescriptor()
	case PipelineTypeTexture:
		return p.createTexturePipelineDescriptor()
	case PipelineTypeText:
		return p.createTextPipelineDescriptor()
	case PipelineTypeBlur:
		return p.createBlurPipelineDescriptor()
	case PipelineTypeRuntimeEffect:
		return p.createRuntimeEffectPipelineDescriptor()
	case PipelineTypeFramebufferBlend:
		return p.createFramebufferBlendPipelineDescriptor()
	default:
		return p.createSolidPipelineDescriptor()
	}
}

// createSolidPipelineDescriptor creates a solid color pipeline descriptor
func (p *PipelineManager) createSolidPipelineDescriptor() *PipelineDescriptor {
	return &PipelineDescriptor{
		Type:           PipelineTypeSolid,
		VertexShader:   p.getVertexShader("solid"),
		FragmentShader: p.getFragmentShader("solid"),
		BlendMode:      render.BlendModeSourceOver,
		Culling:        render.CullModeBack,
		DepthTest:      true,
		DepthWrite:     true,
		Wireframe:      false,
		Samples:        1,
		Label:          "Solid Color Pipeline",
	}
}

// createGradientPipelineDescriptor creates a gradient pipeline descriptor
func (p *PipelineManager) createGradientPipelineDescriptor() *PipelineDescriptor {
	return &PipelineDescriptor{
		Type:           PipelineTypeGradient,
		VertexShader:   p.getVertexShader("gradient"),
		FragmentShader: p.getFragmentShader("gradient"),
		BlendMode:      render.BlendModeSourceOver,
		Culling:        render.CullModeBack,
		DepthTest:      true,
		DepthWrite:     true,
		Wireframe:      false,
		Samples:        1,
		Label:          "Gradient Pipeline",
	}
}

// createTexturePipelineDescriptor creates a texture pipeline descriptor
func (p *PipelineManager) createTexturePipelineDescriptor() *PipelineDescriptor {
	return &PipelineDescriptor{
		Type:           PipelineTypeTexture,
		VertexShader:   p.getVertexShader("texture"),
		FragmentShader: p.getFragmentShader("texture"),
		BlendMode:      render.BlendModeSourceOver,
		Culling:        render.CullModeBack,
		DepthTest:      true,
		DepthWrite:     true,
		Wireframe:      false,
		Samples:        1,
		Label:          "Texture Pipeline",
	}
}

// createTextPipelineDescriptor creates a text pipeline descriptor
func (p *PipelineManager) createTextPipelineDescriptor() *PipelineDescriptor {
	return &PipelineDescriptor{
		Type:           PipelineTypeText,
		VertexShader:   p.getVertexShader("text"),
		FragmentShader: p.getFragmentShader("text"),
		BlendMode:      render.BlendModeSourceOver,
		Culling:        render.CullModeNone,
		DepthTest:      true,
		DepthWrite:     false,
		Wireframe:      false,
		Samples:        4, // Multi-sampling for text
		Label:          "Text Pipeline",
	}
}

// createBlurPipelineDescriptor creates a blur pipeline descriptor
func (p *PipelineManager) createBlurPipelineDescriptor() *PipelineDescriptor {
	return &PipelineDescriptor{
		Type:           PipelineTypeBlur,
		VertexShader:   p.getVertexShader("blur"),
		FragmentShader: p.getFragmentShader("blur"),
		BlendMode:      render.BlendModeSourceOver,
		Culling:        render.CullModeNone,
		DepthTest:      false,
		DepthWrite:     false,
		Wireframe:      false,
		Samples:        1,
		Label:          "Blur Pipeline",
	}
}

// createRuntimeEffectPipelineDescriptor creates a runtime effect pipeline descriptor
func (p *PipelineManager) createRuntimeEffectPipelineDescriptor() *PipelineDescriptor {
	return &PipelineDescriptor{
		Type:           PipelineTypeRuntimeEffect,
		VertexShader:   p.getVertexShader("runtime_effect"),
		FragmentShader: p.getFragmentShader("runtime_effect"),
		BlendMode:      render.BlendModeSourceOver,
		Culling:        render.CullModeBack,
		DepthTest:      true,
		DepthWrite:     true,
		Wireframe:      false,
		Samples:        1,
		Label:          "Runtime Effect Pipeline",
	}
}

// createFramebufferBlendPipelineDescriptor creates a framebuffer blend pipeline descriptor
func (p *PipelineManager) createFramebufferBlendPipelineDescriptor() *PipelineDescriptor {
	return &PipelineDescriptor{
		Type:           PipelineTypeFramebufferBlend,
		VertexShader:   p.getVertexShader("framebuffer_blend"),
		FragmentShader: p.getFragmentShader("framebuffer_blend"),
		BlendMode:      render.BlendModeSourceOver,
		Culling:        render.CullModeNone,
		DepthTest:      false,
		DepthWrite:     false,
		Wireframe:      false,
		Samples:        1,
		Label:          "Framebuffer Blend Pipeline",
	}
}

// getVertexShader returns a vertex shader for the given name
func (p *PipelineManager) getVertexShader(name string) *render.ShaderFunction {
	// TODO: Implement shader loading from library
	return nil
}

// getFragmentShader returns a fragment shader for the given name
func (p *PipelineManager) getFragmentShader(name string) *render.ShaderFunction {
	// TODO: Implement shader loading from library
	return nil
}

// ClearPipelines clears all cached pipelines
func (p *PipelineManager) ClearPipelines() {
	p.pipelines = make(map[PipelineType]render.Pipeline)
}

// GetPipelineCount returns the number of cached pipelines
func (p *PipelineManager) GetPipelineCount() int {
	return len(p.pipelines)
}

// PreloadPipelines preloads all common pipelines
func (p *PipelineManager) PreloadPipelines() {
	commonPipelines := []PipelineType{
		PipelineTypeSolid,
		PipelineTypeGradient,
		PipelineTypeTexture,
		PipelineTypeText,
	}

	for _, pipelineType := range commonPipelines {
		p.GetPipeline(pipelineType)
	}
}

// PipelineKey represents a unique key for pipeline variants
type PipelineKey struct {
	Type        PipelineType
	BlendMode   render.BlendMode
	Samples     int
	HasTexture  bool
	HasGradient bool
	HasText     bool
}

// PipelineCache provides caching for pipeline variants
type PipelineCache struct {
	cache   map[PipelineKey]render.Pipeline
	manager *PipelineManager
}

// NewPipelineCache creates a new pipeline cache
func NewPipelineCache(manager *PipelineManager) *PipelineCache {
	return &PipelineCache{
		cache:   make(map[PipelineKey]render.Pipeline),
		manager: manager,
	}
}

// GetPipelineForKey returns a pipeline for the specified key
func (c *PipelineCache) GetPipelineForKey(key PipelineKey) render.Pipeline {
	if pipeline, exists := c.cache[key]; exists {
		return pipeline
	}

	// Create pipeline variant
	pipeline := c.createPipelineVariant(key)
	c.cache[key] = pipeline

	return pipeline
}

// createPipelineVariant creates a pipeline variant for the given key
func (c *PipelineCache) createPipelineVariant(key PipelineKey) render.Pipeline {
	// TODO: Implement pipeline variant creation
	return c.manager.GetPipeline(key.Type)
}

// Clear clears the pipeline cache
func (c *PipelineCache) Clear() {
	c.cache = make(map[PipelineKey]render.Pipeline)
}
