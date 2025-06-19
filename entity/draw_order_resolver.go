// Package entity - Draw order resolver functionality
package entity

// DrawOrderResolver helps record draw indices in painter's order and sorts
// the draws into an optimized order based on translucency and clips.
// This is inspired by Impeller's DrawOrderResolver.
type DrawOrderResolver interface {
	// AddElement adds an element to the draw order with opacity information.
	AddElement(elementIndex uint, isOpaque bool)

	// PushClip pushes a clip element onto the stack.
	PushClip(elementIndex uint)

	// PopClip pops a clip element from the stack.
	PopClip()

	// Flush flushes the current layer and starts a new one.
	Flush()

	// GetSortedDraws returns the sorted draw order with skip counts.
	GetSortedDraws(opaqueSkipCount, translucentSkipCount uint) []uint
}

// DrawOrderLayer represents a layer in the draw order hierarchy.
type DrawOrderLayer struct {
	// opaqueElements contains opaque elements in this layer.
	opaqueElements []uint
	// translucentElements contains translucent elements in this layer.
	translucentElements []uint
	// dependentElements contains elements that depend on this layer.
	dependentElements []uint
}

// defaultDrawOrderResolver provides the default implementation of DrawOrderResolver.
type defaultDrawOrderResolver struct {
	// drawOrderLayers contains the hierarchy of draw order layers.
	drawOrderLayers []DrawOrderLayer
	// firstRootFlush stores the first root flush if any.
	firstRootFlush *DrawOrderLayer
	// sortedElements stores the sorted element indices.
	sortedElements []uint
}

// NewDrawOrderResolver creates a new draw order resolver.
func NewDrawOrderResolver() DrawOrderResolver {
	return &defaultDrawOrderResolver{
		drawOrderLayers: []DrawOrderLayer{{}}, // Start with one empty layer
		sortedElements:  make([]uint, 0),
	}
}

// AddElement implements DrawOrderResolver.
func (r *defaultDrawOrderResolver) AddElement(elementIndex uint, isOpaque bool) {
	if len(r.drawOrderLayers) == 0 {
		r.drawOrderLayers = append(r.drawOrderLayers, DrawOrderLayer{})
	}

	layer := &r.drawOrderLayers[len(r.drawOrderLayers)-1]
	if isOpaque {
		layer.opaqueElements = append(layer.opaqueElements, elementIndex)
	} else {
		layer.translucentElements = append(layer.translucentElements, elementIndex)
	}
}

// PushClip implements DrawOrderResolver.
func (r *defaultDrawOrderResolver) PushClip(elementIndex uint) {
	if len(r.drawOrderLayers) == 0 {
		r.drawOrderLayers = append(r.drawOrderLayers, DrawOrderLayer{})
	}

	// Add the clip element as a dependent element
	layer := &r.drawOrderLayers[len(r.drawOrderLayers)-1]
	layer.dependentElements = append(layer.dependentElements, elementIndex)

	// Push a new layer
	r.drawOrderLayers = append(r.drawOrderLayers, DrawOrderLayer{})
}

// PopClip implements DrawOrderResolver.
func (r *defaultDrawOrderResolver) PopClip() {
	if len(r.drawOrderLayers) <= 1 {
		return // Can't pop the root layer
	}
	r.drawOrderLayers = r.drawOrderLayers[:len(r.drawOrderLayers)-1]
}

// Flush implements DrawOrderResolver.
func (r *defaultDrawOrderResolver) Flush() {
	if len(r.drawOrderLayers) > 0 && r.firstRootFlush == nil {
		// Store the first root flush
		rootLayer := r.drawOrderLayers[0]
		r.firstRootFlush = &rootLayer
	}

	// Reset to a single empty layer
	r.drawOrderLayers = []DrawOrderLayer{{}}
}

// GetSortedDraws implements DrawOrderResolver.
func (r *defaultDrawOrderResolver) GetSortedDraws(opaqueSkipCount, translucentSkipCount uint) []uint {
	r.sortedElements = r.sortedElements[:0] // Clear but keep capacity

	// Process all layers
	for i := range r.drawOrderLayers {
		r.writeCombinedDraws(&r.drawOrderLayers[i], opaqueSkipCount, translucentSkipCount)
	}

	// Process first root flush if available
	if r.firstRootFlush != nil {
		r.writeCombinedDraws(r.firstRootFlush, opaqueSkipCount, translucentSkipCount)
	}

	return r.sortedElements
}

// writeCombinedDraws writes the combined draws for a layer.
func (r *defaultDrawOrderResolver) writeCombinedDraws(layer *DrawOrderLayer, opaqueSkipCount, translucentSkipCount uint) {
	// Add dependent elements first
	for _, element := range layer.dependentElements {
		r.sortedElements = append(r.sortedElements, element)
	}

	// Add opaque elements in reverse order (back to front)
	opaqueCount := uint(len(layer.opaqueElements))
	if opaqueSkipCount < opaqueCount {
		for i := int(opaqueCount - opaqueSkipCount - 1); i >= 0; i-- {
			r.sortedElements = append(r.sortedElements, layer.opaqueElements[i])
		}
	}

	// Add translucent elements in forward order (front to back)
	translucentCount := uint(len(layer.translucentElements))
	if translucentSkipCount < translucentCount {
		start := translucentSkipCount
		for i := start; i < translucentCount; i++ {
			r.sortedElements = append(r.sortedElements, layer.translucentElements[i])
		}
	}
}

// WriteCombinedDraws writes combined draws to the destination slice.
func (layer *DrawOrderLayer) WriteCombinedDraws(destination *[]uint, opaqueSkipCount, translucentSkipCount uint) {
	// Add dependent elements first
	for _, element := range layer.dependentElements {
		*destination = append(*destination, element)
	}

	// Add opaque elements in reverse order (back to front)
	opaqueCount := uint(len(layer.opaqueElements))
	if opaqueSkipCount < opaqueCount {
		for i := int(opaqueCount - opaqueSkipCount - 1); i >= 0; i-- {
			*destination = append(*destination, layer.opaqueElements[i])
		}
	}

	// Add translucent elements in forward order (front to back)
	translucentCount := uint(len(layer.translucentElements))
	if translucentSkipCount < translucentCount {
		start := translucentSkipCount
		for i := start; i < translucentCount; i++ {
			*destination = append(*destination, layer.translucentElements[i])
		}
	}
}

// GetOpaqueCount returns the number of opaque elements in this layer.
func (layer *DrawOrderLayer) GetOpaqueCount() uint {
	return uint(len(layer.opaqueElements))
}

// GetTranslucentCount returns the number of translucent elements in this layer.
func (layer *DrawOrderLayer) GetTranslucentCount() uint {
	return uint(len(layer.translucentElements))
}

// GetDependentCount returns the number of dependent elements in this layer.
func (layer *DrawOrderLayer) GetDependentCount() uint {
	return uint(len(layer.dependentElements))
}

// IsEmpty returns whether this layer is empty.
func (layer *DrawOrderLayer) IsEmpty() bool {
	return len(layer.opaqueElements) == 0 &&
		len(layer.translucentElements) == 0 &&
		len(layer.dependentElements) == 0
}

// Clear clears all elements from this layer.
func (layer *DrawOrderLayer) Clear() {
	layer.opaqueElements = layer.opaqueElements[:0]
	layer.translucentElements = layer.translucentElements[:0]
	layer.dependentElements = layer.dependentElements[:0]
}
