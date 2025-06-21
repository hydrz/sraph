package geom

import "math"

// Superellipse represents a superellipse shape, extending RoundRect.
// The zero value is not valid.
type Superellipse struct {
	RoundRect
	param SuperellipseParam
}

// NewSuperellipse creates a superellipse with the given rectangle and corner radii.
func NewSuperellipse(rect Rect, radii RoundingRadii) Superellipse {
	param := NewSuperellipseParam(rect, radii)
	return Superellipse{
		RoundRect: NewRoundRect(rect, radii),
		param:     param,
	}
}

// NewSuperellipseOval creates a superellipse in oval form, with radii equal to half the width and height.
func NewSuperellipseOval(rect Rect) Superellipse {
	size := rect.Size()
	return NewSuperellipse(
		rect,
		NewRoundingRadiiFromSizes(size.Scale(0.5)),
	)
}

// NewSuperellipseRadius creates a superellipse with uniform corner radius.
func NewSuperellipseRadius(rect Rect, radius Scalar) Superellipse {
	return NewSuperellipse(
		rect,
		NewRoundingRadii(radius),
	)
}

// NewSuperellipseXY creates a superellipse with separate x and y radii.
func NewSuperellipseXY(rect Rect, xRadius, yRadius Scalar) Superellipse {
	return NewSuperellipse(
		rect,
		NewRoundingRadiiFromSizes(NewSize(xRadius, yRadius)),
	)
}

// NewSuperellipseLTRB creates a superellipse with separate left, top, right, and bottom radii.
func NewSuperellipseLTRB(rect Rect, left, top, right, bottom Scalar) Superellipse {
	return NewSuperellipse(
		rect,
		NewRoundingRadiiLTRB(left, top, right, bottom),
	)
}

// ToApproximateRoundRect returns a rounded rectangle that approximates the superellipse.
// Useful for backends that do not support superellipses.
func (s Superellipse) ToApproximateRoundRect() RoundRect {
	return NewRoundRect(s.RoundRect.Bounds(), s.RoundRect.Radius())
}

// Dispatch emits the path for the superellipse to the given receiver.
// If includeEnd is true, a PathEnd is emitted at the end.
func (s Superellipse) Dispatch(receiver PathReceiver, includeEnd bool) {
	builder := &superellipseBuilder{receiver: receiver}

	var start Point
	if s.param.IsUniform {
		// All corners are the same, so use the same quadrant with different signs.
		start = s.param.TopRight.Offset.Add(
			s.param.TopRight.SignedScale.Mul(NewPoint(0, s.param.TopRight.Top.SemiAxis)),
		)
		receiver.MoveTo(start, true)
		builder.AddQuadrant(s.param.TopRight, false, NewPoint(1, 1))
		builder.AddQuadrant(s.param.TopRight, true, NewPoint(1, -1))
		builder.AddQuadrant(s.param.TopRight, false, NewPoint(-1, -1))
		builder.AddQuadrant(s.param.TopRight, true, NewPoint(-1, 1))
	} else {
		start = s.param.TopRight.Offset.Add(
			s.param.TopRight.SignedScale.Mul(NewPoint(0, s.param.TopRight.Top.SemiAxis)),
		)
		receiver.MoveTo(start, true)
		builder.AddQuadrant(s.param.TopRight, false, NewPoint(1, 1))
		builder.AddQuadrant(s.param.BottomRight, true, NewPoint(1, -1))
		builder.AddQuadrant(s.param.BottomLeft, false, NewPoint(-1, -1))
		builder.AddQuadrant(s.param.TopLeft, true, NewPoint(-1, 1))
	}

	receiver.LineTo(start)
	receiver.Close()

	if includeEnd {
		receiver.PathEnd()
	}
}

// superellipseBuilder assists in building superellipse paths.
type superellipseBuilder struct {
	receiver PathReceiver
}

// AddQuadrant adds a quadrant of the superellipse to the path.
// reverse determines drawing direction; scaleSign flips/scales the quadrant.
func (b superellipseBuilder) AddQuadrant(quadrant SuperellipseQuadrant, reverse bool, scaleSign Point) {
	transform := NewMatrix().Scale(quadrant.SignedScale.Mul(scaleSign)).Translate(quadrant.Offset)
	// If either octant is degenerate (degree < 2), fallback to straight lines.
	if quadrant.Top.Degree < 2 || quadrant.Right.Degree < 2 {
		b.receiver.LineTo(
			quadrant.Top.Offset.Add(NewPoint(quadrant.Top.SemiAxis, quadrant.Top.SemiAxis)).Transform(transform),
		)
		if !reverse {
			b.receiver.LineTo(
				quadrant.Top.Offset.Add(NewPoint(quadrant.Top.SemiAxis, 0)).Transform(transform),
			)
		} else {
			b.receiver.LineTo(
				quadrant.Top.Offset.Add(NewPoint(0, quadrant.Top.SemiAxis)).Transform(transform),
			)
		}
		return
	}
	if !reverse {
		b.AddOctant(quadrant.Top, false, false, transform)
		b.AddOctant(quadrant.Right, true, true, transform)
	} else {
		b.AddOctant(quadrant.Right, false, true, transform)
		b.AddOctant(quadrant.Top, true, false, transform)
	}
}

