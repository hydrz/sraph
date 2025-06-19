package geometry

import (
	"math"

	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// EllipseGeometry represents elliptical geometry
type EllipseGeometry struct {
	BaseGeometry
	bounds    geom.Rect
	segments  int
	transform geom.Matrix
}

// NewEllipseGeometry creates a new ellipse geometry from bounds
func NewEllipseGeometry(bounds geom.Rect) *EllipseGeometry {
	return &EllipseGeometry{
		BaseGeometry: NewBaseGeometry(),
		bounds:       bounds,
		segments:     32, // Default segments
		transform:    geom.NewIdentityMatrix(),
	}
}

// SetBounds sets the ellipse bounds
func (e *EllipseGeometry) SetBounds(bounds geom.Rect) {
	e.bounds = bounds
}

// GetBounds returns the ellipse bounds
func (e *EllipseGeometry) GetBounds() geom.Rect {
	return e.transform.TransformRect(e.bounds)
}

// SetSegments sets the number of segments for tessellation
func (e *EllipseGeometry) SetSegments(segments int) {
	if segments >= 3 {
		e.segments = segments
	}
}

// GetSegments returns the number of segments
func (e *EllipseGeometry) GetSegments() int {
	return e.segments
}

// SetTransform sets the transformation matrix
func (e *EllipseGeometry) SetTransform(transform geom.Matrix) {
	e.transform = transform
}

// GetTransform returns the transformation matrix
func (e *EllipseGeometry) GetTransform() geom.Matrix {
	return e.transform
}

// GetCenter returns the center point of the ellipse
func (e *EllipseGeometry) GetCenter() geom.Point {
	return geom.Point{
		X: e.bounds.X + e.bounds.Width/2,
		Y: e.bounds.Y + e.bounds.Height/2,
	}
}

// GetRadiusX returns the X radius
func (e *EllipseGeometry) GetRadiusX() float32 {
	return e.bounds.Width / 2
}

// GetRadiusY returns the Y radius
func (e *EllipseGeometry) GetRadiusY() float32 {
	return e.bounds.Height / 2
}

// IsAxisAlignedRect returns true if this ellipse is actually a rectangle
func (e *EllipseGeometry) IsAxisAlignedRect() bool {
	// An ellipse is never a rectangle
	return false
}

// CoversArea returns true if this ellipse covers the specified area
func (e *EllipseGeometry) CoversArea(transform geom.Matrix, rect geom.Rect) bool {
	// TODO: Implement proper ellipse-rectangle intersection test
	ellipseBounds := e.GetBounds()
	transformedBounds := transform.TransformRect(ellipseBounds)
	return transformedBounds.Contains(rect)
}

// GetVertexCount returns the number of vertices
func (e *EllipseGeometry) GetVertexCount() int {
	return e.segments + 1 // Center + perimeter vertices
}

// GetIndexCount returns the number of indices
func (e *EllipseGeometry) GetIndexCount() int {
	return e.segments * 3 // One triangle per segment from center
}

// ComputeVertices generates vertex data for rendering
func (e *EllipseGeometry) ComputeVertices() []render.Vertex {
	vertices := make([]render.Vertex, e.GetVertexCount())

	center := e.GetCenter()
	radiusX := e.GetRadiusX()
	radiusY := e.GetRadiusY()

	// Center vertex
	vertices[0] = render.Vertex{
		Position: center,
		TexCoord: geom.Point{X: 0.5, Y: 0.5},
	}

	// Perimeter vertices
	for i := 0; i < e.segments; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(e.segments)
		cos := float32(math.Cos(angle))
		sin := float32(math.Sin(angle))

		x := center.X + cos*radiusX
		y := center.Y + sin*radiusY

		vertices[i+1] = render.Vertex{
			Position: geom.Point{X: x, Y: y},
			TexCoord: geom.Point{
				X: 0.5 + cos*0.5,
				Y: 0.5 + sin*0.5,
			},
		}
	}

	return vertices
}

// GetIndexBuffer generates index data
func (e *EllipseGeometry) GetIndexBuffer() []uint16 {
	indices := make([]uint16, e.GetIndexCount())

	for i := 0; i < e.segments; i++ {
		base := i * 3
		indices[base] = 0                              // Center
		indices[base+1] = uint16(i + 1)                // Current vertex
		indices[base+2] = uint16((i+1)%e.segments + 1) // Next vertex
	}

	return indices
}

// Clone creates a copy of this geometry
func (e *EllipseGeometry) Clone() Geometry {
	return &EllipseGeometry{
		BaseGeometry: e.BaseGeometry.Clone().(BaseGeometry),
		bounds:       e.bounds,
		segments:     e.segments,
		transform:    e.transform,
	}
}

// GetCoverage returns the coverage area for this geometry
func (e *EllipseGeometry) GetCoverage(transform geom.Matrix) geom.Rect {
	combinedTransform := transform.Multiply(e.transform)
	return combinedTransform.TransformRect(e.bounds)
}
