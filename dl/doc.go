// Package dl provides a display list system for 2D graphics rendering.
//
// The display list system allows you to record a sequence of drawing operations
// and replay them efficiently on different rendering backends. This is inspired
// by Flutter's display list architecture but implemented in Go.
//
// # Core Concepts
//
// The display list system consists of several key components:
//
//   - DisplayList: An immutable sequence of recorded drawing operations
//   - DisplayListBuilder: A builder for recording operations into a display list
//   - OpReceiver: Interface for objects that can receive and execute operations
//   - Canvas: A concrete implementation of OpReceiver for immediate rendering
//   - Renderer: Interface for different rendering backends (GPU, CPU, etc.)
//
// # Basic Usage
//
//	// Create a builder to record operations
//	builder := dl.NewDisplayListBuilder()
//
//	// Create paint objects
//	paint := dl.NewPaint()
//	paint.SetColor(dl.ColorRed)
//
//	// Record drawing operations
//	builder.DrawRect(geom.NewRect[dl.Scalar](10, 10, 100, 100), paint)
//	builder.Translate(50, 50)
//	builder.DrawCircle(geom.Point[dl.Scalar]{X: 0, Y: 0}, 25, paint)
//
//	// Build the display list
//	displayList := builder.Build()
//
//	// Render the display list
//	canvas := dl.NewCanvas(geom.NewRect[dl.Scalar](0, 0, 200, 200))
//	displayList.Dispatch(canvas)
//
// # Advanced Features
//
// The display list system supports:
//
//   - Complex transformations (translate, scale, rotate, skew, perspective)
//   - Clipping (rectangle, rounded rectangle, path-based)
//   - Various drawing primitives (lines, rectangles, circles, paths, vertices)
//   - Paint effects (colors, blend modes, anti-aliasing)
//   - Save/restore state management
//   - Bounds tracking and culling
//   - Path construction with bezier curves
//
// # Performance Considerations
//
// Display lists are designed for efficiency:
//
//   - Operations are recorded once and can be replayed multiple times
//   - Bounds information allows for culling operations outside the viewport
//   - Attribute flags help optimize rendering pipeline decisions
//   - Immutable design enables safe sharing across threads
//
// # Integration with Graphics Backends
//
// The system is designed to work with different rendering backends:
//
//   - Implement the Renderer interface for your graphics API
//   - Use RenderContext to manage state and dispatch operations
//   - Operations are backend-agnostic and can target GPU or CPU renderers
//
// For more examples and detailed usage, see the dl_examples.go file.
package dl

// Version information
const (
	// VersionMajor is the major version number.
	VersionMajor = 0
	// VersionMinor is the minor version number.
	VersionMinor = 1
	// VersionPatch is the patch version number.
	VersionPatch = 0
)

// GetVersion returns the version string of the display list package.
func GetVersion() string {
	return "0.1.0"
}

// FeatureFlags represents optional features that may be supported by renderers.
type FeatureFlags uint32

const (
	// FeatureFlagAntiAliasing indicates support for anti-aliased rendering.
	FeatureFlagAntiAliasing FeatureFlags = 1 << iota
	// FeatureFlagAdvancedBlending indicates support for advanced blend modes.
	FeatureFlagAdvancedBlending
	// FeatureFlagPathRendering indicates support for complex path rendering.
	FeatureFlagPathRendering
	// FeatureFlagImageFilters indicates support for image filters.
	FeatureFlagImageFilters
	// FeatureFlagTextRendering indicates support for text rendering.
	FeatureFlagTextRendering
	// FeatureFlagGradients indicates support for gradient rendering.
	FeatureFlagGradients
)

// HasFeature checks if a specific feature flag is set.
func (f FeatureFlags) HasFeature(flag FeatureFlags) bool {
	return (f & flag) != 0
}

// WithFeature returns a new FeatureFlags with the specified flag set.
func (f FeatureFlags) WithFeature(flag FeatureFlags) FeatureFlags {
	return f | flag
}

// WithoutFeature returns a new FeatureFlags with the specified flag cleared.
func (f FeatureFlags) WithoutFeature(flag FeatureFlags) FeatureFlags {
	return f &^ flag
}

