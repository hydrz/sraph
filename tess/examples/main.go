package main

import (
	"fmt"

	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/tess"
)

func main() {
	fmt.Println("Tessellator Examples")
	fmt.Println("===================")

	// Example 1: Basic circle tessellation
	basicCircleExample()

	// Example 2: Stroked circle
	strokedCircleExample()

	// Example 3: Ellipse tessellation
	ellipseExample()

	// Example 4: Path tessellation
	pathTessellationExample()

	// Example 5: LibTess tessellation
	libtessExample()
}

func basicCircleExample() {
	fmt.Println("\n1. Basic Filled Circle:")

	tessellator := tess.NewTessellator[geom.F32]()
	center := geom.Pt[geom.F32](100, 100)
	radius := geom.F32(50)
	pixelRadius := geom.F32(50)

	circleGen := tessellator.FilledCircle(center, radius, pixelRadius)

	fmt.Printf("  Center: %v\n", center)
	fmt.Printf("  Radius: %v\n", radius)
	fmt.Printf("  Vertex count: %d\n", circleGen.GetVertexCount())
	fmt.Printf("  Primitive type: %d\n", circleGen.GetTriangleType())

	var vertices []geom.Point[geom.F32]
	circleGen.GenerateVertices(func(p geom.Point[geom.F32]) {
		vertices = append(vertices, p)
	})

	fmt.Printf("  Generated %d vertices\n", len(vertices))
	if len(vertices) > 0 {
		fmt.Printf("  First vertex (center): %v\n", vertices[0])
	}
	if len(vertices) > 1 {
		fmt.Printf("  Second vertex: %v\n", vertices[1])
	}
}

func strokedCircleExample() {
	fmt.Println("\n2. Stroked Circle:")

	tessellator := tess.NewTessellator[geom.F32]()
	center := geom.Pt[geom.F32](100, 100)
	radius := geom.F32(50)
	halfWidth := geom.F32(5) // 10-pixel wide stroke
	pixelRadius := geom.F32(50)

	circleGen := tessellator.StrokedCircle(center, radius, halfWidth, pixelRadius)

	fmt.Printf("  Center: %v\n", center)
	fmt.Printf("  Radius: %v\n", radius)
	fmt.Printf("  Stroke width: %v\n", halfWidth*2)
	fmt.Printf("  Vertex count: %d\n", circleGen.GetVertexCount())

	var vertices []geom.Point[geom.F32]
	circleGen.GenerateVertices(func(p geom.Point[geom.F32]) {
		vertices = append(vertices, p)
	})

	fmt.Printf("  Generated %d vertices\n", len(vertices))
	if len(vertices) >= 2 {
		fmt.Printf("  First inner vertex: %v\n", vertices[0])
		fmt.Printf("  First outer vertex: %v\n", vertices[1])
	}
}

func ellipseExample() {
	fmt.Println("\n3. Filled Ellipse:")

	tessellator := tess.NewTessellator[geom.F32]()
	bounds := geom.NewRect[geom.F32](50, 50, 150, 100) // 100x50 ellipse
	pixelRadius := geom.F32(50)

	ellipseGen := tessellator.FilledEllipse(bounds, pixelRadius)

	fmt.Printf("  Bounds: %v\n", bounds)
	fmt.Printf("  Center: %v\n", bounds.Center())
	fmt.Printf("  Radii: %v\n", geom.Size[geom.F32]{Width: bounds.Width() / 2, Height: bounds.Height() / 2})
	fmt.Printf("  Vertex count: %d\n", ellipseGen.GetVertexCount())

	var vertices []geom.Point[geom.F32]
	ellipseGen.GenerateVertices(func(p geom.Point[geom.F32]) {
		vertices = append(vertices, p)
	})

	fmt.Printf("  Generated %d vertices\n", len(vertices))
}

