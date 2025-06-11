package geom

import (
	"testing"
)

func TestRect_NewRectAndBasicProperties(t *testing.T) {
	r := NewRect[Int32](1, 2, 5, 6)
	if r.Left() != 1 || r.Top() != 2 || r.Right() != 5 || r.Bottom() != 6 {
		t.Errorf("Rect basic edges failed: got (%v, %v, %v, %v)", r.Left(), r.Top(), r.Right(), r.Bottom())
	}
	if r.X() != 1 || r.Y() != 2 {
		t.Errorf("Rect origin failed: got (%v, %v)", r.X(), r.Y())
	}
	if r.Width() != 4 || r.Height() != 4 {
		t.Errorf("Rect size failed: got (%v, %v)", r.Width(), r.Height())
	}
}

func TestRect_NewRectXYWH(t *testing.T) {
	r := NewRectXYWH[Int32](2, 3, 4, 5)
	if r.Left() != 2 || r.Top() != 3 || r.Right() != 6 || r.Bottom() != 8 {
		t.Errorf("RectXYWH failed: got (%v, %v, %v, %v)", r.Left(), r.Top(), r.Right(), r.Bottom())
	}
	r2 := NewRectXYWH[Int32](2, 3, -4, -5)
	if r2.Left() != -2 || r2.Top() != -2 || r2.Right() != 2 || r2.Bottom() != 3 {
		t.Errorf("RectXYWH negative failed: got (%v, %v, %v, %v)", r2.Left(), r2.Top(), r2.Right(), r2.Bottom())
	}
}

func TestRect_NewRectOriginSize(t *testing.T) {
	origin := NewPoint(Int32(1), Int32(2))
	size := NewSize(Int32(3), Int32(4))
	r := NewRectOriginSize(origin, size)
	if r.Left() != 1 || r.Top() != 2 || r.Right() != 4 || r.Bottom() != 6 {
		t.Errorf("RectOriginSize failed: got (%v, %v, %v, %v)", r.Left(), r.Top(), r.Right(), r.Bottom())
	}
}

func TestRect_NewRectSize(t *testing.T) {
	size := NewSize(Int32(3), Int32(4))
	r := NewRectSize(size)
	if r.Left() != 0 || r.Top() != 0 || r.Right() != 3 || r.Bottom() != 4 {
		t.Errorf("RectSize failed: got (%v, %v, %v, %v)", r.Left(), r.Top(), r.Right(), r.Bottom())
	}
}

func TestRect_BoundingRect(t *testing.T) {
	p1 := NewPoint(Int32(1), Int32(2))
	p2 := NewPoint(Int32(3), Int32(5))
	p3 := NewPoint(Int32(-1), Int32(4))
	r := BoundingRect(p1, p2, p3)
	if r.Left() != -1 || r.Top() != 2 || r.Right() != 3 || r.Bottom() != 5 {
		t.Errorf("BoundingRect failed: got (%v, %v, %v, %v)", r.Left(), r.Top(), r.Right(), r.Bottom())
	}
	empty := BoundingRect[Int32]()
	if !empty.IsEmpty() {
		t.Errorf("BoundingRect empty failed")
	}
}

func TestRect_AreaAndCenter(t *testing.T) {
	r := NewRect[Int32](0, 0, 4, 4)
	if r.Area() != 16 {
		t.Errorf("Area failed: got %v", r.Area())
	}
	center := r.Center()
	if center.X() != 2 || center.Y() != 2 {
		t.Errorf("Center failed: got (%v, %v)", center.X(), center.Y())
	}
}

func TestRect_Positive(t *testing.T) {
	r := NewRect[Int32](5, 6, 1, 2)
	pos := r.Positive()
	if pos.Left() != 1 || pos.Top() != 2 || pos.Right() != 5 || pos.Bottom() != 6 {
		t.Errorf("Positive failed: got (%v, %v, %v, %v)", pos.Left(), pos.Top(), pos.Right(), pos.Bottom())
	}
}

func TestRect_IsEmptyIsSquareIsMaximum(t *testing.T) {
	r := NewRect[Int32](0, 0, 0, 0)
	if !r.IsEmpty() {
		t.Errorf("IsEmpty failed")
	}
	sq := NewRect[Int32](1, 2, 5, 6)
	if !sq.IsSquare() {
		t.Errorf("IsSquare failed")
	}
	max := NewRectMax[Int32]()
	if !max.IsMaximum() {
		t.Errorf("IsMaximum failed")
	}
}

func TestRect_ContainsAndIntersects(t *testing.T) {
	r := NewRect[Int32](0, 0, 10, 10)
	p := NewPoint(Int32(5), Int32(5))
	if !r.Contains(p) {
		t.Errorf("Contains failed")
	}
	if !r.ContainsExclusive(NewPoint(Int32(1), Int32(1))) {
		t.Errorf("ContainsExclusive failed")
	}
	r2 := NewRect[Int32](5, 5, 15, 15)
	if !r.Intersects(r2) {
		t.Errorf("Intersects failed")
	}
	if !r.ContainsRect(NewRect[Int32](2, 2, 8, 8)) {
		t.Errorf("ContainsRect failed")
	}
}

func TestRect_IntersectionUnion(t *testing.T) {
	r1 := NewRect[Int32](0, 0, 10, 10)
	r2 := NewRect[Int32](5, 5, 15, 15)
	inter, ok := r1.Intersection(r2)
	if !ok || inter.Left() != 5 || inter.Top() != 5 || inter.Right() != 10 || inter.Bottom() != 10 {
		t.Errorf("Intersection failed: got %v, ok=%v", inter, ok)
	}
	union := r1.Union(r2)
	if union.Left() != 0 || union.Top() != 0 || union.Right() != 15 || union.Bottom() != 15 {
		t.Errorf("Union failed: got %v", union)
	}
}

func TestRect_ExpandAndExpandPoint(t *testing.T) {
	r := NewRect[Int32](1, 1, 3, 3)
	r2 := r.Expand(1)
	if r2.Left() != 0 || r2.Top() != 0 || r2.Right() != 4 || r2.Bottom() != 4 {
		t.Errorf("Expand failed: got %v", r2)
	}
	p := NewPoint(Int32(5), Int32(5))
	r3 := r.ExpandPoint(p)
	if r3.Right() != 5 || r3.Bottom() != 5 {
		t.Errorf("ExpandPoint failed: got %v", r3)
	}
}

func TestRect_TranslateAndScale(t *testing.T) {
	r := NewRect[Int32](1, 2, 3, 4)
	v := NewVector2(Int32(2), Int32(3))
	r2 := r.Translate(v)
	if r2.Left() != 3 || r2.Top() != 5 {
		t.Errorf("Translate failed: got %v", r2)
	}
	r3 := r.Scale(2)
	if r3.Left() != 2 || r3.Top() != 4 || r3.Right() != 6 || r3.Bottom() != 8 {
		t.Errorf("Scale failed: got %v", r3)
	}
}

func TestRect_Round(t *testing.T) {
	r := NewRect[Float64](1.2, 2.7, 3.5, 4.9)
	ri := r.Round()
	if ri.Left() != 1 || ri.Top() != 3 || ri.Right() != 4 || ri.Bottom() != 5 {
		t.Errorf("Round failed: got %v", ri)
	}
}

func TestRect_String(t *testing.T) {
	r := NewRect[Int32](1, 2, 3, 4)
	s := r.String()
	if s != "((1, 2) => (3, 4))" {
		t.Errorf("String failed: got %s", s)
	}
}
