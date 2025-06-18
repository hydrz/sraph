package geom

type Shear[T Scalar] struct {
	XY T
	XZ T
	YZ T
}

func (s Shear[T]) Eq(o Shear[T]) bool {
	return Eq(s.XY, o.XY) && Eq(s.XZ, o.XZ) && Eq(s.YZ, o.YZ)
}

func (s Shear[T]) String() string {
	return "(" + s.XY.String() + ", " + s.XZ.String() + ", " + s.YZ.String() + ")"
}
