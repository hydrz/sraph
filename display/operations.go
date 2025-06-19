package display

import "github.com/opensraph/sraph/geom"

// Concrete operation types that implement the Operation interface

// SaveOp represents a save operation.
type SaveOp struct{}

// Invoke implements Operation.Invoke for SaveOp.
func (op *SaveOp) Invoke(receiver OpReceiver) {
	receiver.Save()
}

// Bounds implements Operation.Bounds for SaveOp.
func (op *SaveOp) Bounds() *geom.Rect[Scalar] {
	return nil // Save operations don't have bounds
}

// Flags implements Operation.Flags for SaveOp.
func (op *SaveOp) Flags() AttributeFlags {
	return AttrFlagNone
}

// RestoreOp represents a restore operation.
type RestoreOp struct{}

// Invoke implements Operation.Invoke for RestoreOp.
func (op *RestoreOp) Invoke(receiver OpReceiver) {
	receiver.Restore()
}

// Bounds implements Operation.Bounds for RestoreOp.
func (op *RestoreOp) Bounds() *geom.Rect[Scalar] {
	return nil // Restore operations don't have bounds
}

// Flags implements Operation.Flags for RestoreOp.
func (op *RestoreOp) Flags() AttributeFlags {
	return AttrFlagNone
}

// SaveLayerOp represents a save layer operation.
type SaveLayerOp struct {
	Rect  *geom.Rect[Scalar]
	Paint *Paint
}

// Invoke implements Operation.Invoke for SaveLayerOp.
func (op *SaveLayerOp) Invoke(receiver OpReceiver) {
	receiver.SaveLayer(op.Rect, op.Paint)
}

// Bounds implements Operation.Bounds for SaveLayerOp.
func (op *SaveLayerOp) Bounds() *geom.Rect[Scalar] {
	return op.Rect
}

