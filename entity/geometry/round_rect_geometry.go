package geometry

import (
	"math"

	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// RoundRectGeometry represents rounded rectangle geometry
type RoundRectGeometry struct {
	BaseGeometry
	bounds    geom.Rect
	radii     geom.Size // Uniform radii for all corners
	segments  int       // Segments per corner
	transform geom.Matrix
}

// NewRoundRectGeometry creates a new rounded rectangle geometry
func NewRoundRectGeometry(bounds geom.Rect, radii geom.Size) *RoundRectGeometry {
	return &RoundRectGeometry{
		BaseGeometry: NewBaseGeometry(),
		bounds:       bounds,
		radii:        radii,
		segments:     8, // Default segments per corner
		transform:    geom.NewIdentityMatrix(),
	}
}

// SetBounds sets the rectangle bounds
func (r *RoundRectGeometry) SetBounds(bounds geom.Rect) {
	r.bounds = bounds
}

// GetBounds returns the rectangle bounds
func (r *RoundRectGeometry) GetBounds() geom.Rect {
	return r.transform.TransformRect(r.bounds)
}

// SetRadii sets the corner radii
func (r *RoundRectGeometry) SetRadii(radii geom.Size) {
	r.radii = radii
}

// GetRadii returns the corner radii
func (r *RoundRectGeometry) GetRadii() geom.Size {
	return r.radii
}

// SetSegments sets the number of segments per corner
func (r *RoundRectGeometry) SetSegments(segments int) {
	if segments >= 2 {
		r.segments = segments
	}
}

// GetSegments returns the number of segments per corner
func (r *RoundRectGeometry) GetSegments() int {
	return r.segments
}

// SetTransform sets the transformation matrix
func (r *RoundRectGeometry) SetTransform(transform geom.Matrix) {
	r.transform = transform
}

// GetTransform returns the transformation matrix
func (r *RoundRectGeometry) GetTransform() geom.Matrix {
	return r.transform
}

// IsAxisAlignedRect returns true if this is actually a rectangle (no rounding)
func (r *RoundRectGeometry) IsAxisAlignedRect() bool {
	return r.radii.Width <= 0 && r.radii.Height <= 0
}

// CoversArea returns true if this rounded rect covers the specified area
func (r *RoundRectGeometry) CoversArea(transform geom.Matrix, rect geom.Rect) bool {
	// TODO: Implement proper rounded rectangle coverage test
	roundRectBounds := r.GetBounds()
	transformedBounds := transform.TransformRect(roundRectBounds)
	return transformedBounds.Contains(rect)
}

// GetVertexCount returns the number of vertices
func (r *RoundRectGeometry) GetVertexCount() int {
	if r.IsAxisAlignedRect() {
		return 4 // Simple rectangle
	}

	// Center + 4 corners * segments + 4 sides * 2
	return 1 + (4 * r.segments) + (4 * 2)
}

// GetIndexCount returns the number of indices
func (r *RoundRectGeometry) GetIndexCount() int {
	if r.IsAxisAlignedRect() {
		return 6 // Two triangles
	}

	vertexCount := r.GetVertexCount()
	// Triangle fan from center
	return (vertexCount - 1) * 3
}

// ComputeVertices generates vertex data for rendering
func (r *RoundRectGeometry) ComputeVertices() []render.Vertex {
	if r.IsAxisAlignedRect() {
		return r.computeRectVertices()
	}
	return r.computeRoundedVertices()
}

// computeRectVertices generates vertices for a simple rectangle
func (r *RoundRectGeometry) computeRectVertices() []render.Vertex {
	return []render.Vertex{
		{
			Position: geom.Point{X: r.bounds.X, Y: r.bounds.Y},
			TexCoord: geom.Point{X: 0, Y: 0},
		},
		{
			Position: geom.Point{X: r.bounds.X + r.bounds.Width, Y: r.bounds.Y},
			TexCoord: geom.Point{X: 1, Y: 0},
		},
		{
			Position: geom.Point{X: r.bounds.X + r.bounds.Width, Y: r.bounds.Y + r.bounds.Height},
			TexCoord: geom.Point{X: 1, Y: 1},
		},
		{
			Position: geom.Point{X: r.bounds.X, Y: r.bounds.Y + r.bounds.Height},
			TexCoord: geom.Point{X: 0, Y: 1},
		},
	}
}

// computeRoundedVertices generates vertices for a rounded rectangle
func (r *RoundRectGeometry) computeRoundedVertices() []render.Vertex {
	vertices := make([]render.Vertex, 0, r.GetVertexCount())

	// Center vertex
	center := geom.Point{
		X: r.bounds.X + r.bounds.Width/2,
		Y: r.bounds.Y + r.bounds.Height/2,
	}
	vertices = append(vertices, render.Vertex{
		Position: center,
		TexCoord: geom.Point{X: 0.5, Y: 0.5},
	})

	// Clamp radii to fit within rectangle
	maxRadius := math.Min(float64(r.bounds.Width/2), float64(r.bounds.Height/2))
	radiusX := math.Min(float64(r.radii.Width), maxRadius)
	radiusY := math.Min(float64(r.radii.Height), maxRadius)

	// Generate corner vertices
	corners := []geom.Point{
		{X: r.bounds.X + float32(radiusX), Y: r.bounds.Y + float32(radiusY)},                                    // Top-left
		{X: r.bounds.X + r.bounds.Width - float32(radiusX), Y: r.bounds.Y + float32(radiusY)},                   // Top-right
		{X: r.bounds.X + r.bounds.Width - float32(radiusX), Y: r.bounds.Y + r.bounds.Height - float32(radiusY)}, // Bottom-right
		{X: r.bounds.X + float32(radiusX), Y: r.bounds.Y + r.bounds.Height - float32(radiusY)},                  // Bottom-left
	}

	// Angle offsets for each corner
	angleOffsets := []float64{math.Pi, math.Pi / 2, 0, 3 * math.Pi / 2}

	// Generate vertices for each corner
	for cornerIdx, cornerCenter := range corners {
		startAngle := angleOffsets[cornerIdx]
		for i := 0; i <= r.segments; i++ {
			t := float64(i) / float64(r.segments)
			angle := startAngle + t*math.Pi/2 // 90 degrees per corner

			cos := float32(math.Cos(angle))
			sin := float32(math.Sin(angle))

			x := cornerCenter.X + cos*float32(radiusX)
			y := cornerCenter.Y + sin*float32(radiusY)

			// Normalized texture coordinates
			u := (x - r.bounds.X) / r.bounds.Width
			v := (y - r.bounds.Y) / r.bounds.Height

			vertices = append(vertices, render.Vertex{
				Position: geom.Point{X: x, Y: y},
				TexCoord: geom.Point{X: u, Y: v},
			})
		}
	}

	return vertices
}

// GetIndexBuffer generates index data
func (r *RoundRectGeometry) GetIndexBuffer() []uint16 {
	if r.IsAxisAlignedRect() {
		return []uint16{0, 1, 2, 2, 3, 0}
	}

	vertexCount := r.GetVertexCount()
	indices := make([]uint16, 0, r.GetIndexCount())

	// Triangle fan from center
	for i := 1; i < vertexCount-1; i++ {
		indices = append(indices, 0, uint16(i), uint16(i+1))
	}

	// Close the fan
	indices = append(indices, 0, uint16(vertexCount-1), 1)

	return indices
}

// Clone creates a copy of this geometry
func (r *RoundRectGeometry) Clone() Geometry {
	return &RoundRectGeometry{
		BaseGeometry: r.BaseGeometry.Clone().(BaseGeometry),
		bounds:       r.bounds,
		radii:        r.radii,
		segments:     r.segments,
		transform:    r.transform,
	}
}

// GetCoverage returns the coverage area for this geometry
func (r *RoundRectGeometry) GetCoverage(transform geom.Matrix) geom.Rect {
	combinedTransform := transform.Multiply(r.transform)
	return combinedTransform.TransformRect(r.bounds)
}
