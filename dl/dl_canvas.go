package dl

import (
	"github.com/opensraph/sraph/geom"
)

// Canvas provides a concrete implementation of OpReceiver for immediate rendering.
// This represents a rendering context that can execute display list operations.
type Canvas struct {
	// Transform stack for matrix operations
	transformStack []geom.Matrix[Scalar]
	// Clip stack for clipping operations
	clipStack []ClipState
	// Save/restore state stack
	saveStack []SaveState
	// Current render target bounds
	bounds geom.Rect[Scalar]
	// Current device pixel ratio
	devicePixelRatio Scalar
	// Renderer for actual drawing operations
	renderer Renderer
}

// ClipState represents the current clipping state.
type ClipState struct {
	clipType      ClipType
	rect          *geom.Rect[Scalar]
	roundRect     *geom.RoundRect[Scalar]
	path          geom.PathSource[Scalar]
	isAntiAliased bool
}

// ClipType defines the type of clipping operation.
type ClipType uint8

const (
	// ClipTypeRect indicates rectangle clipping.
	ClipTypeRect ClipType = iota
	// ClipTypeRoundRect indicates rounded rectangle clipping.
	ClipTypeRoundRect
	// ClipTypePath indicates path clipping.
	ClipTypePath
)

// SaveState represents the saved state for save/restore operations.
type SaveState struct {
	transformIndex int
	clipIndex      int
}

// NewCanvas creates a new Canvas with the given bounds.
func NewCanvas(bounds geom.Rect[Scalar]) *Canvas {
	return &Canvas{
		transformStack:   []geom.Matrix[Scalar]{geom.NewMatrix[Scalar]()},
		clipStack:        make([]ClipState, 0),
		saveStack:        make([]SaveState, 0),
		bounds:           bounds,
		devicePixelRatio: 1.0,
	}
}

// SetDevicePixelRatio sets the device pixel ratio for the canvas.
func (c *Canvas) SetDevicePixelRatio(dpr Scalar) {
	c.devicePixelRatio = dpr
}

// DevicePixelRatio returns the current device pixel ratio.
func (c *Canvas) DevicePixelRatio() Scalar {
	return c.devicePixelRatio
}

// Bounds returns the canvas bounds.
func (c *Canvas) Bounds() geom.Rect[Scalar] {
	return c.bounds
}

// CurrentTransform returns the current transformation matrix.
func (c *Canvas) CurrentTransform() geom.Matrix[Scalar] {
	if len(c.transformStack) == 0 {
		return geom.NewMatrix[Scalar]()
	}
	return c.transformStack[len(c.transformStack)-1]
}

// Implementation of OpReceiver interface

// Save implements OpReceiver.Save.
func (c *Canvas) Save() {
	state := SaveState{
		transformIndex: len(c.transformStack) - 1,
		clipIndex:      len(c.clipStack),
	}
	c.saveStack = append(c.saveStack, state)

	// Duplicate current transform for the new save level
	current := c.CurrentTransform()
	c.transformStack = append(c.transformStack, current)
}

// Restore implements OpReceiver.Restore.
func (c *Canvas) Restore() {
	if len(c.saveStack) == 0 {
		return // No saved state to restore
	}

	// Get the last saved state
	state := c.saveStack[len(c.saveStack)-1]
	c.saveStack = c.saveStack[:len(c.saveStack)-1]

	// Restore transform stack
	if state.transformIndex < len(c.transformStack) {
		c.transformStack = c.transformStack[:state.transformIndex+1]
	}

	// Restore clip stack
	if state.clipIndex < len(c.clipStack) {
		c.clipStack = c.clipStack[:state.clipIndex]
	}
}

// SaveLayer implements OpReceiver.SaveLayer.
func (c *Canvas) SaveLayer(bounds *geom.Rect[Scalar], paint *Paint) {
	// Save current state first
	c.Save()

	// TODO: Implement layer rendering
	// This would typically involve:
	// 1. Creating an offscreen render target
	// 2. Setting up the layer bounds
	// 3. Applying paint effects to the layer
}

