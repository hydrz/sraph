package wayland

import (
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
	"unsafe"

	"github.com/opensraph/sraph/gio"
	"github.com/rajveermalviya/go-wayland/wayland/client"
	xdg_shell "github.com/rajveermalviya/go-wayland/wayland/stable/xdg-shell"
)

var _ gio.Window = (*WaylandWindow)(nil)

// WaylandWindow represents a Wayland window implementation
type WaylandWindow struct {
	gio.BaseWindow

	driver      *WaylandDriver
	surface     *client.Surface
	xdgSurface  *xdg_shell.Surface
	xdgTopLevel *xdg_shell.Toplevel

	// Buffer management
	currentBuffer *windowBuffer
	buffers       []*windowBuffer
	shmPool       *client.ShmPool

	// Window state tracking
	configured  bool
	needsRedraw bool
	mapped      bool

	// Clipboard support
	clipboardData *gio.DataTransfer

	mu sync.RWMutex
}

type windowBuffer struct {
	buffer *client.Buffer
	data   []byte
	width  int32
	height int32
	stride int32
	busy   bool
	fd     int // Store file descriptor for cleanup
}

// WindowID returns the Wayland window ID
func (w *WaylandWindow) WindowID() gio.WindowID {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return gio.WindowID(uintptr(unsafe.Pointer(w.surface)))
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
		surface.Destroy()
		return nil, fmt.Errorf("wayland: failed to create XDG surface: %w", err)
	}

	// Create toplevel
	toplevel, err := xdgSurface.GetToplevel()
	if err != nil {
		xdgSurface.Destroy()
		surface.Destroy()
		return nil, fmt.Errorf("wayland: failed to create toplevel: %w", err)
	}

	// Create window wrapper
	ww := &WaylandWindow{
		BaseWindow:    bw,
		driver:        driver,
		surface:       surface,
		xdgSurface:    xdgSurface,
		xdgTopLevel:   toplevel,
		needsRedraw:   true,
		clipboardData: gio.NewDataTransfer(),
	}

	// Configure window properties
	if err := ww.init(); err != nil {
		ww.destroy()
		return nil, fmt.Errorf("wayland: failed to initialize window: %w", err)
	}

	return ww, nil
}

