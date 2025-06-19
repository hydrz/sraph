package geometry

import (
	"math"

	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// ArcGeometry represents arc geometry for drawing arcs and pie slices
type ArcGeometry struct {
	BaseGeometry
	ovalBounds    geom.Rect
	startAngle    float32 // Start angle in degrees
	sweepAngle    float32 // Sweep angle in degrees
	includeCenter bool    // Whether to include center point for pie slices
	strokeWidth   float32 // Stroke width (-1 for filled)
	cap           LineCap // Line cap style for stroked arcs
	transform     geom.Matrix
}

// NewArcGeometry creates a new arc geometry
func NewArcGeometry(ovalBounds geom.Rect, startAngle, sweepAngle float32, includeCenter bool) *ArcGeometry {
	return &ArcGeometry{
		BaseGeometry:  NewBaseGeometry(),
		ovalBounds:    ovalBounds,
		startAngle:    startAngle,
		sweepAngle:    sweepAngle,
		includeCenter: includeCenter,
		strokeWidth:   -1.0, // Filled by default
		cap:           LineCapButt,
		transform:     geom.NewIdentityMatrix(),
	}
}

// NewStrokedArcGeometry creates a new stroked arc geometry
func NewStrokedArcGeometry(ovalBounds geom.Rect, startAngle, sweepAngle float32, strokeWidth float32, cap LineCap) *ArcGeometry {
	return &ArcGeometry{
		BaseGeometry:  NewBaseGeometry(),
		ovalBounds:    ovalBounds,
		startAngle:    startAngle,
		sweepAngle:    sweepAngle,
		includeCenter: false, // Stroked arcs don't include center
		strokeWidth:   strokeWidth,
		cap:           cap,
		transform:     geom.NewIdentityMatrix(),
	}
}

// SetOvalBounds sets the bounding rectangle of the oval
func (a *ArcGeometry) SetOvalBounds(bounds geom.Rect) {
	a.ovalBounds = bounds
}

// GetOvalBounds returns the bounding rectangle of the oval
func (a *ArcGeometry) GetOvalBounds() geom.Rect {
	return a.ovalBounds
}

// SetStartAngle sets the start angle in degrees
func (a *ArcGeometry) SetStartAngle(angle float32) {
	a.startAngle = angle
}

// GetStartAngle returns the start angle in degrees
func (a *ArcGeometry) GetStartAngle() float32 {
	return a.startAngle
}

// SetSweepAngle sets the sweep angle in degrees
func (a *ArcGeometry) SetSweepAngle(angle float32) {
	a.sweepAngle = angle
}

// GetSweepAngle returns the sweep angle in degrees
func (a *ArcGeometry) GetSweepAngle() float32 {
	return a.sweepAngle
}

// SetIncludeCenter sets whether to include the center point
func (a *ArcGeometry) SetIncludeCenter(include bool) {
	a.includeCenter = include
}

// GetIncludeCenter returns whether the center point is included
func (a *ArcGeometry) GetIncludeCenter() bool {
	return a.includeCenter
}

// SetStrokeWidth sets the stroke width (-1 for filled)
func (a *ArcGeometry) SetStrokeWidth(width float32) {
	a.strokeWidth = width
}

// GetStrokeWidth returns the stroke width
func (a *ArcGeometry) GetStrokeWidth() float32 {
	return a.strokeWidth
}

// SetCap sets the line cap style
func (a *ArcGeometry) SetCap(cap LineCap) {
	a.cap = cap
}

// GetCap returns the line cap style
func (a *ArcGeometry) GetCap() LineCap {
	return a.cap
}

// SetTransform sets the transformation matrix
func (a *ArcGeometry) SetTransform(transform geom.Matrix) {
	a.transform = transform
}

// GetTransform returns the transformation matrix
func (a *ArcGeometry) GetTransform() geom.Matrix {
	return a.transform
}

// IsStroked returns true if this is a stroked arc
func (a *ArcGeometry) IsStroked() bool {
	return a.strokeWidth >= 0
}

// GetBounds returns the bounds of this geometry
func (a *ArcGeometry) GetBounds() geom.Rect {
	bounds := a.ovalBounds

	// If stroked, expand by stroke width
	if a.IsStroked() {
		halfStroke := a.strokeWidth / 2
		bounds = geom.Rect{
			X:      bounds.X - halfStroke,
			Y:      bounds.Y - halfStroke,
			Width:  bounds.Width + a.strokeWidth,
			Height: bounds.Height + a.strokeWidth,
		}
	}

	return a.transform.TransformRect(bounds)
}

// GetVertexCount returns the estimated number of vertices
func (a *ArcGeometry) GetVertexCount() int {
	// Calculate segments based on sweep angle
	segments := int(math.Abs(float64(a.sweepAngle)) / 5.0) // 5 degrees per segment
	if segments < 3 {
		segments = 3
	}
	if segments > 360 {
		segments = 360
	}

	if a.includeCenter {
		return segments + 1 // Center vertex + arc vertices
	}
	return segments
}

// GetIndexCount returns the estimated number of indices
func (a *ArcGeometry) GetIndexCount() int {
	segments := a.GetVertexCount()
	if a.includeCenter {
		return (segments - 1) * 3 // Triangles from center
	}
	return (segments - 2) * 3 // Triangle fan
}

// ComputeVertices generates vertex data for rendering
func (a *ArcGeometry) ComputeVertices() []render.Vertex {
	segments := a.GetVertexCount()
	vertices := make([]render.Vertex, 0, segments)

	// Calculate arc parameters
	centerX := a.ovalBounds.X + a.ovalBounds.Width/2
	centerY := a.ovalBounds.Y + a.ovalBounds.Height/2
	radiusX := a.ovalBounds.Width / 2
	radiusY := a.ovalBounds.Height / 2

	startRad := float64(a.startAngle) * math.Pi / 180.0
	sweepRad := float64(a.sweepAngle) * math.Pi / 180.0

	// Add center vertex if needed
	if a.includeCenter {
		vertices = append(vertices, render.Vertex{
			Position: geom.Point{X: centerX, Y: centerY},
			TexCoord: geom.Point{X: 0.5, Y: 0.5},
		})
		segments-- // Reduce arc segments count
	}

	// Generate arc vertices
	for i := 0; i <= segments; i++ {
		t := float64(i) / float64(segments)
		angle := startRad + t*sweepRad

		cos := float32(math.Cos(angle))
		sin := float32(math.Sin(angle))

		x := centerX + cos*radiusX
		y := centerY + sin*radiusY

		// Texture coordinates for arc (0.5 + normalized position)
		u := 0.5 + cos*0.5
		v := 0.5 + sin*0.5

		vertices = append(vertices, render.Vertex{
			Position: geom.Point{X: x, Y: y},
			TexCoord: geom.Point{X: u, Y: v},
		})
	}

	return vertices
}

// GetIndexBuffer generates index data
func (a *ArcGeometry) GetIndexBuffer() []uint16 {
	vertices := a.GetVertexCount()
	var indices []uint16

	if a.includeCenter {
		// Triangle fan from center
		for i := 1; i < vertices-1; i++ {
			indices = append(indices, 0, uint16(i), uint16(i+1))
		}
	} else {
		// Triangle fan without center
		for i := 1; i < vertices-2; i++ {
			indices = append(indices, 0, uint16(i), uint16(i+1))
		}
	}

	return indices
}

// Clone creates a copy of this geometry
func (a *ArcGeometry) Clone() Geometry {
	return &ArcGeometry{
		BaseGeometry:  a.BaseGeometry.Clone().(BaseGeometry),
		ovalBounds:    a.ovalBounds,
		startAngle:    a.startAngle,
		sweepAngle:    a.sweepAngle,
		includeCenter: a.includeCenter,
		strokeWidth:   a.strokeWidth,
		cap:           a.cap,
		transform:     a.transform,
	}
}

// GetCoverage returns the coverage area for this geometry
func (a *ArcGeometry) GetCoverage(transform geom.Matrix) geom.Rect {
	combinedTransform := transform.Multiply(a.transform)
	return combinedTransform.TransformRect(a.ovalBounds)
}