// Translate implements OpReceiver.Translate.
func (c *Canvas) Translate(dx, dy Scalar) {
	if len(c.transformStack) == 0 {
		c.transformStack = append(c.transformStack, geom.NewMatrix[Scalar]())
	}

	current := c.transformStack[len(c.transformStack)-1]
	translated := current.Translate2D(geom.Vector2[Scalar]{X: dx, Y: dy})
	c.transformStack[len(c.transformStack)-1] = translated
}

// Scale implements OpReceiver.Scale.
func (c *Canvas) Scale(sx, sy Scalar) {
	if len(c.transformStack) == 0 {
		c.transformStack = append(c.transformStack, geom.NewMatrix[Scalar]())
	}

	current := c.transformStack[len(c.transformStack)-1]
	scaled := current.Scale2D(geom.Vector2[Scalar]{X: sx, Y: sy})
	c.transformStack[len(c.transformStack)-1] = scaled
}

// Rotate implements OpReceiver.Rotate.
func (c *Canvas) Rotate(radians geom.Radians) {
	if len(c.transformStack) == 0 {
		c.transformStack = append(c.transformStack, geom.NewMatrix[Scalar]())
	}

	current := c.transformStack[len(c.transformStack)-1]
	rotated := current.RotateZ(radians)
	c.transformStack[len(c.transformStack)-1] = rotated
}

