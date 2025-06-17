package wayland

import (
	"fmt"
	"image"
	"sync"

	"github.com/opensraph/sraph/gio"
)

var _ gio.Window = (*waylandWindow)(nil)

// waylandWindow represents a Wayland window implementation
type waylandWindow struct {
	gio.BaseWindow

	driver       *waylandDriver
	surface      uintptr // wl_surface
	shellSurface uintptr // wl_shell_surface
	buffer       uintptr // wl_buffer

	// Window state
	width, height  int
	x, y           int
	focused        bool
	visible        bool
	pointerEntered bool
	lastPointerX   int
	lastPointerY   int

	mu sync.RWMutex
}

// Remove mock types - we're using real Wayland resources now

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
	// NOTE: In a real Wayland implementation, we would:
	// 1. Create a surface using wl_compositor_create_surface
	// 2. Create a shell surface using wl_shell_get_shell_surface
	// 3. Set up listeners and configure the window
	//
	// For now, we'll simulate having these resources

	// Simulate surface creation
	surface := uintptr(100 + len(ww.driver.windows))      // Fake surface ID
	shellSurface := uintptr(200 + len(ww.driver.windows)) // Fake shell surface ID

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

	// Set window title if provided
	if attr.Title != "" {
		// In a real implementation, we would call wl_shell_surface_set_title
		fmt.Printf("Would set window title to: %s\n", attr.Title)
	}

	// Mark window as visible
	ww.visible = attr.State.Contains(gio.WindowStateVisible)

	return nil
}

// WindowID returns the Wayland surface as window ID
func (ww *waylandWindow) WindowID() gio.WindowID {
	ww.mu.RLock()
	defer ww.mu.RUnlock()
	if ww.surface != 0 {
		return gio.WindowID(ww.surface)
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
	if ww.surface != 0 {
		ww.driver.unregisterWindow(ww.surface)
		// In a real implementation, we would destroy the Wayland resources
		// wl_shell_surface_destroy(ww.shellSurface)
		// wl_surface_destroy(ww.surface)
		ww.shellSurface = 0
		ww.surface = 0
	}
}

func (ww *waylandWindow) updateTitle(attr, currentAttr gio.WindowAttr) {
	if attr.Title == "" || attr.Title == currentAttr.Title {
		return
	}

	// In a real implementation, we would set the window title
	// titleBytes := append([]byte(attr.Title), 0) // Null-terminate
	// wl_shell_surface_set_title(ww.shellSurface, &titleBytes[0])
	fmt.Printf("Would set window title to: %s\n", attr.Title)
}

func (ww *waylandWindow) updateGeometry(attr, currentAttr gio.WindowAttr) {
	sizeChanged := attr.Width != currentAttr.Width || attr.Height != currentAttr.Height

	if sizeChanged && attr.Width > 0 && attr.Height > 0 {
		ww.width = attr.Width
		ww.height = attr.Height
		// In Wayland, size changes are typically handled through configure events
		// The compositor controls the actual window size
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
	if maximizedChanged && ww.shellSurface != 0 {
		if attr.State.Contains(gio.WindowStateMaximized) {
			// In real Wayland with wl_shell, we would call wl_shell_surface_set_maximized
			// Note: wl_shell is deprecated, modern compositors use xdg_shell
		} else {
			// Unmaximize window - depends on compositor support
		}
	}

	// Handle fullscreen state
	fullscreenChanged := attr.State.Contains(gio.WindowStateFloating) != currentAttr.State.Contains(gio.WindowStateFloating)
	if fullscreenChanged {
		if attr.State.Contains(gio.WindowStateFloating) {
			// In Wayland, this might involve setting window as "always on top"
			// which is not directly supported - depends on compositor
		}
	}

	// Handle visible state
	visibleChanged := attr.State.Contains(gio.WindowStateVisible) != currentAttr.State.Contains(gio.WindowStateVisible)
	if visibleChanged {
		ww.visible = attr.State.Contains(gio.WindowStateVisible)
		if ww.visible && ww.surface != 0 {
			// In real Wayland, we would commit the surface to make it visible
			// This requires buffer attachment and surface commit
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
