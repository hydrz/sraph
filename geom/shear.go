package geom

type Shear[T Scalar] struct {
	XY T
	XZ T
	YZ T
}

func (s Shear[T]) Equal(o Shear[T]) bool {
	return Equal(s.XY, o.XY) && Equal(s.XZ, o.XZ) && Equal(s.YZ, o.YZ)
}

func (s Shear[T]) String() string {
	return "(" + s.XY.String() + ", " + s.XZ.String() + ", " + s.YZ.String() + ")"
}
