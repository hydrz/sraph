package gio

import (
	"sync"
	"time"
)

// Config represents the global configuration for the gio package
type Config struct {
	// Event system configuration
	EventQueueSize    int           // Size of event queue buffer
	EventTimeout      time.Duration // Timeout for event processing
	MaxEventHandlers  int           // Maximum number of event handlers per type
	AsyncEventHandler bool          // Whether to handle events asynchronously

	// Window configuration
	DefaultWindowWidth  int // Default window width
	DefaultWindowHeight int // Default window height

	// Key repeat settings
	KeyRepeatDelay time.Duration // Initial delay before repeat starts
	KeyRepeatRate  time.Duration // Interval between repeats

	mu sync.RWMutex
}

var (
	globalConfig = &Config{
		EventQueueSize:      1000,
		EventTimeout:        time.Second * 5,
		MaxEventHandlers:    100,
		AsyncEventHandler:   true,
		DefaultWindowWidth:  1024,
		DefaultWindowHeight: 768,
		// Default key repeat settings (similar to desktop environments)
		KeyRepeatDelay: 500 * time.Millisecond, // 500ms initial delay
		KeyRepeatRate:  33 * time.Millisecond,  // ~30 Hz repeat rate
	}
)

// GetConfig returns the global configuration
func GetConfig() *Config {
	return globalConfig
}

// SetConfig updates the global configuration
func SetConfig(config *Config) {
	config.mu.Lock()
	defer config.mu.Unlock()
	globalConfig = config
}

// Clone creates a copy of the configuration
func (c *Config) Clone() *Config {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return &Config{
		EventQueueSize:      c.EventQueueSize,
		EventTimeout:        c.EventTimeout,
		MaxEventHandlers:    c.MaxEventHandlers,
		AsyncEventHandler:   c.AsyncEventHandler,
		DefaultWindowWidth:  c.DefaultWindowWidth,
		DefaultWindowHeight: c.DefaultWindowHeight,
		KeyRepeatDelay:      c.KeyRepeatDelay,
		KeyRepeatRate:       c.KeyRepeatRate,
	}
}
