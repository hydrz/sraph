package wayland

import (
	"log"
	"syscall"

	"github.com/opensraph/sraph/gio"
	"github.com/rajveermalviya/go-wayland/wayland/client"
)

// Keyboard event handlers for WaylandWindow
func (w *WaylandWindow) handleKeyboardKey(key uint32, pressed bool, modifiers, serial, time uint32) {
	if w.IsClosed() {
		return
	}

	keyboardEvent := gio.NewKeyboardEvent()
	keyboardEvent.Code = waylandKeycodeToGio(key)
	keyboardEvent.ModifierKey = waylandModifiersToGio(modifiers)

	if pressed {
		keyboardEvent.State = gio.KeyStatePressed
	} else {
		keyboardEvent.State = gio.KeyStateReleased
	}

	if err := w.Publish(keyboardEvent); err != nil {
		log.Printf("wayland: failed to publish keyboard event: %v", err)
	}
}

func (w *WaylandWindow) handleKeyboardModifiers(depressed, latched, locked, group, serial uint32) {
	// Modifier state is tracked but doesn't generate events by itself
	log.Printf("wayland: keyboard modifiers changed: depressed=%x, latched=%x, locked=%x, group=%x",
		depressed, latched, locked, group)
}

func (w *WaylandWindow) handleKeyboardKeymap(format, size uint32, fd int) {
	// TODO: Implement proper keymap handling with libxkbcommon
	log.Printf("wayland: keyboard keymap changed: format=%d, size=%d", format, size)
}

func (w *WaylandWindow) handleKeyboardEnter(keys []byte, serial uint32) {
	if w.IsClosed() {
		return
	}

	// Update window focus state
	w.mu.Lock()
	attr := w.BaseWindow.Attr()
	attr.State |= gio.WindowStateFocused
	w.BaseWindow.SetAttr(attr)
	w.mu.Unlock()

	// Publish window event
	windowEvent := gio.NewWindowEvent()
	windowEvent.Window = w

	if err := w.Publish(windowEvent); err != nil {
		log.Printf("wayland: failed to publish window focus event: %v", err)
	}
}

func (w *WaylandWindow) handleKeyboardLeave(serial uint32) {
	if w.IsClosed() {
		return
	}

	// Update window focus state
	w.mu.Lock()
	attr := w.BaseWindow.Attr()
	attr.State &= ^gio.WindowStateFocused
	w.BaseWindow.SetAttr(attr)
	w.mu.Unlock()

	// Publish window event
	windowEvent := gio.NewWindowEvent()
	windowEvent.Window = w

	if err := w.Publish(windowEvent); err != nil {
		log.Printf("wayland: failed to publish window unfocus event: %v", err)
	}
}

func (w *WaylandWindow) handleKeyboardRepeatInfo(rate, delay int32) {
	log.Printf("wayland: keyboard repeat info: rate=%d, delay=%d", rate, delay)
	// TODO: Use this information for key repeat handling
}

// Driver keyboard event handlers
func (d *WaylandDriver) handleKeyboardKey(e client.KeyboardKeyEvent) {
	d.eventState.keyboardEvent.serial = e.Serial
	d.eventState.keyboardEvent.time = e.Time
	d.eventState.keyboardEvent.key = e.Key
	d.eventState.keyboardEvent.state = e.State

	if d.focusedWindow != nil {
		pressed := e.State == uint32(client.KeyboardKeyStatePressed)
		d.focusedWindow.handleKeyboardKey(e.Key, pressed, d.eventState.keyboardEvent.modifiers, e.Serial, e.Time)
	}
}

func (d *WaylandDriver) handleKeyboardModifiers(e client.KeyboardModifiersEvent) {
	d.eventState.keyboardEvent.modifiers = e.ModsDepressed

	if d.focusedWindow != nil {
		d.focusedWindow.handleKeyboardModifiers(e.ModsDepressed, e.ModsLatched, e.ModsLocked, e.Group, e.Serial)
	}
}

func (d *WaylandDriver) handleKeyboardKeymap(e client.KeyboardKeymapEvent) {
	defer syscall.Close(e.Fd)

	if d.focusedWindow != nil {
		d.focusedWindow.handleKeyboardKeymap(e.Format, e.Size, e.Fd)
	}
}

func (d *WaylandDriver) handleKeyboardEnter(e client.KeyboardEnterEvent) {
	d.updateSerial(e.Serial)

	// Find window for this surface
	if window := d.findWindowBySurface(e.Surface); window != nil {
		d.focusedWindow = window
		window.handleKeyboardEnter(e.Keys, e.Serial)
	}
}

func (d *WaylandDriver) handleKeyboardLeave(e client.KeyboardLeaveEvent) {
	if d.focusedWindow != nil {
		d.focusedWindow.handleKeyboardLeave(e.Serial)
	}
}

func (d *WaylandDriver) handleKeyboardRepeatInfo(e client.KeyboardRepeatInfoEvent) {
	if d.focusedWindow != nil {
		d.focusedWindow.handleKeyboardRepeatInfo(e.Rate, e.Delay)
	}
}