func pathTessellationExample() {
	fmt.Println("\n4. Path Tessellation:")

	// Create a rectangle path
	rect := geom.NewRect[geom.F32](0, 0, 100, 50)
	pathSource := geom.NewRectPathSource(rect)

	fmt.Printf("  Rectangle: %v\n", rect)
	fmt.Printf("  Fill type: %d\n", pathSource.FillType())
	fmt.Printf("  Is convex: %t\n", pathSource.IsConvex())

	// Count storage requirements
	pointCount, contourCount := tess.CountFillStorage(pathSource, 1.0)
	fmt.Printf("  Storage requirements: %d points, %d contours\n", pointCount, contourCount)

	// Convert to segments
	segments := make([]string, 0)
	receiver := &segmentCollector[geom.F32]{segments: &segments}
	tess.PathToFilledSegments(pathSource, receiver)

	fmt.Printf("  Generated segments: %v\n", segments)

	// Convert to vertices
	vertices := make([]geom.Point[geom.F32], 0)
	writer := &vertexCollector[geom.F32]{vertices: &vertices}
	tess.PathToFilledVertices(pathSource, writer, 1.0)

	fmt.Printf("  Generated %d vertices\n", len(vertices))
}

func libtessExample() {
	fmt.Println("\n5. LibTess Tessellation:")

	tessellator := tess.NewTessellatorLibtess[geom.F32]()

	// Create a simple triangle path (more interesting than rectangle)
	rect := geom.NewRect[geom.F32](0, 0, 100, 100)
	pathSource := geom.NewRectPathSource(rect)

	fmt.Printf("  Path bounds: %v\n", pathSource.Bounds())

	// Tessellate
	result := tessellator.Tessellate(pathSource, 1.0, func(vertices []geom.F32, verticesCount int, indices []uint16, indicesCount int) bool {
		fmt.Printf("  Tessellation result:\n")
		fmt.Printf("    Vertex count: %d\n", verticesCount)
		fmt.Printf("    Index count: %d\n", indicesCount)

		// Print first few vertices
		if verticesCount >= 4 {
			fmt.Printf("    First vertex: (%.1f, %.1f)\n", vertices[0], vertices[1])
			fmt.Printf("    Second vertex: (%.1f, %.1f)\n", vertices[2], vertices[3])
		}

		// Print first few indices
		if indicesCount >= 3 {
			fmt.Printf("    First triangle: [%d, %d, %d]\n", indices[0], indices[1], indices[2])
		}

		return true
	})

	fmt.Printf("  Result: %d\n", result)
}

// Helper types for collecting results
type segmentCollector[T geom.Scalar] struct {
	segments *[]string
}

func (s *segmentCollector[T]) BeginContour(origin geom.Point[T], willBeClosed bool) {
	*s.segments = append(*s.segments, "BeginContour")
}

func (s *segmentCollector[T]) RecordLine(p1, p2 geom.Point[T]) {
	*s.segments = append(*s.segments, "Line")
}

func (s *segmentCollector[T]) RecordQuad(p1, cp, p2 geom.Point[T]) {
	*s.segments = append(*s.segments, "Quad")
}

func (s *segmentCollector[T]) RecordConic(p1, cp, p2 geom.Point[T], weight T) {
	*s.segments = append(*s.segments, "Conic")
}

func (s *segmentCollector[T]) RecordCubic(p1, cp1, cp2, p2 geom.Point[T]) {
	*s.segments = append(*s.segments, "Cubic")
}

func (s *segmentCollector[T]) EndContour(origin geom.Point[T], withClose bool) {
	*s.segments = append(*s.segments, "EndContour")
}

type vertexCollector[T geom.Scalar] struct {
	vertices *[]geom.Point[T]
}

func (v *vertexCollector[T]) Write(point geom.Point[T]) {
	*v.vertices = append(*v.vertices, point)
}

func (v *vertexCollector[T]) EndContour() {
	// Mark end of contour if needed
}
