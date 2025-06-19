package contents

import (
	"github.com/opensraph/sraph/display"
	"github.com/opensraph/sraph/geom"
)

// AnonymousContents represents anonymous rendering contents that can be created with custom render and coverage functions
type AnonymousContents interface {
	Contents
}

// RenderProc is a function type for custom rendering operations
type RenderProc func(context ContentContext, entity Entity, pass RenderPass) bool

// CoverageProc is a function type for custom coverage calculations
type CoverageProc func(entity Entity) geom.Rect

// AnonymousContentsImpl implements AnonymousContents
type AnonymousContentsImpl struct {
	*ContentsImpl
	renderProc   RenderProc
	coverageProc CoverageProc
}

// NewAnonymousContents creates a new anonymous contents with custom render and coverage functions
func NewAnonymousContents(renderProc RenderProc, coverageProc CoverageProc) AnonymousContents {
	return &AnonymousContentsImpl{
		ContentsImpl: &ContentsImpl{
			isValid: true,
		},
		renderProc:   renderProc,
		coverageProc: coverageProc,
	}
}

// Render renders the anonymous contents using the custom render function
func (a *AnonymousContentsImpl) Render(context ContentContext, entity Entity, pass RenderPass) bool {
	if a.renderProc == nil {
		return false
	}
	return a.renderProc(context, entity, pass)
}

// GetCoverage returns the coverage using the custom coverage function
func (a *AnonymousContentsImpl) GetCoverage(entity Entity) geom.Rect {
	if a.coverageProc == nil {
		return geom.Rect{}
	}
	return a.coverageProc(entity)
}

// Clone creates a copy of the anonymous contents
func (a *AnonymousContentsImpl) Clone() Contents {
	return &AnonymousContentsImpl{
		ContentsImpl: a.ContentsImpl.Clone().(*ContentsImpl),
		renderProc:   a.renderProc,
		coverageProc: a.coverageProc,
	}
}

// GetColorFilter returns the color filter applied to the contents
func (a *AnonymousContentsImpl) GetColorFilter() display.ColorFilter {
	// TODO: Implement color filter support
	return nil
}

// SetColorFilter sets the color filter for the contents
func (a *AnonymousContentsImpl) SetColorFilter(filter display.ColorFilter) {
	// TODO: Implement color filter support
}
