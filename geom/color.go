package geom

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"
)

// ColorMatrix represents a 4x5 matrix for color transformation.
// Each ColorMatrix can be used to transform RGBA colors in the range [0,1].
// When applied to a color [R, G, B, A], the result is clamped to [0,1] per channel.
type ColorMatrix [20]Scalar

// NewColorMatrix returns an identity ColorMatrix.
func NewColorMatrix() ColorMatrix {
	return ColorMatrix{
		1, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 1, 0,
	}
}

// Color represents an RGBA color with components in [0,1].
// Color implements color.Color. The zero value is fully transparent black.
type Color struct {
	R, G, B, A Scalar
}

// NewColor returns a Color with the given RGBA components in [0,1].
func NewColor(r, g, b, a Scalar) Color {
	return Color{
		R: Clamp(r, 0, 1),
		G: Clamp(g, 0, 1),
		B: Clamp(b, 0, 1),
		A: Clamp(a, 0, 1),
	}
}

// NewColorRGBA8 returns a Color from a 0-255 RGBA color.
func NewColorRGBA8(r, g, b, a uint8) Color {
	return Color{
		R: Scalar(r) / 255.0,
		G: Scalar(g) / 255.0,
		B: Scalar(b) / 255.0,
		A: Scalar(a) / 255.0,
	}
}

// NewColorHex returns a Color from a 32-bit RGBA hex value.
func NewColorHex(hex uint32) Color {
	return Color{
		R: Scalar((hex>>24)&0xff) / 255.0,
		G: Scalar((hex>>16)&0xff) / 255.0,
		B: Scalar((hex>>8)&0xff) / 255.0,
		A: Scalar(hex&0xff) / 255.0,
	}
}

// RandomColor returns a random Color with components in [0,1].
func RandomColor() Color {
	return Color{
		R: Scalar(rand.Float64()),
		G: Scalar(rand.Float64()),
		B: Scalar(rand.Float64()),
		A: 1,
	}.Clamp01()
}

// RGBA implements the [color.Color] interface.
func (c Color) RGBA() (r, g, b, a uint32) {
	// Convert components from [0,1] to [0,0xffff] range
	a = uint32(c.A * 0xffff)
	// Alpha-premultiply RGB components
	r = uint32(c.R * c.A * 0xffff)
	g = uint32(c.G * c.A * 0xffff)
	b = uint32(c.B * c.A * 0xffff)
	return r, g, b, a
}

func (c Color) Go() color.RGBA {
	return color.RGBA{
		R: uint8(math.Round(c.R * 255.0)),
		G: uint8(math.Round(c.G * 255.0)),
		B: uint8(math.Round(c.B * 255.0)),
		A: uint8(math.Round(c.A * 255.0)),
	}
}

// Hex returns the color as a 32-bit RGBA hex value.
func (c Color) Hex() uint32 {
	return (uint32(math.Round(c.R*255.0))&0xff)<<24 |
		(uint32(math.Round(c.G*255.0))&0xff)<<16 |
		(uint32(math.Round(c.B*255.0))&0xff)<<8 |
		(uint32(math.Round(c.A*255.0))&0xff)<<0
}

// ToIColor returns the color as a 32-bit ARGB hex value
func (c Color) ToIColor() uint32 {
	return (uint32(math.Round(c.A*255.0))&0xff)<<24 |
		(uint32(math.Round(c.R*255.0))&0xff)<<16 |
		(uint32(math.Round(c.G*255.0))&0xff)<<8 |
		(uint32(math.Round(c.B*255.0))&0xff)<<0
}

// Equal reports whether c and o are equal within floating-point tolerance.
func (c Color) Equal(o Color) bool {
	return NearlyEqual(c.R, o.R) &&
		NearlyEqual(c.G, o.G) &&
		NearlyEqual(c.B, o.B) &&
		NearlyEqual(c.A, o.A)
}

// Add returns the component-wise sum of c and o.
func (c Color) Add(o Color) Color {
	return Color{c.R + o.R, c.G + o.G, c.B + o.B, c.A + o.A}
}

