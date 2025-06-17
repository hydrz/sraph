package wayland

import (
	"fmt"
	"image"
	"log"
	"sync"
	"time"

	"github.com/opensraph/sraph/gio"
	"github.com/rajveermalviya/go-wayland/wayland/client"
	xdg_shell "github.com/rajveermalviya/go-wayland/wayland/stable/xdg-shell"
)

// WaylandDriver implements the Wayland graphics driver
type WaylandDriver struct {
	display    *client.Display
	registry   *client.Registry
	compositor *client.Compositor
	wmBase     *xdg_shell.WmBase
	shm        *client.Shm
	seat       *client.Seat
	keyboard   *client.Keyboard
	pointer    *client.Pointer
	ctx        *client.Context
	running    bool
	mu         sync.RWMutex
}

// NewWaylandDriver creates a new Wayland driver instance
func NewWaylandDriver() (gio.Driver, error) {
	display, err := client.Connect("")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Wayland display: %w", err)
	}

	driver := &WaylandDriver{
		display: display,
		running: true,
	}

	// Initialize Wayland connection
	if err := driver.initialize(); err != nil {
		display.Destroy()
		return nil, fmt.Errorf("failed to initialize Wayland driver: %w", err)
	}

	return driver, nil
}

// initialize sets up the Wayland connection and global objects
func (d *WaylandDriver) initialize() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Get the context from display
	d.ctx = d.display.Context()

	// Get registry to enumerate global objects
	registry, err := d.display.GetRegistry()
	if err != nil {
		return fmt.Errorf("failed to get registry: %w", err)
	}
	d.registry = registry

	// Set up registry event handler
	d.registry.SetGlobalHandler(func(event client.RegistryGlobalEvent) {
		log.Printf("Found global: %s (version %d, id %d)", event.Interface, event.Version, event.Name)

		switch event.Interface {
		case "wl_compositor":
			// Create and bind compositor
			compositor := client.NewCompositor(d.ctx)
			err := d.registry.Bind(event.Name, event.Interface, event.Version, compositor)
			if err != nil {
				log.Printf("Failed to bind compositor: %v", err)
				return
			}
			d.compositor = compositor
			log.Printf("Bound compositor successfully")

		case "xdg_wm_base":
			// Create and bind XDG window manager
			wmBase := xdg_shell.NewWmBase(d.ctx)
			err := d.registry.Bind(event.Name, event.Interface, event.Version, wmBase)
			if err != nil {
				log.Printf("Failed to bind xdg_wm_base: %v", err)
				return
			}
			d.wmBase = wmBase

			// Set ping handler
			d.wmBase.SetPingHandler(func(pingEvent xdg_shell.WmBasePingEvent) {
				d.wmBase.Pong(pingEvent.Serial)
			})
			log.Printf("Bound xdg_wm_base successfully")

		case "wl_seat":
			// Create and bind seat for input events
			seat := client.NewSeat(d.ctx)
			err := d.registry.Bind(event.Name, event.Interface, event.Version, seat)
			if err != nil {
				log.Printf("Failed to bind seat: %v", err)
				return
			}
			d.seat = seat
			d.setupSeat()
			log.Printf("Bound seat successfully")
			
		case "wl_shm":
			// Create and bind shared memory
			shm := client.NewShm(d.ctx)
			err := d.registry.Bind(event.Name, event.Interface, event.Version, shm)
			if err != nil {
				log.Printf("Failed to bind wl_shm: %v", err)
				return
			}
			d.shm = shm
			log.Printf("Bound wl_shm successfully")
		}
	})

	// Roundtrip to ensure we get all globals
	callback, err := d.display.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync display: %w", err)
	}

	// Set up callback to know when sync is done
	syncDone := make(chan bool, 1)
	callback.SetDoneHandler(func(event client.CallbackDoneEvent) {
		log.Printf("Display sync completed")
		syncDone <- true
	})

	// Flush and wait for events with timeout
	log.Printf("Waiting for registry sync...")
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-syncDone:
			log.Printf("Registry sync completed")
			goto syncComplete
		case <-timeout:
			return fmt.Errorf("timeout waiting for registry sync")
		case <-ticker.C:
			// Continue processing events while waiting
			if err := d.ctx.Dispatch(); err != nil {
				log.Printf("Dispatch error during sync: %v", err)
				// Continue anyway
			}
		}
	}

