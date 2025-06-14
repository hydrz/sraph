package gio

import (
	"github.com/opensraph/sraph/geom"
)

// WindowID represents a unique window identifier
type WindowID uint64

// WindowState represents the state of a window
type WindowState int

const (
	WindowStateNormal WindowState = iota
	WindowStateMinimized
	WindowStateMaximized
	WindowStateFullscreen
)

func (s WindowState) String() string {
	switch s {
	case WindowStateNormal:
		return "normal"
	case WindowStateMinimized:
		return "minimized"
	case WindowStateMaximized:
		return "maximized"
	case WindowStateFullscreen:
		return "fullscreen"
	default:
		return "unknown"
	}
}

// WindowAttribute represents window attributes
type WindowAttribute int

const (
	WindowAttributeResizable WindowAttribute = iota
	WindowAttributeVisible
	WindowAttributeDecorated
	WindowAttributeFocused
	WindowAttributeAutoIconify
	WindowAttributeFloating
	WindowAttributeTransparent
)

func (a WindowAttribute) String() string {
	switch a {
	case WindowAttributeResizable:
		return "resizable"
	case WindowAttributeVisible:
		return "visible"
	case WindowAttributeDecorated:
		return "decorated"
	case WindowAttributeFocused:
		return "focused"
	case WindowAttributeAutoIconify:
		return "auto_iconify"
	case WindowAttributeFloating:
		return "floating"
	case WindowAttributeTransparent:
		return "transparent"
	default:
		return "unknown"
	}
}

// InputMode represents input modes
type InputMode int

const (
	InputModeCursor InputMode = iota
	InputModeStickyKeys
	InputModeStickyMouseButtons
	InputModeRawMouseMotion
)

func (m InputMode) String() string {
	switch m {
	case InputModeCursor:
		return "cursor"
	case InputModeStickyKeys:
		return "sticky_keys"
	case InputModeStickyMouseButtons:
		return "sticky_mouse_buttons"
	case InputModeRawMouseMotion:
		return "raw_mouse_motion"
	default:
		return "unknown"
	}
}

type Window interface {
	// Release closes the window.
	//
	// The behavior of the Window after Release, whether calling its methods or
	// passing it as an argument, is undefined.
	Release()

	// Send a window event to the window.
	Send(event Event)

	// Receive a window event.
	Receive() Event
}

// NewWindowOptions contains window creation hints
type NewWindowOptions struct {
	// Basic window properties
	Title    string
	Size     geom.Size[geom.F32]   // Default size
	Position *geom.Point[geom.F32] // nil means let the system decide

	// Window behavior
	Resizable    bool
	Decorated    bool
	Visible      bool
	Focused      bool
	AutoIconify  bool
	Floating     bool
	Maximized    bool
	Transparent  bool
	CenterCursor bool
	FocusOnShow  bool

	// EventBus
	event EventBus // Optional event bus for window events
}

// DefaultNewWindowOptions returns default window creation hints
func DefaultNewWindowOptions() NewWindowOptions {
	return NewWindowOptions{
		Title:        "Sraph Window",
		Size:         geom.NewSize[geom.F32](800, 600),
		Resizable:    true,
		Decorated:    true,
		Visible:      true,
		Focused:      true,
		AutoIconify:  true,
		CenterCursor: true,
		FocusOnShow:  true,
	}
}
