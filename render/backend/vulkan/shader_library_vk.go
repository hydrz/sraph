package vulkan

import (
	"sync"

	"github.com/opensraph/sraph/render"
	"github.com/vulkan-go/vulkan"
)

// ShaderLibraryVK manages Vulkan shader modules
type ShaderLibraryVK struct {
	device        vulkan.Device
	mutex         sync.RWMutex
	shaderModules map[string]*ShaderFunctionVK
	shaderCache   map[string][]byte
}

// NewShaderLibraryVK creates a new Vulkan shader library
func NewShaderLibraryVK(device vulkan.Device) *ShaderLibraryVK {
	return &ShaderLibraryVK{
		device:        device,
		shaderModules: make(map[string]*ShaderFunctionVK),
		shaderCache:   make(map[string][]byte),
	}
}

// GetFunction retrieves or loads a shader function
func (s *ShaderLibraryVK) GetFunction(name string, stage render.ShaderStage) render.ShaderFunction {
	s.mutex.RLock()
	if function, exists := s.shaderModules[name]; exists {
		s.mutex.RUnlock()
		return function
	}
	s.mutex.RUnlock()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Double-check after acquiring write lock
	if function, exists := s.shaderModules[name]; exists {
		return function
	}

	// Load shader function
	function, err := s.loadShaderFunction(name, stage)
	if err != nil {
		return nil
	}

	s.shaderModules[name] = function
	return function
}

// LoadShader loads a shader from bytecode
func (s *ShaderLibraryVK) LoadShader(name string, bytecode []byte, stage render.ShaderStage) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Cache the bytecode
	s.shaderCache[name] = bytecode

	// Create shader module
	function, err := s.createShaderFunction(name, bytecode, stage)
	if err != nil {
		return err
	}

	s.shaderModules[name] = function
	return nil
}

// CreateShaderModule creates a Vulkan shader module from SPIR-V bytecode
func (s *ShaderLibraryVK) CreateShaderModule(bytecode []byte) (vulkan.ShaderModule, error) {
	createInfo := vulkan.ShaderModuleCreateInfo{
		SType:    vulkan.StructureTypeShaderModuleCreateInfo,
		CodeSize: uint64(len(bytecode)),
		PCode:    bytecode,
	}

	var shaderModule vulkan.ShaderModule
	ret := vulkan.CreateShaderModule(s.device, &createInfo, nil, &shaderModule)
	if ret != vulkan.Success {
		return vulkan.NullHandle, render.NewError("Failed to create shader module")
	}

	return shaderModule, nil
}

// GetShaderStageFlags converts render shader stage to Vulkan stage flags
func (s *ShaderLibraryVK) GetShaderStageFlags(stage render.ShaderStage) vulkan.ShaderStageFlags {
	switch stage {
	case render.ShaderStageVertex:
		return vulkan.ShaderStageFlags(vulkan.ShaderStageVertexBit)
	case render.ShaderStageFragment:
		return vulkan.ShaderStageFlags(vulkan.ShaderStageFragmentBit)
	case render.ShaderStageCompute:
		return vulkan.ShaderStageFlags(vulkan.ShaderStageComputeBit)
	case render.ShaderStageGeometry:
		return vulkan.ShaderStageFlags(vulkan.ShaderStageGeometryBit)
	case render.ShaderStageTessellationControl:
		return vulkan.ShaderStageFlags(vulkan.ShaderStageTessellationControlBit)
	case render.ShaderStageTessellationEvaluation:
		return vulkan.ShaderStageFlags(vulkan.ShaderStageTessellationEvaluationBit)
	default:
		return 0
	}
}

// GetFunction returns a cached shader function
func (s *ShaderLibraryVK) GetCachedFunction(name string) *ShaderFunctionVK {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.shaderModules[name]
}

// RemoveFunction removes a shader function from cache
func (s *ShaderLibraryVK) RemoveFunction(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if function, exists := s.shaderModules[name]; exists {
		function.Cleanup()
		delete(s.shaderModules, name)
	}
	delete(s.shaderCache, name)
}

// GetFunctionCount returns the number of loaded shader functions
func (s *ShaderLibraryVK) GetFunctionCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.shaderModules)
}

// Cleanup cleans up all shader modules
func (s *ShaderLibraryVK) Cleanup() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, function := range s.shaderModules {
		function.Cleanup()
	}

	s.shaderModules = make(map[string]*ShaderFunctionVK)
	s.shaderCache = make(map[string][]byte)
}

// loadShaderFunction loads a shader function by name
func (s *ShaderLibraryVK) loadShaderFunction(name string, stage render.ShaderStage) (*ShaderFunctionVK, error) {
	// Check if bytecode is cached
	bytecode, exists := s.shaderCache[name]
	if !exists {
		return nil, render.NewError("Shader bytecode not found: " + name)
	}

	return s.createShaderFunction(name, bytecode, stage)
}

// createShaderFunction creates a shader function from bytecode
func (s *ShaderLibraryVK) createShaderFunction(name string, bytecode []byte, stage render.ShaderStage) (*ShaderFunctionVK, error) {
	shaderModule, err := s.CreateShaderModule(bytecode)
	if err != nil {
		return nil, err
	}

	return NewShaderFunctionVK(s.device, shaderModule, name, stage), nil
}
