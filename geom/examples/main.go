package main

// var _ geom.PathReceiver[geom.I32] = (*ImagePathReceiver[geom.I32])(nil)

// type ImagePathReceiver[T geom.Scalar] struct {
// 	image image.Image
// }

// func NewImagePathReceiver[T geom.Scalar]() *ImagePathReceiver[T] {
// 	return &ImagePathReceiver[T]{
// 		image: image.NewRGBA(image.Rectangle{
// 			Min: image.Pt(0, 0),
// 			Max: image.Pt(800, 600),
// 		}),
// 	}
// }

// // Close implements geom.PathReceiver.
// func (i *ImagePathReceiver[T]) Close() {
// 	panic("unimplemented")
// }

// // ConicTo implements geom.PathReceiver.
// func (i *ImagePathReceiver[T]) ConicTo(cp geom.Point[T], p2 geom.Point[T], weight float64) bool {
// 	panic("unimplemented")
// }

// // CubicTo implements geom.PathReceiver.
// func (i *ImagePathReceiver[T]) CubicTo(cp1 geom.Point[T], cp2 geom.Point[T], p2 geom.Point[T]) {
// 	panic("unimplemented")
// }

// // LineTo implements geom.PathReceiver.
// func (i *ImagePathReceiver[T]) LineTo(p2 geom.Point[T]) {
// 	panic("unimplemented")
// }

// // MoveTo implements geom.PathReceiver.
// func (i *ImagePathReceiver[T]) MoveTo(p2 geom.Point[T], willBeClosed bool) {
// 	panic("unimplemented")
// }

// // PathEnd implements geom.PathReceiver.
// func (i *ImagePathReceiver[T]) PathEnd() {
// 	panic("unimplemented")
// }

// // QuadTo implements geom.PathReceiver.
// func (i *ImagePathReceiver[T]) QuadTo(cp geom.Point[T], p2 geom.Point[T]) {
// 	panic("unimplemented")
// }

// func main() {
// 	fmt.Println("Starting geometry examples...")

// 	// Create output directory
// 	if err := os.MkdirAll("output", 0755); err != nil {
// 		panic(err)
// 	}

// 	// Example 1: Basic Rectangle
// 	drawBasicRectangle()

// 	// Example 2: Round Rectangle
// 	drawRoundRectangle()

// 	// Example 3: Super Round Rectangle (oval-like)
// 	drawSuperRoundRectangle()

// 	// Example 4: Oval
// 	drawOval()

// 	// Example 5: Rectangles with transformations
// 	drawTransformedRectangles()

// 	// Example 6: Matrix transformations
// 	drawMatrixTransformations()

// 	// Example 7: Gradient simulation
// 	drawGradientSimulation()

// 	// Example 8: Multiple shapes composition
// 	drawCompositeShapes()

// 	// Example 9: RSTransform examples
// 	drawRSTransformExamples()

// 	fmt.Println("All examples generated successfully!")
// }

// // Helper function to save image
// func saveImage(img *image.RGBA, filename string) {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	if err := png.Encode(file, img); err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("Saved: %s\n", filename)
// }
