package vulkan

import (
	"fmt"
	"unsafe"

	"github.com/vulkan-go/vulkan"
)

// TextureVK represents a Vulkan texture implementation
type TextureVK struct {
	context         *ContextVK
	image           vulkan.Image
	imageView       vulkan.ImageView
	deviceMemory    vulkan.DeviceMemory
	format          vulkan.Format
	extent          vulkan.Extent3D
	mipLevels       uint32
	arrayLayers     uint32
	sampleCount     vulkan.SampleCountFlagBits
	usage           vulkan.ImageUsageFlags
	layout          vulkan.ImageLayout
	aspectMask      vulkan.ImageAspectFlags
	ownsImage       bool
	ownsImageView   bool
	isSwapchainImage bool
}

// NewTextureVK creates a new Vulkan texture
func NewTextureVK(context *ContextVK, descriptor *TextureDescriptor) (*TextureVK, error) {
	if context == nil || descriptor == nil {
		return nil, fmt.Errorf("invalid parameters")
	}

	texture := &TextureVK{
		context:         context,
		format:          descriptor.Format,
		extent:          descriptor.Extent,
		mipLevels:       descriptor.MipLevels,
		arrayLayers:     descriptor.ArrayLayers,
		sampleCount:     descriptor.SampleCount,
		usage:           descriptor.Usage,
		layout:          vulkan.ImageLayoutUndefined,
		aspectMask:      getImageAspectFlags(descriptor.Format),
		ownsImage:       true,
		ownsImageView:   true,
		isSwapchainImage: false,
	}

	// Create the image
	if err := texture.createImage(); err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}

	// Allocate and bind memory
	if err := texture.allocateMemory(); err != nil {
		texture.Destroy()
		return nil, fmt.Errorf("failed to allocate memory: %w", err)
	}

	// Create image view
	if err := texture.createImageView(); err != nil {
		texture.Destroy()
		return nil, fmt.Errorf("failed to create image view: %w", err)
	}

	return texture, nil
}

// NewTextureVKFromImage creates a texture from an existing Vulkan image
func NewTextureVKFromImage(context *ContextVK, image vulkan.Image, format vulkan.Format, extent vulkan.Extent3D) (*TextureVK, error) {
	texture := &TextureVK{
		context:         context,
		image:           image,
		format:          format,
		extent:          extent,
		mipLevels:       1,
		arrayLayers:     1,
		sampleCount:     vulkan.SampleCount1Bit,
		usage:           vulkan.ImageUsageColorAttachmentBit,
		layout:          vulkan.ImageLayoutUndefined,
		aspectMask:      getImageAspectFlags(format),
		ownsImage:       false,
		ownsImageView:   true,
		isSwapchainImage: true,
	}

	// Create image view for existing image
	if err := texture.createImageView(); err != nil {
		return nil, fmt.Errorf("failed to create image view: %w", err)
	}

	return texture, nil
}

// createImage creates the Vulkan image
func (t *TextureVK) createImage() error {
	imageInfo := vulkan.ImageCreateInfo{
		SType:         vulkan.StructureTypeImageCreateInfo,
		ImageType:     vulkan.ImageType2d,
		Format:        t.format,
		Extent:        t.extent,
		MipLevels:     t.mipLevels,
		ArrayLayers:   t.arrayLayers,
		Samples:       t.sampleCount,
		Tiling:        vulkan.ImageTilingOptimal,
		Usage:         t.usage,
		SharingMode:   vulkan.SharingModeExclusive,
		InitialLayout: vulkan.ImageLayoutUndefined,
	}

	var image vulkan.Image
	if result := vulkan.CreateImage(t.context.device, &imageInfo, nil, &image); result != vulkan.Success {
		return fmt.Errorf("failed to create image: %s", result)
	}

	t.image = image
	return nil
}

// allocateMemory allocates and binds memory for the image
func (t *TextureVK) allocateMemory() error {
	// Get memory requirements
	var memRequirements vulkan.MemoryRequirements
	vulkan.GetImageMemoryRequirements(t.context.device, t.image, &memRequirements)
	memRequirements.Deref()

	// Find suitable memory type
	memTypeIndex, err := t.context.findMemoryType(memRequirements.MemoryTypeBits,
		vulkan.MemoryPropertyDeviceLocalBit)
	if err != nil {
		return fmt.Errorf("failed to find suitable memory type: %w", err)
	}

	// Allocate memory
	allocInfo := vulkan.MemoryAllocateInfo{
		SType:           vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  memRequirements.Size,
		MemoryTypeIndex: memTypeIndex,
	}

	var deviceMemory vulkan.DeviceMemory
	if result := vulkan.AllocateMemory(t.context.device, &allocInfo, nil, &deviceMemory); result != vulkan.Success {
		return fmt.Errorf("failed to allocate memory: %s", result)
	}

	// Bind memory to image
	if result := vulkan.BindImageMemory(t.context.device, t.image, deviceMemory, 0); result != vulkan.Success {
		vulkan.FreeMemory(t.context.device, deviceMemory, nil)
		return fmt.Errorf("failed to bind memory: %s", result)
	}

	t.deviceMemory = deviceMemory
	return nil
}

