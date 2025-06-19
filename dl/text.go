package dl

import (
	"fmt"

	"github.com/opensraph/sraph/geom"
)

// TextDirection defines the direction of text flow.
type TextDirection uint8

const (
	// TextDirectionLTR represents left-to-right text direction.
	TextDirectionLTR TextDirection = iota
	// TextDirectionRTL represents right-to-left text direction.
	TextDirectionRTL
)

// TextAlign defines how text is aligned within its bounds.
type TextAlign uint8

const (
	// TextAlignLeft aligns text to the left.
	TextAlignLeft TextAlign = iota
	// TextAlignRight aligns text to the right.
	TextAlignRight
	// TextAlignCenter centers text.
	TextAlignCenter
	// TextAlignJustify justifies text (stretches to fill width).
	TextAlignJustify
	// TextAlignStart aligns text to the start (left for LTR, right for RTL).
	TextAlignStart
	// TextAlignEnd aligns text to the end (right for LTR, left for RTL).
	TextAlignEnd
)

// FontWeight defines the weight (boldness) of a font.
type FontWeight uint16

const (
	// FontWeightThin represents thin font weight (100).
	FontWeightThin FontWeight = 100
	// FontWeightExtraLight represents extra light font weight (200).
	FontWeightExtraLight FontWeight = 200
	// FontWeightLight represents light font weight (300).
	FontWeightLight FontWeight = 300
	// FontWeightNormal represents normal font weight (400).
	FontWeightNormal FontWeight = 400
	// FontWeightMedium represents medium font weight (500).
	FontWeightMedium FontWeight = 500
	// FontWeightSemiBold represents semi-bold font weight (600).
	FontWeightSemiBold FontWeight = 600
	// FontWeightBold represents bold font weight (700).
	FontWeightBold FontWeight = 700
	// FontWeightExtraBold represents extra bold font weight (800).
	FontWeightExtraBold FontWeight = 800
	// FontWeightBlack represents black font weight (900).
	FontWeightBlack FontWeight = 900
)

// FontStyle defines the style of a font.
type FontStyle uint8

const (
	// FontStyleNormal represents normal font style.
	FontStyleNormal FontStyle = iota
	// FontStyleItalic represents italic font style.
	FontStyleItalic
	// FontStyleOblique represents oblique font style.
	FontStyleOblique
)

// TextBaseline defines the baseline for text alignment.
type TextBaseline uint8

const (
	// TextBaselineAlphabetic aligns to the alphabetic baseline.
	TextBaselineAlphabetic TextBaseline = iota
	// TextBaselineIdeographic aligns to the ideographic baseline.
	TextBaselineIdeographic
)

// TextDecoration defines text decoration options.
type TextDecoration uint8

const (
	// TextDecorationNone represents no text decoration.
	TextDecorationNone TextDecoration = 0
	// TextDecorationUnderline represents underlined text.
	TextDecorationUnderline TextDecoration = 1 << iota
	// TextDecorationOverline represents overlined text.
	TextDecorationOverline
	// TextDecorationLineThrough represents strike-through text.
	TextDecorationLineThrough
)

// TextDecorationStyle defines the style of text decoration.
type TextDecorationStyle uint8

const (
	// TextDecorationStyleSolid represents solid decoration line.
	TextDecorationStyleSolid TextDecorationStyle = iota
	// TextDecorationStyleDouble represents double decoration line.
	TextDecorationStyleDouble
	// TextDecorationStyleDotted represents dotted decoration line.
	TextDecorationStyleDotted
	// TextDecorationStyleDashed represents dashed decoration line.
	TextDecorationStyleDashed
	// TextDecorationStyleWavy represents wavy decoration line.
	TextDecorationStyleWavy
)

// Font represents a font with specific characteristics.
type Font struct {
	family   string
	size     float32
	weight   FontWeight
	style    FontStyle
	baseline TextBaseline
}

