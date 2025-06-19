package contents

// VerticesContents represents vertex-based rendering contents for custom geometry
type VerticesContents interface {
	Contents

	// SetVertices sets the vertex data
	SetVertices(vertices []Vertex)

	// GetVertices returns the vertex data
	GetVertices() []Vertex

	// SetIndices sets the index data
	SetIndices(indices []uint16)

	// GetIndices returns the index data
	GetIndices() []uint16

	// SetTexture sets the texture for the vertices
	SetTexture(texture Texture)

	// GetTexture returns the texture
	GetTexture() Texture

	// SetBlendMode sets the blend mode
	SetBlendMode(mode BlendMode)

	// GetBlendMode returns the blend mode
	GetBlendMode() BlendMode

	// SetColors sets per-vertex colors
	SetColors(colors []Color)

	// GetColors returns per-vertex colors
	GetColors() []Color
}

// Vertex represents a single vertex with position, texture coordinates, and color
type Vertex struct {
	Position Point
	TexCoord Point
	Color    Color
}

// VerticesContentsImpl implements VerticesContents
type VerticesContentsImpl struct {
	*ContentsImpl
	vertices  []Vertex
	indices   []uint16
	texture   Texture
	blendMode BlendMode
	colors    []Color
}

// NewVerticesContents creates new vertices contents
func NewVerticesContents() VerticesContents {
	return &VerticesContentsImpl{
		ContentsImpl: &ContentsImpl{
			isValid: true,
		},
		vertices:  make([]Vertex, 0),
		indices:   make([]uint16, 0),
		blendMode: BlendModeNormal,
		colors:    make([]Color, 0),
	}
}

// SetVertices sets the vertex data
func (v *VerticesContentsImpl) SetVertices(vertices []Vertex) {
	v.vertices = make([]Vertex, len(vertices))
	copy(v.vertices, vertices)
}

// GetVertices returns the vertex data
func (v *VerticesContentsImpl) GetVertices() []Vertex {
	return v.vertices
}

// SetIndices sets the index data
func (v *VerticesContentsImpl) SetIndices(indices []uint16) {
	v.indices = make([]uint16, len(indices))
	copy(v.indices, indices)
}

// GetIndices returns the index data
func (v *VerticesContentsImpl) GetIndices() []uint16 {
	return v.indices
}

// SetTexture sets the texture for the vertices
func (v *VerticesContentsImpl) SetTexture(texture Texture) {
	v.texture = texture
}

// GetTexture returns the texture
func (v *VerticesContentsImpl) GetTexture() Texture {
	return v.texture
}

// SetBlendMode sets the blend mode
func (v *VerticesContentsImpl) SetBlendMode(mode BlendMode) {
	v.blendMode = mode
}

// GetBlendMode returns the blend mode
func (v *VerticesContentsImpl) GetBlendMode() BlendMode {
	return v.blendMode
}

// SetColors sets per-vertex colors
func (v *VerticesContentsImpl) SetColors(colors []Color) {
	v.colors = make([]Color, len(colors))
	copy(v.colors, colors)
}

// GetColors returns per-vertex colors
func (v *VerticesContentsImpl) GetColors() []Color {
	return v.colors
}

// Render renders the vertices contents
func (v *VerticesContentsImpl) Render(context ContentContext, entity Entity, pass RenderPass) bool {
	if len(v.vertices) == 0 {
		return true
	}

	// TODO: Implement vertex rendering with custom geometry
	return false
}

// GetCoverage returns the coverage bounds of the vertices
func (v *VerticesContentsImpl) GetCoverage(entity Entity) Rect {
	if len(v.vertices) == 0 {
		return Rect{}
	}

	// Calculate bounding box from vertices
	minX, minY := v.vertices[0].Position.X, v.vertices[0].Position.Y
	maxX, maxY := minX, minY

	for _, vertex := range v.vertices {
		if vertex.Position.X < minX {
			minX = vertex.Position.X
		}
		if vertex.Position.X > maxX {
			maxX = vertex.Position.X
		}
		if vertex.Position.Y < minY {
			minY = vertex.Position.Y
		}
		if vertex.Position.Y > maxY {
			maxY = vertex.Position.Y
		}
	}

	return Rect{
		X:      minX,
		Y:      minY,
		Width:  maxX - minX,
		Height: maxY - minY,
	}
}

// Clone creates a copy of the vertices contents
func (v *VerticesContentsImpl) Clone() Contents {
	clone := &VerticesContentsImpl{
		ContentsImpl: v.ContentsImpl.Clone().(*ContentsImpl),
		texture:      v.texture,
		blendMode:    v.blendMode,
	}

	// Deep copy vertices
	clone.vertices = make([]Vertex, len(v.vertices))
	copy(clone.vertices, v.vertices)

	// Deep copy indices
	clone.indices = make([]uint16, len(v.indices))
	copy(clone.indices, v.indices)

	// Deep copy colors
	clone.colors = make([]Color, len(v.colors))
	copy(clone.colors, v.colors)

	return clone
}

// AddTriangle adds a triangle to the vertices
func (v *VerticesContentsImpl) AddTriangle(v1, v2, v3 Vertex) {
	baseIndex := uint16(len(v.vertices))

	v.vertices = append(v.vertices, v1, v2, v3)
	v.indices = append(v.indices, baseIndex, baseIndex+1, baseIndex+2)
}

// AddQuad adds a quad to the vertices
func (v *VerticesContentsImpl) AddQuad(v1, v2, v3, v4 Vertex) {
	baseIndex := uint16(len(v.vertices))

	v.vertices = append(v.vertices, v1, v2, v3, v4)
	v.indices = append(v.indices,
		baseIndex, baseIndex+1, baseIndex+2,
		baseIndex, baseIndex+2, baseIndex+3)
}

// Clear removes all vertices and indices
func (v *VerticesContentsImpl) Clear() {
	v.vertices = v.vertices[:0]
	v.indices = v.indices[:0]
	v.colors = v.colors[:0]
}

// GetTriangleCount returns the number of triangles
func (v *VerticesContentsImpl) GetTriangleCount() int {
	return len(v.indices) / 3
}

// GetVertexCount returns the number of vertices
func (v *VerticesContentsImpl) GetVertexCount() int {
	return len(v.vertices)
}
