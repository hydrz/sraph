# Entity Package

The `entity` package provides a complete entity system for Sraph's rendering pipeline, inspired by Flutter's Impeller entity system.

## Architecture

The entity system is responsible for transforming high-level drawing operations into low-level GPU commands, managing state, resources, and optimization along the way.

### Core Components

#### Entity System
- `entity.go` - Core entity interface and default implementation
- `entity_pass_clip_stack.go` - Clip stack management for entity passes
- `entity_pass_target.go` - Render target management for entity passes
- `entity_playground.go` - Testing and development environment

#### Content System
- `content_context.go` - Graphics context for content rendering
- `contents/` - Directory containing all content types:
  - `contents.go` - Base content interface and implementations
  - `solid_color_contents.go` - Solid color fills
  - `texture_contents.go` - Textured content rendering
  - `gradient_contents.go` - Linear and radial gradients
  - `text_contents.go` - Text rendering using glyph atlases
  - `clip_contents.go` - Clipping operations

#### Rendering Infrastructure
- `render_target_cache.go` - Efficient render target caching
- `inline_pass_context.go` - Inline rendering pass management
- `save_layer_utils.go` - Save layer utilities and coverage computation
- `draw_order_resolver.go` - Draw order optimization

#### Geometry System
- `geometry/` - Directory containing geometric primitives:
  - `geometry.go` - Base geometry types and interfaces
  - `gradient.go` - Gradient geometry definitions
  - `solid.go` - Solid geometry shapes
  - `text.go` - Text geometry layouts
  - `texture.go` - Textured geometry

#### Shader System
- `shaders/` - Directory containing entity shaders:
  - `entity_shaders.go` - WGSL shaders for all entity content types

## Key Features

### Blend Modes
The entity system supports comprehensive blend mode operations including:
- Basic blend modes (SrcOver, Multiply, Screen, etc.)
- Advanced blend modes (ColorDodge, ColorBurn, HardLight, etc.)
- Porter-Duff composition modes

### Content Types
- **Solid Colors**: Efficient solid color fills with anti-aliasing
- **Textures**: Textured content with sampling modes and source/destination rectangles
- **Gradients**: Linear and radial gradients with multiple color stops and tile modes
- **Text**: High-quality text rendering using glyph atlases and font metrics
- **Clipping**: Stencil-based clipping with intersection and difference operations

### Performance Optimizations
- **Render Target Caching**: Efficient reuse of render targets across frames
- **Draw Order Resolution**: Optimal ordering of draw calls for transparency and clipping
- **Resource Management**: Automatic management of GPU resources and memory
- **Clip Stack Management**: Efficient clip state tracking and restoration

### Advanced Features
- **Save Layer Support**: Advanced layer composition with image filters
- **Coverage Computation**: Precise coverage calculation for optimal rendering
- **Anti-aliasing**: High-quality anti-aliased rendering for all content types
- **Playground Environment**: Development and testing utilities

## Usage

The entity system is designed to be used through the high-level widget and element systems, but can also be used directly for custom rendering:

```go
// Create content context
context := entity.NewContentContext(gpuContext, typographerContext)

// Create entity with solid color content
solidContent := contents.NewSolidColorContents()
solidContent.SetColor(geom.ColorRed)
solidContent.SetPath(path)

entity := entity.NewEntity()
entity.SetContents(solidContent)
entity.SetTransform(transform)

// Render entity
entity.Render(context, renderPass)
```

## Dependencies

The entity package depends on:
- `geom` - Geometric primitives and math utilities
- `gpu` - GPU abstraction layer
- `render` - Low-level rendering primitives
- `font` - Typography and text rendering
- `shader` - Shader compilation and management

## Implementation Status

This implementation provides the complete architecture and structure for the entity system, with:
- ✅ Complete interface definitions
- ✅ Core entity and content system
- ✅ Shader definitions for all content types
- ✅ Resource management infrastructure
- ✅ Render target caching
- ✅ Clip stack management
- ⚠️  Implementation placeholders (marked with TODO comments)

The system is designed to be incrementally implemented while maintaining architectural consistency.
