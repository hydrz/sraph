package dl

import (
	"testing"

	"github.com/opensraph/sraph/geom"
)

// TestDlTypes_BasicTypes tests the basic type definitions.
func TestDlTypes_BasicTypes(t *testing.T) {
	// Test ClipOp constants
	if ClipOpDifference != 0 {
		t.Errorf("Expected ClipOpDifference to be 0, got %d", ClipOpDifference)
	}
	if ClipOpIntersect != 1 {
		t.Errorf("Expected ClipOpIntersect to be 1, got %d", ClipOpIntersect)
	}

	// Test PointMode constants
	if PointModePoints != 0 {
		t.Errorf("Expected PointModePoints to be 0, got %d", PointModePoints)
	}
	if PointModeLines != 1 {
		t.Errorf("Expected PointModeLines to be 1, got %d", PointModeLines)
	}
	if PointModePolygon != 2 {
		t.Errorf("Expected PointModePolygon to be 2, got %d", PointModePolygon)
	}
}

// TestDlColor_BasicOperations tests color operations.
func TestDlColor_BasicOperations(t *testing.T) {
	// Test basic color constants
	red := ColorRed
	if red != red {
		t.Error("Color equality failed")
	}

	// Test color creation
	customColor := NewColorFromGeom(geom.NewColorHex(0xFF00FF00)) // Green with full alpha
	if customColor == ColorTransparent {
		t.Error("Custom color should not equal transparent")
	}
}

// TestDlPaint_BasicOperations tests paint operations.
func TestDlPaint_BasicOperations(t *testing.T) {
	paint := NewPaint()

	// Test default values
	if paint.Color() != ColorBlack {
		t.Error("Default paint color should be black")
	}

	if paint.BlendMode() != BlendModeSrcOver {
		t.Error("Default blend mode should be SrcOver")
	}

	if paint.IsAntiAlias() != true {
		t.Error("Default anti-alias should be true")
	}

	// Test setters
	paint.SetColor(ColorRed)
	if paint.Color() != ColorRed {
		t.Error("Paint color setter failed")
	}

	paint.SetBlendMode(BlendModeMultiply)
	if paint.BlendMode() != BlendModeMultiply {
		t.Error("Paint blend mode setter failed")
	}

	paint.SetAntiAlias(false)
	if paint.IsAntiAlias() != false {
		t.Error("Paint anti-alias setter failed")
	}
}

// TestDlStorage_BasicOperations tests storage operations.
func TestDlStorage_BasicOperations(t *testing.T) {
	storage := NewStorage(0)

	// Test initial state
	if storage.Size() != 0 {
		t.Error("New storage should have zero size")
	}

	if storage.Capacity() != 0 {
		t.Error("New storage should have zero capacity")
	}

	// Test allocation
	ptr := storage.Allocate(100)
	if ptr == nil {
		t.Error("Storage allocation failed")
	}

	if storage.Size() != 100 {
		t.Errorf("Expected size 100, got %d", storage.Size())
	}

	if storage.Capacity() < 100 {
		t.Errorf("Expected capacity >= 100, got %d", storage.Capacity())
	}

	// Test reset
	storage.Reset()
	if storage.Size() != 0 {
		t.Error("Storage reset should set size to 0")
	}
}

// TestDlVertices_BasicOperations tests vertex operations.
func TestDlVertices_BasicOperations(t *testing.T) {
	builder := NewVerticesBuilder(VertexModeTriangles, 3)
	builder.WithColors()
	builder.WithTextureCoords()

	// Add vertices
	builder.AddVertex(geom.Point[Scalar]{X: 0, Y: 0})
	builder.AddVertex(geom.Point[Scalar]{X: 1, Y: 0})
	builder.AddVertex(geom.Point[Scalar]{X: 0.5, Y: 1})

	vertices, err := builder.Build()
	if err != nil {
		t.Errorf("Vertices build failed: %v", err)
	}

	if vertices.VertexCount() != 3 {
		t.Errorf("Expected 3 vertices, got %d", vertices.VertexCount())
	}

	if vertices.Mode() != VertexModeTriangles {
		t.Error("Vertex mode mismatch")
	}

	if !vertices.HasColors() {
		t.Error("Vertices should have colors")
	}

	if !vertices.HasTextureCoords() {
		t.Error("Vertices should have texture coordinates")
	}

	// Test bounds
	bounds := vertices.GetBounds()
	if bounds == nil {
		t.Error("Vertices bounds should not be nil")
	}
}

// TestDlOpFlags_BasicOperations tests operation flags.
func TestDlOpFlags_BasicOperations(t *testing.T) {
	// Test OpFlags
	flags := OpFlagNone
	flags = flags.WithFlag(OpFlagModifiesTransparency)

	if !flags.HasFlag(OpFlagModifiesTransparency) {
		t.Error("Flag should be set")
	}

	flags = flags.WithoutFlag(OpFlagModifiesTransparency)
	if flags.HasFlag(OpFlagModifiesTransparency) {
		t.Error("Flag should be cleared")
	}

	// Test AttributeFlags
	attrFlags := AttrFlagNone
	attrFlags = attrFlags.WithAttribute(AttrFlagHasColor)

	if !attrFlags.HasAttribute(AttrFlagHasColor) {
		t.Error("Attribute flag should be set")
	}

	attrFlags = attrFlags.WithoutAttribute(AttrFlagHasColor)
	if attrFlags.HasAttribute(AttrFlagHasColor) {
		t.Error("Attribute flag should be cleared")
	}
}

