package main

import (
	"fmt"

	"github.com/opensraph/sraph/dl"
	"github.com/opensraph/sraph/geom"
)

func main() {
	// Test basic display list creation
	builder := dl.NewDisplayListBuilder()

	// Create a simple paint
	paint := dl.NewPaint()
	paint.SetColor(dl.ColorRed)
	paint.SetStyle(dl.PaintStyleFill)

	// Draw a rectangle
	rect := geom.NewRect[dl.Scalar](10, 10, 100, 100)
	builder.DrawRect(rect, paint)

	// Transform and draw a circle
	builder.Save()
	builder.Translate(50, 50)

	circle := geom.Point[dl.Scalar]{X: 0, Y: 0}
	builder.DrawCircle(circle, 25, paint)
	builder.Restore()

	// Build the display list
	displayList := builder.Build()

	fmt.Printf("Created display list with %d operations\n", len(displayList.Operations()))
	fmt.Printf("Display list bounds: %v\n", displayList.Bounds())

	// Test the mock renderer
	renderer := dl.NewMockRenderer()
	renderContext := dl.NewRenderContext(renderer)

	bounds := geom.NewRect[dl.Scalar](0, 0, 200, 200)
	err := renderContext.BeginFrame(bounds)
	if err != nil {
		fmt.Printf("Error beginning frame: %v\n", err)
		return
	}

	err = renderContext.RenderDisplayList(displayList)
	if err != nil {
		fmt.Printf("Error rendering display list: %v\n", err)
		return
	}

	err = renderContext.EndFrame()
	if err != nil {
		fmt.Printf("Error ending frame: %v\n", err)
		return
	}

	fmt.Println("Rendering completed successfully!")

	// Test spatial indexing
	fmt.Println("\nTesting spatial indexing...")
	rtree := dl.NewRTree(16, 8, -1)

	// Add some rectangles
	for i := 0; i < 10; i++ {
		x := dl.Scalar(i * 20)
		y := dl.Scalar(i * 15)
		rect := geom.NewRect[dl.Scalar](x, y, x+10, y+10)
		rtree.Insert(rect, i)
	}

	// Query for intersections
	queryRect := geom.NewRect[dl.Scalar](25, 20, 75, 60)
	results := rtree.Search(queryRect)
	fmt.Printf("Found %d intersecting rectangles: %v\n", len(results), results)
}
