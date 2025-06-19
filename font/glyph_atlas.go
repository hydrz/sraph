// Package font - Glyph atlas functionality
package font

import (
	"errors"
	"sync"

	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// GlyphAtlasType represents the type of glyph atlas.
type GlyphAtlasType uint8

const (
	// GlyphAtlasTypeAlpha represents an alpha-only atlas for regular text.
	GlyphAtlasTypeAlpha GlyphAtlasType = iota
	// GlyphAtlasTypeColor represents a color atlas for emoji and color fonts.
	GlyphAtlasTypeColor
)

// GlyphAtlas manages a texture atlas for rendered glyphs.
type GlyphAtlas interface {
	// GetType returns the type of this atlas.
	GetType() GlyphAtlasType

	// GetTexture returns the underlying texture.
	GetTexture() render.Texture

	// GetSize returns the size of the atlas in pixels.
	GetSize() geom.Size[geom.I32]

	// FindGlyph looks up a glyph in the atlas.
	FindGlyph(pair FontGlyphPair) (GlyphAtlasEntry, bool)

	// AddGlyph adds a glyph to the atlas and returns its entry.
	AddGlyph(pair FontGlyphPair, bitmap GlyphBitmap) (GlyphAtlasEntry, error)

	// GetPercentFull returns how full the atlas is (0.0 to 1.0).
	GetPercentFull() float32

	// Reset clears the atlas.
	Reset()
}

// GlyphAtlasEntry represents a glyph's location in an atlas.
type GlyphAtlasEntry struct {
	// Bounds is the rectangle containing the glyph in the atlas.
	Bounds geom.Rect[geom.I32]
	// UVBounds is the normalized texture coordinates (0.0 to 1.0).
	UVBounds geom.Rect[geom.F32]
	// IsValid indicates if this entry is valid.
	IsValid bool
}

// GlyphBitmap represents a rendered glyph as a bitmap.
type GlyphBitmap struct {
	// Width is the width of the bitmap in pixels.
	Width int
	// Height is the height of the bitmap in pixels.
	Height int
	// Data is the bitmap data (RGBA for color, A for alpha).
	Data []byte
	// Format is the pixel format.
	Format GlyphBitmapFormat
	// Bounds is the logical bounds of the glyph.
	Bounds geom.Rect[geom.F32]
}

// GlyphBitmapFormat represents the pixel format of a glyph bitmap.
type GlyphBitmapFormat uint8

const (
	// GlyphBitmapFormatAlpha represents 8-bit alpha channel only.
	GlyphBitmapFormatAlpha GlyphBitmapFormat = iota
	// GlyphBitmapFormatRGBA represents 32-bit RGBA.
	GlyphBitmapFormatRGBA
)

// RectanglePacker handles packing rectangles into an atlas.
type RectanglePacker interface {
	// AddRect attempts to add a rectangle and returns its position.
	AddRect(width, height int) (geom.Point[geom.I32], bool)
	// GetPercentFull returns how full the packer is (0.0 to 1.0).
	GetPercentFull() float32
	// Reset clears the packer.
	Reset()
}

// GlyphAtlasContext manages glyph atlas creation and caching.
type GlyphAtlasContext interface {
	// GetGlyphAtlas returns the current atlas.
	GetGlyphAtlas() GlyphAtlas

	// GetAtlasSize returns the size of the atlas.
	GetAtlasSize() geom.Size[int]

	// GetRectPacker returns the rectangle packer.
	GetRectPacker() RectanglePacker

	// UpdateGlyphAtlas updates the atlas with a new one.
	UpdateGlyphAtlas(atlas GlyphAtlas, size geom.Size[int])

	// UpdateRectPacker updates the rectangle packer.
	UpdateRectPacker(packer RectanglePacker)
}

// defaultGlyphAtlas provides the default implementation of GlyphAtlas.
type defaultGlyphAtlas struct {
	atlasType GlyphAtlasType
	texture   render.Texture
	size      geom.Size[geom.I32]
	entries   map[uint64]GlyphAtlasEntry
	packer    RectanglePacker
	mutex     sync.RWMutex
}

// NewGlyphAtlas creates a new glyph atlas.
func NewGlyphAtlas(atlasType GlyphAtlasType, size geom.Size[geom.I32], context render.Context) (GlyphAtlas, error) {
	if context == nil {
		return nil, errors.New("render context is required")
	}

	// TODO: Create texture - this requires implementing texture creation in render package
	// For now, return nil texture
	var texture render.Texture

	return &defaultGlyphAtlas{
		atlasType: atlasType,
		texture:   texture,
		size:      size,
		entries:   make(map[uint64]GlyphAtlasEntry),
		packer:    NewSkylineRectanglePacker(int(size.Width), int(size.Height)),
	}, nil
}

// GetType implements GlyphAtlas.
func (a *defaultGlyphAtlas) GetType() GlyphAtlasType {
	return a.atlasType
}

// GetTexture implements GlyphAtlas.
func (a *defaultGlyphAtlas) GetTexture() render.Texture {
	return a.texture
}

// GetSize implements GlyphAtlas.
func (a *defaultGlyphAtlas) GetSize() geom.Size[geom.I32] {
	return a.size
}

// FindGlyph implements GlyphAtlas.
func (a *defaultGlyphAtlas) FindGlyph(pair FontGlyphPair) (GlyphAtlasEntry, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	entry, exists := a.entries[pair.GetHash()]
	return entry, exists
}

// AddGlyph implements GlyphAtlas.
func (a *defaultGlyphAtlas) AddGlyph(pair FontGlyphPair, bitmap GlyphBitmap) (GlyphAtlasEntry, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Check if already exists
	if entry, exists := a.entries[pair.GetHash()]; exists {
		return entry, nil
	}

	// Pack the glyph
	position, ok := a.packer.AddRect(bitmap.Width, bitmap.Height)
	if !ok {
		return GlyphAtlasEntry{}, errors.New("failed to pack glyph into atlas")
	}

	// Create bounds
	bounds := geom.Rect[int]{
		Left:   position.X,
		Top:    position.Y,
		Right:  position.X + bitmap.Width,
		Bottom: position.Y + bitmap.Height,
	}

	// Create UV bounds
	uvBounds := geom.Rect[geom.F32]{
		Left:   geom.F32(bounds.Left) / geom.F32(a.size.Width),
		Top:    geom.F32(bounds.Top) / geom.F32(a.size.Height),
		Right:  geom.F32(bounds.Right) / geom.F32(a.size.Width),
		Bottom: geom.F32(bounds.Bottom) / geom.F32(a.size.Height),
	}

	entry := GlyphAtlasEntry{
		Bounds:   bounds,
		UVBounds: uvBounds,
		IsValid:  true,
	}

	a.entries[pair.GetHash()] = entry

	// TODO: Upload bitmap data to texture
	// This would require updating the texture with the bitmap data at the specified position

	return entry, nil
}

// GetPercentFull implements GlyphAtlas.
func (a *defaultGlyphAtlas) GetPercentFull() float32 {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.packer.GetPercentFull()
}

// Reset implements GlyphAtlas.
func (a *defaultGlyphAtlas) Reset() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.entries = make(map[uint64]GlyphAtlasEntry)
	a.packer.Reset()
}

