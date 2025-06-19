package display

import (
	"github.com/opensraph/sraph/geom"
)

// Layer represents an off-screen rendering target that can be composited.
// Layers are used for complex rendering operations, caching, and effects.
type Layer interface {
	// Bounds returns the bounds of the layer.
	Bounds() geom.Rect[Scalar]
	// IsOpaque returns true if the layer is fully opaque.
	IsOpaque() bool
	// Canvas returns a canvas for drawing into this layer.
	Canvas() Canvas
	// Snapshot creates a snapshot of the current layer content.
	Snapshot() LayerSnapshot
	// Clear clears the layer with the specified color.
	Clear(color Color)
	// String returns a string representation of the layer.
	String() string
}

// LayerSnapshot represents a snapshot of layer content that can be drawn.
type LayerSnapshot interface {
	// Bounds returns the bounds of the snapshot.
	Bounds() geom.Rect[Scalar]
	// IsOpaque returns true if the snapshot is fully opaque.
	IsOpaque() bool
	// String returns a string representation of the snapshot.
	String() string
}

// LayerTree represents a hierarchical structure of layers for complex scenes.
type LayerTree struct {
	root     *LayerNode
	viewport geom.Rect[Scalar]
	dpr      Scalar // Device pixel ratio
}

// LayerNode represents a node in the layer tree.
type LayerNode struct {
	layer           Layer
	transform       geom.Matrix[Scalar]
	opacity         float32
	children        []*LayerNode
	parent          *LayerNode
	bounds          geom.Rect[Scalar]
	visible         bool
	needsRepainting bool
}

// NewLayerTree creates a new layer tree with the specified viewport.
func NewLayerTree(viewport geom.Rect[Scalar], dpr Scalar) *LayerTree {
	return &LayerTree{
		viewport: viewport,
		dpr:      dpr,
	}
}

// SetRootLayer sets the root layer of the tree.
func (lt *LayerTree) SetRootLayer(layer Layer) {
	lt.root = &LayerNode{
		layer:     layer,
		transform: geom.NewMatrix[Scalar](),
		opacity:   1.0,
		visible:   true,
	}
}

// AddChild adds a child layer node.
func (ln *LayerNode) AddChild(child *LayerNode) {
	child.parent = ln
	ln.children = append(ln.children, child)
}

// RemoveChild removes a child layer node.
func (ln *LayerNode) RemoveChild(child *LayerNode) {
	for i, c := range ln.children {
		if c == child {
			ln.children = append(ln.children[:i], ln.children[i+1:]...)
			child.parent = nil
			break
		}
	}
}

// SetTransform sets the transform for this layer node.
func (ln *LayerNode) SetTransform(transform geom.Matrix[Scalar]) {
	ln.transform = transform
	ln.markNeedsRepainting()
}

// SetOpacity sets the opacity for this layer node.
func (ln *LayerNode) SetOpacity(opacity float32) {
	ln.opacity = opacity
	ln.markNeedsRepainting()
}

// SetVisible sets the visibility for this layer node.
func (ln *LayerNode) SetVisible(visible bool) {
	ln.visible = visible
	ln.markNeedsRepainting()
}

// markNeedsRepainting marks this node and its ancestors as needing repainting.
func (ln *LayerNode) markNeedsRepainting() {
	ln.needsRepainting = true
	if ln.parent != nil {
		ln.parent.markNeedsRepainting()
	}
}

// Paint renders this layer node and its children.
func (ln *LayerNode) Paint(canvas Canvas) {
	if !ln.visible || ln.layer == nil {
		return
	}

	// Save canvas state
	canvas.Save()

	// Apply transform and opacity
	canvas.TransformFullPerspective(ln.transform)

	// Draw the layer
	paint := NewPaint()
	if ln.opacity < 1.0 {
		paint = paint.SetColor(paint.Color().WithAlpha(ln.opacity))
	}

	// TODO: Draw layer content

	// Paint children
	for _, child := range ln.children {
		child.Paint(canvas)
	}

	// Restore canvas state
	canvas.Restore()
}

// OffscreenLayer represents a layer that renders to an off-screen buffer.
type OffscreenLayer struct {
	bounds   geom.Rect[Scalar]
	canvas   *Canvas
	snapshot LayerSnapshot
	isOpaque bool
}

// NewOffscreenLayer creates a new off-screen layer.
func NewOffscreenLayer(bounds geom.Rect[Scalar], isOpaque bool) *OffscreenLayer {
	return &OffscreenLayer{
		bounds:   bounds,
		canvas:   NewCanvas(bounds),
		isOpaque: isOpaque,
	}
}

// Bounds implements Layer interface.
func (ol *OffscreenLayer) Bounds() geom.Rect[Scalar] {
	return ol.bounds
}

// IsOpaque implements Layer interface.
func (ol *OffscreenLayer) IsOpaque() bool {
	return ol.isOpaque
}

// Canvas implements Layer interface.
func (ol *OffscreenLayer) Canvas() Canvas {
	return *ol.canvas
}

// Snapshot implements Layer interface.
func (ol *OffscreenLayer) Snapshot() LayerSnapshot {
	// TODO: Create actual snapshot of layer content
	return &BasicLayerSnapshot{
		bounds:   ol.bounds,
		isOpaque: ol.isOpaque,
	}
}

