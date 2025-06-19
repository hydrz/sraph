# Gio - Cross-Platform Graphics I/O Library

[![Go Reference](https://pkg.go.dev/badge/github.com/opensraph/sraph/gio.svg)](https://pkg.go.dev/github.com/opensraph/sraph/gio)
[![Go Report Card](https://goreportcard.com/badge/github.com/opensraph/sraph/gio)](https://goreportcard.com/report/github.com/opensraph/sraph/gio)

Gio is a cross-platform graphics input/output library written in Go, providing a unified interface for window management, event handling, and graphics operations across different platforms.

## Features

- **Cross-Platform Support**: Works on Linux (X11/Wayland), macOS (Cocoa), Windows (Win32), and mobile platforms
- **Automatic Driver Selection**: Automatically selects the best available graphics driver for your platform
- **Event System**: Comprehensive event handling for keyboard, mouse, wheel, clipboard, and drag-and-drop operations
- **Window Management**: Create and manage windows with flexible attributes and states
- **Type-Safe APIs**: Strong typing with enums and interfaces for better code safety
- **Concurrent Event Processing**: Asynchronous event handling with configurable timeouts
- **Extensible Architecture**: Plugin-based driver system for easy platform extensions

## Supported Platforms

| Platform | Driver | Status |
|----------|--------|--------|
| Linux    | X11    | ✅ Supported |
| Linux    | Wayland| ✅ Supported |
| macOS    | Cocoa  | 🚧 Planned |
| Windows  | Win32  | 🚧 Planned |
| Android  | Mobile | 🚧 Planned |
| iOS      | Mobile | 🚧 Planned |

## Installation

```bash
go get github.com/opensraph/sraph/gio
```

For platform-specific drivers, import them in your application:

```go
import (
    "github.com/opensraph/sraph/gio"
    _ "github.com/opensraph/sraph/gio/x11"     // X11 support
    _ "github.com/opensraph/sraph/gio/wayland" // Wayland support
)
```

## Quick Start

### Basic Window Creation

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/opensraph/sraph/gio"
    _ "github.com/opensraph/sraph/gio/x11"
    _ "github.com/opensraph/sraph/gio/wayland"
)

func main() {
    // Create a window with default settings
    window, err := gio.CreateWindow(gio.NewWindowOptions{
        Title:  "My Application",
        Width:  800,
        Height: 600,
        State:  gio.WindowStateVisible | gio.WindowStateResizable,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Show the window
    if err := window.Show(); err != nil {
        log.Fatal(err)
    }

    // Keep the application running
    time.Sleep(10 * time.Second)
}
```

### Event Handling

```go
// Subscribe to keyboard events
window.Subscribe(gio.EventTypeKeyboard, func(e gio.Event) error {
    if keyEvent, ok := e.(*gio.KeyboardEvent); ok {
        fmt.Printf("Key pressed: %s\n", keyEvent.Key())

        // Exit on Escape key
        if keyEvent.Code == gio.KeyCodeEscape {
            return window.Close()
        }
    }
    return nil
})

// Subscribe to pointer events
window.Subscribe(gio.EventTypePointer, func(e gio.Event) error {
    if pointerEvent, ok := e.(*gio.PointerEvent); ok {
        fmt.Printf("Mouse at: %v\n", pointerEvent.Position)
    }
    return nil
})
```

### Platform Information

```go
// Get platform information
platformInfo := gio.GetPlatformInfo()
fmt.Printf("Platform: %s\n", platformInfo.String())

// List available drivers
drivers := gio.EnumerateDrivers()
for _, driver := range drivers {
    fmt.Printf("Driver: %s (available: %v)\n",
        driver.String(), driver.IsAvailable())
}
```

## Architecture

Gio follows a modular architecture with the following components:

### Core Components

- **Driver System**: Pluggable backend system for different platforms
- **Window Management**: Cross-platform window creation and management
- **Event System**: Type-safe event handling with pub/sub pattern
- **Configuration**: Global configuration management

### Driver Architecture

```
┌─────────────────┐
│   Application   │
└─────────────────┘
         │
┌─────────────────┐
│   Gio Core      │
├─────────────────┤
│ Window Manager  │
│ Event System    │
│ Driver Registry │
└─────────────────┘
         │
┌─────────────────┐
│    Drivers      │
├─────────────────┤
│ X11   Wayland   │
│ Cocoa  Win32    │
└─────────────────┘
```

## API Reference

### Core Types

#### Window
```go
type Window interface {
    BaseWindow
    WindowID() WindowID
}

type BaseWindow interface {
    Width() int
    Height() int
    Title() string
    Position() image.Point
    State() WindowState

    Show() error
    Hide() error
    Close() error

    Subscribe(eventType EventType, handler EventHandler) error
    Unsubscribe(eventType EventType, handler EventHandler)
    Publish(event Event) error
}
```

#### Events
```go
type EventType uint8

const (
    EventTypeWindow
    EventTypeKeyboard
    EventTypePointer
    EventTypeWheel
    EventTypeClipboard
    EventTypeDrag
)
```

#### Window States
```go
type WindowState uint16

const (
    WindowStateClosed
    WindowStateFocused
    WindowStateIconified
    WindowStateMaximized
    WindowStateVisible
    WindowStateResizable
    WindowStateDecorated
    WindowStateFloating
)
```

### Key Functions

#### Window Creation
```go
// Create window with automatic driver selection
func CreateWindow(o ...NewWindowOptions) (Window, error)

// Create window with specific driver
func CreateWindowWithDriver(driverType DriverType, o ...NewWindowOptions) (Window, error)
```

#### Driver Management
```go
// Get best available driver
func GetBestDriver() (Driver, error)

// Get specific driver
func GetDriver(driverType DriverType) (Driver, error)

// Enumerate all drivers
func EnumerateDrivers() []DriverType
```

## Configuration

Gio can be configured globally:

```go
config := gio.GetConfig()
config.EventQueueSize = 2000
config.EventTimeout = time.Second * 10
config.DefaultWindowWidth = 1920
config.DefaultWindowHeight = 1080
gio.SetConfig(config)
```

## Examples

See the [examples](example/) directory for complete examples:

- [Simple Window](example/simple/main.go) - Basic window creation and event handling
- Advanced Graphics - Coming soon
- Multi-Window - Coming soon

## Error Handling

Gio provides structured error handling:

```go
if err := window.Show(); err != nil {
    if gio.IsError(err, gio.ErrorCodeWindowNotInitialized) {
        // Handle specific error
    }
    log.Printf("Window error: %v", err)
}
```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](../CONTRIBUTING.md) for details.

### Development Setup

1. Clone the repository
2. Install Go 1.24+
3. Install platform dependencies:
   - Linux: `libx11-dev`, `libwayland-dev`
   - macOS: Xcode command line tools
   - Windows: Visual Studio Build Tools

### Running Tests

```bash
go test ./...
```

### Building Examples

```bash
cd example/simple
go build -o simple .
./simple
```

## Roadmap

- [ ] Complete Windows (Win32) driver
- [ ] Complete macOS (Cocoa) driver
- [ ] Mobile platform support (Android/iOS)
- [ ] OpenGL/Vulkan integration
- [ ] High-DPI support
- [ ] Gamepad input support
- [ ] IME (Input Method Editor) support
- [ ] Accessibility features

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## Acknowledgments

- Inspired by [Gio UI](https://gioui.org/) and [GLFW](https://www.glfw.org/)
- Thanks to all contributors and the Go community

---

For questions and support, please [open an issue](https://github.com/opensraph/sraph/issues) or join our [discussions](https://github.com/opensraph/sraph/discussions).