// init sets up the window with initial properties
func (w *WaylandWindow) init() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	attr := w.BaseWindow.Attr()

	// Set window title
	if err := w.xdgTopLevel.SetTitle(attr.Title); err != nil {
		return fmt.Errorf("failed to set title: %w", err)
	}

	// Set window app ID
	if err := w.xdgTopLevel.SetAppId(attr.Title); err != nil {
		return fmt.Errorf("failed to set app ID: %w", err)
	}

	// Set minimum size if specified
	if attr.Width > 0 && attr.Height > 0 {
		if err := w.xdgTopLevel.SetMinSize(int32(attr.Width), int32(attr.Height)); err != nil {
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
	w.xdgTopLevel.SetCloseHandler(func(event xdg_shell.ToplevelCloseEvent) {
		w.Close()
	})

	// Handle window configuration changes
	w.xdgTopLevel.SetConfigureHandler(func(event xdg_shell.ToplevelConfigureEvent) {
		w.handleToplevelConfigure(event)
	})

	// Handle XDG surface configuration
	w.xdgSurface.SetConfigureHandler(func(event xdg_shell.SurfaceConfigureEvent) {
		w.handleSurfaceConfigure(event)
	})
}

// handleToplevelConfigure handles toplevel configuration changes
func (w *WaylandWindow) handleToplevelConfigure(event xdg_shell.ToplevelConfigureEvent) {
	w.mu.Lock()
	defer w.mu.Unlock()

	width := event.Width
	height := event.Height

	// Update window attributes if size changed
	if width > 0 && height > 0 {
		attr := w.BaseWindow.Attr()
		if attr.Width != int(width) || attr.Height != int(height) {
			attr.Width = int(width)
			attr.Height = int(height)
			w.BaseWindow.SetAttr(attr)
			w.needsRedraw = true
		}
	}

	// Update window state based on states array
	attr := w.BaseWindow.Attr()
	newState := attr.State

	// Reset certain states
	newState &= ^(gio.WindowStateMaximized | gio.WindowStateFocused)

	// Parse states from event
	for i := 0; i < len(event.States); i += 4 {
		if i+3 < len(event.States) {
			state := uint32(event.States[i]) |
				uint32(event.States[i+1])<<8 |
				uint32(event.States[i+2])<<16 |
				uint32(event.States[i+3])<<24

			switch state {
			case uint32(xdg_shell.ToplevelStateMaximized):
				newState |= gio.WindowStateMaximized
			case uint32(xdg_shell.ToplevelStateFullscreen):
				// Handle fullscreen if needed
			case uint32(xdg_shell.ToplevelStateActivated):
				newState |= gio.WindowStateFocused
			case uint32(xdg_shell.ToplevelStateResizing):
				// Handle resizing state
			}
		}
	}

	if newState != attr.State {
		attr.State = newState
		w.BaseWindow.SetAttr(attr)
	}

	w.configured = true
}

// handleSurfaceConfigure handles XDG surface configuration
func (w *WaylandWindow) handleSurfaceConfigure(event xdg_shell.SurfaceConfigureEvent) {
	// Acknowledge the configure event
	w.xdgSurface.AckConfigure(event.Serial)

	// Trigger render if needed
	if w.configured && w.needsRedraw {
		go func() {
			if err := w.render(); err != nil {
				log.Printf("wayland: render error: %v", err)
			}
		}()
	}
}

// render performs window rendering
func (w *WaylandWindow) render() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.configured || !w.needsRedraw {
		return nil
	}

	attr := w.BaseWindow.Attr()
	width := int32(attr.Width)
	height := int32(attr.Height)

	if width <= 0 || height <= 0 {
		return nil
	}

	// Get or create buffer
	buffer, err := w.getBuffer(width, height)
	if err != nil {
		return fmt.Errorf("failed to get buffer: %w", err)
	}

	// Fill buffer with content
	w.drawFrame(buffer, int(width), int(height))

	// Attach buffer to surface
	if err := w.surface.Attach(buffer.buffer, 0, 0); err != nil {
		return fmt.Errorf("failed to attach buffer: %w", err)
	}

	// Mark entire surface as damaged
	if err := w.surface.Damage(0, 0, width, height); err != nil {
	}

	// Set buffer as busy
	buffer.busy = true

	// Commit surface
	if err := w.surface.Commit(); err != nil {
		return fmt.Errorf("failed to commit surface: %w", err)
	}

	w.currentBuffer = buffer
	w.needsRedraw = false
	w.mapped = true

	return nil
}

// getBuffer gets or creates a buffer for the given size
func (w *WaylandWindow) getBuffer(width, height int32) (*windowBuffer, error) {
	stride := width * 4
	size := stride * height

	// Try to reuse an existing buffer
	for _, buf := range w.buffers {
		if !buf.busy && buf.width == width && buf.height == height {
			return buf, nil
		}
	}

	// Create new buffer
	buffer, err := w.createBuffer(width, height, stride, size)
	if err != nil {
		return nil, err
	}

	// Set release handler
	buffer.buffer.SetReleaseHandler(func(event client.BufferReleaseEvent) {
		buffer.busy = false
	})

	w.buffers = append(w.buffers, buffer)
	return buffer, nil
}

// createBuffer creates a new shared memory buffer
func (w *WaylandWindow) createBuffer(width, height, stride, size int32) (*windowBuffer, error) {
	// Create anonymous file
	fd, err := w.createAnonymousFile(int(size))
	if err != nil {
		return nil, fmt.Errorf("failed to create anonymous file: %w", err)
	}
	defer syscall.Close(fd)

	// Map memory
	data, err := syscall.Mmap(fd, 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("failed to mmap: %w", err)
	}

	// Create or reuse SHM pool
	if w.shmPool == nil {
		pool, err := w.driver.shm.CreatePool(fd, size)
		if err != nil {
			syscall.Munmap(data)
			return nil, fmt.Errorf("failed to create shm pool: %w", err)
		}
		w.shmPool = pool
	}

	// Create buffer from pool
	buffer, err := w.shmPool.CreateBuffer(0, width, height, stride, uint32(client.ShmFormatArgb8888))
	if err != nil {
		syscall.Munmap(data)
		return nil, fmt.Errorf("failed to create buffer: %w", err)
	}

	return &windowBuffer{
		buffer: buffer,
		data:   data,
		width:  width,
		height: height,
		stride: stride,
		busy:   false,
	}, nil
}

