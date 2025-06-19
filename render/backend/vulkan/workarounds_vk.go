package vulkan

import (
	"strings"
	"github.com/vulkan-go/vulkan"
)

// WorkaroundsVK manages device-specific workarounds for Vulkan
type WorkaroundsVK struct {
	deviceProperties vulkan.PhysicalDeviceProperties
	driverInfo       *DriverInfoVK
	
	// Workaround flags
	DisableRobustBufferAccess    bool
	DisableGeometryShaders       bool
	DisableTessellationShaders   bool
	UseExplicitFlushForMappedMemory bool
	AllocateExtraMemoryForBuffers   bool
	DisableSubgroupOperations    bool
	DisableStorageBuffers        bool
	ForceLinearTextures          bool
	DisableAnisotropicFiltering  bool
	UseReducedShaderOptimization bool
	DisableComputeShaders        bool
}

// NewWorkaroundsVK creates a new workarounds manager
func NewWorkaroundsVK(deviceProperties vulkan.PhysicalDeviceProperties) *WorkaroundsVK {
	w := &WorkaroundsVK{
		deviceProperties: deviceProperties,
		driverInfo:       NewDriverInfoVK(deviceProperties),
	}
	
	w.applyWorkarounds()
	return w
}

// applyWorkarounds applies device-specific workarounds
func (w *WorkaroundsVK) applyWorkarounds() {
	deviceName := strings.ToLower(vulkan.ToString(w.deviceProperties.DeviceName[:]))
	vendorID := w.deviceProperties.VendorID
	driverVersion := w.deviceProperties.DriverVersion
	
	// Intel GPU workarounds
	if vendorID == 0x8086 {
		w.applyIntelWorkarounds(deviceName, driverVersion)
	}
	
	// NVIDIA GPU workarounds
	if vendorID == 0x10DE {
		w.applyNVIDIAWorkarounds(deviceName, driverVersion)
	}
	
	// AMD GPU workarounds
	if vendorID == 0x1002 {
		w.applyAMDWorkarounds(deviceName, driverVersion)
	}
	
	// Qualcomm GPU workarounds
	if vendorID == 0x5143 {
		w.applyQualcommWorkarounds(deviceName, driverVersion)
	}
	
	// ARM Mali GPU workarounds
	if vendorID == 0x13B5 {
		w.applyARMWorkarounds(deviceName, driverVersion)
	}
	
	// Mobile/embedded device workarounds
	w.applyMobileWorkarounds(deviceName)
	
	// API version specific workarounds
	w.applyAPIVersionWorkarounds()
}

// applyIntelWorkarounds applies Intel GPU specific workarounds
func (w *WorkaroundsVK) applyIntelWorkarounds(deviceName string, driverVersion uint32) {
	// Intel GPUs often have issues with robust buffer access
	w.DisableRobustBufferAccess = true
	
	// Older Intel integrated GPUs have limited geometry shader support
	if strings.Contains(deviceName, "hd") || strings.Contains(deviceName, "iris") {
		w.DisableGeometryShaders = true
		w.DisableTessellationShaders = true
	}
	
	// Intel drivers sometimes require explicit memory flushes
	w.UseExplicitFlushForMappedMemory = true
	
	// Some Intel GPUs benefit from extra memory allocation
	w.AllocateExtraMemoryForBuffers = true
}

// applyNVIDIAWorkarounds applies NVIDIA GPU specific workarounds
func (w *WorkaroundsVK) applyNVIDIAWorkarounds(deviceName string, driverVersion uint32) {
	// NVIDIA GPUs generally work well, minimal workarounds needed
	
	// Older NVIDIA drivers had subgroup operation issues
	if driverVersion < 400000000 { // Roughly driver version 400.x
		w.DisableSubgroupOperations = true
	}
}

// applyAMDWorkarounds applies AMD GPU specific workarounds
func (w *WorkaroundsVK) applyAMDWorkarounds(deviceName string, driverVersion uint32) {
	// Some older AMD drivers had storage buffer issues
	if driverVersion < 200000000 { // Roughly AMDGPU driver version 20.x
		w.DisableStorageBuffers = true
	}
	
	// AMD GPUs sometimes benefit from linear textures for certain operations
	if strings.Contains(deviceName, "vega") {
		w.ForceLinearTextures = true
	}
}

