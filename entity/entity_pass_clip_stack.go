// Package entity provides a complete entity system for Sraph's rendering pipeline.
// This package is inspired by Impeller's entity system and provides:
//
// - Entity pass management with clip stacks
// - Render target caching and management
// - Content contexts and inline pass contexts
// - Save layer utilities for advanced rendering
//
// The entity system handles the transformation from high-level drawing operations
// to low-level GPU commands, managing state, resources, and optimization along the way.
package entity

import (
	"github.com/opensraph/sraph/geom"
)

// ClipCoverageLayer represents a layer of clip coverage information
type ClipCoverageLayer struct {
	Coverage   *geom.Rect // Optional coverage rect
	ClipHeight uint32     // Height of the clip in the stencil buffer
}

// EntityPassClipStack tracks all clips recorded in the current entity pass stencil.
// These clips are replayed when restoring the backdrop so that the stencil buffer
// is left in an identical state.
type EntityPassClipStack interface {
	// CurrentClipCoverage returns the current clip coverage area
	CurrentClipCoverage() *geom.Rect

	// PushSubpass pushes a new subpass with optional coverage and clip height
	PushSubpass(subpassCoverage *geom.Rect, clipHeight uint32)

	// PopSubpass pops the current subpass
	PopSubpass()

	// HasCoverage returns true if there is any coverage
	HasCoverage() bool

	// RecordClip records a clip operation and returns clip state result
	RecordClip(clipContents Contents, transform geom.Matrix, globalPassPosition geom.Point,
		clipDepth uint32, clipHeightFloor uint32, isAA bool) *ClipStateResult

	// RecordRestore records a restore operation
	RecordRestore(globalPassPosition geom.Point, restoreHeight uint32) *ClipStateResult

	// GetReplayEntities returns all entities that need to be replayed
	GetReplayEntities() []*ReplayResult
}

// ReplayResult contains information needed to replay a clip operation
type ReplayResult struct {
	ClipContents Contents    // The clip contents to replay
	Transform    geom.Matrix // Transform to apply
	ClipCoverage *geom.Rect  // Optional clip coverage
	ClipDepth    uint32      // Depth in the clip stack
}

// ClipStateResult indicates the result of applying clip state
type ClipStateResult struct {
	ShouldRender  bool // Whether the entity should be rendered
	ClipDidChange bool // Whether the clip coverage changed
}

// SubpassState tracks the state of a rendering subpass
type SubpassState struct {
	RenderedClipEntities []*ReplayResult      // Clip entities to be replayed
	ClipCoverage         []*ClipCoverageLayer // Coverage layers for this subpass
}

// DefaultEntityPassClipStack provides a default implementation of EntityPassClipStack
type DefaultEntityPassClipStack struct {
	subpassState    []*SubpassState // Stack of subpass states
	nextReplayIndex uint32          // Index for next replay operation
	initialCoverage geom.Rect       // Initial coverage rectangle
}

// NewEntityPassClipStack creates a new entity pass clip stack with initial coverage
func NewEntityPassClipStack(initialCoverageRect geom.Rect) EntityPassClipStack {
	stack := &DefaultEntityPassClipStack{
		subpassState:    make([]*SubpassState, 0),
		nextReplayIndex: 0,
		initialCoverage: initialCoverageRect,
	}

	// Initialize with root subpass
	stack.subpassState = append(stack.subpassState, &SubpassState{
		RenderedClipEntities: make([]*ReplayResult, 0),
		ClipCoverage: []*ClipCoverageLayer{
			{
				Coverage:   &initialCoverageRect,
				ClipHeight: 0,
			},
		},
	})

	return stack
}

// CurrentClipCoverage returns the current clip coverage area
func (s *DefaultEntityPassClipStack) CurrentClipCoverage() *geom.Rect {
	if len(s.subpassState) == 0 {
		return nil
	}

	currentState := s.getCurrentSubpassState()
	if len(currentState.ClipCoverage) == 0 {
		return nil
	}

	// Return the coverage of the topmost layer
	return currentState.ClipCoverage[len(currentState.ClipCoverage)-1].Coverage
}

// PushSubpass pushes a new subpass with optional coverage and clip height
func (s *DefaultEntityPassClipStack) PushSubpass(subpassCoverage *geom.Rect, clipHeight uint32) {
	newState := &SubpassState{
		RenderedClipEntities: make([]*ReplayResult, 0),
		ClipCoverage: []*ClipCoverageLayer{
			{
				Coverage:   subpassCoverage,
				ClipHeight: clipHeight,
			},
		},
	}
	s.subpassState = append(s.subpassState, newState)
}

// PopSubpass pops the current subpass
func (s *DefaultEntityPassClipStack) PopSubpass() {
	if len(s.subpassState) > 1 {
		s.subpassState = s.subpassState[:len(s.subpassState)-1]
	}
}

// HasCoverage returns true if there is any coverage
func (s *DefaultEntityPassClipStack) HasCoverage() bool {
	coverage := s.CurrentClipCoverage()
	return coverage != nil && !coverage.IsEmpty()
}

// RecordClip records a clip operation and returns clip state result
func (s *DefaultEntityPassClipStack) RecordClip(clipContents Contents, transform geom.Matrix,
	globalPassPosition geom.Point, clipDepth uint32, clipHeightFloor uint32, isAA bool) *ClipStateResult {

	// Create replay result
	replayResult := &ReplayResult{
		ClipContents: clipContents,
		Transform:    transform,
		ClipCoverage: s.CurrentClipCoverage(),
		ClipDepth:    clipDepth,
	}

	// Add to current subpass
	currentState := s.getCurrentSubpassState()
	currentState.RenderedClipEntities = append(currentState.RenderedClipEntities, replayResult)

	// TODO: Implement proper clip coverage computation
	return &ClipStateResult{
		ShouldRender:  true,
		ClipDidChange: true,
	}
}

// RecordRestore records a restore operation
func (s *DefaultEntityPassClipStack) RecordRestore(globalPassPosition geom.Point, restoreHeight uint32) *ClipStateResult {
	// TODO: Implement restore logic
	return &ClipStateResult{
		ShouldRender:  true,
		ClipDidChange: false,
	}
}

// GetReplayEntities returns all entities that need to be replayed
func (s *DefaultEntityPassClipStack) GetReplayEntities() []*ReplayResult {
	var allEntities []*ReplayResult

	for _, state := range s.subpassState {
		allEntities = append(allEntities, state.RenderedClipEntities...)
	}

	return allEntities
}

// getCurrentSubpassState returns the current subpass state
func (s *DefaultEntityPassClipStack) getCurrentSubpassState() *SubpassState {
	if len(s.subpassState) == 0 {
		// This should not happen, but create emergency state
		s.subpassState = append(s.subpassState, &SubpassState{
			RenderedClipEntities: make([]*ReplayResult, 0),
			ClipCoverage:         make([]*ClipCoverageLayer, 0),
		})
	}

	return s.subpassState[len(s.subpassState)-1]
}
