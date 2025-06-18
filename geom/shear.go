package geom

type Shear[T Scalar] struct {
	XY T
	XZ T
	YZ T
}

func (s Shear[T]) Eq(o Shear[T]) bool {
	return NearlyEq(s.XY, o.XY) && NearlyEq(s.XZ, o.XZ) && NearlyEq(s.YZ, o.YZ)
}

func (s Shear[T]) String() string {
	return "(" + s.XY.String() + ", " + s.XZ.String() + ", " + s.YZ.String() + ")"
}
