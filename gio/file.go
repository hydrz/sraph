package gio

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

// DropEffect represents the type of drag-and-drop operation
type DropEffect uint8

const (
	DropEffectNone DropEffect = iota // No effect
	DropEffectCopy
	DropEffectMove
	DropEffectLink
)

// DataTransferItemKind represents the kind of data transfer item
type DataTransferItemKind uint8

const (
	DataTransferItemKindUnknown DataTransferItemKind = iota
	DataTransferItemKindString
	DataTransferItemKindFile
)

// DataTransferItem represents a single item in a data transfer operation
type DataTransferItem struct {
	Kind     DataTransferItemKind // Kind of item (string or file)
	MimeType string               // MIME type
	data     interface{}
}

// GetAsString returns the string representation of the item
func (item *DataTransferItem) GetAsString() (string, error) {

	if item.Kind != DataTransferItemKindString {
		return "", fmt.Errorf("item is not a string type")
	}
	if str, ok := item.data.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("failed to convert data to string")
}

// GetAsFile returns the file representation of the item
func (item *DataTransferItem) GetAsFile() (*os.File, error) {
	if item.Kind != DataTransferItemKindFile {
		return nil, fmt.Errorf("item is not a file type")
	}
	if file, ok := item.data.(*os.File); ok {
		return file, nil
	}
	return nil, fmt.Errorf("failed to convert data to file")
}

// FileList represents a list of files
type FileList struct {
	files []*os.File
}

// NewFileList creates a new FileList
func NewFileList() *FileList {
	return &FileList{
		files: make([]*os.File, 0),
	}
}

// Length returns the number of files in the list
func (fl *FileList) Length() int {
	return len(fl.files)
}

// Item returns the file at the specified index
func (fl *FileList) Item(index int) *os.File {
	if index < 0 || index >= len(fl.files) {
		return nil
	}
	return fl.files[index]
}

// Add adds a file to the list
func (fl *FileList) Add(file *os.File) {
	fl.files = append(fl.files, file)
}

// DataTransfer represents a data transfer operation
type DataTransfer struct {
	dropEffect    DropEffect
	effectAllowed DropEffect
	items         []DataTransferItem
	files         *FileList
	types         []string
	dataMap       map[string]interface{}
}

// NewDataTransfer creates a new DataTransfer instance
func NewDataTransfer() *DataTransfer {
	return &DataTransfer{
		dropEffect:    DropEffectNone,
		effectAllowed: DropEffectCopy,
		items:         make([]DataTransferItem, 0),
		files:         NewFileList(),
		types:         make([]string, 0),
		dataMap:       make(map[string]interface{}),
	}
}

// SetDropEffect sets the drop effect
func (dt *DataTransfer) SetDropEffect(effect DropEffect) {
	dt.dropEffect = effect
}

// GetDropEffect returns the current drop effect
func (dt *DataTransfer) GetDropEffect() DropEffect {
	return dt.dropEffect
}

// SetEffectAllowed sets the allowed effects
func (dt *DataTransfer) SetEffectAllowed(effect DropEffect) {
	dt.effectAllowed = effect
}

// GetEffectAllowed returns the allowed effects
func (dt *DataTransfer) GetEffectAllowed() DropEffect {
	return dt.effectAllowed
}

// SetData sets data for a specific MIME type
func (dt *DataTransfer) SetData(mimeType string, data interface{}) error {
	// Determine item kind based on data type
	var kind DataTransferItemKind
	switch data.(type) {
	case string:
		kind = DataTransferItemKindString
	case *os.File:
		kind = DataTransferItemKindFile
		// Add to files list if it's a file
		if file, ok := data.(*os.File); ok {
			dt.files.Add(file)
		}
	default:
		kind = DataTransferItemKindString
		data = fmt.Sprintf("%v", data)
	}

	// Create and add item
	item := DataTransferItem{
		Kind:     kind,
		MimeType: mimeType,
		data:     data,
	}
	dt.items = append(dt.items, item)

	// Add to types if not already present
	if !dt.containsType(mimeType) {
		dt.types = append(dt.types, mimeType)
	}

	// Store in data map for quick access
	dt.dataMap[mimeType] = data

	return nil
}

