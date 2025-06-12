package gio

type WindowOptions struct {
	// Title is the title of the window.
	Title string
	// Width is the width of the window in pixels.
	Width int
	// Height is the height of the window in pixels.
	Height int
	// Resizable indicates whether the window can be resized by the user.
	Resizable bool
	// Fullscreen indicates whether the window should start in fullscreen mode.
	Fullscreen bool
	// Visible indicates whether the window should be visible when created.
	Visible bool
}

type Window interface {
	// Show displays the window.
	Show() error
	// Hide hides the window.
	Hide() error
	// Close closes the window.
	Close() error
	// Resize resizes the window to the specified width and height.
	Resize(width, height int) error
	// SetTitle sets the title of the window.
	SetTitle(title string) error
	// SetFullscreen sets the window to fullscreen mode.
	SetFullscreen(fullscreen bool) error
	// SetResizable sets whether the window can be resized by the user.
	SetResizable(resizable bool) error
	// SetVisible sets whether the window is visible.
	SetVisible(visible bool) error
	// GetSize returns the current size of the window.
	GetSize() (width, height int, err error)
	// GetTitle returns the current title of the window.
	GetTitle() (title string, err error)
	// GetFullscreen returns whether the window is currently in fullscreen mode.
	GetFullscreen() (fullscreen bool, err error)
	// GetResizable returns whether the window can be resized by the user.
	GetResizable() (resizable bool, err error)
	// GetVisible returns whether the window is currently visible.
	GetVisible() (visible bool, err error)
	// SetPosition sets the position of the window on the screen.
}
