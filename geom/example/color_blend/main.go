package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/opensraph/sraph/geom"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// blendModeInfo contains blend mode names and descriptions
type blendModeInfo struct {
	name        string
	description string
}

// blendModeInfos provides names and descriptions for all blend modes
var blendModeInfos = []blendModeInfo{
	{"Clear", "Alpha=0"}, {"Source", "Show Src"}, {"Destination", "Show Dst"}, {"SrcOver", "Src Over Dst"},
	{"DstOver", "Dst Over Src"}, {"SrcIn", "Src Where Dst"}, {"DstIn", "Dst Where Src"},
	{"SrcOut", "Src Where !Dst"}, {"DstOut", "Dst Where !Src"}, {"SrcATop", "Src Atop Dst"}, {"DstATop", "Dst Atop Src"},
	{"Xor", "Src XOR Dst"}, {"Plus", "Src + Dst"}, {"Modulate", "Src * Dst"},
	{"Screen", "1-(1-S)*(1-D)"}, {"Overlay", "Hard/Soft Mix"}, {"Darken", "min(S,D)"}, {"Lighten", "max(S,D)"},
	{"ColorDodge", "Dodge Effect"}, {"ColorBurn", "Burn Effect"},
	{"HardLight", "Hard Light"}, {"SoftLight", "Soft Light"}, {"Difference", "|S-D|"}, {"Exclusion", "S+D-2*S*D"},
	{"Multiply", "S*D"}, {"Hue", "Hue of S"}, {"Saturation", "Sat of S"}, {"Color", "Color of S"}, {"Luminosity", "Lum of S"},
}

// BlendModeDemo creates a precise test pattern for blend mode visualization
type BlendModeDemo struct {
	rect geom.Rect
}

// At generates optimized test patterns for accurate blend mode verification
func (b *BlendModeDemo) At(x int, y int) color.Color {
	width := int(b.rect.Width())
	height := int(b.rect.Height())

	// Grid layout: 7 columns × 5 rows for 29 blend modes (0-28)
	cols := 7
	rows := 5
	cellW := width / cols
	cellH := height / rows

	col := x / cellW
	row := y / cellH

	if col >= cols {
		col = cols - 1
	}
	if row >= rows {
		row = rows - 1
	}

	modeIdx := row*cols + col
	mode := geom.BlendMode(modeIdx)

	// Show neutral background for undefined modes
	if mode > geom.BlendModeLastMode {
		return geom.NewColor(0.2, 0.2, 0.2, 1.0)
	}

	// Calculate normalized coordinates within cell [0,1]
	fx := float64((x - col*cellW)) / float64(cellW)
	fy := float64((y - row*cellH)) / float64(cellH)

	// Reserve space for text labels at top
	if fy < 0.25 {
		return geom.NewColor(0.1, 0.1, 0.1, 1.0) // Dark area for white text
	}

	// Adjust coordinates for actual test area
	testFy := (fy - 0.25) / 0.75

	// Create standardized test pattern for consistent blend mode evaluation

	// Background pattern: Simple gradient from blue to white
	bgColor := geom.ColorBlue().Lerp(geom.ColorWhite(), geom.Scalar(fx))

	// Add some geometric structure to background
	if fx > 0.3 && fx < 0.7 && testFy > 0.3 && testFy < 0.7 {
		// Central square area with different background
		bgColor = geom.ColorNavy().Lerp(geom.ColorSilver(), geom.Scalar(testFy))
	}

	// Foreground pattern: Semi-transparent shapes with clear boundaries
	var srcColor geom.Color
	var hasForeground bool

	// Left half circle: Red with 70% opacity
	centerLeft := 0.3
	if ((fx-centerLeft)*(fx-centerLeft) + (testFy-0.5)*(testFy-0.5)) < 0.15 {
		srcColor = geom.ColorRed().WithAlpha(0.7)
		hasForeground = true
	}

	// Right half circle: Green with 70% opacity
	centerRight := 0.7
	if ((fx-centerRight)*(fx-centerRight) + (testFy-0.5)*(testFy-0.5)) < 0.15 {
		if hasForeground {
			// Overlap area: Yellow with 80% opacity
			srcColor = geom.ColorYellow().WithAlpha(0.8)
		} else {
			srcColor = geom.ColorGreen().WithAlpha(0.7)
		}
		hasForeground = true
	}

	// Central vertical stripe: Blue with 60% opacity
	if fx > 0.45 && fx < 0.55 {
		if hasForeground {
			// Multiple overlaps: Cyan with 85% opacity
			srcColor = geom.ColorCyan().WithAlpha(0.85)
		} else {
			srcColor = geom.ColorBlue().WithAlpha(0.6)
		}
		hasForeground = true
	}

	// Apply blend mode
	var result geom.Color
	if hasForeground {
		result = bgColor.Blend(srcColor, mode)
	} else {
		result = bgColor
	}

	// Draw cell borders
	const borderWidth = 0.005
	if fx < borderWidth || fx > 1-borderWidth || fy < borderWidth || fy > 1-borderWidth {
		return geom.ColorWhite()
	}

	return result
}

