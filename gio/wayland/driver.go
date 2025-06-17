package wayland

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/opensraph/sraph/gio"
	"github.com/rajveermalviya/go-wayland/wayland/client"
	"github.com/rajveermalviya/go-wayland/wayland/cursor"
	xdg_shell "github.com/rajveermalviya/go-wayland/wayland/stable/xdg-shell"
)

func init() {
	gio.RegisterDriver(gio.DriverTypeWayland, newWaylandDriver)
}

var _ gio.Driver = (*WaylandDriver)(nil)

// WaylandDriver implements the Wayland graphics driver
type WaylandDriver struct {
	// Core Wayland objects
	display     *client.Display
	registry    *client.Registry
	shm         *client.Shm
	compositor  *client.Compositor
	wmBase      *xdg_shell.WmBase
	seat        *client.Seat
	seatVersion uint32

	// Input devices
	keyboard *client.Keyboard
	pointer  *client.Pointer

	// Cursor support
	cursorTheme *cursor.Theme

	// Context and lifecycle
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
	wg      sync.WaitGroup

	// Window management
	windows       map[gio.WindowID]*WaylandWindow
	focusedWindow *WaylandWindow

	// Event handling
	eventState struct {
		pointerEvent  pointerEvent
		keyboardEvent keyboardEvent
	}

	// Data transfer support
	dataDeviceManager *client.DataDeviceManager
	dataDevice        *client.DataDevice // Add missing field
	clipboardManager  *ClipboardManager

	// Serial tracking for clipboard operations
	latestSerial uint32

	mu sync.RWMutex
}

type pointerEvent struct {
	eventMask          int
	surfaceX, surfaceY float64
	button, state      uint32
	time               uint32
	serial             uint32
	axes               [2]struct {
		valid    bool
		value    float64
		discrete int32
	}
	axisSource uint32
}

type keyboardEvent struct {
	serial    uint32
	time      uint32
	key       uint32
	state     uint32
	modifiers uint32
}

// newWaylandDriver creates a new Wayland driver instance
func newWaylandDriver() (gio.Driver, error) {
	display, err := client.Connect("")
	if err != nil {
		return nil, fmt.Errorf("wayland: failed to connect to display: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	driver := &WaylandDriver{
		display: display,
		ctx:     ctx,
		cancel:  cancel,
		running: true,
		windows: make(map[gio.WindowID]*WaylandWindow),
	}

	// Initialize Wayland connection
	if err := driver.init(); err != nil {
		cancel()
		display.Destroy()
		return nil, fmt.Errorf("wayland: failed to initialize driver: %w", err)
	}

	// Start event loop
	driver.wg.Add(1)
	go driver.eventLoop()

	return driver, nil
}

// init sets up the Wayland connection and global objects
func (d *WaylandDriver) init() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Set up display error handler
	d.display.SetErrorHandler(func(event client.DisplayErrorEvent) {
		slog.Error("wayland display error", "objectId", event.ObjectId, "code", event.Code, "message", event.Message)
	})

	// Get registry to enumerate global objects
	registry, err := d.display.GetRegistry()
	if err != nil {
		return fmt.Errorf("failed to get registry: %w", err)
	}
	d.registry = registry

	// Set up registry event handler
	d.registry.SetGlobalHandler(d.handleRegistryGlobal)

	// Wait for global interfaces to be announced
	if err := d.displayRoundTrip(); err != nil {
		return fmt.Errorf("failed to complete display roundtrip: %w", err)
	}

	// Wait for handler events
	if err := d.displayRoundTrip(); err != nil {
		return fmt.Errorf("failed to complete second display roundtrip: %w", err)
	}

	// Verify we have the required globals
	if d.compositor == nil {
		return fmt.Errorf("wl_compositor not available")
	}
	if d.wmBase == nil {
		return fmt.Errorf("xdg_wm_base not available")
	}
	if d.shm == nil {
		return fmt.Errorf("wl_shm not available")
	}

	// Load cursor theme if available
	if d.shm != nil {
		theme, err := cursor.LoadTheme("default", 24, d.shm)
		if err != nil {
			slog.Error("failed to load cursor theme", "error", err)
		} else {
			d.cursorTheme = theme
		}
	}

	slog.Info("wayland driver initialized successfully")
	return nil
}

