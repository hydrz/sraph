package wayland

import (
	"fmt"
	"image"
	"log"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/opensraph/sraph/gio"
	"github.com/rajveermalviya/go-wayland/wayland/client"
	xdg_shell "github.com/rajveermalviya/go-wayland/wayland/stable/xdg-shell"
)

func init() {
	// Register the Wayland driver
	gio.RegisterDriver(gio.DriverTypeWayland, newWaylandDriver)
}

var _ gio.Driver = (*WaylandDriver)(nil)

// newWaylandDriver creates a new Wayland driver instance (matching X11 pattern)
func newWaylandDriver() (gio.Driver, error) {
	display, err := client.Connect("")
	if err != nil {
		return nil, fmt.Errorf("wayland: failed to connect to display: %w", err)
	}

	driver := &WaylandDriver{
		display: display,
		running: true,
	}

	// Initialize Wayland connection
	if err := driver.init(); err != nil {
		display.Destroy()
		return nil, fmt.Errorf("wayland: failed to initialize driver: %w", err)
	}

	return driver, nil
}

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

// init sets up the Wayland connection and global objects
func (d *WaylandDriver) init() error {
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
	if d.shm == nil {
		return fmt.Errorf("wl_shm not available")
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
	bw := gio.NewBaseWindow(opts)
	return newWaylandWindow(d, bw)
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

var _ gio.Window = (*WaylandWindow)(nil)

// WaylandWindow represents a Wayland window implementation
type WaylandWindow struct {
	gio.BaseWindow

	wlDriver   *WaylandDriver
	wlSurface  *client.Surface
	xdgSurface *xdg_shell.Surface
	toplevel   *xdg_shell.Toplevel
	buffer     *client.Buffer
	shmPool    *client.ShmPool

	// Window state tracking
	configured  bool
	needsRedraw bool

	mu sync.RWMutex
}

// WindowID returns the Wayland window ID (surface pointer as uint64)
func (w *WaylandWindow) WindowID() gio.WindowID {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return gio.WindowID(uintptr(unsafe.Pointer(w.wlSurface)))
}

// init sets up the window with initial properties
func (w *WaylandWindow) init() error {
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
	if err := w.wlSurface.Commit(); err != nil {
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
		w.needsRedraw = true

		// Trigger initial render after configuration
		go func() {
			if err := w.Render(); err != nil {
				log.Printf("Failed to render window: %v", err)
			}
		}()
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

	w.needsRedraw = true

	// Trigger render
	if w.configured {
		go func() {
			if err := w.Render(); err != nil {
				log.Printf("Failed to render window: %v", err)
			}
		}()
	}

	// Commit surface to make changes visible
	return w.wlSurface.Commit()
}

// Hide hides the window
func (w *WaylandWindow) Hide() error {
	if err := w.BaseWindow.Hide(); err != nil {
		return err
	}

	// Commit surface to apply changes
	return w.wlSurface.Commit()
}

// Close closes the window and cleans up resources
func (w *WaylandWindow) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.BaseWindow.Close(); err != nil {
		return err
	}

	// Clean up Wayland resources
	if w.buffer != nil {
		w.buffer.Destroy()
		w.buffer = nil
	}
	if w.shmPool != nil {
		w.shmPool.Destroy()
		w.shmPool = nil
	}
	if w.toplevel != nil {
		w.toplevel.Destroy()
		w.toplevel = nil
	}
	if w.xdgSurface != nil {
		w.xdgSurface.Destroy()
		w.xdgSurface = nil
	}
	if w.wlSurface != nil {
		w.wlSurface.Destroy()
		w.wlSurface = nil
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
	if w.wlSurface != nil {
		w.wlSurface.Commit()
	}
}

// newWaylandWindow creates a new Wayland window
func newWaylandWindow(driver *WaylandDriver, bw gio.BaseWindow) (*WaylandWindow, error) {
	driver.mu.RLock()
	defer driver.mu.RUnlock()

	if !driver.running {
		return nil, gio.NewError(gio.ErrorCodeDriverNotFound, "Wayland driver not running", nil)
	}

	if driver.compositor == nil || driver.wmBase == nil {
		return nil, gio.NewError(gio.ErrorCodeDriverInitFailed, "Wayland compositor or window manager not available", nil)
	}

	// Create surface
	surface, err := driver.compositor.CreateSurface()
	if err != nil {
		return nil, fmt.Errorf("wayland: failed to create surface: %w", err)
	}

	// Create XDG surface
	xdgSurface, err := driver.wmBase.GetXdgSurface(surface)
	if err != nil {
		return nil, fmt.Errorf("wayland: failed to create XDG surface: %w", err)
	}

	// Create toplevel
	toplevel, err := xdgSurface.GetToplevel()
	if err != nil {
		return nil, fmt.Errorf("wayland: failed to create toplevel: %w", err)
	}

	// Create window wrapper
	ww := &WaylandWindow{
		BaseWindow:  bw,
		wlDriver:    driver,
		wlSurface:   surface,
		xdgSurface:  xdgSurface,
		toplevel:    toplevel,
		needsRedraw: true,
	}

	// Configure window properties
	if err := ww.init(); err != nil {
		return nil, fmt.Errorf("wayland: failed to initialize window: %w", err)
	}

	return ww, nil
}

// createBuffer creates a simple colored buffer for the window
func (w *WaylandWindow) createBuffer(width, height int32) error {
	if w.wlDriver.shm == nil {
		return fmt.Errorf("wl_shm not available")
	}

	// Calculate buffer size (ARGB8888 format)
	stride := width * 4
	size := stride * height

	// Create anonymous file for shared memory
	fd, err := w.createAnonymousFile(int(size))
	if err != nil {
		return fmt.Errorf("failed to create anonymous file: %w", err)
	}
	defer syscall.Close(fd)

	// Map the memory
	data, err := syscall.Mmap(fd, 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return fmt.Errorf("failed to mmap: %w", err)
	}

	// Fill with a simple gradient pattern
	w.fillBuffer(data, int(width), int(height), int(stride))

	// Unmap the memory
	if err := syscall.Munmap(data); err != nil {
		return fmt.Errorf("failed to munmap: %w", err)
	}

	// Create shared memory pool
	pool, err := w.wlDriver.shm.CreatePool(fd, size)
	if err != nil {
		return fmt.Errorf("failed to create shm pool: %w", err)
	}
	w.shmPool = pool

	// Create buffer from pool
	buffer, err := pool.CreateBuffer(0, width, height, stride, uint32(client.ShmFormatArgb8888))
	if err != nil {
		return fmt.Errorf("failed to create buffer: %w", err)
	}
	w.buffer = buffer

	return nil
}

// createAnonymousFile creates an anonymous file for shared memory
func (w *WaylandWindow) createAnonymousFile(size int) (int, error) {
	// Try to use memfd_create (Linux-specific)
	name := "wayland-buffer"
	fd, _, err := syscall.RawSyscall(319, uintptr(unsafe.Pointer(&[]byte(name)[0])), 0, 0) // SYS_memfd_create
	if err == 0 {
		if err := syscall.Ftruncate(int(fd), int64(size)); err != nil {
			syscall.Close(int(fd))
			return 0, err
		}
		return int(fd), nil
	}
	// Fallback to creating a temporary file
	tmpFile, fileErr := os.CreateTemp("", "wayland-buffer-")
	if fileErr != nil {
		return 0, fileErr
	}

	if err := tmpFile.Truncate(int64(size)); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return 0, err
	}

	fileFd := int(tmpFile.Fd())

	// Unlink the file so it gets deleted when closed
	os.Remove(tmpFile.Name())

	return fileFd, nil
}

// fillBuffer fills the buffer with a simple gradient pattern
func (w *WaylandWindow) fillBuffer(data []byte, width, height, stride int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := y*stride + x*4

			// Create a simple gradient from blue to red
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(128)
			a := uint8(255)

			// ARGB8888 format
			data[offset+0] = b // Blue
			data[offset+1] = g // Green
			data[offset+2] = r // Red
			data[offset+3] = a // Alpha
		}
	}
}

// Render performs a frame render (placeholder for actual rendering)
func (w *WaylandWindow) Render() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.configured || !w.needsRedraw {
		return nil
	}

	attr := w.BaseWindow.Attr()

	// Create buffer if needed
	if w.buffer == nil {
		if err := w.createBuffer(int32(attr.Width), int32(attr.Height)); err != nil {
			return fmt.Errorf("failed to create buffer: %w", err)
		}
	}

	// Attach buffer to surface
	if err := w.wlSurface.Attach(w.buffer, 0, 0); err != nil {
		return fmt.Errorf("failed to attach buffer: %w", err)
	}

	// Mark the entire surface as damaged
	if err := w.wlSurface.Damage(0, 0, int32(attr.Width), int32(attr.Height)); err != nil {
		return fmt.Errorf("failed to damage surface: %w", err)
	}

	// Commit the surface
	if err := w.wlSurface.Commit(); err != nil {
		return fmt.Errorf("failed to commit surface: %w", err)
	}

	w.needsRedraw = false
	log.Printf("Rendered frame for window: %s", w.BaseWindow.Title())
	return nil
}
