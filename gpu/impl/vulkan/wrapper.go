package vulkan

import (
	"fmt"

	"github.com/opensraph/sraph/gpu/impl"
	. "github.com/opensraph/sraph/gpu/webgpu"
)

// VulkanDeviceWrapper wraps VulkanDevice to implement WebGPU Device interface
type VulkanDeviceWrapper struct {
	vulkanDevice *VulkanDevice
	baseDevice   Device // Embed base device functionality
}

// NewVulkanDeviceWrapper creates a wrapper around VulkanDevice
func NewVulkanDeviceWrapper(vulkanDevice *VulkanDevice) Device {
	// Create base device descriptor from Vulkan device
	descriptor := DeviceDescriptor{
		Label: "Vulkan Device",
		RequiredFeatures: []FeatureName{
			FeatureNameDepthClipControl,
			FeatureNameTimestampQuery,
		},
		RequiredLimits: getVulkanDeviceLimits(vulkanDevice),
		DefaultQueue: QueueDescriptor{
			Label: "Vulkan Queue",
		},
	}

	baseDevice := impl.NewDevice(descriptor)

	return &VulkanDeviceWrapper{
		vulkanDevice: vulkanDevice,
		baseDevice:   baseDevice,
	}
}

// CreateBuffer creates a Vulkan-backed buffer
func (vdw *VulkanDeviceWrapper) CreateBuffer(descriptor BufferDescriptor) (Buffer, error) {
	if vdw.vulkanDevice.IsDestroyed() {
		return nil, fmt.Errorf("Vulkan device has been destroyed")
	}

	// Create Vulkan buffer
	vulkanBuffer, err := NewVulkanBuffer(vdw.vulkanDevice, descriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vulkan buffer: %v", err)
	}

	// Wrap in WebGPU Buffer interface
	return NewVulkanBufferWrapper(vulkanBuffer, descriptor), nil
}

// CreateTexture creates a Vulkan-backed texture
func (vdw *VulkanDeviceWrapper) CreateTexture(descriptor TextureDescriptor) (Texture, error) {
	if vdw.vulkanDevice.IsDestroyed() {
		return nil, fmt.Errorf("Vulkan device has been destroyed")
	}

	// Create Vulkan texture
	vulkanTexture, err := NewVulkanTexture(vdw.vulkanDevice, descriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vulkan texture: %v", err)
	}

	// Wrap in WebGPU Texture interface
	return NewVulkanTextureWrapper(vulkanTexture, descriptor), nil
}

// CreateSampler creates a Vulkan-backed sampler
func (vdw *VulkanDeviceWrapper) CreateSampler(descriptor SamplerDescriptor) (Sampler, error) {
	if vdw.vulkanDevice.IsDestroyed() {
		return nil, fmt.Errorf("Vulkan device has been destroyed")
	}

	// Create Vulkan sampler
	vulkanSampler, err := NewVulkanSampler(vdw.vulkanDevice, descriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vulkan sampler: %v", err)
	}

	// Wrap in WebGPU Sampler interface
	return NewVulkanSamplerWrapper(vulkanSampler, descriptor), nil
}

// Delegate other methods to base device
func (vdw *VulkanDeviceWrapper) CreateBindGroup(descriptor BindGroupDescriptor) (BindGroup, error) {
	return vdw.baseDevice.CreateBindGroup(descriptor)
}

func (vdw *VulkanDeviceWrapper) CreateBindGroupLayout(descriptor BindGroupLayoutDescriptor) (BindGroupLayout, error) {
	return vdw.baseDevice.CreateBindGroupLayout(descriptor)
}

func (vdw *VulkanDeviceWrapper) CreateCommandEncoder(descriptor CommandEncoderDescriptor) (CommandEncoder, error) {
	return vdw.baseDevice.CreateCommandEncoder(descriptor)
}

func (vdw *VulkanDeviceWrapper) CreateComputePipeline(descriptor ComputePipelineDescriptor) (ComputePipeline, error) {
	return vdw.baseDevice.CreateComputePipeline(descriptor)
}

