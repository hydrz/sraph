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

/*
*

	explicit Trig(Radians r)
	    : cos(std::cos(r.radians)), sin(std::sin(r.radians)) {}

	/// Construct a Trig object from the given cosine and sine values.
	Trig(double cos, double sin) : cos(cos), sin(sin) {}

	double cos;
	double sin;

	/// @brief  Returns the vector rotated by the represented angle.
	Vector2 operator*(const Vector2& vector) const {
	  return Vector2(static_cast<Scalar>(vector.x * cos - vector.y * sin),
	                 static_cast<Scalar>(vector.x * sin + vector.y * cos));
	}

	/// @brief  Returns the Trig representing the negative version of this angle.
	Trig operator-() const { return Trig(cos, -sin); }

	/// @brief  Returns the corresponding point on a circle of a given |radius|.
	Vector2 operator*(double radius) const {
	  return Vector2(static_cast<Scalar>(cos * radius),
	                 static_cast<Scalar>(sin * radius));
	}

	/// @brief  Returns the corresponding point on an ellipse with the given size.
	Vector2 operator*(const Size& ellipse_radii) const {
	  return Vector2(static_cast<Scalar>(cos * ellipse_radii.width),
	                 static_cast<Scalar>(sin * ellipse_radii.height));
	}

*
*/
func (t Trig[T]) Rotate(v Vector2[T]) Vector2[T] {
	return NewVector2[T](
		v.X()*t.cos-v.Y()*t.sin,
		v.X()*t.sin+v.Y()*t.cos,
	)
}
