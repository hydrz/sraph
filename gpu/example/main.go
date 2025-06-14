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
	buffer := device.CreateBuffer(gpu.BufferDescriptor{
		Size:  4 * 4, // 4 floats
		Usage: gpu.BufferUsageVertex | gpu.BufferUsageCopyDst,
	})

	// 5. Create a shader module (WGSL)
	shaderModule := device.CreateShaderModule(gpu.ShaderModuleDescriptor{
		Label: "simple-shader",
	})

	// 6. Create pipeline layout
	pipelineLayout := device.CreatePipelineLayout(gpu.PipelineLayoutDescriptor{})

	// 7. Create render pipeline
	renderPipeline := device.CreateRenderPipeline(gpu.RenderPipelineDescriptor{
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

	// 8. Create command encoder
	commandEncoder := device.CreateCommandEncoder(gpu.CommandEncoderDescriptor{})

	// 9. Begin render pass (assuming a valid texture view)
	// In a real application, acquire a swapchain texture view here.
	var colorAttachment gpu.RenderPassColorAttachment
	renderPass := commandEncoder.BeginRenderPass(gpu.RenderPassDescriptor{
		ColorAttachments: []gpu.RenderPassColorAttachment{
			colorAttachment,
		},
	})

	// 10. Set pipeline and vertex buffer, then draw
	renderPass.SetPipeline(renderPipeline)
	renderPass.SetVertexBuffer(0, buffer, 0, 4*4)
	renderPass.Draw(3, 1, 0, 0)
	renderPass.End()

	// 11. Finish command buffer
	commandBuffer := commandEncoder.Finish(gpu.CommandBufferDescriptor{})

	// 12. Submit commands to queue
	queue := device.Queue()
	queue.Submit([]gpu.CommandBuffer{commandBuffer})

	// 13. Present (if using a swapchain)
	// surface.Present() // Uncomment and implement as needed

	log.Println("Render submission complete.")
}
