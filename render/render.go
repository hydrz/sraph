// Package render provides the rendering layer for Sraph, inspired by Flutter's rendering system
// and WebGPU/Impeller architecture.
//
// This package provides the core rendering interfaces that bridge between the high-level
// display list and the low-level GPU backend.
package render

import (
	"github.com/opensraph/sraph/display"
	"github.com/opensraph/sraph/geom"
	"github.com/opensraph/sraph/gpu"
)

// Renderer represents the main rendering interface, inspired by Impeller's renderer.
// It manages the rendering context and coordinates with the GPU backend.
type Renderer interface {
	// Render renders a display list to the given surface.
	Render(surface *Surface, displayList *display.DisplayList) error

	// CreateSurface creates a new rendering surface.
	CreateSurface(desc SurfaceDescriptor) (*Surface, error)

	// GetDevice returns the underlying GPU device.
	GetDevice() gpu.Device

	// GetCapabilities returns the renderer capabilities.
	GetCapabilities() RendererCapabilities
}

// Surface represents a render target, similar to Impeller's Surface concept.
type Surface struct {
	texture     gpu.Texture
	textureView gpu.TextureView
	size        geom.Size[geom.F32]
	format      gpu.TextureFormat
	sampleCount uint32
}

// SurfaceDescriptor describes how to create a surface.
type SurfaceDescriptor struct {
	Size        geom.Size[geom.F32]
	Format      gpu.TextureFormat
	Usage       gpu.TextureUsage
	SampleCount uint32
}

// NewSurface creates a new surface.
func NewSurface(device gpu.Device, desc SurfaceDescriptor) (*Surface, error) {
	textureDesc := gpu.TextureDescriptor{
		Size: gpu.Extent3D{
			Width:              uint32(desc.Size.Width),
			Height:             uint32(desc.Size.Height),
			DepthOrArrayLayers: 1,
		},
		Format:        desc.Format,
		Usage:         desc.Usage,
		SampleCount:   desc.SampleCount,
		MipLevelCount: 1,
	}

	texture := device.CreateTexture(textureDesc)
	textureView := texture.CreateView(gpu.TextureViewDescriptor{
		Format: desc.Format,
	})

	return &Surface{
		texture:     texture,
		textureView: textureView,
		size:        desc.Size,
		format:      desc.Format,
		sampleCount: desc.SampleCount,
	}, nil
}

// GetTexture returns the underlying texture.
func (s *Surface) GetTexture() gpu.Texture {
	return s.texture
}

// GetTextureView returns the texture view for rendering.
func (s *Surface) GetTextureView() gpu.TextureView {
	return s.textureView
}

// GetSize returns the surface size.
func (s *Surface) GetSize() geom.Size[geom.F32] {
	return s.size
}

// RendererCapabilities describes what the renderer can do.
type RendererCapabilities struct {
	SupportsAdvancedBlends bool
	MaxTextureSize         uint32
	SupportsCompute        bool
	SupportsTimestamps     bool
}

// ContentContext provides rendering context for entity contents.
// This mirrors Impeller's ContentContext.
type ContentContext struct {
	device          gpu.Device
	queue           gpu.Queue
	renderPipelines map[string]gpu.RenderPipeline
	bindGroups      map[string]gpu.BindGroup
	uniformBuffers  map[string]gpu.Buffer
}

// NewContentContext creates a new content context.
func NewContentContext(device gpu.Device) *ContentContext {
	return &ContentContext{
		device:          device,
		queue:           device.Queue(),
		renderPipelines: make(map[string]gpu.RenderPipeline),
		bindGroups:      make(map[string]gpu.BindGroup),
		uniformBuffers:  make(map[string]gpu.Buffer),
	}
}

// GetDevice returns the GPU device.
func (c *ContentContext) GetDevice() gpu.Device {
	return c.device
}

// GetQueue returns the GPU queue.
func (c *ContentContext) GetQueue() gpu.Queue {
	return c.queue
}

// ImpellerRenderer implements the Renderer interface using Impeller-style architecture.
type ImpellerRenderer struct {
	device       gpu.Device
	context      *ContentContext
	capabilities RendererCapabilities
}

