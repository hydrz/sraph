package geom

import (
	"testing"
)

func TestCreateGradientBuffer_TwoColors(t *testing.T) {
	colors := []Color{
		{R: 1.0, G: 0.0, B: 0.0, A: 1.0}, // Red
		{R: 0.0, G: 1.0, B: 0.0, A: 1.0}, // Green
	}
	stops := []F32{0.0, 1.0}

	result := CreateGradientBuffer(colors, stops)

	if result.TextureSize != 2 {
		t.Errorf("Expected texture size 2, got %d", result.TextureSize)
	}

	if len(result.ColorBytes) != 8 { // 2 colors * 4 bytes (RGBA)
		t.Errorf("Expected 8 color bytes, got %d", len(result.ColorBytes))
	}

	// Check first color (red)
	if result.ColorBytes[0] != 255 || result.ColorBytes[1] != 0 ||
		result.ColorBytes[2] != 0 || result.ColorBytes[3] != 255 {
		t.Errorf("First color incorrect: got [%d, %d, %d, %d], want [255, 0, 0, 255]",
			result.ColorBytes[0], result.ColorBytes[1], result.ColorBytes[2], result.ColorBytes[3])
	}

	// Check second color (green)
	if result.ColorBytes[4] != 0 || result.ColorBytes[5] != 255 ||
		result.ColorBytes[6] != 0 || result.ColorBytes[7] != 255 {
		t.Errorf("Second color incorrect: got [%d, %d, %d, %d], want [0, 255, 0, 255]",
			result.ColorBytes[4], result.ColorBytes[5], result.ColorBytes[6], result.ColorBytes[7])
	}
}

func TestCreateGradientBuffer_MultipleColors(t *testing.T) {
	colors := []Color{
		{R: 1.0, G: 0.0, B: 0.0, A: 1.0}, // Red
		{R: 0.0, G: 1.0, B: 0.0, A: 1.0}, // Green
		{R: 0.0, G: 0.0, B: 1.0, A: 1.0}, // Blue
	}
	stops := []F32{0.0, 0.5, 1.0}

	result := CreateGradientBuffer(colors, stops)

	if !result.IsValid() {
		t.Error("Gradient should be valid")
	}

	if result.TextureSize == 0 {
		t.Error("Texture size should not be zero")
	}

	if len(result.ColorBytes) != int(result.TextureSize)*4 {
		t.Errorf("Color bytes length mismatch: got %d, want %d",
			len(result.ColorBytes), result.TextureSize*4)
	}
}

func TestCreateGradientBuffer_Interpolation(t *testing.T) {
	// Test interpolation between black and white
	colors := []Color{
		{R: 0.0, G: 0.0, B: 0.0, A: 1.0}, // Black
		{R: 1.0, G: 1.0, B: 1.0, A: 1.0}, // White
	}
	stops := []F32{0.0, 1.0}

	result := CreateGradientBuffer(colors, stops)

	// Should have at least the original colors
	if result.TextureSize < 2 {
		t.Errorf("Expected at least 2 texture size, got %d", result.TextureSize)
	}

	// First color should be black
	if result.ColorBytes[0] != 0 || result.ColorBytes[1] != 0 ||
		result.ColorBytes[2] != 0 {
		t.Error("First color should be black")
	}

	// Last color should be white
	lastIndex := (result.TextureSize - 1) * 4
	if result.ColorBytes[lastIndex] != 255 || result.ColorBytes[lastIndex+1] != 255 ||
		result.ColorBytes[lastIndex+2] != 255 {
		t.Error("Last color should be white")
	}
}

func TestCreateGradientBuffer_CloseStops(t *testing.T) {
	// Test with very close stops to check texture size limiting
	colors := []Color{
		{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
		{R: 0.0, G: 1.0, B: 0.0, A: 1.0},
		{R: 0.0, G: 0.0, B: 1.0, A: 1.0},
	}
	stops := []F32{0.0, 0.0001, 1.0} // Very close stops

	result := CreateGradientBuffer(colors, stops)

	// Should limit texture size to reasonable value
	if result.TextureSize > 1024 {
		t.Errorf("Texture size should be limited to 1024, got %d", result.TextureSize)
	}

	if !result.IsValid() {
		t.Error("Gradient should still be valid")
	}
}

func TestCreateGradientBuffer_EmptyInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty input")
		}
	}()

	CreateGradientBuffer([]Color{}, []F32{})
}

func TestCreateGradientBuffer_MismatchedLength(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for mismatched lengths")
		}
	}()

	colors := []Color{{R: 1, G: 0, B: 0, A: 1}}
	stops := []F32{0.0, 1.0} // Different length

	CreateGradientBuffer(colors, stops)
}

