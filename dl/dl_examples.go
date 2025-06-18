package dl

import (
	"fmt"

	"github.com/opensraph/sraph/geom"
)

// Example demonstrates basic usage of the display list system.
func ExampleUsage() {
	// Create a display list builder
	builder := NewDisplayListBuilder()

	// Create some paint objects
	fillPaint := NewPaint()
	fillPaint.SetColor(ColorRed)
	fillPaint.SetStyle(PaintStyleFill)

	strokePaint := NewPaint()
	strokePaint.SetColor(ColorBlue)
	strokePaint.SetStyle(PaintStyleStroke)
	strokePaint.SetStrokeWidth(2.0)

	// Record some drawing operations
	builder.Save()
	builder.Translate(10, 10)
	builder.DrawRect(geom.NewRect[Scalar](0, 0, 100, 100), fillPaint)
	builder.Restore()

	builder.Save()
	builder.Translate(50, 50)
	builder.Scale(1.5, 1.5)
	builder.DrawCircle(geom.Point[Scalar]{X: 0, Y: 0}, 30, strokePaint)
	builder.Restore()

	// Build the display list
	displayList := builder.Build()

	fmt.Printf("Created display list with %d operations\n", displayList.OpCount())
	fmt.Printf("Bounds: %v\n", displayList.Bounds())
	fmt.Printf("Has anti-aliasing: %v\n", displayList.HasAntiAliasing())
}

// ExamplePathBuilding demonstrates path construction and usage.
func ExamplePathBuilding() {
	// Create a path builder
	pb := NewPathBuilder()

	// Build a complex path
	pb.MoveTo(50, 50)
	pb.LineTo(100, 50)
	pb.QuadTo(150, 50, 150, 100)
	pb.LineTo(150, 150)
	pb.CubicTo(150, 200, 100, 200, 50, 150)
	pb.Close()

	// Add a circle to the same path
	pb.AddCircle(geom.Point[Scalar]{X: 200, Y: 100}, 30)

	// Build the final path
	path := pb.Build()

	fmt.Printf("Created path with %d segments\n", len(path.GetSegments()))
	fmt.Printf("Path bounds: %v\n", path.Bounds())

	// Use the path in a display list
	builder := NewDisplayListBuilder()
	paint := NewPaint()
	paint.SetColor(ColorGreen)

	// TODO: This would require implementing PathSource interface
	// builder.DrawPath(path, paint)

	dl := builder.Build()
	fmt.Printf("Display list with path has %d operations\n", dl.OpCount())
}

// ExampleVertexRendering demonstrates vertex-based rendering.
func ExampleVertexRendering() {
	// Create vertices for a triangle
	verticesBuilder := NewVerticesBuilder(VertexModeTriangles, 3)
	verticesBuilder.WithColors()

	// Add triangle vertices with colors
	verticesBuilder.AddVertex(geom.Point[Scalar]{X: 50, Y: 10})
	verticesBuilder.AddVertex(geom.Point[Scalar]{X: 10, Y: 90})
	verticesBuilder.AddVertex(geom.Point[Scalar]{X: 90, Y: 90})

	// Set colors for each vertex
	verticesBuilder.SetColor(0, ColorRed)
	verticesBuilder.SetColor(1, ColorGreen)
	verticesBuilder.SetColor(2, ColorBlue)

	// Build the vertices
	vertices, err := verticesBuilder.Build()
	if err != nil {
		panic(err)
	}

	// Create display list with vertex rendering
	builder := NewDisplayListBuilder()
	paint := NewPaint()

	builder.DrawVertices(vertices, BlendModeSrcOver, paint)

	dl := builder.Build()
	fmt.Printf("Created display list with vertices: %d operations\n", dl.OpCount())
}

// ExampleTransformations demonstrates various transformation operations.
func ExampleTransformations() {
	builder := NewDisplayListBuilder()
	paint := NewPaint()
	paint.SetColor(ColorMagenta)

	rect := geom.NewRect[Scalar](0, 0, 50, 50)

	// Draw original rectangle
	builder.DrawRect(rect, paint)

	// Draw translated rectangle
	builder.Save()
	builder.Translate(60, 0)
	builder.DrawRect(rect, paint)
	builder.Restore()

	// Draw scaled rectangle
	builder.Save()
	builder.Translate(120, 0)
	builder.Scale(1.5, 1.5)
	builder.DrawRect(rect, paint)
	builder.Restore()

	// Draw rotated rectangle
	builder.Save()
	builder.Translate(200, 25)  // Move to center
	builder.Rotate(0.785398)    // 45 degrees in radians
	builder.Translate(-25, -25) // Move back
	builder.DrawRect(rect, paint)
	builder.Restore()

	dl := builder.Build()
	fmt.Printf("Created transformation example with %d operations\n", dl.OpCount())
	fmt.Printf("Bounds: %v\n", dl.Bounds())
}

