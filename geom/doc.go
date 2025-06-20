// Package geom provides comprehensive 2D and 3D geometric primitives, transformations,
// and utilities for graphics programming, game development, and computational geometry.
//
// # Overview
//
// This package offers a complete set of geometric types and operations optimized for
// performance and numerical stability. All geometric types are immutable value types
// that support generic numeric types through the TScalar constraint.
//
// # Core Types
//
// The package is built around several fundamental geometric primitives:
//
//   - Point[T]: 2D points with X, Y coordinates
//   - Vector2[T], Vector3[T], Vector4[T]: N-dimensional vectors
//   - Size[T]: Width and height dimensions
//   - Rect[T]: Axis-aligned rectangles
//   - RoundRect[T]: Rectangles with rounded corners
//   - Matrix[T]: 4x4 transformation matrices
//   - Quaternion[T]: 3D rotation representation
//   - Color: RGBA colors with floating-point precision
//
// # Coordinate System
//
// The library follows DirectX/Vulkan conventions:
//   - Left-handed coordinate system (X right, Y up, Z away from viewer)
//   - Positive rotation is clockwise about the rotation axis
//   - Normalized Device Coordinates (NDC): x,y ∈ [-1,1], z ∈ [0,1]
//   - NDC origin at (0, 0, 0.5) representing the center of the cube
//
// # Basic Usage
//
// Creating and manipulating points:
//
//	p1 := geom.NewPoint(10.0, 20.0)
//	p2 := geom.NewPoint(30.0, 40.0)
//
//	// Vector operations
//	sum := p1.Add(p2)                    // Point addition
//	distance := p1.Distance(p2)          // Euclidean distance
//	normalized := p1.Normalize()         // Unit vector
//	rotated := p1.Rotate(geom.Radians(math.Pi / 4))  // 45° rotation
//
// Working with rectangles:
//
//	rect := geom.NewRect(0.0, 0.0, 100.0, 80.0)
//	center := rect.Center()              // Center point
//	area := rect.Area()                  // Area calculation
//	expanded := rect.Expand(10.0)        // Expand by 10 units
//
//	// Containment testing
//	point := geom.NewPoint(50.0, 40.0)
//	contains := rect.Contains(point)     // Point-in-rectangle test
//
// # Matrix Transformations
//
// The Matrix type provides comprehensive 3D transformation capabilities:
//
//	// Create identity matrix and apply transformations
//	transform := geom.NewMatrix[float64]()
//	result := transform.
//		Translate(geom.Vector3[float64]{X: 10, Y: 20, Z: 0}).
//		Scale(geom.Vector3[float64]{X: 2, Y: 2, Z: 1}).
//		RotateZ(geom.Radians(math.Pi / 4))
//
//	// Transform geometric objects
//	transformedPoint := result.Transform(point)
//
//	// Matrix analysis
//	isInvertible := result.IsInvertible()
//	determinant := result.Determinant()
//	inverse := result.Invert()
//
// # Color Operations
//
// The Color type supports various color spaces and blending modes:
//
//	// Create colors
//	red := geom.ColorRed()
//	custom := geom.NewColorRGBA8(128, 255, 64, 255)
//	fromHex := geom.NewColorHex(0xFF0000)
//
//	// Color arithmetic
//	blended := red.Blend(blue, geom.BlendModeMultiply)
//	interpolated := red.Lerp(blue, 0.5)  // 50% between colors
//
//	// Color space conversions
//	linear := red.SRGBToLinear()
//	premultiplied := red.Premultiply()
//
// # Advanced Features
//
// Rounded rectangles with custom corner radii:
//
//	rect := geom.NewRect(0.0, 0.0, 100.0, 60.0)
//	roundRect := geom.NewRoundRectRadius(rect, 10.0)  // 10px radius
//
//	// Complex corner configurations
//	radii := geom.NewRoundingRadiiLTRB(5, 10, 15, 20)  // Individual corners
//	customRound := geom.NewRoundRect(rect, radii)
//
//	// Containment test accounts for rounded corners
//	inside := roundRect.Contains(point)
//
// Gradients for smooth color transitions:
//
//	stops := []geom.GradientStop{
//		{Color: geom.ColorRed(), Position: 0.0},
//		{Color: geom.ColorGreen(), Position: 0.5},
//		{Color: geom.ColorBlue(), Position: 1.0},
//	}
//
//	linearGradient := geom.NewLinearGradient(stops)
//	gradientData := linearGradient.ToBuffer()  // For GPU upload
//
//	center := geom.NewPoint(50.0, 50.0)
//	radialGradient := geom.NewRadialGradient(center, 30.0, stops)
//
// Quaternion rotations for smooth 3D animations:
//
//	axis := geom.Vector3[float64]{X: 0, Y: 0, Z: 1}
//	angle := geom.Radians(math.Pi / 2)  // 90 degrees
//	quat := geom.NewQuaternionFromAxisAngle(axis, angle)
//
//	// Interpolation between rotations
//	interpolated := quat1.Slerp(quat2, 0.5)
//
//	// Rotate vectors
//	vector := geom.Vector3[float64]{X: 1, Y: 0, Z: 0}
//	rotated := quat.RotateVector3(vector)
//
// # Performance Considerations
//
// All geometric types are designed for high performance:
//   - Value types minimize heap allocations
//   - Immutable design enables safe concurrent access
//   - Generic types allow choosing appropriate numeric precision
//   - Optimized operations for common transformations
//   - SIMD-friendly data layouts where possible
//
// # Numerical Stability
//
// The package provides robust handling of floating-point precision:
//   - Configurable epsilon values for comparisons (Epsilon32, Epsilon64)
//   - NearlyEqual function for tolerance-based equality
//   - IsFinite checks for numerical validity
//   - Stable algorithms for matrix operations
//
// Fixed-point arithmetic for sub-pixel precision:
//
//	fixed := geom.Int26_6(64)  // Represents 1.0 in 26.6 format
//	floatVal := fixed.Float64()  // Convert to float64
//
// # Type Safety
//
// Angular measurements use distinct types to prevent unit confusion:
//
//	degrees := geom.Degrees(45.0)
//	radians := degrees.Radians()        // Explicit conversion
//	rotated := point.Rotate(radians)    // Type-safe API
//
// # Error Handling
//
// Most operations are designed to be infallible, returning sensible defaults:
//   - Zero-length vector normalization returns (1, 0)
//   - Non-invertible matrix inversion returns zero matrix
//   - Division by zero in colors returns transparent black
//
// For validation, use the provided checking functions:
//
//	if !matrix.IsInvertible() {
//		// Handle non-invertible matrix
//	}
//
//	if !color.IsFinite() {
//		// Handle invalid color values
//	}
//
// # Integration
//
// The package integrates well with standard Go libraries:
//
//	// Convert to/from image package types
//	goPoint := point.Go()              // -> image.Point
//	goRect := rect.Go()                // -> image.Rectangle
//	goColor := color.ToRGBA()          // -> color.RGBA
//
//	// Create from Go types
//	pointFromGo := geom.NewPointGo(goPoint)
//	rectFromGo := geom.NewRectGo(goRect)
//	colorFromGo := geom.NewColorGo(goColor)
//
// # Common Patterns
//
// Building transformation hierarchies:
//
//	parentTransform := geom.NewMatrix[float64]().
//		Translate(geom.Vector3[float64]{X: 100, Y: 50, Z: 0}).
//		RotateZ(geom.Radians(math.Pi / 4))
//
//	childLocalTransform := geom.NewMatrix[float64]().
//		Scale(geom.Vector3[float64]{X: 0.5, Y: 0.5, Z: 1})
//
//	// Child world transform = parent * child_local
//	childWorldTransform := parentTransform.Mul(childLocalTransform)
//
// Animation and interpolation:
//
//	// Linear interpolation between positions
//	start := geom.NewPoint(0.0, 0.0)
//	end := geom.NewPoint(100.0, 100.0)
//	t := 0.5  // 50% through animation
//	current := start.Lerp(end, t)
//
//	// Color fading
//	visible := geom.ColorWhite()
//	transparent := visible.WithAlpha(0.0)
//	faded := visible.Lerp(transparent, t)
//
// Viewport and projection transformations:
//
//	// Orthographic projection for 2D rendering
//	screenSize := geom.Size[float64]{Width: 800, Height: 600}
//	orthoMatrix := geom.NewMatrix[float64]().Orthographic(screenSize)
//
//	// Perspective projection for 3D rendering
//	fov := geom.Radians(math.Pi / 3)  // 60 degrees
//	aspectRatio := screenSize.Width / screenSize.Height
//	perspMatrix := geom.NewMatrix[float64]().
//		Perspective(fov, aspectRatio, 0.1, 1000.0)
//
// This package provides the foundation for graphics applications, game engines,
// CAD software, and any application requiring robust geometric computations.
package geom
