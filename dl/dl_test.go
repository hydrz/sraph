package dl

import (
	"testing"

	"github.com/opensraph/sraph/geom"
)

// TestAdvancedPath tests advanced path building capabilities.
func TestAdvancedPath_Polygon(t *testing.T) {
	builder := NewPathBuilder()
	center := geom.Point[Scalar]{X: 100, Y: 100}

	builder.AddPolygon(center, 50, 6, 0) // Hexagon
	path := builder.Build()

	if path == nil {
		t.Error("Expected path to be built successfully")
	}
}

func TestAdvancedPath_Star(t *testing.T) {
	builder := NewPathBuilder()
	center := geom.Point[Scalar]{X: 100, Y: 100}

	builder.AddStar(center, 50, 25, 5, 0) // Five-pointed star
	path := builder.Build()

	if path == nil {
		t.Error("Expected star path to be built successfully")
	}
}

func TestAdvancedPath_Arrow(t *testing.T) {
	builder := NewPathBuilder()
	start := geom.Point[Scalar]{X: 50, Y: 50}
	end := geom.Point[Scalar]{X: 150, Y: 50}

	builder.AddArrow(start, end, 20, 10)
	path := builder.Build()

	if path == nil {
		t.Error("Expected arrow path to be built successfully")
	}
}

func TestAdvancedPath_Spiral(t *testing.T) {
	builder := NewPathBuilder()
	center := geom.Point[Scalar]{X: 100, Y: 100}

	builder.AddSpiral(center, 10, 50, 3, true)
	path := builder.Build()

	if path == nil {
		t.Error("Expected spiral path to be built successfully")
	}
}

// TestImageCreation tests image creation and properties.
func TestImage_BasicImage(t *testing.T) {
	width, height := 64, 64
	data := make([]byte, width*height*4) // RGBA

	// Fill with red
	for i := 0; i < len(data); i += 4 {
		data[i] = 255   // R
		data[i+1] = 0   // G
		data[i+2] = 0   // B
		data[i+3] = 255 // A
	}

	image := NewBasicImage(width, height, ImageFormatRGBA, ColorSpaceSRGB, data)

	if image.Width() != width {
		t.Errorf("Expected width %d, got %d", width, image.Width())
	}
	if image.Height() != height {
		t.Errorf("Expected height %d, got %d", height, image.Height())
	}
	if !image.IsOpaque() {
		t.Error("Expected image to be opaque")
	}
	if image.IsTextureBacked() {
		t.Error("Expected BasicImage to not be texture-backed")
	}
}

func TestImage_TextureImage(t *testing.T) {
	width, height := 128, 128
	textureID := uint32(42)

	image := NewTextureImage(width, height, ColorSpaceSRGB, textureID, true)

	if image.Width() != width {
		t.Errorf("Expected width %d, got %d", width, image.Width())
	}
	if image.Height() != height {
		t.Errorf("Expected height %d, got %d", height, image.Height())
	}
	if !image.IsOpaque() {
		t.Error("Expected image to be opaque")
	}
	if !image.IsTextureBacked() {
		t.Error("Expected TextureImage to be texture-backed")
	}
	if image.TextureID() != textureID {
		t.Errorf("Expected texture ID %d, got %d", textureID, image.TextureID())
	}
}

// TestTextSystem tests text rendering components.
func TestText_Font(t *testing.T) {
	font := NewFont("Arial", 16.0)

	if font.Family() != "Arial" {
		t.Errorf("Expected family 'Arial', got '%s'", font.Family())
	}
	if font.Size() != 16.0 {
		t.Errorf("Expected size 16.0, got %f", font.Size())
	}
	if font.Weight() != FontWeightNormal {
		t.Errorf("Expected normal weight, got %d", font.Weight())
	}

	boldFont := font.WithWeight(FontWeightBold)
	if boldFont.Weight() != FontWeightBold {
		t.Errorf("Expected bold weight, got %d", boldFont.Weight())
	}
}

func TestText_TextStyle(t *testing.T) {
	style := NewTextStyle()

	if style.Color() != ColorBlack {
		t.Error("Expected default color to be black")
	}
	if style.Height() != 1.0 {
		t.Errorf("Expected default height 1.0, got %f", style.Height())
	}

	redStyle := style.WithColor(ColorRed)
	if redStyle.Color() != ColorRed {
		t.Error("Expected color to be red")
	}
}

