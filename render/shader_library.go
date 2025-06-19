package render

import (
	"fmt"
	"sync"

	"github.com/opensraph/sraph/gpu"
)

// shaderLibrary implements ShaderLibrary interface.
type shaderLibrary struct {
	device        gpu.Device
	shaderCache   map[string]gpu.ShaderModule
	functionCache map[string]ShaderFunction
	cacheMutex    sync.RWMutex
}

// NewShaderLibrary creates a new shader library.
func NewShaderLibrary(device gpu.Device) ShaderLibrary {
	return &shaderLibrary{
		device:        device,
		shaderCache:   make(map[string]gpu.ShaderModule),
		functionCache: make(map[string]ShaderFunction),
	}
}

// GetShader implements ShaderLibrary.
func (sl *shaderLibrary) GetShader(name string) (gpu.ShaderModule, error) {
	sl.cacheMutex.RLock()
	if shader, exists := sl.shaderCache[name]; exists {
		sl.cacheMutex.RUnlock()
		return shader, nil
	}
	sl.cacheMutex.RUnlock()

	return nil, fmt.Errorf("shader '%s' not found", name)
}

// CompileShader implements ShaderLibrary.
func (sl *shaderLibrary) CompileShader(source string, stage ShaderStage) (gpu.ShaderModule, error) {
	// Create shader module descriptor
	desc := gpu.ShaderModuleDescriptor{
		Label: fmt.Sprintf("Shader_%s", getStageString(stage)),
	}

	// TODO: In a real implementation, we would compile the shader source
	// For now, we create a placeholder shader module
	shader := sl.device.CreateShaderModule(desc)
	if shader == nil {
		return nil, fmt.Errorf("failed to compile shader for stage %s", getStageString(stage))
	}

	return shader, nil
}

// GetShaderFunction implements ShaderLibrary.
func (sl *shaderLibrary) GetShaderFunction(name string, stage ShaderStage) (ShaderFunction, error) {
	key := fmt.Sprintf("%s_%s", name, getStageString(stage))

	sl.cacheMutex.RLock()
	if function, exists := sl.functionCache[key]; exists {
		sl.cacheMutex.RUnlock()
		return function, nil
	}
	sl.cacheMutex.RUnlock()

	return nil, fmt.Errorf("shader function '%s' for stage %s not found", name, getStageString(stage))
}

// AddShader adds a compiled shader to the library.
func (sl *shaderLibrary) AddShader(name string, shader gpu.ShaderModule) {
	sl.cacheMutex.Lock()
	defer sl.cacheMutex.Unlock()
	sl.shaderCache[name] = shader
}

// AddShaderFunction adds a shader function to the library.
func (sl *shaderLibrary) AddShaderFunction(name string, stage ShaderStage, module gpu.ShaderModule) {
	key := fmt.Sprintf("%s_%s", name, getStageString(stage))

	function := &shaderFunctionImpl{
		name:   name,
		stage:  stage,
		module: module,
	}

	sl.cacheMutex.Lock()
	defer sl.cacheMutex.Unlock()
	sl.functionCache[key] = function
}

// shaderFunctionImpl implements ShaderFunction interface.
type shaderFunctionImpl struct {
	name   string
	stage  ShaderStage
	module gpu.ShaderModule
}

// GetName implements ShaderFunction.
func (sf *shaderFunctionImpl) GetName() string {
	return sf.name
}

// GetStage implements ShaderFunction.
func (sf *shaderFunctionImpl) GetStage() ShaderStage {
	return sf.stage
}

// GetModule implements ShaderFunction.
func (sf *shaderFunctionImpl) GetModule() gpu.ShaderModule {
	return sf.module
}

// getStageString returns the string representation of a shader stage.
func getStageString(stage ShaderStage) string {
	switch stage {
	case ShaderStageVertex:
		return "vertex"
	case ShaderStageFragment:
		return "fragment"
	case ShaderStageCompute:
		return "compute"
	default:
		return "unknown"
	}
}
