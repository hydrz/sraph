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
	refCount  int32
	label     string
	features  []FeatureName
	limits    Limits
	queue     Queue
	destroyed bool
}

// newDevice creates a new WebGPU device
func newDevice(descriptor DeviceDescriptor) Device {
	d := &device{
		refCount:  1,
		label:     descriptor.Label,
		features:  descriptor.RequiredFeatures,
		limits:    descriptor.RequiredLimits,
		queue:     newQueue(descriptor.DefaultQueue),
		destroyed: false,
	}
	return d
}

// CreateBindGroup creates a new bind group
func (d *device) CreateBindGroup(descriptor BindGroupDescriptor) (BindGroup, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	bindGroup := &bindGroup{
		refCount: 1,
		label:    descriptor.Label,
		layout:   descriptor.Layout,
		entries:  descriptor.Entries,
	}

	return bindGroup, nil
}

// CreateBindGroupLayout creates a new bind group layout
func (d *device) CreateBindGroupLayout(descriptor BindGroupLayoutDescriptor) (BindGroupLayout, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	layout := &bindGroupLayout{
		refCount: 1,
		label:    descriptor.Label,
		entries:  descriptor.Entries,
	}

	return layout, nil
}

// CreateBuffer creates a new buffer
func (d *device) CreateBuffer(descriptor BufferDescriptor) (Buffer, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	buffer := &buffer{
		refCount:         1,
		label:            descriptor.Label,
		usage:            descriptor.Usage,
		size:             descriptor.Size,
		mappedAtCreation: descriptor.MappedAtCreation,
		mapState:         BufferMapStateUnmapped,
	}

	if descriptor.MappedAtCreation {
		buffer.mapState = BufferMapStateMapped
		buffer.data = make([]byte, descriptor.Size)
	}

	return buffer, nil
}

// CreateCommandEncoder creates a new command encoder
func (d *device) CreateCommandEncoder(descriptor CommandEncoderDescriptor) (CommandEncoder, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	encoder := &commandEncoder{
		refCount: 1,
		label:    descriptor.Label,
	}

	return encoder, nil
}

// CreateComputePipeline creates a new compute pipeline
func (d *device) CreateComputePipeline(descriptor ComputePipelineDescriptor) (ComputePipeline, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	pipeline := &ComputePipelineImpl{
		refCount: 1,
		label:    descriptor.Label,
		layout:   descriptor.Layout,
		compute:  descriptor.Compute,
	}

	return pipeline, nil
}

// CreateComputePipelineAsync creates a compute pipeline asynchronously
func (d *device) CreateComputePipelineAsync(descriptor ComputePipelineDescriptor, callback CreateComputePipelineAsyncCallbackInfo) Future {
	d.mu.RLock()
	defer d.mu.RUnlock()

	future := Future{
		Id: generateFutureId(),
	}

	// In a real implementation, this would create the pipeline asynchronously
	go func() {
		pipeline := &ComputePipelineImpl{
			refCount: 1,
			label:    descriptor.Label,
			layout:   descriptor.Layout,
			compute:  descriptor.Compute,
		}
		_ = pipeline
	}()

	return future
}

// CreatePipelineLayout creates a new pipeline layout
func (d *device) CreatePipelineLayout(descriptor PipelineLayoutDescriptor) (PipelineLayout, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	layout := &pipelineLayout{
		refCount:         1,
		label:            descriptor.Label,
		bindGroupLayouts: descriptor.BindGroupLayouts,
	}

	return layout, nil
}

// CreateQuerySet creates a new query set
func (d *device) CreateQuerySet(descriptor QuerySetDescriptor) (QuerySet, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	querySet := &querySet{
		refCount: 1,
		label:    descriptor.Label,
		qType:    descriptor.Type,
		count:    descriptor.Count,
	}

	return querySet, nil
}

// CreateRenderBundleEncoder creates a new render bundle encoder
func (d *device) CreateRenderBundleEncoder(descriptor RenderBundleEncoderDescriptor) (RenderBundleEncoder, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	encoder := &renderBundleEncoder{
		refCount:           1,
		label:              descriptor.Label,
		colorFormats:       descriptor.ColorFormats,
		depthStencilFormat: descriptor.DepthStencilFormat,
		sampleCount:        descriptor.SampleCount,
	}

	return encoder, nil
}

// CreateRenderPipeline creates a new render pipeline
func (d *device) CreateRenderPipeline(descriptor RenderPipelineDescriptor) (RenderPipeline, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	pipeline := &RenderPipelineImpl{
		refCount:     1,
		label:        descriptor.Label,
		layout:       descriptor.Layout,
		vertex:       descriptor.Vertex,
		primitive:    descriptor.Primitive,
		depthStencil: descriptor.DepthStencil,
		multisample:  descriptor.Multisample,
		fragment:     descriptor.Fragment,
	}

	return pipeline, nil
}

