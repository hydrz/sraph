// Package gpu implements a Go interface for WebGPU, enabling GPU-accelerated graphics and computation.
//
// Typical WebGPU workflow:
//
// 1. Initialization
//   - CreateInstance: Initialize GPU environment
//   - InstanceRequestAdapter: Select physical GPU adapter
//   - AdapterRequestDevice: Create logical device and obtain queue
//
// 2. Surface & SwapChain (for window rendering)
//   - CreateSurface: Create rendering target (window/canvas)
//   - CreateSwapChain: Manage frame buffers for presentation
//
// 3. Resource Management
//   - CreateBuffer: Create vertex, index, uniform, etc.
//   - QueueWriteBuffer: Upload data to buffer
//   - CreateTexture: Create image or render target resources
//   - QueueWriteTexture: Upload pixel data to texture
//   - CreateSampler: Create texture sampler
//
// 4. Pipeline & Shader Setup
//   - CreateShaderModule: Compile shader (e.g., WGSL)
//   - CreateBindGroupLayout: Define resource binding layout
//   - CreatePipelineLayout: Combine BindGroupLayouts
//   - CreateRenderPipeline: Configure shaders and states
//
// 5. Resource Binding
//   - CreateBindGroup: Bind buffers, textures, samplers to pipeline
//
// 6. Command Encoding
//   - CreateCommandEncoder: Record GPU commands
//   - CommandEncoderBeginRenderPass: Start render pass
//   - RenderPassSetPipeline/SetVertexBuffer/SetIndexBuffer/SetBindGroup: Bind resources
//   - RenderPassDraw/DrawIndexed: Issue draw commands
//
// 7. Submission & Presentation
//   - QueueSubmit: Submit command buffer to GPU queue
//   - SwapChainPresent: Present frame to screen
package gpu
