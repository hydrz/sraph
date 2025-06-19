package entity

import (
	"github.com/opensraph/sraph/entity/shaders"
	"github.com/opensraph/sraph/font"
	"github.com/opensraph/sraph/gpu"
	"github.com/opensraph/sraph/render"
)

// ContentContext provides the graphics context necessary to render entity contents.
// It manages resources related to rendering entities on the GPU, similar to
// Impeller's ContentContext.
type ContentContext struct {
	// gpuContext is the underlying GPU context.
	gpuContext gpu.Context

	// typographerContext provides text rendering capabilities.
	typographerContext font.TypographerContext

	// shaderLibrary manages shaders for content rendering.
	shaderLibrary *render.ShaderLibrary

	// pipelineLibrary manages rendering pipelines.
	pipelineLibrary *render.PipelineLibrary

	// samplerLibrary manages texture samplers.
	samplerLibrary *render.SamplerLibrary

	// renderTargetCache manages render target allocation and caching.
	renderTargetCache RenderTargetCache

	// hostBuffer provides CPU-accessible memory for uniform data.
	hostBuffer *render.HostBuffer

	// resourceAllocator manages GPU memory allocation.
	resourceAllocator *render.ResourceAllocator

	// entityShaderBundle contains all entity shaders.
	entityShaderBundle *shaders.EntityShaderBundle

	// glyphAtlas manages glyph texture atlas for text rendering.
	glyphAtlas *font.GlyphAtlas
}

// NewContentContext creates a new ContentContext with the provided GPU context and typographer.
func NewContentContext(gpuContext gpu.Context, typographerContext font.TypographerContext) *ContentContext {
	context := &ContentContext{
		gpuContext:         gpuContext,
		typographerContext: typographerContext,
		shaderLibrary:      render.NewShaderLibrary(gpuContext),
		pipelineLibrary:    render.NewPipelineLibrary(gpuContext),
		samplerLibrary:     render.NewSamplerLibrary(gpuContext),
		renderTargetCache:  NewRenderTargetCache(gpuContext.GetAllocator(), 4),
		hostBuffer:         render.NewHostBuffer(gpuContext),
		resourceAllocator:  render.NewResourceAllocator(gpuContext),
		entityShaderBundle: shaders.NewEntityShaderBundle(),
		glyphAtlas:         font.NewGlyphAtlas(gpuContext),
	}

	// Initialize shaders
	context.initializeShaders()

	return context
}

// IsValid returns true if the ContentContext is valid and ready for use.
func (c *ContentContext) IsValid() bool {
	return c.gpuContext != nil && c.gpuContext.IsValid()
}

// GetContext returns the underlying GPU context.
func (c *ContentContext) GetContext() gpu.Context {
	return c.gpuContext
}

// GetTypographerContext returns the typographer context for text rendering.
func (c *ContentContext) GetTypographerContext() font.TypographerContext {
	return c.typographerContext
}

// GetShaderLibrary returns the shader library.
func (c *ContentContext) GetShaderLibrary() *render.ShaderLibrary {
	return c.shaderLibrary
}

// GetPipelineLibrary returns the pipeline library.
func (c *ContentContext) GetPipelineLibrary() *render.PipelineLibrary {
	return c.pipelineLibrary
}

// GetSamplerLibrary returns the sampler library.
func (c *ContentContext) GetSamplerLibrary() *render.SamplerLibrary {
	return c.samplerLibrary
}

// GetRenderTargetCache returns the render target cache.
func (c *ContentContext) GetRenderTargetCache() RenderTargetCache {
	return c.renderTargetCache
}

// GetHostBuffer returns the host buffer for uniform data.
func (c *ContentContext) GetHostBuffer() *render.HostBuffer {
	return c.hostBuffer
}

// GetResourceAllocator returns the resource allocator.
func (c *ContentContext) GetResourceAllocator() *render.ResourceAllocator {
	return c.resourceAllocator
}

