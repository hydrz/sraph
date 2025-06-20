// Package geom provides a comprehensive 2D/3D geometry library for graphics programming.
//
// This package offers high-performance geometric primitives and mathematical operations
// with support for generic scalar types, designed for graphics programming, game development,
// and computational geometry applications.
//
// # Core Features
//
// The geom package includes:
//
//   - Scalar types: Scalar, F64, I32, I64, Int, I26_6 (fixed-point)
//   - Angles: Radians and Degrees with automatic conversion
//   - Points and Vectors: 2D, 3D, and 4D vectors with comprehensive operations
//   - Geometric shapes: Rectangles, rounded rectangles, ellipses, superellipses
//   - Transformations: Matrices, quaternions, and efficient 2D transforms (RSTransform)
//   - Colors: RGBA colors with color space conversion and 29 blend modes
//   - Gradients: Linear and radial gradients with texture generation
//   - Path system: Flexible path representation with receivers and sources
//   - Stroke styles: Comprehensive stroke parameters (caps, joins, miter limits)
//   - Wang's formula: Optimal curve subdivision for tessellation
//
// # Type Safety
//
// All geometric types are generic over scalar types, ensuring type safety and performance:
//
//	type Point[T Number] struct {
//		X, Y T
//	}
//
//	type Matrix[T Number] [16]T
//
// The Scalar interface constrains numeric types to those suitable for geometry calculations:
//
//	type Scalar interface {
//		~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
//		Float64() float64
//		String() string
//	}
//
// # Basic Usage
//
// Creating and manipulating points:
//
//	p1 := geom.Pt[geom.Scalar](10, 20)
//	p2 := geom.Pt[geom.Scalar](30, 40)
//	distance := p1.Distance(p2)
//	midpoint := p1.Lerp(p2, 0.5)
//
// Working with rectangles:
//
//	rect := geom.NewRectXYWH[geom.Scalar](0, 0, 100, 200)
//	center := rect.Center()
//	area := rect.Area()
//	contains := rect.Contains(p1)
//
// Matrix transformations:
//
//	matrix := geom.NewMatrix[geom.Scalar]()
//	matrix = matrix.Translate(geom.Vector3[geom.Scalar]{X: 50, Y: 25, Z: 0})
//	matrix = matrix.RotateZ(geom.Degrees(45).Radians())
//	matrix = matrix.Scale(geom.Vector3[geom.Scalar]{X: 1.5, Y: 1.5, Z: 1})
//	transformedPoint := matrix.Transform(p1)
//
// # Color Management
//
// Colors support multiple formats and color space conversions:
//
//	color1 := geom.NewColorRGB8(255, 128, 64)      // 8-bit RGB
//	color2 := geom.NewColorHex(0xFF8040)           // Hex
//	color3 := geom.ColorRed()                      // Predefined
//
//	blended := color1.Blend(color2, geom.BlendModeSrcOver)
//	linear := color1.SRGBToLinear()
//
// # Path System
//
// The path system provides flexible path representation:
//
//	rectPath := geom.NewRectPathSource(rect)
//	ellipsePath := geom.NewEllipsePathSource(bounds)
//
//	// Custom path receiver
//	type MyReceiver struct{}
//	func (r *MyReceiver) MoveTo(p geom.Point[geom.Scalar], willBeClosed bool) { ... }
//	func (r *MyReceiver) LineTo(p geom.Point[geom.Scalar]) { ... }
//	// ... implement other PathReceiver methods
//
//	receiver := &MyReceiver{}
//	rectPath.Dispatch(receiver)
//
// # Performance Considerations
//
// The package is designed for high performance:
//
//   - All operations are allocation-free where possible
//   - Matrix operations use column-major storage for GPU compatibility
//   - Wang's formula provides optimal curve tessellation
//   - Geometric queries are optimized for common cases
//   - Type-specific optimizations for differenT Number types
//
// # Coordinate Systems
//
// The package uses standard mathematical coordinate systems:
//
//   - 2D: Origin at bottom-left, Y-axis pointing up (can be configured)
//   - 3D: Right-handed coordinate system
//   - Matrices: Column-major storage (OpenGL/Vulkan compatible)
//   - Angles: Radians for calculations, Degrees for convenience
//
// # Precision and Tolerances
//
// Floating-point comparisons use configurable epsilon values:
//
//	const (
//		Epsilon32 = 1e-3  // For Scalar comparisons
//		Epsilon64 = 1e-6  // For F64 comparisons
//	)
//
// Use NearlyEqual for safe floating-point equality:
//
//	if geom.NearlyEqual(a, b) {
//		// Values are equal within tolerance
//	}
//
// # Integration
//
// The package integrates well with Go's standard library:
//
//	// Convert to/from standard types
//	goRect := rect.ToGo()                          // image.Rectangle
//	goPoint := point.ToGo()                        // image.Point
//	goColor := color.ToRGBA()                      // color.RGBA
//
//	// From standard types
//	rect := geom.NewRectFromGo[geom.Scalar](goRect)
//	point := geom.NewPointFromGo[geom.Scalar](goPoint)
//	color := geom.NewColorFromRGBA(goColor)
//
// For more detailed examples and API documentation, see the individual type documentation.
package geom
