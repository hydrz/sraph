package dl

import (
	"github.com/opensraph/sraph/geom"
)

// DisplayListBuilder is a concrete implementation of OpRecorder that builds a DisplayList.
type DisplayListBuilder struct {
	// operations holds the recorded operations
	operations []Operation
	// bounds tracks the cumulative bounds of all operations
	bounds geom.Rect[Scalar]
	// hasBounds indicates if bounds have been set
	hasBounds bool
	// transformStack holds the current transformation matrix stack
	transformStack []geom.Matrix[Scalar]
	// clipStack holds the current clipping region stack
	clipStack []ClipEntry
	// paintStack holds the current paint state (for save/restore)
	paintStack []Paint
	// saveCount tracks the number of save operations
	saveCount int
}

// ClipEntry represents a clipping operation in the clip stack.
type ClipEntry struct {
	rect          *geom.Rect[Scalar]
	roundRect     *geom.RoundRect[Scalar]
	path          geom.PathSource[Scalar]
	clipOp        ClipOp
	isAntiAliased bool
}

// NewDisplayListBuilder creates a new DisplayListBuilder.
func NewDisplayListBuilder() *DisplayListBuilder {
	return &DisplayListBuilder{
		operations:     make([]Operation, 0),
		bounds:         geom.Rect[Scalar]{},
		transformStack: []geom.Matrix[Scalar]{geom.NewMatrix[Scalar]()},
		clipStack:      make([]ClipEntry, 0),
		paintStack:     make([]Paint, 0),
		saveCount:      0,
	}
}

// Build finalizes the recording and returns a DisplayList.
func (b *DisplayListBuilder) Build() *DisplayList {
	// Create a copy of operations to ensure immutability
	ops := make([]Operation, len(b.operations))
	copy(ops, b.operations)

	bounds := b.bounds
	if !b.hasBounds {
		// If no bounds were explicitly set, use zero bounds
		bounds = geom.Rect[Scalar]{}
	}

	return NewDisplayList(ops, bounds)
}

// GetBounds returns the current bounds of the recorded operations.
func (b *DisplayListBuilder) GetBounds() geom.Rect[Scalar] {
	return b.bounds
}

// Reset resets the builder to an empty state.
func (b *DisplayListBuilder) Reset() {
	b.operations = b.operations[:0]
	b.bounds = geom.Rect[Scalar]{}
	b.hasBounds = false
	b.transformStack = b.transformStack[:1] // Keep identity matrix
	b.transformStack[0] = geom.NewMatrix[Scalar]()
	b.clipStack = b.clipStack[:0]
	b.paintStack = b.paintStack[:0]
	b.saveCount = 0
}

// updateBounds updates the cumulative bounds with a new rectangle.
func (b *DisplayListBuilder) updateBounds(newBounds geom.Rect[Scalar]) {
	if !b.hasBounds {
		b.bounds = newBounds
		b.hasBounds = true
	} else {
		b.bounds = b.bounds.Union(newBounds)
	}
}

// getCurrentTransform returns the current transformation matrix.
func (b *DisplayListBuilder) getCurrentTransform() geom.Matrix[Scalar] {
	if len(b.transformStack) == 0 {
		return geom.NewMatrix[Scalar]()
	}
	return b.transformStack[len(b.transformStack)-1]
}

// pushTransform pushes a new transformation matrix onto the stack.
func (b *DisplayListBuilder) pushTransform(matrix geom.Matrix[Scalar]) {
	b.transformStack = append(b.transformStack, matrix)
}

// popTransform pops the top transformation matrix from the stack.
func (b *DisplayListBuilder) popTransform() {
	if len(b.transformStack) > 1 {
		b.transformStack = b.transformStack[:len(b.transformStack)-1]
	}
}

// recordOperation adds an operation to the list and updates bounds if necessary.
func (b *DisplayListBuilder) recordOperation(op Operation) {
	b.operations = append(b.operations, op)

	// Update bounds if the operation provides them
	if bounds := op.GetBounds(); bounds != nil {
		// Transform bounds by current transformation
		transform := b.getCurrentTransform()
		transformedBounds := bounds.TransformBounds(transform)
		b.updateBounds(transformedBounds)
	}
}

// Implementation of OpReceiver interface

// Save implements OpReceiver.Save.
func (b *DisplayListBuilder) Save() {
	op := &SaveOp{}
	b.recordOperation(op)
	b.saveCount++
}