// Sub returns the component-wise difference of c and o.
func (c Color) Sub(o Color) Color {
	return Color{c.R - o.R, c.G - o.G, c.B - o.B, c.A - o.A}
}

// Mul returns the component-wise product of c and o.
func (c Color) Mul(o Color) Color {
	return Color{c.R * o.R, c.G * o.G, c.B * o.B, c.A * o.A}
}

// Div returns the component-wise quotient of c and o.
func (c Color) Div(o Color) Color {
	return Color{
		R: c.R / Clamp(o.R, Epsilon64, 1),
		G: c.G / Clamp(o.R, Epsilon64, 1),
		B: c.B / Clamp(o.R, Epsilon64, 1),
		A: c.A / Clamp(o.R, Epsilon64, 1),
	}
}

// Scale returns c with all components multiplied by scale.
func (c Color) Scale(scale Scalar) Color {
	return Color{
		R: c.R * scale,
		G: c.G * scale,
		B: c.B * scale,
		A: c.A * scale,
	}
}

// Clamp01 returns c with all components clamped to [0,1].
func (c Color) Clamp01() Color {
	return Color{
		R: Clamp(c.R, 0, 1),
		G: Clamp(c.G, 0, 1),
		B: Clamp(c.B, 0, 1),
		A: Clamp(c.A, 0, 1),
	}
}

// Premultiply returns c with RGB premultiplied by alpha.
func (c Color) Premultiply() Color {
	return Color{
		R: c.R * c.A,
		G: c.G * c.A,
		B: c.B * c.A,
		A: c.A,
	}
}

// Unpremultiply returns c with RGB unpremultiplied by alpha.
// If alpha is zero or negative, returns fully transparent black.
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

// WithAlpha returns a copy of c with the given alpha value.
func (c Color) WithAlpha(alpha Scalar) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: alpha}
}

// IsTransparent reports whether alpha is zero.
func (c Color) IsTransparent() bool {
	return NearlyEqual(c.A, 0)
}

// IsOpaque reports whether alpha is one.
func (c Color) IsOpaque() bool {
	return NearlyEqual(c.A, 1)
}

// Lerp returns the linear interpolation between c and o by t in [0,1].
func (c Color) Lerp(o Color, t Scalar) Color {
	return Color{
		R: Clamp(c.R+(o.R-c.R)*t, 0, 1),
		G: Clamp(c.G+(o.G-c.G)*t, 0, 1),
		B: Clamp(c.B+(o.B-c.B)*t, 0, 1),
		A: Clamp(c.A+(o.A-c.A)*t, 0, 1),
	}
}

// LinearToSRGB returns c converted from linear to sRGB color space.
func (c Color) LinearToSRGB() Color {
	convert := func(component Scalar) Scalar {
		if component <= 0.0031308 {
			return component * 12.92
		}
		return Scalar(1.055*math.Pow(ToFloat64(component), 1.0/2.4) - 0.055)
	}
	return Color{
		R: convert(c.R),
		G: convert(c.G),
		B: convert(c.B),
		A: c.A,
	}
}

// SRGBToLinear returns c converted from sRGB to linear color space.
func (c Color) SRGBToLinear() Color {
	convert := func(component Scalar) Scalar {
		if component <= 0.04045 {
			return component / 12.92
		}
		return Scalar(math.Pow((ToFloat64(component)+0.055)/1.055, 2.4))
	}
	return Color{
		R: convert(c.R),
		G: convert(c.G),
		B: convert(c.B),
		A: c.A,
	}
}

// ApplyColorMatrix returns c transformed by the given ColorMatrix.
func (c Color) ApplyColorMatrix(matrix ColorMatrix) Color {
	m := matrix
	return Color{
		R: m[0]*c.R + m[1]*c.G + m[2]*c.B + m[3]*c.A + m[4],
		G: m[5]*c.R + m[6]*c.G + m[7]*c.B + m[8]*c.A + m[9],
		B: m[10]*c.R + m[11]*c.G + m[12]*c.B + m[13]*c.A + m[14],
		A: m[15]*c.R + m[16]*c.G + m[17]*c.B + m[18]*c.A + m[19],
	}.Clamp01()
}