// NewFont creates a new font with the specified properties.
func NewFont(family string, size float32) Font {
	return Font{
		family:   family,
		size:     size,
		weight:   FontWeightNormal,
		style:    FontStyleNormal,
		baseline: TextBaselineAlphabetic,
	}
}

// Family returns the font family name.
func (f Font) Family() string {
	return f.family
}

// Size returns the font size.
func (f Font) Size() float32 {
	return f.size
}

// Weight returns the font weight.
func (f Font) Weight() FontWeight {
	return f.weight
}

// Style returns the font style.
func (f Font) Style() FontStyle {
	return f.style
}

// Baseline returns the text baseline.
func (f Font) Baseline() TextBaseline {
	return f.baseline
}

// WithFamily returns a new Font with the specified family.
func (f Font) WithFamily(family string) Font {
	f.family = family
	return f
}

// WithSize returns a new Font with the specified size.
func (f Font) WithSize(size float32) Font {
	f.size = size
	return f
}

// WithWeight returns a new Font with the specified weight.
func (f Font) WithWeight(weight FontWeight) Font {
	f.weight = weight
	return f
}

// WithStyle returns a new Font with the specified style.
func (f Font) WithStyle(style FontStyle) Font {
	f.style = style
	return f
}

// WithBaseline returns a new Font with the specified baseline.
func (f Font) WithBaseline(baseline TextBaseline) Font {
	f.baseline = baseline
	return f
}

// String returns a string representation of the font.
func (f Font) String() string {
	return fmt.Sprintf("Font{%s, %.1fpt, weight: %d, style: %d}",
		f.family, f.size, f.weight, f.style)
}

// TextStyle defines the visual appearance of text.
type TextStyle struct {
	font            Font
	color           Color
	backgroundColor *Color
	decoration      TextDecoration
	decorationColor *Color
	decorationStyle TextDecorationStyle
	letterSpacing   float32
	wordSpacing     float32
	height          float32 // Line height multiplier
	shadows         []Shadow
}

// NewTextStyle creates a new text style with default values.
func NewTextStyle() TextStyle {
	return TextStyle{
		font:            NewFont("", 14.0),
		color:           ColorBlack,
		decoration:      TextDecorationNone,
		decorationStyle: TextDecorationStyleSolid,
		height:          1.0,
	}
}

// Font returns the font.
func (ts TextStyle) Font() Font {
	return ts.font
}

// Color returns the text color.
func (ts TextStyle) Color() Color {
	return ts.color
}

// BackgroundColor returns the background color, if any.
func (ts TextStyle) BackgroundColor() *Color {
	return ts.backgroundColor
}

// Decoration returns the text decoration.
func (ts TextStyle) Decoration() TextDecoration {
	return ts.decoration
}

// DecorationColor returns the decoration color, if any.
func (ts TextStyle) DecorationColor() *Color {
	return ts.decorationColor
}

// DecorationStyle returns the decoration style.
func (ts TextStyle) DecorationStyle() TextDecorationStyle {
	return ts.decorationStyle
}

// LetterSpacing returns the letter spacing.
func (ts TextStyle) LetterSpacing() float32 {
	return ts.letterSpacing
}

// WordSpacing returns the word spacing.
func (ts TextStyle) WordSpacing() float32 {
	return ts.wordSpacing
}

// Height returns the line height multiplier.
func (ts TextStyle) Height() float32 {
	return ts.height
}

// Shadows returns the text shadows.
func (ts TextStyle) Shadows() []Shadow {
	return ts.shadows
}

// WithFont returns a new TextStyle with the specified font.
func (ts TextStyle) WithFont(font Font) TextStyle {
	ts.font = font
	return ts
}

// WithColor returns a new TextStyle with the specified color.
func (ts TextStyle) WithColor(color Color) TextStyle {
	ts.color = color
	return ts
}

// WithBackgroundColor returns a new TextStyle with the specified background color.
func (ts TextStyle) WithBackgroundColor(color Color) TextStyle {
	ts.backgroundColor = &color
	return ts
}

