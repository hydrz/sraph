package geom

import "math"

// Superellipse represents a Superellipse shape, extending Rect.
type Superellipse[T Scalar] struct {
	RoundRect[T]
	param SuperellipseParam[T]
}

// NewSuperellipse creates a new superellipse with the given rectangle and corner radii.
func NewSuperellipse[T Scalar](rect Rect[T], radii RoundingRadii[T]) Superellipse[T] {
	param := NewSuperellipseParam(rect, radii)
	return Superellipse[T]{
		RoundRect: NewRoundRect(rect, radii),
		param:     param,
	}
}

// NewSuperellipseOval creates a new Superellipse that is an oval, with radii Eq to half the width and height of the rectangle.
func NewSuperellipseOval[T Scalar](rect Rect[T]) Superellipse[T] {
	return NewSuperellipse(
		rect,
		NewRoundingRadiiFromSizes(rect.Size().Scale(1/2)),
	)
}

// NewSuperellipseRadius creates a new Superellipse with the given rectangle and uniform corner radius.
func NewSuperellipseRadius[T Scalar](rect Rect[T], radius T) Superellipse[T] {
	return NewSuperellipse(
		rect,
		NewRoundingRadii(radius),
	)
}

// NewSuperellipseXY creates a new Superellipse with the given rectangle and separate x/y radii.
func NewSuperellipseXY[T Scalar](rect Rect[T], xRadius, yRadius T) Superellipse[T] {
	return NewSuperellipse(
		rect,
		NewRoundingRadiiFromSizes(Size[T]{Width: xRadius, Height: yRadius}),
	)
}

// NewSuperellipseLTRB creates a new Superellipse with the given rectangle and separate left, top, right, bottom radii.
func NewSuperellipseLTRB[T Scalar](rect Rect[T], left, top, right, bottom T) Superellipse[T] {
	return NewSuperellipse(
		rect,
		NewRoundingRadiiLTRB(left, top, right, bottom),
	)
}

// ToApproximateRoundRect returns a rounded rectangle that approximates the current superellipse.
// This is used for backends that do not support superellipses.
func (s *Superellipse[T]) ToApproximateRoundRect() RoundRect[T] {
	return NewRoundRect(s.RoundRect.Bounds(), s.RoundRect.Radius())
}

