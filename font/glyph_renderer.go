// Package font - Glyph renderer functionality
package font

import (
	"errors"
	"fmt"

	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// GlyphRenderer handles rendering of glyphs to bitmaps and managing glyph atlases.
// It bridges between the font system and the rendering pipeline.
type GlyphRenderer interface {
	// RenderGlyph renders a glyph to a bitmap.
	RenderGlyph(pair FontGlyphPair, size float32) (GlyphBitmap, error)

	// GetGlyphAtlas returns the glyph atlas for the given type.
	GetGlyphAtlas(atlasType GlyphAtlasType) GlyphAtlas

	// EnsureGlyphInAtlas ensures a glyph is present in the atlas.
	EnsureGlyphInAtlas(pair FontGlyphPair) (GlyphAtlasEntry, error)

	// CreateGlyphAtlas creates a new glyph atlas.
	CreateGlyphAtlas(atlasType GlyphAtlasType, size geom.Size[geom.I32]) (GlyphAtlas, error)

	// GetContext returns the rendering context.
	GetContext() render.Context
}

// GlyphRasterizer handles converting glyph paths to bitmaps.
type GlyphRasterizer interface {
	// RasterizeGlyph rasterizes a glyph path to a bitmap.
	RasterizeGlyph(glyph Glyph, font Font, size float32) (GlyphBitmap, error)

	// GetGlyphPath returns the path for a glyph.
	GetGlyphPath(glyph Glyph, font Font) (geom.Path, error)

	// GetGlyphOutline returns the outline for a glyph.
	GetGlyphOutline(glyph Glyph, font Font) ([]geom.Point[geom.F32], error)
}

// GlyphCache provides caching for rendered glyphs.
type GlyphCache interface {
	// Get retrieves a glyph bitmap from the cache.
	Get(pair FontGlyphPair) (GlyphBitmap, bool)

	// Put stores a glyph bitmap in the cache.
	Put(pair FontGlyphPair, bitmap GlyphBitmap)

	// Remove removes a glyph from the cache.
	Remove(pair FontGlyphPair)

	// Clear clears the cache.
	Clear()

	// Size returns the number of cached glyphs.
	Size() int
}

// defaultGlyphRenderer provides the default implementation of GlyphRenderer.
type defaultGlyphRenderer struct {
	context    render.Context
	alphaAtlas GlyphAtlas
	colorAtlas GlyphAtlas
	rasterizer GlyphRasterizer
	cache      GlyphCache
}

// NewGlyphRenderer creates a new glyph renderer.
func NewGlyphRenderer(context render.Context) (GlyphRenderer, error) {
	if context == nil {
		return nil, errors.New("render context is required")
	}

	// Create initial atlases
	atlasSize := geom.Size[geom.I32]{Width: 512, Height: 512}

	alphaAtlas, err := NewGlyphAtlas(GlyphAtlasTypeAlpha, atlasSize, context)
	if err != nil {
		return nil, fmt.Errorf("failed to create alpha atlas: %w", err)
	}

	colorAtlas, err := NewGlyphAtlas(GlyphAtlasTypeColor, atlasSize, context)
	if err != nil {
		return nil, fmt.Errorf("failed to create color atlas: %w", err)
	}

	return &defaultGlyphRenderer{
		context:    context,
		alphaAtlas: alphaAtlas,
		colorAtlas: colorAtlas,
		rasterizer: NewGlyphRasterizer(),
		cache:      NewGlyphCache(),
	}, nil
}

// RenderGlyph implements GlyphRenderer.
func (r *defaultGlyphRenderer) RenderGlyph(pair FontGlyphPair, size float32) (GlyphBitmap, error) {
	// Check cache first
	if bitmap, exists := r.cache.Get(pair); exists {
		return bitmap, nil
	}

	// Rasterize the glyph
	bitmap, err := r.rasterizer.RasterizeGlyph(pair.Glyph, pair.ScaledFont.Font, size)
	if err != nil {
		return GlyphBitmap{}, fmt.Errorf("failed to rasterize glyph: %w", err)
	}

	// Cache the result
	r.cache.Put(pair, bitmap)

	return bitmap, nil
}

// GetGlyphAtlas implements GlyphRenderer.
func (r *defaultGlyphRenderer) GetGlyphAtlas(atlasType GlyphAtlasType) GlyphAtlas {
	switch atlasType {
	case GlyphAtlasTypeAlpha:
		return r.alphaAtlas
	case GlyphAtlasTypeColor:
		return r.colorAtlas
	default:
		return r.alphaAtlas
	}
}

// EnsureGlyphInAtlas implements GlyphRenderer.
func (r *defaultGlyphRenderer) EnsureGlyphInAtlas(pair FontGlyphPair) (GlyphAtlasEntry, error) {
	// Determine atlas type
	atlasType := GlyphAtlasTypeAlpha // Default for now
	atlas := r.GetGlyphAtlas(atlasType)

	// Check if already in atlas
	if entry, exists := atlas.FindGlyph(pair); exists {
		return entry, nil
	}

	// Render the glyph
	bitmap, err := r.RenderGlyph(pair, 32) // TODO: Use appropriate size
	if err != nil {
		return GlyphAtlasEntry{}, err
	}

	// Add to atlas
	entry, err := atlas.AddGlyph(pair, bitmap)
	if err != nil {
		return GlyphAtlasEntry{}, fmt.Errorf("failed to add glyph to atlas: %w", err)
	}

	return entry, nil
}

// CreateGlyphAtlas implements GlyphRenderer.
func (r *defaultGlyphRenderer) CreateGlyphAtlas(atlasType GlyphAtlasType, size geom.Size[geom.I32]) (GlyphAtlas, error) {
	return NewGlyphAtlas(atlasType, size, r.context)
}

// GetContext implements GlyphRenderer.
func (r *defaultGlyphRenderer) GetContext() render.Context {
	return r.context
}

// defaultGlyphRasterizer provides the default implementation of GlyphRasterizer.
type defaultGlyphRasterizer struct{}

// NewGlyphRasterizer creates a new glyph rasterizer.
func NewGlyphRasterizer() GlyphRasterizer {
	return &defaultGlyphRasterizer{}
}

// RasterizeGlyph implements GlyphRasterizer.
func (r *defaultGlyphRasterizer) RasterizeGlyph(glyph Glyph, font Font, size float32) (GlyphBitmap, error) {
	if font == nil || !font.IsValid() {
		return GlyphBitmap{}, errors.New("invalid font")
	}

	typeface := font.GetTypeface()
	if typeface == nil {
		return GlyphBitmap{}, errors.New("font has no typeface")
	}

	metrics := typeface.GetGlyphMetrics(glyph.Index)

	// Calculate bitmap size
	width := int(metrics.AdvanceWidth + 2) // Add padding
	height := int(font.GetMetrics().LineHeight + 2)

	if width <= 0 || height <= 0 {
		return GlyphBitmap{}, errors.New("invalid glyph dimensions")
	}

	// Create bitmap data
	// TODO: Implement actual rasterization
	// For now, create a simple placeholder bitmap
	data := make([]byte, width*height)
	for i := range data {
		data[i] = 128 // Gray placeholder
	}

	return GlyphBitmap{
		Width:  width,
		Height: height,
		Data:   data,
		Format: GlyphBitmapFormatAlpha,
		Bounds: metrics.BoundingBox,
	}, nil
}

// GetGlyphPath implements GlyphRasterizer.
func (r *defaultGlyphRasterizer) GetGlyphPath(glyph Glyph, font Font) (geom.Path, error) {
	// TODO: Implement path extraction from font data
	return geom.Path{}, errors.New("not implemented")
}

// GetGlyphOutline implements GlyphRasterizer.
func (r *defaultGlyphRasterizer) GetGlyphOutline(glyph Glyph, font Font) ([]geom.Point[geom.F32], error) {
	// TODO: Implement outline extraction from font data
	return nil, errors.New("not implemented")
}

// defaultGlyphCache provides the default implementation of GlyphCache.
type defaultGlyphCache struct {
	glyphs map[uint64]GlyphBitmap
}

// NewGlyphCache creates a new glyph cache.
func NewGlyphCache() GlyphCache {
	return &defaultGlyphCache{
		glyphs: make(map[uint64]GlyphBitmap),
	}
}

// Get implements GlyphCache.
func (c *defaultGlyphCache) Get(pair FontGlyphPair) (GlyphBitmap, bool) {
	bitmap, exists := c.glyphs[pair.GetHash()]
	return bitmap, exists
}

// Put implements GlyphCache.
func (c *defaultGlyphCache) Put(pair FontGlyphPair, bitmap GlyphBitmap) {
	c.glyphs[pair.GetHash()] = bitmap
}

// Remove implements GlyphCache.
func (c *defaultGlyphCache) Remove(pair FontGlyphPair) {
	delete(c.glyphs, pair.GetHash())
}

// Clear implements GlyphCache.
func (c *defaultGlyphCache) Clear() {
	c.glyphs = make(map[uint64]GlyphBitmap)
}

// Size implements GlyphCache.
func (c *defaultGlyphCache) Size() int {
	return len(c.glyphs)
}