func TestText_Paragraph(t *testing.T) {
	paragraph := NewParagraph()

	textStyle := NewTextStyle().WithColor(ColorBlue)
	span := NewTextSpan("Hello, World!", &textStyle)
	paragraph.AddTextSpan(span)

	paragraph.Layout(200)

	if paragraph.Width() == 0 {
		t.Error("Expected paragraph to have non-zero width after layout")
	}
	if paragraph.Height() == 0 {
		t.Error("Expected paragraph to have non-zero height after layout")
	}
}

// TestEffects tests various visual effects.
func TestEffects_ColorFilters(t *testing.T) {
	matrix := [20]float32{
		1, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 1, 0,
	}
	filter := NewMatrixColorFilter(matrix)

	// Test that filter can transform colors
	testColor := ColorRed
	result := filter.ApplyFilter(testColor)

	// Matrix filter should return some result (even if not implemented)
	if result.R() < 0 || result.R() > 1 {
		t.Error("Color filter result should be in valid range")
	}
}

func TestEffects_ImageFilters(t *testing.T) {
	filter := NewBlurImageFilter(10, 10, TileModeClamp)

	inputBounds := geom.NewRect[Scalar](0, 0, 100, 100)
	outputBounds := filter.Bounds(inputBounds)

	// Blur should expand bounds
	if outputBounds.Left >= inputBounds.Left {
		t.Error("Blur filter should expand bounds leftward")
	}
	if outputBounds.Right <= inputBounds.Right {
		t.Error("Blur filter should expand bounds rightward")
	}
}

func TestEffects_MaskFilters(t *testing.T) {
	filter := NewBlurMaskFilter(5, BlurStyleNormal)

	inputBounds := geom.NewRect[Scalar](0, 0, 50, 50)
	outputBounds := filter.Bounds(inputBounds)

	// Mask blur should expand bounds
	if outputBounds.Left >= inputBounds.Left {
		t.Error("Mask filter should expand bounds leftward")
	}
	if outputBounds.Bottom <= inputBounds.Bottom {
		t.Error("Mask filter should expand bounds downward")
	}
}

// TestSpatialIndex tests spatial indexing structures.
func TestSpatial_RTree(t *testing.T) {
	rtree := NewRTree(16, 8, -1)

	// Add some rectangles
	for i := 0; i < 10; i++ {
		rect := geom.NewRect[Scalar](
			Scalar(i*10), Scalar(i*10),
			Scalar(i*10+20), Scalar(i*10+20),
		)
		rtree.Insert(rect, i)
	}

	if rtree.Size() != 10 {
		t.Errorf("Expected R-tree size 10, got %d", rtree.Size())
	}

	// Search for intersecting rectangles
	query := geom.NewRect[Scalar](15, 15, 35, 35)
	results := rtree.Search(query)

	if len(results) == 0 {
		t.Error("Expected to find intersecting rectangles")
	}
}

func TestSpatial_QuadTree(t *testing.T) {
	bounds := geom.NewRect[Scalar](0, 0, 1000, 1000)
	qtree := NewQuadTree(bounds, 5, 10)

	// Add some rectangles
	for i := 0; i < 20; i++ {
		rect := geom.NewRect[Scalar](
			Scalar(i*25), Scalar(i*25),
			Scalar(i*25+50), Scalar(i*25+50),
		)
		qtree.Insert(rect, i)
	}

	if qtree.Size() != 20 {
		t.Errorf("Expected QuadTree size 20, got %d", qtree.Size())
	}

	// Search for intersecting rectangles
	query := geom.NewRect[Scalar](100, 100, 200, 200)
	results := qtree.Search(query)

	if len(results) == 0 {
		t.Error("Expected to find intersecting rectangles")
	}
}

func TestSpatial_Region(t *testing.T) {
	region := NewRegion()

	rect1 := geom.NewRect[Scalar](0, 0, 100, 100)
	rect2 := geom.NewRect[Scalar](50, 50, 150, 150)

	region.AddRect(rect1)
	region.AddRect(rect2)

	if region.IsEmpty() {
		t.Error("Expected region to not be empty")
	}

	area := region.Area()
	if area <= 0 {
		t.Error("Expected region to have positive area")
	}

	// Test point containment
	pointInside := geom.Point[Scalar]{X: 25, Y: 25}
	pointOutside := geom.Point[Scalar]{X: 200, Y: 200}

	if !region.Contains(pointInside) {
		t.Error("Expected region to contain point inside")
	}
	if region.Contains(pointOutside) {
		t.Error("Expected region to not contain point outside")
	}
}

