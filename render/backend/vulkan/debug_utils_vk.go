package vulkan

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

// DebugUtilsVK provides debug utilities for Vulkan debugging
// It manages debug messengers and validation layers
type DebugUtilsVK struct {
	instance         vulkan.Instance
	debugMessenger   vulkan.DebugUtilsMessenger
	validationLayers []string
	enabled          bool
}

// DebugSeverity represents debug message severity levels
type DebugSeverity uint32

const (
	DebugSeverityVerbose DebugSeverity = DebugSeverity(vulkan.DebugUtilsMessageSeverityVerboseBit)
	DebugSeverityInfo    DebugSeverity = DebugSeverity(vulkan.DebugUtilsMessageSeverityInfoBit)
	DebugSeverityWarning DebugSeverity = DebugSeverity(vulkan.DebugUtilsMessageSeverityWarningBit)
	DebugSeverityError   DebugSeverity = DebugSeverity(vulkan.DebugUtilsMessageSeverityErrorBit)
)

// DebugType represents debug message types
type DebugType uint32

const (
	DebugTypeGeneral     DebugType = DebugType(vulkan.DebugUtilsMessageTypeGeneralBit)
	DebugTypeValidation  DebugType = DebugType(vulkan.DebugUtilsMessageTypeValidationBit)
	DebugTypePerformance DebugType = DebugType(vulkan.DebugUtilsMessageTypePerformanceBit)
)

// DebugCallback is the callback function type for debug messages
type DebugCallback func(severity DebugSeverity, msgType DebugType, message string)

// NewDebugUtilsVK creates a new debug utils manager
func NewDebugUtilsVK(instance vulkan.Instance) *DebugUtilsVK {
	return &DebugUtilsVK{
		instance: instance,
		validationLayers: []string{
			"VK_LAYER_KHRONOS_validation",
		},
		enabled: false,
	}
}

// SetupDebugMessenger sets up the debug messenger
func (d *DebugUtilsVK) SetupDebugMessenger(callback DebugCallback) error {
	if !d.enabled {
		return nil
	}

	// TODO: Implement debug messenger setup
	// This requires proper callback function setup with CGO
	return nil
}

// DestroyDebugMessenger destroys the debug messenger
func (d *DebugUtilsVK) DestroyDebugMessenger() {
	if d.debugMessenger != vulkan.NullDebugUtilsMessenger {
		// TODO: Destroy debug messenger
		// vulkan.DestroyDebugUtilsMessenger(d.instance, d.debugMessenger, nil)
		d.debugMessenger = vulkan.NullDebugUtilsMessenger
	}
}

// IsEnabled returns whether debug utils are enabled
func (d *DebugUtilsVK) IsEnabled() bool {
	return d.enabled
}

// SetEnabled enables or disables debug utils
func (d *DebugUtilsVK) SetEnabled(enabled bool) {
	d.enabled = enabled
}

// GetValidationLayers returns the validation layers
func (d *DebugUtilsVK) GetValidationLayers() []string {
	return d.validationLayers
}

// CheckValidationLayerSupport checks if validation layers are supported
func (d *DebugUtilsVK) CheckValidationLayerSupport() (bool, error) {
	var layerCount uint32
	if result := vulkan.EnumerateInstanceLayerProperties(&layerCount, nil); result != vulkan.Success {
		return false, fmt.Errorf("failed to enumerate instance layer properties: %s", result)
	}

	if layerCount == 0 {
		return false, nil
	}

	availableLayers := make([]vulkan.LayerProperties, layerCount)
	if result := vulkan.EnumerateInstanceLayerProperties(&layerCount, availableLayers); result != vulkan.Success {
		return false, fmt.Errorf("failed to get instance layer properties: %s", result)
	}

	// Check if all required validation layers are available
	for _, layerName := range d.validationLayers {
		layerFound := false
		for _, layerProps := range availableLayers {
			if layerName == vulkan.ToString(layerProps.LayerName[:]) {
				layerFound = true
				break
			}
		}
		if !layerFound {
			return false, nil
		}
	}

	return true, nil
}

// SetObjectName sets a debug name for a Vulkan object
func (d *DebugUtilsVK) SetObjectName(device vulkan.Device, objectType vulkan.ObjectType, objectHandle uint64, name string) error {
	if !d.enabled {
		return nil
	}

	// TODO: Implement object naming
	// This requires proper string conversion and extension support
	return nil
}

