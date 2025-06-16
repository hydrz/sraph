package gio

import (
	"context"
	"sync"
	"time"
)

// Context provides cancellation and timeout support for gio operations
type Context interface {
	context.Context

	// WindowID returns the associated window ID if any
	WindowID() WindowID

	// WithWindow creates a new context with window association
	WithWindow(windowID WindowID) Context
}

type gioContext struct {
	context.Context
	windowID WindowID
	mu       sync.RWMutex
}

// NewContext creates a new gio context
func NewContext(parent context.Context) Context {
	if parent == nil {
		parent = context.Background()
	}

	return &gioContext{
		Context: parent,
	}
}

// NewContextWithTimeout creates a new gio context with timeout
func NewContextWithTimeout(parent context.Context, timeout time.Duration) (Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}

	ctx, cancel := context.WithTimeout(parent, timeout)
	return &gioContext{Context: ctx}, cancel
}

// NewContextWithCancel creates a new gio context with cancellation
func NewContextWithCancel(parent context.Context) (Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}

	ctx, cancel := context.WithCancel(parent)
	return &gioContext{Context: ctx}, cancel
}

func (c *gioContext) WindowID() WindowID {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.windowID
}

func (c *gioContext) WithWindow(windowID WindowID) Context {
	c.mu.Lock()
	defer c.mu.Unlock()

	return &gioContext{
		Context:  c.Context,
		windowID: windowID,
	}
}
