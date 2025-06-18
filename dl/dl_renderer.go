package dl

import (
	"github.com/opensraph/sraph/geom"
)

// Renderer defines the interface for rendering display list operations.
// Different backends (GPU, CPU, etc.) can implement this interface.
type Renderer interface {
	// BeginFrame starts a new frame for rendering.
	BeginFrame(bounds geom.Rect[Scalar]) error

	// EndFrame completes the current frame.
	EndFrame() error

	// Clear clears the render target with the specified color.
	Clear(color Color)

	// SetTransform sets the current transformation matrix.
	SetTransform(matrix geom.Matrix[Scalar])

	// PushClip pushes a new clipping region.
	PushClip(clip ClipRegion)

	// PopClip removes the top clipping region.
	PopClip()

	// Rendering operations
	RenderPaint(paint Paint, bounds geom.Rect[Scalar])
	RenderRect(rect geom.Rect[Scalar], paint Paint)
	RenderRoundRect(rrect geom.RoundRect[Scalar], paint Paint)
	RenderCircle(center geom.Point[Scalar], radius Scalar, paint Paint)
	RenderPath(path *Path, paint Paint)
	RenderLine(p0, p1 geom.Point[Scalar], paint Paint)
	RenderVertices(vertices *Vertices, blendMode BlendMode, paint Paint)

	// Resource management
	CreateTexture(width, height int, data []byte) (TextureID, error)
	DeleteTexture(id TextureID)

	// State queries
	GetMaxTextureSize() int
	SupportsAntiAliasing() bool
}

// ClipRegion represents a clipping region.
type ClipRegion struct {
	Type          ClipType
	Rect          *geom.Rect[Scalar]
	RoundRect     *geom.RoundRect[Scalar]
	Path          *Path
	IsAntiAliased bool
	Operation     ClipOp
}

// TextureID represents a unique identifier for a texture resource.
type TextureID uint32

// RenderContext provides rendering context and state management.
type RenderContext struct {
	renderer       Renderer
	transformStack []geom.Matrix[Scalar]
	clipStack      []ClipRegion
	bounds         geom.Rect[Scalar]
}

// NewRenderContext creates a new RenderContext with the given renderer.
func NewRenderContext(renderer Renderer) *RenderContext {
	return &RenderContext{
		renderer:       renderer,
		transformStack: []geom.Matrix[Scalar]{geom.NewMatrix[Scalar]()},
		clipStack:      make([]ClipRegion, 0),
	}
}

// BeginFrame starts rendering a new frame.
func (rc *RenderContext) BeginFrame(bounds geom.Rect[Scalar]) error {
	rc.bounds = bounds
	return rc.renderer.BeginFrame(bounds)
}

// EndFrame completes the current frame.
func (rc *RenderContext) EndFrame() error {
	return rc.renderer.EndFrame()
}

// RenderDisplayList renders a complete display list.
func (rc *RenderContext) RenderDisplayList(dl *DisplayList) error {
	if dl == nil || dl.IsEmpty() {
		return nil
	}

	// Set up initial state
	rc.renderer.SetTransform(rc.getCurrentTransform())

	// Render all operations
	dl.Dispatch(rc)

	return nil
}

// getCurrentTransform returns the current transformation matrix.
func (rc *RenderContext) getCurrentTransform() geom.Matrix[Scalar] {
	if len(rc.transformStack) == 0 {
		return geom.NewMatrix[Scalar]()
	}
	return rc.transformStack[len(rc.transformStack)-1]
}

// Implementation of OpReceiver interface for RenderContext

// Save implements OpReceiver.Save.
func (rc *RenderContext) Save() {
	// Duplicate current transform
	current := rc.getCurrentTransform()
	rc.transformStack = append(rc.transformStack, current)
}

// Restore implements OpReceiver.Restore.
func (rc *RenderContext) Restore() {
	if len(rc.transformStack) > 1 {
		rc.transformStack = rc.transformStack[:len(rc.transformStack)-1]
		rc.renderer.SetTransform(rc.getCurrentTransform())
	}
}