func TestGradientData_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		data    GradientData
		isValid bool
	}{
		{
			name:    "valid gradient",
			data:    GradientData{TextureSize: 10, ColorBytes: make([]uint8, 40)},
			isValid: true,
		},
		{
			name:    "invalid gradient - zero size",
			data:    GradientData{TextureSize: 0, ColorBytes: nil},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.IsValid(); got != tt.isValid {
				t.Errorf("GradientData.IsValid() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestNewLinearGradient(t *testing.T) {
	stops := []GradientStop{
		{Color: ColorRed, Position: 0.0},
		{Color: ColorGreen, Position: 0.5},
		{Color: ColorBlue, Position: 1.0},
	}

	gradient := NewLinearGradient(stops)

	if len(gradient.Stops) != 3 {
		t.Errorf("Expected 3 stops, got %d", len(gradient.Stops))
	}

	// Test ToBuffer
	buffer := gradient.ToBuffer()
	if !buffer.IsValid() {
		t.Error("Linear gradient buffer should be valid")
	}
}

func TestLinearGradient_ToBuffer_EmptyStops(t *testing.T) {
	gradient := NewLinearGradient([]GradientStop{})
	buffer := gradient.ToBuffer()

	if buffer.IsValid() {
		t.Error("Empty gradient should be invalid")
	}
}

func TestNewRadialGradient(t *testing.T) {
	center := Point[F32]{0.5, 0.5}
	radius := F32(1.0)
	stops := []GradientStop{
		{Color: ColorWhite, Position: 0.0},
		{Color: ColorBlack, Position: 1.0},
	}

	gradient := NewRadialGradient(center, radius, stops)

	if !gradient.Center.Eq(center) {
		t.Errorf("Center mismatch: got %v, want %v", gradient.Center, center)
	}

	if gradient.Radius != radius {
		t.Errorf("Radius mismatch: got %v, want %v", gradient.Radius, radius)
	}

	if len(gradient.Stops) != 2 {
		t.Errorf("Expected 2 stops, got %d", len(gradient.Stops))
	}
}

func TestRadialGradient_ToBuffer(t *testing.T) {
	center := Point[F32]{0, 0}
	stops := []GradientStop{
		{Color: ColorRed, Position: 0.0},
		{Color: ColorBlue, Position: 1.0},
	}

	gradient := NewRadialGradient(center, 1.0, stops)
	buffer := gradient.ToBuffer()

	if !buffer.IsValid() {
		t.Error("Radial gradient buffer should be valid")
	}

	if buffer.TextureSize == 0 {
		t.Error("Texture size should not be zero")
	}
}

// Benchmark tests
func BenchmarkCreateGradientBuffer_TwoColors(b *testing.B) {
	colors := []Color{ColorRed, ColorBlue}
	stops := []F32{0.0, 1.0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateGradientBuffer(colors, stops)
	}
}

func BenchmarkCreateGradientBuffer_MultipleColors(b *testing.B) {
	colors := []Color{ColorRed, ColorGreen, ColorBlue, ColorYellow, ColorCyan}
	stops := []F32{0.0, 0.25, 0.5, 0.75, 1.0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateGradientBuffer(colors, stops)
	}
}

func BenchmarkCreateGradientBuffer_LargeGradient(b *testing.B) {
	// Create a gradient with many stops
	numStops := 20
	colors := make([]Color, numStops)
	stops := make([]F32, numStops)

	for i := 0; i < numStops; i++ {
		colors[i] = RandomColor()
		stops[i] = F32(i) / F32(numStops-1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateGradientBuffer(colors, stops)
	}
}

// Test edge cases and error conditions
func TestGradient_EdgeCases(t *testing.T) {
	// Test with alpha variations
	t.Run("alpha_variations", func(t *testing.T) {
		colors := []Color{
			{R: 1, G: 0, B: 0, A: 0.0}, // Transparent red
			{R: 1, G: 0, B: 0, A: 1.0}, // Opaque red
		}
		stops := []F32{0.0, 1.0}

		result := CreateGradientBuffer(colors, stops)
		if !result.IsValid() {
			t.Error("Alpha gradient should be valid")
		}
	})

	// Test with identical stops
	t.Run("identical_stops", func(t *testing.T) {
		colors := []Color{ColorRed, ColorGreen, ColorBlue}
		stops := []F32{0.5, 0.5, 0.5} // All same position

		result := CreateGradientBuffer(colors, stops)
		// Should handle gracefully
		if !result.IsValid() {
			t.Error("Should handle identical stops gracefully")
		}
	})
}

// Integration test with color operations
func TestGradient_Integration(t *testing.T) {
	// Create gradient with blended colors
	baseColor := ColorRed
	blendedColor := baseColor.Blend(ColorBlue, BlendModeMultiply)

	stops := []GradientStop{
		{Color: baseColor, Position: 0.0},
		{Color: blendedColor, Position: 1.0},
	}

	gradient := NewLinearGradient(stops)
	buffer := gradient.ToBuffer()

	if !buffer.IsValid() {
		t.Error("Gradient with blended colors should be valid")
	}

	// Test with color space conversion
	linearColor := ColorRed.SRGBToLinear()
	stops2 := []GradientStop{
		{Color: linearColor, Position: 0.0},
		{Color: linearColor.LinearToSRGB(), Position: 1.0},
	}

	gradient2 := NewLinearGradient(stops2)
	buffer2 := gradient2.ToBuffer()

	if !buffer2.IsValid() {
		t.Error("Gradient with color space conversion should be valid")
	}
}
