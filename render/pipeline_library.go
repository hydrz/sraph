package render

import (
	"fmt"
	"sync"

	"github.com/opensraph/sraph/gpu"
)

// pipelineLibrary implements PipelineLibrary interface.
type pipelineLibrary struct {
	device           gpu.Device
	shaderLibrary    ShaderLibrary
	renderPipelines  map[string]gpu.RenderPipeline
	computePipelines map[string]gpu.ComputePipeline
	cacheMutex       sync.RWMutex
}

// NewPipelineLibrary creates a new pipeline library.
func NewPipelineLibrary(device gpu.Device, shaderLibrary ShaderLibrary) PipelineLibrary {
	return &pipelineLibrary{
		device:           device,
		shaderLibrary:    shaderLibrary,
		renderPipelines:  make(map[string]gpu.RenderPipeline),
		computePipelines: make(map[string]gpu.ComputePipeline),
	}
}

// GetRenderPipeline implements PipelineLibrary.
func (pl *pipelineLibrary) GetRenderPipeline(desc RenderPipelineDescriptor) (gpu.RenderPipeline, error) {
	key := generateRenderPipelineKey(desc)

	pl.cacheMutex.RLock()
	if pipeline, exists := pl.renderPipelines[key]; exists {
		pl.cacheMutex.RUnlock()
		return pipeline, nil
	}
	pl.cacheMutex.RUnlock()

	// Create new pipeline if not cached
	return pl.CreateRenderPipeline(desc)
}

// GetComputePipeline implements PipelineLibrary.
func (pl *pipelineLibrary) GetComputePipeline(desc ComputePipelineDescriptor) (gpu.ComputePipeline, error) {
	key := generateComputePipelineKey(desc)

	pl.cacheMutex.RLock()
	if pipeline, exists := pl.computePipelines[key]; exists {
		pl.cacheMutex.RUnlock()
		return pipeline, nil
	}
	pl.cacheMutex.RUnlock()

	// Create new pipeline if not cached
	return pl.CreateComputePipeline(desc)
}

// CreateRenderPipeline implements PipelineLibrary.
func (pl *pipelineLibrary) CreateRenderPipeline(desc RenderPipelineDescriptor) (gpu.RenderPipeline, error) {
	// Convert to GPU descriptor
	gpuDesc := gpu.RenderPipelineDescriptor{
		Label:  desc.Label,
		Layout: desc.Layout,
		Vertex: gpu.VertexState{
			Module:     desc.Vertex.Module,
			EntryPoint: desc.Vertex.EntryPoint,
			Buffers:    convertVertexBufferLayouts(desc.Vertex.Buffers),
		},
		Primitive: gpu.PrimitiveState{
			Topology:         desc.Primitive.Topology,
			StripIndexFormat: desc.Primitive.StripIndexFormat,
			FrontFace:        desc.Primitive.FrontFace,
			CullMode:         desc.Primitive.CullMode,
		},
		Multisample: gpu.MultisampleState{
			Count:                  desc.Multisample.Count,
			Mask:                   desc.Multisample.Mask,
			AlphaToCoverageEnabled: desc.Multisample.AlphaToCoverageEnabled,
		},
	}

	// Convert depth stencil state if present
	if desc.DepthStencil != nil {
		depthWriteEnabled := gpu.OptionalBoolUndefined
		if desc.DepthStencil.DepthWriteEnabled {
			depthWriteEnabled = gpu.OptionalBoolTrue
		} else {
			depthWriteEnabled = gpu.OptionalBoolFalse
		}

		gpuDesc.DepthStencil = gpu.DepthStencilState{
			Format:              desc.DepthStencil.Format,
			DepthWriteEnabled:   depthWriteEnabled,
			DepthCompare:        desc.DepthStencil.DepthCompare,
			StencilFront:        convertStencilFaceState(desc.DepthStencil.StencilFront),
			StencilBack:         convertStencilFaceState(desc.DepthStencil.StencilBack),
			StencilReadMask:     desc.DepthStencil.StencilReadMask,
			StencilWriteMask:    desc.DepthStencil.StencilWriteMask,
			DepthBias:           desc.DepthStencil.DepthBias,
			DepthBiasSlopeScale: desc.DepthStencil.DepthBiasSlopeScale,
			DepthBiasClamp:      desc.DepthStencil.DepthBiasClamp,
		}
	}

	// Convert fragment state if present
	if desc.Fragment != nil {
		gpuDesc.Fragment = gpu.FragmentState{
			Module:     desc.Fragment.Module,
			EntryPoint: desc.Fragment.EntryPoint,
			Targets:    convertColorTargetStates(desc.Fragment.Targets),
		}
	}

	pipeline := pl.device.CreateRenderPipeline(gpuDesc)
	if pipeline == nil {
		return nil, fmt.Errorf("failed to create render pipeline")
	}

	// Cache the pipeline
	key := generateRenderPipelineKey(desc)
	pl.cacheMutex.Lock()
	pl.renderPipelines[key] = pipeline
	pl.cacheMutex.Unlock()

	return pipeline, nil
}

