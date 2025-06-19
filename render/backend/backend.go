package backend

// Backend defines the interface for graphics backends
type Backend interface {
	// GetName returns the name of the backend
	GetName() string

	// Initialize initializes the backend
	Initialize() error

	// Shutdown shuts down the backend
	Shutdown() error

	// IsSupported returns true if the backend is supported on the current platform
	IsSupported() bool

	// GetCapabilities returns the capabilities of the backend
	GetCapabilities() Capabilities

	// CreateContext creates a new render context
	CreateContext() (Context, error)
}

// BackendType defines the type of graphics backend
type BackendType int

const (
	BackendTypeUnknown BackendType = iota
	BackendTypeOpenGL
	BackendTypeVulkan
	BackendTypeMetal
	BackendTypeD3D11
	BackendTypeD3D12
	BackendTypeWebGPU
)

// Capabilities represents the capabilities of a graphics backend
type Capabilities struct {
	MaxTextureSize            int
	MaxCubeMapTextureSize     int
	Max3DTextureSize          int
	MaxArrayTextureLayers     int
	MaxColorAttachments       int
	MaxSamples                int
	MaxVertexAttributes       int
	MaxVertexBufferBindings   int
	MaxUniformBufferBindings  int
	MaxStorageBufferBindings  int
	MaxTextureBindings        int
	MaxSamplerBindings        int
	SupportsCompute           bool
	SupportsGeometryShader    bool
	SupportsTessellation      bool
	SupportsMultiDrawIndirect bool
	SupportsTimestampQueries  bool
}

// Context represents a render context for a specific backend
type Context interface {
	// GetBackendType returns the backend type
	GetBackendType() BackendType

	// GetCapabilities returns the backend capabilities
	GetCapabilities() Capabilities

	// CreateBuffer creates a new buffer
	CreateBuffer(size int, usage BufferUsage) (Buffer, error)

	// CreateTexture creates a new texture
	CreateTexture(descriptor TextureDescriptor) (Texture, error)

	// CreateSampler creates a new sampler
	CreateSampler(descriptor SamplerDescriptor) (Sampler, error)

	// CreateRenderPipeline creates a new render pipeline
	CreateRenderPipeline(descriptor RenderPipelineDescriptor) (RenderPipeline, error)

	// CreateComputePipeline creates a new compute pipeline
	CreateComputePipeline(descriptor ComputePipelineDescriptor) (ComputePipeline, error)

	// CreateCommandBuffer creates a new command buffer
	CreateCommandBuffer() (CommandBuffer, error)

	// Present presents the rendered content
	Present() error
}

// BackendRegistry manages available graphics backends
type BackendRegistry interface {
	// RegisterBackend registers a new backend
	RegisterBackend(backendType BackendType, backend Backend) error

	// GetBackend returns a backend by type
	GetBackend(backendType BackendType) (Backend, error)

	// GetAvailableBackends returns all available backends
	GetAvailableBackends() []BackendType

	// GetDefaultBackend returns the default backend for the current platform
	GetDefaultBackend() (Backend, error)
}

// BackendRegistryImpl is the default implementation of BackendRegistry
type BackendRegistryImpl struct {
	backends map[BackendType]Backend
}

// NewBackendRegistry creates a new backend registry
func NewBackendRegistry() BackendRegistry {
	return &BackendRegistryImpl{
		backends: make(map[BackendType]Backend),
	}
}

// RegisterBackend registers a new backend
func (r *BackendRegistryImpl) RegisterBackend(backendType BackendType, backend Backend) error {
	if backend == nil {
		return ErrInvalidArgument
	}

	r.backends[backendType] = backend
	return nil
}

// GetBackend returns a backend by type
func (r *BackendRegistryImpl) GetBackend(backendType BackendType) (Backend, error) {
	backend, exists := r.backends[backendType]
	if !exists {
		return nil, ErrResourceNotFound
	}

	return backend, nil
}

// GetAvailableBackends returns all available backends
func (r *BackendRegistryImpl) GetAvailableBackends() []BackendType {
	backends := make([]BackendType, 0, len(r.backends))
	for backendType := range r.backends {
		backends = append(backends, backendType)
	}
	return backends
}

// GetDefaultBackend returns the default backend for the current platform
func (r *BackendRegistryImpl) GetDefaultBackend() (Backend, error) {
	// TODO: Implement platform-specific default backend selection
	for _, backend := range r.backends {
		if backend.IsSupported() {
			return backend, nil
		}
	}

	return nil, ErrResourceNotFound
}

// Global backend registry instance
var globalRegistry BackendRegistry = NewBackendRegistry()

// RegisterBackend registers a backend globally
func RegisterBackend(backendType BackendType, backend Backend) error {
	return globalRegistry.RegisterBackend(backendType, backend)
}

// GetBackend returns a backend by type from the global registry
func GetBackend(backendType BackendType) (Backend, error) {
	return globalRegistry.GetBackend(backendType)
}

// GetAvailableBackends returns all available backends from the global registry
func GetAvailableBackends() []BackendType {
	return globalRegistry.GetAvailableBackends()
}

// GetDefaultBackend returns the default backend from the global registry
func GetDefaultBackend() (Backend, error) {
	return globalRegistry.GetDefaultBackend()
}