// applyQualcommWorkarounds applies Qualcomm Adreno GPU specific workarounds
func (w *WorkaroundsVK) applyQualcommWorkarounds(deviceName string, driverVersion uint32) {
	// Adreno GPUs often have reduced shader optimization needs
	w.UseReducedShaderOptimization = true
	
	// Some Adreno GPUs have limited anisotropic filtering support
	w.DisableAnisotropicFiltering = true
	
	// Memory allocation padding for Adreno GPUs
	w.AllocateExtraMemoryForBuffers = true
}

// applyARMWorkarounds applies ARM Mali GPU specific workarounds
func (w *WorkaroundsVK) applyARMWorkarounds(deviceName string, driverVersion uint32) {
	// Mali GPUs often have limited compute shader support
	w.DisableComputeShaders = true
	
	// Mali GPUs benefit from reduced shader optimization
	w.UseReducedShaderOptimization = true
	
	// Explicit memory management for Mali
	w.UseExplicitFlushForMappedMemory = true
}

// applyMobileWorkarounds applies general mobile device workarounds
func (w *WorkaroundsVK) applyMobileWorkarounds(deviceName string) {
	// Check for mobile device indicators
	if strings.Contains(deviceName, "mali") ||
	   strings.Contains(deviceName, "adreno") ||
	   strings.Contains(deviceName, "powervr") {
		
		// Mobile GPUs often have memory constraints
		w.AllocateExtraMemoryForBuffers = false // Actually allocate less on mobile
		w.UseReducedShaderOptimization = true
		
		// Disable advanced features on mobile
		w.DisableGeometryShaders = true
		w.DisableTessellationShaders = true
		w.DisableComputeShaders = true
	}
}

// applyAPIVersionWorkarounds applies workarounds based on Vulkan API version
func (w *WorkaroundsVK) applyAPIVersionWorkarounds() {
	apiVersion := w.deviceProperties.ApiVersion
	
	// Workarounds for older Vulkan versions
	if apiVersion < vulkan.ApiVersion11 {
		w.DisableSubgroupOperations = true
	}
}

// ShouldDisableRobustBufferAccess returns whether robust buffer access should be disabled
func (w *WorkaroundsVK) ShouldDisableRobustBufferAccess() bool {
	return w.DisableRobustBufferAccess
}

// ShouldDisableGeometryShaders returns whether geometry shaders should be disabled
func (w *WorkaroundsVK) ShouldDisableGeometryShaders() bool {
	return w.DisableGeometryShaders
}

// ShouldUseExplicitFlushForMappedMemory returns whether explicit memory flushing is needed
func (w *WorkaroundsVK) ShouldUseExplicitFlushForMappedMemory() bool {
	return w.UseExplicitFlushForMappedMemory
}

// ShouldAllocateExtraMemoryForBuffers returns whether extra memory should be allocated
func (w *WorkaroundsVK) ShouldAllocateExtraMemoryForBuffers() bool {
	return w.AllocateExtraMemoryForBuffers
}

// GetDriverInfo returns the driver information
func (w *WorkaroundsVK) GetDriverInfo() *DriverInfoVK {
	return w.driverInfo
}

// DriverInfoVK contains driver-specific information
type DriverInfoVK struct {
	VendorID      uint32
	DeviceID      uint32
	DriverVersion uint32
	VendorName    string
	DeviceName    string
	DriverName    string
	APIVersion    uint32
}

// NewDriverInfoVK creates driver info from device properties
func NewDriverInfoVK(properties vulkan.PhysicalDeviceProperties) *DriverInfoVK {
	return &DriverInfoVK{
		VendorID:      properties.VendorID,
		DeviceID:      properties.DeviceID,
		DriverVersion: properties.DriverVersion,
		VendorName:    w.getVendorName(properties.VendorID),
		DeviceName:    vulkan.ToString(properties.DeviceName[:]),
		APIVersion:    properties.ApiVersion,
	}
}

// getVendorName returns the vendor name for a vendor ID
func (w *WorkaroundsVK) getVendorName(vendorID uint32) string {
	switch vendorID {
	case 0x8086:
		return "Intel"
	case 0x10DE:
		return "NVIDIA"
	case 0x1002:
		return "AMD"
	case 0x5143:
		return "Qualcomm"
	case 0x13B5:
		return "ARM"
	default:
		return "Unknown"
	}
}
