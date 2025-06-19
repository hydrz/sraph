package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// SurfaceContextVK manages Vulkan surface and swapchain
type SurfaceContextVK struct {
	context          *ContextVK
	surface          vulkan.Surface
	swapchain        vulkan.Swapchain
	swapchainImages  []vulkan.Image
	swapchainImageViews []vulkan.ImageView
	surfaceFormat    vulkan.SurfaceFormat
	presentMode      vulkan.PresentMode
	extent           vulkan.Extent2D
	imageCount       uint32
	currentImageIndex uint32
	renderPass       *RenderPassVK
	framebuffers     []vulkan.Framebuffer
}

// SurfaceConfig contains surface configuration parameters
type SurfaceConfig struct {
	PreferredFormat vulkan.Format
	PreferredColorSpace vulkan.ColorSpace
	PreferredPresentMode vulkan.PresentMode
	Width           uint32
	Height          uint32
	VSync           bool
}

// NewSurfaceContextVK creates a new Vulkan surface context
func NewSurfaceContextVK(context *ContextVK, surface vulkan.Surface, config *SurfaceConfig) (*SurfaceContextVK, error) {
	if context == nil || surface == vulkan.NullSurface {
		return nil, fmt.Errorf("invalid parameters")
	}

	sc := &SurfaceContextVK{
		context: context,
		surface: surface,
	}

	// Query surface capabilities and formats
	if err := sc.querySurfaceSupport(); err != nil {
		return nil, fmt.Errorf("failed to query surface support: %w", err)
	}

	// Choose surface format
	if err := sc.chooseSurfaceFormat(config); err != nil {
		return nil, fmt.Errorf("failed to choose surface format: %w", err)
	}

	// Choose present mode
	if err := sc.choosePresentMode(config); err != nil {
		return nil, fmt.Errorf("failed to choose present mode: %w", err)
	}

	// Create swapchain
	if err := sc.createSwapchain(config); err != nil {
		return nil, fmt.Errorf("failed to create swapchain: %w", err)
	}

	// Create image views
	if err := sc.createImageViews(); err != nil {
		sc.Destroy()
		return nil, fmt.Errorf("failed to create image views: %w", err)
	}

	return sc, nil
}

// querySurfaceSupport queries surface capabilities and formats
func (sc *SurfaceContextVK) querySurfaceSupport() error {
	// Check surface support (this is a simplified version)
	var supported vulkan.Bool32
	if result := vulkan.GetPhysicalDeviceSurfaceSupport(sc.context.physicalDevice, 
		sc.context.graphicsQueueFamilyIndex, sc.surface, &supported); result != vulkan.Success {
		return fmt.Errorf("failed to check surface support: %s", result)
	}

	if supported == vulkan.False {
		return fmt.Errorf("surface not supported")
	}

	return nil
}

// chooseSurfaceFormat chooses the best surface format
func (sc *SurfaceContextVK) chooseSurfaceFormat(config *SurfaceConfig) error {
	// Query surface formats
	var formatCount uint32
	vulkan.GetPhysicalDeviceSurfaceFormats(sc.context.physicalDevice, sc.surface, &formatCount, nil)
	
	if formatCount == 0 {
		return fmt.Errorf("no surface formats available")
	}

	formats := make([]vulkan.SurfaceFormat, formatCount)
	vulkan.GetPhysicalDeviceSurfaceFormats(sc.context.physicalDevice, sc.surface, &formatCount, formats)

	// Choose preferred format or fallback to first available
	for _, format := range formats {
		if format.Format == config.PreferredFormat && format.ColorSpace == config.PreferredColorSpace {
			sc.surfaceFormat = format
			return nil
		}
	}

	// Use first available format as fallback
	sc.surfaceFormat = formats[0]
	return nil
}

