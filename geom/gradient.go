package geom

import "math"

// LinearGradient defines a linear gradient by a sequence of color stops.
type LinearGradient struct {
	Stops []GradientStop
}

// NewLinearGradient returns a LinearGradient with the provided stops.
func NewLinearGradient(stops []GradientStop) LinearGradient {
	return LinearGradient{Stops: stops}
}

// ToBuffer returns the gradient data for the linear gradient.
// If no stops are present, returns an invalid GradientData.
func (lg LinearGradient) ToBuffer() GradientData {
	if len(lg.Stops) == 0 {
		return GradientData{}
	}

	colors := make([]Color, len(lg.Stops))
	stops := make([]Scalar, len(lg.Stops))

	for i, stop := range lg.Stops {
		colors[i] = stop.Color
		stops[i] = stop.Position
	}

	return CreateGradientBuffer(colors, stops)
}

// RadialGradient defines a radial gradient by center, radius, and color stops.
type RadialGradient struct {
	Center Point  // Center of the radial gradient.
	Radius Scalar // Radius of the radial gradient.
	Stops  []GradientStop
}

// NewRadialGradient returns a RadialGradient with the given center, radius, and stops.
func NewRadialGradient(center Point, radius Scalar, stops []GradientStop) RadialGradient {
	return RadialGradient{
		Center: center,
		Radius: radius,
		Stops:  stops,
	}
}

// ToBuffer returns the gradient data for the radial gradient.
// If no stops are present, returns an invalid GradientData.
func (rg RadialGradient) ToBuffer() GradientData {
	if len(rg.Stops) == 0 {
		return GradientData{}
	}

	colors := make([]Color, len(rg.Stops))
	stops := make([]Scalar, len(rg.Stops))

	for i, stop := range rg.Stops {
		colors[i] = stop.Color
		stops[i] = stop.Position
	}

	return CreateGradientBuffer(colors, stops)
}

// GradientData holds the color data and texture size for a gradient.
// If TextureSize is 0, the gradient is considered invalid.
type GradientData struct {
	ColorBytes  []uint8 // RGBA bytes for the gradient texture.
	TextureSize uint32  // Number of texels in the gradient.
}

// CreateGradientBuffer generates a buffer of interpolated color bytes for a linear gradient.
// The colors and stops slices must have the same length. Panics if not.
// The resulting buffer can be used as a 1D texture for gradient rendering.
func CreateGradientBuffer(colors []Color, stops []Scalar) GradientData {
	if len(stops) != len(colors) {
		panic("stops and colors must have the same length")
	}

	var textureSize uint32
	if len(stops) == 2 {
		textureSize = uint32(len(colors))
	} else {
		minimumDelta := Scalar(1.0) // Minimum delta between stops.
		for i := 1; i < len(stops); i++ {
			value := stops[i] - stops[i-1]
			if value < 0.0001 {
				continue
			}
			if value < minimumDelta {
				minimumDelta = value
			}
		}
		calculated := uint32(math.Round(ToFloat64(1.0/minimumDelta))) + 1
		if calculated > 1024 {
			textureSize = 1024
		} else {
			textureSize = calculated
		}
	}

	data := GradientData{
		ColorBytes:  make([]uint8, 0, textureSize*4),
		TextureSize: textureSize,
	}

	if textureSize == uint32(len(colors)) && len(colors) <= 1024 {
		// Direct mapping: one texel per color.
		for i := 0; i < len(colors); i++ {
			appendColor(colors[i], &data)
		}
	} else {
		// Interpolated mapping.
		previousColor := colors[0]
		previousStop := Scalar(0.0)
		previousColorIndex := 0

		appendColor(previousColor, &data)

		for i := uint32(1); i < textureSize-1; i++ {
			scaledI := Scalar(i) / Scalar(textureSize-1)
			nextColor := colors[previousColorIndex+1]
			nextStop := stops[previousColorIndex+1]

			if NearlyEqual(scaledI, nextStop) {
				appendColor(nextColor, &data)
				previousColor = nextColor
				previousStop = nextStop
				previousColorIndex++
			} else if scaledI < nextStop {
				t := (scaledI - previousStop) / (nextStop - previousStop)
				mixedColor := previousColor.Lerp(nextColor, t)
				appendColor(mixedColor, &data)
			} else {
				previousColor = nextColor
				previousStop = nextStop
				previousColorIndex++

				if previousColorIndex+1 < len(colors) {
					nextColor = colors[previousColorIndex+1]
					nextStop = stops[previousColorIndex+1]
					t := (scaledI - previousStop) / (nextStop - previousStop)
					mixedColor := previousColor.Lerp(nextColor, t)
					appendColor(mixedColor, &data)
				} else {
					appendColor(previousColor, &data)
				}
			}
		}

		appendColor(colors[len(colors)-1], &data)
	}

	return data
}

// appendColor appends the RGBA bytes of color to the gradient data.
func appendColor(color Color, data *GradientData) {
	rgba := color.Go()
	data.ColorBytes = append(data.ColorBytes, rgba.R, rgba.G, rgba.B, rgba.A)
}

// IsValid reports whether the gradient data is valid (TextureSize > 0).
func (g GradientData) IsValid() bool {
	return g.TextureSize > 0
}

// GradientStop represents a color stop in a gradient at a given position in [0,1].
type GradientStop struct {
	Color    Color  // Color at this stop.
	Position Scalar // Position from 0.0 to 1.0.
}
