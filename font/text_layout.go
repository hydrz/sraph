// Package font - Text layout functionality
package font

import (
	"errors"
	"strings"
	"unicode"

	"github.com/opensraph/sraph/geom"
)

// TextLayout handles text shaping and layout operations.
// It converts text strings into positioned glyphs ready for rendering.
type TextLayout interface {
	// LayoutText layouts text using the specified font and returns a text frame.
	LayoutText(text string, font Font, maxWidth float32) (TextFrame, error)

	// LayoutParagraph layouts a paragraph with multiple text runs.
	LayoutParagraph(paragraph Paragraph, maxWidth float32) (TextFrame, error)

	// MeasureText measures the size of text without creating a full layout.
	MeasureText(text string, font Font) geom.Size[geom.F32]

	// GetLineHeight returns the line height for the given font.
	GetLineHeight(font Font) float32

	// WordWrap performs word wrapping on text.
	WordWrap(text string, font Font, maxWidth float32) []string
}

// Paragraph represents a collection of text runs with styling.
type Paragraph interface {
	// GetTextRuns returns all text runs in this paragraph.
	GetTextRuns() []TextRun

	// AddTextRun adds a text run to this paragraph.
	AddTextRun(run TextRun)

	// GetText returns the complete text content.
	GetText() string

	// GetBounds returns the bounding box of the paragraph.
	GetBounds() geom.Rect[geom.F32]

	// GetAlignment returns the text alignment.
	GetAlignment() TextAlignment

	// SetAlignment sets the text alignment.
	SetAlignment(alignment TextAlignment)
}

// TextAlignment represents text alignment options.
type TextAlignment int

const (
	// TextAlignmentLeft aligns text to the left.
	TextAlignmentLeft TextAlignment = iota
	// TextAlignmentCenter centers text.
	TextAlignmentCenter
	// TextAlignmentRight aligns text to the right.
	TextAlignmentRight
	// TextAlignmentJustify justifies text.
	TextAlignmentJustify
)

// TextDirection represents text direction.
type TextDirection int

const (
	// TextDirectionLTR represents left-to-right text.
	TextDirectionLTR TextDirection = iota
	// TextDirectionRTL represents right-to-left text.
	TextDirectionRTL
)

// LineBreaker handles line breaking and word wrapping.
type LineBreaker interface {
	// BreakLines breaks text into lines based on the given constraints.
	BreakLines(text string, font Font, maxWidth float32) []LineBreak

	// FindBreakOpportunities finds potential line break positions.
	FindBreakOpportunities(text string) []int
}

// LineBreak represents a line break in text.
type LineBreak struct {
	// Start is the start position in the text.
	Start int
	// End is the end position in the text.
	End int
	// Width is the width of the line.
	Width float32
	// IsHardBreak indicates if this is a hard line break.
	IsHardBreak bool
}

// defaultTextLayout provides the default implementation of TextLayout.
type defaultTextLayout struct {
	lineBreaker LineBreaker
}

// NewTextLayout creates a new text layout.
func NewTextLayout() TextLayout {
	return &defaultTextLayout{
		lineBreaker: NewLineBreaker(),
	}
}

// LayoutText implements TextLayout.
func (l *defaultTextLayout) LayoutText(text string, font Font, maxWidth float32) (TextFrame, error) {
	if font == nil || !font.IsValid() {
		return nil, errors.New("invalid font")
	}

	// Create a text run
	run := NewTextRun(font)
	if run == nil {
		return nil, errors.New("failed to create text run")
	}

	// Shape the text into glyphs
	err := l.shapeText(text, font, run)
	if err != nil {
		return nil, err
	}

	// Create a text frame
	frame := NewTextFrame()
	if frame == nil {
		return nil, errors.New("failed to create text frame")
	}

	frame.AddRun(run)
	return frame, nil
}

// LayoutParagraph implements TextLayout.
func (l *defaultTextLayout) LayoutParagraph(paragraph Paragraph, maxWidth float32) (TextFrame, error) {
	if paragraph == nil {
		return nil, errors.New("paragraph is nil")
	}

	frame := NewTextFrame()
	if frame == nil {
		return nil, errors.New("failed to create text frame")
	}

	runs := paragraph.GetTextRuns()
	for _, run := range runs {
		if run != nil && run.IsValid() {
			frame.AddRun(run)
		}
	}

	return frame, nil
}

// MeasureText implements TextLayout.
func (l *defaultTextLayout) MeasureText(text string, font Font) geom.Size[geom.F32] {
	if font == nil || !font.IsValid() {
		return geom.Size[geom.F32]{}
	}

	width := geom.F32(0)
	height := geom.F32(font.GetMetrics().LineHeight)

	typeface := font.GetTypeface()
	if typeface == nil {
		return geom.Size[geom.F32]{Width: width, Height: height}
	}

	for _, r := range text {
		if glyphIndex, ok := typeface.GetGlyphForCodepoint(r); ok {
			metrics := typeface.GetGlyphMetrics(glyphIndex)
			width += geom.F32(metrics.AdvanceWidth)
		}
	}

	return geom.Size[geom.F32]{Width: width, Height: height}
}

