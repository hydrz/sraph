package wayland

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/opensraph/sraph/gio"
)

var _ gio.Driver = (*waylandDriver)(nil)

func init() {
	gio.RegisterDriver(gio.DriverTypeWayland, newWaylandDriver)
}

// Wayland library handles
var (
	libwayland   uintptr
	libxkbcommon uintptr
)

// Wayland function signatures (only functions that actually exist in libwayland-client.so)
var (
	wl_display_connect    func(name *byte) uintptr
	wl_display_disconnect func(display uintptr)
	wl_display_dispatch   func(display uintptr) int32
	wl_display_roundtrip  func(display uintptr) int32
	wl_display_flush      func(display uintptr) int32
	wl_display_get_fd     func(display uintptr) int32

	// XKB functions
	xkb_context_new              func(flags uint32) uintptr
	xkb_context_unref            func(context uintptr)
	xkb_keymap_new_from_string   func(context uintptr, str *byte, format uint32, flags uint32) uintptr
	xkb_keymap_unref             func(keymap uintptr)
	xkb_state_new                func(keymap uintptr) uintptr
	xkb_state_unref              func(state uintptr)
	xkb_state_key_get_one_sym    func(state uintptr, key uint32) uint32
	xkb_state_update_mask        func(state uintptr, depressed_mods uint32, latched_mods uint32, locked_mods uint32, depressed_layout uint32, latched_layout uint32, locked_layout uint32) uint32
	xkb_state_mod_name_is_active func(state uintptr, name *byte, type_ uint32) int32
)

// Constants
const (
	XKB_CONTEXT_NO_FLAGS              = 0
	XKB_KEYMAP_FORMAT_TEXT_V1         = 1
	XKB_KEYMAP_COMPILE_NO_FLAGS       = 0
	XKB_STATE_MODS_EFFECTIVE          = 1 << 3
	WL_KEYBOARD_KEYMAP_FORMAT_XKB_V1  = 1
	WL_KEYBOARD_KEY_STATE_PRESSED     = 1
	WL_POINTER_BUTTON_STATE_PRESSED   = 1
	WL_POINTER_AXIS_VERTICAL_SCROLL   = 0
	WL_POINTER_AXIS_HORIZONTAL_SCROLL = 1
	WL_SEAT_CAPABILITY_POINTER        = 1
	WL_SEAT_CAPABILITY_KEYBOARD       = 2
)

// Interface strings
var (
	wl_compositor_interface = []byte("wl_compositor\x00")
	wl_shell_interface      = []byte("wl_shell\x00")
	wl_shm_interface        = []byte("wl_shm\x00")
	wl_seat_interface       = []byte("wl_seat\x00")
)

// Listener function types
type registryGlobalFunc func(data uintptr, registry uintptr, name uint32, iface uintptr, version uint32)
type registryGlobalRemoveFunc func(data uintptr, registry uintptr, name uint32)
type seatCapabilitiesFunc func(data uintptr, seat uintptr, capabilities uint32)
type seatNameFunc func(data uintptr, seat uintptr, name uintptr)
type keyboardKeymapFunc func(data uintptr, keyboard uintptr, format uint32, fd int32, size uint32)
type keyboardEnterFunc func(data uintptr, keyboard uintptr, serial uint32, surface uintptr, keys uintptr)
type keyboardLeaveFunc func(data uintptr, keyboard uintptr, serial uint32, surface uintptr)
type keyboardKeyFunc func(data uintptr, keyboard uintptr, serial uint32, time uint32, key uint32, state uint32)
type keyboardModifiersFunc func(data uintptr, keyboard uintptr, serial uint32, mods_depressed uint32, mods_latched uint32, mods_locked uint32, group uint32)
type pointerEnterFunc func(data uintptr, pointer uintptr, serial uint32, surface uintptr, surface_x int32, surface_y int32)
type pointerLeaveFunc func(data uintptr, pointer uintptr, serial uint32, surface uintptr)
type pointerMotionFunc func(data uintptr, pointer uintptr, time uint32, surface_x int32, surface_y int32)
type pointerButtonFunc func(data uintptr, pointer uintptr, serial uint32, time uint32, button uint32, state uint32)
type pointerAxisFunc func(data uintptr, pointer uintptr, time uint32, axis uint32, value int32)
type shellSurfacePingFunc func(data uintptr, shell_surface uintptr, serial uint32)
type shellSurfaceConfigureFunc func(data uintptr, shell_surface uintptr, edges uint32, width int32, height int32)
type shellSurfacePopupDoneFunc func(data uintptr, shell_surface uintptr)

