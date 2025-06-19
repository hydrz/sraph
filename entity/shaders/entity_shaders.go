// Package shaders provides shader definitions and utilities for entity rendering.
// This package contains WGSL shaders and related utilities for various entity content types.
package shaders

// Entity shader bundle containing all shaders used by the entity system
const (
	// SolidColorVertexShader - Vertex shader for solid color rendering
	SolidColorVertexShader = `
struct VertexInput {
    @location(0) position: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) position: vec4<f32>,
}

struct Uniforms {
    transform: mat4x4<f32>,
    color: vec4<f32>,
    opacity: f32,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;

@vertex
fn vs_main(input: VertexInput) -> VertexOutput {
    var output: VertexOutput;
    output.position = uniforms.transform * vec4<f32>(input.position, 0.0, 1.0);
    return output;
}
`

	// SolidColorFragmentShader - Fragment shader for solid color rendering
	SolidColorFragmentShader = `
struct Uniforms {
    transform: mat4x4<f32>,
    color: vec4<f32>,
    opacity: f32,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;

@fragment
fn fs_main() -> @location(0) vec4<f32> {
    return vec4<f32>(uniforms.color.rgb, uniforms.color.a * uniforms.opacity);
}
`

	// TextureVertexShader - Vertex shader for texture rendering
	TextureVertexShader = `
struct VertexInput {
    @location(0) position: vec2<f32>,
    @location(1) tex_coord: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) position: vec4<f32>,
    @location(0) tex_coord: vec2<f32>,
}

struct Uniforms {
    transform: mat4x4<f32>,
    opacity: f32,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;

@vertex
fn vs_main(input: VertexInput) -> VertexOutput {
    var output: VertexOutput;
    output.position = uniforms.transform * vec4<f32>(input.position, 0.0, 1.0);
    output.tex_coord = input.tex_coord;
    return output;
}
`

	// TextureFragmentShader - Fragment shader for texture rendering
	TextureFragmentShader = `
struct Uniforms {
    transform: mat4x4<f32>,
    opacity: f32,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;
@group(0) @binding(1) var texture_sampler: sampler;
@group(0) @binding(2) var texture_2d: texture_2d<f32>;

@fragment
fn fs_main(@location(0) tex_coord: vec2<f32>) -> @location(0) vec4<f32> {
    let tex_color = textureSample(texture_2d, texture_sampler, tex_coord);
    return vec4<f32>(tex_color.rgb, tex_color.a * uniforms.opacity);
}
`

	// LinearGradientVertexShader - Vertex shader for linear gradient rendering
	LinearGradientVertexShader = `
struct VertexInput {
    @location(0) position: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) position: vec4<f32>,
    @location(0) local_coord: vec2<f32>,
}

struct Uniforms {
    transform: mat4x4<f32>,
    start_point: vec2<f32>,
    end_point: vec2<f32>,
    opacity: f32,
    tile_mode: i32,
    stop_count: i32,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;

@vertex
fn vs_main(input: VertexInput) -> VertexOutput {
    var output: VertexOutput;
    output.position = uniforms.transform * vec4<f32>(input.position, 0.0, 1.0);
    output.local_coord = input.position;
    return output;
}
`

	// LinearGradientFragmentShader - Fragment shader for linear gradient rendering
	LinearGradientFragmentShader = `
struct GradientStop {
    position: f32,
    color: vec4<f32>,
}

struct Uniforms {
    transform: mat4x4<f32>,
    start_point: vec2<f32>,
    end_point: vec2<f32>,
    opacity: f32,
    tile_mode: i32,
    stop_count: i32,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;
@group(0) @binding(1) var<storage, read> gradient_stops: array<GradientStop>;

fn interpolate_gradient(t: f32) -> vec4<f32> {
    if (uniforms.stop_count <= 0) {
        return vec4<f32>(0.0, 0.0, 0.0, 1.0);
    }

    if (uniforms.stop_count == 1) {
        return gradient_stops[0].color;
    }

    // Find the two stops to interpolate between
    for (var i: i32 = 0; i < uniforms.stop_count - 1; i++) {
        let stop1 = gradient_stops[i];
        let stop2 = gradient_stops[i + 1];

        if (t >= stop1.position && t <= stop2.position) {
            let factor = (t - stop1.position) / (stop2.position - stop1.position);
            return mix(stop1.color, stop2.color, factor);
        }
    }

    // Handle edge cases
    if (t < gradient_stops[0].position) {
        return gradient_stops[0].color;
    } else {
        return gradient_stops[uniforms.stop_count - 1].color;
    }
}

@fragment
fn fs_main(@location(0) local_coord: vec2<f32>) -> @location(0) vec4<f32> {
    let gradient_vector = uniforms.end_point - uniforms.start_point;
    let point_vector = local_coord - uniforms.start_point;

    let gradient_length_sq = dot(gradient_vector, gradient_vector);
    if (gradient_length_sq == 0.0) {
        return gradient_stops[0].color;
    }

    let t = dot(point_vector, gradient_vector) / gradient_length_sq;

    // Apply tile mode
    var final_t = t;
    switch (uniforms.tile_mode) {
        case 0: { // Clamp
            final_t = clamp(t, 0.0, 1.0);
        }
        case 1: { // Repeat
            final_t = fract(t);
        }
        case 2: { // Mirror
            let mod_t = fract(t);
            if (floor(t) % 2.0 == 0.0) {
                final_t = mod_t;
            } else {
                final_t = 1.0 - mod_t;
            }
        }
        default: { // Decal (treat as clamp for now)
            final_t = clamp(t, 0.0, 1.0);
        }
    }

    let color = interpolate_gradient(final_t);
    return vec4<f32>(color.rgb, color.a * uniforms.opacity);
}
`

	// TextVertexShader - Vertex shader for text rendering
	TextVertexShader = `
struct VertexInput {
    @location(0) position: vec2<f32>,
    @location(1) tex_coord: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) position: vec4<f32>,
    @location(0) tex_coord: vec2<f32>,
}

struct Uniforms {
    transform: mat4x4<f32>,
    color: vec4<f32>,
    opacity: f32,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;

@vertex
fn vs_main(input: VertexInput) -> VertexOutput {
    var output: VertexOutput;
    output.position = uniforms.transform * vec4<f32>(input.position, 0.0, 1.0);
    output.tex_coord = input.tex_coord;
    return output;
}
`

	// TextFragmentShader - Fragment shader for text rendering
	TextFragmentShader = `
struct Uniforms {
    transform: mat4x4<f32>,
    color: vec4<f32>,
    opacity: f32,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;
@group(0) @binding(1) var glyph_sampler: sampler;
@group(0) @binding(2) var glyph_atlas: texture_2d<f32>;

@fragment
fn fs_main(@location(0) tex_coord: vec2<f32>) -> @location(0) vec4<f32> {
    let alpha = textureSample(glyph_atlas, glyph_sampler, tex_coord).r;
    return vec4<f32>(uniforms.color.rgb, alpha * uniforms.color.a * uniforms.opacity);
}
`

	// StencilVertexShader - Vertex shader for stencil operations
	StencilVertexShader = `
struct VertexInput {
    @location(0) position: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) position: vec4<f32>,
}

struct Uniforms {
    transform: mat4x4<f32>,
    clip_op: i32,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;

@vertex
fn vs_main(input: VertexInput) -> VertexOutput {
    var output: VertexOutput;
    output.position = uniforms.transform * vec4<f32>(input.position, 0.0, 1.0);
    return output;
}
`

	// StencilFragmentShader - Fragment shader for stencil operations
	StencilFragmentShader = `
struct Uniforms {
    transform: mat4x4<f32>,
    clip_op: i32,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;

@fragment
fn fs_main() -> @location(0) vec4<f32> {
    // Stencil operations typically don't write to color buffer
    discard;
}
`
)

