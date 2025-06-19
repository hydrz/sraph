// Package gpu provides WebGPU implementation interfaces for Sraph.
//
// This package implements the WebGPU standard API in Go, providing cross-platform
// access to modern GPU functionality. It serves as the foundation for Sraph's
// rendering capabilities.
package gpu

import (
	"errors"
)

//go:generate go run ./gen -yaml ./gen/webgpu.yml -go wgpu.go

var (
	// ErrAdapterNotFound indicates no suitable adapter was found
	ErrAdapterNotFound = errors.New("no suitable adapter found")

	// ErrDeviceLost indicates the device was lost
	ErrDeviceLost = errors.New("device lost")

	// ErrOutOfMemory indicates insufficient memory
	ErrOutOfMemory = errors.New("out of memory")
)

// Global instance for WebGPU operations
var instance Instance

// Initialize initializes the GPU subsystem.
// This should be called before any other GPU operations.
func Initialize() error {
	// TODO: Initialize the WebGPU instance
	// This would typically involve loading native libraries
	// and setting up the WebGPU context
	return nil
}

// Shutdown shuts down the GPU subsystem and cleans up resources.
func Shutdown() {
	// TODO: Clean up global resources
}

// RequestAdapter requests a GPU adapter with the specified options.
func RequestAdapter(descriptor RequestAdapterOptions) (Adapter, error) {
	if instance == nil {
		return nil, errors.New("GPU not initialized")
	}
	return instance.RequestAdapter(descriptor)
}

// GetPreferredCanvasFormat returns the preferred texture format for canvas rendering.
func GetPreferredCanvasFormat() TextureFormat {
	// TODO: Query the actual preferred format from the platform
	// For now, return a common format
	return TextureFormatBGRA8Unorm
}

// WGSLLanguageFeatures returns the supported WGSL language features.
func WGSLLanguageFeatures() (*SupportedWGSLLanguageFeatures, error) {
	if instance == nil {
		return nil, errors.New("GPU not initialized")
	}
	return instance.WGSLLanguageFeatures()
}

// EnumerateAdapters returns all available adapters.
func EnumerateAdapters() ([]Adapter, error) {
	// TODO: Implement adapter enumeration
	// This would query all available GPU adapters on the system
	return nil, errors.New("not implemented")
}

// CreateInstance creates a WebGPU instance with the given descriptor.
func CreateInstance(descriptor InstanceDescriptor) Instance {
	// TODO: Create actual instance based on descriptor
	return instance
}
