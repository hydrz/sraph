package gio

import "strings"

//go:generate go tool stringer -type=KeyCode -trimprefix=KeyCode -output=keyboard_string.go

type KeyCode uint16

// KeyCode constants represent the physical key on the keyboard
const (
	KeyCodeUnknown KeyCode = 0x0000 // Unknown key

	// Control keys
	KeyCodeEscape    KeyCode = 0x0001 // Escape
	KeyCodeTab       KeyCode = 0x000F // Tab
	KeyCodeCapsLock  KeyCode = 0x003A // CapsLock
	KeyCodeBackspace KeyCode = 0x000E // Backspace
	KeyCodeEnter     KeyCode = 0x001C // Enter
	KeyCodeSpace     KeyCode = 0x0039 // Space

	// Modifier keys
	KeyCodeShiftLeft    KeyCode = 0x002A // Left Shift
	KeyCodeShiftRight   KeyCode = 0x0036 // Right Shift
	KeyCodeControlLeft  KeyCode = 0x001D // Left Control
	KeyCodeControlRight KeyCode = 0xE01D // Right Control
	KeyCodeAltLeft      KeyCode = 0x0038 // Left Alt
	KeyCodeAltRight     KeyCode = 0xE038 // Right Alt
	KeyCodeMetaLeft     KeyCode = 0xE05B // Left Meta (Windows/Command)
	KeyCodeMetaRight    KeyCode = 0xE05C // Right Meta (Windows/Command)

	// Number row keys
	KeyCodeDigit1 KeyCode = 0x0002 // 1
	KeyCodeDigit2 KeyCode = 0x0003 // 2
	KeyCodeDigit3 KeyCode = 0x0004 // 3
	KeyCodeDigit4 KeyCode = 0x0005 // 4
	KeyCodeDigit5 KeyCode = 0x0006 // 5
	KeyCodeDigit6 KeyCode = 0x0007 // 6
	KeyCodeDigit7 KeyCode = 0x0008 // 7
	KeyCodeDigit8 KeyCode = 0x0009 // 8
	KeyCodeDigit9 KeyCode = 0x000A // 9
	KeyCodeDigit0 KeyCode = 0x000B // 0
	KeyCodeMinus  KeyCode = 0x000C // -
	KeyCodeEqual  KeyCode = 0x000D // =

	// Top row function keys
	KeyCodeF1  KeyCode = 0x003B // F1
	KeyCodeF2  KeyCode = 0x003C // F2
	KeyCodeF3  KeyCode = 0x003D // F3
	KeyCodeF4  KeyCode = 0x003E // F4
	KeyCodeF5  KeyCode = 0x003F // F5
	KeyCodeF6  KeyCode = 0x0040 // F6
	KeyCodeF7  KeyCode = 0x0041 // F7
	KeyCodeF8  KeyCode = 0x0042 // F8
	KeyCodeF9  KeyCode = 0x0043 // F9
	KeyCodeF10 KeyCode = 0x0044 // F10
	KeyCodeF11 KeyCode = 0x0057 // F11
	KeyCodeF12 KeyCode = 0x0058 // F12

	// Alphabet keys
	KeyCodeKeyA KeyCode = 0x001E // A
	KeyCodeKeyB KeyCode = 0x0030 // B
	KeyCodeKeyC KeyCode = 0x002E // C
	KeyCodeKeyD KeyCode = 0x0020 // D
	KeyCodeKeyE KeyCode = 0x0012 // E
	KeyCodeKeyF KeyCode = 0x0021 // F
	KeyCodeKeyG KeyCode = 0x0022 // G
	KeyCodeKeyH KeyCode = 0x0023 // H
	KeyCodeKeyI KeyCode = 0x0017 // I
	KeyCodeKeyJ KeyCode = 0x0024 // J
	KeyCodeKeyK KeyCode = 0x0025 // K
	KeyCodeKeyL KeyCode = 0x0026 // L
	KeyCodeKeyM KeyCode = 0x0032 // M
	KeyCodeKeyN KeyCode = 0x0031 // N
	KeyCodeKeyO KeyCode = 0x0018 // O
	KeyCodeKeyP KeyCode = 0x0019 // P
	KeyCodeKeyQ KeyCode = 0x0010 // Q
	KeyCodeKeyR KeyCode = 0x0013 // R
	KeyCodeKeyS KeyCode = 0x001F // S
	KeyCodeKeyT KeyCode = 0x0014 // T
	KeyCodeKeyU KeyCode = 0x0016 // U
	KeyCodeKeyV KeyCode = 0x002F // V
	KeyCodeKeyW KeyCode = 0x0011 // W
	KeyCodeKeyX KeyCode = 0x002D // X
	KeyCodeKeyY KeyCode = 0x0015 // Y
	KeyCodeKeyZ KeyCode = 0x002C // Z

	// Symbol keys
	KeyCodeBracketLeft  KeyCode = 0x001A // [
	KeyCodeBracketRight KeyCode = 0x001B // ]
	KeyCodeBackslash    KeyCode = 0x002B // \
	KeyCodeSemicolon    KeyCode = 0x0027 // ;
	KeyCodeQuote        KeyCode = 0x0028 // '
	KeyCodeBackquote    KeyCode = 0x0029 // `
	KeyCodeComma        KeyCode = 0x0033 // ,
	KeyCodePeriod       KeyCode = 0x0034 // .
	KeyCodeSlash        KeyCode = 0x0035 // /

	// Navigation keys
	KeyCodeInsert     KeyCode = 0xE052 // Insert
	KeyCodeDelete     KeyCode = 0xE053 // Delete
	KeyCodeHome       KeyCode = 0xE047 // Home
	KeyCodeEnd        KeyCode = 0xE04F // End
	KeyCodePageUp     KeyCode = 0xE049 // PageUp
	KeyCodePageDown   KeyCode = 0xE051 // PageDown
	KeyCodeArrowUp    KeyCode = 0xE048 // Up Arrow
	KeyCodeArrowDown  KeyCode = 0xE050 // Down Arrow
	KeyCodeArrowLeft  KeyCode = 0xE04B // Left Arrow
	KeyCodeArrowRight KeyCode = 0xE04D // Right Arrow

	// Numpad keys
	KeyCodeNumLock        KeyCode = 0xE045 // NumLock
	KeyCodeNumpadDivide   KeyCode = 0xE035 // /
	KeyCodeNumpadMultiply KeyCode = 0x0037 // *
	KeyCodeNumpadSubtract KeyCode = 0x004A // -
	KeyCodeNumpadAdd      KeyCode = 0x004E // +
	KeyCodeNumpadEnter    KeyCode = 0xE01C // Enter
	KeyCodeNumpad1        KeyCode = 0x004F // 1
	KeyCodeNumpad2        KeyCode = 0x0050 // 2
	KeyCodeNumpad3        KeyCode = 0x0051 // 3
	KeyCodeNumpad4        KeyCode = 0x004B // 4
	KeyCodeNumpad5        KeyCode = 0x004C // 5
	KeyCodeNumpad6        KeyCode = 0x004D // 6
	KeyCodeNumpad7        KeyCode = 0x0047 // 7
	KeyCodeNumpad8        KeyCode = 0x0048 // 8
	KeyCodeNumpad9        KeyCode = 0x0049 // 9
	KeyCodeNumpad0        KeyCode = 0x0052 // 0
	KeyCodeNumpadDecimal  KeyCode = 0x0053 // .
	KeyCodeNumpadEqual    KeyCode = 0x0059 // = (Numpad)
	KeyCodeNumpadComma    KeyCode = 0x007E // , (Numpad)

	// Lock keys
	KeyCodeScrollLock KeyCode = 0x0046 // ScrollLock
	KeyCodePause      KeyCode = 0x0045 // Pause/Break

	// Print/Screen keys
	KeyCodePrintScreen KeyCode = 0xE037 // PrintScreen

	// Application keys
	KeyCodeContextMenu KeyCode = 0xE05D // ContextMenu

	// Media keys (common subset)
	KeyCodeMediaTrackPrevious KeyCode = 0xE010 // MediaTrackPrevious
	KeyCodeMediaTrackNext     KeyCode = 0xE019 // MediaTrackNext
	KeyCodeMediaPlayPause     KeyCode = 0xE022 // MediaPlayPause
	KeyCodeMediaStop          KeyCode = 0xE024 // MediaStop
	KeyCodeAudioVolumeMute    KeyCode = 0xE020 // AudioVolumeMute
	KeyCodeAudioVolumeDown    KeyCode = 0xE02E // AudioVolumeDown
	KeyCodeAudioVolumeUp      KeyCode = 0xE030 // AudioVolumeUp

	// Browser keys (common subset)
	KeyCodeBrowserHome      KeyCode = 0xE032 // BrowserHome
	KeyCodeBrowserSearch    KeyCode = 0xE065 // BrowserSearch
	KeyCodeBrowserFavorites KeyCode = 0xE066 // BrowserFavorites
	KeyCodeBrowserRefresh   KeyCode = 0xE067 // BrowserRefresh
	KeyCodeBrowserStop      KeyCode = 0xE068 // BrowserStop
	KeyCodeBrowserForward   KeyCode = 0xE069 // BrowserForward
	KeyCodeBrowserBack      KeyCode = 0xE06A // BrowserBack

	// International keys
	KeyCodeIntlBackslash KeyCode = 0x0056 // IntlBackslash
	KeyCodeIntlRo        KeyCode = 0x0073 // IntlRo
	KeyCodeIntlYen       KeyCode = 0x007D // IntlYen

	// Language keys
	KeyCodeLang1 KeyCode = 0x0072 // Lang1 (Korean/Hanja)
	KeyCodeLang2 KeyCode = 0x0071 // Lang2 (Korean/Han/Yeong)
	KeyCodeLang3 KeyCode = 0x0078 // Lang3
	KeyCodeLang4 KeyCode = 0x0077 // Lang4

	// IME keys
	KeyCodeKanaMode   KeyCode = 0x0070 // KanaMode
	KeyCodeConvert    KeyCode = 0x0079 // Convert
	KeyCodeNonConvert KeyCode = 0x007B // NonConvert

	// Power keys
	KeyCodePower  KeyCode = 0xE05E // Power
	KeyCodeSleep  KeyCode = 0xE05F // Sleep
	KeyCodeWakeUp KeyCode = 0xE063 // WakeUp

	// Application launch keys
	KeyCodeLaunchApp1  KeyCode = 0xE06B // LaunchApp1
	KeyCodeLaunchApp2  KeyCode = 0xE021 // LaunchApp2
	KeyCodeLaunchMail  KeyCode = 0xE06C // LaunchMail
	KeyCodeMediaSelect KeyCode = 0xE06D // MediaSelect
)

