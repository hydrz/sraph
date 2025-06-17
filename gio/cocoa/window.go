package cocoa

import (
	"fmt"
	"image"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/opensraph/sraph/gio"
)

var _ gio.Window = (*cocoaWindow)(nil)

// cocoaWindow represents a Cocoa window implementation
type cocoaWindow struct {
	gio.BaseWindow

	driver   *cocoaDriver
	nsWindow uintptr // NSWindow instance
	delegate uintptr // NSWindowDelegate instance
	view     uintptr // NSView instance

	mu sync.RWMutex
}

// windowDelegateClass holds the Objective-C class for window delegate
var windowDelegateClass uintptr

func newCocoaWindow(driver *cocoaDriver, bw gio.BaseWindow) (*cocoaWindow, error) {
	cw := &cocoaWindow{
		driver:     driver,
		BaseWindow: bw,
	}

	if err := cw.init(); err != nil {
		return nil, fmt.Errorf("cocoa: cocoaWindow.init failed: %v", err)
	}

	return cw, nil
}

func (cw *cocoaWindow) init() error {
	// Initialize window delegate class if not already done
	if windowDelegateClass == 0 {
		if err := cw.initWindowDelegateClass(); err != nil {
			return fmt.Errorf("failed to initialize window delegate class: %v", err)
		}
	}

	// Create NSWindow
	if err := cw.createNSWindow(); err != nil {
		return fmt.Errorf("failed to create NSWindow: %v", err)
	}

	// Create window delegate
	if err := cw.createWindowDelegate(); err != nil {
		return fmt.Errorf("failed to create window delegate: %v", err)
	}

	// Create content view
	if err := cw.createContentView(); err != nil {
		return fmt.Errorf("failed to create content view: %v", err)
	}

	// Configure window
	cw.configureWindow()

	return nil
}

func (cw *cocoaWindow) initWindowDelegateClass() error {
	// Create a custom NSWindowDelegate class
	superClass := objc_getClass(cString("NSObject"))
	if superClass == 0 {
		return fmt.Errorf("failed to get NSObject class")
	}

	className := "GioWindowDelegate"
	class := objc_allocateClassPair(superClass, cString(className), 0)
	if class == 0 {
		return fmt.Errorf("failed to allocate class pair")
	}

	// Add instance variable to store window reference
	class_addIvar(class,
		cString("window"),
		8,             // sizeof(void*)
		3,             // log2(8) alignment
		cString("^v")) // pointer type encoding

	// Add delegate methods
	cw.addDelegateMethods(class)

	// Register the class
	objc_registerClassPair(class)
	windowDelegateClass = class

	return nil
}

func (cw *cocoaWindow) addDelegateMethods(class uintptr) {
	// windowShouldClose: method
	shouldCloseImpl := func(self, cmd, sender uintptr) uintptr {
		// Get window reference from instance variable
		var windowPtr uintptr
		object_getInstanceVariable(self, cString("window"), &windowPtr)

		if windowPtr != 0 {
			window := (*cocoaWindow)(unsafe.Pointer(windowPtr))
			window.handleWindowShouldClose()
		}
		return 1 // YES - allow close
	}

	class_addMethod(class,
		sel_windowShouldClose,
		purego.NewCallback(shouldCloseImpl),
		cString("c@:@")) // returns char, takes id, SEL, id

	// windowDidResize: method
	didResizeImpl := func(self, cmd, notification uintptr) {
		var windowPtr uintptr
		object_getInstanceVariable(self, cString("window"), &windowPtr)

		if windowPtr != 0 {
			window := (*cocoaWindow)(unsafe.Pointer(windowPtr))
			window.handleWindowDidResize()
		}
	}

	class_addMethod(class,
		sel_windowDidResize,
		purego.NewCallback(didResizeImpl),
		cString("v@:@")) // void, takes id, SEL, id

	// windowDidBecomeKey: method
	didBecomeKeyImpl := func(self, cmd, notification uintptr) {
		var windowPtr uintptr
		object_getInstanceVariable(self, cString("window"), &windowPtr)

		if windowPtr != 0 {
			window := (*cocoaWindow)(unsafe.Pointer(windowPtr))
			window.handleWindowDidBecomeKey()
		}
	}

	class_addMethod(class,
		sel_windowDidBecomeKey,
		purego.NewCallback(didBecomeKeyImpl),
		cString("v@:@"))

	// windowDidResignKey: method
	didResignKeyImpl := func(self, cmd, notification uintptr) {
		var windowPtr uintptr
		object_getInstanceVariable(self, cString("window"), &windowPtr)

		if windowPtr != 0 {
			window := (*cocoaWindow)(unsafe.Pointer(windowPtr))
			window.handleWindowDidResignKey()
		}
	}

	class_addMethod(class,
		sel_windowDidResignKey,
		purego.NewCallback(didResignKeyImpl),
		cString("v@:@"))
}

