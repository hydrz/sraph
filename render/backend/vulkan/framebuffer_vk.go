package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// FramebufferVK wraps Vulkan framebuffer
// It represents a collection of attachments for rendering
type FramebufferVK struct {
	device      vulkan.Device
	framebuffer vulkan.Framebuffer
	renderPass  vulkan.RenderPass
	imageViews  []vulkan.ImageView
	width       uint32
	height      uint32
	layers      uint32
}

// FramebufferCreateInfo contains framebuffer creation parameters
type FramebufferCreateInfo struct {
	RenderPass vulkan.RenderPass
	ImageViews []vulkan.ImageView
	Width      uint32
	Height     uint32
	Layers     uint32
}

// NewFramebufferVK creates a new Vulkan framebuffer
func NewFramebufferVK(device vulkan.Device, createInfo *FramebufferCreateInfo) (*FramebufferVK, error) {
	if createInfo == nil {
		return nil, fmt.Errorf("framebuffer create info cannot be nil")
	}

	vkCreateInfo := vulkan.FramebufferCreateInfo{
		SType:           vulkan.StructureTypeFramebufferCreateInfo,
		RenderPass:      createInfo.RenderPass,
		AttachmentCount: uint32(len(createInfo.ImageViews)),
		Width:           createInfo.Width,
		Height:          createInfo.Height,
		Layers:          createInfo.Layers,
	}

	if len(createInfo.ImageViews) > 0 {
		vkCreateInfo.PAttachments = createInfo.ImageViews
	}

	var framebuffer vulkan.Framebuffer
	if result := vulkan.CreateFramebuffer(device, &vkCreateInfo, nil, &framebuffer); result != vulkan.Success {
		return nil, fmt.Errorf("failed to create framebuffer: %s", result)
	}

	return &FramebufferVK{
		device:      device,
		framebuffer: framebuffer,
		renderPass:  createInfo.RenderPass,
		imageViews:  createInfo.ImageViews,
		width:       createInfo.Width,
		height:      createInfo.Height,
		layers:      createInfo.Layers,
	}, nil
}

// GetFramebuffer returns the Vulkan framebuffer handle
func (fb *FramebufferVK) GetFramebuffer() vulkan.Framebuffer {
	return fb.framebuffer
}

// GetRenderPass returns the associated render pass
func (fb *FramebufferVK) GetRenderPass() vulkan.RenderPass {
	return fb.renderPass
}

// GetImageViews returns the image views
func (fb *FramebufferVK) GetImageViews() []vulkan.ImageView {
	return fb.imageViews
}

// GetWidth returns the framebuffer width
func (fb *FramebufferVK) GetWidth() uint32 {
	return fb.width
}

// GetHeight returns the framebuffer height
func (fb *FramebufferVK) GetHeight() uint32 {
	return fb.height
}

// GetLayers returns the framebuffer layer count
func (fb *FramebufferVK) GetLayers() uint32 {
	return fb.layers
}

// GetSize returns the framebuffer dimensions
func (fb *FramebufferVK) GetSize() (uint32, uint32) {
	return fb.width, fb.height
}

// Destroy destroys the framebuffer
func (fb *FramebufferVK) Destroy() {
	if fb.framebuffer != vulkan.NullFramebuffer {
		vulkan.DestroyFramebuffer(fb.device, fb.framebuffer, nil)
		fb.framebuffer = vulkan.NullFramebuffer
	}
}

// FramebufferCacheVK manages framebuffer caching
// It provides efficient creation and reuse of framebuffers
type FramebufferCacheVK struct {
	device       vulkan.Device
	framebuffers map[string]*FramebufferVK
}

// NewFramebufferCacheVK creates a new framebuffer cache
func NewFramebufferCacheVK(device vulkan.Device) *FramebufferCacheVK {
	return &FramebufferCacheVK{
		device:       device,
		framebuffers: make(map[string]*FramebufferVK),
	}
}

// GetFramebuffer gets or creates a framebuffer
func (fc *FramebufferCacheVK) GetFramebuffer(createInfo *FramebufferCreateInfo) (*FramebufferVK, error) {
	key := fc.generateKey(createInfo)

	if fb, exists := fc.framebuffers[key]; exists {
		return fb, nil
	}

	fb, err := NewFramebufferVK(fc.device, createInfo)
	if err != nil {
		return nil, err
	}

	fc.framebuffers[key] = fb
	return fb, nil
}

// generateKey generates a cache key for framebuffer
func (fc *FramebufferCacheVK) generateKey(createInfo *FramebufferCreateInfo) string {
	// TODO: Implement proper key generation based on create info
	return fmt.Sprintf("fb_%p_%dx%dx%d_%d",
		createInfo.RenderPass,
		createInfo.Width,
		createInfo.Height,
		createInfo.Layers,
		len(createInfo.ImageViews))
}

// Clear clears the framebuffer cache
func (fc *FramebufferCacheVK) Clear() {
	for _, fb := range fc.framebuffers {
		fb.Destroy()
	}
	fc.framebuffers = make(map[string]*FramebufferVK)
}

// GetCacheSize returns the number of cached framebuffers
func (fc *FramebufferCacheVK) GetCacheSize() int {
	return len(fc.framebuffers)
}

