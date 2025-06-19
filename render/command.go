package render

// Command represents a generic render command
type Command interface {
	// Execute executes the command
	Execute() error

	// IsValid returns true if the command is valid
	IsValid() bool

	// GetType returns the command type
	GetType() CommandType

	// GetLabel returns the debug label for this command
	GetLabel() string

	// SetLabel sets a debug label for this command
	SetLabel(label string)
}

// CommandType defines the type of render command
type CommandType int

const (
	CommandTypeUnknown CommandType = iota
	CommandTypeDrawIndexed
	CommandTypeDraw
	CommandTypeDispatch
	CommandTypeCopyBuffer
	CommandTypeCopyTexture
	CommandTypeGenerateMips
	CommandTypeBeginRenderPass
	CommandTypeEndRenderPass
)

// CommandImpl is the base implementation for all commands
type CommandImpl struct {
	commandType CommandType
	label       string
	isValid     bool
}

// Execute executes the command (base implementation)
func (c *CommandImpl) Execute() error {
	// Base implementation does nothing
	return nil
}

// IsValid returns true if the command is valid
func (c *CommandImpl) IsValid() bool {
	return c.isValid
}

// GetType returns the command type
func (c *CommandImpl) GetType() CommandType {
	return c.commandType
}

// GetLabel returns the debug label for this command
func (c *CommandImpl) GetLabel() string {
	return c.label
}

// SetLabel sets a debug label for this command
func (c *CommandImpl) SetLabel(label string) {
	c.label = label
}

// DrawIndexedCommand represents an indexed draw command
type DrawIndexedCommand struct {
	*CommandImpl
	indexCount    int
	instanceCount int
	firstIndex    int
	baseVertex    int
	firstInstance int
}

// NewDrawIndexedCommand creates a new indexed draw command
func NewDrawIndexedCommand(indexCount, instanceCount, firstIndex, baseVertex, firstInstance int) *DrawIndexedCommand {
	return &DrawIndexedCommand{
		CommandImpl: &CommandImpl{
			commandType: CommandTypeDrawIndexed,
			isValid:     true,
		},
		indexCount:    indexCount,
		instanceCount: instanceCount,
		firstIndex:    firstIndex,
		baseVertex:    baseVertex,
		firstInstance: firstInstance,
	}
}

// Execute executes the indexed draw command
func (c *DrawIndexedCommand) Execute() error {
	// TODO: Implement indexed drawing
	return nil
}

// DrawCommand represents a draw command
type DrawCommand struct {
	*CommandImpl
	vertexCount   int
	instanceCount int
	firstVertex   int
	firstInstance int
}

// NewDrawCommand creates a new draw command
func NewDrawCommand(vertexCount, instanceCount, firstVertex, firstInstance int) *DrawCommand {
	return &DrawCommand{
		CommandImpl: &CommandImpl{
			commandType: CommandTypeDraw,
			isValid:     true,
		},
		vertexCount:   vertexCount,
		instanceCount: instanceCount,
		firstVertex:   firstVertex,
		firstInstance: firstInstance,
	}
}

// Execute executes the draw command
func (c *DrawCommand) Execute() error {
	// TODO: Implement drawing
	return nil
}

// DispatchCommand represents a compute dispatch command
type DispatchCommand struct {
	*CommandImpl
	groupCountX int
	groupCountY int
	groupCountZ int
}

// NewDispatchCommand creates a new dispatch command
func NewDispatchCommand(groupCountX, groupCountY, groupCountZ int) *DispatchCommand {
	return &DispatchCommand{
		CommandImpl: &CommandImpl{
			commandType: CommandTypeDispatch,
			isValid:     true,
		},
		groupCountX: groupCountX,
		groupCountY: groupCountY,
		groupCountZ: groupCountZ,
	}
}

// Execute executes the dispatch command
func (c *DispatchCommand) Execute() error {
	// TODO: Implement compute dispatch
	return nil
}

// CopyBufferCommand represents a buffer copy command
type CopyBufferCommand struct {
	*CommandImpl
	sourceBuffer      Buffer
	destinationBuffer Buffer
	sourceOffset      int
	destinationOffset int
	size              int
}

// NewCopyBufferCommand creates a new buffer copy command
func NewCopyBufferCommand(source, destination Buffer, sourceOffset, destOffset, size int) *CopyBufferCommand {
	return &CopyBufferCommand{
		CommandImpl: &CommandImpl{
			commandType: CommandTypeCopyBuffer,
			isValid:     source != nil && destination != nil && size > 0,
		},
		sourceBuffer:      source,
		destinationBuffer: destination,
		sourceOffset:      sourceOffset,
		destinationOffset: destOffset,
		size:              size,
	}
}

// Execute executes the buffer copy command
func (c *CopyBufferCommand) Execute() error {
	if !c.IsValid() {
		return ErrInvalidState
	}

	return c.destinationBuffer.CopyFromBuffer(c.sourceBuffer, c.sourceOffset, c.destinationOffset, c.size)
}

// CopyTextureCommand represents a texture copy command
type CopyTextureCommand struct {
	*CommandImpl
	sourceTexture      Texture
	destinationTexture Texture
}

// NewCopyTextureCommand creates a new texture copy command
func NewCopyTextureCommand(source, destination Texture) *CopyTextureCommand {
	return &CopyTextureCommand{
		CommandImpl: &CommandImpl{
			commandType: CommandTypeCopyTexture,
			isValid:     source != nil && destination != nil,
		},
		sourceTexture:      source,
		destinationTexture: destination,
	}
}

// Execute executes the texture copy command
func (c *CopyTextureCommand) Execute() error {
	if !c.IsValid() {
		return ErrInvalidState
	}

	return c.destinationTexture.CopyFromTexture(c.sourceTexture)
}
