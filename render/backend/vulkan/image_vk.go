package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// ImageVK wraps Vulkan image
// It represents a multidimensional array of data
type ImageVK struct {
	device      vulkan.Device
	image       vulkan.Image
	memory      vulkan.DeviceMemory
	imageView   vulkan.ImageView
	format      vulkan.Format
	width       uint32
	height      uint32
	depth       uint32
	mipLevels   uint32
	arrayLayers uint32
	samples     vulkan.SampleCountFlags
	usage       vulkan.ImageUsageFlags
	layout      vulkan.ImageLayout
	owned       bool // Whether this ImageVK owns the image (should destroy it)
}

// ImageCreateInfo contains image creation parameters
type ImageCreateInfo struct {
	Width       uint32
	Height      uint32
	Depth       uint32
	MipLevels   uint32
	ArrayLayers uint32
	Format      vulkan.Format
	Tiling      vulkan.ImageTiling
	Usage       vulkan.ImageUsageFlags
	Samples     vulkan.SampleCountFlags
	Properties  vulkan.MemoryPropertyFlags
}

// NewImageVK creates a new Vulkan image
func NewImageVK(device vulkan.Device, allocator *AllocatorVK, createInfo *ImageCreateInfo) (*ImageVK, error) {
	if createInfo == nil {
		return nil, fmt.Errorf("image create info cannot be nil")
	}

	// Create image
	vkCreateInfo := vulkan.ImageCreateInfo{
		SType:     vulkan.StructureTypeImageCreateInfo,
		ImageType: getImageType(createInfo.Width, createInfo.Height, createInfo.Depth),
		Format:    createInfo.Format,
		Extent: vulkan.Extent3D{
			Width:  createInfo.Width,
			Height: createInfo.Height,
			Depth:  createInfo.Depth,
		},
		MipLevels:     createInfo.MipLevels,
		ArrayLayers:   createInfo.ArrayLayers,
		Samples:       createInfo.Samples,
		Tiling:        createInfo.Tiling,
		Usage:         createInfo.Usage,
		SharingMode:   vulkan.SharingModeExclusive,
		InitialLayout: vulkan.ImageLayoutUndefined,
	}

	var image vulkan.Image
	if result := vulkan.CreateImage(device, &vkCreateInfo, nil, &image); result != vulkan.Success {
		return nil, fmt.Errorf("failed to create image: %s", result)
	}

	// Get memory requirements
	memReqs := allocator.GetImageMemoryRequirements(image)

	// Allocate memory
	allocation, err := allocator.AllocateMemory(memReqs.Size, memReqs.MemoryTypeBits, createInfo.Properties)
	if err != nil {
		vulkan.DestroyImage(device, image, nil)
		return nil, err
	}

	// Bind image memory
	if err := allocator.BindImageMemory(image, allocation); err != nil {
		allocator.FreeMemory(allocation)
		vulkan.DestroyImage(device, image, nil)
		return nil, err
	}

	return &ImageVK{
		device:      device,
		image:       image,
		memory:      allocation.GetMemory(),
		format:      createInfo.Format,
		width:       createInfo.Width,
		height:      createInfo.Height,
		depth:       createInfo.Depth,
		mipLevels:   createInfo.MipLevels,
		arrayLayers: createInfo.ArrayLayers,
		samples:     createInfo.Samples,
		usage:       createInfo.Usage,
		layout:      vulkan.ImageLayoutUndefined,
		owned:       true,
	}, nil
}

// NewImageVKFromHandle creates an ImageVK from existing Vulkan image handle
func NewImageVKFromHandle(device vulkan.Device, image vulkan.Image, format vulkan.Format, width, height uint32) *ImageVK {
	return &ImageVK{
		device:      device,
		image:       image,
		format:      format,
		width:       width,
		height:      height,
		depth:       1,
		mipLevels:   1,
		arrayLayers: 1,
		samples:     vulkan.SampleCountFlags(vulkan.SampleCount1Bit),
		layout:      vulkan.ImageLayoutUndefined,
		owned:       false, // Don't destroy externally owned image
	}
}

