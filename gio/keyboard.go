package gio

import (
	"time"

	"golang.org/x/text/language"
)

type KeyCode uint16

//go:generate go tool stringer -type=KeyCode -trimprefix=KeyCode -output=keyboard_string.go

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

// KeyLocation represents the location of the key on the keyboard
type KeyLocation uint8

// Key location constants following W3C standard
const (
	KeyLocationStandard KeyLocation = 0x00
	KeyLocationLeft     KeyLocation = 0x01 // Left key (e.g., left Shift, left Control)
	KeyLocationRight    KeyLocation = 0x02 // Right key (e.g., right Shift, right Control)
	KeyLocationNumpad   KeyLocation = 0x03 // Numpad key (e.g., numpad Enter)
)

type ModifierKey uint16

// ModifierKey constants represent the modifier keys
const (
	ModifierKeyNone       ModifierKey = iota      // No modifier key
	ModifierKeyAlt                    = 1 << iota // Alt key
	ModifierKeyAltGraph                           // AltGraph key (right Alt on some keyboards)
	ModifierKeyCapsLock                           // CapsLock key
	ModifierKeyControl                            // Control key
	ModifierKeyFn                                 // Fn key (function key, not standard)
	ModifierKeyFnLock                             // FnLock key (function lock, not standard)
	ModifierKeyHyper                              // Hyper key (not standard, often used in custom keyboards)
	ModifierKeyMeta                               // Meta key (often the Windows key)
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

// KeyboardEvent represents a keyboard event following W3C standard
type KeyboardEvent struct {
	BaseEvent
	Location    KeyLocation  // The location of the key on the keyboard
	ModifierKey ModifierKey  // Bitmask of modifier keys pressed
	Repeat      bool         // Whether the key is being held down
	Locale      language.Tag // The locale identifier
	Code        KeyCode      // The code value of the key pressed
}

// NewKeyboardEvent creates a new keyboard event
func NewKeyboardEvent(code KeyCode, modifierKey ModifierKey, repeat bool) *KeyboardEvent {
	return &KeyboardEvent{
		BaseEvent: BaseEvent{
			eventType: EventTypeKeyboard,
			time:      time.Now(),
		},
		Code:        code,
		ModifierKey: modifierKey,
		Repeat:      repeat,
		Location:    KeyLocationStandard, // Default location
		Locale:      language.Und,        // Default locale (undefined)
	}
}

// The KeyboardEvent interface's key read-only property returns the value of the key pressed by the user, taking into consideration the state of modifier keys such as Shift as well as the keyboard locale and layout.
func (e *KeyboardEvent) Key() string {
	// Check if Shift modifier is active
	isShiftPressed := e.ModifierKey.Contains(ModifierKeyShift)
	isCapsLockActive := e.ModifierKey.Contains(ModifierKeyCapsLock)

	switch e.Code {
	// Number row keys
	case KeyCodeDigit1:
		if isShiftPressed {
			return "!"
		}
		return "1"
	case KeyCodeDigit2:
		if isShiftPressed {
			return "@"
		}
		return "2"
	case KeyCodeDigit3:
		if isShiftPressed {
			return "#"
		}
		return "3"
	case KeyCodeDigit4:
		if isShiftPressed {
			return "$"
		}
		return "4"
	case KeyCodeDigit5:
		if isShiftPressed {
			return "%"
		}
		return "5"
	case KeyCodeDigit6:
		if isShiftPressed {
			return "^"
		}
		return "6"
	case KeyCodeDigit7:
		if isShiftPressed {
			return "&"
		}
		return "7"
	case KeyCodeDigit8:
		if isShiftPressed {
			return "*"
		}
		return "8"
	case KeyCodeDigit9:
		if isShiftPressed {
			return "("
		}
		return "9"
	case KeyCodeDigit0:
		if isShiftPressed {
			return ")"
		}
		return "0"

	// Alphabet keys - consider both Shift and CapsLock
	case KeyCodeKeyA:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "A"
		}
		return "a"
	case KeyCodeKeyB:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "B"
		}
		return "b"
	case KeyCodeKeyC:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "C"
		}
		return "c"
	case KeyCodeKeyD:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "D"
		}
		return "d"
	case KeyCodeKeyE:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "E"
		}
		return "e"
	case KeyCodeKeyF:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "F"
		}
		return "f"
	case KeyCodeKeyG:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "G"
		}
		return "g"
	case KeyCodeKeyH:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "H"
		}
		return "h"
	case KeyCodeKeyI:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "I"
		}
		return "i"
	case KeyCodeKeyJ:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "J"
		}
		return "j"
	case KeyCodeKeyK:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "K"
		}
		return "k"
	case KeyCodeKeyL:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "L"
		}
		return "l"
	case KeyCodeKeyM:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "M"
		}
		return "m"
	case KeyCodeKeyN:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "N"
		}
		return "n"
	case KeyCodeKeyO:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "O"
		}
		return "o"
	case KeyCodeKeyP:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "P"
		}
		return "p"
	case KeyCodeKeyQ:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "Q"
		}
		return "q"
	case KeyCodeKeyR:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "R"
		}
		return "r"
	case KeyCodeKeyS:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "S"
		}
		return "s"
	case KeyCodeKeyT:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "T"
		}
		return "t"
	case KeyCodeKeyU:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "U"
		}
		return "u"
	case KeyCodeKeyV:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "V"
		}
		return "v"
	case KeyCodeKeyW:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "W"
		}
		return "w"
	case KeyCodeKeyX:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "X"
		}
		return "x"
	case KeyCodeKeyY:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "Y"
		}
		return "y"
	case KeyCodeKeyZ:
		if (isShiftPressed && !isCapsLockActive) || (!isShiftPressed && isCapsLockActive) {
			return "Z"
		}
		return "z"

	// Symbol keys
	case KeyCodeMinus:
		if isShiftPressed {
			return "_"
		}
		return "-"
	case KeyCodeEqual:
		if isShiftPressed {
			return "+"
		}
		return "="
	case KeyCodeBracketLeft:
		if isShiftPressed {
			return "{"
		}
		return "["
	case KeyCodeBracketRight:
		if isShiftPressed {
			return "}"
		}
		return "]"
	case KeyCodeBackslash:
		if isShiftPressed {
			return "|"
		}
		return "\\"
	case KeyCodeSemicolon:
		if isShiftPressed {
			return ":"
		}
		return ";"
	case KeyCodeQuote:
		if isShiftPressed {
			return "\""
		}
		return "'"
	case KeyCodeBackquote:
		if isShiftPressed {
			return "~"
		}
		return "`"
	case KeyCodeComma:
		if isShiftPressed {
			return "<"
		}
		return ","
	case KeyCodePeriod:
		if isShiftPressed {
			return ">"
		}
		return "."
	case KeyCodeSlash:
		if isShiftPressed {
			return "?"
		}
		return "/"

	// Special keys - return their standard names
	case KeyCodeSpace:
		return " "
	case KeyCodeTab:
		return "Tab"
	case KeyCodeEnter, KeyCodeNumpadEnter:
		return "Enter"
	case KeyCodeBackspace:
		return "Backspace"
	case KeyCodeDelete:
		return "Delete"
	case KeyCodeEscape:
		return "Escape"
	case KeyCodeArrowLeft:
		return "ArrowLeft"
	case KeyCodeArrowRight:
		return "ArrowRight"
	case KeyCodeArrowUp:
		return "ArrowUp"
	case KeyCodeArrowDown:
		return "ArrowDown"
	case KeyCodeHome:
		return "Home"
	case KeyCodeEnd:
		return "End"
	case KeyCodePageUp:
		return "PageUp"
	case KeyCodePageDown:
		return "PageDown"
	case KeyCodeInsert:
		return "Insert"

	// Function keys
	case KeyCodeF1:
		return "F1"
	case KeyCodeF2:
		return "F2"
	case KeyCodeF3:
		return "F3"
	case KeyCodeF4:
		return "F4"
	case KeyCodeF5:
		return "F5"
	case KeyCodeF6:
		return "F6"
	case KeyCodeF7:
		return "F7"
	case KeyCodeF8:
		return "F8"
	case KeyCodeF9:
		return "F9"
	case KeyCodeF10:
		return "F10"
	case KeyCodeF11:
		return "F11"
	case KeyCodeF12:
		return "F12"

	// Modifier keys
	case KeyCodeShiftLeft, KeyCodeShiftRight:
		return "Shift"
	case KeyCodeControlLeft, KeyCodeControlRight:
		return "Control"
	case KeyCodeAltLeft, KeyCodeAltRight:
		return "Alt"
	case KeyCodeMetaLeft, KeyCodeMetaRight:
		return "Meta"

	// Lock keys
	case KeyCodeCapsLock:
		return "CapsLock"
	case KeyCodeNumLock:
		return "NumLock"
	case KeyCodeScrollLock:
		return "ScrollLock"

	// Numpad keys
	case KeyCodeNumpad0:
		return "0"
	case KeyCodeNumpad1:
		return "1"
	case KeyCodeNumpad2:
		return "2"
	case KeyCodeNumpad3:
		return "3"
	case KeyCodeNumpad4:
		return "4"
	case KeyCodeNumpad5:
		return "5"
	case KeyCodeNumpad6:
		return "6"
	case KeyCodeNumpad7:
		return "7"
	case KeyCodeNumpad8:
		return "8"
	case KeyCodeNumpad9:
		return "9"
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

	// Media and other special keys
	case KeyCodePause:
		return "Pause"
	case KeyCodePrintScreen:
		return "PrintScreen"
	case KeyCodeContextMenu:
		return "ContextMenu"

	default:
		// For unknown keys, return the string representation of the KeyCode
		return e.Code.String()
	}
}
