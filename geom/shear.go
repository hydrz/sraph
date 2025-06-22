package geom

type Shear interface {
	// XY returns the shear factor in the XY plane.
	XY() Scalar
	// XZ returns the shear factor in the XZ plane.
	XZ() Scalar
	// YZ returns the shear factor in the YZ plane.
	YZ() Scalar
	// Equal reports whether this shear and another are equal within floating-point tolerance.
	Equal(o Shear) bool
	// String returns a string representation of the shear factors.
	String() string
}

type shear[T Number] struct {
	xy T
	xz T
	yz T
}

// NewShear creates a new shear with the specified factors in the XY, XZ, and YZ planes.
func NewShear[T Number](xy, xz, yz T) Shear {
	return shear[T]{xy: xy, xz: xz, yz: yz}
}

// XY returns the shear factor in the XY plane.
func (s shear[T]) XY() Scalar {
	return ToScalar(s.xy)
}

// XZ returns the shear factor in the XZ plane.
func (s shear[T]) XZ() Scalar {
	return ToScalar(s.xz)
}

// YZ returns the shear factor in the YZ plane.
func (s shear[T]) YZ() Scalar {
	return ToScalar(s.yz)
}

func (s shear[T]) Equal(o Shear) bool {
	return NearlyEqual(s.xy, T(o.XY())) &&
		NearlyEqual(s.xz, T(o.XZ())) &&
		NearlyEqual(s.yz, T(o.YZ()))
}

func (s shear[T]) String() string {
	return "(" + ToString(s.xy) + ", " + ToString(s.xz) + ", " + ToString(s.yz) + ")"
}
