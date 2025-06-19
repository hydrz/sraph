package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// SwapchainVK wraps Vulkan swapchain
// It manages the presentation images and surface
type SwapchainVK struct {
	device            vulkan.Device
	swapchain         vulkan.Swapchain
	surface           vulkan.Surface
	format            vulkan.SurfaceFormat
	presentMode       vulkan.PresentMode
	extent            vulkan.Extent2D
	images            []vulkan.Image
	imageViews        []vulkan.ImageView
	currentImageIndex uint32
	imageCount        uint32
}

// SwapchainSupportDetails contains swapchain support information
type SwapchainSupportDetails struct {
	Capabilities vulkan.SurfaceCapabilities
	Formats      []vulkan.SurfaceFormat
	PresentModes []vulkan.PresentMode
}

// NewSwapchainVK creates a new Vulkan swapchain
func NewSwapchainVK(device vulkan.Device, physicalDevice vulkan.PhysicalDevice, surface vulkan.Surface, width, height uint32, oldSwapchain vulkan.Swapchain) (*SwapchainVK, error) {
	// Query swapchain support
	support, err := querySwapchainSupport(physicalDevice, surface)
	if err != nil {
		return nil, err
	}

	// Choose surface format
	surfaceFormat := chooseSurfaceFormat(support.Formats)

	// Choose present mode
	presentMode := choosePresentMode(support.PresentModes)

	// Choose extent
	extent := chooseExtent(support.Capabilities, width, height)

	// Choose image count
	imageCount := support.Capabilities.MinImageCount + 1
	if support.Capabilities.MaxImageCount > 0 && imageCount > support.Capabilities.MaxImageCount {
		imageCount = support.Capabilities.MaxImageCount
	}

	// Create swapchain
	createInfo := vulkan.SwapchainCreateInfo{
		SType:            vulkan.StructureTypeSwapchainCreateInfo,
		Surface:          surface,
		MinImageCount:    imageCount,
		ImageFormat:      surfaceFormat.Format,
		ImageColorSpace:  surfaceFormat.ColorSpace,
		ImageExtent:      extent,
		ImageArrayLayers: 1,
		ImageUsage:       vulkan.ImageUsageFlags(vulkan.ImageUsageColorAttachmentBit),
		PreTransform:     support.Capabilities.CurrentTransform,
		CompositeAlpha:   vulkan.CompositeAlphaFlags(vulkan.CompositeAlphaOpaqueBit),
		PresentMode:      presentMode,
		Clipped:          vulkan.True,
		OldSwapchain:     oldSwapchain,
	}

	// TODO: Handle queue family indices for concurrent/exclusive sharing mode

	var swapchain vulkan.Swapchain
	if result := vulkan.CreateSwapchain(device, &createInfo, nil, &swapchain); result != vulkan.Success {
		return nil, fmt.Errorf("failed to create swapchain: %s", result)
	}

	// Get swapchain images
	var actualImageCount uint32
	vulkan.GetSwapchainImages(device, swapchain, &actualImageCount, nil)
	images := make([]vulkan.Image, actualImageCount)
	vulkan.GetSwapchainImages(device, swapchain, &actualImageCount, images)

	// Create image views
	imageViews := make([]vulkan.ImageView, actualImageCount)
	for i := range images {
		imageViews[i], err = createImageView(device, images[i], surfaceFormat.Format, vulkan.ImageAspectFlags(vulkan.ImageAspectColorBit))
		if err != nil {
			// Clean up created image views
			for j := 0; j < i; j++ {
				vulkan.DestroyImageView(device, imageViews[j], nil)
			}
			vulkan.DestroySwapchain(device, swapchain, nil)
			return nil, err
		}
	}

	return &SwapchainVK{
		device:            device,
		swapchain:         swapchain,
		surface:           surface,
		format:            surfaceFormat,
		presentMode:       presentMode,
		extent:            extent,
		images:            images,
		imageViews:        imageViews,
		currentImageIndex: 0,
		imageCount:        actualImageCount,
	}, nil
}