// CreateRenderPipelineAsync creates a render pipeline asynchronously
func (d *device) CreateRenderPipelineAsync(descriptor RenderPipelineDescriptor, callback CreateRenderPipelineAsyncCallbackInfo) Future {
	d.mu.RLock()
	defer d.mu.RUnlock()

	future := Future{
		Id: generateFutureId(),
	}

	// In a real implementation, this would create the pipeline asynchronously
	go func() {
		pipeline := &RenderPipelineImpl{
			refCount:     1,
			label:        descriptor.Label,
			layout:       descriptor.Layout,
			vertex:       descriptor.Vertex,
			primitive:    descriptor.Primitive,
			depthStencil: descriptor.DepthStencil,
			multisample:  descriptor.Multisample,
			fragment:     descriptor.Fragment,
		}
		_ = pipeline
	}()

	return future
}

// CreateSampler creates a new sampler
func (d *device) CreateSampler(descriptor SamplerDescriptor) (Sampler, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	sampler := &sampler{
		refCount:      1,
		label:         descriptor.Label,
		addressModeU:  descriptor.AddressModeU,
		addressModeV:  descriptor.AddressModeV,
		addressModeW:  descriptor.AddressModeW,
		magFilter:     descriptor.MagFilter,
		minFilter:     descriptor.MinFilter,
		mipmapFilter:  descriptor.MipmapFilter,
		lodMinClamp:   descriptor.LodMinClamp,
		lodMaxClamp:   descriptor.LodMaxClamp,
		compare:       descriptor.Compare,
		maxAnisotropy: descriptor.MaxAnisotropy,
	}

	return sampler, nil
}

// CreateShaderModule creates a new shader module
func (d *device) CreateShaderModule(descriptor ShaderModuleDescriptor) (ShaderModule, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	shaderModule := &shaderModule{
		refCount: 1,
		label:    descriptor.Label,
	}

	return shaderModule, nil
}

// CreateTexture creates a new texture
func (d *device) CreateTexture(descriptor TextureDescriptor) (Texture, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	texture := &texture{
		refCount:      1,
		label:         descriptor.Label,
		usage:         descriptor.Usage,
		dimension:     descriptor.Dimension,
		size:          descriptor.Size,
		format:        descriptor.Format,
		mipLevelCount: descriptor.MipLevelCount,
		sampleCount:   descriptor.SampleCount,
		viewFormats:   descriptor.ViewFormats,
	}

	return texture, nil
}

// Destroy destroys the device
func (d *device) Destroy() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.destroyed = true
	return nil
}

// GetAdapterInfo gets adapter information
func (d *device) GetAdapterInfo(adapterInfo AdapterInfo) (Status, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return StatusError, fmt.Errorf("device has been destroyed")
	}

	// Return default adapter info
	adapterInfo.Vendor = "WebGPU Implementation"
	adapterInfo.Device = "Generic Device"
	adapterInfo.Description = "WebGPU Device"

	return StatusSuccess, nil
}

// GetFeatures retrieves supported features
func (d *device) GetFeatures(features SupportedFeatures) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return fmt.Errorf("device has been destroyed")
	}

	features.Features = append(features.Features, d.features...)
	return nil
}

// GetLimits retrieves device limits
func (d *device) GetLimits(limits Limits) (Status, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return StatusError, fmt.Errorf("device has been destroyed")
	}

	limits = d.limits
	return StatusSuccess, nil
}

// GetLostFuture gets the device lost future
func (d *device) GetLostFuture() (Future, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return Future{}, fmt.Errorf("device has been destroyed")
	}

	future := Future{
		Id: generateFutureId(),
	}

	return future, nil
}

// GetQueue gets the device queue
func (d *device) GetQueue() (Queue, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return nil, fmt.Errorf("device has been destroyed")
	}

	return d.queue, nil
}

// HasFeature checks if a feature is supported
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

// PopErrorScope pops an error scope
func (d *device) PopErrorScope(callback PopErrorScopeCallbackInfo) Future {
	d.mu.RLock()
	defer d.mu.RUnlock()

	future := Future{
		Id: generateFutureId(),
	}

	// In a real implementation, this would pop error scope asynchronously
	return future
}

// PushErrorScope pushes an error scope
func (d *device) PushErrorScope(filter ErrorFilter) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.destroyed {
		return fmt.Errorf("device has been destroyed")
	}

	// In a real implementation, this would push error scope
	return nil
}

// SetLabel sets the device label
func (d *device) SetLabel(label string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.destroyed {
		return fmt.Errorf("device has been destroyed")
	}

	d.label = label
	return nil
}

// AddRef increments the reference count
func (d *device) AddRef() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.destroyed {
		return fmt.Errorf("device has been destroyed")
	}

	d.refCount++
	return nil
}

// Release decrements the reference count and destroys if zero
func (d *device) Release() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.destroyed {
		return fmt.Errorf("device has been destroyed")
	}

	d.refCount--
	if d.refCount <= 0 {
		d.destroyed = true
	}

	return nil
}
