package render

// VertexBufferBuilder helps build vertex buffers efficiently
type VertexBufferBuilder interface {
	// SetVertexCount sets the number of vertices
	SetVertexCount(count int)

	// AddAttribute adds a vertex attribute
	AddAttribute(name string, format VertexFormat, offset int)

	// SetStride sets the vertex stride
	SetStride(stride int)

	// SetData sets the vertex data
	SetData(data []byte)

	// Build builds the vertex buffer
	Build() (Buffer, error)

	// Reset resets the builder state
	Reset()
}

// VertexBufferBuilderImpl is the default implementation of VertexBufferBuilder
type VertexBufferBuilderImpl struct {
	vertexCount int
	attributes  []VertexAttribute
	stride      int
	data        []byte
	context     Context
}

// NewVertexBufferBuilder creates a new vertex buffer builder
func NewVertexBufferBuilder(context Context) VertexBufferBuilder {
	return &VertexBufferBuilderImpl{
		context:    context,
		attributes: make([]VertexAttribute, 0),
	}
}

// SetVertexCount sets the number of vertices
func (b *VertexBufferBuilderImpl) SetVertexCount(count int) {
	b.vertexCount = count
}

// AddAttribute adds a vertex attribute
func (b *VertexBufferBuilderImpl) AddAttribute(name string, format VertexFormat, offset int) {
	attribute := VertexAttribute{
		Location: len(b.attributes),
		Format:   format,
		Offset:   offset,
	}
	b.attributes = append(b.attributes, attribute)
}

// SetStride sets the vertex stride
func (b *VertexBufferBuilderImpl) SetStride(stride int) {
	b.stride = stride
}

// SetData sets the vertex data
func (b *VertexBufferBuilderImpl) SetData(data []byte) {
	b.data = make([]byte, len(data))
	copy(b.data, data)
}

// Build builds the vertex buffer
func (b *VertexBufferBuilderImpl) Build() (Buffer, error) {
	if b.context == nil {
		return nil, ErrInvalidState
	}

	if len(b.data) == 0 {
		return nil, ErrInvalidArgument
	}

	// Calculate expected data size
	expectedSize := b.vertexCount * b.stride
	if len(b.data) < expectedSize {
		return nil, ErrInvalidArgument
	}

	// Create buffer with vertex usage
	buffer, err := b.context.CreateBuffer(len(b.data), BufferUsageVertex)
	if err != nil {
		return nil, err
	}

	// Copy data to buffer
	data, err := buffer.Map()
	if err != nil {
		return nil, err
	}

	copy(data, b.data)

	err = buffer.Unmap()
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

// Reset resets the builder state
func (b *VertexBufferBuilderImpl) Reset() {
	b.vertexCount = 0
	b.attributes = b.attributes[:0]
	b.stride = 0
	b.data = nil
}

// IndexBufferBuilder helps build index buffers efficiently
type IndexBufferBuilder interface {
	// SetIndexCount sets the number of indices
	SetIndexCount(count int)

	// SetIndexFormat sets the index format (16-bit or 32-bit)
	SetIndexFormat(format IndexFormat)

	// SetData sets the index data
	SetData(data []byte)

	// SetIndices16 sets 16-bit index data
	SetIndices16(indices []uint16)

	// SetIndices32 sets 32-bit index data
	SetIndices32(indices []uint32)

	// Build builds the index buffer
	Build() (Buffer, error)

	// Reset resets the builder state
	Reset()
}

// IndexFormat defines the format of index data
type IndexFormat int

const (
	IndexFormat16 IndexFormat = iota
	IndexFormat32
)

// IndexBufferBuilderImpl is the default implementation of IndexBufferBuilder
type IndexBufferBuilderImpl struct {
	indexCount  int
	indexFormat IndexFormat
	data        []byte
	context     Context
}

// NewIndexBufferBuilder creates a new index buffer builder
func NewIndexBufferBuilder(context Context) IndexBufferBuilder {
	return &IndexBufferBuilderImpl{
		context:     context,
		indexFormat: IndexFormat16,
	}
}

// SetIndexCount sets the number of indices
func (b *IndexBufferBuilderImpl) SetIndexCount(count int) {
	b.indexCount = count
}

// SetIndexFormat sets the index format (16-bit or 32-bit)
func (b *IndexBufferBuilderImpl) SetIndexFormat(format IndexFormat) {
	b.indexFormat = format
}

// SetData sets the index data
func (b *IndexBufferBuilderImpl) SetData(data []byte) {
	b.data = make([]byte, len(data))
	copy(b.data, data)
}

// SetIndices16 sets 16-bit index data
func (b *IndexBufferBuilderImpl) SetIndices16(indices []uint16) {
	b.indexFormat = IndexFormat16
	b.indexCount = len(indices)

	// Convert to bytes
	b.data = make([]byte, len(indices)*2)
	for i, index := range indices {
		b.data[i*2] = byte(index)
		b.data[i*2+1] = byte(index >> 8)
	}
}

// SetIndices32 sets 32-bit index data
func (b *IndexBufferBuilderImpl) SetIndices32(indices []uint32) {
	b.indexFormat = IndexFormat32
	b.indexCount = len(indices)

	// Convert to bytes
	b.data = make([]byte, len(indices)*4)
	for i, index := range indices {
		b.data[i*4] = byte(index)
		b.data[i*4+1] = byte(index >> 8)
		b.data[i*4+2] = byte(index >> 16)
		b.data[i*4+3] = byte(index >> 24)
	}
}

// Build builds the index buffer
func (b *IndexBufferBuilderImpl) Build() (Buffer, error) {
	if b.context == nil {
		return nil, ErrInvalidState
	}

	if len(b.data) == 0 {
		return nil, ErrInvalidArgument
	}

	// Create buffer with index usage
	buffer, err := b.context.CreateBuffer(len(b.data), BufferUsageIndex)
	if err != nil {
		return nil, err
	}

	// Copy data to buffer
	data, err := buffer.Map()
	if err != nil {
		return nil, err
	}

	copy(data, b.data)

	err = buffer.Unmap()
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

// Reset resets the builder state
func (b *IndexBufferBuilderImpl) Reset() {
	b.indexCount = 0
	b.indexFormat = IndexFormat16
	b.data = nil
}