// Destroy destroys the framebuffer cache
func (fc *FramebufferCacheVK) Destroy() {
	fc.Clear()
}

// FramebufferBuilderVK helps build framebuffers
// It provides a builder pattern for framebuffer creation
type FramebufferBuilderVK struct {
	device     vulkan.Device
	renderPass vulkan.RenderPass
	imageViews []vulkan.ImageView
	width      uint32
	height     uint32
	layers     uint32
}

// NewFramebufferBuilderVK creates a new framebuffer builder
func NewFramebufferBuilderVK(device vulkan.Device) *FramebufferBuilderVK {
	return &FramebufferBuilderVK{
		device: device,
		layers: 1, // Default to single layer
	}
}

// SetRenderPass sets the render pass
func (fb *FramebufferBuilderVK) SetRenderPass(renderPass vulkan.RenderPass) *FramebufferBuilderVK {
	fb.renderPass = renderPass
	return fb
}

// AddImageView adds an image view attachment
func (fb *FramebufferBuilderVK) AddImageView(imageView vulkan.ImageView) *FramebufferBuilderVK {
	fb.imageViews = append(fb.imageViews, imageView)
	return fb
}

// SetImageViews sets all image view attachments
func (fb *FramebufferBuilderVK) SetImageViews(imageViews []vulkan.ImageView) *FramebufferBuilderVK {
	fb.imageViews = imageViews
	return fb
}

// SetSize sets the framebuffer dimensions
func (fb *FramebufferBuilderVK) SetSize(width, height uint32) *FramebufferBuilderVK {
	fb.width = width
	fb.height = height
	return fb
}

// SetLayers sets the number of layers
func (fb *FramebufferBuilderVK) SetLayers(layers uint32) *FramebufferBuilderVK {
	fb.layers = layers
	return fb
}

// Build builds the framebuffer
func (fb *FramebufferBuilderVK) Build() (*FramebufferVK, error) {
	if fb.renderPass == vulkan.NullRenderPass {
		return nil, fmt.Errorf("render pass must be set")
	}

	if fb.width == 0 || fb.height == 0 {
		return nil, fmt.Errorf("framebuffer dimensions must be greater than 0")
	}

	createInfo := &FramebufferCreateInfo{
		RenderPass: fb.renderPass,
		ImageViews: fb.imageViews,
		Width:      fb.width,
		Height:     fb.height,
		Layers:     fb.layers,
	}

	return NewFramebufferVK(fb.device, createInfo)
}

// Reset resets the builder to its initial state
func (fb *FramebufferBuilderVK) Reset() *FramebufferBuilderVK {
	fb.renderPass = vulkan.NullRenderPass
	fb.imageViews = fb.imageViews[:0]
	fb.width = 0
	fb.height = 0
	fb.layers = 1
	return fb
}

// AttachmentInfo contains attachment information for framebuffer
type AttachmentInfo struct {
	ImageView vulkan.ImageView
	Format    vulkan.Format
	Samples   vulkan.SampleCountFlags
	LoadOp    vulkan.AttachmentLoadOp
	StoreOp   vulkan.AttachmentStoreOp
	Layout    vulkan.ImageLayout
}

// FramebufferManagerVK manages framebuffer lifecycle
// It provides high-level framebuffer management
type FramebufferManagerVK struct {
	device             vulkan.Device
	cache              *FramebufferCacheVK
	builder            *FramebufferBuilderVK
	activeFramebuffers []*FramebufferVK
}

// NewFramebufferManagerVK creates a new framebuffer manager
func NewFramebufferManagerVK(device vulkan.Device) *FramebufferManagerVK {
	return &FramebufferManagerVK{
		device:             device,
		cache:              NewFramebufferCacheVK(device),
		builder:            NewFramebufferBuilderVK(device),
		activeFramebuffers: make([]*FramebufferVK, 0),
	}
}

// CreateFramebuffer creates a new framebuffer
func (fm *FramebufferManagerVK) CreateFramebuffer(createInfo *FramebufferCreateInfo) (*FramebufferVK, error) {
	fb, err := fm.cache.GetFramebuffer(createInfo)
	if err != nil {
		return nil, err
	}

	fm.activeFramebuffers = append(fm.activeFramebuffers, fb)
	return fb, nil
}

// GetBuilder returns the framebuffer builder
func (fm *FramebufferManagerVK) GetBuilder() *FramebufferBuilderVK {
	return fm.builder.Reset()
}

// GetActiveFramebufferCount returns the number of active framebuffers
func (fm *FramebufferManagerVK) GetActiveFramebufferCount() int {
	return len(fm.activeFramebuffers)
}

// GetCacheSize returns the cache size
func (fm *FramebufferManagerVK) GetCacheSize() int {
	return fm.cache.GetCacheSize()
}

// Cleanup removes unused framebuffers
func (fm *FramebufferManagerVK) Cleanup() {
	// TODO: Implement reference counting and cleanup logic
}

// Destroy destroys the framebuffer manager
func (fm *FramebufferManagerVK) Destroy() {
	fm.cache.Destroy()
	fm.activeFramebuffers = fm.activeFramebuffers[:0]
}
