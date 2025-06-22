# Geom - High-Performance 2D/3D Geometry Library for Go

A comprehensive, high-performance geometry library for Go that provides essential mathematical primitives for 2D and 3D graphics, game development, and computational geometry applications.

## Features

### Core Types
- **Point/Vector2**: 2D points and vectors with comprehensive arithmetic operations
- **Vector3/Vector4**: 3D and 4D vectors for spatial calculations
- **Size**: 2D dimensions with scaling and mathematical operations
- **Rect**: Axis-aligned rectangles with intersection, union, and transformation operations
- **RoundRect**: Rectangles with configurable corner radii
- **Matrix**: 4×4 transformation matrices with full 3D transformation support
- **Quaternion**: Efficient 3D rotation representation

### Advanced Geometry
- **Color**: RGBA color representation with blending modes and color space conversions
- **Gradient**: Linear and radial gradients with customizable color stops
- **Path**: Vector path operations and rendering primitives
- **Superellipse**: Smooth rounded rectangles with advanced corner handling
- **Stroke**: Configurable stroke styles for path rendering

### Mathematical Operations
- **Transformations**: Translation, rotation, scaling, shearing, and perspective projection
- **Interpolation**: Linear interpolation (lerp) for smooth animations
- **Geometric Queries**: Intersection testing, containment checks, distance calculations
- **Color Blending**: 29 different blend modes including Porter-Duff compositing
- **Curve Subdivision**: Wang's formula for optimal curve tessellation

### Performance Features
- **Generic Types**: Type-safe operations with Go 1.18+ generics
- **Memory Efficient**: Minimal allocations with value-based types
- **SIMD-Ready**: Designed for future SIMD optimizations
- **Numerical Stability**: Robust floating-point operations with epsilon comparisons

## Installation

```bash
go get github.com/opensraph/sraph/geom
```

## Quick Start

### Basic 2D Operations

```go
package main

import (
    "fmt"
    "github.com/opensraph/sraph/geom"
)

func main() {
    // Create points and perform vector arithmetic
    p1 := geom.NewPoint[float64](10, 20)
    p2 := geom.NewPoint[float64](30, 40)

    sum := p1.Add(p2)
    distance := p1.Distance(p2)

    fmt.Printf("Sum: %v, Distance: %.2f\n", sum, distance)

    // Rectangle operations
    rect := geom.NewRectXYWH[float64](0, 0, 100, 50)
    center := rect.Center()
    area := rect.Area()

    fmt.Printf("Center: %v, Area: %.2f\n", center, area)

    // Check containment
    point := geom.NewPoint[float64](25, 25)
    contains := rect.Contains(point)
    fmt.Printf("Rectangle contains point: %v\n", contains)
}
```

### 3D Transformations

```go
package main

import (
    "math"
    "github.com/opensraph/sraph/geom"
)

func main() {
    // Create a transformation matrix
    matrix := geom.NewMatrix().
        Translate(geom.NewVector3[float64](10, 0, 0)).
        RotateY(geom.Radians(math.Pi / 4)).
        Scale(geom.NewVector3[float64](2, 2, 2))

    // Transform a 3D point
    point := geom.NewVector3[float64](1, 1, 1)
    transformed := point.Transform(matrix)

    fmt.Printf("Original: %v\n", point)
    fmt.Printf("Transformed: %v\n", transformed)

    // Quaternion rotation
    axis := geom.NewVector3[float64](0, 1, 0)
    quat := geom.NewQuaternionFromAxisAngle(axis, geom.Radians(math.Pi/2))
    rotated := quat.RotateVector3(point)

    fmt.Printf("Quaternion rotated: %v\n", rotated)
}
```

### Color and Gradients

```go
package main

import (
    "github.com/opensraph/sraph/geom"
)

func main() {
    // Create colors
    red := geom.ColorRed()
    blue := geom.ColorBlue()

    // Blend colors with different modes
    blended := red.Blend(blue, geom.BlendModeMultiply)
    fmt.Printf("Blended color: %v\n", blended)

    // Create a gradient
    stops := []geom.GradientStop{
        {Color: geom.ColorRed(), Position: 0.0},
        {Color: geom.ColorYellow(), Position: 0.5},
        {Color: geom.ColorBlue(), Position: 1.0},
    }

    gradient := geom.NewLinearGradient(stops)
    buffer := gradient.ToBuffer()

    fmt.Printf("Gradient texture size: %d\n", buffer.TextureSize)
}
```

### Rounded Rectangles

