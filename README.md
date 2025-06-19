# Sraph

Sraph is a UI toolkit for building beautiful, natively compiled applications for mobile, web, desktop, and embedded devices using the Go programming language.

## Directory Structure

```
sraph/
├── version/                     # Version info
│   └── version.go
│
├── geom/                        # Mathematical primitives and geometry operations
│   ├── scalar.go
│   ├── point.go
│   ├── vector.go                # Vector operations
│   ├── rect.go
│   ├── matrix.go
│   ├── color.go
│   ├── path.go
│   ├── rstransform.go
│   └── ...
│
├── tess/                        # Path tessellation and geometry processin
│
├── display/                     # Display list and drawing operations
│
├── render/                      # High-level rendering pipeline
│
├── entity/                      # Renderable entities and content
│
├── gpu/                         # GPU backends and hardware abstraction
│   ├── gpu.go
│   ├── wgpu.go                  # WebGPU standard API
│   ├── backend.go               # Backend interface and registration
│   ├── vulkan/                  # Vulkan backend support
│   ├── metal/                   # Metal backend support
│   ├── gles/                    # OpenGL ES backend support
│   └── gen/                     # Code generation tools
│       └── main.go
│
├── shader/                      # Shader compilation and management
│   ├── library.go
│   ├── compiler.go
│   ├── bundle.go
│   └── builtin/                 # Built-in shaders
│       ├── solid.wgsl
│       ├── texture.wgsl
│       └── gradient.wgsl
│
├── font/                        # Font loading and text layout
│
├── gio/                         # Window management and platform integration
│   ├── gio.go
│   ├── driver.go                # Platform driver registration
│   ├── window.go
│   ├── surface.go
│   ├── event.go
│   ├── key.go
│   ├── pointer.go
│   ├── clipboard.go
│   ├── wayland/                 # Wayland support
│   ├── x11/                     # X11 support
│   ├── win32/                   # Windows support
│   ├── cocoa/                   # macOS support
│   └── android/                 # Android support
│
├── widget/                      # UI widgets and layout
│
├── app/                         # Application framework and runtime
│
│
├── examples/                    # Example applications
│   ├── hello/
│   ├── shapes/
│   ├── text/
└──   └── animation/

```

## Architecture Overview

```
flowchart TD

%% Application Layer
subgraph "Application"
    App["<b>app</b><br/>(engine, scheduler, state)"]
    Examples["<b>examples</b><br/>(hello, shapes, text, animation)"]
end

%% UI Layer
subgraph "UI"
    Widget["<b>widget</b><br/>(button, text, container, layout)"]
end

%% Content Layer
subgraph "Content"
    Display["<b>display</b><br/>(DisplayList, Builder, Canvas, Paint)"]
    Entity["<b>entity</b><br/>(solid, texture, gradient, text)"]
    Font["<b>font</b><br/>(font, glyph, atlas, layout, paragraph)"]
end

%% Geometry & Tessellation Layer
subgraph "Geometry & Tessellation"
    Geom["<b>geom</b><br/>(scalar, point, matrix, path)"]
    Tess["<b>tess</b><br/>(path_tess, polygon, stroke, curve)"]
end

%% Rendering Pipeline Layer
subgraph "Rendering Pipeline"
    Render["<b>render</b><br/>(renderer, context, pipeline, surface)"]
end

%% Graphics Backend Layer
subgraph "Graphics Backend"
    Gpu["<b>gpu</b><br/>(webgpu, backend, vulkan, metal, gles)"]
    Shader["<b>shader</b><br/>(shader_comp, wgsl, builtin)"]
end

%% Platform Layer
subgraph "Platform"
    Gio["<b>gio</b><br/>(window, wayland, x11, win32, cocoa, android)"]
end

%% System Layer
Sys["<b>System</b><br/>(OS APIs & GPU Drivers)"]

%% Main dependency chain
App --> Widget
Widget --> Display
Display --> Entity
Display --> Font
Display --> Geom
Display --> Tess
Entity --> Geom
Font --> Geom
Display --> Render
Geom --> Tess
Tess --> Render
Render --> Gpu
Render --> Shader
Gpu --> Gio
Gio --> Sys
Shader --> Gpu

%% Supplementary
Examples -.-> App

%% Styling
classDef layer fill:#f8fbff,stroke:#89b4fa,stroke-width:2px;
classDef node fill:#f3f8fd,stroke:#7ea6e0,stroke-width:1.5px,font-size:15px;
classDef example fill:#faf3e3,stroke:#dbb97a,stroke-width:1.2px,stroke-dasharray: 5 3,font-size:14px;

class App,Widget,Display,Entity,Font,Geom,Tess,Render,Gpu,Shader,Gio,Sys node;
class Examples example;
```

## Features

- High-performance 2D/3D geometry and rendering
- Display list architecture for efficient drawing
- Advanced text, image, and visual effects support

## Getting Started

```bash
go get github.com/opensraph/sraph
```

## Changelog

[![release](https://github.com/opensraph/sraph/actions/workflows/release.yml/badge.svg)](https://github.com/opensraph/sraph/releases)


## Contributing

We welcome contributions! Please read the [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[![License: MIT](https://opensource.org/licenses/MIT)
This project is licensed under the terms of the MIT license. See the [LICENSE](LICENSE) file for details.