func (cw *cocoaWindow) createNSWindow() error {
	attr := cw.Attr()

	// Default size if not specified
	width, height := float64(attr.Width), float64(attr.Height)
	if width <= 0 {
		width = 1024
	}
	if height <= 0 {
		height = 768
	}

	// Create NSRect for content
	contentRect := cw.driver.nsRectMake(0, 0, width, height)

	// Window style mask
	styleMask := uintptr(1 << 0) // NSWindowStyleMaskTitled
	styleMask |= uintptr(1 << 1) // NSWindowStyleMaskClosable
	styleMask |= uintptr(1 << 2) // NSWindowStyleMaskMiniaturizable

	if attr.State.Contains(gio.WindowStateResizable) {
		styleMask |= uintptr(1 << 3) // NSWindowStyleMaskResizable
	}

	// Allocate NSWindow
	windowClass := objc_getClass(cString("NSWindow"))
	window := objc_msgSend(windowClass, sel_alloc)

	// Initialize NSWindow
	cw.nsWindow = objc_msgSend(window,
		sel_initWithContentRect,
		uintptr(unsafe.Pointer(&contentRect)),
		styleMask,
		0, // NSBackingStoreBuffered
		0) // defer=NO

	if cw.nsWindow == 0 {
		return fmt.Errorf("failed to create NSWindow")
	}

	return nil
}

func (cw *cocoaWindow) createWindowDelegate() error {
	// Create delegate instance
	delegate := objc_msgSend(windowDelegateClass, sel_alloc)
	cw.delegate = objc_msgSend(delegate, sel_init)

	if cw.delegate == 0 {
		return fmt.Errorf("failed to create window delegate")
	}

	// Set window reference in delegate
	windowPtr := uintptr(unsafe.Pointer(cw))
	object_setInstanceVariable(cw.delegate, cString("window"), windowPtr)

	// Set delegate on window
	objc_msgSend(cw.nsWindow, sel_setDelegate, cw.delegate)

	return nil
}

func (cw *cocoaWindow) createContentView() error {
	// Get content view from window
	cw.view = objc_msgSend(cw.nsWindow, sel_contentView)

	if cw.view == 0 {
		return fmt.Errorf("failed to get content view")
	}

	// Set view to accept first responder
	objc_msgSend(cw.view, sel_setAcceptsFirstResponder, 1)

	return nil
}

func (cw *cocoaWindow) configureWindow() {
	attr := cw.Attr()

	// Set window title
	if attr.Title != "" {
		title := cw.driver.nsStringWithGoString(attr.Title)
		objc_msgSend(cw.nsWindow, sel_setTitle, title)
	}

	// Position window
	if attr.Position.X != 0 || attr.Position.Y != 0 {
		point := cw.driver.nsPointMake(float64(attr.Position.X), float64(attr.Position.Y))
		objc_msgSend(cw.nsWindow, sel_setFrameOrigin, uintptr(unsafe.Pointer(&point)))
	}

	// Show window if visible
	if attr.State.Contains(gio.WindowStateVisible) {
		objc_msgSend(cw.nsWindow, sel_makeKeyAndOrderFront, 0)
	}
}

// WindowID returns the Cocoa window ID
func (cw *cocoaWindow) WindowID() gio.WindowID {
	cw.mu.RLock()
	defer cw.mu.RUnlock()
	return gio.WindowID(cw.nsWindow)
}

func (cw *cocoaWindow) SetAttr(attr gio.WindowAttr) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	currentAttr := cw.Attr()
	if currentAttr.Equal(attr) {
		return
	}

	// Handle window closure first - early return
	if attr.State.Contains(gio.WindowStateClosed) {
		cw.handleWindowClose()
		cw.BaseWindow.SetAttr(attr)
		return
	}

	// Update basic properties
	cw.updateTitle(attr, currentAttr)
	cw.updateGeometry(attr, currentAttr)

	// Handle state changes
	cw.updateWindowStates(attr, currentAttr)

	// Update base attributes last
	cw.BaseWindow.SetAttr(attr)
}

func (cw *cocoaWindow) handleWindowClose() {
	if cw.nsWindow != 0 {
		objc_msgSend(cw.nsWindow, sel_close)
	}
}

func (cw *cocoaWindow) updateTitle(attr, currentAttr gio.WindowAttr) {
	if attr.Title == "" || attr.Title == currentAttr.Title {
		return
	}

	title := cw.driver.nsStringWithGoString(attr.Title)
	objc_msgSend(cw.nsWindow, sel_setTitle, title)
}

