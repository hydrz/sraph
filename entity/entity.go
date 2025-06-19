// Package entity provides renderable entities and content for the Sraph UI toolkit.
//
// This package defines the entity system inspired by Flutter's Impeller, where
// entities combine graphical content with properties like transformation,
// blend mode, and stencil clip depth.
package entity

import (
	"github.com/opensraph/sraph/geom"
)

// BlendMode defines how pixels are combined during rendering.
// Based on Impeller's BlendMode enum and WebGPU blend operations.
type BlendMode uint8

const (
	// Basic blend modes supported by most graphics pipelines
	BlendModeClear BlendMode = iota
	BlendModeSrc
	BlendModeDst
	BlendModeSrcOver // Default
	BlendModeDstOver
	BlendModeSrcIn
	BlendModeDstIn
	BlendModeSrcOut
	BlendModeDstOut
	BlendModeSrcATop
	BlendModeDstATop
	BlendModeXor
	BlendModePlus
	BlendModeModulate

	// Advanced blend modes (may require special handling)
	BlendModeScreen
	BlendModeOverlay
	BlendModeDarken
	BlendModeLighten
	BlendModeColorDodge
	BlendModeColorBurn
	BlendModeHardLight
	BlendModeSoftLight
	BlendModeDifference
	BlendModeExclusion
	BlendModeMultiply
	BlendModeHue
	BlendModeSaturation
	BlendModeColor
	BlendModeLuminosity
)

// TileMode defines how to handle colors outside the source bounds.
// Based on Impeller's Entity::TileMode.
type TileMode uint8

const (
	TileModeClamp  TileMode = iota // Replicate edge color
	TileModeRepeat                 // Repeat the pattern
	TileModeMirror                 // Mirror the pattern
	TileModeDecal                  // Transparent outside bounds
)

// ClipOperation defines how clipping regions are combined.
type ClipOperation uint8

const (
	ClipOperationDifference ClipOperation = iota
	ClipOperationIntersect
)

// RenderingMode defines how entities are transformed during rendering.
// Based on Impeller's Entity::RenderingMode.
type RenderingMode uint8

const (
	// Direct mode: Entity's transform is used as local-to-screen transform
	RenderingModeDirect RenderingMode = iota

	// Subpass modes: Entity is in screen space, transform handling varies
	RenderingModeSubpassAppendSnapshot
	RenderingModeSubpassPrependSnapshot
)

// Contents represents the graphical content that can be rendered.
// This is inspired by Impeller's Contents interface.
type Contents interface {
	// Render performs the actual rendering of this content.
	Render(ctx *ContentContext, entity *Entity, pass RenderPass) bool

	// GetCoverage returns the coverage area of this content for the given entity.
	GetCoverage(entity *Entity) *geom.Rect[geom.F32]

	// CanInheritOpacity returns whether this content can inherit opacity from parent.
	CanInheritOpacity(entity *Entity) bool

	// Clone creates a copy of this content.
	Clone() Contents
}

// Entity represents a renderable object within the Impeller scene.
// An Entity combines graphical content (Contents) with properties like
// transformation, blend mode, and stencil clip depth.
type Entity struct {
	contents  Contents
	transform geom.Matrix[geom.F32]
	blendMode BlendMode
	clipDepth geom.F32
	visible   bool
	zIndex    int
}

// NewEntity creates a new Entity with the given contents.
func NewEntity() *Entity {
	return &Entity{
		transform: geom.NewMatrix[geom.F32](), // Identity matrix
		blendMode: BlendModeSrcOver,
		clipDepth: 0.0,
		visible:   true,
		zIndex:    0,
	}
}

// GetTransform returns the transformation matrix for this entity.
func (e *Entity) GetTransform() geom.Matrix[geom.F32] {
	return e.transform
}

// SetTransform sets the transformation matrix for this entity.
func (e *Entity) SetTransform(transform geom.Matrix[geom.F32]) {
	e.transform = transform
}

// GetBlendMode returns the blend mode for this entity.
func (e *Entity) GetBlendMode() BlendMode {
	return e.blendMode
}

// SetBlendMode sets the blend mode for this entity.
func (e *Entity) SetBlendMode(blendMode BlendMode) {
	e.blendMode = blendMode
}

// GetClipDepth returns the stencil clip depth for this entity.
func (e *Entity) GetClipDepth() geom.F32 {
	return e.clipDepth
}

// SetClipDepth sets the stencil clip depth for this entity.
func (e *Entity) SetClipDepth(clipDepth geom.F32) {
	e.clipDepth = clipDepth
}

// IsVisible returns whether this entity is visible.
func (e *Entity) IsVisible() bool {
	return e.visible
}

// SetVisible sets the visibility of this entity.
func (e *Entity) SetVisible(visible bool) {
	e.visible = visible
}

// GetZIndex returns the z-index (depth) of this entity.
func (e *Entity) GetZIndex() int {
	return e.zIndex
}

// SetZIndex sets the z-index (depth) of this entity.
func (e *Entity) SetZIndex(zIndex int) {
	e.zIndex = zIndex
}

// GetContents returns the contents of this entity.
func (e *Entity) GetContents() Contents {
	return e.contents
}

// SetContents sets the contents of this entity.
func (e *Entity) SetContents(contents Contents) {
	e.contents = contents
}

// Render renders this entity using the given context and render pass.
func (e *Entity) Render(ctx *ContentContext, pass RenderPass) bool {
	if !e.visible || e.contents == nil {
		return true
	}
	return e.contents.Render(ctx, e, pass)
}

// GetCoverage returns the coverage area of this entity.
func (e *Entity) GetCoverage() *geom.Rect[geom.F32] {
	if e.contents == nil {
		return nil
	}
	return e.contents.GetCoverage(e)
}

// ContentContext provides rendering context for Contents.
// This mirrors Impeller's ContentContext.
type ContentContext struct {
	// Add context fields as needed for rendering
	// This will be expanded based on GPU backend requirements
}

// RenderPass represents a rendering pass.
// This will be implemented by the GPU backend.
type RenderPass interface {
	// GetOrthographicTransform returns the orthographic transformation matrix
	GetOrthographicTransform() geom.Matrix[geom.F32]

	// GetRenderTargetSize returns the size of the render target
	GetRenderTargetSize() geom.Size[geom.F32]
}
