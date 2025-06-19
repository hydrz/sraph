// Package entity provides gradient entity implementation.
package geometry

import (
	"github.com/opensraph/sraph/geom"
)

// GradientEntity represents an entity that renders a gradient.
type GradientEntity struct {
	*BaseEntity
	gradient GradientType
	colors   []geom.Color
	stops    []geom.F32
}

// GradientType defines the type of gradient.
type GradientType uint8

const (
	// GradientTypeLinear represents a linear gradient.
	GradientTypeLinear GradientType = iota
	// GradientTypeRadial represents a radial gradient.
	GradientTypeRadial
	// GradientTypeSweep represents a sweep/conical gradient.
	GradientTypeSweep
)

// NewGradientEntity creates a new gradient entity.
func NewGradientEntity(bounds geom.Rect[geom.F32], gradientType GradientType) *GradientEntity {
	entity := &GradientEntity{
		BaseEntity: NewBaseEntity(),
		gradient:   gradientType,
		colors:     make([]geom.Color, 0),
		stops:      make([]geom.F32, 0),
	}
	entity.SetBounds(bounds)
	return entity
}

// GradientType returns the type of gradient.
func (g *GradientEntity) GradientType() GradientType {
	return g.gradient
}

// SetGradientType sets the type of gradient.
func (g *GradientEntity) SetGradientType(gradientType GradientType) {
	g.gradient = gradientType
}

// Colors returns the gradient colors.
func (g *GradientEntity) Colors() []geom.Color {
	return g.colors
}

// Stops returns the gradient stops.
func (g *GradientEntity) Stops() []geom.F32 {
	return g.stops
}

// AddColorStop adds a color stop to the gradient.
func (g *GradientEntity) AddColorStop(stop geom.F32, color geom.Color) {
	g.stops = append(g.stops, stop)
	g.colors = append(g.colors, color)
}

// ClearColorStops clears all color stops.
func (g *GradientEntity) ClearColorStops() {
	g.colors = g.colors[:0]
	g.stops = g.stops[:0]
}

// Render implements RenderableEntity.
func (g *GradientEntity) Render(ctx RenderContext) error {
	// TODO: Implement gradient rendering
	// This would typically involve:
	// 1. Setting up gradient shader uniforms
	// 2. Passing color stops and gradient parameters
	// 3. Drawing a quad with gradient coordinates
	// 4. Applying the entity's transformation matrix
	return nil
}

// LinearGradientEntity represents a linear gradient entity.
type LinearGradientEntity struct {
	*GradientEntity
	startPoint geom.Point[geom.F32]
	endPoint   geom.Point[geom.F32]
}

// NewLinearGradientEntity creates a new linear gradient entity.
func NewLinearGradientEntity(bounds geom.Rect[geom.F32], startPoint, endPoint geom.Point[geom.F32]) *LinearGradientEntity {
	return &LinearGradientEntity{
		GradientEntity: NewGradientEntity(bounds, GradientTypeLinear),
		startPoint:     startPoint,
		endPoint:       endPoint,
	}
}

// StartPoint returns the start point of the linear gradient.
func (l *LinearGradientEntity) StartPoint() geom.Point[geom.F32] {
	return l.startPoint
}

// SetStartPoint sets the start point of the linear gradient.
func (l *LinearGradientEntity) SetStartPoint(point geom.Point[geom.F32]) {
	l.startPoint = point
}

// EndPoint returns the end point of the linear gradient.
func (l *LinearGradientEntity) EndPoint() geom.Point[geom.F32] {
	return l.endPoint
}

// SetEndPoint sets the end point of the linear gradient.
func (l *LinearGradientEntity) SetEndPoint(point geom.Point[geom.F32]) {
	l.endPoint = point
}

// Render implements RenderableEntity.
func (l *LinearGradientEntity) Render(ctx RenderContext) error {
	// TODO: Implement linear gradient rendering
	return nil
}

// RadialGradientEntity represents a radial gradient entity.
type RadialGradientEntity struct {
	*GradientEntity
	center geom.Point[geom.F32]
	radius geom.F32
	focal  geom.Point[geom.F32] // Focal point for elliptical gradients
}

// NewRadialGradientEntity creates a new radial gradient entity.
func NewRadialGradientEntity(bounds geom.Rect[geom.F32], center geom.Point[geom.F32], radius geom.F32) *RadialGradientEntity {
	return &RadialGradientEntity{
		GradientEntity: NewGradientEntity(bounds, GradientTypeRadial),
		center:         center,
		radius:         radius,
		focal:          center, // Default focal point is the center
	}
}

// Center returns the center point of the radial gradient.
func (r *RadialGradientEntity) Center() geom.Point[geom.F32] {
	return r.center
}

// SetCenter sets the center point of the radial gradient.
func (r *RadialGradientEntity) SetCenter(center geom.Point[geom.F32]) {
	r.center = center
}

// Radius returns the radius of the radial gradient.
func (r *RadialGradientEntity) Radius() geom.F32 {
	return r.radius
}

// SetRadius sets the radius of the radial gradient.
func (r *RadialGradientEntity) SetRadius(radius geom.F32) {
	r.radius = radius
}

// Focal returns the focal point of the radial gradient.
func (r *RadialGradientEntity) Focal() geom.Point[geom.F32] {
	return r.focal
}

// SetFocal sets the focal point of the radial gradient.
func (r *RadialGradientEntity) SetFocal(focal geom.Point[geom.F32]) {
	r.focal = focal
}

// Render implements RenderableEntity.
func (r *RadialGradientEntity) Render(ctx RenderContext) error {
	// TODO: Implement radial gradient rendering
	return nil
}

// SweepGradientEntity represents a sweep/conical gradient entity.
type SweepGradientEntity struct {
	*GradientEntity
	center     geom.Point[geom.F32]
	startAngle geom.F32 // Start angle in radians
	endAngle   geom.F32 // End angle in radians
}

// NewSweepGradientEntity creates a new sweep gradient entity.
func NewSweepGradientEntity(bounds geom.Rect[geom.F32], center geom.Point[geom.F32], startAngle, endAngle geom.F32) *SweepGradientEntity {
	return &SweepGradientEntity{
		GradientEntity: NewGradientEntity(bounds, GradientTypeSweep),
		center:         center,
		startAngle:     startAngle,
		endAngle:       endAngle,
	}
}

// Center returns the center point of the sweep gradient.
func (s *SweepGradientEntity) Center() geom.Point[geom.F32] {
	return s.center
}

// SetCenter sets the center point of the sweep gradient.
func (s *SweepGradientEntity) SetCenter(center geom.Point[geom.F32]) {
	s.center = center
}

// StartAngle returns the start angle of the sweep gradient in radians.
func (s *SweepGradientEntity) StartAngle() geom.F32 {
	return s.startAngle
}

// SetStartAngle sets the start angle of the sweep gradient in radians.
func (s *SweepGradientEntity) SetStartAngle(angle geom.F32) {
	s.startAngle = angle
}

// EndAngle returns the end angle of the sweep gradient in radians.
func (s *SweepGradientEntity) EndAngle() geom.F32 {
	return s.endAngle
}

// SetEndAngle sets the end angle of the sweep gradient in radians.
func (s *SweepGradientEntity) SetEndAngle(angle geom.F32) {
	s.endAngle = angle
}
