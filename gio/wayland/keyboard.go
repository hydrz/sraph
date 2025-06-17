package wayland

import (
	"log/slog"
	"syscall"

	"github.com/opensraph/sraph/gio"
	"github.com/rajveermalviya/go-wayland/wayland/client"
)

// publishKeyboardEvent publishes keyboard event with error handling
func (w *WaylandWindow) publishKeyboardEvent(event *gio.KeyboardEvent) {
	if w.IsClosed() {
		return
	}

	if err := w.Publish(event); err != nil {
		slog.Error("failed to publish keyboard event", "error", err)
	}
}

// publishWindowEvent publishes window event with error handling
func (w *WaylandWindow) publishWindowEvent() {
	if w.IsClosed() {
		return
	}

	windowEvent := gio.NewWindowEvent()
	windowEvent.Window = w

	if err := w.Publish(windowEvent); err != nil {
		slog.Error("failed to publish window event", "error", err)
	}
}

// Keyboard event handlers for WaylandWindow
func (w *WaylandWindow) handleKeyboardKey(key uint32, pressed bool, modifiers, serial, time uint32) {
	keyboardEvent := gio.NewKeyboardEvent()
	keyboardEvent.Code = waylandKeycodeToGio(key)
	keyboardEvent.ModifierKey = waylandModifiersToGio(modifiers)

	if pressed {
		keyboardEvent.State = gio.KeyStatePressed
	} else {
		keyboardEvent.State = gio.KeyStateReleased
	}

	slog.Debug("publishing keyboard event",
		"key", key,
		"code", keyboardEvent.Code,
		"modifiers", keyboardEvent.ModifierKey,
		"state", keyboardEvent.State)

	w.publishKeyboardEvent(keyboardEvent)
}

func (w *WaylandWindow) handleKeyboardModifiers(depressed, latched, locked, group, serial uint32) {
	// Modifier state is tracked but doesn't generate events by itself
	slog.Debug("keyboard modifiers changed",
		"depressed", depressed,
		"latched", latched,
		"locked", locked,
		"group", group)
}

func (w *WaylandWindow) handleKeyboardKeymap(format, size uint32, fd int) {
	// TODO: Implement proper keymap handling with libxkbcommon
	slog.Debug("keyboard keymap changed", "format", format, "size", size)
}

func (w *WaylandWindow) handleKeyboardEnter(keys []byte, serial uint32) {
	// Update window focus state
	w.updateFocusState(true)
	w.publishWindowEvent()
}

func (w *WaylandWindow) handleKeyboardLeave(serial uint32) {
	// Update window focus state
	w.updateFocusState(false)
	w.publishWindowEvent()
}

// updateFocusState updates window focus state safely
func (w *WaylandWindow) updateFocusState(focused bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	attr := w.BaseWindow.Attr()
	if focused {
		attr.State |= gio.WindowStateFocused
	} else {
		attr.State &= ^gio.WindowStateFocused
	}
	w.BaseWindow.SetAttr(attr)
}

func (w *WaylandWindow) handleKeyboardRepeatInfo(rate, delay int32) {
	slog.Debug("keyboard repeat info", "rate", rate, "delay", delay)
	// TODO: Use this information for key repeat handling
}

// Driver keyboard event handlers
func (d *WaylandDriver) handleKeyboardKey(e client.KeyboardKeyEvent) {
	d.updateSerial(e.Serial)

	d.eventState.keyboardEvent.serial = e.Serial
	d.eventState.keyboardEvent.time = e.Time
	d.eventState.keyboardEvent.key = e.Key
	d.eventState.keyboardEvent.state = e.State

	slog.Debug("keyboard key event",
		"key", e.Key,
		"state", e.State,
		"modifiers", d.eventState.keyboardEvent.modifiers,
		"focused_window", d.focusedWindow != nil)

	if d.focusedWindow != nil {
		pressed := e.State == uint32(client.KeyboardKeyStatePressed)

		// For combination keys, use the current modifier state as-is
		// Don't modify it based on the current key being pressed/released
		currentModifiers := d.eventState.keyboardEvent.modifiers

		slog.Debug("dispatching keyboard event",
			"key", e.Key,
			"pressed", pressed,
			"modifiers", currentModifiers,
			"window_id", d.focusedWindow.WindowID())

		d.focusedWindow.handleKeyboardKey(e.Key, pressed, currentModifiers, e.Serial, e.Time)
	}
}

