package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// TextureSourceVK represents a Vulkan texture source implementation
type TextureSourceVK struct {
	context      *ContextVK
	image        vulkan.Image
	imageView    vulkan.ImageView
	deviceMemory vulkan.DeviceMemory
	format       vulkan.Format
	extent       vulkan.Extent2D
	mipLevels    uint32
	ownsImage    bool
}

// NewTextureSourceVK creates a new Vulkan texture source
func NewTextureSourceVK(context *ContextVK, descriptor *TextureSourceDescriptor) (*TextureSourceVK, error) {
	if context == nil || descriptor == nil {
		return nil, fmt.Errorf("invalid parameters")
	}

	source := &TextureSourceVK{
		context:   context,
		format:    descriptor.Format,
		extent:    descriptor.Extent,
		mipLevels: descriptor.MipLevels,
		ownsImage: true,
	}

	// Create the texture source
	if err := source.createImageAndView(); err != nil {
		return nil, fmt.Errorf("failed to create texture source: %w", err)
	}

	return source, nil
}

// NewTextureSourceVKFromExisting creates a texture source from existing Vulkan objects
func NewTextureSourceVKFromExisting(context *ContextVK, image vulkan.Image, imageView vulkan.ImageView, format vulkan.Format, extent vulkan.Extent2D) *TextureSourceVK {
	return &TextureSourceVK{
		context:   context,
		image:     image,
		imageView: imageView,
		format:    format,
		extent:    extent,
		mipLevels: 1,
		ownsImage: false,
	}
}

// createImageAndView creates the image and image view
func (ts *TextureSourceVK) createImageAndView() error {
	// Create image
	imageInfo := vulkan.ImageCreateInfo{
		SType:     vulkan.StructureTypeImageCreateInfo,
		ImageType: vulkan.ImageType2d,
		Format:    ts.format,
		Extent: vulkan.Extent3D{
			Width:  ts.extent.Width,
			Height: ts.extent.Height,
			Depth:  1,
		},
		MipLevels:     ts.mipLevels,
		ArrayLayers:   1,
		Samples:       vulkan.SampleCount1Bit,
		Tiling:        vulkan.ImageTilingOptimal,
		Usage:         vulkan.ImageUsageSampledBit | vulkan.ImageUsageTransferDstBit,
		SharingMode:   vulkan.SharingModeExclusive,
		InitialLayout: vulkan.ImageLayoutUndefined,
	}

	var image vulkan.Image
	if result := vulkan.CreateImage(ts.context.device, &imageInfo, nil, &image); result != vulkan.Success {
		return fmt.Errorf("failed to create image: %s", result)
	}
	ts.image = image

	// Allocate memory for the image
	if err := ts.allocateMemory(); err != nil {
		return fmt.Errorf("failed to allocate memory: %w", err)
	}

	// Create image view
	imageViewInfo := vulkan.ImageViewCreateInfo{
		SType:    vulkan.StructureTypeImageViewCreateInfo,
		Image:    ts.image,
		ViewType: vulkan.ImageViewType2d,
		Format:   ts.format,
		Components: vulkan.ComponentMapping{
			R: vulkan.ComponentSwizzleIdentity,
			G: vulkan.ComponentSwizzleIdentity,
			B: vulkan.ComponentSwizzleIdentity,
			A: vulkan.ComponentSwizzleIdentity,
		},
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     vulkan.ImageAspectColorBit,
			BaseMipLevel:   0,
			LevelCount:     ts.mipLevels,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}

	var imageView vulkan.ImageView
	if result := vulkan.CreateImageView(ts.context.device, &imageViewInfo, nil, &imageView); result != vulkan.Success {
		vulkan.DestroyImage(ts.context.device, ts.image, nil)
		return fmt.Errorf("failed to create image view: %s", result)
	}
	ts.imageView = imageView

	return nil
}

// allocateMemory allocates and binds memory for the image
func (ts *TextureSourceVK) allocateMemory() error {
	// Get memory requirements
	var memRequirements vulkan.MemoryRequirements
	vulkan.GetImageMemoryRequirements(ts.context.device, ts.image, &memRequirements)
	memRequirements.Deref()

	// Find suitable memory type (TODO: implement proper memory type selection)
	memTypeIndex := uint32(0) // Placeholder

	// Allocate memory
	allocInfo := vulkan.MemoryAllocateInfo{
		SType:           vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  memRequirements.Size,
		MemoryTypeIndex: memTypeIndex,
	}

	var deviceMemory vulkan.DeviceMemory
	if result := vulkan.AllocateMemory(ts.context.device, &allocInfo, nil, &deviceMemory); result != vulkan.Success {
		return fmt.Errorf("failed to allocate memory: %s", result)
	}

	// Bind memory to image
	if result := vulkan.BindImageMemory(ts.context.device, ts.image, deviceMemory, 0); result != vulkan.Success {
		vulkan.FreeMemory(ts.context.device, deviceMemory, nil)
		return fmt.Errorf("failed to bind memory: %s", result)
	}

	ts.deviceMemory = deviceMemory
	return nil
}

// UploadData uploads data to the texture source
func (ts *TextureSourceVK) UploadData(data []byte, width, height uint32) error {
	// TODO: Implement data upload via staging buffer
	return fmt.Errorf("texture data upload not implemented")
}

// GetImage returns the Vulkan image handle
func (ts *TextureSourceVK) GetImage() vulkan.Image {
	return ts.image
}

// GetImageView returns the Vulkan image view handle
func (ts *TextureSourceVK) GetImageView() vulkan.ImageView {
	return ts.imageView
}

// GetFormat returns the texture format
func (ts *TextureSourceVK) GetFormat() vulkan.Format {
	return ts.format
}

// GetExtent returns the texture extent
func (ts *TextureSourceVK) GetExtent() vulkan.Extent2D {
	return ts.extent
}

// GetMipLevels returns the number of mip levels
func (ts *TextureSourceVK) GetMipLevels() uint32 {
	return ts.mipLevels
}

// Destroy destroys the texture source and frees resources
func (ts *TextureSourceVK) Destroy() {
	if ts.context == nil {
		return
	}

	if ts.imageView != vulkan.NullImageView {
		vulkan.DestroyImageView(ts.context.device, ts.imageView, nil)
		ts.imageView = vulkan.NullImageView
	}

	if ts.ownsImage && ts.image != vulkan.NullImage {
		vulkan.DestroyImage(ts.context.device, ts.image, nil)
		ts.image = vulkan.NullImage
	}

	if ts.deviceMemory != vulkan.NullDeviceMemory {
		vulkan.FreeMemory(ts.context.device, ts.deviceMemory, nil)
		ts.deviceMemory = vulkan.NullDeviceMemory
	}
}

// TextureSourceDescriptor describes texture source creation parameters
type TextureSourceDescriptor struct {
	Format    vulkan.Format
	Extent    vulkan.Extent2D
	MipLevels uint32
	Usage     vulkan.ImageUsageFlags
}
