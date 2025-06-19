package geom

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"
)

// ColorMatrix represents a 4x5 matrix for color transformation.
//
//	[ a, b, c, d, e ]
//	[ f, g, h, i, j ]
//	[ k, l, m, n, o ]
//	[ p, q, r, s, t ]
//
// When applied to a color [R, G, B, A], the resulting color is computed as:
//
//	R’ = a*R + b*G + c*B + d*A + e;
//	G’ = f*R + g*G + h*B + i*A + j;
//	B’ = k*R + l*G + m*B + n*A + o;
//	A’ = p*R + q*G + r*B + s*A + t;
//
// That resulting color [R’, G’, B’, A’] then has each channel clamped to the 0
// to 1 range.
type ColorMatrix [20]F32

func NewColorMatrix() ColorMatrix {
	return ColorMatrix{
		1, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 1, 0,
	}
}

// Color represents an RGBA color with components in the range [0, 1].
// It implements Go's standard color.Color interface.
type Color struct {
	R, G, B, A F32
}

// NewColor creates a new Color with the given RGBA components.
func NewColor(r, g, b, a F32) Color {
	return Color{R: r, G: g, B: b, A: a}
}

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

// NewColorFromGo creates a Color from any color.Color.
func NewColorFromGo(c color.Color) Color {
	r, g, b, a := c.RGBA()
	return Color{
		R: F32(r) / 65535.0,
		G: F32(g) / 65535.0,
		B: F32(b) / 65535.0,
		A: F32(a) / 65535.0,
	}
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

// NewColorFromRGBA creates a Color from Go's standard color.RGBA.
func NewColorFromRGBA(c color.RGBA) Color {
	return Color{
		R: F32(c.R) / 255.0,
		G: F32(c.G) / 255.0,
		B: F32(c.B) / 255.0,
		A: F32(c.A) / 255.0,
	}
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

// ToR8G8B8A8 converts color to 8-bit RGBA array.
func (c Color) ToR8G8B8A8() [4]uint8 {
	return [4]uint8{
		uint8(Clamp(c.R*255+0.5, 0, 255)),
		uint8(Clamp(c.G*255+0.5, 0, 255)),
		uint8(Clamp(c.B*255+0.5, 0, 255)),
		uint8(Clamp(c.A*255+0.5, 0, 255)),
	}
}

// ToARGB converts color to ARGB 32-bit value.
func (c Color) ToARGB() uint32 {
	rgba := c.ToR8G8B8A8()
	return uint32(rgba[3])<<24 | uint32(rgba[0])<<16 | uint32(rgba[1])<<8 | uint32(rgba[2])
}

// ToIColor converts color to 32-bit ARGB representation.
func (c Color) ToIColor() uint32 {
	return c.ToARGB()
}

// Eq compares two colors for Eqity within tolerance.
func (c Color) Eq(other Color) bool {
	return ScalarEq(c.R, other.R) &&
		ScalarEq(c.G, other.G) &&
		ScalarEq(c.B, other.B) &&
		ScalarEq(c.A, other.A)
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

// Div performs component-wise division.
func (c Color) Div(other Color) Color {
	return Color{
		R: c.R / other.R,
		G: c.G / other.G,
		B: c.B / other.B,
		A: c.A / other.A,
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

// LinearToSRGB converts color from linear space to sRGB space.
func (c Color) LinearToSRGB() Color {
	convert := func(component F32) F32 {
		if component <= 0.0031308 {
			return component * 12.92
		}
		return F32(1.055*math.Pow(component.Float64(), 1.0/2.4) - 0.055)
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
		return F32(math.Pow((component.Float64()+0.055)/1.055, 2.4))
	}

	return Color{
		R: convert(c.R),
		G: convert(c.G),
		B: convert(c.B),
		A: c.A,
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
			R: F32(math.Min(result.R.Float64(), 1.0)),
			G: F32(math.Min(result.G.Float64(), 1.0)),
			B: F32(math.Min(result.B.Float64(), 1.0)),
			A: F32(math.Min(result.A.Float64(), 1.0)),
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
		screenSrc := dstRGB.Scale(2).Sub(Vector3[F32]{X: 1, Y: 1, Z: 1})
		screen := screenSrc.Add(srcRGB).Sub(screenSrc.Mul(srcRGB))
		multiply := srcRGB.Mul(dstRGB.Scale(2))
		blendResult := componentChoose(multiply, screen, dstRGB, 0.5)
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeDarken:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := Vector3[F32]{
			X: F32(math.Min(dstRGB.X.Float64(), srcRGB.X.Float64())),
			Y: F32(math.Min(dstRGB.Y.Float64(), srcRGB.Y.Float64())),
			Z: F32(math.Min(dstRGB.Z.Float64(), srcRGB.Z.Float64())),
		}
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeLighten:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		blendResult := Vector3[F32]{
			X: F32(math.Max(dstRGB.X.Float64(), srcRGB.X.Float64())),
			Y: F32(math.Max(dstRGB.Y.Float64(), srcRGB.Y.Float64())),
			Z: F32(math.Max(dstRGB.Z.Float64(), srcRGB.Z.Float64())),
		}
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeColorDodge:
		blendFunc := func(d, s F32) F32 {
			if d < kEhCloseEnough {
				return 0.0
			}
			if 1.0-s < kEhCloseEnough {
				return 1.0
			}
			return F32(math.Min(1.0, d.Float64()/(1.0-s.Float64())))
		}
		blendResult := Vector3[F32]{
			X: blendFunc(dst.R, src.R),
			Y: blendFunc(dst.G, src.G),
			Z: blendFunc(dst.B, src.B),
		}
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeColorBurn:
		blendFunc := func(d, s F32) F32 {
			if 1.0-d < kEhCloseEnough {
				return 1.0
			}
			if s < kEhCloseEnough {
				return 0.0
			}
			return 1.0 - F32(math.Min(1.0, (1.0-d.Float64())/s.Float64()))
		}
		blendResult := Vector3[F32]{
			X: blendFunc(dst.R, src.R),
			Y: blendFunc(dst.G, src.G),
			Z: blendFunc(dst.B, src.B),
		}
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeHardLight:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		screenSrc := srcRGB.Scale(2).Sub(Vector3[F32]{X: 1, Y: 1, Z: 1})
		screen := screenSrc.Add(dstRGB).Sub(screenSrc.Mul(dstRGB))
		multiply := dstRGB.Mul(srcRGB.Scale(2))
		blendResult := componentChoose(multiply, screen, srcRGB, 0.5)
		return applyBlendedColor(dst, src, blendResult).Unpremultiply()
	case BlendModeSoftLight:
		dstRGB := toRGB(dst)
		srcRGB := toRGB(src)
		D := componentChoose(
			Vector3[F32]{
				X: ((dstRGB.X*16-12)*dstRGB.X + 4) * dstRGB.X,
				Y: ((dstRGB.Y*16-12)*dstRGB.Y + 4) * dstRGB.Y,
				Z: ((dstRGB.Z*16-12)*dstRGB.Z + 4) * dstRGB.Z,
			},
			Vector3[F32]{
				X: F32(math.Sqrt(dstRGB.X.Float64())),
				Y: F32(math.Sqrt(dstRGB.Y.Float64())),
				Z: F32(math.Sqrt(dstRGB.Z.Float64())),
			},
			dstRGB,
			0.25,
		)
		case1 := dstRGB.Sub(Vector3[F32]{X: 1, Y: 1, Z: 1}.Sub(srcRGB.Scale(2)).Mul(dstRGB).Mul(Vector3[F32]{X: 1, Y: 1, Z: 1}.Sub(dstRGB)))
		case2 := dstRGB.Add(srcRGB.Scale(2).Sub(Vector3[F32]{X: 1, Y: 1, Z: 1}).Mul(D.Sub(dstRGB)))
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
func luminosity(color Vector3[F32]) F32 {
	return color.X*0.3 + color.Y*0.59 + color.Z*0.11
}

func saturation(color Vector3[F32]) F32 {
	return max(color.X, color.Y, color.Z) -
		min(color.X, color.Y, color.Z)
}

func clipColor(color Vector3[F32]) Vector3[F32] {
	lum := luminosity(color)
	mn := min(color.X, color.Y, color.Z)
	mx := max(color.X, color.Y, color.Z)
	const kEhCloseEnough = 1e-6

	if mn < 0 {
		diff := lum - mn + kEhCloseEnough
		if diff != 0 {
			factor := lum / diff
			color = Vector3[F32]{
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
			color = Vector3[F32]{
				X: lum + (color.X-lum)*factor,
				Y: lum + (color.Y-lum)*factor,
				Z: lum + (color.Z-lum)*factor,
			}
		}
	}

	return color
}

func setLuminosity(color Vector3[F32], lum F32) Vector3[F32] {
	relativeLum := lum - luminosity(color)
	return clipColor(Vector3[F32]{
		X: color.X + relativeLum,
		Y: color.Y + relativeLum,
		Z: color.Z + relativeLum,
	})
}

func setSaturation(color Vector3[F32], sat F32) Vector3[F32] {
	mn := min(color.X, color.Y, color.Z)
	mx := max(color.X, color.Y, color.Z)
	if mn < mx {
		factor := sat / (mx - mn)
		return Vector3[F32]{
			X: (color.X - mn) * factor,
			Y: (color.Y - mn) * factor,
			Z: (color.Z - mn) * factor,
		}
	}
	return Vector3[F32]{}
}

func componentChoose(a, b, value Vector3[F32], cutoff F32) Vector3[F32] {
	result := Vector3[F32]{}
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

func toRGB(c Color) Vector3[F32] {
	return Vector3[F32]{X: c.R, Y: c.G, Z: c.B}
}

func fromRGB(rgb Vector3[F32], alpha F32) Color {
	return Color{R: rgb.X, G: rgb.Y, B: rgb.Z, A: alpha}
}

func applyBlendedColor(dst, src Color, blendResult Vector3[F32]) Color {
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

// Additional predefined colors matching C++ implementation
func ColorAliceBlue() Color            { return NewColorRGB8(240, 248, 255) }
func ColorAntiqueWhite() Color         { return NewColorRGB8(250, 235, 215) }
func ColorAqua() Color                 { return NewColorRGB8(0, 255, 255) }
func ColorAquaMarine() Color           { return NewColorRGB8(127, 255, 212) }
func ColorAzure() Color                { return NewColorRGB8(240, 255, 255) }
func ColorBeige() Color                { return NewColorRGB8(245, 245, 220) }
func ColorBisque() Color               { return NewColorRGB8(255, 228, 196) }
func ColorBlanchedAlmond() Color       { return NewColorRGB8(255, 235, 205) }
func ColorBlueViolet() Color           { return NewColorRGB8(138, 43, 226) }
func ColorBrown() Color                { return NewColorRGB8(165, 42, 42) }
func ColorBurlyWood() Color            { return NewColorRGB8(222, 184, 135) }
func ColorCadetBlue() Color            { return NewColorRGB8(95, 158, 160) }
func ColorChartreuse() Color           { return NewColorRGB8(127, 255, 0) }
func ColorChocolate() Color            { return NewColorRGB8(210, 105, 30) }
func ColorCoral() Color                { return NewColorRGB8(255, 127, 80) }
func ColorCornflowerBlue() Color       { return NewColorRGB8(100, 149, 237) }
func ColorCornsilk() Color             { return NewColorRGB8(255, 248, 220) }
func ColorCrimson() Color              { return NewColorRGB8(220, 20, 60) }
func ColorDarkBlue() Color             { return NewColorRGB8(0, 0, 139) }
func ColorDarkCyan() Color             { return NewColorRGB8(0, 139, 139) }
func ColorDarkGoldenrod() Color        { return NewColorRGB8(184, 134, 11) }
func ColorDarkGray() Color             { return NewColorRGB8(169, 169, 169) }
func ColorDarkGreen() Color            { return NewColorRGB8(0, 100, 0) }
func ColorDarkGrey() Color             { return NewColorRGB8(169, 169, 169) }
func ColorDarkKhaki() Color            { return NewColorRGB8(189, 183, 107) }
func ColorDarkMagenta() Color          { return NewColorRGB8(139, 0, 139) }
func ColorDarkOliveGreen() Color       { return NewColorRGB8(85, 107, 47) }
func ColorDarkOrange() Color           { return NewColorRGB8(255, 140, 0) }
func ColorDarkOrchid() Color           { return NewColorRGB8(153, 50, 204) }
func ColorDarkRed() Color              { return NewColorRGB8(139, 0, 0) }
func ColorDarkSalmon() Color           { return NewColorRGB8(233, 150, 122) }
func ColorDarkSeagreen() Color         { return NewColorRGB8(143, 188, 143) }
func ColorDarkSlateBlue() Color        { return NewColorRGB8(72, 61, 139) }
func ColorDarkSlateGray() Color        { return NewColorRGB8(47, 79, 79) }
func ColorDarkSlateGrey() Color        { return NewColorRGB8(47, 79, 79) }
func ColorDarkTurquoise() Color        { return NewColorRGB8(0, 206, 209) }
func ColorDarkViolet() Color           { return NewColorRGB8(148, 0, 211) }
func ColorDeepPink() Color             { return NewColorRGB8(255, 20, 147) }
func ColorDeepSkyBlue() Color          { return NewColorRGB8(0, 191, 255) }
func ColorDimGray() Color              { return NewColorRGB8(105, 105, 105) }
func ColorDimGrey() Color              { return NewColorRGB8(105, 105, 105) }
func ColorDodgerBlue() Color           { return NewColorRGB8(30, 144, 255) }
func ColorFirebrick() Color            { return NewColorRGB8(178, 34, 34) }
func ColorFloralWhite() Color          { return NewColorRGB8(255, 250, 240) }
func ColorForestGreen() Color          { return NewColorRGB8(34, 139, 34) }
func ColorFuchsia() Color              { return NewColorRGB8(255, 0, 255) }
func ColorGainsboro() Color            { return NewColorRGB8(220, 220, 220) }
func ColorGhostwhite() Color           { return NewColorRGB8(248, 248, 255) }
func ColorGold() Color                 { return NewColorRGB8(255, 215, 0) }
func ColorGoldenrod() Color            { return NewColorRGB8(218, 165, 32) }
func ColorGreenYellow() Color          { return NewColorRGB8(173, 255, 47) }
func ColorHoneydew() Color             { return NewColorRGB8(240, 255, 240) }
func ColorHotPink() Color              { return NewColorRGB8(255, 105, 180) }
func ColorIndianRed() Color            { return NewColorRGB8(205, 92, 92) }
func ColorIndigo() Color               { return NewColorRGB8(75, 0, 130) }
func ColorIvory() Color                { return NewColorRGB8(255, 255, 240) }
func ColorKhaki() Color                { return NewColorRGB8(240, 230, 140) }
func ColorLavender() Color             { return NewColorRGB8(230, 230, 250) }
func ColorLavenderBlush() Color        { return NewColorRGB8(255, 240, 245) }
func ColorLawnGreen() Color            { return NewColorRGB8(124, 252, 0) }
func ColorLemonChiffon() Color         { return NewColorRGB8(255, 250, 205) }
func ColorLightBlue() Color            { return NewColorRGB8(173, 216, 230) }
func ColorLightCoral() Color           { return NewColorRGB8(240, 128, 128) }
func ColorLightCyan() Color            { return NewColorRGB8(224, 255, 255) }
func ColorLightGoldenrodYellow() Color { return NewColorRGB8(250, 250, 210) }
func ColorLightGray() Color            { return NewColorRGB8(211, 211, 211) }
func ColorLightGreen() Color           { return NewColorRGB8(144, 238, 144) }
func ColorLightGrey() Color            { return NewColorRGB8(211, 211, 211) }
func ColorLightPink() Color            { return NewColorRGB8(255, 182, 193) }
func ColorLightSalmon() Color          { return NewColorRGB8(255, 160, 122) }
func ColorLightSeaGreen() Color        { return NewColorRGB8(32, 178, 170) }
func ColorLightSkyBlue() Color         { return NewColorRGB8(135, 206, 250) }
func ColorLightSlateGray() Color       { return NewColorRGB8(119, 136, 153) }
func ColorLightSlateGrey() Color       { return NewColorRGB8(119, 136, 153) }
func ColorLightSteelBlue() Color       { return NewColorRGB8(176, 196, 222) }
func ColorLightYellow() Color          { return NewColorRGB8(255, 255, 224) }
func ColorLime() Color                 { return NewColorRGB8(0, 255, 0) }
func ColorLimeGreen() Color            { return NewColorRGB8(50, 205, 50) }
func ColorLinen() Color                { return NewColorRGB8(250, 240, 230) }
func ColorMaroon() Color               { return NewColorRGB8(128, 0, 0) }
func ColorMediumAquamarine() Color     { return NewColorRGB8(102, 205, 170) }
func ColorMediumBlue() Color           { return NewColorRGB8(0, 0, 205) }
func ColorMediumOrchid() Color         { return NewColorRGB8(186, 85, 211) }
func ColorMediumPurple() Color         { return NewColorRGB8(147, 112, 219) }
func ColorMediumSeagreen() Color       { return NewColorRGB8(60, 179, 113) }
func ColorMediumSlateBlue() Color      { return NewColorRGB8(123, 104, 238) }
func ColorMediumSpringGreen() Color    { return NewColorRGB8(0, 250, 154) }
func ColorMediumTurquoise() Color      { return NewColorRGB8(72, 209, 204) }
func ColorMediumVioletRed() Color      { return NewColorRGB8(199, 21, 133) }
func ColorMidnightBlue() Color         { return NewColorRGB8(25, 25, 112) }
func ColorMintCream() Color            { return NewColorRGB8(245, 255, 250) }
func ColorMistyRose() Color            { return NewColorRGB8(255, 228, 225) }
func ColorMoccasin() Color             { return NewColorRGB8(255, 228, 181) }
func ColorNavajoWhite() Color          { return NewColorRGB8(255, 222, 173) }
func ColorNavy() Color                 { return NewColorRGB8(0, 0, 128) }
func ColorOldLace() Color              { return NewColorRGB8(253, 245, 230) }
func ColorOlive() Color                { return NewColorRGB8(128, 128, 0) }
func ColorOliveDrab() Color            { return NewColorRGB8(107, 142, 35) }
func ColorOrangeRed() Color            { return NewColorRGB8(255, 69, 0) }
func ColorOrchid() Color               { return NewColorRGB8(218, 112, 214) }
func ColorPaleGoldenrod() Color        { return NewColorRGB8(238, 232, 170) }
func ColorPaleGreen() Color            { return NewColorRGB8(152, 251, 152) }
func ColorPaleTurquoise() Color        { return NewColorRGB8(175, 238, 238) }
func ColorPaleVioletRed() Color        { return NewColorRGB8(219, 112, 147) }
func ColorPapayaWhip() Color           { return NewColorRGB8(255, 239, 213) }
func ColorPeachpuff() Color            { return NewColorRGB8(255, 218, 185) }
func ColorPeru() Color                 { return NewColorRGB8(205, 133, 63) }
func ColorPlum() Color                 { return NewColorRGB8(221, 160, 221) }
func ColorPowderBlue() Color           { return NewColorRGB8(176, 224, 230) }
func ColorRosyBrown() Color            { return NewColorRGB8(188, 143, 143) }
func ColorRoyalBlue() Color            { return NewColorRGB8(65, 105, 225) }
func ColorSaddleBrown() Color          { return NewColorRGB8(139, 69, 19) }
func ColorSalmon() Color               { return NewColorRGB8(250, 128, 114) }
func ColorSandyBrown() Color           { return NewColorRGB8(244, 164, 96) }
func ColorSeagreen() Color             { return NewColorRGB8(46, 139, 87) }
func ColorSeashell() Color             { return NewColorRGB8(255, 245, 238) }
func ColorSienna() Color               { return NewColorRGB8(160, 82, 45) }
func ColorSilver() Color               { return NewColorRGB8(192, 192, 192) }
func ColorSkyBlue() Color              { return NewColorRGB8(135, 206, 235) }
func ColorSlateBlue() Color            { return NewColorRGB8(106, 90, 205) }
func ColorSlateGray() Color            { return NewColorRGB8(112, 128, 144) }
func ColorSlateGrey() Color            { return NewColorRGB8(112, 128, 144) }
func ColorSnow() Color                 { return NewColorRGB8(255, 250, 250) }
func ColorSpringGreen() Color          { return NewColorRGB8(0, 255, 127) }
func ColorSteelBlue() Color            { return NewColorRGB8(70, 130, 180) }
func ColorTan() Color                  { return NewColorRGB8(210, 180, 140) }
func ColorTeal() Color                 { return NewColorRGB8(0, 128, 128) }
func ColorThistle() Color              { return NewColorRGB8(216, 191, 216) }
func ColorTomato() Color               { return NewColorRGB8(255, 99, 71) }
func ColorTurquoise() Color            { return NewColorRGB8(64, 224, 208) }
func ColorViolet() Color               { return NewColorRGB8(238, 130, 238) }
func ColorWheat() Color                { return NewColorRGB8(245, 222, 179) }
func ColorWhitesmoke() Color           { return NewColorRGB8(245, 245, 245) }
func ColorYellowGreen() Color          { return NewColorRGB8(154, 205, 50) }
