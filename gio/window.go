package gio

import (
	"fmt"
	"image"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

// WindowID represents a unique identifier for a window
type WindowID uint64

// WindowState represents window attributes
type WindowState uint16

const (
	WindowStateUnknown WindowState = iota
	WindowStateClosed  WindowState = 1 << iota
	WindowStateFocused
	WindowStateIconified
	WindowStateMaximized
	WindowStateVisible
	WindowStateHovered
	WindowStateResizable
	WindowStateDecorated
	WindowStateFloating
	WindowStateAutoIconify
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
	Position      image.Point
	State         WindowState
}

func (wa WindowAttr) Equal(other WindowAttr) bool {
	return wa.Title == other.Title &&
		wa.Width == other.Width &&
		wa.Height == other.Height &&
		wa.Position == other.Position &&
		wa.State == other.State
}

func (wa WindowAttr) Apply(o ...WindowAttr) WindowAttr {
	if len(o) == 0 {
		return wa
	}

	result := wa
	for _, attr := range o {
		if attr.Title != "" {
			result.Title = attr.Title
		}
		if attr.Width > 0 {
			result.Width = attr.Width
		}
		if attr.Height > 0 {
			result.Height = attr.Height
		}
		if !attr.Position.Eq(image.Point{}) {
			result.Position = attr.Position
		}
		// Only merge state if it's not the default unknown state
		if attr.State != WindowStateUnknown {
			// Replace state instead of ORing to avoid conflicting states
			result.State = attr.State
		}
	}

	return result
}

// String returns a string representation of the WindowAttr
func (wa WindowAttr) String() string {
	return "(" +
		"Title: '" + wa.Title + "', " +
		"Width: " + strconv.Itoa(wa.Width) + ", " +
		"Height: " + strconv.Itoa(wa.Height) + ", " +
		"Position: " + wa.Position.String() + ", " +
		"State: " + fmt.Sprintf("%016b", wa.State) + ")"
}

type NewWindowOptions = WindowAttr

func DefaultNewWindowOptions() NewWindowOptions {
	config := GetConfig()
	return NewWindowOptions{
		Title:    "Window",
		Width:    config.DefaultWindowWidth,
		Height:   config.DefaultWindowHeight,
		Position: image.Point{X: 100, Y: 100},
		State:    WindowStateFocused | WindowStateVisible | WindowStateResizable | WindowStateDecorated,
	}
}

type WindowHandler unsafe.Pointer

// Window interface with context support
type Window interface {
	BaseWindow
	WindowID() WindowID
}

// BaseWindow interface with improved lifecycle management
type BaseWindow interface {
	Width() int
	Height() int
	Title() string
	Position() image.Point
	State() WindowState

	Show() error
	Hide() error
	Close() error

	Publish(event Event) error
	Subscribe(eventType EventType, handler EventHandler) error
	Unsubscribe(eventType EventType, handler EventHandler)

	Attr() WindowAttr
	SetAttr(attr WindowAttr)

	// Status check methods
	IsClosed() bool
	IsFocused() bool
	IsIconified() bool
	IsMaximized() bool
	IsVisible() bool
	IsResizable() bool
	IsDecorated() bool
	IsFloating() bool
	IsAutoIconify() bool
}

func NewBaseWindow(o NewWindowOptions) BaseWindow {
	attr := DefaultNewWindowOptions().Apply(o)
	config := GetConfig()

	bw := &baseWindow{
		attr:           attr,
		eventBus:       NewEventBus(),
		pressedKeys:    make(map[KeyCode]bool),
		lastKeyPress:   make(map[KeyCode]time.Time),
		keyRepeatDelay: config.KeyRepeatDelay,
		keyRepeatRate:  config.KeyRepeatRate,
	}

	go bw.startKeyboardCleanup()
	return bw
}

var _ BaseWindow = (*baseWindow)(nil)

type baseWindow struct {
	attr      WindowAttr
	eventBus  *EventBus
	destroyed bool
	mu        sync.RWMutex

	pressedKeys map[KeyCode]bool

	// Key repeat detection
	lastKeyPress   map[KeyCode]time.Time // Track last key press time
	keyRepeatDelay time.Duration         // Initial delay before repeat starts
	keyRepeatRate  time.Duration         // Interval between repeats
}

// Title implements BaseWindow.
func (w *baseWindow) Title() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.Title
}

// Width implements BaseWindow.
func (w *baseWindow) Width() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.Width
}

func (w *baseWindow) Attr() WindowAttr {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr
}

// Height implements Window.
func (w *baseWindow) Height() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.Height
}

// Position returns the window position.
func (w *baseWindow) Position() image.Point {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.Position
}