// WithDecoration returns a new TextStyle with the specified decoration.
func (ts TextStyle) WithDecoration(decoration TextDecoration) TextStyle {
	ts.decoration = decoration
	return ts
}

// WithDecorationColor returns a new TextStyle with the specified decoration color.
func (ts TextStyle) WithDecorationColor(color Color) TextStyle {
	ts.decorationColor = &color
	return ts
}

// WithDecorationStyle returns a new TextStyle with the specified decoration style.
func (ts TextStyle) WithDecorationStyle(style TextDecorationStyle) TextStyle {
	ts.decorationStyle = style
	return ts
}

// WithLetterSpacing returns a new TextStyle with the specified letter spacing.
func (ts TextStyle) WithLetterSpacing(spacing float32) TextStyle {
	ts.letterSpacing = spacing
	return ts
}

// WithWordSpacing returns a new TextStyle with the specified word spacing.
func (ts TextStyle) WithWordSpacing(spacing float32) TextStyle {
	ts.wordSpacing = spacing
	return ts
}

// WithHeight returns a new TextStyle with the specified line height.
func (ts TextStyle) WithHeight(height float32) TextStyle {
	ts.height = height
	return ts
}

// WithShadows returns a new TextStyle with the specified shadows.
func (ts TextStyle) WithShadows(shadows []Shadow) TextStyle {
	ts.shadows = shadows
	return ts
}

// Shadow represents a text shadow.
type Shadow struct {
	offset     geom.Point[Scalar]
	blurRadius float32
	color      Color
}

// NewShadow creates a new shadow with the specified properties.
func NewShadow(offset geom.Point[Scalar], blurRadius float32, color Color) Shadow {
	return Shadow{
		offset:     offset,
		blurRadius: blurRadius,
		color:      color,
	}
}

// Offset returns the shadow offset.
func (s Shadow) Offset() geom.Point[Scalar] {
	return s.offset
}

// BlurRadius returns the shadow blur radius.
func (s Shadow) BlurRadius() float32 {
	return s.blurRadius
}

// Color returns the shadow color.
func (s Shadow) Color() Color {
	return s.color
}

// TextSpan represents a span of text with consistent styling.
type TextSpan struct {
	text     string
	style    *TextStyle
	children []*TextSpan
}

// NewTextSpan creates a new text span.
func NewTextSpan(text string, style *TextStyle) *TextSpan {
	return &TextSpan{
		text:  text,
		style: style,
	}
}

// Text returns the text content.
func (ts *TextSpan) Text() string {
	return ts.text
}

// Style returns the text style.
func (ts *TextSpan) Style() *TextStyle {
	return ts.style
}

// Children returns the child text spans.
func (ts *TextSpan) Children() []*TextSpan {
	return ts.children
}

// AddChild adds a child text span.
func (ts *TextSpan) AddChild(child *TextSpan) {
	ts.children = append(ts.children, child)
}

// Paragraph represents a paragraph of styled text.
type Paragraph struct {
	spans    []*TextSpan
	style    ParagraphStyle
	maxWidth float32
	built    bool
	// Layout information
	lines    []TextLine
	height   float32
	width    float32
	baseline float32
}

// ParagraphStyle defines paragraph-level styling.
type ParagraphStyle struct {
	textAlign     TextAlign
	textDirection TextDirection
	maxLines      *int
	ellipsis      *string
	fontSize      float32
	fontWeight    FontWeight
	fontStyle     FontStyle
	fontFamily    string
	height        float32
	textBaseline  TextBaseline
}

// NewParagraphStyle creates a new paragraph style with default values.
func NewParagraphStyle() ParagraphStyle {
	return ParagraphStyle{
		textAlign:     TextAlignStart,
		textDirection: TextDirectionLTR,
		fontSize:      14.0,
		fontWeight:    FontWeightNormal,
		fontStyle:     FontStyleNormal,
		height:        1.0,
		textBaseline:  TextBaselineAlphabetic,
	}
}

