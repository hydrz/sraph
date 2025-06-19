package render

import (
	"fmt"
	"math"
)

// TextureUtil provides utility functions for texture operations
type TextureUtil struct{}

// CalculateMipLevels calculates the number of mip levels for given dimensions
func CalculateMipLevels(width, height uint32) uint32 {
	if width == 0 || height == 0 {
		return 0
	}

	maxDimension := width
	if height > maxDimension {
		maxDimension = height
	}

	return uint32(math.Floor(math.Log2(float64(maxDimension)))) + 1
}

// CalculateMipDimensions calculates the dimensions for a specific mip level
func CalculateMipDimensions(width, height uint32, level uint32) (uint32, uint32) {
	mipWidth := width >> level
	mipHeight := height >> level

	// Ensure minimum size of 1x1
	if mipWidth == 0 {
		mipWidth = 1
	}
	if mipHeight == 0 {
		mipHeight = 1
	}

	return mipWidth, mipHeight
}

// CalculateTextureSize calculates the total size in bytes for a texture
func CalculateTextureSize(width, height uint32, format PixelFormat, mipLevels uint32) uint32 {
	if mipLevels == 0 {
		mipLevels = CalculateMipLevels(width, height)
	}

	bytesPerPixel := GetPixelFormatBytesPerPixel(format)
	var totalSize uint32 = 0

	for level := uint32(0); level < mipLevels; level++ {
		mipWidth, mipHeight := CalculateMipDimensions(width, height, level)
		mipSize := mipWidth * mipHeight * bytesPerPixel
		totalSize += mipSize
	}

	return totalSize
}

// GetPixelFormatBytesPerPixel returns the number of bytes per pixel for a format
func GetPixelFormatBytesPerPixel(format PixelFormat) uint32 {
	switch format {
	case PixelFormatR8Unorm:
		return 1
	case PixelFormatRG8Unorm:
		return 2
	case PixelFormatRGBA8Unorm:
		return 4
	case PixelFormatBGRA8Unorm:
		return 4
	case PixelFormatR16Float:
		return 2
	case PixelFormatRG16Float:
		return 4
	case PixelFormatRGBA16Float:
		return 8
	case PixelFormatR32Float:
		return 4
	case PixelFormatRG32Float:
		return 8
	case PixelFormatRGBA32Float:
		return 16
	case PixelFormatDepth32Float:
		return 4
	case PixelFormatStencil8:
		return 1
	default:
		return 4 // Default to 4 bytes for unknown formats
	}
}

// IsDepthFormat returns true if the format is a depth format
func IsDepthFormat(format PixelFormat) bool {
	switch format {
	case PixelFormatDepth32Float:
		return true
	default:
		return false
	}
}

// IsStencilFormat returns true if the format is a stencil format
func IsStencilFormat(format PixelFormat) bool {
	switch format {
	case PixelFormatStencil8:
		return true
	default:
		return false
	}
}

// IsColorFormat returns true if the format is a color format
func IsColorFormat(format PixelFormat) bool {
	return !IsDepthFormat(format) && !IsStencilFormat(format)
}

// ValidateTextureDimensions validates texture dimensions
func ValidateTextureDimensions(width, height uint32, textureType TextureType) error {
	if width == 0 || height == 0 {
		return fmt.Errorf("texture dimensions must be greater than zero")
	}

	// Check for power-of-two requirements for certain operations
	if !IsPowerOfTwo(width) || !IsPowerOfTwo(height) {
		// Warning: Non-power-of-two textures may have limitations
	}

	// Check maximum texture size (platform-dependent)
	const maxTextureSize = 16384 // Common maximum
	if width > maxTextureSize || height > maxTextureSize {
		return fmt.Errorf("texture dimensions exceed maximum size of %d", maxTextureSize)
	}

	return nil
}

// IsPowerOfTwo checks if a number is a power of two
func IsPowerOfTwo(n uint32) bool {
	return n > 0 && (n&(n-1)) == 0
}

// NextPowerOfTwo returns the next power of two greater than or equal to n
func NextPowerOfTwo(n uint32) uint32 {
	if n == 0 {
		return 1
	}

	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++

	return n
}

// TextureCreateInfo contains information for creating a texture
type TextureCreateInfo struct {
	Width       uint32
	Height      uint32
	Depth       uint32
	MipLevels   uint32
	Format      PixelFormat
	Type        TextureType
	Usage       TextureUsage
	SampleCount uint32
}

// Validate validates the texture creation info
func (info *TextureCreateInfo) Validate() error {
	if err := ValidateTextureDimensions(info.Width, info.Height, info.Type); err != nil {
		return err
	}

	if info.MipLevels == 0 {
		info.MipLevels = CalculateMipLevels(info.Width, info.Height)
	}

	if info.SampleCount == 0 {
		info.SampleCount = 1
	}

	return nil
}
