# Wayland Backend for Gio Graphics System

This directory contains the Wayland backend implementation for the Gio graphics system. The Wayland backend provides native Wayland support for creating and managing windows on Linux systems using the Wayland display server protocol.

## Features

### Core Features
- ✅ **Window Creation and Management**: Create and manage native Wayland windows
- ✅ **Event Handling**: Comprehensive keyboard, mouse, and scroll event support
- ✅ **Window States**: Support for window focus, visibility, maximization, and more
- ✅ **Multi-Window Support**: Handle multiple windows simultaneously
- ✅ **Driver Registration**: Automatic registration with the Gio driver system

### Window Operations
- ✅ **Basic Operations**: Show, hide, close, resize, move windows
- ✅ **State Management**: Focus, visibility, maximization, minimization
- ✅ **Attribute Control**: Title, size, position, and window states
- ✅ **Event Subscription**: Subscribe to keyboard, mouse, and window events

### Input Handling
- ✅ **Keyboard Input**: Full keyboard event support with modifier keys
- ✅ **Mouse Input**: Mouse button events with position tracking
- ✅ **Scroll Events**: Mouse wheel and touchpad scrolling
- ✅ **Pointer Tracking**: Enter/leave and motion events

## Architecture

### Driver Structure
```
wayland/
├── driver.go          # Main Wayland driver implementation
├── window.go          # Wayland window implementation
├── example/           # Example applications
│   └── main.go       # Basic Wayland backend test
└── README.md         # This file
```

### Core Components

#### WaylandDriver (`driver.go`)
- Manages Wayland display connection and global objects
- Handles event loop and input device management
- Provides window creation and event translation
- Manages keyboard mapping and input method context

#### WaylandWindow (`window.go`)
- Implements the `gio.Window` interface for Wayland
- Manages individual window state and properties
- Handles window-specific events and operations
- Provides surface management and rendering context

## Usage

### Basic Example

```go
package main

import (
    "github.com/opensraph/sraph/gio"
    _ "github.com/opensraph/sraph/gio/wayland" // Import to register driver
)

func main() {
    // Create window with Wayland driver
    window, err := gio.CreateWindowWithDriver(gio.DriverTypeWayland, gio.NewWindowOptions{
        Title:  "My Wayland App",
        Width:  800,
        Height: 600,
        State:  gio.WindowStateVisible | gio.WindowStateFocused,
    })
    if err != nil {
        panic(err)
    }
    defer window.Close()

    // Subscribe to events
    window.Subscribe(gio.EventTypeKeyboard, func(e gio.Event) error {
        if keyEvent, ok := e.(*gio.KeyboardEvent); ok {
            fmt.Printf("Key: %v\n", keyEvent.Code)
        }
        return nil
    })

    // Keep window open
    select {}
}
```

### Automatic Driver Selection

The Wayland driver is automatically selected as the best driver on Linux systems when `WAYLAND_DISPLAY` environment variable is set:

```go
// Automatically chooses best driver (Wayland on Wayland sessions)
window, err := gio.CreateWindow(gio.NewWindowOptions{
    Title:  "Auto-Selected Driver",
    Width:  800,
    Height: 600,
})
```

## Implementation Details

### Mock Implementation

The current implementation is a **mock/stub implementation** that demonstrates the complete API structure without requiring actual Wayland libraries. This approach:

- ✅ **Provides Complete API**: All interfaces and methods are implemented
- ✅ **Demonstrates Architecture**: Shows how a real Wayland backend would be structured
- ✅ **Enables Testing**: Allows testing of the driver system without Wayland dependencies
- ✅ **Platform Agnostic**: Works on any platform for development and testing

### Real Wayland Implementation

To convert this to a full Wayland implementation, you would need to:

1. **Add CGO Bindings**: Include actual Wayland client library bindings
2. **Implement Protocols**: Add support for core Wayland protocols (wl_compositor, wl_shell, etc.)
3. **Buffer Management**: Implement shared memory buffers and rendering
4. **XKB Integration**: Add real XKB keyboard handling
5. **Protocol Extensions**: Support for additional Wayland protocols as needed

### Dependencies for Real Implementation

A full implementation would require:
```
wayland-client >= 1.18
wayland-protocols >= 1.24
xkbcommon >= 1.0
```

## Event System

### Supported Events

- **Keyboard Events**: Key press/release with modifier support
- **Mouse Events**: Button press/release with position tracking
- **Wheel Events**: Scroll events with delta values
- **Window Events**: Focus, resize, move, close events
- **Pointer Events**: Enter/leave and motion tracking

### Event Flow

1. **Native Events**: Wayland events received from compositor
2. **Translation**: Convert to Gio event format
3. **Dispatching**: Route events to appropriate window
4. **Subscription**: Deliver to registered event handlers

## Testing

### Running Tests

```bash
cd wayland/example
go build -o wayland-test
./wayland-test
```

### Expected Output

```
=== Wayland Driver Test ===
Platform Info: Platform: linux/amd64, Available drivers: [Wayland X11], Best driver: Wayland

Available drivers:
  - Wayland (available: true)
  - X11 (available: true)

Creating window with Wayland driver...
Window created successfully with ID: 100
Window size: 800x600
Window title: Wayland Test Window
Window visible: true
Window focused: true

Testing window operations...
Resizing window to 1024x768...
Changing window title...
Maximizing window...
Unmaximizing window...
Window operations completed successfully!
```

## Status

### Current Status: ✅ **Mock Implementation Complete**

The Wayland backend provides a complete mock implementation that demonstrates:
- Full API compatibility with the Gio graphics system
- Proper driver registration and selection
- Complete event handling architecture
- Window management and state control
- Multi-window support

### Next Steps for Production

1. **CGO Integration**: Add actual Wayland client library bindings
2. **Protocol Implementation**: Implement core Wayland protocols
3. **Rendering Integration**: Add buffer management and rendering support
4. **Testing**: Comprehensive testing on real Wayland compositors
5. **Documentation**: Complete API documentation and examples

## Contributing

When contributing to the Wayland backend:

1. **Maintain API Compatibility**: Ensure all changes maintain compatibility with the `gio.Driver` and `gio.Window` interfaces
2. **Test Thoroughly**: Test both mock and any real implementations
3. **Follow Conventions**: Use the same patterns as the X11 backend for consistency
4. **Document Changes**: Update this README and add inline documentation

## License

This Wayland backend is part of the Sraph graphics system and follows the same license as the parent project.
