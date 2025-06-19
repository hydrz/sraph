package geom

import (
	"math"
	"testing"
)

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
	if !ps.Bounds().Eq(rect) {
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
	if !ps.Bounds().Eq(rect) {
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
		Offset:         Point[F64]{0.0, 0.0},
		SemiAxis:       F64(20),
		Degree:         F64(4),
		CircleStart:    Point[F64]{10.0, 10.0},
		CircleCenter:   Point[F64]{0.0, 0.0},
		CircleMaxAngle: Radians(math.PiOver2),
	}
	_ = builder.circularArcPoints(octant)
	_ = builder.superellipseArcPoints(octant)
	_ = builder.superellipseBezierFactors(4)
}

func TestSuperellipse_FindCircleCenterAndReplaceNaN(t *testing.T) {
	a := Point[F64]{0.0, 0.0}
	b := Point[F64]{10.0, 0.0}
	r := F64(5)
	center := findCircleCenter(a, b, r)
	if math.IsNaN(center.X.Float64()) || math.IsNaN(center.Y.Float64()) {
		t.Error("findCircleCenter returned NaN")
	}
	v := Point[F64]{F64(math.NaN()), 2}
	def := Size[F64]{1, 3}
	res := replaceNaNWithDefault(v, def)
	if !ScalarEq(res.X, F64(1)) || !ScalarEq(res.Y, F64(2)) {
		t.Error("replaceNaNWithDefault did not replace NaN as expected")
	}
}