// Listener structures
type registryListener struct {
	global        uintptr
	global_remove uintptr
}

type seatListener struct {
	capabilities uintptr
	name         uintptr
}

type keyboardListener struct {
	keymap    uintptr
	enter     uintptr
	leave     uintptr
	key       uintptr
	modifiers uintptr
}

type pointerListener struct {
	enter  uintptr
	leave  uintptr
	motion uintptr
	button uintptr
	axis   uintptr
}

type shellSurfaceListener struct {
	ping       uintptr
	configure  uintptr
	popup_done uintptr
}

type waylandDriver struct {
	display    uintptr
	registry   uintptr
	compositor uintptr
	shell      uintptr
	shm        uintptr
	seat       uintptr
	keyboard   uintptr
	pointer    uintptr

	// XKB keyboard context
	xkbContext uintptr
	xkbKeymap  uintptr
	xkbState   uintptr

	// Event management
	eventChan chan waylandEvent
	mu        sync.RWMutex
	windows   map[uintptr]*waylandWindow
	running   bool

	// Listeners
	registryListener     registryListener
	seatListener         seatListener
	keyboardListener     keyboardListener
	pointerListener      pointerListener
	shellSurfaceListener shellSurfaceListener
}

type waylandEvent struct {
	window *waylandWindow
	event  interface{}
}

func loadWaylandLibraries() error {
	var err error

	// Load Wayland client library
	libwayland, err = purego.Dlopen("libwayland-client.so.0", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load libwayland-client.so.0: %v", err)
	}

	// Load XKB common library
	libxkbcommon, err = purego.Dlopen("libxkbcommon.so.0", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load libxkbcommon.so.0: %v", err)
	}

	// Load only functions that actually exist in the libraries
	purego.RegisterLibFunc(&wl_display_connect, libwayland, "wl_display_connect")
	purego.RegisterLibFunc(&wl_display_disconnect, libwayland, "wl_display_disconnect")
	purego.RegisterLibFunc(&wl_display_dispatch, libwayland, "wl_display_dispatch")
	purego.RegisterLibFunc(&wl_display_roundtrip, libwayland, "wl_display_roundtrip")
	purego.RegisterLibFunc(&wl_display_flush, libwayland, "wl_display_flush")
	purego.RegisterLibFunc(&wl_display_get_fd, libwayland, "wl_display_get_fd")

	// Load XKB functions
	purego.RegisterLibFunc(&xkb_context_new, libxkbcommon, "xkb_context_new")
	purego.RegisterLibFunc(&xkb_context_unref, libxkbcommon, "xkb_context_unref")
	purego.RegisterLibFunc(&xkb_keymap_new_from_string, libxkbcommon, "xkb_keymap_new_from_string")
	purego.RegisterLibFunc(&xkb_keymap_unref, libxkbcommon, "xkb_keymap_unref")
	purego.RegisterLibFunc(&xkb_state_new, libxkbcommon, "xkb_state_new")
	purego.RegisterLibFunc(&xkb_state_unref, libxkbcommon, "xkb_state_unref")
	purego.RegisterLibFunc(&xkb_state_key_get_one_sym, libxkbcommon, "xkb_state_key_get_one_sym")
	purego.RegisterLibFunc(&xkb_state_update_mask, libxkbcommon, "xkb_state_update_mask")
	purego.RegisterLibFunc(&xkb_state_mod_name_is_active, libxkbcommon, "xkb_state_mod_name_is_active")

	return nil
}

