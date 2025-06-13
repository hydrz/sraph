package main

import (
	"log"

	gpu "github.com/opensraph/sraph/gpu"
)

func main() {
	// 2. Request GPU adapter
	adapter, err := gpu.RequestAdapter(gpu.RequestAdapterOptions{})
	if err != nil {
		log.Fatalf("Failed to request GPU adapter: %v", err)
	}

	// 3. Request logical device
	device, err := adapter.RequestDevice(gpu.DeviceDescriptor{})
	if err != nil {
		log.Fatalf("Failed to request GPU device: %v", err)
	}

	// 4. Create a buffer
	buffer, err := device.CreateBuffer(gpu.BufferDescriptor{
		Size:  4 * 4, // 4 floats
		Usage: gpu.BufferUsageVertex | gpu.BufferUsageCopyDst,
	})
	if err != nil {
		log.Fatalf("Failed to create buffer: %v", err)
	}

	// 5. Create a shader module (WGSL)
	shaderModule, err := device.CreateShaderModule(gpu.ShaderModuleDescriptor{
		Label: "simple-shader",
	})
	if err != nil {
		log.Fatalf("Failed to create shader module: %v", err)
	}

	// 6. Create pipeline layout
	pipelineLayout, err := device.CreatePipelineLayout(gpu.PipelineLayoutDescriptor{})
	if err != nil {
		log.Fatalf("Failed to create pipeline layout: %v", err)
	}

	// 7. Create render pipeline
	renderPipeline, err := device.CreateRenderPipeline(gpu.RenderPipelineDescriptor{
		Layout: pipelineLayout,
		Vertex: gpu.VertexState{
			Module:     shaderModule,
			EntryPoint: "vs_main",
		},
		Fragment: gpu.FragmentState{
			Module:     shaderModule,
			EntryPoint: "fs_main",
			Targets: []gpu.ColorTargetState{
				{Format: gpu.TextureFormatBGRA8Unorm},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create render pipeline: %v", err)
	}

	// 8. Create command encoder
	commandEncoder, err := device.CreateCommandEncoder(gpu.CommandEncoderDescriptor{})
	if err != nil {
		log.Fatalf("Failed to create command encoder: %v", err)
	}

	// 9. Begin render pass (assuming a valid texture view)
	// In a real application, acquire a swapchain texture view here.
	var colorAttachment gpu.RenderPassColorAttachment
	renderPass, err := commandEncoder.BeginRenderPass(gpu.RenderPassDescriptor{
		ColorAttachments: []gpu.RenderPassColorAttachment{
			colorAttachment,
		},
	})
	if err != nil {
		log.Fatalf("Failed to begin render pass: %v", err)
	}

	// 10. Set pipeline and vertex buffer, then draw
	if err := renderPass.SetPipeline(renderPipeline); err != nil {
		log.Fatalf("Failed to set pipeline: %v", err)
	}
	if err := renderPass.SetVertexBuffer(0, buffer, 0, 4*4); err != nil {
		log.Fatalf("Failed to set vertex buffer: %v", err)
	}
	if err := renderPass.Draw(3, 1, 0, 0); err != nil {
		log.Fatalf("Failed to issue draw: %v", err)
	}
	if err := renderPass.End(); err != nil {
		log.Fatalf("Failed to end render pass: %v", err)
	}

	// 11. Finish command buffer
	commandBuffer, err := commandEncoder.Finish(gpu.CommandBufferDescriptor{})
	if err != nil {
		log.Fatalf("Failed to finish command buffer: %v", err)
	}

	// 12. Submit commands to queue
	queue, err := device.Queue()
	if err != nil {
		log.Fatalf("Failed to get queue: %v", err)
	}
	if err := queue.Submit([]gpu.CommandBuffer{commandBuffer}); err != nil {
		log.Fatalf("Failed to submit command buffer: %v", err)
	}

	// 13. Present (if using a swapchain)
	// surface.Present() // Uncomment and implement as needed

	log.Println("Render submission complete.")
}
