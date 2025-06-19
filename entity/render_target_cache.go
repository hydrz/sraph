package entity

import (
	"github.com/opensraph/sraph/gpu"
	"github.com/opensraph/sraph/render"
)

// RenderTargetCache caches render targets for efficient reuse across frames.
// It implements the RenderTargetAllocator interface and manages texture lifecycle.
type RenderTargetCache interface {
	// Start begins a new frame, marking all cached targets as unused
	Start()

	// End completes the current frame and releases unused targets
	End()

	// DisableCache temporarily disables caching
	DisableCache()

	// EnableCache re-enables caching
	EnableCache()

	// CreateOffscreen creates or reuses an offscreen render target
	CreateOffscreen(
		context gpu.Context,
		size gpu.Size,
		mipCount int,
		label string,
		colorAttachmentConfig *render.AttachmentConfig,
		stencilAttachmentConfig *render.AttachmentConfig,
		existingColorTexture gpu.Texture,
		existingDepthStencilTexture gpu.Texture,
	) *render.RenderTarget

	// CreateOffscreenMSAA creates or reuses an MSAA offscreen render target
	CreateOffscreenMSAA(
		context gpu.Context,
		size gpu.Size,
		mipCount int,
		label string,
		colorAttachmentConfig *render.AttachmentConfigMSAA,
		stencilAttachmentConfig *render.AttachmentConfig,
		existingColorMSAATexture gpu.Texture,
		existingColorResolveTexture gpu.Texture,
		existingDepthStencilTexture gpu.Texture,
	) *render.RenderTarget

	// CachedTextureCount returns the number of cached textures
	CachedTextureCount() int
}

// RenderTargetData holds cached render target information
type RenderTargetData struct {
	UsedThisFrame       bool                       // Whether this target was used in the current frame
	KeepAliveFrameCount uint32                     // Number of frames to keep alive
	Config              *render.RenderTargetConfig // Configuration used to create this target
	RenderTarget        *render.RenderTarget       // The cached render target
}

// DefaultRenderTargetCache provides a default implementation of RenderTargetCache
type DefaultRenderTargetCache struct {
	allocator           gpu.Allocator       // GPU memory allocator
	renderTargetData    []*RenderTargetData // Cached render target data
	keepAliveFrameCount uint32              // Default keep alive frame count
	cacheDisabledCount  uint32              // Counter for cache disable requests
}

// NewRenderTargetCache creates a new render target cache
func NewRenderTargetCache(allocator gpu.Allocator, keepAliveFrameCount uint32) RenderTargetCache {
	if keepAliveFrameCount == 0 {
		keepAliveFrameCount = 4 // Default value from Impeller
	}

	return &DefaultRenderTargetCache{
		allocator:           allocator,
		renderTargetData:    make([]*RenderTargetData, 0),
		keepAliveFrameCount: keepAliveFrameCount,
		cacheDisabledCount:  0,
	}
}

// Start begins a new frame, marking all cached targets as unused
func (c *DefaultRenderTargetCache) Start() {
	for _, data := range c.renderTargetData {
		data.UsedThisFrame = false
	}
}

// End completes the current frame and releases unused targets
func (c *DefaultRenderTargetCache) End() {
	// Remove targets that haven't been used and have expired
	filtered := make([]*RenderTargetData, 0, len(c.renderTargetData))

	for _, data := range c.renderTargetData {
		if data.UsedThisFrame {
			// Reset keep alive counter for used targets
			data.KeepAliveFrameCount = c.keepAliveFrameCount
			filtered = append(filtered, data)
		} else {
			// Decrement keep alive counter for unused targets
			if data.KeepAliveFrameCount > 0 {
				data.KeepAliveFrameCount--
				if data.KeepAliveFrameCount > 0 {
					filtered = append(filtered, data)
				}
				// If KeepAliveFrameCount reaches 0, the target will be discarded
			}
		}
	}

	c.renderTargetData = filtered
}

// DisableCache temporarily disables caching
func (c *DefaultRenderTargetCache) DisableCache() {
	c.cacheDisabledCount++
}

// EnableCache re-enables caching
func (c *DefaultRenderTargetCache) EnableCache() {
	if c.cacheDisabledCount > 0 {
		c.cacheDisabledCount--
	}
}

// CacheEnabled returns true if caching is currently enabled
func (c *DefaultRenderTargetCache) CacheEnabled() bool {
	return c.cacheDisabledCount == 0
}