// GetSwapchain returns the Vulkan swapchain handle
func (s *SwapchainVK) GetSwapchain() vulkan.Swapchain {
	return s.swapchain
}

// GetImageCount returns the number of swapchain images
func (s *SwapchainVK) GetImageCount() uint32 {
	return s.imageCount
}

// GetImages returns the swapchain images
func (s *SwapchainVK) GetImages() []vulkan.Image {
	return s.images
}

// GetImageViews returns the swapchain image views
func (s *SwapchainVK) GetImageViews() []vulkan.ImageView {
	return s.imageViews
}

// GetFormat returns the swapchain surface format
func (s *SwapchainVK) GetFormat() vulkan.SurfaceFormat {
	return s.format
}

// GetExtent returns the swapchain extent
func (s *SwapchainVK) GetExtent() vulkan.Extent2D {
	return s.extent
}

// GetPresentMode returns the present mode
func (s *SwapchainVK) GetPresentMode() vulkan.PresentMode {
	return s.presentMode
}

// GetCurrentImageIndex returns the current image index
func (s *SwapchainVK) GetCurrentImageIndex() uint32 {
	return s.currentImageIndex
}

// AcquireNextImage acquires the next image from the swapchain
func (s *SwapchainVK) AcquireNextImage(timeout uint64, semaphore vulkan.Semaphore, fence vulkan.Fence) (uint32, vulkan.Result) {
	var imageIndex uint32
	result := vulkan.AcquireNextImage(s.device, s.swapchain, timeout, semaphore, fence, &imageIndex)
	if result == vulkan.Success || result == vulkan.Suboptimal {
		s.currentImageIndex = imageIndex
	}
	return imageIndex, result
}

// Present presents the current image
func (s *SwapchainVK) Present(queue vulkan.Queue, waitSemaphores []vulkan.Semaphore, imageIndex uint32) vulkan.Result {
	presentInfo := vulkan.PresentInfo{
		SType:              vulkan.StructureTypePresentInfo,
		WaitSemaphoreCount: uint32(len(waitSemaphores)),
		SwapchainCount:     1,
		PSwapchains:        []vulkan.Swapchain{s.swapchain},
		PImageIndices:      []uint32{imageIndex},
	}

	if len(waitSemaphores) > 0 {
		presentInfo.PWaitSemaphores = waitSemaphores
	}

	return vulkan.QueuePresent(queue, &presentInfo)
}

// Recreate recreates the swapchain with new dimensions
func (s *SwapchainVK) Recreate(physicalDevice vulkan.PhysicalDevice, width, height uint32) error {
	oldSwapchain := s.swapchain

	newSwapchain, err := NewSwapchainVK(s.device, physicalDevice, s.surface, width, height, oldSwapchain)
	if err != nil {
		return err
	}

	// Destroy old resources
	s.Destroy()

	// Update with new swapchain
	*s = *newSwapchain

	// Destroy old swapchain
	if oldSwapchain != vulkan.NullSwapchain {
		vulkan.DestroySwapchain(s.device, oldSwapchain, nil)
	}

	return nil
}

// Destroy destroys the swapchain and its resources
func (s *SwapchainVK) Destroy() {
	// Destroy image views
	for _, imageView := range s.imageViews {
		if imageView != vulkan.NullImageView {
			vulkan.DestroyImageView(s.device, imageView, nil)
		}
	}
	s.imageViews = s.imageViews[:0]

	// Destroy swapchain
	if s.swapchain != vulkan.NullSwapchain {
		vulkan.DestroySwapchain(s.device, s.swapchain, nil)
		s.swapchain = vulkan.NullSwapchain
	}

	s.images = s.images[:0]
	s.imageCount = 0
}