syncComplete:

	// Verify we have the required globals
	if d.compositor == nil {
		return fmt.Errorf("wl_compositor not available")
	}
	if d.wmBase == nil {
		return fmt.Errorf("xdg_wm_base not available")
	}

	log.Printf("Wayland driver initialized successfully")
	return nil
}

// setupSeat configures input devices from the seat
func (d *WaylandDriver) setupSeat() {
	if d.seat == nil {
		return
	}

	d.seat.SetCapabilitiesHandler(func(event client.SeatCapabilitiesEvent) {
		// Setup keyboard if available
		if event.Capabilities&uint32(client.SeatCapabilityKeyboard) != 0 {
			keyboard, err := d.seat.GetKeyboard()
			if err == nil {
				d.keyboard = keyboard
				d.setupKeyboard()
			}
		}

		// Setup pointer if available
		if event.Capabilities&uint32(client.SeatCapabilityPointer) != 0 {
			pointer, err := d.seat.GetPointer()
			if err == nil {
				d.pointer = pointer
				d.setupPointer()
			}
		}
	})
}

// setupKeyboard configures keyboard event handlers
func (d *WaylandDriver) setupKeyboard() {
	if d.keyboard == nil {
		return
	}

	d.keyboard.SetKeyHandler(func(event client.KeyboardKeyEvent) {
		// Convert Wayland key event to gio keyboard event
		keyEvent := gio.NewKeyboardEvent()
		keyEvent.Code = gio.KeyCode(event.Key)
		if event.State == uint32(client.KeyboardKeyStatePressed) {
			keyEvent.State = gio.KeyStatePressed
		} else {
			keyEvent.State = gio.KeyStateReleased
		}
		// Publish event to active window
		// TODO: Route to appropriate window
	})

	d.keyboard.SetModifiersHandler(func(event client.KeyboardModifiersEvent) {
		// Handle modifier key changes
		// TODO: Update modifier state
	})
}

// setupPointer configures pointer event handlers
func (d *WaylandDriver) setupPointer() {
	if d.pointer == nil {
		return
	}

	d.pointer.SetButtonHandler(func(event client.PointerButtonEvent) {
		// Convert Wayland pointer button event to gio pointer event
		pointerEvent := gio.NewPointerEvent()
		pointerEvent.PointerType = gio.PointerTypeMouse
		pointerEvent.MouseButton = gio.MouseButton(event.Button)
		if event.State == uint32(client.PointerButtonStatePressed) {
			pointerEvent.KeyState = gio.KeyStatePressed
		} else {
			pointerEvent.KeyState = gio.KeyStateReleased
		}
		// TODO: Route to appropriate window
	})

	d.pointer.SetMotionHandler(func(event client.PointerMotionEvent) {
		// Handle pointer motion
		pointerEvent := gio.NewPointerEvent()
		pointerEvent.PointerType = gio.PointerTypeMouse
		pointerEvent.Position = image.Point{
			X: int(event.SurfaceX),
			Y: int(event.SurfaceY),
		}
		// TODO: Route to appropriate window
	})

	d.pointer.SetEnterHandler(func(event client.PointerEnterEvent) {
		// Handle pointer enter surface
		// TODO: Set focus to window
	})

	d.pointer.SetLeaveHandler(func(event client.PointerLeaveEvent) {
		// Handle pointer leave surface
		// TODO: Remove focus from window
	})
}

// Type returns the driver type
func (d *WaylandDriver) Type() gio.DriverType {
	return gio.DriverTypeWayland
}

