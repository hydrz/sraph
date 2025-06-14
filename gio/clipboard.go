package gio

// ClipboardFormat represents the format of clipboard data
type ClipboardFormat int

const (
	ClipboardFormatText ClipboardFormat = iota
	ClipboardFormatImage
	ClipboardFormatHTML
	ClipboardFormatRTF
	ClipboardFormatFiles
)

func (f ClipboardFormat) String() string {
	switch f {
	case ClipboardFormatText:
		return "text/plain"
	case ClipboardFormatImage:
		return "image/png"
	case ClipboardFormatHTML:
		return "text/html"
	case ClipboardFormatRTF:
		return "text/rtf"
	case ClipboardFormatFiles:
		return "application/x-file-list"
	default:
		return "unknown"
	}
}

// ClipboardData represents clipboard data
type ClipboardData struct {
	Format ClipboardFormat
	Data   []byte
	Text   string   // For text format
	Files  []string // For files format
}