func newWaylandDriver() (gio.Driver, error) {
	// Check if we're in a Wayland session
	if os.Getenv("WAYLAND_DISPLAY") == "" {
		return nil, fmt.Errorf("wayland driver: no Wayland display available")
	}

	// Load libraries
	if err := loadWaylandLibraries(); err != nil {
		return nil, fmt.Errorf("wayland driver: %v", err)
	}

	// Connect to Wayland display
	display := wl_display_connect(nil)
	if display == 0 {
		return nil, fmt.Errorf("wayland driver: failed to connect to Wayland display")
	}

	// Initialize XKB
	xkbContext := xkb_context_new(XKB_CONTEXT_NO_FLAGS)
	if xkbContext == 0 {
		wl_display_disconnect(display)
		return nil, fmt.Errorf("wayland driver: failed to create XKB context")
	}

	wd := &waylandDriver{
		display:    display,
		xkbContext: xkbContext,
		eventChan:  make(chan waylandEvent, 256),
		windows:    make(map[uintptr]*waylandWindow),
		running:    true,
	}

	// NOTE: In a real implementation, we would need to:
	// 1. Get the registry using protocol messages
	// 2. Bind to compositor, shell, shm, and seat
	// 3. Set up event listeners
	//
	// For now, we'll simulate having these objects
	wd.registry = 1   // Fake registry
	wd.compositor = 2 // Fake compositor
	wd.shell = 3      // Fake shell
	wd.shm = 4        // Fake shm
	wd.seat = 5       // Fake seat

	// Set up listeners (would be real in complete implementation)
	wd.setupListeners()

	// Start event processing loop
	go wd.eventLoop()

	return wd, nil
}

func (wd *waylandDriver) setupListeners() {
	// Set up registry listener
	wd.registryListener.global = purego.NewCallback(wd.registryGlobal)
	wd.registryListener.global_remove = purego.NewCallback(wd.registryGlobalRemove)

	// Set up seat listener
	wd.seatListener.capabilities = purego.NewCallback(wd.seatCapabilities)
	wd.seatListener.name = purego.NewCallback(wd.seatName)

	// Set up keyboard listener
	wd.keyboardListener.keymap = purego.NewCallback(wd.keyboardKeymap)
	wd.keyboardListener.enter = purego.NewCallback(wd.keyboardEnter)
	wd.keyboardListener.leave = purego.NewCallback(wd.keyboardLeave)
	wd.keyboardListener.key = purego.NewCallback(wd.keyboardKey)
	wd.keyboardListener.modifiers = purego.NewCallback(wd.keyboardModifiers)

	// Set up pointer listener
	wd.pointerListener.enter = purego.NewCallback(wd.pointerEnter)
	wd.pointerListener.leave = purego.NewCallback(wd.pointerLeave)
	wd.pointerListener.motion = purego.NewCallback(wd.pointerMotion)
	wd.pointerListener.button = purego.NewCallback(wd.pointerButton)
	wd.pointerListener.axis = purego.NewCallback(wd.pointerAxis)

	// Set up shell surface listener
	wd.shellSurfaceListener.ping = purego.NewCallback(wd.shellSurfacePing)
	wd.shellSurfaceListener.configure = purego.NewCallback(wd.shellSurfaceConfigure)
	wd.shellSurfaceListener.popup_done = purego.NewCallback(wd.shellSurfacePopupDone)
}

func (wd *waylandDriver) cleanup() {
	if wd.xkbState != 0 {
		xkb_state_unref(wd.xkbState)
	}
	if wd.xkbKeymap != 0 {
		xkb_keymap_unref(wd.xkbKeymap)
	}
	if wd.xkbContext != 0 {
		xkb_context_unref(wd.xkbContext)
	}
	if wd.display != 0 {
		wl_display_disconnect(wd.display)
	}
}