// GetImage returns the Vulkan image handle
func (img *ImageVK) GetImage() vulkan.Image {
	return img.image
}

// GetImageView returns the image view (creates if not exists)
func (img *ImageVK) GetImageView() vulkan.ImageView {
	if img.imageView == vulkan.NullImageView {
		// Create default image view
		aspectMask := getImageAspectMask(img.format)
		imageView, err := img.CreateImageView(vulkan.ImageViewType2d, aspectMask, 0, img.mipLevels, 0, img.arrayLayers)
		if err == nil {
			img.imageView = imageView
		}
	}
	return img.imageView
}

// GetFormat returns the image format
func (img *ImageVK) GetFormat() vulkan.Format {
	return img.format
}

// GetWidth returns the image width
func (img *ImageVK) GetWidth() uint32 {
	return img.width
}

// GetHeight returns the image height
func (img *ImageVK) GetHeight() uint32 {
	return img.height
}

// GetDepth returns the image depth
func (img *ImageVK) GetDepth() uint32 {
	return img.depth
}

// GetMipLevels returns the number of mip levels
func (img *ImageVK) GetMipLevels() uint32 {
	return img.mipLevels
}

// GetArrayLayers returns the number of array layers
func (img *ImageVK) GetArrayLayers() uint32 {
	return img.arrayLayers
}

// GetSamples returns the sample count
func (img *ImageVK) GetSamples() vulkan.SampleCountFlags {
	return img.samples
}

// GetUsage returns the image usage flags
func (img *ImageVK) GetUsage() vulkan.ImageUsageFlags {
	return img.usage
}

// GetLayout returns the current image layout
func (img *ImageVK) GetLayout() vulkan.ImageLayout {
	return img.layout
}

// SetLayout sets the current image layout
func (img *ImageVK) SetLayout(layout vulkan.ImageLayout) {
	img.layout = layout
}

// GetSize returns the image dimensions
func (img *ImageVK) GetSize() (uint32, uint32, uint32) {
	return img.width, img.height, img.depth
}

// CreateImageView creates an image view for this image
func (img *ImageVK) CreateImageView(viewType vulkan.ImageViewType, aspectMask vulkan.ImageAspectFlags, baseMipLevel, levelCount, baseArrayLayer, layerCount uint32) (vulkan.ImageView, error) {
	createInfo := vulkan.ImageViewCreateInfo{
		SType:    vulkan.StructureTypeImageViewCreateInfo,
		Image:    img.image,
		ViewType: viewType,
		Format:   img.format,
		Components: vulkan.ComponentMapping{
			R: vulkan.ComponentSwizzleIdentity,
			G: vulkan.ComponentSwizzleIdentity,
			B: vulkan.ComponentSwizzleIdentity,
			A: vulkan.ComponentSwizzleIdentity,
		},
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     aspectMask,
			BaseMipLevel:   baseMipLevel,
			LevelCount:     levelCount,
			BaseArrayLayer: baseArrayLayer,
			LayerCount:     layerCount,
		},
	}

	var imageView vulkan.ImageView
	if result := vulkan.CreateImageView(img.device, &createInfo, nil, &imageView); result != vulkan.Success {
		return vulkan.NullImageView, fmt.Errorf("failed to create image view: %s", result)
	}

	return imageView, nil
}

// TransitionLayout transitions the image layout
func (img *ImageVK) TransitionLayout(commandBuffer vulkan.CommandBuffer, oldLayout, newLayout vulkan.ImageLayout, aspectMask vulkan.ImageAspectFlags) {
	barrier := vulkan.ImageMemoryBarrier{
		SType:               vulkan.StructureTypeImageMemoryBarrier,
		OldLayout:           oldLayout,
		NewLayout:           newLayout,
		SrcQueueFamilyIndex: vulkan.QueueFamilyIgnored,
		DstQueueFamilyIndex: vulkan.QueueFamilyIgnored,
		Image:               img.image,
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     aspectMask,
			BaseMipLevel:   0,
			LevelCount:     img.mipLevels,
			BaseArrayLayer: 0,
			LayerCount:     img.arrayLayers,
		},
	}

	// Set access masks based on layouts
	sourceStage, destinationStage := getPipelineStageFlags(oldLayout, newLayout, &barrier)

	vulkan.CmdPipelineBarrier(
		commandBuffer,
		sourceStage,
		destinationStage,
		0,
		0, nil,
		0, nil,
		1, []vulkan.ImageMemoryBarrier{barrier})

	img.layout = newLayout
}

