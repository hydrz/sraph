package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/opensraph/sraph/geom"
)

func main() {
	fmt.Println("Starting geometry examples...")

	// Create output directory
	if err := os.MkdirAll("output", 0755); err != nil {
		panic(err)
	}

	// Example 1: Basic Rectangle
	drawBasicRectangle()

	// Example 2: Round Rectangle
	drawRoundRectangle()

	// Example 3: Super Round Rectangle (oval-like)
	drawSuperRoundRectangle()

	// Example 4: Oval
	drawOval()

	// Example 5: Rectangles with transformations
	drawTransformedRectangles()

	// Example 6: Matrix transformations
	drawMatrixTransformations()

	// Example 7: Gradient simulation
	drawGradientSimulation()

	// Example 8: Multiple shapes composition
	drawCompositeShapes()

	// Example 9: RSTransform examples
	drawRSTransformExamples()

	fmt.Println("All examples generated successfully!")
}

// Example 1: Basic Rectangle
func drawBasicRectangle() {
	fmt.Println("Drawing basic rectangle...")

	// Create a simple rectangle
	rect := geom.NewRect[geom.F32](50, 50, 250, 150)
	source := geom.NewRectPathSource(rect)

	// Render to image
	img := geom.RasterizePathSourceWithColors(
		source, 300, 200,
		color.RGBA{0, 0, 0, 255},       // Black stroke
		color.RGBA{100, 150, 200, 255}, // Blue fill
		color.RGBA{255, 255, 255, 255}, // White background
	)

	saveImage(img, "output/01_basic_rectangle.png")
}

// Example 2: Round Rectangle
func drawRoundRectangle() {
	fmt.Println("Drawing round rectangle...")

	rect := geom.NewRect[geom.F32](50, 50, 250, 150)
	roundRect := geom.NewRoundRectRadius(rect, 20)
	source := geom.NewRoundRectPathSource(roundRect)

	img := geom.RasterizePathSourceWithColors(
		source, 300, 200,
		color.RGBA{0, 0, 0, 255},
		color.RGBA{200, 100, 150, 255}, // Pink fill
		color.RGBA{255, 255, 255, 255},
	)

	saveImage(img, "output/02_round_rectangle.png")
}

// Example 3: Super Round Rectangle with different corner radii
func drawSuperRoundRectangle() {
	fmt.Println("Drawing super round rectangle...")

	rect := geom.NewRect[geom.F32](50, 50, 250, 150)
	// Create round rect with different radii for each corner
	// NewRoundRectLTRB(rect, left, top, right, bottom) - these are the radius values for each corner
	roundRect := geom.NewRoundRectLTRB(rect, 30, 10, 40, 25)
	source := geom.NewRoundRectPathSource(roundRect)

	img := geom.RasterizePathSourceWithColors(
		source, 300, 200,
		color.RGBA{0, 0, 0, 255},
		color.RGBA{150, 200, 100, 255}, // Green fill
		color.RGBA{255, 255, 255, 255},
	)

	saveImage(img, "output/03_super_round_rectangle.png")
}

// Example 4: Oval (ellipse)
func drawOval() {
	fmt.Println("Drawing oval...")

	rect := geom.NewRect[geom.F32](50, 50, 250, 150)
	oval := geom.NewRoundRectOval(rect)
	source := geom.NewRoundRectPathSource(oval)

	img := geom.RasterizePathSourceWithColors(
		source, 300, 200,
		color.RGBA{0, 0, 0, 255},
		color.RGBA{255, 200, 100, 255}, // Orange fill
		color.RGBA{255, 255, 255, 255},
	)

	saveImage(img, "output/04_oval.png")
}

