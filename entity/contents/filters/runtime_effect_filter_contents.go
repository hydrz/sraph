package filters

import (
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/render"
)

// RuntimeEffectFilterContents represents a filter that applies custom shader effects
type RuntimeEffectFilterContents struct {
	FilterContentsBase
	ShaderSource  string
	Uniforms      map[string]interface{}
	InputTextures []*render.Texture
	OutputTexture *render.Texture
}

// NewRuntimeEffectFilterContents creates a new runtime effect filter
func NewRuntimeEffectFilterContents(shaderSource string) *RuntimeEffectFilterContents {
	return &RuntimeEffectFilterContents{
		ShaderSource: shaderSource,
		Uniforms:     make(map[string]interface{}),
	}
}

// WithUniform sets a uniform value for the shader
func (f *RuntimeEffectFilterContents) WithUniform(name string, value interface{}) *RuntimeEffectFilterContents {
	f.Uniforms[name] = value
	return f
}

// WithInputTexture adds an input texture for the shader
func (f *RuntimeEffectFilterContents) WithInputTexture(texture *render.Texture) *RuntimeEffectFilterContents {
	f.InputTextures = append(f.InputTextures, texture)
	return f
}

// WithOutputTexture sets the output texture
func (f *RuntimeEffectFilterContents) WithOutputTexture(texture *render.Texture) *RuntimeEffectFilterContents {
	f.OutputTexture = texture
	return f
}

// Render renders the runtime effect filter
func (f *RuntimeEffectFilterContents) Render(renderer render.ContentRenderer, entity render.Entity) bool {
	// TODO: Implement custom shader effect rendering
	// This involves compiling and executing custom shader code
	return false
}

// GetCoverage returns the coverage area of the filter
func (f *RuntimeEffectFilterContents) GetCoverage(entity render.Entity) geom.Rect {
	// Runtime effects can potentially change coverage, but default to entity bounds
	return entity.GetBounds()
}

// Clone creates a copy of the filter
func (f *RuntimeEffectFilterContents) Clone() FilterContents {
	uniforms := make(map[string]interface{})
	for k, v := range f.Uniforms {
		uniforms[k] = v
	}

	inputTextures := make([]*render.Texture, len(f.InputTextures))
	copy(inputTextures, f.InputTextures)

	return &RuntimeEffectFilterContents{
		ShaderSource:  f.ShaderSource,
		Uniforms:      uniforms,
		InputTextures: inputTextures,
		OutputTexture: f.OutputTexture,
	}
}

// SetShader sets the shader source code
func (f *RuntimeEffectFilterContents) SetShader(source string) {
	f.ShaderSource = source
}

// GetShader returns the shader source code
func (f *RuntimeEffectFilterContents) GetShader() string {
	return f.ShaderSource
}

// GetUniform returns a uniform value
func (f *RuntimeEffectFilterContents) GetUniform(name string) (interface{}, bool) {
	value, exists := f.Uniforms[name]
	return value, exists
}

// GetInputTextures returns all input textures
func (f *RuntimeEffectFilterContents) GetInputTextures() []*render.Texture {
	return f.InputTextures
}

// ClearInputTextures removes all input textures
func (f *RuntimeEffectFilterContents) ClearInputTextures() {
	f.InputTextures = f.InputTextures[:0]
}
