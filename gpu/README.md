# GPU Package - WebGPU Go Interface

A Go interface for WebGPU, enabling GPU-accelerated graphics and computation. This package provides a comprehensive Go API for WebGPU operations, allowing developers to leverage modern GPU capabilities for rendering and compute tasks.

## Features

- **Complete WebGPU API Coverage**: Full implementation of WebGPU specification
- **Cross-Platform Support**: Works across different operating systems and GPU backends
- **Type-Safe Interface**: Strongly typed Go interfaces for all WebGPU operations
- **Memory Management**: Proper resource lifecycle management
- **Error Handling**: Comprehensive error reporting and handling
- **Performance Optimized**: Minimal overhead between Go and native WebGPU calls

## Supported Backends

- **Vulkan**: High-performance graphics and compute API
- **Metal**: Apple's graphics API for macOS and iOS
- **DirectX 12**: Microsoft's graphics API for Windows
- **OpenGL/OpenGL ES**: Legacy graphics API support

## Installation

```bash
go get github.com/opensraph/sraph/gpu
```

## Quick Start

### Basic Setup

```go
package main

import (
    "fmt"
    "log"

    "github.com/opensraph/sraph/gpu"
)

func main() {
    // Request a GPU adapter
    adapter, err := gpu.RequestAdapter(gpu.RequestAdapterOptions{
        PowerPreference: gpu.PowerPreferenceHighPerformance,
    })
    if err != nil {
        log.Fatal("Failed to request adapter:", err)
    }

    // Get adapter information
    info, err := adapter.Info()
    if err != nil {
        log.Fatal("Failed to get adapter info:", err)
    }

    fmt.Printf("Using GPU: %s (%s)\n", info.Device, info.Vendor)

    // Request a device
    device, err := adapter.RequestDevice(gpu.DeviceDescriptor{
        Label: "Main Device",
    })
    if err != nil {
        log.Fatal("Failed to request device:", err)
    }

    // Get the command queue
    queue := device.Queue()

    fmt.Println("WebGPU setup complete!")
}
```

### Creating and Using Buffers

```go
// Create a buffer
buffer := device.CreateBuffer(gpu.BufferDescriptor{
    Label: "Vertex Buffer",
    Size:  1024,
    Usage: gpu.BufferUsageVertex | gpu.BufferUsageCopyDst,
})

// Write data to buffer
data := []float32{1.0, 2.0, 3.0, 4.0}
queue.WriteBuffer(buffer, 0, unsafe.Pointer(&data[0]), uintptr(len(data)*4))

// Clean up
defer buffer.Destroy()
```

### Basic Render Pipeline

```go
// Create shader module
shaderModule := device.CreateShaderModule(gpu.ShaderModuleDescriptor{
    Label: "Basic Shader",
    // Add shader source here
})

// Create render pipeline
pipeline := device.CreateRenderPipeline(gpu.RenderPipelineDescriptor{
    Label: "Basic Pipeline",
    Vertex: gpu.VertexState{
        Module:     shaderModule,
        EntryPoint: "vs_main",
    },
    Fragment: gpu.FragmentState{
        Module:     shaderModule,
        EntryPoint: "fs_main",
        Targets: []gpu.ColorTargetState{
            {
                Format:    gpu.TextureFormatBGRA8UnormSrgb,
                WriteMask: gpu.ColorWriteMaskAll,
            },
        },
    },
    Primitive: gpu.PrimitiveState{
        Topology: gpu.PrimitiveTopologyTriangleList,
    },
})
```

## API Overview

### Core Workflow

The typical WebGPU workflow in Go follows these steps:

1. **Initialization**
   ```go
   adapter, _ := gpu.RequestAdapter(options)
   device, _ := adapter.RequestDevice(descriptor)
   queue := device.Queue()
   ```

2. **Resource Creation**
   ```go
   buffer := device.CreateBuffer(bufferDesc)
   texture := device.CreateTexture(textureDesc)
   sampler := device.CreateSampler(samplerDesc)
   ```

3. **Pipeline Setup**
   ```go
   shaderModule := device.CreateShaderModule(shaderDesc)
   pipeline := device.CreateRenderPipeline(pipelineDesc)
   bindGroup := device.CreateBindGroup(bindGroupDesc)
   ```

4. **Command Recording**
   ```go
   encoder := device.CreateCommandEncoder(encoderDesc)
   renderPass := encoder.BeginRenderPass(renderPassDesc)
   renderPass.SetPipeline(pipeline)
   renderPass.SetBindGroup(0, bindGroup, nil)
   renderPass.Draw(3, 1, 0, 0)
   renderPass.End()
   commandBuffer := encoder.Finish(gpu.CommandBufferDescriptor{})
   ```

