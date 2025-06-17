package cocoa

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/opensraph/sraph/gio"
)

var _ gio.Driver = (*cocoaDriver)(nil)

func init() {
	runtime.LockOSThread()
	gio.RegisterDriver(gio.DriverTypeCocoa, newCocoaDriver)
}

// Cocoa framework handles
var (
	libFoundation   uintptr
	libAppKit       uintptr
	libCoreGraphics uintptr
)

// Objective-C runtime functions
var (
	objc_getClass              func(name *byte) uintptr
	objc_msgSend               func(receiver uintptr, selector uintptr, args ...uintptr) uintptr
	objc_msgSend_fpret         func(receiver uintptr, selector uintptr, args ...uintptr) float64
	objc_msgSend_stret         func(ret unsafe.Pointer, receiver uintptr, selector uintptr, args ...uintptr)
	sel_registerName           func(name *byte) uintptr
	class_createInstance       func(cls uintptr, extraBytes uintptr) uintptr
	objc_allocateClassPair     func(superclass uintptr, name *byte, extraBytes uintptr) uintptr
	objc_registerClassPair     func(cls uintptr)
	class_addMethod            func(cls uintptr, name uintptr, imp uintptr, types *byte) bool
	class_addIvar              func(cls uintptr, name *byte, size uintptr, alignment uint8, types *byte) bool
	object_getInstanceVariable func(obj uintptr, name *byte, outValue *uintptr) uintptr
	object_setInstanceVariable func(obj uintptr, name *byte, value uintptr) uintptr
)

// Core Graphics functions
var (
	CGMainDisplayID            func() uint32
	CGDisplayBounds            func(display uint32) CGRect
	CGEventCreateKeyboardEvent func(source uintptr, virtualKey uint16, keyDown bool) uintptr
	CGEventPost                func(tap uint32, event uintptr)
	CGEventCreateMouseEvent    func(source uintptr, mouseType uint32, mouseCursorPosition CGPoint, mouseButton uint32) uintptr
)

// Foundation/AppKit functions (using objc_msgSend)
var (
	NSApp                  uintptr // [NSApplication sharedApplication]
	NSApplicationClass     uintptr
	NSWindowClass          uintptr
	NSViewClass            uintptr
	NSEventClass           uintptr
	NSStringClass          uintptr
	NSAutoreleasePoolClass uintptr
)

// Selectors
var (
	sel_sharedApplication         uintptr
	sel_setActivationPolicy       uintptr
	sel_activateIgnoringOtherApps uintptr
	sel_run                       uintptr
	sel_stop                      uintptr
	sel_nextEventMatchingMask     uintptr
	sel_sendEvent                 uintptr
	sel_updateWindows             uintptr

	sel_alloc                    uintptr
	sel_init                     uintptr
	sel_initWithContentRect      uintptr
	sel_setTitle                 uintptr
	sel_makeKeyAndOrderFront     uintptr
	sel_orderOut                 uintptr
	sel_close                    uintptr
	sel_setDelegate              uintptr
	sel_frame                    uintptr
	sel_setFrame                 uintptr
	sel_center                   uintptr
	sel_contentView              uintptr
	sel_setContentSize           uintptr
	sel_setFrameOrigin           uintptr
	sel_makeKeyWindow            uintptr
	sel_zoom                     uintptr
	sel_setAcceptsFirstResponder uintptr

	// Window delegate selectors
	sel_windowShouldClose  uintptr
	sel_windowDidResize    uintptr
	sel_windowDidBecomeKey uintptr
	sel_windowDidResignKey uintptr

	// Event selectors
	sel_type             uintptr
	sel_locationInWindow uintptr
	sel_keyCode          uintptr
	sel_modifierFlags    uintptr
	sel_isARepeat        uintptr
	sel_scrollingDeltaX  uintptr
	sel_scrollingDeltaY  uintptr

	sel_stringWithUTF8String uintptr
	sel_autorelease          uintptr
	sel_release              uintptr
	sel_retain               uintptr
)

