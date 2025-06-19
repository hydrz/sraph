// package display provides a display list system for 2D graphics rendering.
//
// The display list system allows you to record a sequence of drawing operations
// and replay them efficiently on different rendering backends. This is inspired
// by Flutter's display list architecture but implemented in Go.
//
// # Core Concepts
//
// The display list system consists of several key components:
//
//   - DisplayList: An immutable sequence of recorded drawing operations
//   - DisplayListBuilder: A builder for recording operations into a display list
//   - OpReceiver: Interface for objects that can receive and execute operations
//   - Canvas: A concrete implementation of OpReceiver for immediate rendering
//   - Renderer: Interface for different rendering backends (GPU, CPU, etc.)
//
// # Basic Usage
//
//	// Create a builder to record operations
//	builder := dl.NewDisplayListBuilder()
//
//	// Create paint objects
//	paint := dl.NewPaint()
//	paint.SetColor(dl.ColorRed)
//
//	// Record drawing operations
//	builder.DrawRect(geom.NewRect[dl.Scalar](10, 10, 100, 100), paint)
//	builder.Translate(50, 50)
//	builder.DrawCircle(geom.Point[dl.Scalar]{X: 0, Y: 0}, 25, paint)
//
//	// Build the display list
//	displayList := builder.Build()
//
//	// Render the display list
//	canvas := dl.NewCanvas(geom.NewRect[dl.Scalar](0, 0, 200, 200))
//	displayList.Dispatch(canvas)
//
// # Advanced Features
//
// The display list system supports:
//
//   - Complex transformations (translate, scale, rotate, skew, perspective)
//   - Clipping (rectangle, rounded rectangle, path-based)
//   - Various drawing primitives (lines, rectangles, circles, paths, vertices)
//   - Paint effects (colors, blend modes, anti-aliasing)
//   - Save/restore state management
//   - Bounds tracking and culling
//   - Path construction with bezier curves
//
// # Performance Considerations
//
// Display lists are designed for efficiency:
//
//   - Operations are recorded once and can be replayed multiple times
//   - Bounds information allows for culling operations outside the viewport
//   - Attribute flags help optimize rendering pipeline decisions
//   - Immutable design enables safe sharing across threads
//
// # Integration with Graphics Backends
//
// The system is designed to work with different rendering backends:
//
//   - Implement the Renderer interface for your graphics API
//   - Use RenderContext to manage state and dispatch operations
//   - Operations are backend-agnostic and can target GPU or CPU renderers
//
// For more examples and detailed usage, see the dl_examples.go file.
package display

// TODO: Future features to implement
//
// The following features are planned for future versions:
//
//   - Text rendering support (fonts, text layout, text effects)
//   - Image rendering support (textures, image filters, nine-patch)
//   - Advanced effects (shadows, gradients, mask filters)
//   - Animation support (interpolation, easing functions)
//   - Spatial indexing (R-tree for efficient culling)
//   - Serialization (save/load display lists)
//   - Multi-threading support (parallel rendering)
//   - GPU-accelerated backends (Vulkan, Metal, DirectX)
//   - Vector graphics import/export (SVG, PDF)
//   - Performance profiling tools
//   - Memory pool optimizations
//   - Tessellation improvements
//   - Advanced clipping modes
//   - Layer effects and compositing
