package geom

import (
	"math"
	"testing"
)

func TestCubicSubdivisions(t *testing.T) {
	tests := []struct {
		name           string
		scaleFactor    Scalar
		p0, p1, p2, p3 Point[Scalar]
		wantMin        Scalar // minimum expected subdivisions
	}{
		{
			name:        "straight line cubic",
			scaleFactor: 1.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](1, 0),
			p2:          NewPoint[Scalar](2, 0),
			p3:          NewPoint[Scalar](3, 0),
			wantMin:     0,
		},
		{
			name:        "curved cubic with scale factor 1",
			scaleFactor: 1.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](1, 1),
			p2:          NewPoint[Scalar](2, 1),
			p3:          NewPoint[Scalar](3, 0),
			wantMin:     1,
		},
		{
			name:        "curved cubic with higher scale factor",
			scaleFactor: 2.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](1, 1),
			p2:          NewPoint[Scalar](2, 1),
			p3:          NewPoint[Scalar](3, 0),
			wantMin:     1,
		},
		{
			name:        "sharp curve cubic",
			scaleFactor: 1.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](0, 10),
			p2:          NewPoint[Scalar](0, 10),
			p3:          NewPoint[Scalar](10, 0),
			wantMin:     3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CubicSubdivisions(tt.scaleFactor, tt.p0, tt.p1, tt.p2, tt.p3)
			if got < tt.wantMin {
				t.Errorf("CubicSubdivisions() = %v, want >= %v", got, tt.wantMin)
			}
			if got < 0 {
				t.Errorf("CubicSubdivisions() = %v, subdivisions should not be negative", got)
			}
		})
	}
}

func TestQuadraticSubdivisions(t *testing.T) {
	tests := []struct {
		name        string
		scaleFactor Scalar
		p0, p1, p2  Point[Scalar]
		wantMin     Scalar
	}{
		{
			name:        "straight line quadratic",
			scaleFactor: 1.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](1, 0),
			p2:          NewPoint[Scalar](2, 0),
			wantMin:     0,
		},
		{
			name:        "curved quadratic with scale factor 1",
			scaleFactor: 1.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](1, 1),
			p2:          NewPoint[Scalar](2, 0),
			wantMin:     1,
		},
		{
			name:        "curved quadratic with higher scale factor",
			scaleFactor: 2.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](1, 1),
			p2:          NewPoint[Scalar](2, 0),
			wantMin:     1,
		},
		{
			name:        "sharp curve quadratic",
			scaleFactor: 1.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](0, 10),
			p2:          NewPoint[Scalar](10, 0),
			wantMin:     3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QuadraticSubdivisions(tt.scaleFactor, tt.p0, tt.p1, tt.p2)
			if got < tt.wantMin {
				t.Errorf("QuadraticSubdivisions() = %v, want >= %v", got, tt.wantMin)
			}
			if got < 0 {
				t.Errorf("QuadraticSubdivisions() = %v, subdivisions should not be negative", got)
			}
		})
	}
}

func TestConicSubdivisions(t *testing.T) {
	tests := []struct {
		name        string
		scaleFactor Scalar
		p0, p1, p2  Point[Scalar]
		weight      Scalar
		wantMin     Scalar
	}{
		{
			name:        "straight line conic with weight 1",
			scaleFactor: 1.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](1, 0),
			p2:          NewPoint[Scalar](2, 0),
			weight:      1.0,
			wantMin:     0,
		},
		{
			name:        "curved conic with weight 1",
			scaleFactor: 1.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](1, 1),
			p2:          NewPoint[Scalar](2, 0),
			weight:      1.0,
			wantMin:     1,
		},
		{
			name:        "curved conic with weight 0.5",
			scaleFactor: 1.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](1, 1),
			p2:          NewPoint[Scalar](2, 0),
			weight:      0.5,
			wantMin:     1,
		},
		{
			name:        "curved conic with weight 2",
			scaleFactor: 1.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](1, 1),
			p2:          NewPoint[Scalar](2, 0),
			weight:      2.0,
			wantMin:     1,
		},
		{
			name:        "sharp curve conic",
			scaleFactor: 1.0,
			p0:          NewPoint[Scalar](0, 0),
			p1:          NewPoint[Scalar](0, 10),
			p2:          NewPoint[Scalar](10, 0),
			weight:      1.0,
			wantMin:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConicSubdivisions(tt.scaleFactor, tt.p0, tt.p1, tt.p2, tt.weight)
			if got < tt.wantMin {
				t.Errorf("ConicSubdivisions() = %v, want >= %v", got, tt.wantMin)
			}
			if got < 0 {
				t.Errorf("ConicSubdivisions() = %v, subdivisions should not be negative", got)
			}
			if math.IsNaN(ToFloat64(got)) || math.IsInf(ToFloat64(got), 0) {
				t.Errorf("ConicSubdivisions() = %v, result should be finite", got)
			}
		})
	}
}
