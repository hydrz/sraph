package geom

var _ PathReceiver[F32] = (*testPathReceiver[F32])(nil)

// testPathReceiver is a generic PathReceiver for test assertions.
type testPathReceiver[T Scalar] struct {
	ops      []string
	moves    []Point[T]
	lines    []Point[T]
	quads    [][2]Point[T]
	conics   [][3]interface{}
	cubics   [][3]Point[T]
	closed   int
	pathEnds int
}

func (d *testPathReceiver[T]) MoveTo(p2 Point[T], willBeClosed bool) {
	d.ops = append(d.ops, "MoveTo")
	d.moves = append(d.moves, p2)
}
func (d *testPathReceiver[T]) LineTo(p2 Point[T]) {
	d.ops = append(d.ops, "LineTo")
	d.lines = append(d.lines, p2)
}
func (d *testPathReceiver[T]) QuadTo(cp, p2 Point[T]) {
	d.ops = append(d.ops, "QuadTo")
	d.quads = append(d.quads, [2]Point[T]{cp, p2})
}
func (d *testPathReceiver[T]) ConicTo(cp, p2 Point[T], weight float64) bool {
	d.ops = append(d.ops, "ConicTo")
	d.conics = append(d.conics, [3]interface{}{cp, p2, weight})
	return true
}
func (d *testPathReceiver[T]) CubicTo(cp1, cp2, p2 Point[T]) {
	d.ops = append(d.ops, "CubicTo")
	d.cubics = append(d.cubics, [3]Point[T]{cp1, cp2, p2})
}
func (d *testPathReceiver[T]) Close() {
	d.ops = append(d.ops, "Close")
	d.closed++
}
func (d *testPathReceiver[T]) PathEnd() {
	d.ops = append(d.ops, "PathEnd")
	d.pathEnds++
}
