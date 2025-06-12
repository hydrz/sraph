package impl

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

// MemoryManager tracks WebGPU memory usage
type MemoryManager struct {
	mu               sync.RWMutex
	totalAllocated   uint64
	bufferAllocated  uint64
	textureAllocated uint64
	enabled          bool
}

var globalMemoryManager = &MemoryManager{
	enabled: false,
}

// EnableMemoryTracking enables memory usage tracking
func EnableMemoryTracking() {
	globalMemoryManager.mu.Lock()
	defer globalMemoryManager.mu.Unlock()
	globalMemoryManager.enabled = true
}

// DisableMemoryTracking disables memory usage tracking
func DisableMemoryTracking() {
	globalMemoryManager.mu.Lock()
	defer globalMemoryManager.mu.Unlock()
	globalMemoryManager.enabled = false
}

// TrackBufferAllocation tracks buffer memory allocation
func TrackBufferAllocation(size uint64) {
	if !globalMemoryManager.enabled {
		return
	}

	atomic.AddUint64(&globalMemoryManager.totalAllocated, size)
	atomic.AddUint64(&globalMemoryManager.bufferAllocated, size)
}

// TrackBufferDeallocation tracks buffer memory deallocation
func TrackBufferDeallocation(size uint64) {
	if !globalMemoryManager.enabled {
		return
	}

	atomic.AddUint64(&globalMemoryManager.totalAllocated, ^(size - 1)) // Subtract
	atomic.AddUint64(&globalMemoryManager.bufferAllocated, ^(size - 1))
}

// TrackTextureAllocation tracks texture memory allocation
func TrackTextureAllocation(size uint64) {
	if !globalMemoryManager.enabled {
		return
	}

	atomic.AddUint64(&globalMemoryManager.totalAllocated, size)
	atomic.AddUint64(&globalMemoryManager.textureAllocated, size)
}

// TrackTextureDeallocation tracks texture memory deallocation
func TrackTextureDeallocation(size uint64) {
	if !globalMemoryManager.enabled {
		return
	}

	atomic.AddUint64(&globalMemoryManager.totalAllocated, ^(size - 1))
	atomic.AddUint64(&globalMemoryManager.textureAllocated, ^(size - 1))
}

// GetMemoryUsage returns current memory usage statistics
func GetMemoryUsage() (total, buffer, texture uint64) {
	if !globalMemoryManager.enabled {
		return 0, 0, 0
	}

	return atomic.LoadUint64(&globalMemoryManager.totalAllocated),
		atomic.LoadUint64(&globalMemoryManager.bufferAllocated),
		atomic.LoadUint64(&globalMemoryManager.textureAllocated)
}

// CalculateTextureSize calculates the memory size of a texture
func CalculateTextureSize(descriptor TextureDescriptor) uint64 {
	blockSize := uint64(GetTextureFormatBlockSize(descriptor.Format))

	// Calculate total size considering all mip levels
	var totalSize uint64
	for mip := uint32(0); mip < descriptor.MipLevelCount; mip++ {
		mipWidth := max(1, descriptor.Size.Width>>mip)
		mipHeight := max(1, descriptor.Size.Height>>mip)
		mipDepth := descriptor.Size.DepthOrArrayLayers

		if descriptor.Dimension == TextureDimension3D {
			mipDepth = max(1, descriptor.Size.DepthOrArrayLayers>>mip)
		}

		mipSize := uint64(mipWidth) * uint64(mipHeight) * uint64(mipDepth) * blockSize
		totalSize += mipSize
	}

	return totalSize * uint64(descriptor.SampleCount)
}

// SafePointerAccess provides safe access to mapped buffer memory
type SafePointerAccess struct {
	ptr    unsafe.Pointer
	size   uintptr
	buffer Buffer
}

// NewSafePointerAccess creates a safe pointer accessor
func NewSafePointerAccess(buffer Buffer, ptr unsafe.Pointer, size uintptr) *SafePointerAccess {
	return &SafePointerAccess{
		ptr:    ptr,
		size:   size,
		buffer: buffer,
	}
}

// Read safely reads data from the buffer
func (spa *SafePointerAccess) Read(offset uintptr, dst []byte) error {
	if offset+uintptr(len(dst)) > spa.size {
		return fmt.Errorf("read exceeds buffer bounds")
	}

	// Check if buffer is still mapped
	state, err := spa.buffer.GetMapState()
	if err != nil {
		return err
	}
	if state != BufferMapStateMapped {
		return fmt.Errorf("buffer is not mapped")
	}

	src := (*[1 << 30]byte)(unsafe.Pointer(uintptr(spa.ptr) + offset))[:len(dst):len(dst)]
	copy(dst, src)
	return nil
}

// Write safely writes data to the buffer
func (spa *SafePointerAccess) Write(offset uintptr, src []byte) error {
	if offset+uintptr(len(src)) > spa.size {
		return fmt.Errorf("write exceeds buffer bounds")
	}

	// Check if buffer is still mapped
	state, err := spa.buffer.GetMapState()
	if err != nil {
		return err
	}
	if state != BufferMapStateMapped {
		return fmt.Errorf("buffer is not mapped")
	}

	dst := (*[1 << 30]byte)(unsafe.Pointer(uintptr(spa.ptr) + offset))[:len(src):len(src)]
	copy(dst, src)
	return nil
}
