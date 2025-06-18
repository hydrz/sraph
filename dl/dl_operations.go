package dl

import "github.com/opensraph/sraph/geom"

// Concrete operation types that implement the Operation interface

// SaveOp represents a save operation.
type SaveOp struct{}

// Invoke implements Operation.Invoke for SaveOp.
func (op *SaveOp) Invoke(receiver OpReceiver) {
	receiver.Save()
}

// GetBounds implements Operation.GetBounds for SaveOp.
func (op *SaveOp) GetBounds() *geom.Rect[Scalar] {
	return nil // Save operations don't have bounds
}

// GetFlags implements Operation.GetFlags for SaveOp.
func (op *SaveOp) GetFlags() AttributeFlags {
	return AttrFlagNone
}

// RestoreOp represents a restore operation.
type RestoreOp struct{}

// Invoke implements Operation.Invoke for RestoreOp.
func (op *RestoreOp) Invoke(receiver OpReceiver) {
	receiver.Restore()
}

// GetBounds implements Operation.GetBounds for RestoreOp.
func (op *RestoreOp) GetBounds() *geom.Rect[Scalar] {
	return nil // Restore operations don't have bounds
}

// GetFlags implements Operation.GetFlags for RestoreOp.
func (op *RestoreOp) GetFlags() AttributeFlags {
	return AttrFlagNone
}

// SaveLayerOp represents a save layer operation.
type SaveLayerOp struct {
	Bounds *geom.Rect[Scalar]
	Paint  *Paint
}

// Invoke implements Operation.Invoke for SaveLayerOp.
func (op *SaveLayerOp) Invoke(receiver OpReceiver) {
	receiver.SaveLayer(op.Bounds, op.Paint)
}

// GetBounds implements Operation.GetBounds for SaveLayerOp.
func (op *SaveLayerOp) GetBounds() *geom.Rect[Scalar] {
	return op.Bounds
}

// GetFlags implements Operation.GetFlags for SaveLayerOp.
func (op *SaveLayerOp) GetFlags() AttributeFlags {
	if op.Paint != nil {
		return GetAttributeFlags(*op.Paint)
	}
	return AttrFlagNone
}

// TranslateOp represents a translate operation.
type TranslateOp struct {
	DX, DY Scalar
}

// Invoke implements Operation.Invoke for TranslateOp.
func (op *TranslateOp) Invoke(receiver OpReceiver) {
	receiver.Translate(op.DX, op.DY)
}

// GetBounds implements Operation.GetBounds for TranslateOp.
func (op *TranslateOp) GetBounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// GetFlags implements Operation.GetFlags for TranslateOp.
func (op *TranslateOp) GetFlags() AttributeFlags {
	return AttrFlagNone
}

// ScaleOp represents a scale operation.
type ScaleOp struct {
	SX, SY Scalar
}

// Invoke implements Operation.Invoke for ScaleOp.
func (op *ScaleOp) Invoke(receiver OpReceiver) {
	receiver.Scale(op.SX, op.SY)
}

// GetBounds implements Operation.GetBounds for ScaleOp.
func (op *ScaleOp) GetBounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// GetFlags implements Operation.GetFlags for ScaleOp.
func (op *ScaleOp) GetFlags() AttributeFlags {
	return AttrFlagNone
}

// RotateOp represents a rotate operation.
type RotateOp struct {
	Radians geom.Radians
}

// Invoke implements Operation.Invoke for RotateOp.
func (op *RotateOp) Invoke(receiver OpReceiver) {
	receiver.Rotate(op.Radians)
}

// GetBounds implements Operation.GetBounds for RotateOp.
func (op *RotateOp) GetBounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// GetFlags implements Operation.GetFlags for RotateOp.
func (op *RotateOp) GetFlags() AttributeFlags {
	return AttrFlagNone
}

// SkewOp represents a skew operation.
type SkewOp struct {
	SX, SY Scalar
}

// Invoke implements Operation.Invoke for SkewOp.
func (op *SkewOp) Invoke(receiver OpReceiver) {
	receiver.Skew(op.SX, op.SY)
}

// GetBounds implements Operation.GetBounds for SkewOp.
func (op *SkewOp) GetBounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// GetFlags implements Operation.GetFlags for SkewOp.
func (op *SkewOp) GetFlags() AttributeFlags {
	return AttrFlagNone
}

// Transform2DAffineOp represents a 2D affine transform operation.
type Transform2DAffineOp struct {
	MXX, MXY, MYX, MYY, MXT, MYT Scalar
}