// Cocoa structures
type CGPoint struct {
	X, Y float64
}

type CGSize struct {
	Width, Height float64
}

type CGRect struct {
	Origin CGPoint
	Size   CGSize
}

type NSRect CGRect
type NSPoint CGPoint
type NSSize CGSize

// Constants
const (
	NSApplicationActivationPolicyRegular    = 0
	NSApplicationActivationPolicyAccessory  = 1
	NSApplicationActivationPolicyProhibited = 2

	NSEventMaskAny            = ^uint64(0)
	NSEventMaskKeyDown        = 1 << 10
	NSEventMaskKeyUp          = 1 << 11
	NSEventMaskMouseMoved     = 1 << 5
	NSEventMaskLeftMouseDown  = 1 << 1
	NSEventMaskLeftMouseUp    = 1 << 2
	NSEventMaskRightMouseDown = 1 << 3
	NSEventMaskRightMouseUp   = 1 << 4
	NSEventMaskScrollWheel    = 1 << 22

	NSWindowStyleMaskTitled         = 1 << 0
	NSWindowStyleMaskClosable       = 1 << 1
	NSWindowStyleMaskMiniaturizable = 1 << 2
	NSWindowStyleMaskResizable      = 1 << 3

	NSBackingStoreBuffered = 2

	kCGEventTapOptionDefault = 0
	kCGHIDEventTap           = 0
)

type cocoaDriver struct {
	app  uintptr // NSApplication
	pool uintptr // NSAutoreleasePool

	// Event management
	eventChan chan cocoaEvent
	mu        sync.RWMutex
	windows   map[uintptr]*cocoaWindow
	running   bool
}

type cocoaEvent struct {
	window *cocoaWindow
	event  interface{}
}

func loadCocoaFrameworks() error {
	var err error

	// Load Objective-C runtime
	libobjc, err := purego.Dlopen("libobjc.A.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load libobjc: %v", err)
	}

	// Load Foundation framework
	libFoundation, err = purego.Dlopen("/System/Library/Frameworks/Foundation.framework/Foundation", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load Foundation framework: %v", err)
	}

	// Load AppKit framework
	libAppKit, err = purego.Dlopen("/System/Library/Frameworks/AppKit.framework/AppKit", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load AppKit framework: %v", err)
	}

	// Load Core Graphics framework
	libCoreGraphics, err = purego.Dlopen("/System/Library/Frameworks/CoreGraphics.framework/CoreGraphics", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load CoreGraphics framework: %v", err)
	}

	// Register Objective-C runtime functions
	purego.RegisterLibFunc(&objc_getClass, libobjc, "objc_getClass")
	purego.RegisterLibFunc(&objc_msgSend, libobjc, "objc_msgSend")
	purego.RegisterLibFunc(&objc_msgSend_fpret, libobjc, "objc_msgSend_fpret")
	purego.RegisterLibFunc(&objc_msgSend_stret, libobjc, "objc_msgSend_stret")
	purego.RegisterLibFunc(&sel_registerName, libobjc, "sel_registerName")
	purego.RegisterLibFunc(&class_createInstance, libobjc, "class_createInstance")
	purego.RegisterLibFunc(&objc_allocateClassPair, libobjc, "objc_allocateClassPair")
	purego.RegisterLibFunc(&objc_registerClassPair, libobjc, "objc_registerClassPair")
	purego.RegisterLibFunc(&class_addMethod, libobjc, "class_addMethod")
	purego.RegisterLibFunc(&class_addIvar, libobjc, "class_addIvar")
	purego.RegisterLibFunc(&object_getInstanceVariable, libobjc, "object_getInstanceVariable")
	purego.RegisterLibFunc(&object_setInstanceVariable, libobjc, "object_setInstanceVariable")

	// Register Core Graphics functions
	purego.RegisterLibFunc(&CGMainDisplayID, libCoreGraphics, "CGMainDisplayID")
	purego.RegisterLibFunc(&CGDisplayBounds, libCoreGraphics, "CGDisplayBounds")
	purego.RegisterLibFunc(&CGEventCreateKeyboardEvent, libCoreGraphics, "CGEventCreateKeyboardEvent")
	purego.RegisterLibFunc(&CGEventPost, libCoreGraphics, "CGEventPost")
	purego.RegisterLibFunc(&CGEventCreateMouseEvent, libCoreGraphics, "CGEventCreateMouseEvent")

	return nil
}

