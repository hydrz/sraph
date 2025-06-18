package geom

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
)

// BlendMode represents different blending modes for color composition.
type BlendMode uint8

const (
	// Basic blend modes (compatible with most graphics pipelines)
	BlendModeClear BlendMode = iota
	BlendModeSrc
	BlendModeDst
	BlendModeSrcOver
	BlendModeDstOver
	BlendModeSrcIn
	BlendModeDstIn
	BlendModeSrcOut
	BlendModeDstOut
	BlendModeSrcATop
	BlendModeDstATop
	BlendModeXor
	BlendModePlus
	BlendModeModulate

	// Advanced blend modes (require special handling)
	BlendModeScreen
	BlendModeOverlay
	BlendModeDarken
	BlendModeLighten
	BlendModeColorDodge
	BlendModeColorBurn
	BlendModeHardLight
	BlendModeSoftLight
	BlendModeDifference
	BlendModeExclusion
	BlendModeMultiply
	BlendModeHue
	BlendModeSaturation
	BlendModeColor
	BlendModeLuminosity

	BlendModeLastMode    = BlendModeLuminosity
	BlendModeDefaultMode = BlendModeSrcOver
)

// String returns the string representation of BlendMode.
func (b BlendMode) String() string {
	names := []string{
		"Clear", "Src", "Dst", "SrcOver", "DstOver", "SrcIn", "DstIn",
		"SrcOut", "DstOut", "SrcATop", "DstATop", "Xor", "Plus", "Modulate",
		"Screen", "Overlay", "Darken", "Lighten", "ColorDodge", "ColorBurn",
		"HardLight", "SoftLight", "Difference", "Exclusion", "Multiply",
		"Hue", "Saturation", "Color", "Luminosity",
	}
	if int(b) < len(names) {
		return names[b]
	}
	return fmt.Sprintf("BlendMode(%d)", b)
}

// ColorMatrix represents a 4x5 matrix for color transformation.
// Matrix format: [R', G', B', A'] = [R, G, B, A, 1] * Matrix
type ColorMatrix [20]F32

// Color represents an RGBA color with components in the range [0, 1].
// It implements Go's standard color.Color interface.
type Color struct {
	R, G, B, A F32
}

// NewColor creates a new Color with the given RGBA components.
func NewColor(r, g, b, a F32) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// NewColorFromRGBA creates a Color from Go's standard color.RGBA.
func NewColorFromRGBA(c color.RGBA) Color {
	return Color{
		R: F32(c.R) / 255.0,
		G: F32(c.G) / 255.0,
		B: F32(c.B) / 255.0,
		A: F32(c.A) / 255.0,
	}
}

// NewColorFromColor creates a Color from any color.Color.
func NewColorFromColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	return Color{
		R: F32(r) / 65535.0,
		G: F32(g) / 65535.0,
		B: F32(b) / 65535.0,
		A: F32(a) / 65535.0,
	}
}

// RGBA implements color.Color interface.
// Returns the alpha-premultiplied red, green, blue and alpha values.
func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R*65535 + 0.5)
	g = uint32(c.G*65535 + 0.5)
	b = uint32(c.B*65535 + 0.5)
	a = uint32(c.A*65535 + 0.5)
	return
}

// ToRGBA converts Color to Go's standard color.RGBA.
func (c Color) ToRGBA() color.RGBA {
	return color.RGBA{
		R: uint8(Clamp(c.R*255+0.5, 0, 255)),
		G: uint8(Clamp(c.G*255+0.5, 0, 255)),
		B: uint8(Clamp(c.B*255+0.5, 0, 255)),
		A: uint8(Clamp(c.A*255+0.5, 0, 255)),
	}
}

// String returns a string representation of the color.
func (c Color) String() string {
	return fmt.Sprintf("Color(%.3f, %.3f, %.3f, %.3f)", c.R, c.G, c.B, c.A)
}

