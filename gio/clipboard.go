package gio

import "sync"

// ClipboardFormat represents the format of clipboard data
type ClipboardFormat int

const (
	ClipboardFormatText ClipboardFormat = iota
	ClipboardFormatImage
	ClipboardFormatHTML
	ClipboardFormatRTF
	ClipboardFormatFiles
)

func (f ClipboardFormat) String() string {
	switch f {
	case ClipboardFormatText:
		return "text/plain"
	case ClipboardFormatImage:
		return "image/png"
	case ClipboardFormatHTML:
		return "text/html"
	case ClipboardFormatRTF:
		return "text/rtf"
	case ClipboardFormatFiles:
		return "application/x-file-list"
	default:
		return "unknown"
	}
}

// ClipboardData represents clipboard data
type ClipboardData struct {
	Format ClipboardFormat
	Data   []byte
	Text   string   // For text format
	Files  []string // For files format
}

// EventHandler defines the function signature for event handlers.
type EventHandler func(Event)

// EventDispatcher manages event handlers and dispatches events.
type EventDispatcher struct {
	mu       sync.RWMutex
	handlers map[EventType][]EventHandler
}

// NewEventDispatcher creates a new EventDispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Register adds an event handler for a specific event type.
func (d *EventDispatcher) Register(eventType EventType, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Unregister removes an event handler for a specific event type.
func (d *EventDispatcher) Unregister(eventType EventType, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	handlers := d.handlers[eventType]
	for i, h := range handlers {
		if &h == &handler {
			d.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// Dispatch sends an event to all registered handlers for its type.
func (d *EventDispatcher) Dispatch(event Event) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, handler := range d.handlers[event.Type()] {
		// Call handlers in a new goroutine to avoid blocking.
		go handler(event)
	}
}
