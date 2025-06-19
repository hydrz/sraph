package render

import (
	"github.com/opensraph/sraph/gpu"
)

// BlitPass represents a pass that performs blit (block transfer) operations.
// Blit passes are used to copy data between textures and buffers efficiently
// without the need for a full render pass.
type BlitPass interface {
	// SetLabel sets a debug label for the blit pass
	SetLabel(label string)

	// GetLabel returns the debug label
	GetLabel() string

	// AddCopy adds a copy operation between textures
	AddCopy(source gpu.Texture, destination gpu.Texture, sourceOrigin, destOrigin gpu.Origin3D, size gpu.Extent3D)

	// AddCopyBuffer adds a copy operation between buffers
	AddCopyBuffer(source gpu.Buffer, destination gpu.Buffer, sourceOffset, destOffset, size uint64)

	// AddCopyBufferToTexture adds a copy operation from buffer to texture
	AddCopyBufferToTexture(source gpu.Buffer, destination gpu.Texture,
		bufferOffset uint64, textureOrigin gpu.Origin3D, textureSize gpu.Extent3D)

	// AddCopyTextureToBuffer adds a copy operation from texture to buffer
	AddCopyTextureToBuffer(source gpu.Texture, destination gpu.Buffer,
		textureOrigin gpu.Origin3D, textureSize gpu.Extent3D, bufferOffset uint64)

	// GenerateMipmap generates mipmaps for a texture
	GenerateMipmap(texture gpu.Texture)

	// End finalizes the blit pass and submits all operations
	End() bool

	// IsValid returns true if the blit pass is valid
	IsValid() bool
}

// BlitOperation represents a single blit operation
type BlitOperation struct {
	// Type of the blit operation
	Type BlitOperationType

	// Source resource
	Source BlitResource

	// Destination resource
	Destination BlitResource

	// Copy parameters
	CopyParams BlitCopyParams
}

// BlitOperationType defines the type of blit operation
type BlitOperationType int

const (
	BlitOperationTypeCopyTexture BlitOperationType = iota
	BlitOperationTypeCopyBuffer
	BlitOperationTypeCopyBufferToTexture
	BlitOperationTypeCopyTextureToBuffer
	BlitOperationTypeGenerateMipmap
)

// BlitResource represents a resource used in blit operations
type BlitResource struct {
	// Texture resource (if applicable)
	Texture gpu.Texture

	// Buffer resource (if applicable)
	Buffer gpu.Buffer

	// Mip level for textures
	MipLevel uint32

	// Array layer for texture arrays
	ArrayLayer uint32
}

// BlitCopyParams contains parameters for copy operations
type BlitCopyParams struct {
	// Source offset/origin
	SourceOrigin gpu.Origin3D

	// Destination offset/origin
	DestinationOrigin gpu.Origin3D

	// Size of the copy operation
	Size gpu.Extent3D

	// Buffer offset for buffer operations
	SourceBufferOffset      uint64
	DestinationBufferOffset uint64
}

// DefaultBlitPass provides a default implementation of BlitPass
type DefaultBlitPass struct {
	// Debug label
	label string

	// List of blit operations
	operations []BlitOperation

	// GPU context
	context gpu.Context

	// Whether the pass has been ended
	ended bool
}

// NewBlitPass creates a new blit pass
func NewBlitPass(context gpu.Context) BlitPass {
	return &DefaultBlitPass{
		label:      "",
		operations: make([]BlitOperation, 0),
		context:    context,
		ended:      false,
	}
}

// SetLabel sets a debug label for the blit pass
func (p *DefaultBlitPass) SetLabel(label string) {
	p.label = label
}

// GetLabel returns the debug label
func (p *DefaultBlitPass) GetLabel() string {
	return p.label
}

