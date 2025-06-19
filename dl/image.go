package dl

import (
	"fmt"
	"image"

	"github.com/opensraph/sraph/geom"
)

// Image represents a raster image that can be drawn on a canvas.
// It provides an abstraction over different image formats and backends.
type Image interface {
	// Width returns the width of the image in pixels.
	Width() int
	// Height returns the height of the image in pixels.
	Height() int
	// Size returns the size of the image.
	Size() geom.Size[geom.I32]
	// Bounds returns the bounds of the image.
	Bounds() geom.Rect[geom.I32]
	// IsOpaque returns true if the image is fully opaque.
	IsOpaque() bool
	// IsTextureBacked returns true if the image is backed by a GPU texture.
	IsTextureBacked() bool
	// ColorSpace returns the color space of the image.
	ColorSpace() ColorSpace
	// String returns a string representation of the image.
	String() string
}

// ImageFormat represents different image formats.
type ImageFormat uint8

const (
	// ImageFormatUnknown represents an unknown image format.
	ImageFormatUnknown ImageFormat = iota
	// ImageFormatRGBA represents RGBA format.
	ImageFormatRGBA
	// ImageFormatBGRA represents BGRA format.
	ImageFormatBGRA
	// ImageFormatRGB represents RGB format.
	ImageFormatRGB
	// ImageFormatGray represents grayscale format.
	ImageFormatGray
	// ImageFormatAlpha represents alpha-only format.
	ImageFormatAlpha
)

// BasicImage represents a basic CPU-backed image implementation.
type BasicImage struct {
	width      int
	height     int
	format     ImageFormat
	colorSpace ColorSpace
	data       []byte
	isOpaque   bool
}

// NewBasicImage creates a new basic image with the specified properties.
func NewBasicImage(width, height int, format ImageFormat, colorSpace ColorSpace, data []byte) *BasicImage {
	return &BasicImage{
		width:      width,
		height:     height,
		format:     format,
		colorSpace: colorSpace,
		data:       data,
		isOpaque:   determineOpacity(data, format),
	}
}

// NewBasicImageFromGoImage creates a basic image from a Go image.Image.
func NewBasicImageFromGoImage(img image.Image) *BasicImage {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Convert to RGBA format
	data := make([]byte, width*height*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			offset := (y*width + x) * 4
			data[offset] = byte(r >> 8)
			data[offset+1] = byte(g >> 8)
			data[offset+2] = byte(b >> 8)
			data[offset+3] = byte(a >> 8)
		}
	}

	return NewBasicImage(width, height, ImageFormatRGBA, ColorSpaceSRGB, data)
}

// Width implements Image interface.
func (img *BasicImage) Width() int {
	return img.width
}

// Height implements Image interface.
func (img *BasicImage) Height() int {
	return img.height
}

// Size implements Image interface.
func (img *BasicImage) Size() geom.Size[geom.I32] {
	return geom.Size[geom.I32]{Width: geom.I32(img.width), Height: geom.I32(img.height)}
}

// Bounds implements Image interface.
func (img *BasicImage) Bounds() geom.Rect[geom.I32] {
	return geom.NewRect(0, 0, geom.I32(img.width), geom.I32(img.height))
}

// IsOpaque implements Image interface.
func (img *BasicImage) IsOpaque() bool {
	return img.isOpaque
}

// IsTextureBacked implements Image interface.
func (img *BasicImage) IsTextureBacked() bool {
	return false // BasicImage is CPU-backed
}

// ColorSpace implements Image interface.
func (img *BasicImage) ColorSpace() ColorSpace {
	return img.colorSpace
}

// String implements Image interface.
func (img *BasicImage) String() string {
	return fmt.Sprintf("BasicImage{%dx%d, format: %d, opaque: %t}",
		img.width, img.height, img.format, img.isOpaque)
}

// Format returns the image format.
func (img *BasicImage) Format() ImageFormat {
	return img.format
}

// Data returns the raw image data.
func (img *BasicImage) Data() []byte {
	return img.data
}

// BytesPerPixel returns the number of bytes per pixel for the format.
func (img *BasicImage) BytesPerPixel() int {
	switch img.format {
	case ImageFormatRGBA, ImageFormatBGRA:
		return 4
	case ImageFormatRGB:
		return 3
	case ImageFormatGray, ImageFormatAlpha:
		return 1
	default:
		return 0
	}
}

// determineOpacity analyzes image data to determine if it's fully opaque.
func determineOpacity(data []byte, format ImageFormat) bool {
	switch format {
	case ImageFormatRGBA, ImageFormatBGRA:
		// Check alpha channel (every 4th byte for RGBA, every 4th byte for BGRA)
		for i := 3; i < len(data); i += 4 {
			if data[i] < 255 {
				return false
			}
		}
		return true
	case ImageFormatRGB, ImageFormatGray:
		// RGB and Gray formats don't have alpha, so they're opaque
		return true
	case ImageFormatAlpha:
		// Alpha-only format, check if all values are 255
		for _, alpha := range data {
			if alpha < 255 {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// TextureImage represents a GPU texture-backed image.
type TextureImage struct {
	width      int
	height     int
	colorSpace ColorSpace
	textureID  uint32
	isOpaque   bool
}

// NewTextureImage creates a new texture-backed image.
func NewTextureImage(width, height int, colorSpace ColorSpace, textureID uint32, isOpaque bool) *TextureImage {
	return &TextureImage{
		width:      width,
		height:     height,
		colorSpace: colorSpace,
		textureID:  textureID,
		isOpaque:   isOpaque,
	}
}

// Width implements Image interface.
func (img *TextureImage) Width() int {
	return img.width
}

// Height implements Image interface.
func (img *TextureImage) Height() int {
	return img.height
}

// Size implements Image interface.
func (img *TextureImage) Size() geom.Size[geom.I32] {
	return geom.Size[geom.I32]{Width: geom.I32(img.width), Height: geom.I32(img.height)}
}

// Bounds implements Image interface.
func (img *TextureImage) Bounds() geom.Rect[geom.I32] {
	return geom.NewRect(0, 0, geom.I32(img.width), geom.I32(img.height))
}

// IsOpaque implements Image interface.
func (img *TextureImage) IsOpaque() bool {
	return img.isOpaque
}

// IsTextureBacked implements Image interface.
func (img *TextureImage) IsTextureBacked() bool {
	return true
}

// ColorSpace implements Image interface.
func (img *TextureImage) ColorSpace() ColorSpace {
	return img.colorSpace
}

// String implements Image interface.
func (img *TextureImage) String() string {
	return fmt.Sprintf("TextureImage{%dx%d, textureID: %d, opaque: %t}",
		img.width, img.height, img.textureID, img.isOpaque)
}

// TextureID returns the GPU texture ID.
func (img *TextureImage) TextureID() uint32 {
	return img.textureID
}

// ImageDecoder provides functionality to decode images from various sources.
type ImageDecoder struct {
	// TODO: Add decoder implementation
}

// DecodeFromBytes decodes an image from byte data.
func (d *ImageDecoder) DecodeFromBytes(data []byte) (Image, error) {
	// TODO: Implement image decoding from bytes
	return nil, fmt.Errorf("image decoding not implemented")
}

// DecodeFromFile decodes an image from a file.
func (d *ImageDecoder) DecodeFromFile(filename string) (Image, error) {
	// TODO: Implement image decoding from file
	return nil, fmt.Errorf("image decoding not implemented")
}
