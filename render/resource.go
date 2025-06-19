package render

import (
	"errors"
)

// Common render errors
var (
	ErrInvalidState     = errors.New("invalid state")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrResourceNotFound = errors.New("resource not found")
	ErrOperationFailed  = errors.New("operation failed")
)

// Resource represents a base GPU resource
type Resource interface {
	// GetLabel returns the debug label for this resource
	GetLabel() string

	// SetLabel sets a debug label for this resource
	SetLabel(label string)

	// IsValid returns true if the resource is valid
	IsValid() bool
}

// ResourceImpl is the base implementation for all GPU resources
type ResourceImpl struct {
	label   string
	isValid bool
}

// GetLabel returns the debug label for this resource
func (r *ResourceImpl) GetLabel() string {
	return r.label
}

// SetLabel sets a debug label for this resource
func (r *ResourceImpl) SetLabel(label string) {
	r.label = label
}

// IsValid returns true if the resource is valid
func (r *ResourceImpl) IsValid() bool {
	return r.isValid
}

// Sampler represents a texture sampler
type Sampler interface {
	Resource

	// GetMinFilter returns the minification filter
	GetMinFilter() Filter

	// GetMagFilter returns the magnification filter
	GetMagFilter() Filter

	// GetMipFilter returns the mip filter
	GetMipFilter() MipFilter

	// GetAddressModeU returns the address mode for U coordinate
	GetAddressModeU() AddressMode

	// GetAddressModeV returns the address mode for V coordinate
	GetAddressModeV() AddressMode

	// GetAddressModeW returns the address mode for W coordinate
	GetAddressModeW() AddressMode
}

// Filter defines texture filtering modes
type Filter int

const (
	FilterNearest Filter = iota
	FilterLinear
)

// MipFilter defines mip filtering modes
type MipFilter int

const (
	MipFilterNearest MipFilter = iota
	MipFilterLinear
	MipFilterNotMipmapped
)

// AddressMode defines texture address modes
type AddressMode int

const (
	AddressModeClampToEdge AddressMode = iota
	AddressModeRepeat
	AddressModeMirrorRepeat
)

// SamplerDescriptor describes the configuration for a sampler
type SamplerDescriptor struct {
	MinFilter    Filter
	MagFilter    Filter
	MipFilter    MipFilter
	AddressModeU AddressMode
	AddressModeV AddressMode
	AddressModeW AddressMode
}

// SamplerImpl is the default implementation of Sampler
type SamplerImpl struct {
	*ResourceImpl
	descriptor SamplerDescriptor
}

// NewSampler creates a new sampler with the specified descriptor
func NewSampler(descriptor SamplerDescriptor) Sampler {
	return &SamplerImpl{
		ResourceImpl: &ResourceImpl{
			label:   "",
			isValid: true,
		},
		descriptor: descriptor,
	}
}

// GetMinFilter returns the minification filter
func (s *SamplerImpl) GetMinFilter() Filter {
	return s.descriptor.MinFilter
}

// GetMagFilter returns the magnification filter
func (s *SamplerImpl) GetMagFilter() Filter {
	return s.descriptor.MagFilter
}

// GetMipFilter returns the mip filter
func (s *SamplerImpl) GetMipFilter() MipFilter {
	return s.descriptor.MipFilter
}

// GetAddressModeU returns the address mode for U coordinate
func (s *SamplerImpl) GetAddressModeU() AddressMode {
	return s.descriptor.AddressModeU
}

// GetAddressModeV returns the address mode for V coordinate
func (s *SamplerImpl) GetAddressModeV() AddressMode {
	return s.descriptor.AddressModeV
}

// GetAddressModeW returns the address mode for W coordinate
func (s *SamplerImpl) GetAddressModeW() AddressMode {
	return s.descriptor.AddressModeW
}