// Destroy destroys the image and its resources
func (img *ImageVK) Destroy() {
	if img.imageView != vulkan.NullImageView {
		vulkan.DestroyImageView(img.device, img.imageView, nil)
		img.imageView = vulkan.NullImageView
	}

	if img.owned && img.image != vulkan.NullImage {
		vulkan.DestroyImage(img.device, img.image, nil)
		img.image = vulkan.NullImage
	}

	if img.memory != vulkan.NullDeviceMemory {
		vulkan.FreeMemory(img.device, img.memory, nil)
		img.memory = vulkan.NullDeviceMemory
	}
}

// getImageType determines the image type based on dimensions
func getImageType(width, height, depth uint32) vulkan.ImageType {
	if depth > 1 {
		return vulkan.ImageType3d
	} else if height > 1 {
		return vulkan.ImageType2d
	} else {
		return vulkan.ImageType1d
	}
}

// getImageAspectMask returns the appropriate aspect mask for a format
func getImageAspectMask(format vulkan.Format) vulkan.ImageAspectFlags {
	switch format {
	case vulkan.FormatD16Unorm, vulkan.FormatD32Sfloat:
		return vulkan.ImageAspectFlags(vulkan.ImageAspectDepthBit)
	case vulkan.FormatS8Uint:
		return vulkan.ImageAspectFlags(vulkan.ImageAspectStencilBit)
	case vulkan.FormatD16UnormS8Uint, vulkan.FormatD24UnormS8Uint, vulkan.FormatD32SfloatS8Uint:
		return vulkan.ImageAspectFlags(vulkan.ImageAspectDepthBit | vulkan.ImageAspectStencilBit)
	default:
		return vulkan.ImageAspectFlags(vulkan.ImageAspectColorBit)
	}
}

