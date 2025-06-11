package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/hydrz/sraph/geom"
)

func main() {
	fmt.Println("Image PathReceiver Example")

	// Create examples
	examples := []func(){
		rectangleExample,
		ellipseExample,
		complexPathExample,
		animationExample,
	}

	for i, example := range examples {
		fmt.Printf("Running example %d...\n", i+1)
		example()
	}

	fmt.Println("All examples completed successfully!")
}

// rectangleExample demonstrates basic rectangle rendering
func rectangleExample() {
	width, height := 400, 300

	// Create a rectangle path source
	rect := geom.NewRect[geom.F32](50, 50, 300, 200)
	rectSource := geom.NewRectPathSource(rect)

	// Render with custom colors
	img := geom.RasterizePathSourceWithColors(
		rectSource,
		width, height,
		color.RGBA{255, 0, 0, 255},     // Red stroke
		color.RGBA{0, 255, 0, 128},     // Semi-transparent green fill
		color.RGBA{255, 255, 255, 255}, // White background
	)

	saveImage(img, "rectangle_example.png")
	fmt.Println("  Rectangle example saved as rectangle_example.png")
}

// ellipseExample demonstrates ellipse/circle rendering
func ellipseExample() {
	width, height := 400, 400

	// Create multiple ellipses
	receiver := geom.NewImagePathReceiver[geom.F32](width, height)
	receiver.Clear(color.RGBA{240, 240, 240, 255}) // Light gray background

	// Large circle
	circle1 := geom.NewRect[geom.F32](50, 50, 300, 300)
	ellipseSource1 := geom.NewEllipsePathSource(circle1)
	receiver.SetStrokeColor(color.RGBA{0, 0, 255, 255})   // Blue stroke
	receiver.SetFillColor(color.RGBA{173, 216, 230, 255}) // Light blue fill
	ellipseSource1.Dispatch(receiver)
	receiver.Fill()

	// Smaller ellipse
	ellipse2 := geom.NewRect[geom.F32](120, 100, 160, 200)
	ellipseSource2 := geom.NewEllipsePathSource(ellipse2)
	receiver.SetStrokeColor(color.RGBA{255, 165, 0, 255}) // Orange stroke
	receiver.SetFillColor(color.RGBA{255, 255, 0, 180})   // Semi-transparent yellow fill
	ellipseSource2.Dispatch(receiver)
	receiver.Fill()

	saveImage(receiver.Image(), "ellipse_example.png")
	fmt.Println("  Ellipse example saved as ellipse_example.png")
}

// complexPathExample demonstrates custom path creation
func complexPathExample() {
	width, height := 500, 500
	receiver := geom.NewImagePathReceiver[geom.F32](width, height)
	receiver.Clear(color.RGBA{32, 32, 32, 255})             // Dark background
	receiver.SetStrokeColor(color.RGBA{255, 255, 255, 255}) // White stroke

	// Draw a star shape using custom PathReceiver calls
	center := geom.NewPoint[geom.F32](250, 250)
	outerRadius := geom.F32(100)
	innerRadius := geom.F32(40)
	points := 5

	// Calculate star points
	var starPoints []geom.Point[geom.F32]
	for i := 0; i < points*2; i++ {
		angle := float64(i) * math.Pi / float64(points)
		radius := outerRadius
		if i%2 == 1 {
			radius = innerRadius
		}

		x := center.X() + geom.F32(math.Cos(angle))*radius
		y := center.Y() + geom.F32(math.Sin(angle))*radius
		starPoints = append(starPoints, geom.NewPoint(x, y))
	}

	// Draw the star
	if len(starPoints) > 0 {
		receiver.MoveTo(starPoints[0], true)
		for i := 1; i < len(starPoints); i++ {
			receiver.LineTo(starPoints[i])
		}
		receiver.Close()

		// Set fill color and fill the star
		receiver.SetFillColor(color.RGBA{255, 215, 0, 255}) // Gold fill
		receiver.Fill()
		receiver.PathEnd()
	}

	// Add some decorative curves
	receiver.SetStrokeColor(color.RGBA{255, 105, 180, 255}) // Hot pink

	// Draw curved lines around the star
	for i := 0; i < 8; i++ {
		angle := float64(i) * math.Pi / 4
		startRadius := geom.F32(120)
		endRadius := geom.F32(180)

		start := geom.NewPoint(
			center.X()+geom.F32(math.Cos(angle))*startRadius,
			center.Y()+geom.F32(math.Sin(angle))*startRadius,
		)
		end := geom.NewPoint(
			center.X()+geom.F32(math.Cos(angle))*endRadius,
			center.Y()+geom.F32(math.Sin(angle))*endRadius,
		)

		// Control point for curve
		controlAngle := angle + math.Pi/8
		control := geom.NewPoint(
			center.X()+geom.F32(math.Cos(controlAngle))*((startRadius+endRadius)/2),
			center.Y()+geom.F32(math.Sin(controlAngle))*((startRadius+endRadius)/2),
		)

		receiver.MoveTo(start, false)
		receiver.QuadTo(control, end)
		receiver.PathEnd()
	}

	saveImage(receiver.Image(), "complex_path_example.png")
	fmt.Println("  Complex path example saved as complex_path_example.png")
}

