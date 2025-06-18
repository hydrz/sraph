# Tessellator Package

The `tess` package provides comprehensive tessellation functionality for 2D geometric shapes and paths. It includes both optimized shape generators and general-purpose path tessellation capabilities.

## Features

### Shape Tessellation
- **Filled Circles**: Generate triangle fan vertices for filled circles
- **Stroked Circles**: Generate triangle strip vertices for circle outlines
- **Filled Ellipses**: Generate vertices for elliptical shapes
- **Round Cap Lines**: Generate vertices for lines with rounded end caps
- **Arc Generation**: Support for filled and stroked arc segments

### Path Tessellation
- **Path to Segments**: Convert path sources to line/curve segments
- **Path to Vertices**: Convert paths to triangle vertices for rendering
- **Storage Estimation**: Calculate memory requirements for tessellation
- **Path Pruning**: Automatic cleanup of degenerate segments

### LibTess Integration
- **Arbitrary Tessellation**: Handle complex, concave polygons
- **Multiple Fill Rules**: Support for non-zero and even-odd fill rules
- **Callback-based Results**: Flexible result handling

## Core Types

### Tessellator
The main tessellator class that provides optimized shape generation:
```go
tessellator := tess.NewTessellator[geom.F32]()

// Generate a filled circle
circleGen := tessellator.FilledCircle(center, radius, pixelRadius)
circleGen.GenerateVertices(func(p geom.Point[geom.F32]) {
    // Process vertex
})
```

### Path Tessellation Functions
Convert path sources to renderable vertices:
```go
// Convert path to segments
tess.PathToFilledSegments(pathSource, segmentReceiver)

// Convert path to vertices
tess.PathToFilledVertices(pathSource, vertexWriter, scale)

// Count storage requirements
pointCount, contourCount := tess.CountFillStorage(pathSource, scale)
```

### LibTess Tessellator
For complex polygon tessellation:
```go
tessellator := tess.NewTessellatorLibtess[geom.F32]()
result := tessellator.Tessellate(pathSource, tolerance, func(vertices []geom.F32, verticesCount int, indices []uint16, indicesCount int) bool {
    // Process tessellation result
    return true
})
```

## Interfaces

### VertexWriter
Interface for collecting tessellated vertices:
```go
type VertexWriter[T geom.Scalar] interface {
    Write(point geom.Point[T])
    EndContour()
}
```

### SegmentReceiver
Interface for receiving path segments:
```go
type SegmentReceiver[T geom.Scalar] interface {
    BeginContour(origin geom.Point[T], willBeClosed bool)
    RecordLine(p1, p2 geom.Point[T])
    RecordQuad(p1, cp, p2 geom.Point[T])
    RecordConic(p1, cp, p2 geom.Point[T], weight T)
    RecordCubic(p1, cp1, cp2, p2 geom.Point[T])
    EndContour(origin geom.Point[T], withClose bool)
}
```

### VertexGenerator
Interface for shape-specific vertex generation:
```go
type VertexGenerator[T geom.Scalar] interface {
    GetTriangleType() PrimitiveType
    GetVertexCount() int
    GenerateVertices(proc TessellatedVertexProc[T])
}
```

## Usage Examples

### Basic Circle Tessellation
```go
tessellator := tess.NewTessellator[geom.F32]()
center := geom.Pt[geom.F32](100, 100)
radius := geom.F32(50)
pixelRadius := geom.F32(50)

circleGen := tessellator.FilledCircle(center, radius, pixelRadius)
var vertices []geom.Point[geom.F32]
circleGen.GenerateVertices(func(p geom.Point[geom.F32]) {
    vertices = append(vertices, p)
})
```

### Path Tessellation
```go
// Create a path source (rectangle, round rect, etc.)
rect := geom.NewRect[geom.F32](0, 0, 100, 100)
pathSource := geom.NewRectPathSource(rect)

// Convert to vertices
var vertices []geom.Point[geom.F32]
writer := &vertexWriter{vertices: &vertices}
tess.PathToFilledVertices(pathSource, writer, 1.0)
```

### Complex Polygon Tessellation
```go
tessellator := tess.NewTessellatorLibtess[geom.F32]()
result := tessellator.Tessellate(pathSource, 1.0, func(vertices []geom.F32, verticesCount int, indices []uint16, indicesCount int) bool {
    // vertices contains x,y pairs: [x0, y0, x1, y1, ...]
    // indices contains triangle indices into the vertex array
    for i := 0; i < indicesCount; i += 3 {
        // Process triangle: indices[i], indices[i+1], indices[i+2]
    }
    return true
})
```

## Performance Considerations

- The tessellator caches trigonometric values for improved performance
- Circle and ellipse tessellation is optimized for common cases
- Path tessellation includes automatic segment pruning
- Storage requirements can be pre-calculated to avoid reallocations

## Integration with Geometry Package

The tessellator works seamlessly with the `geom` package:
- Uses `geom.Point[T]` for all vertex positions
- Supports `geom.PathSource` interface for path tessellation
- Compatible with `geom.Scalar` numeric types
- Integrates with geometric primitives like rectangles and ellipses

## Thread Safety

The tessellator objects are **not thread-safe**. Create separate instances for concurrent use or implement external synchronization.