func (cw *cocoaWindow) updateGeometry(attr, currentAttr gio.WindowAttr) {
	sizeChanged := attr.Width != currentAttr.Width || attr.Height != currentAttr.Height
	posChanged := attr.Position.X != currentAttr.Position.X || attr.Position.Y != currentAttr.Position.Y

	if sizeChanged && attr.Width > 0 && attr.Height > 0 {
		size := cw.driver.nsSizeMake(float64(attr.Width), float64(attr.Height))
		objc_msgSend(cw.nsWindow, sel_setContentSize, uintptr(unsafe.Pointer(&size)))
	}

	if posChanged {
		point := cw.driver.nsPointMake(float64(attr.Position.X), float64(attr.Position.Y))
		objc_msgSend(cw.nsWindow, sel_setFrameOrigin, uintptr(unsafe.Pointer(&point)))
	}
}

func (cw *cocoaWindow) updateWindowStates(attr, currentAttr gio.WindowAttr) {
	// Handle visible state
	visibleChanged := attr.State.Contains(gio.WindowStateVisible) != currentAttr.State.Contains(gio.WindowStateVisible)
	if visibleChanged {
		if attr.State.Contains(gio.WindowStateVisible) {
			objc_msgSend(cw.nsWindow, sel_makeKeyAndOrderFront, 0)
		} else {
			objc_msgSend(cw.nsWindow, sel_orderOut, 0)
		}
	}

	// Handle maximized state
	maximizedChanged := attr.State.Contains(gio.WindowStateMaximized) != currentAttr.State.Contains(gio.WindowStateMaximized)
	if maximizedChanged {
		if attr.State.Contains(gio.WindowStateMaximized) {
			objc_msgSend(cw.nsWindow, sel_zoom, 0)
		}
		// Note: No direct way to un-maximize, user needs to click zoom button
	}

	// Handle focus state
	focusChanged := attr.State.Contains(gio.WindowStateFocused) != currentAttr.State.Contains(gio.WindowStateFocused)
	if focusChanged && attr.State.Contains(gio.WindowStateFocused) {
		objc_msgSend(cw.nsWindow, sel_makeKeyWindow)
	}
}

// Event handlers called from delegate methods
func (cw *cocoaWindow) handleWindowShouldClose() {
	attr := cw.Attr()
	attr.State |= gio.WindowStateClosed
	cw.BaseWindow.SetAttr(attr)
}

func (cw *cocoaWindow) handleWindowDidResize() {
	// Get current window frame
	frame := cw.getWindowFrame()

	attr := cw.Attr()
	changed := false

	// Update size if changed
	if attr.Width != int(frame.Size.Width) || attr.Height != int(frame.Size.Height) {
		attr.Width = int(frame.Size.Width)
		attr.Height = int(frame.Size.Height)
		changed = true
	}

	if changed {
		cw.BaseWindow.SetAttr(attr)
	}
}

func (cw *cocoaWindow) handleWindowDidBecomeKey() {
	attr := cw.Attr()
	if !attr.State.Contains(gio.WindowStateFocused) {
		attr.State |= gio.WindowStateFocused
		cw.BaseWindow.SetAttr(attr)
	}
}

func (cw *cocoaWindow) handleWindowDidResignKey() {
	attr := cw.Attr()
	if attr.State.Contains(gio.WindowStateFocused) {
		attr.State &^= gio.WindowStateFocused
		cw.BaseWindow.SetAttr(attr)
	}
}

// Helper methods
func (cw *cocoaWindow) getWindowFrame() NSRect {
	var frame NSRect
	// Get frame from window - this is simplified
	// Real implementation would need to handle returned structs properly
	objc_msgSend_stret(unsafe.Pointer(&frame), cw.nsWindow, sel_frame)
	return frame
}

// Event translation methods
func (cw *cocoaWindow) translateKeyCode(keyCode uint16) gio.KeyCode {
	return cw.driver.translateKeyCode(keyCode)
}

func (cw *cocoaWindow) translateModifiers(modifierFlags uint64) gio.ModifierKey {
	return cw.driver.translateModifiers(modifierFlags)
}

func (cw *cocoaWindow) translateMouseButton(buttonNumber int) gio.MouseButton {
	return cw.driver.translateMouseButton(buttonNumber)
}

