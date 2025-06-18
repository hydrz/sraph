package tess

import (
	"github.com/opensraph/sraph/geom"
)

// Generator functions for different shapes

// generateFilledCircle generates vertices for a filled circle.
func generateFilledCircle[T geom.Scalar](trigs *Trigs[T], data EllipticalVertexGeneratorData[T], proc TessellatedVertexProc[T]) {
	center := data.ReferenceCenters[0]
	radius := data.Radii.Width

	// Generate center point first for triangle fan
	proc(center)

	// Generate circle perimeter points
	for i := 0; i < trigs.Size(); i++ {
		trig := trigs.Get(i)
		point := trig.CirclePoint(radius).Add(center)
		proc(point)
	}

	// Close the circle by repeating the first perimeter point
	if trigs.Size() > 0 {
		trig := trigs.Get(0)
		point := trig.CirclePoint(radius).Add(center)
		proc(point)
	}
}

// generateStrokedCircle generates vertices for a stroked circle.
func generateStrokedCircle[T geom.Scalar](trigs *Trigs[T], data EllipticalVertexGeneratorData[T], proc TessellatedVertexProc[T]) {
	center := data.ReferenceCenters[0]
	radius := data.Radii.Width
	halfWidth := data.HalfWidth

	innerRadius := radius - halfWidth
	outerRadius := radius + halfWidth

	// Generate alternating inner and outer points for triangle strip
	for i := 0; i < trigs.Size(); i++ {
		trig := trigs.Get(i)

		innerPoint := trig.CirclePoint(innerRadius).Add(center)
		outerPoint := trig.CirclePoint(outerRadius).Add(center)

		proc(innerPoint)
		proc(outerPoint)
	}

	// Close the stroke by repeating the first points
	if trigs.Size() > 0 {
		trig := trigs.Get(0)
		innerPoint := trig.CirclePoint(innerRadius).Add(center)
		outerPoint := trig.CirclePoint(outerRadius).Add(center)

		proc(innerPoint)
		proc(outerPoint)
	}
}

// generateFilledEllipse generates vertices for a filled ellipse.
func generateFilledEllipse[T geom.Scalar](trigs *Trigs[T], data EllipticalVertexGeneratorData[T], proc TessellatedVertexProc[T]) {
	center := data.ReferenceCenters[0]
	radii := data.Radii

	// Generate center point first for triangle fan
	proc(center)

	// Generate ellipse perimeter points
	for i := 0; i < trigs.Size(); i++ {
		trig := trigs.Get(i)
		point := trig.EllipsePoint(radii).Add(center)
		proc(point)
	}

	// Close the ellipse by repeating the first perimeter point
	if trigs.Size() > 0 {
		trig := trigs.Get(0)
		point := trig.EllipsePoint(radii).Add(center)
		proc(point)
	}
}

// generateRoundCapLine generates vertices for a line with round caps.
func generateRoundCapLine[T geom.Scalar](trigs *Trigs[T], data EllipticalVertexGeneratorData[T], proc TessellatedVertexProc[T]) {
	p0 := data.ReferenceCenters[0]
	p1 := data.ReferenceCenters[1]
	radius := data.Radii.Width

	// Calculate line direction and perpendicular
	lineVec := p1.Sub(p0)
	lineLength := lineVec.Length()

	if lineLength == 0 {
		// Degenerate line, just draw a circle at p0
		generateFilledCircle(trigs, EllipticalVertexGeneratorData[T]{
			ReferenceCenters: [2]geom.Point[T]{p0, p0},
			Radii:            data.Radii,
			HalfWidth:        0,
		}, proc)
		return
	}

	lineDir := lineVec.Scale(1.0 / lineLength)
	perpDir := geom.Point[T]{X: -lineDir.Y, Y: lineDir.X}

	// Generate start cap (half circle)
	for i := trigs.Size() / 2; i < trigs.Size(); i++ {
		trig := trigs.Get(i)
		// Rotate the trig to align with line direction
		point := geom.Point[T]{
			X: trig.Cos*perpDir.X - trig.Sin*lineDir.X,
			Y: trig.Cos*perpDir.Y - trig.Sin*lineDir.Y,
		}.Scale(radius).Add(p0)
		proc(point)
	}

	// Generate line body
	proc(p0.Add(perpDir.Scale(radius)))
	proc(p1.Add(perpDir.Scale(radius)))
	proc(p0.Sub(perpDir.Scale(radius)))
	proc(p1.Sub(perpDir.Scale(radius)))

	// Generate end cap (half circle)
	for i := 0; i <= trigs.Size()/2; i++ {
		trig := trigs.Get(i)
		// Rotate the trig to align with line direction
		point := geom.Point[T]{
			X: trig.Cos*perpDir.X + trig.Sin*lineDir.X,
			Y: trig.Cos*perpDir.Y + trig.Sin*lineDir.Y,
		}.Scale(radius).Add(p1)
		proc(point)
	}
}

