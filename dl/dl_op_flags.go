package dl

// OpFlags represents flags that describe the characteristics of display list operations.
type OpFlags uint32

// Operation flag constants
const (
	// OpFlagNone indicates no special flags.
	OpFlagNone OpFlags = 0

	// OpFlagModifiesTransparency indicates the operation modifies transparency.
	OpFlagModifiesTransparency OpFlags = 1 << 0

	// OpFlagReadsDestination indicates the operation reads from the destination.
	OpFlagReadsDestination OpFlags = 1 << 1

	// OpFlagFlood indicates the operation affects the entire clip region.
	OpFlagFlood OpFlags = 1 << 2

	// OpFlagNonDestructive indicates the operation is non-destructive.
	OpFlagNonDestructive OpFlags = 1 << 3

	// OpFlagCheap indicates the operation is computationally cheap.
	OpFlagCheap OpFlags = 1 << 4

	// OpFlagExpensive indicates the operation is computationally expensive.
	OpFlagExpensive OpFlags = 1 << 5
)

// HasFlag returns true if the specified flag is set.
func (f OpFlags) HasFlag(flag OpFlags) bool {
	return (f & flag) != 0
}

// WithFlag returns a new OpFlags with the specified flag set.
func (f OpFlags) WithFlag(flag OpFlags) OpFlags {
	return f | flag
}

// WithoutFlag returns a new OpFlags with the specified flag cleared.
func (f OpFlags) WithoutFlag(flag OpFlags) OpFlags {
	return f &^ flag
}

// Combine combines multiple flags.
func (f OpFlags) Combine(other OpFlags) OpFlags {
	return f | other
}

// AttributeFlags represents flags that describe paint attributes.
type AttributeFlags uint32

// Attribute flag constants
const (
	// AttrFlagNone indicates no attributes.
	AttrFlagNone AttributeFlags = 0

	// AttrFlagHasColor indicates the paint has a color.
	AttrFlagHasColor AttributeFlags = 1 << 0

	// AttrFlagHasColorSource indicates the paint has a color source (shader).
	AttrFlagHasColorSource AttributeFlags = 1 << 1

	// AttrFlagHasColorFilter indicates the paint has a color filter.
	AttrFlagHasColorFilter AttributeFlags = 1 << 2

	// AttrFlagHasImageFilter indicates the paint has an image filter.
	AttrFlagHasImageFilter AttributeFlags = 1 << 3

	// AttrFlagHasMaskFilter indicates the paint has a mask filter.
	AttrFlagHasMaskFilter AttributeFlags = 1 << 4

	// AttrFlagIsStroke indicates the paint is in stroke mode.
	AttrFlagIsStroke AttributeFlags = 1 << 5

	// AttrFlagIsAntiAlias indicates anti-aliasing is enabled.
	AttrFlagIsAntiAlias AttributeFlags = 1 << 6

	// AttrFlagHasText indicates the operation involves text rendering.
	AttrFlagHasText AttributeFlags = 1 << 7

	// AttrFlagHasImage indicates the operation involves image rendering.
	AttrFlagHasImage AttributeFlags = 1 << 8

	// AttrFlagModifiesTransparentBlack indicates the operation modifies transparent black pixels.
	AttrFlagModifiesTransparentBlack AttributeFlags = 1 << 9
)

// HasAttribute returns true if the specified attribute flag is set.
func (f AttributeFlags) HasAttribute(flag AttributeFlags) bool {
	return (f & flag) != 0
}

// WithAttribute returns a new AttributeFlags with the specified flag set.
func (f AttributeFlags) WithAttribute(flag AttributeFlags) AttributeFlags {
	return f | flag
}

// WithoutAttribute returns a new AttributeFlags with the specified flag cleared.
func (f AttributeFlags) WithoutAttribute(flag AttributeFlags) AttributeFlags {
	return f &^ flag
}

// GetOpFlags returns the operation flags for a given paint configuration.
func GetOpFlags(paint Paint) OpFlags {
	flags := OpFlagNone

	// Check if the operation modifies transparency
	if paint.BlendMode() != BlendModeSrcOver {
		flags = flags.WithFlag(OpFlagModifiesTransparency)
	}

	// Check if the operation reads from destination
	switch paint.BlendMode() {
	case BlendModeDstOver, BlendModeSrcIn, BlendModeDstIn,
		BlendModeSrcOut, BlendModeDstOut, BlendModeSrcATop,
		BlendModeDstATop, BlendModeXor, BlendModePlus,
		BlendModeScreen, BlendModeOverlay, BlendModeDarken,
		BlendModeLighten, BlendModeColorDodge, BlendModeColorBurn,
		BlendModeHardLight, BlendModeSoftLight, BlendModeDifference,
		BlendModeExclusion, BlendModeMultiply, BlendModeHue,
		BlendModeSaturation, BlendModeColor, BlendModeLuminosity:
		flags = flags.WithFlag(OpFlagReadsDestination)
	}

	// Anti-aliasing operations are typically more expensive
	if paint.IsAntiAlias() {
		flags = flags.WithFlag(OpFlagExpensive)
	} else {
		flags = flags.WithFlag(OpFlagCheap)
	}

	return flags
}

// GetAttributeFlags returns the attribute flags for a given paint.
func GetAttributeFlags(paint Paint) AttributeFlags {
	flags := AttrFlagNone

	// Always has color
	flags = flags.WithAttribute(AttrFlagHasColor)

	// Check stroke mode
	if paint.IsStroke() {
		flags = flags.WithAttribute(AttrFlagIsStroke)
	}

	// Check anti-aliasing
	if paint.IsAntiAlias() {
		flags = flags.WithAttribute(AttrFlagIsAntiAlias)
	}

	return flags
}