// AddOctant adds an octant curve to the path.
// reverse determines curve direction; flip swaps x/y axes; mat is applied.
func (b superellipseBuilder) AddOctant(octant SuperellipseOctant, reverse, flip bool, mat Matrix) {
	mat = mat.Mul(
		NewMatrix().Translate(octant.Offset),
	)

	if flip {
		// Flip the octant by swapping x and y axes.
		mat = mat.Mul(Matrix{
			0, 1, 0, 0,
			1, 0, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		})
	}

	circlePoints := b.circularArcPoints(octant)
	sePoints := b.arcPoints(octant)

	if !reverse {
		b.receiver.CubicTo(
			sePoints[1].Transform(mat),
			sePoints[2].Transform(mat),
			sePoints[3].Transform(mat),
		)
		b.receiver.CubicTo(
			circlePoints[1].Transform(mat),
			circlePoints[2].Transform(mat),
			circlePoints[3].Transform(mat),
		)
	} else {
		b.receiver.CubicTo(
			circlePoints[2].Transform(mat),
			circlePoints[1].Transform(mat),
			circlePoints[0].Transform(mat),
		)
		b.receiver.CubicTo(
			sePoints[2].Transform(mat),
			sePoints[1].Transform(mat),
			sePoints[0].Transform(mat),
		)
	}

}

// circularArcPoints returns the four control points for the circular arc segment of the octant.
func (b superellipseBuilder) circularArcPoints(octant SuperellipseOctant) Quad {
	startVector := octant.CircleStart.Sub(octant.CircleCenter)
	endVector := startVector.Rotate(Radians(-octant.CircleMaxAngle))
	circleEnd := octant.CircleCenter.Add(endVector)
	startTangent := NewPoint(startVector.Y(), -startVector.X()).Normalize()
	endTangent := NewPoint(-endVector.Y(), endVector.X()).Normalize()
	bezierFactor := Scalar(math.Tan(ToFloat64(octant.CircleMaxAngle) / 4 * 4 / 3))
	radius := startVector.Length()

	return Quad{
		octant.CircleStart,
		octant.CircleStart.Add(startTangent.Scale(bezierFactor * radius)),
		circleEnd.Add(endTangent.Scale(bezierFactor * radius)),
		circleEnd,
	}
}

// arcPoints returns the four control points for the superellipse arc segment of the octant.
func (b superellipseBuilder) arcPoints(octant SuperellipseOctant) Quad {
	start := NewPoint(0, octant.SemiAxis)
	end := octant.CircleStart
	startTangent := NewPoint(1, 0)
	circleStartVector := octant.CircleStart.Sub(octant.CircleCenter)
	endTangent := NewPoint(-circleStartVector.Y(), circleStartVector.X()).Normalize()
	factors := b.bezierFactors(octant.SemiAxis)
	return Quad{
		start,
		start.Add(startTangent.Scale(factors[0] * octant.SemiAxis)),
		end.Add(endTangent.Scale(factors[1] * octant.SemiAxis)),
		end,
	}
}