// GetLineHeight implements TextLayout.
func (l *defaultTextLayout) GetLineHeight(font Font) float32 {
	if font == nil || !font.IsValid() {
		return 0
	}
	return font.GetMetrics().LineHeight
}

// WordWrap implements TextLayout.
func (l *defaultTextLayout) WordWrap(text string, font Font, maxWidth float32) []string {
	if font == nil || !font.IsValid() {
		return []string{text}
	}

	lines := make([]string, 0)
	words := strings.Fields(text)

	if len(words) == 0 {
		return lines
	}

	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		size := l.MeasureText(testLine, font)
		if float32(size.Width) <= maxWidth || currentLine == "" {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// shapeText shapes text into glyphs and adds them to the text run.
func (l *defaultTextLayout) shapeText(text string, font Font, run TextRun) error {
	typeface := font.GetTypeface()
	if typeface == nil {
		return errors.New("font has no typeface")
	}

	x := geom.F32(0)
	y := geom.F32(0)

	for _, r := range text {
		glyphIndex, ok := typeface.GetGlyphForCodepoint(r)
		if !ok {
			continue
		}

		glyph := Glyph{
			Index: glyphIndex,
			Type:  GlyphTypePath, // Default assume path
		}

		position := geom.Point[geom.F32]{X: x, Y: y}
		run.AddGlyph(glyph, position)

		// Advance position
		metrics := typeface.GetGlyphMetrics(glyphIndex)
		x += geom.F32(metrics.AdvanceWidth)
	}

	return nil
}

// defaultParagraph provides the default implementation of Paragraph.
type defaultParagraph struct {
	runs      []TextRun
	alignment TextAlignment
	bounds    geom.Rect[geom.F32]
}

// NewParagraph creates a new paragraph.
func NewParagraph() Paragraph {
	return &defaultParagraph{
		runs:      make([]TextRun, 0),
		alignment: TextAlignmentLeft,
		bounds:    geom.Rect[geom.F32]{},
	}
}

// GetTextRuns implements Paragraph.
func (p *defaultParagraph) GetTextRuns() []TextRun {
	return p.runs
}

// AddTextRun implements Paragraph.
func (p *defaultParagraph) AddTextRun(run TextRun) {
	if run != nil && run.IsValid() {
		p.runs = append(p.runs, run)
	}
}

// GetText implements Paragraph.
func (p *defaultParagraph) GetText() string {
	var text strings.Builder
	for _, run := range p.runs {
		// TODO: Extract text from run
		// This would require storing the original text with the run
	}
	return text.String()
}

// GetBounds implements Paragraph.
func (p *defaultParagraph) GetBounds() geom.Rect[geom.F32] {
	return p.bounds
}

// GetAlignment implements Paragraph.
func (p *defaultParagraph) GetAlignment() TextAlignment {
	return p.alignment
}

// SetAlignment implements Paragraph.
func (p *defaultParagraph) SetAlignment(alignment TextAlignment) {
	p.alignment = alignment
}

// defaultLineBreaker provides the default implementation of LineBreaker.
type defaultLineBreaker struct{}

// NewLineBreaker creates a new line breaker.
func NewLineBreaker() LineBreaker {
	return &defaultLineBreaker{}
}

// BreakLines implements LineBreaker.
func (b *defaultLineBreaker) BreakLines(text string, font Font, maxWidth float32) []LineBreak {
	breaks := make([]LineBreak, 0)

	if text == "" {
		return breaks
	}

	layout := NewTextLayout()
	words := strings.Fields(text)

	start := 0
	currentLine := ""

	for i, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		size := layout.MeasureText(testLine, font)
		if float32(size.Width) <= maxWidth || currentLine == "" {
			currentLine = testLine
		} else {
			// Add line break
			if currentLine != "" {
				end := start + len(currentLine)
				breaks = append(breaks, LineBreak{
					Start:       start,
					End:         end,
					Width:       float32(layout.MeasureText(currentLine, font).Width),
					IsHardBreak: false,
				})
				start = end + 1 // +1 for space
			}
			currentLine = word
		}

		// Check for hard line breaks
		if strings.Contains(word, "\n") {
			parts := strings.Split(word, "\n")
			for j, part := range parts {
				if j > 0 {
					// Hard line break
					end := start + len(part)
					breaks = append(breaks, LineBreak{
						Start:       start,
						End:         end,
						Width:       float32(layout.MeasureText(part, font).Width),
						IsHardBreak: true,
					})
					start = end + 1
				}
			}
		}
	}

	// Add final line
	if currentLine != "" {
		end := start + len(currentLine)
		breaks = append(breaks, LineBreak{
			Start:       start,
			End:         end,
			Width:       float32(layout.MeasureText(currentLine, font).Width),
			IsHardBreak: false,
		})
	}

	return breaks
}

// FindBreakOpportunities implements LineBreaker.
func (b *defaultLineBreaker) FindBreakOpportunities(text string) []int {
	opportunities := make([]int, 0)

	for i, r := range text {
		if unicode.IsSpace(r) || r == '-' || r == '/' {
			opportunities = append(opportunities, i)
		}
	}

	return opportunities
}
