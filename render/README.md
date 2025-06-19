# Render Package

The render package provides the core rendering infrastructure for the sraph graphics framework. It is inspired by Flutter's Impeller renderer architecture and provides a modern, GPU-accelerated rendering system.

## Architecture Overview

The render package is organized into several key components:

### Core Components

- **Context**: The main rendering context that manages GPU resources and state
- **CommandBuffer**: Encodes and manages rendering commands
- **CommandQueue**: Manages command buffer execution
- **RenderPass**: Manages render passes with render targets
- **Pipeline**: Manages render and compute pipelines
- **Resources**: Manages GPU resources like buffers, textures, and samplers

### Resource Management

- **Buffer**: GPU buffer resources for vertex, index, and uniform data
- **Texture**: GPU texture resources for images and render targets
- **Sampler**: Texture sampling configuration
- **Resource Pool**: Efficient resource pooling for performance

### Rendering Pipeline

- **RenderTarget**: Defines render targets with color/depth/stencil attachments
- **Surface**: Manages presentable surfaces
- **Vertex Descriptors**: Defines vertex attribute layouts
- **Shader Functions**: Manages vertex and fragment shaders

### Backend Support

The render package supports multiple graphics backends through the `backend` subpackage:

- OpenGL
- Vulkan
- Metal
- DirectX 11/12
- WebGPU

## Key Features

### Modern GPU Architecture
- Command buffer-based rendering
- Multiple render passes
- Compute shader support
- Resource state management

### Performance Optimizations
- Resource pooling
- Efficient memory management
- Batch rendering support
- GPU-driven rendering

### Cross-Platform Support
- Abstracted backend interface
- Platform-specific optimizations
- Automatic backend selection

## Usage Example

```go
// Create a rendering context
context := render.NewContext()

// Create a render target
renderTarget := render.NewRenderTarget(render.RenderTargetConfig{
    Width:  800,
    Height: 600,
})

// Create a command buffer
cmdBuffer := context.CreateCommandBuffer()

// Begin render pass
renderPass := cmdBuffer.BeginRenderPass(renderTarget)

// Set pipeline and draw
renderPass.SetPipeline(myPipeline)
renderPass.DrawIndexed(indexCount, instanceCount, 0, 0, 0)

// End render pass and submit
cmdBuffer.EndRenderPass()
cmdBuffer.Submit()
```

## File Structure

```
render/
├── render.go              # Main render package interface
├── context.go             # Render context implementation
├── command_buffer.go      # Command buffer management
├── command_queue.go       # Command queue management
├── render_pass.go         # Render pass management
├── pipeline.go            # Pipeline management
├── pipeline_library.go    # Pipeline caching
├── shader_function.go     # Shader function management
├── shader_library.go      # Shader caching
├── vertex_descriptor.go   # Vertex attribute descriptions
├── buffer.go              # Buffer resource management
├── resource.go            # Base resource management
├── surface.go             # Surface management
├── render_target.go       # Render target management
├── capabilities.go        # GPU capabilities
├── resource_allocator.go  # Resource allocation
├── sampler_library.go     # Sampler caching
├── blit_command.go        # Blit operations
├── blit_pass.go           # Blit pass management
├── compute_pass.go        # Compute pass management
├── compute_pipeline_descriptor.go # Compute pipeline configuration
├── pipeline_descriptor.go # Render pipeline configuration
├── command.go             # Individual render commands
├── pool.go                # Resource pooling
├── snapshot.go            # Render state snapshots
├── vertex_buffer_builder.go # Vertex buffer utilities
└── backend/
    └── backend.go         # Backend abstraction
```

## Design Principles

1. **Performance First**: Optimized for modern GPU architectures
2. **Cross-Platform**: Abstracted backend interface for portability
3. **Resource Efficiency**: Automatic resource management and pooling
4. **Developer Friendly**: Clean, intuitive API design
5. **Extensible**: Modular architecture for easy extension

## Integration with Entity System

The render package works closely with the entity system to provide:

- Entity rendering through the entity renderer
- Content-based rendering pipeline
- Efficient batching and sorting
- GPU-accelerated effects and filters

## TODO

This architecture provides the foundation for:
- [ ] Backend implementations (OpenGL, Vulkan, Metal, etc.)
- [ ] Shader compilation and management
- [ ] Advanced rendering techniques
- [ ] Performance profiling and debugging tools
- [ ] Integration with the entity system
- [ ] Comprehensive test coverage
3. **Entity System** - Converts display list operations to renderable entities
4. **Render Passes** - Manages GPU command recording and execution

