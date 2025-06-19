# Sraph Geometry Package

A comprehensive, high-performance 2D/3D geometry library for Go, designed for graphics programming, game development, and computational geometry applications.

## Overview

The `geom` package provides a complete set of geometric primitives and mathematical operations with support for generic scalar types. It's built with performance and type safety in mind, offering both integer and floating-point arithmetic with configurable precision.

## Features

### Core Types
- **Scalar Types**: `F32`, `F64`, `I32`, `I64`, `Int`, `I26_6` (fixed-point)
- **Angles**: `Radians`, `Degrees` with automatic conversion
- **Points & Vectors**: 2D/3D/4D vectors with comprehensive operations
- **Geometric Shapes**: Rectangles, rounded rectangles, ellipses, superellipses
- **Transformations**: Matrices, quaternions, RSTransform for efficient 2D transforms

### Advanced Features
- **Color Management**: RGBA colors with color space conversion (sRGB ↔ Linear)
- **Gradients**: Linear and radial gradients with texture generation
- **Path System**: Flexible path representation with receivers and sources
- **Stroke Styles**: Comprehensive stroke parameters (caps, joins, miter limits)
- **Blend Modes**: 29 blend modes including Porter-Duff and advanced modes
- **Wang's Formula**: Curve subdivision for optimal tessellation

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/opensraph/sraph/geom"
)

func main() {
    // Create points and perform operations
    p1 := geom.Pt[geom.F32](10, 20)
    p2 := geom.Pt[geom.F32](30, 40)

    distance := p1.Distance(p2)
    midpoint := p1.Lerp(p2, 0.5)

    fmt.Printf("Distance: %.2f\n", distance)
    fmt.Printf("Midpoint: %s\n", midpoint)

    // Create and manipulate rectangles
    rect := geom.NewRectXYWH[geom.F32](0, 0, 100, 200)
    center := rect.Center()
    area := rect.Area()

    // Create rounded rectangle
    roundRect := geom.NewRoundRectRadius(rect, 10)

    // Matrix transformations
    transform := geom.NewMatrix[geom.F32]()
    transform = transform.Translate(geom.Vector2[geom.F32]{X: 50, Y: 25})
    transform = transform.RotateZ(geom.Radians(0.5))
    transform = transform.Scale(geom.Vector2[geom.F32]{X: 1.5, Y: 1.5})

    transformedRect := rect.TransformBounds(transform)
}
```

## Core Components

### Scalar Types and Precision

```go
// Different scalar types for different use cases
type F32 float32  // 32-bit floating point
type F64 float64  // 64-bit floating point
type I32 int32    // 32-bit signed integer
type I26_6 int32  // 26.6 fixed-point for precise typography

// Configurable epsilon for floating-point comparisons
const (
    Epsilon32 = 1e-3  // For F32
    Epsilon64 = 1e-6  // For F64
)

// Safe equality comparison
if geom.ScalarEq(a, b) {
    // Values are equal within tolerance
}
```

### Points and Vectors

```go
// 2D operations
point := geom.Pt[geom.F32](10, 20)
vector := geom.Vector2[geom.F32]{X: 5, Y: 10}

// Vector operations
length := point.Length()
normalized := point.Normalize()
dotProduct := point.Dot(vector)
crossProduct := point.Cross(vector)

// 3D vectors
vec3 := geom.Vector3[geom.F32]{X: 1, Y: 2, Z: 3}
cross3D := vec3.Cross(geom.Vector3[geom.F32]{X: 4, Y: 5, Z: 6})

// 4D vectors for homogeneous coordinates
vec4 := geom.Vector4[geom.F32]{X: 1, Y: 2, Z: 3, W: 1}
```

### Rectangles and Shapes

```go
// Create rectangles in different ways
rect1 := geom.NewRect[geom.F32](0, 0, 100, 200)           // LTRB
rect2 := geom.NewRectXYWH[geom.F32](10, 20, 80, 60)       // XYWH
rect3 := geom.NewRectOriginSize(origin, size)              // Origin + Size

// Rectangle operations
union := rect1.Union(rect2)
intersection := rect1.Intersect(rect2)
contains := rect1.Contains(point)
overlaps := rect1.Intersects(rect2)

