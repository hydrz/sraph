package geometry

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// PointFieldGeometry represents a geometry that renders a field of points
// Used for particle systems, point clouds, or other point-based graphics
type PointFieldGeometry struct {
	baseGeometry
	Points     []geom.Point
	PointSizes []float32
	Colors     []geom.Color
	Texture    *render.Texture
}

// NewPointFieldGeometry creates a new point field geometry with the given points
func NewPointFieldGeometry(points []geom.Point) *PointFieldGeometry {
	return &PointFieldGeometry{
		Points: points,
	}
}

// WithPointSizes sets custom sizes for each point
func (g *PointFieldGeometry) WithPointSizes(sizes []float32) *PointFieldGeometry {
	g.PointSizes = sizes
	return g
}

// WithColors sets custom colors for each point
func (g *PointFieldGeometry) WithColors(colors []geom.Color) *PointFieldGeometry {
	g.Colors = colors
	return g
}

// WithTexture sets a texture to be used for point sprites
func (g *PointFieldGeometry) WithTexture(texture *render.Texture) *PointFieldGeometry {
	g.Texture = texture
	return g
}

// GetPositionBuffer generates the position buffer for rendering
func (g *PointFieldGeometry) GetPositionBuffer(transform geom.Matrix, offset geom.Point) []float32 {
	// TODO: Implement position buffer generation for point field
	return nil
}

// GetBounds returns the axis-aligned bounding box of all points
func (g *PointFieldGeometry) GetBounds() geom.Rect {
	if len(g.Points) == 0 {
		return geom.Rect{}
	}

	minX, minY := g.Points[0].X, g.Points[0].Y
	maxX, maxY := g.Points[0].X, g.Points[0].Y

	for _, point := range g.Points[1:] {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	return geom.Rect{
		Origin: geom.Point{X: minX, Y: minY},
		Size:   geom.Size{Width: maxX - minX, Height: maxY - minY},
	}
}

// ComputeUVs generates UV coordinates for point sprites if texture is used
func (g *PointFieldGeometry) ComputeUVs(transform geom.Matrix) ([]geom.Point, bool) {
	// TODO: Implement UV computation for point field
	return nil, false
}

// GetIndexBuffer returns the index buffer for rendering
func (g *PointFieldGeometry) GetIndexBuffer() []uint16 {
	// Point field uses direct point rendering, no indices needed
	return nil
}
