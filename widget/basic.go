// Package main provides core widget implementations for the Sraph UI toolkit.
package widget

import (
	"github.com/opensraph/sraph/font"
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// Text is a widget that displays a string of text with a single style.
type Text struct {
	key   Key
	data  string
	style *TextStyle
}

// NewText creates a new Text widget.
func NewText(data string, style *TextStyle) *Text {
	return &Text{
		data:  data,
		style: style,
	}
}

// Build implements Widget.
func (t *Text) Build(context BuildContext) Widget {
	// Text is a leaf widget, so it doesn't build children
	return t
}

// Key implements Widget.
func (t *Text) Key() Key {
	return t.key
}

// SetKey implements Widget.
func (t *Text) SetKey(key Key) {
	t.key = key
}

// Data returns the text data.
func (t *Text) Data() string {
	return t.data
}

// Style returns the text style.
func (t *Text) Style() *TextStyle {
	return t.style
}

// TextStyle defines styling for text widgets.
type TextStyle struct {
	Color      geom.Color
	FontSize   float32
	FontWeight font.Weight
	FontStyle  font.Style
	FontFamily string
}

// Container is a widget that contains a single child and applies
// decoration, positioning, and sizing constraints.
type Container struct {
	key        Key
	child      Widget
	width      *float32
	height     *float32
	padding    EdgeInsets
	margin     EdgeInsets
	color      *geom.Color
	decoration Decoration
	alignment  Alignment
}

// NewContainer creates a new Container widget.
func NewContainer() *Container {
	return &Container{}
}

// Build implements Widget.
func (c *Container) Build(context BuildContext) Widget {
	return c.child
}

// Key implements Widget.
func (c *Container) Key() Key {
	return c.key
}

// SetKey implements Widget.
func (c *Container) SetKey(key Key) {
	c.key = key
}

// Child implements Container.
func (c *Container) Child() Widget {
	return c.child
}

// SetChild implements Container.
func (c *Container) SetChild(child Widget) {
	c.child = child
}

// WithWidth sets the width constraint.
func (c *Container) WithWidth(width float32) *Container {
	c.width = &width
	return c
}

// WithHeight sets the height constraint.
func (c *Container) WithHeight(height float32) *Container {
	c.height = &height
	return c
}

// WithPadding sets the padding.
func (c *Container) WithPadding(padding EdgeInsets) *Container {
	c.padding = padding
	return c
}

// WithMargin sets the margin.
func (c *Container) WithMargin(margin EdgeInsets) *Container {
	c.margin = margin
	return c
}

// WithColor sets the background color.
func (c *Container) WithColor(color geom.Color) *Container {
	c.color = &color
	return c
}

// WithAlignment sets the alignment of the child.
func (c *Container) WithAlignment(alignment Alignment) *Container {
	c.alignment = alignment
	return c
}

// EdgeInsets represents insets from each of the four sides.
type EdgeInsets struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

// All creates EdgeInsets with the same value on all sides.
func (EdgeInsets) All(value float32) EdgeInsets {
	return EdgeInsets{
		Top:    value,
		Right:  value,
		Bottom: value,
		Left:   value,
	}
}

// Symmetric creates EdgeInsets with symmetric horizontal and vertical values.
func (EdgeInsets) Symmetric(horizontal, vertical float32) EdgeInsets {
	return EdgeInsets{
		Top:    vertical,
		Right:  horizontal,
		Bottom: vertical,
		Left:   horizontal,
	}
}

// Only creates EdgeInsets with specific values for each side.
func (EdgeInsets) Only(top, right, bottom, left float32) EdgeInsets {
	return EdgeInsets{
		Top:    top,
		Right:  right,
		Bottom: bottom,
		Left:   left,
	}
}

// Decoration defines how to paint a decoration.
type Decoration interface {
	// CreateBoxPainter creates a box painter for this decoration.
	CreateBoxPainter() BoxPainter
}

// BoxPainter paints a decoration on a render object.
type BoxPainter interface {
	// Paint paints the decoration.
	Paint(canvas render.Canvas, rect geom.Rect[geom.F32])
}

// Alignment specifies how to align a child within its parent.
type Alignment struct {
	X, Y float32
}

var (
	// Common alignment constants
	AlignmentTopLeft      = Alignment{X: -1.0, Y: -1.0}
	AlignmentTopCenter    = Alignment{X: 0.0, Y: -1.0}
	AlignmentTopRight     = Alignment{X: 1.0, Y: -1.0}
	AlignmentCenterLeft   = Alignment{X: -1.0, Y: 0.0}
	AlignmentCenter       = Alignment{X: 0.0, Y: 0.0}
	AlignmentCenterRight  = Alignment{X: 1.0, Y: 0.0}
	AlignmentBottomLeft   = Alignment{X: -1.0, Y: 1.0}
	AlignmentBottomCenter = Alignment{X: 0.0, Y: 1.0}
	AlignmentBottomRight  = Alignment{X: 1.0, Y: 1.0}
)

// Column is a widget that displays its children in a vertical array.
type Column struct {
	key                Key
	children           []Widget
	mainAxisAlignment  MainAxisAlignment
	crossAxisAlignment CrossAxisAlignment
	mainAxisSize       MainAxisSize
}

// NewColumn creates a new Column widget.
func NewColumn(children ...Widget) *Column {
	return &Column{
		children:           children,
		mainAxisAlignment:  MainAxisAlignmentStart,
		crossAxisAlignment: CrossAxisAlignmentCenter,
		mainAxisSize:       MainAxisSizeMax,
	}
}

// Build implements Widget.
func (c *Column) Build(context BuildContext) Widget {
	// Column is a multi-child widget
	return c
}

// Key implements Widget.
func (c *Column) Key() Key {
	return c.key
}

// SetKey implements Widget.
func (c *Column) SetKey(key Key) {
	c.key = key
}

// Children implements MultiChildWidget.
func (c *Column) Children() []Widget {
	return c.children
}

// SetChildren implements MultiChildWidget.
func (c *Column) SetChildren(children []Widget) {
	c.children = children
}

// AddChild implements MultiChildWidget.
func (c *Column) AddChild(child Widget) {
	c.children = append(c.children, child)
}

// RemoveChild implements MultiChildWidget.
func (c *Column) RemoveChild(child Widget) {
	for i, c := range c.children {
		if c == child {
			c.children = append(c.children[:i], c.children[i+1:]...)
			break
		}
	}
}

// WithMainAxisAlignment sets the main axis alignment.
func (c *Column) WithMainAxisAlignment(alignment MainAxisAlignment) *Column {
	c.mainAxisAlignment = alignment
	return c
}

// WithCrossAxisAlignment sets the cross axis alignment.
func (c *Column) WithCrossAxisAlignment(alignment CrossAxisAlignment) *Column {
	c.crossAxisAlignment = alignment
	return c
}

// Row is a widget that displays its children in a horizontal array.
type Row struct {
	key                Key
	children           []Widget
	mainAxisAlignment  MainAxisAlignment
	crossAxisAlignment CrossAxisAlignment
	mainAxisSize       MainAxisSize
}

// NewRow creates a new Row widget.
func NewRow(children ...Widget) *Row {
	return &Row{
		children:           children,
		mainAxisAlignment:  MainAxisAlignmentStart,
		crossAxisAlignment: CrossAxisAlignmentCenter,
		mainAxisSize:       MainAxisSizeMax,
	}
}

// Build implements Widget.
func (r *Row) Build(context BuildContext) Widget {
	// Row is a multi-child widget
	return r
}

// Key implements Widget.
func (r *Row) Key() Key {
	return r.key
}

// SetKey implements Widget.
func (r *Row) SetKey(key Key) {
	r.key = key
}

// Children implements MultiChildWidget.
func (r *Row) Children() []Widget {
	return r.children
}

// SetChildren implements MultiChildWidget.
func (r *Row) SetChildren(children []Widget) {
	r.children = children
}

// AddChild implements MultiChildWidget.
func (r *Row) AddChild(child Widget) {
	r.children = append(r.children, child)
}

// RemoveChild implements MultiChildWidget.
func (r *Row) RemoveChild(child Widget) {
	for i, c := range r.children {
		if c == child {
			r.children = append(r.children[:i], r.children[i+1:]...)
			break
		}
	}
}

// MainAxisAlignment defines how children are aligned along the main axis.
type MainAxisAlignment int

const (
	MainAxisAlignmentStart        MainAxisAlignment = iota // Place children at the start
	MainAxisAlignmentEnd                                   // Place children at the end
	MainAxisAlignmentCenter                                // Place children at the center
	MainAxisAlignmentSpaceBetween                          // Space evenly with no space at ends
	MainAxisAlignmentSpaceAround                           // Space evenly with half space at ends
	MainAxisAlignmentSpaceEvenly                           // Space evenly including ends
)

// CrossAxisAlignment defines how children are aligned along the cross axis.
type CrossAxisAlignment int

const (
	CrossAxisAlignmentStart    CrossAxisAlignment = iota // Align to the start
	CrossAxisAlignmentEnd                                // Align to the end
	CrossAxisAlignmentCenter                             // Align to the center
	CrossAxisAlignmentStretch                            // Stretch to fill
	CrossAxisAlignmentBaseline                           // Align to text baseline
)

// MainAxisSize defines how much space the widget should occupy in the main axis.
type MainAxisSize int

const (
	MainAxisSizeMin MainAxisSize = iota // Use minimum space needed
	MainAxisSizeMax                     // Use maximum space available
)

// Expanded is a widget that expands a child of a Row, Column, or Flex.
type Expanded struct {
	key   Key
	child Widget
	flex  int
}

// NewExpanded creates a new Expanded widget.
func NewExpanded(child Widget) *Expanded {
	return &Expanded{
		child: child,
		flex:  1,
	}
}

// Build implements Widget.
func (e *Expanded) Build(context BuildContext) Widget {
	return e.child
}

// Key implements Widget.
func (e *Expanded) Key() Key {
	return e.key
}

// SetKey implements Widget.
func (e *Expanded) SetKey(key Key) {
	e.key = key
}

// Child implements Container.
func (e *Expanded) Child() Widget {
	return e.child
}

// SetChild implements Container.
func (e *Expanded) SetChild(child Widget) {
	e.child = child
}

// WithFlex sets the flex factor.
func (e *Expanded) WithFlex(flex int) *Expanded {
	e.flex = flex
	return e
}

// Flex returns the flex factor.
func (e *Expanded) Flex() int {
	return e.flex
}