// Eq compares two colors for Eqity within tolerance.
func (c Color) Eq(other Color) bool {
	return NearlyEq(c.R, other.R) &&
		NearlyEq(c.G, other.G) &&
		NearlyEq(c.B, other.B) &&
		NearlyEq(c.A, other.A)
}

// Add performs component-wise addition.
func (c Color) Add(other Color) Color {
	return Color{
		R: c.R + other.R,
		G: c.G + other.G,
		B: c.B + other.B,
		A: c.A + other.A,
	}
}

// Sub performs component-wise subtraction.
func (c Color) Sub(other Color) Color {
	return Color{
		R: c.R - other.R,
		G: c.G - other.G,
		B: c.B - other.B,
		A: c.A - other.A,
	}
}

// Mul performs component-wise multiplication.
func (c Color) Mul(other Color) Color {
	return Color{
		R: c.R * other.R,
		G: c.G * other.G,
		B: c.B * other.B,
		A: c.A * other.A,
	}
}

// Scale multiplies all components by a scalar.
func (c Color) Scale(scale F32) Color {
	return Color{
		R: c.R * scale,
		G: c.G * scale,
		B: c.B * scale,
		A: c.A * scale,
	}
}

// Clamp01 clamps all color components to the range [0, 1].
func (c Color) Clamp01() Color {
	return Color{
		R: Clamp(c.R, 0, 1),
		G: Clamp(c.G, 0, 1),
		B: Clamp(c.B, 0, 1),
		A: Clamp(c.A, 0, 1),
	}
}

// Premultiply returns the color with RGB premultiplied by alpha.
func (c Color) Premultiply() Color {
	return Color{
		R: c.R * c.A,
		G: c.G * c.A,
		B: c.B * c.A,
		A: c.A,
	}
}

// Unpremultiply returns the color with RGB unpremultiplied by alpha.
func (c Color) Unpremultiply() Color {
	if c.A <= 0 {
		return Color{}
	}
	return Color{
		R: c.R / c.A,
		G: c.G / c.A,
		B: c.B / c.A,
		A: c.A,
	}
}

// WithAlpha returns a new color with the specified alpha value.
func (c Color) WithAlpha(alpha F32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: alpha}
}

// IsTransparent returns true if the alpha component is zero.
func (c Color) IsTransparent() bool {
	return c.A == 0
}

// IsOpaque returns true if the alpha component is one.
func (c Color) IsOpaque() bool {
	return c.A == 1
}

// Lerp performs linear interpolation between two colors.
func (c Color) Lerp(other Color, t F32) Color {
	return Color{
		R: c.R + (other.R-c.R)*t,
		G: c.G + (other.G-c.G)*t,
		B: c.B + (other.B-c.B)*t,
		A: c.A + (other.A-c.A)*t,
	}
}

// ApplyColorMatrix applies a color transformation matrix.
func (c Color) ApplyColorMatrix(matrix ColorMatrix) Color {
	m := matrix
	return Color{
		R: m[0]*c.R + m[1]*c.G + m[2]*c.B + m[3]*c.A + m[4],
		G: m[5]*c.R + m[6]*c.G + m[7]*c.B + m[8]*c.A + m[9],
		B: m[10]*c.R + m[11]*c.G + m[12]*c.B + m[13]*c.A + m[14],
		A: m[15]*c.R + m[16]*c.G + m[17]*c.B + m[18]*c.A + m[19],
	}.Clamp01()
}