func initializeCocoaClasses() error {
	// Get classes
	NSApplicationClass = objc_getClass(cString("NSApplication"))
	NSWindowClass = objc_getClass(cString("NSWindow"))
	NSViewClass = objc_getClass(cString("NSView"))
	NSEventClass = objc_getClass(cString("NSEvent"))
	NSStringClass = objc_getClass(cString("NSString"))
	NSAutoreleasePoolClass = objc_getClass(cString("NSAutoreleasePool"))

	if NSApplicationClass == 0 || NSWindowClass == 0 || NSViewClass == 0 ||
		NSEventClass == 0 || NSStringClass == 0 || NSAutoreleasePoolClass == 0 {
		return fmt.Errorf("failed to get required Cocoa classes")
	}

	// Register selectors
	sel_sharedApplication = sel_registerName(cString("sharedApplication"))
	sel_setActivationPolicy = sel_registerName(cString("setActivationPolicy:"))
	sel_activateIgnoringOtherApps = sel_registerName(cString("activateIgnoringOtherApps:"))
	sel_run = sel_registerName(cString("run"))
	sel_stop = sel_registerName(cString("stop:"))
	sel_nextEventMatchingMask = sel_registerName(cString("nextEventMatchingMask:untilDate:inMode:dequeue:"))
	sel_sendEvent = sel_registerName(cString("sendEvent:"))
	sel_updateWindows = sel_registerName(cString("updateWindows"))

	sel_alloc = sel_registerName(cString("alloc"))
	sel_init = sel_registerName(cString("init"))
	sel_initWithContentRect = sel_registerName(cString("initWithContentRect:styleMask:backing:defer:"))
	sel_setTitle = sel_registerName(cString("setTitle:"))
	sel_makeKeyAndOrderFront = sel_registerName(cString("makeKeyAndOrderFront:"))
	sel_orderOut = sel_registerName(cString("orderOut:"))
	sel_close = sel_registerName(cString("close"))
	sel_setDelegate = sel_registerName(cString("setDelegate:"))
	sel_frame = sel_registerName(cString("frame"))
	sel_setFrame = sel_registerName(cString("setFrame:display:"))
	sel_center = sel_registerName(cString("center"))
	sel_contentView = sel_registerName(cString("contentView"))
	sel_setContentSize = sel_registerName(cString("setContentSize:"))
	sel_setFrameOrigin = sel_registerName(cString("setFrameOrigin:"))
	sel_makeKeyWindow = sel_registerName(cString("makeKeyWindow"))
	sel_zoom = sel_registerName(cString("zoom:"))
	sel_setAcceptsFirstResponder = sel_registerName(cString("setAcceptsFirstResponder:"))

	// Window delegate selectors
	sel_windowShouldClose = sel_registerName(cString("windowShouldClose:"))
	sel_windowDidResize = sel_registerName(cString("windowDidResize:"))
	sel_windowDidBecomeKey = sel_registerName(cString("windowDidBecomeKey:"))
	sel_windowDidResignKey = sel_registerName(cString("windowDidResignKey:"))

	// Event selectors
	sel_type = sel_registerName(cString("type"))
	sel_locationInWindow = sel_registerName(cString("locationInWindow"))
	sel_keyCode = sel_registerName(cString("keyCode"))
	sel_modifierFlags = sel_registerName(cString("modifierFlags"))
	sel_isARepeat = sel_registerName(cString("isARepeat"))
	sel_scrollingDeltaX = sel_registerName(cString("scrollingDeltaX"))
	sel_scrollingDeltaY = sel_registerName(cString("scrollingDeltaY"))

	sel_stringWithUTF8String = sel_registerName(cString("stringWithUTF8String:"))
	sel_autorelease = sel_registerName(cString("autorelease"))
	sel_release = sel_registerName(cString("release"))
	sel_retain = sel_registerName(cString("retain"))

	return nil
}