// CreateWindow creates a new Wayland window
func (d *WaylandDriver) CreateWindow(opts gio.NewWindowOptions) (gio.Window, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if !d.running {
		return nil, gio.NewError(gio.ErrorCodeDriverNotFound, "Wayland driver not running", nil)
	}

	if d.compositor == nil || d.wmBase == nil {
		return nil, gio.NewError(gio.ErrorCodeDriverInitFailed, "Wayland compositor or window manager not available", nil)
	}

	// Create surface
	surface, err := d.compositor.CreateSurface()
	if err != nil {
		return nil, fmt.Errorf("failed to create surface: %w", err)
	}

	// Create XDG surface
	xdgSurface, err := d.wmBase.GetXdgSurface(surface)
	if err != nil {
		return nil, fmt.Errorf("failed to create XDG surface: %w", err)
	}

	// Create toplevel
	toplevel, err := xdgSurface.GetToplevel()
	if err != nil {
		return nil, fmt.Errorf("failed to create toplevel: %w", err)
	}

	// Create window wrapper
	window := &WaylandWindow{
		BaseWindow: gio.NewBaseWindow(opts),
		driver:     d,
		surface:    surface,
		xdgSurface: xdgSurface,
		toplevel:   toplevel,
		id:         gio.WindowID(time.Now().UnixNano()),
	}

	// Configure window properties
	if err := window.configure(); err != nil {
		return nil, fmt.Errorf("failed to configure window: %w", err)
	}

	return window, nil
}

// Run starts the Wayland event loop
func (d *WaylandDriver) Run() error {
	if !d.running {
		return gio.NewError(gio.ErrorCodeInvalidWindowState, "driver not running", nil)
	}

	log.Printf("Starting Wayland event loop")

	// Main event loop
	for d.running {
		// Dispatch pending events
		if err := d.ctx.Dispatch(); err != nil {
			log.Printf("Wayland dispatch error: %v", err)
			// Don't return error immediately, try to continue
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// Small sleep to prevent busy loop
		time.Sleep(1 * time.Millisecond)
	}

	log.Printf("Wayland event loop stopped")
	return nil
}

// Shutdown stops the Wayland driver
func (d *WaylandDriver) Shutdown() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.running = false

	// Clean up resources
	if d.keyboard != nil {
		// Keyboard doesn't have Destroy method, just set to nil
		d.keyboard = nil
	}
	if d.pointer != nil {
		// Pointer doesn't have Destroy method, just set to nil
		d.pointer = nil
	}
	if d.seat != nil {
		// Seat doesn't have Destroy method, just set to nil
		d.seat = nil
	}
	if d.wmBase != nil {
		d.wmBase.Destroy()
	}
	if d.compositor != nil {
		d.compositor.Destroy()
	}
	if d.registry != nil {
		d.registry.Destroy()
	}
	if d.display != nil {
		d.display.Destroy()
	}
	if d.ctx != nil {
		d.ctx.Close()
	}

	return nil
}

// WaylandWindow represents a Wayland window
type WaylandWindow struct {
	gio.BaseWindow
	driver     *WaylandDriver
	surface    *client.Surface
	xdgSurface *xdg_shell.Surface
	toplevel   *xdg_shell.Toplevel
	id         gio.WindowID
	configured bool
	mu         sync.RWMutex
}

// WindowID returns the unique window identifier
func (w *WaylandWindow) WindowID() gio.WindowID {
	return w.id
}

// configure sets up the window with initial properties
func (w *WaylandWindow) configure() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	attr := w.BaseWindow.Attr()

	// Set window title
	if err := w.toplevel.SetTitle(attr.Title); err != nil {
		return fmt.Errorf("failed to set title: %w", err)
	}

	// Set window app ID (use title as default)
	if err := w.toplevel.SetAppId(attr.Title); err != nil {
		return fmt.Errorf("failed to set app ID: %w", err)
	}

	// Set window size constraints
	if attr.Width > 0 && attr.Height > 0 {
		if err := w.toplevel.SetMinSize(int32(attr.Width), int32(attr.Height)); err != nil {
			return fmt.Errorf("failed to set min size: %w", err)
		}
	}

	// Set up event handlers
	w.setupEventHandlers()

	// Commit the surface to make it visible
	if err := w.surface.Commit(); err != nil {
		return fmt.Errorf("failed to commit surface: %w", err)
	}

	return nil
}

