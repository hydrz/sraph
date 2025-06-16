package wayland

import (
	"fmt"
	"image"
	"sync"
	"unsafe"

	"github.com/opensraph/sraph/gio"
)

var _ gio.Window = (*waylandWindow)(nil)

// waylandWindow represents a Wayland window implementation
type waylandWindow struct {
	gio.BaseWindow

	driver       *waylandDriver
	surface      *wlSurface
	shellSurface *wlShellSurface
	buffer       *wlBuffer

	// Mock Wayland window state
	width, height   int
	x, y            int
	focused         bool
	visible         bool
	pointerEntered  bool
	lastPointerX    int
	lastPointerY    int

	mu sync.RWMutex
}

// Mock Wayland types for shell surface and buffer
type wlShellSurface struct {
	ptr unsafe.Pointer
}

type wlBuffer struct {
	ptr unsafe.Pointer
}

func newWaylandWindow(driver *waylandDriver, bw gio.BaseWindow) (*waylandWindow, error) {
	ww := &waylandWindow{
		driver:     driver,
		BaseWindow: bw,
	}

	if err := ww.init(); err != nil {
		return nil, fmt.Errorf("wayland window: init failed: %v", err)
	}

	return ww, nil
}

func (ww *waylandWindow) init() error {
	// Mock Wayland surface creation
	surface := &wlSurface{ptr: unsafe.Pointer(uintptr(100))}
	shellSurface := &wlShellSurface{ptr: unsafe.Pointer(uintptr(101))}

	ww.mu.Lock()
	ww.surface = surface
	ww.shellSurface = shellSurface
	ww.mu.Unlock()

	// Register window with driver
	ww.driver.registerWindow(surface, ww)

	// Set initial attributes
	attr := ww.Attr()
	ww.width = attr.Width
	ww.height = attr.Height
	ww.x = attr.Position.X
	ww.y = attr.Position.Y

	if ww.width <= 0 {
		ww.width = 800
	}
	if ww.height <= 0 {
		ww.height = 600
	}

	// Mock making the window visible
	ww.visible = attr.State.Contains(gio.WindowStateVisible)

	return nil
}

// WindowID returns the Wayland surface as window ID
func (ww *waylandWindow) WindowID() gio.WindowID {
	ww.mu.RLock()
	defer ww.mu.RUnlock()
	if ww.surface != nil {
		return gio.WindowID(uintptr(ww.surface.ptr))
	}
	return 0
}

func (ww *waylandWindow) SetAttr(attr gio.WindowAttr) {
	ww.mu.Lock()
	defer ww.mu.Unlock()

	currentAttr := ww.Attr()
	if currentAttr.Equal(attr) {
		return
	}

	// Handle window closure first
	if attr.State.Contains(gio.WindowStateClosed) {
		ww.handleWindowClose()
		ww.BaseWindow.SetAttr(attr)
		return
	}

	// Update basic properties
	ww.updateTitle(attr, currentAttr)
	ww.updateGeometry(attr, currentAttr)
	ww.updateWindowStates(attr, currentAttr)

	// Update base attributes
	ww.BaseWindow.SetAttr(attr)
}

func (ww *waylandWindow) handleWindowClose() {
	if ww.surface != nil {
		ww.driver.unregisterWindow(ww.surface)
		// In real Wayland, we would call wl_surface_destroy
		ww.surface = nil
	}
	if ww.shellSurface != nil {
		// In real Wayland, we would call wl_shell_surface_destroy
		ww.shellSurface = nil
	}
}

func (ww *waylandWindow) updateTitle(attr, currentAttr gio.WindowAttr) {
	if attr.Title == "" || attr.Title == currentAttr.Title {
		return
	}
	// In real Wayland, we would call wl_shell_surface_set_title
	// For now, we just store it in the base window
}

func (ww *waylandWindow) updateGeometry(attr, currentAttr gio.WindowAttr) {
	sizeChanged := attr.Width != currentAttr.Width || attr.Height != currentAttr.Height

	if sizeChanged && attr.Width > 0 && attr.Height > 0 {
		ww.width = attr.Width
		ww.height = attr.Height
		// In real Wayland, we would need to create a new buffer and attach it
		ww.handleConfigure(ww.width, ww.height)
	}

	posChanged := attr.Position.X != currentAttr.Position.X || attr.Position.Y != currentAttr.Position.Y
	if posChanged {
		ww.x = attr.Position.X
		ww.y = attr.Position.Y
		// Note: In Wayland, client-side positioning is limited
		// Position is mostly controlled by the compositor
	}
}

