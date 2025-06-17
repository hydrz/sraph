package gio

import (
	"fmt"
	"image"
	"sync"
	"testing"
	"time"
)

func TestDefaultNewWindowOptions(t *testing.T) {
	// Save original config
	originalConfig := GetConfig()
	defer SetConfig(originalConfig)

	// Set test config
	testConfig := &Config{
		DefaultWindowWidth:  800,
		DefaultWindowHeight: 600,
		KeyRepeatDelay:      100 * time.Millisecond,
		KeyRepeatRate:       50 * time.Millisecond,
	}
	SetConfig(testConfig)

	opts := DefaultNewWindowOptions()

	if opts.Title != "Window" {
		t.Errorf("Expected title 'Window', got '%s'", opts.Title)
	}

	if opts.Width != 800 {
		t.Errorf("Expected width 800, got %d", opts.Width)
	}

	if opts.Height != 600 {
		t.Errorf("Expected height 600, got %d", opts.Height)
	}

	expectedState := WindowStateFocused | WindowStateVisible | WindowStateResizable | WindowStateDecorated
	if opts.State != expectedState {
		t.Errorf("Expected state %d, got %d", expectedState, opts.State)
	}
}

func TestWindowAttrEqual(t *testing.T) {
	attr1 := WindowAttr{
		Title:    "Test Window",
		Width:    800,
		Height:   600,
		Position: image.Point{X: 100, Y: 200},
		State:    WindowStateVisible,
	}

	attr2 := WindowAttr{
		Title:    "Test Window",
		Width:    800,
		Height:   600,
		Position: image.Point{X: 100, Y: 200},
		State:    WindowStateVisible,
	}

	if !attr1.Equal(attr2) {
		t.Error("Expected equal attributes to be equal")
	}

	attr2.Title = "Different Title"
	if attr1.Equal(attr2) {
		t.Error("Expected different attributes not to be equal")
	}
}

func TestWindowAttrApply(t *testing.T) {
	base := WindowAttr{
		Title:    "Base Window",
		Width:    800,
		Height:   600,
		Position: image.Point{X: 100, Y: 200},
		State:    WindowStateVisible,
	}

	override := WindowAttr{
		Title:  "Override Title",
		Width:  1024,
		Height: 768,
	}

	result := base.Apply(override)

	if result.Title != "Override Title" {
		t.Errorf("Expected title 'Override Title', got '%s'", result.Title)
	}

	if result.Width != 1024 {
		t.Errorf("Expected width 1024, got %d", result.Width)
	}

	if result.Height != 768 {
		t.Errorf("Expected height 768, got %d", result.Height)
	}

	// Should keep original position and state
	if result.Position != base.Position {
		t.Errorf("Expected position %v, got %v", base.Position, result.Position)
	}

	if result.State != base.State {
		t.Errorf("Expected state %d, got %d", base.State, result.State)
	}
}

func TestWindowStateContains(t *testing.T) {
	state := WindowStateVisible | WindowStateResizable | WindowStateDecorated

	if !state.Contains(WindowStateVisible) {
		t.Error("Expected state to contain WindowStateVisible")
	}

	if !state.Contains(WindowStateResizable) {
		t.Error("Expected state to contain WindowStateResizable")
	}

	if state.Contains(WindowStateClosed) {
		t.Error("Expected state not to contain WindowStateClosed")
	}
}

func TestWindowStateDiff(t *testing.T) {
	state := WindowStateVisible | WindowStateResizable | WindowStateDecorated

	result := state.Diff(WindowStateVisible)
	expected := WindowStateResizable | WindowStateDecorated

	if result != expected {
		t.Errorf("Expected diff result %d, got %d", expected, result)
	}
}

func TestNewBaseWindow(t *testing.T) {
	opts := NewWindowOptions{
		Title:    "Test Window",
		Width:    1024,
		Height:   768,
		Position: image.Point{X: 50, Y: 100},
		State:    WindowStateVisible | WindowStateResizable,
	}

	window := NewBaseWindow(opts)

	if window.Title() != "Test Window" {
		t.Errorf("Expected title 'Test Window', got '%s'", window.Title())
	}

	if window.Width() != 1024 {
		t.Errorf("Expected width 1024, got %d", window.Width())
	}

	if window.Height() != 768 {
		t.Errorf("Expected height 768, got %d", window.Height())
	}

	expectedPos := image.Point{X: 50, Y: 100}
	if window.Position() != expectedPos {
		t.Errorf("Expected position %v, got %v", expectedPos, window.Position())
	}
}

