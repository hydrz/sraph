// Package font - Typography context functionality
package font

import (
	"errors"
	"sync"

	"github.com/opensraph/sraph/render"
)

// TypographerContext provides the graphics context necessary to render text.
// It manages resources related to rendering text on the GPU, inspired by Impeller's TypographerContext.
type TypographerContext interface {
	// IsValid returns whether this context is valid.
	IsValid() bool

	// CreateGlyphAtlasContext creates a glyph atlas context for the given type.
	CreateGlyphAtlasContext(atlasType GlyphAtlasType) (GlyphAtlasContext, error)

	// GetGlyphRenderer returns the glyph renderer.
	GetGlyphRenderer() GlyphRenderer

	// GetRenderContext returns the underlying render context.
	GetRenderContext() render.Context

	// CreateTextFrame creates a text frame from the given text and font.
	CreateTextFrame(text string, font Font) (TextFrame, error)

	// RenderTextFrame renders a text frame.
	RenderTextFrame(frame TextFrame, pass render.RenderPass) error
}

// LazyGlyphAtlas provides lazy loading of glyph atlases.
// It creates atlases on demand and manages their lifecycle.
type LazyGlyphAtlas interface {
	// GetGlyphAtlas returns the glyph atlas, creating it if necessary.
	GetGlyphAtlas(atlasType GlyphAtlasType) (GlyphAtlas, error)

	// HasGlyphAtlas returns whether an atlas of the given type exists.
	HasGlyphAtlas(atlasType GlyphAtlasType) bool

	// CreateGlyphAtlas creates a new glyph atlas of the given type.
	CreateGlyphAtlas(atlasType GlyphAtlasType) (GlyphAtlas, error)

	// GetTypographerContext returns the typographer context.
	GetTypographerContext() TypographerContext
}

// TextRenderer handles rendering of text frames to render passes.
type TextRenderer interface {
	// RenderTextFrame renders a text frame to a render pass.
	RenderTextFrame(frame TextFrame, pass render.RenderPass) error

	// RenderTextRun renders a text run to a render pass.
	RenderTextRun(run TextRun, pass render.RenderPass) error

	// PrepareGlyphs ensures all glyphs in the frame are loaded in atlases.
	PrepareGlyphs(frame TextFrame) error

	// GetShader returns the text rendering shader.
	GetShader() render.Shader
}

// defaultTypographerContext provides the default implementation of TypographerContext.
type defaultTypographerContext struct {
	renderContext render.Context
	glyphRenderer GlyphRenderer
	textLayout    TextLayout
	textRenderer  TextRenderer
	isValid       bool
	mutex         sync.RWMutex
}

// NewTypographerContext creates a new typographer context.
func NewTypographerContext(renderContext render.Context) (TypographerContext, error) {
	if renderContext == nil {
		return nil, errors.New("render context is required")
	}

	glyphRenderer, err := NewGlyphRenderer(renderContext)
	if err != nil {
		return nil, err
	}

	return &defaultTypographerContext{
		renderContext: renderContext,
		glyphRenderer: glyphRenderer,
		textLayout:    NewTextLayout(),
		textRenderer:  NewTextRenderer(renderContext),
		isValid:       true,
	}, nil
}

// IsValid implements TypographerContext.
func (c *defaultTypographerContext) IsValid() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isValid
}

// CreateGlyphAtlasContext implements TypographerContext.
func (c *defaultTypographerContext) CreateGlyphAtlasContext(atlasType GlyphAtlasType) (GlyphAtlasContext, error) {
	return NewGlyphAtlasContext(atlasType, c.renderContext)
}

// GetGlyphRenderer implements TypographerContext.
func (c *defaultTypographerContext) GetGlyphRenderer() GlyphRenderer {
	return c.glyphRenderer
}

// GetRenderContext implements TypographerContext.
func (c *defaultTypographerContext) GetRenderContext() render.Context {
	return c.renderContext
}

// CreateTextFrame implements TypographerContext.
func (c *defaultTypographerContext) CreateTextFrame(text string, font Font) (TextFrame, error) {
	return c.textLayout.LayoutText(text, font, 0) // No width limit for single line
}

// RenderTextFrame implements TypographerContext.
func (c *defaultTypographerContext) RenderTextFrame(frame TextFrame, pass render.RenderPass) error {
	return c.textRenderer.RenderTextFrame(frame, pass)
}

