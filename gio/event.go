package gio

import (
	"context"
	"image"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/text/language"
)

type EventType uint8

const (
	EventTypeUnknown   EventType = iota
	EventTypeWindow              // Window events (e.g., resize, move, close)
	EventTypeKeyboard            // Keyboard events
	EventTypeMouse               // Mouse events
	EventTypeWheel               // Mouse wheel events
	EventTypeClipboard           // Clipboard events
	EventTypeDrag                // Drag and drop events
)

type Event interface {
	Type() EventType
	Time() time.Time
	ID() uint64 // Unique event ID
}

type baseEvent struct {
	eventType EventType
	time      time.Time
	id        uint64
}

var eventIDCounter uint64

func NewBaseEvent(eventType EventType) Event {
	return &baseEvent{
		eventType: eventType,
		time:      time.Now(),
		id:        atomic.AddUint64(&eventIDCounter, 1),
	}
}

func (e *baseEvent) Type() EventType { return e.eventType }
func (e *baseEvent) Time() time.Time { return e.time }
func (e *baseEvent) ID() uint64      { return e.id }

type WindowEvent struct {
	Event
	Window BaseWindow
}

func NewWindowEvent() *WindowEvent {
	return &WindowEvent{
		Event: NewBaseEvent(EventTypeWindow),
	}
}

type ClipboardEvent struct {
	Event
	Data *DataTransfer // Data associated with the clipboard event
}

func NewClipboardEvent() *ClipboardEvent {
	return &ClipboardEvent{
		Event: NewBaseEvent(EventTypeClipboard),
		Data:  NewDataTransfer(),
	}
}

type DragEvent struct {
	Event
	Data *DataTransfer // Data associated with the drag event
}

func NewDragEvent() *DragEvent {
	return &DragEvent{
		Event: NewBaseEvent(EventTypeDrag),
		Data:  NewDataTransfer(),
	}
}

// KeyboardEvent represents a keyboard event following W3C standard
type KeyboardEvent struct {
	Event
	Location    KeyLocation  // The location of the key on the keyboard
	ModifierKey ModifierKey  // Bitmask of modifier keys pressed
	Repeat      bool         // Whether the key is being held down
	Locale      language.Tag // The locale identifier
	Code        KeyCode      // The code value of the key pressed
}

// NewKeyboardEvent creates a new keyboard event
func NewKeyboardEvent() *KeyboardEvent {
	return &KeyboardEvent{
		Event: NewBaseEvent(EventTypeKeyboard),
	}
}

type TouchEvent struct {
	Event
	ChangedTouches TouchList   // List of touches that changed since the last event
	ModifierKey    ModifierKey // Modifier keys pressed during the event
	TargetTouches  TouchList   // List of touches currently on the target
	Touches        TouchList   // List of all touches on the screen
}

// NewTouchEvent creates a new touch event
func NewTouchEvent() *TouchEvent {
	return &TouchEvent{
		Event: NewBaseEvent(EventTypeUnknown), // Use a more specific type if needed
	}
}

type MouseEvent struct {
	Event
	Button      MouseButton // The mouse button that was pressed
	ModifierKey ModifierKey // Modifier keys pressed during the event
	Position    image.Point // Position of the mouse event relative to the viewport
}

func NewMouseEvent() *MouseEvent {
	return &MouseEvent{
		Event: NewBaseEvent(EventTypeMouse),
	}
}

type PointerEvent struct {
	*MouseEvent
	PointerID          int         // Unique identifier for the pointer
	Pressure           float64     // Pressure applied by the pointer (0.0 to 1.0)
	TangentialPressure float64     // Pressure applied tangentially (0.0 to 1.0)
	TiltX              float64     // Tilt angle in the X direction (degrees)
	TiltY              float64     // Tilt angle in the Y direction (degrees)
	Twist              float64     // Twist angle (degrees)
	Width              float64     // Width of the pointer contact area (CSS pixels)
	Height             float64     // Height of the pointer contact area (CSS pixels)
	PointerType        PointerType // Type of pointer (mouse, pen, touch)
	IsPrimary          bool        // Whether this is the primary pointer for the device
}

func NewPointerEvent() *PointerEvent {
	return &PointerEvent{
		MouseEvent: NewMouseEvent(),
	}
}

type WheelEvent struct {
	Event
	DeltaX      float64     // Horizontal scroll delta
	DeltaY      float64     // Vertical scroll delta
	ModifierKey ModifierKey // Modifier keys pressed during the event
	Position    image.Point // Position of the wheel event relative to the viewport
}

func NewWheelEvent() *WheelEvent {
	return &WheelEvent{
		Event: NewBaseEvent(EventTypeWheel),
	}
}

type EventHandler func(ctx Context, e Event) error

// EventBus provides improved event handling with context support
type EventBus struct {
	subscribers map[EventType][]EventHandler
	eventQueue  chan Event
	config      *Config
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

func NewEventBus() *EventBus {
	config := GetConfig()
	ctx, cancel := context.WithCancel(context.Background())

	eb := &EventBus{
		subscribers: make(map[EventType][]EventHandler),
		eventQueue:  make(chan Event, config.EventQueueSize),
		config:      config.Clone(),
		ctx:         ctx,
		cancel:      cancel,
	}

	// Start event processing goroutine
	eb.wg.Add(1)
	go eb.processEvents()

	return eb
}

func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	handlers := eb.subscribers[eventType]
	if len(handlers) >= eb.config.MaxEventHandlers {
		return ErrEventHandlerFull
	}

	eb.subscribers[eventType] = append(handlers, handler)
	return nil
}

func (eb *EventBus) Unsubscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	handlers, exists := eb.subscribers[eventType]
	if !exists {
		return
	}

	// Remove handler (simplified approach)
	for i := len(handlers) - 1; i >= 0; i-- {
		eb.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
		break
	}

	if len(eb.subscribers[eventType]) == 0 {
		delete(eb.subscribers, eventType)
	}
}

func (eb *EventBus) Publish(event Event) error {
	select {
	case eb.eventQueue <- event:
		return nil
	case <-time.After(eb.config.EventTimeout):
		return ErrOperationTimeout
	case <-eb.ctx.Done():
		return eb.ctx.Err()
	}
}

func (eb *EventBus) Close() {
	eb.cancel()
	eb.wg.Wait()
	close(eb.eventQueue)
}

func (eb *EventBus) processEvents() {
	defer eb.wg.Done()

	for {
		select {
		case event, ok := <-eb.eventQueue:
			if !ok {
				return
			}
			eb.handleEvent(event)
		case <-eb.ctx.Done():
			return
		}
	}
}

func (eb *EventBus) handleEvent(event Event) {
	eb.mu.RLock()
	handlers, exists := eb.subscribers[event.Type()]
	if !exists {
		eb.mu.RUnlock()
		return
	}

	// Copy handlers to avoid holding lock during execution
	handlersCopy := make([]EventHandler, len(handlers))
	copy(handlersCopy, handlers)
	eb.mu.RUnlock()

	ctx := NewContext(eb.ctx)

	if eb.config.AsyncEventHandler {
		// Execute handlers asynchronously
		for _, handler := range handlersCopy {
			go func(h EventHandler) {
				if err := h(ctx, event); err != nil {
					// Log error or handle it appropriately
					_ = err
				}
			}(handler)
		}
	} else {
		// Execute handlers synchronously
		for _, handler := range handlersCopy {
			if err := handler(ctx, event); err != nil {
				// Log error or handle it appropriately
				_ = err
			}
		}
	}
}