// defaultGlyphAtlasContext provides the default implementation of GlyphAtlasContext.
type defaultGlyphAtlasContext struct {
	atlas  GlyphAtlas
	size   geom.Size[int]
	packer RectanglePacker
	mutex  sync.RWMutex
}

// NewGlyphAtlasContext creates a new glyph atlas context.
func NewGlyphAtlasContext(atlasType GlyphAtlasType, context render.Context) (GlyphAtlasContext, error) {
	initialSize := geom.Size[int]{Width: 512, Height: 512}

	atlas, err := NewGlyphAtlas(atlasType, initialSize, context)
	if err != nil {
		return nil, err
	}

	return &defaultGlyphAtlasContext{
		atlas:  atlas,
		size:   initialSize,
		packer: NewSkylineRectanglePacker(initialSize.Width, initialSize.Height),
	}, nil
}

// GetGlyphAtlas implements GlyphAtlasContext.
func (c *defaultGlyphAtlasContext) GetGlyphAtlas() GlyphAtlas {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.atlas
}

// GetAtlasSize implements GlyphAtlasContext.
func (c *defaultGlyphAtlasContext) GetAtlasSize() geom.Size[int] {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.size
}

// GetRectPacker implements GlyphAtlasContext.
func (c *defaultGlyphAtlasContext) GetRectPacker() RectanglePacker {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.packer
}

// UpdateGlyphAtlas implements GlyphAtlasContext.
func (c *defaultGlyphAtlasContext) UpdateGlyphAtlas(atlas GlyphAtlas, size geom.Size[int]) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.atlas = atlas
	c.size = size
}

// UpdateRectPacker implements GlyphAtlasContext.
func (c *defaultGlyphAtlasContext) UpdateRectPacker(packer RectanglePacker) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.packer = packer
}

// Helper functions

func getTextureFormat(atlasType GlyphAtlasType) render.TextureFormat {
	switch atlasType {
	case GlyphAtlasTypeAlpha:
		return render.TextureFormatR8Unorm
	case GlyphAtlasTypeColor:
		return render.TextureFormatRGBA8Unorm
	default:
		return render.TextureFormatR8Unorm
	}
}