// Blend returns the result of blending c with src using the given BlendMode.
// For unsupported modes, BlendModeSrcOver is used.
func (c Color) Blend(src Color, mode BlendMode) Color {
	dst := c
	const kEhCloseEnough = 1e-6

	switch mode {
	case BlendModeClear:
		return ColorBlackTransparent()
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
	case BlendModeSrcOut:
		// r = s * (1-da)
		return src.Premultiply().Scale(1 - dst.A).Unpremultiply()
	case BlendModeDstOut:
		// r = d * (1-sa)
		return dst.Premultiply().Scale(1 - src.A).Unpremultiply()
	case BlendModeSrcATop:
		// r = s*da + d*(1-sa)
		return src.Premultiply().Scale(dst.A).Add(dst.Premultiply().Scale(1 - src.A)).Unpremultiply()
	case BlendModeDstATop:
		// r = d*sa + s*(1-da)
		return dst.Premultiply().Scale(src.A).Add(src.Premultiply().Scale(1 - dst.A)).Unpremultiply()
	case BlendModeXor:
		// r = s*(1-da) + d*(1-sa)
		return src.Premultiply().Scale(1 - dst.A).Add(dst.Premultiply().Scale(1 - src.A)).Unpremultiply()
	case BlendModePlus:
		// r = min(s + d, 1)
		result := src.Premultiply().Add(dst.Premultiply())
		return Color{
			R: Scalar(math.Min(ToFloat64(result.R), 1.0)),
			G: Scalar(math.Min(ToFloat64(result.G), 1.0)),
			B: Scalar(math.Min(ToFloat64(result.B), 1.0)),
			A: Scalar(math.Min(ToFloat64(result.A), 1.0)),
		}.Unpremultiply()
	case BlendModeModulate:
		// r = s*d
		return src.Premultiply().Mul(dst.Premultiply()).Unpremultiply()
	case BlendModeMultiply:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := dstRGB.Mul(srcRGB)
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeScreen:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := srcRGB.Add(dstRGB).Sub(srcRGB.Mul(dstRGB))
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeOverlay:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		screenSrc := dstRGB.Scale(2).Sub(Vector3[Scalar]{X: 1, Y: 1, Z: 1})
		screen := screenSrc.Add(srcRGB).Sub(screenSrc.Mul(srcRGB))
		multiply := srcRGB.Mul(dstRGB.Scale(2))
		blendResult := componentChoose(multiply, screen, dstRGB, 0.5)
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeDarken:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := Vector3[Scalar]{
			X: Scalar(math.Min(ToFloat64(dstRGB.X), ToFloat64(srcRGB.X))),
			Y: Scalar(math.Min(ToFloat64(dstRGB.Y), ToFloat64(srcRGB.Y))),
			Z: Scalar(math.Min(ToFloat64(dstRGB.Z), ToFloat64(srcRGB.Z))),
		}
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeLighten:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := Vector3[Scalar]{
			X: Scalar(math.Max(ToFloat64(dstRGB.X), ToFloat64(srcRGB.X))),
			Y: Scalar(math.Max(ToFloat64(dstRGB.Y), ToFloat64(srcRGB.Y))),
			Z: Scalar(math.Max(ToFloat64(dstRGB.Z), ToFloat64(srcRGB.Z))),
		}
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeColorDodge:
		blendFunc := func(d, s Scalar) Scalar {
			if d < kEhCloseEnough {
				return 0.0
			}
			if 1.0-s < kEhCloseEnough {
				return 1.0
			}
			return Scalar(math.Min(1.0, ToFloat64(d)/(1.0-ToFloat64(s))))
		}
		blendResult := Vector3[Scalar]{
			X: blendFunc(dst.R, src.R),
			Y: blendFunc(dst.G, src.G),
			Z: blendFunc(dst.B, src.B),
		}
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeColorBurn:
		blendFunc := func(d, s Scalar) Scalar {
			if 1.0-d < kEhCloseEnough {
				return 1.0
			}
			if s < kEhCloseEnough {
				return 0.0
			}
			return 1.0 - Scalar(math.Min(1.0, (1.0-ToFloat64(d))/ToFloat64(s)))
		}
		blendResult := Vector3[Scalar]{
			X: blendFunc(dst.R, src.R),
			Y: blendFunc(dst.G, src.G),
			Z: blendFunc(dst.B, src.B),
		}
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeHardLight:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		screenSrc := srcRGB.Scale(2).Sub(Vector3[Scalar]{X: 1, Y: 1, Z: 1})
		screen := screenSrc.Add(dstRGB).Sub(screenSrc.Mul(dstRGB))
		multiply := dstRGB.Mul(srcRGB.Scale(2))
		blendResult := componentChoose(multiply, screen, srcRGB, 0.5)
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeSoftLight:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		D := componentChoose(
			Vector3[Scalar]{
				X: ((dstRGB.X*16-12)*dstRGB.X + 4) * dstRGB.X,
				Y: ((dstRGB.Y*16-12)*dstRGB.Y + 4) * dstRGB.Y,
				Z: ((dstRGB.Z*16-12)*dstRGB.Z + 4) * dstRGB.Z,
			},
			Vector3[Scalar]{
				X: Scalar(math.Sqrt(ToFloat64(dstRGB.X))),
				Y: Scalar(math.Sqrt(ToFloat64(dstRGB.Y))),
				Z: Scalar(math.Sqrt(ToFloat64(dstRGB.Z))),
			},
			dstRGB,
			0.25,
		)
		case1 := dstRGB.Sub(Vector3[Scalar]{X: 1, Y: 1, Z: 1}.Sub(srcRGB.Scale(2)).Mul(dstRGB).Mul(Vector3[Scalar]{X: 1, Y: 1, Z: 1}.Sub(dstRGB)))
		case2 := dstRGB.Add(srcRGB.Scale(2).Sub(Vector3[Scalar]{X: 1, Y: 1, Z: 1}).Mul(D.Sub(dstRGB)))
		blendResult := componentChoose(case1, case2, srcRGB, 0.5)
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeDifference:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := dstRGB.Sub(srcRGB).Abs()
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeExclusion:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := dstRGB.Add(srcRGB).Sub(dstRGB.Mul(srcRGB).Scale(2))
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeHue:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := setLuminosity(setSaturation(srcRGB, saturation(dstRGB)), luminosity(dstRGB))
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeSaturation:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := setLuminosity(setSaturation(dstRGB, saturation(srcRGB)), luminosity(dstRGB))
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeColor:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := setLuminosity(srcRGB, luminosity(dstRGB))
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeLuminosity:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := setLuminosity(dstRGB, luminosity(srcRGB))
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	default:
		// For unsupported modes, default to source over
		return c.Blend(src, BlendModeSrcOver)
	}
}

