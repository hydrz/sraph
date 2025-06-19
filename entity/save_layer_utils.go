package entity

import (
	"math"

	"github.com/opensraph/sraph/geom"
)

// SaveLayerUtils provides utilities for computing save layer coverage and managing
// save layer operations in the rendering pipeline.

// DefaultSizeThreshold is the default threshold for size difference comparison
const DefaultSizeThreshold = 0.3

// ComputeSaveLayerCoverage computes the coverage of a subpass in the global coordinate space.
//
// Parameters:
//   - contentCoverage: The computed coverage of the contents of the save layer.
//     This value may be empty if the save layer has no contents, or Rect::Maximum
//     if the contents are unbounded (like a destructive blend).
//   - effectTransform: The CTM of the subpass.
//   - coverageLimit: The current clip coverage. Used to bound the subpass size.
//   - imageFilter: A subpass image filter, or nil.
//   - floodOutputCoverage: Whether the coverage should be flooded to clip coverage
//     regardless of input coverage. Should be set to true when the restore Paint
//     has a destructive blend mode.
//   - floodInputCoverage: Whether the content coverage should be flooded. Should
//     be set to true if the paint has a backdrop filter or if there is a
//     transparent black effecting color filter.
//
// The coverage computation expects contentCoverage to be in the child coordinate space.
// effectTransform is used to transform this back into the global coordinate space.
// A return value of nil indicates that the coverage is empty or otherwise does not
// intersect with the parent coverage limit and should be discarded.
func ComputeSaveLayerCoverage(
	contentCoverage geom.Rect,
	effectTransform geom.Matrix,
	coverageLimit geom.Rect,
	imageFilter FilterContents,
	floodOutputCoverage bool,
	floodInputCoverage bool,
) *geom.Rect {

	coverage := contentCoverage

	// Flood input coverage if requested
	if floodInputCoverage {
		coverage = coverageLimit
	}

	// Apply image filter expansion if present
	if imageFilter != nil {
		// Get filter coverage expansion
		filterCoverage := imageFilter.GetCoverage(coverage)
		if filterCoverage != nil {
			coverage = *filterCoverage
		}
	}

	// Transform coverage to global coordinate space
	transformedCoverage := effectTransform.TransformRect(coverage)

	// Intersect with coverage limit
	finalCoverage := transformedCoverage.Intersection(coverageLimit)

	// Check if coverage is empty
	if finalCoverage.IsEmpty() {
		return nil
	}

	// Flood output coverage if requested
	if floodOutputCoverage {
		finalCoverage = coverageLimit
	}

	return &finalCoverage
}

// SizeDifferenceUnderThreshold checks if the size difference between two sizes
// is under the specified threshold
func SizeDifferenceUnderThreshold(a, b geom.Size, threshold float64) bool {
	if threshold <= 0 {
		threshold = DefaultSizeThreshold
	}

	widthDiff := math.Abs(float64(a.Width - b.Width))
	heightDiff := math.Abs(float64(a.Height - b.Height))

	maxWidth := math.Max(float64(a.Width), float64(b.Width))
	maxHeight := math.Max(float64(a.Height), float64(b.Height))

	if maxWidth == 0 && maxHeight == 0 {
		return true
	}

	if maxWidth == 0 {
		return heightDiff/maxHeight < threshold
	}

	if maxHeight == 0 {
		return widthDiff/maxWidth < threshold
	}

	return (widthDiff/maxWidth < threshold) && (heightDiff/maxHeight < threshold)
}

// FilterContents represents a filter that can be applied to content
type FilterContents interface {
	// GetCoverage returns the coverage area affected by this filter
	GetCoverage(input geom.Rect) *geom.Rect

	// IsValid returns true if the filter is valid
	IsValid() bool
}

// SaveLayerOptions holds options for save layer operations
type SaveLayerOptions struct {
	// Bounds optional bounds for the save layer
	Bounds *geom.Rect

	// Paint to apply when restoring the layer
	Paint *Paint

	// ImageFilter to apply to the layer contents
	ImageFilter FilterContents

	// BackdropFilter to apply to the backdrop
	BackdropFilter FilterContents

	// FloodOutputCoverage whether to flood output coverage
	FloodOutputCoverage bool

	// FloodInputCoverage whether to flood input coverage
	FloodInputCoverage bool
}

// Paint represents painting properties for save layer restoration
type Paint struct {
	// BlendMode for compositing
	BlendMode BlendMode

	// Alpha transparency value (0.0 to 1.0)
	Alpha float32

	// ColorFilter to apply to the paint
	ColorFilter ColorFilter
}

// ColorFilter represents a color transformation filter
type ColorFilter interface {
	// EffectsTransparentBlack returns true if this filter affects transparent black pixels
	EffectsTransparentBlack() bool

	// IsValid returns true if the filter is valid
	IsValid() bool
}

// SaveLayerHelper provides helper methods for save layer operations
type SaveLayerHelper struct{}

// NewSaveLayerHelper creates a new save layer helper
func NewSaveLayerHelper() *SaveLayerHelper {
	return &SaveLayerHelper{}
}

// ComputeOptimalBounds computes optimal bounds for a save layer operation
func (h *SaveLayerHelper) ComputeOptimalBounds(
	contentCoverage geom.Rect,
	effectTransform geom.Matrix,
	coverageLimit geom.Rect,
	options *SaveLayerOptions,
) *geom.Rect {

	if options == nil {
		options = &SaveLayerOptions{}
	}

	// Determine if we need to flood coverage based on paint properties
	floodOutput := options.FloodOutputCoverage
	floodInput := options.FloodInputCoverage

	if options.Paint != nil {
		// Check if blend mode is destructive
		if IsDestructiveBlendMode(options.Paint.BlendMode) {
			floodOutput = true
		}

		// Check if color filter affects transparent black
		if options.Paint.ColorFilter != nil && options.Paint.ColorFilter.EffectsTransparentBlack() {
			floodInput = true
		}
	}

	// Check for backdrop filter
	if options.BackdropFilter != nil {
		floodInput = true
	}

	return ComputeSaveLayerCoverage(
		contentCoverage,
		effectTransform,
		coverageLimit,
		options.ImageFilter,
		floodOutput,
		floodInput,
	)
}

// IsDestructiveBlendMode returns true if the blend mode is destructive
func IsDestructiveBlendMode(mode BlendMode) bool {
	switch mode {
	case BlendModeClear, BlendModeSrc, BlendModeDstOut, BlendModeSrcOut:
		return true
	default:
		return false
	}
}
