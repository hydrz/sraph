// Gradient shader for rendering linear and radial gradients
// Supports multiple gradient stops and various gradient types

struct VertexInput {
    @location(0) position: vec2<f32>,
    @location(1) local_coords: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(0) local_coords: vec2<f32>,
}

struct GradientUniforms {
    transform: mat4x4<f32>,
    start_color: vec4<f32>,
    end_color: vec4<f32>,
    start_point: vec2<f32>,
    end_point: vec2<f32>,
    gradient_type: i32, // 0 = linear, 1 = radial
    stop_count: i32,    // Number of gradient stops (future enhancement)
}

@group(0) @binding(0)
var<uniform> gradient: GradientUniforms;

// Vertex shader
@vertex
fn vs_main(in: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    out.clip_position = gradient.transform * vec4<f32>(in.position, 0.0, 1.0);
    out.local_coords = in.local_coords;
    return out;
}

// Fragment shader
@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {
    var t: f32;

    if (gradient.gradient_type == 0) {
        // Linear gradient
        let gradient_vec = gradient.end_point - gradient.start_point;
        let point_vec = in.local_coords - gradient.start_point;

        // Project point onto gradient vector
        t = dot(point_vec, gradient_vec) / dot(gradient_vec, gradient_vec);
    } else {
        // Radial gradient
        let center = gradient.start_point;
        let radius = length(gradient.end_point - gradient.start_point);
        let distance = length(in.local_coords - center);
        t = distance / radius;
    }

    // Clamp t to [0, 1] and interpolate colors
    let clamped_t = clamp(t, 0.0, 1.0);
    return mix(gradient.start_color, gradient.end_color, clamped_t);
}
