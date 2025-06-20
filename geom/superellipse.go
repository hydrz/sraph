package geom

import "math"

// Superellipse represents a superellipse shape, extending RoundRect.
// The zero value is not valid.
type Superellipse[T TScalar] struct {
	RoundRect[T]
	param SuperellipseParam[T]
}

// NewSuperellipse creates a superellipse with the given rectangle and corner radii.
func NewSuperellipse[T TScalar](rect Rect[T], radii RoundingRadii[T]) Superellipse[T] {
	param := NewSuperellipseParam(rect, radii)
	return Superellipse[T]{
		RoundRect: NewRoundRect(rect, radii),
		param:     param,
	}
}

// NewSuperellipseOval creates a superellipse in oval form, with radii equal to half the width and height.
func NewSuperellipseOval[T TScalar](rect Rect[T]) Superellipse[T] {
	size := rect.Size()
	halfSize := Size[T]{Width: size.Width / T(2), Height: size.Height / T(2)}
	return NewSuperellipse(
		rect,
		NewRoundingRadiiFromSizes(halfSize),
	)
}

// NewSuperellipseRadius creates a superellipse with uniform corner radius.
func NewSuperellipseRadius[T TScalar](rect Rect[T], radius T) Superellipse[T] {
	return NewSuperellipse(
		rect,
		NewRoundingRadii(radius),
	)
}

// NewSuperellipseXY creates a superellipse with separate x and y radii.
func NewSuperellipseXY[T TScalar](rect Rect[T], xRadius, yRadius T) Superellipse[T] {
	return NewSuperellipse(
		rect,
		NewRoundingRadiiFromSizes(Size[T]{Width: xRadius, Height: yRadius}),
	)
}

// NewSuperellipseLTRB creates a superellipse with separate left, top, right, and bottom radii.
func NewSuperellipseLTRB[T TScalar](rect Rect[T], left, top, right, bottom T) Superellipse[T] {
	return NewSuperellipse(
		rect,
		NewRoundingRadiiLTRB(left, top, right, bottom),
	)
}

// ToApproximateRoundRect returns a rounded rectangle that approximates the superellipse.
// Useful for backends that do not support superellipses.
func (s Superellipse[T]) ToApproximateRoundRect() RoundRect[T] {
	return NewRoundRect(s.RoundRect.Bounds(), s.RoundRect.Radius())
}

// Dispatch emits the path for the superellipse to the given receiver.
// If includeEnd is true, a PathEnd is emitted at the end.
func (s Superellipse[T]) Dispatch(receiver PathReceiver[T], includeEnd bool) {
	builder := &superellipseBuilder[T]{receiver: receiver}

	var start Point[T]
	if s.param.IsUniform {
		// All corners are the same, so use the same quadrant with different signs.
		start = s.param.TopRight.Offset.Add(
			s.param.TopRight.SignedScale.Mul(Point[T]{X: 0, Y: s.param.TopRight.Top.SemiAxis}),
		)
		receiver.MoveTo(start, true)
		builder.AddQuadrant(s.param.TopRight, false, Point[T]{X: 1, Y: 1})
		builder.AddQuadrant(s.param.TopRight, true, Point[T]{X: 1, Y: -1})
		builder.AddQuadrant(s.param.TopRight, false, Point[T]{X: -1, Y: -1})
		builder.AddQuadrant(s.param.TopRight, true, Point[T]{X: -1, Y: 1})
	} else {
		start = s.param.TopRight.Offset.Add(
			s.param.TopRight.SignedScale.Mul(Point[T]{X: 0, Y: s.param.TopRight.Top.SemiAxis}),
		)
		receiver.MoveTo(start, true)
		builder.AddQuadrant(s.param.TopRight, false, Point[T]{X: 1, Y: 1})
		builder.AddQuadrant(s.param.BottomRight, true, Point[T]{X: 1, Y: -1})
		builder.AddQuadrant(s.param.BottomLeft, false, Point[T]{X: -1, Y: -1})
		builder.AddQuadrant(s.param.TopLeft, true, Point[T]{X: -1, Y: 1})
	}

	receiver.LineTo(start)
	receiver.Close()

	if includeEnd {
		receiver.PathEnd()
	}
}

