# Display List (dl) Package

A comprehensive Go implementation of a 2D graphics display list system, inspired by Flutter's display list architecture. This package provides high-performance 2D graphics rendering with support for complex visual effects, text rendering, image processing, and spatial indexing.

## Overview

The display list package provides a complete 2D graphics system that records, optimizes, and replays drawing operations efficiently. It's designed to be backend-agnostic while providing advanced features like text layout, image processing, visual effects, and spatial indexing for high-performance graphics applications.

## Key Features

### Core Display List System
- **DisplayList**: Immutable sequence of recorded drawing operations with bounds tracking
- **DisplayListBuilder**: Records operations with automatic optimization and attribute tracking
- **OpReceiver**: Interface for executing operations on different backends
- **Canvas**: Immediate-mode rendering with transform/clip stacks and layer support
- **Renderer**: Pluggable backend interface for GPU/CPU rendering

### Advanced Graphics Primitives
- **Paint**: Comprehensive drawing attributes (color, effects, blend modes, styles)
- **Path**: Complex shapes with Boolean operations and advanced path effects
- **Image**: Multi-format image support (CPU/GPU backed) with advanced sampling
- **Text**: Rich text with mixed styles, layout, and internationalization support
- **Vertices**: Efficient triangle-based rendering
- **Color**: RGBA color representation with color space support
- **Transformations**: Matrix-based 2D/3D transformations

### Visual Effects & Filters
- **ColorSource**: Gradients, patterns, and procedural color generation
- **ColorFilter**: Matrix-based color transformations
- **ImageFilter**: Blur effects and advanced image processing
- **MaskFilter**: Alpha channel effects
- **PathEffect**: Dash patterns and path modifications
- **BackdropFilter**: Background blur and effects

### Spatial Indexing & Optimization
- **R-Tree**: Efficient spatial indexing for culling and intersection queries
- **QuadTree**: Alternative spatial indexing structure
- **Region**: Non-overlapping rectangle collections
- **RenderCache**: Caching for expensive rendering operations

### Text System
- **Font**: Font family, size, weight, and style management
- **TextStyle**: Rich text styling with colors, decorations, and effects
- **Paragraph**: Complex text layout with mixed styles
- **TextSpan**: Hierarchical text structure support

### Layer System
- **Layer**: Off-screen rendering targets for complex compositing
- **LayerTree**: Hierarchical layer management
- **CompositingLayer**: Advanced compositing with backdrop filters

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ DisplayList     │    │ Canvas/Renderer │    │ Graphics Backend│
│ Builder         │───▶│                 │───▶│ (GPU/CPU/etc)   │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ Records ops     │    │ Executes ops    │    │ Actual rendering│
│ Immutable result│    │ Manages state   │    │ Platform-specific│
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Basic Usage

```go
package main

import (
    "github.com/opensraph/sraph/dl"
    "github.com/opensraph/sraph/geom"
)


func main() {
    // Create a builder to record operations
    builder := dl.NewDisplayListBuilder()

    // Create paint for drawing
    paint := dl.NewPaint()
    paint = paint.SetColor(dl.ColorRed)

    // Record drawing operations
    builder.Save()
    builder.Translate(10, 10)
    builder.DrawRect(geom.NewRect[dl.Scalar](0, 0, 100, 100), paint)
    builder.Restore()

    builder.DrawCircle(geom.Point[dl.Scalar]{X: 150, Y: 75}, 25, paint)

    // Build immutable display list
    displayList := builder.Build()

    // Execute on a canvas
    canvas := dl.NewCanvas(geom.NewRect[dl.Scalar](0, 0, 300, 200))
    displayList.Dispatch(canvas)
}
```

## Advanced Features

### Complex Path Construction
```go
pb := dl.NewPathBuilder()

// Basic shapes
pb.MoveTo(50, 50)
pb.LineTo(100, 50)
pb.QuadTo(150, 50, 150, 100)
pb.Close()

// Advanced shapes
pb.AddCircle(geom.Point[dl.Scalar]{X: 200, Y: 100}, 30)
pb.AddPolygon(geom.Point[dl.Scalar]{X: 300, Y: 100}, 40, 6, 0) // Hexagon
pb.AddStar(geom.Point[dl.Scalar]{X: 400, Y: 100}, 50, 25, 5, 0) // Star
pb.AddSpiral(geom.Point[dl.Scalar]{X: 500, Y: 100}, 10, 40, 3, true) // Spiral

path := pb.Build()
```

