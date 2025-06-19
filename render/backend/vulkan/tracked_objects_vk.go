package vulkan

import (
	"sync"

	"github.com/vulkan-go/vulkan"
)

// TrackedObjectsVK manages Vulkan objects that need to be tracked for lifetime management
type TrackedObjectsVK struct {
	mutex           sync.RWMutex
	buffers         []vulkan.Buffer
	images          []vulkan.Image
	imageViews      []vulkan.ImageView
	pipelines       []vulkan.Pipeline
	descriptorSets  []vulkan.DescriptorSet
	commandBuffers  []vulkan.CommandBuffer
	fences          []vulkan.Fence
	semaphores      []vulkan.Semaphore
	framebuffers    []vulkan.Framebuffer
	renderPasses    []vulkan.RenderPass
	descriptorPools []vulkan.DescriptorPool
	pipelineLayouts []vulkan.PipelineLayout
	shaderModules   []vulkan.ShaderModule
}

// NewTrackedObjectsVK creates a new tracked objects manager
func NewTrackedObjectsVK() *TrackedObjectsVK {
	return &TrackedObjectsVK{}
}

// TrackBuffer tracks a buffer for cleanup
func (t *TrackedObjectsVK) TrackBuffer(buffer vulkan.Buffer) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.buffers = append(t.buffers, buffer)
}

// TrackImage tracks an image for cleanup
func (t *TrackedObjectsVK) TrackImage(image vulkan.Image) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.images = append(t.images, image)
}

// TrackImageView tracks an image view for cleanup
func (t *TrackedObjectsVK) TrackImageView(imageView vulkan.ImageView) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.imageViews = append(t.imageViews, imageView)
}

// TrackPipeline tracks a pipeline for cleanup
func (t *TrackedObjectsVK) TrackPipeline(pipeline vulkan.Pipeline) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.pipelines = append(t.pipelines, pipeline)
}

// TrackDescriptorSet tracks a descriptor set for cleanup
func (t *TrackedObjectsVK) TrackDescriptorSet(descriptorSet vulkan.DescriptorSet) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.descriptorSets = append(t.descriptorSets, descriptorSet)
}

// TrackCommandBuffer tracks a command buffer for cleanup
func (t *TrackedObjectsVK) TrackCommandBuffer(commandBuffer vulkan.CommandBuffer) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.commandBuffers = append(t.commandBuffers, commandBuffer)
}

// TrackFence tracks a fence for cleanup
func (t *TrackedObjectsVK) TrackFence(fence vulkan.Fence) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.fences = append(t.fences, fence)
}

// TrackSemaphore tracks a semaphore for cleanup
func (t *TrackedObjectsVK) TrackSemaphore(semaphore vulkan.Semaphore) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.semaphores = append(t.semaphores, semaphore)
}

// TrackFramebuffer tracks a framebuffer for cleanup
func (t *TrackedObjectsVK) TrackFramebuffer(framebuffer vulkan.Framebuffer) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.framebuffers = append(t.framebuffers, framebuffer)
}

// TrackRenderPass tracks a render pass for cleanup
func (t *TrackedObjectsVK) TrackRenderPass(renderPass vulkan.RenderPass) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.renderPasses = append(t.renderPasses, renderPass)
}

// TrackDescriptorPool tracks a descriptor pool for cleanup
func (t *TrackedObjectsVK) TrackDescriptorPool(descriptorPool vulkan.DescriptorPool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.descriptorPools = append(t.descriptorPools, descriptorPool)
}

// TrackPipelineLayout tracks a pipeline layout for cleanup
func (t *TrackedObjectsVK) TrackPipelineLayout(pipelineLayout vulkan.PipelineLayout) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.pipelineLayouts = append(t.pipelineLayouts, pipelineLayout)
}

// TrackShaderModule tracks a shader module for cleanup
func (t *TrackedObjectsVK) TrackShaderModule(shaderModule vulkan.ShaderModule) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.shaderModules = append(t.shaderModules, shaderModule)
}

// Clear clears all tracked objects
func (t *TrackedObjectsVK) Clear() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.buffers = t.buffers[:0]
	t.images = t.images[:0]
	t.imageViews = t.imageViews[:0]
	t.pipelines = t.pipelines[:0]
	t.descriptorSets = t.descriptorSets[:0]
	t.commandBuffers = t.commandBuffers[:0]
	t.fences = t.fences[:0]
	t.semaphores = t.semaphores[:0]
	t.framebuffers = t.framebuffers[:0]
	t.renderPasses = t.renderPasses[:0]
	t.descriptorPools = t.descriptorPools[:0]
	t.pipelineLayouts = t.pipelineLayouts[:0]
	t.shaderModules = t.shaderModules[:0]
}

// GetTrackedCount returns the total number of tracked objects
func (t *TrackedObjectsVK) GetTrackedCount() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return len(t.buffers) + len(t.images) + len(t.imageViews) +
		len(t.pipelines) + len(t.descriptorSets) + len(t.commandBuffers) +
		len(t.fences) + len(t.semaphores) + len(t.framebuffers) +
		len(t.renderPasses) + len(t.descriptorPools) + len(t.pipelineLayouts) +
		len(t.shaderModules)
}