// superellipseBuilder assists in building superellipse paths.
type superellipseBuilder[T TScalar] struct {
	receiver PathReceiver[T]
}

// AddQuadrant adds a quadrant of the superellipse to the path.
// reverse determines drawing direction; scaleSign flips/scales the quadrant.
func (b superellipseBuilder[T]) AddQuadrant(quadrant SuperellipseQuadrant[T], reverse bool, scaleSign Point[T]) {
	transform := NewMatrix[T]().Scale(quadrant.SignedScale.Mul(scaleSign)).Translate(quadrant.Offset)
	// If either octant is degenerate (degree < 2), fallback to straight lines.
	if quadrant.Top.Degree < 2 || quadrant.Right.Degree < 2 {
		b.receiver.LineTo(
			transform.transformPoint(
				quadrant.Top.Offset.Add(Point[T]{X: quadrant.Top.SemiAxis, Y: quadrant.Top.SemiAxis}),
			),
		)
		if !reverse {
			b.receiver.LineTo(
				transform.transformPoint(
					quadrant.Top.Offset.Add(Point[T]{X: quadrant.Top.SemiAxis, Y: 0}),
				),
			)
		} else {
			b.receiver.LineTo(
				transform.transformPoint(
					quadrant.Top.Offset.Add(Point[T]{X: 0, Y: quadrant.Top.SemiAxis}),
				),
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
// reverse determines curve direction; flip swaps x/y axes; externalTransform is applied.
func (b superellipseBuilder[T]) AddOctant(octant SuperellipseOctant[T], reverse, flip bool, externalTransform Matrix[T]) {
	transform := externalTransform.Mul(
		NewMatrix[T]().Translate(octant.Offset),
	)

	if flip {
		// Flip the octant by swapping x and y axes.
		transform = transform.Mul(Matrix[T]{
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
			transform.transformPoint(sePoints[1]),
			transform.transformPoint(sePoints[2]),
			transform.transformPoint(sePoints[3]),
		)
		b.receiver.CubicTo(
			transform.transformPoint(circlePoints[1]),
			transform.transformPoint(circlePoints[2]),
			transform.transformPoint(circlePoints[3]),
		)
	} else {
		b.receiver.CubicTo(
			transform.transformPoint(circlePoints[2]),
			transform.transformPoint(circlePoints[1]),
			transform.transformPoint(circlePoints[0]),
		)
		b.receiver.CubicTo(
			transform.transformPoint(sePoints[2]),
			transform.transformPoint(sePoints[1]),
			transform.transformPoint(sePoints[0]),
		)
	}

}

// circularArcPoints returns the four control points for the circular arc segment of the octant.
func (b superellipseBuilder[T]) circularArcPoints(octant SuperellipseOctant[T]) [4]Point[T] {
	startVector := octant.CircleStart.Sub(octant.CircleCenter)
	endVector := startVector.Rotate(Radians(-octant.CircleMaxAngle))
	circleEnd := octant.CircleCenter.Add(endVector)
	startTangent := Point[T]{X: startVector.Y, Y: -startVector.X}.Normalize()
	endTangent := Point[T]{X: -endVector.Y, Y: endVector.X}.Normalize()
	bezierFactor := T(math.Tan(ToFloat64(octant.CircleMaxAngle) / 4 * 4 / 3))
	radius := startVector.Length()

	return [4]Point[T]{
		octant.CircleStart,
		octant.CircleStart.Add(startTangent.Scale(bezierFactor * radius)),
		circleEnd.Add(endTangent.Scale(bezierFactor * radius)),
		circleEnd,
	}
}

// arcPoints returns the four control points for the superellipse arc segment of the octant.
func (b superellipseBuilder[T]) arcPoints(octant SuperellipseOctant[T]) [4]Point[T] {
	start := Point[T]{X: 0, Y: octant.SemiAxis}
	end := octant.CircleStart
	startTangent := Point[T]{X: 1, Y: 0}
	circleStartVector := octant.CircleStart.Sub(octant.CircleCenter)
	endTangent := Point[T]{X: -circleStartVector.Y, Y: circleStartVector.X}.Normalize()
	factors := b.bezierFactors(octant.SemiAxis)
	return [4]Point[T]{
		start,
		start.Add(startTangent.Scale(factors[0] * octant.SemiAxis)),
		end.Add(endTangent.Scale(factors[1] * octant.SemiAxis)),
		end,
	}
}

// bezierFactors returns the Bezier control factors for a given superellipse degree.
// Uses precomputed values and interpolation for best accuracy.
func (b superellipseBuilder[T]) bezierFactors(n T) [2]T {
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
	numRecords := T(len(precomputedVariables))
	step := T(1)
	minN := T(2)
	maxN := minN + (numRecords-1)*step

	if n >= maxN {
		// Use heuristic formula for large n.
		return [2]T{
			T(1.07 - math.Exp(1.307649835)*math.Pow(ToFloat64(n), -0.8568516731)),
			T(-0.01 + math.Exp(-0.9287690322)*math.Pow(ToFloat64(n), -0.6120901398)),
		}
	}

	steps := Clamp((n-minN)/step, 0, numRecords-1)
	left := int(Clamp(math.Floor(ToFloat64(steps)), 0, ToFloat64(numRecords-2)))
	frac := float32(steps - T(left))
	return [2]T{
		T((1-frac)*precomputedVariables[left][0] + frac*precomputedVariables[left+1][0]),
		T((1-frac)*precomputedVariables[left][1] + frac*precomputedVariables[left+1][1]),
	}
}

// SuperellipsePathSource implements PathSource for a single superellipse.
type SuperellipsePathSource[T TScalar] struct {
	superellipse Superellipse[T]
}

// NewSuperellipsePathSource returns a PathSource for a superellipse.
func NewSuperellipsePathSource[T TScalar](superellipse Superellipse[T]) PathSource[T] {
	return &SuperellipsePathSource[T]{superellipse: superellipse}
}

// FillType returns the fill rule for the superellipse path.
func (s SuperellipsePathSource[T]) FillType() FillType {
	return FillTypeNonZero
}

// Bounds returns the bounding rectangle of the superellipse.
func (s SuperellipsePathSource[T]) Bounds() Rect[T] {
	return s.superellipse.Bounds()
}

// IsConvex reports whether the superellipse is convex.
func (s SuperellipsePathSource[T]) IsConvex() bool {
	return true
}

// Dispatch emits the path for the superellipse to the receiver.
func (s SuperellipsePathSource[T]) Dispatch(receiver PathReceiver[T]) {
	s.superellipse.Dispatch(receiver, true)
}

// DiffSuperellipsePathSource implements PathSource for the difference of two superellipses.
type DiffSuperellipsePathSource[T TScalar] struct {
	outter Superellipse[T]
	inner  Superellipse[T]
}

// NewDiffSuperellipsePathSource returns a PathSource for the difference of two superellipses.
func NewDiffSuperellipsePathSource[T TScalar](outter, inner Superellipse[T]) PathSource[T] {
	return &DiffSuperellipsePathSource[T]{outter: outter, inner: inner}
}

// FillType returns the fill rule for the difference path.
func (d *DiffSuperellipsePathSource[T]) FillType() FillType {
	return FillTypeEvenOdd
}

// IsConvex reports whether the difference of two superellipses is convex.
func (d *DiffSuperellipsePathSource[T]) IsConvex() bool {
	return false
}

// Bounds returns the bounding rectangle of the outer superellipse.
func (d *DiffSuperellipsePathSource[T]) Bounds() Rect[T] {
	return d.outter.Bounds()
}

// Dispatch emits the path for the difference of two superellipses to the receiver.
func (d *DiffSuperellipsePathSource[T]) Dispatch(receiver PathReceiver[T]) {
	d.outter.Dispatch(receiver, false)
	d.inner.Dispatch(receiver, true)
}

// SuperellipseOctant holds parameters for drawing a square-like rounded superellipse octant.
//
// A Degree of 0 means the radius is 0 and this octant is a square of size SemiAxis at Offset.
type SuperellipseOctant[T TScalar] struct {
	Offset         Point[T] // Center of the octant, relative to origin.
	SemiAxis       T        // Semi-axis length.
	Degree         T        // Degree of the superellipse. 0 means square.
	MaxTheta       T        // Range of the parametric "theta".
	CircleStart    Point[T] // Start point of the circular arc, relative to Offset.
	CircleCenter   Point[T] // Center of the circular arc, relative to Offset.
	CircleMaxAngle Radians  // Angular span of the circular arc, in radians.
}

// SuperellipseQuadrant holds parameters for a quadrant of a rounded superellipse.
//
// Used to define a quadrant of an arbitrary rounded superellipse.
type SuperellipseQuadrant[T TScalar] struct {
	Offset      Point[T] // Center of the quadrant, relative to origin.
	SignedScale Point[T] // Scaling factor for normalization and flipping.
	Top         SuperellipseOctant[T]
	Right       SuperellipseOctant[T]
}

// SuperellipseParam expands input parameters for a rounded superellipse to drawing variables.
type SuperellipseParam[T TScalar] struct {
	TopRight    SuperellipseQuadrant[T]
	BottomRight SuperellipseQuadrant[T]
	BottomLeft  SuperellipseQuadrant[T]
	TopLeft     SuperellipseQuadrant[T]
	IsUniform   bool // If true, only TopRight is populated.
}

// NewSuperellipseParam computes drawing parameters for a superellipse with given bounds and radii.
func NewSuperellipseParam[T TScalar](bounds Rect[T], radii RoundingRadii[T]) SuperellipseParam[T] {
	if radii.IsUniform() && !radii.TopRight.IsZero() {
		sq := newSuperellipseQuadrant(bounds.Center(), bounds.TopRight(), radii.TopRight, Size[T]{-1, 1})
		return SuperellipseParam[T]{
			TopRight:  sq,
			IsUniform: radii.IsUniform(),
		}
	}

	split := func(left, right, ratioLeft, ratioRight T) T {
		if ratioLeft == 0 && ratioRight == 0 {
			return (left + right) / 2
		}
		return (left*ratioLeft + right*ratioLeft) / (ratioLeft + ratioRight)
	}

	topSplit := split(bounds.Left, bounds.Right, radii.TopLeft.Width, radii.TopRight.Width)
	rightSplit := split(bounds.Top, bounds.Bottom, radii.TopRight.Height, radii.BottomRight.Height)
	bottomSplit := split(bounds.Left, bounds.Right, radii.BottomLeft.Width, radii.BottomRight.Width)
	leftSplit := split(bounds.Top, bounds.Bottom, radii.TopLeft.Height, radii.BottomLeft.Height)

	return SuperellipseParam[T]{
		TopRight: newSuperellipseQuadrant(
			Point[T]{topSplit, rightSplit}, bounds.TopRight(), radii.TopRight, Size[T]{1, -1},
		),
		BottomRight: newSuperellipseQuadrant(
			Point[T]{bottomSplit, rightSplit}, bounds.BottomRight(), radii.BottomRight, Size[T]{1, 1},
		),
		BottomLeft: newSuperellipseQuadrant(
			Point[T]{bottomSplit, leftSplit}, bounds.BottomLeft(), radii.BottomLeft, Size[T]{-1, 1},
		),
		TopLeft: newSuperellipseQuadrant(
			Point[T]{topSplit, leftSplit}, bounds.TopLeft(), radii.TopLeft, Size[T]{-1, -1},
		),
		IsUniform: false,
	}
}

// NewSuperellipseParamRadius computes drawing parameters for a superellipse with uniform radius.
func NewSuperellipseParamRadius[T TScalar](bounds Rect[T], radius T) SuperellipseParam[T] {
	sq := newSuperellipseQuadrant(bounds.Center(), bounds.TopRight(), Size[T]{radius, radius}, Size[T]{-1, 1})
	return SuperellipseParam[T]{
		TopRight:  sq,
		IsUniform: true,
	}
}

// newSuperellipseQuadrant computes parameters for a quadrant of a rounded superellipse with asymmetrical radii.
// center is the quadrant center, corner is the corner point, sign indicates the quadrant orientation.
func newSuperellipseQuadrant[T TScalar](center Point[T], corner Point[T], in_radii Size[T], sign Size[T]) SuperellipseQuadrant[T] {

	centerVector := corner.Sub(center)
	radii := in_radii.Abs()
	normRadius := radii.MinDimension()
	var forwardScale Size[T]
	if normRadius == 0 {
		forwardScale = Size[T]{1, 1}
	} else {
		forwardScale = radii.Scale(1 / normRadius)
	}

	normHalfSize := centerVector.Abs().DivSize(forwardScale)
	signedScale := replaceNaNWithDefault(centerVector.Div(normHalfSize), sign)

	c := normHalfSize.X - normHalfSize.Y
	return SuperellipseQuadrant[T]{
		Offset:      center,
		SignedScale: signedScale,
		Top:         newSuperellipseOctant(Point[T]{0, -c}, normHalfSize.X, normRadius),
		Right:       newSuperellipseOctant(Point[T]{c, 0}, normHalfSize.Y, normRadius),
	}

}

// newSuperellipseOctant computes parameters for a square-like rounded superellipse with a symmetrical radius.
// center is the octant center, a is the semi-axis, radius is the corner radius.
func newSuperellipseOctant[T TScalar](center Point[T], a T, radius T) SuperellipseOctant[T] {
	const gapFactor = 0.29289321881 // 1-cos(pi/4)

	if radius <= 0 {
		return SuperellipseOctant[T]{
			Offset:         center,
			SemiAxis:       a,
			Degree:         0,
			CircleStart:    Point[T]{a, a},
			CircleCenter:   Point[T]{0, 0},
			CircleMaxAngle: Radians(0),
		}
	}

	af := ToFloat64(a)

	ratio := a * 2 / radius
	g := T(gapFactor * ToFloat64(radius))
	precomputedVars := superellipseComputeNAndXj(ToFloat64(ratio))
	n := precomputedVars[0]
	xJ := precomputedVars[1] * af
	yJ := math.Pow(1-math.Pow(precomputedVars[1], n), 1/n) * af
	maxTheta := math.Asin(math.Pow(precomputedVars[1], n/2))
	tanPhiJ := math.Pow(xJ/yJ, n-1)
	d := (xJ - tanPhiJ*yJ) / (1 - tanPhiJ)
	R := T((ToFloat64(a) - d - ToFloat64(g)) * math.Sqrt(2))

	PointM := Point[T]{a - g, a - g}
	pointJ := Point[T]{T(xJ), T(yJ)}
	var circleCenter Point[T]
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

	return SuperellipseOctant[T]{
		Offset:         center,
		SemiAxis:       a,
		Degree:         T(n),
		MaxTheta:       T(maxTheta),
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
func findCircleCenter[T TScalar](a Point[T], b Point[T], r T) Point[T] {
	aToB := b.Sub(a)
	m := a.Add(b).Scale(T(1) / T(2))
	cToM := Point[T]{-aToB.Y, aToB.X}
	distanceAM := aToB.Length() / 2
	distanceCM := T(math.Sqrt(ToFloat64(r*r - distanceAM*distanceAM)))
	return m.Sub(cToM.Normalize().Scale(distanceCM))
}

// replaceNaNWithDefault replaces NaN values in v with the corresponding values from defaultValue.
func replaceNaNWithDefault[T TScalar](v Point[T], defaultValue Size[T]) Point[T] {
	x := defaultValue.Width
	y := defaultValue.Height

	if !math.IsNaN(ToFloat64(v.X)) {
		x = v.X
	}

	if !math.IsNaN(ToFloat64(v.Y)) {
		y = v.Y
	}

	return Point[T]{x, y}
}