type ModifierKey uint16

// ModifierKey constants represent the modifier keys
const (
	ModifierKeyNone       ModifierKey = 1 << iota // No modifier key
	ModifierKeyAlt                                // Alt key (Alt key on Windows/Linux, Option key on macOS)
	ModifierKeyAltGraph                           // AltGraph key (right Alt on some keyboards)
	ModifierKeyCapsLock                           // CapsLock key
	ModifierKeyCtrl                               // Control key (Control key on Windows/Linux/macOS)
	ModifierKeyFn                                 // Fn key (function key, not standard)
	ModifierKeyFnLock                             // FnLock key (function lock, not standard)
	ModifierKeyHyper                              // Hyper key (not standard, often used in custom keyboards)
	ModifierKeyMeta                               // Meta key (Windows key on Windows/Linux, Command key on macOS)
	ModifierKeyNumLock                            // NumLock key
	ModifierKeyScrollLock                         // ScrollLock key
	ModifierKeyShift                              // Shift key
	ModifierKeySuper                              // Super key (often the Windows key)
	ModifierKeySymbol                             // Symbol key (not standard, often used in custom keyboards)
	ModifierKeySymbolLock                         // SymbolLock key (not standard, often used in custom keyboards)
)

// Contains checks if a specific modifier key is active
func (m *ModifierKey) Contains(key ModifierKey) bool {
	return (*m & key) != 0
}