// Flags implements Operation.Flags for SaveLayerOp.
func (op *SaveLayerOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for TranslateOp.
func (op *TranslateOp) Bounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// Flags implements Operation.Flags for TranslateOp.
func (op *TranslateOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for ScaleOp.
func (op *ScaleOp) Bounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// Flags implements Operation.Flags for ScaleOp.
func (op *ScaleOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for RotateOp.
func (op *RotateOp) Bounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// Flags implements Operation.Flags for RotateOp.
func (op *RotateOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for SkewOp.
func (op *SkewOp) Bounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// Flags implements Operation.Flags for SkewOp.
func (op *SkewOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for Transform2DAffineOp.
func (op *Transform2DAffineOp) Bounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// Flags implements Operation.Flags for Transform2DAffineOp.
func (op *Transform2DAffineOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for TransformFullPerspectiveOp.
func (op *TransformFullPerspectiveOp) Bounds() *geom.Rect[Scalar] {
	return nil // Transform operations don't have bounds
}

// Flags implements Operation.Flags for TransformFullPerspectiveOp.
func (op *TransformFullPerspectiveOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for ClipRectOp.
func (op *ClipRectOp) Bounds() *geom.Rect[Scalar] {
	return &op.Rect
}

// Flags implements Operation.Flags for ClipRectOp.
func (op *ClipRectOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for ClipRRectOp.
func (op *ClipRRectOp) Bounds() *geom.Rect[Scalar] {
	bounds := op.RRect.Bounds()
	return &bounds
}

// Flags implements Operation.Flags for ClipRRectOp.
func (op *ClipRRectOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for ClipPathOp.
func (op *ClipPathOp) Bounds() *geom.Rect[Scalar] {
	// TODO: Implement path bounds calculation
	return nil
}

// Flags implements Operation.Flags for ClipPathOp.
func (op *ClipPathOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for DrawPaintOp.
func (op *DrawPaintOp) Bounds() *geom.Rect[Scalar] {
	return nil // Paint operations cover the entire canvas
}

// Flags implements Operation.Flags for DrawPaintOp.
func (op *DrawPaintOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for DrawColorOp.
func (op *DrawColorOp) Bounds() *geom.Rect[Scalar] {
	return nil // Color operations cover the entire canvas
}

// Flags implements Operation.Flags for DrawColorOp.
func (op *DrawColorOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for DrawLineOp.
func (op *DrawLineOp) Bounds() *geom.Rect[Scalar] {
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

// Flags implements Operation.Flags for DrawLineOp.
func (op *DrawLineOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for DrawRectOp.
func (op *DrawRectOp) Bounds() *geom.Rect[Scalar] {
	return &op.Rect
}

// Flags implements Operation.Flags for DrawRectOp.
func (op *DrawRectOp) Flags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawOvalOp represents a draw oval operation.
type DrawOvalOp struct {
	Rect  geom.Rect[Scalar]
	Paint Paint
}

// Invoke implements Operation.Invoke for DrawOvalOp.
func (op *DrawOvalOp) Invoke(receiver OpReceiver) {
	receiver.DrawOval(op.Rect, op.Paint)
}

// Bounds implements Operation.Bounds for DrawOvalOp.
func (op *DrawOvalOp) Bounds() *geom.Rect[Scalar] {
	return &op.Rect
}

// Flags implements Operation.Flags for DrawOvalOp.
func (op *DrawOvalOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for DrawCircleOp.
func (op *DrawCircleOp) Bounds() *geom.Rect[Scalar] {
	bounds := geom.NewRect(
		op.Center.X-op.Radius, op.Center.Y-op.Radius,
		op.Center.X+op.Radius, op.Center.Y+op.Radius,
	)
	return &bounds
}

// Flags implements Operation.Flags for DrawCircleOp.
func (op *DrawCircleOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for DrawRRectOp.
func (op *DrawRRectOp) Bounds() *geom.Rect[Scalar] {
	bounds := op.RRect.Bounds()
	return &bounds
}

// Flags implements Operation.Flags for DrawRRectOp.
func (op *DrawRRectOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for DrawDRRectOp.
func (op *DrawDRRectOp) Bounds() *geom.Rect[Scalar] {
	bounds := op.Outer.Bounds()
	return &bounds
}

// Flags implements Operation.Flags for DrawDRRectOp.
func (op *DrawDRRectOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for DrawPathOp.
func (op *DrawPathOp) Bounds() *geom.Rect[Scalar] {
	// TODO: Implement path bounds calculation
	return nil
}

// Flags implements Operation.Flags for DrawPathOp.
func (op *DrawPathOp) Flags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawArcOp represents a draw arc operation.
type DrawArcOp struct {
	Rect      geom.Rect[Scalar]
	Start     Scalar
	Sweep     Scalar
	UseCenter bool
	Paint     Paint
}

// Invoke implements Operation.Invoke for DrawArcOp.
func (op *DrawArcOp) Invoke(receiver OpReceiver) {
	receiver.DrawArc(op.Rect, op.Start, op.Sweep, op.UseCenter, op.Paint)
}

// Bounds implements Operation.Bounds for DrawArcOp.
func (op *DrawArcOp) Bounds() *geom.Rect[Scalar] {
	return &op.Rect
}

// Flags implements Operation.Flags for DrawArcOp.
func (op *DrawArcOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for DrawPointsOp.
func (op *DrawPointsOp) Bounds() *geom.Rect[Scalar] {
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

// Flags implements Operation.Flags for DrawPointsOp.
func (op *DrawPointsOp) Flags() AttributeFlags {
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

// Bounds implements Operation.Bounds for DrawVerticesOp.
func (op *DrawVerticesOp) Bounds() *geom.Rect[Scalar] {
	if op.Vertices != nil {
		return op.Vertices.Bounds()
	}
	return nil
}

// Flags implements Operation.Flags for DrawVerticesOp.
func (op *DrawVerticesOp) Flags() AttributeFlags {
	return GetAttributeFlags(op.Paint)
}

// DrawImageOp represents a draw image operation.
type DrawImageOp struct {
	Image Image
	Point geom.Point[Scalar]
	Paint Paint
}

// Invoke implements Operation.Invoke for DrawImageOp.
func (op *DrawImageOp) Invoke(receiver OpReceiver) {
	receiver.DrawImage(op.Image, op.Point, op.Paint)
}

// Bounds implements Operation.Bounds for DrawImageOp.
func (op *DrawImageOp) Bounds() *geom.Rect[Scalar] {
	if op.Image != nil {
		bounds := geom.NewRect(
			op.Point.X,
			op.Point.Y,
			op.Point.X+Scalar(op.Image.Width()),
			op.Point.Y+Scalar(op.Image.Height()),
		)
		return &bounds
	}
	return nil
}

// Flags implements Operation.Flags for DrawImageOp.
func (op *DrawImageOp) Flags() AttributeFlags {
	return GetAttributeFlags(op.Paint) | AttrFlagHasImage
}

// DrawImageWithSamplingOp represents a draw image with sampling operation.
type DrawImageWithSamplingOp struct {
	Image    Image
	Point    geom.Point[Scalar]
	Sampling SamplingOptions
	Paint    Paint
}

// Invoke implements Operation.Invoke for DrawImageWithSamplingOp.
func (op *DrawImageWithSamplingOp) Invoke(receiver OpReceiver) {
	receiver.DrawImageWithSampling(op.Image, op.Point, op.Sampling, op.Paint)
}

// Bounds implements Operation.Bounds for DrawImageWithSamplingOp.
func (op *DrawImageWithSamplingOp) Bounds() *geom.Rect[Scalar] {
	if op.Image != nil {
		bounds := geom.NewRect(
			op.Point.X,
			op.Point.Y,
			op.Point.X+Scalar(op.Image.Width()),
			op.Point.Y+Scalar(op.Image.Height()),
		)
		return &bounds
	}
	return nil
}

// Flags implements Operation.Flags for DrawImageWithSamplingOp.
func (op *DrawImageWithSamplingOp) Flags() AttributeFlags {
	return GetAttributeFlags(op.Paint) | AttrFlagHasImage
}

// DrawParagraphOp represents a draw paragraph operation.
type DrawParagraphOp struct {
	Paragraph *Paragraph
	Point     geom.Point[Scalar]
}

// Invoke implements Operation.Invoke for DrawParagraphOp.
func (op *DrawParagraphOp) Invoke(receiver OpReceiver) {
	receiver.DrawParagraph(op.Paragraph, op.Point)
}

// Bounds implements Operation.Bounds for DrawParagraphOp.
func (op *DrawParagraphOp) Bounds() *geom.Rect[Scalar] {
	if op.Paragraph != nil {
		bounds := geom.NewRect(
			op.Point.X,
			op.Point.Y,
			op.Point.X+Scalar(op.Paragraph.Width()),
			op.Point.Y+Scalar(op.Paragraph.Height()),
		)
		return &bounds
	}
	return nil
}

// Flags implements Operation.Flags for DrawParagraphOp.
func (op *DrawParagraphOp) Flags() AttributeFlags {
	return AttrFlagHasText
}

// DrawShadowOp represents a draw shadow operation.
type DrawShadowOp struct {
	Path             geom.PathSource[Scalar]
	Color            Color
	Elevation        Scalar
	Transparent      bool
	DevicePixelRatio Scalar
}

// Invoke implements Operation.Invoke for DrawShadowOp.
func (op *DrawShadowOp) Invoke(receiver OpReceiver) {
	receiver.DrawShadow(op.Path, op.Color, op.Elevation, op.Transparent, op.DevicePixelRatio)
}

// Bounds implements Operation.Bounds for DrawShadowOp.
func (op *DrawShadowOp) Bounds() *geom.Rect[Scalar] {
	bounds := op.Path.Bounds()
	// Expand bounds for shadow offset and blur
	shadowOffset := op.Elevation * 0.5 // Simplified shadow calculation
	expandedBounds := geom.NewRect(
		bounds.Left-shadowOffset,
		bounds.Top-shadowOffset,
		bounds.Right+shadowOffset,
		bounds.Bottom+shadowOffset,
	)
	return &expandedBounds
}

// Flags implements Operation.Flags for DrawShadowOp.
func (op *DrawShadowOp) Flags() AttributeFlags {
	return AttrFlagNone
}
