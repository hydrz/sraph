package display

import (
	"github.com/opensraph/sraph/geom"
)

// DisplayList represents a recorded sequence of graphics operations.
// It provides an immutable sequence of drawing commands that can be
// played back on different canvases or contexts.
type DisplayList struct {
	// operations contains the sequence of recorded operations
	operations []Operation
	// bounds contains the bounding rectangle of all operations
	bounds geom.Rect[Scalar]
	// totalOpCount is the total number of operations in the list
	totalOpCount int
	// hasAntiAliasing indicates if any operations use anti-aliasing
	hasAntiAliasing bool
	// hasText indicates if any operations draw text
	hasText bool
	// hasImage indicates if any operations draw images
	hasImage bool
	// modifiesTransparentBlack indicates if the list modifies transparent black pixels
	modifiesTransparentBlack bool
	// rtree is an optional spatial index for efficient culling (TODO: implement)
	// rtree *RTree
}

// Operation represents a single drawing operation that can be recorded and played back.
type Operation interface {
	// Invoke executes this operation on the given receiver
	Invoke(receiver OpReceiver)
	// Bounds returns the bounds affected by this operation, if available
	Bounds() *geom.Rect[Scalar]
	// Flags returns the attribute flags for this operation
	Flags() AttributeFlags
}

// NewDisplayList creates a new DisplayList with the given operations and bounds.
func NewDisplayList(operations []Operation, bounds geom.Rect[Scalar]) *DisplayList {
	dl := &DisplayList{
		operations:   make([]Operation, len(operations)),
		bounds:       bounds,
		totalOpCount: len(operations),
	}

	copy(dl.operations, operations)

	// Analyze operations to set flags
	for _, op := range operations {
		flags := op.Flags()
		if flags.HasAttribute(AttrFlagIsAntiAlias) {
			dl.hasAntiAliasing = true
		}
		if flags.HasAttribute(AttrFlagHasText) {
			dl.hasText = true
		}
		if flags.HasAttribute(AttrFlagHasImage) {
			dl.hasImage = true
		}
		if flags.HasAttribute(AttrFlagModifiesTransparentBlack) {
			dl.modifiesTransparentBlack = true
		}
	}

	return dl
}

// Bounds returns the bounding rectangle of all operations in the display list.
func (dl *DisplayList) Bounds() geom.Rect[Scalar] {
	return dl.bounds
}

// Operations returns the operations in the display list.
func (dl *DisplayList) Operations() []Operation {
	return dl.operations
}

// HasNonTrivialBlendMode returns true if any operations use non-trivial blend modes.
func (dl *DisplayList) HasNonTrivialBlendMode() bool {
	// TODO: Implement proper blend mode checking
	return false
}

// OpCount returns the total number of operations in the display list.
func (dl *DisplayList) OpCount() int {
	return dl.totalOpCount
}

// HasAntiAliasing returns true if any operations in the list use anti-aliasing.
func (dl *DisplayList) HasAntiAliasing() bool {
	return dl.hasAntiAliasing
}

// HasText returns true if any operations in the list draw text.
func (dl *DisplayList) HasText() bool {
	return dl.hasText
}

// HasImage returns true if any operations in the list draw images.
func (dl *DisplayList) HasImage() bool {
	return dl.hasImage
}

// ModifiesTransparentBlack returns true if the list modifies transparent black pixels.
func (dl *DisplayList) ModifiesTransparentBlack() bool {
	return dl.modifiesTransparentBlack
}

// IsEmpty returns true if the display list contains no operations.
func (dl *DisplayList) IsEmpty() bool {
	return dl.totalOpCount == 0
}

// Dispatch plays back all operations in the display list on the given receiver.
// This is the main method for rendering a display list.
func (dl *DisplayList) Dispatch(receiver OpReceiver) {
	for _, op := range dl.operations {
		op.Invoke(receiver)
	}
}

// DispatchWithCulling plays back operations in the display list on the given receiver,
// but only those that intersect with the given cull rectangle.
// This can improve performance by skipping operations outside the visible area.
func (dl *DisplayList) DispatchWithCulling(receiver OpReceiver, cullRect geom.Rect[Scalar]) {
	for _, op := range dl.operations {
		if bounds := op.Bounds(); bounds != nil {
			if !bounds.Intersects(cullRect) {
				continue
			}
		}
		op.Invoke(receiver)
	}
}

// Clone creates a deep copy of the display list.
func (dl *DisplayList) Clone() *DisplayList {
	clonedOps := make([]Operation, len(dl.operations))
	copy(clonedOps, dl.operations)

	return &DisplayList{
		operations:               clonedOps,
		bounds:                   dl.bounds,
		totalOpCount:             dl.totalOpCount,
		hasAntiAliasing:          dl.hasAntiAliasing,
		hasText:                  dl.hasText,
		hasImage:                 dl.hasImage,
		modifiesTransparentBlack: dl.modifiesTransparentBlack,
	}
}

// Equals returns true if this display list is equivalent to another.
func (dl *DisplayList) Equals(other *DisplayList) bool {
	if other == nil {
		return false
	}

	if dl.totalOpCount != other.totalOpCount {
		return false
	}

	if dl.bounds != other.bounds {
		return false
	}

	// Deep comparison would require Operation.Equals method
	// For now, we compare basic properties
	return dl.hasAntiAliasing == other.hasAntiAliasing &&
		dl.hasText == other.hasText &&
		dl.hasImage == other.hasImage &&
		dl.modifiesTransparentBlack == other.modifiesTransparentBlack
}
