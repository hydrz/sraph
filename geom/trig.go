package geom

import "math"

type Trig[T Scalar] struct {
	Cos T
	Sin T
}

func NewTrig[T Scalar](r Radians[T]) Trig[T] {
	return Trig[T]{
		Cos: T(math.Cos(r.Float64())),
		Sin: T(math.Sin(r.Float64())),
	}
}

// Returns the vector rotated by the represented angle.
func (t Trig[T]) Rotate(v Vector2[T]) Vector2[T] {
	return Vector2[T]{
		X: v.X*t.Cos - v.Y*t.Sin,
		Y: v.X*t.Sin + v.Y*t.Cos,
	}
}

// Returns the Trig representing the negative version of this angle.
func (t Trig[T]) Neg() Trig[T] {
	return Trig[T]{
		Cos: t.Cos,
		Sin: -t.Sin,
	}
}

// Returns the corresponding point on a circle of a given |radius|.
func (t Trig[T]) CirclePoint(radius T) Point[T] {
	return Point[T]{
		t.Cos * radius,
		t.Sin * radius,
	}
}

// Returns the corresponding point on an ellipse with the given size.
func (t Trig[T]) EllipsePoint(ellipseRadii Size[T]) Point[T] {
	return Point[T]{
		t.Cos * ellipseRadii.Width,
		t.Sin * ellipseRadii.Height,
	}
}
