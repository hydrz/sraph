package geom

import "testing"

func TestRoundingRadii_NewAndProperties(t *testing.T) {
	r := NewRoundingRadii[Float32](5)
	if !r.TopLeft().Equal(NewSize[Float32](5, 5)) ||
		!r.TopRight().Equal(NewSize[Float32](5, 5)) ||
		!r.BottomLeft().Equal(NewSize[Float32](5, 5)) ||
		!r.BottomRight().Equal(NewSize[Float32](5, 5)) {
		t.Errorf("NewRoundingRadii: all corners should be (5,5)")
	}
	if !r.IsUniform() {
		t.Errorf("IsUniform: should be true for all corners equal")
	}
	if !r.IsFinite() {
		t.Errorf("IsFinite: should be true for finite radii")
	}
	if r.IsEmpty() {
		t.Errorf("IsEmpty: should be false for nonzero radii")
	}
}

func TestRoundingRadiiLTRB(t *testing.T) {
	r := NewRoundingRadiiLTRB[Float32](1, 2, 3, 4)
	if !r.TopLeft().Equal(NewSize[Float32](1, 2)) ||
		!r.TopRight().Equal(NewSize[Float32](3, 2)) ||
		!r.BottomLeft().Equal(NewSize[Float32](1, 4)) ||
		!r.BottomRight().Equal(NewSize[Float32](3, 4)) {
		t.Errorf("NewRoundingRadii4: corners not as expected")
	}
	if r.IsUniform() {
		t.Errorf("IsUniform: should be false for unequal corners")
	}
}

func TestRoundingRadii_IsEmpty(t *testing.T) {
	r := NewRoundingRadii[Float32](0)
	if !r.IsEmpty() {
		t.Errorf("IsEmpty: should be true for all zero radii")
	}
	r2 := NewRoundingRadiiLTRB[Float32](0, 1, 0, 0)
	if r2.IsEmpty() {
		t.Errorf("IsEmpty: should be false if any corner is nonzero")
	}
}

func TestRoundingRadii_Scale(t *testing.T) {
	r := NewRoundingRadiiLTRB[Float32](1, 2, 3, 4)
	s := r.Scale(2)
	if !s.TopLeft().Equal(NewSize[Float32](2, 4)) ||
		!s.TopRight().Equal(NewSize[Float32](6, 4)) ||
		!s.BottomLeft().Equal(NewSize[Float32](2, 8)) ||
		!s.BottomRight().Equal(NewSize[Float32](6, 8)) {
		t.Errorf("Scale: scaling not correct")
	}
}

func TestRoundingRadii_ScaleToFit(t *testing.T) {
	r := NewRoundingRadii[Float32](10)
	bounds := NewRect[Float32](0, 0, 15, 20)
	scaled := r.ScaleToFit(bounds)
	sumTop := scaled.TopLeft().Width() + scaled.TopRight().Width()
	sumLeft := scaled.TopLeft().Height() + scaled.BottomLeft().Height()
	if sumTop > bounds.Width()+1e-5 || sumLeft > bounds.Height()+1e-5 {
		t.Errorf("ScaleToFit: radii sum exceeds bounds")
	}
}

func TestRoundingRadii_Equal(t *testing.T) {
	r1 := NewRoundingRadiiLTRB[Float32](1, 2, 3, 4)
	r2 := NewRoundingRadiiLTRB[Float32](1, 2, 3, 4)
	r3 := NewRoundingRadiiLTRB[Float32](1, 2, 3, 5)
	if !r1.Equal(r2) {
		t.Errorf("Equal: should be true for identical radii")
	}
	if r1.Equal(r3) {
		t.Errorf("Equal: should be false for different radii")
	}
}

func TestRoundingRadii_String(t *testing.T) {
	r := NewRoundingRadii[Float32](1)
	s := r.String()
	if s == "" {
		t.Errorf("String: should not be empty")
	}
}