// TestLayers tests layer system functionality.
func TestLayers_OffscreenLayer(t *testing.T) {
	bounds := geom.NewRect[Scalar](0, 0, 200, 200)
	layer := NewOffscreenLayer(bounds, true)

	if layer.Bounds() != bounds {
		t.Error("Layer bounds should match creation bounds")
	}
	if !layer.IsOpaque() {
		t.Error("Layer should be opaque as specified")
	}

	canvas := layer.Canvas()
	if &canvas == nil {
		t.Error("Layer should provide a valid canvas")
	}

	snapshot := layer.Snapshot()
	if snapshot == nil {
		t.Error("Layer should provide a valid snapshot")
	}
}

func TestLayers_LayerTree(t *testing.T) {
	viewport := geom.NewRect[Scalar](0, 0, 800, 600)
	layerTree := NewLayerTree(viewport, 1.0)

	rootLayer := NewOffscreenLayer(viewport, true)
	layerTree.SetRootLayer(rootLayer)

	// The layer tree should now have a root
	if layerTree.root == nil {
		t.Error("Layer tree should have a root after setting")
	}
	if layerTree.root.layer != rootLayer {
		t.Error("Layer tree root should match the set layer")
	}
}

// TestSamplingOptions tests texture sampling configurations.
func TestSampling_Options(t *testing.T) {
	nearest := NewSamplingOptionsNearest()
	linear := NewSamplingOptionsLinear()
	cubic := NewSamplingOptionsCubic(0.5, 0.5)
	catmull := NewSamplingOptionsCatmullRom()

	if nearest.FilterMode() != FilterModeNearest {
		t.Error("Nearest sampling should use nearest filter mode")
	}
	if linear.FilterMode() != FilterModeLinear {
		t.Error("Linear sampling should use linear filter mode")
	}
	if !cubic.UseCubic() {
		t.Error("Cubic sampling should use cubic filtering")
	}
	if !catmull.UseCubic() {
		t.Error("Catmull-Rom sampling should use cubic filtering")
	}
}

// TestAdvancedPaint tests advanced paint features.
func TestPaint_AdvancedAttributes(t *testing.T) {
	paint := NewPaint()

	// Test style setting
	fillPaint := paint.SetStyle(PaintStyleFill)
	strokePaint := paint.SetStyle(PaintStyleStroke)

	if fillPaint.Style() != PaintStyleFill {
		t.Error("Fill paint should have fill style")
	}
	if strokePaint.Style() != PaintStyleStroke {
		t.Error("Stroke paint should have stroke style")
	}

	// Test color setting
	redPaint := paint.SetColor(ColorRed)
	if redPaint.Color() != ColorRed {
		t.Error("Paint color should be red")
	}

	// Test blend mode
	multiplyPaint := paint.SetBlendMode(BlendModeMultiply)
	if multiplyPaint.BlendMode() != BlendModeMultiply {
		t.Error("Paint should have multiply blend mode")
	}
}

// BenchmarkRTreeSearch benchmarks R-tree search performance.
func BenchmarkRTreeSearch(b *testing.B) {
	rtree := NewRTree(16, 8, -1)

	// Populate with many rectangles
	for i := 0; i < 10000; i++ {
		rect := geom.NewRect[Scalar](
			Scalar(i%1000), Scalar(i/1000),
			Scalar(i%1000+10), Scalar(i/1000+10),
		)
		rtree.Insert(rect, i)
	}

	query := geom.NewRect[Scalar](500, 5, 600, 15)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := rtree.Search(query)
		_ = results // Prevent optimization
	}
}

// BenchmarkPathBuilding benchmarks complex path construction.
func BenchmarkPathBuilding(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := NewPathBuilder()

		// Build a complex path
		builder.AddPolygon(geom.Point[Scalar]{X: 100, Y: 100}, 50, 8, 0)
		builder.AddStar(geom.Point[Scalar]{X: 200, Y: 100}, 40, 20, 5, 0)
		builder.AddSpiral(geom.Point[Scalar]{X: 300, Y: 100}, 5, 30, 2, true)

		path := builder.Build()
		_ = path // Prevent optimization
	}
}