### Vertex-Based Rendering
```go
builder := dl.NewVerticesBuilder(dl.VertexModeTriangles, 3)
builder.WithColors().WithTextureCoords()

// Add vertices with positions, colors, and texture coordinates
builder.AddVertexFull(
    geom.Point[dl.Scalar]{X: 0, Y: 0},
    geom.Point[dl.Scalar]{X: 0, Y: 0}, // UV
    dl.ColorRed,
)
builder.AddVertexFull(
    geom.Point[dl.Scalar]{X: 1, Y: 0},
    geom.Point[dl.Scalar]{X: 1, Y: 0}, // UV
    dl.ColorGreen,
)
builder.AddVertexFull(
    geom.Point[dl.Scalar]{X: 0.5, Y: 1},
    geom.Point[dl.Scalar]{X: 0.5, Y: 1}, // UV
    dl.ColorBlue,
)

vertices, _ := builder.Build()
```

### Rich Text Rendering
```go
// Create text styles
headerStyle := dl.NewTextStyle().
    WithFont(dl.NewFont("Arial", 24).WithWeight(dl.FontWeightBold)).
    WithColor(dl.ColorBlack)

bodyStyle := dl.NewTextStyle().
    WithFont(dl.NewFont("Arial", 16)).
    WithColor(dl.ColorGray)

// Build paragraph with mixed styles
paragraphBuilder := dl.NewParagraphBuilder(dl.NewParagraphStyle())
paragraphBuilder.PushStyle(headerStyle)
paragraphBuilder.AddText("Header Text\n")
paragraphBuilder.PopStyle()

paragraphBuilder.PushStyle(bodyStyle)
paragraphBuilder.AddText("This is body text with different styling.")
paragraphBuilder.PopStyle()

paragraph := paragraphBuilder.Build()
paragraph.Layout(300) // Set width constraint

// Draw the paragraph
builder.DrawParagraph(paragraph, geom.Point[dl.Scalar]{X: 10, Y: 10})
```

### Visual Effects
```go
paint := dl.NewPaint()

// Add blur effect
blurFilter := dl.NewBlurImageFilter(5, 5, dl.TileModeClamp)
paint = paint.SetImageFilter(blurFilter)

// Add color transformation
colorMatrix := [20]float32{
    0.393, 0.769, 0.189, 0, 0, // Red channel
    0.349, 0.686, 0.168, 0, 0, // Green channel
    0.272, 0.534, 0.131, 0, 0, // Blue channel
    0,     0,     0,     1, 0, // Alpha channel
}
colorFilter := dl.NewMatrixColorFilter(colorMatrix)
paint = paint.SetColorFilter(colorFilter)

// Add drop shadow
maskFilter := dl.NewBlurMaskFilter(3, dl.BlurStyleNormal)
paint = paint.SetMaskFilter(maskFilter)
```

### Spatial Indexing for Performance
```go
// Create spatial index for efficient culling
rtree := dl.NewRTree(16, 8, -1)

// Add drawable objects to index
for i, drawable := range drawables {
    rtree.Insert(drawable.Bounds(), i)
}

// Query for visible objects
viewportBounds := geom.NewRect[dl.Scalar](0, 0, 800, 600)
visibleIds := rtree.Search(viewportBounds)

// Only render visible objects
for _, id := range visibleIds {
    drawables[id].Render(builder)
}
```

### Layer-Based Compositing
```go
// Create off-screen layer
layerBounds := geom.NewRect[dl.Scalar](0, 0, 200, 200)
layer := dl.NewOffscreenLayer(layerBounds, false)

// Draw to layer
layerCanvas := layer.Canvas()
layerCanvas.DrawRect(layerBounds, backgroundPaint)
layerCanvas.DrawCircle(center, radius, circlePaint)

// Apply backdrop filter
backdropFilter := dl.NewBlurBackdropFilter(10, 10, dl.TileModeClamp)
compositing := dl.NewCompositingLayer(layer, dl.BlendModeMultiply, 0.8)
compositing.SetBackdropFilter(backdropFilter)

// Composite layer onto main canvas
compositing.Composite(mainCanvas, backdrop)
```

### Custom Renderer Implementation
```go
type MyRenderer struct {
    // Implementation specific fields
}

func (r *MyRenderer) BeginFrame(bounds geom.Rect[dl.Scalar]) error {
    // Initialize frame rendering
    return nil
}

func (r *MyRenderer) RenderRect(rect geom.Rect[dl.Scalar], paint dl.Paint) {
    // Custom rectangle rendering logic
}

func (r *MyRenderer) RenderPath(path *dl.Path, paint dl.Paint) {
    // Custom path rendering logic
}

// Implement all required Renderer interface methods...

// Use custom renderer
renderer := &MyRenderer{}
context := dl.NewRenderContext(renderer)
context.RenderDisplayList(displayList)
```

## Supported Operations

### Drawing Operations
- **Basic Shapes**: Rectangles, circles, ovals, rounded rectangles
- **Complex Paths**: Bezier curves, polygons, stars, spirals
- **Lines & Points**: Single lines, point arrays, polylines
- **Text**: Rich text with mixed styles and layout
- **Images**: Raster images with sampling options
- **Vertices**: Triangle meshes with colors and textures