// Bounds returns the image boundaries
func (b *BlendModeDemo) Bounds() image.Rectangle {
	return b.rect.Go()
}

// ColorModel returns the color model
func (b *BlendModeDemo) ColorModel() color.Model {
	return color.RGBAModel
}

// drawLabels renders blend mode names and descriptions on each cell
func drawLabels(img draw.Image, rect geom.Rect) {
	cols := 7
	rows := 5
	cellW := int(rect.Width()) / cols
	cellH := int(rect.Height()) / rows

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}), // White text
		Face: basicfont.Face7x13,
	}

	// Draw labels for each blend mode
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			idx := row*cols + col
			if idx >= len(blendModeInfos) {
				continue
			}

			info := blendModeInfos[idx]

			// Position text in reserved area at top of cell
			x := col*cellW + 4
			y := row*cellH + 12

			// Draw blend mode name
			drawer.Dot = fixed.P(x, y)
			drawer.DrawString(info.name)

			// Draw mode index
			drawer.Dot = fixed.P(x, y+14)
			drawer.DrawString(fmt.Sprintf("#%d", idx))

			// Draw description
			drawer.Dot = fixed.P(x, y+28)
			drawer.DrawString(info.description)
		}
	}
}

// drawLegend adds explanatory information
func drawLegend(img draw.Image, rect geom.Rect) {
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
		Face: basicfont.Face7x13,
	}

	// Position legend in bottom-right area
	legendX := int(rect.Width()) - 650
	legendY := int(rect.Height()) - 120

	// Title
	drawer.Dot = fixed.P(legendX, legendY)
	drawer.DrawString("Blend Mode Test Pattern")

	// Pattern description
	drawer.Dot = fixed.P(legendX, legendY+20)
	drawer.DrawString("Background: Blue->White gradient with navy square")

	drawer.Dot = fixed.P(legendX, legendY+35)
	drawer.DrawString("Foreground: Red circle (left) + Green circle (right) + Blue stripe (center)")

	drawer.Dot = fixed.P(legendX, legendY+50)
	drawer.DrawString("Overlap areas show blend interactions")

	// Color legend with actual test colors
	legendY += 70

	// Background sample
	bgRect := image.Rect(legendX, legendY, legendX+50, legendY+15)
	bgGrad := geom.ColorBlue().Lerp(geom.ColorWhite(), 0.5)
	for y := bgRect.Min.Y; y < bgRect.Max.Y; y++ {
		for x := bgRect.Min.X; x < bgRect.Max.X; x++ {
			img.Set(x, y, bgGrad)
		}
	}
	drawer.Dot = fixed.P(legendX+55, legendY+12)
	drawer.DrawString("Background Color")

	// Foreground samples
	legendY += 20

	// Red circle sample
	redRect := image.Rect(legendX, legendY, legendX+15, legendY+15)
	redColor := geom.ColorRed().WithAlpha(0.7)
	for y := redRect.Min.Y; y < redRect.Max.Y; y++ {
		for x := redRect.Min.X; x < redRect.Max.X; x++ {
			img.Set(x, y, redColor)
		}
	}

	// Green circle sample
	greenRect := image.Rect(legendX+20, legendY, legendX+35, legendY+15)
	greenColor := geom.ColorGreen().WithAlpha(0.7)
	for y := greenRect.Min.Y; y < greenRect.Max.Y; y++ {
		for x := greenRect.Min.X; x < greenRect.Max.X; x++ {
			img.Set(x, y, greenColor)
		}
	}

	// Blue stripe sample
	blueRect := image.Rect(legendX+40, legendY, legendX+50, legendY+15)
	blueColor := geom.ColorBlue().WithAlpha(0.6)
	for y := blueRect.Min.Y; y < blueRect.Max.Y; y++ {
		for x := blueRect.Min.X; x < blueRect.Max.X; x++ {
			img.Set(x, y, blueColor)
		}
	}

	drawer.Dot = fixed.P(legendX+55, legendY+12)
	drawer.DrawString("Foreground Colors (70%, 70%, 60% alpha)")
}

