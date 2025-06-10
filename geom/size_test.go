package geom

import (
	"math"
	"testing"
)

func TestSizeBasicOps(t *testing.T) {
	a := NewSize(F32(3), F32(4))
	b := NewSize(F32(1), F32(2))

	// Add
	c := a.Add(b)
	if !c.Equal(NewSize(F32(4), F32(6))) {
		t.Errorf("Add failed: got %v", c)
	}

	// Sub
	c = a.Sub(b)
	if !c.Equal(NewSize(F32(2), F32(2))) {
		t.Errorf("Sub failed: got %v", c)
	}

	// Mul
	c = a.Mul(b)
	if !c.Equal(NewSize(F32(3), F32(8))) {
		t.Errorf("Mul failed: got %v", c)
	}

	// MulScalar
	c = a.MulScalar(F32(2))
	if !c.Equal(NewSize(F32(6), F32(8))) {
		t.Errorf("MulScalar failed: got %v", c)
	}

	// Div
	c = a.Div(b)
	if !c.Equal(NewSize(F32(3), F32(2))) {
		t.Errorf("Div failed: got %v", c)
	}

	// Neg
	c = a.Neg()
	if !c.Equal(NewSize(F32(-3), F32(-4))) {
		t.Errorf("Neg failed: got %v", c)
	}

	// Scale
	c = a.Scale(F32(2), F32(3))
	if !c.Equal(NewSize(F32(6), F32(12))) {
		t.Errorf("Scale failed: got %v", c)
	}
}

func TestSizeMinMaxDimension(t *testing.T) {
	a := NewSize(F32(3), F32(4))
	if a.MinDimension() != F32(3) {
		t.Errorf("MinDimension failed: got %v", a.MinDimension())
	}
	if a.MaxDimension() != F32(4) {
		t.Errorf("MaxDimension failed: got %v", a.MaxDimension())
	}
}

func TestSizeAreaAbs(t *testing.T) {
	a := NewSize(F32(-3), F32(4))
	if a.Area() != F32(-12) {
		t.Errorf("Area failed: got %v", a.Area())
	}
	abs := a.Abs()
	if !abs.Equal(NewSize(F32(3), F32(4))) {
		t.Errorf("Abs failed: got %v", abs)
	}
}

func TestSizeFloorCeilRound(t *testing.T) {
	a := NewSize(F32(3.7), F32(-4.3))
	if !a.Floor().Equal(NewSize(F32(3), F32(-5))) {
		t.Errorf("Floor failed: got %v", a.Floor())
	}
	if !a.Ceil().Equal(NewSize(F32(4), F32(-4))) {
		t.Errorf("Ceil failed: got %v", a.Ceil())
	}
	if !a.Round().Equal(NewSize(F32(4), F32(-4))) {
		t.Errorf("Round failed: got %v", a.Round())
	}
}

func TestSizeIsZeroIsSquare(t *testing.T) {
	a := NewSize(F32(0), F32(0))
	if !a.IsZero() {
		t.Errorf("IsZero failed")
	}
	b := NewSize(F32(2), F32(2))
	if !b.IsSquare() {
		t.Errorf("IsSquare failed")
	}
	c := NewSize(F32(2), F32(3))
	if c.IsSquare() {
		t.Errorf("IsSquare should be false")
	}
}

func TestSizeIsFiniteIsInfinite(t *testing.T) {
	a := NewSize(F32(1), F32(2))
	if !a.IsFinite() {
		t.Errorf("IsFinite failed")
	}
	b := NewSize(F32(float32(math.Inf(1))), F32(2))
	if !b.IsInfinite() {
		t.Errorf("IsInfinite failed")
	}
}

func TestSizeMinMax(t *testing.T) {
	a := NewSize(F32(3), F32(4))
	b := NewSize(F32(2), F32(5))
	c := NewSize(F32(4), F32(1))
	min := a.Min(b, c)
	max := a.Max(b, c)
	if !min.Equal(NewSize(F32(2), F32(1))) {
		t.Errorf("Min failed: got %v", min)
	}
	if !max.Equal(NewSize(F32(4), F32(5))) {
		t.Errorf("Max failed: got %v", max)
	}
}

func TestSizeMipCount(t *testing.T) {
	a := NewSize(F32(8), F32(4))
	if a.MipCount() != 4 {
		t.Errorf("MipCount failed: got %v", a.MipCount())
	}
	b := NewSize(F32(1), F32(1))
	if b.MipCount() != 1 {
		t.Errorf("MipCount failed for 1x1: got %v", b.MipCount())
	}
}

func TestSizeString(t *testing.T) {
	a := NewSize(F32(3), F32(4))
	str := a.String()
	if str != "Size(3, 4)" {
		t.Errorf("String failed: got %v", str)
	}
}