// NewImpellerRenderer creates a new Impeller-style renderer.
func NewImpellerRenderer(device gpu.Device) *ImpellerRenderer {
	return &ImpellerRenderer{
		device:  device,
		context: NewContentContext(device),
		capabilities: RendererCapabilities{
			SupportsAdvancedBlends: true, // TODO: Query actual capabilities
			MaxTextureSize:         8192, // TODO: Query actual capabilities
			SupportsCompute:        true, // TODO: Query actual capabilities
			SupportsTimestamps:     true, // TODO: Query actual capabilities
		},
	}
}

// Render implements Renderer.
func (r *ImpellerRenderer) Render(surface *Surface, displayList *display.DisplayList) error {
	commandEncoder := r.device.CreateCommandEncoder(gpu.CommandEncoderDescriptor{
		Label: "Render Pass",
	})

	renderPass := commandEncoder.BeginRenderPass(gpu.RenderPassDescriptor{
		ColorAttachments: []gpu.RenderPassColorAttachment{
			{
				View:       surface.GetTextureView(),
				LoadOp:     gpu.LoadOpClear,
				StoreOp:    gpu.StoreOpStore,
				ClearValue: gpu.Color{R: 0.0, G: 0.0, B: 0.0, A: 1.0},
			},
		},
	})

	// TODO: Properly integrate with display list
	// For now, we just end the render pass
	renderPass.End()
	commandBuffer := commandEncoder.Finish(gpu.CommandBufferDescriptor{})
	r.device.Queue().Submit([]gpu.CommandBuffer{commandBuffer})

	return nil
}

// CreateSurface implements Renderer.
func (r *ImpellerRenderer) CreateSurface(desc SurfaceDescriptor) (*Surface, error) {
	return NewSurface(r.device, desc)
}

// GetDevice implements Renderer.
func (r *ImpellerRenderer) GetDevice() gpu.Device {
	return r.device
}

// GetCapabilities implements Renderer.
func (r *ImpellerRenderer) GetCapabilities() RendererCapabilities {
	return r.capabilities
}

// RenderObject is an object in the render tree.
// It provides layout, painting, hit testing, and compositing behavior.
type RenderObject interface {
	// Layout computes the layout for this render object.
	Layout(constraints Constraints)

	// Paint paints this render object into the given context.
	Paint(context *PaintingContext, offset geom.Point[geom.F32])

	// HitTest determines whether the given point is within this render object.
	HitTest(result *HitTestResult, position geom.Point[geom.F32]) bool

	// Size returns the size of this render object.
	Size() geom.Size[geom.F32]

	// Parent returns the parent of this render object.
	Parent() RenderObject

	// AdoptChild adds a child to this render object.
	AdoptChild(child RenderObject)

	// DropChild removes a child from this render object.
	DropChild(child RenderObject)

	// MarkNeedsLayout marks this render object as needing layout.
	MarkNeedsLayout()

	// MarkNeedsPaint marks this render object as needing paint.
	MarkNeedsPaint()
}

// Constraints define layout constraints for render objects.
type Constraints interface {
	// IsTight returns true if there is exactly one size that satisfies the constraints.
	IsTight() bool

	// IsNormalized returns true if the constraints are in canonical form.
	IsNormalized() bool
}

// BoxConstraints describe layout constraints for rectangular render objects.
type BoxConstraints struct {
	MinWidth  geom.F32
	MaxWidth  geom.F32
	MinHeight geom.F32
	MaxHeight geom.F32
}

// PaintingContext provides a context for painting render objects.
type PaintingContext struct {
	canvas     *display.DisplayListBuilder
	clipBounds geom.Rect[geom.F32]
}

// NewPaintingContext creates a new painting context.
func NewPaintingContext(canvas *display.DisplayListBuilder, clipBounds geom.Rect[geom.F32]) *PaintingContext {
	return &PaintingContext{
		canvas:     canvas,
		clipBounds: clipBounds,
	}
}

// HitTestResult stores the results of hit testing.
type HitTestResult struct {
	// Path stores the path of render objects that were hit.
	Path []RenderObject

	// LocalPosition is the position in the coordinate system of the hit object.
	LocalPosition geom.Point[geom.F32]
}
