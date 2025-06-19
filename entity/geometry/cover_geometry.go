package geometry

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// CoverGeometry represents geometry that covers the entire render pass area
type CoverGeometry struct {
	BaseGeometry
	transform geom.Matrix
}

// NewCoverGeometry creates a new cover geometry
func NewCoverGeometry() *CoverGeometry {
	return &CoverGeometry{
		BaseGeometry: NewBaseGeometry(),
		transform:    geom.NewIdentityMatrix(),
	}
}

// SetTransform sets the transformation matrix
func (c *CoverGeometry) SetTransform(transform geom.Matrix) {
	c.transform = transform
}

// GetTransform returns the transformation matrix
func (c *CoverGeometry) GetTransform() geom.Matrix {
	return c.transform
}

// GetBounds returns the bounds of this geometry (covers entire area)
func (c *CoverGeometry) GetBounds() geom.Rect {
	// Cover geometry typically covers the entire render target
	// The actual bounds are determined at render time
	return geom.Rect{
		X:      -1e6, // Very large bounds
		Y:      -1e6,
		Width:  2e6,
		Height: 2e6,
	}
}

// GetVertexCount returns the number of vertices (quad)
func (c *CoverGeometry) GetVertexCount() int {
	return 4
}

// GetIndexCount returns the number of indices (two triangles)
func (c *CoverGeometry) GetIndexCount() int {
	return 6
}

// ComputeVertices generates vertex data for a full-screen quad
func (c *CoverGeometry) ComputeVertices() []render.Vertex {
	// Full-screen quad in normalized device coordinates
	return []render.Vertex{
		{
			Position: geom.Point{X: -1, Y: -1},
			TexCoord: geom.Point{X: 0, Y: 0},
		},
		{
			Position: geom.Point{X: 1, Y: -1},
			TexCoord: geom.Point{X: 1, Y: 0},
		},
		{
			Position: geom.Point{X: 1, Y: 1},
			TexCoord: geom.Point{X: 1, Y: 1},
		},
		{
			Position: geom.Point{X: -1, Y: 1},
			TexCoord: geom.Point{X: 0, Y: 1},
		},
	}
}

// ComputeVerticesForRenderTarget generates vertices for specific render target size
func (c *CoverGeometry) ComputeVerticesForRenderTarget(targetSize geom.Size) []render.Vertex {
	return []render.Vertex{
		{
			Position: geom.Point{X: 0, Y: 0},
			TexCoord: geom.Point{X: 0, Y: 0},
		},
		{
			Position: geom.Point{X: targetSize.Width, Y: 0},
			TexCoord: geom.Point{X: 1, Y: 0},
		},
		{
			Position: geom.Point{X: targetSize.Width, Y: targetSize.Height},
			TexCoord: geom.Point{X: 1, Y: 1},
		},
		{
			Position: geom.Point{X: 0, Y: targetSize.Height},
			TexCoord: geom.Point{X: 0, Y: 1},
		},
	}
}

// GetIndexBuffer generates index data
func (c *CoverGeometry) GetIndexBuffer() []uint16 {
	return []uint16{
		0, 1, 2, // First triangle
		2, 3, 0, // Second triangle
	}
}

// CoversArea returns true since this geometry covers any area
func (c *CoverGeometry) CoversArea(transform geom.Matrix, rect geom.Rect) bool {
	return true
}

// CanApplyMaskFilter returns whether mask filters can be applied
func (c *CoverGeometry) CanApplyMaskFilter() bool {
	return false // Cover geometry typically doesn't support mask filters
}

// Clone creates a copy of this geometry
func (c *CoverGeometry) Clone() Geometry {
	return &CoverGeometry{
		BaseGeometry: c.BaseGeometry.Clone().(BaseGeometry),
		transform:    c.transform,
	}
}

// GetCoverage returns the coverage area for this geometry
func (c *CoverGeometry) GetCoverage(transform geom.Matrix) geom.Rect {
	// Cover geometry covers everything
	return c.GetBounds()
}