// Rounded rectangles with different corner radii
roundRect := geom.NewRoundRectLTRB(rect, 5, 10, 15, 20)
isOval := roundRect.IsOval()

// Superellipses for smooth, organic shapes
superellipse := geom.NewSuperellipseRadius(rect, 15)
```

### Matrix Transformations

```go
// Create transformation matrices
matrix := geom.NewMatrix[geom.F32]()

// Apply transformations (operations are chainable)
matrix = matrix.Translate(geom.Vector3[geom.F32]{X: 100, Y: 50, Z: 0})
matrix = matrix.RotateZ(geom.Degrees(45).Radians())
matrix = matrix.Scale(geom.Vector3[geom.F32]{X: 2, Y: 2, Z: 1})

// Check matrix properties
isIdentity := matrix.IsIdentity()
isInvertible := matrix.IsInvertible()
hasTranslation := matrix.HasTranslation()
isAxisAligned := matrix.IsAxisAligned()

// Transform geometry
transformedPoint := matrix.Transform(point)
transformedRect := rect.TransformBounds(matrix)

// Matrix decomposition
decomp := matrix.Decompose()
translation := decomp.Translation
rotation := decomp.Rotation
scale := decomp.Scale
```

### Colors and Gradients

```go
// Create colors in various formats
color1 := geom.NewColorRGB8(255, 128, 64)                    // 8-bit RGB
color2 := geom.NewColorHex(0xFF8040)                         // Hex
color3 := geom.NewColor[geom.F32](1.0, 0.5, 0.25, 1.0)     // Float RGBA
color4 := geom.ColorRed()                                    // Predefined colors

// Color operations
blended := color1.Blend(color2, geom.BlendModeSrcOver)
interpolated := color1.Lerp(color2, 0.5)
premultiplied := color1.Premultiply()

// Color space conversion
linear := color1.SRGBToLinear()
srgb := linear.LinearToSRGB()

// Create gradients
stops := []geom.GradientStop[geom.F32]{
    {Color: geom.ColorRed(), Position: 0.0},
    {Color: geom.ColorBlue(), Position: 1.0},
}
gradient := geom.NewLinearGradient(stops)
gradientData := gradient.ToBuffer()
```

### Path System

```go
// Create path sources for different shapes
rectPath := geom.NewRectPathSource(rect)
ellipsePath := geom.NewEllipsePathSource(bounds)
roundRectPath := geom.NewRoundRectPathSource(roundRect)

// Custom path receiver
type MyPathReceiver struct{}

func (r *MyPathReceiver) MoveTo(p geom.Point[geom.F32], willBeClosed bool) {
    // Handle move to operation
}

func (r *MyPathReceiver) LineTo(p geom.Point[geom.F32]) {
    // Handle line to operation
}

func (r *MyPathReceiver) QuadTo(cp, p geom.Point[geom.F32]) {
    // Handle quadratic curve
}

func (r *MyPathReceiver) CubicTo(cp1, cp2, p geom.Point[geom.F32]) {
    // Handle cubic curve
}

func (r *MyPathReceiver) Close() {
    // Handle path close
}

func (r *MyPathReceiver) PathEnd() {
    // Handle path end
}

// Dispatch path to receiver
receiver := &MyPathReceiver{}
rectPath.Dispatch(receiver)
```

### Quaternions for 3D Rotations

```go
// Create quaternions
axis := geom.Vector3[geom.F32]{X: 0, Y: 1, Z: 0}
angle := geom.Degrees(90).Radians()
quat := geom.NewQuaternionFromAxisAngle(axis, angle)

// Quaternion operations
normalized := quat.Normalize()
inverted := quat.Invert()
interpolated := quat.Slerp(otherQuat, 0.5)

// Rotate vectors
vector := geom.Vector3[geom.F32]{X: 1, Y: 0, Z: 0}
rotated := quat.RotateVector3(vector)

// Convert to matrix
rotationMatrix := geom.NewMatrix[geom.F32]().RotateQuat(quat)
```

## Blend Modes

The package supports 29 different blend modes for color composition:

```go
// Porter-Duff modes
geom.BlendModeClear
geom.BlendModeSrc
geom.BlendModeDst
geom.BlendModeSrcOver  // Default
geom.BlendModeDstOver
geom.BlendModeSrcIn
geom.BlendModeDstIn
geom.BlendModeSrcOut
geom.BlendModeDstOut
geom.BlendModeSrcATop
geom.BlendModeDstATop
geom.BlendModeXor
geom.BlendModePlus
geom.BlendModeModulate

