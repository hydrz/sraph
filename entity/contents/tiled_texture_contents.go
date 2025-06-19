package contents

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// TiledTextureContents provides tiled texture rendering capabilities
type TiledTextureContents struct {
	BaseContents
	texture     render.Texture
	srcRect     geom.Rect
	dstRect     geom.Rect
	tileMode    TileMode
	filterMode  FilterMode
	opacity     float32
	transform   geom.Matrix
	colorFilter *ColorFilter
	blendMode   render.BlendMode
}

// TileMode defines how textures are tiled
type TileMode int

const (
	// TileModeClamp clamp texture coordinates to edge
	TileModeClamp TileMode = iota
	// TileModeRepeat repeat texture
	TileModeRepeat
	// TileModeMirror mirror texture on repeat
	TileModeMirror
	// TileModeDecal decal mode (transparent outside bounds)
	TileModeDecal
)

// FilterMode defines texture filtering
type FilterMode int

const (
	// FilterModeNearest nearest neighbor filtering
	FilterModeNearest FilterMode = iota
	// FilterModeLinear linear filtering
	FilterModeLinear
	// FilterModeCubic cubic filtering
	FilterModeCubic
)

// ColorFilter defines color transformation for textures
type ColorFilter struct {
	Matrix [4][5]float32 // 4x5 color matrix for RGBA + offset
	Mode   ColorFilterMode
}

// ColorFilterMode defines color filter modes
type ColorFilterMode int

const (
	// ColorFilterModeMatrix apply color matrix
	ColorFilterModeMatrix ColorFilterMode = iota
	// ColorFilterModeMultiply multiply colors
	ColorFilterModeMultiply
	// ColorFilterModeScreen screen blend colors
	ColorFilterModeScreen
)

// NewTiledTextureContents creates a new tiled texture contents
func NewTiledTextureContents(texture render.Texture, srcRect, dstRect geom.Rect) *TiledTextureContents {
	return &TiledTextureContents{
		BaseContents: NewBaseContents(),
		texture:      texture,
		srcRect:      srcRect,
		dstRect:      dstRect,
		tileMode:     TileModeRepeat,
		filterMode:   FilterModeLinear,
		opacity:      1.0,
		transform:    geom.NewIdentityMatrix(),
		colorFilter:  nil,
		blendMode:    render.BlendModeSourceOver,
	}
}

// SetTexture sets the texture to be tiled
func (t *TiledTextureContents) SetTexture(texture render.Texture) {
	t.texture = texture
}

// GetTexture returns the current texture
func (t *TiledTextureContents) GetTexture() render.Texture {
	return t.texture
}

// SetSrcRect sets the source rectangle within the texture
func (t *TiledTextureContents) SetSrcRect(rect geom.Rect) {
	t.srcRect = rect
}

// GetSrcRect returns the current source rectangle
func (t *TiledTextureContents) GetSrcRect() geom.Rect {
	return t.srcRect
}

// SetDstRect sets the destination rectangle for tiling
func (t *TiledTextureContents) SetDstRect(rect geom.Rect) {
	t.dstRect = rect
}

// GetDstRect returns the current destination rectangle
func (t *TiledTextureContents) GetDstRect() geom.Rect {
	return t.dstRect
}

// SetTileMode sets the tiling mode
func (t *TiledTextureContents) SetTileMode(mode TileMode) {
	t.tileMode = mode
}

// GetTileMode returns the current tile mode
func (t *TiledTextureContents) GetTileMode() TileMode {
	return t.tileMode
}

// SetFilterMode sets the texture filtering mode
func (t *TiledTextureContents) SetFilterMode(mode FilterMode) {
	t.filterMode = mode
}

// GetFilterMode returns the current filter mode
func (t *TiledTextureContents) GetFilterMode() FilterMode {
	return t.filterMode
}

// SetOpacity sets the opacity for the tiled texture
func (t *TiledTextureContents) SetOpacity(opacity float32) {
	t.opacity = opacity
}

// GetOpacity returns the current opacity
func (t *TiledTextureContents) GetOpacity() float32 {
	return t.opacity
}

// SetTransform sets the transformation matrix
func (t *TiledTextureContents) SetTransform(transform geom.Matrix) {
	t.transform = transform
}

// GetTransform returns the current transformation matrix
func (t *TiledTextureContents) GetTransform() geom.Matrix {
	return t.transform
}