// String returns the string representation of the ModifierKey
func (m ModifierKey) String() string {
	if m == 0 {
		return "None"
	}
	keys := []struct {
		flag ModifierKey
		name string
	}{
		{ModifierKeyAlt, "Alt"},
		{ModifierKeyAltGraph, "AltGraph"},
		{ModifierKeyCapsLock, "CapsLock"},
		{ModifierKeyCtrl, "Ctrl"},
		{ModifierKeyFn, "Fn"},
		{ModifierKeyFnLock, "FnLock"},
		{ModifierKeyHyper, "Hyper"},
		{ModifierKeyMeta, "Meta"},
		{ModifierKeyNumLock, "NumLock"},
		{ModifierKeyScrollLock, "ScrollLock"},
		{ModifierKeyShift, "Shift"},
		{ModifierKeySuper, "Super"},
		{ModifierKeySymbol, "Symbol"},
		{ModifierKeySymbolLock, "SymbolLock"},
	}
	var names []string
	for _, k := range keys {
		if m&k.flag != 0 {
			names = append(names, k.name)
		}
	}
	if len(names) == 0 {
		return "Unknown"
	}
	return strings.Join(names, "+")
}

type KeyState uint8

// KeyboardState constants represent the state of the keyboard
const (
	KeyStateUnknown  KeyState = iota // Unknown state
	KeyStatePressed                  // Key is pressed
	KeyStateReleased                 // Key is released
)