// Restore implements OpReceiver.Restore.
func (b *DisplayListBuilder) Restore() {
	if b.saveCount > 0 {
		op := &RestoreOp{}
		b.recordOperation(op)
		b.saveCount--
		// Pop transform stack
		b.popTransform()
	}
}

// SaveLayer implements OpReceiver.SaveLayer.
func (b *DisplayListBuilder) SaveLayer(bounds *geom.Rect[Scalar], paint *Paint) {
	op := &SaveLayerOp{
		Bounds: bounds,
		Paint:  paint,
	}
	b.recordOperation(op)
	b.saveCount++
}

// Translate implements OpReceiver.Translate.
func (b *DisplayListBuilder) Translate(dx, dy Scalar) {
	op := &TranslateOp{DX: dx, DY: dy}
	b.recordOperation(op)

	// Update transformation stack
	current := b.getCurrentTransform()
	translated := current.Translate2D(geom.Vector2[Scalar]{X: dx, Y: dy})
	b.transformStack[len(b.transformStack)-1] = translated
}

// Scale implements OpReceiver.Scale.
func (b *DisplayListBuilder) Scale(sx, sy Scalar) {
	op := &ScaleOp{SX: sx, SY: sy}
	b.recordOperation(op)

	// Update transformation stack
	current := b.getCurrentTransform()
	scaled := current.Scale2D(geom.Vector2[Scalar]{X: sx, Y: sy})
	b.transformStack[len(b.transformStack)-1] = scaled
}

// Rotate implements OpReceiver.Rotate.
func (b *DisplayListBuilder) Rotate(radians geom.Radians) {
	op := &RotateOp{Radians: radians}
	b.recordOperation(op)

	// Update transformation stack
	current := b.getCurrentTransform()
	rotated := current.RotateZ(radians)
	b.transformStack[len(b.transformStack)-1] = rotated
}

// Skew implements OpReceiver.Skew.
func (b *DisplayListBuilder) Skew(sx, sy Scalar) {
	op := &SkewOp{SX: sx, SY: sy}
	b.recordOperation(op)

	// Update transformation stack
	current := b.getCurrentTransform()
	skewed := current.Skew(sx, sy)
	b.transformStack[len(b.transformStack)-1] = skewed
}

// Transform2DAffine implements OpReceiver.Transform2DAffine.
func (b *DisplayListBuilder) Transform2DAffine(mxx, mxy, myx, myy, mxt, myt Scalar) {
	op := &Transform2DAffineOp{
		MXX: mxx, MXY: mxy,
		MYX: myx, MYY: myy,
		MXT: mxt, MYT: myt,
	}
	b.recordOperation(op)

	// Update transformation stack
	current := b.getCurrentTransform()
	affineMatrix := geom.Matrix[Scalar]{
		mxx, mxy, mxt,
		myx, myy, myt,
		0, 0, 1,
	}
	transformed := current.Mul(affineMatrix)
	b.transformStack[len(b.transformStack)-1] = transformed
}

// TransformFullPerspective implements OpReceiver.TransformFullPerspective.
func (b *DisplayListBuilder) TransformFullPerspective(matrix geom.Matrix[Scalar]) {
	op := &TransformFullPerspectiveOp{Matrix: matrix}
	b.recordOperation(op)

	// Update transformation stack
	current := b.getCurrentTransform()
	transformed := current.Mul(matrix)
	b.transformStack[len(b.transformStack)-1] = transformed
}

// ClipRect implements OpReceiver.ClipRect.
func (b *DisplayListBuilder) ClipRect(rect geom.Rect[Scalar], clipOp ClipOp, isAntiAliased bool) {
	op := &ClipRectOp{
		Rect:          rect,
		ClipOp:        clipOp,
		IsAntiAliased: isAntiAliased,
	}
	b.recordOperation(op)

	// Update clip stack
	b.clipStack = append(b.clipStack, ClipEntry{
		rect:          &rect,
		clipOp:        clipOp,
		isAntiAliased: isAntiAliased,
	})
}

// ClipRRect implements OpReceiver.ClipRRect.
func (b *DisplayListBuilder) ClipRRect(rrect geom.RoundRect[Scalar], clipOp ClipOp, isAntiAliased bool) {
	op := &ClipRRectOp{
		RRect:         rrect,
		ClipOp:        clipOp,
		IsAntiAliased: isAntiAliased,
	}
	b.recordOperation(op)

	// Update clip stack
	b.clipStack = append(b.clipStack, ClipEntry{
		roundRect:     &rrect,
		clipOp:        clipOp,
		isAntiAliased: isAntiAliased,
	})
}

