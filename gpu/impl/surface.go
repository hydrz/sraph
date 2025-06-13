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
	label        string
	configured   bool
	config       SurfaceConfiguration
	capabilities SurfaceCapabilities
	destroyed    bool
}

// newSurface creates a new WebGPU surface
func newSurface(label string) Surface {
	return &surface{
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

// NewSurface creates a new WebGPU surface (public factory function)
func NewSurface(descriptor SurfaceDescriptor) Surface {
	return newSurface(descriptor.Label)
}

// Configure configures the surface
func (s *surface) Configure(config SurfaceConfiguration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.destroyed {
		return fmt.Errorf("surface has been destroyed")
	}

	// Validate configuration
	if config.Device == nil {
		return fmt.Errorf("device cannot be nil")
	}
	if config.Width == 0 || config.Height == 0 {
		return fmt.Errorf("width and height must be greater than 0")
	}

	// Check if format is supported
	formatSupported := false
	for _, format := range s.capabilities.Formats {
		if format == config.Format {
			formatSupported = true
			break
		}
	}
	if !formatSupported {
		return fmt.Errorf("format %v is not supported", config.Format)
	}

	s.config = config
	s.configured = true
	return nil
}

// Capabilities implements Surface.Capabilities.
func (s *surface) Capabilities(adapter Adapter) (*SurfaceCapabilities, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.destroyed {
		return nil, fmt.Errorf("surface has been destroyed")
	}

	caps := s.capabilities
	return &caps, nil
}

// CurrentTexture implements Surface.CurrentTexture.
func (s *surface) CurrentTexture() (*SurfaceTexture, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.destroyed {
		return nil, fmt.Errorf("surface has been destroyed")
	}

	if !s.configured {
		return nil, fmt.Errorf("surface is not configured")
	}

	texture := &texture{
		label:         "Surface Texture",
		usage:         s.config.Usage,
		dimension:     TextureDimension2D,
		size:          Extent3D{Width: s.config.Width, Height: s.config.Height, DepthOrArrayLayers: 1},
		format:        s.config.Format,
		mipLevelCount: 1,
		sampleCount:   1,
	}

	return &SurfaceTexture{
		Texture: texture,
		Status:  SurfaceGetCurrentTextureStatusSuccessOptimal,
	}, nil
}

// Present implements Surface.Present.
func (s *surface) Present() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.destroyed {
		return fmt.Errorf("surface has been destroyed")
	}

	if !s.configured {
		return fmt.Errorf("surface is not configured")
	}

	// In a real implementation, this would present the frame to the display
	return nil
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
