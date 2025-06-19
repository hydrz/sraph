// Package entity provides texture entity implementation.
package entity

import (
	"github.com/opensraph/sraph/geom"
)

// TextureEntity represents an entity that renders a texture.
type TextureEntity struct {
	*BaseEntity
	textureID TextureID
	uvRect    geom.Rect[geom.F32] // UV coordinates for texture sampling
	tintColor geom.Color          // Color to tint the texture
}

// TextureID represents a unique identifier for a texture resource.
type TextureID uint32

// NewTextureEntity creates a new texture entity.
func NewTextureEntity(bounds geom.Rect[geom.F32], textureID TextureID) *TextureEntity {
	entity := &TextureEntity{
		BaseEntity: NewBaseEntity(),
		textureID:  textureID,
		uvRect:     geom.NewRect[geom.F32](0, 0, 1, 1), // Full texture by default
		tintColor:  geom.ColorWhite(),                  // No tint by default
	}
	entity.SetBounds(bounds)
	return entity
}

// TextureID returns the texture ID for this entity.
func (t *TextureEntity) TextureID() TextureID {
	return t.textureID
}

// SetTextureID sets the texture ID for this entity.
func (t *TextureEntity) SetTextureID(textureID TextureID) {
	t.textureID = textureID
}

// UVRect returns the UV rectangle for texture sampling.
func (t *TextureEntity) UVRect() geom.Rect[geom.F32] {
	return t.uvRect
}

// SetUVRect sets the UV rectangle for texture sampling.
func (t *TextureEntity) SetUVRect(uvRect geom.Rect[geom.F32]) {
	t.uvRect = uvRect
}

// TintColor returns the tint color applied to the texture.
func (t *TextureEntity) TintColor() geom.Color {
	return t.tintColor
}

// SetTintColor sets the tint color applied to the texture.
func (t *TextureEntity) SetTintColor(color geom.Color) {
	t.tintColor = color
}

// Render implements RenderableEntity.
func (t *TextureEntity) Render(ctx RenderContext) error {
	// TODO: Implement texture rendering
	// This would typically involve:
	// 1. Binding the texture resource
	// 2. Setting up shader uniforms with UV mapping and tint color
	// 3. Drawing a quad with texture coordinates
	// 4. Applying the entity's transformation matrix
	return nil
}

// SpriteEntity represents a sprite (sub-region of a texture atlas).
type SpriteEntity struct {
	*TextureEntity
	atlasID    TextureID
	spriteRect geom.Rect[geom.F32] // Rectangle in atlas coordinates
}

// NewSpriteEntity creates a new sprite entity from a texture atlas.
func NewSpriteEntity(bounds geom.Rect[geom.F32], atlasID TextureID, spriteRect geom.Rect[geom.F32]) *SpriteEntity {
	entity := &SpriteEntity{
		TextureEntity: NewTextureEntity(bounds, atlasID),
		atlasID:       atlasID,
		spriteRect:    spriteRect,
	}

	// Set UV coordinates based on sprite rectangle in atlas
	// Assuming atlas coordinates are normalized (0-1)
	entity.SetUVRect(spriteRect)

	return entity
}

// AtlasID returns the texture atlas ID.
func (s *SpriteEntity) AtlasID() TextureID {
	return s.atlasID
}

// SetAtlasID sets the texture atlas ID.
func (s *SpriteEntity) SetAtlasID(atlasID TextureID) {
	s.atlasID = atlasID
	s.SetTextureID(atlasID)
}

// SpriteRect returns the sprite rectangle in atlas coordinates.
func (s *SpriteEntity) SpriteRect() geom.Rect[geom.F32] {
	return s.spriteRect
}

// SetSpriteRect sets the sprite rectangle in atlas coordinates.
func (s *SpriteEntity) SetSpriteRect(spriteRect geom.Rect[geom.F32]) {
	s.spriteRect = spriteRect
	s.SetUVRect(spriteRect)
}

// NinePatchEntity represents a nine-patch texture entity for scalable UI elements.
type NinePatchEntity struct {
	*TextureEntity
	borderLeft   geom.F32
	borderTop    geom.F32
	borderRight  geom.F32
	borderBottom geom.F32
}

// NewNinePatchEntity creates a new nine-patch entity.
func NewNinePatchEntity(
	bounds geom.Rect[geom.F32],
	textureID TextureID,
	borderLeft, borderTop, borderRight, borderBottom geom.F32,
) *NinePatchEntity {
	return &NinePatchEntity{
		TextureEntity: NewTextureEntity(bounds, textureID),
		borderLeft:    borderLeft,
		borderTop:     borderTop,
		borderRight:   borderRight,
		borderBottom:  borderBottom,
	}
}

// BorderLeft returns the left border size.
func (n *NinePatchEntity) BorderLeft() geom.F32 {
	return n.borderLeft
}

// BorderTop returns the top border size.
func (n *NinePatchEntity) BorderTop() geom.F32 {
	return n.borderTop
}

// BorderRight returns the right border size.
func (n *NinePatchEntity) BorderRight() geom.F32 {
	return n.borderRight
}

// BorderBottom returns the bottom border size.
func (n *NinePatchEntity) BorderBottom() geom.F32 {
	return n.borderBottom
}

// SetBorders sets all border sizes.
func (n *NinePatchEntity) SetBorders(left, top, right, bottom geom.F32) {
	n.borderLeft = left
	n.borderTop = top
	n.borderRight = right
	n.borderBottom = bottom
}

// Render implements RenderableEntity.
func (n *NinePatchEntity) Render(ctx RenderContext) error {
	// TODO: Implement nine-patch rendering
	// This involves dividing the texture into 9 regions and scaling appropriately
	return nil
}

// TiledTextureEntity represents a texture that is tiled across the entity bounds.
type TiledTextureEntity struct {
	*TextureEntity
	tileSize geom.Size[geom.F32]
	offset   geom.Point[geom.F32]
}

// NewTiledTextureEntity creates a new tiled texture entity.
func NewTiledTextureEntity(bounds geom.Rect[geom.F32], textureID TextureID, tileSize geom.Size[geom.F32]) *TiledTextureEntity {
	return &TiledTextureEntity{
		TextureEntity: NewTextureEntity(bounds, textureID),
		tileSize:      tileSize,
		offset:        geom.Point[geom.F32]{X: 0, Y: 0},
	}
}

// TileSize returns the size of each tile.
func (t *TiledTextureEntity) TileSize() geom.Size[geom.F32] {
	return t.tileSize
}

// SetTileSize sets the size of each tile.
func (t *TiledTextureEntity) SetTileSize(size geom.Size[geom.F32]) {
	t.tileSize = size
}

// Offset returns the tiling offset.
func (t *TiledTextureEntity) Offset() geom.Point[geom.F32] {
	return t.offset
}

// SetOffset sets the tiling offset.
func (t *TiledTextureEntity) SetOffset(offset geom.Point[geom.F32]) {
	t.offset = offset
}

// Render implements RenderableEntity.
func (t *TiledTextureEntity) Render(ctx RenderContext) error {
	// TODO: Implement tiled texture rendering
	return nil
}