// SaveLayer implements OpReceiver.SaveLayer.
func (rc *RenderContext) SaveLayer(bounds *geom.Rect[Scalar], paint *Paint) {
	// TODO: Implement layer rendering
	rc.Save()
}

// Translate implements OpReceiver.Translate.
func (rc *RenderContext) Translate(dx, dy Scalar) {
	current := rc.getCurrentTransform()
	translated := current.Translate2D(geom.Vector2[Scalar]{X: dx, Y: dy})
	rc.transformStack[len(rc.transformStack)-1] = translated
	rc.renderer.SetTransform(translated)
}

// Scale implements OpReceiver.Scale.
func (rc *RenderContext) Scale(sx, sy Scalar) {
	current := rc.getCurrentTransform()
	scaled := current.Scale2D(geom.Vector2[Scalar]{X: sx, Y: sy})
	rc.transformStack[len(rc.transformStack)-1] = scaled
	rc.renderer.SetTransform(scaled)
}

// Rotate implements OpReceiver.Rotate.
func (rc *RenderContext) Rotate(radians geom.Radians) {
	current := rc.getCurrentTransform()
	rotated := current.RotateZ(geom.Radians(radians))
	rc.transformStack[len(rc.transformStack)-1] = rotated
	rc.renderer.SetTransform(rotated)
}

