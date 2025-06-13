package vulkan

import (
	"fmt"
	"sync"
	"unsafe"

	. "github.com/opensraph/sraph/gpu/wgpu"
)

// VulkanSurface represents a Vulkan surface
type VulkanSurface struct {
	mu         sync.RWMutex
	handle     uintptr // VkSurfaceKHR handle
	instance   *VulkanInstance
	swapchain  *VulkanSwapchain
	configured bool
	config     SurfaceConfiguration
	destroyed  bool
}

// VulkanSwapchain represents a Vulkan swapchain
type VulkanSwapchain struct {
	mu           sync.RWMutex
	handle       uintptr // VkSwapchainKHR handle
	device       *VulkanDevice
	surface      *VulkanSurface
	images       []uintptr // VkImage handles
	imageViews   []uintptr // VkImageView handles
	currentImage uint32
	config       SurfaceConfiguration
	destroyed    bool
}

// NewVulkanSurface creates a new Vulkan surface
func NewVulkanSurface(instance *VulkanInstance, platformData unsafe.Pointer) (*VulkanSurface, error) {
	if instance == nil {
		return nil, fmt.Errorf("instance cannot be nil")
	}

	surface := &VulkanSurface{
		instance: instance,
	}

	handle, err := instance.CreateSurface(platformData)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vulkan surface: %v", err)
	}

	surface.handle = handle
	return surface, nil
}

// Configure configures the Vulkan surface with a swapchain
func (vs *VulkanSurface) Configure(config SurfaceConfiguration, vulkanDevice *VulkanDevice) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if vs.destroyed {
		return fmt.Errorf("surface has been destroyed")
	}

	if vulkanDevice == nil {
		return fmt.Errorf("device cannot be nil")
	}

	// Validate configuration
	if config.Width == 0 || config.Height == 0 {
		return fmt.Errorf("width and height must be greater than 0")
	}

	// Create or recreate swapchain
	swapchain, err := NewVulkanSwapchain(vulkanDevice, vs, config)
	if err != nil {
		return fmt.Errorf("failed to create swapchain: %v", err)
	}

	// Destroy old swapchain if exists
	if vs.swapchain != nil {
		vs.swapchain.Destroy()
	}

	vs.swapchain = swapchain
	vs.config = config
	vs.configured = true

	return nil
}

// GetCurrentTexture gets the current swapchain texture
func (vs *VulkanSurface) GetCurrentTexture() (*VulkanTexture, error) {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if vs.destroyed {
		return nil, fmt.Errorf("surface has been destroyed")
	}

	if !vs.configured || vs.swapchain == nil {
		return nil, fmt.Errorf("surface is not configured")
	}

	return vs.swapchain.GetCurrentTexture()
}

// Present presents the current frame
func (vs *VulkanSurface) Present() error {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if vs.destroyed {
		return fmt.Errorf("surface has been destroyed")
	}

	if !vs.configured || vs.swapchain == nil {
		return fmt.Errorf("surface is not configured")
	}

	return vs.swapchain.Present()
}

// GetHandle returns the surface handle
func (vs *VulkanSurface) GetHandle() uintptr {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	return vs.handle
}

// IsConfigured checks if the surface is configured
func (vs *VulkanSurface) IsConfigured() bool {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	return vs.configured && vs.swapchain != nil
}

// GetConfig returns the current surface configuration
func (vs *VulkanSurface) GetConfig() SurfaceConfiguration {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	return vs.config
}

// Destroy destroys the Vulkan surface
func (vs *VulkanSurface) Destroy() error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if vs.destroyed {
		return nil
	}

	// Destroy swapchain first
	if vs.swapchain != nil {
		vs.swapchain.Destroy()
		vs.swapchain = nil
	}

	// In a real implementation, this would call vkDestroySurfaceKHR
	vs.handle = 0
	vs.destroyed = true

	return nil
}

