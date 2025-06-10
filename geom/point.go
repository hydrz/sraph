package geom

import (
	"math"
)

// Point represents a 2D point with X and Y coordinates.
// The interface provides a set of arithmetic, comparison, and utility operations for manipulating and querying point values.
// All operations return a new point and do not modify the original point.
type Point[T Scalar] interface {
	// X returns the X coordinate.
	X() T
	// Y returns the Y coordinate.
	Y() T

	// Add returns the sum of this point and ano.
	Add(o Point[T]) Point[T]
	// Sub returns the difference of this point and ano.
	Sub(o Point[T]) Point[T]
	// Mul returns the element-wise product of this point and ano.
	Mul(o Point[T]) Point[T]
	// Div returns the element-wise division of this point by ano.
	Div(o Point[T]) Point[T]

	// Neg returns the negated point.
	Neg() Point[T]
	// Abs returns the point with absolute values.
	Abs() Point[T]
	// Floor returns the point with floor applied to each coordinate.
	Floor() Point[T]
	// Ceil returns the point with ceil applied to each coordinate.
	Ceil() Point[T]
	// Round returns the point with round applied to each coordinate.
	Round() Point[T]

	// Equal returns true if this point equals ano.
	Equal(o Point[T]) bool
	// Min returns the element-wise minimum with ano point.
	Min(o ...Point[T]) Point[T]
	// Max returns the element-wise maximum with ano point.
	Max(o ...Point[T]) Point[T]

	// IsFinite returns true if both coordinates are finite.
	IsFinite() bool

	// IsZero returns true if both coordinates are zero.
	IsZero() bool

	// Translate returns a new point translated by the given vector.
	Translate(vector Point[T]) Point[T]

	// Scalar returns the point scaled by a scalar.
	Scale(scale T) Point[T]

	// Rotate rotates this point around the origin by the given angle in radians.
	Rotate(angle Radians) Point[T]

	// Normalize returns a normalized version of this point.
	//
	// Normalization scales the point to have a length of 1, while maintaining its direction.
	// This is useful in many applications, such as graphics and physics, where you need a unit vector.
	// If the point is zero, it returns a default unit vector (1, 0).
	// Formula: normalized = (x / length, y / length).
	Normalize() Point[T]

	// Dot calculates the dot product of this point and ano.
	//
	// The dot product is used to determine the angle relationship between two vectors (orthogonal, same direction, opposite direction, etc.).
	// If the dot product is 0, the two vectors are orthogonal.
	// Formula: dot = x1*x2 + y1*y2.
	Dot(o Point[T]) T

	// Cross calculates the cross product of this point and ano.
	//
	// The cross product is used to determine the relative direction of two vectors (clockwise/counterclockwise).
	// A positive result means counterclockwise, negative means clockwise, and zero means collinear.
	// Formula: cross = x1*y2 - y1*x2.
	Cross(o Point[T]) T

	// Distance returns the Euclidean distance between this point and ano.
	//
	// Returns the standard distance between two points, in the same unit as the coordinates.
	// Formula: distance = sqrt((x2 - x1)^2 + (y2 - y1)^2).
	Distance(o Point[T]) T

	// DistanceSquared returns the squared Euclidean distance between this point and ano.
	//
	// Commonly used for distance comparisons to avoid the square root operation and improve efficiency.
	// Formula: distanceSquared = (x2 - x1)^2 + (y2 - y1)^2.
	DistanceSquared(o Point[T]) T

	// Length returns the distance from this point to the origin (0,0), i.e., the vector magnitude.
	//
	// Returns the magnitude (norm) of this point (vector).
	// Formula: length = sqrt(x^2 + y^2).
	Length() T

	// LengthSquared returns the squared distance from this point to the origin (0,0).
	//
	// Commonly used to compare vector magnitudes without needing a square root.
	// Formula: lengthSquared = x^2 + y^2.
	LengthSquared() T

	// Reflect calculates the reflection vector of this point (vector) about the specified axis (vector).
	//
	// Commonly used in physics simulations (such as light reflection, collision bounce, etc.)
	// Formula: v' = v - 2 * (v·n) * n, where v is the original vector and n is the normal (axis).
	// For example, if you want to reflect a point across the x-axis, you would use Point(0, 1) as the axis.
	Reflect(axis Point[T]) Point[T]

	// Lerp performs linear interpolation between this point and ano.
	//
	// Used for animation, smooth transitions, etc.
	// Formula: lerp = this + (o - this) * t, where t is the interpolation factor (0 to 1).
	Lerp(o Point[T], t T) Point[T]

	// String returns a string representation of the point. like "Point(x, y)".
	String() string
}

func NewPoint[T Scalar](x, y T) Point[T] {
	return &point[T]{x: x, y: y}
}

type point[T Scalar] struct {
	x T
	y T
}

// X implements Point.
func (p point[T]) X() T {
	return p.x
}

// Y implements Point.
func (p point[T]) Y() T {
	return p.y
}

// Add implements Point.
func (p point[T]) Add(o Point[T]) Point[T] {
	return &point[T]{x: p.x + o.X(), y: p.y + o.Y()}
}

// Sub implements Point.
func (p point[T]) Sub(o Point[T]) Point[T] {
	return &point[T]{x: p.x - o.X(), y: p.y - o.Y()}
}

