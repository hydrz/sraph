package geom

import (
	"math"
	"testing"
)

func TestSuperellipse_Dispatch_UniformAndNonUniform(t *testing.T) {
	rect := NewRect(0.0, 0.0, 100.0, 100.0)
	tests := []struct {
		name  string
		radii RoundingRadii
	}{
		{"Uniform", NewRoundingRadii(20.0)},
		{"NonUniform", NewRoundingRadiiLTRB(10.0, 20.0, 30.0, 40.0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			se := NewSuperellipse(rect, tt.radii)
			receiver := &testPathReceiver{}
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
	rect := NewRect(0.0, 0.0, 100.0, 100.0)
	se := NewSuperellipse(rect, NewRoundingRadii(10.0))
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
	receiver := &testPathReceiver{}
	ps.Dispatch(receiver)
	if len(receiver.moves) == 0 {
		t.Error("SuperellipsePathSource Dispatch did not call MoveTo")
	}
}

func TestSuperellipse_DiffPathSource(t *testing.T) {
	rect := NewRect(0.0, 0.0, 100.0, 100.0)
	se1 := NewSuperellipse(rect, NewRoundingRadii(10.0))
	se2 := NewSuperellipse(rect, NewRoundingRadii(5.0))
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
	receiver := &testPathReceiver{}
	ps.Dispatch(receiver)
	if len(receiver.moves) == 0 {
		t.Error("DiffSuperellipsePathSource Dispatch did not call MoveTo")
	}
}

func TestSuperellipse_ParamUniformAndNonUniform(t *testing.T) {
	rect := NewRect(0.0, 0.0, 100.0, 100.0)
	param := NewSuperellipseParam(rect, NewRoundingRadii(10.0))
	if !param.IsUniform {
		t.Error("SuperellipseParam should be uniform for uniform radii")
	}
	param2 := NewSuperellipseParam(rect, NewRoundingRadiiLTRB(10.0, 20.0, 30.0, 40.0))
	if param2.IsUniform {
		t.Error("SuperellipseParam should not be uniform for non-uniform radii")
	}
}

func TestSuperellipse_InternalHelpers(t *testing.T) {
	builder := &superellipseBuilder{}
	octant := SuperellipseOctant{
		Offset:         NewPoint[Scalar](0.0, 0.0),
		SemiAxis:       20,
		Degree:         4,
		CircleStart:    NewPoint[Scalar](10.0, 10.0),
		CircleCenter:   NewPoint[Scalar](0.0, 0.0),
		CircleMaxAngle: Radians(math.Pi / 2),
	}
	_ = builder.circularArcPoints(octant)
	_ = builder.arcPoints(octant)
	_ = builder.bezierFactors(4)
}

func TestSuperellipse_FindCircleCenterAndReplaceNaN(t *testing.T) {
	a := NewPoint[Scalar](0.0, 0.0)
	b := NewPoint[Scalar](10.0, 0.0)
	r := Scalar(5.0)
	center := findCircleCenter(a, b, r)
	if math.IsNaN(ToFloat64(center.X())) || math.IsNaN(ToFloat64(center.Y())) {
		t.Error("findCircleCenter returned NaN")
	}
	v := NewPoint(Scalar(math.NaN()), 2)
	def := NewSize[Scalar](1, 3)
	res := replaceNaNWithDefault(v, def)
	if !NearlyEqual(res.X(), 1) || !NearlyEqual(res.Y(), 2) {
		t.Error("replaceNaNWithDefault did not replace NaN as expected")
	}
}
