package gpu

import (
	"context"
	"fmt"
	"sync"
)

var _ Texture = (*texture)(nil)

// texture implements the Texture interface
type texture struct {
	mu            sync.RWMutex
	refCount      int32
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

// newTexture creates a new WebGPU texture
func newTexture(descriptor TextureDescriptor) Texture {
	return &texture{
		refCount:      1,
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
func (t *texture) CreateView(ctx context.Context, descriptor TextureViewDescriptor) (TextureView, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return nil, fmt.Errorf("texture has been destroyed")
	}

	view := &TextureViewImpl{
		refCount:        1,
		label:           descriptor.Label,
		format:          descriptor.Format,
		dimension:       descriptor.Dimension,
		baseMipLevel:    descriptor.BaseMipLevel,
		mipLevelCount:   descriptor.MipLevelCount,
		baseArrayLayer:  descriptor.BaseArrayLayer,
		arrayLayerCount: descriptor.ArrayLayerCount,
		aspect:          descriptor.Aspect,
		usage:           descriptor.Usage,
		texture:         t,
	}

	return view, nil
}

// Destroy destroys the texture
func (t *texture) Destroy(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.destroyed {
		return fmt.Errorf("texture has been destroyed")
	}

	t.destroyed = true
	return nil
}

// GetDepthOrArrayLayers gets the depth or array layers
func (t *texture) GetDepthOrArrayLayers(ctx context.Context) (uint32, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return 0, fmt.Errorf("texture has been destroyed")
	}

	return t.size.DepthOrArrayLayers, nil
}

// GetDimension gets the texture dimension
func (t *texture) GetDimension(ctx context.Context) (TextureDimension, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return TextureDimensionUndefined, fmt.Errorf("texture has been destroyed")
	}

	return t.dimension, nil
}

// GetFormat gets the texture format
func (t *texture) GetFormat(ctx context.Context) (TextureFormat, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return TextureFormatUndefined, fmt.Errorf("texture has been destroyed")
	}

	return t.format, nil
}

// GetHeight gets the texture height
func (t *texture) GetHeight(ctx context.Context) (uint32, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return 0, fmt.Errorf("texture has been destroyed")
	}

	return t.size.Height, nil
}

// GetMipLevelCount gets the mip level count
func (t *texture) GetMipLevelCount(ctx context.Context) (uint32, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return 0, fmt.Errorf("texture has been destroyed")
	}

	return t.mipLevelCount, nil
}

// GetSampleCount gets the sample count
func (t *texture) GetSampleCount(ctx context.Context) (uint32, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return 0, fmt.Errorf("texture has been destroyed")
	}

	return t.sampleCount, nil
}

// GetUsage gets the texture usage
func (t *texture) GetUsage(ctx context.Context) (TextureUsage, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return TextureUsageNone, fmt.Errorf("texture has been destroyed")
	}

	return t.usage, nil
}

// GetWidth gets the texture width
func (t *texture) GetWidth(ctx context.Context) (uint32, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.destroyed {
		return 0, fmt.Errorf("texture has been destroyed")
	}

	return t.size.Width, nil
}

// SetLabel sets the texture label
func (t *texture) SetLabel(ctx context.Context, label string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.destroyed {
		return fmt.Errorf("texture has been destroyed")
	}

	t.label = label
	return nil
}

// AddRef increments the reference count
func (t *texture) AddRef(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.destroyed {
		return fmt.Errorf("texture has been destroyed")
	}

	t.refCount++
	return nil
}

// Release decrements the reference count and destroys if zero
func (t *texture) Release(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.destroyed {
		return fmt.Errorf("texture has been destroyed")
	}

	t.refCount--
	if t.refCount <= 0 {
		t.destroyed = true
	}

	return nil
}

// TextureViewImpl implements the TextureView interface
type TextureViewImpl struct {
	mu              sync.RWMutex
	refCount        int32
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
func (tv *TextureViewImpl) SetLabel(ctx context.Context, label string) error {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	if tv.destroyed {
		return fmt.Errorf("texture view has been destroyed")
	}

	tv.label = label
	return nil
}

// AddRef increments the reference count
func (tv *TextureViewImpl) AddRef(ctx context.Context) error {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	if tv.destroyed {
		return fmt.Errorf("texture view has been destroyed")
	}

	tv.refCount++
	return nil
}

// Release decrements the reference count and destroys if zero
func (tv *TextureViewImpl) Release(ctx context.Context) error {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	if tv.destroyed {
		return fmt.Errorf("texture view has been destroyed")
	}

	tv.refCount--
	if tv.refCount <= 0 {
		tv.destroyed = true
	}

	return nil
}