// generateFilledArcFan generates vertices for a filled arc using triangle fan.
func generateFilledArcFan[T geom.Scalar](trigs *Trigs[T], iteration *ArcIteration[T], ovalBounds geom.Rect[T], useCenter bool, proc TessellatedVertexProc[T]) {
	center := ovalBounds.Center()
	radii := geom.Size[T]{
		Width:  ovalBounds.Width() / 2,
		Height: ovalBounds.Height() / 2,
	}

	if useCenter {
		proc(center)
	}

	// Generate start point
	startPoint := geom.Point[T]{
		X: iteration.Start.X * radii.Width,
		Y: iteration.Start.Y * radii.Height,
	}.Add(center)
	proc(startPoint)

	// Generate quadrant points
	for _, quadrant := range iteration.Quadrants {
		for i := quadrant.StartIndex; i < quadrant.EndIndex; i++ {
			trig := trigs.Get(i)
			rotatedTrig := geom.Point[T]{
				X: trig.Cos*quadrant.Axis.X - trig.Sin*quadrant.Axis.Y,
				Y: trig.Cos*quadrant.Axis.Y + trig.Sin*quadrant.Axis.X,
			}
			point := geom.Point[T]{
				X: rotatedTrig.X * radii.Width,
				Y: rotatedTrig.Y * radii.Height,
			}.Add(center)
			proc(point)
		}
	}

	// Generate end point
	endPoint := geom.Point[T]{
		X: iteration.End.X * radii.Width,
		Y: iteration.End.Y * radii.Height,
	}.Add(center)
	proc(endPoint)
}

// generateFilledArcStrip generates vertices for a filled arc using triangle strip.
func generateFilledArcStrip[T geom.Scalar](trigs *Trigs[T], iteration *ArcIteration[T], ovalBounds geom.Rect[T], useCenter bool, proc TessellatedVertexProc[T]) {
	// For now, use the same implementation as fan but without center
	generateFilledArcFan(trigs, iteration, ovalBounds, false, proc)
}

// generateStrokedArc generates vertices for a stroked arc.
func generateStrokedArc[T geom.Scalar](trigs *Trigs[T], iteration *ArcIteration[T], ovalBounds geom.Rect[T], halfWidth T, cap geom.StrokeCap, proc TessellatedVertexProc[T]) {
	center := ovalBounds.Center()
	radii := geom.Size[T]{
		Width:  ovalBounds.Width() / 2,
		Height: ovalBounds.Height() / 2,
	}

	innerRadii := geom.Size[T]{
		Width:  radii.Width - halfWidth,
		Height: radii.Height - halfWidth,
	}
	outerRadii := geom.Size[T]{
		Width:  radii.Width + halfWidth,
		Height: radii.Height + halfWidth,
	}

	// Generate start points
	startInner := geom.Point[T]{
		X: iteration.Start.X * innerRadii.Width,
		Y: iteration.Start.Y * innerRadii.Height,
	}.Add(center)
	startOuter := geom.Point[T]{
		X: iteration.Start.X * outerRadii.Width,
		Y: iteration.Start.Y * outerRadii.Height,
	}.Add(center)

	proc(startInner)
	proc(startOuter)

	// Generate quadrant points
	for _, quadrant := range iteration.Quadrants {
		for i := quadrant.StartIndex; i < quadrant.EndIndex; i++ {
			trig := trigs.Get(i)
			rotatedTrig := geom.Point[T]{
				X: trig.Cos*quadrant.Axis.X - trig.Sin*quadrant.Axis.Y,
				Y: trig.Cos*quadrant.Axis.Y + trig.Sin*quadrant.Axis.X,
			}

			innerPoint := geom.Point[T]{
				X: rotatedTrig.X * innerRadii.Width,
				Y: rotatedTrig.Y * innerRadii.Height,
			}.Add(center)
			outerPoint := geom.Point[T]{
				X: rotatedTrig.X * outerRadii.Width,
				Y: rotatedTrig.Y * outerRadii.Height,
			}.Add(center)

			proc(innerPoint)
			proc(outerPoint)
		}
	}

	// Generate end points
	endInner := geom.Point[T]{
		X: iteration.End.X * innerRadii.Width,
		Y: iteration.End.Y * innerRadii.Height,
	}.Add(center)
	endOuter := geom.Point[T]{
		X: iteration.End.X * outerRadii.Width,
		Y: iteration.End.Y * outerRadii.Height,
	}.Add(center)

	proc(endInner)
	proc(endOuter)

	// TODO: Add cap generation based on cap type
}
