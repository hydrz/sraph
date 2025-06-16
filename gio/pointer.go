package gio

type PointerType uint8

// PointerType constants represent the type of pointer
const (
	PointerTypeUnknown PointerType = iota // Unknown pointer type
	PointerTypeMouse                      // Mouse pointer
	PointerTypePen                        // Pen pointer
	PointerTypeTouch                      // Touch pointer
)