// Dispatch overrides RoundRect.Dispatch to handle superellipse-specific logic.
func (s *Superellipse[T]) Dispatch(receiver PathReceiver[T], includeEnd bool) {
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

// superellipseBuilder is a helper for building superellipse paths.
type superellipseBuilder[T Scalar] struct {
	receiver PathReceiver[T]
}

// AddQuadrant adds a quadrant of the superellipse to the path.
// reverse: whether to reverse the drawing direction.
// scaleSign: sign vector for scaling and flipping.
func (b *superellipseBuilder[T]) AddQuadrant(quadrant SuperellipseQuadrant[T], reverse bool, scaleSign Point[T]) {
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
// reverse: whether to reverse the curve direction.
// flip: whether to flip the octant (swap x/y).
// externalTransform: transformation matrix to apply.
func (b *superellipseBuilder[T]) AddOctant(octant SuperellipseOctant[T], reverse, flip bool, externalTransform Matrix[T]) {
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
	sePoints := b.superellipseArcPoints(octant)

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
func (b *superellipseBuilder[T]) circularArcPoints(octant SuperellipseOctant[T]) [4]Point[T] {
	startVector := octant.CircleStart.Sub(octant.CircleCenter)
	endVector := startVector.Rotate(Radians(-octant.CircleMaxAngle))
	circleEnd := octant.CircleCenter.Add(endVector)
	startTangent := Point[T]{X: startVector.Y, Y: -startVector.X}.Normalize()
	endTangent := Point[T]{X: -endVector.Y, Y: endVector.X}.Normalize()
	bezierFactor := math.Tan(octant.CircleMaxAngle.Float64() / 4 * 4 / 3)
	radius := startVector.Length()

	return [4]Point[T]{
		octant.CircleStart,
		octant.CircleStart.Add(startTangent.Scale(T(bezierFactor) * radius)),
		circleEnd.Add(endTangent.Scale(T(bezierFactor) * radius)),
		circleEnd,
	}
}

// superellipseArcPoints returns the four control points for the superellipse arc segment of the octant.
func (b *superellipseBuilder[T]) superellipseArcPoints(octant SuperellipseOctant[T]) [4]Point[T] {
	start := Point[T]{X: 0, Y: octant.SemiAxis}
	end := octant.CircleStart
	startTangent := Point[T]{X: 1, Y: 0}
	circleStartVector := octant.CircleStart.Sub(octant.CircleCenter)
	endTangent := Point[T]{X: -circleStartVector.Y, Y: circleStartVector.X}.Normalize()
	factors := b.superellipseBezierFactors(octant.SemiAxis.Float64())
	return [4]Point[T]{
		start,
		start.Add(startTangent.Scale(T(factors[0] * octant.SemiAxis.Float64()))),
		end.Add(endTangent.Scale(T(factors[1] * octant.SemiAxis.Float64()))),
		end,
	}
}

// superellipseBezierFactors returns the Bezier control factors for a given superellipse degree.
// Uses precomputed values and interpolation for best accuracy.
func (b *superellipseBuilder[T]) superellipseBezierFactors(n float64) [2]float64 {
	kPrecomputedVariables := [][2]float64{
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
	numRecords := len(kPrecomputedVariables)
	step := 1
	minN := 2
	maxN := minN + (numRecords-1)*step

	if int(n) >= maxN {
		// Use heuristic formula for large n.
		return [2]float64{
			1.07 - math.Exp(1.307649835)*math.Pow(n, -0.8568516731),
			-0.01 + math.Exp(-0.9287690322)*math.Pow(n, -0.6120901398),
		}
	}

	steps := Clamp((n-float64(minN))/float64(step), 0, float64(numRecords-1))
	left := Clamp(int(math.Floor(steps)), 0, numRecords-2)
	frac := steps - float64(left)
	return [2]float64{
		(1-frac)*kPrecomputedVariables[left][0] + frac*kPrecomputedVariables[left+1][0],
		(1-frac)*kPrecomputedVariables[left][1] + frac*kPrecomputedVariables[left+1][1],
	}
}

// SuperellipsePathSource is a PathSource implementation for a single superellipse.
type SuperellipsePathSource[T Scalar] struct {
	superellipse Superellipse[T]
}

// NewSuperellipsePathSource creates a new PathSource for a superellipse.
func NewSuperellipsePathSource[T Scalar](superellipse Superellipse[T]) PathSource[T] {
	return &SuperellipsePathSource[T]{superellipse: superellipse}
}

// FillType implements PathSource.
func (s *SuperellipsePathSource[T]) FillType() FillType {
	return FillTypeNonZero
}

// Bounds implements PathSource.
func (s *SuperellipsePathSource[T]) Bounds() Rect[T] {
	return s.superellipse.Bounds()
}

// IsConvex implements PathSource.
func (s *SuperellipsePathSource[T]) IsConvex() bool {
	return true // Superellipses are convex shapes.
}

// Dispatch implements PathSource.
func (s *SuperellipsePathSource[T]) Dispatch(receiver PathReceiver[T]) {
	s.superellipse.Dispatch(receiver, true)
}

// DiffSuperellipsePathSource is a PathSource for the difference of two superellipses.
type DiffSuperellipsePathSource[T Scalar] struct {
	outter Superellipse[T]
	inner  Superellipse[T]
}

// NewDiffSuperellipsePathSource creates a new PathSource for the difference of two superellipses.
func NewDiffSuperellipsePathSource[T Scalar](outter, inner Superellipse[T]) PathSource[T] {
	return &DiffSuperellipsePathSource[T]{outter: outter, inner: inner}
}

// FillType implements PathSource.
func (d *DiffSuperellipsePathSource[T]) FillType() FillType {
	return FillTypeEvenOdd
}

// IsConvex implements PathSource.
func (d *DiffSuperellipsePathSource[T]) IsConvex() bool {
	return false // Difference of superellipses is not convex.
}

// Bounds implements PathSource.
func (d *DiffSuperellipsePathSource[T]) Bounds() Rect[T] {
	return d.outter.Bounds()
}

// Dispatch implements PathSource.
func (d *DiffSuperellipsePathSource[T]) Dispatch(receiver PathReceiver[T]) {
	// Draw the outer superellipse, then the inner one in reverse for even-odd fill.
	d.outter.Dispatch(receiver, false)
	d.inner.Dispatch(receiver, true)
}

// SuperellipseOctant holds parameters for drawing a square-like rounded superellipse octant.
//
// A Degree of 0 means that the radius is 0, and this octant is a square
// of size SemiAxis at Offset. All other fields are ignored in this case.
type SuperellipseOctant[T Scalar] struct {
	// Offset of the octant's center from the origin.
	// All other coordinates in this struct are relative to this point.
	Offset Point[T]

	// SemiAxis is the semi-axis length of the superellipse.
	SemiAxis T
	// Degree is the degree of the superellipse.
	// If 0, this octant is a square of size SemiAxis.
	Degree T
	// MaxTheta is the range of the parameter "theta" used to define the superellipse curve.
	// "theta" is not the angle but the implicit parameter in the curve's parametric equation.
	MaxTheta T

	// CircleStart is the coordinate of the top left end of the circular arc, relative to Offset.
	CircleStart Point[T]
	// CircleCenter is the center of the circular arc, relative to Offset.
	CircleCenter Point[T]
	// CircleMaxAngle is the angular span of the circular arc, in radians.
	CircleMaxAngle Radians
}

// SuperellipseQuadrant holds parameters for a quadrant of a rounded superellipse.
//
// This struct is used to define a quadrant of an arbitrary rounded superellipse.
type SuperellipseQuadrant[T Scalar] struct {
	// Offset of the quadrant's center from the origin.
	// All other coordinates in this struct are relative to this point.
	Offset Point[T]

	// SignedScale is the scaling factor used to transform a normalized rounded superellipse
	// back to its original, unnormalized shape.
	//
	// Normalization means adjusting the original curve (possibly with asymmetrical corner sizes)
	// into a symmetrical one by reducing the longer radius to match the shorter one.
	// For example, to draw a rounded superellipse with size (200, 300) and radii (20, 10),
	// the function first draws a normalized RSE with size (100, 300) and radii (10, 10),
	// then scales it by (2x, 1x) to restore the original proportions.
	//
	// Normalization also flips the curve to the first quadrant (positive x and y)
	// if it originally resides in another quadrant. This affects the signs of SignedScale.
	SignedScale Point[T]

	// The two octants that make up this quadrant after normalization.
	Top   SuperellipseOctant[T]
	Right SuperellipseOctant[T]
}

// SuperellipseParam expands input parameters for a rounded superellipse to drawing variables.
type SuperellipseParam[T Scalar] struct {
	// The four quadrants that make up the full contour.
	TopRight    SuperellipseQuadrant[T]
	BottomRight SuperellipseQuadrant[T]
	BottomLeft  SuperellipseQuadrant[T]
	TopLeft     SuperellipseQuadrant[T]
	// If IsUniform is true, only TopRight is populated.
	IsUniform bool
}

func NewSuperellipseParam[T Scalar](bounds Rect[T], radii RoundingRadii[T]) SuperellipseParam[T] {
	if radii.IsUniform() && !radii.TopRight.IsZero() {
		sq := newSuperellipseQuadrant(bounds.Center(), bounds.RightTop(), radii.TopRight, Size[T]{-1, 1})
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
			Point[T]{topSplit, rightSplit}, bounds.RightTop(), radii.TopRight, Size[T]{1, -1},
		),
		BottomRight: newSuperellipseQuadrant(
			Point[T]{bottomSplit, rightSplit}, bounds.RightBottom(), radii.BottomRight, Size[T]{1, 1},
		),
		BottomLeft: newSuperellipseQuadrant(
			Point[T]{bottomSplit, leftSplit}, bounds.LeftBottom(), radii.BottomLeft, Size[T]{-1, 1},
		),
		TopLeft: newSuperellipseQuadrant(
			Point[T]{topSplit, leftSplit}, bounds.LeftTop(), radii.TopLeft, Size[T]{-1, -1},
		),
		IsUniform: false,
	}
}

func NewSuperellipseParamRadius[T Scalar](bounds Rect[T], radius T) SuperellipseParam[T] {
	sq := newSuperellipseQuadrant(bounds.Center(), bounds.RightTop(), Size[T]{radius, radius}, Size[T]{-1, 1})
	return SuperellipseParam[T]{
		TopRight:  sq,
		IsUniform: true,
	}
}

// Compute parameters for a quadrant of a rounded superellipse with asymmetrical
// radii.
//
// The `corner` is the coordinate of the corner point in the same coordinate
// space as `center`, which specifies the half size of the bounding box.
//
// The `sign` is a vector of {±1, ±1} that specifies which quadrant the curve
// should be, which should have the same sign as `corner - center` except that
// the latter may have a 0.
func newSuperellipseQuadrant[T Scalar](center Point[T], corner Point[T], in_radii Size[T], sign Size[T]) SuperellipseQuadrant[T] {

	centerVector := corner.Sub(center)
	radii := in_radii.Abs()
	// The prefix "norm" is short for "normalized".
	//
	// Be extra careful to avoid NaNs in cases that some coordinates of `in_radii`
	// or `corner_vector` are zero.
	normRadius := radii.MinDimension()
	var forwardScale Size[T]
	if normRadius == 0 {
		forwardScale = Size[T]{1, 1}
	} else {
		forwardScale = radii.Scale(1 / normRadius)
	}

	normHalfSize := centerVector.Abs().DivSize(forwardScale)
	signedScale := replaceNaNWithDefault(centerVector.Div(normHalfSize), sign)

	// Each quadrant curve is composed of two octant curves, each of which belongs
	// to a square-like rounded rectangle. For the two octants to connect at the
	// circular arc, the centers these two square-like rounded rectangle must be
	// offset from the quadrant center by a same distance in different directions.
	// The distance is denoted as `c`.
	c := normHalfSize.X - normHalfSize.Y
	return SuperellipseQuadrant[T]{
		Offset:      center,
		SignedScale: signedScale,
		Top:         newSuperellipseOctant(Point[T]{0, -c}, normHalfSize.X, normRadius),
		Right:       newSuperellipseOctant(Point[T]{c, 0}, normHalfSize.Y, normRadius),
	}

}

// Compute parameters for a square-like rounded superellipse with a symmetrical radius.
//
// The following figure shows the first quadrant of a square-like rounded
// superellipse.
//
//	       superellipse
//	 A     ↓            circular arc
//	 ---------...._J   ↙
//	 |           /   `⟍ M (where x=y)
//	 |          /     ⟋ ⟍
//	 |         /   ⟋     \
//	 |        / ⟋         |
//	 |       ᜱD           |
//	 |     ⟋              |
//	 |  ⟋                 |
//	 |⟋                   |
//	 +--------------------| A'
//	O
//	 ←-------- a ---------→
func newSuperellipseOctant[T Scalar](center Point[T], a T, radius T) SuperellipseOctant[T] {
	// gapFactor is used to calculate the "gap", which is the distance from the midpoint
	// of the curved corners to the nearest sides of the bounding box.
	//
	// For symmetrical corner radii, the midpoint is where the circular arc intersects
	// its quadrant bisector. For asymmetrical radii, the midpoint is transformed accordingly.
	//
	// Experiments show the gap is linear with respect to the corner radius on that dimension.
	const gapFactor = 0.29289321881 // 1-cos(pi/4)

	if radius <= 0 {
		// If the radius is zero, return a square octant.
		return SuperellipseOctant[T]{
			Offset:         center,
			SemiAxis:       a,
			Degree:         0,
			CircleStart:    Point[T]{a, a},
			CircleCenter:   Point[T]{0, 0},
			CircleMaxAngle: Radians(0),
		}
	}

	af := a.Float64()

	ratio := a * 2 / radius
	g := T(gapFactor * radius.Float64())
	precomputedVars := superellipseComputeNAndXj(ratio.Float64())
	n := precomputedVars[0]
	xJ := precomputedVars[1] * af
	yJ := math.Pow(1-math.Pow(precomputedVars[1], n), 1/n) * af
	maxTheta := math.Asin(math.Pow(precomputedVars[1], n/2))
	tanPhiJ := math.Pow(xJ/yJ, n-1)
	d := (xJ - tanPhiJ*yJ) / (1 - tanPhiJ)
	R := T((a.Float64() - d - g.Float64()) * math.Sqrt(2))

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

func superellipseComputeNAndXj(ratio float64) [2]float64 {
	const (
		minRatio = 2.00

		// The curve is split into three parts:
		// * The first part uses a denser lookup table.
		// * The second part uses a sparser lookup table.
		// * The third part uses a straight line.
		firstStepInverse = 10 // = 1/0.10
		firstMaxRatio    = 2.50
		firstNumRecords  = 6

		secondStepInverse = 2 // = 1/0.50
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
		kNumRecords = len(precomputedVariables) // Number of precomputed records
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

// Find the center of the circle that passes the given two points and have the
// given radius.
//
// Denote the middle point of A and B as M. The key is to find the center of
// the circle.
//
//	       A --__
//	        /  ⟍ `、
//	       /   M  ⟍\
//	      /       ⟋  B
//	     /     ⟋   ↗
//	    /   ⟋
//	   / ⟋    r
//	C ᜱ  ↙
func findCircleCenter[T Scalar](a Point[T], b Point[T], r T) Point[T] {
	aToB := b.Sub(a)
	m := a.Add(b).Scale(1 / 2)
	cToM := Point[T]{-aToB.Y, aToB.X}
	distanceAM := aToB.Length() / 2
	distanceCM := math.Sqrt(float64(r*r - distanceAM*distanceAM))
	return m.Sub(cToM.Normalize().Scale(T(distanceCM)))
}

func replaceNaNWithDefault[T Scalar](v Point[T], defaultValue Size[T]) Point[T] {
	x := defaultValue.Width
	y := defaultValue.Height

	if !math.IsNaN(v.X.Float64()) {
		x = v.X
	}

	if !math.IsNaN(v.Y.Float64()) {
		y = v.Y
	}

	return Point[T]{x, y}
}
