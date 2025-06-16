package gio

import (
	"fmt"
	"image"
	"sync"

	"github.com/opensraph/sraph/gpu"
)

// WindowState represents window attributes
type WindowState uint16

const (
	WindowStateUnknown                WindowState = iota      // Unknown attribute
	WindowStateFocused                WindowState = 1 << iota // Window is focused
	WindowStateIconified                                      // Window is minimized
	WindowStateMaximized                                      // Window is maximized
	WindowStateVisible                                        // Window is visible
	WindowStateHovered                                        // Cursor is over the window
	WindowStateResizable                                      // Window is resizable
	WindowStateDecorated                                      // Window has decorations (title bar, borders, etc.)
	WindowStateFloating                                       // Window is always on top
	WindowStateAutoIconify                                    // Fullscreen windows auto-iconify on focus loss
	WindowStateCenterCursor                                   // Cursor is centered over fullscreen windows
	WindowStateTransparentFramebuffer                         // Framebuffer is transparent
	WindowStateFocusOnShow                                    // Window gets focus when shown
	WindowStateScaleToMonitor                                 // Window content area scales with monitor content scale
)

func (ws WindowState) Contains(state WindowState) bool {
	return ws&state != 0
}

func (ws WindowState) Diff(state WindowState) WindowState {
	return ws &^ state
}

// WindowAttr represents window creation attributes
type WindowAttr struct {
	Title         string
	Width, Height int
	Position      image.Point // Position of the window on the screen
	State         WindowState // Initial state of the window
}

func (wa WindowAttr) Equal(other WindowAttr) bool {
	return wa.Title == other.Title &&
		wa.Width == other.Width &&
		wa.Height == other.Height &&
		wa.Position.Eq(other.Position) &&
		wa.State.Diff(other.State) == 0
}

func (wa WindowAttr) Apply(o ...WindowAttr) WindowAttr {
	if len(o) == 0 {
		return wa
	}

	for _, attr := range o {
		if attr.Title != "" {
			wa.Title = attr.Title
		}
		if attr.Width > 0 {
			wa.Width = attr.Width
		}
		if attr.Height > 0 {
			wa.Height = attr.Height
		}
		if !attr.Position.Eq(image.Point{}) {
			wa.Position = attr.Position
		}
		if attr.State != WindowStateUnknown {
			wa.State = wa.State | attr.State
		}
	}

	return wa
}

type PlatformWindow interface {
	// SetAttr sets the window attributes.
	SetAttr(attr WindowAttr) error
	// Event returns the event bus for the window.
	Event() *EventBus
	// Surface returns the GPU surface associated with the window.
	Surface() (gpu.Surface, error)

	Close() error
}

type Window interface {
	Width() int
	Height() int
	Title() string
	Position() image.Point
	State() WindowState

	Show() error
	Hide() error
	Close() error

	// Surface returns the GPU surface associated with the window.
	Surface() (gpu.Surface, error)

	// Event returns the event bus for the window.
	Event() *EventBus
}

type NewWindowOptions = WindowAttr

// DefaultNewWindowOptions returns the default window attributes.
func DefaultNewWindowOptions() NewWindowOptions {
	return NewWindowOptions{
		Title:    "Window",
		Width:    800,
		Height:   600,
		Position: image.Point{X: 100, Y: 100},
		State:    WindowStateFocused | WindowStateVisible | WindowStateResizable | WindowStateDecorated,
	}
}

// NewWindow creates a new window with the given attributes.
func NewWindow(pw PlatformWindow, options ...NewWindowOptions) (Window, error) {
	o := DefaultNewWindowOptions()
	o = o.Apply(options...)

	if err := pw.SetAttr(o); err != nil {
		if closeErr := pw.Close(); closeErr != nil {
			return nil, closeErr
		}
		return nil, err
	}

	surface, err := pw.Surface()
	if err != nil {
		if closeErr := pw.Close(); closeErr != nil {
			return nil, closeErr
		}
		return nil, err
	}

	return &window{
		attr:    o,
		pw:      pw,
		event:   pw.Event(),
		surface: surface,
	}, nil
}

type window struct {
	attr    WindowAttr
	pw      PlatformWindow
	event   *EventBus
	surface gpu.Surface
	mu      sync.Mutex // Mutex to protect against concurrent access
	closed  bool       // Flag to indicate if the window is closed
}

// Close implements Window.
func (w *window) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return nil // Already closed
	}
	w.closed = true

	if err := w.pw.Close(); err != nil {
		return fmt.Errorf("failed to close window: %w", err)
	}

	return nil
}

// Event implements Window.
func (w *window) Event() *EventBus {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return nil // Return nil if the window is closed
	}

	return w.event
}

// Height implements Window.
func (w *window) Height() int {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return 0 // Return 0 if the window is closed
	}

	return w.attr.Height
}

// Hide implements Window.
func (w *window) Hide() error {
	panic("unimplemented")
}

// Position implements Window.
func (w *window) Position() image.Point {
	panic("unimplemented")
}

// Show implements Window.
func (w *window) Show() error {
	panic("unimplemented")
}

// State implements Window.
func (w *window) State() WindowState {
	panic("unimplemented")
}

// Surface implements Window.
func (w *window) Surface() (gpu.Surface, error) {
	panic("unimplemented")
}

// Title implements Window.
func (w *window) Title() string {
	panic("unimplemented")
}

// Width implements Window.
func (w *window) Width() int {
	panic("unimplemented")
}
