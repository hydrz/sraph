// Package font - Text frame functionality
package font

import (
	"github.com/opensraph/sraph/geom"
)

// TextFrame represents a collection of shaped text runs.
// This is the primary entry point for text rendering in Sraph,
// inspired by Impeller's TextFrame architecture.
type TextFrame interface {
	// IsValid returns whether this text frame is valid.
	IsValid() bool

	// GetBounds returns the bounding box of the entire text frame.
	GetBounds() geom.Rect[geom.F32]

	// GetRunCount returns the number of text runs in this frame.
	GetRunCount() int

	// GetRuns returns all text runs in this frame.
	GetRuns() []TextRun

	// GetAtlasType returns the atlas type needed for this frame.
	GetAtlasType() AtlasType

	// HasColor returns whether this frame contains color glyphs.
	HasColor() bool

	// AddRun adds a text run to this frame.
	AddRun(run TextRun) bool
}

// AtlasType represents the type of glyph atlas needed.
type AtlasType uint8

const (
	// AtlasTypeAlpha indicates an alpha-only atlas for regular text.
	AtlasTypeAlpha AtlasType = iota
	// AtlasTypeColor indicates a color atlas for emoji and color fonts.
	AtlasTypeColor
)

// NewTextFrame creates a new empty text frame.
func NewTextFrame() TextFrame {
	return &defaultTextFrame{
		runs:      make([]TextRun, 0),
		isValid:   true,
		bounds:    geom.Rect[geom.F32]{},
		boundsSet: false,
		hasColor:  false,
	}
}

// NewTextFrameWithRuns creates a new text frame with the specified runs.
func NewTextFrameWithRuns(runs []TextRun, hasColor bool) TextFrame {
	frame := &defaultTextFrame{
		runs:      make([]TextRun, len(runs)),
		isValid:   true,
		bounds:    geom.Rect[geom.F32]{},
		boundsSet: false,
		hasColor:  hasColor,
	}

	copy(frame.runs, runs)
	return frame
}

// defaultTextFrame provides the default implementation of TextFrame.
type defaultTextFrame struct {
	runs      []TextRun
	isValid   bool
	bounds    geom.Rect[geom.F32]
	boundsSet bool
	hasColor  bool
}

// IsValid implements TextFrame.
func (tf *defaultTextFrame) IsValid() bool {
	return tf.isValid
}

// GetBounds implements TextFrame.
func (tf *defaultTextFrame) GetBounds() geom.Rect[geom.F32] {
	if !tf.IsValid() || len(tf.runs) == 0 {
		return geom.Rect[geom.F32]{}
	}

	if tf.boundsSet {
		return tf.bounds
	}

	// Calculate bounds by union of all run bounds
	if len(tf.runs) == 0 {
		tf.bounds = geom.Rect[geom.F32]{}
		tf.boundsSet = true
		return tf.bounds
	}

	// Initialize with first run bounds
	tf.bounds = tf.runs[0].GetBounds()

	// Union with remaining run bounds
	for i := 1; i < len(tf.runs); i++ {
		runBounds := tf.runs[i].GetBounds()
		tf.bounds = tf.bounds.Union(runBounds)
	}

	tf.boundsSet = true
	return tf.bounds
}

// GetRunCount implements TextFrame.
func (tf *defaultTextFrame) GetRunCount() int {
	return len(tf.runs)
}

// GetRuns implements TextFrame.
func (tf *defaultTextFrame) GetRuns() []TextRun {
	// Return a copy to prevent external modification
	result := make([]TextRun, len(tf.runs))
	copy(result, tf.runs)
	return result
}

// GetAtlasType implements TextFrame.
func (tf *defaultTextFrame) GetAtlasType() AtlasType {
	if tf.hasColor {
		return AtlasTypeColor
	}
	return AtlasTypeAlpha
}

// HasColor implements TextFrame.
func (tf *defaultTextFrame) HasColor() bool {
	return tf.hasColor
}

// AddRun implements TextFrame.
func (tf *defaultTextFrame) AddRun(run TextRun) bool {
	if !tf.IsValid() || run == nil || !run.IsValid() {
		return false
	}

	tf.runs = append(tf.runs, run)

	// Invalidate cached bounds
	tf.boundsSet = false
	return true
}

// RoundScaledFontSize rounds a scaled font size to a rational number.
// This is used for consistent glyph caching across different scales.
func RoundScaledFontSize(scale float32) Rational {
	// Round to nearest 0.125 (1/8) to reduce cache pressure
	const precision = 8
	roundedNumerator := int32(scale*precision + 0.5)
	return Rational{
		Numerator:   roundedNumerator,
		Denominator: precision,
	}
}

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