// State returns the window state.
func (w *baseWindow) State() WindowState {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.State
}

// SetAttr sets the window attributes with better validation
func (w *baseWindow) SetAttr(attr WindowAttr) {
	if w.IsClosed() {
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	oldAttr := w.attr
	newAttr := w.attr.Apply(attr)

	// Validate state consistency
	if newAttr.State.Contains(WindowStateClosed) && newAttr.State.Contains(WindowStateVisible) {
		// Remove visible state if window is closed
		newAttr.State = newAttr.State.Diff(WindowStateVisible)
	}

	if oldAttr.Equal(newAttr) {
		return
	}

	w.attr = newAttr

	// Publish event without holding lock
	go func() {
		if err := w.publishAttrChangeEvent(); err != nil {
			// Log error appropriately
			_ = err
		}
	}()
}

// Show shows the window if supported by the platform.
func (w *baseWindow) Show() error {
	if w.IsClosed() {
		return ErrWindowNotInitialized
	}

	w.mu.Lock()
	if w.attr.State.Contains(WindowStateVisible) {
		w.mu.Unlock()
		return nil
	}

	w.attr.State |= WindowStateVisible
	w.mu.Unlock()

	return w.publishAttrChangeEvent()
}

// Hide hides the window if supported by the platform.
func (w *baseWindow) Hide() error {
	if w.IsClosed() {
		return ErrWindowNotInitialized
	}

	w.mu.Lock()
	if !w.attr.State.Contains(WindowStateVisible) {
		w.mu.Unlock()
		return nil
	}

	w.attr.State &= ^WindowStateVisible
	w.mu.Unlock()

	return w.publishAttrChangeEvent()
}

// Close implements Window.
func (w *baseWindow) Close() error {
	if w.IsClosed() {
		return ErrWindowNotInitialized
	}

	w.mu.Lock()
	if w.attr.State.Contains(WindowStateClosed) {
		w.mu.Unlock()
		return nil
	}

	w.attr.State |= WindowStateClosed
	w.mu.Unlock()

	return w.publishAttrChangeEvent()
}

// Publish publishes an event to the window's event bus
func (w *baseWindow) Publish(event Event) error {
	if w.IsClosed() || w.eventBus == nil {
		return ErrWindowNotInitialized
	}

	// Enhanced keyboard event processing
	if keyEvent, ok := event.(*KeyboardEvent); ok {
		processedEvent := w.processKeyboardEvent(keyEvent)
		if processedEvent == nil {
			// Event was filtered out (e.g., fake key release during repeat)
			return nil
		}
		return w.eventBus.Publish(processedEvent)
	}

	return w.eventBus.Publish(event)
}

func (w *baseWindow) Subscribe(eventType EventType, handler EventHandler) error {
	if w.IsClosed() || w.eventBus == nil {
		return ErrWindowNotInitialized
	}
	return w.eventBus.Subscribe(eventType, handler)
}

func (w *baseWindow) Unsubscribe(eventType EventType, handler EventHandler) {
	if w.IsClosed() || w.eventBus == nil {
		return
	}
	w.eventBus.Unsubscribe(eventType, handler)
}

// publishAttrChangeEvent publishes a window attribute change event
func (w *baseWindow) publishAttrChangeEvent() error {
	if w.eventBus == nil {
		return ErrWindowNotInitialized
	}

	event := NewWindowEvent()
	event.Window = w
	return w.eventBus.Publish(event)
}

// IsAutoIconify implements BaseWindow.
func (w *baseWindow) IsAutoIconify() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.State.Contains(WindowStateAutoIconify)
}

// IsClosed implements BaseWindow.
func (w *baseWindow) IsClosed() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.State.Contains(WindowStateClosed)
}

// IsDecorated implements BaseWindow.
func (w *baseWindow) IsDecorated() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.State.Contains(WindowStateDecorated)
}

// IsFloating implements BaseWindow.
func (w *baseWindow) IsFloating() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.State.Contains(WindowStateFloating)
}

// IsFocused implements BaseWindow.
func (w *baseWindow) IsFocused() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.State.Contains(WindowStateFocused)
}

// IsIconified implements BaseWindow.
func (w *baseWindow) IsIconified() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.State.Contains(WindowStateIconified)
}

// IsMaximized implements BaseWindow.
func (w *baseWindow) IsMaximized() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.State.Contains(WindowStateMaximized)
}

// IsResizable implements BaseWindow.
func (w *baseWindow) IsResizable() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.State.Contains(WindowStateResizable)
}

// IsVisible implements BaseWindow.
func (w *baseWindow) IsVisible() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.attr.State.Contains(WindowStateVisible)
}