// Example 5: Rectangles with basic transformations
func drawTransformedRectangles() {
	fmt.Println("Drawing transformed rectangles...")

	// Create image receiver manually for multiple shapes
	receiver := geom.NewImagePathReceiver[geom.F32](400, 300)
	receiver.Clear(color.RGBA{255, 255, 255, 255})

	// Original rectangle
	rect1 := geom.NewRect[geom.F32](50, 50, 150, 100)
	receiver.SetFillColor(color.RGBA{200, 100, 100, 255})
	drawRectToReceiver(rect1, receiver)

	// Scaled rectangle
	rect2 := rect1.Scale(1.5)
	receiver.SetFillColor(color.RGBA{100, 200, 100, 255})
	drawRectToReceiver(rect2, receiver)

	// Translated rectangle
	rect3 := rect1.TranslateXY(200, 100)
	receiver.SetFillColor(color.RGBA{100, 100, 200, 255})
	drawRectToReceiver(rect3, receiver)

	// Expanded rectangle
	rect4 := rect1.Expand(20).TranslateXY(0, 150)
	receiver.SetFillColor(color.RGBA{200, 200, 100, 255})
	drawRectToReceiver(rect4, receiver)

	receiver.Fill()
	saveImage(receiver.Image(), "output/05_transformed_rectangles.png")
}

// Example 6: Matrix transformations
func drawMatrixTransformations() {
	fmt.Println("Drawing matrix transformations...")

	receiver := geom.NewImagePathReceiver[geom.F32](400, 300)
	receiver.Clear(color.RGBA{240, 240, 240, 255})

	// Original rectangle
	baseRect := geom.NewRect[geom.F32](100, 100, 200, 150)

	// Create transformation matrix
	matrix := geom.NewMatrix[geom.F32]()

	// Apply rotation around center
	center := baseRect.Center()
	matrix = matrix.Translate2D(geom.NewVector2(center.X(), center.Y()))
	matrix = matrix.RotateZ(geom.NewRadians[geom.F32](0.5)) // ~30 degrees
	matrix = matrix.Translate2D(geom.NewVector2(-center.X(), -center.Y()))

	// Transform rectangle corners
	corners := baseRect.Points()
	transformedCorners := make([]geom.Point[geom.F32], 4)
	for i, corner := range corners {
		transformedCorners[i] = matrix.TransformPoint(corner)
	}

	// Draw transformed shape by connecting points
	receiver.SetStrokeColor(color.RGBA{255, 0, 0, 255})
	receiver.SetFillColor(color.RGBA{255, 100, 100, 128})

	receiver.MoveTo(transformedCorners[0], true)
	for i := 1; i < len(transformedCorners); i++ {
		receiver.LineTo(transformedCorners[i])
	}
	receiver.Close()
	receiver.PathEnd()

	// Draw original rectangle for comparison
	receiver.SetStrokeColor(color.RGBA{0, 0, 255, 255})
	receiver.SetFillColor(color.RGBA{100, 100, 255, 128})
	drawRectToReceiver(baseRect, receiver)

	receiver.Fill()
	saveImage(receiver.Image(), "output/06_matrix_transformations.png")
}

// Example 7: Gradient simulation using multiple rectangles
func drawGradientSimulation() {
	fmt.Println("Drawing gradient simulation...")

	receiver := geom.NewImagePathReceiver[geom.F32](300, 200)
	receiver.Clear(color.RGBA{255, 255, 255, 255})

	// Create gradient effect with multiple rectangles
	width := geom.F32(300)
	height := geom.F32(200)
	steps := 50

	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		x := geom.F32(i) * width / geom.F32(steps)
		rectWidth := width/geom.F32(steps) + 1 // Slightly overlap to avoid gaps

		// Interpolate colors from blue to red
		r := uint8(t * 255)
		g := uint8((1 - t) * 100)
		b := uint8((1 - t) * 255)

		rect := geom.NewRect(x, 0, x+rectWidth, height)
		receiver.SetFillColor(color.RGBA{r, g, b, 255})
		drawRectToReceiver(rect, receiver)
	}

	receiver.Fill()
	saveImage(receiver.Image(), "output/07_gradient_simulation.png")
}