func (ks KeyState) String() string {
	switch ks {
	case KeyStateUnknown:
		return "Unknown"
	case KeyStatePressed:
		return "Pressed"
	case KeyStateReleased:
		return "Released"
	default:
		return "UnknownState"
	}
}

// KeyLocation represents the location of the key on the keyboard
type KeyLocation uint8

// Key location constants following W3C standard
const (
	KeyLocationStandard KeyLocation = 0x00
	KeyLocationLeft     KeyLocation = 0x01 // Left key (e.g., left Shift, left Control)
	KeyLocationRight    KeyLocation = 0x02 // Right key (e.g., right Shift, right Control)
	KeyLocationNumpad   KeyLocation = 0x03 // Numpad key (e.g., numpad Enter)
)

func (kl KeyLocation) String() string {
	switch kl {
	case KeyLocationStandard:
		return "Standard"
	case KeyLocationLeft:
		return "Left"
	case KeyLocationRight:
		return "Right"
	case KeyLocationNumpad:
		return "Numpad"
	default:
		return "Unknown"
	}
}

// Key returns the value of the key pressed by the user
func (e *KeyboardEvent) Key() string {
	return keyToString(e.Code, e.ModifierKey)
}

// keyToString converts a KeyCode to its string representation considering modifiers
func keyToString(code KeyCode, modifiers ModifierKey) string {
	isShift := modifiers.Contains(ModifierKeyShift)
	isCaps := modifiers.Contains(ModifierKeyCapsLock)

	// Symbol and number row
	symbols := map[KeyCode][2]string{
		KeyCodeDigit1: {"1", "!"}, KeyCodeDigit2: {"2", "@"}, KeyCodeDigit3: {"3", "#"},
		KeyCodeDigit4: {"4", "$"}, KeyCodeDigit5: {"5", "%"}, KeyCodeDigit6: {"6", "^"},
		KeyCodeDigit7: {"7", "&"}, KeyCodeDigit8: {"8", "*"}, KeyCodeDigit9: {"9", "("},
		KeyCodeDigit0: {"0", ")"}, KeyCodeMinus: {"-", "_"}, KeyCodeEqual: {"=", "+"},
		KeyCodeBracketLeft: {"[", "{"}, KeyCodeBracketRight: {"]", "}"},
		KeyCodeBackslash: {"\\", "|"}, KeyCodeSemicolon: {";", ":"}, KeyCodeQuote: {"'", "\""},
		KeyCodeBackquote: {"`", "~"}, KeyCodeComma: {",", "<"}, KeyCodePeriod: {".", ">"},
		KeyCodeSlash: {"/", "?"},
	}
	if pair, ok := symbols[code]; ok {
		if isShift {
			return pair[1]
		}
		return pair[0]
	}

	// Alphabet
	if code >= KeyCodeKeyA && code <= KeyCodeKeyZ {
		ch := 'a' + rune(code-KeyCodeKeyA)
		if (isShift && !isCaps) || (!isShift && isCaps) {
			return string(ch - 32) // Uppercase
		}
		return string(ch)
	}

	// Numpad
	if code >= KeyCodeNumpad0 && code <= KeyCodeNumpad9 {
		return string('0' + rune(code-KeyCodeNumpad0))
	}
	switch code {
	case KeyCodeNumpadDecimal:
		return "."
	case KeyCodeNumpadDivide:
		return "/"
	case KeyCodeNumpadMultiply:
		return "*"
	case KeyCodeNumpadSubtract:
		return "-"
	case KeyCodeNumpadAdd:
		return "+"
	case KeyCodeNumpadEqual:
		return "="
	case KeyCodeNumpadComma:
		return ","
	}

	// Special keys
	special := map[KeyCode]string{
		KeyCodeSpace: " ",
		KeyCodeTab:   "Tab",
		KeyCodeEnter: "Enter", KeyCodeNumpadEnter: "Enter",
		KeyCodeBackspace: "Backspace",
		KeyCodeDelete:    "Delete",
		KeyCodeEscape:    "Escape",
		KeyCodeArrowLeft: "ArrowLeft", KeyCodeArrowRight: "ArrowRight",
		KeyCodeArrowUp: "ArrowUp", KeyCodeArrowDown: "ArrowDown",
		KeyCodeHome: "Home", KeyCodeEnd: "End",
		KeyCodePageUp: "PageUp", KeyCodePageDown: "PageDown",
		KeyCodeInsert: "Insert",
		KeyCodeF1:     "F1", KeyCodeF2: "F2", KeyCodeF3: "F3", KeyCodeF4: "F4",
		KeyCodeF5: "F5", KeyCodeF6: "F6", KeyCodeF7: "F7", KeyCodeF8: "F8",
		KeyCodeF9: "F9", KeyCodeF10: "F10", KeyCodeF11: "F11", KeyCodeF12: "F12",
		KeyCodeShiftLeft: "Shift", KeyCodeShiftRight: "Shift",
		KeyCodeControlLeft: "Control", KeyCodeControlRight: "Control",
		KeyCodeAltLeft: "Alt", KeyCodeAltRight: "Alt",
		KeyCodeMetaLeft: "Meta", KeyCodeMetaRight: "Meta",
		KeyCodeCapsLock: "CapsLock", KeyCodeNumLock: "NumLock", KeyCodeScrollLock: "ScrollLock",
		KeyCodePause: "Pause", KeyCodePrintScreen: "PrintScreen", KeyCodeContextMenu: "ContextMenu",
	}
	if name, ok := special[code]; ok {
		return name
	}

	// Fallback to generated string
	return code.String()
}
