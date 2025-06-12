package impl

import (
	"fmt"
	"sync"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ Surface = (*surface)(nil)

// surface implements the Surface interface
type surface struct {
	mu           sync.RWMutex
	refCount     int32
	label        string
	configured   bool
	config       SurfaceConfiguration
	capabilities SurfaceCapabilities
	destroyed    bool
}

// newSurface creates a new WebGPU surface
func newSurface(label string) Surface {
	return &surface{
		refCount:   1,
		label:      label,
		configured: false,
		capabilities: SurfaceCapabilities{
			Usages: TextureUsageRenderAttachment | TextureUsageCopySrc,
			Formats: []TextureFormat{
				TextureFormatBGRA8Unorm,
				TextureFormatBGRA8UnormSrgb,
				TextureFormatRGBA8Unorm,
				TextureFormatRGBA8UnormSrgb,
			},
			PresentModes: []PresentMode{
				PresentModeFifo,
				PresentModeImmediate,
				PresentModeMailbox,
			},
			AlphaModes: []CompositeAlphaMode{
				CompositeAlphaModeOpaque,
				CompositeAlphaModePremultiplied,
			},
		},
	}
}

// Configure configures the surface
func (s *surface) Configure(config SurfaceConfiguration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.destroyed {
		return fmt.Errorf("surface has been destroyed")
	}

	s.config = config
	s.configured = true
	return nil
}

// GetCapabilities gets surface capabilities
func (s *surface) GetCapabilities(adapter Adapter, capabilities SurfaceCapabilities) (Status, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.destroyed {
		return StatusError, fmt.Errorf("surface has been destroyed")
	}

	capabilities.Usages = s.capabilities.Usages
	capabilities.Formats = append(capabilities.Formats, s.capabilities.Formats...)
	capabilities.PresentModes = append(capabilities.PresentModes, s.capabilities.PresentModes...)
	capabilities.AlphaModes = append(capabilities.AlphaModes, s.capabilities.AlphaModes...)

	return StatusSuccess, nil
}

// GetCurrentTexture gets the current surface texture
func (s *surface) GetCurrentTexture(surfaceTexture SurfaceTexture) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.destroyed {
		return fmt.Errorf("surface has been destroyed")
	}

	if !s.configured {
		return fmt.Errorf("surface is not configured")
	}

	// Create a new texture for the current frame
	texture := &texture{
		refCount:      1,
		label:         "Surface Texture",
		usage:         s.config.Usage,
		dimension:     TextureDimension2D,
		size:          Extent3D{Width: s.config.Width, Height: s.config.Height, DepthOrArrayLayers: 1},
		format:        s.config.Format,
		mipLevelCount: 1,
		sampleCount:   1,
	}

	surfaceTexture.Texture = texture
	surfaceTexture.Status = SurfaceGetCurrentTextureStatusSuccessOptimal

	return nil
}

// Present presents the current frame
func (s *surface) Present() (Status, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.destroyed {
		return StatusError, fmt.Errorf("surface has been destroyed")
	}

	if !s.configured {
		return StatusError, fmt.Errorf("surface is not configured")
	}

	// In a real implementation, this would present the frame to the display
	return StatusSuccess, nil
}

// SetLabel sets the surface label
func (s *surface) SetLabel(label string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.destroyed {
		return fmt.Errorf("surface has been destroyed")
	}

	s.label = label
	return nil
}

// Unconfigure unconfigures the surface
func (s *surface) Unconfigure() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.destroyed {
		return fmt.Errorf("surface has been destroyed")
	}

	s.configured = false
	return nil
}

// AddRef increments the reference count
func (s *surface) AddRef() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.destroyed {
		return fmt.Errorf("surface has been destroyed")
	}

	s.refCount++
	return nil
}

// Release decrements the reference count and destroys if zero
func (s *surface) Release() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.destroyed {
		return fmt.Errorf("surface has been destroyed")
	}

	s.refCount--
	if s.refCount <= 0 {
		s.destroyed = true
	}

	return nil
}
