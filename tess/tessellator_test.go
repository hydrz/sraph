package tess

import (
	"testing"

	"github.com/opensraph/sraph/geom"
)

// Test basic tessellator functionality
func TestTessellator(t *testing.T) {
	tessellator := NewTessellator[geom.F32]()
	if tessellator == nil {
		t.Fatal("NewTessellator returned nil")
	}

	// Test filled circle generation
	center := geom.Pt[geom.F32](50, 50)
	radius := geom.F32(25)
	pixelRadius := geom.F32(25)

	circleGen := tessellator.FilledCircle(center, radius, pixelRadius)
	if circleGen == nil {
		t.Fatal("FilledCircle returned nil")
	}

	if circleGen.GetTriangleType() != PrimitiveTypeTriangleStrip {
		t.Error("Expected triangle strip primitive type")
	}

	vertexCount := circleGen.GetVertexCount()
	if vertexCount <= 0 {
		t.Error("Expected positive vertex count")
	}

	// Test vertex generation
	var vertices []geom.Point[geom.F32]
	circleGen.GenerateVertices(func(p geom.Point[geom.F32]) {
		vertices = append(vertices, p)
	})

	if len(vertices) == 0 {
		t.Error("Expected vertices to be generated")
	}

	// First vertex should be the center for filled circle
	if !vertices[0].Eq(center) {
		t.Errorf("First vertex should be center, got %v, expected %v", vertices[0], center)
	}
}

// Test stroked circle generation
func TestStrokedCircle(t *testing.T) {
	tessellator := NewTessellator[geom.F32]()
	center := geom.Pt[geom.F32](50, 50)
	radius := geom.F32(25)
	halfWidth := geom.F32(2)
	pixelRadius := geom.F32(25)

	circleGen := tessellator.StrokedCircle(center, radius, halfWidth, pixelRadius)
	if circleGen == nil {
		t.Fatal("StrokedCircle returned nil")
	}

	// Test vertex generation
	var vertices []geom.Point[geom.F32]
	circleGen.GenerateVertices(func(p geom.Point[geom.F32]) {
		vertices = append(vertices, p)
	})

	if len(vertices) == 0 {
		t.Error("Expected vertices to be generated")
	}

	// For stroked circles, we should have pairs of inner/outer vertices
	if len(vertices)%2 != 0 {
		t.Error("Expected even number of vertices for stroked circle")
	}
}

// Test path tessellation
func TestPathTessellation(t *testing.T) {
	// Create a simple rectangle path source
	rect := geom.NewRect[geom.F32](0, 0, 100, 100)
	pathSource := geom.NewRectPathSource(rect)

	// Test path to filled segments
	var segments []string
	receiver := &testSegmentReceiver[geom.F32]{segments: &segments}
	PathToFilledSegments(pathSource, receiver)

	if len(segments) == 0 {
		t.Error("Expected path segments to be generated")
	}
}

// Test storage counting
func TestStorageCounting(t *testing.T) {
	rect := geom.NewRect[geom.F32](0, 0, 100, 100)
	pathSource := geom.NewRectPathSource(rect)

	pointCount, contourCount := CountFillStorage(pathSource, 1.0)

	if pointCount <= 0 {
		t.Error("Expected positive point count")
	}

	if contourCount <= 0 {
		t.Error("Expected positive contour count")
	}
}

// Test libtess tessellator
func TestLibtessTessellator(t *testing.T) {
	tessellator := NewTessellatorLibtess[geom.F32]()
	if tessellator == nil {
		t.Fatal("NewTessellatorLibtess returned nil")
	}

	// Create a simple rectangle path source
	rect := geom.NewRect[geom.F32](0, 0, 100, 100)
	pathSource := geom.NewRectPathSource(rect)

	// Test tessellation
	var callbackInvoked bool
	result := tessellator.Tessellate(pathSource, 1.0, func(vertices []geom.F32, verticesCount int, indices []uint16, indicesCount int) bool {
		callbackInvoked = true
		if verticesCount <= 0 {
			t.Error("Expected positive vertex count")
		}
		return true
	})

	if result != ResultSuccess {
		t.Errorf("Expected successful tessellation, got %v", result)
	}

	if !callbackInvoked {
		t.Error("Expected callback to be invoked")
	}
}

// Test helper types
type testSegmentReceiver[T geom.Scalar] struct {
	segments *[]string
}

func (r *testSegmentReceiver[T]) BeginContour(origin geom.Point[T], willBeClosed bool) {
	*r.segments = append(*r.segments, "BeginContour")
}

func (r *testSegmentReceiver[T]) RecordLine(p1, p2 geom.Point[T]) {
	*r.segments = append(*r.segments, "RecordLine")
}

func (r *testSegmentReceiver[T]) RecordQuad(p1, cp, p2 geom.Point[T]) {
	*r.segments = append(*r.segments, "RecordQuad")
}

func (r *testSegmentReceiver[T]) RecordConic(p1, cp, p2 geom.Point[T], weight T) {
	*r.segments = append(*r.segments, "RecordConic")
}

func (r *testSegmentReceiver[T]) RecordCubic(p1, cp1, cp2, p2 geom.Point[T]) {
	*r.segments = append(*r.segments, "RecordCubic")
}

func (r *testSegmentReceiver[T]) EndContour(origin geom.Point[T], withClose bool) {
	*r.segments = append(*r.segments, "EndContour")
}
