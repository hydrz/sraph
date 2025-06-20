package geom

type Shear[T Number] struct {
	XY T
	XZ T
	YZ T
}

func (s Shear[T]) Eq(o Shear[T]) bool {
	return NearlyEqual(s.XY, o.XY) && NearlyEqual(s.XZ, o.XZ) && NearlyEqual(s.YZ, o.YZ)
}

func (s Shear[T]) String() string {
	return "(" + ToString(s.XY) + ", " + ToString(s.XZ) + ", " + ToString(s.YZ) + ")"
}
