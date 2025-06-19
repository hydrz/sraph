# Sraph Font Package

The `font` package provides comprehensive typography and text rendering functionality for the Sraph UI toolkit. It is inspired by Flutter's Impeller typographer architecture but adapted for Go and Sraph's needs.

## Architecture Overview

The font package is structured around several key components that work together to provide complete text rendering capabilities:

### Core Components

#### 1. Typeface and Font System
- **Typeface**: Represents the intrinsic properties of a font family (loaded from font files)
- **Font**: Combines a typeface with size, metrics, and styling modifications
- **TypefaceLoader**: Handles loading typefaces from various sources (files, memory)
- **FontManager**: Manages font registration, discovery, and caching

#### 2. Text Layout and Shaping
- **TextLayout**: Handles text shaping and layout operations
- **TextRun**: Represents a collection of positioned glyphs from a specific font
- **TextFrame**: Represents a collection of shaped text runs (main entry point for text rendering)
- **Paragraph**: Manages multiple text runs with styling and alignment
- **LineBreaker**: Handles line breaking and word wrapping

#### 3. Glyph Management and Rendering
- **GlyphAtlas**: Manages texture atlases for rendered glyphs
- **GlyphRenderer**: Handles rendering of glyphs to bitmaps and atlas management
- **GlyphRasterizer**: Converts glyph paths to bitmaps
- **GlyphCache**: Provides caching for rendered glyphs
- **RectanglePacker**: Efficiently packs glyphs into atlas textures

#### 4. Typography Context
- **TypographerContext**: Provides graphics context for text rendering
- **LazyGlyphAtlas**: Lazy loading of glyph atlases
- **TextRenderer**: Handles rendering of text frames to render passes

## Key Features

### Font Loading and Management
- Support for multiple font formats (TTF, OTF, WOFF, WOFF2)
- Font registration and discovery system
- Family-based font matching with weight and style fallbacks
- Font caching and performance optimization

### Text Layout and Shaping
- Unicode text shaping and glyph positioning
- Line breaking and word wrapping
- Paragraph layout with alignment support
- Multi-run text frames for rich text

### Glyph Rendering
- Efficient glyph atlas management
- Alpha and color glyph support (for emoji)
- Subpixel positioning for crisp text
- Skyline-based rectangle packing for optimal atlas usage

### Performance Optimization
- Lazy loading of resources
- Comprehensive caching at multiple levels
- GPU-accelerated rendering pipeline integration
- Memory-efficient glyph storage

## File Structure

```
font/
├── font.go                 # Core types and interfaces
├── typeface_loader.go      # Font loading functionality
├── font_manager.go         # Font management and registration
├── text_run.go            # Text run implementation
├── text_frame.go          # Text frame implementation
├── text_layout.go         # Text layout and shaping
├── glyph_atlas.go         # Glyph atlas management
├── glyph_renderer.go      # Glyph rendering and caching
├── rectangle_packer.go    # Rectangle packing algorithms
└── typographer_context.go # Typography context and rendering
```

## Usage Examples

### Basic Font Loading and Text Rendering

```go
// Create font manager
fontManager := font.NewFontManager()

// Load and register a font
err := fontManager.RegisterFontFromFile("path/to/font.ttf", "MyFont")
if err != nil {
    log.Fatal(err)
}

// Get the font
myFont, ok := fontManager.GetFont("MyFont")
if !ok {
    log.Fatal("Font not found")
}

// Create text layout
layout := font.NewTextLayout()

// Layout text
frame, err := layout.LayoutText("Hello, World!", myFont, 300)
if err != nil {
    log.Fatal(err)
}

// Render the text frame (requires render context)
// err = typographerContext.RenderTextFrame(frame, renderPass)
```

### Advanced Typography

```go
// Create a paragraph with multiple text runs
paragraph := font.NewParagraph()
paragraph.SetAlignment(font.TextAlignmentCenter)

// Add different styled runs
run1 := font.NewTextRun(boldFont)
run2 := font.NewTextRun(italicFont)

// Layout the paragraph
frame, err := layout.LayoutParagraph(paragraph, 400)
if err != nil {
    log.Fatal(err)
}
```

## Integration with Sraph

The font package integrates seamlessly with other Sraph components:

- **Render Package**: Text frames are rendered using the render pipeline
- **Display Package**: Text elements are part of the display list system
- **Geom Package**: Uses geometric types for positioning and measurements
- **GPU Package**: Leverages GPU acceleration for glyph rendering

## Implementation Status

The font package provides a complete architecture with the following implementation status:

- ✅ **Core interfaces and types**: Complete
- ✅ **Font loading and management**: Basic implementation
- ✅ **Text layout framework**: Complete structure
- ✅ **Glyph atlas management**: Core functionality
- ⚠️ **Actual font parsing**: TODO (placeholder implementation)
- ⚠️ **Glyph rasterization**: TODO (placeholder implementation)
- ⚠️ **Advanced text shaping**: TODO (basic implementation)
- ⚠️ **GPU integration**: TODO (interface defined)

## TODOs and Future Work

### High Priority
1. **Font Parsing**: Implement actual TTF/OTF parsing for glyph extraction
2. **Glyph Rasterization**: Implement path-to-bitmap conversion
3. **GPU Integration**: Complete integration with render pipeline
4. **Text Shaping**: Implement proper Unicode text shaping

### Medium Priority
1. **Advanced Features**: Ligatures, kerning, complex scripts
2. **Performance**: Optimize atlas packing and caching
3. **Font Fallback**: Implement comprehensive fallback system
4. **Metrics**: Improve font metrics calculation

### Low Priority
1. **Font Subsetting**: Support for font subsetting
2. **Variable Fonts**: Support for variable font technology
3. **Color Fonts**: Enhanced color font support
4. **Hinting**: Font hinting for better rendering at small sizes

## Design Principles

The font package follows these key design principles:

1. **Modularity**: Clear separation of concerns with well-defined interfaces
2. **Performance**: Lazy loading and comprehensive caching
3. **Extensibility**: Easy to extend with new font formats and features
4. **Go Idioms**: Follows Go best practices and conventions
5. **GPU-First**: Designed for GPU-accelerated rendering
6. **Unicode Support**: Full Unicode text processing capabilities

This architecture provides a solid foundation for high-quality text rendering in the Sraph UI toolkit while maintaining flexibility for future enhancements and optimizations.