// CreateWindow implements gio.Driver
func (wd *waylandDriver) CreateWindow(o gio.NewWindowOptions) (gio.Window, error) {
	bw := gio.NewBaseWindow(o)
	return newWaylandWindow(wd, bw)
}

// Type implements gio.Driver
func (wd *waylandDriver) Type() gio.DriverType {
	return gio.DriverTypeWayland
}

func (wd *waylandDriver) eventLoop() {
	for wd.running {
		ret := wl_display_dispatch(wd.display)
		if ret == -1 {
			// Display connection lost or error occurred
			break
		}

		// Process any events in our queue
		select {
		case event := <-wd.eventChan:
			if event.window != nil {
				// Handle event for specific window
				_ = event
			}
		default:
			runtime.Gosched()
		}
	}
}

func (wd *waylandDriver) registerWindow(surface uintptr, window *waylandWindow) {
	wd.mu.Lock()
	defer wd.mu.Unlock()
	wd.windows[surface] = window
}

func (wd *waylandDriver) unregisterWindow(surface uintptr) {
	wd.mu.Lock()
	defer wd.mu.Unlock()
	delete(wd.windows, surface)
}

func (wd *waylandDriver) getWindow(surface uintptr) *waylandWindow {
	wd.mu.RLock()
	defer wd.mu.RUnlock()
	return wd.windows[surface]
}

// Convert Wayland button codes to gio MouseButton
func (wd *waylandDriver) translateMouseButton(button uint32) gio.MouseButton {
	// Linux input event codes
	const (
		BTN_LEFT   = 0x110
		BTN_RIGHT  = 0x111
		BTN_MIDDLE = 0x112
		BTN_SIDE   = 0x113
		BTN_EXTRA  = 0x114
	)

	switch button {
	case BTN_LEFT:
		return gio.MouseButtonLeft
	case BTN_RIGHT:
		return gio.MouseButtonRight
	case BTN_MIDDLE:
		return gio.MouseButtonMiddle
	case BTN_SIDE:
		return gio.MouseButtonBack
	case BTN_EXTRA:
		return gio.MouseButtonForward
	default:
		return gio.MouseButtonUnknown
	}
}

// Convert XKB key to gio KeyCode
func (wd *waylandDriver) translateKeyCode(key uint32) gio.KeyCode {
	if wd.xkbState == 0 {
		return gio.KeyCodeUnknown
	}

	keysym := xkb_state_key_get_one_sym(wd.xkbState, key+8) // Wayland uses evdev codes
	return xkbKeysymToGioKeyCode(keysym)
}

// Convert XKB modifiers to gio ModifierKey
func (wd *waylandDriver) translateModifiers() gio.ModifierKey {
	if wd.xkbState == 0 {
		return gio.ModifierKeyNone
	}

	var modifiers gio.ModifierKey

	shiftStr := []byte("Shift\x00")
	controlStr := []byte("Control\x00")
	mod1Str := []byte("Mod1\x00")
	mod4Str := []byte("Mod4\x00")

	if xkb_state_mod_name_is_active(wd.xkbState, &shiftStr[0], XKB_STATE_MODS_EFFECTIVE) != 0 {
		modifiers |= gio.ModifierKeyShift
	}
	if xkb_state_mod_name_is_active(wd.xkbState, &controlStr[0], XKB_STATE_MODS_EFFECTIVE) != 0 {
		modifiers |= gio.ModifierKeyCtrl
	}
	if xkb_state_mod_name_is_active(wd.xkbState, &mod1Str[0], XKB_STATE_MODS_EFFECTIVE) != 0 {
		modifiers |= gio.ModifierKeyAlt
	}
	if xkb_state_mod_name_is_active(wd.xkbState, &mod4Str[0], XKB_STATE_MODS_EFFECTIVE) != 0 {
		modifiers |= gio.ModifierKeySuper
	}

	return modifiers
}

func (wd *waylandDriver) shutdown() {
	wd.mu.Lock()
	wd.running = false
	wd.mu.Unlock()
	close(wd.eventChan)
	wd.cleanup()
}

