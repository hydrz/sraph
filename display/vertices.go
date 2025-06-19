package display

import (
	"fmt"

	"github.com/opensraph/sraph/geom"
)

// VertexMode defines how vertices are connected to form primitives.
type VertexMode uint8

const (
	// VertexModeTriangles treats each group of 3 vertices as a separate triangle.
	VertexModeTriangles VertexMode = iota
	// VertexModeTriangleStrip connects vertices as a triangle strip.
	VertexModeTriangleStrip
	// VertexModeTriangleFan connects vertices as a triangle fan.
	VertexModeTriangleFan
)

// Vertices represents a collection of vertices for rendering.
type Vertices struct {
	mode          VertexMode
	positions     []geom.Point[Scalar]
	textureCoords []geom.Point[Scalar]
	colors        []Color
	indices       []uint16
}

// VerticesBuilder helps construct Vertices objects.
type VerticesBuilder struct {
	mode          VertexMode
	positions     []geom.Point[Scalar]
	textureCoords []geom.Point[Scalar]
	colors        []Color
	indices       []uint16
	hasTexCoords  bool
	hasColors     bool
	hasIndices    bool
}

// NewVerticesBuilder creates a new VerticesBuilder.
func NewVerticesBuilder(mode VertexMode, vertexCount int) *VerticesBuilder {
	return &VerticesBuilder{
		mode:      mode,
		positions: make([]geom.Point[Scalar], 0, vertexCount),
	}
}

// WithTextureCoords enables texture coordinates for the vertices.
func (b *VerticesBuilder) WithTextureCoords() *VerticesBuilder {
	if !b.hasTexCoords {
		b.hasTexCoords = true
		b.textureCoords = make([]geom.Point[Scalar], 0, cap(b.positions))
	}
	return b
}

// WithColors enables colors for the vertices.
func (b *VerticesBuilder) WithColors() *VerticesBuilder {
	if !b.hasColors {
		b.hasColors = true
		b.colors = make([]Color, 0, cap(b.positions))
	}
	return b
}

// WithIndices enables indices for the vertices.
func (b *VerticesBuilder) WithIndices(indexCount int) *VerticesBuilder {
	if !b.hasIndices {
		b.hasIndices = true
		b.indices = make([]uint16, 0, indexCount)
	}
	return b
}

// AddVertex adds a vertex with position only.
func (b *VerticesBuilder) AddVertex(position geom.Point[Scalar]) *VerticesBuilder {
	b.positions = append(b.positions, position)

	// Add default values for optional components
	if b.hasTexCoords {
		b.textureCoords = append(b.textureCoords, geom.Point[Scalar]{})
	}
	if b.hasColors {
		b.colors = append(b.colors, ColorWhite)
	}

	return b
}

// AddVertexWithTexCoord adds a vertex with position and texture coordinate.
func (b *VerticesBuilder) AddVertexWithTexCoord(position, texCoord geom.Point[Scalar]) *VerticesBuilder {
	b.positions = append(b.positions, position)

	if b.hasTexCoords {
		b.textureCoords = append(b.textureCoords, texCoord)
	}
	if b.hasColors {
		b.colors = append(b.colors, ColorWhite)
	}

	return b
}

// AddVertexWithColor adds a vertex with position and color.
func (b *VerticesBuilder) AddVertexWithColor(position geom.Point[Scalar], color Color) *VerticesBuilder {
	b.positions = append(b.positions, position)

	if b.hasTexCoords {
		b.textureCoords = append(b.textureCoords, geom.Point[Scalar]{})
	}
	if b.hasColors {
		b.colors = append(b.colors, color)
	}

	return b
}

// AddVertexFull adds a vertex with position, texture coordinate, and color.
func (b *VerticesBuilder) AddVertexFull(position, texCoord geom.Point[Scalar], color Color) *VerticesBuilder {
	b.positions = append(b.positions, position)

	if b.hasTexCoords {
		b.textureCoords = append(b.textureCoords, texCoord)
	}
	if b.hasColors {
		b.colors = append(b.colors, color)
	}

	return b
}

