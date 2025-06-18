package dl

import (
	"fmt"
	"unsafe"
)

// Storage manages a buffer for display list operations.
// It provides efficient allocation and deallocation of memory for display list data.
type Storage struct {
	data     []byte
	size     int
	capacity int
}

// NewStorage creates a new Storage with the specified initial capacity.
func NewStorage(initialCapacity int) *Storage {
	if initialCapacity <= 0 {
		initialCapacity = 1024 // Default capacity
	}
	return &Storage{
		data:     make([]byte, initialCapacity),
		size:     0,
		capacity: initialCapacity,
	}
}

// NewStorageEmpty creates a new empty Storage.
func NewStorageEmpty() *Storage {
	return &Storage{
		data:     nil,
		size:     0,
		capacity: 0,
	}
}

// Base returns a pointer to the base of the storage data.
func (s *Storage) Base() unsafe.Pointer {
	if len(s.data) == 0 {
		return nil
	}
	return unsafe.Pointer(&s.data[0])
}

// Size returns the current size of used data in bytes.
func (s *Storage) Size() int {
	return s.size
}

// Capacity returns the total capacity of the storage in bytes.
func (s *Storage) Capacity() int {
	return s.capacity
}

// Allocate allocates the specified number of bytes and returns a pointer to the allocated memory.
func (s *Storage) Allocate(bytes int) unsafe.Pointer {
	if bytes <= 0 {
		return nil
	}

	// Ensure we have enough capacity
	if s.size+bytes > s.capacity {
		s.realloc(s.size + bytes)
	}

	// Return pointer to the allocated memory
	ptr := unsafe.Pointer(&s.data[s.size])
	s.size += bytes

	return ptr
}

// Reset resets the storage to empty state without deallocating memory.
func (s *Storage) Reset() {
	s.size = 0
}

// Clear clears the storage and deallocates memory.
func (s *Storage) Clear() {
	s.data = nil
	s.size = 0
	s.capacity = 0
}

// realloc reallocates the storage to at least the specified capacity.
func (s *Storage) realloc(minCapacity int) {
	// Calculate new capacity (power of 2 growth)
	newCapacity := s.capacity
	if newCapacity == 0 {
		newCapacity = 1024
	}

	for newCapacity < minCapacity {
		newCapacity *= 2
	}

	// Allocate new buffer
	newData := make([]byte, newCapacity)

	// Copy existing data
	if s.size > 0 && s.data != nil {
		copy(newData, s.data[:s.size])
	}

	// Update storage
	s.data = newData
	s.capacity = newCapacity
}

// Clone creates a copy of the storage.
func (s *Storage) Clone() *Storage {
	newStorage := &Storage{
		size:     s.size,
		capacity: s.capacity,
	}

	if s.capacity > 0 {
		newStorage.data = make([]byte, s.capacity)
		if s.size > 0 {
			copy(newStorage.data, s.data[:s.size])
		}
	}

	return newStorage
}

// String returns a string representation of the storage.
func (s *Storage) String() string {
	return fmt.Sprintf("Storage{Size: %d, Capacity: %d}", s.size, s.capacity)
}

// WriteTo writes data to the storage at the current position.
func (s *Storage) WriteTo(data []byte) {
	if len(data) == 0 {
		return
	}

	ptr := s.Allocate(len(data))
	if ptr != nil {
		// Copy data to allocated memory
		copy((*[1 << 30]byte)(ptr)[:len(data)], data)
	}
}

// WriteValue writes a value of any type to the storage.
func (s *Storage) WriteValue(value interface{}) {
	switch v := value.(type) {
	case uint8:
		s.WriteTo([]byte{v})
	case uint16:
		data := (*[2]byte)(unsafe.Pointer(&v))[:]
		s.WriteTo(data)
	case uint32:
		data := (*[4]byte)(unsafe.Pointer(&v))[:]
		s.WriteTo(data)
	case uint64:
		data := (*[8]byte)(unsafe.Pointer(&v))[:]
		s.WriteTo(data)
	case float32:
		data := (*[4]byte)(unsafe.Pointer(&v))[:]
		s.WriteTo(data)
	case float64:
		data := (*[8]byte)(unsafe.Pointer(&v))[:]
		s.WriteTo(data)
	default:
		// For other types, we would need reflection or type-specific handling
		panic(fmt.Sprintf("unsupported type for WriteValue: %T", value))
	}
}
