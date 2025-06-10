# Sraph Geometry Package

The `geom` package provides a comprehensive set of 2D and 3D geometric primitives and operations for the Sraph graphics rendering engine. It includes points, vectors, matrices, colors, gradients, and various geometric shapes with efficient implementations optimized for graphics applications.

## Table of Contents

- [Installation](#installation)
- [Core Types](#core-types)
- [Basic Geometric Primitives](#basic-geometric-primitives)
- [Colors and Gradients](#colors-and-gradients)
- [Transformations](#transformations)
- [Shapes and Paths](#shapes-and-paths)
- [Performance Features](#performance-features)
- [Examples](#examples)
- [API Reference](#api-reference)

## Installation

```go
import "github.com/opensraph/sraph/geom"
```

## Core Types

### Scalar Type
All geometric operations use `Scalar` (float32) for consistency and performance:

```go
type Scalar = float32
```

### Number Interface
Generic numeric operations support multiple types through the `Number` interface:

```go
type Number = matht.Number  // Supports int32, int64, float32, float64
```

## Basic Geometric Primitives

### Points and Vectors

#### 2D Points
```go
// Create points
p1 := geom.NewPoint(10.0, 20.0)
p2 := geom.NewPointZero()

// Arithmetic operations
p3 := p1.Add(p2)
p4 := p1.Sub(p2)
p5 := p1.MulScalar(2.0)

// Distance and length
distance := p1.Distance(p2)
length := p1.Length()

// Normalization
normalized := p1.Normalize()

// Rotation
rotated := p1.Rotate(geom.Radians(math.Pi / 4)) // 45 degrees
```

#### 3D Vectors
```go
// Create 3D vectors
v1 := geom.NewVector3(1.0, 2.0, 3.0)
v2 := geom.NewVector3Zero()

// Vector operations
dotProduct := v1.Dot(v2)
crossProduct := v1.Cross(v2)
normalized := v1.Normalize()

// Transform to other types
point := v1.ToPoint()  // Projects to 2D
color := geom.Vector3FromColor(geom.ColorRed)
```

#### 4D Vectors
```go
// Create 4D vectors (useful for colors and homogeneous coordinates)
v4 := geom.NewVector4(1.0, 2.0, 3.0, 1.0)
fromColor := geom.Vector4FromColor(geom.ColorBlue)

// Extract components
xyz := v4.XYZ()  // Returns Vector3
xy := v4.XY()    // Returns Vector2
```

### Rectangles and Sizes

#### Rectangles
```go
// Create rectangles
rect1 := geom.NewRect(0, 0, 100, 200)  // left, top, right, bottom
rect2 := geom.NewRectXYWH(10, 20, 50, 100)  // x, y, width, height

// Rectangle operations
area := rect1.Area()
center := rect1.Center()
size := rect1.Size()

// Containment and intersection
contains := rect1.Contains(geom.NewPoint(50, 50))
intersects := rect1.IntersectsWithRect(rect2)
intersection := rect1.Intersect(rect2)
union := rect1.Union(rect2)

// Transformations
expanded := rect1.ExpandAll(10)  // Expand by 10 units on all sides
shifted := rect1.Shift(5, 5)    // Translate by (5, 5)
```

#### Sizes
```go
// Create sizes
size1 := geom.NewSize(100, 200)
size2 := geom.NewSizeDim(50)  // Square size 50x50

// Size operations
area := size1.Area()
aspectRatio := size1.AspectRatio()
scaled := size1.Mul(2.0)

// State checking
isEmpty := size1.IsEmpty()
isSquare := size1.IsSquare()
```

## Colors and Gradients

### Colors
```go
// Create colors (RGBA values in range [0,1])
red := geom.NewColorRGB(1.0, 0.0, 0.0)
blue := geom.NewColor(0.0, 0.0, 1.0, 0.8)  // With alpha
gray := geom.NewColorGray(0.5, 1.0)

// From 8-bit values
color8bit := geom.NewRGBA8(255, 128, 64, 255)

// From hex values
hexColor := geom.NewColorFromHex(0xFF8040FF, true)  // With alpha
rgbHex := geom.NewColorFromHex(0xFF8040, false)     // RGB only

// Color operations
blended := red.Blend(blue, 0.5)  // 50% blend
premult := red.Premultiply()     // Premultiplied alpha
hsv := red.ToHSV()               // Convert to HSV

// Predefined colors
black := geom.ColorBlack
white := geom.ColorWhite
transparent := geom.ColorTransparent
```

### Gradients
```go
// Simple two-color gradient
gradient := geom.NewGradientTwoColor(geom.ColorRed, geom.ColorBlue)

// Multi-stop gradient
stops := []geom.GradientStop{
    geom.NewGradientStop(0.0, geom.ColorRed),
    geom.NewGradientStop(0.5, geom.ColorYellow),
    geom.NewGradientStop(1.0, geom.ColorBlue),
}
multiGradient := geom.NewGradient(stops)

// Rainbow gradient
rainbow := geom.NewRainbowGradient()

// Sample colors from gradient
color := gradient.SampleColor(0.25)  // Get color at 25% position

// Convert to texture data for GPU rendering
textureData := gradient.ToGradientData()
```

## Transformations

### Matrices
```go
// Create transformation matrices
identity := geom.NewMatrixIdentity()
translation := geom.NewTranslationMatrix(geom.NewVector3(10, 20, 0))
rotation := geom.NewRotationZMatrix(geom.Radians(math.Pi / 4))
scale := geom.NewScaleMatrix(geom.NewVector3(2, 2, 1))

// Combine transformations
combined := translation.Multiply(rotation).Multiply(scale)

// Transform points and vectors
point := geom.NewPoint(1, 1)
transformed := combined.TransformPoint(point)

vector := geom.NewVector2(1, 0)
rotatedVector := combined.TransformVector2(vector)

// Matrix properties
determinant := combined.Determinant()
inverse := combined.Inverse()
isIdentity := combined.IsIdentity()
```

### RS Transforms (Rotation, Scale, Translation)
```go
// Optimized 2D transform for sprites
rs := geom.NewRSTransformFromOriginScaleRotation(
    geom.NewPoint(100, 100),  // origin
    2.0,                      // scale
    geom.Radians(math.Pi/4),  // rotation
)

// Transform points efficiently
point := geom.NewPoint(10, 10)
transformed := rs.TransformPoint(point)

// Combine transforms
rs2 := geom.NewRSTransformTranslation(geom.NewPoint(50, 50))
combined := rs.Compose(rs2)

// Convert to matrix if needed
matrix := rs.ToMatrix()
```

### Quaternions (3D Rotations)
```go
// Create quaternions
identity := geom.NewQuaternionIdentity()
axisAngle := geom.NewQuaternionFromAxisAngle(
    geom.NewVector3(0, 1, 0),  // Y-axis
    geom.Radians(math.Pi/2),   // 90 degrees
)

// Quaternion operations
q1 := geom.NewQuaternion(0, 0, 0, 1)
q2 := geom.NewQuaternion(0, 1, 0, 0)
combined := q1.Mul(q2)  // Compose rotations

// Rotate vectors
vector := geom.NewVector3(1, 0, 0)
rotated := q1.RotateVector(vector)

// Interpolation
interpolated := q1.Slerp(q2, 0.5)  // Smooth interpolation

// Convert to matrix
rotationMatrix := q1.ToMatrix()
```

## Shapes and Paths

### Round Rectangles
```go
// Create round rectangles
rect := geom.NewRect(0, 0, 100, 100)

// Uniform corner radius
roundRect := geom.NewRoundRectFromRectRadius(rect, 10)

// Different radii per corner
radii := geom.NewRoundingRadiiCorners(
    geom.NewSize(5, 5),   // top-left
    geom.NewSize(10, 10), // top-right
    geom.NewSize(15, 15), // bottom-left
    geom.NewSize(20, 20), // bottom-right
)
customRoundRect := geom.NewRoundRect(rect, radii)

// Shape properties
isOval := roundRect.IsOval()
isRect := roundRect.IsRect()
containsPoint := roundRect.Contains(geom.NewPoint(50, 50))

// Transform shapes
scaled := roundRect.Scale(2.0)
translated := roundRect.Shift(10, 10)
```

### Path Sources
```go
// Rectangle path
rectPath := geom.NewRectPathSource(geom.NewRect(0, 0, 100, 100))

// Round rectangle path
roundRectPath := geom.NewRoundRectPathSource(roundRect)

// Ellipse path
ellipsePath := geom.NewEllipsePathSource(geom.NewRect(0, 0, 100, 100))

// Use paths (example with a hypothetical path renderer)
// renderer.DrawPath(rectPath, paint)
```

## Performance Features

### Half-Precision Floats
For GPU optimization and memory efficiency:

```go
// Convert to half-precision
halfValue := geom.NewHalf(3.14159)
halfVector := geom.NewHalfVector2(1.0, 2.0)

// Batch conversions
fullPrecision := []geom.Scalar{1.0, 2.0, 3.0, 4.0}
halfPrecision := geom.Float32SliceToHalf(fullPrecision)

// Convert back
restored := geom.HalfSliceToFloat32(halfPrecision)
```

### Separated Vectors
For efficient polyline processing:

```go
// Create separated vector (direction + magnitude)
vector := geom.NewVector2(3.0, 4.0)
separated := geom.NewSeparatedVector2FromVector(vector)

// Access components
direction := separated.Direction  // Normalized direction
magnitude := separated.Magnitude  // Length

// Efficient operations
rotated := separated.Rotate(geom.Radians(math.Pi / 4))
scaled := separated.Scale(2.0)

// Reconstruct vector
reconstructed := separated.Vector()
```

## Examples

### Basic 2D Graphics Setup
```go
func setup2DGraphics() {
    // Create viewport
    viewport := geom.NewRect(0, 0, 800, 600)

    // Create orthographic projection
    projection := geom.NewOrtho2DMatrix(0, 800, 600, 0)

    // Create view transform
    view := geom.NewMatrixIdentity()

    // Combine transformations
    mvp := projection.Multiply(view)

    fmt.Printf("Viewport: %s\n", viewport)
    fmt.Printf("MVP Matrix: %s\n", mvp)
}
```

### Color Manipulation
```go
func colorExample() {
    // Create base colors
    red := geom.ColorRed
    blue := geom.ColorBlue

    // Create gradient
    gradient := geom.NewGradientTwoColor(red, blue)

    // Sample colors along gradient
    for i := 0; i <= 10; i++ {
        t := float32(i) / 10.0
        color := gradient.SampleColor(t)

        // Convert to different formats
        rgba8 := color.ToRGBA8()
        hex := color.ToHex(true)

        fmt.Printf("t=%.1f: RGBA8=%v, Hex=%s\n", t, rgba8, hex)
    }
}
```

### Shape Intersection
```go
func shapeIntersection() {
    // Create shapes
    rect1 := geom.NewRect(0, 0, 100, 100)
    rect2 := geom.NewRect(50, 50, 150, 150)

    // Create round rectangles
    round1 := geom.NewRoundRectFromRectRadius(rect1, 10)
    round2 := geom.NewRoundRectFromRectRadius(rect2, 15)

    // Test containment
    testPoint := geom.NewPoint(75, 75)
    in1 := round1.Contains(testPoint)
    in2 := round2.Contains(testPoint)

    fmt.Printf("Point %s: in shape1=%t, in shape2=%t\n", testPoint, in1, in2)

    // Rectangle intersection
    intersection := rect1.Intersect(rect2)
    fmt.Printf("Rectangle intersection: %s\n", intersection)
}
```

### Animation with Interpolation
```go
func animationExample() {
    // Start and end transforms
    start := geom.NewRSTransformFromOriginScaleRotation(
        geom.NewPoint(0, 0), 1.0, 0,
    )
    end := geom.NewRSTransformFromOriginScaleRotation(
        geom.NewPoint(100, 100), 2.0, geom.Radians(math.Pi),
    )

    // Animate over time
    steps := 10
    for i := 0; i <= steps; i++ {
        t := float32(i) / float32(steps)

        // Interpolate components
        scale := 1.0 + t*1.0  // 1.0 to 2.0
        rotation := t * float32(math.Pi)  // 0 to π
        position := geom.NewPoint(0, 0).Lerp(geom.NewPoint(100, 100), t)

        // Create interpolated transform
        current := geom.NewRSTransformFromOriginScaleRotation(
            position, scale, geom.Radians(rotation),
        )

        fmt.Printf("Step %d: %s\n", i, current)
    }
}
```

## API Reference

### Constants
- `Epsilon`: Small value for floating-point comparisons (1e-3)
- `Sqrt2Over2`: √2/2, useful for conic sections

### Core Types
- `Scalar`: Primary floating-point type (float32)
- `Point`, `PointI32`, `PointI64`: 2D points with different numeric types
- `Vector2`, `Vector3`, `Vector4`: Vector types
- `Size`, `SizeI32`, `SizeI64`: Size types
- `Rect`, `RectI32`, `RectI64`: Rectangle types

### Colors and Visual Elements
- `Color`: RGBA color with floating-point components
- `Gradient`, `GradientStop`: Gradient definitions
- `GradientData`: GPU-ready gradient texture data

### Geometric Shapes
- `RoundRect`: Rectangle with rounded corners
- `RoundingRadii`: Corner radius specifications
- `RoundSuperellipse`: Advanced rounded rectangle with superellipse curves

### Transformations
- `Matrix`: 4x4 transformation matrix
- `RSTransform`: Optimized 2D rotation/scale/translation
- `Quaternion`: 3D rotation representation
- `MatrixDecomposition`: Decomposed transformation components

### Path and Rendering
- `PathSource`: Interface for renderable paths
- `PathReceiver`: Interface for path command processing
- `FillType`: Path filling rules (NonZero, Odd)

### Performance Types
- `Half`, `HalfVector2`, `HalfVector3`, `HalfVector4`: Half-precision types
- `SeparatedVector2`: Direction/magnitude vector representation

### Utility Types
- `BlurParameters`: Gaussian blur parameters
- `StrokeParameters`: Stroke rendering parameters
- `Radians`, `Degrees`: Angle types
- `Rational`: Rational number representation

For complete API documentation, see the [package documentation](https://pkg.go.dev/github.com/opensraph/sraph/geom).

## Contributing

This package is part of the Sraph graphics engine. Contributions should follow the project's coding standards and include appropriate tests.

## License

See the main Sraph project for license information.
