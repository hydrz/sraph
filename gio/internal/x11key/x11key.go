//go:generate go run gen.go

// x11key contains X11 numeric codes for the keyboard and mouse.
package x11key

import (
	"unicode"

	"github.com/jezek/xgb/xproto"
	"github.com/opensraph/sraph/gio"
)

// These constants come from /usr/include/X11/X.h
const (
	ShiftMask   = 1 << 0
	LockMask    = 1 << 1
	ControlMask = 1 << 2
	Mod1Mask    = 1 << 3
	Mod2Mask    = 1 << 4
	Mod3Mask    = 1 << 5
	Mod4Mask    = 1 << 6
	Mod5Mask    = 1 << 7
	Button1Mask = 1 << 8
	Button2Mask = 1 << 9
	Button3Mask = 1 << 10
	Button4Mask = 1 << 11
	Button5Mask = 1 << 12
)

type KeysymTable struct {
	Table [256][6]uint32

	NumLockMod, ModeSwitchMod, ISOLevel3ShiftMod uint16
}

func (t *KeysymTable) Lookup(detail xproto.Keycode, state uint16) (rune, gio.KeyCode) {
	te := t.Table[detail][0:2]
	if state&t.ModeSwitchMod != 0 {
		te = t.Table[detail][2:4]
	}
	if state&t.ISOLevel3ShiftMod != 0 {
		te = t.Table[detail][4:6]
	}

	// The key event's rune depends on whether the shift key is down.
	unshifted := rune(te[0])
	r := unshifted
	if state&t.NumLockMod != 0 && isKeypad(te[1]) {
		if state&ShiftMask == 0 {
			r = rune(te[1])
		}
	} else if state&ShiftMask != 0 {
		r = rune(te[1])
		// In X11, a zero keysym when shift is down means to use what the
		// keysym is when shift is up.
		if r == 0 {
			r = unshifted
		}
	}

	// The key event's code is independent of whether the shift key is down.
	var c gio.KeyCode
	if 0 <= unshifted && unshifted < 0x80 {
		c = asciiKeycodes[unshifted]
		if state&LockMask != 0 {
			r = unicode.ToUpper(r)
		}
	} else if kk, isKeypad := keypadKeysyms[r]; isKeypad {
		r, c = kk.rune, kk.code
	} else if nuk := nonUnicodeKeycodes[unshifted]; nuk != gio.KeyCodeUnknown {
		r, c = -1, nuk
	} else {
		r = keysymCodePoints[r]
		if state&LockMask != 0 {
			r = unicode.ToUpper(r)
		}
	}

	return r, c
}

func isKeypad(keysym uint32) bool {
	return keysym >= 0xff80 && keysym <= 0xffbd
}

func KeyModifiers(state uint16) (m gio.ModifierKey) {
	if state&ShiftMask != 0 {
		m |= gio.ModifierKeyShift
	}
	if state&ControlMask != 0 {
		m |= gio.ModifierKeyCtrl
	}
	if state&Mod1Mask != 0 {
		m |= gio.ModifierKeyAlt
	}
	if state&Mod4Mask != 0 {
		m |= gio.ModifierKeyMeta
	}
	return m
}