// ExampleClipping demonstrates clipping operations.
func ExampleClipping() {
	builder := NewDisplayListBuilder()

	fillPaint := NewPaint()
	fillPaint.SetColor(ColorYellow)

	// Set up clipping region
	clipRect := geom.NewRect[Scalar](25, 25, 75, 75)
	builder.ClipRect(clipRect, ClipOpIntersect, true)

	// Draw something that will be clipped
	largeRect := geom.NewRect[Scalar](0, 0, 100, 100)
	builder.DrawRect(largeRect, fillPaint)

	dl := builder.Build()
	fmt.Printf("Created clipping example with %d operations\n", dl.OpCount())
}

// ExampleComplexScene demonstrates a more complex scene composition.
func ExampleComplexScene() {
	builder := NewDisplayListBuilder()

	// Background
	bgPaint := NewPaint()
	bgPaint.SetColor(ColorWhite)
	builder.DrawPaint(bgPaint)

	// Sun
	sunPaint := NewPaint()
	sunPaint.SetColor(ColorYellow)
	builder.DrawCircle(geom.Point[Scalar]{X: 200, Y: 50}, 25, sunPaint)

	// Ground
	groundPaint := NewPaint()
	groundPaint.SetColor(ColorGreen)
	groundRect := geom.NewRect[Scalar](0, 150, 300, 200)
	builder.DrawRect(groundRect, groundPaint)

	// House
	housePaint := NewPaint()
	housePaint.SetColor(NewColorRGB(0.8, 0.6, 0.4)) // Brown
	houseRect := geom.NewRect[Scalar](50, 100, 150, 150)
	builder.DrawRect(houseRect, housePaint)

	// Roof (triangle using path)
	roofPaint := NewPaint()
	roofPaint.SetColor(ColorRed)

	pb := NewPathBuilder()
	pb.MoveTo(40, 100)  // Left edge
	pb.LineTo(100, 60)  // Peak
	pb.LineTo(160, 100) // Right edge
	pb.Close()

	// TODO: This would require implementing PathSource interface
	// builder.DrawPath(pb.Build(), roofPaint)

	// Windows
	windowPaint := NewPaint()
	windowPaint.SetColor(ColorCyan)

	window1 := geom.NewRect[Scalar](70, 110, 85, 125)
	window2 := geom.NewRect[Scalar](115, 110, 130, 125)
	builder.DrawRect(window1, windowPaint)
	builder.DrawRect(window2, windowPaint)

	// Door
	doorPaint := NewPaint()
	doorPaint.SetColor(NewColorRGB(0.4, 0.2, 0.0)) // Dark brown

	doorRect := geom.NewRect[Scalar](90, 130, 110, 150)
	builder.DrawRect(doorRect, doorPaint)

	dl := builder.Build()
	fmt.Printf("Created complex scene with %d operations\n", dl.OpCount())
	fmt.Printf("Scene bounds: %v\n", dl.Bounds())
	fmt.Printf("Has anti-aliasing: %v\n", dl.HasAntiAliasing())
}

// MockRenderer provides a simple mock implementation of the Renderer interface for testing.
type MockRenderer struct {
	operationCount   int
	currentTransform geom.Matrix[Scalar]
	clipDepth        int
}

// NewMockRenderer creates a new MockRenderer.
func NewMockRenderer() *MockRenderer {
	return &MockRenderer{
		currentTransform: geom.NewMatrix[Scalar](),
	}
}

// BeginFrame implements Renderer.BeginFrame.
func (r *MockRenderer) BeginFrame(bounds geom.Rect[Scalar]) error {
	r.operationCount = 0
	fmt.Printf("Mock: Begin frame with bounds %v\n", bounds)
	return nil
}

// EndFrame implements Renderer.EndFrame.
func (r *MockRenderer) EndFrame() error {
	fmt.Printf("Mock: End frame, total operations: %d\n", r.operationCount)
	return nil
}