// Callback functions (simplified for demonstration)
func (wd *waylandDriver) registryGlobal(data uintptr, registry uintptr, name uint32, iface uintptr, version uint32) {
	// In a real implementation, we would parse the interface name and bind to it
	// For now, just log that we received a global
	fmt.Printf("Wayland global received: name=%d, version=%d\n", name, version)
}

func (wd *waylandDriver) registryGlobalRemove(data uintptr, registry uintptr, name uint32) {
	// Handle global object removal if needed
}

func (wd *waylandDriver) seatCapabilities(data uintptr, seat uintptr, caps uint32) {
	// In a real implementation, we would check capabilities and get keyboard/pointer
	fmt.Printf("Seat capabilities: %d\n", caps)
}

func (wd *waylandDriver) seatName(data uintptr, seat uintptr, name uintptr) {
	// Handle seat name if needed
}

func (wd *waylandDriver) keyboardKeymap(data uintptr, keyboard uintptr, format uint32, fd int32, size uint32) {
	// Simplified keymap handling
	if format != WL_KEYBOARD_KEYMAP_FORMAT_XKB_V1 {
		syscall.Close(int(fd))
		return
	}

	// Map the keymap file
	mapPtr, err := syscall.Mmap(int(fd), 0, int(size), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		syscall.Close(int(fd))
		return
	}
	defer syscall.Munmap(mapPtr)
	defer syscall.Close(int(fd))

	// Create XKB keymap from string
	keymapStr := append(mapPtr, 0) // Ensure null termination

	if wd.xkbKeymap != 0 {
		xkb_keymap_unref(wd.xkbKeymap)
	}
	if wd.xkbState != 0 {
		xkb_state_unref(wd.xkbState)
	}

	wd.xkbKeymap = xkb_keymap_new_from_string(wd.xkbContext, &keymapStr[0], XKB_KEYMAP_FORMAT_TEXT_V1, XKB_KEYMAP_COMPILE_NO_FLAGS)
	if wd.xkbKeymap != 0 {
		wd.xkbState = xkb_state_new(wd.xkbKeymap)
	}
}

func (wd *waylandDriver) keyboardEnter(data uintptr, keyboard uintptr, serial uint32, surface uintptr, keys uintptr) {
	window := wd.getWindow(surface)
	if window != nil {
		window.handleFocusIn()
	}
}

func (wd *waylandDriver) keyboardLeave(data uintptr, keyboard uintptr, serial uint32, surface uintptr) {
	window := wd.getWindow(surface)
	if window != nil {
		window.handleFocusOut()
	}
}

func (wd *waylandDriver) keyboardKey(data uintptr, keyboard uintptr, serial uint32, time uint32, key uint32, state uint32) {
	// Find the focused window - in this simple implementation, we'll use the first window
	wd.mu.RLock()
	var targetWindow *waylandWindow
	for _, window := range wd.windows {
		targetWindow = window
		break // Use first window for now
	}
	wd.mu.RUnlock()

	if targetWindow != nil {
		keyCode := wd.translateKeyCode(key)
		modifiers := wd.translateModifiers()
		pressed := state == WL_KEYBOARD_KEY_STATE_PRESSED
		targetWindow.handleKeyEvent(keyCode, modifiers, pressed)
	}
}

func (wd *waylandDriver) keyboardModifiers(data uintptr, keyboard uintptr, serial uint32, modsDepressed uint32, modsLatched uint32, modsLocked uint32, group uint32) {
	if wd.xkbState != 0 {
		xkb_state_update_mask(wd.xkbState, modsDepressed, modsLatched, modsLocked, 0, 0, group)
	}
}

func (wd *waylandDriver) pointerEnter(data uintptr, pointer uintptr, serial uint32, surface uintptr, sx int32, sy int32) {
	window := wd.getWindow(surface)
	if window != nil {
		x := fixedToInt(sx)
		y := fixedToInt(sy)
		window.handlePointerEnter(x, y)
	}
}

