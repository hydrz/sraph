// Package entity provides text entity implementation.
package entity

import (
	"github.com/opensraph/sraph/geom"
)

// TextEntity represents an entity that renders text.
type TextEntity struct {
	*BaseEntity
	text          string
	fontID        FontID
	fontSize      geom.F32
	color         geom.Color
	alignment     TextAlignment
	lineHeight    geom.F32
	letterSpacing geom.F32
	wordWrap      bool
	maxWidth      geom.F32
}

// FontID represents a unique identifier for a font resource.
type FontID uint32

// TextAlignment defines text alignment options.
type TextAlignment uint8

const (
	// TextAlignmentLeft aligns text to the left.
	TextAlignmentLeft TextAlignment = iota
	// TextAlignmentCenter centers text horizontally.
	TextAlignmentCenter
	// TextAlignmentRight aligns text to the right.
	TextAlignmentRight
	// TextAlignmentJustify justifies text (spreads to full width).
	TextAlignmentJustify
)

// NewTextEntity creates a new text entity.
func NewTextEntity(bounds geom.Rect[geom.F32], text string, fontID FontID, fontSize geom.F32, color geom.Color) *TextEntity {
	entity := &TextEntity{
		BaseEntity:    NewBaseEntity(),
		text:          text,
		fontID:        fontID,
		fontSize:      fontSize,
		color:         color,
		alignment:     TextAlignmentLeft,
		lineHeight:    1.2,
		letterSpacing: 0,
		wordWrap:      true,
		maxWidth:      bounds.Width(),
	}
	entity.SetBounds(bounds)
	return entity
}

// Text returns the text content.
func (t *TextEntity) Text() string {
	return t.text
}

// SetText sets the text content.
func (t *TextEntity) SetText(text string) {
	t.text = text
}

// FontID returns the font ID for this text entity.
func (t *TextEntity) FontID() FontID {
	return t.fontID
}

// SetFontID sets the font ID for this text entity.
func (t *TextEntity) SetFontID(fontID FontID) {
	t.fontID = fontID
}

// FontSize returns the font size.
func (t *TextEntity) FontSize() geom.F32 {
	return t.fontSize
}

// SetFontSize sets the font size.
func (t *TextEntity) SetFontSize(size geom.F32) {
	t.fontSize = size
}

// Color returns the text color.
func (t *TextEntity) Color() geom.Color {
	return t.color
}

// SetColor sets the text color.
func (t *TextEntity) SetColor(color geom.Color) {
	t.color = color
}

// Alignment returns the text alignment.
func (t *TextEntity) Alignment() TextAlignment {
	return t.alignment
}

// SetAlignment sets the text alignment.
func (t *TextEntity) SetAlignment(alignment TextAlignment) {
	t.alignment = alignment
}

// LineHeight returns the line height multiplier.
func (t *TextEntity) LineHeight() geom.F32 {
	return t.lineHeight
}

// SetLineHeight sets the line height multiplier.
func (t *TextEntity) SetLineHeight(lineHeight geom.F32) {
	t.lineHeight = lineHeight
}

// LetterSpacing returns the letter spacing.
func (t *TextEntity) LetterSpacing() geom.F32 {
	return t.letterSpacing
}

// SetLetterSpacing sets the letter spacing.
func (t *TextEntity) SetLetterSpacing(spacing geom.F32) {
	t.letterSpacing = spacing
}

// WordWrap returns whether word wrapping is enabled.
func (t *TextEntity) WordWrap() bool {
	return t.wordWrap
}

// SetWordWrap sets whether word wrapping is enabled.
func (t *TextEntity) SetWordWrap(wrap bool) {
	t.wordWrap = wrap
}

// MaxWidth returns the maximum width for text layout.
func (t *TextEntity) MaxWidth() geom.F32 {
	return t.maxWidth
}

// SetMaxWidth sets the maximum width for text layout.
func (t *TextEntity) SetMaxWidth(width geom.F32) {
	t.maxWidth = width
}

// Render implements RenderableEntity.
func (t *TextEntity) Render(ctx RenderContext) error {
	// TODO: Implement text rendering
	// This would typically involve:
	// 1. Loading the font and measuring the text
	// 2. Performing text layout (line breaking, alignment)
	// 3. Rendering glyphs to texture atlas
	// 4. Drawing textured quads for each glyph
	// 5. Applying the entity's transformation matrix
	return nil
}

// RichTextEntity represents an entity that can render rich text with multiple styles.
type RichTextEntity struct {
	*BaseEntity
	segments []TextSegment
}

// TextSegment represents a segment of text with specific styling.
type TextSegment struct {
	Text          string
	FontID        FontID
	FontSize      geom.F32
	Color         geom.Color
	Bold          bool
	Italic        bool
	Underline     bool
	Strikethrough bool
	LetterSpacing geom.F32
}