// String returns a string representation of the color.
func (c Color) String() string {
	return fmt.Sprintf("R=%.2f,G=%.2f,B=%.2f,A=%.2f", c.R, c.G, c.B, c.A)
}

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

// Helper functions for HSV blend modes
func luminosity(color Vector3[Scalar]) Scalar {
	return color.X*0.3 + color.Y*0.59 + color.Z*0.11
}

func saturation(color Vector3[Scalar]) Scalar {
	return max(color.X, color.Y, color.Z) -
		min(color.X, color.Y, color.Z)
}

func clipColor(color Vector3[Scalar]) Vector3[Scalar] {
	lum := luminosity(color)
	mn := min(color.X, color.Y, color.Z)
	mx := max(color.X, color.Y, color.Z)
	const kEhCloseEnough = 1e-6

	if mn < 0 {
		diff := lum - mn + kEhCloseEnough
		if diff != 0 {
			factor := lum / diff
			color = Vector3[Scalar]{
				X: lum + (color.X-lum)*factor,
				Y: lum + (color.Y-lum)*factor,
				Z: lum + (color.Z-lum)*factor,
			}
		}
	}

	if mx > 1 {
		diff := mx - lum + kEhCloseEnough
		if diff != 0 {
			factor := (1 - lum) / diff
			color = Vector3[Scalar]{
				X: lum + (color.X-lum)*factor,
				Y: lum + (color.Y-lum)*factor,
				Z: lum + (color.Z-lum)*factor,
			}
		}
	}

	return color
}

