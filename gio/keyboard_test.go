package gio

import "testing"

func TestKeyToString(t *testing.T) {
	testCases := []struct {
		name      string
		code      KeyCode
		modifiers ModifierKey
		expected  string
	}{
		// Basic alphabet keys without modifiers
		{"Key A lowercase", KeyCodeKeyA, ModifierKeyNone, "a"},
		{"Key A uppercase with shift", KeyCodeKeyA, ModifierKeyShift, "A"},
		{"Key A uppercase with caps lock", KeyCodeKeyA, ModifierKeyCapsLock, "A"},
		{"Key A lowercase with shift and caps lock", KeyCodeKeyA, ModifierKeyShift | ModifierKeyCapsLock, "a"},

		// Number row
		{"Digit 1", KeyCodeDigit1, ModifierKeyNone, "1"},
		{"Digit 1 with shift", KeyCodeDigit1, ModifierKeyShift, "!"},
		{"Digit 5", KeyCodeDigit5, ModifierKeyNone, "5"},
		{"Digit 5 with shift", KeyCodeDigit5, ModifierKeyShift, "%"},

		// Symbol keys
		{"Minus", KeyCodeMinus, ModifierKeyNone, "-"},
		{"Minus with shift", KeyCodeMinus, ModifierKeyShift, "_"},
		{"Equal", KeyCodeEqual, ModifierKeyNone, "="},
		{"Equal with shift", KeyCodeEqual, ModifierKeyShift, "+"},

		// Special keys
		{"Space", KeyCodeSpace, ModifierKeyNone, " "},
		{"Tab", KeyCodeTab, ModifierKeyNone, "Tab"},
		{"Enter", KeyCodeEnter, ModifierKeyNone, "Enter"},
		{"Backspace", KeyCodeBackspace, ModifierKeyNone, "Backspace"},

		// Function keys
		{"F1", KeyCodeF1, ModifierKeyNone, "F1"},
		{"F12", KeyCodeF12, ModifierKeyNone, "F12"},

		// Arrow keys
		{"Arrow Left", KeyCodeArrowLeft, ModifierKeyNone, "ArrowLeft"},
		{"Arrow Up", KeyCodeArrowUp, ModifierKeyNone, "ArrowUp"},

		// Modifier keys
		{"Shift Left", KeyCodeShiftLeft, ModifierKeyNone, "Shift"},
		{"Control Right", KeyCodeControlRight, ModifierKeyNone, "Control"},
		{"Alt Left", KeyCodeAltLeft, ModifierKeyNone, "Alt"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := keyToString(tc.code, tc.modifiers)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestModifierKeyContains(t *testing.T) {
	combined := ModifierKeyShift | ModifierKeyCtrl | ModifierKeyAlt

	testCases := []struct {
		name     string
		modifier ModifierKey
		check    ModifierKey
		expected bool
	}{
		{"Contains Shift", combined, ModifierKeyShift, true},
		{"Contains Ctrl", combined, ModifierKeyCtrl, true},
		{"Contains Alt", combined, ModifierKeyAlt, true},
		{"Does not contain Meta", combined, ModifierKeyMeta, false},
		{"None contains nothing", ModifierKeyNone, ModifierKeyShift, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.modifier.Contains(tc.check)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestKeyStateString(t *testing.T) {
	testCases := []struct {
		state    KeyState
		expected string
	}{
		{KeyStateUnknown, "Unknown"},
		{KeyStatePressed, "Pressed"},
		{KeyStateReleased, "Released"},
	}

	for _, tc := range testCases {
		result := tc.state.String()
		if result != tc.expected {
			t.Errorf("Expected '%s', got '%s'", tc.expected, result)
		}
	}
}

func TestKeyLocationString(t *testing.T) {
	testCases := []struct {
		location KeyLocation
		expected string
	}{
		{KeyLocationStandard, "Standard"},
		{KeyLocationLeft, "Left"},
		{KeyLocationRight, "Right"},
		{KeyLocationNumpad, "Numpad"},
	}

	for _, tc := range testCases {
		result := tc.location.String()
		if result != tc.expected {
			t.Errorf("Expected '%s', got '%s'", tc.expected, result)
		}
	}
}

func TestKeyboardEventKey(t *testing.T) {
	event := &KeyboardEvent{
		Code:        KeyCodeKeyA,
		ModifierKey: ModifierKeyShift,
	}

	result := event.Key()
	expected := "A"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
