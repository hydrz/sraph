package display

import (
	"math"

	"github.com/opensraph/sraph/geom"
)

// Utility functions and helpers for display list operations

// RectFromXYWH creates a rectangle from x, y, width, height.
func RectFromXYWH(x, y, width, height Scalar) geom.Rect[Scalar] {
	return geom.NewRectXYWH(x, y, width, height)
}

// RectFromCenter creates a rectangle centered at the given point with the specified size.
func RectFromCenter(center geom.Point[Scalar], size geom.Size[Scalar]) geom.Rect[Scalar] {
	halfWidth := size.Width / 2
	halfHeight := size.Height / 2
	return geom.NewRect(
		center.X-halfWidth, center.Y-halfHeight,
		center.X+halfWidth, center.Y+halfHeight,
	)
}

// RectFromPoints creates a rectangle that encompasses the given points.
func RectFromPoints(points []geom.Point[Scalar]) geom.Rect[Scalar] {
	if len(points) == 0 {
		return geom.Rect[Scalar]{}
	}

	min := points[0]
	max := points[0]

	for _, pt := range points[1:] {
		if pt.X < min.X {
			min.X = pt.X
		}
		if pt.Y < min.Y {
			min.Y = pt.Y
		}
		if pt.X > max.X {
			max.X = pt.X
		}
		if pt.Y > max.Y {
			max.Y = pt.Y
		}
	}

	return geom.NewRect(min.X, min.Y, max.X, max.Y)
}

// DegreesToRadians converts degrees to radians.
func DegreesToRadians(degrees Scalar) Scalar {
	return degrees * (3.14159265359 / 180.0)
}

// RadiansToDegrees converts radians to degrees.
func RadiansToDegrees(radians Scalar) Scalar {
	return radians * (180.0 / 3.14159265359)
}

// LerpScalar linearly interpolates between two scalar values.
func LerpScalar(a, b, t Scalar) Scalar {
	return a + (b-a)*t
}

// LerpPoint linearly interpolates between two points.
func LerpPoint(a, b geom.Point[Scalar], t Scalar) geom.Point[Scalar] {
	return geom.Point[Scalar]{
		X: LerpScalar(a.X, b.X, t),
		Y: LerpScalar(a.Y, b.Y, t),
	}
}

// LerpColor linearly interpolates between two colors.
func LerpColor(a, b Color, t Scalar) Color {

	return NewColor(
		float32(LerpScalar(a.GeomColor().R, b.GeomColor().R, t)),
		float32(LerpScalar(a.GeomColor().G, b.GeomColor().G, t)),
		float32(LerpScalar(a.GeomColor().B, b.GeomColor().B, t)),
		float32(LerpScalar(a.GeomColor().A, b.GeomColor().A, t)),
	)
}

// DistanceSquared calculates the squared distance between two points.
func DistanceSquared(a, b geom.Point[Scalar]) Scalar {
	dx := b.X - a.X
	dy := b.Y - a.Y
	return dx*dx + dy*dy
}

// Distance calculates the distance between two points.
func Distance(a, b geom.Point[Scalar]) Scalar {
	return Scalar(math.Sqrt(float64(DistanceSquared(a, b))))
}

// Normalize normalizes a vector to unit length.
func Normalize(v geom.Point[Scalar]) geom.Point[Scalar] {
	length := Distance(geom.Point[Scalar]{}, v)
	if length == 0 {
		return geom.Point[Scalar]{}
	}
	return geom.Point[Scalar]{X: v.X / length, Y: v.Y / length}
}

// RotatePoint rotates a point around the origin by the given angle in radians.
func RotatePoint(point geom.Point[Scalar], radians Scalar) geom.Point[Scalar] {
	cos := Scalar(math.Cos(float64(radians)))
	sin := Scalar(math.Sin(float64(radians)))

	return geom.Point[Scalar]{
		X: point.X*cos - point.Y*sin,
		Y: point.X*sin + point.Y*cos,
	}
}

// RotatePointAround rotates a point around another point by the given angle.
func RotatePointAround(point, center geom.Point[Scalar], radians Scalar) geom.Point[Scalar] {
	// Translate to origin
	translated := geom.Point[Scalar]{X: point.X - center.X, Y: point.Y - center.Y}
	// Rotate
	rotated := RotatePoint(translated, radians)
	// Translate back
	return geom.Point[Scalar]{X: rotated.X + center.X, Y: rotated.Y + center.Y}
}

