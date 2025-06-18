package dl

import "github.com/opensraph/sraph/geom"

// OpReceiver defines the interface for receiving display list operations.
// This is the core interface that defines all the drawing operations
// that can be stored in a display list.
type OpReceiver interface {
	// Save/Restore operations
	Save()
	Restore()
	SaveLayer(bounds *geom.Rect[Scalar], paint *Paint)

	// Transform operations
	Translate(dx, dy Scalar)
	Scale(sx, sy Scalar)
	Rotate(radians geom.Radians)
	Skew(sx, sy Scalar)
	Transform2DAffine(mxx, mxy, myx, myy, mxt, myt Scalar)
	TransformFullPerspective(matrix geom.Matrix[Scalar])

	// Clip operations
	ClipRect(rect geom.Rect[Scalar], clipOp ClipOp, isAntiAliased bool)
	ClipRRect(rrect geom.RoundRect[Scalar], clipOp ClipOp, isAntiAliased bool)
	ClipPath(path geom.PathSource[Scalar], clipOp ClipOp, isAntiAliased bool)

	// Drawing operations - Basic shapes
	DrawPaint(paint Paint)
	DrawColor(color Color, blendMode BlendMode)
	DrawLine(p0, p1 geom.Point[Scalar], paint Paint)
	DrawRect(rect geom.Rect[Scalar], paint Paint)
	DrawOval(bounds geom.Rect[Scalar], paint Paint)
	DrawCircle(center geom.Point[Scalar], radius Scalar, paint Paint)
	DrawRRect(rrect geom.RoundRect[Scalar], paint Paint)
	DrawDRRect(outer, inner geom.RoundRect[Scalar], paint Paint)
	DrawPath(path geom.PathSource[Scalar], paint Paint)

	// Drawing operations - Advanced
	DrawArc(bounds geom.Rect[Scalar], start, sweep Scalar, useCenter bool, paint Paint)
	DrawPoints(mode PointMode, points []geom.Point[Scalar], paint Paint)
	DrawVertices(vertices *Vertices, blendMode BlendMode, paint Paint)

	// Text operations (placeholders for future implementation)
	// DrawTextBlob(blob *TextBlob, x, y Scalar, paint Paint)
	// DrawTextFrame(frame *TextFrame, x, y Scalar, paint Paint)

	// Image operations (placeholders for future implementation)
	// DrawImage(image *Image, x, y Scalar, sampling SamplingOptions, paint *Paint)
	// DrawImageRect(image *Image, src, dst geom.Rect[Scalar], sampling SamplingOptions, paint *Paint, constraint SrcRectConstraint)
	// DrawImageNine(image *Image, center geom.Rect[Scalar], dst geom.Rect[Scalar], filter FilterMode, paint *Paint)

	// Shadow operations (placeholders for future implementation)
	// DrawShadow(path geom.PathSource[Scalar], color Color, elevation Scalar, transparent bool, dpr Scalar)
}

// OpRecorder is an interface for objects that can record display list operations.
type OpRecorder interface {
	OpReceiver

	// Build finalizes the recording and returns a display list.
	Build() *DisplayList

	// GetBounds returns the current bounds of the recorded operations.
	GetBounds() geom.Rect[Scalar]

	// Reset resets the recorder to an empty state.
	Reset()
}
