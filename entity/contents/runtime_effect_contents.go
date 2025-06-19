package contents

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// RuntimeEffect represents a runtime effect shader configuration
type RuntimeEffect struct {
	// ShaderSource contains the shader source code
	ShaderSource string
	// UniformData stores uniform variables for the shader
	UniformData map[string]interface{}
	// Textures contains input textures for the shader
	Textures []render.Texture
}

// RuntimeEffectContents provides runtime shader effect rendering capabilities
type RuntimeEffectContents struct {
	BaseContents
	effect       *RuntimeEffect
	geometry     Geometry
	uniformData  map[string]interface{}
	childTexture render.Texture
	opacity      float32
}

// NewRuntimeEffectContents creates a new runtime effect contents instance
func NewRuntimeEffectContents(effect *RuntimeEffect, geometry Geometry) *RuntimeEffectContents {
	return &RuntimeEffectContents{
		BaseContents: NewBaseContents(),
		effect:       effect,
		geometry:     geometry,
		uniformData:  make(map[string]interface{}),
		opacity:      1.0,
	}
}

// SetEffect updates the runtime effect
func (r *RuntimeEffectContents) SetEffect(effect *RuntimeEffect) {
	r.effect = effect
}

// GetEffect returns the current runtime effect
func (r *RuntimeEffectContents) GetEffect() *RuntimeEffect {
	return r.effect
}

// SetUniform sets a uniform variable value
func (r *RuntimeEffectContents) SetUniform(name string, value interface{}) {
	r.uniformData[name] = value
}

// GetUniform gets a uniform variable value
func (r *RuntimeEffectContents) GetUniform(name string) interface{} {
	return r.uniformData[name]
}

// SetChildTexture sets the input texture for the effect
func (r *RuntimeEffectContents) SetChildTexture(texture render.Texture) {
	r.childTexture = texture
}

// GetChildTexture returns the input texture
func (r *RuntimeEffectContents) GetChildTexture() render.Texture {
	return r.childTexture
}

// SetOpacity sets the opacity for the effect
func (r *RuntimeEffectContents) SetOpacity(opacity float32) {
	r.opacity = opacity
}

// GetOpacity returns the current opacity
func (r *RuntimeEffectContents) GetOpacity() float32 {
	return r.opacity
}

// GetGeometry returns the geometry for this contents
func (r *RuntimeEffectContents) GetGeometry() Geometry {
	return r.geometry
}

// Render implements the Contents interface
func (r *RuntimeEffectContents) Render(context *ContentContext, entity *Entity, pass *RenderPass) bool {
	// TODO: Implement runtime effect rendering
	// This would involve:
	// 1. Compiling the runtime shader if needed
	// 2. Setting up uniform data
	// 3. Binding input textures
	// 4. Executing the render command
	return false
}

// GetBounds returns the bounds of this contents
func (r *RuntimeEffectContents) GetBounds() geom.Rect {
	if r.geometry != nil {
		return r.geometry.GetBounds()
	}
	return geom.Rect{}
}

// Clone creates a copy of this contents
func (r *RuntimeEffectContents) Clone() Contents {
	clone := &RuntimeEffectContents{
		BaseContents: r.BaseContents.Clone().(BaseContents),
		effect:       r.effect, // Shallow copy of effect
		geometry:     r.geometry,
		childTexture: r.childTexture,
		opacity:      r.opacity,
		uniformData:  make(map[string]interface{}),
	}

	// Deep copy uniform data
	for k, v := range r.uniformData {
		clone.uniformData[k] = v
	}

	return clone
}

// RuntimeEffectFactory creates runtime effect contents instances
type RuntimeEffectFactory struct{}

// CreateRuntimeEffect creates a new runtime effect contents
func (f *RuntimeEffectFactory) CreateRuntimeEffect(effect *RuntimeEffect, geometry Geometry) Contents {
	return NewRuntimeEffectContents(effect, geometry)
}

// ValidateEffect validates a runtime effect configuration
func (f *RuntimeEffectFactory) ValidateEffect(effect *RuntimeEffect) error {
	// TODO: Implement effect validation
	// Check shader syntax, uniform types, etc.
	return nil
}