// bezierFactors returns the Bezier control factors for a given superellipse degree.
// Uses precomputed values and interpolation for best accuracy.
func (b superellipseBuilder) bezierFactors(n Scalar) [2]Scalar {
	precomputedVariables := [][2]float32{
		{0.01339448, 0.05994973}, // n=2.0
		{0.13664115, 0.13592082}, // n=3.0
		{0.24545546, 0.14099516}, // n=4.0
		{0.32353151, 0.12808021}, // n=5.0
		{0.39093068, 0.11726264}, // n=6.0
		{0.44847800, 0.10808278}, // n=7.0
		{0.49817452, 0.10026175}, // n=8.0
		{0.54105583, 0.09344429}, // n=9.0
		{0.57812578, 0.08748984}, // n=10.0
		{0.61050961, 0.08224722}, // n=11.0
		{0.63903989, 0.07759639}, // n=12.0
		{0.66416338, 0.07346530}, // n=13.0
		{0.68675338, 0.06974996}, // n=14.0
		{0.70678034, 0.06529512}, // n=15.0
	}
	numRecords := Scalar(len(precomputedVariables))
	step := Scalar(1)
	minN := Scalar(2)
	maxN := minN + (numRecords-1)*step

	if n >= maxN {
		// Use heuristic formula for large n.
		return [2]Scalar{
			Scalar(1.07 - math.Exp(1.307649835)*math.Pow(ToFloat64(n), -0.8568516731)),
			Scalar(-0.01 + math.Exp(-0.9287690322)*math.Pow(ToFloat64(n), -0.6120901398)),
		}
	}

	steps := Clamp((n-minN)/step, 0, numRecords-1)
	left := int(Clamp(math.Floor(ToFloat64(steps)), 0, ToFloat64(numRecords-2)))
	frac := float32(steps - Scalar(left))
	return [2]Scalar{
		Scalar((1-frac)*precomputedVariables[left][0] + frac*precomputedVariables[left+1][0]),
		Scalar((1-frac)*precomputedVariables[left][1] + frac*precomputedVariables[left+1][1]),
	}
}

// SuperellipsePathSource implements PathSource for a single superellipse.
type SuperellipsePathSource struct {
	superellipse Superellipse
}

// NewSuperellipsePathSource returns a PathSource for a superellipse.
func NewSuperellipsePathSource(superellipse Superellipse) PathSource {
	return &SuperellipsePathSource{superellipse: superellipse}
}

// FillType returns the fill rule for the superellipse path.
func (s SuperellipsePathSource) FillType() FillType {
	return FillTypeNonZero
}

// Bounds returns the bounding rectangle of the superellipse.
func (s SuperellipsePathSource) Bounds() Rect {
	return s.superellipse.Bounds()
}

// IsConvex reports whether the superellipse is convex.
func (s SuperellipsePathSource) IsConvex() bool {
	return true
}

// Dispatch emits the path for the superellipse to the receiver.
func (s SuperellipsePathSource) Dispatch(receiver PathReceiver) {
	s.superellipse.Dispatch(receiver, true)
}

// DiffSuperellipsePathSource implements PathSource for the difference of two superellipses.
type DiffSuperellipsePathSource struct {
	outter Superellipse
	inner  Superellipse
}

// NewDiffSuperellipsePathSource returns a PathSource for the difference of two superellipses.
func NewDiffSuperellipsePathSource(outter, inner Superellipse) PathSource {
	return &DiffSuperellipsePathSource{outter: outter, inner: inner}
}

// FillType returns the fill rule for the difference path.
func (d *DiffSuperellipsePathSource) FillType() FillType {
	return FillTypeEvenOdd
}

// IsConvex reports whether the difference of two superellipses is convex.
func (d *DiffSuperellipsePathSource) IsConvex() bool {
	return false
}

// Bounds returns the bounding rectangle of the outer superellipse.
func (d *DiffSuperellipsePathSource) Bounds() Rect {
	return d.outter.Bounds()
}

// Dispatch emits the path for the difference of two superellipses to the receiver.
func (d *DiffSuperellipsePathSource) Dispatch(receiver PathReceiver) {
	d.outter.Dispatch(receiver, false)
	d.inner.Dispatch(receiver, true)
}

// SuperellipseOctant holds parameters for drawing a square-like rounded superellipse octant.
//
// A Degree of 0 means the radius is 0 and this octant is a square of size SemiAxis at Offset.
type SuperellipseOctant struct {
	Offset         Point   // Center of the octant, relative to origin.
	SemiAxis       Scalar  // Semi-axis length.
	Degree         Scalar  // Degree of the superellipse. 0 means square.
	MaxTheta       Scalar  // Range of the parametric "theta".
	CircleStart    Point   // Start point of the circular arc, relative to Offset.
	CircleCenter   Point   // Center of the circular arc, relative to Offset.
	CircleMaxAngle Radians // Angular span of the circular arc, in radians.
}

// SuperellipseQuadrant holds parameters for a quadrant of a rounded superellipse.
//
// Used to define a quadrant of an arbitrary rounded superellipse.
type SuperellipseQuadrant struct {
	Offset      Point // Center of the quadrant, relative to origin.
	SignedScale Point // Scaling factor for normalization and flipping.
	Top         SuperellipseOctant
	Right       SuperellipseOctant
}

