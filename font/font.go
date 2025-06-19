// Package font provides typography and text rendering functionality for the Sraph UI toolkit.
//
// This package handles font loading, glyph management, text shaping, and layout operations.
// It is inspired by Impeller's typographer architecture but adapted for Go and Sraph's needs.
package font

import (
	"fmt"

	"github.com/opensraph/sraph/geom"
)

// Rational represents a rational number (fraction).
type Rational struct {
	// Numerator is the numerator of the fraction.
	Numerator int32
	// Denominator is the denominator of the fraction.
	Denominator int32
}

// ToFloat32 converts the rational number to a float32.
func (r Rational) ToFloat32() float32 {
	if r.Denominator == 0 {
		return 0
	}
	return float32(r.Numerator) / float32(r.Denominator)
}

// NewRational creates a new rational number.
func NewRational(numerator, denominator int32) Rational {
	return Rational{
		Numerator:   numerator,
		Denominator: denominator,
	}
}

// Typeface represents the intrinsic properties of a font family.
// It describes the font file and its capabilities but without size or styling modifications.
type Typeface interface {
	// IsValid returns whether this typeface is valid and can be used.
	IsValid() bool

	// GetHash returns a hash value for this typeface, used for caching and comparison.
	GetHash() uint64

	// GetGlyphCount returns the number of glyphs in this typeface.
	GetGlyphCount() int

	// GetUnitsPerEM returns the number of font units per em square.
	GetUnitsPerEM() int

	// GetAscender returns the ascender in font units.
	GetAscender() float32

	// GetDescender returns the descender in font units.
	GetDescender() float32

	// GetCapHeight returns the cap height in font units.
	GetCapHeight() float32

	// GetXHeight returns the x-height in font units.
	GetXHeight() float32

	// GetGlyphForCodepoint returns the glyph ID for a given Unicode codepoint.
	GetGlyphForCodepoint(codepoint rune) (GlyphIndex, bool)

	// GetGlyphMetrics returns metrics for a specific glyph.
	GetGlyphMetrics(glyph GlyphIndex) GlyphMetrics

	// GetKerning returns the kerning adjustment between two glyphs.
	GetKerning(left, right GlyphIndex) float32
}

// AxisAlignment determines the axis along which there is subpixel positioning.
type AxisAlignment uint8

const (
	// AxisAlignmentNone indicates no subpixel positioning.
	AxisAlignmentNone AxisAlignment = iota
	// AxisAlignmentX indicates subpixel positioning along the X axis.
	AxisAlignmentX
	// AxisAlignmentY indicates subpixel positioning along the Y axis.
	AxisAlignmentY
	// AxisAlignmentAll indicates subpixel positioning along both axes.
	AxisAlignmentAll
)

// FontMetrics describes the metrics of a font at a specific size.
type FontMetrics struct {
	// Size is the font size in points.
	Size float32
	// Ascent is the distance from the baseline to the top of the tallest glyph.
	Ascent float32
	// Descent is the distance from the baseline to the bottom of the lowest glyph.
	Descent float32
	// Leading is the additional space between lines.
	Leading float32
	// LineHeight is the total line height (ascent + descent + leading).
	LineHeight float32
	// UnderlinePosition is the position of underlines relative to the baseline.
	UnderlinePosition float32
	// UnderlineThickness is the thickness of underlines.
	UnderlineThickness float32
}

// Font describes a typeface along with any modifications to its intrinsic properties.
// It combines a Typeface with size, metrics, and axis alignment information.
type Font interface {
	// IsValid returns whether this font is valid and can be used.
	IsValid() bool

	// GetTypeface returns the underlying typeface.
	GetTypeface() Typeface

	// GetMetrics returns the font metrics.
	GetMetrics() FontMetrics

	// GetAxisAlignment returns the axis alignment for subpixel positioning.
	GetAxisAlignment() AxisAlignment

	// GetHash returns a hash value for this font, used for caching and comparison.
	GetHash() uint64

	// IsEqual returns whether this font is equal to another font.
	IsEqual(other Font) bool
}

// GlyphIndex represents the index of a glyph within a typeface.
type GlyphIndex uint16

// GlyphType represents the type of glyph data.
type GlyphType uint8

const (
	// GlyphTypePath indicates the glyph is a vector path.
	GlyphTypePath GlyphType = iota
	// GlyphTypeBitmap indicates the glyph is a bitmap.
	GlyphTypeBitmap
)

