package gio

// KeyCode definitions
//
//go:generate go run golang.org/x/tools/cmd/stringer  -type=KeyCode,KeyAction,MouseButton,ModifierKey,JoystickAction --trimprefix=Key,KeyAction,MouseButton,Mod,Joystick --output=keyboard_string.go
type KeyCode int32

const (
	KeyUnknown KeyCode = -1

	// Printable keys
	KeySpace      KeyCode = 32
	KeyApostrophe KeyCode = 39 // '
	KeyComma      KeyCode = 44 // ,
	KeyMinus      KeyCode = 45 // -
	KeyPeriod     KeyCode = 46 // .
	KeySlash      KeyCode = 47 // /

	Key0 KeyCode = 48
	Key1 KeyCode = 49
	Key2 KeyCode = 50
	Key3 KeyCode = 51
	Key4 KeyCode = 52
	Key5 KeyCode = 53
	Key6 KeyCode = 54
	Key7 KeyCode = 55
	Key8 KeyCode = 56
	Key9 KeyCode = 57

	KeySemicolon KeyCode = 59 // ;
	KeyEqual     KeyCode = 61 // =

	KeyA KeyCode = 65
	KeyB KeyCode = 66
	KeyC KeyCode = 67
	KeyD KeyCode = 68
	KeyE KeyCode = 69
	KeyF KeyCode = 70
	KeyG KeyCode = 71
	KeyH KeyCode = 72
	KeyI KeyCode = 73
	KeyJ KeyCode = 74
	KeyK KeyCode = 75
	KeyL KeyCode = 76
	KeyM KeyCode = 77
	KeyN KeyCode = 78
	KeyO KeyCode = 79
	KeyP KeyCode = 80
	KeyQ KeyCode = 81
	KeyR KeyCode = 82
	KeyS KeyCode = 83
	KeyT KeyCode = 84
	KeyU KeyCode = 85
	KeyV KeyCode = 86
	KeyW KeyCode = 87
	KeyX KeyCode = 88
	KeyY KeyCode = 89
	KeyZ KeyCode = 90

	KeyLeftBracket  KeyCode = 91 // [
	KeyBackslash    KeyCode = 92 // \
	KeyRightBracket KeyCode = 93 // ]
	KeyGraveAccent  KeyCode = 96 // `

	KeyWorld1 KeyCode = 161 // non-US #1
	KeyWorld2 KeyCode = 162 // non-US #2

	// Function keys
	KeyEscape      KeyCode = 256
	KeyEnter       KeyCode = 257
	KeyTab         KeyCode = 258
	KeyBackspace   KeyCode = 259
	KeyInsert      KeyCode = 260
	KeyDelete      KeyCode = 261
	KeyRight       KeyCode = 262
	KeyLeft        KeyCode = 263
	KeyDown        KeyCode = 264
	KeyUp          KeyCode = 265
	KeyPageUp      KeyCode = 266
	KeyPageDown    KeyCode = 267
	KeyHome        KeyCode = 268
	KeyEnd         KeyCode = 269
	KeyCapsLock    KeyCode = 280
	KeyScrollLock  KeyCode = 281
	KeyNumLock     KeyCode = 282
	KeyPrintScreen KeyCode = 283
	KeyPause       KeyCode = 284

	KeyF1  KeyCode = 290
	KeyF2  KeyCode = 291
	KeyF3  KeyCode = 292
	KeyF4  KeyCode = 293
	KeyF5  KeyCode = 294
	KeyF6  KeyCode = 295
	KeyF7  KeyCode = 296
	KeyF8  KeyCode = 297
	KeyF9  KeyCode = 298
	KeyF10 KeyCode = 299
	KeyF11 KeyCode = 300
	KeyF12 KeyCode = 301
	KeyF13 KeyCode = 302
	KeyF14 KeyCode = 303
	KeyF15 KeyCode = 304
	KeyF16 KeyCode = 305
	KeyF17 KeyCode = 306
	KeyF18 KeyCode = 307
	KeyF19 KeyCode = 308
	KeyF20 KeyCode = 309
	KeyF21 KeyCode = 310
	KeyF22 KeyCode = 311
	KeyF23 KeyCode = 312
	KeyF24 KeyCode = 313
	KeyF25 KeyCode = 314

	// Keypad
	KeyKp0        KeyCode = 320
	KeyKp1        KeyCode = 321
	KeyKp2        KeyCode = 322
	KeyKp3        KeyCode = 323
	KeyKp4        KeyCode = 324
	KeyKp5        KeyCode = 325
	KeyKp6        KeyCode = 326
	KeyKp7        KeyCode = 327
	KeyKp8        KeyCode = 328
	KeyKp9        KeyCode = 329
	KeyKpDecimal  KeyCode = 330
	KeyKpDivide   KeyCode = 331
	KeyKpMultiply KeyCode = 332
	KeyKpSubtract KeyCode = 333
	KeyKpAdd      KeyCode = 334
	KeyKpEnter    KeyCode = 335
	KeyKpEqual    KeyCode = 336

	// Modifier keys
	KeyLeftShift    KeyCode = 340
	KeyLeftControl  KeyCode = 341
	KeyLeftAlt      KeyCode = 342
	KeyLeftSuper    KeyCode = 343
	KeyRightShift   KeyCode = 344
	KeyRightControl KeyCode = 345
	KeyRightAlt     KeyCode = 346
	KeyRightSuper   KeyCode = 347
	KeyMenu         KeyCode = 348
)

// Mouse button definitions
type MouseButton int32

const (
	MouseButtonLeft   MouseButton = 0
	MouseButtonRight  MouseButton = 1
	MouseButtonMiddle MouseButton = 2
	MouseButton1      MouseButton = MouseButtonLeft
	MouseButton2      MouseButton = MouseButtonRight
	MouseButton3      MouseButton = MouseButtonMiddle
	MouseButton4      MouseButton = 3
	MouseButton5      MouseButton = 4
	MouseButton6      MouseButton = 5
	MouseButton7      MouseButton = 6
	MouseButton8      MouseButton = 7
)

// KeyAction definitions
type KeyAction int32

const (
	KeyActionRelease KeyAction = 0
	KeyActionPress   KeyAction = 1
	KeyActionRepeat  KeyAction = 2
)

// Modifier key definitions
//

type ModifierKey int32

const (
	ModShift    ModifierKey = 0x0001
	ModControl  ModifierKey = 0x0002
	ModAlt      ModifierKey = 0x0004
	ModSuper    ModifierKey = 0x0008
	ModCapsLock ModifierKey = 0x0010
	ModNumLock  ModifierKey = 0x0020
)

// Joystick action definitions
type JoystickAction int32

const (
	JoystickConnected    JoystickAction = 0x00040001
	JoystickDisconnected JoystickAction = 0x00040002
)
