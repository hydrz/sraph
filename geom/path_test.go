package geom

// testPathReceiver is a PathReceiver for test assertions.
type testPathReceiver struct {
	ops      []string
	moves    []Point
	lines    []Point
	quads    [][2]Point
	conics   [][3]interface{}
	cubics   [][3]Point
	closed   int
	pathEnds int
}

func (d *testPathReceiver) MoveTo(p2 Point, willBeClosed bool) {
	d.ops = append(d.ops, "MoveTo")
	d.moves = append(d.moves, p2)
}
func (d *testPathReceiver) LineTo(p2 Point) {
	d.ops = append(d.ops, "LineTo")
	d.lines = append(d.lines, p2)
}
func (d *testPathReceiver) QuadTo(cp, p2 Point) {
	d.ops = append(d.ops, "QuadTo")
	d.quads = append(d.quads, [2]Point{cp, p2})
}
func (d *testPathReceiver) ConicTo(cp, p2 Point, weight float64) bool {
	d.ops = append(d.ops, "ConicTo")
	d.conics = append(d.conics, [3]interface{}{cp, p2, weight})
	return true
}
func (d *testPathReceiver) CubicTo(cp1, cp2, p2 Point) {
	d.ops = append(d.ops, "CubicTo")
	d.cubics = append(d.cubics, [3]Point{cp1, cp2, p2})
}
func (d *testPathReceiver) Close() {
	d.ops = append(d.ops, "Close")
	d.closed++
}
func (d *testPathReceiver) PathEnd() {
	d.ops = append(d.ops, "PathEnd")
	d.pathEnds++
}
