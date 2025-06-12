package gpu

import (
	"context"
	"testing"
)

// TestInterfaceImplementations verifies that all concrete types implement their respective interfaces
func TestInterfaceImplementations(t *testing.T) {
	ctx := context.Background()

	// Test Instance interface
	instanceDesc := InstanceDescriptor{}
	instance, err := gpu.CreateInstance(ctx, instanceDesc)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	if instance == nil {
		t.Fatal("CreateInstance returned nil")
	}

	// Test Adapter interface through Instance
	adapterOptions := RequestAdapterOptions{}
	adapterFuture := instance.RequestAdapter(ctx, adapterOptions)
	if adapterFuture.Id == 0 {
		t.Fatal("RequestAdapter returned invalid future")
	}

	// Test Surface interface
	surfaceDesc := SurfaceDescriptor{Label: "test-surface"}
	surface, err := instance.CreateSurface(ctx, surfaceDesc)
	if err != nil {
		t.Fatalf("CreateSurface failed: %v", err)
	}
	if surface == nil {
		t.Fatal("CreateSurface returned nil")
	}

	// For testing purposes, create mock implementations
	_ = &adapter{refCount: 1}
	device := &device{refCount: 1}

	// Test Device interface methods
	bufferDesc := BufferDescriptor{Size: 1024, Usage: BufferUsageVertex}
	buffer, err := device.CreateBuffer(ctx, bufferDesc)
	if err != nil {
		t.Fatalf("CreateBuffer failed: %v", err)
	}
	if buffer == nil {
		t.Fatal("CreateBuffer returned nil")
	}

	// Test Texture interface
	textureDesc := TextureDescriptor{
		Size:   Extent3D{Width: 256, Height: 256, DepthOrArrayLayers: 1},
		Format: TextureFormatRGBA8Unorm,
		Usage:  TextureUsageTextureBinding,
	}
	texture, err := device.CreateTexture(ctx, textureDesc)
	if err != nil {
		t.Fatalf("CreateTexture failed: %v", err)
	}
	if texture == nil {
		t.Fatal("CreateTexture returned nil")
	}

	// Test Sampler interface
	samplerDesc := SamplerDescriptor{
		AddressModeU: AddressModeClampToEdge,
		AddressModeV: AddressModeClampToEdge,
		MagFilter:    FilterModeLinear,
		MinFilter:    FilterModeLinear,
	}
	sampler, err := device.CreateSampler(ctx, samplerDesc)
	if err != nil {
		t.Fatalf("CreateSampler failed: %v", err)
	}
	if sampler == nil {
		t.Fatal("CreateSampler returned nil")
	}

	// Test ShaderModule interface
	shaderDesc := ShaderModuleDescriptor{Label: "test-shader"}
	shader, err := device.CreateShaderModule(ctx, shaderDesc)
	if err != nil {
		t.Fatalf("CreateShaderModule failed: %v", err)
	}
	if shader == nil {
		t.Fatal("CreateShaderModule returned nil")
	}

	// Test PipelineLayout interface
	pipelineLayoutDesc := PipelineLayoutDescriptor{Label: "test-layout"}
	pipelineLayout, err := device.CreatePipelineLayout(ctx, pipelineLayoutDesc)
	if err != nil {
		t.Fatalf("CreatePipelineLayout failed: %v", err)
	}
	if pipelineLayout == nil {
		t.Fatal("CreatePipelineLayout returned nil")
	}

	// Test QuerySet interface
	querySetDesc := QuerySetDescriptor{
		Type:  QueryTypeOcclusion,
		Count: 16,
	}
	querySet, err := device.CreateQuerySet(ctx, querySetDesc)
	if err != nil {
		t.Fatalf("CreateQuerySet failed: %v", err)
	}
	if querySet == nil {
		t.Fatal("CreateQuerySet returned nil")
	}

	// Test BindGroupLayout interface
	bindGroupLayoutDesc := BindGroupLayoutDescriptor{Label: "test-bind-group-layout"}
	bindGroupLayout, err := device.CreateBindGroupLayout(ctx, bindGroupLayoutDesc)
	if err != nil {
		t.Fatalf("CreateBindGroupLayout failed: %v", err)
	}
	if bindGroupLayout == nil {
		t.Fatal("CreateBindGroupLayout returned nil")
	}

	// Test Queue interface
	queue, err := device.GetQueue(ctx)
	if err != nil {
		t.Fatalf("GetQueue failed: %v", err)
	}
	if queue == nil {
		t.Fatal("GetQueue returned nil")
	}

	// Test CommandEncoder interface
	commandEncoderDesc := CommandEncoderDescriptor{Label: "test-encoder"}
	commandEncoder, err := device.CreateCommandEncoder(ctx, commandEncoderDesc)
	if err != nil {
		t.Fatalf("CreateCommandEncoder failed: %v", err)
	}
	if commandEncoder == nil {
		t.Fatal("CreateCommandEncoder returned nil")
	}

	t.Log("All interface implementations verified successfully")
}

// TestReferenceCountingPattern verifies that reference counting works correctly
func TestReferenceCountingPattern(t *testing.T) {
	ctx := context.Background()

	// Test with Buffer implementation
	bufferDesc := BufferDescriptor{Size: 1024, Usage: BufferUsageVertex}
	buffer := NewBuffer(bufferDesc)

	// Initial reference count should be 1
	err := buffer.AddRef(ctx)
	if err != nil {
		t.Fatalf("AddRef failed: %v", err)
	}

	// Release once (should still be alive)
	err = buffer.Release(ctx)
	if err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	// Release again (should destroy)
	err = buffer.Release(ctx)
	if err != nil {
		t.Fatalf("Final release failed: %v", err)
	}

	// Further operations should fail
	err = buffer.AddRef(ctx)
	if err == nil {
		t.Fatal("AddRef should have failed on destroyed buffer")
	}

	t.Log("Reference counting pattern verified successfully")
}

// TestAsyncOperations verifies that Future-based operations work
func TestAsyncOperations(t *testing.T) {
	ctx := context.Background()

	// Test GPU creation and instance creation
	instanceDesc := InstanceDescriptor{}
	instance, err := gpu.CreateInstance(ctx, instanceDesc)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	// Test async adapter request
	adapterOptions := RequestAdapterOptions{}
	future := instance.RequestAdapter(ctx, adapterOptions)
	if future.Id == 0 {
		t.Fatal("RequestAdapter returned invalid future ID")
	}

	// Test shader compilation info future
	device := &device{refCount: 1}
	shaderDesc := ShaderModuleDescriptor{Label: "test-shader"}
	shader, err := device.CreateShaderModule(ctx, shaderDesc)
	if err != nil {
		t.Fatalf("CreateShaderModule failed: %v", err)
	}

	compilationFuture := shader.GetCompilationInfo(ctx)
	if compilationFuture.Id == 0 {
		t.Fatal("GetCompilationInfo returned invalid future ID")
	}

	t.Log("Async operations verified successfully")
}
