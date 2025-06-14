package gio

// CursorMode represents the cursor input mode
type CursorMode int

const (
	CursorModeNormal CursorMode = iota
	CursorModeHidden
	CursorModeDisabled
)

func (m CursorMode) String() string {
	switch m {
	case CursorModeNormal:
		return "normal"
	case CursorModeHidden:
		return "hidden"
	case CursorModeDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}

// StandardCursor represents standard system cursor shapes
type StandardCursor int

const (
	ArrowCursor StandardCursor = iota
	IBeamCursor
	CrosshairCursor
	HandCursor
	HResizeCursor
	VResizeCursor
	ResizeCursor
)

func (c StandardCursor) String() string {
	switch c {
	case ArrowCursor:
		return "arrow"
	case IBeamCursor:
		return "ibeam"
	case CrosshairCursor:
		return "crosshair"
	case HandCursor:
		return "hand"
	case HResizeCursor:
		return "hresize"
	case VResizeCursor:
		return "vresize"
	case ResizeCursor:
		return "resize"
	default:
		return "unknown"
	}
}

// Cursor represents a cursor object
type Cursor interface {
	Destroy() error
}
