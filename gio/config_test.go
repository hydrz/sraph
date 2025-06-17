package gio

import (
	"testing"
	"time"
)

func TestGetConfig(t *testing.T) {
	config := GetConfig()
	if config == nil {
		t.Fatal("GetConfig() returned nil")
	}

	// Test default values
	if config.EventQueueSize != 1000 {
		t.Errorf("Expected EventQueueSize 1000, got %d", config.EventQueueSize)
	}

	if config.DefaultWindowWidth != 1024 {
		t.Errorf("Expected DefaultWindowWidth 1024, got %d", config.DefaultWindowWidth)
	}

	if config.KeyRepeatDelay != 500*time.Millisecond {
		t.Errorf("Expected KeyRepeatDelay 500ms, got %v", config.KeyRepeatDelay)
	}

	if config.KeyRepeatRate != 33*time.Millisecond {
		t.Errorf("Expected KeyRepeatRate 33ms, got %v", config.KeyRepeatRate)
	}
}

func TestSetConfig(t *testing.T) {
	// Save original config
	originalConfig := GetConfig().Clone()

	// Create new config
	newConfig := &Config{
		EventQueueSize:      2000,
		EventTimeout:        10 * time.Second,
		MaxEventHandlers:    200,
		AsyncEventHandler:   false,
		DefaultWindowWidth:  1920,
		DefaultWindowHeight: 1080,
		KeyRepeatDelay:      400 * time.Millisecond,
		KeyRepeatRate:       25 * time.Millisecond,
	}

	// Set new config
	SetConfig(newConfig)

	// Verify config was set
	config := GetConfig()
	if config.EventQueueSize != 2000 {
		t.Errorf("Expected EventQueueSize 2000, got %d", config.EventQueueSize)
	}

	if config.DefaultWindowWidth != 1920 {
		t.Errorf("Expected DefaultWindowWidth 1920, got %d", config.DefaultWindowWidth)
	}

	if config.KeyRepeatDelay != 400*time.Millisecond {
		t.Errorf("Expected KeyRepeatDelay 400ms, got %v", config.KeyRepeatDelay)
	}

	// Restore original config
	SetConfig(originalConfig)
}

func TestConfigClone(t *testing.T) {
	original := GetConfig()
	cloned := original.Clone()

	// Test that values are equal
	if cloned.EventQueueSize != original.EventQueueSize {
		t.Error("Clone should have same EventQueueSize")
	}

	if cloned.KeyRepeatDelay != original.KeyRepeatDelay {
		t.Error("Clone should have same KeyRepeatDelay")
	}

	// Test that it's a separate instance
	cloned.EventQueueSize = 9999
	if original.EventQueueSize == 9999 {
		t.Error("Modifying clone should not affect original")
	}
}
