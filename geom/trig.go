package geom

import "math"

// Trig represents the cosine and sine of an angle.
// Each Trig value can be used for efficient repeated 2D rotations.
// The zero value is valid and represents a zero angle (Cos=1, Sin=0).
type Trig[T TScalar] struct {
	Cos T // Cosine of the angle.
	Sin T // Sine of the angle.
}

// NewTrig returns a Trig representing the cosine and sine of the given angle in radians.
func NewTrig[T TScalar](r Radians) Trig[T] {
	return Trig[T]{
		Cos: T(math.Cos(ToFloat64(r))),
		Sin: T(math.Sin(ToFloat64(r))),
	}
}

// Rotate returns the vector v rotated by the angle represented by t.
func (t Trig[T]) Rotate(v Vector2[T]) Vector2[T] {
	return Vector2[T]{
		X: v.X*t.Cos - v.Y*t.Sin,
		Y: v.X*t.Sin + v.Y*t.Cos,
	}
}

// Neg returns the Trig representing the negative of the current angle.
func (t Trig[T]) Neg() Trig[T] {
	return Trig[T]{
		Cos: t.Cos,
		Sin: -t.Sin,
	}
}

// CirclePoint returns the point on a circle of the given radius at the angle represented by t.
func (t Trig[T]) CirclePoint(radius T) Point[T] {
	return Point[T]{
		t.Cos * radius,
		t.Sin * radius,
	}
}

// EllipsePoint returns the point on an ellipse with the given radii at the angle represented by t.
func (t Trig[T]) EllipsePoint(ellipseRadii Size[T]) Point[T] {
	return Point[T]{
		t.Cos * ellipseRadii.Width,
		t.Sin * ellipseRadii.Height,
	}
}