// choosePresentMode chooses the best present mode
func (sc *SurfaceContextVK) choosePresentMode(config *SurfaceConfig) error {
	// Query present modes
	var presentModeCount uint32
	vulkan.GetPhysicalDeviceSurfacePresentModes(sc.context.physicalDevice, sc.surface, &presentModeCount, nil)
	
	if presentModeCount == 0 {
		return fmt.Errorf("no present modes available")
	}

	presentModes := make([]vulkan.PresentMode, presentModeCount)
	vulkan.GetPhysicalDeviceSurfacePresentModes(sc.context.physicalDevice, sc.surface, &presentModeCount, presentModes)

	// Choose preferred present mode
	preferredMode := config.PreferredPresentMode
	if config.VSync {
		preferredMode = vulkan.PresentModeFifo
	}

	for _, mode := range presentModes {
		if mode == preferredMode {
			sc.presentMode = mode
			return nil
		}
	}

	// FIFO is guaranteed to be available
	sc.presentMode = vulkan.PresentModeFifo
	return nil
}

// createSwapchain creates the Vulkan swapchain
func (sc *SurfaceContextVK) createSwapchain(config *SurfaceConfig) error {
	// Query surface capabilities
	var capabilities vulkan.SurfaceCapabilities
	if result := vulkan.GetPhysicalDeviceSurfaceCapabilities(sc.context.physicalDevice, sc.surface, &capabilities); result != vulkan.Success {
		return fmt.Errorf("failed to get surface capabilities: %s", result)
	}
	capabilities.Deref()

	// Choose extent
	if capabilities.CurrentExtent.Width != ^uint32(0) {
		sc.extent = capabilities.CurrentExtent
	} else {
		sc.extent = vulkan.Extent2D{
			Width:  clampUint32(config.Width, capabilities.MinImageExtent.Width, capabilities.MaxImageExtent.Width),
			Height: clampUint32(config.Height, capabilities.MinImageExtent.Height, capabilities.MaxImageExtent.Height),
		}
	}

	// Choose image count
	sc.imageCount = capabilities.MinImageCount + 1
	if capabilities.MaxImageCount > 0 && sc.imageCount > capabilities.MaxImageCount {
		sc.imageCount = capabilities.MaxImageCount
	}

	// Create swapchain
	swapchainInfo := vulkan.SwapchainCreateInfo{
		SType:            vulkan.StructureTypeSwapchainCreateInfo,
		Surface:          sc.surface,
		MinImageCount:    sc.imageCount,
		ImageFormat:      sc.surfaceFormat.Format,
		ImageColorSpace:  sc.surfaceFormat.ColorSpace,
		ImageExtent:      sc.extent,
		ImageArrayLayers: 1,
		ImageUsage:       vulkan.ImageUsageColorAttachmentBit,
		ImageSharingMode: vulkan.SharingModeExclusive,
		PreTransform:     capabilities.CurrentTransform,
		CompositeAlpha:   vulkan.CompositeAlphaOpaqueBit,
		PresentMode:      sc.presentMode,
		Clipped:          vulkan.True,
		OldSwapchain:     vulkan.NullSwapchain,
	}

	var swapchain vulkan.Swapchain
	if result := vulkan.CreateSwapchain(sc.context.device, &swapchainInfo, nil, &swapchain); result != vulkan.Success {
		return fmt.Errorf("failed to create swapchain: %s", result)
	}

	sc.swapchain = swapchain

	// Get swapchain images
	vulkan.GetSwapchainImages(sc.context.device, sc.swapchain, &sc.imageCount, nil)
	sc.swapchainImages = make([]vulkan.Image, sc.imageCount)
	vulkan.GetSwapchainImages(sc.context.device, sc.swapchain, &sc.imageCount, sc.swapchainImages)

	return nil
}