// Clear implements Layer interface.
func (ol *OffscreenLayer) Clear(color Color) {
	// TODO: Clear layer with color
}

// String implements Layer interface.
func (ol *OffscreenLayer) String() string {
	return "OffscreenLayer{...}"
}

// BasicLayerSnapshot represents a basic implementation of LayerSnapshot.
type BasicLayerSnapshot struct {
	bounds   geom.Rect[Scalar]
	isOpaque bool
}

// Bounds implements LayerSnapshot interface.
func (bls *BasicLayerSnapshot) Bounds() geom.Rect[Scalar] {
	return bls.bounds
}

// IsOpaque implements LayerSnapshot interface.
func (bls *BasicLayerSnapshot) IsOpaque() bool {
	return bls.isOpaque
}

// String implements LayerSnapshot interface.
func (bls *BasicLayerSnapshot) String() string {
	return "BasicLayerSnapshot{...}"
}

// RenderCache provides caching for expensive rendering operations.
type RenderCache struct {
	entries map[string]CacheEntry
	maxSize int
	size    int
}

// CacheEntry represents a cached rendering result.
type CacheEntry struct {
	snapshot LayerSnapshot
	bounds   geom.Rect[Scalar]
	hash     string
	accessed int64 // Timestamp for LRU
}

// NewRenderCache creates a new render cache with the specified maximum size.
func NewRenderCache(maxSize int) *RenderCache {
	return &RenderCache{
		entries: make(map[string]CacheEntry),
		maxSize: maxSize,
	}
}

// Get retrieves a cached entry by key.
func (rc *RenderCache) Get(key string) (LayerSnapshot, bool) {
	entry, exists := rc.entries[key]
	if exists {
		// TODO: Update access time for LRU
		return entry.snapshot, true
	}
	return nil, false
}

// Put stores a snapshot in the cache.
func (rc *RenderCache) Put(key string, snapshot LayerSnapshot, bounds geom.Rect[Scalar]) {
	// TODO: Implement cache eviction if needed
	rc.entries[key] = CacheEntry{
		snapshot: snapshot,
		bounds:   bounds,
		hash:     key,
		// accessed: current timestamp
	}
	rc.size++
}

// Clear clears all entries from the cache.
func (rc *RenderCache) Clear() {
	rc.entries = make(map[string]CacheEntry)
	rc.size = 0
}

// Size returns the current number of cached entries.
func (rc *RenderCache) Size() int {
	return rc.size
}

// BackdropFilter represents a filter that can be applied to backdrop content.
type BackdropFilter interface {
	// Apply applies the filter to the backdrop.
	Apply(backdrop LayerSnapshot) LayerSnapshot
	// Bounds returns the bounds that this filter would expand content to.
	Bounds(inputBounds geom.Rect[Scalar]) geom.Rect[Scalar]
	// String returns a string representation of the filter.
	String() string
}

// BlurBackdropFilter applies a blur effect to backdrop content.
type BlurBackdropFilter struct {
	sigmaX   float32
	sigmaY   float32
	tileMode TileMode
}

// NewBlurBackdropFilter creates a new blur backdrop filter.
func NewBlurBackdropFilter(sigmaX, sigmaY float32, tileMode TileMode) *BlurBackdropFilter {
	return &BlurBackdropFilter{
		sigmaX:   sigmaX,
		sigmaY:   sigmaY,
		tileMode: tileMode,
	}
}

// Apply implements BackdropFilter interface.
func (bbf *BlurBackdropFilter) Apply(backdrop LayerSnapshot) LayerSnapshot {
	// TODO: Implement backdrop blur
	return backdrop
}

// Bounds implements BackdropFilter interface.
func (bbf *BlurBackdropFilter) Bounds(inputBounds geom.Rect[Scalar]) geom.Rect[Scalar] {
	expansion := Scalar(3.0 * max(bbf.sigmaX, bbf.sigmaY))
	return geom.NewRect(
		inputBounds.Left-expansion,
		inputBounds.Top-expansion,
		inputBounds.Right+expansion,
		inputBounds.Bottom+expansion,
	)
}

// String implements BackdropFilter interface.
func (bbf *BlurBackdropFilter) String() string {
	return "BlurBackdropFilter{...}"
}

// CompositingLayer handles complex compositing operations.
type CompositingLayer struct {
	layer          Layer
	blendMode      BlendMode
	opacity        float32
	clipBounds     *geom.Rect[Scalar]
	backdropFilter BackdropFilter
}

// NewCompositingLayer creates a new compositing layer.
func NewCompositingLayer(layer Layer, blendMode BlendMode, opacity float32) *CompositingLayer {
	return &CompositingLayer{
		layer:     layer,
		blendMode: blendMode,
		opacity:   opacity,
	}
}

// SetBackdropFilter sets the backdrop filter for this compositing layer.
func (cl *CompositingLayer) SetBackdropFilter(filter BackdropFilter) {
	cl.backdropFilter = filter
}

// SetClipBounds sets the clip bounds for this compositing layer.
func (cl *CompositingLayer) SetClipBounds(bounds geom.Rect[Scalar]) {
	cl.clipBounds = &bounds
}

// Composite composites this layer onto the destination canvas.
func (cl *CompositingLayer) Composite(canvas Canvas, backdrop LayerSnapshot) {
	// TODO: Implement complex compositing with backdrop filters
}
