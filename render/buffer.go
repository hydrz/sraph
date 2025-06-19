package render

// Buffer represents a GPU buffer resource
type Buffer interface {
	Resource

	// GetSize returns the size of the buffer in bytes
	GetSize() int

	// GetUsage returns the buffer usage flags
	GetUsage() BufferUsage

	// Map maps the buffer for CPU access
	Map() ([]byte, error)

	// Unmap unmaps the buffer from CPU access
	Unmap() error

	// CopyFromBuffer copies data from another buffer
	CopyFromBuffer(source Buffer, sourceOffset, destOffset, size int) error
}

// BufferUsage defines how a buffer can be used
type BufferUsage int

const (
	BufferUsageVertex BufferUsage = 1 << iota
	BufferUsageIndex
	BufferUsageUniform
	BufferUsageStorage
	BufferUsageTransferSrc
	BufferUsageTransferDst
)

// BufferImpl is the default implementation of Buffer
type BufferImpl struct {
	*ResourceImpl
	size  int
	usage BufferUsage
	data  []byte
}

// NewBuffer creates a new buffer with the specified size and usage
func NewBuffer(size int, usage BufferUsage) Buffer {
	return &BufferImpl{
		ResourceImpl: &ResourceImpl{
			label:   "",
			isValid: true,
		},
		size:  size,
		usage: usage,
		data:  make([]byte, size),
	}
}

// GetSize returns the size of the buffer in bytes
func (b *BufferImpl) GetSize() int {
	return b.size
}

// GetUsage returns the buffer usage flags
func (b *BufferImpl) GetUsage() BufferUsage {
	return b.usage
}

// Map maps the buffer for CPU access
func (b *BufferImpl) Map() ([]byte, error) {
	if !b.IsValid() {
		return nil, ErrInvalidState
	}
	return b.data, nil
}

// Unmap unmaps the buffer from CPU access
func (b *BufferImpl) Unmap() error {
	// TODO: Implement buffer unmapping
	return nil
}

// CopyFromBuffer copies data from another buffer
func (b *BufferImpl) CopyFromBuffer(source Buffer, sourceOffset, destOffset, size int) error {
	if !b.IsValid() || source == nil || !source.IsValid() {
		return ErrInvalidState
	}

	sourceData, err := source.Map()
	if err != nil {
		return err
	}

	if sourceOffset+size > len(sourceData) || destOffset+size > len(b.data) {
		return ErrInvalidArgument
	}

	copy(b.data[destOffset:destOffset+size], sourceData[sourceOffset:sourceOffset+size])
	return nil
}

// Texture represents a GPU texture resource
type TextureType int

const (
	TextureType1D TextureType = iota
	TextureType2D
	TextureType3D
	TextureTypeCube
)

// PixelFormat defines the pixel format of a texture
type PixelFormat int

const (
	PixelFormatRGBA8Unorm PixelFormat = iota
	PixelFormatBGRA8Unorm
	PixelFormatR8Unorm
	PixelFormatRG8Unorm
	PixelFormatR16Float
	PixelFormatRG16Float
	PixelFormatRGBA16Float
	PixelFormatR32Float
	PixelFormatRG32Float
	PixelFormatRGBA32Float
	PixelFormatDepth32Float
	PixelFormatStencil8
)

// TextureDescriptor describes the configuration for a texture
type TextureDescriptor struct {
	Type        TextureType
	Format      PixelFormat
	Width       int
	Height      int
	Depth       int
	MipLevels   int
	SampleCount int
	Usage       TextureUsage
	StorageMode StorageMode
}

// TextureUsage defines how a texture can be used
type TextureUsage int

const (
	TextureUsageShaderRead TextureUsage = 1 << iota
	TextureUsageShaderWrite
	TextureUsageRenderTarget
	TextureUsageTransferSrc
	TextureUsageTransferDst
)

// Texture represents a GPU texture resource
type Texture interface {
	Resource

	// GetType returns the texture type
	GetType() TextureType

	// GetFormat returns the pixel format
	GetFormat() PixelFormat

	// GetSize returns the texture dimensions
	GetSize() (int, int, int)

	// GetMipLevels returns the number of mip levels
	GetMipLevels() int

	// GetSampleCount returns the sample count for MSAA
	GetSampleCount() int

	// GetUsage returns the texture usage flags
	GetUsage() TextureUsage

	// CopyFromTexture copies data from another texture
	CopyFromTexture(source Texture) error
}

// TextureImpl is the default implementation of Texture
type TextureImpl struct {
	*ResourceImpl
	descriptor TextureDescriptor
}

// NewTexture creates a new texture with the specified descriptor
func NewTexture(descriptor TextureDescriptor) Texture {
	return &TextureImpl{
		ResourceImpl: &ResourceImpl{
			label:   "",
			isValid: true,
		},
		descriptor: descriptor,
	}
}

// GetType returns the texture type
func (t *TextureImpl) GetType() TextureType {
	return t.descriptor.Type
}

// GetFormat returns the pixel format
func (t *TextureImpl) GetFormat() PixelFormat {
	return t.descriptor.Format
}

// GetSize returns the texture dimensions
func (t *TextureImpl) GetSize() (int, int, int) {
	return t.descriptor.Width, t.descriptor.Height, t.descriptor.Depth
}

// GetMipLevels returns the number of mip levels
func (t *TextureImpl) GetMipLevels() int {
	return t.descriptor.MipLevels
}

// GetSampleCount returns the sample count for MSAA
func (t *TextureImpl) GetSampleCount() int {
	return t.descriptor.SampleCount
}

// GetUsage returns the texture usage flags
func (t *TextureImpl) GetUsage() TextureUsage {
	return t.descriptor.Usage
}

// CopyFromTexture copies data from another texture
func (t *TextureImpl) CopyFromTexture(source Texture) error {
	if !t.IsValid() || source == nil || !source.IsValid() {
		return ErrInvalidState
	}

	// TODO: Implement texture copying
	return nil
}