// NewRichTextEntity creates a new rich text entity.
func NewRichTextEntity(bounds geom.Rect[geom.F32]) *RichTextEntity {
	entity := &RichTextEntity{
		BaseEntity: NewBaseEntity(),
		segments:   make([]TextSegment, 0),
	}
	entity.SetBounds(bounds)
	return entity
}

// AddSegment adds a text segment to the rich text entity.
func (r *RichTextEntity) AddSegment(segment TextSegment) {
	r.segments = append(r.segments, segment)
}

// Segments returns all text segments.
func (r *RichTextEntity) Segments() []TextSegment {
	return r.segments
}

// ClearSegments clears all text segments.
func (r *RichTextEntity) ClearSegments() {
	r.segments = r.segments[:0]
}

// Render implements RenderableEntity.
func (r *RichTextEntity) Render(ctx RenderContext) error {
	// TODO: Implement rich text rendering
	// This involves rendering each segment with its own styling
	return nil
}

// EditableTextEntity represents an entity for editable text input.
type EditableTextEntity struct {
	*TextEntity
	cursorPosition int
	selectionStart int
	selectionEnd   int
	cursorVisible  bool
	placeholder    string
	readOnly       bool
}

// NewEditableTextEntity creates a new editable text entity.
func NewEditableTextEntity(bounds geom.Rect[geom.F32], fontID FontID, fontSize geom.F32, color geom.Color) *EditableTextEntity {
	return &EditableTextEntity{
		TextEntity:     NewTextEntity(bounds, "", fontID, fontSize, color),
		cursorPosition: 0,
		selectionStart: 0,
		selectionEnd:   0,
		cursorVisible:  true,
		placeholder:    "",
		readOnly:       false,
	}
}

// CursorPosition returns the current cursor position.
func (e *EditableTextEntity) CursorPosition() int {
	return e.cursorPosition
}

// SetCursorPosition sets the cursor position.
func (e *EditableTextEntity) SetCursorPosition(position int) {
	if position < 0 {
		position = 0
	}
	if position > len(e.text) {
		position = len(e.text)
	}
	e.cursorPosition = position
}

// HasSelection returns whether there is an active text selection.
func (e *EditableTextEntity) HasSelection() bool {
	return e.selectionStart != e.selectionEnd
}

// GetSelection returns the start and end of the current selection.
func (e *EditableTextEntity) GetSelection() (int, int) {
	return e.selectionStart, e.selectionEnd
}

// SetSelection sets the text selection.
func (e *EditableTextEntity) SetSelection(start, end int) {
	if start < 0 {
		start = 0
	}
	if end > len(e.text) {
		end = len(e.text)
	}
	if start > end {
		start, end = end, start
	}
	e.selectionStart = start
	e.selectionEnd = end
}

// ClearSelection clears the current text selection.
func (e *EditableTextEntity) ClearSelection() {
	e.selectionStart = e.cursorPosition
	e.selectionEnd = e.cursorPosition
}

// CursorVisible returns whether the cursor is visible.
func (e *EditableTextEntity) CursorVisible() bool {
	return e.cursorVisible
}

// SetCursorVisible sets whether the cursor is visible.
func (e *EditableTextEntity) SetCursorVisible(visible bool) {
	e.cursorVisible = visible
}

// Placeholder returns the placeholder text.
func (e *EditableTextEntity) Placeholder() string {
	return e.placeholder
}

// SetPlaceholder sets the placeholder text.
func (e *EditableTextEntity) SetPlaceholder(placeholder string) {
	e.placeholder = placeholder
}

// ReadOnly returns whether the text is read-only.
func (e *EditableTextEntity) ReadOnly() bool {
	return e.readOnly
}

// SetReadOnly sets whether the text is read-only.
func (e *EditableTextEntity) SetReadOnly(readOnly bool) {
	e.readOnly = readOnly
}

// InsertText inserts text at the current cursor position.
func (e *EditableTextEntity) InsertText(text string) {
	if e.readOnly {
		return
	}

	// Replace selection if any
	if e.HasSelection() {
		e.DeleteSelection()
	}

	// Insert text
	before := e.text[:e.cursorPosition]
	after := e.text[e.cursorPosition:]
	e.text = before + text + after
	e.cursorPosition += len(text)
	e.ClearSelection()
}

// DeleteSelection deletes the currently selected text.
func (e *EditableTextEntity) DeleteSelection() {
	if !e.HasSelection() || e.readOnly {
		return
	}

	before := e.text[:e.selectionStart]
	after := e.text[e.selectionEnd:]
	e.text = before + after
	e.cursorPosition = e.selectionStart
	e.ClearSelection()
}

// Render implements RenderableEntity.
func (e *EditableTextEntity) Render(ctx RenderContext) error {
	// TODO: Implement editable text rendering
	// This includes rendering the text, cursor, and selection highlight
	return nil
}