// defaultLazyGlyphAtlas provides the default implementation of LazyGlyphAtlas.
type defaultLazyGlyphAtlas struct {
	typographerContext TypographerContext
	alphaContext       GlyphAtlasContext
	colorContext       GlyphAtlasContext
	mutex              sync.RWMutex
}

// NewLazyGlyphAtlas creates a new lazy glyph atlas.
func NewLazyGlyphAtlas(typographerContext TypographerContext) LazyGlyphAtlas {
	return &defaultLazyGlyphAtlas{
		typographerContext: typographerContext,
	}
}

// GetGlyphAtlas implements LazyGlyphAtlas.
func (a *defaultLazyGlyphAtlas) GetGlyphAtlas(atlasType GlyphAtlasType) (GlyphAtlas, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	var context GlyphAtlasContext
	var err error

	switch atlasType {
	case GlyphAtlasTypeAlpha:
		if a.alphaContext == nil {
			a.alphaContext, err = a.typographerContext.CreateGlyphAtlasContext(atlasType)
			if err != nil {
				return nil, err
			}
		}
		context = a.alphaContext
	case GlyphAtlasTypeColor:
		if a.colorContext == nil {
			a.colorContext, err = a.typographerContext.CreateGlyphAtlasContext(atlasType)
			if err != nil {
				return nil, err
			}
		}
		context = a.colorContext
	default:
		return nil, errors.New("unknown atlas type")
	}

	return context.GetGlyphAtlas(), nil
}

// HasGlyphAtlas implements LazyGlyphAtlas.
func (a *defaultLazyGlyphAtlas) HasGlyphAtlas(atlasType GlyphAtlasType) bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	switch atlasType {
	case GlyphAtlasTypeAlpha:
		return a.alphaContext != nil
	case GlyphAtlasTypeColor:
		return a.colorContext != nil
	default:
		return false
	}
}

// CreateGlyphAtlas implements LazyGlyphAtlas.
func (a *defaultLazyGlyphAtlas) CreateGlyphAtlas(atlasType GlyphAtlasType) (GlyphAtlas, error) {
	context, err := a.typographerContext.CreateGlyphAtlasContext(atlasType)
	if err != nil {
		return nil, err
	}

	return context.GetGlyphAtlas(), nil
}

// GetTypographerContext implements LazyGlyphAtlas.
func (a *defaultLazyGlyphAtlas) GetTypographerContext() TypographerContext {
	return a.typographerContext
}

// defaultTextRenderer provides the default implementation of TextRenderer.
type defaultTextRenderer struct {
	renderContext render.Context
	shader        render.Shader
}

// NewTextRenderer creates a new text renderer.
func NewTextRenderer(renderContext render.Context) TextRenderer {
	// TODO: Load text rendering shader
	var shader render.Shader

	return &defaultTextRenderer{
		renderContext: renderContext,
		shader:        shader,
	}
}

// RenderTextFrame implements TextRenderer.
func (r *defaultTextRenderer) RenderTextFrame(frame TextFrame, pass render.RenderPass) error {
	if frame == nil || !frame.IsValid() {
		return errors.New("invalid text frame")
	}

	// Prepare all glyphs
	err := r.PrepareGlyphs(frame)
	if err != nil {
		return err
	}

	// Render each run
	runs := frame.GetRuns()
	for _, run := range runs {
		err := r.RenderTextRun(run, pass)
		if err != nil {
			return err
		}
	}

	return nil
}

// RenderTextRun implements TextRenderer.
func (r *defaultTextRenderer) RenderTextRun(run TextRun, pass render.RenderPass) error {
	if run == nil || !run.IsValid() {
		return errors.New("invalid text run")
	}

	// TODO: Implement actual rendering
	// This would involve:
	// 1. Getting glyph atlas entries for all glyphs
	// 2. Creating vertex buffers with glyph positions and UV coordinates
	// 3. Binding the atlas texture
	// 4. Drawing the glyphs

	return nil
}

// PrepareGlyphs implements TextRenderer.
func (r *defaultTextRenderer) PrepareGlyphs(frame TextFrame) error {
	// TODO: Ensure all glyphs are loaded in atlases
	return nil
}

// GetShader implements TextRenderer.
func (r *defaultTextRenderer) GetShader() render.Shader {
	return r.shader
}