// String returns a string representation of the base window
func (w *baseWindow) String() string {
	return w.Attr().String()
}

// processKeyboardEvent handles keyboard events with repeat detection
func (w *baseWindow) processKeyboardEvent(event *KeyboardEvent) Event {
	if w.IsClosed() {
		return nil
	}

	switch event.State {
	case KeyStatePressed:
		return w.handleKeyPress(event)
	case KeyStateReleased:
		return w.handleKeyRelease(event)
	default:
		// Unknown state, pass through
		return event
	}
}

// handleKeyPress processes key press events with repeat detection
func (w *baseWindow) handleKeyPress(event *KeyboardEvent) Event {
	keyCode := event.Code
	currentTime := time.Now()

	w.mu.Lock()
	defer w.mu.Unlock()

	// Detect key repeat
	lastPress, exists := w.lastKeyPress[keyCode]
	wasPressed := w.pressedKeys[keyCode]

	if exists && wasPressed {
		// This is a potential repeat
		timeSinceLastPress := currentTime.Sub(lastPress)
		event.Repeat = w.isRepeatTiming(timeSinceLastPress)
	} else {
		// First press or new press after release
		event.Repeat = false
	}

	// Update tracking state
	w.lastKeyPress[keyCode] = currentTime
	w.pressedKeys[keyCode] = true

	return event
}

// handleKeyRelease processes key release events with repeat detection
func (w *baseWindow) handleKeyRelease(event *KeyboardEvent) Event {
	keyCode := event.Code

	w.mu.Lock()
	defer w.mu.Unlock()

	// Check if this is a real release or part of auto-repeat
	if !w.isRealKeyRelease(keyCode) {
		// This is a fake release during auto-repeat, ignore it
		return nil
	}

	// Real key release - clean up state
	delete(w.lastKeyPress, keyCode)
	delete(w.pressedKeys, keyCode)
	event.Repeat = false

	return event
}

// isRepeatTiming determines if the timing indicates a key repeat
func (w *baseWindow) isRepeatTiming(timeSinceLastPress time.Duration) bool {
	// Check if we're in the initial delay period
	if timeSinceLastPress < w.keyRepeatDelay {
		// Too soon for any repeat
		return timeSinceLastPress > 10*time.Millisecond
	}

	// We're past the initial delay, check repeat rate
	repeatPhaseTime := timeSinceLastPress - w.keyRepeatDelay
	expectedRepeats := repeatPhaseTime / w.keyRepeatRate
	actualTime := w.keyRepeatDelay + (expectedRepeats * w.keyRepeatRate)

	// Allow some tolerance for timing variations
	tolerance := 15 * time.Millisecond
	return timeSinceLastPress >= actualTime-tolerance &&
		timeSinceLastPress <= actualTime+tolerance
}

// isRealKeyRelease determines if a key release is genuine (not auto-repeat artifact)
func (w *baseWindow) isRealKeyRelease(keyCode KeyCode) bool {
	// Check if key was actually being tracked as pressed
	if !w.pressedKeys[keyCode] {
		return true // Not tracked as pressed, likely real
	}

	lastPress, exists := w.lastKeyPress[keyCode]
	if !exists {
		return true // No press record, likely real
	}

	timeSincePress := time.Since(lastPress)

	// If release comes too quickly during potential repeat phase, it's likely fake
	if timeSincePress > w.keyRepeatDelay {
		// We're in repeat phase - check if this timing matches auto-repeat pattern
		repeatPhaseTime := timeSincePress - w.keyRepeatDelay

		// Auto-repeat releases happen at very specific intervals
		tolerance := 10 * time.Millisecond
		expectedRepeats := repeatPhaseTime / w.keyRepeatRate
		nextRepeatTime := w.keyRepeatDelay + (expectedRepeats * w.keyRepeatRate)

		// If we're very close to an expected repeat time, this might be fake
		return timeSincePress < nextRepeatTime-tolerance ||
			timeSincePress > nextRepeatTime+tolerance
	}

	// Before repeat delay, releases are usually real
	return true
}

// cleanupOldKeyStates removes stale key states
func (w *baseWindow) cleanupOldKeyStates() {
	if w.IsClosed() {
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	currentTime := time.Now()
	maxAge := 5 * time.Second

	for keyCode, lastPress := range w.lastKeyPress {
		if currentTime.Sub(lastPress) > maxAge {
			delete(w.lastKeyPress, keyCode)
			delete(w.pressedKeys, keyCode)
		}
	}
}

// startKeyboardCleanup starts a background goroutine to clean up stale key states
func (w *baseWindow) startKeyboardCleanup() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if w.IsClosed() {
				return
			}
			w.cleanupOldKeyStates()
		}
	}()
}