// createImageView creates the image view
func (t *TextureVK) createImageView() error {
	imageViewInfo := vulkan.ImageViewCreateInfo{
		SType:    vulkan.StructureTypeImageViewCreateInfo,
		Image:    t.image,
		ViewType: vulkan.ImageViewType2d,
		Format:   t.format,
		Components: vulkan.ComponentMapping{
			R: vulkan.ComponentSwizzleIdentity,
			G: vulkan.ComponentSwizzleIdentity,
			B: vulkan.ComponentSwizzleIdentity,
			A: vulkan.ComponentSwizzleIdentity,
		},
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     t.aspectMask,
			BaseMipLevel:   0,
			LevelCount:     t.mipLevels,
			BaseArrayLayer: 0,
			LayerCount:     t.arrayLayers,
		},
	}

	var imageView vulkan.ImageView
	if result := vulkan.CreateImageView(t.context.device, &imageViewInfo, nil, &imageView); result != vulkan.Success {
		return fmt.Errorf("failed to create image view: %s", result)
	}

	t.imageView = imageView
	return nil
}

// TransitionLayout transitions the image layout
func (t *TextureVK) TransitionLayout(oldLayout, newLayout vulkan.ImageLayout, commandBuffer vulkan.CommandBuffer) error {
	barrier := vulkan.ImageMemoryBarrier{
		SType:               vulkan.StructureTypeImageMemoryBarrier,
		OldLayout:           oldLayout,
		NewLayout:           newLayout,
		SrcQueueFamilyIndex: vulkan.QueueFamilyIgnored,
		DstQueueFamilyIndex: vulkan.QueueFamilyIgnored,
		Image:               t.image,
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     t.aspectMask,
			BaseMipLevel:   0,
			LevelCount:     t.mipLevels,
			BaseArrayLayer: 0,
			LayerCount:     t.arrayLayers,
		},
	}

	// Configure pipeline stages and access masks based on layout transition
	var sourceStage, destinationStage vulkan.PipelineStageFlags

	if oldLayout == vulkan.ImageLayoutUndefined && newLayout == vulkan.ImageLayoutTransferDstOptimal {
		barrier.SrcAccessMask = 0
		barrier.DstAccessMask = vulkan.AccessTransferWriteBit
		sourceStage = vulkan.PipelineStageTopOfPipeBit
		destinationStage = vulkan.PipelineStageTransferBit
	} else if oldLayout == vulkan.ImageLayoutTransferDstOptimal && newLayout == vulkan.ImageLayoutShaderReadOnlyOptimal {
		barrier.SrcAccessMask = vulkan.AccessTransferWriteBit
		barrier.DstAccessMask = vulkan.AccessShaderReadBit
		sourceStage = vulkan.PipelineStageTransferBit
		destinationStage = vulkan.PipelineStageFragmentShaderBit
	} else {
		return fmt.Errorf("unsupported layout transition")
	}

	vulkan.CmdPipelineBarrier(commandBuffer, sourceStage, destinationStage, 0, 0, nil, 0, nil, 1, []vulkan.ImageMemoryBarrier{barrier})

	t.layout = newLayout
	return nil
}

// GetImage returns the Vulkan image handle
func (t *TextureVK) GetImage() vulkan.Image {
	return t.image
}

// GetImageView returns the Vulkan image view handle
func (t *TextureVK) GetImageView() vulkan.ImageView {
	return t.imageView
}

// GetFormat returns the texture format
func (t *TextureVK) GetFormat() vulkan.Format {
	return t.format
}

// GetExtent returns the texture extent
func (t *TextureVK) GetExtent() vulkan.Extent3D {
	return t.extent
}

// GetMipLevels returns the number of mip levels
func (t *TextureVK) GetMipLevels() uint32 {
	return t.mipLevels
}

// GetLayout returns the current image layout
func (t *TextureVK) GetLayout() vulkan.ImageLayout {
	return t.layout
}

// IsSwapchainImage returns true if this is a swapchain image
func (t *TextureVK) IsSwapchainImage() bool {
	return t.isSwapchainImage
}

// Destroy destroys the texture and frees resources
func (t *TextureVK) Destroy() {
	if t.context == nil {
		return
	}

	if t.ownsImageView && t.imageView != vulkan.NullImageView {
		vulkan.DestroyImageView(t.context.device, t.imageView, nil)
		t.imageView = vulkan.NullImageView
	}

	if t.ownsImage && t.image != vulkan.NullImage {
		vulkan.DestroyImage(t.context.device, t.image, nil)
		t.image = vulkan.NullImage
	}

	if t.deviceMemory != vulkan.NullDeviceMemory {
		vulkan.FreeMemory(t.context.device, t.deviceMemory, nil)
		t.deviceMemory = vulkan.NullDeviceMemory
	}
}

// TextureDescriptor describes texture creation parameters
type TextureDescriptor struct {
	Format      vulkan.Format
	Extent      vulkan.Extent3D
	MipLevels   uint32
	ArrayLayers uint32
	SampleCount vulkan.SampleCountFlagBits
	Usage       vulkan.ImageUsageFlags
}

// getImageAspectFlags returns appropriate aspect flags for the format
func getImageAspectFlags(format vulkan.Format) vulkan.ImageAspectFlags {
	switch format {
	case vulkan.FormatD16Unorm, vulkan.FormatD32Sfloat:
		return vulkan.ImageAspectDepthBit
	case vulkan.FormatD16UnormS8Uint, vulkan.FormatD24UnormS8Uint, vulkan.FormatD32SfloatS8Uint:
		return vulkan.ImageAspectDepthBit | vulkan.ImageAspectStencilBit
	case vulkan.FormatS8Uint:
		return vulkan.ImageAspectStencilBit
	default:
		return vulkan.ImageAspectColorBit
	}
}