// getPipelineStageFlags returns pipeline stage flags for layout transition
func getPipelineStageFlags(oldLayout, newLayout vulkan.ImageLayout, barrier *vulkan.ImageMemoryBarrier) (vulkan.PipelineStageFlags, vulkan.PipelineStageFlags) {
	var sourceStage, destinationStage vulkan.PipelineStageFlags

	switch oldLayout {
	case vulkan.ImageLayoutUndefined:
		barrier.SrcAccessMask = 0
		sourceStage = vulkan.PipelineStageFlags(vulkan.PipelineStageTopOfPipeBit)
	case vulkan.ImageLayoutColorAttachmentOptimal:
		barrier.SrcAccessMask = vulkan.AccessFlags(vulkan.AccessColorAttachmentWriteBit)
		sourceStage = vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit)
	case vulkan.ImageLayoutDepthStencilAttachmentOptimal:
		barrier.SrcAccessMask = vulkan.AccessFlags(vulkan.AccessDepthStencilAttachmentWriteBit)
		sourceStage = vulkan.PipelineStageFlags(vulkan.PipelineStageEarlyFragmentTestsBit)
	case vulkan.ImageLayoutShaderReadOnlyOptimal:
		barrier.SrcAccessMask = vulkan.AccessFlags(vulkan.AccessShaderReadBit)
		sourceStage = vulkan.PipelineStageFlags(vulkan.PipelineStageFragmentShaderBit)
	case vulkan.ImageLayoutTransferSrcOptimal:
		barrier.SrcAccessMask = vulkan.AccessFlags(vulkan.AccessTransferReadBit)
		sourceStage = vulkan.PipelineStageFlags(vulkan.PipelineStageTransferBit)
	case vulkan.ImageLayoutTransferDstOptimal:
		barrier.SrcAccessMask = vulkan.AccessFlags(vulkan.AccessTransferWriteBit)
		sourceStage = vulkan.PipelineStageFlags(vulkan.PipelineStageTransferBit)
	case vulkan.ImageLayoutPresentSrc:
		barrier.SrcAccessMask = 0
		sourceStage = vulkan.PipelineStageFlags(vulkan.PipelineStageBottomOfPipeBit)
	default:
		sourceStage = vulkan.PipelineStageFlags(vulkan.PipelineStageAllCommandsBit)
	}

	switch newLayout {
	case vulkan.ImageLayoutColorAttachmentOptimal:
		barrier.DstAccessMask = vulkan.AccessFlags(vulkan.AccessColorAttachmentWriteBit)
		destinationStage = vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit)
	case vulkan.ImageLayoutDepthStencilAttachmentOptimal:
		barrier.DstAccessMask = vulkan.AccessFlags(vulkan.AccessDepthStencilAttachmentWriteBit)
		destinationStage = vulkan.PipelineStageFlags(vulkan.PipelineStageEarlyFragmentTestsBit)
	case vulkan.ImageLayoutShaderReadOnlyOptimal:
		barrier.DstAccessMask = vulkan.AccessFlags(vulkan.AccessShaderReadBit)
		destinationStage = vulkan.PipelineStageFlags(vulkan.PipelineStageFragmentShaderBit)
	case vulkan.ImageLayoutTransferSrcOptimal:
		barrier.DstAccessMask = vulkan.AccessFlags(vulkan.AccessTransferReadBit)
		destinationStage = vulkan.PipelineStageFlags(vulkan.PipelineStageTransferBit)
	case vulkan.ImageLayoutTransferDstOptimal:
		barrier.DstAccessMask = vulkan.AccessFlags(vulkan.AccessTransferWriteBit)
		destinationStage = vulkan.PipelineStageFlags(vulkan.PipelineStageTransferBit)
	case vulkan.ImageLayoutPresentSrc:
		barrier.DstAccessMask = 0
		destinationStage = vulkan.PipelineStageFlags(vulkan.PipelineStageBottomOfPipeBit)
	default:
		destinationStage = vulkan.PipelineStageFlags(vulkan.PipelineStageAllCommandsBit)
	}

	return sourceStage, destinationStage
}

// ImageManagerVK manages image lifecycle
// It provides high-level image management
type ImageManagerVK struct {
	device    vulkan.Device
	allocator *AllocatorVK
	images    []*ImageVK
}

// NewImageManagerVK creates a new image manager
func NewImageManagerVK(device vulkan.Device, allocator *AllocatorVK) *ImageManagerVK {
	return &ImageManagerVK{
		device:    device,
		allocator: allocator,
		images:    make([]*ImageVK, 0),
	}
}

// CreateImage creates a new image
func (im *ImageManagerVK) CreateImage(createInfo *ImageCreateInfo) (*ImageVK, error) {
	image, err := NewImageVK(im.device, im.allocator, createInfo)
	if err != nil {
		return nil, err
	}

	im.images = append(im.images, image)
	return image, nil
}

// CreateImage2D creates a 2D image
func (im *ImageManagerVK) CreateImage2D(width, height uint32, format vulkan.Format, usage vulkan.ImageUsageFlags, properties vulkan.MemoryPropertyFlags) (*ImageVK, error) {
	createInfo := &ImageCreateInfo{
		Width:       width,
		Height:      height,
		Depth:       1,
		MipLevels:   1,
		ArrayLayers: 1,
		Format:      format,
		Tiling:      vulkan.ImageTilingOptimal,
		Usage:       usage,
		Samples:     vulkan.SampleCountFlags(vulkan.SampleCount1Bit),
		Properties:  properties,
	}

	return im.CreateImage(createInfo)
}

// GetImageCount returns the number of managed images
func (im *ImageManagerVK) GetImageCount() int {
	return len(im.images)
}

// Destroy destroys all managed images
func (im *ImageManagerVK) Destroy() {
	for _, image := range im.images {
		image.Destroy()
	}
	im.images = im.images[:0]
}