```go
package main

import (
    "github.com/opensraph/sraph/geom"
)

func main() {
    // Create a rounded rectangle
    rect := geom.NewRect[float64](0, 0, 100, 60)
    radii := geom.NewRoundingRadii[float64](10) // 10px radius on all corners
    roundRect := geom.NewRoundRect(rect, radii)

    // Check if point is inside rounded rectangle
    point := geom.NewPoint[float64](50, 30)
    inside := roundRect.Contains(point)

    fmt.Printf("Point inside rounded rect: %v\n", inside)

    // Create rounded rect with different corner radii
    customRadii := geom.NewRoundingRadiiLTRB[float64](5, 10, 15, 20)
    customRoundRect := geom.NewRoundRect(rect, customRadii)

    fmt.Printf("Custom rounded rect bounds: %v\n", customRoundRect.Bounds())
}
```

## API Documentation

### Core Types

#### Point/Vector2
2D points and vectors with arithmetic operations, transformations, and geometric queries.

**Key Methods:**
- `Add(Point) Point` - Vector addition
- `Distance(Point) Scalar` - Euclidean distance
- `Normalize() Point` - Unit vector
- `Transform(Matrix) Point` - Matrix transformation
- `Rotate(Radians) Point` - Rotation around origin

#### Rect
Axis-aligned rectangles with comprehensive operations.

**Key Methods:**
- `Contains(Point) bool` - Point containment test
- `Intersect(Rect) Rect` - Rectangle intersection
- `Union(Rect) Rect` - Rectangle union
- `Transform(Matrix) Quad` - Matrix transformation
- `Scale(Scalar) Rect` - Uniform scaling

#### Matrix
4×4 transformation matrices for 3D graphics.

**Key Methods:**
- `Translate(Vector3) Matrix` - Translation transformation
- `Scale(Vector3) Matrix` - Scaling transformation
- `RotateX/Y/Z(Radians) Matrix` - Axis rotations
- `Perspective(Radians, Scalar, Scalar, Scalar) Matrix` - Perspective projection
- `Mul(Matrix) Matrix` - Matrix multiplication

#### Color
RGBA color with advanced blending and color space operations.

**Key Methods:**
- `Blend(Color, BlendMode) Color` - Color blending
- `Lerp(Color, Scalar) Color` - Linear interpolation
- `Premultiply() Color` - Alpha premultiplication
- `SRGBToLinear() Color` - Color space conversion

### Supported Blend Modes

The library supports 29 different blend modes:

- **Porter-Duff**: Clear, Source, Destination, SrcOver, DstOver, SrcIn, DstIn, SrcOut, DstOut, SrcATop, DstATop, Xor
- **Mathematical**: Plus, Modulate, Screen, Multiply
- **Component**: Darken, Lighten, ColorDodge, ColorBurn, HardLight, SoftLight, Difference, Exclusion
- **HSL**: Hue, Saturation, Color, Luminosity
- **Advanced**: Overlay

## Performance Considerations

- **Value Types**: All types are designed as values to minimize heap allocations
- **Generic Implementation**: Type-safe operations without boxing overhead
- **Numerical Stability**: Uses epsilon comparisons for floating-point operations
- **Memory Layout**: Optimized for cache-friendly access patterns
- **Minimal Dependencies**: Only depends on Go standard library and golang.org/x/image

## Coordinate Systems

### 2D Coordinates
- **Origin**: Top-left corner (0, 0)
- **X-axis**: Points right (positive direction)
- **Y-axis**: Points down (positive direction)
- **Rotation**: Clockwise positive

### 3D Coordinates (DirectX/Vulkan Convention)
- **X-axis**: Points right
- **Y-axis**: Points up
- **Z-axis**: Points away from viewer (into screen)
- **Rotation**: Clockwise about rotation axis
- **NDC**: x,y ∈ [-1,1], z ∈ [0,1]

## Examples

See the [examples](example/) directory for comprehensive usage examples:

- **Color Blending**: Demonstrates all blend modes with visual test patterns
- **Transformations**: 3D matrix operations and quaternion rotations
- **Geometric Queries**: Intersection testing and containment checks

## Requirements

- Go 1.21 or later (uses built-in `min`/`max` functions)
- No external dependencies

## Contributing

Contributions are welcome! Please ensure:

1. **Code Quality**: Follow Go best practices and maintain existing code style
2. **Testing**: Add comprehensive tests for new features
3. **Documentation**: Update documentation and examples as needed
4. **Performance**: Consider performance implications of changes

---

**Note**: This library prioritizes correctness, performance, and ease of use. It's designed to be a solid foundation for graphics and computational geometry applications in Go.