// QuadraticBezier calculates a point on a quadratic bezier curve.
func QuadraticBezier(p0, p1, p2 geom.Point[Scalar], t Scalar) geom.Point[Scalar] {
	oneMinusT := 1 - t
	return geom.Point[Scalar]{
		X: oneMinusT*oneMinusT*p0.X + 2*oneMinusT*t*p1.X + t*t*p2.X,
		Y: oneMinusT*oneMinusT*p0.Y + 2*oneMinusT*t*p1.Y + t*t*p2.Y,
	}
}

// CubicBezier calculates a point on a cubic bezier curve.
func CubicBezier(p0, p1, p2, p3 geom.Point[Scalar], t Scalar) geom.Point[Scalar] {
	oneMinusT := 1 - t
	oneMinusT2 := oneMinusT * oneMinusT
	oneMinusT3 := oneMinusT2 * oneMinusT
	t2 := t * t
	t3 := t2 * t

	return geom.Point[Scalar]{
		X: oneMinusT3*p0.X + 3*oneMinusT2*t*p1.X + 3*oneMinusT*t2*p2.X + t3*p3.X,
		Y: oneMinusT3*p0.Y + 3*oneMinusT2*t*p1.Y + 3*oneMinusT*t2*p2.Y + t3*p3.Y,
	}
}

// TransformBuilder provides a fluent interface for building transformation matrices.
type TransformBuilder struct {
	matrix geom.Matrix[Scalar]
}

// NewTransformBuilder creates a new TransformBuilder with identity matrix.
func NewTransformBuilder() *TransformBuilder {
	return &TransformBuilder{
		matrix: geom.NewMatrix[Scalar](),
	}
}

// Translate adds a translation to the transformation.
func (tb *TransformBuilder) Translate(dx, dy Scalar) *TransformBuilder {
	tb.matrix = tb.matrix.Translate(geom.Vector2[Scalar]{X: dx, Y: dy})
	return tb
}

// Scale adds a scale to the transformation.
func (tb *TransformBuilder) Scale(sx, sy Scalar) *TransformBuilder {
	tb.matrix = tb.matrix.Scale(geom.Vector2[Scalar]{X: sx, Y: sy})
	return tb
}

// Rotate adds a rotation to the transformation.
func (tb *TransformBuilder) Rotate(radians Scalar) *TransformBuilder {
	tb.matrix = tb.matrix.RotateZ(geom.Radians(radians))
	return tb
}

// RotateDegrees adds a rotation in degrees to the transformation.
func (tb *TransformBuilder) RotateDegrees(degrees Scalar) *TransformBuilder {
	return tb.Rotate(DegreesToRadians(degrees))
}

// Build returns the final transformation matrix.
func (tb *TransformBuilder) Build() geom.Matrix[Scalar] {
	return tb.matrix
}

// Reset resets the transformation to identity.
func (tb *TransformBuilder) Reset() *TransformBuilder {
	tb.matrix = geom.NewMatrix[Scalar]()
	return tb
}

// PaintBuilder provides a fluent interface for building paint objects.
type PaintBuilder struct {
	paint Paint
}

// NewPaintBuilder creates a new PaintBuilder.
func NewPaintBuilder() *PaintBuilder {
	return &PaintBuilder{
		paint: NewPaint(),
	}
}

// WithColor sets the paint color.
func (pb *PaintBuilder) WithColor(color Color) *PaintBuilder {
	pb.paint.SetColor(color)
	return pb
}

// WithBlendMode sets the blend mode.
func (pb *PaintBuilder) WithBlendMode(mode BlendMode) *PaintBuilder {
	pb.paint.SetBlendMode(mode)
	return pb
}

// WithAntiAlias sets anti-aliasing.
func (pb *PaintBuilder) WithAntiAlias(enabled bool) *PaintBuilder {
	pb.paint.SetAntiAlias(enabled)
	return pb
}

// Build returns the final paint object.
func (pb *PaintBuilder) Build() Paint {
	return pb.paint
}

