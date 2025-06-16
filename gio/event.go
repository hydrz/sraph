package gio

import (
	"sync"
	"time"
)

type EventType uint8

const (
	EventTypeUnknown  EventType = iota
	EventTypeKeyboard           // Keyboard events
	EventTypeMouse              // Mouse events
	EventTypeWheel              // Mouse wheel events
	EventTypeWindow             // Window events (resize, close, etc.)
	EventTypeExpose             // Window redraw events
)

type Event interface {
	Type() EventType
	Time() time.Time
}

type BaseEvent struct {
	eventType EventType
	time      time.Time
}

func (e *BaseEvent) Type() EventType { return e.eventType }
func (e *BaseEvent) Time() time.Time { return e.time }

type EventHandler func(e Event)

type EventBus struct {
	subscribers map[EventType][]*EventHandler
	mu          sync.Mutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]*EventHandler),
	}
}
func (eb *EventBus) Subscribe(eventType EventType, handler *EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if _, exists := eb.subscribers[eventType]; !exists {
		eb.subscribers[eventType] = []*EventHandler{}
	}
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}
func (eb *EventBus) Unsubscribe(eventType EventType, handler *EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if handlers, exists := eb.subscribers[eventType]; exists {
		for i, h := range handlers {
			if h == handler {
				eb.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
		if len(eb.subscribers[eventType]) == 0 {
			delete(eb.subscribers, eventType)
		}
	}
}

func (eb *EventBus) Send(event Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if handlers, exists := eb.subscribers[event.Type()]; exists {
		for _, handler := range handlers {
			go (*handler)(event)
		}
	}
}