// Clear implements Renderer.Clear.
func (r *MockRenderer) Clear(color Color) {
	r.operationCount++
	fmt.Printf("Mock: Clear with color %v\n", color)
}

// SetTransform implements Renderer.SetTransform.
func (r *MockRenderer) SetTransform(matrix geom.Matrix[Scalar]) {
	r.currentTransform = matrix
	fmt.Printf("Mock: Set transform\n")
}

// PushClip implements Renderer.PushClip.
func (r *MockRenderer) PushClip(clip ClipRegion) {
	r.clipDepth++
	fmt.Printf("Mock: Push clip (depth: %d)\n", r.clipDepth)
}

// PopClip implements Renderer.PopClip.
func (r *MockRenderer) PopClip() {
	if r.clipDepth > 0 {
		r.clipDepth--
	}
	fmt.Printf("Mock: Pop clip (depth: %d)\n", r.clipDepth)
}

// RenderPaint implements Renderer.RenderPaint.
func (r *MockRenderer) RenderPaint(paint Paint, bounds geom.Rect[Scalar]) {
	r.operationCount++
	fmt.Printf("Mock: Render paint fill\n")
}

// RenderRect implements Renderer.RenderRect.
func (r *MockRenderer) RenderRect(rect geom.Rect[Scalar], paint Paint) {
	r.operationCount++
	fmt.Printf("Mock: Render rect %v\n", rect)
}

// RenderRoundRect implements Renderer.RenderRoundRect.
func (r *MockRenderer) RenderRoundRect(rrect geom.RoundRect[Scalar], paint Paint) {
	r.operationCount++
	fmt.Printf("Mock: Render round rect\n")
}

// RenderCircle implements Renderer.RenderCircle.
func (r *MockRenderer) RenderCircle(center geom.Point[Scalar], radius Scalar, paint Paint) {
	r.operationCount++
	fmt.Printf("Mock: Render circle at %v with radius %v\n", center, radius)
}

// RenderPath implements Renderer.RenderPath.
func (r *MockRenderer) RenderPath(path *Path, paint Paint) {
	r.operationCount++
	fmt.Printf("Mock: Render path with %d segments\n", len(path.GetSegments()))
}

// RenderLine implements Renderer.RenderLine.
func (r *MockRenderer) RenderLine(p0, p1 geom.Point[Scalar], paint Paint) {
	r.operationCount++
	fmt.Printf("Mock: Render line from %v to %v\n", p0, p1)
}

// RenderVertices implements Renderer.RenderVertices.
func (r *MockRenderer) RenderVertices(vertices *Vertices, blendMode BlendMode, paint Paint) {
	r.operationCount++
	fmt.Printf("Mock: Render vertices (%d positions)\n", vertices.VertexCount())
}

// CreateTexture implements Renderer.CreateTexture.
func (r *MockRenderer) CreateTexture(width, height int, data []byte) (TextureID, error) {
	fmt.Printf("Mock: Create texture %dx%d\n", width, height)
	return TextureID(1), nil
}

// DeleteTexture implements Renderer.DeleteTexture.
func (r *MockRenderer) DeleteTexture(id TextureID) {
	fmt.Printf("Mock: Delete texture %d\n", id)
}

// MaxTextureSize implements Renderer.MaxTextureSize.
func (r *MockRenderer) MaxTextureSize() int {
	return 4096
}

// SupportsAntiAliasing implements Renderer.SupportsAntiAliasing.
func (r *MockRenderer) SupportsAntiAliasing() bool {
	return true
}

// ExampleRenderingPipeline demonstrates the complete rendering pipeline.
func ExampleRenderingPipeline() {
	// Create a mock renderer
	renderer := NewMockRenderer()

	// Create render context
	renderContext := NewRenderContext(renderer)

	// Create a simple display list
	builder := NewDisplayListBuilder()
	paint := NewPaint()
	paint.SetColor(ColorRed)

	builder.DrawRect(geom.NewRect[Scalar](10, 10, 100, 100), paint)
	builder.Translate(50, 50)
	builder.DrawCircle(geom.Point[Scalar]{X: 0, Y: 0}, 25, paint)

	dl := builder.Build()

	// Render the display list
	bounds := geom.NewRect[Scalar](0, 0, 200, 200)
	renderContext.BeginFrame(bounds)
	renderContext.RenderDisplayList(dl)
	renderContext.EndFrame()

	fmt.Printf("Rendering pipeline completed\n")
}