// Invoke implements Operation.Invoke for Transform2DAffineOp.
func (op *Transform2DAffineOp) Invoke(receiver OpReceiver) {
	receiver.Transform2DAffine(op.MXX, op.MXY, op.MYX, op.MYY, op.MXT, op.MYT)
}

// GetBounds implements Operation.GetBounds for Transform2DAffineOp.
func (op *Transform2DAffineOp) GetBounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// GetFlags implements Operation.GetFlags for Transform2DAffineOp.
func (op *Transform2DAffineOp) GetFlags() AttributeFlags {
	return AttrFlagNone
}

// TransformFullPerspectiveOp represents a full perspective transform operation.
type TransformFullPerspectiveOp struct {
	Matrix geom.Matrix[Scalar]
}

// Invoke implements Operation.Invoke for TransformFullPerspectiveOp.
func (op *TransformFullPerspectiveOp) Invoke(receiver OpReceiver) {
	receiver.TransformFullPerspective(op.Matrix)
}

// GetBounds implements Operation.GetBounds for TransformFullPerspectiveOp.
func (op *TransformFullPerspectiveOp) GetBounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// GetFlags implements Operation.GetFlags for TransformFullPerspectiveOp.
func (op *TransformFullPerspectiveOp) GetFlags() AttributeFlags {
	return AttrFlagNone
}

// ClipRectOp represents a clip rectangle operation.
type ClipRectOp struct {
	Rect          geom.Rect[Scalar]
	ClipOp        ClipOp
	IsAntiAliased bool
}

// Invoke implements Operation.Invoke for ClipRectOp.
func (op *ClipRectOp) Invoke(receiver OpReceiver) {
	receiver.ClipRect(op.Rect, op.ClipOp, op.IsAntiAliased)
}

// GetBounds implements Operation.GetBounds for ClipRectOp.
func (op *ClipRectOp) GetBounds() *geom.Rect[Scalar] {
	return &op.Rect
}

// GetFlags implements Operation.GetFlags for ClipRectOp.
func (op *ClipRectOp) GetFlags() AttributeFlags {
	flags := AttrFlagNone
	if op.IsAntiAliased {
		flags = flags.WithAttribute(AttrFlagIsAntiAlias)
	}
	return flags
}

// ClipRRectOp represents a clip rounded rectangle operation.
type ClipRRectOp struct {
	RRect         geom.RoundRect[Scalar]
	ClipOp        ClipOp
	IsAntiAliased bool
}

// Invoke implements Operation.Invoke for ClipRRectOp.
func (op *ClipRRectOp) Invoke(receiver OpReceiver) {
	receiver.ClipRRect(op.RRect, op.ClipOp, op.IsAntiAliased)
}

// GetBounds implements Operation.GetBounds for ClipRRectOp.
func (op *ClipRRectOp) GetBounds() *geom.Rect[Scalar] {
	bounds := op.RRect.Bounds()
	return &bounds
}

// GetFlags implements Operation.GetFlags for ClipRRectOp.
func (op *ClipRRectOp) GetFlags() AttributeFlags {
	flags := AttrFlagNone
	if op.IsAntiAliased {
		flags = flags.WithAttribute(AttrFlagIsAntiAlias)
	}
	return flags
}

// ClipPathOp represents a clip path operation.
type ClipPathOp struct {
	Path          geom.PathSource[Scalar]
	ClipOp        ClipOp
	IsAntiAliased bool
}

// Invoke implements Operation.Invoke for ClipPathOp.
func (op *ClipPathOp) Invoke(receiver OpReceiver) {
	receiver.ClipPath(op.Path, op.ClipOp, op.IsAntiAliased)
}

// GetBounds implements Operation.GetBounds for ClipPathOp.
func (op *ClipPathOp) GetBounds() *geom.Rect[Scalar] {
	// TODO: Implement path bounds calculation
	return nil
}

// GetFlags implements Operation.GetFlags for ClipPathOp.
func (op *ClipPathOp) GetFlags() AttributeFlags {
	flags := AttrFlagNone
	if op.IsAntiAliased {
		flags = flags.WithAttribute(AttrFlagIsAntiAlias)
	}
	return flags
}

// DrawPaintOp represents a draw paint operation.
type DrawPaintOp struct {
	Paint Paint
}

// Invoke implements Operation.Invoke for DrawPaintOp.
func (op *DrawPaintOp) Invoke(receiver OpReceiver) {
	receiver.DrawPaint(op.Paint)
}

// GetBounds implements Operation.GetBounds for DrawPaintOp.
func (op *DrawPaintOp) GetBounds() *geom.Rect[Scalar] {
	return nil // Paint operations cover the entire canvas
}