// GetData retrieves data for a specific MIME type
func (dt *DataTransfer) GetData(mimeType string) (interface{}, bool) {
	data, exists := dt.dataMap[mimeType]
	return data, exists
}

// GetTypes returns all available MIME types
func (dt *DataTransfer) GetTypes() []string {
	result := make([]string, len(dt.types))
	copy(result, dt.types)
	return result
}

// ClearData removes data for a specific MIME type or all data if no type specified
func (dt *DataTransfer) ClearData(mimeType ...string) {
	if len(mimeType) == 0 {
		// Clear all data
		dt.items = make([]DataTransferItem, 0)
		dt.files = NewFileList()
		dt.types = make([]string, 0)
		dt.dataMap = make(map[string]interface{})
		return
	}

	targetType := mimeType[0]

	// Remove from items
	newItems := make([]DataTransferItem, 0)
	for _, item := range dt.items {
		if item.MimeType != targetType {
			newItems = append(newItems, item)
		}
	}
	dt.items = newItems

	// Remove from types
	newTypes := make([]string, 0)
	for _, t := range dt.types {
		if t != targetType {
			newTypes = append(newTypes, t)
		}
	}
	dt.types = newTypes

	// Remove from data map
	delete(dt.dataMap, targetType)
}

// GetFiles returns the file list
func (dt *DataTransfer) GetFiles() *FileList {
	return dt.files
}

// GetItems returns all data transfer items
func (dt *DataTransfer) GetItems() []DataTransferItem {
	result := make([]DataTransferItem, len(dt.items))
	copy(result, dt.items)
	return result
}

// AddFile adds a file by path
func (dt *DataTransfer) AddFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	// Detect MIME type
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return dt.SetData(mimeType, file)
}

// GetText retrieves text data
func (dt *DataTransfer) GetText() (string, bool) {
	if data, exists := dt.GetData("text/plain"); exists {
		if str, ok := data.(string); ok {
			return str, true
		}
	}
	return "", false
}

// WriteTo writes all text data to an io.Writer
func (dt *DataTransfer) WriteTo(w io.Writer) (int64, error) {
	var totalBytes int64

	for _, item := range dt.items {
		if item.Kind == DataTransferItemKindString {
			if str, err := item.GetAsString(); err == nil {
				n, err := w.Write([]byte(str))
				totalBytes += int64(n)
				if err != nil {
					return totalBytes, err
				}
			}
		}
	}

	return totalBytes, nil
}

// containsType checks if a MIME type is already in the types list
func (dt *DataTransfer) containsType(mimeType string) bool {
	for _, t := range dt.types {
		if strings.EqualFold(t, mimeType) {
			return true
		}
	}
	return false
}

// String returns a string representation of the DataTransfer
func (dt *DataTransfer) String() string {
	var sb strings.Builder
	sb.WriteString("DataTransfer:\n")
	sb.WriteString(fmt.Sprintf("DropEffect: %d\n", dt.dropEffect))
	sb.WriteString(fmt.Sprintf("EffectAllowed: %d\n", dt.effectAllowed))
	sb.WriteString("Items:\n")
	for _, item := range dt.items {
		sb.WriteString(fmt.Sprintf("- Kind: %d, MIME Type: %s\n", item.Kind, item.MimeType))
	}
	sb.WriteString("Files:\n")
	for i, file := range dt.files.files {
		if file != nil {
			sb.WriteString(fmt.Sprintf("- File %d: %s\n", i+1, file.Name()))
		} else {
			sb.WriteString(fmt.Sprintf("- File %d: <nil>\n", i+1))
		}
	}
	return sb.String()
}
