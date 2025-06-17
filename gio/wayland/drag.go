package wayland

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/opensraph/sraph/gio"
	"github.com/rajveermalviya/go-wayland/wayland/client"
)

// DragOperation represents an ongoing drag operation
type DragOperation struct {
	source  *client.DataSource
	data    *gio.DataTransfer
	surface *client.Surface
	serial  uint32
	window  *WaylandWindow
}

// StartDrag initiates a drag operation
func (w *WaylandWindow) StartDrag(data *gio.DataTransfer, serial uint32) (*DragOperation, error) {
	if w.driver.dataDeviceManager == nil {
		return nil, fmt.Errorf("data device manager not available")
	}

	if w.driver.dataDevice == nil {
		return nil, fmt.Errorf("data device not available")
	}

	// Create data source
	dataSource, err := w.driver.dataDeviceManager.CreateDataSource()
	if err != nil {
		return nil, fmt.Errorf("failed to create data source: %w", err)
	}

	// Offer MIME types
	types := data.GetTypes()
	for _, mimeType := range types {
		if err := dataSource.Offer(mimeType); err != nil {
			slog.Error("failed to offer MIME type", "mimeType", mimeType, "error", err)
		}
	}

	// Create drag operation
	dragOp := &DragOperation{
		source:  dataSource,
		data:    data,
		surface: w.surface,
		serial:  serial,
		window:  w,
	}

	// Set up data source handlers
	dataSource.SetSendHandler(func(event client.DataSourceSendEvent) {
		dragOp.handleSendData(event.MimeType, event.Fd)
	})

	dataSource.SetCancelledHandler(func(event client.DataSourceCancelledEvent) {
		dragOp.cleanup()
	})

	dataSource.SetDndDropPerformedHandler(func(event client.DataSourceDndDropPerformedEvent) {
		slog.Debug("drag drop performed")
	})

	dataSource.SetDndFinishedHandler(func(event client.DataSourceDndFinishedEvent) {
		dragOp.cleanup()
	})

	// Start drag and drop
	err = w.driver.dataDevice.StartDrag(dataSource, w.surface, nil, serial)
	if err != nil {
		dataSource.Destroy()
		return nil, fmt.Errorf("failed to start drag: %w", err)
	}

	return dragOp, nil
}

// handleSendData sends drag data when requested
func (do *DragOperation) handleSendData(mimeType string, fd int) {
	defer func() {
		if err := os.NewFile(uintptr(fd), "drag-data").Close(); err != nil {
			slog.Error("failed to close drag data fd", "error", err)
		}
	}()

	if do.data == nil {
		return
	}

	// Get data for the requested MIME type
	data, exists := do.data.GetData(mimeType)
	if !exists {
		slog.Warn("requested MIME type not available for drag", "mimeType", mimeType)
		return
	}

	// Convert data to bytes
	var dataBytes []byte
	switch v := data.(type) {
	case string:
		dataBytes = []byte(v)
	case []byte:
		dataBytes = v
	default:
		dataBytes = []byte(fmt.Sprintf("%v", v))
	}

	// Write data to file descriptor
	file := os.NewFile(uintptr(fd), "drag-data")
	if _, err := file.Write(dataBytes); err != nil {
		slog.Error("failed to write drag data", "error", err)
	}
}

// cleanup cleans up the drag operation
func (do *DragOperation) cleanup() {
	if do.source != nil {
		do.source.Destroy()
		do.source = nil
	}
}

// HandleDragEnter processes drag enter events
func (w *WaylandWindow) HandleDragEnter(dataOffer *client.DataOffer, x, y float64, surface *client.Surface) {
	if w.IsClosed() {
		return
	}

	dragEvent := gio.NewDragEvent()
	dragEvent.Data = gio.NewDataTransfer()

	// Set up data offer handlers to collect MIME types
	mimeTypes := make([]string, 0)
	dataOffer.SetOfferHandler(func(event client.DataOfferOfferEvent) {
		mimeTypes = append(mimeTypes, event.MimeType)
		// Add placeholder data - actual data retrieved on drop
		dragEvent.Data.SetData(event.MimeType, nil)
	})

	// Publish drag event
	if err := w.Publish(dragEvent); err != nil {
		slog.Error("failed to publish drag enter event", "error", err)
	}
}

// HandleDragMotion processes drag motion events
func (w *WaylandWindow) HandleDragMotion(x, y float64, time uint32) {
	if w.IsClosed() {
		return
	}

	// Create pointer event for drag motion
	pointerEvent := gio.NewPointerEvent()
	pointerEvent.PointerType = gio.PointerTypeMouse
	pointerEvent.Position.X = int(x)
	pointerEvent.Position.Y = int(y)

	if err := w.Publish(pointerEvent); err != nil {
		slog.Error("failed to publish drag motion event", "error", err)
	}
}

// HandleDragLeave processes drag leave events
func (w *WaylandWindow) HandleDragLeave() {
	if w.IsClosed() {
		return
	}

	// Create empty drag event to signal leave
	dragEvent := gio.NewDragEvent()
	dragEvent.Data = gio.NewDataTransfer()

	if err := w.Publish(dragEvent); err != nil {
		slog.Error("failed to publish drag leave event", "error", err)
	}
}

// HandleDrop processes drop events
func (w *WaylandWindow) HandleDrop(dataOffer *client.DataOffer) {
	if w.IsClosed() {
		return
	}

	dragEvent := gio.NewDragEvent()
	dragEvent.Data = gio.NewDataTransfer()

	// Accept the drop and retrieve data
	// This is a simplified version - full implementation would handle multiple MIME types
	mimeTypes := []string{"text/plain", "text/uri-list"}

	for _, mimeType := range mimeTypes {
		go w.retrieveDropData(dataOffer, mimeType, dragEvent.Data)
	}

	// Publish drop event
	if err := w.Publish(dragEvent); err != nil {
		slog.Error("failed to publish drop event", "error", err)
	}
}

// retrieveDropData retrieves data for a specific MIME type from drop
func (w *WaylandWindow) retrieveDropData(dataOffer *client.DataOffer, mimeType string, dataTransfer *gio.DataTransfer) {
	// Create pipe for data transfer
	r, writer, err := os.Pipe()
	if err != nil {
		slog.Error("failed to create pipe for drop data", "error", err)
		return
	}
	defer r.Close()
	defer writer.Close()

	// Receive data
	if err := dataOffer.Receive(mimeType, int(writer.Fd())); err != nil {
		slog.Error("failed to receive drop data", "error", err)
		return
	}

	// Close write end so read will complete
	writer.Close()

	// Read data
	buf := make([]byte, 4096)
	n, err := r.Read(buf)
	if err != nil && err.Error() != "EOF" {
		slog.Error("failed to read drop data", "error", err)
		return
	}

	// Store data
	if n > 0 {
		dataTransfer.SetData(mimeType, string(buf[:n]))
	}
}
