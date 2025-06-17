package wayland

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"syscall"
	"time"

	"github.com/opensraph/sraph/gio"
	"github.com/rajveermalviya/go-wayland/wayland/client"
)

// ClipboardManager handles Wayland clipboard operations
type ClipboardManager struct {
	driver            *WaylandDriver
	dataDevice        *client.DataDevice
	dataDeviceManager *client.DataDeviceManager
	currentSource     *client.DataSource
	selection         *gio.DataTransfer

	// Callback for incoming clipboard data
	onDataReceived func(*gio.DataTransfer)
}

// NewClipboardManager creates a new clipboard manager
func NewClipboardManager(driver *WaylandDriver) *ClipboardManager {
	return &ClipboardManager{
		driver: driver,
	}
}

// Initialize sets up the clipboard manager with Wayland protocols
func (cm *ClipboardManager) Initialize() error {
	// Use the data device manager from the driver
	cm.dataDeviceManager = cm.driver.dataDeviceManager
	if cm.dataDeviceManager == nil {
		return fmt.Errorf("data device manager not available")
	}

	// Use the data device from the driver
	cm.dataDevice = cm.driver.dataDevice
	if cm.dataDevice != nil {
		cm.setupDataDeviceHandlers()
	}

	return nil
}

// setupDataDeviceHandlers configures event handlers for data device
func (cm *ClipboardManager) setupDataDeviceHandlers() {
	if cm.dataDevice == nil {
		return
	}

	// Handle data offers (incoming clipboard data)
	cm.dataDevice.SetDataOfferHandler(func(event client.DataDeviceDataOfferEvent) {
		cm.handleDataOffer(event.Id)
	})

	// Handle selection changes
	cm.dataDevice.SetSelectionHandler(func(event client.DataDeviceSelectionEvent) {
		cm.handleSelection(event.Id)
	})

	// Handle enter events for drag and drop
	cm.dataDevice.SetEnterHandler(func(event client.DataDeviceEnterEvent) {
		cm.handleDragEnter(event)
	})

	// Handle leave events for drag and drop
	cm.dataDevice.SetLeaveHandler(func(event client.DataDeviceLeaveEvent) {
		cm.handleDragLeave()
	})

	// Handle motion events for drag and drop
	cm.dataDevice.SetMotionHandler(func(event client.DataDeviceMotionEvent) {
		cm.handleDragMotion(event)
	})

	// Handle drop events
	cm.dataDevice.SetDropHandler(func(event client.DataDeviceDropEvent) {
		cm.handleDrop()
	})
}

// SetClipboardData sets data to the system clipboard
func (cm *ClipboardManager) SetClipboardData(data *gio.DataTransfer) error {
	if cm.dataDeviceManager == nil {
		return fmt.Errorf("data device manager not initialized")
	}

	// Create data source
	dataSource, err := cm.dataDeviceManager.CreateDataSource()
	if err != nil {
		return fmt.Errorf("failed to create data source: %w", err)
	}

	// Store reference to current source
	if cm.currentSource != nil {
		cm.currentSource.Destroy()
	}
	cm.currentSource = dataSource
	cm.selection = data

	// Offer MIME types
	types := data.GetTypes()
	for _, mimeType := range types {
		if err := dataSource.Offer(mimeType); err != nil {
			slog.Error("failed to offer MIME type", "mimeType", mimeType, "error", err)
		}
	}

	// Set up send handler
	dataSource.SetSendHandler(func(event client.DataSourceSendEvent) {
		cm.handleSendData(event.MimeType, event.Fd)
	})

	// Set up cancelled handler
	dataSource.SetCancelledHandler(func(event client.DataSourceCancelledEvent) {
		if cm.currentSource == dataSource {
			cm.currentSource = nil
			cm.selection = nil
		}
		dataSource.Destroy()
	})

	// Set selection
	if cm.dataDevice != nil {
		serial := cm.getLatestSerial() // Would need to track latest serial from input events
		if err := cm.dataDevice.SetSelection(dataSource, serial); err != nil {
			dataSource.Destroy()
			return fmt.Errorf("failed to set selection: %w", err)
		}
	}

	return nil
}

// GetClipboardData requests data from the system clipboard
func (cm *ClipboardManager) GetClipboardData(callback func(*gio.DataTransfer)) error {
	if cm.dataDevice == nil {
		return fmt.Errorf("data device not initialized")
	}

	cm.onDataReceived = callback

	// Request selection data - this will trigger data offer events
	// The actual data retrieval happens in the event handlers
	return nil
}

// handleDataOffer processes incoming data offers
func (cm *ClipboardManager) handleDataOffer(dataOffer *client.DataOffer) {
	if dataOffer == nil {
		return
	}

	mimeTypes := make([]string, 0)

	// Set up offer handler to collect MIME types
	dataOffer.SetOfferHandler(func(event client.DataOfferOfferEvent) {
		mimeTypes = append(mimeTypes, event.MimeType)
	})

	// Wait for all offers to be received (simplified)
	time.Sleep(10 * time.Millisecond)

	// Create data transfer with available types
	dataTransfer := gio.NewDataTransfer()

	// Request data for each MIME type
	for _, mimeType := range mimeTypes {
		go cm.requestDataForType(dataOffer, mimeType, dataTransfer)
	}
}