// AddCopy adds a copy operation between textures
func (p *DefaultBlitPass) AddCopy(source gpu.Texture, destination gpu.Texture, sourceOrigin, destOrigin gpu.Origin3D, size gpu.Extent3D) {
	if p.ended {
		return
	}

	operation := BlitOperation{
		Type: BlitOperationTypeCopyTexture,
		Source: BlitResource{
			Texture: source,
		},
		Destination: BlitResource{
			Texture: destination,
		},
		CopyParams: BlitCopyParams{
			SourceOrigin:      sourceOrigin,
			DestinationOrigin: destOrigin,
			Size:              size,
		},
	}

	p.operations = append(p.operations, operation)
}

// AddCopyBuffer adds a copy operation between buffers
func (p *DefaultBlitPass) AddCopyBuffer(source gpu.Buffer, destination gpu.Buffer, sourceOffset, destOffset, size uint64) {
	if p.ended {
		return
	}

	operation := BlitOperation{
		Type: BlitOperationTypeCopyBuffer,
		Source: BlitResource{
			Buffer: source,
		},
		Destination: BlitResource{
			Buffer: destination,
		},
		CopyParams: BlitCopyParams{
			SourceBufferOffset:      sourceOffset,
			DestinationBufferOffset: destOffset,
			Size: gpu.Extent3D{
				Width:  uint32(size),
				Height: 1,
				Depth:  1,
			},
		},
	}

	p.operations = append(p.operations, operation)
}

// AddCopyBufferToTexture adds a copy operation from buffer to texture
func (p *DefaultBlitPass) AddCopyBufferToTexture(source gpu.Buffer, destination gpu.Texture,
	bufferOffset uint64, textureOrigin gpu.Origin3D, textureSize gpu.Extent3D) {
	if p.ended {
		return
	}

	operation := BlitOperation{
		Type: BlitOperationTypeCopyBufferToTexture,
		Source: BlitResource{
			Buffer: source,
		},
		Destination: BlitResource{
			Texture: destination,
		},
		CopyParams: BlitCopyParams{
			SourceBufferOffset: bufferOffset,
			DestinationOrigin:  textureOrigin,
			Size:               textureSize,
		},
	}

	p.operations = append(p.operations, operation)
}

// AddCopyTextureToBuffer adds a copy operation from texture to buffer
func (p *DefaultBlitPass) AddCopyTextureToBuffer(source gpu.Texture, destination gpu.Buffer,
	textureOrigin gpu.Origin3D, textureSize gpu.Extent3D, bufferOffset uint64) {
	if p.ended {
		return
	}

	operation := BlitOperation{
		Type: BlitOperationTypeCopyTextureToBuffer,
		Source: BlitResource{
			Texture: source,
		},
		Destination: BlitResource{
			Buffer: destination,
		},
		CopyParams: BlitCopyParams{
			SourceOrigin:            textureOrigin,
			DestinationBufferOffset: bufferOffset,
			Size:                    textureSize,
		},
	}

	p.operations = append(p.operations, operation)
}

// GenerateMipmap generates mipmaps for a texture
func (p *DefaultBlitPass) GenerateMipmap(texture gpu.Texture) {
	if p.ended {
		return
	}

	operation := BlitOperation{
		Type: BlitOperationTypeGenerateMipmap,
		Source: BlitResource{
			Texture: texture,
		},
		Destination: BlitResource{
			Texture: texture,
		},
	}

	p.operations = append(p.operations, operation)
}

// End finalizes the blit pass and submits all operations
func (p *DefaultBlitPass) End() bool {
	if p.ended {
		return false
	}

	p.ended = true

	// TODO: Submit all operations to the GPU context
	// This would involve encoding all the blit operations into GPU commands

	return true
}

// IsValid returns true if the blit pass is valid
func (p *DefaultBlitPass) IsValid() bool {
	return p.context != nil && !p.ended
}

// GetOperations returns all blit operations (for testing/debugging)
func (p *DefaultBlitPass) GetOperations() []BlitOperation {
	return p.operations
}
