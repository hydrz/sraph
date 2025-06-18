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

// GetDevicePixelRatio returns the current device pixel ratio.
func (c *Canvas) GetDevicePixelRatio() Scalar {
	return c.devicePixelRatio
}

// GetBounds returns the canvas bounds.
func (c *Canvas) GetBounds() geom.Rect[Scalar] {
	return c.bounds
}

// GetCurrentTransform returns the current transformation matrix.
func (c *Canvas) GetCurrentTransform() geom.Matrix[Scalar] {
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
	current := c.GetCurrentTransform()
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
