package impl

import (
	"fmt"
	"sync"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

var _ Device = (*device)(nil)

// device implements the Device interface
type device struct {
	mu        sync.RWMutex
	label     string
	features  []FeatureName
	limits    Limits
	queue     Queue
	destroyed bool
}

// NewDevice creates a new WebGPU device
func NewDevice(descriptor DeviceDescriptor) Device {
	d := &device{
		label:     descriptor.Label,
		features:  descriptor.RequiredFeatures,
		limits:    descriptor.RequiredLimits,
		queue:     newQueue(descriptor.DefaultQueue),
		destroyed: false,
	}
	return d
}

// CreateBindGroup implements Device.CreateBindGroup.
func (d *device) CreateBindGroup(descriptor BindGroupDescriptor) (BindGroup, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	bindGroup := &bindGroup{
		label:   descriptor.Label,
		layout:  descriptor.Layout,
		entries: descriptor.Entries,
	}

	return bindGroup, nil
}

// CreateBindGroupLayout implements Device.CreateBindGroupLayout.
func (d *device) CreateBindGroupLayout(descriptor BindGroupLayoutDescriptor) (BindGroupLayout, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	layout := &bindGroupLayout{

		label:   descriptor.Label,
		entries: descriptor.Entries,
	}

	return layout, nil
}

// CreateBuffer implements Device.CreateBuffer.
func (d *device) CreateBuffer(descriptor BufferDescriptor) (Buffer, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return NewBuffer(descriptor), nil
}

// CreateCommandEncoder implements Device.CreateCommandEncoder.
func (d *device) CreateCommandEncoder(descriptor CommandEncoderDescriptor) (CommandEncoder, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	encoder := &commandEncoder{

		label: descriptor.Label,
	}

	return encoder, nil
}

// CreateComputePipeline implements Device.CreateComputePipeline.
func (d *device) CreateComputePipeline(descriptor ComputePipelineDescriptor) (ComputePipeline, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return NewComputePipeline(descriptor), nil
}

// CreatePipelineLayout implements Device.CreatePipelineLayout.
func (d *device) CreatePipelineLayout(descriptor PipelineLayoutDescriptor) (PipelineLayout, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return NewPipelineLayout(descriptor), nil
}

// CreateQuerySet implements Device.CreateQuerySet.
func (d *device) CreateQuerySet(descriptor QuerySetDescriptor) (QuerySet, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return NewQuerySet(descriptor), nil
}

// CreateRenderBundleEncoder implements Device.CreateRenderBundleEncoder.
func (d *device) CreateRenderBundleEncoder(descriptor RenderBundleEncoderDescriptor) (RenderBundleEncoder, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return NewRenderBundleEncoder(descriptor), nil
}

// CreateRenderPipeline implements Device.CreateRenderPipeline.
func (d *device) CreateRenderPipeline(descriptor RenderPipelineDescriptor) (RenderPipeline, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return NewRenderPipeline(descriptor), nil
}

// CreateSampler implements Device.CreateSampler.
func (d *device) CreateSampler(descriptor SamplerDescriptor) (Sampler, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return NewSampler(descriptor), nil
}

// CreateShaderModule implements Device.CreateShaderModule.
func (d *device) CreateShaderModule(descriptor ShaderModuleDescriptor) (ShaderModule, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return NewShaderModule(descriptor), nil
}

// CreateTexture implements Device.CreateTexture.
func (d *device) CreateTexture(descriptor TextureDescriptor) (Texture, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return NewTexture(descriptor), nil
}

// Destroy implements Device.Destroy.
func (d *device) Destroy() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.destroyed = true
	return nil
}

// AdapterInfo implements Device.AdapterInfo.
func (d *device) AdapterInfo() (*AdapterInfo, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	info := &AdapterInfo{
		Vendor:      "WebGPU Implementation",
		Device:      "Generic Device",
		Description: "WebGPU Device",
	}
	return info, nil
}

// Features implements Device.Features.
func (d *device) Features() (*SupportedFeatures, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return &SupportedFeatures{Features: append([]FeatureName{}, d.features...)}, nil
}

// Limits implements Device.Limits.
func (d *device) Limits() (*Limits, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	limits := d.limits
	return &limits, nil
}

// LostFuture implements Device.LostFuture.
func (d *device) LostFuture() (Future, error) {
	return Future{}, fmt.Errorf("LostFuture is not implemented in this version of the device")
}

// Queue implements Device.Queue.
func (d *device) Queue() (Queue, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return d.queue, nil
}

// HasFeature implements Device.HasFeature.
func (d *device) HasFeature(feature FeatureName) (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return false, fmt.Errorf("device has been destroyed")
	}

	for _, f := range d.features {
		if f == feature {
			return true, nil
		}
	}
	return false, nil
}

// PopErrorScope implements Device.PopErrorScope.
func (d *device) PopErrorScope() error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return fmt.Errorf("device has been destroyed")
	}

	// TODO: Implement error scope stack logic if needed
	return nil
}

// PushErrorScope implements Device.PushErrorScope.
func (d *device) PushErrorScope(filter ErrorFilter) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return fmt.Errorf("device has been destroyed")
	}

	// TODO: Implement error scope stack logic if needed
	return nil
}

// SetLabel implements Device.SetLabel.
func (d *device) SetLabel(label string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.destroyed {
		return fmt.Errorf("device has been destroyed")
	}

	d.label = label
	return nil
}