// Simplified event handling - these would be called from the driver's event loop
func (cw *cocoaWindow) handleEvent(event uintptr) {
	eventType := objc_msgSend(event, sel_type)

	switch eventType {
	case 1: // NSEventTypeLeftMouseDown
		cw.handleMouseDown(event, gio.MouseButtonLeft)
	case 2: // NSEventTypeLeftMouseUp
		cw.handleMouseUp(event, gio.MouseButtonLeft)
	case 3: // NSEventTypeRightMouseDown
		cw.handleMouseDown(event, gio.MouseButtonRight)
	case 4: // NSEventTypeRightMouseUp
		cw.handleMouseUp(event, gio.MouseButtonRight)
	case 5: // NSEventTypeMouseMoved
		cw.handleMouseMoved(event)
	case 10: // NSEventTypeKeyDown
		cw.handleKeyDown(event)
	case 11: // NSEventTypeKeyUp
		cw.handleKeyUp(event)
	case 22: // NSEventTypeScrollWheel
		cw.handleScrollWheel(event)
	}
}

func (cw *cocoaWindow) handleMouseDown(event uintptr, button gio.MouseButton) {
	location := cw.getEventLocation(event)
	modifiers := cw.getEventModifiers(event)

	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = button
	mouseEvent.ModifierKey = cw.translateModifiers(modifiers)
	mouseEvent.Position = image.Point{X: int(location.X), Y: int(location.Y)}

	cw.Publish(mouseEvent)
}

func (cw *cocoaWindow) handleMouseUp(event uintptr, button gio.MouseButton) {
	location := cw.getEventLocation(event)
	modifiers := cw.getEventModifiers(event)

	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = button
	mouseEvent.ModifierKey = cw.translateModifiers(modifiers)
	mouseEvent.Position = image.Point{X: int(location.X), Y: int(location.Y)}

	cw.Publish(mouseEvent)
}

func (cw *cocoaWindow) handleMouseMoved(event uintptr) {
	location := cw.getEventLocation(event)
	modifiers := cw.getEventModifiers(event)

	mouseEvent := gio.NewMouseEvent()
	mouseEvent.Button = gio.MouseButtonUnknown
	mouseEvent.ModifierKey = cw.translateModifiers(modifiers)
	mouseEvent.Position = image.Point{X: int(location.X), Y: int(location.Y)}

	cw.Publish(mouseEvent)
}

func (cw *cocoaWindow) handleKeyDown(event uintptr) {
	keyCode := cw.getEventKeyCode(event)
	modifiers := cw.getEventModifiers(event)
	isRepeat := cw.getEventIsRepeat(event)

	keyEvent := gio.NewKeyboardEvent()
	keyEvent.Code = cw.translateKeyCode(keyCode)
	keyEvent.ModifierKey = cw.translateModifiers(modifiers)
	keyEvent.Repeat = isRepeat

	cw.Publish(keyEvent)
}

func (cw *cocoaWindow) handleKeyUp(event uintptr) {
	keyCode := cw.getEventKeyCode(event)
	modifiers := cw.getEventModifiers(event)

	keyEvent := gio.NewKeyboardEvent()
	keyEvent.Code = cw.translateKeyCode(keyCode)
	keyEvent.ModifierKey = cw.translateModifiers(modifiers)
	keyEvent.Repeat = false

	cw.Publish(keyEvent)
}

func (cw *cocoaWindow) handleScrollWheel(event uintptr) {
	deltaX := cw.getEventScrollDeltaX(event)
	deltaY := cw.getEventScrollDeltaY(event)
	location := cw.getEventLocation(event)
	modifiers := cw.getEventModifiers(event)

	wheelEvent := gio.NewWheelEvent()
	wheelEvent.DeltaX = deltaX
	wheelEvent.DeltaY = deltaY
	wheelEvent.ModifierKey = cw.translateModifiers(modifiers)
	wheelEvent.Position = image.Point{X: int(location.X), Y: int(location.Y)}

	cw.Publish(wheelEvent)
}

// Helper methods to extract event properties
func (cw *cocoaWindow) getEventLocation(event uintptr) NSPoint {
	var point NSPoint
	objc_msgSend_stret(unsafe.Pointer(&point), event, sel_locationInWindow)
	return point
}

func (cw *cocoaWindow) getEventModifiers(event uintptr) uint64 {
	return uint64(objc_msgSend(event, sel_modifierFlags))
}

func (cw *cocoaWindow) getEventKeyCode(event uintptr) uint16 {
	return uint16(objc_msgSend(event, sel_keyCode))
}

func (cw *cocoaWindow) getEventIsRepeat(event uintptr) bool {
	return objc_msgSend(event, sel_isARepeat) != 0
}

func (cw *cocoaWindow) getEventScrollDeltaX(event uintptr) float64 {
	return objc_msgSend_fpret(event, sel_scrollingDeltaX)
}

func (cw *cocoaWindow) getEventScrollDeltaY(event uintptr) float64 {
	return objc_msgSend_fpret(event, sel_scrollingDeltaY)
}
