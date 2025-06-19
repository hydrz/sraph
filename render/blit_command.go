package render

import (
	"github.com/opensraph/sraph/gpu"
)

// BlitCommand represents a blit (block transfer) operation command.
// Blit operations are used to copy data between textures and buffers efficiently.
type BlitCommand interface {
	// GetSource returns the source resource for the blit operation
	GetSource() BlitSource

	// GetDestination returns the destination resource for the blit operation
	GetDestination() BlitDestination

	// GetBlitRegion returns the region to be copied
	GetBlitRegion() BlitRegion

	// IsValid returns true if the blit command is valid
	IsValid() bool
}

// BlitSource represents the source of a blit operation
type BlitSource interface {
	// GetTexture returns the source texture (if applicable)
	GetTexture() gpu.Texture

	// GetBuffer returns the source buffer (if applicable)
	GetBuffer() gpu.Buffer

	// GetRegion returns the source region
	GetRegion() BlitRegion
}

// BlitDestination represents the destination of a blit operation
type BlitDestination interface {
	// GetTexture returns the destination texture (if applicable)
	GetTexture() gpu.Texture

	// GetBuffer returns the destination buffer (if applicable)
	GetBuffer() gpu.Buffer

	// GetRegion returns the destination region
	GetRegion() BlitRegion
}

// BlitRegion defines a region for blit operations
type BlitRegion struct {
	// Origin of the region
	Origin gpu.Origin3D

	// Size of the region
	Size gpu.Extent3D

	// Mip level for texture operations
	MipLevel uint32

	// Array layer for texture arrays
	ArrayLayer uint32
}

// BlitType defines the type of blit operation
type BlitType int

const (
	BlitTypeTextureToTexture BlitType = iota
	BlitTypeBufferToTexture
	BlitTypeTextureToBuffer
	BlitTypeBufferToBuffer
)

// DefaultBlitCommand provides a default implementation of BlitCommand
type DefaultBlitCommand struct {
	source      BlitSource
	destination BlitDestination
	blitRegion  BlitRegion
	blitType    BlitType
}

// NewBlitCommand creates a new blit command
func NewBlitCommand(source BlitSource, destination BlitDestination, region BlitRegion) BlitCommand {
	return &DefaultBlitCommand{
		source:      source,
		destination: destination,
		blitRegion:  region,
		blitType:    determineBlitType(source, destination),
	}
}

// GetSource returns the source resource for the blit operation
func (c *DefaultBlitCommand) GetSource() BlitSource {
	return c.source
}

// GetDestination returns the destination resource for the blit operation
func (c *DefaultBlitCommand) GetDestination() BlitDestination {
	return c.destination
}

// GetBlitRegion returns the region to be copied
func (c *DefaultBlitCommand) GetBlitRegion() BlitRegion {
	return c.blitRegion
}

// IsValid returns true if the blit command is valid
func (c *DefaultBlitCommand) IsValid() bool {
	return c.source != nil && c.destination != nil
}

// GetBlitType returns the type of blit operation
func (c *DefaultBlitCommand) GetBlitType() BlitType {
	return c.blitType
}

// determineBlitType determines the blit type based on source and destination
func determineBlitType(source BlitSource, destination BlitDestination) BlitType {
	sourceHasTexture := source.GetTexture() != nil
	sourceHasBuffer := source.GetBuffer() != nil
	destHasTexture := destination.GetTexture() != nil
	destHasBuffer := destination.GetBuffer() != nil

	if sourceHasTexture && destHasTexture {
		return BlitTypeTextureToTexture
	} else if sourceHasBuffer && destHasTexture {
		return BlitTypeBufferToTexture
	} else if sourceHasTexture && destHasBuffer {
		return BlitTypeTextureToBuffer
	} else {
		return BlitTypeBufferToBuffer
	}
}

// DefaultBlitSource provides a default implementation of BlitSource
type DefaultBlitSource struct {
	texture gpu.Texture
	buffer  gpu.Buffer
	region  BlitRegion
}

// NewBlitSourceFromTexture creates a blit source from a texture
func NewBlitSourceFromTexture(texture gpu.Texture, region BlitRegion) BlitSource {
	return &DefaultBlitSource{
		texture: texture,
		buffer:  nil,
		region:  region,
	}
}

// NewBlitSourceFromBuffer creates a blit source from a buffer
func NewBlitSourceFromBuffer(buffer gpu.Buffer, region BlitRegion) BlitSource {
	return &DefaultBlitSource{
		texture: nil,
		buffer:  buffer,
		region:  region,
	}
}

// GetTexture returns the source texture
func (s *DefaultBlitSource) GetTexture() gpu.Texture {
	return s.texture
}

// GetBuffer returns the source buffer
func (s *DefaultBlitSource) GetBuffer() gpu.Buffer {
	return s.buffer
}

// GetRegion returns the source region
func (s *DefaultBlitSource) GetRegion() BlitRegion {
	return s.region
}

// DefaultBlitDestination provides a default implementation of BlitDestination
type DefaultBlitDestination struct {
	texture gpu.Texture
	buffer  gpu.Buffer
	region  BlitRegion
}

// NewBlitDestinationFromTexture creates a blit destination from a texture
func NewBlitDestinationFromTexture(texture gpu.Texture, region BlitRegion) BlitDestination {
	return &DefaultBlitDestination{
		texture: texture,
		buffer:  nil,
		region:  region,
	}
}

// NewBlitDestinationFromBuffer creates a blit destination from a buffer
func NewBlitDestinationFromBuffer(buffer gpu.Buffer, region BlitRegion) BlitDestination {
	return &DefaultBlitDestination{
		texture: nil,
		buffer:  buffer,
		region:  region,
	}
}

// GetTexture returns the destination texture
func (d *DefaultBlitDestination) GetTexture() gpu.Texture {
	return d.texture
}

// GetBuffer returns the destination buffer
func (d *DefaultBlitDestination) GetBuffer() gpu.Buffer {
	return d.buffer
}

// GetRegion returns the destination region
func (d *DefaultBlitDestination) GetRegion() BlitRegion {
	return d.region
}