// createAnonymousFile creates an anonymous file for shared memory
func (w *WaylandWindow) createAnonymousFile(size int) (int, error) {
	// Try to use memfd_create (Linux-specific)
	name := "wayland-buffer"
	fd, _, err := syscall.RawSyscall(319, uintptr(unsafe.Pointer(&[]byte(name)[0])), 0, 0)
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

// drawFrame fills the buffer with content
func (w *WaylandWindow) drawFrame(buffer *windowBuffer, width, height int) {
	data := buffer.data
	stride := int(buffer.stride)

	// Create a simple gradient pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := y*stride + x*4

			// Create gradient from blue to red
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(128)
			a := uint8(255)

			// ARGB8888 format (little endian)
			data[offset+0] = b // Blue
			data[offset+1] = g // Green
			data[offset+2] = r // Red
			data[offset+3] = a // Alpha
		}
	}
}

// SetAttr updates window attributes
func (w *WaylandWindow) SetAttr(attr gio.WindowAttr) {
	w.BaseWindow.SetAttr(attr)

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.xdgTopLevel == nil {
		return
	}

	currentAttr := w.BaseWindow.Attr()

	// Update title if changed
	if attr.Title != "" && attr.Title != currentAttr.Title {
		w.xdgTopLevel.SetTitle(attr.Title)
	}

	// Handle state changes
	if attr.State.Contains(gio.WindowStateMaximized) && !currentAttr.State.Contains(gio.WindowStateMaximized) {
		w.xdgTopLevel.SetMaximized()
	} else if !attr.State.Contains(gio.WindowStateMaximized) && currentAttr.State.Contains(gio.WindowStateMaximized) {
		w.xdgTopLevel.UnsetMaximized()
	}

	// Trigger redraw if size changed
	if attr.Width != currentAttr.Width || attr.Height != currentAttr.Height {
		w.needsRedraw = true
	}

	// Commit changes
	if w.surface != nil {
		w.surface.Commit()
	}
}

// Close closes the window and cleans up resources
func (w *WaylandWindow) Close() error {
	if err := w.BaseWindow.Close(); err != nil {
		return err
	}

	w.destroy()

	// Remove from driver's window list
	w.driver.removeWindow(w.WindowID())

	return nil
}

// destroy cleans up all Wayland resources
func (w *WaylandWindow) destroy() {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Clean up buffers
	for _, buffer := range w.buffers {
		if buffer.buffer != nil {
			buffer.buffer.Destroy()
		}
		if len(buffer.data) > 0 {
			syscall.Munmap(buffer.data)
		}
	}
	w.buffers = nil
	w.currentBuffer = nil

	// Clean up SHM pool
	if w.shmPool != nil {
		w.shmPool.Destroy()
		w.shmPool = nil
	}

	// Clean up XDG objects
	if w.xdgTopLevel != nil {
		w.xdgTopLevel.Destroy()
		w.xdgTopLevel = nil
	}
	if w.xdgSurface != nil {
		w.xdgSurface.Destroy()
		w.xdgSurface = nil
	}
	if w.surface != nil {
		w.surface.Destroy()
		w.surface = nil
	}
}

// Show makes the window visible
func (w *WaylandWindow) Show() error {
	if err := w.BaseWindow.Show(); err != nil {
		return err
	}

	w.mu.Lock()
	w.needsRedraw = true
	w.mu.Unlock()

	// Trigger render if configured
	if w.configured {
		go func() {
			if err := w.render(); err != nil {
				log.Printf("wayland: render error: %v", err)
			}
		}()
	}

	return nil
}
