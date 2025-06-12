package impl

import (
	"fmt"
	"time"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

// Public WebGPU implementation API

// Initialize initializes the WebGPU implementation
func Initialize() error {
	// Initialize global managers with proper ordering
	globalCallbackManager = NewCallbackManager()
	globalCallbackRegistry = newCallbackRegistry(globalCallbackManager)
	globalErrorManager = NewErrorManager()
	InitializeFutureManager()

	return nil
}

// Shutdown shuts down the WebGPU implementation
func Shutdown() error {
	// Shutdown in reverse order with proper synchronization
	if globalFutureManager != nil {
		globalFutureManager.Shutdown()
	}

	if globalCallbackManager != nil {
		globalCallbackManager.Shutdown()
	}

	return nil
}

// Version returns the implementation version
func Version() string {
	return "WebGPU Implementation v1.0.0"
}

// GetFeatureLevel returns the supported feature level
func GetFeatureLevel() FeatureLevel {
	return FeatureLevelCore
}

// CreateGPU creates the main GPU interface
func CreateGPU() GPU {
	return NewGPU()
}

// CreateDefaultDevice creates a device with default settings
func CreateDefaultDevice(adapter Adapter) Future {
	descriptor := DeviceDescriptor{
		Label: "Default Device",
		RequiredFeatures: []FeatureName{
			FeatureNameDepthClipControl,
		},
		RequiredLimits: getDefaultLimits(),
		DefaultQueue: QueueDescriptor{
			Label: "Default Queue",
		},
	}

	callback := RequestDeviceCallbackInfo{
		Mode: CallbackModeAllowSpontaneous,
		Callback: func(status RequestDeviceStatus, device Device, message string) {
			if status != RequestDeviceStatusSuccess {
				fmt.Printf("Device creation failed: %s\n", message)
			}
		},
	}

	return adapter.RequestDevice(descriptor, callback)
}

// Utility functions for common operations

// CreateDefaultTexture creates a 2D texture with common settings
func CreateDefaultTexture(device Device, width, height uint32, format TextureFormat) (Texture, error) {
	descriptor := TextureDescriptor{
		Label:     "Default Texture",
		Usage:     TextureUsageTextureBinding | TextureUsageRenderAttachment,
		Dimension: TextureDimension2D,
		Size: Extent3D{
			Width:              width,
			Height:             height,
			DepthOrArrayLayers: 1,
		},
		Format:        format,
		MipLevelCount: 1,
		SampleCount:   1,
	}

	return device.CreateTexture(descriptor)
}

// CreateDefaultBuffer creates a buffer with common settings
func CreateDefaultBuffer(device Device, size uint64, usage BufferUsage) (Buffer, error) {
	descriptor := BufferDescriptor{
		Label: "Default Buffer",
		Usage: usage,
		Size:  size,
	}

	return device.CreateBuffer(descriptor)
}

// CreateDefaultSampler creates a sampler with common settings
func CreateDefaultSampler(device Device) (Sampler, error) {
	descriptor := SamplerDescriptor{
		Label:        "Default Sampler",
		AddressModeU: AddressModeClampToEdge,
		AddressModeV: AddressModeClampToEdge,
		AddressModeW: AddressModeClampToEdge,
		MagFilter:    FilterModeLinear,
		MinFilter:    FilterModeLinear,
		MipmapFilter: MipmapFilterModeLinear,
	}

	return device.CreateSampler(descriptor)
}

// CreateInstanceWithDefaults creates an instance with default settings
func CreateInstanceWithDefaults() Instance {
	descriptor := InstanceDescriptor{
		RequiredFeatures: []InstanceFeatureName{
			InstanceFeatureNameTimedWaitAnyEnable,
		},
		RequiredLimits: InstanceLimits{
			TimedWaitAnyMaxCount: 64,
		},
	}
	return NewInstance(descriptor)
}

// Helper functions for debugging

// EnableAllDebugging enables all debugging features
func EnableAllDebugging() {
	EnableDebug()
	EnableResourceTracking()
}

// DisableAllDebugging disables all debugging features
func DisableAllDebugging() {
	DisableDebug()
	DisableResourceTracking()
}

// GetAllReports returns all debugging reports
func GetAllReports() map[string]string {
	return map[string]string{
		"debug":       GetDebugReport(),
		"performance": GetPerformanceReport(),
		"resources":   GetResourceReport(),
	}
}

// Helper functions for texture creation with validation
func CreateValidatedTexture(device Device, descriptor TextureDescriptor) (Texture, error) {
	// Get device limits
	var limits Limits
	if status, err := device.GetLimits(limits); err != nil || status != StatusSuccess {
		return nil, fmt.Errorf("failed to get device limits: %v", err)
	}

	// Validate texture descriptor
	if err := ValidateTextureDescriptor(descriptor, limits); err != nil {
		return nil, fmt.Errorf("texture validation failed: %v", err)
	}

	return device.CreateTexture(descriptor)
}

// Helper function for buffer creation with validation
func CreateValidatedBuffer(device Device, descriptor BufferDescriptor) (Buffer, error) {
	// Get device limits
	var limits Limits
	if status, err := device.GetLimits(limits); err != nil || status != StatusSuccess {
		return nil, fmt.Errorf("failed to get device limits: %v", err)
	}

	// Validate buffer descriptor
	if err := ValidateBufferDescriptor(descriptor, limits); err != nil {
		return nil, fmt.Errorf("buffer validation failed: %v", err)
	}

	return device.CreateBuffer(descriptor)
}

// ProcessEvents processes all pending events
func ProcessEvents() error {
	if globalCallbackManager == nil {
		return fmt.Errorf("WebGPU not initialized")
	}

	globalCallbackManager.ProcessEvents()
	return nil
}

// WaitForEvents waits for events to be available for processing
func WaitForEvents(timeoutNS uint64) bool {
	if globalCallbackManager == nil {
		return false
	}

	timeout := time.Duration(timeoutNS) * time.Nanosecond
	return globalCallbackManager.WaitForEvents(timeout)
}

// GetPendingEventCount returns the number of pending events
func GetPendingEventCount() int32 {
	if globalCallbackManager == nil {
		return 0
	}

	return globalCallbackManager.GetPendingEventCount()
}

// CleanupExpiredFutures cleans up futures that completed more than maxAgeMS milliseconds ago
func CleanupExpiredFutures(maxAgeMS uint64) int {
	if globalFutureManager == nil {
		return 0
	}

	maxAge := time.Duration(maxAgeMS) * time.Millisecond
	return globalFutureManager.CleanupExpiredFutures(maxAge)
}

// GetFutureStats returns statistics about tracked futures
func GetFutureStats() (total, completed int) {
	if globalFutureManager == nil {
		return 0, 0
	}

	return globalFutureManager.GetFutureCount()
}
