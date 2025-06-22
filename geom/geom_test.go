package geom

import (
	"testing"
)

func TestPointInt26_6_BasicOps(t *testing.T) {
	a := NewPoint[Int26_6](1<<6, 2<<6) // (1.0, 2.0) in fixed-point representation
	b := NewPoint[Int26_6](1<<5, 1<<5) // (0.5, 0.5)

	sum := a.Add(b)
	sumExpected := NewPoint[Int26_6](1<<6+1<<5, 2<<6+1<<5) // (1.0+0.5, 2.0+0.5) = (1.5, 2.5)
	if !sum.Equal(sumExpected) {
		t.Errorf("Add failed: got %v, want %v", sum, sumExpected)
	}

	diff := a.Sub(b)
	diffExpected := NewPoint[Int26_6](1<<6-1<<5, 2<<6-1<<5) // (1.0-0.5, 2.0-0.5) = (0.5, 1.5)
	if !diff.Equal(diffExpected) {
		t.Errorf("Sub failed: got %v, want %v", diff, diffExpected)
	}

	mul := a.Mul(b)
	mulExpected := NewPoint[Int26_6](1<<5, 1<<6) // (1.0*0.5, 2.0*0.5) = (0.5, 1.0)
	if !mul.Equal(mulExpected) {
		t.Errorf("Mul failed: got %v, want %v", mul, mulExpected)
	}

	div := a.Div(b)
	divExpected := NewPoint[Int26_6](2<<6, 4<<6) // (1.0/0.5, 2.0/0.5) = (2.0, 4.0)
	if !div.Equal(divExpected) {
		t.Errorf("Div failed: got %v, want %v", div, divExpected)
	}

	neg := a.Neg()
	negExpected := NewPoint[Int26_6](-(1 << 6), -(2 << 6)) // (-1.0, -2.0)
	if !neg.Equal(negExpected) {
		t.Errorf("Neg failed: got %v, want %v", neg, negExpected)
	}
}

func TestPointInt26_6_LengthAndNormalize(t *testing.T) {
	p := NewPoint[Int26_6](1<<6, 0) // (1.0, 0.0)
	if !NearlyEqual(p.Length(), 1.0) {
		t.Errorf("Length failed: got %v, want 1.0", p.Length())
	}
	n := p.Normalize()
	if !NearlyEqual(n.Length(), 1.0) {
		t.Errorf("Normalize failed: got %v, want length 1.0", n.Length())
	}
}

func TestPointInt26_6_String(t *testing.T) {
	p := NewPoint[Int26_6](1<<6+1<<4, -(1<<6 + 65))
	s := p.String()
	want := "(1:16, -2:01)" // Expected string representation
	if s != want {
		t.Errorf("String failed: got %v, want %v", s, want)
	}
}