// SuperellipseParam expands input parameters for a rounded superellipse to drawing variables.
type SuperellipseParam struct {
	TopRight    SuperellipseQuadrant
	BottomRight SuperellipseQuadrant
	BottomLeft  SuperellipseQuadrant
	TopLeft     SuperellipseQuadrant
	IsUniform   bool // If true, only TopRight is populated.
}

// NewSuperellipseParam computes drawing parameters for a superellipse with given bounds and radii.
func NewSuperellipseParam(bounds Rect, radii RoundingRadii) SuperellipseParam {
	if radii.IsUniform() && !radii.TopRight().IsZero() {
		sq := newSuperellipseQuadrant(bounds.Center(), bounds.TopRight(), radii.TopRight(), NewSize(-1, 1))
		return SuperellipseParam{
			TopRight:  sq,
			IsUniform: radii.IsUniform(),
		}
	}

	split := func(left, right, ratioLeft, ratioRight Scalar) Scalar {
		if ratioLeft == 0 && ratioRight == 0 {
			return (left + right) / 2
		}
		return (left*ratioLeft + right*ratioLeft) / (ratioLeft + ratioRight)
	}

	topSplit := split(bounds.Left(), bounds.Right(), radii.TopLeft().Width(), radii.TopRight().Width())
	rightSplit := split(bounds.Top(), bounds.Bottom(), radii.TopRight().Height(), radii.BottomRight().Height())
	bottomSplit := split(bounds.Left(), bounds.Right(), radii.BottomLeft().Width(), radii.BottomRight().Width())
	leftSplit := split(bounds.Top(), bounds.Bottom(), radii.TopLeft().Height(), radii.BottomLeft().Height())

	return SuperellipseParam{
		TopRight: newSuperellipseQuadrant(
			NewPoint(topSplit, rightSplit), bounds.TopRight(), radii.TopRight(), NewSize(1, -1),
		),
		BottomRight: newSuperellipseQuadrant(
			NewPoint(bottomSplit, rightSplit), bounds.BottomRight(), radii.BottomRight(), NewSize(1, 1),
		),
		BottomLeft: newSuperellipseQuadrant(
			NewPoint(bottomSplit, leftSplit), bounds.BottomLeft(), radii.BottomLeft(), NewSize(-1, 1),
		),
		TopLeft: newSuperellipseQuadrant(
			NewPoint(topSplit, leftSplit), bounds.TopLeft(), radii.TopLeft(), NewSize(-1, -1),
		),
		IsUniform: false,
	}
}

// newSuperellipseQuadrant computes parameters for a quadrant of a rounded superellipse with asymmetrical radii.
// center is the quadrant center, corner is the corner point, sign indicates the quadrant orientation.
func newSuperellipseQuadrant(center Point, corner Point, in_radii Size, sign Size) SuperellipseQuadrant {

	centerVector := corner.Sub(center)
	radii := in_radii.Abs()
	normRadius := radii.MinDimension()
	var forwardScale Size
	if normRadius == 0 {
		forwardScale = NewSize(1, 1)
	} else {
		forwardScale = radii.Scale(1 / normRadius)
	}

	normHalfSize := centerVector.Abs().DivSize(forwardScale)
	signedScale := replaceNaNWithDefault(centerVector.Div(normHalfSize), sign)

	c := normHalfSize.X() - normHalfSize.Y()
	return SuperellipseQuadrant{
		Offset:      center,
		SignedScale: signedScale,
		Top:         newSuperellipseOctant(NewPoint(0, -c), normHalfSize.X(), normRadius),
		Right:       newSuperellipseOctant(NewPoint(c, 0), normHalfSize.Y(), normRadius),
	}

}