// createImageViews creates image views for swapchain images
func (sc *SurfaceContextVK) createImageViews() error {
	sc.swapchainImageViews = make([]vulkan.ImageView, len(sc.swapchainImages))

	for i, image := range sc.swapchainImages {
		imageViewInfo := vulkan.ImageViewCreateInfo{
			SType:    vulkan.StructureTypeImageViewCreateInfo,
			Image:    image,
			ViewType: vulkan.ImageViewType2d,
			Format:   sc.surfaceFormat.Format,
			Components: vulkan.ComponentMapping{
				R: vulkan.ComponentSwizzleIdentity,
				G: vulkan.ComponentSwizzleIdentity,
				B: vulkan.ComponentSwizzleIdentity,
				A: vulkan.ComponentSwizzleIdentity,
			},
			SubresourceRange: vulkan.ImageSubresourceRange{
				AspectMask:     vulkan.ImageAspectColorBit,
				BaseMipLevel:   0,
				LevelCount:     1,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
		}

		var imageView vulkan.ImageView
		if result := vulkan.CreateImageView(sc.context.device, &imageViewInfo, nil, &imageView); result != vulkan.Success {
			// Cleanup already created image views
			for j := 0; j < i; j++ {
				vulkan.DestroyImageView(sc.context.device, sc.swapchainImageViews[j], nil)
			}
			return fmt.Errorf("failed to create image view %d: %s", i, result)
		}

		sc.swapchainImageViews[i] = imageView
	}

	return nil
}

// AcquireNextImage acquires the next swapchain image
func (sc *SurfaceContextVK) AcquireNextImage(semaphore vulkan.Semaphore, fence vulkan.Fence) (uint32, error) {
	var imageIndex uint32
	result := vulkan.AcquireNextImage(sc.context.device, sc.swapchain, ^uint64(0), semaphore, fence, &imageIndex)
	
	if result == vulkan.ErrorOutOfDate || result == vulkan.Suboptimal {
		return imageIndex, fmt.Errorf("swapchain out of date")
	} else if result != vulkan.Success {
		return imageIndex, fmt.Errorf("failed to acquire next image: %s", result)
	}

	sc.currentImageIndex = imageIndex
	return imageIndex, nil
}

// PresentImage presents the current image
func (sc *SurfaceContextVK) PresentImage(waitSemaphores []vulkan.Semaphore) error {
	presentInfo := vulkan.PresentInfo{
		SType:              vulkan.StructureTypePresentInfo,
		WaitSemaphoreCount: uint32(len(waitSemaphores)),
		PWaitSemaphores:    waitSemaphores,
		SwapchainCount:     1,
		PSwapchains:        []vulkan.Swapchain{sc.swapchain},
		PImageIndices:      []uint32{sc.currentImageIndex},
	}

	result := vulkan.QueuePresent(sc.context.presentQueue, &presentInfo)
	if result == vulkan.ErrorOutOfDate || result == vulkan.Suboptimal {
		return fmt.Errorf("swapchain out of date")
	} else if result != vulkan.Success {
		return fmt.Errorf("failed to present image: %s", result)
	}

	return nil
}

// GetCurrentImage returns the current swapchain image
func (sc *SurfaceContextVK) GetCurrentImage() vulkan.Image {
	if sc.currentImageIndex < uint32(len(sc.swapchainImages)) {
		return sc.swapchainImages[sc.currentImageIndex]
	}
	return vulkan.NullImage
}

// GetCurrentImageView returns the current swapchain image view
func (sc *SurfaceContextVK) GetCurrentImageView() vulkan.ImageView {
	if sc.currentImageIndex < uint32(len(sc.swapchainImageViews)) {
		return sc.swapchainImageViews[sc.currentImageIndex]
	}
	return vulkan.NullImageView
}

// GetExtent returns the swapchain extent
func (sc *SurfaceContextVK) GetExtent() vulkan.Extent2D {
	return sc.extent
}

// GetFormat returns the swapchain format
func (sc *SurfaceContextVK) GetFormat() vulkan.Format {
	return sc.surfaceFormat.Format
}

// GetImageCount returns the number of swapchain images
func (sc *SurfaceContextVK) GetImageCount() uint32 {
	return sc.imageCount
}

// Destroy destroys the surface context and frees resources
func (sc *SurfaceContextVK) Destroy() {
	if sc.context == nil {
		return
	}

	// Destroy framebuffers
	for _, framebuffer := range sc.framebuffers {
		vulkan.DestroyFramebuffer(sc.context.device, framebuffer, nil)
	}

	// Destroy render pass
	if sc.renderPass != nil {
		sc.renderPass.Destroy()
	}

	// Destroy image views
	for _, imageView := range sc.swapchainImageViews {
		vulkan.DestroyImageView(sc.context.device, imageView, nil)
	}

	// Destroy swapchain
	if sc.swapchain != vulkan.NullSwapchain {
		vulkan.DestroySwapchain(sc.context.device, sc.swapchain, nil)
	}

	// Note: Surface should be destroyed by the caller
}

// clampUint32 clamps a value between min and max
func clampUint32(value, min, max uint32) uint32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
