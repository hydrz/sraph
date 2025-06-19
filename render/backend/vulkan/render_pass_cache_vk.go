package vulkan

import (
	"fmt"
	"sync"

	"github.com/vulkan-go/vulkan"
)

// RenderPassCacheVK manages a cache of Vulkan render passes
type RenderPassCacheVK struct {
	context      *ContextVK
	cache        map[RenderPassKey]*CachedRenderPass
	mutex        sync.RWMutex
	maxCacheSize int
}

// RenderPassKey uniquely identifies a render pass configuration
type RenderPassKey struct {
	ColorFormats   []vulkan.Format
	DepthFormat    vulkan.Format
	Samples        vulkan.SampleCountFlagBits
	LoadOps        []vulkan.AttachmentLoadOp
	StoreOps       []vulkan.AttachmentStoreOp
	InitialLayouts []vulkan.ImageLayout
	FinalLayouts   []vulkan.ImageLayout
}

// CachedRenderPass represents a cached render pass with usage tracking
type CachedRenderPass struct {
	RenderPass *RenderPassVK
	Key        RenderPassKey
	RefCount   int
	LastUsed   int64
}

// NewRenderPassCacheVK creates a new render pass cache
func NewRenderPassCacheVK(context *ContextVK, maxCacheSize int) *RenderPassCacheVK {
	return &RenderPassCacheVK{
		context:      context,
		cache:        make(map[RenderPassKey]*CachedRenderPass),
		maxCacheSize: maxCacheSize,
	}
}

// GetRenderPass retrieves or creates a render pass for the given configuration
func (cache *RenderPassCacheVK) GetRenderPass(descriptor *RenderPassCacheDescriptor) (*RenderPassVK, error) {
	key := cache.createKey(descriptor)

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Check if render pass exists in cache
	if cached, exists := cache.cache[key]; exists {
		cached.RefCount++
		// TODO: Update LastUsed timestamp
		return cached.RenderPass, nil
	}

	// Create new render pass
	renderPass, err := cache.createRenderPass(descriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to create render pass: %w", err)
	}

	// Add to cache
	cached := &CachedRenderPass{
		RenderPass: renderPass,
		Key:        key,
		RefCount:   1,
		// TODO: Set LastUsed timestamp
	}

	cache.cache[key] = cached

	// Cleanup cache if necessary
	cache.cleanupCache()

	return renderPass, nil
}

// ReleaseRenderPass decreases the reference count for a render pass
func (cache *RenderPassCacheVK) ReleaseRenderPass(renderPass *RenderPassVK) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Find the cached render pass
	for _, cached := range cache.cache {
		if cached.RenderPass == renderPass {
			cached.RefCount--
			if cached.RefCount <= 0 {
				// Remove from cache and destroy
				delete(cache.cache, cached.Key)
				cached.RenderPass.Destroy()
			}
			return
		}
	}
}

// createKey creates a cache key from the descriptor
func (cache *RenderPassCacheVK) createKey(descriptor *RenderPassCacheDescriptor) RenderPassKey {
	key := RenderPassKey{
		ColorFormats:   make([]vulkan.Format, len(descriptor.ColorAttachments)),
		LoadOps:        make([]vulkan.AttachmentLoadOp, len(descriptor.ColorAttachments)),
		StoreOps:       make([]vulkan.AttachmentStoreOp, len(descriptor.ColorAttachments)),
		InitialLayouts: make([]vulkan.ImageLayout, len(descriptor.ColorAttachments)),
		FinalLayouts:   make([]vulkan.ImageLayout, len(descriptor.ColorAttachments)),
		Samples:        descriptor.Samples,
	}

	// Copy color attachment information
	for i, attachment := range descriptor.ColorAttachments {
		key.ColorFormats[i] = attachment.Format
		key.LoadOps[i] = attachment.LoadOp
		key.StoreOps[i] = attachment.StoreOp
		key.InitialLayouts[i] = attachment.InitialLayout
		key.FinalLayouts[i] = attachment.FinalLayout
	}

	// Set depth attachment information
	if descriptor.DepthAttachment != nil {
		key.DepthFormat = descriptor.DepthAttachment.Format
	}

	return key
}

// createRenderPass creates a new render pass from the descriptor
func (cache *RenderPassCacheVK) createRenderPass(descriptor *RenderPassCacheDescriptor) (*RenderPassVK, error) {
	builder := NewRenderPassBuilderVK(cache.context)
	builder.SetExtent(descriptor.Extent.Width, descriptor.Extent.Height)

	if descriptor.Samples != vulkan.SampleCount1Bit {
		builder.SetMultisampling(descriptor.Samples)
	}

	// Add color attachments
	for _, attachment := range descriptor.ColorAttachments {
		builder.AddColorAttachment(attachment.Format, attachment.LoadOp, attachment.StoreOp)
	}

	// Add depth attachment if present
	if descriptor.DepthAttachment != nil {
		builder.AddDepthAttachment(descriptor.DepthAttachment.Format,
			descriptor.DepthAttachment.LoadOp, descriptor.DepthAttachment.StoreOp)
	}

	return builder.Build()
}

// cleanupCache removes least recently used items if cache is full
func (cache *RenderPassCacheVK) cleanupCache() {
	if len(cache.cache) <= cache.maxCacheSize {
		return
	}

	// TODO: Implement LRU cleanup based on LastUsed timestamp
	// For now, just remove items with zero reference count
	for key, cached := range cache.cache {
		if cached.RefCount <= 0 {
			delete(cache.cache, key)
			cached.RenderPass.Destroy()
			if len(cache.cache) <= cache.maxCacheSize {
				break
			}
		}
	}
}

// Clear removes all render passes from the cache
func (cache *RenderPassCacheVK) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	for key, cached := range cache.cache {
		cached.RenderPass.Destroy()
		delete(cache.cache, key)
	}
}

// GetCacheSize returns the current number of cached render passes
func (cache *RenderPassCacheVK) GetCacheSize() int {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	return len(cache.cache)
}

// GetMaxCacheSize returns the maximum cache size
func (cache *RenderPassCacheVK) GetMaxCacheSize() int {
	return cache.maxCacheSize
}

// SetMaxCacheSize sets the maximum cache size
func (cache *RenderPassCacheVK) SetMaxCacheSize(size int) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.maxCacheSize = size
	cache.cleanupCache()
}

// Destroy destroys the cache and all cached render passes
func (cache *RenderPassCacheVK) Destroy() {
	cache.Clear()
}

// RenderPassCacheDescriptor describes a render pass for caching
type RenderPassCacheDescriptor struct {
	ColorAttachments []ColorAttachmentDescriptor
	DepthAttachment  *DepthAttachmentDescriptor
	Extent           vulkan.Extent2D
	Samples          vulkan.SampleCountFlagBits
}

// ColorAttachmentDescriptor describes a color attachment
type ColorAttachmentDescriptor struct {
	Format        vulkan.Format
	LoadOp        vulkan.AttachmentLoadOp
	StoreOp       vulkan.AttachmentStoreOp
	InitialLayout vulkan.ImageLayout
	FinalLayout   vulkan.ImageLayout
}

// DepthAttachmentDescriptor describes a depth attachment
type DepthAttachmentDescriptor struct {
	Format        vulkan.Format
	LoadOp        vulkan.AttachmentLoadOp
	StoreOp       vulkan.AttachmentStoreOp
	InitialLayout vulkan.ImageLayout
	FinalLayout   vulkan.ImageLayout
}
