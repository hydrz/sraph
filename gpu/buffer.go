package gpu

import (
	"fmt"
	"sync"
	"unsafe"
)

var _ Buffer = (*buffer)(nil)

// buffer implements the Buffer interface
type buffer struct {
	mu               sync.RWMutex
	refCount         int32
	label            string
	usage            BufferUsage
	size             uint64
	mappedAtCreation bool
	mapState         BufferMapState
	data             []byte
	destroyed        bool
}

// NewBuffer creates a new WebGPU buffer
func NewBuffer(descriptor BufferDescriptor) Buffer {
	buffer := &buffer{
		refCount:         1,
		label:            descriptor.Label,
		usage:            descriptor.Usage,
		size:             descriptor.Size,
		mappedAtCreation: descriptor.MappedAtCreation,
		mapState:         BufferMapStateUnmapped,
	}

	if descriptor.MappedAtCreation {
		buffer.mapState = BufferMapStateMapped
		buffer.data = make([]byte, descriptor.Size)
	}

	return buffer
}

// Destroy destroys the buffer
func (b *buffer) Destroy() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.destroyed {
		return fmt.Errorf("buffer has been destroyed")
	}

	b.destroyed = true
	b.data = nil
	return nil
}

// GetConstMappedRange gets a constant mapped range
func (b *buffer) GetConstMappedRange(offset uintptr, size uintptr) (unsafe.Pointer, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.destroyed {
		return nil, fmt.Errorf("buffer has been destroyed")
	}

	if b.mapState != BufferMapStateMapped {
		return nil, fmt.Errorf("buffer is not mapped")
	}

	if offset+size > uintptr(len(b.data)) {
		return nil, fmt.Errorf("offset and size exceed buffer bounds")
	}

	return unsafe.Pointer(&b.data[offset]), nil
}

// GetMappedRange gets a mapped range
func (b *buffer) GetMappedRange(offset uintptr, size uintptr) (unsafe.Pointer, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.destroyed {
		return nil, fmt.Errorf("buffer has been destroyed")
	}

	if b.mapState != BufferMapStateMapped {
		return nil, fmt.Errorf("buffer is not mapped")
	}

	if offset+size > uintptr(len(b.data)) {
		return nil, fmt.Errorf("offset and size exceed buffer bounds")
	}

	return unsafe.Pointer(&b.data[offset]), nil
}

// GetMapState gets the buffer map state
func (b *buffer) GetMapState() (BufferMapState, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.destroyed {
		return BufferMapStateUnmapped, fmt.Errorf("buffer has been destroyed")
	}

	return b.mapState, nil
}

// GetSize gets the buffer size
func (b *buffer) GetSize() (uint64, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.destroyed {
		return 0, fmt.Errorf("buffer has been destroyed")
	}

	return b.size, nil
}

// GetUsage gets the buffer usage
func (b *buffer) GetUsage() (BufferUsage, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.destroyed {
		return BufferUsageNone, fmt.Errorf("buffer has been destroyed")
	}

	return b.usage, nil
}

// MapAsync maps the buffer asynchronously
func (b *buffer) MapAsync(mode MapMode, offset uintptr, size uintptr, callback BufferMapCallbackInfo) Future {
	b.mu.Lock()
	defer b.mu.Unlock()

	future := Future{
		Id: generateFutureId(),
	}

	if b.destroyed {
		return future
	}

	// Simulate async mapping
	go func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		if !b.destroyed && b.mapState == BufferMapStateUnmapped {
			b.mapState = BufferMapStateMapped
			if b.data == nil {
				b.data = make([]byte, b.size)
			}
		}
	}()

	return future
}

// ReadMappedRange reads from a mapped range
func (b *buffer) ReadMappedRange(offset uintptr, data unsafe.Pointer, size uintptr) (Status, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.destroyed {
		return StatusError, fmt.Errorf("buffer has been destroyed")
	}

	if b.mapState != BufferMapStateMapped {
		return StatusError, fmt.Errorf("buffer is not mapped")
	}

	if offset+size > uintptr(len(b.data)) {
		return StatusError, fmt.Errorf("offset and size exceed buffer bounds")
	}

	// Copy data from buffer to provided pointer
	src := b.data[offset : offset+size]
	dst := (*[1 << 30]byte)(data)[:size:size]
	copy(dst, src)

	return StatusSuccess, nil
}

// SetLabel sets the buffer label
func (b *buffer) SetLabel(label string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.destroyed {
		return fmt.Errorf("buffer has been destroyed")
	}

	b.label = label
	return nil
}

// Unmap unmaps the buffer
func (b *buffer) Unmap() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.destroyed {
		return fmt.Errorf("buffer has been destroyed")
	}

	if b.mapState == BufferMapStateMapped {
		b.mapState = BufferMapStateUnmapped
	}

	return nil
}

// WriteMappedRange writes to a mapped range
func (b *buffer) WriteMappedRange(offset uintptr, data unsafe.Pointer, size uintptr) (Status, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.destroyed {
		return StatusError, fmt.Errorf("buffer has been destroyed")
	}

	if b.mapState != BufferMapStateMapped {
		return StatusError, fmt.Errorf("buffer is not mapped")
	}

	if offset+size > uintptr(len(b.data)) {
		return StatusError, fmt.Errorf("offset and size exceed buffer bounds")
	}

	// Copy data from provided pointer to buffer
	src := (*[1 << 30]byte)(data)[:size:size]
	dst := b.data[offset : offset+size]
	copy(dst, src)

	return StatusSuccess, nil
}

// AddRef increments the reference count
func (b *buffer) AddRef() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.destroyed {
		return fmt.Errorf("buffer has been destroyed")
	}

	b.refCount++
	return nil
}

// Release decrements the reference count and destroys if zero
func (b *buffer) Release() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.destroyed {
		return fmt.Errorf("buffer has been destroyed")
	}

	b.refCount--
	if b.refCount <= 0 {
		b.destroyed = true
		b.data = nil
	}

	return nil
}
