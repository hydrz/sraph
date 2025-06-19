package geometry

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// LineGeometry represents line geometry
type LineGeometry struct {
	BaseGeometry
	startPoint geom.Point
	endPoint   geom.Point
	width      float32
	cap        LineCap
	transform  geom.Matrix
}

// LineCap defines the end cap style for lines
type LineCap int

const (
	// LineCapButt flat end cap
	LineCapButt LineCap = iota
	// LineCapRound rounded end cap
	LineCapRound
	// LineCapSquare square end cap
	LineCapSquare
)

// NewLineGeometry creates a new line geometry
func NewLineGeometry(start, end geom.Point, width float32) *LineGeometry {
	return &LineGeometry{
		BaseGeometry: NewBaseGeometry(),
		startPoint:   start,
		endPoint:     end,
		width:        width,
		cap:          LineCapButt,
		transform:    geom.NewIdentityMatrix(),
	}
}

// SetStartPoint sets the line start point
func (l *LineGeometry) SetStartPoint(point geom.Point) {
	l.startPoint = point
}

// GetStartPoint returns the line start point
func (l *LineGeometry) GetStartPoint() geom.Point {
	return l.startPoint
}

// SetEndPoint sets the line end point
func (l *LineGeometry) SetEndPoint(point geom.Point) {
	l.endPoint = point
}

// GetEndPoint returns the line end point
func (l *LineGeometry) GetEndPoint() geom.Point {
	return l.endPoint
}

// SetWidth sets the line width
func (l *LineGeometry) SetWidth(width float32) {
	l.width = width
}

// GetWidth returns the line width
func (l *LineGeometry) GetWidth() float32 {
	return l.width
}

// SetCap sets the line cap style
func (l *LineGeometry) SetCap(cap LineCap) {
	l.cap = cap
}

// GetCap returns the line cap style
func (l *LineGeometry) GetCap() LineCap {
	return l.cap
}

// SetTransform sets the transformation matrix
func (l *LineGeometry) SetTransform(transform geom.Matrix) {
	l.transform = transform
}

// GetTransform returns the current transformation matrix
func (l *LineGeometry) GetTransform() geom.Matrix {
	return l.transform
}

// GetBounds returns the bounds of this geometry
func (l *LineGeometry) GetBounds() geom.Rect {
	// Calculate line bounds including width
	minX := l.startPoint.X
	maxX := l.endPoint.X
	minY := l.startPoint.Y
	maxY := l.endPoint.Y

	if minX > maxX {
		minX, maxX = maxX, minX
	}
	if minY > maxY {
		minY, maxY = maxY, minY
	}

	// Expand by half width
	halfWidth := l.width / 2
	rect := geom.Rect{
		X:      minX - halfWidth,
		Y:      minY - halfWidth,
		Width:  maxX - minX + l.width,
		Height: maxY - minY + l.width,
	}

	return l.transform.TransformRect(rect)
}

// GetVertexCount returns the number of vertices
func (l *LineGeometry) GetVertexCount() int {
	// Line rendered as quad
	return 4
}

// GetIndexCount returns the number of indices
func (l *LineGeometry) GetIndexCount() int {
	// Two triangles for the quad
	return 6
}

// ComputeVertices generates vertex data for rendering
func (l *LineGeometry) ComputeVertices() []render.Vertex {
	// Calculate line direction and perpendicular
	dx := l.endPoint.X - l.startPoint.X
	dy := l.endPoint.Y - l.startPoint.Y
	length := geom.Sqrt(dx*dx + dy*dy)

	if length == 0 {
		// Degenerate line, return empty vertices
		return []render.Vertex{}
	}

	// Normalize direction
	dx /= length
	dy /= length

	// Perpendicular direction
	perpX := -dy
	perpY := dx

	// Half width offset
	halfWidth := l.width / 2
	offsetX := perpX * halfWidth
	offsetY := perpY * halfWidth

	vertices := []render.Vertex{
		{
			Position: geom.Point{
				X: l.startPoint.X + offsetX,
				Y: l.startPoint.Y + offsetY,
			},
			TexCoord: geom.Point{X: 0, Y: 0},
		},
		{
			Position: geom.Point{
				X: l.startPoint.X - offsetX,
				Y: l.startPoint.Y - offsetY,
			},
			TexCoord: geom.Point{X: 0, Y: 1},
		},
		{
			Position: geom.Point{
				X: l.endPoint.X - offsetX,
				Y: l.endPoint.Y - offsetY,
			},
			TexCoord: geom.Point{X: 1, Y: 1},
		},
		{
			Position: geom.Point{
				X: l.endPoint.X + offsetX,
				Y: l.endPoint.Y + offsetY,
			},
			TexCoord: geom.Point{X: 1, Y: 0},
		},
	}

	return vertices
}