5. **Execution**
   ```go
   queue.Submit([]gpu.CommandBuffer{commandBuffer})
   ```

### Key Interfaces

- **Instance**: Entry point for WebGPU operations
- **Adapter**: Represents a physical GPU
- **Device**: Logical GPU device for resource creation
- **Queue**: Command submission and data transfer
- **Buffer**: GPU memory buffers
- **Texture**: 2D/3D image data
- **Pipeline**: Shader programs and GPU state
- **BindGroup**: Resource binding for shaders

## Surface Integration

For window rendering, create and configure a surface:

```go
// Create surface (platform-specific)
surface, err := instance.CreateSurface(gpu.SurfaceDescriptor{
    Label: "Main Window",
    // Platform-specific surface source
})

// Configure surface
surface.Configure(gpu.SurfaceConfiguration{
    Device:      device,
    Format:      gpu.GetPreferredCanvasFormat(),
    Usage:       gpu.TextureUsageRenderAttachment,
    Width:       800,
    Height:      600,
    PresentMode: gpu.PresentModeFifo,
})

// Get current texture for rendering
surfaceTexture := surface.CurrentTexture()
defer surface.Present()
```

## Error Handling

The package provides comprehensive error handling:

```go
// Check for device loss
future := device.LostFuture()
// Handle device lost events appropriately

// Push/pop error scopes for debugging
device.PushErrorScope(gpu.ErrorFilterValidation)
// ... GPU operations ...
device.PopErrorScope()
```

## Memory Management

Proper resource cleanup is essential:

```go
// Always destroy resources when done
defer buffer.Destroy()
defer texture.Destroy()
defer device.Destroy()

// Unmap buffers after use
defer buffer.Unmap()

// Release command buffers after submission
// (handled automatically by the implementation)
```

## Features and Limits

Query device capabilities:

```go
// Check supported features
features := device.Features()
if device.HasFeature(gpu.FeatureNameTimestampQuery) {
    // Use timestamp queries
}

// Check device limits
limits, err := device.Limits()
if err == nil {
    fmt.Printf("Max texture size: %d\n", limits.MaxTextureDimension2D)
}
```

## Compute Shaders

Example compute pipeline usage:

```go
// Create compute pipeline
computePipeline := device.CreateComputePipeline(gpu.ComputePipelineDescriptor{
    Label: "Compute Pipeline",
    Compute: gpu.ComputeState{
        Module:     computeShader,
        EntryPoint: "cs_main",
    },
})

// Dispatch compute work
encoder := device.CreateCommandEncoder(gpu.CommandEncoderDescriptor{})
computePass := encoder.BeginComputePass(gpu.ComputePassDescriptor{})
computePass.SetPipeline(computePipeline)
computePass.SetBindGroup(0, computeBindGroup, nil)
computePass.DispatchWorkgroups(64, 1, 1)
computePass.End()
```

## Best Practices

1. **Resource Lifecycle**: Always destroy resources when no longer needed
2. **Error Handling**: Check all error returns and handle appropriately
3. **Buffer Management**: Use appropriate buffer usage flags
4. **Pipeline State**: Minimize pipeline state changes for better performance
5. **Command Batching**: Batch similar operations in command buffers
6. **Memory Alignment**: Follow GPU memory alignment requirements

## Examples

See the `examples/` directory for complete working examples:

- **Triangle**: Basic triangle rendering
- **Texture**: Texture loading and sampling
- **Compute**: Compute shader examples
- **Instancing**: Instance rendering techniques

## Platform Support

| Platform | Backend | Status      |
| -------- | ------- | ----------- |
| Windows  | D3D12   | ✅ Supported |
| Windows  | Vulkan  | ✅ Supported |
| macOS    | Metal   | ✅ Supported |
| Linux    | Vulkan  | ✅ Supported |
| Web      | WebGPU  | ✅ Supported |

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the BSD 3-Clause License. See LICENSE file for details.

## Acknowledgments

- WebGPU specification by the W3C GPU for the Web Community Group
- WebGPU-Native project for the C API specification
- Contributors to the WebGPU ecosystem

## Links

- [WebGPU Specification](https://www.w3.org/TR/webgpu/)
- [WebGPU Shading Language](https://www.w3.org/TR/WGSL/)
- [Learn WebGPU](https://eliemichel.github.io/LearnWebGPU/)
- [WebGPU Samples](https://webgpu.github.io/webgpu-samples/)
