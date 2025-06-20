package geom

import (
	"math"
	"testing"
)

func TestColor_RGBA(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		wantR uint32
		wantG uint32
		wantB uint32
		wantA uint32
	}{
		{
			name:  "white opaque",
			color: Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			wantR: 65535, wantG: 65535, wantB: 65535, wantA: 65535,
		},
		{
			name:  "black opaque",
			color: Color{R: 0.0, G: 0.0, B: 0.0, A: 1.0},
			wantR: 0, wantG: 0, wantB: 0, wantA: 65535,
		},
		{
			name:  "red half alpha",
			color: Color{R: 1.0, G: 0.0, B: 0.0, A: 0.5},
			wantR: 32768, wantG: 0, wantB: 0, wantA: 32768,
		},
		{
			name:  "transparent",
			color: Color{R: 0.5, G: 0.5, B: 0.5, A: 0.0},
			wantR: 0, wantG: 0, wantB: 0, wantA: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotG, gotB, gotA := tt.color.RGBA()
			// Allow small rounding differences
			tolerance := uint32(1)
			if diff := absDiffUint32(gotR, tt.wantR); diff > tolerance {
				t.Errorf("Color.RGBA() R = %v, want %v", gotR, tt.wantR)
			}
			if diff := absDiffUint32(gotG, tt.wantG); diff > tolerance {
				t.Errorf("Color.RGBA() G = %v, want %v", gotG, tt.wantG)
			}
			if diff := absDiffUint32(gotB, tt.wantB); diff > tolerance {
				t.Errorf("Color.RGBA() B = %v, want %v", gotB, tt.wantB)
			}
			if diff := absDiffUint32(gotA, tt.wantA); diff > tolerance {
				t.Errorf("Color.RGBA() A = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}

func TestNewColorRGBA8(t *testing.T) {
	tests := []struct {
		name       string
		r, g, b, a uint8
		want       Color
	}{
		{
			name: "white opaque",
			r:    255, g: 255, b: 255, a: 255,
			want: Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
		},
		{
			name: "red transparent",
			r:    255, g: 0, b: 0, a: 0,
			want: Color{R: 1.0, G: 0.0, B: 0.0, A: 0.0},
		},
		{
			name: "half values",
			r:    128, g: 128, b: 128, a: 128,
			want: Color{R: 128.0 / 255.0, G: 128.0 / 255.0, B: 128.0 / 255.0, A: 128.0 / 255.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewColorRGBA8(tt.r, tt.g, tt.b, tt.a)
			if !got.Equal(tt.want) {
				t.Errorf("NewColorRGBA8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewColorHex(t *testing.T) {
	tests := []struct {
		name string
		hex  uint32
		want Color
	}{
		{
			name: "white",
			hex:  0xFFFFFFFF,
			want: Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
		},
		{
			name: "red",
			hex:  0xFF0000FF,
			want: Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
		},
		{
			name: "green",
			hex:  0x00FF00FF,
			want: Color{R: 0.0, G: 1.0, B: 0.0, A: 1.0},
		},
		{
			name: "blue",
			hex:  0x0000FFFF,
			want: Color{R: 0.0, G: 0.0, B: 1.0, A: 1.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewColorHex(tt.hex)
			if !got.Equal(tt.want) {
				t.Errorf("NewColorHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_Add(t *testing.T) {
	c1 := Color{R: 0.2, G: 0.3, B: 0.4, A: 0.5}
	c2 := Color{R: 0.1, G: 0.2, B: 0.3, A: 0.4}
	want := Color{R: 0.3, G: 0.5, B: 0.7, A: 0.9}

	got := c1.Add(c2)
	if !got.Equal(want) {
		t.Errorf("Color.Add() = %v, want %v", got, want)
	}
}

func TestColor_Sub(t *testing.T) {
	c1 := Color{R: 0.5, G: 0.6, B: 0.7, A: 0.8}
	c2 := Color{R: 0.1, G: 0.2, B: 0.3, A: 0.4}
	want := Color{R: 0.4, G: 0.4, B: 0.4, A: 0.4}

	got := c1.Sub(c2)
	if !got.Equal(want) {
		t.Errorf("Color.Sub() = %v, want %v", got, want)
	}
}

func TestColor_Mul(t *testing.T) {
	c1 := Color{R: 0.5, G: 0.6, B: 0.8, A: 1.0}
	c2 := Color{R: 0.2, G: 0.5, B: 0.25, A: 0.5}
	want := Color{R: 0.1, G: 0.3, B: 0.2, A: 0.5}

	got := c1.Mul(c2)
	if !got.Equal(want) {
		t.Errorf("Color.Mul() = %v, want %v", got, want)
	}
}

func TestColor_Scale(t *testing.T) {
	c := Color{R: 0.2, G: 0.4, B: 0.6, A: 0.8}
	scale := Scalar(2.0)
	want := Color{R: 0.4, G: 0.8, B: 1.2, A: 1.6}

	got := c.Scale(scale)
	if !got.Equal(want) {
		t.Errorf("Color.Scale() = %v, want %v", got, want)
	}
}

func TestColor_Clamp01(t *testing.T) {
	tests := []struct {
		name  string
		input Color
		want  Color
	}{
		{
			name:  "already clamped",
			input: Color{R: 0.5, G: 0.5, B: 0.5, A: 0.5},
			want:  Color{R: 0.5, G: 0.5, B: 0.5, A: 0.5},
		},
		{
			name:  "over 1.0",
			input: Color{R: 1.5, G: 2.0, B: 1.2, A: 1.1},
			want:  Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
		},
		{
			name:  "under 0.0",
			input: Color{R: -0.5, G: -1.0, B: -0.2, A: -0.1},
			want:  Color{R: 0.0, G: 0.0, B: 0.0, A: 0.0},
		},
		{
			name:  "mixed range",
			input: Color{R: -0.1, G: 0.5, B: 1.5, A: 0.8},
			want:  Color{R: 0.0, G: 0.5, B: 1.0, A: 0.8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Clamp01()
			if !got.Equal(tt.want) {
				t.Errorf("Color.Clamp01() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_Premultiply(t *testing.T) {
	tests := []struct {
		name  string
		input Color
		want  Color
	}{
		{
			name:  "opaque color",
			input: Color{R: 0.5, G: 0.6, B: 0.7, A: 1.0},
			want:  Color{R: 0.5, G: 0.6, B: 0.7, A: 1.0},
		},
		{
			name:  "half alpha",
			input: Color{R: 0.8, G: 0.6, B: 0.4, A: 0.5},
			want:  Color{R: 0.4, G: 0.3, B: 0.2, A: 0.5},
		},
		{
			name:  "transparent",
			input: Color{R: 1.0, G: 1.0, B: 1.0, A: 0.0},
			want:  Color{R: 0.0, G: 0.0, B: 0.0, A: 0.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Premultiply()
			if !got.Equal(tt.want) {
				t.Errorf("Color.Premultiply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_Unpremultiply(t *testing.T) {
	tests := []struct {
		name  string
		input Color
		want  Color
	}{
		{
			name:  "opaque color",
			input: Color{R: 0.5, G: 0.6, B: 0.7, A: 1.0},
			want:  Color{R: 0.5, G: 0.6, B: 0.7, A: 1.0},
		},
		{
			name:  "half alpha premultiplied",
			input: Color{R: 0.4, G: 0.3, B: 0.2, A: 0.5},
			want:  Color{R: 0.8, G: 0.6, B: 0.4, A: 0.5},
		},
		{
			name:  "zero alpha",
			input: Color{R: 0.5, G: 0.5, B: 0.5, A: 0.0},
			want:  Color{R: 0.0, G: 0.0, B: 0.0, A: 0.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Unpremultiply()
			if !got.Equal(tt.want) {
				t.Errorf("Color.Unpremultiply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_WithAlpha(t *testing.T) {
	c := Color{R: 0.5, G: 0.6, B: 0.7, A: 1.0}
	newAlpha := Scalar(0.3)
	want := Color{R: 0.5, G: 0.6, B: 0.7, A: 0.3}

	got := c.WithAlpha(newAlpha)
	if !got.Equal(want) {
		t.Errorf("Color.WithAlpha() = %v, want %v", got, want)
	}
}

func TestColor_IsTransparent(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  bool
	}{
		{
			name:  "fully transparent",
			color: Color{R: 1.0, G: 1.0, B: 1.0, A: 0.0},
			want:  true,
		},
		{
			name:  "opaque",
			color: Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			want:  false,
		},
		{
			name:  "semi-transparent",
			color: Color{R: 1.0, G: 1.0, B: 1.0, A: 0.5},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.IsTransparent(); got != tt.want {
				t.Errorf("Color.IsTransparent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_IsOpaque(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  bool
	}{
		{
			name:  "fully opaque",
			color: Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			want:  true,
		},
		{
			name:  "transparent",
			color: Color{R: 1.0, G: 1.0, B: 1.0, A: 0.0},
			want:  false,
		},
		{
			name:  "semi-transparent",
			color: Color{R: 1.0, G: 1.0, B: 1.0, A: 0.5},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.IsOpaque(); got != tt.want {
				t.Errorf("Color.IsOpaque() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_Lerp(t *testing.T) {
	c1 := Color{R: 0.0, G: 0.0, B: 0.0, A: 0.0}
	c2 := Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0}

	tests := []struct {
		name string
		t    Scalar
		want Color
	}{
		{
			name: "t=0 returns first color",
			t:    0.0,
			want: c1,
		},
		{
			name: "t=1 returns second color",
			t:    1.0,
			want: c2,
		},
		{
			name: "t=0.5 returns midpoint",
			t:    0.5,
			want: Color{R: 0.5, G: 0.5, B: 0.5, A: 0.5},
		},
		{
			name: "t=0.25 returns quarter",
			t:    0.25,
			want: Color{R: 0.25, G: 0.25, B: 0.25, A: 0.25},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c1.Lerp(c2, tt.t)
			if !got.Equal(tt.want) {
				t.Errorf("Color.Lerp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_Blend(t *testing.T) {
	dst := Color{R: 0.5, G: 0.5, B: 0.5, A: 1.0}
	src := Color{R: 1.0, G: 0.0, B: 0.0, A: 0.5}

	tests := []struct {
		name string
		mode BlendMode
	}{
		{"clear", BlendModeClear},
		{"src", BlendModeSrc},
		{"dst", BlendModeDst},
		{"src_over", BlendModeSrcOver},
		{"dst_over", BlendModeDstOver},
		{"src_in", BlendModeSrcIn},
		{"dst_in", BlendModeDstIn},
		{"multiply", BlendModeMultiply},
		{"screen", BlendModeScreen},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dst.Blend(src, tt.mode)

			// Basic sanity checks
			if !result.IsFinite() {
				t.Errorf("Color.Blend() result is not finite: %v", result)
			}

			// Mode-specific checks
			switch tt.mode {
			case BlendModeClear:
				want := Color{R: 0, G: 0, B: 0, A: 0}
				if !result.Equal(want) {
					t.Errorf("BlendModeClear should return zero color, got %v", result)
				}
			case BlendModeSrc:
				if !result.Equal(src) {
					t.Errorf("BlendModeSrc should return source color, got %v, want %v", result, src)
				}
			case BlendModeDst:
				if !result.Equal(dst) {
					t.Errorf("BlendModeDst should return destination color, got %v, want %v", result, dst)
				}
			}
		})
	}
}

func TestBlendMode_String(t *testing.T) {
	tests := []struct {
		mode BlendMode
		want string
	}{
		{BlendModeClear, "Clear"},
		{BlendModeSrc, "Src"},
		{BlendModeDst, "Dst"},
		{BlendModeSrcOver, "SrcOver"},
		{BlendModeMultiply, "Multiply"},
		{BlendMode(255), "BlendMode(255)"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("BlendMode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPredefinedColors(t *testing.T) {
	tests := []struct {
		name  string
		color Color
	}{
		{"ColorWhite", ColorWhite()},
		{"ColorBlack", ColorBlack()},
		{"ColorRed", ColorRed()},
		{"ColorGreen", ColorGreen()},
		{"ColorBlue", ColorBlue()},
		{"ColorTransparent", ColorWhiteTransparent()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check that predefined colors are valid
			if !tt.color.IsFinite() {
				t.Errorf("%s is not finite: %v", tt.name, tt.color)
			}

			// Check specific color properties
			switch tt.name {
			case "ColorWhite":
				if tt.color.R != 1.0 || tt.color.G != 1.0 || tt.color.B != 1.0 || tt.color.A != 1.0 {
					t.Errorf("ColorWhite should be (1,1,1,1), got %v", tt.color)
				}
			case "ColorBlack":
				if tt.color.R != 0.0 || tt.color.G != 0.0 || tt.color.B != 0.0 || tt.color.A != 1.0 {
					t.Errorf("ColorBlack should be (0,0,0,1), got %v", tt.color)
				}
			case "ColorTransparent":
				if !tt.color.IsTransparent() {
					t.Errorf("ColorTransparent should be transparent, got %v", tt.color)
				}
			}
		})
	}
}

func TestRandomColor(t *testing.T) {
	// Test that RandomColor generates valid colors
	for i := 0; i < 10; i++ {
		c := RandomColor()

		if !c.IsFinite() {
			t.Errorf("RandomColor() generated non-finite color: %v", c)
		}

		if !c.IsOpaque() {
			t.Errorf("RandomColor() should be opaque, got alpha: %v", c.A)
		}

		// RGB components should be in [0, 1] range
		if c.R < 0 || c.R > 1 || c.G < 0 || c.G > 1 || c.B < 0 || c.B > 1 {
			t.Errorf("RandomColor() components out of range: %v", c)
		}
	}
}

// Edge cases and error conditions
func TestColor_EdgeCases(t *testing.T) {
	// Test with extreme values
	extremeColor := Color{R: math.MaxFloat32, G: -math.MaxFloat32, B: 0, A: 1}

	// Should not panic
	_ = extremeColor.Clamp01()
	_ = extremeColor.Scale(0.5)
	_ = extremeColor.Add(ColorWhite())

	// Test premultiply/unpremultiply cycle
	original := Color{R: 0.8, G: 0.6, B: 0.4, A: 0.7}
	cycled := original.Premultiply().Unpremultiply()

	if !original.Equal(cycled) {
		t.Errorf("Premultiply/Unpremultiply cycle failed: original=%v, cycled=%v", original, cycled)
	}
}

// Helper functions
func absDiffUint32(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}

// Check if color components are finite
func (c Color) IsFinite() bool {
	return !math.IsInf(ToFloat64(c.R), 0) && !math.IsNaN(ToFloat64(c.R)) &&
		!math.IsInf(ToFloat64(c.G), 0) && !math.IsNaN(ToFloat64(c.G)) &&
		!math.IsInf(ToFloat64(c.B), 0) && !math.IsNaN(ToFloat64(c.B)) &&
		!math.IsInf(ToFloat64(c.A), 0) && !math.IsNaN(ToFloat64(c.A))
}