// These constants come from /usr/include/X11/{keysymdef,XF86keysym}.h
const (
	xkISOLeftTab = 0xfe20
	xkBackSpace  = 0xff08
	xkTab        = 0xff09
	xkReturn     = 0xff0d
	xkEscape     = 0xff1b
	xkMultiKey   = 0xff20
	xkHome       = 0xff50
	xkLeft       = 0xff51
	xkUp         = 0xff52
	xkRight      = 0xff53
	xkDown       = 0xff54
	xkPageUp     = 0xff55
	xkPageDown   = 0xff56
	xkEnd        = 0xff57
	xkInsert     = 0xff63
	xkMenu       = 0xff67
	xkHelp       = 0xff6a

	xkNumLock        = 0xff7f
	xkKeypadEnter    = 0xff8d
	xkKeypadHome     = 0xff95
	xkKeypadLeft     = 0xff96
	xkKeypadUp       = 0xff97
	xkKeypadRight    = 0xff98
	xkKeypadDown     = 0xff99
	xkKeypadPageUp   = 0xff9a
	xkKeypadPageDown = 0xff9b
	xkKeypadEnd      = 0xff9c
	xkKeypadInsert   = 0xff9e
	xkKeypadDelete   = 0xff9f
	xkKeypadEqual    = 0xffbd
	xkKeypadMultiply = 0xffaa
	xkKeypadAdd      = 0xffab
	xkKeypadSubtract = 0xffad
	xkKeypadDecimal  = 0xffae
	xkKeypadDivide   = 0xffaf
	xkKeypad0        = 0xffb0
	xkKeypad1        = 0xffb1
	xkKeypad2        = 0xffb2
	xkKeypad3        = 0xffb3
	xkKeypad4        = 0xffb4
	xkKeypad5        = 0xffb5
	xkKeypad6        = 0xffb6
	xkKeypad7        = 0xffb7
	xkKeypad8        = 0xffb8
	xkKeypad9        = 0xffb9

	xkF1       = 0xffbe
	xkF2       = 0xffbf
	xkF3       = 0xffc0
	xkF4       = 0xffc1
	xkF5       = 0xffc2
	xkF6       = 0xffc3
	xkF7       = 0xffc4
	xkF8       = 0xffc5
	xkF9       = 0xffc6
	xkF10      = 0xffc7
	xkF11      = 0xffc8
	xkF12      = 0xffc9
	xkShiftL   = 0xffe1
	xkShiftR   = 0xffe2
	xkControlL = 0xffe3
	xkControlR = 0xffe4
	xkCapsLock = 0xffe5
	xkAltL     = 0xffe9
	xkAltR     = 0xffea
	xkSuperL   = 0xffeb
	xkSuperR   = 0xffec
	xkDelete   = 0xffff

	xf86xkAudioLowerVolume = 0x1008ff11
	xf86xkAudioMute        = 0x1008ff12
	xf86xkAudioRaiseVolume = 0x1008ff13
)

// nonUnicodeKeycodes maps from those xproto.Keysym values (converted to runes)
// that do not correspond to a Unicode code point, such as "Page Up", "F1" or
// "Left Shift", to gio.KeyCode values.
var nonUnicodeKeycodes = map[rune]gio.KeyCode{
	xkISOLeftTab: gio.KeyCodeTab,
	xkBackSpace:  gio.KeyCodeBackspace,
	xkTab:        gio.KeyCodeTab,
	xkReturn:     gio.KeyCodeEnter,
	xkEscape:     gio.KeyCodeEscape,
	xkHome:       gio.KeyCodeHome,
	xkLeft:       gio.KeyCodeArrowLeft,
	xkUp:         gio.KeyCodeArrowUp,
	xkRight:      gio.KeyCodeArrowRight,
	xkDown:       gio.KeyCodeArrowDown,
	xkPageUp:     gio.KeyCodePageUp,
	xkPageDown:   gio.KeyCodePageDown,
	xkEnd:        gio.KeyCodeEnd,
	xkInsert:     gio.KeyCodeInsert,
	xkMenu:       gio.KeyCodeContextMenu,
	xkHelp:       gio.KeyCodeUnknown, // No direct equivalent in current KeyCode
	xkNumLock:    gio.KeyCodeNumLock,
	xkMultiKey:   gio.KeyCodeUnknown, // No direct equivalent in current KeyCode

	xkKeypadEnter:    gio.KeyCodeNumpadEnter,
	xkKeypadHome:     gio.KeyCodeHome,
	xkKeypadLeft:     gio.KeyCodeArrowLeft,
	xkKeypadUp:       gio.KeyCodeArrowUp,
	xkKeypadRight:    gio.KeyCodeArrowRight,
	xkKeypadDown:     gio.KeyCodeArrowDown,
	xkKeypadPageUp:   gio.KeyCodePageUp,
	xkKeypadPageDown: gio.KeyCodePageDown,
	xkKeypadEnd:      gio.KeyCodeEnd,
	xkKeypadInsert:   gio.KeyCodeInsert,
	xkKeypadDelete:   gio.KeyCodeDelete,

	xkF1:  gio.KeyCodeF1,
	xkF2:  gio.KeyCodeF2,
	xkF3:  gio.KeyCodeF3,
	xkF4:  gio.KeyCodeF4,
	xkF5:  gio.KeyCodeF5,
	xkF6:  gio.KeyCodeF6,
	xkF7:  gio.KeyCodeF7,
	xkF8:  gio.KeyCodeF8,
	xkF9:  gio.KeyCodeF9,
	xkF10: gio.KeyCodeF10,
	xkF11: gio.KeyCodeF11,
	xkF12: gio.KeyCodeF12,

	xkShiftL:   gio.KeyCodeShiftLeft,
	xkShiftR:   gio.KeyCodeShiftRight,
	xkControlL: gio.KeyCodeControlLeft,
	xkControlR: gio.KeyCodeControlRight,
	xkCapsLock: gio.KeyCodeCapsLock,
	xkAltL:     gio.KeyCodeAltLeft,
	xkAltR:     gio.KeyCodeAltRight,
	xkSuperL:   gio.KeyCodeMetaLeft,
	xkSuperR:   gio.KeyCodeMetaRight,

	xkDelete: gio.KeyCodeDelete,

	xf86xkAudioRaiseVolume: gio.KeyCodeAudioVolumeUp,
	xf86xkAudioLowerVolume: gio.KeyCodeAudioVolumeDown,
	xf86xkAudioMute:        gio.KeyCodeAudioVolumeMute,
}

