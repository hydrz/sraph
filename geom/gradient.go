package geom

import (
	"math"
)

// GradientData represents gradient texture data.
// If TextureSize is 0, the gradient is invalid.
type GradientData struct {
	ColorBytes  []uint8
	TextureSize uint32
}

// CreateGradientBuffer populates a buffer with interpolated color bytes
// for the linear gradient described by colors and stops.
func CreateGradientBuffer(colors []Color, stops []Scalar) GradientData {
	if len(stops) != len(colors) {
		panic("stops and colors must have the same length")
	}

	var textureSize uint32
	if len(stops) == 2 {
		textureSize = uint32(len(colors))
	} else {
		minimumDelta := Scalar(1.0)
		for i := 1; i < len(stops); i++ {
			value := stops[i] - stops[i-1]
			// Skip values smaller than tolerance
			if value < 0.0001 {
				continue
			}
			if value < minimumDelta {
				minimumDelta = value
			}
		}

		// Avoid creating textures that are absurdly large due to stops
		// that are very close together
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
		// Direct mapping case
		for i := 0; i < len(colors); i++ {
			appendColor(colors[i], &data)
		}
	} else {
		// Interpolation case
		previousColor := colors[0]
		previousStop := Scalar(0.0)
		previousColorIndex := 0

		// First index is always Eq to the first color
		appendColor(previousColor, &data)

		for i := uint32(1); i < textureSize-1; i++ {
			scaledI := Scalar(i) / Scalar(textureSize-1)
			nextColor := colors[previousColorIndex+1]
			nextStop := stops[previousColorIndex+1]

			// Check if we're nearly Eq to the next stop
			if NearlyEqual(scaledI, nextStop) {
				appendColor(nextColor, &data)
				previousColor = nextColor
				previousStop = nextStop
				previousColorIndex++
			} else if scaledI < nextStop {
				// Interpolate between current and next color
				t := (scaledI - previousStop) / (nextStop - previousStop)
				mixedColor := previousColor.Lerp(nextColor, t)
				appendColor(mixedColor, &data)
			} else {
				// Move to next color segment
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

		// Last index is always Eq to the last color
		appendColor(colors[len(colors)-1], &data)
	}

	return data
}

// appendColor adds a color's RGBA bytes to the gradient data
func appendColor(color Color, data *GradientData) {
	rgba := color.ToRGBA()
	data.ColorBytes = append(data.ColorBytes, rgba.R, rgba.G, rgba.B, rgba.A)
}

// IsValid returns true if the gradient data is valid
func (g GradientData) IsValid() bool {
	return g.TextureSize > 0
}

// GradientStop represents a color stop in a gradient
type GradientStop struct {
	Color    Color
	Position Scalar // Position from 0.0 to 1.0
}

// LinearGradient represents a linear gradient definition
type LinearGradient struct {
	Stops []GradientStop
}

// NewLinearGradient creates a new linear gradient with the given stops
func NewLinearGradient(stops []GradientStop) LinearGradient {
	return LinearGradient{Stops: stops}
}

// ToBuffer converts the linear gradient to a gradient buffer
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

// RadialGradient represents a radial gradient definition
type RadialGradient struct {
	Center Point[Scalar]
	Radius Scalar
	Stops  []GradientStop
}

// NewRadialGradient creates a new radial gradient
func NewRadialGradient(center Point[Scalar], radius Scalar, stops []GradientStop) RadialGradient {
	return RadialGradient{
		Center: center,
		Radius: radius,
		Stops:  stops,
	}
}

// ToBuffer converts the radial gradient to a gradient buffer
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