func setLuminosity(color Vector3[Scalar], lum Scalar) Vector3[Scalar] {
	relativeLum := lum - luminosity(color)
	return clipColor(Vector3[Scalar]{
		X: color.X + relativeLum,
		Y: color.Y + relativeLum,
		Z: color.Z + relativeLum,
	})
}

func setSaturation(color Vector3[Scalar], sat Scalar) Vector3[Scalar] {
	mn := min(color.X, color.Y, color.Z)
	mx := max(color.X, color.Y, color.Z)
	if mn < mx {
		factor := sat / (mx - mn)
		return Vector3[Scalar]{
			X: (color.X - mn) * factor,
			Y: (color.Y - mn) * factor,
			Z: (color.Z - mn) * factor,
		}
	}
	return Vector3[Scalar]{}
}

func componentChoose(a, b, value Vector3[Scalar], cutoff Scalar) Vector3[Scalar] {
	result := Vector3[Scalar]{}
	if value.X > cutoff {
		result.X = b.X
	} else {
		result.X = a.X
	}
	if value.Y > cutoff {
		result.Y = b.Y
	} else {
		result.Y = a.Y
	}
	if value.Z > cutoff {
		result.Z = b.Z
	} else {
		result.Z = a.Z
	}
	return result
}

func toRGB(c Color) Vector3[Scalar] {
	return Vector3[Scalar]{X: c.R, Y: c.G, Z: c.B}
}

func fromRGB(rgb Vector3[Scalar], alpha Scalar) Color {
	return Color{R: rgb.X, G: rgb.Y, B: rgb.Z, A: alpha}
}

func applyBlendedColor(dst, src Color, blendResult Vector3[Scalar]) Color {
	dst = dst.Premultiply()

	// Use the blended color for areas where the source and destination colors overlap
	blended := fromRGB(blendResult, src.A*dst.A).Premultiply()
	// Use the original source color for any remaining non-overlapping areas
	srcPremul := src.Premultiply().Scale(1 - dst.A)
	src = blended.Add(srcPremul)

	// Source-over composite the blended source color atop the destination
	return src.Add(dst.Scale(1 - src.A))
}

// Predefined colors
func ColorWhite() Color            { return Color{R: 1, G: 1, B: 1, A: 1} }
func ColorBlack() Color            { return Color{R: 0, G: 0, B: 0, A: 1} }
func ColorTransparent() Color      { return Color{R: 0, G: 0, B: 0, A: 0} }
func ColorWhiteTransparent() Color { return Color{R: 1, G: 1, B: 1, A: 0} }
func ColorBlackTransparent() Color { return Color{R: 0, G: 0, B: 0, A: 0} }
func ColorRed() Color              { return Color{R: 1, G: 0, B: 0, A: 1} }
func ColorGreen() Color            { return Color{R: 0, G: 1, B: 0, A: 1} }
func ColorBlue() Color             { return Color{R: 0, G: 0, B: 1, A: 1} }
func ColorYellow() Color           { return Color{R: 1, G: 1, B: 0, A: 1} }
func ColorCyan() Color             { return Color{R: 0, G: 1, B: 1, A: 1} }
func ColorMagenta() Color          { return Color{R: 1, G: 0, B: 1, A: 1} }