// Mul implements Point.
func (p point[T]) Mul(o Point[T]) Point[T] {
	return &point[T]{x: p.x * o.X(), y: p.y * o.Y()}
}

// Div implements Point.
func (p point[T]) Div(o Point[T]) Point[T] {
	return &point[T]{x: p.x / o.X(), y: p.y / o.Y()}
}

// Neg implements Point.
func (p point[T]) Neg() Point[T] {
	return &point[T]{x: -p.x, y: -p.y}
}

// Abs implements Point.
func (p point[T]) Abs() Point[T] {
	var zero T
	x := p.x
	y := p.y
	if x < zero {
		x = -x
	}
	if y < zero {
		y = -y
	}
	return &point[T]{x: x, y: y}
}

// Floor implements Point.
func (p point[T]) Floor() Point[T] {
	return &point[T]{
		x: T(math.Floor(float64(p.x))),
		y: T(math.Floor(float64(p.y))),
	}
}

// Ceil implements Point.
func (p point[T]) Ceil() Point[T] {
	return &point[T]{
		x: T(math.Ceil(float64(p.x))),
		y: T(math.Ceil(float64(p.y))),
	}
}

// Round implements Point.
func (p point[T]) Round() Point[T] {
	return &point[T]{
		x: T(math.Round(float64(p.x))),
		y: T(math.Round(float64(p.y))),
	}
}

// Equal implements Point.
func (p point[T]) Equal(o Point[T]) bool {
	return Equal(p.x, o.X()) && Equal(p.y, o.Y())
}

// Min implements Point.
func (p point[T]) Min(o ...Point[T]) Point[T] {
	// Returns the element-wise minimum with o points.
	if len(o) == 0 {
		return p
	}
	minX := p.x
	minY := p.y
	for _, o := range o {
		if o.X() < minX {
			minX = o.X()
		}
		if o.Y() < minY {
			minY = o.Y()
		}
	}
	return &point[T]{x: minX, y: minY}
}

// Max implements Point.
func (p point[T]) Max(o ...Point[T]) Point[T] {
	// Returns the element-wise maximum with o points.
	if len(o) == 0 {
		return p
	}
	maxX := p.x
	maxY := p.y
	for _, o := range o {
		if o.X() > maxX {
			maxX = o.X()
		}
		if o.Y() > maxY {
			maxY = o.Y()
		}
	}
	return &point[T]{x: maxX, y: maxY}
}

// IsFinite implements Point.
func (p point[T]) IsFinite() bool {
	return IsFinite(p.x) && IsFinite(p.y)
}

// IsZero implements Point.
func (p point[T]) IsZero() bool {
	var zero T
	return p.x == zero && p.y == zero
}

// Translate implements Point.
func (p point[T]) Translate(vector Vector2[T]) Point[T] {
	return &point[T]{x: p.x + vector.X(), y: p.y + vector.Y()}
}

// Scale implements Point.
func (p point[T]) Scale(scale T) Point[T] {
	return &point[T]{x: p.x * scale, y: p.y * scale}
}

// Rotate implements Point.
func (p point[T]) Rotate(angle Radians) Point[T] {
	cos := T(math.Cos(float64(angle)))
	sin := T(math.Sin(float64(angle)))
	return &point[T]{
		x: p.x*cos - p.y*sin,
		y: p.x*sin + p.y*cos,
	}
}

// Normalize implements Point.
func (p point[T]) Normalize() Point[T] {
	len := p.Length()
	var zero T
	if len == zero {
		return &point[T]{x: 1, y: 0}
	}
	return &point[T]{x: p.x / len, y: p.y / len}
}

// Dot implements Point.
func (p point[T]) Dot(o Point[T]) T {
	return p.x*o.X() + p.y*o.Y()
}

// Cross implements Point.
func (p point[T]) Cross(o Point[T]) T {
	return p.x*o.Y() - p.y*o.X()
}

// Distance implements Point.
func (p point[T]) Distance(o Point[T]) T {
	dx := p.x - o.X()
	dy := p.y - o.Y()
	return T(math.Sqrt(float64(dx*dx + dy*dy)))
}

// DistanceSquared implements Point.
func (p point[T]) DistanceSquared(o Point[T]) T {
	dx := p.x - o.X()
	dy := p.y - o.Y()
	return dx*dx + dy*dy
}

// Length implements Point.
func (p point[T]) Length() T {
	return T(math.Sqrt(float64(p.x*p.x + p.y*p.y)))
}

// LengthSquared implements Point.
func (p point[T]) LengthSquared() T {
	return p.x*p.x + p.y*p.y
}

// Reflect implements Point.
func (p point[T]) Reflect(axis Point[T]) Point[T] {
	axis = axis.Normalize()
	dot := p.Dot(axis)
	return p.Sub(axis.Scale(dot).Scale(2))
}

// Lerp implements Point.
func (p point[T]) Lerp(o Point[T], t T) Point[T] {
	return &point[T]{
		x: p.x + (o.X()-p.x)*t,
		y: p.y + (o.Y()-p.y)*t,
	}
}

// String implements Point.
func (p point[T]) String() string {
	return "Point(" + p.x.String() + ", " + p.y.String() + ")"
}
