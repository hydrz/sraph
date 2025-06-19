// Package main demonstrates usage of the Sraph font package
package main

import (
	"fmt"
	"log"

	"github.com/opensraph/sraph/font"
	"github.com/opensraph/sraph/geom"
)

func main() {
	fmt.Println("Sraph Font Package Demo")

	// Demo 1: Font Manager
	demoFontManager()

	// Demo 2: Text Layout
	demoTextLayout()

	// Demo 3: Glyph Atlas
	demoGlyphAtlas()

	// Demo 4: Typography Context
	demoTypographyContext()
}

func demoFontManager() {
	fmt.Println("\n=== Font Manager Demo ===")

	// Create font manager
	fontManager := font.NewFontManager()

	// Load fonts would typically be done like this:
	// err := fontManager.RegisterFontFromFile("assets/fonts/Roboto-Regular.ttf", "Roboto")
	// if err != nil {
	//     log.Printf("Failed to load font: %v", err)
	// }

	// For demo purposes, let's create a mock font
	// This would normally be done by the font loader
	mockTypeface := createMockTypeface()
	mockFont := font.DefaultFont(mockTypeface, 16, font.AxisAlignmentNone)

	err := fontManager.RegisterFont(mockFont, "MockFont")
	if err != nil {
		log.Printf("Failed to register font: %v", err)
		return
	}

	// Get registered fonts
	fonts := fontManager.GetRegisteredFonts()
	fmt.Printf("Registered fonts: %v\n", fonts)

	// Get font by name
	retrievedFont, ok := fontManager.GetFont("MockFont")
	if ok {
		fmt.Printf("Retrieved font: %v\n", retrievedFont.IsValid())
	}
}

func demoTextLayout() {
	fmt.Println("\n=== Text Layout Demo ===")

	// Create text layout
	layout := font.NewTextLayout()

	// Create mock font
	mockTypeface := createMockTypeface()
	mockFont := font.DefaultFont(mockTypeface, 16, font.AxisAlignmentNone)

	// Layout text
	text := "Hello, Sraph World!"
	frame, err := layout.LayoutText(text, mockFont, 300)
	if err != nil {
		log.Printf("Failed to layout text: %v", err)
		return
	}

	if frame != nil && frame.IsValid() {
		fmt.Printf("Text frame created successfully\n")
		fmt.Printf("Run count: %d\n", frame.GetRunCount())
		fmt.Printf("Has color: %v\n", frame.HasColor())
	}

	// Measure text
	size := layout.MeasureText(text, mockFont)
	fmt.Printf("Text size: %.2fx%.2f\n", float32(size.Width), float32(size.Height))

	// Word wrap demo
	longText := "This is a very long text that should be wrapped across multiple lines when the maximum width is exceeded."
	lines := layout.WordWrap(longText, mockFont, 200)
	fmt.Printf("Word wrap result (%d lines):\n", len(lines))
	for i, line := range lines {
		fmt.Printf("  Line %d: %s\n", i+1, line)
	}
}

func demoGlyphAtlas() {
	fmt.Println("\n=== Glyph Atlas Demo ===")

	// Note: This would normally require a render context
	// For demo purposes, we'll show the interface usage

	// Create rectangle packer
	packer := font.NewSkylineRectanglePacker(512, 512)

	// Pack some rectangles
	positions := []struct {
		width, height int
		expected      bool
	}{
		{32, 32, true},
		{64, 48, true},
		{16, 24, true},
		{512, 512, false}, // Should fail - too big
	}

	for i, pos := range positions {
		point, success := packer.AddRect(pos.width, pos.height)
		fmt.Printf("Rect %d (%dx%d): success=%v", i+1, pos.width, pos.height, success)
		if success {
			fmt.Printf(", position=(%d,%d)", int32(point.X), int32(point.Y))
		}
		fmt.Println()
	}

	fmt.Printf("Atlas %.1f%% full\n", packer.GetPercentFull()*100)
}

func demoTypographyContext() {
	fmt.Println("\n=== Typography Context Demo ===")

	// Note: This would normally require a render context
	// For demo purposes, we'll show the interface usage

	fmt.Println("Typography context would be created with:")
	fmt.Println("  typographerContext, err := font.NewTypographerContext(renderContext)")
	fmt.Println("  if err != nil { ... }")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  frame, err := typographerContext.CreateTextFrame(text, font)")
	fmt.Println("  err = typographerContext.RenderTextFrame(frame, renderPass)")
}

// createMockTypeface creates a mock typeface for demonstration
func createMockTypeface() font.Typeface {
	// Create mock font data
	mockData := []byte("MOCK_FONT_DATA")

	// Create typeface info
	info := font.TypefaceInfo{
		FamilyName: "MockFont",
		StyleName:  "Regular",
		Weight:     font.WeightNormal,
		Style:      font.StyleNormal,
		Format:     font.FontFormatTTF,
		IsValid:    true,
	}

	// Create typeface
	return font.NewDefaultTypeface(mockData, info)
}

// Additional demo functions

func demoFontDescriptor() {
	fmt.Println("\n=== Font Descriptor Demo ===")

	descriptor := font.FontDescriptor{
		FamilyName:    "Roboto",
		Weight:        font.WeightBold,
		Style:         font.StyleItalic,
		Size:          18,
		AxisAlignment: font.AxisAlignmentX,
	}

	fmt.Printf("Font descriptor: %+v\n", descriptor)
	fmt.Printf("Weight: %s\n", descriptor.Weight.String())
	fmt.Printf("Style: %s\n", descriptor.Style.String())
}

func demoGlyphTypes() {
	fmt.Println("\n=== Glyph Types Demo ===")

	// Create some glyphs
	pathGlyph := font.Glyph{
		Index: 65, // 'A'
		Type:  font.GlyphTypePath,
	}

	bitmapGlyph := font.Glyph{
		Index: 128512, // Emoji
		Type:  font.GlyphTypeBitmap,
	}

	fmt.Printf("Path glyph: Index=%d, Type=%d\n", pathGlyph.Index, pathGlyph.Type)
	fmt.Printf("Bitmap glyph: Index=%d, Type=%d\n", bitmapGlyph.Index, bitmapGlyph.Type)

	// Create glyph metrics
	metrics := font.GlyphMetrics{
		AdvanceWidth:     12.5,
		AdvanceHeight:    16.0,
		LeftSideBearing:  1.0,
		RightSideBearing: 1.5,
		BoundingBox: geom.Rect[geom.F32]{
			Left:   geom.F32(0),
			Top:    geom.F32(0),
			Right:  geom.F32(12),
			Bottom: geom.F32(16),
		},
	}

	fmt.Printf("Glyph metrics: %+v\n", metrics)
}

func init() {
	// Run additional demos
	demoFontDescriptor()
	demoGlyphTypes()
}