// Advanced blend modes
geom.BlendModeScreen
geom.BlendModeOverlay
geom.BlendModeDarken
geom.BlendModeLighten
geom.BlendModeColorDodge
geom.BlendModeColorBurn
geom.BlendModeHardLight
geom.BlendModeSoftLight
geom.BlendModeDifference
geom.BlendModeExclusion
geom.BlendModeMultiply

// HSV blend modes
geom.BlendModeHue
geom.BlendModeSaturation
geom.BlendModeColor
geom.BlendModeLuminosity
```

## Performance Features

### Wang's Formula for Curve Tessellation

```go
// Optimal subdivision count for smooth curves
scaleFactor := transform.MaxBasisLengthXY()

// For cubic curves
subdivisions := geom.CubicSubdivisions(scaleFactor, p0, p1, p2, p3)

// For quadratic curves
subdivisions := geom.QuadraticSubdivisions(scaleFactor, p0, p1, p2)

// For conic curves
subdivisions := geom.ConicSubdivisions(scaleFactor, p0, p1, p2, weight)
```

### Efficient 2D Transformations

```go
// RSTransform for efficient sprite transformations
transform := geom.NewRSTransform(origin, scale, rotation)

// Generate transformed quad
quad := transform.QuadSize(spriteSize)
bounds := transform.BoundsSize(spriteSize)

// Check if transformation is axis-aligned for optimization
if transform.IsAxisAligned() {
    // Use faster axis-aligned rendering path
}
```

## Type Safety and Generics

The library extensively uses Go generics for type safety and performance:

```go
// All geometric types are generic over scalar types
type Point[T Scalar] struct {
    X, Y T
}

type Matrix[T Scalar] [16]T

type Color struct {
    R, G, B, A geom.F32  // Colors use F32 for consistency
}

// The Scalar interface ensures type compatibility
type Scalar interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
    Float64() float64
    String() string
}
```

## Integration Examples

### With Graphics Libraries

```go
// Convert to standard Go types
goRect := rect.ToGo()                    // image.Rectangle
goPoint := point.ToGo()                  // image.Point
goColor := color.ToRGBA()                // color.RGBA

// From standard Go types
rect := geom.NewRectFromGo[geom.F32](goRect)
point := geom.NewPointFromGo[geom.F32](goPoint)
color := geom.NewColorFromRGBA(goColor)
```

### Custom Rendering Pipeline

```go
type Renderer struct {
    transform geom.Matrix[geom.F32]
}

func (r *Renderer) DrawRect(rect geom.Rect[geom.F32], color geom.Color) {
    // Transform rectangle
    corners := rect.Transform(r.transform)

    // Convert to render format
    vertices := make([]float32, 8)
    for i, corner := range corners {
        vertices[i*2] = float32(corner.X)
        vertices[i*2+1] = float32(corner.Y)
    }

    // Submit to GPU...
}

func (r *Renderer) DrawPath(path geom.PathSource[geom.F32]) {
    receiver := &r.pathReceiver
    path.Dispatch(receiver)
}
```

## Best Practices

1. **Choose Appropriate Scalar Types**:
   - Use `F32` for general graphics work
   - Use `F64` for high-precision calculations
   - Use `I32` for pixel-perfect integer coordinates
   - Use `I26_6` for typography and precise measurements

2. **Leverage Type Safety**:
   - Use strongly-typed angles (`Radians`/`Degrees`)
   - Prefer `ScalarEq()` over `==` for floating-point comparisons
   - Use generic types consistently throughout your codebase

3. **Optimize Transformations**:
   - Check matrix properties before expensive operations
   - Use `RSTransform` for simple 2D transformations
   - Cache transformation matrices when possible

4. **Memory Management**:
   - Reuse geometric objects when possible
   - Use value types (structs) instead of pointers for small objects
   - Consider object pooling for frequently created/destroyed geometry

## Dependencies

- Go 1.21+ (requires generics support)
- Standard library only
