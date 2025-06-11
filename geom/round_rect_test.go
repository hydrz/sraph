package geom

import "testing"

func TestRoundRect_IsRect(t *testing.T) {
	rect := NewRect[Float32](0, 0, 10, 10)
	radii := NewRoundingRadii[Float32](0)
	rr := NewRoundRect(rect, radii)
	if !rr.IsRect() {
		t.Errorf("IsRect: should be true for zero radii")
	}
}

func TestRoundRect_IsOval(t *testing.T) {
	rect := NewRect[Float32](0, 0, 10, 10)
	radii := NewRoundingRadii[Float32](5)
	rr := NewRoundRect(rect, radii)
	if !rr.IsOval() {
		t.Errorf("IsOval: should be true for uniform radii == half width/height")
	}
}

func TestRoundRect_Contains(t *testing.T) {
	rect := NewRect[Float32](0, 0, 10, 10)
	radii := NewRoundingRadii[Float32](2)
	rr := NewRoundRect(rect, radii)
	inside := NewPoint[Float32](5, 5)
	outside := NewPoint[Float32](20, 20)
	if !rr.Contains(inside) {
		t.Errorf("Contains: should be true for inside point")
	}
	if rr.Contains(outside) {
		t.Errorf("Contains: should be false for outside point")
	}
}

func TestRoundRect_Dispatch(t *testing.T) {
	rect := NewRect[Float32](0, 0, 10, 10)
	radii := NewRoundingRadii[Float32](2)
	rr := NewRoundRect(rect, radii)
	receiver := &dummyReceiver[Float32]{}
	rr.Dispatch(receiver, true)
	if len(receiver.ops) == 0 {
		t.Errorf("Dispatch: should call receiver methods")
	}
}

func TestRoundRectPathSource(t *testing.T) {
	rect := NewRect[Float32](0, 0, 10, 10)
	radii := NewRoundingRadii[Float32](2)
	rr := NewRoundRect(rect, radii)
	src := NewRoundRectPathSource(rr)
	if src.FillType() != FillTypeNonZero {
		t.Errorf("RoundRectPathSource: FillType should be NonZero")
	}
	if !src.IsConvex() {
		t.Errorf("RoundRectPathSource: IsConvex should be true")
	}
	if !src.Bounds().Equal(rect) {
		t.Errorf("RoundRectPathSource: Bounds mismatch")
	}
	receiver := &dummyReceiver[Float32]{}
	src.Dispatch(receiver)
	if len(receiver.ops) == 0 {
		t.Errorf("RoundRectPathSource: Dispatch should call receiver")
	}
}

func TestDiffRoundRectPathSource(t *testing.T) {
	rect := NewRect[Float32](0, 0, 10, 10)
	radii1 := NewRoundingRadii[Float32](2)
	radii2 := NewRoundingRadii[Float32](1)
	rr1 := NewRoundRect(rect, radii1)
	rr2 := NewRoundRect(rect, radii2)
	src := NewDiffRoundRectPathSource(rr1, rr2)
	if src.FillType() != FillTypeEvenOdd {
		t.Errorf("DiffRoundRectPathSource: FillType should be EvenOdd")
	}
	if src.IsConvex() {
		t.Errorf("DiffRoundRectPathSource: IsConvex should be false")
	}
	if !src.Bounds().Equal(rect) {
		t.Errorf("DiffRoundRectPathSource: Bounds mismatch")
	}
	receiver := &dummyReceiver[Float32]{}
	src.Dispatch(receiver)
	if len(receiver.ops) == 0 {
		t.Errorf("DiffRoundRectPathSource: Dispatch should call receiver")
	}
}