func TestBaseWindowStateQueries(t *testing.T) {
	opts := NewWindowOptions{
		State: WindowStateVisible | WindowStateResizable | WindowStateDecorated,
	}

	window := NewBaseWindow(opts)

	if !window.IsVisible() {
		t.Error("Expected window to be visible")
	}

	if !window.IsResizable() {
		t.Error("Expected window to be resizable")
	}

	if !window.IsDecorated() {
		t.Error("Expected window to be decorated")
	}

	if window.IsClosed() {
		t.Error("Expected window not to be closed")
	}

	if window.IsMaximized() {
		t.Error("Expected window not to be maximized")
	}
}

func TestBaseWindowSetAttr(t *testing.T) {
	window := NewBaseWindow(NewWindowOptions{
		Title:  "Original Title",
		Width:  800,
		Height: 600,
	})

	newAttr := WindowAttr{
		Title:  "New Title",
		Width:  1024,
		Height: 768,
		State:  WindowStateMaximized,
	}

	window.SetAttr(newAttr)

	if window.Title() != "New Title" {
		t.Errorf("Expected title 'New Title', got '%s'", window.Title())
	}

	if window.Width() != 1024 {
		t.Errorf("Expected width 1024, got %d", window.Width())
	}

	if window.Height() != 768 {
		t.Errorf("Expected height 768, got %d", window.Height())
	}

	if !window.IsMaximized() {
		t.Error("Expected window to be maximized")
	}
}

func TestBaseWindowShowHide(t *testing.T) {
	window := NewBaseWindow(NewWindowOptions{
		State: WindowStateResizable, // Not visible initially
	})

	if window.IsVisible() {
		t.Error("Expected window not to be visible initially")
	}

	err := window.Show()
	if err != nil {
		t.Errorf("Unexpected error showing window: %v", err)
	}

	if !window.IsVisible() {
		t.Error("Expected window to be visible after Show()")
	}

	err = window.Hide()
	if err != nil {
		t.Errorf("Unexpected error hiding window: %v", err)
	}

	if window.IsVisible() {
		t.Error("Expected window not to be visible after Hide()")
	}
}

func TestBaseWindowClose(t *testing.T) {
	window := NewBaseWindow(NewWindowOptions{
		State: WindowStateVisible,
	})

	if window.IsClosed() {
		t.Error("Expected window not to be closed initially")
	}

	err := window.Close()
	if err != nil {
		t.Errorf("Unexpected error closing window: %v", err)
	}

	if !window.IsClosed() {
		t.Error("Expected window to be closed after Close()")
	}

	// Operations on closed window should return error
	err = window.Show()
	if err != ErrWindowNotInitialized {
		t.Errorf("Expected ErrWindowNotInitialized, got %v", err)
	}
}

