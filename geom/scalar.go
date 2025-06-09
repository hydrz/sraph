package geom

// Scalar represents a quantity that has only magnitude (numerical value) and no direction, commonly used for coordinates and lengths.
// In computer graphics, the following types are typically used to represent scalar values:
// - float32: Used for most graphics computations, providing sufficient precision and performance.
// - float64: Used for calculations requiring higher precision, such as physical simulations and advanced rendering.
// - Half-precision float: Used to reduce memory usage and improve performance in scenarios where high precision is not required.
// - Integer: Used for discrete values such as pixel coordinates and indices.
// - Fixed-point: Used in scenarios that require high precision, such as font rendering.
type Scalar interface {
	// Add returns the sum of the current scalar and another scalar.
	Add(o Scalar) Scalar
	// Sub returns the difference between the current scalar and another scalar.
	Sub(o Scalar) Scalar
	// Mul returns the product of the current scalar and another scalar.
	Mul(o Scalar) Scalar
	// Div returns the quotient of the current scalar divided by another scalar.
	Div(o Scalar) Scalar
	// Neg returns the negated value of the scalar.
	Neg() Scalar

	// Equal returns true if the current scalar is equal to another scalar.
	Equal(o Scalar) bool
	// Less returns true if the current scalar is less than another scalar.
	Less(o Scalar) bool
	// LessEqual returns true if the current scalar is less than or equal to another scalar.
	LessEqual(o Scalar) bool
	// Greater returns true if the current scalar is greater than another scalar.
	Greater(o Scalar) bool
	// GreaterEqual returns true if the current scalar is greater than or equal to another scalar.
	GreaterEqual(o Scalar) bool
	// NotEqual returns true if the current scalar is not equal to another scalar.
	NotEqual(o Scalar) bool

	// Min returns the minimum value among the current scalar and the provided scalars.
	Min(o ...Scalar) Scalar
	// Max returns the maximum value among the current scalar and the provided scalars.
	Max(o ...Scalar) Scalar

	// IsNan returns true if the scalar is NaN (not a number).
	IsNan() bool
	// IsInf returns true if the scalar is infinite.
	IsInf() bool
	// IsZero returns true if the scalar is zero.
	IsZero() bool
	// IsFinite returns true if the scalar is a finite number (not NaN or infinite).
	IsFinite() bool
	// Clamp returns the scalar clamped between the given minimum and maximum values.
	Clamp(min, max Scalar) Scalar

	// Abs returns the absolute value of the scalar.
	Abs() Scalar
	// Pow returns the current scalar raised to the power of another scalar.
	Pow(o Scalar) Scalar
	// Sqrt returns the square root of the scalar.
	Sqrt() Scalar

	// Floor returns the greatest integer value less than or equal to the scalar.
	Floor() Scalar
	// Ceil returns the smallest integer value greater than or equal to the scalar.
	Ceil() Scalar
	// Round returns the nearest integer to the scalar.
	Round() Scalar

	// Sin returns the sine of the scalar (in radians).
	Sin() Scalar
	// Cos returns the cosine of the scalar (in radians).
	Cos() Scalar
	// Tan returns the tangent of the scalar (in radians).
	Tan() Scalar
	// Acos returns the arccosine of the scalar (in radians).
	Acos() Scalar
	// Atan returns the arctangent of the scalar (in radians).
	Atan() Scalar
	// Atan2 returns the arctangent of the quotient of the current scalar and another scalar.
	Atan2(o Scalar) Scalar

	// String returns the string representation of the scalar.
	String() string
}
