package contents

// FramebufferBlendContents represents contents that blend with the current framebuffer
type FramebufferBlendContents interface {
	ColorSourceContents

	// SetBlendMode sets the blend mode for framebuffer blending
	SetBlendMode(mode BlendMode)

	// GetBlendMode returns the current blend mode
	GetBlendMode() BlendMode

	// SetSourceContents sets the source contents to blend
	SetSourceContents(contents Contents)

	// GetSourceContents returns the source contents
	GetSourceContents() Contents

	// SetDestinationContents sets the destination contents
	SetDestinationContents(contents Contents)

	// GetDestinationContents returns the destination contents
	GetDestinationContents() Contents
}

// FramebufferBlendContentsImpl implements FramebufferBlendContents
type FramebufferBlendContentsImpl struct {
	*ColorSourceContentsImpl
	blendMode           BlendMode
	sourceContents      Contents
	destinationContents Contents
	requiresReadback    bool
}

// NewFramebufferBlendContents creates new framebuffer blend contents
func NewFramebufferBlendContents() FramebufferBlendContents {
	return &FramebufferBlendContentsImpl{
		ColorSourceContentsImpl: NewColorSourceContents(),
		blendMode:               BlendModeNormal,
		requiresReadback:        false,
	}
}

// SetBlendMode sets the blend mode for framebuffer blending
func (f *FramebufferBlendContentsImpl) SetBlendMode(mode BlendMode) {
	f.blendMode = mode
	f.requiresReadback = f.isAdvancedBlendMode(mode)
}

// GetBlendMode returns the current blend mode
func (f *FramebufferBlendContentsImpl) GetBlendMode() BlendMode {
	return f.blendMode
}

// SetSourceContents sets the source contents to blend
func (f *FramebufferBlendContentsImpl) SetSourceContents(contents Contents) {
	f.sourceContents = contents
}

// GetSourceContents returns the source contents
func (f *FramebufferBlendContentsImpl) GetSourceContents() Contents {
	return f.sourceContents
}

// SetDestinationContents sets the destination contents
func (f *FramebufferBlendContentsImpl) SetDestinationContents(contents Contents) {
	f.destinationContents = contents
}

// GetDestinationContents returns the destination contents
func (f *FramebufferBlendContentsImpl) GetDestinationContents() Contents {
	return f.destinationContents
}

// Render renders the framebuffer blend contents
func (f *FramebufferBlendContentsImpl) Render(context ContentContext, entity Entity, pass RenderPass) bool {
	if f.sourceContents == nil {
		return false
	}

	// TODO: Implement framebuffer blending with readback if needed
	if f.requiresReadback {
		return f.renderAdvancedBlend(context, entity, pass)
	}

	return f.renderNormalBlend(context, entity, pass)
}

// renderAdvancedBlend handles advanced blend modes that require framebuffer readback
func (f *FramebufferBlendContentsImpl) renderAdvancedBlend(context ContentContext, entity Entity, pass RenderPass) bool {
	// TODO: Implement advanced blending with framebuffer readback
	return false
}

// renderNormalBlend handles normal blend modes
func (f *FramebufferBlendContentsImpl) renderNormalBlend(context ContentContext, entity Entity, pass RenderPass) bool {
	// TODO: Implement normal blending
	return false
}

// GetCoverage returns the coverage of the blend contents
func (f *FramebufferBlendContentsImpl) GetCoverage(entity Entity) Rect {
	if f.sourceContents == nil {
		return Rect{}
	}

	sourceCoverage := f.sourceContents.GetCoverage(entity)

	if f.destinationContents != nil {
		destCoverage := f.destinationContents.GetCoverage(entity)
		return sourceCoverage.Union(destCoverage)
	}

	return sourceCoverage
}

// Clone creates a copy of the framebuffer blend contents
func (f *FramebufferBlendContentsImpl) Clone() Contents {
	clone := &FramebufferBlendContentsImpl{
		ColorSourceContentsImpl: f.ColorSourceContentsImpl,
		blendMode:               f.blendMode,
		requiresReadback:        f.requiresReadback,
	}

	if f.sourceContents != nil {
		clone.sourceContents = f.sourceContents.Clone()
	}

	if f.destinationContents != nil {
		clone.destinationContents = f.destinationContents.Clone()
	}

	return clone
}

// isAdvancedBlendMode checks if the blend mode requires advanced blending
func (f *FramebufferBlendContentsImpl) isAdvancedBlendMode(mode BlendMode) bool {
	switch mode {
	case BlendModeMultiply, BlendModeScreen, BlendModeOverlay,
		BlendModeDarken, BlendModeLighten, BlendModeColorDodge,
		BlendModeColorBurn, BlendModeHardLight, BlendModeSoftLight,
		BlendModeDifference, BlendModeExclusion, BlendModeHue,
		BlendModeSaturation, BlendModeColor, BlendModeLuminosity:
		return true
	default:
		return false
	}
}

// RequiresReadback returns whether the blend mode requires framebuffer readback
func (f *FramebufferBlendContentsImpl) RequiresReadback() bool {
	return f.requiresReadback
}
