package dl

import (
	"fmt"

	"github.com/opensraph/sraph/geom"
)

// ComplexRenderingExample demonstrates advanced rendering capabilities.
func ComplexRenderingExample() {
	// Create a display list builder
	builder := NewDisplayListBuilder()

	// Set up complex transformations
	builder.Save()
	builder.Scale(2.0, 2.0)
	builder.Rotate(0.785398) // 45 degrees in radians
	builder.Translate(50, 50)

	// Create a complex path with multiple shapes
	pathBuilder := NewPathBuilder()

	// Add a star shape
	center := geom.Point[Scalar]{X: 100, Y: 100}
	pathBuilder.AddStar(center, 50, 25, 5, 0)

	// Add a spiral
	pathBuilder.AddSpiral(geom.Point[Scalar]{X: 200, Y: 100}, 10, 40, 3, true)

	// Add wavy line
	pathBuilder.AddWavyLine(
		geom.Point[Scalar]{X: 50, Y: 200},
		geom.Point[Scalar]{X: 250, Y: 200},
		10, 2)

	path := pathBuilder.Build()

	// Create paint with gradient
	gradient := NewLinearGradient(
		geom.Point[Scalar]{X: 0, Y: 0},
		geom.Point[Scalar]{X: 100, Y: 100},
		[]Color{ColorRed, ColorYellow, ColorBlue},
		[]float32{0.0, 0.5, 1.0},
		TileModeClamp,
	)

	paint := NewPaint().
		SetStyle(PaintStyleFill).
		SetAntiAlias(true)
	// TODO: Set gradient as color source
	// paint = paint.SetColorSource(gradient)

	// Add blur effect
	blurFilter := NewBlurImageFilter(5.0, 5.0, TileModeClamp)
	// TODO: Set image filter on paint
	// paint = paint.SetImageFilter(blurFilter)

	// Draw the complex path
	builder.DrawPath(path, paint)

	// Create text with styling
	textStyle := NewTextStyle().
		WithFont(NewFont("Arial", 24).WithWeight(FontWeightBold)).
		WithColor(ColorBlack).
		WithShadows([]Shadow{
			NewShadow(geom.Point[Scalar]{X: 2, Y: 2}, 4, ColorGray),
		})

	paragraphStyle := NewParagraphStyle()
	paragraphBuilder := NewParagraphBuilder(paragraphStyle)
	paragraphBuilder.PushStyle(textStyle)
	paragraphBuilder.AddText("Complex Rendering Example")
	paragraph := paragraphBuilder.Build()

	paragraph.Layout(300)
	builder.DrawParagraph(paragraph, geom.Point[Scalar]{X: 50, Y: 300})

	// Add image with effects
	// Create a simple test image
	imageData := make([]byte, 64*64*4) // 64x64 RGBA
	for i := 0; i < len(imageData); i += 4 {
		imageData[i] = 255   // R
		imageData[i+1] = 128 // G
		imageData[i+2] = 64  // B
		imageData[i+3] = 255 // A
	}

	testImage := NewBasicImage(64, 64, ImageFormatRGBA, ColorSpaceSRGB, imageData)

	// Draw image with sampling
	sampling := NewSamplingOptionsCubic(0.5, 0.5)
	imagePaint := NewPaint().SetAntiAlias(true)
	builder.DrawImageWithSampling(testImage, geom.Point[Scalar]{X: 300, Y: 50}, sampling, imagePaint)

	// Create vertices for custom geometry
	verticesBuilder := NewVerticesBuilder(VertexModeTriangles, 3)
	verticesBuilder.WithColors()

	verticesBuilder.AddVertex(geom.Point[Scalar]{X: 350, Y: 200})
	verticesBuilder.AddVertex(geom.Point[Scalar]{X: 400, Y: 280})
	verticesBuilder.AddVertex(geom.Point[Scalar]{X: 300, Y: 280})

	verticesBuilder.SetColor(0, ColorRed)
	verticesBuilder.SetColor(1, ColorGreen)
	verticesBuilder.SetColor(2, ColorBlue)

	vertices, err := verticesBuilder.Build()
	if err == nil {
		verticesPaint := NewPaint().SetAntiAlias(true)
		builder.DrawVertices(vertices, BlendModeSrcOver, verticesPaint)
	}

	// Add shadow effects
	shadowPath := NewPathBuilder()
	shadowPath.AddRect(geom.NewRect[Scalar](100, 350, 200, 400))
	shadowPathBuilt := shadowPath.Build()

	builder.DrawShadow(shadowPathBuilt, ColorBlack, 10.0, false, 1.0)

	// Restore transformation
	builder.Restore()

	// Build the final display list
	displayList := builder.Build()

	fmt.Printf("Created complex display list with %d operations\n", len(displayList.Operations()))
	fmt.Printf("Display list bounds: %v\n", displayList.Bounds())
	fmt.Printf("Has non-trivial blend modes: %t\n", displayList.HasNonTrivialBlendMode())
	fmt.Printf("Modifies transparent black: %t\n", displayList.ModifiesTransparentBlack())
}

