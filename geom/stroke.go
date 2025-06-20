package geom

// StrokeCap specifies the style of the stroke's end caps.
type StrokeCap uint8

const (
	StrokeCapButt   StrokeCap = iota // Stroke ends with a flat edge at the endpoint.
	StrokeCapRound                   // Stroke ends with a half-circle centered on the endpoint.
	StrokeCapSquare                  // Stroke ends with a square extending beyond the endpoint.
)

// StrokeJoin specifies how two connected segments are joined.
type StrokeJoin uint8

const (
	StrokeJoinMiter StrokeJoin = iota // Segments are joined with a sharp or extended corner.
	StrokeJoinRound                   // Segments are joined with a rounded corner.
	StrokeJoinBevel                   // Segments are joined with a beveled (flattened) corner.
)

// StrokeStyle describes how to render the outline of a path or shape.
// The zero value is valid and represents a stroke with zero width and default join/cap styles.
type StrokeStyle[T TScalar] struct {
	Width      T          // Stroke width.
	Cap        StrokeCap  // End cap style.
	Join       StrokeJoin // Join style between segments.
	MiterLimit T          // Miter join limit; controls when a miter is replaced by a bevel.
}

// Equal reports whether s and other have identical stroke parameters.
func (s StrokeStyle[T]) Equal(other StrokeStyle[T]) bool {
	return s.Width == other.Width &&
		s.Cap == other.Cap &&
		s.Join == other.Join &&
		s.MiterLimit == other.MiterLimit
}