// setupEventHandlers configures window event handlers
func (w *WaylandWindow) setupEventHandlers() {
	// Handle window close requests
	w.toplevel.SetCloseHandler(func(event xdg_shell.ToplevelCloseEvent) {
		w.Close()
	})

	// Handle window configuration changes
	w.toplevel.SetConfigureHandler(func(event xdg_shell.ToplevelConfigureEvent) {
		w.mu.Lock()
		defer w.mu.Unlock()

		// Update window attributes based on configuration
		if event.Width > 0 && event.Height > 0 {
			attr := w.BaseWindow.Attr()
			attr.Width = int(event.Width)
			attr.Height = int(event.Height)
			w.BaseWindow.SetAttr(attr)
		}

		// Update window state
		attr := w.BaseWindow.Attr()
		newState := attr.State

		// Check states array for current state
		for _, stateData := range event.States {
			// States is a byte array, convert to uint32 to compare with constants
			state := uint32(stateData)
			switch state {
			case uint32(xdg_shell.ToplevelStateMaximized):
				newState |= gio.WindowStateMaximized
			case uint32(xdg_shell.ToplevelStateFullscreen):
				// Handle fullscreen state
			case uint32(xdg_shell.ToplevelStateActivated):
				newState |= gio.WindowStateFocused
			}
		}

		if newState != attr.State {
			attr.State = newState
			w.BaseWindow.SetAttr(attr)
		}

		w.configured = true
	})

	// Handle XDG surface configuration
	w.xdgSurface.SetConfigureHandler(func(event xdg_shell.SurfaceConfigureEvent) {
		// Acknowledge the configure event
		w.xdgSurface.AckConfigure(event.Serial)
	})
}

// Show makes the window visible
func (w *WaylandWindow) Show() error {
	if err := w.BaseWindow.Show(); err != nil {
		return err
	}

	// Commit surface to make changes visible
	return w.surface.Commit()
}

// Hide hides the window
func (w *WaylandWindow) Hide() error {
	if err := w.BaseWindow.Hide(); err != nil {
		return err
	}

	// Commit surface to apply changes
	return w.surface.Commit()
}

// Close closes the window and cleans up resources
func (w *WaylandWindow) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.BaseWindow.Close(); err != nil {
		return err
	}

	// Clean up Wayland resources
	if w.toplevel != nil {
		w.toplevel.Destroy()
		w.toplevel = nil
	}
	if w.xdgSurface != nil {
		w.xdgSurface.Destroy()
		w.xdgSurface = nil
	}
	if w.surface != nil {
		w.surface.Destroy()
		w.surface = nil
	}

	return nil
}

// SetAttr updates window attributes
func (w *WaylandWindow) SetAttr(attr gio.WindowAttr) {
	w.BaseWindow.SetAttr(attr)

	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.toplevel == nil {
		return
	}

	// Update title if changed
	if w.BaseWindow.Title() != attr.Title {
		w.toplevel.SetTitle(attr.Title)
	}

	// Handle maximize/minimize state changes
	if attr.State.Contains(gio.WindowStateMaximized) && !w.BaseWindow.State().Contains(gio.WindowStateMaximized) {
		w.toplevel.SetMaximized()
	} else if !attr.State.Contains(gio.WindowStateMaximized) && w.BaseWindow.State().Contains(gio.WindowStateMaximized) {
		w.toplevel.UnsetMaximized()
	}

	// Commit changes
	if w.surface != nil {
		w.surface.Commit()
	}
}

// init registers the Wayland driver
func init() {
	gio.RegisterDriver(gio.DriverTypeWayland, func() (gio.Driver, error) {
		return NewWaylandDriver()
	})
}
