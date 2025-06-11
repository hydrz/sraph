package geom

// StrokeCap defines the style of the stroke's end caps.
type StrokeCap uint8

const (
	StrokeCapButt   StrokeCap = iota // Flat cap at the end of the stroke.
	StrokeCapRound                   // Semi-circular cap at the end of the stroke.
	StrokeCapSquare                  // Square cap extends beyond the end of the stroke.
)

// StrokeJoin defines the style of the join between two connected segments.
type StrokeJoin uint8

const (
	StrokeJoinMiter StrokeJoin = iota // Sharp corner or extended join.
	StrokeJoinRound                   // Rounded join at the corner.
	StrokeJoinBevel                   // Beveled (flattened) join at the corner.
)

// StrokeStyle describes the parameters for stroking a path or geometric object.
// T is a numeric type that represents the stroke width and miter limit.
type StrokeStyle[T Scalar] struct {
	Width      T          // Stroke width.
	Cap        StrokeCap  // Style of the stroke's end caps.
	Join       StrokeJoin // Style of the join between segments.
	MiterLimit T          // Limit for miter joins; controls when a miter is replaced by a bevel.
}

// Equal returns true if all fields of two StrokeStyle values are equal.
func (s StrokeStyle[T]) Equal(other StrokeStyle[T]) bool {
	return s.Width == other.Width &&
		s.Cap == other.Cap &&
		s.Join == other.Join &&
		s.MiterLimit == other.MiterLimit
}