// animationExample demonstrates creating multiple frames
func animationExample() {
	width, height := 300, 300
	frames := 8

	fmt.Printf("  Creating %d animation frames...\n", frames)

	for frame := 0; frame < frames; frame++ {
		receiver := geom.NewImagePathReceiver[geom.F32](width, height)
		receiver.Clear(color.RGBA{255, 255, 255, 255}) // White background

		// Animated rotating squares
		center := geom.NewPoint[geom.F32](150, 150)

		for i := 0; i < 3; i++ {
			// Calculate rotation angle for this frame and square
			baseAngle := float64(frame) * math.Pi / 4 // Base rotation
			offset := float64(i) * math.Pi * 2 / 3    // Offset for each square
			angle := geom.NewRadians[geom.F32](geom.F32(baseAngle + offset))

			// Create square corners
			size := geom.F32(30 + i*20)
			halfSize := size / 2

			corners := []geom.Point[geom.F32]{
				geom.NewPoint(-halfSize, -halfSize),
				geom.NewPoint(halfSize, -halfSize),
				geom.NewPoint(halfSize, halfSize),
				geom.NewPoint(-halfSize, halfSize),
			}

			// Rotate and translate corners
			var rotatedCorners []geom.Point[geom.F32]
			for _, corner := range corners {
				rotated := corner.Rotate(angle)
				translated := rotated.Add(center)
				rotatedCorners = append(rotatedCorners, translated)
			}

			// Set colors based on square index
			colors := []color.RGBA{
				{255, 0, 0, 200}, // Red
				{0, 255, 0, 200}, // Green
				{0, 0, 255, 200}, // Blue
			}

			receiver.SetStrokeColor(color.RGBA{0, 0, 0, 255})
			receiver.SetFillColor(colors[i])

			// Draw the square
			receiver.MoveTo(rotatedCorners[0], true)
			for j := 1; j < len(rotatedCorners); j++ {
				receiver.LineTo(rotatedCorners[j])
			}
			receiver.Close()
			receiver.Fill()
			receiver.PathEnd()
		}

		filename := fmt.Sprintf("animation_frame_%02d.png", frame)
		saveImage(receiver.Image(), filename)
	}

	fmt.Printf("  Animation frames saved as animation_frame_*.png\n")
}

// Helper function to save images
func saveImage(img *image.RGBA, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create file %s: %v", filename, err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Printf("Failed to encode PNG %s: %v", filename, err)
		return
	}
}

// demonstrateInterfaces shows how the ImagePathReceiver implements multiple interfaces
func demonstrateInterfaces() {
	receiver := geom.NewImagePathReceiver[geom.F32](200, 200)

	// PathReceiver interface usage
	var pathReceiver geom.PathReceiver[geom.F32] = receiver
	pathReceiver.MoveTo(geom.NewPoint[geom.F32](10, 10), false)
	pathReceiver.LineTo(geom.NewPoint[geom.F32](190, 190))
	pathReceiver.PathEnd()

	// image.Image interface usage
	var img image.Image = receiver
	bounds := img.Bounds()
	colorModel := img.ColorModel()
	pixelColor := img.At(100, 100)

	fmt.Printf("Image bounds: %v\n", bounds)
	fmt.Printf("Color model: %v\n", colorModel)
	fmt.Printf("Pixel at (100,100): %v\n", pixelColor)
}

// benchmarkExample demonstrates performance characteristics
func benchmarkExample() {
	fmt.Println("  Running performance test...")

	width, height := 1000, 1000
	receiver := geom.NewImagePathReceiver[geom.F32](width, height)
	receiver.Clear(color.RGBA{255, 255, 255, 255})

	// Draw many small rectangles
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			x := geom.F32(i * 10)
			y := geom.F32(j * 10)
			rect := geom.NewRect[geom.F32](x, y, x+8, y+8)

			// Alternate colors
			if (i+j)%2 == 0 {
				receiver.SetFillColor(color.RGBA{255, 0, 0, 255})
			} else {
				receiver.SetFillColor(color.RGBA{0, 0, 255, 255})
			}

			receiver.MoveTo(rect.LeftTop(), true)
			receiver.LineTo(rect.RightTop())
			receiver.LineTo(rect.RightBottom())
			receiver.LineTo(rect.LeftBottom())
			receiver.Close()
			receiver.Fill()
			receiver.PathEnd()
		}
	}

	saveImage(receiver.Image(), "performance_test.png")
	fmt.Println("  Performance test completed - 10,000 rectangles rendered")
}