// GetFlags implements Operation.GetFlags for DrawPaintOp.
func (op *DrawPaintOp) GetFlags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawColorOp represents a draw color operation.
type DrawColorOp struct {
	Color     Color
	BlendMode BlendMode
}

// Invoke implements Operation.Invoke for DrawColorOp.
func (op *DrawColorOp) Invoke(receiver OpReceiver) {
	receiver.DrawColor(op.Color, op.BlendMode)
}

// GetBounds implements Operation.GetBounds for DrawColorOp.
func (op *DrawColorOp) GetBounds() *geom.Rect[Scalar] {
	return nil // Color operations cover the entire canvas
}

// GetFlags implements Operation.GetFlags for DrawColorOp.
func (op *DrawColorOp) GetFlags() AttributeFlags {
	return AttrFlagHasColor
}

// DrawLineOp represents a draw line operation.
type DrawLineOp struct {
	P0, P1 geom.Point[Scalar]
	Paint  Paint
}

// Invoke implements Operation.Invoke for DrawLineOp.
func (op *DrawLineOp) Invoke(receiver OpReceiver) {
	receiver.DrawLine(op.P0, op.P1, op.Paint)
}

// GetBounds implements Operation.GetBounds for DrawLineOp.
func (op *DrawLineOp) GetBounds() *geom.Rect[Scalar] {
	minX := op.P0.X
	if op.P1.X < minX {
		minX = op.P1.X
	}
	maxX := op.P0.X
	if op.P1.X > maxX {
		maxX = op.P1.X
	}
	minY := op.P0.Y
	if op.P1.Y < minY {
		minY = op.P1.Y
	}
	maxY := op.P0.Y
	if op.P1.Y > maxY {
		maxY = op.P1.Y
	}

	bounds := geom.NewRect(minX, minY, maxX, maxY)
	return &bounds
}

// GetFlags implements Operation.GetFlags for DrawLineOp.
func (op *DrawLineOp) GetFlags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawRectOp represents a draw rectangle operation.
type DrawRectOp struct {
	Rect  geom.Rect[Scalar]
	Paint Paint
}

// Invoke implements Operation.Invoke for DrawRectOp.
func (op *DrawRectOp) Invoke(receiver OpReceiver) {
	receiver.DrawRect(op.Rect, op.Paint)
}

// GetBounds implements Operation.GetBounds for DrawRectOp.
func (op *DrawRectOp) GetBounds() *geom.Rect[Scalar] {
	return &op.Rect
}

// GetFlags implements Operation.GetFlags for DrawRectOp.
func (op *DrawRectOp) GetFlags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawOvalOp represents a draw oval operation.
type DrawOvalOp struct {
	Bounds geom.Rect[Scalar]
	Paint  Paint
}

// Invoke implements Operation.Invoke for DrawOvalOp.
func (op *DrawOvalOp) Invoke(receiver OpReceiver) {
	receiver.DrawOval(op.Bounds, op.Paint)
}

// GetBounds implements Operation.GetBounds for DrawOvalOp.
func (op *DrawOvalOp) GetBounds() *geom.Rect[Scalar] {
	return &op.Bounds
}

// GetFlags implements Operation.GetFlags for DrawOvalOp.
func (op *DrawOvalOp) GetFlags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawCircleOp represents a draw circle operation.
type DrawCircleOp struct {
	Center geom.Point[Scalar]
	Radius Scalar
	Paint  Paint
}

// Invoke implements Operation.Invoke for DrawCircleOp.
func (op *DrawCircleOp) Invoke(receiver OpReceiver) {
	receiver.DrawCircle(op.Center, op.Radius, op.Paint)
}

// GetBounds implements Operation.GetBounds for DrawCircleOp.
func (op *DrawCircleOp) GetBounds() *geom.Rect[Scalar] {
	bounds := geom.NewRect(
		op.Center.X-op.Radius, op.Center.Y-op.Radius,
		op.Center.X+op.Radius, op.Center.Y+op.Radius,
	)
	return &bounds
}

// GetFlags implements Operation.GetFlags for DrawCircleOp.
func (op *DrawCircleOp) GetFlags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawRRectOp represents a draw rounded rectangle operation.
type DrawRRectOp struct {
	RRect geom.RoundRect[Scalar]
	Paint Paint
}

// Invoke implements Operation.Invoke for DrawRRectOp.
func (op *DrawRRectOp) Invoke(receiver OpReceiver) {
	receiver.DrawRRect(op.RRect, op.Paint)
}