// querySwapchainSupport queries swapchain support details
func querySwapchainSupport(physicalDevice vulkan.PhysicalDevice, surface vulkan.Surface) (*SwapchainSupportDetails, error) {
	var capabilities vulkan.SurfaceCapabilities
	if result := vulkan.GetPhysicalDeviceSurfaceCapabilities(physicalDevice, surface, &capabilities); result != vulkan.Success {
		return nil, fmt.Errorf("failed to get surface capabilities: %s", result)
	}

	// Get surface formats
	var formatCount uint32
	vulkan.GetPhysicalDeviceSurfaceFormats(physicalDevice, surface, &formatCount, nil)
	formats := make([]vulkan.SurfaceFormat, formatCount)
	if formatCount > 0 {
		vulkan.GetPhysicalDeviceSurfaceFormats(physicalDevice, surface, &formatCount, formats)
	}

	// Get present modes
	var presentModeCount uint32
	vulkan.GetPhysicalDeviceSurfacePresentModes(physicalDevice, surface, &presentModeCount, nil)
	presentModes := make([]vulkan.PresentMode, presentModeCount)
	if presentModeCount > 0 {
		vulkan.GetPhysicalDeviceSurfacePresentModes(physicalDevice, surface, &presentModeCount, presentModes)
	}

	return &SwapchainSupportDetails{
		Capabilities: capabilities,
		Formats:      formats,
		PresentModes: presentModes,
	}, nil
}

// chooseSurfaceFormat chooses the best surface format
func chooseSurfaceFormat(availableFormats []vulkan.SurfaceFormat) vulkan.SurfaceFormat {
	// Prefer sRGB format
	for _, format := range availableFormats {
		if format.Format == vulkan.FormatB8g8r8a8Srgb && format.ColorSpace == vulkan.ColorSpaceSrgbNonlinear {
			return format
		}
	}

	// Fall back to first available format
	if len(availableFormats) > 0 {
		return availableFormats[0]
	}

	// Default format
	return vulkan.SurfaceFormat{
		Format:     vulkan.FormatB8g8r8a8Srgb,
		ColorSpace: vulkan.ColorSpaceSrgbNonlinear,
	}
}

// choosePresentMode chooses the best present mode
func choosePresentMode(availablePresentModes []vulkan.PresentMode) vulkan.PresentMode {
	// Prefer mailbox mode for triple buffering
	for _, mode := range availablePresentModes {
		if mode == vulkan.PresentModeMailbox {
			return mode
		}
	}

	// Fall back to FIFO (always available)
	return vulkan.PresentModeFifo
}

// chooseExtent chooses the swap extent
func chooseExtent(capabilities vulkan.SurfaceCapabilities, width, height uint32) vulkan.Extent2D {
	if capabilities.CurrentExtent.Width != ^uint32(0) {
		return capabilities.CurrentExtent
	}

	extent := vulkan.Extent2D{
		Width:  width,
		Height: height,
	}

	// Clamp to supported range
	if extent.Width < capabilities.MinImageExtent.Width {
		extent.Width = capabilities.MinImageExtent.Width
	} else if extent.Width > capabilities.MaxImageExtent.Width {
		extent.Width = capabilities.MaxImageExtent.Width
	}

	if extent.Height < capabilities.MinImageExtent.Height {
		extent.Height = capabilities.MinImageExtent.Height
	} else if extent.Height > capabilities.MaxImageExtent.Height {
		extent.Height = capabilities.MaxImageExtent.Height
	}

	return extent
}

// createImageView creates an image view
func createImageView(device vulkan.Device, image vulkan.Image, format vulkan.Format, aspectFlags vulkan.ImageAspectFlags) (vulkan.ImageView, error) {
	createInfo := vulkan.ImageViewCreateInfo{
		SType:    vulkan.StructureTypeImageViewCreateInfo,
		Image:    image,
		ViewType: vulkan.ImageViewType2d,
		Format:   format,
		Components: vulkan.ComponentMapping{
			R: vulkan.ComponentSwizzleIdentity,
			G: vulkan.ComponentSwizzleIdentity,
			B: vulkan.ComponentSwizzleIdentity,
			A: vulkan.ComponentSwizzleIdentity,
		},
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     aspectFlags,
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}

	var imageView vulkan.ImageView
	if result := vulkan.CreateImageView(device, &createInfo, nil, &imageView); result != vulkan.Success {
		return vulkan.NullImageView, fmt.Errorf("failed to create image view: %s", result)
	}

	return imageView, nil
}