// Glyph represents a single glyph with its index and type.
type Glyph struct {
	// Index is the glyph index in the typeface.
	Index GlyphIndex
	// Type indicates whether this is a path or bitmap glyph.
	Type GlyphType
}

// GlyphMetrics contains metrics for a specific glyph.
type GlyphMetrics struct {
	// AdvanceWidth is the horizontal advance width of the glyph.
	AdvanceWidth float32
	// AdvanceHeight is the vertical advance height of the glyph.
	AdvanceHeight float32
	// LeftSideBearing is the left side bearing of the glyph.
	LeftSideBearing float32
	// RightSideBearing is the right side bearing of the glyph.
	RightSideBearing float32
	// BoundingBox is the bounding box of the glyph.
	BoundingBox geom.Rect[geom.F32]
}

// ScaledFont represents a font and a scale, used as a key in glyph atlases.
type ScaledFont struct {
	// Font is the base font.
	Font Font
	// Scale is the scaling factor applied to the font.
	Scale Rational
}

// GetHash returns a hash value for this scaled font.
func (sf ScaledFont) GetHash() uint64 {
	// TODO: Implement proper hash combining
	return sf.Font.GetHash() ^ uint64(sf.Scale.Numerator)<<32 ^ uint64(sf.Scale.Denominator)
}

// SubpixelPosition represents all possible positions for subpixel alignment.
type SubpixelPosition uint8

const (
	// Subpixel00 represents no subpixel offset (0/4, 0/4).
	Subpixel00 SubpixelPosition = iota
	// Subpixel01 represents offset (0/4, 1/4).
	Subpixel01
	// Subpixel02 represents offset (0/4, 2/4).
	Subpixel02
	// Subpixel03 represents offset (0/4, 3/4).
	Subpixel03
	// Subpixel10 represents offset (1/4, 0/4).
	Subpixel10
	// Subpixel11 represents offset (1/4, 1/4).
	Subpixel11
	// Subpixel12 represents offset (1/4, 2/4).
	Subpixel12
	// Subpixel13 represents offset (1/4, 3/4).
	Subpixel13
	// Subpixel20 represents offset (2/4, 0/4).
	Subpixel20
	// Subpixel21 represents offset (2/4, 1/4).
	Subpixel21
	// Subpixel22 represents offset (2/4, 2/4).
	Subpixel22
	// Subpixel23 represents offset (2/4, 3/4).
	Subpixel23
	// Subpixel30 represents offset (3/4, 0/4).
	Subpixel30
	// Subpixel31 represents offset (3/4, 1/4).
	Subpixel31
	// Subpixel32 represents offset (3/4, 2/4).
	Subpixel32
	// Subpixel33 represents offset (3/4, 3/4).
	Subpixel33
)

// FontGlyphPair represents a combination of a scaled font and glyph.
// This is used as a key for caching rendered glyphs in atlases.
type FontGlyphPair struct {
	// ScaledFont is the font and scale information.
	ScaledFont ScaledFont
	// Glyph is the glyph information.
	Glyph Glyph
	// SubpixelPosition is the subpixel positioning information.
	SubpixelPosition SubpixelPosition
}

// GetHash returns a hash value for this font-glyph pair.
func (fgp FontGlyphPair) GetHash() uint64 {
	// TODO: Implement proper hash combining
	return fgp.ScaledFont.GetHash() ^ uint64(fgp.Glyph.Index)<<16 ^ uint64(fgp.Glyph.Type)<<8 ^ uint64(fgp.SubpixelPosition)
}

// GlyphProperties contains styling properties for a glyph.
type GlyphProperties struct {
	// Color is the color to render the glyph.
	Color geom.Color
	// Stroke contains optional stroke parameters.
	Stroke *StrokeParameters
}

// StrokeParameters defines how to stroke a glyph outline.
type StrokeParameters struct {
	// Width is the stroke width.
	Width float32
	// Cap is the line cap style.
	Cap geom.StrokeCap
	// Join is the line join style.
	Join geom.StrokeJoin
	// MiterLimit is the miter limit for miter joins.
	MiterLimit float32
}

// Weight represents font weight values (compatible with CSS font-weight).
type Weight int

