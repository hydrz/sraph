package geom

import (
	"math"
	"testing"
)

func TestSuperellipse_NewSuperellipseVariants(t *testing.T) {
	type args struct {
		rect  Rect[F64]
		radii RoundingRadii[F64]
	}
	tests := []struct {
		name  string
		rect  Rect[F64]
		radii RoundingRadii[F64]
	}{
		{
			name:  "Uniform radius",
			rect:  NewRect[F64](0.0, 0.0, 100.0, 50.0),
			radii: NewRoundingRadii[F64](10.0),
		},
		{
			name:  "Zero radius",
			rect:  NewRect[F64](0.0, 0.0, 100.0, 50.0),
			radii: NewRoundingRadii[F64](0.0),
		},
		{
			name:  "Non-uniform radii",
			rect:  NewRect[F64](0.0, 0.0, 100.0, 100.0),
			radii: NewRoundingRadiiLTRB[F64](10.0, 20.0, 30.0, 40.0),
		},
		{
			name:  "Negative radius (should handle gracefully)",
			rect:  NewRect[F64](0.0, 0.0, 100.0, 50.0),
			radii: NewRoundingRadii[F64](-5.0),
		},
		{
			name:  "Empty rect",
			rect:  NewRect[F64](0.0, 0.0, 0.0, 0.0),
			radii: NewRoundingRadii[F64](10.0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			se := NewSuperellipse(tt.rect, tt.radii)
			if se == nil {
				t.Fatal("NewSuperellipse returned nil")
			}
			if !se.Bounds().Equal(tt.rect) {
				t.Errorf("Superellipse bounds mismatch: got %v, want %v", se.Bounds(), tt.rect)
			}
			rr := se.ToApproximateRoundRect()
			if rr == nil {
				t.Error("ToApproximateRoundRect returned nil")
			}
		})
	}
}

func TestSuperellipse_Constructors(t *testing.T) {
	rect := NewRect[F64](0.0, 0.0, 80.0, 40.0)
	if se := NewSuperellipseOval(rect); se == nil {
		t.Error("NewSuperellipseOval returned nil")
	}
	rect2 := NewRect[F64](0.0, 0.0, 60.0, 60.0)
	if se := NewSuperellipseRadius(rect2, F64(15)); se == nil {
		t.Error("NewSuperellipseRadius returned nil")
	}
	rect3 := NewRect[F64](0.0, 0.0, 120.0, 60.0)
	if se := NewSuperellipseXY(rect3, F64(20), F64(10)); se == nil {
		t.Error("NewSuperellipseXY returned nil")
	}
	rect4 := NewRect[F64](0.0, 0.0, 100.0, 100.0)
	if se := NewSuperellipseLTRB(rect4, F64(10), F64(20), F64(30), F64(40)); se == nil {
		t.Error("NewSuperellipseLTRB returned nil")
	}
}

func TestSuperellipse_Dispatch_UniformAndNonUniform(t *testing.T) {
	rect := NewRect[F64](0.0, 0.0, 100.0, 100.0)
	tests := []struct {
		name  string
		radii RoundingRadii[F64]
	}{
		{"Uniform", NewRoundingRadii[F64](20.0)},
		{"NonUniform", NewRoundingRadiiLTRB[F64](10.0, 20.0, 30.0, 40.0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			se := NewSuperellipse(rect, tt.radii)
			receiver := &testPathReceiver[F64]{}
			se.Dispatch(receiver, true)
			if len(receiver.moves) == 0 {
				t.Error("Dispatch did not call MoveTo")
			}
			if receiver.closed == 0 {
				t.Error("Dispatch did not call Close")
			}
			if receiver.pathEnds == 0 {
				t.Error("Dispatch did not call PathEnd")
			}
		})
	}
}

func TestSuperellipse_PathSourceVariants(t *testing.T) {
	rect := NewRect[F64](0.0, 0.0, 100.0, 100.0)
	se := NewSuperellipse(rect, NewRoundingRadii[F64](10.0))
	ps := NewSuperellipsePathSource(se)
	if ps.FillType() != FillTypeNonZero {
		t.Error("SuperellipsePathSource FillType should be FillTypeNonZero")
	}
	if !ps.Bounds().Equal(rect) {
		t.Error("SuperellipsePathSource Bounds mismatch")
	}
	if !ps.IsConvex() {
		t.Error("SuperellipsePathSource should be convex")
	}
	receiver := &testPathReceiver[F64]{}
	ps.Dispatch(receiver)
	if len(receiver.moves) == 0 {
		t.Error("SuperellipsePathSource Dispatch did not call MoveTo")
	}
}

func TestSuperellipse_DiffPathSource(t *testing.T) {
	rect := NewRect[F64](0.0, 0.0, 100.0, 100.0)
	se1 := NewSuperellipse(rect, NewRoundingRadii[F64](10.0))
	se2 := NewSuperellipse(rect, NewRoundingRadii[F64](5.0))
	ps := NewDiffSuperellipsePathSource(se1, se2)
	if ps.FillType() != FillTypeEvenOdd {
		t.Error("DiffSuperellipsePathSource FillType should be FillTypeEvenOdd")
	}
	if ps.IsConvex() {
		t.Error("DiffSuperellipsePathSource should not be convex")
	}
	if !ps.Bounds().Equal(rect) {
		t.Error("DiffSuperellipsePathSource Bounds mismatch")
	}
	receiver := &testPathReceiver[F64]{}
	ps.Dispatch(receiver)
	if len(receiver.moves) == 0 {
		t.Error("DiffSuperellipsePathSource Dispatch did not call MoveTo")
	}
}

func TestSuperellipse_ParamUniformAndNonUniform(t *testing.T) {
	rect := NewRect[F64](0.0, 0.0, 100.0, 100.0)
	param := NewSuperellipseParam(rect, NewRoundingRadii[F64](10.0))
	if !param.IsUniform {
		t.Error("SuperellipseParam should be uniform for uniform radii")
	}
	param2 := NewSuperellipseParam(rect, NewRoundingRadiiLTRB[F64](10.0, 20.0, 30.0, 40.0))
	if param2.IsUniform {
		t.Error("SuperellipseParam should not be uniform for non-uniform radii")
	}
}

func TestSuperellipse_InternalHelpers(t *testing.T) {
	builder := &superellipseBuilder[F64]{}
	octant := SuperellipseOctant[F64]{
		Offset:         NewPoint[F64](0.0, 0.0),
		SemiAxis:       F64(20),
		Degree:         F64(4),
		CircleStart:    NewPoint[F64](10.0, 10.0),
		CircleCenter:   NewPoint[F64](0.0, 0.0),
		CircleMaxAngle: NewRadians[F64](math.Pi / 2),
	}
	_ = builder.circularArcPoints(octant)
	_ = builder.superellipseArcPoints(octant)
	_ = builder.superellipseBezierFactors(4)
}

func TestSuperellipse_FindCircleCenterAndReplaceNaN(t *testing.T) {
	a := NewPoint[F64](0.0, 0.0)
	b := NewPoint[F64](10.0, 0.0)
	r := F64(5)
	center := findCircleCenter(a, b, r)
	if math.IsNaN(center.X().Float64()) || math.IsNaN(center.Y().Float64()) {
		t.Error("findCircleCenter returned NaN")
	}
	v := NewPoint[F64](F64(math.NaN()), F64(2))
	def := NewSize(F64(1), F64(3))
	res := replaceNaNWithDefault(v, def)
	if !Equal(res.X(), F64(1)) || !Equal(res.Y(), F64(2)) {
		t.Error("replaceNaNWithDefault did not replace NaN as expected")
	}
}
