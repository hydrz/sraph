package contents

// GradientGenerator provides functionality for generating gradient textures
type GradientGenerator interface {
	// CreateGradientTexture creates a texture containing gradient data
	CreateGradientTexture(colors []Color, stops []float32) (Texture, error)

	// CreateLinearGradientTexture creates a texture for linear gradients
	CreateLinearGradientTexture(colors []Color, stops []float32, width, height int) (Texture, error)

	// CreateRadialGradientTexture creates a texture for radial gradients
	CreateRadialGradientTexture(colors []Color, stops []float32, radius int) (Texture, error)

	// CreateConicalGradientTexture creates a texture for conical gradients
	CreateConicalGradientTexture(colors []Color, stops []float32, width, height int) (Texture, error)

	// CreateSweepGradientTexture creates a texture for sweep gradients
	CreateSweepGradientTexture(colors []Color, stops []float32, width, height int) (Texture, error)
}

// GradientGeneratorImpl implements GradientGenerator
type GradientGeneratorImpl struct {
	context        ContentContext
	maxTextureSize int
}

// NewGradientGenerator creates a new gradient generator
func NewGradientGenerator(context ContentContext) GradientGenerator {
	return &GradientGeneratorImpl{
		context:        context,
		maxTextureSize: 1024, // Default max texture size
	}
}

// CreateGradientTexture creates a texture containing gradient data
func (g *GradientGeneratorImpl) CreateGradientTexture(colors []Color, stops []float32) (Texture, error) {
	if len(colors) == 0 {
		return nil, ErrInvalidArgument
	}

	// Default to linear gradient texture
	return g.CreateLinearGradientTexture(colors, stops, 256, 1)
}

// CreateLinearGradientTexture creates a texture for linear gradients
func (g *GradientGeneratorImpl) CreateLinearGradientTexture(colors []Color, stops []float32, width, height int) (Texture, error) {
	if len(colors) == 0 {
		return nil, ErrInvalidArgument
	}

	// Generate normalized stops if not provided
	normalizedStops := stops
	if len(stops) == 0 {
		normalizedStops = make([]float32, len(colors))
		for i := range normalizedStops {
			normalizedStops[i] = float32(i) / float32(len(colors)-1)
		}
	}

	// Generate texture data
	textureData := make([]byte, width*height*4) // RGBA

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			t := float32(x) / float32(width-1)
			color := g.interpolateColor(colors, normalizedStops, t)

			offset := (y*width + x) * 4
			textureData[offset] = byte(color.R * 255)
			textureData[offset+1] = byte(color.G * 255)
			textureData[offset+2] = byte(color.B * 255)
			textureData[offset+3] = byte(color.A * 255)
		}
	}

	// TODO: Create texture from data using context
	return nil, nil
}

// CreateRadialGradientTexture creates a texture for radial gradients
func (g *GradientGeneratorImpl) CreateRadialGradientTexture(colors []Color, stops []float32, radius int) (Texture, error) {
	if len(colors) == 0 {
		return nil, ErrInvalidArgument
	}

	size := radius * 2

	// Generate normalized stops if not provided
	normalizedStops := stops
	if len(stops) == 0 {
		normalizedStops = make([]float32, len(colors))
		for i := range normalizedStops {
			normalizedStops[i] = float32(i) / float32(len(colors)-1)
		}
	}

	// Generate texture data
	textureData := make([]byte, size*size*4) // RGBA
	center := float32(radius)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float32(x) - center
			dy := float32(y) - center
			distance := sqrt(dx*dx + dy*dy)
			t := distance / float32(radius)

			color := g.interpolateColor(colors, normalizedStops, t)

			offset := (y*size + x) * 4
			textureData[offset] = byte(color.R * 255)
			textureData[offset+1] = byte(color.G * 255)
			textureData[offset+2] = byte(color.B * 255)
			textureData[offset+3] = byte(color.A * 255)
		}
	}

	// TODO: Create texture from data using context
	return nil, nil
}

// CreateConicalGradientTexture creates a texture for conical gradients
func (g *GradientGeneratorImpl) CreateConicalGradientTexture(colors []Color, stops []float32, width, height int) (Texture, error) {
	if len(colors) == 0 {
		return nil, ErrInvalidArgument
	}

	// Generate normalized stops if not provided
	normalizedStops := stops
	if len(stops) == 0 {
		normalizedStops = make([]float32, len(colors))
		for i := range normalizedStops {
			normalizedStops[i] = float32(i) / float32(len(colors)-1)
		}
	}

	// Generate texture data
	textureData := make([]byte, width*height*4) // RGBA
	centerX := float32(width) / 2
	centerY := float32(height) / 2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := float32(x) - centerX
			dy := float32(y) - centerY
			angle := atan2(dy, dx)
			if angle < 0 {
				angle += 2 * PI
			}
			t := angle / (2 * PI)

			color := g.interpolateColor(colors, normalizedStops, t)

			offset := (y*width + x) * 4
			textureData[offset] = byte(color.R * 255)
			textureData[offset+1] = byte(color.G * 255)
			textureData[offset+2] = byte(color.B * 255)
			textureData[offset+3] = byte(color.A * 255)
		}
	}

	// TODO: Create texture from data using context
	return nil, nil
}

// CreateSweepGradientTexture creates a texture for sweep gradients
func (g *GradientGeneratorImpl) CreateSweepGradientTexture(colors []Color, stops []float32, width, height int) (Texture, error) {
	// Sweep gradient is similar to conical gradient
	return g.CreateConicalGradientTexture(colors, stops, width, height)
}

// interpolateColor interpolates between colors using the given stops
func (g *GradientGeneratorImpl) interpolateColor(colors []Color, stops []float32, t float32) Color {
	if len(colors) == 0 {
		return Color{R: 0, G: 0, B: 0, A: 1}
	}

	if len(colors) == 1 {
		return colors[0]
	}

	// Clamp t to [0, 1]
	if t <= 0 {
		return colors[0]
	}
	if t >= 1 {
		return colors[len(colors)-1]
	}

	// Find the color segment
	for i := 0; i < len(stops)-1; i++ {
		if t >= stops[i] && t <= stops[i+1] {
			// Interpolate between colors
			segmentT := (t - stops[i]) / (stops[i+1] - stops[i])
			color1 := colors[i]
			color2 := colors[i+1]

			return Color{
				R: color1.R + segmentT*(color2.R-color1.R),
				G: color1.G + segmentT*(color2.G-color1.G),
				B: color1.B + segmentT*(color2.B-color1.B),
				A: color1.A + segmentT*(color2.A-color1.A),
			}
		}
	}

	// Return last color if no segment found
	return colors[len(colors)-1]
}

// Math helper functions (these would typically be imported)
func sqrt(x float32) float32 {
	// TODO: Use proper math library
	return x // Placeholder
}

func atan2(y, x float32) float32 {
	// TODO: Use proper math library
	return 0 // Placeholder
}

const PI = 3.14159265359
