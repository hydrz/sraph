package render

import (
	"github.com/opensraph/sraph/display"
)

// StorageMode defines how texture memory is stored and accessed
type StorageMode int

const (
	StorageModeShared StorageMode = iota
	StorageModePrivate
	StorageModeMemoryless
)

// LoadAction defines what happens to existing contents when a render pass begins
type LoadAction int

const (
	LoadActionDontCare LoadAction = iota
	LoadActionLoad
	LoadActionClear
)

// StoreAction defines what happens to render target contents when a render pass ends
type StoreAction int

const (
	StoreActionDontCare StoreAction = iota
	StoreActionStore
	StoreActionMultisampleResolve
	StoreActionStoreAndMultisampleResolve
)

// RenderTargetConfig defines the configuration for a render target
type RenderTargetConfig struct {
	Width           int
	Height          int
	MipCount        int
	HasMSAA         bool
	HasDepthStencil bool
}

// Hash returns a hash value for the render target configuration
func (c RenderTargetConfig) Hash() uint64 {
	// TODO: Implement proper hash function
	return uint64(c.Width*c.Height + c.MipCount)
}

// AttachmentConfig defines configuration for a render target attachment
type AttachmentConfig struct {
	StorageMode StorageMode
	LoadAction  LoadAction
	StoreAction StoreAction
	ClearColor  display.Color
}

// AttachmentConfigMSAA defines configuration for MSAA render target attachment
type AttachmentConfigMSAA struct {
	StorageMode        StorageMode
	ResolveStorageMode StorageMode
	LoadAction         LoadAction
	StoreAction        StoreAction
	ClearColor         display.Color
}

// ColorAttachment represents a color attachment in a render target
type ColorAttachment struct {
	Texture        Texture
	ResolveTexture Texture
	Config         AttachmentConfig
}

// DepthAttachment represents a depth attachment in a render target
type DepthAttachment struct {
	Texture Texture
	Config  AttachmentConfig
}

// StencilAttachment represents a stencil attachment in a render target
type StencilAttachment struct {
	Texture Texture
	Config  AttachmentConfig
}

// RenderTarget represents a render target with color, depth, and stencil attachments
type RenderTarget interface {
	// GetSize returns the size of the render target
	GetSize() (int, int)

	// IsValid checks if the render target is valid
	IsValid() bool

	// GetColorAttachment returns the color attachment at the specified index
	GetColorAttachment(index int) ColorAttachment

	// GetDepthAttachment returns the depth attachment
	GetDepthAttachment() DepthAttachment

	// GetStencilAttachment returns the stencil attachment
	GetStencilAttachment() StencilAttachment

	// HasColorAttachment checks if there's a color attachment at the specified index
	HasColorAttachment(index int) bool

	// HasDepthAttachment checks if there's a depth attachment
	HasDepthAttachment() bool

	// HasStencilAttachment checks if there's a stencil attachment
	HasStencilAttachment() bool

	// GetSampleCount returns the sample count for MSAA
	GetSampleCount() int
}

// RenderTargetImpl is the default implementation of RenderTarget
type RenderTargetImpl struct {
	config            RenderTargetConfig
	colorAttachments  map[int]ColorAttachment
	depthAttachment   *DepthAttachment
	stencilAttachment *StencilAttachment
	isValid           bool
}

// NewRenderTarget creates a new render target with the specified configuration
func NewRenderTarget(config RenderTargetConfig) RenderTarget {
	return &RenderTargetImpl{
		config:           config,
		colorAttachments: make(map[int]ColorAttachment),
		isValid:          true,
	}
}

// GetSize returns the size of the render target
func (rt *RenderTargetImpl) GetSize() (int, int) {
	return rt.config.Width, rt.config.Height
}

// IsValid checks if the render target is valid
func (rt *RenderTargetImpl) IsValid() bool {
	return rt.isValid
}

// GetColorAttachment returns the color attachment at the specified index
func (rt *RenderTargetImpl) GetColorAttachment(index int) ColorAttachment {
	return rt.colorAttachments[index]
}

// GetDepthAttachment returns the depth attachment
func (rt *RenderTargetImpl) GetDepthAttachment() DepthAttachment {
	if rt.depthAttachment != nil {
		return *rt.depthAttachment
	}
	return DepthAttachment{}
}

// GetStencilAttachment returns the stencil attachment
func (rt *RenderTargetImpl) GetStencilAttachment() StencilAttachment {
	if rt.stencilAttachment != nil {
		return *rt.stencilAttachment
	}
	return StencilAttachment{}
}

// HasColorAttachment checks if there's a color attachment at the specified index
func (rt *RenderTargetImpl) HasColorAttachment(index int) bool {
	_, exists := rt.colorAttachments[index]
	return exists
}

// HasDepthAttachment checks if there's a depth attachment
func (rt *RenderTargetImpl) HasDepthAttachment() bool {
	return rt.depthAttachment != nil
}

// HasStencilAttachment checks if there's a stencil attachment
func (rt *RenderTargetImpl) HasStencilAttachment() bool {
	return rt.stencilAttachment != nil
}

// GetSampleCount returns the sample count for MSAA
func (rt *RenderTargetImpl) GetSampleCount() int {
	if rt.config.HasMSAA {
		return 4 // Default MSAA sample count
	}
	return 1
}

// SetColorAttachment sets a color attachment at the specified index
func (rt *RenderTargetImpl) SetColorAttachment(index int, attachment ColorAttachment) {
	rt.colorAttachments[index] = attachment
}

// SetDepthAttachment sets the depth attachment
func (rt *RenderTargetImpl) SetDepthAttachment(attachment DepthAttachment) {
	rt.depthAttachment = &attachment
}

// SetStencilAttachment sets the stencil attachment
func (rt *RenderTargetImpl) SetStencilAttachment(attachment StencilAttachment) {
	rt.stencilAttachment = &attachment
}
