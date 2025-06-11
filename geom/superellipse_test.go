package geom

import (
	"math"
	"testing"
)

func TestSuperellipse_NewSuperellipseVariants(t *testing.T) {
	type args struct {
		rect  Rect[Float64]
		radii RoundingRadii[Float64]
	}
	tests := []struct {
		name  string
		rect  Rect[Float64]
		radii RoundingRadii[Float64]
	}{
		{
			name:  "Uniform radius",
			rect:  NewRect(Float64(0), Float64(0), Float64(100), Float64(50)),
			radii: NewRoundingRadii(Float64(10)),
		},
		{
			name:  "Zero radius",
			rect:  NewRect(Float64(0), Float64(0), Float64(100), Float64(50)),
			radii: NewRoundingRadii(Float64(0)),
		},
		{
			name:  "Non-uniform radii",
			rect:  NewRect(Float64(0), Float64(0), Float64(100), Float64(100)),
			radii: NewRoundingRadiiLTRB(Float64(10), Float64(20), Float64(30), Float64(40)),
		},
		{
			name:  "Negative radius (should handle gracefully)",
			rect:  NewRect(Float64(0), Float64(0), Float64(100), Float64(50)),
			radii: NewRoundingRadii(Float64(-5)),
		},
		{
			name:  "Empty rect",
			rect:  NewRect(Float64(0), Float64(0), Float64(0), Float64(0)),
			radii: NewRoundingRadii(Float64(10)),
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
	rect := NewRect(Float64(0), Float64(0), Float64(80), Float64(40))
	if se := NewSuperellipseOval(rect); se == nil {
		t.Error("NewSuperellipseOval returned nil")
	}
	rect2 := NewRect(Float64(0), Float64(0), Float64(60), Float64(60))
	if se := NewSuperellipseRadius(rect2, Float64(15)); se == nil {
		t.Error("NewSuperellipseRadius returned nil")
	}
	rect3 := NewRect(Float64(0), Float64(0), Float64(120), Float64(60))
	if se := NewSuperellipseXY(rect3, Float64(20), Float64(10)); se == nil {
		t.Error("NewSuperellipseXY returned nil")
	}
	rect4 := NewRect(Float64(0), Float64(0), Float64(100), Float64(100))
	if se := NewSuperellipseLTRB(rect4, Float64(10), Float64(20), Float64(30), Float64(40)); se == nil {
		t.Error("NewSuperellipseLTRB returned nil")
	}
}

func TestSuperellipse_Dispatch_UniformAndNonUniform(t *testing.T) {
	rect := NewRect(Float64(0), Float64(0), Float64(100), Float64(100))
	tests := []struct {
		name  string
		radii RoundingRadii[Float64]
	}{
		{"Uniform", NewRoundingRadii(Float64(20))},
		{"NonUniform", NewRoundingRadiiLTRB(Float64(10), Float64(20), Float64(30), Float64(40))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			se := NewSuperellipse(rect, tt.radii)
			receiver := &testPathReceiver[Float64]{}
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
	rect := NewRect(Float64(0), Float64(0), Float64(100), Float64(100))
	se := NewSuperellipse(rect, NewRoundingRadii(Float64(10)))
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
	receiver := &testPathReceiver[Float64]{}
	ps.Dispatch(receiver)
	if len(receiver.moves) == 0 {
		t.Error("SuperellipsePathSource Dispatch did not call MoveTo")
	}
}

func TestSuperellipse_DiffPathSource(t *testing.T) {
	rect := NewRect(Float64(0), Float64(0), Float64(100), Float64(100))
	se1 := NewSuperellipse(rect, NewRoundingRadii(Float64(10)))
	se2 := NewSuperellipse(rect, NewRoundingRadii(Float64(5)))
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
	receiver := &testPathReceiver[Float64]{}
	ps.Dispatch(receiver)
	if len(receiver.moves) == 0 {
		t.Error("DiffSuperellipsePathSource Dispatch did not call MoveTo")
	}
}

func TestSuperellipse_ParamUniformAndNonUniform(t *testing.T) {
	rect := NewRect(Float64(0), Float64(0), Float64(100), Float64(100))
	param := NewSuperellipseParam(rect, NewRoundingRadii(Float64(10)))
	if !param.IsUniform {
		t.Error("SuperellipseParam should be uniform for uniform radii")
	}
	param2 := NewSuperellipseParam(rect, NewRoundingRadiiLTRB(Float64(10), Float64(20), Float64(30), Float64(40)))
	if param2.IsUniform {
		t.Error("SuperellipseParam should not be uniform for non-uniform radii")
	}
}

func TestSuperellipse_InternalHelpers(t *testing.T) {
	builder := &superellipseBuilder[Float64]{}
	octant := SuperellipseOctant[Float64]{
		Offset:         NewPoint(Float64(0), Float64(0)),
		SemiAxis:       Float64(20),
		Degree:         Float64(4),
		CircleStart:    NewPoint(Float64(10), Float64(10)),
		CircleCenter:   NewPoint(Float64(0), Float64(0)),
		CircleMaxAngle: Radians(math.Pi / 2),
	}
	_ = builder.circularArcPoints(octant)
	_ = builder.superellipseArcPoints(octant)
	_ = builder.superellipseBezierFactors(4)
}

func TestSuperellipse_FindCircleCenterAndReplaceNaN(t *testing.T) {
	a := NewPoint(Float64(0), Float64(0))
	b := NewPoint(Float64(10), Float64(0))
	r := Float64(5)
	center := findCircleCenter(a, b, r)
	if math.IsNaN(center.X().Float64()) || math.IsNaN(center.Y().Float64()) {
		t.Error("findCircleCenter returned NaN")
	}
	v := NewPoint(Float64(math.NaN()), Float64(2))
	def := NewSize(Float64(1), Float64(3))
	res := replaceNaNWithDefault(v, def)
	if !Equal(res.X(), Float64(1)) || !Equal(res.Y(), Float64(2)) {
		t.Error("replaceNaNWithDefault did not replace NaN as expected")
	}
}