// Skew implements OpReceiver.Skew.
func (c *Canvas) Skew(sx, sy Scalar) {
	if len(c.transformStack) == 0 {
		c.transformStack = append(c.transformStack, geom.NewMatrix[Scalar]())
	}

	// TODO: Implement proper skew transformation
	// For now, we'll create a basic skew matrix manually
	skewMatrix := geom.Matrix[Scalar]{
		1, sx, 0, 0,
		sy, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	current := c.transformStack[len(c.transformStack)-1]
	skewed := current.Mul(skewMatrix)
	c.transformStack[len(c.transformStack)-1] = skewed
}

// Transform2DAffine implements OpReceiver.Transform2DAffine.
func (c *Canvas) Transform2DAffine(mxx, mxy, myx, myy, mxt, myt Scalar) {
	if len(c.transformStack) == 0 {
		c.transformStack = append(c.transformStack, geom.NewMatrix[Scalar]())
	}

	affineMatrix := geom.Matrix[Scalar]{
		mxx, mxy, mxt, 0,
		myx, myy, myt, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	current := c.transformStack[len(c.transformStack)-1]
	transformed := current.Mul(affineMatrix)
	c.transformStack[len(c.transformStack)-1] = transformed
}

// TransformFullPerspective implements OpReceiver.TransformFullPerspective.
func (c *Canvas) TransformFullPerspective(matrix geom.Matrix[Scalar]) {
	if len(c.transformStack) == 0 {
		c.transformStack = append(c.transformStack, geom.NewMatrix[Scalar]())
	}

	current := c.transformStack[len(c.transformStack)-1]
	transformed := current.Mul(matrix)
	c.transformStack[len(c.transformStack)-1] = transformed
}

// ClipRect implements OpReceiver.ClipRect.
func (c *Canvas) ClipRect(rect geom.Rect[Scalar], clipOp ClipOp, isAntiAliased bool) {
	clipState := ClipState{
		clipType:      ClipTypeRect,
		rect:          &rect,
		isAntiAliased: isAntiAliased,
	}
	c.clipStack = append(c.clipStack, clipState)

	// TODO: Apply actual clipping to the rendering context
}

// ClipRRect implements OpReceiver.ClipRRect.
func (c *Canvas) ClipRRect(rrect geom.RoundRect[Scalar], clipOp ClipOp, isAntiAliased bool) {
	clipState := ClipState{
		clipType:      ClipTypeRoundRect,
		roundRect:     &rrect,
		isAntiAliased: isAntiAliased,
	}
	c.clipStack = append(c.clipStack, clipState)

	// TODO: Apply actual clipping to the rendering context
}

// ClipPath implements OpReceiver.ClipPath.
func (c *Canvas) ClipPath(path geom.PathSource[Scalar], clipOp ClipOp, isAntiAliased bool) {
	clipState := ClipState{
		clipType:      ClipTypePath,
		path:          path,
		isAntiAliased: isAntiAliased,
	}
	c.clipStack = append(c.clipStack, clipState)

	// TODO: Apply actual clipping to the rendering context
}

// Drawing operations (these would be implemented by specific rendering backends)

// DrawPaint implements OpReceiver.DrawPaint.
func (c *Canvas) DrawPaint(paint Paint) {
	// TODO: Implement paint drawing
	// This would fill the entire canvas with the paint
}

// DrawColor implements OpReceiver.DrawColor.
func (c *Canvas) DrawColor(color Color, blendMode BlendMode) {
	// TODO: Implement color drawing
	// This would fill the entire canvas with the color using the blend mode
}

// DrawLine implements OpReceiver.DrawLine.
func (c *Canvas) DrawLine(p0, p1 geom.Point[Scalar], paint Paint) {
	// TODO: Implement line drawing
	// This would draw a line from p0 to p1 with the given paint
}

// DrawRect implements OpReceiver.DrawRect.
func (c *Canvas) DrawRect(rect geom.Rect[Scalar], paint Paint) {
	// TODO: Implement rectangle drawing
	// This would draw a rectangle with the given paint
}

// DrawOval implements OpReceiver.DrawOval.
func (c *Canvas) DrawOval(bounds geom.Rect[Scalar], paint Paint) {
	// TODO: Implement oval drawing
	// This would draw an oval within the given bounds
}

// DrawCircle implements OpReceiver.DrawCircle.
func (c *Canvas) DrawCircle(center geom.Point[Scalar], radius Scalar, paint Paint) {
	// TODO: Implement circle drawing
	// This would draw a circle at the given center with the specified radius
}

// DrawRRect implements OpReceiver.DrawRRect.
func (c *Canvas) DrawRRect(rrect geom.RoundRect[Scalar], paint Paint) {
	// TODO: Implement rounded rectangle drawing
	// This would draw a rounded rectangle with the given paint
}

// DrawDRRect implements OpReceiver.DrawDRRect.
func (c *Canvas) DrawDRRect(outer, inner geom.RoundRect[Scalar], paint Paint) {
	// TODO: Implement double rounded rectangle drawing
	// This would draw the outer rrect with the inner rrect cut out
}

// DrawPath implements OpReceiver.DrawPath.
func (c *Canvas) DrawPath(path geom.PathSource[Scalar], paint Paint) {
	// TODO: Implement path drawing
	// This would draw the path with the given paint
}

// DrawArc implements OpReceiver.DrawArc.
func (c *Canvas) DrawArc(bounds geom.Rect[Scalar], start, sweep Scalar, useCenter bool, paint Paint) {
	// TODO: Implement arc drawing
	// This would draw an arc within the bounds from start angle for sweep radians
}

// DrawPoints implements OpReceiver.DrawPoints.
func (c *Canvas) DrawPoints(mode PointMode, points []geom.Point[Scalar], paint Paint) {
	// TODO: Implement points drawing
	// This would draw points according to the specified mode
}

// DrawVertices implements OpReceiver.DrawVertices.
func (c *Canvas) DrawVertices(vertices *Vertices, blendMode BlendMode, paint Paint) {
	// TODO: Implement vertices drawing
	// This would draw the vertices as triangles with the given blend mode and paint
}

// DrawImage draws an image at the specified position.
func (c *Canvas) DrawImage(image Image, position geom.Point[Scalar], paint Paint) {
	// TODO: Implement image drawing
	// This would draw the image at the given position with the paint
}

// DrawImageRect draws an image scaled to fit the destination rectangle.
func (c *Canvas) DrawImageRect(image Image, src, dst geom.Rect[Scalar], paint Paint, constraint SrcRectConstraint) {
	// TODO: Implement image rectangle drawing
	// This would draw the image from src rect to dst rect
}

// DrawImageNine draws a nine-patch image.
func (c *Canvas) DrawImageNine(image Image, center geom.Rect[Scalar], dst geom.Rect[Scalar], paint Paint) {
	// TODO: Implement nine-patch image drawing
	// This would draw the image as a nine-patch to fit the destination
}

// DrawImageWithSampling draws an image with specific sampling options.
func (c *Canvas) DrawImageWithSampling(image Image, position geom.Point[Scalar], sampling SamplingOptions, paint Paint) {
	// TODO: Implement image drawing with sampling
}

// DrawParagraph draws a text paragraph at the specified position.
func (c *Canvas) DrawParagraph(paragraph *Paragraph, position geom.Point[Scalar]) {
	// TODO: Implement paragraph drawing
	paragraph.Paint(*c, position)
}

// DrawTextBlob draws a text blob at the specified position.
func (c *Canvas) DrawTextBlob(text string, position geom.Point[Scalar], font Font, paint Paint) {
	// TODO: Implement text blob drawing
	// This would draw text with the given font and paint
}

// DrawShadow draws a shadow for the given path.
func (c *Canvas) DrawShadow(path geom.PathSource[Scalar], color Color, elevation float32, transparentOccluder, dpr Scalar) {
	// TODO: Implement shadow drawing
	// This would draw a shadow based on the path elevation and lighting
}

// DrawAtlas draws multiple sprites from a texture atlas.
func (c *Canvas) DrawAtlas(atlas Image, transforms []geom.RSTransform[Scalar], texCoords []geom.Rect[Scalar], colors []Color, blendMode BlendMode, paint Paint) {
	// TODO: Implement atlas drawing
	// This would draw multiple sprites efficiently from an atlas texture
}

// DrawPatch draws a patch (cubic Bezier patch) for advanced 3D-style rendering.
func (c *Canvas) DrawPatch(cubics []geom.Point[Scalar], colors []Color, texCoords []geom.Point[Scalar], blendMode BlendMode, paint Paint) {
	// TODO: Implement patch drawing
	// This would draw a cubic Bezier patch with interpolated colors and textures
}

// ComputeShadowBounds computes the bounds of a shadow for the given path.
func (c *Canvas) ComputeShadowBounds(path geom.PathSource[Scalar], elevation float32, dpr Scalar, transform geom.Matrix[Scalar]) geom.Rect[Scalar] {
	// TODO: Implement shadow bounds computation
	// This would calculate the bounds that a shadow would occupy
	return geom.Rect[Scalar]{}
}

// FlushPendingOperations flushes any pending drawing operations to the underlying renderer.
func (c *Canvas) FlushPendingOperations() {
	// TODO: Implement operation flushing
	// This would ensure all queued operations are executed
}

// IsRecording returns true if this canvas is recording operations rather than rendering.
func (c *Canvas) IsRecording() bool {
	return false // Canvas renders immediately, doesn't record
}

// LocalClipBounds returns the current local clip bounds.
func (c *Canvas) LocalClipBounds() geom.Rect[Scalar] {
	// TODO: Implement local clip bounds calculation
	return c.bounds
}

// DeviceClipBounds returns the current device clip bounds.
func (c *Canvas) DeviceClipBounds() geom.Rect[Scalar] {
	// TODO: Implement device clip bounds calculation
	transform := c.CurrentTransform()
	return c.bounds.TransformBounds(transform)
}

// DrawDisplayList draws another display list into this canvas.
func (c *Canvas) DrawDisplayList(displayList *DisplayList, opacity float32) {
	// TODO: Implement display list drawing
	// This would execute all operations from another display list
}

// DrawLayer draws a layer with the given paint.
func (c *Canvas) DrawLayer(layer Layer, paint Paint) {
	// TODO: Implement layer drawing
	// This would draw a pre-rendered layer
}

// CreateLayer creates a new layer for off-screen rendering.
func (c *Canvas) CreateLayer(bounds geom.Rect[Scalar]) Layer {
	// TODO: Implement layer creation
	// This would create a new render target
	return nil
}

// CanvasDrawingStyle defines various drawing style options.
type CanvasDrawingStyle struct {
	compositingOperation  BlendMode
	globalAlpha           float32
	imageSmoothingEnabled bool
	shadowBlur            float32
	shadowColor           Color
	shadowOffsetX         float32
	shadowOffsetY         float32
}

// NewCanvasDrawingStyle creates a new drawing style with default values.
func NewCanvasDrawingStyle() CanvasDrawingStyle {
	return CanvasDrawingStyle{
		compositingOperation:  BlendModeSrcOver,
		globalAlpha:           1.0,
		imageSmoothingEnabled: true,
		shadowColor:           ColorTransparent,
	}
}