func newCocoaDriver() (gio.Driver, error) {
	// Check if we're on macOS
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("cocoa driver: only supported on macOS")
	}

	// Load frameworks
	if err := loadCocoaFrameworks(); err != nil {
		return nil, fmt.Errorf("cocoa driver: %v", err)
	}

	// Initialize classes and selectors
	if err := initializeCocoaClasses(); err != nil {
		return nil, fmt.Errorf("cocoa driver: %v", err)
	}

	// Create autorelease pool
	pool := objc_msgSend(NSAutoreleasePoolClass, sel_alloc)
	pool = objc_msgSend(pool, sel_init)

	// Get shared application
	app := objc_msgSend(NSApplicationClass, sel_sharedApplication)
	if app == 0 {
		return nil, fmt.Errorf("cocoa driver: failed to get NSApplication")
	}

	// Set activation policy
	objc_msgSend(app, sel_setActivationPolicy, NSApplicationActivationPolicyRegular)

	cd := &cocoaDriver{
		app:       app,
		pool:      pool,
		eventChan: make(chan cocoaEvent, 256),
		windows:   make(map[uintptr]*cocoaWindow),
		running:   true,
	}

	// Start event processing loop
	go cd.eventLoop()

	return cd, nil
}

// CreateWindow implements gio.Driver
func (cd *cocoaDriver) CreateWindow(o gio.NewWindowOptions) (gio.Window, error) {
	bw := gio.NewBaseWindow(o)
	return newCocoaWindow(cd, bw)
}

// Type implements gio.Driver
func (cd *cocoaDriver) Type() gio.DriverType {
	return gio.DriverTypeCocoa
}

func (cd *cocoaDriver) eventLoop() {
	for cd.running {
		// Check for events in queue
		select {
		case event := <-cd.eventChan:
			if event.window != nil {
				// Handle event for specific window
				_ = event
			}
		default:
			// Process Cocoa events
			cd.processCocoaEvents()
			runtime.Gosched()
		}
	}
}

func (cd *cocoaDriver) processCocoaEvents() {
	// Get the default run loop mode (NSDefaultRunLoopMode)
	defaultMode := nsString("NSDefaultRunLoopMode")

	// Poll for events without blocking
	event := objc_msgSend(cd.app, sel_nextEventMatchingMask,
		uintptr(NSEventMaskAny),
		0, // untilDate (nil for non-blocking)
		defaultMode,
		1) // dequeue: YES

	if event != 0 {
		// Send event to application for processing
		objc_msgSend(cd.app, sel_sendEvent, event)

		// Update windows
		objc_msgSend(cd.app, sel_updateWindows)
	}
}

func (cd *cocoaDriver) registerWindow(window uintptr, cocoaWin *cocoaWindow) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	cd.windows[window] = cocoaWin
}

func (cd *cocoaDriver) unregisterWindow(window uintptr) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	delete(cd.windows, window)
}

func (cd *cocoaDriver) getWindow(window uintptr) *cocoaWindow {
	cd.mu.RLock()
	defer cd.mu.RUnlock()
	return cd.windows[window]
}

func (cd *cocoaDriver) shutdown() {
	cd.mu.Lock()
	cd.running = false
	cd.mu.Unlock()

	// Stop the application
	objc_msgSend(cd.app, sel_stop, cd.app)

	// Release autorelease pool
	if cd.pool != 0 {
		objc_msgSend(cd.pool, sel_release)
	}

	close(cd.eventChan)
}