// Example 8: Multiple shapes composition
func drawCompositeShapes() {
	fmt.Println("Drawing composite shapes...")

	receiver := geom.NewImagePathReceiver[geom.F32](400, 300)
	receiver.Clear(color.RGBA{250, 250, 250, 255})

	// Draw multiple overlapping round rectangles
	colors := []color.RGBA{
		{255, 100, 100, 180}, // Semi-transparent red
		{100, 255, 100, 180}, // Semi-transparent green
		{100, 100, 255, 180}, // Semi-transparent blue
		{255, 255, 100, 180}, // Semi-transparent yellow
	}

	for i, col := range colors {
		offsetX := geom.F32(i * 30)
		offsetY := geom.F32(i * 20)

		rect := geom.NewRect(50+offsetX, 50+offsetY, 200+offsetX, 150+offsetY)
		roundRect := geom.NewRoundRectRadius(rect, 25)

		receiver.SetFillColor(col)
		receiver.SetStrokeColor(color.RGBA{0, 0, 0, 255})

		roundRect.Dispatch(receiver, false)
		receiver.PathEnd()
	}

	receiver.Fill()
	saveImage(receiver.Image(), "output/08_composite_shapes.png")
}

// Example 9: RSTransform examples
func drawRSTransformExamples() {
	fmt.Println("Drawing RSTransform examples...")

	receiver := geom.NewImagePathReceiver[geom.F32](400, 300)
	receiver.Clear(color.RGBA{245, 245, 245, 255})

	// Base rectangle size
	width := geom.F32(60)
	height := geom.F32(40)

	// Create multiple RSTransforms
	transforms := []struct {
		origin geom.Point[geom.F32]
		scale  geom.F32
		angle  geom.Radians[geom.F32]
		color  color.RGBA
	}{
		{geom.NewPoint[geom.F32](100, 100), 1.0, geom.NewRadians[geom.F32](0.0), color.RGBA{255, 100, 100, 255}},
		{geom.NewPoint[geom.F32](200, 100), 1.2, geom.NewRadians[geom.F32](0.5), color.RGBA{100, 255, 100, 255}},
		{geom.NewPoint[geom.F32](300, 100), 0.8, geom.NewRadians[geom.F32](1.0), color.RGBA{100, 100, 255, 255}},
		{geom.NewPoint[geom.F32](150, 200), 1.5, geom.NewRadians[geom.F32](1.5), color.RGBA{255, 255, 100, 255}},
		{geom.NewPoint[geom.F32](250, 200), 1.0, geom.NewRadians[geom.F32](2.0), color.RGBA{255, 100, 255, 255}},
	}

	for _, tf := range transforms {
		rst := geom.NewRSTransform(tf.origin, tf.scale, tf.angle)
		quad := rst.Quad(width, height)

		receiver.SetFillColor(tf.color)
		receiver.SetStrokeColor(color.RGBA{0, 0, 0, 255})

		// Draw the transformed quad
		receiver.MoveTo(quad[0], true) // UpperLeft
		receiver.LineTo(quad[1])       // UpperRight
		receiver.LineTo(quad[3])       // LowerRight
		receiver.LineTo(quad[2])       // LowerLeft
		receiver.Close()
		receiver.PathEnd()
	}

	receiver.Fill()
	saveImage(receiver.Image(), "output/09_rstransform_examples.png")
}

// Helper function to draw a rectangle to receiver
func drawRectToReceiver(rect geom.Rect[geom.F32], receiver *geom.ImagePathReceiver[geom.F32]) {
	receiver.MoveTo(rect.LeftTop(), true)
	receiver.LineTo(rect.RightTop())
	receiver.LineTo(rect.RightBottom())
	receiver.LineTo(rect.LeftBottom())
	receiver.Close()
	receiver.PathEnd()
}

// Helper function to save image
func saveImage(img *image.RGBA, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	fmt.Printf("Saved: %s\n", filename)
}
