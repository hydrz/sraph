package surface

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/opensraph/sraph/gpu"
	"github.com/rajveermalviya/go-wayland/wayland/client"
)

// WaylandSurfaceManager manages WebGPU surface creation for Wayland windows
type WaylandSurfaceManager struct {
	instance gpu.Instance
	adapter  gpu.Adapter
	device   gpu.Device
}

// WaylandSurfaceHandle contains the native Wayland surface handles
type WaylandSurfaceHandle struct {
	Display *client.Display
	Surface *client.Surface
}

// NewWaylandSurfaceManager creates a new Wayland surface manager
func NewWaylandSurfaceManager() (*WaylandSurfaceManager, error) {
	// Get available backends
	backends := gpu.GetAllBackends()
	if len(backends) == 0 {
		return nil, fmt.Errorf("no GPU backends available")
	}

	// Use the first available backend
	backend := backends[0]

	// Create GPU instance
	instance := backend.CreateInstance(gpu.InstanceDescriptor{})
	if instance == nil {
		return nil, fmt.Errorf("failed to create GPU instance")
	}

	return &WaylandSurfaceManager{
		instance: instance,
	}, nil
}

// CreateSurface creates a WebGPU surface from Wayland surface handles
func (m *WaylandSurfaceManager) CreateSurface(handle *WaylandSurfaceHandle, label string) (gpu.Surface, error) {
	if handle == nil || handle.Display == nil || handle.Surface == nil {
		return nil, fmt.Errorf("invalid Wayland surface handle")
	}

	// Create surface descriptor with Wayland-specific chain
	surfaceDesc := gpu.SurfaceDescriptor{
		Label: label,
		// Note: The actual chaining with SurfaceSourceWaylandSurface needs to be implemented
		// in the GPU backend. For now, we create a basic surface.
	}

	// TODO: Add proper Wayland surface chaining
	// This would require extending the SurfaceDescriptor to support chained structures
	// like SurfaceSourceWaylandSurface with Display and Surface pointers

	gpuSurface, err := m.instance.CreateSurface(surfaceDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to create GPU surface: %w", err)
	}

	log.Printf("Created WebGPU surface with label: %s", label)
	return gpuSurface, nil
}

// CreateAdapter requests a compatible GPU adapter for the given surface
func (m *WaylandSurfaceManager) CreateAdapter(surface gpu.Surface, powerPreference gpu.PowerPreference) error {
	adapter, err := m.instance.RequestAdapter(gpu.RequestAdapterOptions{
		PowerPreference:   powerPreference,
		CompatibleSurface: surface,
	})
	if err != nil {
		return fmt.Errorf("failed to request GPU adapter: %w", err)
	}

	m.adapter = adapter
	log.Printf("Created GPU adapter")
	return nil
}

// CreateDevice requests a GPU device from the adapter
func (m *WaylandSurfaceManager) CreateDevice(label string, requiredFeatures []gpu.FeatureName) error {
	if m.adapter == nil {
		return fmt.Errorf("adapter not created")
	}

	device, err := m.adapter.RequestDevice(gpu.DeviceDescriptor{
		Label:            label,
		RequiredFeatures: requiredFeatures,
	})
	if err != nil {
		return fmt.Errorf("failed to request GPU device: %w", err)
	}

	m.device = device
	log.Printf("Created GPU device with label: %s", label)
	return nil
}

// GetDevice returns the created GPU device
func (m *WaylandSurfaceManager) GetDevice() gpu.Device {
	return m.device
}

// GetAdapter returns the created GPU adapter
func (m *WaylandSurfaceManager) GetAdapter() gpu.Adapter {
	return m.adapter
}

// GetInstance returns the GPU instance
func (m *WaylandSurfaceManager) GetInstance() gpu.Instance {
	return m.instance
}

// CreateWaylandSurfaceWithPointers creates a WebGPU surface with raw pointers
// This is a helper function for when we need to pass raw pointers to the WebGPU backend
func CreateWaylandSurfaceWithPointers(instance gpu.Instance, displayPtr, surfacePtr unsafe.Pointer, label string) (gpu.Surface, error) {
	// Create surface descriptor
	surfaceDesc := gpu.SurfaceDescriptor{
		Label: label,
		// Note: Need to implement proper chaining mechanism in the GPU package
		// to attach Wayland-specific surface source data (display and surface pointers)
	}

	// For now, create a basic surface
	// TODO: Implement proper chaining to pass Wayland-specific data
	// The waylandSource should be chained to surfaceDesc when the GPU backend supports it
	surface, err := instance.CreateSurface(surfaceDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to create Wayland surface: %w", err)
	}

	log.Printf("Created Wayland WebGPU surface with pointers: display=%p, surface=%p", displayPtr, surfacePtr)
	return surface, nil
}
