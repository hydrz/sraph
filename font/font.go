// Package font provides font loading and text layout functionality for the Sraph UI toolkit.
//
// This package handles font loading, glyph management, text measurement, and
// advanced text layout operations including paragraph layout and text shaping.
package font

import (
	"io"

	"github.com/opensraph/sraph/geom"
)

// Font represents a loaded font with various weights and styles.
type Font interface {
	// Name returns the font family name.
	Name() string

	// Weight returns the font weight.
	Weight() Weight

	// Style returns the font style.
	Style() Style

	// Size returns the font size in points.
	Size() float32

	// SetSize sets the font size in points.
	SetSize(size float32)

	// Ascent returns the font ascent in font units.
	Ascent() float32

	// Descent returns the font descent in font units.
	Descent() float32

	// LineHeight returns the recommended line height.
	LineHeight() float32

	// MeasureText measures the dimensions of the given text.
	MeasureText(text string) geom.Size[geom.F32]

	// GetGlyph returns the glyph for the given rune.
	GetGlyph(r rune) (Glyph, bool)
}

// Weight represents font weight values.
type Weight int

const (
	WeightThin       Weight = 100
	WeightExtraLight Weight = 200
	WeightLight      Weight = 300
	WeightNormal     Weight = 400
	WeightMedium     Weight = 500
	WeightSemiBold   Weight = 600
	WeightBold       Weight = 700
	WeightExtraBold  Weight = 800
	WeightBlack      Weight = 900
)

// Style represents font style values.
type Style int

const (
	StyleNormal Style = iota
	StyleItalic
	StyleOblique
)

// Glyph represents a single glyph in a font.
type Glyph interface {
	// ID returns the glyph ID.
	ID() uint32

	// Rune returns the Unicode rune this glyph represents.
	Rune() rune

	// Advance returns the advance width of this glyph.
	Advance() float32

	// Bounds returns the bounding box of this glyph.
	Bounds() geom.Rect[geom.F32]

	// Path returns the path data for this glyph.
	Path() geom.PathSource[geom.F32]
}

// Atlas represents a texture atlas containing rendered glyphs.
type Atlas interface {
	// Size returns the size of the atlas texture.
	Size() (width, height int)

	// Add adds a glyph to the atlas and returns its texture coordinates.
	Add(glyph Glyph, size float32) (TextureCoords, error)

	// Get returns the texture coordinates for a glyph if it exists in the atlas.
	Get(glyphID uint32, size float32) (TextureCoords, bool)

	// Texture returns the atlas texture data.
	Texture() []byte

	// Clear clears all glyphs from the atlas.
	Clear()
}

// TextureCoords represents texture coordinates for a glyph in an atlas.
type TextureCoords struct {
	// UV coordinates in the atlas texture (0-1 range).
	U1, V1, U2, V2 float32

	// Offset from the baseline to position the glyph.
	OffsetX, OffsetY float32

	// Size of the glyph in pixels.
	Width, Height float32
}

// Layout represents a text layout engine for complex text rendering.
type Layout interface {
	// SetFont sets the font for text layout.
	SetFont(font Font)

	// SetText sets the text to be laid out.
	SetText(text string)

	// SetMaxWidth sets the maximum width for text wrapping.
	SetMaxWidth(width float32)

	// SetAlignment sets the text alignment.
	SetAlignment(align Alignment)

	// Layout performs the text layout and returns the result.
	Layout() *LayoutResult
}

// Alignment represents text alignment options.
type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
	AlignJustify
)

// LayoutResult contains the result of text layout.
type LayoutResult struct {
	// Lines contains the laid out text lines.
	Lines []Line

	// Bounds contains the bounding box of the entire text.
	Bounds geom.Rect[geom.F32]

	// Baseline is the y-coordinate of the first line's baseline.
	Baseline float32
}

// Line represents a single line of laid out text.
type Line struct {
	// Glyphs contains the glyphs in this line.
	Glyphs []PositionedGlyph

	// Bounds contains the bounding box of this line.
	Bounds geom.Rect[geom.F32]

	// Baseline is the y-coordinate of this line's baseline.
	Baseline float32
}