// createReferenceImage creates a reference pattern for comparison
func createReferenceImage() {
	// Create a reference showing expected results for key blend modes
	imgRect := geom.NewRect[float64](0, 0, 800, 600)
	rgba := image.NewRGBA(image.Rect(0, 0, int(imgRect.Width()), int(imgRect.Height())))

	// Fill with test pattern
	bg := geom.ColorBlue().Lerp(geom.ColorWhite(), 0.5)
	fg := geom.ColorRed().WithAlpha(0.7)

	drawer := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
		Face: basicfont.Face7x13,
	}

	y := 50
	testModes := []struct {
		mode geom.BlendMode
		name string
	}{
		{geom.BlendModeClear, "Clear (should be transparent)"},
		{geom.BlendModeSrc, "Source (should show red)"},
		{geom.BlendModeDst, "Destination (should show blue-white)"},
		{geom.BlendModeSrcOver, "SrcOver (red over blue-white)"},
		{geom.BlendModeMultiply, "Multiply (darker)"},
		{geom.BlendModeScreen, "Screen (brighter)"},
	}

	for _, test := range testModes {
		// Draw test strip
		for x := 50; x < 400; x++ {
			result := bg.Blend(fg, test.mode)
			for dy := 0; dy < 30; dy++ {
				rgba.Set(x, y+dy, result)
			}
		}

		// Draw label
		drawer.Dot = fixed.P(420, y+20)
		drawer.DrawString(test.name)

		y += 40
	}

	// Save reference
	f, _ := os.Create("output/blend_reference.png")
	defer f.Close()
	png.Encode(f, rgba)
}

func main() {
	// Create high-resolution image for detailed analysis
	imgRect := geom.NewRect[float64](0, 0, 1800, 1200)
	demo := &BlendModeDemo{rect: imgRect}

	// Generate the test pattern
	rgba := image.NewRGBA(demo.Bounds())
	draw.Draw(rgba, rgba.Bounds(), demo, image.Point{}, draw.Src)

	// Add labels and legend
	drawLabels(rgba, imgRect)
	drawLegend(rgba, imgRect)

	// Save main test image
	f, err := os.Create("output/blend_modes_test.png")
	if err != nil {
		panic(fmt.Sprintf("Failed to create output file: %v", err))
	}
	defer f.Close()

	if err := png.Encode(f, rgba); err != nil {
		panic(fmt.Sprintf("Failed to encode PNG: %v", err))
	}

	// Create reference comparison image
	createReferenceImage()

	fmt.Println("Blend mode test images saved:")
	fmt.Println("- blend_modes_test.png: Complete test grid")
	fmt.Println("- blend_reference.png: Reference comparisons")
	fmt.Printf("Generated test patterns for %d blend modes\n", len(blendModeInfos))

	// Print verification guidelines
	fmt.Println("\nVerification Guidelines:")
	fmt.Println("1. Clear mode should show transparent/black areas")
	fmt.Println("2. Source mode should show only foreground colors")
	fmt.Println("3. Destination mode should show only background")
	fmt.Println("4. SrcOver should show foreground over background")
	fmt.Println("5. Multiply should darken colors")
	fmt.Println("6. Screen should brighten colors")
	fmt.Println("7. Compare with reference image for accuracy")
}