// GetIndexBuffer generates index data
func (l *LineGeometry) GetIndexBuffer() []uint16 {
	return []uint16{
		0, 1, 2, // First triangle
		2, 3, 0, // Second triangle
	}
}

// Clone creates a copy of this geometry
func (l *LineGeometry) Clone() Geometry {
	return &LineGeometry{
		BaseGeometry: l.BaseGeometry.Clone().(BaseGeometry),
		startPoint:   l.startPoint,
		endPoint:     l.endPoint,
		width:        l.width,
		cap:          l.cap,
		transform:    l.transform,
	}
}

// EllipseGeometry represents elliptical geometry
type EllipseGeometry struct {
	BaseGeometry
	center    geom.Point
	radiusX   float32
	radiusY   float32
	segments  int
	transform geom.Matrix
}

// NewEllipseGeometry creates a new elliptical geometry
func NewEllipseGeometry(center geom.Point, radiusX, radiusY float32, segments int) *EllipseGeometry {
	if segments < 3 {
		segments = 32
	}

	return &EllipseGeometry{
		BaseGeometry: NewBaseGeometry(),
		center:       center,
		radiusX:      radiusX,
		radiusY:      radiusY,
		segments:     segments,
		transform:    geom.NewIdentityMatrix(),
	}
}

// SetCenter sets the ellipse center
func (e *EllipseGeometry) SetCenter(center geom.Point) {
	e.center = center
}

// GetCenter returns the ellipse center
func (e *EllipseGeometry) GetCenter() geom.Point {
	return e.center
}

// SetRadiusX sets the X radius
func (e *EllipseGeometry) SetRadiusX(radius float32) {
	e.radiusX = radius
}

// GetRadiusX returns the X radius
func (e *EllipseGeometry) GetRadiusX() float32 {
	return e.radiusX
}

// SetRadiusY sets the Y radius
func (e *EllipseGeometry) SetRadiusY(radius float32) {
	e.radiusY = radius
}

// GetRadiusY returns the Y radius
func (e *EllipseGeometry) GetRadiusY() float32 {
	return e.radiusY
}

// SetSegments sets the number of segments
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

// GetTransform returns the current transformation matrix
func (e *EllipseGeometry) GetTransform() geom.Matrix {
	return e.transform
}

// GetBounds returns the bounds of this geometry
func (e *EllipseGeometry) GetBounds() geom.Rect {
	rect := geom.Rect{
		X:      e.center.X - e.radiusX,
		Y:      e.center.Y - e.radiusY,
		Width:  e.radiusX * 2,
		Height: e.radiusY * 2,
	}
	return e.transform.TransformRect(rect)
}

// GetVertexCount returns the number of vertices
func (e *EllipseGeometry) GetVertexCount() int {
	return e.segments + 1 // Center vertex + perimeter vertices
}

// GetIndexCount returns the number of indices
func (e *EllipseGeometry) GetIndexCount() int {
	return e.segments * 3 // Each segment forms a triangle with center
}

// ComputeVertices generates vertex data for rendering
func (e *EllipseGeometry) ComputeVertices() []render.Vertex {
	vertices := make([]render.Vertex, e.GetVertexCount())

	// Center vertex
	vertices[0] = render.Vertex{
		Position: e.center,
		TexCoord: geom.Point{X: 0.5, Y: 0.5},
	}

	// Perimeter vertices
	for i := 0; i < e.segments; i++ {
		angle := float32(i) * 2.0 * geom.Pi / float32(e.segments)
		cos := geom.Cos(angle)
		sin := geom.Sin(angle)

		vertices[i+1] = render.Vertex{
			Position: geom.Point{
				X: e.center.X + cos*e.radiusX,
				Y: e.center.Y + sin*e.radiusY,
			},
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
		indices[base+1] = uint16(i + 1)                // Current perimeter vertex
		indices[base+2] = uint16((i+1)%e.segments + 1) // Next perimeter vertex
	}

	return indices
}

// Clone creates a copy of this geometry
func (e *EllipseGeometry) Clone() Geometry {
	return &EllipseGeometry{
		BaseGeometry: e.BaseGeometry.Clone().(BaseGeometry),
		center:       e.center,
		radiusX:      e.radiusX,
		radiusY:      e.radiusY,
		segments:     e.segments,
		transform:    e.transform,
	}
}