// Additional predefined colors
func ColorAliceBlue() Color            { return NewColorRGBA8(240, 248, 255, 255) }
func ColorAntiqueWhite() Color         { return NewColorRGBA8(250, 235, 215, 255) }
func ColorAqua() Color                 { return NewColorRGBA8(0, 255, 255, 255) }
func ColorAquaMarine() Color           { return NewColorRGBA8(127, 255, 212, 255) }
func ColorAzure() Color                { return NewColorRGBA8(240, 255, 255, 255) }
func ColorBeige() Color                { return NewColorRGBA8(245, 245, 220, 255) }
func ColorBisque() Color               { return NewColorRGBA8(255, 228, 196, 255) }
func ColorBlanchedAlmond() Color       { return NewColorRGBA8(255, 235, 205, 255) }
func ColorBlueViolet() Color           { return NewColorRGBA8(138, 43, 226, 255) }
func ColorBrown() Color                { return NewColorRGBA8(165, 42, 42, 255) }
func ColorBurlyWood() Color            { return NewColorRGBA8(222, 184, 135, 255) }
func ColorCadetBlue() Color            { return NewColorRGBA8(95, 158, 160, 255) }
func ColorChartreuse() Color           { return NewColorRGBA8(127, 255, 0, 255) }
func ColorChocolate() Color            { return NewColorRGBA8(210, 105, 30, 255) }
func ColorCoral() Color                { return NewColorRGBA8(255, 127, 80, 255) }
func ColorCornflowerBlue() Color       { return NewColorRGBA8(100, 149, 237, 255) }
func ColorCornsilk() Color             { return NewColorRGBA8(255, 248, 220, 255) }
func ColorCrimson() Color              { return NewColorRGBA8(220, 20, 60, 255) }
func ColorDarkBlue() Color             { return NewColorRGBA8(0, 0, 139, 255) }
func ColorDarkCyan() Color             { return NewColorRGBA8(0, 139, 139, 255) }
func ColorDarkGoldenrod() Color        { return NewColorRGBA8(184, 134, 11, 255) }
func ColorDarkGray() Color             { return NewColorRGBA8(169, 169, 169, 255) }
func ColorDarkGreen() Color            { return NewColorRGBA8(0, 100, 0, 255) }
func ColorDarkGrey() Color             { return NewColorRGBA8(169, 169, 169, 255) }
func ColorDarkKhaki() Color            { return NewColorRGBA8(189, 183, 107, 255) }
func ColorDarkMagenta() Color          { return NewColorRGBA8(139, 0, 139, 255) }
func ColorDarkOliveGreen() Color       { return NewColorRGBA8(85, 107, 47, 255) }
func ColorDarkOrange() Color           { return NewColorRGBA8(255, 140, 0, 255) }
func ColorDarkOrchid() Color           { return NewColorRGBA8(153, 50, 204, 255) }
func ColorDarkRed() Color              { return NewColorRGBA8(139, 0, 0, 255) }
func ColorDarkSalmon() Color           { return NewColorRGBA8(233, 150, 122, 255) }
func ColorDarkSeagreen() Color         { return NewColorRGBA8(143, 188, 143, 255) }
func ColorDarkSlateBlue() Color        { return NewColorRGBA8(72, 61, 139, 255) }
func ColorDarkSlateGray() Color        { return NewColorRGBA8(47, 79, 79, 255) }
func ColorDarkSlateGrey() Color        { return NewColorRGBA8(47, 79, 79, 255) }
func ColorDarkTurquoise() Color        { return NewColorRGBA8(0, 206, 209, 255) }
func ColorDarkViolet() Color           { return NewColorRGBA8(148, 0, 211, 255) }
func ColorDeepPink() Color             { return NewColorRGBA8(255, 20, 147, 255) }
func ColorDeepSkyBlue() Color          { return NewColorRGBA8(0, 191, 255, 255) }
func ColorDimGray() Color              { return NewColorRGBA8(105, 105, 105, 255) }
func ColorDimGrey() Color              { return NewColorRGBA8(105, 105, 105, 255) }
func ColorDodgerBlue() Color           { return NewColorRGBA8(30, 144, 255, 255) }
func ColorFirebrick() Color            { return NewColorRGBA8(178, 34, 34, 255) }
func ColorFloralWhite() Color          { return NewColorRGBA8(255, 250, 240, 255) }
func ColorForestGreen() Color          { return NewColorRGBA8(34, 139, 34, 255) }
func ColorFuchsia() Color              { return NewColorRGBA8(255, 0, 255, 255) }
func ColorGainsboro() Color            { return NewColorRGBA8(220, 220, 220, 255) }
func ColorGhostwhite() Color           { return NewColorRGBA8(248, 248, 255, 255) }
func ColorGold() Color                 { return NewColorRGBA8(255, 215, 0, 255) }
func ColorGoldenrod() Color            { return NewColorRGBA8(218, 165, 32, 255) }
func ColorGreenYellow() Color          { return NewColorRGBA8(173, 255, 47, 255) }
func ColorHoneydew() Color             { return NewColorRGBA8(240, 255, 240, 255) }
func ColorHotPink() Color              { return NewColorRGBA8(255, 105, 180, 255) }
func ColorIndianRed() Color            { return NewColorRGBA8(205, 92, 92, 255) }
func ColorIndigo() Color               { return NewColorRGBA8(75, 0, 130, 255) }
func ColorIvory() Color                { return NewColorRGBA8(255, 255, 240, 255) }
func ColorKhaki() Color                { return NewColorRGBA8(240, 230, 140, 255) }
func ColorLavender() Color             { return NewColorRGBA8(230, 230, 250, 255) }
func ColorLavenderBlush() Color        { return NewColorRGBA8(255, 240, 245, 255) }
func ColorLawnGreen() Color            { return NewColorRGBA8(124, 252, 0, 255) }
func ColorLemonChiffon() Color         { return NewColorRGBA8(255, 250, 205, 255) }
func ColorLightBlue() Color            { return NewColorRGBA8(173, 216, 230, 255) }
func ColorLightCoral() Color           { return NewColorRGBA8(240, 128, 128, 255) }
func ColorLightCyan() Color            { return NewColorRGBA8(224, 255, 255, 255) }
func ColorLightGoldenrodYellow() Color { return NewColorRGBA8(250, 250, 210, 255) }
func ColorLightGray() Color            { return NewColorRGBA8(211, 211, 211, 255) }
func ColorLightGreen() Color           { return NewColorRGBA8(144, 238, 144, 255) }
func ColorLightGrey() Color            { return NewColorRGBA8(211, 211, 211, 255) }
func ColorLightPink() Color            { return NewColorRGBA8(255, 182, 193, 255) }
func ColorLightSalmon() Color          { return NewColorRGBA8(255, 160, 122, 255) }
func ColorLightSeaGreen() Color        { return NewColorRGBA8(32, 178, 170, 255) }
func ColorLightSkyBlue() Color         { return NewColorRGBA8(135, 206, 250, 255) }
func ColorLightSlateGray() Color       { return NewColorRGBA8(119, 136, 153, 255) }
func ColorLightSlateGrey() Color       { return NewColorRGBA8(119, 136, 153, 255) }
func ColorLightSteelBlue() Color       { return NewColorRGBA8(176, 196, 222, 255) }
func ColorLightYellow() Color          { return NewColorRGBA8(255, 255, 224, 255) }
func ColorLime() Color                 { return NewColorRGBA8(0, 255, 0, 255) }
func ColorLimeGreen() Color            { return NewColorRGBA8(50, 205, 50, 255) }
func ColorLinen() Color                { return NewColorRGBA8(250, 240, 230, 255) }
func ColorMaroon() Color               { return NewColorRGBA8(128, 0, 0, 255) }
func ColorMediumAquamarine() Color     { return NewColorRGBA8(102, 205, 170, 255) }
func ColorMediumBlue() Color           { return NewColorRGBA8(0, 0, 205, 255) }
func ColorMediumOrchid() Color         { return NewColorRGBA8(186, 85, 211, 255) }
func ColorMediumPurple() Color         { return NewColorRGBA8(147, 112, 219, 255) }
func ColorMediumSeagreen() Color       { return NewColorRGBA8(60, 179, 113, 255) }
func ColorMediumSlateBlue() Color      { return NewColorRGBA8(123, 104, 238, 255) }
func ColorMediumSpringGreen() Color    { return NewColorRGBA8(0, 250, 154, 255) }
func ColorMediumTurquoise() Color      { return NewColorRGBA8(72, 209, 204, 255) }
func ColorMediumVioletRed() Color      { return NewColorRGBA8(199, 21, 133, 255) }
func ColorMidnightBlue() Color         { return NewColorRGBA8(25, 25, 112, 255) }
func ColorMintCream() Color            { return NewColorRGBA8(245, 255, 250, 255) }
func ColorMistyRose() Color            { return NewColorRGBA8(255, 228, 225, 255) }
func ColorMoccasin() Color             { return NewColorRGBA8(255, 228, 181, 255) }
func ColorNavajoWhite() Color          { return NewColorRGBA8(255, 222, 173, 255) }
func ColorNavy() Color                 { return NewColorRGBA8(0, 0, 128, 255) }
func ColorOldLace() Color              { return NewColorRGBA8(253, 245, 230, 255) }
func ColorOlive() Color                { return NewColorRGBA8(128, 128, 0, 255) }
func ColorOliveDrab() Color            { return NewColorRGBA8(107, 142, 35, 255) }
func ColorOrangeRed() Color            { return NewColorRGBA8(255, 69, 0, 255) }
func ColorOrchid() Color               { return NewColorRGBA8(218, 112, 214, 255) }
func ColorPaleGoldenrod() Color        { return NewColorRGBA8(238, 232, 170, 255) }
func ColorPaleGreen() Color            { return NewColorRGBA8(152, 251, 152, 255) }
func ColorPaleTurquoise() Color        { return NewColorRGBA8(175, 238, 238, 255) }
func ColorPaleVioletRed() Color        { return NewColorRGBA8(219, 112, 147, 255) }
func ColorPapayaWhip() Color           { return NewColorRGBA8(255, 239, 213, 255) }
func ColorPeachpuff() Color            { return NewColorRGBA8(255, 218, 185, 255) }
func ColorPeru() Color                 { return NewColorRGBA8(205, 133, 63, 255) }
func ColorPlum() Color                 { return NewColorRGBA8(221, 160, 221, 255) }
func ColorPowderBlue() Color           { return NewColorRGBA8(176, 224, 230, 255) }
func ColorRosyBrown() Color            { return NewColorRGBA8(188, 143, 143, 255) }
func ColorRoyalBlue() Color            { return NewColorRGBA8(65, 105, 225, 255) }
func ColorSaddleBrown() Color          { return NewColorRGBA8(139, 69, 19, 255) }
func ColorSalmon() Color               { return NewColorRGBA8(250, 128, 114, 255) }
func ColorSandyBrown() Color           { return NewColorRGBA8(244, 164, 96, 255) }
func ColorSeagreen() Color             { return NewColorRGBA8(46, 139, 87, 255) }
func ColorSeashell() Color             { return NewColorRGBA8(255, 245, 238, 255) }
func ColorSienna() Color               { return NewColorRGBA8(160, 82, 45, 255) }
func ColorSilver() Color               { return NewColorRGBA8(192, 192, 192, 255) }
func ColorSkyBlue() Color              { return NewColorRGBA8(135, 206, 235, 255) }
func ColorSlateBlue() Color            { return NewColorRGBA8(106, 90, 205, 255) }
func ColorSlateGray() Color            { return NewColorRGBA8(112, 128, 144, 255) }
func ColorSlateGrey() Color            { return NewColorRGBA8(112, 128, 144, 255) }
func ColorSnow() Color                 { return NewColorRGBA8(255, 250, 250, 255) }
func ColorSpringGreen() Color          { return NewColorRGBA8(0, 255, 127, 255) }
func ColorSteelBlue() Color            { return NewColorRGBA8(70, 130, 180, 255) }
func ColorTan() Color                  { return NewColorRGBA8(210, 180, 140, 255) }
func ColorTeal() Color                 { return NewColorRGBA8(0, 128, 128, 255) }
func ColorThistle() Color              { return NewColorRGBA8(216, 191, 216, 255) }
func ColorTomato() Color               { return NewColorRGBA8(255, 99, 71, 255) }
func ColorTurquoise() Color            { return NewColorRGBA8(64, 224, 208, 255) }
func ColorViolet() Color               { return NewColorRGBA8(238, 130, 238, 255) }
func ColorWheat() Color                { return NewColorRGBA8(245, 222, 179, 255) }
func ColorWhitesmoke() Color           { return NewColorRGBA8(245, 245, 245, 255) }
func ColorYellowGreen() Color          { return NewColorRGBA8(154, 205, 50, 255) }
