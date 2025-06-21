package geom

import "math"

// Trig defines the interface for representing the cosine and sine of an angle.
// Each Trig value can be used for efficient repeated 2D rotations.
// The zero value is valid and represents a zero angle (Cos=1, Sin=0).
type Trig interface {
	// Cos returns the cosine of the angle.
	Cos() Scalar
	// Sin returns the sine of the angle.
	Sin() Scalar
	// Rotate returns the vector v rotated by this angle.
	Rotate(v Vector2) Vector2
	// Neg returns the Trig representing the negative of the current angle.
	Neg() Trig
	// CirclePoint returns the point on a circle of the given radius at this angle.
	CirclePoint(radius Scalar) Point
	// EllipsePoint returns the point on an ellipse with the given radii at this angle.
	EllipsePoint(ellipseRadii Size) Point
}

// trig represents the cosine and sine of an angle.
type trig struct {
	cos, sin Scalar
}

// NewTrig returns a Trig representing the cosine and sine of the given angle in radians.
func NewTrig(r Radians) Trig {
	return trig{
		cos: Scalar(math.Cos(ToFloat64(r))),
		sin: Scalar(math.Sin(ToFloat64(r))),
	}
}

// Cos implements Trig.Cos.
func (t trig) Cos() Scalar { return t.cos }

// Sin implements Trig.Sin.
func (t trig) Sin() Scalar { return t.sin }

// Rotate implements Trig.Rotate.
func (t trig) Rotate(v Vector2) Vector2 {
	return NewVector2(
		v.X()*t.cos-v.Y()*t.sin,
		v.X()*t.sin+v.Y()*t.cos,
	)
}

// Neg implements Trig.Neg.
func (t trig) Neg() Trig {
	return trig{
		cos: t.cos,
		sin: -t.sin,
	}
}

// CirclePoint implements Trig.CirclePoint.
func (t trig) CirclePoint(radius Scalar) Point {
	return NewPoint(
		t.cos*radius,
		t.sin*radius,
	)
}

// EllipsePoint implements Trig.EllipsePoint.
func (t trig) EllipsePoint(ellipseRadii Size) Point {
	return NewPoint(
		t.cos*ellipseRadii.Width(),
		t.sin*ellipseRadii.Height(),
	)
}