// SetColorFilter sets the color filter
func (t *TiledTextureContents) SetColorFilter(filter *ColorFilter) {
	t.colorFilter = filter
}

// GetColorFilter returns the current color filter
func (t *TiledTextureContents) GetColorFilter() *ColorFilter {
	return t.colorFilter
}

// SetBlendMode sets the blend mode
func (t *TiledTextureContents) SetBlendMode(mode render.BlendMode) {
	t.blendMode = mode
}

// GetBlendMode returns the current blend mode
func (t *TiledTextureContents) GetBlendMode() render.BlendMode {
	return t.blendMode
}

// Render implements the Contents interface
func (t *TiledTextureContents) Render(context *ContentContext, entity *Entity, pass *RenderPass) bool {
	// TODO: Implement tiled texture rendering
	// This would involve:
	// 1. Setting up texture sampling with specified tile mode
	// 2. Applying filtering mode
	// 3. Computing tiling parameters based on src/dst rects
	// 4. Applying color filter if specified
	// 5. Rendering with specified opacity and blend mode
	// 6. Applying transformation matrix
	return false
}

// GetBounds returns the bounds of this contents
func (t *TiledTextureContents) GetBounds() geom.Rect {
	// Transform the destination rectangle
	transformedRect := t.transform.TransformRect(t.dstRect)
	return transformedRect
}

// Clone creates a copy of this contents
func (t *TiledTextureContents) Clone() Contents {
	clone := &TiledTextureContents{
		BaseContents: t.BaseContents.Clone().(BaseContents),
		texture:      t.texture,
		srcRect:      t.srcRect,
		dstRect:      t.dstRect,
		tileMode:     t.tileMode,
		filterMode:   t.filterMode,
		opacity:      t.opacity,
		transform:    t.transform,
		blendMode:    t.blendMode,
	}

	// Deep copy color filter if it exists
	if t.colorFilter != nil {
		clone.colorFilter = &ColorFilter{
			Matrix: t.colorFilter.Matrix,
			Mode:   t.colorFilter.Mode,
		}
	}

	return clone
}

// GetCoverage returns the coverage area for this contents
func (t *TiledTextureContents) GetCoverage(transform geom.Matrix) geom.Rect {
	combinedTransform := transform.Multiply(t.transform)
	return combinedTransform.TransformRect(t.dstRect)
}

// TiledTextureFactory creates tiled texture contents instances
type TiledTextureFactory struct{}

// CreateTiledTexture creates a new tiled texture contents
func (f *TiledTextureFactory) CreateTiledTexture(texture render.Texture, srcRect, dstRect geom.Rect) Contents {
	return NewTiledTextureContents(texture, srcRect, dstRect)
}

// CreateTiledTextureWithMode creates a new tiled texture with specified tile mode
func (f *TiledTextureFactory) CreateTiledTextureWithMode(texture render.Texture, srcRect, dstRect geom.Rect, tileMode TileMode) Contents {
	contents := NewTiledTextureContents(texture, srcRect, dstRect)
	contents.SetTileMode(tileMode)
	return contents
}

// CreateTiledTextureWithFilter creates a new tiled texture with specified filter mode
func (f *TiledTextureFactory) CreateTiledTextureWithFilter(texture render.Texture, srcRect, dstRect geom.Rect, tileMode TileMode, filterMode FilterMode) Contents {
	contents := NewTiledTextureContents(texture, srcRect, dstRect)
	contents.SetTileMode(tileMode)
	contents.SetFilterMode(filterMode)
	return contents
}

// NewColorFilter creates a new color filter with the specified matrix
func NewColorFilter(matrix [4][5]float32, mode ColorFilterMode) *ColorFilter {
	return &ColorFilter{
		Matrix: matrix,
		Mode:   mode,
	}
}

// NewIdentityColorFilter creates an identity color filter (no change)
func NewIdentityColorFilter() *ColorFilter {
	return &ColorFilter{
		Matrix: [4][5]float32{
			{1, 0, 0, 0, 0}, // Red
			{0, 1, 0, 0, 0}, // Green
			{0, 0, 1, 0, 0}, // Blue
			{0, 0, 0, 1, 0}, // Alpha
		},
		Mode: ColorFilterModeMatrix,
	}
}