func (vdw *VulkanDeviceWrapper) CreateComputePipelineAsync(descriptor ComputePipelineDescriptor, callback CreateComputePipelineAsyncCallbackInfo) Future {
	return vdw.baseDevice.CreateComputePipelineAsync(descriptor, callback)
}

func (vdw *VulkanDeviceWrapper) CreatePipelineLayout(descriptor PipelineLayoutDescriptor) (PipelineLayout, error) {
	return vdw.baseDevice.CreatePipelineLayout(descriptor)
}

func (vdw *VulkanDeviceWrapper) CreateQuerySet(descriptor QuerySetDescriptor) (QuerySet, error) {
	return vdw.baseDevice.CreateQuerySet(descriptor)
}

func (vdw *VulkanDeviceWrapper) CreateRenderBundleEncoder(descriptor RenderBundleEncoderDescriptor) (RenderBundleEncoder, error) {
	return vdw.baseDevice.CreateRenderBundleEncoder(descriptor)
}

func (vdw *VulkanDeviceWrapper) CreateRenderPipeline(descriptor RenderPipelineDescriptor) (RenderPipeline, error) {
	return vdw.baseDevice.CreateRenderPipeline(descriptor)
}

func (vdw *VulkanDeviceWrapper) CreateRenderPipelineAsync(descriptor RenderPipelineDescriptor, callback CreateRenderPipelineAsyncCallbackInfo) Future {
	return vdw.baseDevice.CreateRenderPipelineAsync(descriptor, callback)
}

func (vdw *VulkanDeviceWrapper) CreateShaderModule(descriptor ShaderModuleDescriptor) (ShaderModule, error) {
	return vdw.baseDevice.CreateShaderModule(descriptor)
}

func (vdw *VulkanDeviceWrapper) Destroy() error {
	err1 := vdw.vulkanDevice.Destroy()
	err2 := vdw.baseDevice.Destroy()

	if err1 != nil {
		return err1
	}
	return err2
}

// ...existing code... (delegate remaining methods to baseDevice)

func (vdw *VulkanDeviceWrapper) GetAdapterInfo(adapterInfo AdapterInfo) (Status, error) {
	return vdw.baseDevice.GetAdapterInfo(adapterInfo)
}

func (vdw *VulkanDeviceWrapper) GetFeatures(features SupportedFeatures) error {
	return vdw.baseDevice.GetFeatures(features)
}

func (vdw *VulkanDeviceWrapper) GetLimits(limits Limits) (Status, error) {
	return vdw.baseDevice.GetLimits(limits)
}

func (vdw *VulkanDeviceWrapper) GetLostFuture() (Future, error) {
	return vdw.baseDevice.GetLostFuture()
}

func (vdw *VulkanDeviceWrapper) GetQueue() (Queue, error) {
	return vdw.baseDevice.GetQueue()
}

func (vdw *VulkanDeviceWrapper) HasFeature(feature FeatureName) (bool, error) {
	return vdw.baseDevice.HasFeature(feature)
}

func (vdw *VulkanDeviceWrapper) PopErrorScope(callback PopErrorScopeCallbackInfo) Future {
	return vdw.baseDevice.PopErrorScope(callback)
}

func (vdw *VulkanDeviceWrapper) PushErrorScope(filter ErrorFilter) error {
	return vdw.baseDevice.PushErrorScope(filter)
}

func (vdw *VulkanDeviceWrapper) SetLabel(label string) error {
	return vdw.baseDevice.SetLabel(label)
}

func (vdw *VulkanDeviceWrapper) AddRef() error {
	return vdw.baseDevice.AddRef()
}

func (vdw *VulkanDeviceWrapper) Release() error {
	return vdw.baseDevice.Release()
}

// Helper function to convert Vulkan device limits
func getVulkanDeviceLimits(vulkanDevice *VulkanDevice) Limits {
	physicalDevice := vulkanDevice.GetPhysicalDevice()

	return Limits{
		MaxTextureDimension2D:   physicalDevice.Properties.MaxTextureSize,
		MaxBufferSize:           physicalDevice.Properties.MaxBufferSize,
		MaxBindGroups:           8,
		MaxBindingsPerBindGroup: 1000,
		// ...existing code... (other default limits)
	}
}
