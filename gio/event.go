package gio

import (
	"time"

	"github.com/opensraph/sraph/geom"
)

type EventType int

const (
	EventTypeUnknown EventType = iota
	// Window events
	EventTypeWindowClose
	EventTypeWindowResize
	EventTypeWindowMove
	EventTypeWindowFocus
	EventTypeWindowIconify
	EventTypeWindowMaximize
	EventTypeWindowRefresh
	EventTypeWindowFramebufferResize
	EventTypeWindowContentScale

	// Input events
	EventTypeKey
	EventTypeChar
	EventTypeMouseButton
	EventTypeMouseMove
	EventTypeMouseEnter
	EventTypeScroll

	// Monitor events
	EventTypeMonitor

	// Joystick events
	EventTypeJoystick

	// Drop events
	EventTypeDrop
)

func (e EventType) String() string {
	switch e {
	case EventTypeWindowClose:
		return "window_close"
	case EventTypeWindowResize:
		return "window_resize"
	case EventTypeWindowMove:
		return "window_move"
	case EventTypeWindowFocus:
		return "window_focus"
	case EventTypeWindowIconify:
		return "window_iconify"
	case EventTypeWindowMaximize:
		return "window_maximize"
	case EventTypeWindowRefresh:
		return "window_refresh"
	case EventTypeWindowFramebufferResize:
		return "window_framebuffer_resize"
	case EventTypeWindowContentScale:
		return "window_content_scale"
	case EventTypeKey:
		return "key_event"
	case EventTypeChar:
		return "char_event"
	case EventTypeMouseButton:
		return "mouse_button_event"
	case EventTypeMouseMove:
		return "mouse_move_event"
	case EventTypeMouseEnter:
		return "mouse_enter_event"
	case EventTypeScroll:
		return "scroll_event"
	case EventTypeJoystick:
		return "joystick_event"
	case EventTypeDrop:
		return "drop_event"
	default:
		return "unknown_event_type"
	}
}

type Event interface {
	Type() EventType
	Time() time.Time
	WindowID() uint64
}

// BaseEvent provides common event functionality
type BaseEvent struct {
	EventType EventType
	Time      time.Time
	WinID     uint64
}

func (e BaseEvent) Type() EventType {
	return e.EventType
}

func (e BaseEvent) Timestamp() time.Time {
	return e.Time
}

func (e BaseEvent) WindowID() uint64 {
	return e.WinID
}

// Window Events
type WindowCloseEvent struct {
	BaseEvent
}

type WindowResizeEvent struct {
	BaseEvent
	Size geom.Size[geom.F32]
}

type WindowMoveEvent struct {
	BaseEvent
	Position geom.Point[geom.F32]
}

type WindowFocusEvent struct {
	BaseEvent
	Focused bool
}

type WindowIconifyEvent struct {
	BaseEvent
	Iconified bool
}

type WindowMaximizeEvent struct {
	BaseEvent
	Maximized bool
}

type WindowRefreshEvent struct {
	BaseEvent
}

type WindowFramebufferResizeEvent struct {
	BaseEvent
	Size geom.Size[geom.F32]
}

type WindowContentScaleEvent struct {
	BaseEvent
	Scale geom.Point[geom.F32]
}

// Input Events
type KeyEvent struct {
	BaseEvent
	Key      KeyCode
	Scancode int
	Action   KeyAction
	Mods     ModifierKey
}

type CharEvent struct {
	BaseEvent
	Codepoint uint32
}

type MouseButtonEvent struct {
	BaseEvent
	Button MouseButton
	Action KeyAction
	Mods   ModifierKey
}

type MouseMoveEvent struct {
	BaseEvent
	Position geom.Point[geom.F32]
}

type MouseEnterEvent struct {
	BaseEvent
	Entered bool
}

type ScrollEvent struct {
	BaseEvent
	Offset geom.Point[geom.F32]
}

// Joystick Events
type JoystickEvent struct {
	BaseEvent
	Joystick int
	Action   JoystickAction
}

// Drop Events
type DropEvent struct {
	BaseEvent
	Paths []string
}
