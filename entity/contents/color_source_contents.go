package contents

// ColorSourceContents is a base interface for contents that provide color sources like gradients and solid colors
type ColorSourceContents interface {
	Contents

	// SetGeometry sets the geometry for the color source
	SetGeometry(geometry Geometry)

	// GetGeometry returns the current geometry
	GetGeometry() Geometry

	// SetInvertColors sets whether to invert colors
	SetInvertColors(invert bool)

	// GetInvertColors returns whether colors are inverted
	GetInvertColors() bool

	// SetOpacity sets the opacity for the color source
	SetOpacity(opacity float32)

	// GetOpacity returns the current opacity
	GetOpacity() float32
}

// ColorSourceContentsImpl provides base implementation for color source contents
type ColorSourceContentsImpl struct {
	*ContentsImpl
	geometry     Geometry
	invertColors bool
	opacity      float32
}

// NewColorSourceContents creates a new color source contents base
func NewColorSourceContents() *ColorSourceContentsImpl {
	return &ColorSourceContentsImpl{
		ContentsImpl: &ContentsImpl{
			isValid: true,
		},
		opacity: 1.0,
	}
}

// SetGeometry sets the geometry for the color source
func (c *ColorSourceContentsImpl) SetGeometry(geometry Geometry) {
	c.geometry = geometry
}

// GetGeometry returns the current geometry
func (c *ColorSourceContentsImpl) GetGeometry() Geometry {
	return c.geometry
}

// SetInvertColors sets whether to invert colors
func (c *ColorSourceContentsImpl) SetInvertColors(invert bool) {
	c.invertColors = invert
}

// GetInvertColors returns whether colors are inverted
func (c *ColorSourceContentsImpl) GetInvertColors() bool {
	return c.invertColors
}

// SetOpacity sets the opacity for the color source
func (c *ColorSourceContentsImpl) SetOpacity(opacity float32) {
	c.opacity = opacity
}

// GetOpacity returns the current opacity
func (c *ColorSourceContentsImpl) GetOpacity() float32 {
	return c.opacity
}

// ApplyColorModifiers applies color modifications like inversion and opacity
func (c *ColorSourceContentsImpl) ApplyColorModifiers(color Color) Color {
	result := color

	if c.invertColors {
		result.R = 1.0 - result.R
		result.G = 1.0 - result.G
		result.B = 1.0 - result.B
	}

	result.A *= c.opacity

	return result
}

// GetCoverage returns the coverage from the geometry
func (c *ColorSourceContentsImpl) GetCoverage(entity Entity) Rect {
	if c.geometry == nil {
		return Rect{}
	}
	return c.geometry.GetCoverage(entity.GetTransform())
}

// CanInheritOpacity returns whether this content can inherit opacity from parent
func (c *ColorSourceContentsImpl) CanInheritOpacity(entity Entity) bool {
	return true
}

// SetInheritedOpacity sets inherited opacity from parent
func (c *ColorSourceContentsImpl) SetInheritedOpacity(opacity float32) {
	c.opacity *= opacity
}