// SetObjectTag sets a debug tag for a Vulkan object
func (d *DebugUtilsVK) SetObjectTag(device vulkan.Device, objectType vulkan.ObjectType, objectHandle uint64, tagName uint64, tagData []byte) error {
	if !d.enabled {
		return nil
	}

	// TODO: Implement object tagging
	return nil
}

// BeginDebugLabel begins a debug label region
func (d *DebugUtilsVK) BeginDebugLabel(commandBuffer vulkan.CommandBuffer, labelName string, color [4]float32) {
	if !d.enabled {
		return
	}

	// TODO: Implement debug label begin
}

// EndDebugLabel ends a debug label region
func (d *DebugUtilsVK) EndDebugLabel(commandBuffer vulkan.CommandBuffer) {
	if !d.enabled {
		return
	}

	// TODO: Implement debug label end
}

// InsertDebugLabel inserts a debug label
func (d *DebugUtilsVK) InsertDebugLabel(commandBuffer vulkan.CommandBuffer, labelName string, color [4]float32) {
	if !d.enabled {
		return
	}

	// TODO: Implement debug label insert
}

// PrintDebugMessage prints a debug message
func (d *DebugUtilsVK) PrintDebugMessage(severity DebugSeverity, msgType DebugType, message string) {
	if !d.enabled {
		return
	}

	severityStr := d.severityToString(severity)
	typeStr := d.typeToString(msgType)

	fmt.Printf("[%s][%s] %s\n", severityStr, typeStr, message)
}

// severityToString converts severity to string
func (d *DebugUtilsVK) severityToString(severity DebugSeverity) string {
	switch severity {
	case DebugSeverityVerbose:
		return "VERBOSE"
	case DebugSeverityInfo:
		return "INFO"
	case DebugSeverityWarning:
		return "WARNING"
	case DebugSeverityError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// typeToString converts type to string
func (d *DebugUtilsVK) typeToString(msgType DebugType) string {
	switch msgType {
	case DebugTypeGeneral:
		return "GENERAL"
	case DebugTypeValidation:
		return "VALIDATION"
	case DebugTypePerformance:
		return "PERFORMANCE"
	default:
		return "UNKNOWN"
	}
}

// Destroy destroys the debug utils
func (d *DebugUtilsVK) Destroy() {
	d.DestroyDebugMessenger()
}

// ValidationFeaturesVK manages validation features
type ValidationFeaturesVK struct {
	enabledFeatures  []vulkan.ValidationFeatureEnable
	disabledFeatures []vulkan.ValidationFeatureDisable
}

// NewValidationFeaturesVK creates a new validation features manager
func NewValidationFeaturesVK() *ValidationFeaturesVK {
	return &ValidationFeaturesVK{
		enabledFeatures:  make([]vulkan.ValidationFeatureEnable, 0),
		disabledFeatures: make([]vulkan.ValidationFeatureDisable, 0),
	}
}

// EnableFeature enables a validation feature
func (vf *ValidationFeaturesVK) EnableFeature(feature vulkan.ValidationFeatureEnable) {
	vf.enabledFeatures = append(vf.enabledFeatures, feature)
}

// DisableFeature disables a validation feature
func (vf *ValidationFeaturesVK) DisableFeature(feature vulkan.ValidationFeatureDisable) {
	vf.disabledFeatures = append(vf.disabledFeatures, feature)
}

// GetEnabledFeatures returns enabled validation features
func (vf *ValidationFeaturesVK) GetEnabledFeatures() []vulkan.ValidationFeatureEnable {
	return vf.enabledFeatures
}

// GetDisabledFeatures returns disabled validation features
func (vf *ValidationFeaturesVK) GetDisabledFeatures() []vulkan.ValidationFeatureDisable {
	return vf.disabledFeatures
}

// GetValidationFeaturesInfo returns validation features info structure
func (vf *ValidationFeaturesVK) GetValidationFeaturesInfo() *vulkan.ValidationFeatures {
	if len(vf.enabledFeatures) == 0 && len(vf.disabledFeatures) == 0 {
		return nil
	}

	// TODO: Return proper validation features info
	return nil
}