// newSuperellipseOctant computes parameters for a square-like rounded superellipse with a symmetrical radius.
// center is the octant center, a is the semi-axis, radius is the corner radius.
func newSuperellipseOctant(center Point, a Scalar, radius Scalar) SuperellipseOctant {
	const gapFactor = 0.29289321881 // 1-cos(pi/4)

	if radius <= 0 {
		return SuperellipseOctant{
			Offset:         center,
			SemiAxis:       a,
			Degree:         0,
			CircleStart:    NewPoint(a, a),
			CircleCenter:   NewPoint(0, 0),
			CircleMaxAngle: Radians(0),
		}
	}

	af := ToFloat64(a)

	ratio := a * 2 / radius
	g := Scalar(gapFactor * ToFloat64(radius))
	precomputedVars := superellipseComputeNAndXj(ToFloat64(ratio))
	n := precomputedVars[0]
	xJ := precomputedVars[1] * af
	yJ := math.Pow(1-math.Pow(precomputedVars[1], n), 1/n) * af
	maxTheta := math.Asin(math.Pow(precomputedVars[1], n/2))
	tanPhiJ := math.Pow(xJ/yJ, n-1)
	d := (xJ - tanPhiJ*yJ) / (1 - tanPhiJ)
	R := Scalar((ToFloat64(a) - d - ToFloat64(g)) * math.Sqrt(2))

	PointM := NewPoint(a-g, a-g)
	pointJ := NewPoint(Scalar(xJ), Scalar(yJ))
	var circleCenter Point
	if radius == 0 {
		circleCenter = PointM
	} else {
		circleCenter = findCircleCenter(pointJ, PointM, R)
	}
	var circleMaxAngle Radians
	if radius == 0 {
		circleMaxAngle = Radians(0)
	} else {
		circleMaxAngle = (PointM.Sub(circleCenter)).AngleTo(pointJ.Sub(circleCenter))
	}

	return SuperellipseOctant{
		Offset:         center,
		SemiAxis:       a,
		Degree:         Scalar(n),
		MaxTheta:       Scalar(maxTheta),
		CircleStart:    pointJ,
		CircleCenter:   circleCenter,
		CircleMaxAngle: circleMaxAngle,
	}
}

// superellipseComputeNAndXj computes the degree n and normalized x_j for a given ratio.
func superellipseComputeNAndXj(ratio float64) [2]float64 {
	const (
		minRatio          = 2.00
		firstStepInverse  = 10
		firstMaxRatio     = 2.50
		firstNumRecords   = 6
		secondStepInverse = 2
		secondMaxRatio    = 5.00
		thirdNSlope       = 1.559599389
		thirdKxjSlope     = 0.522807185
	)

	var (
		precomputedVariables = [][2]float64{
			{2.00000000, 1.13276676}, // ratio=2.00
			{2.18349805, 1.20311921}, // ratio=2.10
			{2.33888662, 1.28698796}, // ratio=2.20
			{2.48660575, 1.36351941}, // ratio=2.30
			{2.62226596, 1.44717976}, // ratio=2.40
			{2.75148990, 1.53385819}, // ratio=2.50
			{3.36298265, 1.98288283}, // ratio=3.00
			{4.08649929, 2.23811846}, // ratio=3.50
			{4.85481134, 2.47563463}, // ratio=4.00
			{5.62945551, 2.72948597}, // ratio=4.50
			{6.43023796, 2.98020421}, // ratio=5.00
		}
		kNumRecords = len(precomputedVariables)
	)

	if ratio > secondMaxRatio {
		n := thirdNSlope*(ratio-secondMaxRatio) +
			precomputedVariables[kNumRecords-1][0]
		k_xJ := thirdKxjSlope*(ratio-secondMaxRatio) +
			precomputedVariables[kNumRecords-1][1]
		return [2]float64{n, 1 - 1/k_xJ}
	}
	ratio = Clamp(ratio, minRatio, secondMaxRatio)
	var steps float64
	if ratio < firstMaxRatio {
		steps = (ratio - minRatio) * firstStepInverse
	} else {
		steps =
			(ratio-firstMaxRatio)*secondStepInverse + firstNumRecords - 1
	}

	left := Clamp(math.Floor(steps), 0, float64(kNumRecords-2))
	frac := steps - left
	n := (1-frac)*precomputedVariables[int(left)][0] +
		frac*precomputedVariables[int(left)+1][0]
	k_xJ := (1-frac)*precomputedVariables[int(left)][1] +
		frac*precomputedVariables[int(left)+1][1]
	return [2]float64{n, 1 - 1/k_xJ}
}

// findCircleCenter returns the center of the circle passing through points a and b with radius r.
func findCircleCenter(a Point, b Point, r Scalar) Point {
	aToB := b.Sub(a)
	m := a.Add(b).Scale(0.5)
	cToM := NewPoint(-aToB.Y(), aToB.X())
	distanceAM := aToB.Length() / 2
	distanceCM := Scalar(math.Sqrt(ToFloat64(r*r - distanceAM*distanceAM)))
	return m.Sub(cToM.Normalize().Scale(distanceCM))
}

// replaceNaNWithDefault replaces NaN values in v with the corresponding values from defaultValue.
func replaceNaNWithDefault(v Point, defaultValue Size) Point {
	x := defaultValue.Width()
	y := defaultValue.Height()

	if !math.IsNaN(ToFloat64(v.X())) {
		x = v.X()
	}

	if !math.IsNaN(ToFloat64(v.Y())) {
		y = v.Y()
	}

	return NewPoint(x, y)
}