// TextLine represents a line of text within a paragraph.
type TextLine struct {
	spans    []*TextSpan
	bounds   geom.Rect[Scalar]
	baseline float32
	ascent   float32
	descent  float32
	leading  float32
}

// NewParagraph creates a new paragraph.
func NewParagraph() *Paragraph {
	return &Paragraph{
		style: NewParagraphStyle(),
	}
}

// AddTextSpan adds a text span to the paragraph.
func (p *Paragraph) AddTextSpan(span *TextSpan) {
	p.spans = append(p.spans, span)
	p.built = false
}

// Layout performs text layout with the specified constraints.
func (p *Paragraph) Layout(maxWidth float32) {
	p.maxWidth = maxWidth
	p.layout()
	p.built = true
}

// Height returns the height of the laid out paragraph.
func (p *Paragraph) Height() float32 {
	if !p.built {
		return 0
	}
	return p.height
}

// Width returns the width of the laid out paragraph.
func (p *Paragraph) Width() float32 {
	if !p.built {
		return 0
	}
	return p.width
}

// Baseline returns the baseline of the first line.
func (p *Paragraph) Baseline() float32 {
	if !p.built {
		return 0
	}
	return p.baseline
}

// Lines returns the text lines.
func (p *Paragraph) Lines() []TextLine {
	return p.lines
}

// Paint renders the paragraph at the specified position.
func (p *Paragraph) Paint(canvas Canvas, position geom.Point[Scalar]) {
	// TODO: Implement paragraph painting
}

// layout performs the text layout algorithm.
func (p *Paragraph) layout() {
	// TODO: Implement text layout algorithm
	// This is a complex process that involves:
	// 1. Text shaping (converting text to glyphs)
	// 2. Line breaking
	// 3. Bidirectional text processing
	// 4. Font fallback
	// 5. Positioning glyphs

	// For now, create a simple single-line layout
	if len(p.spans) > 0 {
		p.lines = []TextLine{
			{
				spans:    p.spans,
				bounds:   geom.NewRect[Scalar](0, 0, Scalar(p.maxWidth), Scalar(p.style.fontSize*p.style.height)),
				baseline: p.style.fontSize * 0.8, // Approximate baseline
				ascent:   p.style.fontSize * 0.8,
				descent:  p.style.fontSize * 0.2,
				leading:  0,
			},
		}
		p.height = p.style.fontSize * p.style.height
		p.width = p.maxWidth
		p.baseline = p.style.fontSize * 0.8
	}
}

// ParagraphBuilder helps build paragraphs with styled text.
type ParagraphBuilder struct {
	style   ParagraphStyle
	spans   []*TextSpan
	current *TextStyle
	stack   []*TextStyle
}

// NewParagraphBuilder creates a new paragraph builder.
func NewParagraphBuilder(style ParagraphStyle) *ParagraphBuilder {
	return &ParagraphBuilder{
		style: style,
	}
}

// PushStyle pushes a new text style onto the style stack.
func (pb *ParagraphBuilder) PushStyle(style TextStyle) {
	pb.stack = append(pb.stack, pb.current)
	pb.current = &style
}

// PopStyle pops the current text style from the style stack.
func (pb *ParagraphBuilder) PopStyle() {
	if len(pb.stack) > 0 {
		pb.current = pb.stack[len(pb.stack)-1]
		pb.stack = pb.stack[:len(pb.stack)-1]
	} else {
		pb.current = nil
	}
}

// AddText adds text with the current style.
func (pb *ParagraphBuilder) AddText(text string) {
	span := NewTextSpan(text, pb.current)
	pb.spans = append(pb.spans, span)
}

// Build creates the final paragraph.
func (pb *ParagraphBuilder) Build() *Paragraph {
	paragraph := NewParagraph()
	paragraph.style = pb.style
	paragraph.spans = pb.spans
	return paragraph
}
