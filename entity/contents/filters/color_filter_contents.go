package filters

import (
	"github.com/opensraph/sraph/geom"
)

// ColorFilterContents provides color filter rendering capabilities
type ColorFilterContents struct {
	BaseFilterContents
	source      FilterInput
	colorMatrix geom.ColorMatrix
	filterType  ColorFilterType
}

// ColorFilterType defines the type of color filter
type ColorFilterType int

const (
	// ColorFilterTypeMatrix matrix-based color filter
	ColorFilterTypeMatrix ColorFilterType = iota
	// ColorFilterTypeHueRotate hue rotation filter
	ColorFilterTypeHueRotate
	// ColorFilterTypeSaturate saturation filter
	ColorFilterTypeSaturate
	// ColorFilterTypeBrightness brightness filter
	ColorFilterTypeBrightness
	// ColorFilterTypeContrast contrast filter
	ColorFilterTypeContrast
	// ColorFilterTypeInvert invert colors filter
	ColorFilterTypeInvert
	// ColorFilterTypeGrayscale grayscale filter
	ColorFilterTypeGrayscale
	// ColorFilterTypeSepia sepia filter
	ColorFilterTypeSepia
)

// NewColorFilterContents creates a new color filter contents
func NewColorFilterContents(source FilterInput, colorMatrix geom.ColorMatrix) *ColorFilterContents {
	return &ColorFilterContents{
		BaseFilterContents: NewBaseFilterContents(),
		source:             source,
		colorMatrix:        colorMatrix,
		filterType:         ColorFilterTypeMatrix,
	}
}

// SetSource sets the input source
func (c *ColorFilterContents) SetSource(source FilterInput) {
	c.source = source
}

// GetSource returns the input source
func (c *ColorFilterContents) GetSource() FilterInput {
	return c.source
}

// SetColorMatrix sets the color transformation matrix
func (c *ColorFilterContents) SetColorMatrix(matrix geom.ColorMatrix) {
	c.colorMatrix = matrix
}

// GetColorMatrix returns the current color matrix
func (c *ColorFilterContents) GetColorMatrix() geom.ColorMatrix {
	return c.colorMatrix
}

// SetFilterType sets the color filter type
func (c *ColorFilterContents) SetFilterType(filterType ColorFilterType) {
	c.filterType = filterType
}

// GetFilterType returns the current filter type
func (c *ColorFilterContents) GetFilterType() ColorFilterType {
	return c.filterType
}

// Render implements the FilterContents interface
func (c *ColorFilterContents) Render(context *FilterContext, entity *Entity, pass *RenderPass) bool {
	// TODO: Implement color filter rendering
	// This would involve:
	// 1. Applying the color matrix transformation
	// 2. Handling different filter types
	// 3. Optimizing for common filter operations
	return false
}

// GetBounds returns the bounds of this filter
func (c *ColorFilterContents) GetBounds() geom.Rect {
	return c.source.GetBounds()
}

// Clone creates a copy of this filter
func (c *ColorFilterContents) Clone() FilterContents {
	return &ColorFilterContents{
		BaseFilterContents: c.BaseFilterContents.Clone().(BaseFilterContents),
		source:             c.source,
		colorMatrix:        c.colorMatrix,
		filterType:         c.filterType,
	}
}

// GetCoverage returns the coverage area for this filter
func (c *ColorFilterContents) GetCoverage(transform geom.Matrix) geom.Rect {
	return c.GetBounds()
}

// CreateHueRotateFilter creates a hue rotation color filter
func CreateHueRotateFilter(source FilterInput, degrees float32) *ColorFilterContents {
	matrix := geom.CreateHueRotationMatrix(degrees)
	filter := NewColorFilterContents(source, matrix)
	filter.SetFilterType(ColorFilterTypeHueRotate)
	return filter
}

// CreateSaturateFilter creates a saturation color filter
func CreateSaturateFilter(source FilterInput, saturation float32) *ColorFilterContents {
	matrix := geom.CreateSaturationMatrix(saturation)
	filter := NewColorFilterContents(source, matrix)
	filter.SetFilterType(ColorFilterTypeSaturate)
	return filter
}

// CreateBrightnessFilter creates a brightness color filter
func CreateBrightnessFilter(source FilterInput, brightness float32) *ColorFilterContents {
	matrix := geom.CreateBrightnessMatrix(brightness)
	filter := NewColorFilterContents(source, matrix)
	filter.SetFilterType(ColorFilterTypeBrightness)
	return filter
}

// CreateContrastFilter creates a contrast color filter
func CreateContrastFilter(source FilterInput, contrast float32) *ColorFilterContents {
	matrix := geom.CreateContrastMatrix(contrast)
	filter := NewColorFilterContents(source, matrix)
	filter.SetFilterType(ColorFilterTypeContrast)
	return filter
}

// CreateInvertFilter creates an invert color filter
func CreateInvertFilter(source FilterInput) *ColorFilterContents {
	matrix := geom.CreateInvertMatrix()
	filter := NewColorFilterContents(source, matrix)
	filter.SetFilterType(ColorFilterTypeInvert)
	return filter
}

// CreateGrayscaleFilter creates a grayscale color filter
func CreateGrayscaleFilter(source FilterInput) *ColorFilterContents {
	matrix := geom.CreateGrayscaleMatrix()
	filter := NewColorFilterContents(source, matrix)
	filter.SetFilterType(ColorFilterTypeGrayscale)
	return filter
}

// CreateSepiaFilter creates a sepia color filter
func CreateSepiaFilter(source FilterInput) *ColorFilterContents {
	matrix := geom.CreateSepiaMatrix()
	filter := NewColorFilterContents(source, matrix)
	filter.SetFilterType(ColorFilterTypeSepia)
	return filter
}