// Skew implements OpReceiver.Skew.
func (rc *RenderContext) Skew(sx, sy Scalar) {
	skewMatrix := geom.Matrix[Scalar]{
		1, sx, 0, 0,
		sy, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	current := rc.getCurrentTransform()
	skewed := current.Mul(skewMatrix)
	rc.transformStack[len(rc.transformStack)-1] = skewed
	rc.renderer.SetTransform(skewed)
}

// Transform2DAffine implements OpReceiver.Transform2DAffine.
func (rc *RenderContext) Transform2DAffine(mxx, mxy, myx, myy, mxt, myt Scalar) {
	affineMatrix := geom.Matrix[Scalar]{
		mxx, mxy, mxt, 0,
		myx, myy, myt, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	current := rc.getCurrentTransform()
	transformed := current.Mul(affineMatrix)
	rc.transformStack[len(rc.transformStack)-1] = transformed
	rc.renderer.SetTransform(transformed)
}

// TransformFullPerspective implements OpReceiver.TransformFullPerspective.
func (rc *RenderContext) TransformFullPerspective(matrix geom.Matrix[Scalar]) {
	current := rc.getCurrentTransform()
	transformed := current.Mul(matrix)
	rc.transformStack[len(rc.transformStack)-1] = transformed
	rc.renderer.SetTransform(transformed)
}

// ClipRect implements OpReceiver.ClipRect.
func (rc *RenderContext) ClipRect(rect geom.Rect[Scalar], clipOp ClipOp, isAntiAliased bool) {
	clip := ClipRegion{
		Type:          ClipTypeRect,
		Rect:          &rect,
		IsAntiAliased: isAntiAliased,
		Operation:     clipOp,
	}
	rc.clipStack = append(rc.clipStack, clip)
	rc.renderer.PushClip(clip)
}

// ClipRRect implements OpReceiver.ClipRRect.
func (rc *RenderContext) ClipRRect(rrect geom.RoundRect[Scalar], clipOp ClipOp, isAntiAliased bool) {
	clip := ClipRegion{
		Type:          ClipTypeRoundRect,
		RoundRect:     &rrect,
		IsAntiAliased: isAntiAliased,
		Operation:     clipOp,
	}
	rc.clipStack = append(rc.clipStack, clip)
	rc.renderer.PushClip(clip)
}

// ClipPath implements OpReceiver.ClipPath.
func (rc *RenderContext) ClipPath(path geom.PathSource[Scalar], clipOp ClipOp, isAntiAliased bool) {
	// TODO: Convert PathSource to Path
	clip := ClipRegion{
		Type:          ClipTypePath,
		Path:          nil, // TODO: Convert path
		IsAntiAliased: isAntiAliased,
		Operation:     clipOp,
	}
	rc.clipStack = append(rc.clipStack, clip)
	rc.renderer.PushClip(clip)
}

// DrawPaint implements OpReceiver.DrawPaint.
func (rc *RenderContext) DrawPaint(paint Paint) {
	rc.renderer.RenderPaint(paint, rc.bounds)
}

// DrawColor implements OpReceiver.DrawColor.
func (rc *RenderContext) DrawColor(color Color, blendMode BlendMode) {
	paint := NewPaint()
	paint.SetColor(color)
	paint.SetBlendMode(blendMode)
	rc.renderer.RenderPaint(paint, rc.bounds)
}

// DrawLine implements OpReceiver.DrawLine.
func (rc *RenderContext) DrawLine(p0, p1 geom.Point[Scalar], paint Paint) {
	rc.renderer.RenderLine(p0, p1, paint)
}

// DrawRect implements OpReceiver.DrawRect.
func (rc *RenderContext) DrawRect(rect geom.Rect[Scalar], paint Paint) {
	rc.renderer.RenderRect(rect, paint)
}

// DrawOval implements OpReceiver.DrawOval.
func (rc *RenderContext) DrawOval(bounds geom.Rect[Scalar], paint Paint) {
	// Create an ellipse path and render it
	pb := NewPathBuilder()
	pb.AddEllipse(bounds)
	path := pb.Build()
	rc.renderer.RenderPath(path, paint)
}

// DrawCircle implements OpReceiver.DrawCircle.
func (rc *RenderContext) DrawCircle(center geom.Point[Scalar], radius Scalar, paint Paint) {
	rc.renderer.RenderCircle(center, radius, paint)
}

// DrawRRect implements OpReceiver.DrawRRect.
func (rc *RenderContext) DrawRRect(rrect geom.RoundRect[Scalar], paint Paint) {
	rc.renderer.RenderRoundRect(rrect, paint)
}

// DrawDRRect implements OpReceiver.DrawDRRect.
func (rc *RenderContext) DrawDRRect(outer, inner geom.RoundRect[Scalar], paint Paint) {
	// TODO: Implement double round rect rendering
	// This typically involves rendering the outer shape and then subtracting the inner
	rc.renderer.RenderRoundRect(outer, paint)
}

// DrawPath implements OpReceiver.DrawPath.
func (rc *RenderContext) DrawPath(path geom.PathSource[Scalar], paint Paint) {
	// TODO: Convert PathSource to Path
	// For now, create an empty path
	pb := NewPathBuilder()
	builtPath := pb.Build()
	rc.renderer.RenderPath(builtPath, paint)
}

// DrawArc implements OpReceiver.DrawArc.
func (rc *RenderContext) DrawArc(bounds geom.Rect[Scalar], start, sweep Scalar, useCenter bool, paint Paint) {
	// TODO: Implement arc rendering by creating a path
	pb := NewPathBuilder()
	// Arc implementation would go here
	path := pb.Build()
	rc.renderer.RenderPath(path, paint)
}

// DrawPoints implements OpReceiver.DrawPoints.
func (rc *RenderContext) DrawPoints(mode PointMode, points []geom.Point[Scalar], paint Paint) {
	// TODO: Implement point rendering based on mode
	switch mode {
	case PointModePoints:
		// Draw each point as a small circle
		for _, pt := range points {
			rc.renderer.RenderCircle(pt, 1, paint) // 1 pixel radius
		}
	case PointModeLines:
		// Draw lines between consecutive pairs
		for i := 0; i < len(points)-1; i += 2 {
			if i+1 < len(points) {
				rc.renderer.RenderLine(points[i], points[i+1], paint)
			}
		}
	case PointModePolygon:
		// Draw lines between consecutive points
		for i := 0; i < len(points)-1; i++ {
			rc.renderer.RenderLine(points[i], points[i+1], paint)
		}
	}
}

// DrawVertices implements OpReceiver.DrawVertices.
func (rc *RenderContext) DrawVertices(vertices *Vertices, blendMode BlendMode, paint Paint) {
	rc.renderer.RenderVertices(vertices, blendMode, paint)
}
