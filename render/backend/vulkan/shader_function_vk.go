package vulkan

import (
	"github.com/opensraph/sraph/render"
	"github.com/vulkan-go/vulkan"
)

// ShaderFunctionVK represents a Vulkan shader function
type ShaderFunctionVK struct {
	device       vulkan.Device
	shaderModule vulkan.ShaderModule
	name         string
	entryPoint   string
	stage        render.ShaderStage
	stageInfo    vulkan.PipelineShaderStageCreateInfo
}

// NewShaderFunctionVK creates a new Vulkan shader function
func NewShaderFunctionVK(device vulkan.Device, shaderModule vulkan.ShaderModule, name string, stage render.ShaderStage) *ShaderFunctionVK {
	function := &ShaderFunctionVK{
		device:       device,
		shaderModule: shaderModule,
		name:         name,
		entryPoint:   "main", // Default SPIR-V entry point
		stage:        stage,
	}

	function.setupStageInfo()
	return function
}

// GetName returns the shader function name
func (s *ShaderFunctionVK) GetName() string {
	return s.name
}

// GetStage returns the shader stage
func (s *ShaderFunctionVK) GetStage() render.ShaderStage {
	return s.stage
}

// GetEntryPoint returns the shader entry point
func (s *ShaderFunctionVK) GetEntryPoint() string {
	return s.entryPoint
}

// SetEntryPoint sets the shader entry point
func (s *ShaderFunctionVK) SetEntryPoint(entryPoint string) {
	s.entryPoint = entryPoint
	s.setupStageInfo()
}

// GetShaderModule returns the Vulkan shader module
func (s *ShaderFunctionVK) GetShaderModule() vulkan.ShaderModule {
	return s.shaderModule
}

// GetStageInfo returns the pipeline shader stage create info
func (s *ShaderFunctionVK) GetStageInfo() vulkan.PipelineShaderStageCreateInfo {
	return s.stageInfo
}

// Cleanup cleans up the shader function resources
func (s *ShaderFunctionVK) Cleanup() {
	if s.shaderModule != vulkan.NullHandle {
		vulkan.DestroyShaderModule(s.device, s.shaderModule, nil)
		s.shaderModule = vulkan.NullHandle
	}
}

// setupStageInfo sets up the pipeline shader stage create info
func (s *ShaderFunctionVK) setupStageInfo() {
	s.stageInfo = vulkan.PipelineShaderStageCreateInfo{
		SType:  vulkan.StructureTypePipelineShaderStageCreateInfo,
		Stage:  s.convertStageToVulkan(),
		Module: s.shaderModule,
		PName:  s.entryPoint,
	}
}

// convertStageToVulkan converts render shader stage to Vulkan stage
func (s *ShaderFunctionVK) convertStageToVulkan() vulkan.ShaderStageFlagBits {
	switch s.stage {
	case render.ShaderStageVertex:
		return vulkan.ShaderStageVertexBit
	case render.ShaderStageFragment:
		return vulkan.ShaderStageFragmentBit
	case render.ShaderStageCompute:
		return vulkan.ShaderStageComputeBit
	case render.ShaderStageGeometry:
		return vulkan.ShaderStageGeometryBit
	case render.ShaderStageTessellationControl:
		return vulkan.ShaderStageTessellationControlBit
	case render.ShaderStageTessellationEvaluation:
		return vulkan.ShaderStageTessellationEvaluationBit
	default:
		return vulkan.ShaderStageVertexBit // Default fallback
	}
}

// GetSpecializationInfo returns specialization constants info
func (s *ShaderFunctionVK) GetSpecializationInfo() *vulkan.SpecializationInfo {
	// TODO: Implement specialization constants support
	return nil
}

// SetSpecializationConstants sets specialization constants
func (s *ShaderFunctionVK) SetSpecializationConstants(constants map[uint32]interface{}) {
	// TODO: Implement specialization constants
}

// Validate validates the shader function
func (s *ShaderFunctionVK) Validate() error {
	if s.shaderModule == vulkan.NullHandle {
		return render.NewError("Invalid shader module")
	}
	if s.entryPoint == "" {
		return render.NewError("Entry point not set")
	}
	return nil
}

// Clone creates a copy of the shader function
func (s *ShaderFunctionVK) Clone() *ShaderFunctionVK {
	return &ShaderFunctionVK{
		device:       s.device,
		shaderModule: s.shaderModule, // Shared reference
		name:         s.name,
		entryPoint:   s.entryPoint,
		stage:        s.stage,
		stageInfo:    s.stageInfo,
	}
}
