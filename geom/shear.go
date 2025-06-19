package geom

type Shear[T Scalar] struct {
	XY T
	XZ T
	YZ T
}

func (s Shear[T]) Eq(o Shear[T]) bool {
	return ScalarEq(s.XY, o.XY) && ScalarEq(s.XZ, o.XZ) && ScalarEq(s.YZ, o.YZ)
}

func (s Shear[T]) String() string {
	return "(" + s.XY.String() + ", " + s.XZ.String() + ", " + s.YZ.String() + ")"
}