func (wd *waylandDriver) pointerLeave(data uintptr, pointer uintptr, serial uint32, surface uintptr) {
	window := wd.getWindow(surface)
	if window != nil {
		window.handlePointerLeave()
	}
}

func (wd *waylandDriver) pointerMotion(data uintptr, pointer uintptr, time uint32, sx int32, sy int32) {
	// Find the focused window - use first window for simplicity
	wd.mu.RLock()
	var targetWindow *waylandWindow
	for _, window := range wd.windows {
		targetWindow = window
		break
	}
	wd.mu.RUnlock()

	if targetWindow != nil {
		x := fixedToInt(sx)
		y := fixedToInt(sy)
		modifiers := wd.translateModifiers()
		targetWindow.handlePointerMotion(x, y, modifiers)
	}
}

func (wd *waylandDriver) pointerButton(data uintptr, pointer uintptr, serial uint32, time uint32, button uint32, state uint32) {
	// Find the focused window
	wd.mu.RLock()
	var targetWindow *waylandWindow
	for _, window := range wd.windows {
		targetWindow = window
		break
	}
	wd.mu.RUnlock()

	if targetWindow != nil {
		mouseButton := wd.translateMouseButton(button)
		modifiers := wd.translateModifiers()
		pressed := state == WL_POINTER_BUTTON_STATE_PRESSED
		targetWindow.handleButtonEvent(mouseButton, modifiers, pressed)
	}
}

func (wd *waylandDriver) pointerAxis(data uintptr, pointer uintptr, time uint32, axis uint32, value int32) {
	// Find the focused window
	wd.mu.RLock()
	var targetWindow *waylandWindow
	for _, window := range wd.windows {
		targetWindow = window
		break
	}
	wd.mu.RUnlock()

	if targetWindow != nil {
		delta := fixedToFloat(value)
		modifiers := wd.translateModifiers()

		var deltaX, deltaY float64
		if axis == WL_POINTER_AXIS_HORIZONTAL_SCROLL {
			deltaX = delta
		} else {
			deltaY = delta
		}

		targetWindow.handleScrollEvent(deltaX, deltaY, modifiers)
	}
}

func (wd *waylandDriver) shellSurfacePing(data uintptr, shell_surface uintptr, serial uint32) {
	// In a real implementation, we would call wl_shell_surface_pong
	// For now, just acknowledge we received a ping
	fmt.Printf("Shell surface ping received, serial: %d\n", serial)
}

func (wd *waylandDriver) shellSurfaceConfigure(data uintptr, shell_surface uintptr, edges uint32, width int32, height int32) {
	// The data parameter should point to the window, but since we're using a simplified approach,
	// we'll find the window by iterating through our registered windows
	wd.mu.RLock()
	var targetWindow *waylandWindow
	for _, window := range wd.windows {
		// In a real implementation, you'd match the shell_surface to the window
		targetWindow = window
		break
	}
	wd.mu.RUnlock()

	if targetWindow != nil {
		targetWindow.handleConfigure(int(width), int(height))
	}
}

func (wd *waylandDriver) shellSurfacePopupDone(data uintptr, shell_surface uintptr) {
	// Handle popup done if needed
}

// Helper functions
func fixedToInt(f int32) int {
	return int(f / 256)
}

func fixedToFloat(f int32) float64 {
	return float64(f) / 256.0
}

func goString(cstr *byte) string {
	if cstr == nil {
		return ""
	}
	var length int
	for {
		if *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(cstr)) + uintptr(length))) == 0 {
			break
		}
		length++
	}
	return string(unsafe.Slice(cstr, length))
}