// NewVulkanSwapchain creates a new Vulkan swapchain
func NewVulkanSwapchain(device *VulkanDevice, surface *VulkanSurface, config SurfaceConfiguration) (*VulkanSwapchain, error) {
	if device == nil {
		return nil, fmt.Errorf("device cannot be nil")
	}
	if surface == nil {
		return nil, fmt.Errorf("surface cannot be nil")
	}

	swapchain := &VulkanSwapchain{
		device:  device,
		surface: surface,
		config:  config,
	}

	if err := swapchain.createSwapchain(); err != nil {
		return nil, fmt.Errorf("failed to create swapchain: %v", err)
	}

	return swapchain, nil
}

// createSwapchain creates the actual Vulkan swapchain
func (vsc *VulkanSwapchain) createSwapchain() error {
	// In a real implementation, this would call vkCreateSwapchainKHR
	vsc.handle = uintptr(0xAABBCCDD + int(vsc.config.Width)*int(vsc.config.Height))

	// Create mock swapchain images
	imageCount := uint32(3) // Triple buffering
	vsc.images = make([]uintptr, imageCount)
	vsc.imageViews = make([]uintptr, imageCount)

	for i := uint32(0); i < imageCount; i++ {
		vsc.images[i] = uintptr(0xDDCCBBAA + int(i)*0x1000)
		vsc.imageViews[i] = uintptr(0xEEDDCCBB + int(i)*0x1000)
	}

	return nil
}

// GetCurrentTexture gets the current swapchain texture
func (vsc *VulkanSwapchain) GetCurrentTexture() (*VulkanTexture, error) {
	vsc.mu.Lock()
	defer vsc.mu.Unlock()

	if vsc.destroyed {
		return nil, fmt.Errorf("swapchain has been destroyed")
	}

	if len(vsc.images) == 0 {
		return nil, fmt.Errorf("no swapchain images available")
	}

	// In a real implementation, this would call vkAcquireNextImageKHR
	imageIndex := vsc.currentImage % uint32(len(vsc.images))

	// Create a texture wrapper for the swapchain image
	descriptor := TextureDescriptor{
		Label:     "Swapchain Texture",
		Usage:     TextureUsageRenderAttachment | TextureUsageCopySrc,
		Dimension: TextureDimension2D,
		Size: Extent3D{
			Width:              vsc.config.Width,
			Height:             vsc.config.Height,
			DepthOrArrayLayers: 1,
		},
		Format:        vsc.config.Format,
		MipLevelCount: 1,
		SampleCount:   1,
	}

	texture := &VulkanTexture{
		handle:     vsc.images[imageIndex],
		device:     vsc.device,
		descriptor: descriptor,
		allocation: nil, // Swapchain images don't have separate allocations
	}

	return texture, nil
}

// Present presents the current frame
func (vsc *VulkanSwapchain) Present() error {
	vsc.mu.Lock()
	defer vsc.mu.Unlock()

	if vsc.destroyed {
		return fmt.Errorf("swapchain has been destroyed")
	}

	if len(vsc.images) == 0 {
		return fmt.Errorf("no swapchain images available")
	}

	// In a real implementation, this would call vkQueuePresentKHR
	vsc.currentImage = (vsc.currentImage + 1) % uint32(len(vsc.images))

	return nil
}

// GetImageCount returns the number of swapchain images
func (vsc *VulkanSwapchain) GetImageCount() uint32 {
	vsc.mu.RLock()
	defer vsc.mu.RUnlock()
	return uint32(len(vsc.images))
}

// GetCurrentImageIndex returns the current image index
func (vsc *VulkanSwapchain) GetCurrentImageIndex() uint32 {
	vsc.mu.RLock()
	defer vsc.mu.RUnlock()
	return vsc.currentImage
}

// GetHandle returns the swapchain handle
func (vsc *VulkanSwapchain) GetHandle() uintptr {
	vsc.mu.RLock()
	defer vsc.mu.RUnlock()
	return vsc.handle
}

// Destroy destroys the Vulkan swapchain
func (vsc *VulkanSwapchain) Destroy() error {
	vsc.mu.Lock()
	defer vsc.mu.Unlock()

	if vsc.destroyed {
		return nil
	}

	// In a real implementation, this would destroy image views and call vkDestroySwapchainKHR
	vsc.handle = 0
	vsc.images = nil
	vsc.imageViews = nil
	vsc.destroyed = true

	return nil
}
