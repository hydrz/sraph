package render

// Snapshot represents a snapshot of the current render state
// This is useful for saving and restoring render state during complex rendering operations
type Snapshot interface {
	// GetRenderTarget returns the render target at the time of snapshot
	GetRenderTarget() RenderTarget

	// GetViewport returns the viewport at the time of snapshot
	GetViewport() Viewport

	// GetScissorRect returns the scissor rectangle at the time of snapshot
	GetScissorRect() Rect

	// GetPipeline returns the active pipeline at the time of snapshot
	GetPipeline() Pipeline

	// IsValid returns true if the snapshot is valid
	IsValid() bool

	// GetTimestamp returns the timestamp when the snapshot was taken
	GetTimestamp() int64
}

// Viewport represents a viewport rectangle
type Viewport struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
	MinZ   float32
	MaxZ   float32
}

// Rect represents a rectangle
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// SnapshotImpl is the default implementation of Snapshot
type SnapshotImpl struct {
	renderTarget RenderTarget
	viewport     Viewport
	scissorRect  Rect
	pipeline     Pipeline
	timestamp    int64
	isValid      bool
}

// NewSnapshot creates a new snapshot with the current render state
func NewSnapshot(renderTarget RenderTarget, viewport Viewport, scissorRect Rect, pipeline Pipeline) Snapshot {
	return &SnapshotImpl{
		renderTarget: renderTarget,
		viewport:     viewport,
		scissorRect:  scissorRect,
		pipeline:     pipeline,
		timestamp:    getCurrentTimestamp(),
		isValid:      true,
	}
}

// GetRenderTarget returns the render target at the time of snapshot
func (s *SnapshotImpl) GetRenderTarget() RenderTarget {
	return s.renderTarget
}

// GetViewport returns the viewport at the time of snapshot
func (s *SnapshotImpl) GetViewport() Viewport {
	return s.viewport
}

// GetScissorRect returns the scissor rectangle at the time of snapshot
func (s *SnapshotImpl) GetScissorRect() Rect {
	return s.scissorRect
}

// GetPipeline returns the active pipeline at the time of snapshot
func (s *SnapshotImpl) GetPipeline() Pipeline {
	return s.pipeline
}

// IsValid returns true if the snapshot is valid
func (s *SnapshotImpl) IsValid() bool {
	return s.isValid
}

// GetTimestamp returns the timestamp when the snapshot was taken
func (s *SnapshotImpl) GetTimestamp() int64 {
	return s.timestamp
}

// getCurrentTimestamp returns the current timestamp
func getCurrentTimestamp() int64 {
	// TODO: Implement proper timestamp generation
	return 0
}

// SnapshotManager manages snapshots for efficient state management
type SnapshotManager interface {
	// TakeSnapshot takes a snapshot of the current render state
	TakeSnapshot(renderTarget RenderTarget, viewport Viewport, scissorRect Rect, pipeline Pipeline) Snapshot

	// RestoreSnapshot restores the render state from a snapshot
	RestoreSnapshot(snapshot Snapshot) error

	// GetSnapshotCount returns the number of snapshots managed
	GetSnapshotCount() int

	// ClearSnapshots clears all managed snapshots
	ClearSnapshots()
}

// SnapshotManagerImpl is the default implementation of SnapshotManager
type SnapshotManagerImpl struct {
	snapshots    []Snapshot
	maxSnapshots int
}

// NewSnapshotManager creates a new snapshot manager
func NewSnapshotManager(maxSnapshots int) SnapshotManager {
	return &SnapshotManagerImpl{
		snapshots:    make([]Snapshot, 0),
		maxSnapshots: maxSnapshots,
	}
}

// TakeSnapshot takes a snapshot of the current render state
func (sm *SnapshotManagerImpl) TakeSnapshot(renderTarget RenderTarget, viewport Viewport, scissorRect Rect, pipeline Pipeline) Snapshot {
	snapshot := NewSnapshot(renderTarget, viewport, scissorRect, pipeline)

	// Add to managed snapshots
	if len(sm.snapshots) >= sm.maxSnapshots {
		// Remove oldest snapshot if at max capacity
		sm.snapshots = sm.snapshots[1:]
	}
	sm.snapshots = append(sm.snapshots, snapshot)

	return snapshot
}

// RestoreSnapshot restores the render state from a snapshot
func (sm *SnapshotManagerImpl) RestoreSnapshot(snapshot Snapshot) error {
	if snapshot == nil || !snapshot.IsValid() {
		return ErrInvalidArgument
	}

	// TODO: Implement render state restoration
	return nil
}

// GetSnapshotCount returns the number of snapshots managed
func (sm *SnapshotManagerImpl) GetSnapshotCount() int {
	return len(sm.snapshots)
}

// ClearSnapshots clears all managed snapshots
func (sm *SnapshotManagerImpl) ClearSnapshots() {
	sm.snapshots = sm.snapshots[:0]
}