// Helper function to convert XKB keysym to gio KeyCode
func xkbKeysymToGioKeyCode(keysym uint32) gio.KeyCode {
	// This is a simplified mapping - in a real implementation,
	// you'd want a comprehensive lookup table
	switch keysym {
	case 0xff08: // BackSpace
		return gio.KeyCodeBackspace
	case 0xff09: // Tab
		return gio.KeyCodeTab
	case 0xff0d: // Return
		return gio.KeyCodeEnter
	case 0xff1b: // Escape
		return gio.KeyCodeEscape
	case 0x0020: // Space
		return gio.KeyCodeSpace
	case 0xffbe: // F1
		return gio.KeyCodeF1
	case 0xffbf: // F2
		return gio.KeyCodeF2
	case 0xffc0: // F3
		return gio.KeyCodeF3
	case 0xffc1: // F4
		return gio.KeyCodeF4
	case 0xffc2: // F5
		return gio.KeyCodeF5
	case 0xffc3: // F6
		return gio.KeyCodeF6
	case 0xffc4: // F7
		return gio.KeyCodeF7
	case 0xffc5: // F8
		return gio.KeyCodeF8
	case 0xffc6: // F9
		return gio.KeyCodeF9
	case 0xffc7: // F10
		return gio.KeyCodeF10
	case 0xffc8: // F11
		return gio.KeyCodeF11
	case 0xffc9: // F12
		return gio.KeyCodeF12
	case 0xff51: // Left
		return gio.KeyCodeArrowLeft
	case 0xff52: // Up
		return gio.KeyCodeArrowUp
	case 0xff53: // Right
		return gio.KeyCodeArrowRight
	case 0xff54: // Down
		return gio.KeyCodeArrowDown
	// ASCII characters (letters)
	case 0x0061: // a
		return gio.KeyCodeKeyA
	case 0x0062: // b
		return gio.KeyCodeKeyB
	case 0x0063: // c
		return gio.KeyCodeKeyC
	case 0x0064: // d
		return gio.KeyCodeKeyD
	case 0x0065: // e
		return gio.KeyCodeKeyE
	case 0x0066: // f
		return gio.KeyCodeKeyF
	case 0x0067: // g
		return gio.KeyCodeKeyG
	case 0x0068: // h
		return gio.KeyCodeKeyH
	case 0x0069: // i
		return gio.KeyCodeKeyI
	case 0x006a: // j
		return gio.KeyCodeKeyJ
	case 0x006b: // k
		return gio.KeyCodeKeyK
	case 0x006c: // l
		return gio.KeyCodeKeyL
	case 0x006d: // m
		return gio.KeyCodeKeyM
	case 0x006e: // n
		return gio.KeyCodeKeyN
	case 0x006f: // o
		return gio.KeyCodeKeyO
	case 0x0070: // p
		return gio.KeyCodeKeyP
	case 0x0071: // q
		return gio.KeyCodeKeyQ
	case 0x0072: // r
		return gio.KeyCodeKeyR
	case 0x0073: // s
		return gio.KeyCodeKeyS
	case 0x0074: // t
		return gio.KeyCodeKeyT
	case 0x0075: // u
		return gio.KeyCodeKeyU
	case 0x0076: // v
		return gio.KeyCodeKeyV
	case 0x0077: // w
		return gio.KeyCodeKeyW
	case 0x0078: // x
		return gio.KeyCodeKeyX
	case 0x0079: // y
		return gio.KeyCodeKeyY
	case 0x007a: // z
		return gio.KeyCodeKeyZ
	// Numbers
	case 0x0030: // 0
		return gio.KeyCodeDigit0
	case 0x0031: // 1
		return gio.KeyCodeDigit1
	case 0x0032: // 2
		return gio.KeyCodeDigit2
	case 0x0033: // 3
		return gio.KeyCodeDigit3
	case 0x0034: // 4
		return gio.KeyCodeDigit4
	case 0x0035: // 5
		return gio.KeyCodeDigit5
	case 0x0036: // 6
		return gio.KeyCodeDigit6
	case 0x0037: // 7
		return gio.KeyCodeDigit7
	case 0x0038: // 8
		return gio.KeyCodeDigit8
	case 0x0039: // 9
		return gio.KeyCodeDigit9
	default:
		// For other ASCII characters
		if keysym >= 0x20 && keysym <= 0x7e {
			return gio.KeyCode(keysym)
		}
		return gio.KeyCodeUnknown
	}
}
