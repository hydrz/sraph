package dl

import (
	"fmt"

	"github.com/opensraph/sraph/geom"
)

// ColorSpace defines the color space for colors.
type ColorSpace uint8

const (
	// ColorSpaceSRGB represents the sRGB color space.
	ColorSpaceSRGB ColorSpace = iota
	// ColorSpaceExtendedSRGB represents the extended sRGB color space.
	ColorSpaceExtendedSRGB
	// ColorSpaceDisplayP3 represents the Display P3 color space.
	ColorSpaceDisplayP3
)

// Color represents a color in a specific color space.
// It wraps the geom.Color with additional display list specific functionality.
type Color struct {
	color      geom.Color
	colorSpace ColorSpace
}

// NewColor creates a new Color with the specified RGBA values in the sRGB color space.
func NewColor(r, g, b, a float32) Color {
	return Color{
		color:      geom.NewColor(geom.F32(r), geom.F32(g), geom.F32(b), geom.F32(a)),
		colorSpace: ColorSpaceSRGB,
	}
}

// NewColorFromGeom creates a new Color from a geom.Color in the sRGB color space.
func NewColorFromGeom(color geom.Color) Color {
	return Color{
		color:      color,
		colorSpace: ColorSpaceSRGB,
	}
}

// NewColorWithSpace creates a new Color with the specified RGBA values and color space.
func NewColorWithSpace(r, g, b, a float32, colorSpace ColorSpace) Color {
	return Color{
		color:      geom.NewColor(geom.F32(r), geom.F32(g), geom.F32(b), geom.F32(a)),
		colorSpace: colorSpace,
	}
}

//

// R returns the red component of the color.
func (c Color) R() float32 {
	return float32(c.color.R)
}

// G returns the green component of the color.
func (c Color) G() float32 {
	return float32(c.color.G)
}

// B returns the blue component of the color.
func (c Color) B() float32 {
	return float32(c.color.B)
}

// A returns the alpha component of the color.
func (c Color) A() float32 {
	return float32(c.color.A)
}

// ColorSpace returns the color space of the color.
func (c Color) ColorSpace() ColorSpace {
	return c.colorSpace
}

// GeomColor returns the underlying geom.Color.
func (c Color) GeomColor() geom.Color {
	return c.color
}

// WithAlpha returns a new Color with the specified alpha value.
func (c Color) WithAlpha(alpha float32) Color {
	return Color{
		color:      c.color.WithAlpha(geom.F32(alpha)),
		colorSpace: c.colorSpace,
	}
}

// IsTransparent returns true if the color is fully transparent.
func (c Color) IsTransparent() bool {
	return c.color.IsTransparent()
}

// IsOpaque returns true if the color is fully opaque.
func (c Color) IsOpaque() bool {
	return c.color.IsOpaque()
}

// String returns a string representation of the color.
func (c Color) String() string {
	return fmt.Sprintf("Color{RGBA: %s, ColorSpace: %d}", c.color.String(), c.colorSpace)
}

// Predefined colors in sRGB color space
var (
	// ColorTransparent represents a fully transparent color.
	ColorTransparent = NewColor(0, 0, 0, 0)
	// ColorBlack represents black color.
	ColorBlack = NewColor(0, 0, 0, 1)
	// ColorWhite represents white color.
	ColorWhite = NewColor(1, 1, 1, 1)
	// ColorRed represents red color.
	ColorRed = NewColor(1, 0, 0, 1)
	// ColorGreen represents green color.
	ColorGreen = NewColor(0, 1, 0, 1)
	// ColorBlue represents blue color.
	ColorBlue = NewColor(0, 0, 1, 1)
	// ColorYellow represents yellow color.
	ColorYellow = NewColor(1, 1, 0, 1)
	// ColorMagenta represents magenta color.
	ColorMagenta = NewColor(1, 0, 1, 1)
	// ColorCyan represents cyan color.
	ColorCyan = NewColor(0, 1, 1, 1)
)

// NewColorRGB creates a new Color with RGB values (alpha = 1.0).
func NewColorRGB(r, g, b float32) Color {
	return NewColor(r, g, b, 1.0)
}
