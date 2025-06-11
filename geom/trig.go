package geom

import "math"

type Trig[T Scalar] struct {
	cos T
	sin T
}

func NewTrig[T Scalar](r Radians[T]) Trig[T] {
	return Trig[T]{
		cos: T(math.Cos(r.Float64())),
		sin: T(math.Sin(r.Float64())),
	}
}

func NewTrigCosSin[T Scalar](cos, sin float64) Trig[T] {
	return Trig[T]{
		cos: T(cos),
		sin: T(sin),
	}
}

// Returns the vector rotated by the represented angle.
func (t Trig[T]) Rotate(v Vector2[T]) Vector2[T] {
	return NewVector2(
		v.X()*t.cos-v.Y()*t.sin,
		v.X()*t.sin+v.Y()*t.cos,
	)
}

// Returns the Trig representing the negative version of this angle.
func (t Trig[T]) Neg() Trig[T] {
	return Trig[T]{
		cos: t.cos,
		sin: -t.sin,
	}
}

// Returns the corresponding point on a circle of a given |radius|.
func (t Trig[T]) CirclePoint(radius T) Point[T] {
	return NewPoint(
		t.cos*radius,
		t.sin*radius,
	)
}

// Returns the corresponding point on an ellipse with the given size.
func (t Trig[T]) EllipsePoint(ellipseRadii Size[T]) Point[T] {
	return NewPoint(
		t.cos*ellipseRadii.Width(),
		t.sin*ellipseRadii.Height(),
	)
}