// requestDataForType requests specific MIME type data
func (cm *ClipboardManager) requestDataForType(dataOffer *client.DataOffer, mimeType string, dataTransfer *gio.DataTransfer) {
	// Create pipe for data transfer
	r, w, err := os.Pipe()
	if err != nil {
		slog.Error("failed to create pipe", "error", err)
		return
	}
	defer r.Close()
	defer w.Close()

	// Receive data
	if err := dataOffer.Receive(mimeType, int(w.Fd())); err != nil {
		slog.Error("failed to receive data", "error", err)
		return
	}

	// Close write end so read will complete
	w.Close()

	// Read all data
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		slog.Error("failed to read clipboard data", "error", err)
		return
	}

	// Store data in transfer object
	if err := dataTransfer.SetData(mimeType, buf.String()); err != nil {
		slog.Error("failed to set clipboard data", "error", err)
		return
	}

	// Notify callback if this is the first/main type
	if cm.onDataReceived != nil && mimeType == "text/plain" {
		go cm.onDataReceived(dataTransfer)
	}
}

// handleSelection processes selection changes
func (cm *ClipboardManager) handleSelection(dataOffer *client.DataOffer) {
	if dataOffer != nil {
		cm.handleDataOffer(dataOffer)
	}
}

// handleSendData sends clipboard data when requested
func (cm *ClipboardManager) handleSendData(mimeType string, fd int) {
	defer syscall.Close(fd)

	if cm.selection == nil {
		return
	}

	// Get data for the requested MIME type
	data, exists := cm.selection.GetData(mimeType)
	if !exists {
		slog.Warn("requested MIME type not available", "mimeType", mimeType)
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
	file := os.NewFile(uintptr(fd), "clipboard")
	defer file.Close()

	if _, err := file.Write(dataBytes); err != nil {
		slog.Error("failed to write clipboard data", "error", err)
	}
}

// getLatestSerial returns the latest serial number from input events
func (cm *ClipboardManager) getLatestSerial() uint32 {
	return cm.driver.getLatestSerial()
}

// Drag and drop handlers
func (cm *ClipboardManager) handleDragEnter(event client.DataDeviceEnterEvent) {
	// Handle drag enter - would create drag event
	slog.Debug("drag enter", "surface", event.Surface)
}

func (cm *ClipboardManager) handleDragLeave() {
	// Handle drag leave
	slog.Debug("drag leave")
}

func (cm *ClipboardManager) handleDragMotion(event client.DataDeviceMotionEvent) {
	// Handle drag motion
	slog.Debug("drag motion", "x", event.X, "y", event.Y)
}

func (cm *ClipboardManager) handleDrop() {
	// Handle drop event
	slog.Debug("drop occurred")
}

// Cleanup releases clipboard resources
func (cm *ClipboardManager) Cleanup() {
	if cm.currentSource != nil {
		cm.currentSource.Destroy()
		cm.currentSource = nil
	}

	if cm.dataDevice != nil {
		// Data device is managed by the driver
		cm.dataDevice = nil
	}

	if cm.dataDeviceManager != nil {
		cm.dataDeviceManager.Destroy()
		cm.dataDeviceManager = nil
	}
}

// SetClipboardData sets data to the system clipboard (window method)
func (w *WaylandWindow) SetClipboardData(data *gio.DataTransfer) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.IsClosed() {
		return gio.ErrWindowNotInitialized
	}

	// Use driver's clipboard manager
	if w.driver.clipboardManager != nil {
		return w.driver.clipboardManager.SetClipboardData(data)
	}

	// Fallback to local storage
	w.clipboardData = data
	slog.Debug("clipboard data set (local storage fallback)")
	return nil
}

// GetClipboardData requests data from the system clipboard (window method)
func (w *WaylandWindow) GetClipboardData() error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.IsClosed() {
		return gio.ErrWindowNotInitialized
	}

	// Use driver's clipboard manager
	if w.driver.clipboardManager != nil {
		return w.driver.clipboardManager.GetClipboardData(func(data *gio.DataTransfer) {
			clipboardEvent := gio.NewClipboardEvent()
			clipboardEvent.Data = data

			if err := w.Publish(clipboardEvent); err != nil {
				slog.Error("failed to publish clipboard event", "error", err)
			}
		})
	}

	// Fallback to local storage
	if w.clipboardData != nil {
		clipboardEvent := gio.NewClipboardEvent()
		clipboardEvent.Data = w.clipboardData

		go func() {
			if err := w.Publish(clipboardEvent); err != nil {
				slog.Error("failed to publish clipboard event", "error", err)
			}
		}()
	}

	return nil
}

// handleDataOffer handles incoming data offers (drag and drop)
func (w *WaylandWindow) handleDataOffer(mimeTypes []string) {
	if w.IsClosed() {
		return
	}

	dragEvent := gio.NewDragEvent()

	// Create data transfer object with available MIME types
	for _, mimeType := range mimeTypes {
		dragEvent.Data.SetData(mimeType, nil) // Data will be retrieved on drop
	}

	if err := w.Publish(dragEvent); err != nil {
		slog.Error("failed to publish drag event", "error", err)
	}
}
