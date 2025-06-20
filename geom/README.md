# Geom - High-Performance 2D/3D Geometry Library for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/opensraph/sraph/geom.svg)](https://pkg.go.dev/github.com/opensraph/sraph/geom)
[![Go Report Card](https://goreportcard.com/badge/github.com/opensraph/sraph/geom)](https://goreportcard.com/report/github.com/opensraph/sraph/geom)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive, high-performance geometry library for Go, providing robust 2D and 3D geometric primitives, transformations, and utilities for graphics programming, game development, and computational geometry applications.

## Features

### Core Geometric Primitives
- **Points & Vectors**: 2D/3D points and vectors with full arithmetic operations
- **Rectangles**: Axis-aligned rectangles with comprehensive manipulation methods
- **Rounded Rectangles**: Rectangles with customizable corner radii
- **Sizes**: Width/height dimensions with scaling and arithmetic operations
- **Colors**: RGBA color manipulation with multiple color spaces and blending modes

### Mathematical Operations
- **Matrix Transformations**: 4x4 matrices for 3D transformations (translation, rotation, scaling, perspective)
- **Quaternions**: Efficient 3D rotation representation with interpolation
- **Angular Measurements**: Type-safe radians and degrees with conversions
- **Fixed-Point Numbers**: Int26_6 for sub-pixel precision graphics

### Advanced Features
- **Path Generation**: Bezier curves, lines, and complex path construction
- **Gradients**: Linear and radial gradients with color interpolation
- **Color Blending**: 29 blend modes including Porter-Duff and separable blend modes
- **Collision Detection**: Point-in-shape testing for rectangles and rounded rectangles
- **Coordinate Systems**: Support for different graphics API conventions (DirectX/Vulkan, OpenGL)

## Installation

```bash
go get github.com/opensraph/sraph/geom
```

## Quick Start

### Basic Geometry Operations

```go
package main

import (
    "fmt"
    "github.com/opensraph/sraph/geom"
)

func main() {
    // Create points and perform operations
    p1 := geom.Point[float64]{X: 10, Y: 20}
    p2 := geom.Point[float64]{X: 30, Y: 40}

    // Vector arithmetic
    sum := p1.Add(p2)
    distance := p1.Distance(p2)

    fmt.Printf("Sum: %v, Distance: %.2f\n", sum, distance)

    // Create and manipulate rectangles
    rect := geom.NewRect[float64](0, 0, 100, 80)
    center := rect.Center()
    area := rect.Area()

    fmt.Printf("Center: %v, Area: %.2f\n", center, area)

    // Check containment
    point := geom.Point[float64]{X: 50, Y: 40}
    contains := rect.Contains(point)
    fmt.Printf("Rectangle contains point: %v\n", contains)
}
```

### Matrix Transformations

```go
// Create transformation matrices
transform := geom.NewMatrix[float64]()

// Apply transformations in sequence
result := transform.
    Translate(geom.Vector3[float64]{X: 10, Y: 20, Z: 0}).
    Scale(geom.Vector3[float64]{X: 2, Y: 2, Z: 1}).
    RotateZ(geom.Radians(math.Pi / 4)) // 45 degrees

// Transform points
point := geom.Point[float64]{X: 1, Y: 1}
transformed := result.Transform(point).(geom.Point[float64])

// Check matrix properties
isInvertible := result.IsInvertible()
determinant := result.Determinant()
```

### Color Operations

```go
// Create colors
red := geom.ColorRed()
blue := geom.ColorBlue()
customColor := geom.NewColorRGBA8(128, 255, 64, 255)

// Color blending
blended := red.Blend(blue, geom.BlendModeMultiply)

// Color space conversions
linear := red.SRGBToLinear()
srgb := linear.LinearToSRGB()

// Color interpolation
interpolated := red.Lerp(blue, 0.5) // 50% between red and blue
```

### Rounded Rectangles and Paths

```go
// Create rounded rectangle
rect := geom.NewRect[float64](0, 0, 100, 60)
roundRect := geom.NewRoundRectRadius(rect, 10) // 10-pixel corner radius

// Check if point is inside rounded rectangle
point := geom.Point[float64]{X: 95, Y: 5}
inside := roundRect.Contains(point)

// Generate path data
type MyPathReceiver struct {
    commands []string
}

func (r *MyPathReceiver) MoveTo(p geom.Point[float64], willBeClosed bool) {
    r.commands = append(r.commands, fmt.Sprintf("MoveTo(%.2f, %.2f)", p.X, p.Y))
}
// ... implement other PathReceiver methods

receiver := &MyPathReceiver{}
roundRect.Dispatch(receiver, true)
```

### Gradients

```go
// Create linear gradient
stops := []geom.GradientStop{
    {Color: geom.ColorRed(), Position: 0.0},
    {Color: geom.ColorGreen(), Position: 0.5},
    {Color: geom.ColorBlue(), Position: 1.0},
}

linearGradient := geom.NewLinearGradient(stops)
gradientData := linearGradient.ToBuffer()

// Create radial gradient
center := geom.Point[float64]{X: 50, Y: 50}
radialGradient := geom.NewRadialGradient(center, 30, stops)
```

## API Reference

### Core Types

#### Scalar Types
- `Scalar` - Primary floating-point type (alias for float64)
- `TScalar` - Generic constraint for numeric types
- `Radians` / `Degrees` - Type-safe angular measurements
- `Int26_6` - Fixed-point number for sub-pixel precision

#### Geometric Primitives
- `Point[T]` - 2D point with X, Y coordinates
- `Vector2[T]` / `Vector3[T]` / `Vector4[T]` - N-dimensional vectors
- `Size[T]` - Width and height dimensions
- `Rect[T]` - Axis-aligned rectangle
- `RoundRect[T]` - Rectangle with rounded corners
- `Quad[T]` - Array of 4 points representing a quadrilateral

#### Transformations
- `Matrix[T]` - 4x4 transformation matrix
- `Quaternion[T]` - Quaternion for 3D rotations
- `RSTransform[T]` - Rotation, scale, and translation transform

#### Colors and Graphics
- `Color` - RGBA color with floating-point components
- `BlendMode` - Enumeration of blending modes
- `ColorMatrix` - 4x5 matrix for color transformations

### Key Methods

#### Point Operations
```go
// Arithmetic
p1.Add(p2)          // Vector addition
p1.Sub(p2)          // Vector subtraction
p1.Scale(factor)    // Scalar multiplication
p1.Dot(p2)          // Dot product
p1.Cross(p2)        // 2D cross product

// Geometric
p1.Distance(p2)     // Euclidean distance
p1.Length()         // Vector magnitude
p1.Normalize()      // Unit vector
p1.Rotate(angle)    // Rotation around origin
p1.Lerp(p2, t)      // Linear interpolation
```

#### Rectangle Operations
```go
// Construction
geom.NewRect(left, top, right, bottom)
geom.NewRectXYWH(x, y, width, height)
geom.BoundingRect(points...)

// Properties
rect.Width()        // Rectangle width
rect.Height()       // Rectangle height
rect.Area()         // Rectangle area
rect.Center()       // Center point
rect.IsEmpty()      // Check if empty
rect.IsSquare()     // Check if square

// Transformations
rect.Translate(vector)    // Move rectangle
rect.Scale(factor)        // Scale rectangle
rect.Expand(amount)       // Expand by amount
rect.Union(other)         // Union with another rectangle
rect.Intersect(other)     // Intersection with another rectangle

// Testing
rect.Contains(point)      // Point containment
rect.ContainsRect(other)  // Rectangle containment
rect.Intersects(other)    // Intersection test
```

#### Matrix Operations
```go
// Construction
geom.NewMatrix[T]()  // Identity matrix

// Transformations
matrix.Translate(vector)     // Add translation
matrix.Scale(vector)         // Add scaling
matrix.Rotate(angle, axis)   // Add rotation
matrix.RotateX/Y/Z(angle)   // Axis-specific rotation

// Properties
matrix.IsIdentity()      // Check if identity
matrix.IsInvertible()    // Check if invertible
matrix.Determinant()     // Matrix determinant
matrix.Transpose()       // Matrix transpose
matrix.Invert()          // Matrix inverse

// Application
matrix.Transform(geometry)           // Transform geometric objects
matrix.TransformDirection(vector)    // Transform direction vectors
```

#### Color Operations
```go
// Construction
geom.NewColor(r, g, b, a)
geom.NewColorRGB8(r, g, b)
geom.NewColorHex(0xFF0000)  // Red

// Arithmetic
color1.Add(color2)
color1.Mul(color2)
color1.Scale(factor)
color1.Lerp(color2, t)

// Blending
color1.Blend(color2, geom.BlendModeMultiply)

// Color Space
color.SRGBToLinear()
color.LinearToSRGB()
color.Premultiply()
color.Unpremultiply()
```

## Performance Considerations

### Memory Efficiency
- All geometric types are value types (structs) for minimal heap allocation
- Methods return new values rather than modifying in place (immutable design)
- Generic types allow choosing appropriate numeric precision

### Computational Efficiency
- Optimized matrix operations with specialized methods for common transformations
- Fast path detection for axis-aligned and identity transformations
- Efficient collision detection algorithms

### Numerical Stability
- Configurable epsilon values for floating-point comparisons
- Robust handling of edge cases (zero vectors, singular matrices)
- Fixed-point arithmetic support for graphics applications

## Use Cases

### Graphics Programming
- 2D and 3D rendering pipelines
- UI layout and animation systems
- Canvas and vector graphics libraries
- Texture mapping and coordinate transformations

### Game Development
- Physics simulations and collision detection
- Camera systems and view transformations
- Sprite positioning and animation
- Particle systems

### CAD and Engineering
- Geometric modeling and analysis
- Technical drawing applications
- Measurement and annotation tools
- Coordinate system conversions

### Image Processing
- Color space conversions and corrections
- Geometric image transformations
- Gradient generation and blending
- Alpha compositing operations

## Coordinate System Conventions

This library follows the **DirectX/Vulkan convention** for 3D graphics:
- **Left-handed coordinate system**: X right, Y up, Z away from viewer
- **Positive rotation**: Clockwise about the rotation axis
- **NDC bounds**: x,y ∈ [-1,1], z ∈ [0,1] where 0 is near plane, 1 is far plane
- **NDC origin**: (0, 0, 0.5) representing the center of the NDC cube

For OpenGL compatibility, appropriate conversion utilities are provided.

## Testing and Quality Assurance

- Comprehensive unit tests covering all public APIs
- Property-based testing for mathematical operations
- Benchmark tests for performance-critical operations
- Fuzzing tests for numerical stability
- Cross-platform compatibility testing

## Contributing

Contributions are welcome! Please ensure that:
- All new features include comprehensive tests
- Code follows Go formatting standards (`gofmt`)
- Documentation is updated for new public APIs
- Performance implications are considered and benchmarked

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by geometry libraries from Skia, Flutter, and other graphics frameworks
- Mathematical foundations based on computer graphics and computational geometry literature
- Color science implementations following CIE and W3C standards