// handleRegistryGlobal handles global interface announcements
func (d *WaylandDriver) handleRegistryGlobal(event client.RegistryGlobalEvent) {
	slog.Info("found global interface", "interface", event.Interface, "version", event.Version)

	switch event.Interface {
	case "wl_compositor":
		d.bindCompositor(event)
	case "xdg_wm_base":
		d.bindWmBase(event)
	case "wl_seat":
		d.bindSeat(event)
	case "wl_shm":
		d.bindShm(event)
	case "wl_data_device_manager":
		d.bindDataDeviceManager(event)
	}
}

// bindCompositor binds the compositor interface
func (d *WaylandDriver) bindCompositor(event client.RegistryGlobalEvent) {
	compositor := client.NewCompositor(d.display.Context())
	err := d.registry.Bind(event.Name, event.Interface, event.Version, compositor)
	if err != nil {
		slog.Error("failed to bind compositor", "error", err)
		return
	}
	d.compositor = compositor
}

// bindWmBase binds the XDG window manager interface
func (d *WaylandDriver) bindWmBase(event client.RegistryGlobalEvent) {
	wmBase := xdg_shell.NewWmBase(d.display.Context())
	err := d.registry.Bind(event.Name, event.Interface, event.Version, wmBase)
	if err != nil {
		slog.Error("failed to bind xdg_wm_base", "error", err)
		return
	}
	d.wmBase = wmBase

	// Set ping handler
	d.wmBase.SetPingHandler(func(pingEvent xdg_shell.WmBasePingEvent) {
		if err := d.wmBase.Pong(pingEvent.Serial); err != nil {
			slog.Error("failed to respond to ping", "error", err)
		}
	})
}

// bindSeat binds the seat interface for input
func (d *WaylandDriver) bindSeat(event client.RegistryGlobalEvent) {
	seat := client.NewSeat(d.display.Context())
	err := d.registry.Bind(event.Name, event.Interface, event.Version, seat)
	if err != nil {
		slog.Error("failed to bind seat", "error", err)
		return
	}
	d.seat = seat
	d.seatVersion = event.Version

	// Set up seat event handlers
	d.seat.SetCapabilitiesHandler(d.handleSeatCapabilities)
	d.seat.SetNameHandler(func(event client.SeatNameEvent) {
		slog.Info("seat name", "name", event.Name)
	})
}

// bindShm binds the shared memory interface
func (d *WaylandDriver) bindShm(event client.RegistryGlobalEvent) {
	shm := client.NewShm(d.display.Context())
	err := d.registry.Bind(event.Name, event.Interface, event.Version, shm)
	if err != nil {
		slog.Error("failed to bind wl_shm", "error", err)
		return
	}
	d.shm = shm

	d.shm.SetFormatHandler(func(formatEvent client.ShmFormatEvent) {
		slog.Debug("supported pixel format", "format", client.ShmFormat(formatEvent.Format))
	})
}

// bindDataDeviceManager binds the data device manager interface
func (d *WaylandDriver) bindDataDeviceManager(event client.RegistryGlobalEvent) {
	dataDeviceManager := client.NewDataDeviceManager(d.display.Context())
	err := d.registry.Bind(event.Name, event.Interface, event.Version, dataDeviceManager)
	if err != nil {
		slog.Error("failed to bind data device manager", "error", err)
		return
	}
	d.dataDeviceManager = dataDeviceManager

	// Create data device if we have a seat
	if d.seat != nil {
		dataDevice, err := dataDeviceManager.GetDataDevice(d.seat)
		if err != nil {
			slog.Error("failed to create data device", "error", err)
		} else {
			d.dataDevice = dataDevice
		}
	}

	// Initialize clipboard manager
	d.clipboardManager = NewClipboardManager(d)
	if err := d.clipboardManager.Initialize(); err != nil {
		slog.Error("failed to initialize clipboard manager", "error", err)
	}
}

// eventLoop runs the main Wayland event loop
func (d *WaylandDriver) eventLoop() {
	defer d.wg.Done()

	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			if !d.running {
				return
			}

			// Dispatch pending events
			if err := d.display.Context().Dispatch(); err != nil {
				slog.Error("dispatch error", "error", err)
				continue
			}
		}
	}
}

