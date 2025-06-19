# Shader Package

The shader package provides comprehensive shader management and compilation capabilities for the Sraph UI toolkit. It handles loading, compiling, bundling, and managing GPU shaders including built-in shaders for common rendering operations.

## Features

### Shader Compilation
- **WGSL Compiler**: Full support for WebGPU Shading Language (WGSL)
- **Preprocessing**: Support for includes, defines, and template processing
- **Optimization**: Multiple optimization levels (None, Basic, Full)
- **Validation**: Syntax and semantic validation of shader code
- **Error Reporting**: Detailed compilation diagnostics

### Shader Management
- **Shader Library**: Centralized management of all shaders
- **Built-in Shaders**: Pre-compiled shaders for common operations (solid color, texture, gradient)
- **Dynamic Loading**: Load shaders from files or embedded sources
- **Resource Management**: Automatic GPU resource cleanup
- **Shader Programs**: Combine vertex and fragment shaders into programs

### Shader Bundling
- **Bundle Creation**: Create packages of multiple shaders
- **Asset Management**: Manage shader assets and dependencies
- **Serialization**: Save and load shader bundles to/from files
- **Registry System**: Global registry for shader bundles

### Utilities
- **Type Detection**: Automatic shader type detection from source
- **Uniform Extraction**: Parse uniform variables from shader source
- **Template Processing**: Support for shader templates with replacements
- **Validation Tools**: Comprehensive source validation

## Architecture

```
shader/
├── builtin.go          # Built-in shader sources and metadata
├── builtin/            # Built-in shader source files (.wgsl)
│   ├── solid.wgsl      # Solid color rendering
│   ├── texture.wgsl    # Texture sampling
│   └── gradient.wgsl   # Linear gradient rendering
├── bundle.go           # Shader bundling and asset management
├── compiler.go         # WGSL compiler implementation
├── library.go          # Shader management and loading
└── utils.go            # Utilities and helper functions
```

## Usage Examples

### Basic Shader Loading
```go
// Create shader manager
device := gpu.CreateDevice()
manager := shader.NewDefaultShaderManager(device)

// Load shader from file
shader, err := manager.LoadShaderFromFile("vertex.wgsl", shader.ShaderTypeVertex)
if err != nil {
    log.Fatal(err)
}

// Get built-in shader
solidShader := manager.GetBuiltinShader("solid")
```

### Shader Compilation
```go
// Create compiler
compiler := shader.NewWGSLCompiler()
compiler.SetOptimizationLevel(shader.OptimizationLevelBasic)
compiler.AddDefine("USE_TEXTURE", "1")

// Compile shader
source := shader.ShaderSource{
    Name:   "custom",
    Source: wgslCode,
    Type:   shader.ShaderTypeFragment,
}

compiled, err := compiler.Compile(source)
if err != nil {
    log.Fatal(err)
}
```

### Shader Bundling
```go
// Create bundle builder
compiler := shader.NewWGSLCompiler()
builder := shader.NewDefaultBundleBuilder(compiler)

// Add shaders to bundle
builder.SetName("MyShaders")
builder.AddShaderFromFile("vertex", "shaders/vertex.wgsl", shader.ShaderTypeVertex)
builder.AddShaderFromFile("fragment", "shaders/fragment.wgsl", shader.ShaderTypeFragment)

// Build bundle
bundle, err := builder.Build()
if err != nil {
    log.Fatal(err)
}

// Use shaders from bundle
vertexShader, err := bundle.GetShader("vertex")
```

### Shader Programs
```go
// Create shader program from vertex and fragment shaders
program, err := manager.CreateProgram(vertexShader, fragmentShader)
if err != nil {
    log.Fatal(err)
}

// Use program for rendering
pipeline := program.RenderPipeline()
```

## Built-in Shaders

### Solid Color Shader
- **Purpose**: Renders solid colors and simple shapes
- **Uniforms**: Color value
- **Usage**: UI elements, solid fills

### Texture Shader
- **Purpose**: Renders textured quads and sprites
- **Uniforms**: Texture sampler, texture coordinates
- **Usage**: Images, textured surfaces

### Gradient Shader
- **Purpose**: Renders linear gradients
- **Uniforms**: Start/end colors, gradient points
- **Usage**: Gradient backgrounds, smooth transitions

## API Reference

### Core Interfaces

#### ShaderManager
Main interface for shader operations:
- `LoadShader(source, type)` - Load shader from source code
- `LoadShaderFromFile(filename, type)` - Load shader from file
- `GetBuiltinShader(name)` - Get built-in shader by name
- `CreateProgram(vertex, fragment)` - Create shader program

#### Compiler
Shader compilation interface:
- `Compile(source)` - Compile shader source
- `PreprocessShader(source)` - Preprocess with includes/defines
- `ValidateShader(source)` - Validate shader syntax
- `OptimizeShader(shader)` - Optimize compiled shader

#### Bundle
Shader bundle interface:
- `GetShader(name)` - Get shader from bundle
- `ListShaders()` - List all shaders in bundle
- `HasShader(name)` - Check if shader exists

### Key Types

#### ShaderType
- `ShaderTypeVertex` - Vertex shader
- `ShaderTypeFragment` - Fragment/pixel shader
- `ShaderTypeCompute` - Compute shader

#### OptimizationLevel
- `OptimizationLevelNone` - No optimization
- `OptimizationLevelBasic` - Basic optimizations
- `OptimizationLevelFull` - Aggressive optimizations

#### UniformType
- `UniformTypeFloat`, `UniformTypeVec2/3/4` - Scalar and vector types
- `UniformTypeMat2/3/4` - Matrix types
- `UniformTypeSampler2D`, `UniformTypeSamplerCube` - Texture samplers

## Integration

The shader package integrates with:
- **GPU Package**: Low-level WebGPU interface
- **Render Package**: High-level rendering pipeline
- **Entity Package**: Renderable objects and materials
- **Display Package**: Display list rendering

## Future Enhancements

- Shader caching and hot-reloading
- More built-in shaders (phong, PBR, etc.)
- Shader variant system for different quality levels
- Cross-compilation to other shader languages
- Shader debugging and profiling tools