func (d *WaylandDriver) handleKeyboardModifiers(e client.KeyboardModifiersEvent) {
	d.updateSerial(e.Serial)

	// Update modifiers state - this will be used for subsequent key events
	d.eventState.keyboardEvent.modifiers = e.ModsDepressed

	slog.Debug("keyboard modifiers changed",
		"depressed", e.ModsDepressed,
		"latched", e.ModsLatched,
		"locked", e.ModsLocked,
		"group", e.Group)

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

// releaseKeyboard releases the keyboard interface
func (d *WaylandDriver) releaseKeyboard() {
	if d.keyboard != nil && d.seatVersion >= 3 {
		if err := d.keyboard.Release(); err != nil {
			slog.Error("failed to release keyboard", "error", err)
		}
	}
	d.keyboard = nil
	slog.Info("keyboard interface released")
}

// Wayland keycode to gio KeyCode lookup table - Fixed mappings
var waylandKeycodeMap = map[uint32]gio.KeyCode{
	// Control keys
	1:  gio.KeyCodeEscape,
	14: gio.KeyCodeBackspace,
	15: gio.KeyCodeTab,
	28: gio.KeyCodeEnter,
	57: gio.KeyCodeSpace,

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

	// Top row alphabet keys (QWERTY)
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

	// Middle row alphabet keys (ASDF)
	30: gio.KeyCodeKeyA,
	31: gio.KeyCodeKeyS,
	32: gio.KeyCodeKeyD,
	33: gio.KeyCodeKeyF,
	34: gio.KeyCodeKeyG,
	35: gio.KeyCodeKeyH,
	36: gio.KeyCodeKeyJ,
	37: gio.KeyCodeKeyK,
	38: gio.KeyCodeKeyL,

	// Bottom row alphabet keys (ZXCV)
	44: gio.KeyCodeKeyZ,
	45: gio.KeyCodeKeyX,
	46: gio.KeyCodeKeyC,
	47: gio.KeyCodeKeyV,
	48: gio.KeyCodeKeyB,
	49: gio.KeyCodeKeyN,
	50: gio.KeyCodeKeyM,

	// Arrow keys
	103: gio.KeyCodeArrowUp,
	105: gio.KeyCodeArrowLeft,
	106: gio.KeyCodeArrowRight,
	108: gio.KeyCodeArrowDown,

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

	// Special keys
	58:  gio.KeyCodeCapsLock,
	69:  gio.KeyCodeNumLock,
	70:  gio.KeyCodeScrollLock,
	102: gio.KeyCodeHome,
	104: gio.KeyCodePageUp,
	107: gio.KeyCodeEnd,
	109: gio.KeyCodePageDown,
	110: gio.KeyCodeInsert,
	111: gio.KeyCodeDelete,

	// Punctuation keys
	12: gio.KeyCodeMinus,        // -
	13: gio.KeyCodeEqual,        // =
	26: gio.KeyCodeBracketLeft,  // [
	27: gio.KeyCodeBracketRight, // ]
	39: gio.KeyCodeSemicolon,    // ;
	40: gio.KeyCodeQuote,        // '
	41: gio.KeyCodeBackquote,    // `
	43: gio.KeyCodeBackslash,    // \
	51: gio.KeyCodeComma,        // ,
	52: gio.KeyCodePeriod,       // .
	53: gio.KeyCodeSlash,        // /

	// Numpad keys
	71: gio.KeyCodeNumpad7,
	72: gio.KeyCodeNumpad8,
	73: gio.KeyCodeNumpad9,
	74: gio.KeyCodeNumpadSubtract,
	75: gio.KeyCodeNumpad4,
	76: gio.KeyCodeNumpad5,
	77: gio.KeyCodeNumpad6,
	78: gio.KeyCodeNumpadAdd,
	79: gio.KeyCodeNumpad1,
	80: gio.KeyCodeNumpad2,
	81: gio.KeyCodeNumpad3,
	82: gio.KeyCodeNumpad0,
	83: gio.KeyCodeNumpadDecimal,
	96: gio.KeyCodeNumpadEnter,
	98: gio.KeyCodeNumpadDivide,
	55: gio.KeyCodeNumpadMultiply,

	// Modifier keys
	42:  gio.KeyCodeShiftLeft,    // Left Shift
	54:  gio.KeyCodeShiftRight,   // Right Shift
	29:  gio.KeyCodeControlLeft,  // Left Ctrl
	97:  gio.KeyCodeControlRight, // Right Ctrl
	56:  gio.KeyCodeAltLeft,      // Left Alt
	100: gio.KeyCodeAltRight,     // Right Alt
	125: gio.KeyCodeMetaLeft,     // Left Meta/Super
	126: gio.KeyCodeMetaRight,    // Right Meta/Super
}

func waylandKeycodeToGio(key uint32) gio.KeyCode {
	// Check if the key is in the fixed mapping
	if code, ok := waylandKeycodeMap[key]; ok {
		return code
	}

	// If not found, return a default value (unknown key)
	slog.Warn("unknown Wayland keycode", "key", key)
	return gio.KeyCodeUnknown
}

// Wayland modifier to gio ModifierKey lookup table - Fixed bit positions
var waylandModifierMap = []struct {
	waylandMask uint32
	gioModifier gio.ModifierKey
}{
	{1 << 0, gio.ModifierKeyShift},      // Shift
	{1 << 1, gio.ModifierKeyCapsLock},   // Caps Lock
	{1 << 2, gio.ModifierKeyCtrl},       // Ctrl
	{1 << 3, gio.ModifierKeyAlt},        // Alt
	{1 << 4, gio.ModifierKeyNumLock},    // Num Lock
	{1 << 5, gio.ModifierKeyScrollLock}, // Scroll Lock (if available)
	{1 << 6, gio.ModifierKeyMeta},       // Meta/Super/Windows key
}

func waylandModifiersToGio(modifiers uint32) gio.ModifierKey {
	var gioModifiers gio.ModifierKey
	for _, mod := range waylandModifierMap {
		if modifiers&mod.waylandMask != 0 {
			gioModifiers |= mod.gioModifier
		}
	}
	return gioModifiers
}

// getModifierBitForKey returns the modifier bit for a given key code
func getModifierBitForKey(key uint32) uint32 {
	switch key {
	case 42, 54: // Left Shift, Right Shift
		return 1 << 0
	case 58: // Caps Lock
		return 1 << 1
	case 29, 97: // Left Ctrl, Right Ctrl
		return 1 << 2
	case 56, 100: // Left Alt, Right Alt
		return 1 << 3
	case 69: // Num Lock
		return 1 << 4
	case 70: // Scroll Lock
		return 1 << 5
	case 125, 126: // Left Meta/Super, Right Meta/Super
		return 1 << 6
	default:
		return 0
	}
}