// displayRoundTrip performs a display roundtrip
func (d *WaylandDriver) displayRoundTrip() error {
	callback, err := d.display.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync display: %w", err)
	}
	defer callback.Destroy()

	done := make(chan bool, 1)
	callback.SetDoneHandler(func(event client.CallbackDoneEvent) {
		done <- true
	})

	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-timeout:
			return fmt.Errorf("timeout waiting for display roundtrip")
		case <-ticker.C:
			if err := d.display.Context().Dispatch(); err != nil {
				return fmt.Errorf("dispatch error during roundtrip: %w", err)
			}
		}
	}
}

// Type returns the driver type
func (d *WaylandDriver) Type() gio.DriverType {
	return gio.DriverTypeWayland
}

// CreateWindow creates a new Wayland window
func (d *WaylandDriver) CreateWindow(opts gio.NewWindowOptions) (gio.Window, error) {
	if !d.running {
		return nil, gio.NewError(gio.ErrorCodeDriverNotFound, "driver not running", nil)
	}

	bw := gio.NewBaseWindow(opts)
	window, err := newWaylandWindow(d, bw)
	if err != nil {
		return nil, err
	}

	// Register window
	d.mu.Lock()
	d.windows[window.WindowID()] = window
	d.mu.Unlock()

	return window, nil
}

// handleSeatCapabilities handles seat capability changes
func (d *WaylandDriver) handleSeatCapabilities(event client.SeatCapabilitiesEvent) {
	// Handle pointer capability
	havePointer := (event.Capabilities & uint32(client.SeatCapabilityPointer)) != 0
	if havePointer && d.pointer == nil {
		d.attachPointer()
	} else if !havePointer && d.pointer != nil {
		d.releasePointer()
	}

	// Handle keyboard capability
	haveKeyboard := (event.Capabilities & uint32(client.SeatCapabilityKeyboard)) != 0
	if haveKeyboard && d.keyboard == nil {
		d.attachKeyboard()
	} else if !haveKeyboard && d.keyboard != nil {
		d.releaseKeyboard()
	}
}

// attachKeyboard sets up the keyboard interface with better error handling
func (d *WaylandDriver) attachKeyboard() {
	keyboard, err := d.seat.GetKeyboard()
	if err != nil {
		slog.Error("failed to get keyboard", "error", err)
		return
	}
	d.keyboard = keyboard

	// Set up keyboard event handlers - order matters for proper modifier tracking
	d.keyboard.SetModifiersHandler(d.handleKeyboardModifiers)
	d.keyboard.SetKeyHandler(d.handleKeyboardKey)
	d.keyboard.SetKeymapHandler(d.handleKeyboardKeymap)
	d.keyboard.SetEnterHandler(d.handleKeyboardEnter)
	d.keyboard.SetLeaveHandler(d.handleKeyboardLeave)
	d.keyboard.SetRepeatInfoHandler(d.handleKeyboardRepeatInfo)

	// Initialize modifier state
	d.eventState.keyboardEvent.modifiers = 0

	slog.Info("keyboard interface registered")
}

// attachPointer sets up the pointer interface
func (d *WaylandDriver) attachPointer() {
	pointer, err := d.seat.GetPointer()
	if err != nil {
		slog.Error("failed to get pointer", "error", err)
		return
	}
	d.pointer = pointer

	// Set up pointer event handlers
	d.pointer.SetEnterHandler(d.handlePointerEnter)
	d.pointer.SetLeaveHandler(d.handlePointerLeave)
	d.pointer.SetMotionHandler(d.handlePointerMotion)
	d.pointer.SetButtonHandler(d.handlePointerButton)
	d.pointer.SetAxisHandler(d.handlePointerAxis)

	// Add missing frame handler - this is crucial for event processing
	d.pointer.SetFrameHandler(d.handlePointerFrame)

	// Add optional axis handlers if available
	if d.seatVersion >= 5 {
		d.pointer.SetAxisSourceHandler(d.handlePointerAxisSource)
		d.pointer.SetAxisStopHandler(d.handlePointerAxisStop)
		d.pointer.SetAxisDiscreteHandler(d.handlePointerAxisDiscrete)
	}

	slog.Info("pointer interface registered")
}