// ClipPath implements OpReceiver.ClipPath.
func (b *DisplayListBuilder) ClipPath(path geom.PathSource[Scalar], clipOp ClipOp, isAntiAliased bool) {
	op := &ClipPathOp{
		Path:          path,
		ClipOp:        clipOp,
		IsAntiAliased: isAntiAliased,
	}
	b.recordOperation(op)

	// Update clip stack
	b.clipStack = append(b.clipStack, ClipEntry{
		path:          path,
		clipOp:        clipOp,
		isAntiAliased: isAntiAliased,
	})
}

// DrawPaint implements OpReceiver.DrawPaint.
func (b *DisplayListBuilder) DrawPaint(paint Paint) {
	op := &DrawPaintOp{Paint: paint}
	b.recordOperation(op)
}

// DrawColor implements OpReceiver.DrawColor.
func (b *DisplayListBuilder) DrawColor(color Color, blendMode BlendMode) {
	op := &DrawColorOp{Color: color, BlendMode: blendMode}
	b.recordOperation(op)
}

// DrawLine implements OpReceiver.DrawLine.
func (b *DisplayListBuilder) DrawLine(p0, p1 geom.Point[Scalar], paint Paint) {
	op := &DrawLineOp{P0: p0, P1: p1, Paint: paint}
	b.recordOperation(op)
}

// DrawRect implements OpReceiver.DrawRect.
func (b *DisplayListBuilder) DrawRect(rect geom.Rect[Scalar], paint Paint) {
	op := &DrawRectOp{Rect: rect, Paint: paint}
	b.recordOperation(op)
}

// DrawOval implements OpReceiver.DrawOval.
func (b *DisplayListBuilder) DrawOval(bounds geom.Rect[Scalar], paint Paint) {
	op := &DrawOvalOp{Bounds: bounds, Paint: paint}
	b.recordOperation(op)
}

// DrawCircle implements OpReceiver.DrawCircle.
func (b *DisplayListBuilder) DrawCircle(center geom.Point[Scalar], radius Scalar, paint Paint) {
	op := &DrawCircleOp{Center: center, Radius: radius, Paint: paint}
	b.recordOperation(op)
}

// DrawRRect implements OpReceiver.DrawRRect.
func (b *DisplayListBuilder) DrawRRect(rrect geom.RoundRect[Scalar], paint Paint) {
	op := &DrawRRectOp{RRect: rrect, Paint: paint}
	b.recordOperation(op)
}

// DrawDRRect implements OpReceiver.DrawDRRect.
func (b *DisplayListBuilder) DrawDRRect(outer, inner geom.RoundRect[Scalar], paint Paint) {
	op := &DrawDRRectOp{Outer: outer, Inner: inner, Paint: paint}
	b.recordOperation(op)
}

// DrawPath implements OpReceiver.DrawPath.
func (b *DisplayListBuilder) DrawPath(path geom.PathSource[Scalar], paint Paint) {
	op := &DrawPathOp{Path: path, Paint: paint}
	b.recordOperation(op)
}

// DrawArc implements OpReceiver.DrawArc.
func (b *DisplayListBuilder) DrawArc(bounds geom.Rect[Scalar], start, sweep Scalar, useCenter bool, paint Paint) {
	op := &DrawArcOp{
		Bounds:    bounds,
		Start:     start,
		Sweep:     sweep,
		UseCenter: useCenter,
		Paint:     paint,
	}
	b.recordOperation(op)
}

// DrawPoints implements OpReceiver.DrawPoints.
func (b *DisplayListBuilder) DrawPoints(mode PointMode, points []geom.Point[Scalar], paint Paint) {
	op := &DrawPointsOp{Mode: mode, Points: points, Paint: paint}
	b.recordOperation(op)
}

// DrawVertices implements OpReceiver.DrawVertices.
func (b *DisplayListBuilder) DrawVertices(vertices *Vertices, blendMode BlendMode, paint Paint) {
	op := &DrawVerticesOp{Vertices: vertices, BlendMode: blendMode, Paint: paint}
	b.recordOperation(op)
}