## Key Components

### Rendering Context

The `Context` interface provides the foundation for all rendering operations:

```go
// Create a rendering context
device := gpu.GetDevice() // Assume you have a GPU device
context, err := render.NewImpellerContext(device)
if err != nil {
    log.Fatal(err)
}
defer context.Shutdown()
```

### Renderer

The main `Renderer` interface provides high-level rendering functionality:

```go
// Create a renderer
renderer, err := render.NewImpellerRenderer(device)
if err != nil {
    log.Fatal(err)
}
defer renderer.Shutdown()

// Create a surface for rendering
surfaceDesc := render.SurfaceDescriptor{
    Label:       "Main Surface",
    Size:        geom.Size[geom.F32]{Width: 800, Height: 600},
    Format:      gpu.TextureFormatBGRA8Unorm,
    Usage:       gpu.TextureUsageRenderAttachment,
    SampleCount: 1,
}

surface, err := renderer.CreateSurface(surfaceDesc)
if err != nil {
    log.Fatal(err)
}

// Render a display list
err = renderer.Render(surface, displayList)
if err != nil {
    log.Fatal(err)
}
```

### Resource Management

The render package includes specialized managers for different types of GPU resources:

#### Resource Allocator
```go
allocator := context.GetResourceAllocator()

// Create a buffer
bufferDesc := render.BufferDescriptor{
    Label: "Vertex Buffer",
    Size:  1024,
    Usage: gpu.BufferUsageVertex,
}
buffer, err := allocator.CreateBuffer(bufferDesc)

// Create a texture
textureDesc := render.TextureDescriptor{
    Label:         "Color Texture",
    Size:          gpu.Extent3D{Width: 512, Height: 512, DepthOrArrayLayers: 1},
    Format:        gpu.TextureFormatRGBA8Unorm,
    Usage:         gpu.TextureUsageTextureBinding | gpu.TextureUsageRenderAttachment,
    SampleCount:   1,
    MipLevelCount: 1,
    Dimension:     gpu.TextureDimension2D,
}
texture, err := allocator.CreateTexture(textureDesc)
```

#### Shader Library
```go
shaderLibrary := context.GetShaderLibrary()

// Compile a shader
shaderSource := `
    @vertex
    fn vs_main() -> @builtin(position) vec4<f32> {
        return vec4<f32>(0.0, 0.0, 0.0, 1.0);
    }
`
shader, err := shaderLibrary.CompileShader(shaderSource, render.ShaderStageVertex)
```

#### Pipeline Library
```go
pipelineLibrary := context.GetPipelineLibrary()

// Create a render pipeline
pipelineDesc := render.RenderPipelineDescriptor{
    Label: "Basic Pipeline",
    Vertex: render.VertexState{
        Module:     vertexShader,
        EntryPoint: "vs_main",
    },
    Fragment: &render.FragmentState{
        Module:     fragmentShader,
        EntryPoint: "fs_main",
        Targets: []render.ColorTargetState{
            {
                Format:    gpu.TextureFormatBGRA8Unorm,
                WriteMask: gpu.ColorWriteMaskAll,
            },
        },
    },
    Primitive: render.PrimitiveState{
        Topology: gpu.PrimitiveTopologyTriangleList,
        CullMode: gpu.CullModeBack,
    },
}

pipeline, err := pipelineLibrary.CreateRenderPipeline(pipelineDesc)
```

### Entity System

The entity system converts display list operations into GPU-renderable entities:

```go
// Create entities manually
bounds := geom.Rect[geom.F32]{Left: 0, Top: 0, Right: 100, Bottom: 100}
color := display.NewColor(1.0, 0.0, 0.0, 1.0) // Red
entity := render.NewSolidColorEntity(color, bounds)

// Set transformation
transform := geom.NewMatrix[geom.F32]().
    Translate(geom.Vector3[geom.F32]{X: 50, Y: 50, Z: 0}).
    Scale(geom.Vector3[geom.F32]{X: 2, Y: 2, Z: 1})
entity.SetTransform(transform)
```

### Render Passes

Render passes manage GPU command recording:

