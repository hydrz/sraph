package impl

import (
	"fmt"
	"sync"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ Texture = (*texture)(nil)
var _ TextureView = (*TextureViewImpl)(nil)

// texture implements the Texture interface
type texture struct {
	mu            sync.RWMutex
	label         string
	usage         TextureUsage
	dimension     TextureDimension
	size          Extent3D
	format        TextureFormat
	mipLevelCount uint32
	sampleCount   uint32
	viewFormats   []TextureFormat
	destroyed     bool
}

// NewTexture creates a new WebGPU texture
func NewTexture(descriptor TextureDescriptor) Texture {
	return &texture{
		label:         descriptor.Label,
		usage:         descriptor.Usage,
		dimension:     descriptor.Dimension,
		size:          descriptor.Size,
		format:        descriptor.Format,
		mipLevelCount: descriptor.MipLevelCount,
		sampleCount:   descriptor.SampleCount,
		viewFormats:   descriptor.ViewFormats,
	}
}

// CreateView creates a texture view
func (t *texture) CreateView(descriptor TextureViewDescriptor) (TextureView, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return nil, fmt.Errorf("texture has been destroyed")
	}

	view := NewTextureView(descriptor, t)

	return view, nil
}

// NewTextureView creates a new texture view (factory function)
func NewTextureView(descriptor TextureViewDescriptor, texture *texture) TextureView {
	return &TextureViewImpl{
		label:           descriptor.Label,
		format:          descriptor.Format,
		dimension:       descriptor.Dimension,
		baseMipLevel:    descriptor.BaseMipLevel,
		mipLevelCount:   descriptor.MipLevelCount,
		baseArrayLayer:  descriptor.BaseArrayLayer,
		arrayLayerCount: descriptor.ArrayLayerCount,
		aspect:          descriptor.Aspect,
		usage:           descriptor.Usage,
		texture:         texture,
	}
}

// Destroy destroys the texture
func (t *texture) Destroy() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.destroyed {
		return fmt.Errorf("texture has been destroyed")
	}

	t.destroyed = true
	return nil
}

// DepthOrArrayLayers implements Texture.DepthOrArrayLayers.
func (t *texture) DepthOrArrayLayers() (uint32, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return 0, fmt.Errorf("texture has been destroyed")
	}

	return t.size.DepthOrArrayLayers, nil
}

// Dimension implements Texture.Dimension.
func (t *texture) Dimension() (TextureDimension, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return TextureDimensionUndefined, fmt.Errorf("texture has been destroyed")
	}

	return t.dimension, nil
}

// Format implements Texture.Format.
func (t *texture) Format() (TextureFormat, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return TextureFormatUndefined, fmt.Errorf("texture has been destroyed")
	}

	return t.format, nil
}

// Height implements Texture.Height.
func (t *texture) Height() (uint32, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return 0, fmt.Errorf("texture has been destroyed")
	}

	return t.size.Height, nil
}

// MipLevelCount implements Texture.MipLevelCount.
func (t *texture) MipLevelCount() (uint32, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return 0, fmt.Errorf("texture has been destroyed")
	}

	return t.mipLevelCount, nil
}

// SampleCount implements Texture.SampleCount.
func (t *texture) SampleCount() (uint32, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return 0, fmt.Errorf("texture has been destroyed")
	}

	return t.sampleCount, nil
}

// Usage implements Texture.Usage.
func (t *texture) Usage() (TextureUsage, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return TextureUsageNone, fmt.Errorf("texture has been destroyed")
	}

	return t.usage, nil
}

// Width implements Texture.Width.
func (t *texture) Width() (uint32, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return 0, fmt.Errorf("texture has been destroyed")
	}

	return t.size.Width, nil
}

// SetLabel sets the texture label
func (t *texture) SetLabel(label string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.destroyed {
		return fmt.Errorf("texture has been destroyed")
	}

	t.label = label
	return nil
}

// TextureViewImpl implements the TextureView interface
type TextureViewImpl struct {
	mu              sync.RWMutex
	label           string
	format          TextureFormat
	dimension       TextureViewDimension
	baseMipLevel    uint32
	mipLevelCount   uint32
	baseArrayLayer  uint32
	arrayLayerCount uint32
	aspect          TextureAspect
	usage           TextureUsage
	texture         *texture
	destroyed       bool
}

// SetLabel sets the texture view label
func (tv *TextureViewImpl) SetLabel(label string) error {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	if tv.destroyed {
		return fmt.Errorf("texture view has been destroyed")
	}

	tv.label = label
	return nil
}
