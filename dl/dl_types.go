package dl

import "github.com/opensraph/sraph/geom"

// ClipOp defines the operation to apply when clipping.
type ClipOp uint8

const (
	// ClipOpDifference subtracts the clip path from the existing clip.
	ClipOpDifference ClipOp = iota
	// ClipOpIntersect intersects the clip path with the existing clip.
	ClipOpIntersect
)

// PointMode defines how points are rendered.
type PointMode uint8

const (
	// PointModePoints draws each point separately.
	PointModePoints PointMode = iota
	// PointModeLines draws each separate pair of points as a line segment.
	PointModeLines
	// PointModePolygon draws each pair of overlapping points as a line segment.
	PointModePolygon
)

// SrcRectConstraint defines source rectangle constraint for drawing operations.
type SrcRectConstraint uint8

const (
	// SrcRectConstraintStrict enforces strict source rectangle bounds.
	SrcRectConstraintStrict SrcRectConstraint = iota
	// SrcRectConstraintFast allows relaxed source rectangle bounds for performance.
	SrcRectConstraintFast
)

// FilterMode defines the filtering mode for texture sampling.
type FilterMode uint8

const (
	// FilterModeNearest uses single sample point (nearest neighbor).
	FilterModeNearest FilterMode = iota
	// FilterModeLinear interpolates between 2x2 sample points (bilinear interpolation).
	FilterModeLinear
)

// ImageSampling defines the image sampling mode.
type ImageSampling uint8

const (
	// ImageSamplingNearestNeighbor uses nearest neighbor sampling.
	ImageSamplingNearestNeighbor ImageSampling = iota
	// ImageSamplingLinear uses linear sampling.
	ImageSamplingLinear
	// ImageSamplingMipmapLinear uses mipmap linear sampling.
	ImageSamplingMipmapLinear
	// ImageSamplingCubic uses cubic sampling.
	ImageSamplingCubic
)

// TileMode defines how to repeat, fold, or omit colors outside of the
// typically defined range of the source of the colors.
type TileMode uint8

const (
	// TileModeClamp clamps to the edge colors.
	TileModeClamp TileMode = iota
	// TileModeRepeat repeats the pattern.
	TileModeRepeat
	// TileModeMirror mirrors the pattern.
	TileModeMirror
	// TileModeDecal uses transparent pixels outside the range.
	TileModeDecal
)

// Scalar represents the primary scalar type used in display list operations.
type Scalar = geom.F32

// PaintStyle defines the style of painting (fill or stroke)
type PaintStyle uint8

const (
	// PaintStyleFill fills the shape
	PaintStyleFill PaintStyle = iota
	// PaintStyleStroke strokes the outline
	PaintStyleStroke
)