func TestBaseWindowEventHandling(t *testing.T) {
	window := NewBaseWindow(NewWindowOptions{})

	var receivedEvents []Event
	var mu sync.Mutex

	handler := func(e Event) error {
		mu.Lock()
		receivedEvents = append(receivedEvents, e)
		mu.Unlock()
		return nil
	}

	err := window.Subscribe(EventTypeKeyboard, handler)
	if err != nil {
		t.Errorf("Unexpected error subscribing to events: %v", err)
	}

	// Create and publish a keyboard event
	keyEvent := NewKeyboardEvent()
	keyEvent.Code = KeyCodeKeyA
	keyEvent.State = KeyStatePressed

	err = window.Publish(keyEvent)
	if err != nil {
		t.Errorf("Unexpected error publishing event: %v", err)
	}

	// Give some time for event processing
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	eventCount := len(receivedEvents)
	mu.Unlock()

	if eventCount != 1 {
		t.Errorf("Expected 1 event, got %d", eventCount)
	}

	// Test unsubscribe
	window.Unsubscribe(EventTypeKeyboard, handler)

	// Publish another event
	err = window.Publish(keyEvent)
	if err != nil {
		t.Errorf("Unexpected error publishing event: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	finalEventCount := len(receivedEvents)
	mu.Unlock()

	// Should still be 1 since we unsubscribed
	if finalEventCount != 1 {
		t.Errorf("Expected 1 event after unsubscribe, got %d", finalEventCount)
	}
}

func TestBaseWindowConcurrentAccess(t *testing.T) {
	window := NewBaseWindow(NewWindowOptions{
		Title:  "Concurrent Test",
		Width:  800,
		Height: 600,
	})

	const numGoroutines = 100
	var wg sync.WaitGroup

	// Test concurrent attribute changes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			attr := WindowAttr{
				Title:  fmt.Sprintf("Title %d", id),
				Width:  800 + id,
				Height: 600 + id,
			}

			window.SetAttr(attr)

			// Read operations
			_ = window.Title()
			_ = window.Width()
			_ = window.Height()
			_ = window.IsVisible()
		}(i)
	}

	wg.Wait()

	// Window should still be functional
	if window.IsClosed() {
		t.Error("Window should not be closed after concurrent access")
	}
}

func TestBaseWindowAttrString(t *testing.T) {
	attr := WindowAttr{
		Title:    "Test Window",
		Width:    1024,
		Height:   768,
		Position: image.Point{X: 100, Y: 200},
		State:    WindowStateVisible | WindowStateResizable,
	}

	str := attr.String()

	// Check that string contains key information
	if !contains(str, "Test Window") {
		t.Error("String representation should contain title")
	}

	if !contains(str, "1024") {
		t.Error("String representation should contain width")
	}

	if !contains(str, "768") {
		t.Error("String representation should contain height")
	}
}

func TestBaseWindowStateConsistency(t *testing.T) {
	window := NewBaseWindow(NewWindowOptions{})

	// Test that closed window removes visible state
	attr := WindowAttr{
		State: WindowStateClosed | WindowStateVisible,
	}

	window.SetAttr(attr)

	// Should not be visible if closed
	if window.IsVisible() {
		t.Error("Closed window should not be visible")
	}

	if !window.IsClosed() {
		t.Error("Window should be closed")
	}
}

func TestBaseWindowEmptyApply(t *testing.T) {
	base := WindowAttr{
		Title:  "Original",
		Width:  800,
		Height: 600,
	}

	// Apply with no options should return original
	result := base.Apply()

	if !result.Equal(base) {
		t.Error("Apply with no arguments should return original")
	}
}

func TestBaseWindowMultipleApply(t *testing.T) {
	base := WindowAttr{
		Title:  "Original",
		Width:  800,
		Height: 600,
	}

	override1 := WindowAttr{
		Title: "First Override",
		Width: 1024,
	}

	override2 := WindowAttr{
		Title:  "Second Override",
		Height: 768,
	}

	result := base.Apply(override1, override2)

	if result.Title != "Second Override" {
		t.Errorf("Expected final title 'Second Override', got '%s'", result.Title)
	}

	if result.Width != 1024 {
		t.Errorf("Expected width 1024, got %d", result.Width)
	}

	if result.Height != 768 {
		t.Errorf("Expected height 768, got %d", result.Height)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			hasSubstring(s, substr))))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func BenchmarkBaseWindowSetAttr(b *testing.B) {
	window := NewBaseWindow(NewWindowOptions{})

	attr := WindowAttr{
		Title:  "Benchmark Window",
		Width:  1920,
		Height: 1080,
		State:  WindowStateMaximized,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		window.SetAttr(attr)
	}
}

func BenchmarkBaseWindowStateQueries(b *testing.B) {
	window := NewBaseWindow(NewWindowOptions{
		State: WindowStateVisible | WindowStateResizable | WindowStateDecorated,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = window.IsVisible()
		_ = window.IsResizable()
		_ = window.IsDecorated()
		_ = window.IsClosed()
		_ = window.IsMaximized()
	}
}