func (ww *waylandWindow) updateWindowStates(attr, currentAttr gio.WindowAttr) {
	// Handle maximized state
	maximizedChanged := attr.State.Contains(gio.WindowStateMaximized) != currentAttr.State.Contains(gio.WindowStateMaximized)
	if maximizedChanged {
		if attr.State.Contains(gio.WindowStateMaximized) {
			// In real Wayland, call wl_shell_surface_set_maximized
		} else {
			// Unmaximize window
		}
	}

	// Handle fullscreen state
	fullscreenChanged := attr.State.Contains(gio.WindowStateFloating) != currentAttr.State.Contains(gio.WindowStateFloating)
	if fullscreenChanged {
		if attr.State.Contains(gio.WindowStateFloating) {
			// In real Wayland, this might involve setting window as "always on top"
			// which is not directly supported - depends on compositor
		}
	}

	// Handle visible state
	visibleChanged := attr.State.Contains(gio.WindowStateVisible) != currentAttr.State.Contains(gio.WindowStateVisible)
	if visibleChanged {
		ww.visible = attr.State.Contains(gio.WindowStateVisible)
		if ww.visible {
			// In real Wayland, we would commit the surface
		} else {
			// Hide window - in Wayland, typically done by not committing
		}
	}

	// Handle focus state
	focusChanged := attr.State.Contains(gio.WindowStateFocused) != currentAttr.State.Contains(gio.WindowStateFocused)
	if focusChanged {
		ww.focused = attr.State.Contains(gio.WindowStateFocused)
		// In Wayland, focus is controlled by the compositor
		// Client can request focus but it's not guaranteed
	}
}

// Event handling methods

func (ww *waylandWindow) handleFocusIn() {
	ww.mu.Lock()
	if !ww.focused {
		ww.focused = true
		attr := ww.Attr()
		attr.State |= gio.WindowStateFocused
		ww.BaseWindow.SetAttr(attr)
	}
	ww.mu.Unlock()
}

func (ww *waylandWindow) handleFocusOut() {
	ww.mu.Lock()
	if ww.focused {
		ww.focused = false
		attr := ww.Attr()
		attr.State &^= gio.WindowStateFocused
		ww.BaseWindow.SetAttr(attr)
	}
	ww.mu.Unlock()
}

func (ww *waylandWindow) handleKeyEvent(keyCode gio.KeyCode, modifiers gio.ModifierKey, pressed bool) {
	keyEvent := gio.NewKeyboardEvent()
	keyEvent.Code = keyCode
	keyEvent.ModifierKey = modifiers
	keyEvent.Repeat = false // TODO: implement repeat detection

	ww.Publish(keyEvent)
}

func (ww *waylandWindow) handlePointerEnter(x, y int) {
	ww.mu.Lock()
	ww.pointerEntered = true
	ww.lastPointerX = x
	ww.lastPointerY = y
	ww.mu.Unlock()

	// Generate mouse move event
	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = gio.MouseButtonUnknown
	mouseEvent.ModifierKey = gio.ModifierKeyNone
	mouseEvent.Position = image.Point{X: x, Y: y}

	ww.Publish(mouseEvent)
}

func (ww *waylandWindow) handlePointerLeave() {
	ww.mu.Lock()
	ww.pointerEntered = false
	ww.mu.Unlock()
}

func (ww *waylandWindow) handlePointerMotion(x, y int, modifiers gio.ModifierKey) {
	ww.mu.Lock()
	ww.lastPointerX = x
	ww.lastPointerY = y
	ww.mu.Unlock()

	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = gio.MouseButtonUnknown
	mouseEvent.ModifierKey = modifiers
	mouseEvent.Position = image.Point{X: x, Y: y}

	ww.Publish(mouseEvent)
}

func (ww *waylandWindow) handleButtonEvent(button gio.MouseButton, modifiers gio.ModifierKey, pressed bool) {
	ww.mu.RLock()
	x, y := ww.lastPointerX, ww.lastPointerY
	ww.mu.RUnlock()

	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = button
	mouseEvent.ModifierKey = modifiers
	mouseEvent.Position = image.Point{X: x, Y: y}

	ww.Publish(mouseEvent)
}

func (ww *waylandWindow) handleScrollEvent(deltaX, deltaY float64, modifiers gio.ModifierKey) {
	ww.mu.RLock()
	x, y := ww.lastPointerX, ww.lastPointerY
	ww.mu.RUnlock()

	wheelEvent := gio.NewWheelEvent()
	wheelEvent.DeltaX = deltaX
	wheelEvent.DeltaY = deltaY
	wheelEvent.ModifierKey = modifiers
	wheelEvent.Position = image.Point{X: x, Y: y}

	ww.Publish(wheelEvent)
}

func (ww *waylandWindow) handleConfigure(width, height int) {
	ww.mu.Lock()
	sizeChanged := ww.width != width || ww.height != height
	if sizeChanged {
		ww.width = width
		ww.height = height
		
		attr := ww.Attr()
		attr.Width = width
		attr.Height = height
		ww.BaseWindow.SetAttr(attr)
	}
	ww.mu.Unlock()
}
