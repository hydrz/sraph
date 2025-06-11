package geom

var _ PathReceiver[Float32] = (*dummyReceiver[Float32])(nil)

type dummyReceiver[T Scalar] struct {
	ops []string
}

func (d *dummyReceiver[T]) MoveTo(p2 Point[T], willBeClosed bool) {
	d.ops = append(d.ops, "MoveTo")
}
func (d *dummyReceiver[T]) LineTo(p2 Point[T]) {
	d.ops = append(d.ops, "LineTo")
}
func (d *dummyReceiver[T]) QuadTo(cp, p2 Point[T]) {
	d.ops = append(d.ops, "QuadTo")
}
func (d *dummyReceiver[T]) ConicTo(cp, p2 Point[T], weight float64) bool {
	d.ops = append(d.ops, "ConicTo")
	return true
}
func (d *dummyReceiver[T]) CubicTo(cp1, cp2, p2 Point[T]) {
	d.ops = append(d.ops, "CubicTo")
}
func (d *dummyReceiver[T]) Close() {
	d.ops = append(d.ops, "Close")
}
func (d *dummyReceiver[T]) PathEnd() {
	d.ops = append(d.ops, "PathEnd")
}