// attachKeyboard sets up the keyboard interface
func (d *WaylandDriver) attachKeyboard() {
	keyboard, err := d.seat.GetKeyboard()
	if err != nil {
		log.Printf("wayland: failed to get keyboard: %v", err)
		return
	}
	d.keyboard = keyboard

	// Set up keyboard event handlers
	d.keyboard.SetKeyHandler(d.handleKeyboardKey)
	d.keyboard.SetModifiersHandler(d.handleKeyboardModifiers)
	d.keyboard.SetKeymapHandler(d.handleKeyboardKeymap)
	d.keyboard.SetEnterHandler(d.handleKeyboardEnter)
	d.keyboard.SetLeaveHandler(d.handleKeyboardLeave)
	d.keyboard.SetRepeatInfoHandler(d.handleKeyboardRepeatInfo)

	log.Printf("wayland: keyboard interface registered")
}

// releaseKeyboard releases the keyboard interface
func (d *WaylandDriver) releaseKeyboard() {
	if d.keyboard != nil && d.seatVersion >= 3 {
		if err := d.keyboard.Release(); err != nil {
			log.Printf("wayland: failed to release keyboard: %v", err)
		}
	}
	d.keyboard = nil
	log.Printf("wayland: keyboard interface released")
}

// waylandKeycodeToGio converts Wayland keycode to gio KeyCode
func waylandKeycodeToGio(key uint32) gio.KeyCode {
	// Wayland uses Linux input event codes + 8
	// Convert back to standard scan codes
	scancode := key - 8

	if gioKeyCode, exists := waylandKeycodeMap[scancode]; exists {
		return gioKeyCode
	}

	return gio.KeyCodeUnknown
}

// waylandModifiersToGio converts Wayland modifiers to gio ModifierKey
func waylandModifiersToGio(modifiers uint32) gio.ModifierKey {
	var result gio.ModifierKey

	for _, mapping := range waylandModifierMap {
		if modifiers&mapping.waylandMask != 0 {
			result |= mapping.gioModifier
		}
	}

	return result
}

// Wayland keycode to gio KeyCode lookup table
var waylandKeycodeMap = map[uint32]gio.KeyCode{
	// Control keys
	1:  gio.KeyCodeEscape,
	14: gio.KeyCodeBackspace,
	15: gio.KeyCodeTab,
	28: gio.KeyCodeEnter,
	57: gio.KeyCodeSpace,

	// Modifier keys
	42: gio.KeyCodeShiftLeft,
	54: gio.KeyCodeShiftRight,

	// Number row
	2:  gio.KeyCodeDigit1,
	3:  gio.KeyCodeDigit2,
	4:  gio.KeyCodeDigit3,
	5:  gio.KeyCodeDigit4,
	6:  gio.KeyCodeDigit5,
	7:  gio.KeyCodeDigit6,
	8:  gio.KeyCodeDigit7,
	9:  gio.KeyCodeDigit8,
	10: gio.KeyCodeDigit9,
	11: gio.KeyCodeDigit0,

	// Alphabet keys
	16: gio.KeyCodeKeyQ,
	17: gio.KeyCodeKeyW,
	18: gio.KeyCodeKeyE,
	19: gio.KeyCodeKeyR,
	20: gio.KeyCodeKeyT,
	21: gio.KeyCodeKeyY,
	22: gio.KeyCodeKeyU,
	23: gio.KeyCodeKeyI,
	24: gio.KeyCodeKeyO,
	25: gio.KeyCodeKeyP,
	30: gio.KeyCodeKeyA,
	31: gio.KeyCodeKeyS,
	32: gio.KeyCodeKeyD,
	33: gio.KeyCodeKeyF,
	34: gio.KeyCodeKeyG,
	35: gio.KeyCodeKeyH,
	36: gio.KeyCodeKeyJ,
	37: gio.KeyCodeKeyK,
	38: gio.KeyCodeKeyL,
	44: gio.KeyCodeKeyZ,
	45: gio.KeyCodeKeyX,
	46: gio.KeyCodeKeyC,
	47: gio.KeyCodeKeyV,
	48: gio.KeyCodeKeyB,
	49: gio.KeyCodeKeyN,
	50: gio.KeyCodeKeyM,

	// Function keys
	59: gio.KeyCodeF1,
	60: gio.KeyCodeF2,
	61: gio.KeyCodeF3,
	62: gio.KeyCodeF4,
	63: gio.KeyCodeF5,
	64: gio.KeyCodeF6,
	65: gio.KeyCodeF7,
	66: gio.KeyCodeF8,
	67: gio.KeyCodeF9,
	68: gio.KeyCodeF10,
	87: gio.KeyCodeF11,
	88: gio.KeyCodeF12,

	// Other keys
	58: gio.KeyCodeCapsLock,
}

// Wayland modifier to gio ModifierKey lookup table
var waylandModifierMap = []struct {
	waylandMask uint32
	gioModifier gio.ModifierKey
}{
	{1 << 0, gio.ModifierKeyShift},    // WL_KEYBOARD_MODIFIER_SHIFT
	{1 << 1, gio.ModifierKeyCapsLock}, // WL_KEYBOARD_MODIFIER_CAPS
	{1 << 2, gio.ModifierKeyCtrl},     // WL_KEYBOARD_MODIFIER_CTRL
	{1 << 3, gio.ModifierKeyAlt},      // WL_KEYBOARD_MODIFIER_ALT
	{1 << 4, gio.ModifierKeyNumLock},  // WL_KEYBOARD_MODIFIER_NUM
	{1 << 6, gio.ModifierKeyMeta},     // WL_KEYBOARD_MODIFIER_LOGO
}
