package render

// ComputePipelineDescriptor describes the configuration for a compute pipeline
type ComputePipelineDescriptor struct {
	Label                         string
	ComputeFunction               ShaderFunction
	MaxTotalThreadsPerThreadgroup int
	ThreadgroupSizeX              int
	ThreadgroupSizeY              int
	ThreadgroupSizeZ              int
}

// IsValid checks if the compute pipeline descriptor is valid
func (d *ComputePipelineDescriptor) IsValid() bool {
	return d.ComputeFunction != nil && d.MaxTotalThreadsPerThreadgroup > 0
}

// ComputePipeline represents a compute pipeline for compute shaders
type ComputePipeline interface {
	Pipeline

	// GetDescriptor returns the compute pipeline descriptor
	GetDescriptor() ComputePipelineDescriptor

	// GetMaxTotalThreadsPerThreadgroup returns the maximum total threads per threadgroup
	GetMaxTotalThreadsPerThreadgroup() int

	// GetThreadgroupSize returns the threadgroup size in each dimension
	GetThreadgroupSize() (int, int, int)
}

// ComputePipelineImpl is the default implementation of ComputePipeline
type ComputePipelineImpl struct {
	*PipelineImpl
	descriptor ComputePipelineDescriptor
}

// NewComputePipeline creates a new compute pipeline with the specified descriptor
func NewComputePipeline(descriptor ComputePipelineDescriptor) ComputePipeline {
	return &ComputePipelineImpl{
		PipelineImpl: &PipelineImpl{
			label:   descriptor.Label,
			isValid: descriptor.IsValid(),
		},
		descriptor: descriptor,
	}
}

// GetDescriptor returns the compute pipeline descriptor
func (p *ComputePipelineImpl) GetDescriptor() ComputePipelineDescriptor {
	return p.descriptor
}

// GetMaxTotalThreadsPerThreadgroup returns the maximum total threads per threadgroup
func (p *ComputePipelineImpl) GetMaxTotalThreadsPerThreadgroup() int {
	return p.descriptor.MaxTotalThreadsPerThreadgroup
}

// GetThreadgroupSize returns the threadgroup size in each dimension
func (p *ComputePipelineImpl) GetThreadgroupSize() (int, int, int) {
	return p.descriptor.ThreadgroupSizeX,
		p.descriptor.ThreadgroupSizeY,
		p.descriptor.ThreadgroupSizeZ
}