// Convert Cocoa key codes to gio KeyCode
func (cd *cocoaDriver) translateKeyCode(keyCode uint16) gio.KeyCode {
	// macOS virtual key codes
	switch keyCode {
	case 0x35: // kVK_Escape
		return gio.KeyCodeEscape
	case 0x30: // kVK_Tab
		return gio.KeyCodeTab
	case 0x39: // kVK_CapsLock
		return gio.KeyCodeCapsLock
	case 0x33: // kVK_Delete
		return gio.KeyCodeBackspace
	case 0x24: // kVK_Return
		return gio.KeyCodeEnter
	case 0x31: // kVK_Space
		return gio.KeyCodeSpace

	// Function keys
	case 0x7A: // kVK_F1
		return gio.KeyCodeF1
	case 0x78: // kVK_F2
		return gio.KeyCodeF2
	case 0x63: // kVK_F3
		return gio.KeyCodeF3
	case 0x76: // kVK_F4
		return gio.KeyCodeF4
	case 0x60: // kVK_F5
		return gio.KeyCodeF5
	case 0x61: // kVK_F6
		return gio.KeyCodeF6
	case 0x62: // kVK_F7
		return gio.KeyCodeF7
	case 0x64: // kVK_F8
		return gio.KeyCodeF8
	case 0x65: // kVK_F9
		return gio.KeyCodeF9
	case 0x6D: // kVK_F10
		return gio.KeyCodeF10
	case 0x67: // kVK_F11
		return gio.KeyCodeF11
	case 0x6F: // kVK_F12
		return gio.KeyCodeF12

	// Numbers
	case 0x12: // kVK_ANSI_1
		return gio.KeyCodeDigit1
	case 0x13: // kVK_ANSI_2
		return gio.KeyCodeDigit2
	case 0x14: // kVK_ANSI_3
		return gio.KeyCodeDigit3
	case 0x15: // kVK_ANSI_4
		return gio.KeyCodeDigit4
	case 0x17: // kVK_ANSI_5
		return gio.KeyCodeDigit5
	case 0x16: // kVK_ANSI_6
		return gio.KeyCodeDigit6
	case 0x1A: // kVK_ANSI_7
		return gio.KeyCodeDigit7
	case 0x1C: // kVK_ANSI_8
		return gio.KeyCodeDigit8
	case 0x19: // kVK_ANSI_9
		return gio.KeyCodeDigit9
	case 0x1D: // kVK_ANSI_0
		return gio.KeyCodeDigit0

	// Letters
	case 0x00: // kVK_ANSI_A
		return gio.KeyCodeKeyA
	case 0x0B: // kVK_ANSI_B
		return gio.KeyCodeKeyB
	case 0x08: // kVK_ANSI_C
		return gio.KeyCodeKeyC
	case 0x02: // kVK_ANSI_D
		return gio.KeyCodeKeyD
	case 0x0E: // kVK_ANSI_E
		return gio.KeyCodeKeyE
	case 0x03: // kVK_ANSI_F
		return gio.KeyCodeKeyF
	case 0x05: // kVK_ANSI_G
		return gio.KeyCodeKeyG
	case 0x04: // kVK_ANSI_H
		return gio.KeyCodeKeyH
	case 0x22: // kVK_ANSI_I
		return gio.KeyCodeKeyI
	case 0x26: // kVK_ANSI_J
		return gio.KeyCodeKeyJ
	case 0x28: // kVK_ANSI_K
		return gio.KeyCodeKeyK
	case 0x25: // kVK_ANSI_L
		return gio.KeyCodeKeyL
	case 0x2E: // kVK_ANSI_M
		return gio.KeyCodeKeyM
	case 0x2D: // kVK_ANSI_N
		return gio.KeyCodeKeyN
	case 0x1F: // kVK_ANSI_O
		return gio.KeyCodeKeyO
	case 0x23: // kVK_ANSI_P
		return gio.KeyCodeKeyP
	case 0x0C: // kVK_ANSI_Q
		return gio.KeyCodeKeyQ
	case 0x0F: // kVK_ANSI_R
		return gio.KeyCodeKeyR
	case 0x01: // kVK_ANSI_S
		return gio.KeyCodeKeyS
	case 0x11: // kVK_ANSI_T
		return gio.KeyCodeKeyT
	case 0x20: // kVK_ANSI_U
		return gio.KeyCodeKeyU
	case 0x09: // kVK_ANSI_V
		return gio.KeyCodeKeyV
	case 0x0D: // kVK_ANSI_W
		return gio.KeyCodeKeyW
	case 0x07: // kVK_ANSI_X
		return gio.KeyCodeKeyX
	case 0x10: // kVK_ANSI_Y
		return gio.KeyCodeKeyY
	case 0x06: // kVK_ANSI_Z
		return gio.KeyCodeKeyZ

	// Arrow keys
	case 0x7E: // kVK_UpArrow
		return gio.KeyCodeArrowUp
	case 0x7D: // kVK_DownArrow
		return gio.KeyCodeArrowDown
	case 0x7B: // kVK_LeftArrow
		return gio.KeyCodeArrowLeft
	case 0x7C: // kVK_RightArrow
		return gio.KeyCodeArrowRight

	default:
		return gio.KeyCodeUnknown
	}
}