// GetGlyphAtlas returns the glyph atlas for text rendering.
func (c *ContentContext) GetGlyphAtlas() *font.GlyphAtlas {
	return c.glyphAtlas
}

// GetSolidColorShader returns the solid color shader.
func (c *ContentContext) GetSolidColorShader() *render.Shader {
	return c.shaderLibrary.GetShader("entity_solid_color")
}

// GetTextureShader returns the texture shader.
func (c *ContentContext) GetTextureShader() *render.Shader {
	return c.shaderLibrary.GetShader("entity_texture")
}

// GetLinearGradientShader returns the linear gradient shader.
func (c *ContentContext) GetLinearGradientShader() *render.Shader {
	return c.shaderLibrary.GetShader("entity_linear_gradient")
}

// GetTextShader returns the text shader.
func (c *ContentContext) GetTextShader() *render.Shader {
	return c.shaderLibrary.GetShader("entity_text")
}

// GetStencilShader returns the stencil shader.
func (c *ContentContext) GetStencilShader() *render.Shader {
	return c.shaderLibrary.GetShader("entity_stencil")
}

// CreateCommandBuffer creates a new command buffer.
func (c *ContentContext) CreateCommandBuffer() *render.CommandBuffer {
	return c.gpuContext.CreateCommandBuffer()
}

// CreateTexture creates a new texture with the given descriptor.
func (c *ContentContext) CreateTexture(descriptor *gpu.TextureDescriptor) gpu.Texture {
	return c.gpuContext.CreateTexture(descriptor)
}

// CreateSampler creates a new sampler with the given descriptor.
func (c *ContentContext) CreateSampler(descriptor *gpu.SamplerDescriptor) gpu.Sampler {
	return c.gpuContext.CreateSampler(descriptor)
}

// CreateBuffer creates a new buffer with the given descriptor.
func (c *ContentContext) CreateBuffer(descriptor *gpu.BufferDescriptor) gpu.Buffer {
	return c.gpuContext.CreateBuffer(descriptor)
}

// Flush submits all pending GPU work.
func (c *ContentContext) Flush() bool {
	return c.gpuContext.Flush()
}

// Shutdown cleanly shuts down the ContentContext and releases resources.
func (c *ContentContext) Shutdown() {
	if c.glyphAtlas != nil {
		c.glyphAtlas.Shutdown()
	}

	if c.hostBuffer != nil {
		c.hostBuffer.Shutdown()
	}

	if c.resourceAllocator != nil {
		c.resourceAllocator.Shutdown()
	}

	if c.gpuContext != nil {
		c.gpuContext.Shutdown()
	}
}

// initializeShaders loads and compiles all entity shaders.
func (c *ContentContext) initializeShaders() {
	// Load solid color shader
	c.shaderLibrary.CreateShader("entity_solid_color",
		c.entityShaderBundle.GetVertexShader(shaders.ShaderTypeSolidColor),
		c.entityShaderBundle.GetFragmentShader(shaders.ShaderTypeSolidColor))

	// Load texture shader
	c.shaderLibrary.CreateShader("entity_texture",
		c.entityShaderBundle.GetVertexShader(shaders.ShaderTypeTexture),
		c.entityShaderBundle.GetFragmentShader(shaders.ShaderTypeTexture))

	// Load linear gradient shader
	c.shaderLibrary.CreateShader("entity_linear_gradient",
		c.entityShaderBundle.GetVertexShader(shaders.ShaderTypeLinearGradient),
		c.entityShaderBundle.GetFragmentShader(shaders.ShaderTypeLinearGradient))

	// Load text shader
	c.shaderLibrary.CreateShader("entity_text",
		c.entityShaderBundle.GetVertexShader(shaders.ShaderTypeText),
		c.entityShaderBundle.GetFragmentShader(shaders.ShaderTypeText))

	// Load stencil shader
	c.shaderLibrary.CreateShader("entity_stencil",
		c.entityShaderBundle.GetVertexShader(shaders.ShaderTypeStencil),
		c.entityShaderBundle.GetFragmentShader(shaders.ShaderTypeStencil))
}