// CreateComputePipeline implements PipelineLibrary.
func (pl *pipelineLibrary) CreateComputePipeline(desc ComputePipelineDescriptor) (gpu.ComputePipeline, error) {
	gpuDesc := gpu.ComputePipelineDescriptor{
		Label:  desc.Label,
		Layout: desc.Layout,
		Compute: gpu.ComputeState{
			Module:     desc.Compute.Module,
			EntryPoint: desc.Compute.EntryPoint,
		},
	}

	pipeline := pl.device.CreateComputePipeline(gpuDesc)
	if pipeline == nil {
		return nil, fmt.Errorf("failed to create compute pipeline")
	}

	// Cache the pipeline
	key := generateComputePipelineKey(desc)
	pl.cacheMutex.Lock()
	pl.computePipelines[key] = pipeline
	pl.cacheMutex.Unlock()

	return pipeline, nil
}

// Helper functions for conversion

func convertVertexBufferLayouts(layouts []VertexBufferLayout) []gpu.VertexBufferLayout {
	result := make([]gpu.VertexBufferLayout, len(layouts))
	for i, layout := range layouts {
		result[i] = gpu.VertexBufferLayout{
			ArrayStride: layout.ArrayStride,
			StepMode:    layout.StepMode,
			Attributes:  convertVertexAttributes(layout.Attributes),
		}
	}
	return result
}

func convertVertexAttributes(attributes []VertexAttribute) []gpu.VertexAttribute {
	result := make([]gpu.VertexAttribute, len(attributes))
	for i, attr := range attributes {
		result[i] = gpu.VertexAttribute{
			Format:         attr.Format,
			Offset:         attr.Offset,
			ShaderLocation: attr.ShaderLocation,
		}
	}
	return result
}

func convertStencilFaceState(state StencilFaceState) gpu.StencilFaceState {
	return gpu.StencilFaceState{
		Compare:     state.Compare,
		FailOp:      state.FailOp,
		DepthFailOp: state.DepthFailOp,
		PassOp:      state.PassOp,
	}
}

func convertColorTargetStates(targets []ColorTargetState) []gpu.ColorTargetState {
	result := make([]gpu.ColorTargetState, len(targets))
	for i, target := range targets {
		result[i] = gpu.ColorTargetState{
			Format:    target.Format,
			WriteMask: target.WriteMask,
		}
		if target.Blend != nil {
			result[i].Blend = gpu.BlendState{
				Color: gpu.BlendComponent{
					Operation: target.Blend.Color.Operation,
					SrcFactor: target.Blend.Color.SrcFactor,
					DstFactor: target.Blend.Color.DstFactor,
				},
				Alpha: gpu.BlendComponent{
					Operation: target.Blend.Alpha.Operation,
					SrcFactor: target.Blend.Alpha.SrcFactor,
					DstFactor: target.Blend.Alpha.DstFactor,
				},
			}
		}
	}
	return result
}

// generateRenderPipelineKey creates a unique key for render pipeline caching.
func generateRenderPipelineKey(desc RenderPipelineDescriptor) string {
	// TODO: Implement proper hash key generation based on descriptor content
	return fmt.Sprintf("render_%s", desc.Label)
}

// generateComputePipelineKey creates a unique key for compute pipeline caching.
func generateComputePipelineKey(desc ComputePipelineDescriptor) string {
	// TODO: Implement proper hash key generation based on descriptor content
	return fmt.Sprintf("compute_%s", desc.Label)
}
