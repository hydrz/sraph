package geom

import "testing"

func TestRoundingRadii_NewAndProperties(t *testing.T) {
	r := NewRoundingRadii[F32](5.0)
	if !r.TopLeft.Eq(Size[F32]{5, 5}) ||
		!r.TopRight.Eq(Size[F32]{5, 5}) ||
		!r.BottomLeft.Eq(Size[F32]{5, 5}) ||
		!r.BottomRight.Eq(Size[F32]{5, 5}) {
		t.Errorf("NewRoundingRadii: all corners should be (5,5)")
	}
	if !r.IsUniform() {
		t.Errorf("IsUniform: should be true for all corners Eq")
	}
	if !r.IsFinite() {
		t.Errorf("IsFinite: should be true for finite radii")
	}
	if r.IsEmpty() {
		t.Errorf("IsEmpty: should be false for nonzero radii")
	}
}

func TestRoundingRadiiLTRB(t *testing.T) {
	r := NewRoundingRadiiLTRB[F32](1.0, 2.0, 3.0, 4.0)
	if !r.TopLeft.Eq(Size[F32]{1, 2}) ||
		!r.TopRight.Eq(Size[F32]{3, 2}) ||
		!r.BottomLeft.Eq(Size[F32]{1, 4}) ||
		!r.BottomRight.Eq(Size[F32]{3, 4}) {
		t.Errorf("NewRoundingRadii4: corners not as expected")
	}
	if r.IsUniform() {
		t.Errorf("IsUniform: should be false for unEq corners")
	}
}

func TestRoundingRadii_IsEmpty(t *testing.T) {
	r := NewRoundingRadii[F32](0.0)
	if !r.IsEmpty() {
		t.Errorf("IsEmpty: should be true for all zero radii")
	}
	r2 := NewRoundingRadiiLTRB[F32](0.0, 1.0, 0.0, 0.0)
	if r2.IsEmpty() {
		t.Errorf("IsEmpty: should be false if any corner is nonzero")
	}
}

func TestRoundingRadii_Scale(t *testing.T) {
	r := NewRoundingRadiiLTRB[F32](1.0, 2.0, 3.0, 4.0)
	s := r.Scale(2.0)
	if !s.TopLeft.Eq(Size[F32]{2, 4}) ||
		!s.TopRight.Eq(Size[F32]{6, 4}) ||
		!s.BottomLeft.Eq(Size[F32]{2, 8}) ||
		!s.BottomRight.Eq(Size[F32]{6, 8}) {
		t.Errorf("Scale: scaling not correct")
	}
}

func TestRoundingRadii_ScaleToFit(t *testing.T) {
	r := NewRoundingRadii[F32](10.0)
	bounds := NewRect[F32](0.0, 0.0, 15.0, 20.0)
	scaled := r.ScaleToFit(bounds)
	sumTop := scaled.TopLeft.Width + scaled.TopRight.Width
	sumLeft := scaled.TopLeft.Height + scaled.BottomLeft.Height
	if sumTop > bounds.Width()+1e-5 || sumLeft > bounds.Height()+1e-5 {
		t.Errorf("ScaleToFit: radii sum exceeds bounds")
	}
}

func TestRoundingRadii_Eq(t *testing.T) {
	r1 := NewRoundingRadiiLTRB[F32](1.0, 2.0, 3.0, 4.0)
	r2 := NewRoundingRadiiLTRB[F32](1.0, 2.0, 3.0, 4.0)
	r3 := NewRoundingRadiiLTRB[F32](1.0, 2.0, 3.0, 5.0)
	if !r1.Eq(r2) {
		t.Errorf("Eq: should be true for identical radii")
	}
	if r1.Eq(r3) {
		t.Errorf("Eq: should be false for different radii")
	}
}

func TestRoundingRadii_String(t *testing.T) {
	r := NewRoundingRadii[F32](1.0)
	s := r.String()
	if s == "" {
		t.Errorf("String: should not be empty")
	}
}