// CreateOffscreen creates or reuses an offscreen render target
func (c *DefaultRenderTargetCache) CreateOffscreen(
	context gpu.Context,
	size gpu.Size,
	mipCount int,
	label string,
	colorAttachmentConfig *render.AttachmentConfig,
	stencilAttachmentConfig *render.AttachmentConfig,
	existingColorTexture gpu.Texture,
	existingDepthStencilTexture gpu.Texture,
) *render.RenderTarget {

	if !c.CacheEnabled() {
		// Cache disabled, create new target directly
		return c.createOffscreenDirect(
			context, size, mipCount, label,
			colorAttachmentConfig, stencilAttachmentConfig,
			existingColorTexture, existingDepthStencilTexture,
		)
	}

	// Create configuration for cache lookup
	config := &render.RenderTargetConfig{
		Size:                           size,
		MipCount:                       mipCount,
		ColorAttachmentConfig:          colorAttachmentConfig,
		StencilAttachmentConfig:        stencilAttachmentConfig,
		HasExistingColorTexture:        existingColorTexture != nil,
		HasExistingDepthStencilTexture: existingDepthStencilTexture != nil,
	}

	// Look for matching cached target
	for _, data := range c.renderTargetData {
		if c.configMatches(data.Config, config) {
			data.UsedThisFrame = true
			return data.RenderTarget
		}
	}

	// No matching target found, create new one
	target := c.createOffscreenDirect(
		context, size, mipCount, label,
		colorAttachmentConfig, stencilAttachmentConfig,
		existingColorTexture, existingDepthStencilTexture,
	)

	if target != nil {
		// Add to cache
		data := &RenderTargetData{
			UsedThisFrame:       true,
			KeepAliveFrameCount: c.keepAliveFrameCount,
			Config:              config,
			RenderTarget:        target,
		}
		c.renderTargetData = append(c.renderTargetData, data)
	}

	return target
}

// CreateOffscreenMSAA creates or reuses an MSAA offscreen render target
func (c *DefaultRenderTargetCache) CreateOffscreenMSAA(
	context gpu.Context,
	size gpu.Size,
	mipCount int,
	label string,
	colorAttachmentConfig *render.AttachmentConfigMSAA,
	stencilAttachmentConfig *render.AttachmentConfig,
	existingColorMSAATexture gpu.Texture,
	existingColorResolveTexture gpu.Texture,
	existingDepthStencilTexture gpu.Texture,
) *render.RenderTarget {

	// TODO: Implement MSAA target creation and caching
	// For now, create a regular offscreen target
	return c.CreateOffscreen(
		context, size, mipCount, label,
		nil, // Convert MSAA config to regular config
		stencilAttachmentConfig,
		existingColorResolveTexture, // Use resolve texture as color texture
		existingDepthStencilTexture,
	)
}

// CachedTextureCount returns the number of cached textures
func (c *DefaultRenderTargetCache) CachedTextureCount() int {
	return len(c.renderTargetData)
}

// createOffscreenDirect creates a new offscreen render target without caching
func (c *DefaultRenderTargetCache) createOffscreenDirect(
	context gpu.Context,
	size gpu.Size,
	mipCount int,
	label string,
	colorAttachmentConfig *render.AttachmentConfig,
	stencilAttachmentConfig *render.AttachmentConfig,
	existingColorTexture gpu.Texture,
	existingDepthStencilTexture gpu.Texture,
) *render.RenderTarget {

	// TODO: Implement actual render target creation
	// This is a placeholder implementation
	target := render.NewRenderTarget()
	if target == nil {
		return nil
	}

	// Set up color attachment
	if existingColorTexture != nil {
		target.SetColorTexture(existingColorTexture)
	} else {
		// Create new color texture
		colorTexture := c.createColorTexture(context, size, mipCount, colorAttachmentConfig)
		if colorTexture != nil {
			target.SetColorTexture(colorTexture)
		}
	}

	// Set up depth/stencil attachment
	if existingDepthStencilTexture != nil {
		target.SetDepthStencilTexture(existingDepthStencilTexture)
	} else if stencilAttachmentConfig != nil {
		// Create new depth/stencil texture
		depthStencilTexture := c.createDepthStencilTexture(context, size, stencilAttachmentConfig)
		if depthStencilTexture != nil {
			target.SetDepthStencilTexture(depthStencilTexture)
		}
	}

	return target
}

// createColorTexture creates a new color texture
func (c *DefaultRenderTargetCache) createColorTexture(
	context gpu.Context,
	size gpu.Size,
	mipCount int,
	config *render.AttachmentConfig,
) gpu.Texture {

	// TODO: Implement texture creation using allocator
	// This is a placeholder
	return nil
}

// createDepthStencilTexture creates a new depth/stencil texture
func (c *DefaultRenderTargetCache) createDepthStencilTexture(
	context gpu.Context,
	size gpu.Size,
	config *render.AttachmentConfig,
) gpu.Texture {

	// TODO: Implement depth/stencil texture creation
	// This is a placeholder
	return nil
}

// configMatches checks if two render target configurations match
func (c *DefaultRenderTargetCache) configMatches(a, b *render.RenderTargetConfig) bool {
	if a == nil || b == nil {
		return a == b
	}

	return a.Size.Width == b.Size.Width &&
		a.Size.Height == b.Size.Height &&
		a.MipCount == b.MipCount &&
		a.HasExistingColorTexture == b.HasExistingColorTexture &&
		a.HasExistingDepthStencilTexture == b.HasExistingDepthStencilTexture
	// TODO: Add more detailed config comparison
}