// asciiKeycodes maps lower-case ASCII runes to gio.KeyCode values.
var asciiKeycodes = [0x80]gio.KeyCode{
	'a': gio.KeyCodeKeyA,
	'b': gio.KeyCodeKeyB,
	'c': gio.KeyCodeKeyC,
	'd': gio.KeyCodeKeyD,
	'e': gio.KeyCodeKeyE,
	'f': gio.KeyCodeKeyF,
	'g': gio.KeyCodeKeyG,
	'h': gio.KeyCodeKeyH,
	'i': gio.KeyCodeKeyI,
	'j': gio.KeyCodeKeyJ,
	'k': gio.KeyCodeKeyK,
	'l': gio.KeyCodeKeyL,
	'm': gio.KeyCodeKeyM,
	'n': gio.KeyCodeKeyN,
	'o': gio.KeyCodeKeyO,
	'p': gio.KeyCodeKeyP,
	'q': gio.KeyCodeKeyQ,
	'r': gio.KeyCodeKeyR,
	's': gio.KeyCodeKeyS,
	't': gio.KeyCodeKeyT,
	'u': gio.KeyCodeKeyU,
	'v': gio.KeyCodeKeyV,
	'w': gio.KeyCodeKeyW,
	'x': gio.KeyCodeKeyX,
	'y': gio.KeyCodeKeyY,
	'z': gio.KeyCodeKeyZ,

	'1': gio.KeyCodeDigit1,
	'2': gio.KeyCodeDigit2,
	'3': gio.KeyCodeDigit3,
	'4': gio.KeyCodeDigit4,
	'5': gio.KeyCodeDigit5,
	'6': gio.KeyCodeDigit6,
	'7': gio.KeyCodeDigit7,
	'8': gio.KeyCodeDigit8,
	'9': gio.KeyCodeDigit9,
	'0': gio.KeyCodeDigit0,

	' ':  gio.KeyCodeSpace,
	'-':  gio.KeyCodeMinus,
	'=':  gio.KeyCodeEqual,
	'[':  gio.KeyCodeBracketLeft,
	']':  gio.KeyCodeBracketRight,
	'\\': gio.KeyCodeBackslash,
	';':  gio.KeyCodeSemicolon,
	'\'': gio.KeyCodeQuote,
	'`':  gio.KeyCodeBackquote,
	',':  gio.KeyCodeComma,
	'.':  gio.KeyCodePeriod,
	'/':  gio.KeyCodeSlash,
}

type keypadKeysym struct {
	rune rune
	code gio.KeyCode
}

var keypadKeysyms = map[rune]keypadKeysym{
	xkKeypadEqual:    {'=', gio.KeyCodeNumpadEqual},
	xkKeypadMultiply: {'*', gio.KeyCodeNumpadMultiply},
	xkKeypadAdd:      {'+', gio.KeyCodeNumpadAdd},
	xkKeypadSubtract: {'-', gio.KeyCodeNumpadSubtract},
	xkKeypadDecimal:  {'.', gio.KeyCodeNumpadDecimal},
	xkKeypadDivide:   {'/', gio.KeyCodeNumpadDivide},
	xkKeypad0:        {'0', gio.KeyCodeNumpad0},
	xkKeypad1:        {'1', gio.KeyCodeNumpad1},
	xkKeypad2:        {'2', gio.KeyCodeNumpad2},
	xkKeypad3:        {'3', gio.KeyCodeNumpad3},
	xkKeypad4:        {'4', gio.KeyCodeNumpad4},
	xkKeypad5:        {'5', gio.KeyCodeNumpad5},
	xkKeypad6:        {'6', gio.KeyCodeNumpad6},
	xkKeypad7:        {'7', gio.KeyCodeNumpad7},
	xkKeypad8:        {'8', gio.KeyCodeNumpad8},
	xkKeypad9:        {'9', gio.KeyCodeNumpad9},
}