// releasePointer releases the pointer interface
func (d *WaylandDriver) releasePointer() {
	if d.pointer != nil && d.seatVersion >= 3 {
		if err := d.pointer.Release(); err != nil {
			slog.Error("failed to release pointer", "error", err)
		}
	}
	d.pointer = nil
	slog.Info("pointer interface released")
}

// resetPointerEvent resets event state
func (d *WaylandDriver) resetPointerEvent(pe *pointerEvent) {
	d.eventState.pointerEvent = pointerEvent{
		surfaceX: pe.surfaceX,
		surfaceY: pe.surfaceY,
	}
}

// updateSerial updates the latest serial number for clipboard operations
func (d *WaylandDriver) updateSerial(serial uint32) {
	if serial > d.latestSerial {
		d.latestSerial = serial
	}
}

// getLatestSerial returns the latest serial number
func (d *WaylandDriver) getLatestSerial() uint32 {
	return d.latestSerial
}

// Shutdown stops the Wayland driver
func (d *WaylandDriver) Shutdown() error {
	d.mu.Lock()
	d.running = false
	d.mu.Unlock()

	// Cancel context to stop event loop
	d.cancel()
	d.wg.Wait()

	// Clean up resources
	d.mu.Lock()
	defer d.mu.Unlock()

	// Close all windows
	for _, window := range d.windows {
		if window != nil {
			window.Close()
		}
	}
	d.windows = make(map[gio.WindowID]*WaylandWindow)

	// Release input devices
	d.releasePointer()
	d.releaseKeyboard()

	// Clean up cursor theme
	if d.cursorTheme != nil {
		d.cursorTheme.Destroy()
		d.cursorTheme = nil
	}

	// Clean up Wayland objects
	if d.wmBase != nil {
		d.wmBase.Destroy()
		d.wmBase = nil
	}
	if d.compositor != nil {
		d.compositor.Destroy()
		d.compositor = nil
	}
	if d.shm != nil {
		d.shm.Destroy()
		d.shm = nil
	}
	if d.registry != nil {
		d.registry.Destroy()
		d.registry = nil
	}
	if d.display != nil {
		d.display.Destroy()
		d.display = nil
	}

	// Clean up clipboard manager
	if d.clipboardManager != nil {
		d.clipboardManager.Cleanup()
		d.clipboardManager = nil
	}

	// Clean up data device manager
	if d.dataDeviceManager != nil {
		d.dataDeviceManager.Destroy()
		d.dataDeviceManager = nil
	}

	return nil
}

// removeWindow removes a window from tracking
func (d *WaylandDriver) removeWindow(windowID gio.WindowID) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if window, exists := d.windows[windowID]; exists {
		if d.focusedWindow == window {
			d.focusedWindow = nil
		}
		delete(d.windows, windowID)
	}
}

// findWindowBySurface finds a window by its surface
func (d *WaylandDriver) findWindowBySurface(surface *client.Surface) *WaylandWindow {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, window := range d.windows {
		if window.surface == surface {
			return window
		}
	}
	return nil
}

// setCursor sets the cursor for the pointer
func (d *WaylandDriver) setCursor(serial uint32, name string) {
	if d.cursorTheme == nil || d.pointer == nil {
		return
	}

	cursor := d.cursorTheme.GetCursor(name)
	if cursor == nil {
		return
	}

	image := cursor.Images[0]

	surface, err := d.compositor.CreateSurface()
	if err != nil {
		slog.Error("failed to create cursor surface", "error", err)
		return
	}

	buffer, err := image.GetBuffer()
	if err != nil {
		slog.Error("failed to get cursor buffer", "error", err)
		return
	}

	if buffer != nil {
		if err := surface.Attach(buffer, 0, 0); err != nil {
			slog.Error("failed to attach cursor buffer", "error", err)
			return
		}
		if err := surface.Damage(0, 0, int32(image.Width), int32(image.Height)); err != nil {
			slog.Error("failed to damage cursor surface", "error", err)
		}
		if err := surface.Commit(); err != nil {
			slog.Error("failed to commit cursor surface", "error", err)
			return
		}

		hotspotX := int32(image.HotspotX)
		hotspotY := int32(image.HotspotY)
		if err := d.pointer.SetCursor(serial, surface, hotspotX, hotspotY); err != nil {
			slog.Error("failed to set cursor", "error", err)
		}
	}
}
