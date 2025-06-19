package vulkan

import (
	"github.com/opensraph/sraph/render"
	"github.com/vulkan-go/vulkan"
)

// PipelineVK represents a Vulkan graphics pipeline
type PipelineVK struct {
	device           vulkan.Device
	pipeline         vulkan.Pipeline
	pipelineLayout   vulkan.PipelineLayout
	descriptor       render.GraphicsPipelineDescriptor
	vertexDescriptor *VertexDescriptorVK
}

// NewPipelineVK creates a new Vulkan graphics pipeline
func NewPipelineVK(device vulkan.Device, pipeline vulkan.Pipeline) *PipelineVK {
	return &PipelineVK{
		device:   device,
		pipeline: pipeline,
	}
}

// GetDescriptor returns the pipeline descriptor
func (p *PipelineVK) GetDescriptor() render.PipelineDescriptor {
	return p.descriptor
}

// BindToPipelineLayout binds the pipeline to a pipeline layout
func (p *PipelineVK) BindToPipelineLayout(layout vulkan.PipelineLayout) {
	p.pipelineLayout = layout
}

// GetPipeline returns the underlying Vulkan pipeline
func (p *PipelineVK) GetPipeline() vulkan.Pipeline {
	return p.pipeline
}

// GetPipelineLayout returns the pipeline layout
func (p *PipelineVK) GetPipelineLayout() vulkan.PipelineLayout {
	return p.pipelineLayout
}

// GetVertexDescriptor returns the vertex descriptor
func (p *PipelineVK) GetVertexDescriptor() *VertexDescriptorVK {
	return p.vertexDescriptor
}

// SetVertexDescriptor sets the vertex descriptor
func (p *PipelineVK) SetVertexDescriptor(descriptor *VertexDescriptorVK) {
	p.vertexDescriptor = descriptor
}

// Cleanup cleans up the pipeline resources
func (p *PipelineVK) Cleanup() {
	if p.pipeline != vulkan.NullHandle {
		vulkan.DestroyPipeline(p.device, p.pipeline, nil)
		p.pipeline = vulkan.NullHandle
	}
	if p.pipelineLayout != vulkan.NullHandle {
		vulkan.DestroyPipelineLayout(p.device, p.pipelineLayout, nil)
		p.pipelineLayout = vulkan.NullHandle
	}
}

// ComputePipelineVK represents a Vulkan compute pipeline
type ComputePipelineVK struct {
	device         vulkan.Device
	pipeline       vulkan.Pipeline
	pipelineLayout vulkan.PipelineLayout
	descriptor     render.ComputePipelineDescriptor
}

// NewComputePipelineVK creates a new Vulkan compute pipeline
func NewComputePipelineVK(device vulkan.Device, pipeline vulkan.Pipeline) *ComputePipelineVK {
	return &ComputePipelineVK{
		device:   device,
		pipeline: pipeline,
	}
}

// GetDescriptor returns the compute pipeline descriptor
func (p *ComputePipelineVK) GetDescriptor() render.ComputePipelineDescriptor {
	return p.descriptor
}

// BindToPipelineLayout binds the pipeline to a pipeline layout
func (p *ComputePipelineVK) BindToPipelineLayout(layout vulkan.PipelineLayout) {
	p.pipelineLayout = layout
}

// GetPipeline returns the underlying Vulkan pipeline
func (p *ComputePipelineVK) GetPipeline() vulkan.Pipeline {
	return p.pipeline
}

// GetPipelineLayout returns the pipeline layout
func (p *ComputePipelineVK) GetPipelineLayout() vulkan.PipelineLayout {
	return p.pipelineLayout
}

// Cleanup cleans up the compute pipeline resources
func (p *ComputePipelineVK) Cleanup() {
	if p.pipeline != vulkan.NullHandle {
		vulkan.DestroyPipeline(p.device, p.pipeline, nil)
		p.pipeline = vulkan.NullHandle
	}
	if p.pipelineLayout != vulkan.NullHandle {
		vulkan.DestroyPipelineLayout(p.device, p.pipelineLayout, nil)
		p.pipelineLayout = vulkan.NullHandle
	}
}

// PipelineLibraryVK manages a collection of Vulkan pipelines
type PipelineLibraryVK struct {
	device           vulkan.Device
	cache            *PipelineCacheVK
	pipelines        map[string]*PipelineVK
	computePipelines map[string]*ComputePipelineVK
}

// NewPipelineLibraryVK creates a new pipeline library
func NewPipelineLibraryVK(device vulkan.Device, cache *PipelineCacheVK) *PipelineLibraryVK {
	return &PipelineLibraryVK{
		device:           device,
		cache:            cache,
		pipelines:        make(map[string]*PipelineVK),
		computePipelines: make(map[string]*ComputePipelineVK),
	}
}

// GetPipeline gets or creates a graphics pipeline
func (p *PipelineLibraryVK) GetPipeline(descriptor render.GraphicsPipelineDescriptor) (*PipelineVK, error) {
	return p.cache.GetGraphicsPipeline(descriptor)
}

// GetComputePipeline gets or creates a compute pipeline
func (p *PipelineLibraryVK) GetComputePipeline(descriptor render.ComputePipelineDescriptor) (*ComputePipelineVK, error) {
	return p.cache.GetComputePipeline(descriptor)
}

// Cleanup cleans up the pipeline library
func (p *PipelineLibraryVK) Cleanup() {
	for _, pipeline := range p.pipelines {
		pipeline.Cleanup()
	}
	for _, pipeline := range p.computePipelines {
		pipeline.Cleanup()
	}
}