// RenderingHints provides hints to renderers about preferred rendering quality vs performance.
type RenderingHints struct {
	// PreferQuality indicates that quality should be prioritized over performance.
	PreferQuality bool
	// EnableCaching indicates that intermediate results should be cached when possible.
	EnableCaching bool
	// MaxComplexity provides a hint about the maximum acceptable complexity for operations.
	MaxComplexity int
}

// DefaultRenderingHints returns the default rendering hints.
func DefaultRenderingHints() RenderingHints {
	return RenderingHints{
		PreferQuality: true,
		EnableCaching: true,
		MaxComplexity: 1000,
	}
}

// Statistics provides information about display list usage and performance.
type Statistics struct {
	// TotalOperations is the total number of operations recorded.
	TotalOperations int
	// DrawOperations is the number of drawing operations.
	DrawOperations int
	// TransformOperations is the number of transformation operations.
	TransformOperations int
	// ClipOperations is the number of clipping operations.
	ClipOperations int
	// StateOperations is the number of save/restore operations.
	StateOperations int
	// BoundsArea is the total area covered by the display list bounds.
	BoundsArea Scalar
}

// NewStatistics creates a new Statistics object.
func NewStatistics() *Statistics {
	return &Statistics{}
}

// Reset resets all statistics counters to zero.
func (s *Statistics) Reset() {
	s.TotalOperations = 0
	s.DrawOperations = 0
	s.TransformOperations = 0
	s.ClipOperations = 0
	s.StateOperations = 0
	s.BoundsArea = 0
}

// RecordOperation records statistics for an operation.
func (s *Statistics) RecordOperation(op Operation) {
	s.TotalOperations++

	// Classify operation types based on the operation
	// This would need to be implemented based on the actual operation types
	// For now, we'll use a simple heuristic
	switch op.(type) {
	case *SaveOp, *RestoreOp, *SaveLayerOp:
		s.StateOperations++
	case *TranslateOp, *ScaleOp, *RotateOp, *SkewOp, *Transform2DAffineOp, *TransformFullPerspectiveOp:
		s.TransformOperations++
	case *ClipRectOp, *ClipRRectOp, *ClipPathOp:
		s.ClipOperations++
	default:
		s.DrawOperations++
	}
}

// Global statistics instance
var globalStats = NewStatistics()

// GetGlobalStatistics returns the global statistics instance.
func GetGlobalStatistics() *Statistics {
	return globalStats
}

// Error types for display list operations
type (
	// InvalidOperationError is returned when an invalid operation is attempted.
	InvalidOperationError struct {
		Message string
	}

	// ResourceError is returned when a resource-related error occurs.
	ResourceError struct {
		Message string
		Code    int
	}

	// RenderingError is returned when a rendering error occurs.
	RenderingError struct {
		Message string
		Backend string
	}
)

// Error implements the error interface for InvalidOperationError.
func (e *InvalidOperationError) Error() string {
	return "invalid operation: " + e.Message
}

// Error implements the error interface for ResourceError.
func (e *ResourceError) Error() string {
	return "resource error: " + e.Message
}

// Error implements the error interface for RenderingError.
func (e *RenderingError) Error() string {
	return "rendering error (" + e.Backend + "): " + e.Message
}

// Debug and development helpers

// DebugMode controls whether debug information is collected and logged.
var DebugMode = false

// SetDebugMode enables or disables debug mode.
func SetDebugMode(enabled bool) {
	DebugMode = enabled
}

// IsDebugMode returns true if debug mode is enabled.
func IsDebugMode() bool {
	return DebugMode
}

// TODO: Future features to implement
//
// The following features are planned for future versions:
//
//   - Text rendering support (fonts, text layout, text effects)
//   - Image rendering support (textures, image filters, nine-patch)
//   - Advanced effects (shadows, gradients, mask filters)
//   - Animation support (interpolation, easing functions)
//   - Spatial indexing (R-tree for efficient culling)
//   - Serialization (save/load display lists)
//   - Multi-threading support (parallel rendering)
//   - GPU-accelerated backends (Vulkan, Metal, DirectX)
//   - Vector graphics import/export (SVG, PDF)
//   - Performance profiling tools
//   - Memory pool optimizations
//   - Tessellation improvements
//   - Advanced clipping modes
//   - Layer effects and compositing