// Blend blends this color with another using the specified blend mode.
func (c Color) Blend(src Color, mode BlendMode) Color {
	dst := c
	switch mode {
	case BlendModeClear:
		return Color{}
	case BlendModeSrc:
		return src
	case BlendModeDst:
		return dst
	case BlendModeSrcOver:
		// r = s + (1-sa)*d
		return src.Premultiply().Add(dst.Premultiply().Scale(1 - src.A)).Unpremultiply()
	case BlendModeDstOver:
		// r = d + (1-da)*s
		return dst.Premultiply().Add(src.Premultiply().Scale(1 - dst.A)).Unpremultiply()
	case BlendModeSrcIn:
		// r = s * da
		return src.Premultiply().Scale(dst.A).Unpremultiply()
	case BlendModeDstIn:
		// r = d * sa
		return dst.Premultiply().Scale(src.A).Unpremultiply()
	case BlendModeMultiply:
		return Color{
			R: dst.R * src.R,
			G: dst.G * src.G,
			B: dst.B * src.B,
			A: dst.A * src.A,
		}
	case BlendModeScreen:
		return Color{
			R: src.R + dst.R - src.R*dst.R,
			G: src.G + dst.G - src.G*dst.G,
			B: src.B + dst.B - src.B*dst.B,
			A: src.A + dst.A - src.A*dst.A,
		}
	default:
		// For unsupported modes, default to source over
		return c.Blend(src, BlendModeSrcOver)
	}
}

// Predefined colors
var (
	ColorWhite            = Color{R: 1, G: 1, B: 1, A: 1}
	ColorBlack            = Color{R: 0, G: 0, B: 0, A: 1}
	ColorTransparent      = Color{R: 0, G: 0, B: 0, A: 0}
	ColorWhiteTransparent = Color{R: 1, G: 1, B: 1, A: 0}
	ColorRed              = Color{R: 1, G: 0, B: 0, A: 1}
	ColorGreen            = Color{R: 0, G: 1, B: 0, A: 1}
	ColorBlue             = Color{R: 0, G: 0, B: 1, A: 1}
	ColorYellow           = Color{R: 1, G: 1, B: 0, A: 1}
	ColorCyan             = Color{R: 0, G: 1, B: 1, A: 1}
	ColorMagenta          = Color{R: 1, G: 0, B: 1, A: 1}
)

// NewColorRGBA8 creates a Color from 8-bit RGBA values.
func NewColorRGBA8(r, g, b, a uint8) Color {
	return Color{
		R: F32(r) / 255.0,
		G: F32(g) / 255.0,
		B: F32(b) / 255.0,
		A: F32(a) / 255.0,
	}
}

// NewColorRGB8 creates an opaque Color from 8-bit RGB values.
func NewColorRGB8(r, g, b uint8) Color {
	return NewColorRGBA8(r, g, b, 255)
}

// NewColorHex creates a Color from a hexadecimal color value.
func NewColorHex(hex uint32) Color {
	return NewColorRGBA8(
		uint8((hex>>16)&0xFF),
		uint8((hex>>8)&0xFF),
		uint8(hex&0xFF),
		255,
	)
}

// NewColorHexA creates a Color from a hexadecimal color value with alpha.
func NewColorHexA(hex uint32) Color {
	return NewColorRGBA8(
		uint8((hex>>24)&0xFF),
		uint8((hex>>16)&0xFF),
		uint8((hex>>8)&0xFF),
		uint8(hex&0xFF),
	)
}

// RandomColor generates a random opaque color.
func RandomColor() Color {
	return Color{
		R: F32(rand.Float32()),
		G: F32(rand.Float32()),
		B: F32(rand.Float32()),
		A: 1.0,
	}
}

// LinearToSRGB converts color from linear space to sRGB space.
func (c Color) LinearToSRGB() Color {
	convert := func(component F32) F32 {
		if component <= 0.0031308 {
			return component * 12.92
		}
		return F32(1.055*math.Pow(float64(component), 1.0/2.4) - 0.055)
	}

	return Color{
		R: convert(c.R),
		G: convert(c.G),
		B: convert(c.B),
		A: c.A,
	}
}

// SRGBToLinear converts color from sRGB space to linear space.
func (c Color) SRGBToLinear() Color {
	convert := func(component F32) F32 {
		if component <= 0.04045 {
			return component / 12.92
		}
		return F32(math.Pow((float64(component)+0.055)/1.055, 2.4))
	}

	return Color{
		R: convert(c.R),
		G: convert(c.G),
		B: convert(c.B),
		A: c.A,
	}
}