// PathHelper provides utility functions for common path operations.
type PathHelper struct{}

// CreateRoundedRectPath creates a path for a rounded rectangle.
func (PathHelper) CreateRoundedRectPath(rect geom.Rect[Scalar], radius Scalar) *Path {
	pb := NewPathBuilder()
	pb.AddRoundRect(geom.NewRoundRectRadius(rect, radius))
	return pb.Build()
}

// CreateStarPath creates a path for a star shape.
func (PathHelper) CreateStarPath(center geom.Point[Scalar], outerRadius, innerRadius Scalar, points int) *Path {
	pb := NewPathBuilder()

	angleStep := 2 * 3.14159265359 / Scalar(points*2)

	for i := 0; i < points*2; i++ {
		angle := Scalar(i) * angleStep
		radius := outerRadius
		if i%2 == 1 {
			radius = innerRadius
		}

		x := center.X + radius*Scalar(math.Cos(float64(angle)))
		y := center.Y + radius*Scalar(math.Sin(float64(angle)))

		if i == 0 {
			pb.MoveTo(x, y)
		} else {
			pb.LineTo(x, y)
		}
	}

	pb.Close()
	return pb.Build()
}

// CreateArrowPath creates a path for an arrow shape.
func (PathHelper) CreateArrowPath(start, end geom.Point[Scalar], arrowheadLength, arrowheadWidth Scalar) *Path {
	pb := NewPathBuilder()

	// Calculate direction vector
	dx := end.X - start.X
	dy := end.Y - start.Y
	length := Distance(start, end)

	if length == 0 {
		return pb.Build()
	}

	// Normalize direction
	dirX := dx / length
	dirY := dy / length

	// Perpendicular vector
	perpX := -dirY
	perpY := dirX

	// Arrow shaft
	pb.MoveTo(start.X, start.Y)
	pb.LineTo(end.X, end.Y)

	// Arrowhead
	arrowBaseX := end.X - dirX*arrowheadLength
	arrowBaseY := end.Y - dirY*arrowheadLength

	arrowLeft := geom.Point[Scalar]{
		X: arrowBaseX + perpX*arrowheadWidth/2,
		Y: arrowBaseY + perpY*arrowheadWidth/2,
	}
	arrowRight := geom.Point[Scalar]{
		X: arrowBaseX - perpX*arrowheadWidth/2,
		Y: arrowBaseY - perpY*arrowheadWidth/2,
	}

	pb.MoveTo(arrowLeft.X, arrowLeft.Y)
	pb.LineTo(end.X, end.Y)
	pb.LineTo(arrowRight.X, arrowRight.Y)

	return pb.Build()
}

// BoundsHelper provides utility functions for bounds calculations.
type BoundsHelper struct{}

// Union combines multiple rectangles into a single bounding rectangle.
func (BoundsHelper) Union(rects ...geom.Rect[Scalar]) geom.Rect[Scalar] {
	if len(rects) == 0 {
		return geom.Rect[Scalar]{}
	}

	result := rects[0]
	for _, rect := range rects[1:] {
		result = result.Union(rect)
	}

	return result
}

// Intersect finds the intersection of multiple rectangles.
func (BoundsHelper) Intersect(rects ...geom.Rect[Scalar]) geom.Rect[Scalar] {
	if len(rects) == 0 {
		return geom.Rect[Scalar]{}
	}

	result := rects[0]
	for _, rect := range rects[1:] {
		result = result.Intersect(rect)
		// If intersection becomes empty, return empty rect
		if result.IsEmpty() {
			return geom.Rect[Scalar]{}
		}
	}

	return result
}

// Expand expands a rectangle by the given amount in all directions.
func (BoundsHelper) Expand(rect geom.Rect[Scalar], amount Scalar) geom.Rect[Scalar] {
	return geom.NewRect(
		rect.Left-amount, rect.Top-amount,
		rect.Right+amount, rect.Bottom+amount,
	)
}

// Contract contracts a rectangle by the given amount in all directions.
func (BoundsHelper) Contract(rect geom.Rect[Scalar], amount Scalar) geom.Rect[Scalar] {
	return geom.NewRect(
		rect.Left+amount, rect.Top+amount,
		rect.Right-amount, rect.Bottom-amount,
	)
}