```go
// Create a render target
colorAttachment := render.ColorAttachment{
    View:       surface.GetTextureView(),
    LoadOp:     gpu.LoadOpClear,
    StoreOp:    gpu.StoreOpStore,
    ClearValue: gpu.Color{R: 0.0, G: 0.0, B: 0.0, A: 1.0},
}

renderTarget := &render.RenderTarget{
    colorAttachments: []render.ColorAttachment{colorAttachment},
    size:            surface.GetSize(),
    sampleCount:     1,
}

// Create and use a render pass
renderPass, err := render.NewRenderPass(context, renderTarget, "Custom Pass")
if err != nil {
    log.Fatal(err)
}

// Set pipeline and draw
renderPass.SetPipeline(pipeline)
renderPass.SetVertexBuffer(0, vertexBuffer, 0, 256)
renderPass.Draw(3, 1, 0, 0) // Draw a triangle

// Finish and submit
commandBuffer := renderPass.End()
context.GetQueue().Submit([]gpu.CommandBuffer{commandBuffer})
```

## Content Rendering

The `ContentRenderer` handles conversion of display lists to renderable entities:

```go
contentRenderer := render.NewContentRenderer(context)

// Render a display list
err := contentRenderer.RenderDisplayList(displayList, renderTarget)
if err != nil {
    log.Fatal(err)
}
```

## Specialized Renderers

### Solid Fill Renderer
```go
solidRenderer := render.NewSolidFillRenderer(context)
err := solidRenderer.RenderSolidFill(renderPass, bounds, color, transform)
```

### Text Renderer
```go
textRenderer := render.NewTextRenderer(context)
style := render.TextStyle{
    FontSize:   16.0,
    Color:      display.NewColor(0.0, 0.0, 0.0, 1.0),
    FontWeight: render.FontWeightNormal,
    FontStyle:  render.FontStyleNormal,
}
err := textRenderer.RenderText(renderPass, "Hello, World!", position, style, transform)
```

### Image Renderer
```go
imageRenderer := render.NewImageRenderer(context)
srcRect := geom.Rect[geom.F32]{Left: 0, Top: 0, Right: 256, Bottom: 256}
dstRect := geom.Rect[geom.F32]{Left: 50, Top: 50, Right: 306, Bottom: 306}
err := imageRenderer.RenderImage(renderPass, texture, srcRect, dstRect, paint, transform)
```

## Capabilities and Features

Query rendering capabilities:

```go
capabilities := context.GetCapabilities()

if capabilities.SupportsAdvancedBlends {
    // Use advanced blend modes
}

if capabilities.SupportsCompute {
    // Use compute shaders for complex effects
}

maxTextureSize := capabilities.MaxTextureSize
maxSamples := capabilities.MaxSampleCount
```

## Performance Considerations

1. **Resource Pooling**: The render package automatically pools frequently used resources like buffers and textures.

2. **Pipeline Caching**: Render pipelines are cached to avoid expensive recompilation.

3. **Batch Rendering**: The entity system automatically batches similar rendering operations.

4. **GPU Memory Management**: The resource allocator efficiently manages GPU memory allocation and deallocation.

## Integration with Display Lists

The render package seamlessly integrates with Sraph's display list system:

```go
// Create a display list builder
builder := display.NewDisplayListBuilder()

// Add operations
builder.DrawRect(rect, paint)
builder.DrawText("Hello", position, textStyle)
builder.DrawImage(image, srcRect, dstRect, paint)

// Build the display list
displayList := builder.Build()

// Render it
err := renderer.Render(surface, displayList)
```

## Error Handling

The render package provides comprehensive error handling:

```go
if err != nil {
    switch {
    case errors.Is(err, render.ErrInvalidSurface):
        // Handle invalid surface
    case errors.Is(err, render.ErrOutOfMemory):
        // Handle GPU memory exhaustion
    case errors.Is(err, render.ErrDeviceLost):
        // Handle device loss and recovery
    default:
        // Handle other errors
    }
}
```

## Best Practices

1. **Resource Lifetime Management**: Always properly dispose of resources when they're no longer needed.

2. **Context Sharing**: Share rendering contexts across multiple renderers when possible.

3. **Surface Sizing**: Ensure render surfaces match the actual display size for optimal performance.

4. **Pipeline Reuse**: Reuse render pipelines across multiple render operations to minimize GPU state changes.

5. **Memory Monitoring**: Monitor GPU memory usage, especially when dealing with large textures or buffers.

## Future Enhancements

The render package is designed to be extensible. Future enhancements may include:

- Advanced post-processing effects
- Real-time ray tracing integration
- Machine learning acceleration
- Cross-platform compute shader support
- Enhanced text rendering with complex scripts

## See Also

- [Display Package Documentation](../display/README.md)
- [GPU Package Documentation](../gpu/README.md)
- [Geometry Package Documentation](../geom/README.md)
