package impl

import (
	"fmt"
	"sync"
	"unsafe"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ Buffer = (*buffer)(nil)

// buffer implements the Buffer interface
type buffer struct {
	mu               sync.RWMutex
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

// Destroy implements Buffer.Destroy.
// Destroys the buffer, releasing any resources associated with it.
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

// ConstMappedRange implements Buffer.ConstMappedRange.
// Gets a constant mapped range from the buffer.
func (b *buffer) ConstMappedRange(offset uintptr, size uintptr) (unsafe.Pointer, error) {
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

// MappedRange implements Buffer.MappedRange.
// Gets a mapped range from the buffer.
func (b *buffer) MappedRange(offset uintptr, size uintptr) (unsafe.Pointer, error) {
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

// MapState implements Buffer.MapState.
// Gets the buffer map state.
func (b *buffer) MapState() (BufferMapState, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.destroyed {
		return BufferMapStateUnmapped, fmt.Errorf("buffer has been destroyed")
	}

	return b.mapState, nil
}

// Size implements Buffer.Size.
// Gets the buffer size.
func (b *buffer) Size() (uint64, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.destroyed {
		return 0, fmt.Errorf("buffer has been destroyed")
	}

	return b.size, nil
}

// Usage implements Buffer.Usage.
// Gets the buffer usage.
func (b *buffer) Usage() (BufferUsage, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.destroyed {
		return BufferUsageNone, fmt.Errorf("buffer has been destroyed")
	}

	return b.usage, nil
}

// ReadMappedRange implements Buffer.ReadMappedRange.
// Reads from a mapped range in the buffer.
func (b *buffer) ReadMappedRange(offset uintptr, data *unsafe.Pointer, size uintptr) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.destroyed {
		return fmt.Errorf("buffer has been destroyed")
	}

	if b.mapState != BufferMapStateMapped {
		return fmt.Errorf("buffer is not mapped")
	}

	if offset+size > uintptr(len(b.data)) {
		return fmt.Errorf("offset and size exceed buffer bounds")
	}

	src := b.data[offset : offset+size]
	dst := (*[1 << 30]byte)(*data)[:size:size]
	copy(dst, src)

	return nil
}

// SetLabel implements Buffer.SetLabel.
// Sets the buffer label.
func (b *buffer) SetLabel(label string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.destroyed {
		return fmt.Errorf("buffer has been destroyed")
	}

	b.label = label
	return nil
}

// Unmap implements Buffer.Unmap.
// Unmaps the buffer, making it no longer accessible.
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

// WriteMappedRange implements Buffer.WriteMappedRange.
// Writes to a mapped range in the buffer.
func (b *buffer) WriteMappedRange(offset uintptr, data unsafe.Pointer, size uintptr) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.destroyed {
		return fmt.Errorf("buffer has been destroyed")
	}

	if b.mapState != BufferMapStateMapped {
		return fmt.Errorf("buffer is not mapped")
	}

	if offset+size > uintptr(len(b.data)) {
		return fmt.Errorf("offset and size exceed buffer bounds")
	}

	src := (*[1 << 30]byte)(data)[:size:size]
	dst := b.data[offset : offset+size]
	copy(dst, src)

	return nil
}