// PerformanceExample demonstrates performance optimization techniques.
func PerformanceExample() {
	// Create a cache for expensive operations
	cache := NewRenderCache(100)

	// Create an R-tree for spatial indexing
	rtree := NewRTree(16, 8, -1)

	// Add many rectangles to the R-tree
	for i := 0; i < 1000; i++ {
		x := float32(i%100) * 10
		y := float32(i/100) * 10
		rect := geom.NewRect[Scalar](Scalar(x), Scalar(y), Scalar(x+8), Scalar(y+8))
		rtree.Insert(rect, i)
	}

	// Perform spatial queries
	queryRect := geom.NewRect[Scalar](50, 50, 150, 150)
	results := rtree.Search(queryRect)

	fmt.Printf("R-tree contains %d rectangles\n", rtree.Size())
	fmt.Printf("Query found %d intersecting rectangles\n", len(results))

	// Demonstrate region operations
	region1 := NewRegionFromRect(geom.NewRect[Scalar](0, 0, 100, 100))
	region2 := NewRegionFromRect(geom.NewRect[Scalar](50, 50, 150, 150))

	intersection := region1.Intersect(region2)
	union := region1.Union(region2)

	fmt.Printf("Region intersection area: %f\n", intersection.Area())
	fmt.Printf("Region union area: %f\n", union.Area())

	// Layer tree example
	viewport := geom.NewRect[Scalar](0, 0, 800, 600)
	layerTree := NewLayerTree(viewport, 1.0)

	// Create layers
	backgroundLayer := NewOffscreenLayer(viewport, true)
	foregroundLayer := NewOffscreenLayer(
		geom.NewRect[Scalar](100, 100, 300, 300), false)

	layerTree.SetRootLayer(backgroundLayer)

	fmt.Printf("Created layer tree with viewport: %v\n", viewport)
}

// TextRenderingExample demonstrates advanced text rendering.
func TextRenderingExample() {
	// Create a paragraph with mixed styles
	paragraphStyle := NewParagraphStyle()
	builder := NewParagraphBuilder(paragraphStyle)

	// Regular text
	regularStyle := NewTextStyle().
		WithFont(NewFont("Arial", 16)).
		WithColor(ColorBlack)
	builder.PushStyle(regularStyle)
	builder.AddText("This is ")

	// Bold text
	boldStyle := regularStyle.
		WithFont(regularStyle.Font().WithWeight(FontWeightBold)).
		WithColor(ColorRed)
	builder.PushStyle(boldStyle)
	builder.AddText("bold ")
	builder.PopStyle()

	// Italic text with shadow
	italicStyle := regularStyle.
		WithFont(regularStyle.Font().WithStyle(FontStyleItalic)).
		WithColor(ColorBlue).
		WithShadows([]Shadow{
			NewShadow(geom.Point[Scalar]{X: 1, Y: 1}, 2, ColorGray),
		})
	builder.PushStyle(italicStyle)
	builder.AddText("italic ")
	builder.PopStyle()

	builder.AddText("and normal text.")

	paragraph := builder.Build()
	paragraph.Layout(300)

	fmt.Printf("Paragraph height: %f\n", paragraph.Height())
	fmt.Printf("Paragraph width: %f\n", paragraph.Width())
	fmt.Printf("Number of lines: %d\n", len(paragraph.Lines()))
}

