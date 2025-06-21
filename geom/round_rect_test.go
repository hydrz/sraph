package geom

import "testing"

func TestRoundRect_IsRect(t *testing.T) {
	tests := []struct {
		name  string
		rect  Rect
		radii RoundingRadii
		want  bool
	}{
		{"Zero radii", NewRect(0, 0, 10, 10), NewRoundingRadii(0), true},
		{"Non-zero radii", NewRect(0, 0, 10, 10), NewRoundingRadii(2), false},
		{"Empty rect", NewRect(0, 0, 0, 0), NewRoundingRadii(0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := NewRoundRect(tt.rect, tt.radii)
			if got := rr.IsRect(); got != tt.want {
				t.Errorf("IsRect: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundRect_IsOval(t *testing.T) {
	tests := []struct {
		name  string
		rect  Rect
		radii RoundingRadii
		want  bool
	}{
		{"Oval", NewRect(0, 0, 10, 10), NewRoundingRadii(5), true},
		{"Not oval", NewRect(0, 0, 10, 10), NewRoundingRadii(2), false},
		{"Empty rect", NewRect(0, 0, 0, 0), NewRoundingRadii(0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := NewRoundRect(tt.rect, tt.radii)
			if got := rr.IsOval(); got != tt.want {
				t.Errorf("IsOval: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundRect_Contains(t *testing.T) {
	rect := NewRect(0, 0, 10, 10)
	radii := NewRoundingRadii(2)
	rr := NewRoundRect(rect, radii)
	tests := []struct {
		name  string
		point Point
		want  bool
	}{
		{"Inside", NewPoint(5, 5), true},
		{"Outside", NewPoint(20, 20), false},
		{"On corner", NewPoint(2, 2), true},
		{"On rounded edge", NewPoint(2, 0), true},
		{"On straight edge", NewPoint(5, 0), true},
		{"Outside rounded edge", NewPoint(0, 0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rr.Contains(tt.point); got != tt.want {
				t.Errorf("Contains(%v): got %v, want %v", tt.point, got, tt.want)
			}
		})
	}
}

func TestRoundRect_Dispatch(t *testing.T) {
	rect := NewRect(0, 0, 10, 10)
	radii := NewRoundingRadii(2)
	rr := NewRoundRect(rect, radii)
	receiver := &testPathReceiver{}
	rr.Dispatch(receiver, true)
	if len(receiver.ops) == 0 {
		t.Errorf("Dispatch: should call receiver methods")
	}
}

func TestRoundRectPathSource(t *testing.T) {
	rect := NewRect(0, 0, 10, 10)
	radii := NewRoundingRadii(2)
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
	receiver := &testPathReceiver{}
	src.Dispatch(receiver)
	if len(receiver.ops) == 0 {
		t.Errorf("RoundRectPathSource: Dispatch should call receiver")
	}
}

func TestDiffRoundRectPathSource(t *testing.T) {
	rect := NewRect(0, 0, 10, 10)
	radii1 := NewRoundingRadii(2)
	radii2 := NewRoundingRadii(1)
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
	receiver := &testPathReceiver{}
	src.Dispatch(receiver)
	if len(receiver.ops) == 0 {
		t.Errorf("DiffRoundRectPathSource: Dispatch should call receiver")
	}
}