// AddIndex adds an index to the index buffer.
func (b *VerticesBuilder) AddIndex(index uint16) *VerticesBuilder {
	if b.hasIndices {
		b.indices = append(b.indices, index)
	}
	return b
}

// AddTriangle adds a triangle using indices.
func (b *VerticesBuilder) AddTriangle(i0, i1, i2 uint16) *VerticesBuilder {
	if b.hasIndices {
		b.indices = append(b.indices, i0, i1, i2)
	}
	return b
}

// SetColor sets the color for the vertex at the specified index.
func (b *VerticesBuilder) SetColor(index int, color Color) *VerticesBuilder {
	if b.hasColors && index >= 0 && index < len(b.colors) {
		b.colors[index] = color
	}
	return b
}

// Build creates the final Vertices object.
func (b *VerticesBuilder) Build() (*Vertices, error) {
	if len(b.positions) == 0 {
		return nil, fmt.Errorf("vertices must have at least one position")
	}

	vertices := &Vertices{
		mode:      b.mode,
		positions: make([]geom.Point[Scalar], len(b.positions)),
	}

	copy(vertices.positions, b.positions)

	if b.hasTexCoords && len(b.textureCoords) == len(b.positions) {
		vertices.textureCoords = make([]geom.Point[Scalar], len(b.textureCoords))
		copy(vertices.textureCoords, b.textureCoords)
	}

	if b.hasColors && len(b.colors) == len(b.positions) {
		vertices.colors = make([]Color, len(b.colors))
		copy(vertices.colors, b.colors)
	}

	if b.hasIndices && len(b.indices) > 0 {
		vertices.indices = make([]uint16, len(b.indices))
		copy(vertices.indices, b.indices)
	}

	return vertices, nil
}

// Mode returns the vertex mode.
func (v *Vertices) Mode() VertexMode {
	return v.mode
}

// PositionCount returns the number of position vertices.
func (v *Vertices) PositionCount() int {
	return len(v.positions)
}

// VertexCount returns the number of position vertices (alias for PositionCount).
func (v *Vertices) VertexCount() int {
	return len(v.positions)
}

// Positions returns the position vertices.
func (v *Vertices) Positions() []geom.Point[Scalar] {
	return v.positions
}

// HasTextureCoords returns true if the vertices have texture coordinates.
func (v *Vertices) HasTextureCoords() bool {
	return len(v.textureCoords) > 0
}

// TextureCoords returns the texture coordinate vertices.
func (v *Vertices) TextureCoords() []geom.Point[Scalar] {
	return v.textureCoords
}

// HasColors returns true if the vertices have colors.
func (v *Vertices) HasColors() bool {
	return len(v.colors) > 0
}

// Colors returns the color vertices.
func (v *Vertices) Colors() []Color {
	return v.colors
}

// HasIndices returns true if the vertices have indices.
func (v *Vertices) HasIndices() bool {
	return len(v.indices) > 0
}

// IndexCount returns the number of indices.
func (v *Vertices) IndexCount() int {
	return len(v.indices)
}

// Indices returns the indices.
func (v *Vertices) Indices() []uint16 {
	return v.indices
}

// String returns a string representation of the vertices.
func (v *Vertices) String() string {
	return fmt.Sprintf("Vertices{Mode: %d, Positions: %d, TexCoords: %t, Colors: %t, Indices: %d}",
		v.mode, len(v.positions), v.HasTextureCoords(), v.HasColors(), len(v.indices))
}

// Bounds returns the bounding rectangle of all vertex positions.
func (v *Vertices) Bounds() *geom.Rect[Scalar] {
	if len(v.positions) == 0 {
		return nil
	}

	min := v.positions[0]
	max := v.positions[0]

	for _, pos := range v.positions[1:] {
		if pos.X < min.X {
			min.X = pos.X
		}
		if pos.Y < min.Y {
			min.Y = pos.Y
		}
		if pos.X > max.X {
			max.X = pos.X
		}
		if pos.Y > max.Y {
			max.Y = pos.Y
		}
	}

	bounds := geom.NewRect(min.X, min.Y, max.X, max.Y)
	return &bounds
}