// EffectsExample demonstrates various visual effects.
func EffectsExample() {
	// Color filters
	matrixFilter := NewMatrixColorFilter([20]float32{
		1, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 1, 0,
	})

	// Image filters
	blurFilter := NewBlurImageFilter(10, 10, TileModeClamp)

	// Mask filters
	maskFilter := NewBlurMaskFilter(5, BlurStyleNormal)

	// Path effects
	dashEffect := NewDashPathEffect([]float32{10, 5, 5, 5}, 0)

	fmt.Printf("Created color filter: %s\n", matrixFilter.String())
	fmt.Printf("Created blur filter: %s\n", blurFilter.String())
	fmt.Printf("Created mask filter: %s\n", maskFilter.String())
	fmt.Printf("Created dash effect: %s\n", dashEffect.String())

	// Backdrop filters for glassmorphism effects
	backdropBlur := NewBlurBackdropFilter(20, 20, TileModeClamp)

	// Compositing layers
	layer := NewOffscreenLayer(geom.NewRect[Scalar](0, 0, 200, 200), false)
	compositingLayer := NewCompositingLayer(layer, BlendModeMultiply, 0.8)
	compositingLayer.SetBackdropFilter(backdropBlur)

	fmt.Printf("Created compositing layer with backdrop filter\n")
}

// ImageProcessingExample demonstrates image manipulation capabilities.
func ImageProcessingExample() {
	// Create test images
	width, height := 256, 256
	imageData := make([]byte, width*height*4)

	// Generate a gradient pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := (y*width + x) * 4
			imageData[offset] = byte(x)             // R
			imageData[offset+1] = byte(y)           // G
			imageData[offset+2] = byte((x + y) / 2) // B
			imageData[offset+3] = 255               // A
		}
	}

	image := NewBasicImage(width, height, ImageFormatRGBA, ColorSpaceSRGB, imageData)

	fmt.Printf("Created image: %s\n", image.String())
	fmt.Printf("Image is opaque: %t\n", image.IsOpaque())
	fmt.Printf("Image bytes per pixel: %d\n", image.BytesPerPixel())

	// Sampling options for different quality levels
	nearestSampling := NewSamplingOptionsNearest()
	linearSampling := NewSamplingOptionsLinear()
	cubicSampling := NewSamplingOptionsCubic(0.5, 0.5)
	catmullSampling := NewSamplingOptionsCatmullRom()

	fmt.Printf("Nearest sampling: filter=%d\n", nearestSampling.FilterMode())
	fmt.Printf("Linear sampling: filter=%d\n", linearSampling.FilterMode())
	fmt.Printf("Cubic sampling: filter=%d, cubic=%t\n",
		cubicSampling.FilterMode(), cubicSampling.UseCubic())
	fmt.Printf("Catmull-Rom sampling: filter=%d\n", catmullSampling.FilterMode())
}

// RunAllExamples runs all the example functions.
func RunAllExamples() {
	fmt.Println("=== Complex Rendering Example ===")
	ComplexRenderingExample()

	fmt.Println("\n=== Performance Example ===")
	PerformanceExample()

	fmt.Println("\n=== Text Rendering Example ===")
	TextRenderingExample()

	fmt.Println("\n=== Effects Example ===")
	EffectsExample()

	fmt.Println("\n=== Image Processing Example ===")
	ImageProcessingExample()

	fmt.Println("\nAll examples completed successfully!")
}
