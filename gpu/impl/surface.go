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