// GetBounds implements Operation.GetBounds for DrawRRectOp.
func (op *DrawRRectOp) GetBounds() *geom.Rect[Scalar] {
	bounds := op.RRect.Bounds()
	return &bounds
}

// GetFlags implements Operation.GetFlags for DrawRRectOp.
func (op *DrawRRectOp) GetFlags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawDRRectOp represents a draw double rounded rectangle operation.
type DrawDRRectOp struct {
	Outer geom.RoundRect[Scalar]
	Inner geom.RoundRect[Scalar]
	Paint Paint
}

// Invoke implements Operation.Invoke for DrawDRRectOp.
func (op *DrawDRRectOp) Invoke(receiver OpReceiver) {
	receiver.DrawDRRect(op.Outer, op.Inner, op.Paint)
}

// GetBounds implements Operation.GetBounds for DrawDRRectOp.
func (op *DrawDRRectOp) GetBounds() *geom.Rect[Scalar] {
	bounds := op.Outer.Bounds()
	return &bounds
}

// GetFlags implements Operation.GetFlags for DrawDRRectOp.
func (op *DrawDRRectOp) GetFlags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawPathOp represents a draw path operation.
type DrawPathOp struct {
	Path  geom.PathSource[Scalar]
	Paint Paint
}

// Invoke implements Operation.Invoke for DrawPathOp.
func (op *DrawPathOp) Invoke(receiver OpReceiver) {
	receiver.DrawPath(op.Path, op.Paint)
}

// GetBounds implements Operation.GetBounds for DrawPathOp.
func (op *DrawPathOp) GetBounds() *geom.Rect[Scalar] {
	// TODO: Implement path bounds calculation
	return nil
}

// GetFlags implements Operation.GetFlags for DrawPathOp.
func (op *DrawPathOp) GetFlags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawArcOp represents a draw arc operation.
type DrawArcOp struct {
	Bounds    geom.Rect[Scalar]
	Start     Scalar
	Sweep     Scalar
	UseCenter bool
	Paint     Paint
}

// Invoke implements Operation.Invoke for DrawArcOp.
func (op *DrawArcOp) Invoke(receiver OpReceiver) {
	receiver.DrawArc(op.Bounds, op.Start, op.Sweep, op.UseCenter, op.Paint)
}

// GetBounds implements Operation.GetBounds for DrawArcOp.
func (op *DrawArcOp) GetBounds() *geom.Rect[Scalar] {
	return &op.Bounds
}

// GetFlags implements Operation.GetFlags for DrawArcOp.
func (op *DrawArcOp) GetFlags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawPointsOp represents a draw points operation.
type DrawPointsOp struct {
	Mode   PointMode
	Points []geom.Point[Scalar]
	Paint  Paint
}

// Invoke implements Operation.Invoke for DrawPointsOp.
func (op *DrawPointsOp) Invoke(receiver OpReceiver) {
	receiver.DrawPoints(op.Mode, op.Points, op.Paint)
}

// GetBounds implements Operation.GetBounds for DrawPointsOp.
func (op *DrawPointsOp) GetBounds() *geom.Rect[Scalar] {
	if len(op.Points) == 0 {
		return nil
	}

	min := op.Points[0]
	max := op.Points[0]

	for _, pt := range op.Points[1:] {
		if pt.X < min.X {
			min.X = pt.X
		}
		if pt.Y < min.Y {
			min.Y = pt.Y
		}
		if pt.X > max.X {
			max.X = pt.X
		}
		if pt.Y > max.Y {
			max.Y = pt.Y
		}
	}

	bounds := geom.NewRect(min.X, min.Y, max.X, max.Y)
	return &bounds
}

// GetFlags implements Operation.GetFlags for DrawPointsOp.
func (op *DrawPointsOp) GetFlags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawVerticesOp represents a draw vertices operation.
type DrawVerticesOp struct {
	Vertices  *Vertices
	BlendMode BlendMode
	Paint     Paint
}

// Invoke implements Operation.Invoke for DrawVerticesOp.
func (op *DrawVerticesOp) Invoke(receiver OpReceiver) {
	receiver.DrawVertices(op.Vertices, op.BlendMode, op.Paint)
}

// GetBounds implements Operation.GetBounds for DrawVerticesOp.
func (op *DrawVerticesOp) GetBounds() *geom.Rect[Scalar] {
	if op.Vertices != nil {
		return op.Vertices.GetBounds()
	}
	return nil
}

// GetFlags implements Operation.GetFlags for DrawVerticesOp.
func (op *DrawVerticesOp) GetFlags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}