// Convert Cocoa modifiers to gio ModifierKey
func (cd *cocoaDriver) translateModifiers(modifiers uint64) gio.ModifierKey {
	var gioMods gio.ModifierKey

	// NSEventModifierFlagShift = 1 << 17
	if modifiers&(1<<17) != 0 {
		gioMods |= gio.ModifierKeyShift
	}

	// NSEventModifierFlagControl = 1 << 18
	if modifiers&(1<<18) != 0 {
		gioMods |= gio.ModifierKeyCtrl
	}

	// NSEventModifierFlagOption = 1 << 19 (Alt)
	if modifiers&(1<<19) != 0 {
		gioMods |= gio.ModifierKeyAlt
	}

	// NSEventModifierFlagCommand = 1 << 20 (Cmd/Meta)
	if modifiers&(1<<20) != 0 {
		gioMods |= gio.ModifierKeyMeta
	}

	return gioMods
}

// Convert Cocoa mouse button to gio MouseButton
func (cd *cocoaDriver) translateMouseButton(buttonNumber int) gio.MouseButton {
	switch buttonNumber {
	case 0: // Left button
		return gio.MouseButtonLeft
	case 1: // Right button
		return gio.MouseButtonRight
	case 2: // Middle button
		return gio.MouseButtonMiddle
	case 3: // Button 4
		return gio.MouseButtonBack
	case 4: // Button 5
		return gio.MouseButtonForward
	default:
		return gio.MouseButtonUnknown
	}
}

// Helper functions
func cString(s string) *byte {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return &b[0]
}

func nsString(s string) uintptr {
	cstr := cString(s)
	return objc_msgSend(NSStringClass, sel_stringWithUTF8String, uintptr(unsafe.Pointer(cstr)))
}

// Helper methods for creating NSRect, NSPoint, NSSize
func (cd *cocoaDriver) nsRectMake(x, y, w, h float64) NSRect {
	return NSRect{
		Origin: CGPoint{X: x, Y: y},
		Size:   CGSize{Width: w, Height: h},
	}
}

func (cd *cocoaDriver) nsPointMake(x, y float64) NSPoint {
	return NSPoint{X: x, Y: y}
}

func (cd *cocoaDriver) nsSizeMake(w, h float64) NSSize {
	return NSSize{Width: w, Height: h}
}

func (cd *cocoaDriver) nsStringWithGoString(s string) uintptr {
	return nsString(s)
}

type nsRect NSRect
type nsPoint NSPoint
type nsSize NSSize