// PositionedGlyph represents a glyph positioned in a text layout.
type PositionedGlyph struct {
	// Glyph is the glyph data.
	Glyph Glyph

	// Position is the position of this glyph.
	Position geom.Point[geom.F32]

	// Size is the size at which this glyph should be rendered.
	Size float32
}

// Paragraph represents a paragraph layout engine for advanced text rendering.
type Paragraph interface {
	// SetStyle sets the paragraph style.
	SetStyle(style *ParagraphStyle)

	// AddText adds text with the given style to the paragraph.
	AddText(text string, style *TextStyle)

	// Layout performs paragraph layout with the given constraints.
	Layout(constraints Constraints) *ParagraphResult
}

// ParagraphStyle defines styling for an entire paragraph.
type ParagraphStyle struct {
	// TextAlign specifies the text alignment.
	TextAlign Alignment

	// TextDirection specifies the text direction.
	TextDirection TextDirection

	// MaxLines specifies the maximum number of lines (0 for unlimited).
	MaxLines int

	// Ellipsis specifies the ellipsis string for text overflow.
	Ellipsis string

	// LineHeight specifies the line height multiplier.
	LineHeight float32
}

// TextStyle defines styling for a span of text.
type TextStyle struct {
	// Font specifies the font to use.
	Font Font

	// Color specifies the text color.
	Color geom.Color

	// FontSize specifies the font size in points.
	FontSize float32

	// FontWeight specifies the font weight.
	FontWeight Weight

	// FontStyle specifies the font style.
	FontStyle Style

	// Decoration specifies text decoration.
	Decoration TextDecoration

	// DecorationColor specifies the decoration color.
	DecorationColor geom.Color

	// DecorationThickness specifies the decoration thickness.
	DecorationThickness float32
}

// TextDirection represents text direction.
type TextDirection int

const (
	TextDirectionLTR TextDirection = iota // Left to right
	TextDirectionRTL                      // Right to left
)

// TextDecoration represents text decoration options.
type TextDecoration int

const (
	TextDecorationNone TextDecoration = iota
	TextDecorationUnderline
	TextDecorationOverline
	TextDecorationLineThrough
)

// Constraints represents layout constraints for paragraph layout.
type Constraints struct {
	// MaxWidth specifies the maximum width for the paragraph.
	MaxWidth float32

	// MaxHeight specifies the maximum height for the paragraph.
	MaxHeight float32
}

// ParagraphResult contains the result of paragraph layout.
type ParagraphResult struct {
	// Lines contains the laid out lines.
	Lines []Line

	// Bounds contains the bounding box of the entire paragraph.
	Bounds geom.Rect[geom.F32]

	// DidExceedMaxLines indicates if the text exceeded the maximum lines.
	DidExceedMaxLines bool
}

// FontLoader provides font loading functionality.
type FontLoader interface {
	// LoadFont loads a font from the given reader.
	LoadFont(r io.Reader) (Font, error)

	// LoadFontFromFile loads a font from a file path.
	LoadFontFromFile(path string) (Font, error)

	// LoadSystemFont loads a system font by name.
	LoadSystemFont(name string) (Font, error)

	// ListSystemFonts returns a list of available system fonts.
	ListSystemFonts() []string
}

// DefaultFontLoader returns the default font loader implementation.
func DefaultFontLoader() FontLoader {
	// TODO: Implement default font loader
	return nil
}

// NewAtlas creates a new glyph atlas with the given size.
func NewAtlas(width, height int) Atlas {
	// TODO: Implement atlas creation
	return nil
}

// NewLayout creates a new text layout engine.
func NewLayout() Layout {
	// TODO: Implement layout creation
	return nil
}

// NewParagraph creates a new paragraph layout engine.
func NewParagraph() Paragraph {
	// TODO: Implement paragraph creation
	return nil
}
