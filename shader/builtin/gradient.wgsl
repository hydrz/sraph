// Gradient vertex shader
struct VertexInput {
    @location(0) position: vec2<f32>,
    @location(1) local_coords: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(0) local_coords: vec2<f32>,
}

struct GradientUniforms {
    start_color: vec4<f32>,
    end_color: vec4<f32>,
    start_point: vec2<f32>,
    end_point: vec2<f32>,
}

@group(0) @binding(0)
var<uniform> gradient: GradientUniforms;

@vertex
fn vs_main(in: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    out.clip_position = vec4<f32>(in.position, 0.0, 1.0);
    out.local_coords = in.local_coords;
    return out;
}

// Linear gradient fragment shader
@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {
    let gradient_vec = gradient.end_point - gradient.start_point;
    let point_vec = in.local_coords - gradient.start_point;

    // Project point onto gradient vector
    let t = dot(point_vec, gradient_vec) / dot(gradient_vec, gradient_vec);
    let clamped_t = clamp(t, 0.0, 1.0);

    // Interpolate between start and end colors
    return mix(gradient.start_color, gradient.end_color, clamped_t);
}