### Transformation Operations
- **2D Transforms**: Translate, scale, rotate, skew
- **3D Transforms**: Full perspective transformation matrices
- **Efficient 2D**: RSTransform for optimized sprite rendering

### Clipping Operations
- **Rectangle Clipping**: Axis-aligned and rotated rectangles
- **Rounded Rectangle Clipping**: Smooth corner clipping
- **Path Clipping**: Complex shape clipping with Boolean operations

### State Management
- **Save/Restore**: Hierarchical state management
- **Layer Support**: Off-screen rendering targets
- **Attribute Tracking**: Automatic optimization hints

## Performance Features

- **Bounds Tracking**: Automatic calculation of operation bounds for culling
- **Attribute Flags**: Operation metadata for rendering optimization
- **Immutable Design**: Safe sharing across threads
- **Memory Pooling**: Efficient storage management
- **State Batching**: Minimize state changes during rendering
- **Spatial Indexing**: R-tree and QuadTree for fast intersection queries
- **Wang's Formula**: Optimal curve tessellation

## Current Implementation Status

### ✅ Completed
- [x] Core display list architecture
- [x] Complete shape rendering (rect, circle, path, etc.)
- [x] Advanced path construction with complex shapes
- [x] Comprehensive paint system with effects
- [x] Matrix-based transformations
- [x] Clipping operations (rect, round rect, path)
- [x] Save/restore state management
- [x] Vertex-based rendering with colors/textures
- [x] Storage management and operation recording
- [x] Canvas abstraction with state tracking
- [x] Renderer interface for multiple backends
- [x] Text system foundation (fonts, styles, paragraphs)
- [x] Image system (CPU/GPU backed)
- [x] Visual effects (filters, blend modes)
- [x] Spatial indexing (R-tree, QuadTree)
- [x] Layer system for advanced compositing
- [x] Utility functions and helpers
- [x] Comprehensive test coverage

### 🚧 In Progress
- [ ] Text layout engine implementation
- [ ] Image loading and format support
- [ ] GPU backend implementations
- [ ] Performance optimizations
- [ ] Advanced effects implementation

### 📋 Planned
- [ ] Animation support and interpolation
- [ ] SVG import/export capabilities
- [ ] Advanced text features (bidirectional, shaping)
- [ ] GPU compute shader effects
- [ ] Multi-threading support
- [ ] Memory optimization
- [ ] Profiling and debugging tools
- [ ] Documentation improvements

## Dependencies

- `github.com/opensraph/sraph/geom`: Geometry types and operations

## Testing

```bash
# Run all tests
go test ./dl/...

# Run with coverage
go test -cover ./dl/...

# Run benchmarks
go test -bench=. ./dl/...

# Run specific test categories
go test -run TestAdvanced ./dl/...
```

## Examples

The package includes comprehensive examples:

- **Basic Usage**: Simple drawing operations and display list creation
- **Complex Scenes**: Multi-layered scenes with various effects
- **Path Construction**: Advanced path building with geometric shapes
- **Text Rendering**: Rich text with mixed styles and layouts
- **Vertex Rendering**: Triangle-based rendering with interpolated attributes
- **Spatial Indexing**: Performance optimization with spatial structures
- **Custom Renderers**: Implementing custom rendering backends

See `examples.go` and `dl_test.go` for detailed usage examples.

## Performance Considerations

### Optimization Tips
1. **Use Spatial Indexing**: For scenes with many objects, use R-tree or QuadTree
2. **Batch Similar Operations**: Group similar drawing calls together
3. **Leverage Bounds Tracking**: Display lists automatically track bounds for culling
4. **Cache Complex Paths**: Reuse path objects for repeated shapes
5. **Use Layers Judiciously**: Layers are powerful but have memory overhead
6. **Profile Your Usage**: Use the included benchmarks to identify bottlenecks

### Memory Management
- Display lists are immutable once built, enabling safe sharing
- Builders can be reset and reused to minimize allocations
- Geometry objects use value semantics where possible
- Large resources (images, complex paths) should be cached

### Threading
- Display lists are thread-safe for reading once built
- Builders should not be shared across threads during recording
- Canvas operations should be performed on a single thread
- Renderers may have specific threading requirements

## Contributing

This implementation is part of the hydrz/sraph project. Contributions are welcome!

### Key Design Principles

1. **Backend Agnostic**: Operations are recorded independent of rendering backend
2. **Immutable**: Display lists are immutable once built for thread safety
3. **Efficient**: Optimized for both recording and playback performance
4. **Extensible**: Easy to add new operations and rendering backends
5. **Compatible**: API design inspired by proven systems (Flutter, Skia)

### Development Guidelines

- Follow Go best practices and conventions
- Use English for all comments and documentation
- Ensure comprehensive test coverage for new features
- Include benchmarks for performance-critical code
- Update documentation for API changes

## License

[Add license information here]
