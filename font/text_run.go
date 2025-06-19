// Package font - Text run functionality
package font

import (
	"github.com/opensraph/sraph/geom"
)

// TextRun represents a collection of positioned glyphs from a specific font.
// This is inspired by Impeller's TextRun but adapted for Sraph's architecture.
type TextRun interface {
	// IsValid returns whether this text run is valid.
	IsValid() bool

	// GetFont returns the font used for this text run.
	GetFont() Font

	// GetGlyphCount returns the number of glyphs in this run.
	GetGlyphCount() int

	// GetGlyphPositions returns all glyph positions in this run.
	GetGlyphPositions() []GlyphPosition

	// AddGlyph adds a glyph at the specified position to this run.
	AddGlyph(glyph Glyph, position geom.Point[geom.F32]) bool

	// GetBounds returns the bounding box of this text run.
	GetBounds() geom.Rect[geom.F32]
}

// GlyphPosition represents a glyph positioned within a text run.
type GlyphPosition struct {
	// Glyph is the glyph data.
	Glyph Glyph
	// Position is the position of this glyph relative to the run's origin.
	Position geom.Point[geom.F32]
}

// NewTextRun creates a new text run with the specified font.
func NewTextRun(font Font) TextRun {
	if font == nil || !font.IsValid() {
		return nil
	}

	return &defaultTextRun{
		font:      font,
		glyphs:    make([]GlyphPosition, 0),
		isValid:   true,
		bounds:    geom.Rect[geom.F32]{},
		boundsSet: false,
	}
}

// NewTextRunWithGlyphs creates a new text run with the specified font and glyphs.
func NewTextRunWithGlyphs(font Font, glyphs []GlyphPosition) TextRun {
	if font == nil || !font.IsValid() {
		return nil
	}

	run := &defaultTextRun{
		font:      font,
		glyphs:    make([]GlyphPosition, len(glyphs)),
		isValid:   true,
		bounds:    geom.Rect[geom.F32]{},
		boundsSet: false,
	}

	copy(run.glyphs, glyphs)
	return run
}

// defaultTextRun provides the default implementation of TextRun.
type defaultTextRun struct {
	font      Font
	glyphs    []GlyphPosition
	isValid   bool
	bounds    geom.Rect[geom.F32]
	boundsSet bool
}

// IsValid implements TextRun.
func (tr *defaultTextRun) IsValid() bool {
	return tr.isValid && tr.font != nil && tr.font.IsValid()
}

// GetFont implements TextRun.
func (tr *defaultTextRun) GetFont() Font {
	return tr.font
}

// GetGlyphCount implements TextRun.
func (tr *defaultTextRun) GetGlyphCount() int {
	return len(tr.glyphs)
}

// GetGlyphPositions implements TextRun.
func (tr *defaultTextRun) GetGlyphPositions() []GlyphPosition {
	// Return a copy to prevent external modification
	result := make([]GlyphPosition, len(tr.glyphs))
	copy(result, tr.glyphs)
	return result
}

// AddGlyph implements TextRun.
func (tr *defaultTextRun) AddGlyph(glyph Glyph, position geom.Point[geom.F32]) bool {
	if !tr.IsValid() {
		return false
	}

	tr.glyphs = append(tr.glyphs, GlyphPosition{
		Glyph:    glyph,
		Position: position,
	})

	// Invalidate cached bounds
	tr.boundsSet = false
	return true
}

// GetBounds implements TextRun.
func (tr *defaultTextRun) GetBounds() geom.Rect[geom.F32] {
	if !tr.IsValid() || len(tr.glyphs) == 0 {
		return geom.Rect[geom.F32]{}
	}

	if tr.boundsSet {
		return tr.bounds
	}

	// Calculate bounds by examining all glyph positions
	if len(tr.glyphs) == 0 {
		tr.bounds = geom.Rect[geom.F32]{}
		tr.boundsSet = true
		return tr.bounds
	}

	// Initialize with first glyph
	firstPos := tr.glyphs[0].Position
	firstGlyph := tr.glyphs[0].Glyph
	firstMetrics := tr.font.GetTypeface().GetGlyphMetrics(firstGlyph.Index)

	minX := firstPos.X + firstMetrics.BoundingBox.Left
	minY := firstPos.Y + firstMetrics.BoundingBox.Top
	maxX := firstPos.X + firstMetrics.BoundingBox.Right
	maxY := firstPos.Y + firstMetrics.BoundingBox.Bottom

	// Extend bounds with remaining glyphs
	for i := 1; i < len(tr.glyphs); i++ {
		pos := tr.glyphs[i].Position
		glyph := tr.glyphs[i].Glyph
		metrics := tr.font.GetTypeface().GetGlyphMetrics(glyph.Index)

		glyphMinX := pos.X + metrics.BoundingBox.Left
		glyphMinY := pos.Y + metrics.BoundingBox.Top
		glyphMaxX := pos.X + metrics.BoundingBox.Right
		glyphMaxY := pos.Y + metrics.BoundingBox.Bottom

		if glyphMinX < minX {
			minX = glyphMinX
		}
		if glyphMinY < minY {
			minY = glyphMinY
		}
		if glyphMaxX > maxX {
			maxX = glyphMaxX
		}
		if glyphMaxY > maxY {
			maxY = glyphMaxY
		}
	}

	tr.bounds = geom.Rect[geom.F32]{
		Left:   minX,
		Top:    minY,
		Right:  maxX,
		Bottom: maxY,
	}
	tr.boundsSet = true

	return tr.bounds
}