const (
	// WeightThin represents weight 100.
	WeightThin Weight = 100
	// WeightExtraLight represents weight 200.
	WeightExtraLight Weight = 200
	// WeightLight represents weight 300.
	WeightLight Weight = 300
	// WeightNormal represents weight 400 (regular).
	WeightNormal Weight = 400
	// WeightMedium represents weight 500.
	WeightMedium Weight = 500
	// WeightSemiBold represents weight 600.
	WeightSemiBold Weight = 600
	// WeightBold represents weight 700.
	WeightBold Weight = 700
	// WeightExtraBold represents weight 800.
	WeightExtraBold Weight = 800
	// WeightBlack represents weight 900.
	WeightBlack Weight = 900
)

// Style represents font style values.
type Style int

const (
	// StyleNormal represents normal (upright) style.
	StyleNormal Style = iota
	// StyleItalic represents italic style.
	StyleItalic
	// StyleOblique represents oblique style.
	StyleOblique
)

// DefaultFont creates a font from a typeface with the given size and alignment.
func DefaultFont(typeface Typeface, size float32, alignment AxisAlignment) Font {
	if typeface == nil || !typeface.IsValid() {
		return nil
	}

	metrics := calculateFontMetrics(typeface, size)
	return &defaultFont{
		typeface:  typeface,
		metrics:   metrics,
		alignment: alignment,
	}
}

// calculateFontMetrics calculates font metrics for a given typeface and size.
func calculateFontMetrics(typeface Typeface, size float32) FontMetrics {
	unitsPerEM := float32(typeface.GetUnitsPerEM())
	scale := size / unitsPerEM

	ascent := typeface.GetAscender() * scale
	descent := typeface.GetDescender() * scale
	leading := 0.0 // TODO: Calculate proper leading

	return FontMetrics{
		Size:               size,
		Ascent:             ascent,
		Descent:            -descent, // Descent is typically negative
		Leading:            float32(leading),
		LineHeight:         ascent - descent + float32(leading),
		UnderlinePosition:  -size * 0.1, // TODO: Get from font
		UnderlineThickness: size * 0.05, // TODO: Get from font
	}
}

// defaultFont provides the default implementation of Font.
type defaultFont struct {
	typeface  Typeface
	metrics   FontMetrics
	alignment AxisAlignment
	hash      uint64
	hashValid bool
}

// IsValid implements Font.
func (f *defaultFont) IsValid() bool {
	return f.typeface != nil && f.typeface.IsValid()
}

// GetTypeface implements Font.
func (f *defaultFont) GetTypeface() Typeface {
	return f.typeface
}

// GetMetrics implements Font.
func (f *defaultFont) GetMetrics() FontMetrics {
	return f.metrics
}

// GetAxisAlignment implements Font.
func (f *defaultFont) GetAxisAlignment() AxisAlignment {
	return f.alignment
}

// GetHash implements Font.
func (f *defaultFont) GetHash() uint64 {
	if !f.hashValid {
		// TODO: Implement proper hash combining
		f.hash = f.typeface.GetHash() ^
			uint64(f.metrics.Size*1000) ^
			uint64(f.alignment)<<32
		f.hashValid = true
	}
	return f.hash
}

// IsEqual implements Font.
func (f *defaultFont) IsEqual(other Font) bool {
	if other == nil {
		return false
	}

	otherDefault, ok := other.(*defaultFont)
	if !ok {
		return false
	}

	return f.typeface.GetHash() == otherDefault.typeface.GetHash() &&
		f.metrics.Size == otherDefault.metrics.Size &&
		f.alignment == otherDefault.alignment
}

// String returns a string representation of the font weight.
func (w Weight) String() string {
	switch w {
	case WeightThin:
		return "Thin"
	case WeightExtraLight:
		return "ExtraLight"
	case WeightLight:
		return "Light"
	case WeightNormal:
		return "Normal"
	case WeightMedium:
		return "Medium"
	case WeightSemiBold:
		return "SemiBold"
	case WeightBold:
		return "Bold"
	case WeightExtraBold:
		return "ExtraBold"
	case WeightBlack:
		return "Black"
	default:
		return fmt.Sprintf("Weight(%d)", int(w))
	}
}

// String returns a string representation of the font style.
func (s Style) String() string {
	switch s {
	case StyleNormal:
		return "Normal"
	case StyleItalic:
		return "Italic"
	case StyleOblique:
		return "Oblique"
	default:
		return fmt.Sprintf("Style(%d)", int(s))
	}
}
