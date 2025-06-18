package dl

import "github.com/opensraph/sraph/geom"

// BlendMode defines how new pixels are combined with existing pixels.
// This is a wrapper around geom.BlendMode for display list specific functionality.
type BlendMode = geom.BlendMode

// Predefined blend modes (re-export from geom package)
const (
	// BlendModeClear clears the destination.
	BlendModeClear = geom.BlendModeClear
	// BlendModeSrc replaces the destination with the source.
	BlendModeSrc = geom.BlendModeSrc
	// BlendModeDst keeps the destination unchanged.
	BlendModeDst = geom.BlendModeDst
	// BlendModeSrcOver composites the source over the destination.
	BlendModeSrcOver = geom.BlendModeSrcOver
	// BlendModeDstOver composites the destination over the source.
	BlendModeDstOver = geom.BlendModeDstOver
	// BlendModeSrcIn shows the source where the destination is opaque.
	BlendModeSrcIn = geom.BlendModeSrcIn
	// BlendModeDstIn shows the destination where the source is opaque.
	BlendModeDstIn = geom.BlendModeDstIn
	// BlendModeSrcOut shows the source where the destination is transparent.
	BlendModeSrcOut = geom.BlendModeSrcOut
	// BlendModeDstOut shows the destination where the source is transparent.
	BlendModeDstOut = geom.BlendModeDstOut
	// BlendModeSrcATop composites the source over the destination, but only where the destination is opaque.
	BlendModeSrcATop = geom.BlendModeSrcATop
	// BlendModeDstATop composites the destination over the source, but only where the source is opaque.
	BlendModeDstATop = geom.BlendModeDstATop
	// BlendModeXor shows the source where the destination is transparent and vice versa.
	BlendModeXor = geom.BlendModeXor
	// BlendModePlus adds the source and destination.
	BlendModePlus = geom.BlendModePlus
	// BlendModeModulate multiplies the source and destination.
	BlendModeModulate = geom.BlendModeModulate
	// BlendModeScreen performs screen blending.
	BlendModeScreen = geom.BlendModeScreen
	// BlendModeOverlay performs overlay blending.
	BlendModeOverlay = geom.BlendModeOverlay
	// BlendModeDarken performs darken blending.
	BlendModeDarken = geom.BlendModeDarken
	// BlendModeLighten performs lighten blending.
	BlendModeLighten = geom.BlendModeLighten
	// BlendModeColorDodge performs color dodge blending.
	BlendModeColorDodge = geom.BlendModeColorDodge
	// BlendModeColorBurn performs color burn blending.
	BlendModeColorBurn = geom.BlendModeColorBurn
	// BlendModeHardLight performs hard light blending.
	BlendModeHardLight = geom.BlendModeHardLight
	// BlendModeSoftLight performs soft light blending.
	BlendModeSoftLight = geom.BlendModeSoftLight
	// BlendModeDifference performs difference blending.
	BlendModeDifference = geom.BlendModeDifference
	// BlendModeExclusion performs exclusion blending.
	BlendModeExclusion = geom.BlendModeExclusion
	// BlendModeMultiply performs multiply blending.
	BlendModeMultiply = geom.BlendModeMultiply
	// BlendModeHue performs hue blending.
	BlendModeHue = geom.BlendModeHue
	// BlendModeSaturation performs saturation blending.
	BlendModeSaturation = geom.BlendModeSaturation
	// BlendModeColor performs color blending.
	BlendModeColor = geom.BlendModeColor
	// BlendModeLuminosity performs luminosity blending.
	BlendModeLuminosity = geom.BlendModeLuminosity
)