// ShaderType represents the type of shader
type ShaderType int

const (
	ShaderTypeSolidColor ShaderType = iota
	ShaderTypeTexture
	ShaderTypeLinearGradient
	ShaderTypeRadialGradient
	ShaderTypeText
	ShaderTypeStencil
)

// EntityShaderBundle contains all shaders for entity rendering
type EntityShaderBundle struct {
	solidColorVertex       string
	solidColorFragment     string
	textureVertex          string
	textureFragment        string
	linearGradientVertex   string
	linearGradientFragment string
	textVertex             string
	textFragment           string
	stencilVertex          string
	stencilFragment        string
}

// NewEntityShaderBundle creates a new entity shader bundle
func NewEntityShaderBundle() *EntityShaderBundle {
	return &EntityShaderBundle{
		solidColorVertex:       SolidColorVertexShader,
		solidColorFragment:     SolidColorFragmentShader,
		textureVertex:          TextureVertexShader,
		textureFragment:        TextureFragmentShader,
		linearGradientVertex:   LinearGradientVertexShader,
		linearGradientFragment: LinearGradientFragmentShader,
		textVertex:             TextVertexShader,
		textFragment:           TextFragmentShader,
		stencilVertex:          StencilVertexShader,
		stencilFragment:        StencilFragmentShader,
	}
}

// GetVertexShader returns the vertex shader for the given type
func (b *EntityShaderBundle) GetVertexShader(shaderType ShaderType) string {
	switch shaderType {
	case ShaderTypeSolidColor:
		return b.solidColorVertex
	case ShaderTypeTexture:
		return b.textureVertex
	case ShaderTypeLinearGradient:
		return b.linearGradientVertex
	case ShaderTypeText:
		return b.textVertex
	case ShaderTypeStencil:
		return b.stencilVertex
	default:
		return ""
	}
}

// GetFragmentShader returns the fragment shader for the given type
func (b *EntityShaderBundle) GetFragmentShader(shaderType ShaderType) string {
	switch shaderType {
	case ShaderTypeSolidColor:
		return b.solidColorFragment
	case ShaderTypeTexture:
		return b.textureFragment
	case ShaderTypeLinearGradient:
		return b.linearGradientFragment
	case ShaderTypeText:
		return b.textFragment
	case ShaderTypeStencil:
		return b.stencilFragment
	default:
		return ""
	}
}