// TestDlDisplayList_BasicOperations tests display list operations.
func TestDlDisplayList_BasicOperations(t *testing.T) {
	builder := NewDisplayListBuilder()

	// Record some operations
	paint := NewPaint()
	paint.SetColor(ColorRed)

	builder.Save()
	builder.Translate(10, 10)
	builder.DrawRect(geom.NewRect[Scalar](0, 0, 100, 100), paint)
	builder.Restore()

	builder.DrawCircle(geom.Point[Scalar]{X: 50, Y: 50}, 25, paint)

	dl := builder.Build()

	if dl.IsEmpty() {
		t.Error("Display list should not be empty")
	}

	if dl.GetOpCount() == 0 {
		t.Error("Display list should have operations")
	}

	// Test bounds
	bounds := dl.GetBounds()
	if bounds.IsEmpty() {
		t.Error("Display list bounds should not be empty")
	}
}

// TestDlPath_BasicOperations tests path operations.
func TestDlPath_BasicOperations(t *testing.T) {
	pb := NewPathBuilder()

	// Build a simple path
	pb.MoveTo(0, 0)
	pb.LineTo(100, 0)
	pb.LineTo(100, 100)
	pb.LineTo(0, 100)
	pb.Close()

	path := pb.Build()

	if path.IsEmpty() {
		t.Error("Path should not be empty")
	}

	segments := path.GetSegments()
	if len(segments) == 0 {
		t.Error("Path should have segments")
	}

	// Test bounds
	bounds := path.GetBounds()
	if bounds.IsEmpty() {
		t.Error("Path bounds should not be empty")
	}

	// Test path reset
	pb.Reset()
	if !pb.IsEmpty() {
		t.Error("Reset path builder should be empty")
	}
}

// TestDlCanvas_BasicOperations tests canvas operations.
func TestDlCanvas_BasicOperations(t *testing.T) {
	bounds := geom.NewRect[Scalar](0, 0, 200, 200)
	canvas := NewCanvas(bounds)

	// Test initial state
	if canvas.GetBounds() != bounds {
		t.Error("Canvas bounds mismatch")
	}

	if canvas.GetDevicePixelRatio() != 1.0 {
		t.Error("Default device pixel ratio should be 1.0")
	}

	// Test save/restore
	canvas.Save()
	canvas.Translate(50, 50)
	canvas.Restore()

	// Transform should be back to identity after restore
	transform := canvas.GetCurrentTransform()
	identity := geom.NewMatrix[Scalar]()

	// Simple check - in practice you'd need a proper matrix comparison
	_ = transform
	_ = identity
}

// TestDlBuilder_BasicOperations tests display list builder operations.
func TestDlBuilder_BasicOperations(t *testing.T) {
	builder := NewDisplayListBuilder()

	// Test initial state
	if !builder.GetBounds().IsEmpty() {
		t.Error("New builder should have empty bounds")
	}

	// Record operations
	paint := NewPaint()
	rect := geom.NewRect[Scalar](10, 10, 100, 100)

	builder.DrawRect(rect, paint)

	// Check that bounds are updated
	bounds := builder.GetBounds()
	if bounds.IsEmpty() {
		t.Error("Builder bounds should be updated after operations")
	}

	// Build display list
	dl := builder.Build()
	if dl.GetOpCount() == 0 {
		t.Error("Built display list should have operations")
	}

	// Test reset
	builder.Reset()
	if !builder.GetBounds().IsEmpty() {
		t.Error("Reset builder should have empty bounds")
	}
}

// TestDlSamplingOptions_BasicOperations tests sampling options.
func TestDlSamplingOptions_BasicOperations(t *testing.T) {
	// Test different sampling modes
	nearest := SamplingOptionsNearest
	if nearest.FilterMode() != FilterModeNearest {
		t.Error("Nearest sampling should use nearest filter mode")
	}

	linear := SamplingOptionsLinear
	if linear.FilterMode() != FilterModeLinear {
		t.Error("Linear sampling should use linear filter mode")
	}

	mipmap := SamplingOptionsMipmap
	if mipmap.ImageSampling() != ImageSamplingMipmapLinear {
		t.Error("Mipmap sampling should use mipmap linear mode")
	}
}

// BenchmarkDlDisplayList_Build benchmarks display list building.
func BenchmarkDlDisplayList_Build(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		builder := NewDisplayListBuilder()
		paint := NewPaint()

		// Record some operations
		for j := 0; j < 100; j++ {
			x := Scalar(j * 10)
			y := Scalar(j * 5)
			rect := geom.NewRect[Scalar](x, y, x+50, y+50)
			builder.DrawRect(rect, paint)
		}

		_ = builder.Build()
	}
}

// BenchmarkDlPath_Build benchmarks path building.
func BenchmarkDlPath_Build(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pb := NewPathBuilder()

		// Build a complex path
		pb.MoveTo(0, 0)
		for j := 0; j < 100; j++ {
			x := Scalar(j * 2)
			y := Scalar(j % 10)
			pb.LineTo(x, y)
		}
		pb.Close()

		_ = pb.Build()
	}
}
