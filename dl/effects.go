package dl

import (
	"github.com/opensraph/sraph/geom"
)

// ColorSource represents a source of color information for painting.
// This can be a solid color, gradient, pattern, or other color generator.
type ColorSource interface {
	// Color returns the color at the specified position.
	Color(x, y float32) Color
	// IsOpaque returns true if the color source is fully opaque.
	IsOpaque() bool
	// String returns a string representation of the color source.
	String() string
}

// ColorFilter represents a filter that can transform colors.
type ColorFilter interface {
	// ApplyFilter applies the color filter to the input color.
	ApplyFilter(color Color) Color
	// String returns a string representation of the color filter.
	String() string
}

// ImageFilter represents a filter that can be applied to rendered content.
type ImageFilter interface {
	// Bounds returns the bounds that this filter would expand content to.
	Bounds(inputBounds geom.Rect[Scalar]) geom.Rect[Scalar]
	// String returns a string representation of the image filter.
	String() string
}

// MaskFilter represents a filter that can be applied to the alpha channel.
type MaskFilter interface {
	// Bounds returns the bounds that this mask filter would expand content to.
	Bounds(inputBounds geom.Rect[Scalar]) geom.Rect[Scalar]
	// String returns a string representation of the mask filter.
	String() string
}

// PathEffect represents an effect that can be applied to paths.
type PathEffect interface {
	// ApplyEffect applies the path effect to the input path.
	ApplyEffect(path *Path) *Path
	// String returns a string representation of the path effect.
	String() string
}

// Gradient represents a color gradient.
type Gradient struct {
	colors    []Color
	positions []float32
	tileMode  TileMode
}

// NewLinearGradient creates a new linear gradient.
func NewLinearGradient(start, end geom.Point[Scalar], colors []Color, positions []float32, tileMode TileMode) *Gradient {
	return &Gradient{
		colors:    colors,
		positions: positions,
		tileMode:  tileMode,
	}
}

// Color implements ColorSource for Gradient.
func (g *Gradient) Color(x, y float32) Color {
	// TODO: Implement gradient color calculation
	if len(g.colors) == 0 {
		return ColorTransparent
	}
	return g.colors[0]
}

// IsOpaque implements ColorSource for Gradient.
func (g *Gradient) IsOpaque() bool {
	for _, color := range g.colors {
		if !color.IsOpaque() {
			return false
		}
	}
	return true
}

// String implements ColorSource for Gradient.
func (g *Gradient) String() string {
	return "Gradient{...}"
}

// SolidColorSource represents a solid color source.
type SolidColorSource struct {
	color Color
}

// NewSolidColorSource creates a new solid color source.
func NewSolidColorSource(color Color) *SolidColorSource {
	return &SolidColorSource{color: color}
}

// Color implements ColorSource for SolidColorSource.
func (s *SolidColorSource) Color(x, y float32) Color {
	return s.color
}

// IsOpaque implements ColorSource for SolidColorSource.
func (s *SolidColorSource) IsOpaque() bool {
	return s.color.IsOpaque()
}

// String implements ColorSource for SolidColorSource.
func (s *SolidColorSource) String() string {
	return "SolidColorSource{" + s.color.String() + "}"
}

// MatrixColorFilter applies a color matrix transformation.
type MatrixColorFilter struct {
	matrix [20]float32 // 4x5 color transformation matrix
}

// NewMatrixColorFilter creates a new matrix color filter.
func NewMatrixColorFilter(matrix [20]float32) *MatrixColorFilter {
	return &MatrixColorFilter{matrix: matrix}
}

// ApplyFilter implements ColorFilter for MatrixColorFilter.
func (m *MatrixColorFilter) ApplyFilter(color Color) Color {
	// TODO: Implement matrix color transformation
	return color
}

// String implements ColorFilter for MatrixColorFilter.
func (m *MatrixColorFilter) String() string {
	return "MatrixColorFilter{...}"
}

// BlurImageFilter applies a Gaussian blur effect.
type BlurImageFilter struct {
	sigmaX, sigmaY float32
	tileMode       TileMode
}

// NewBlurImageFilter creates a new blur image filter.
func NewBlurImageFilter(sigmaX, sigmaY float32, tileMode TileMode) *BlurImageFilter {
	return &BlurImageFilter{
		sigmaX:   sigmaX,
		sigmaY:   sigmaY,
		tileMode: tileMode,
	}
}

// Bounds implements ImageFilter for BlurImageFilter.
func (b *BlurImageFilter) Bounds(inputBounds geom.Rect[Scalar]) geom.Rect[Scalar] {
	// Blur expands bounds by approximately 3 * sigma
	expansion := Scalar(3.0 * max(b.sigmaX, b.sigmaY))
	return geom.NewRect(
		inputBounds.Left-expansion,
		inputBounds.Top-expansion,
		inputBounds.Right+expansion,
		inputBounds.Bottom+expansion,
	)
}

// String implements ImageFilter for BlurImageFilter.
func (b *BlurImageFilter) String() string {
	return "BlurImageFilter{...}"
}

// BlurMaskFilter applies a blur to the alpha channel.
type BlurMaskFilter struct {
	sigma float32
	style BlurStyle
}

// BlurStyle defines the style of blur for mask filters.
type BlurStyle uint8

const (
	// BlurStyleNormal applies blur to both sides of the edge.
	BlurStyleNormal BlurStyle = iota
	// BlurStyleSolid applies blur only to the transparent pixels.
	BlurStyleSolid
	// BlurStyleOuter applies blur only outside the original shape.
	BlurStyleOuter
	// BlurStyleInner applies blur only inside the original shape.
	BlurStyleInner
)

// NewBlurMaskFilter creates a new blur mask filter.
func NewBlurMaskFilter(sigma float32, style BlurStyle) *BlurMaskFilter {
	return &BlurMaskFilter{
		sigma: sigma,
		style: style,
	}
}

// Bounds implements MaskFilter for BlurMaskFilter.
func (b *BlurMaskFilter) Bounds(inputBounds geom.Rect[Scalar]) geom.Rect[Scalar] {
	expansion := Scalar(3.0 * b.sigma)
	return geom.NewRect(
		inputBounds.Left-expansion,
		inputBounds.Top-expansion,
		inputBounds.Right+expansion,
		inputBounds.Bottom+expansion,
	)
}

// String implements MaskFilter for BlurMaskFilter.
func (b *BlurMaskFilter) String() string {
	return "BlurMaskFilter{...}"
}

// DashPathEffect creates dashed lines.
type DashPathEffect struct {
	intervals []float32
	phase     float32
}

// NewDashPathEffect creates a new dash path effect.
func NewDashPathEffect(intervals []float32, phase float32) *DashPathEffect {
	return &DashPathEffect{
		intervals: intervals,
		phase:     phase,
	}
}

// ApplyEffect implements PathEffect for DashPathEffect.
func (d *DashPathEffect) ApplyEffect(path *Path) *Path {
	// TODO: Implement dash path effect
	return path
}

// String implements PathEffect for DashPathEffect.
func (d *DashPathEffect) String() string {
	return "DashPathEffect{...}"
}
