package dl

import (
	"fmt"

	"github.com/opensraph/sraph/geom"
)

// RTree implements an R-tree spatial index for efficient spatial queries.
// It stores bounding rectangles with associated IDs for fast intersection queries.
type RTree struct {
	root      *rtreeNode
	size      int
	maxNodes  int
	minNodes  int
	invalidID int
}

// rtreeNode represents a node in the R-tree.
type rtreeNode struct {
	bounds   geom.Rect[Scalar]
	children []*rtreeNode
	entries  []rtreeEntry
	isLeaf   bool
	level    int
}

// rtreeEntry represents an entry in an R-tree leaf node.
type rtreeEntry struct {
	bounds geom.Rect[Scalar]
	id     int
}

// NewRTree creates a new R-tree with the specified parameters.
func NewRTree(maxNodes, minNodes, invalidID int) *RTree {
	if maxNodes < 2 {
		maxNodes = 16
	}
	if minNodes < 1 {
		minNodes = maxNodes / 2
	}

	return &RTree{
		maxNodes:  maxNodes,
		minNodes:  minNodes,
		invalidID: invalidID,
	}
}

// NewRTreeFromRects creates an R-tree from a slice of rectangles and IDs.
func NewRTreeFromRects(rects []geom.Rect[Scalar], ids []int, invalidID int) *RTree {
	tree := NewRTree(16, 8, invalidID)

	for i, rect := range rects {
		id := invalidID
		if i < len(ids) {
			id = ids[i]
		}
		tree.Insert(rect, id)
	}

	return tree
}

// Insert adds a rectangle with the given ID to the R-tree.
func (rt *RTree) Insert(bounds geom.Rect[Scalar], id int) {
	entry := rtreeEntry{bounds: bounds, id: id}

	if rt.root == nil {
		rt.root = &rtreeNode{
			bounds:  bounds,
			entries: []rtreeEntry{entry},
			isLeaf:  true,
			level:   0,
		}
		rt.size = 1
		return
	}

	rt.insert(rt.root, entry)
	rt.size++
}

// Search finds all rectangles that intersect with the query rectangle.
func (rt *RTree) Search(query geom.Rect[Scalar]) []int {
	if rt.root == nil {
		return nil
	}

	var results []int
	rt.search(rt.root, query, &results)
	return results
}

// SearchAndConsolidateRects searches for intersecting rectangles and returns
// a consolidated list of non-overlapping rectangles.
func (rt *RTree) SearchAndConsolidateRects(query geom.Rect[Scalar], deband bool) []geom.Rect[Scalar] {
	ids := rt.Search(query)
	if len(ids) == 0 {
		return nil
	}

	// TODO: Implement rectangle consolidation algorithm
	// This would merge overlapping rectangles into non-overlapping ones
	var rects []geom.Rect[Scalar]

	// For now, just return the query rectangle
	rects = append(rects, query)

	return rects
}

// Size returns the number of entries in the R-tree.
func (rt *RTree) Size() int {
	return rt.size
}

// Bounds returns the overall bounds of all entries in the R-tree.
func (rt *RTree) Bounds() geom.Rect[Scalar] {
	if rt.root == nil {
		return geom.Rect[Scalar]{}
	}
	return rt.root.bounds
}

// Clear removes all entries from the R-tree.
func (rt *RTree) Clear() {
	rt.root = nil
	rt.size = 0
}

// String returns a string representation of the R-tree.
func (rt *RTree) String() string {
	return fmt.Sprintf("RTree{size: %d, bounds: %v}", rt.size, rt.Bounds())
}

// insert inserts an entry into the R-tree at the appropriate node.
func (rt *RTree) insert(node *rtreeNode, entry rtreeEntry) {
	if node.isLeaf {
		node.entries = append(node.entries, entry)
		node.bounds = expandBounds(node.bounds, entry.bounds)

		// Check if node needs to be split
		if len(node.entries) > rt.maxNodes {
			rt.splitLeafNode(node)
		}
	} else {
		// Find best child to insert into
		bestChild := rt.chooseSubtree(node, entry.bounds)
		rt.insert(bestChild, entry)

		// Update bounds
		node.bounds = expandBounds(node.bounds, entry.bounds)

		// Check if child was split and handle overflow
		if len(node.children) > rt.maxNodes {
			rt.splitInternalNode(node)
		}
	}
}

// search recursively searches for intersecting rectangles.
func (rt *RTree) search(node *rtreeNode, query geom.Rect[Scalar], results *[]int) {
	if !intersects(node.bounds, query) {
		return
	}

	if node.isLeaf {
		for _, entry := range node.entries {
			if intersects(entry.bounds, query) && entry.id != rt.invalidID {
				*results = append(*results, entry.id)
			}
		}
	} else {
		for _, child := range node.children {
			rt.search(child, query, results)
		}
	}
}

// chooseSubtree selects the best child node for inserting a rectangle.
func (rt *RTree) chooseSubtree(node *rtreeNode, bounds geom.Rect[Scalar]) *rtreeNode {
	if len(node.children) == 0 {
		return nil
	}

	best := node.children[0]
	minEnlargement := enlargementArea(best.bounds, bounds)

	for _, child := range node.children[1:] {
		enlargement := enlargementArea(child.bounds, bounds)
		if enlargement < minEnlargement {
			minEnlargement = enlargement
			best = child
		}
	}

	return best
}

// splitLeafNode splits a leaf node that has too many entries.
func (rt *RTree) splitLeafNode(node *rtreeNode) {
	// TODO: Implement proper R-tree node splitting algorithm
	// This would use techniques like quadratic split or linear split

	// Simple split for now - divide entries in half
	mid := len(node.entries) / 2

	// Create new node with second half of entries
	newNode := &rtreeNode{
		entries: make([]rtreeEntry, len(node.entries)-mid),
		isLeaf:  true,
		level:   node.level,
	}
	copy(newNode.entries, node.entries[mid:])

	// Update original node with first half
	node.entries = node.entries[:mid]

	// Recalculate bounds
	node.bounds = calculateBounds(node)
	newNode.bounds = calculateBounds(newNode)

	// If this is the root, create a new root
	if node == rt.root {
		rt.root = &rtreeNode{
			children: []*rtreeNode{node, newNode},
			isLeaf:   false,
			level:    node.level + 1,
		}
		rt.root.bounds = expandBounds(node.bounds, newNode.bounds)
	}
}

// splitInternalNode splits an internal node that has too many children.
func (rt *RTree) splitInternalNode(node *rtreeNode) {
	// TODO: Implement internal node splitting
	// Similar to leaf splitting but for child nodes
}

// Helper functions

// intersects checks if two rectangles intersect.
func intersects(a, b geom.Rect[Scalar]) bool {
	return a.Left <= b.Right && b.Left <= a.Right &&
		a.Top <= b.Bottom && b.Top <= a.Bottom
}

// expandBounds expands bounds to include another rectangle.
func expandBounds(bounds, other geom.Rect[Scalar]) geom.Rect[Scalar] {
	return geom.NewRect(
		min(bounds.Left, other.Left),
		min(bounds.Top, other.Top),
		max(bounds.Right, other.Right),
		max(bounds.Bottom, other.Bottom),
	)
}

// enlargementArea calculates the area increase needed to include a rectangle.
func enlargementArea(bounds, other geom.Rect[Scalar]) Scalar {
	expanded := expandBounds(bounds, other)
	originalArea := (bounds.Right - bounds.Left) * (bounds.Bottom - bounds.Top)
	expandedArea := (expanded.Right - expanded.Left) * (expanded.Bottom - expanded.Top)
	return expandedArea - originalArea
}

// calculateBounds calculates the bounding rectangle for a node.
func calculateBounds(node *rtreeNode) geom.Rect[Scalar] {
	if node.isLeaf && len(node.entries) > 0 {
		bounds := node.entries[0].bounds
		for _, entry := range node.entries[1:] {
			bounds = expandBounds(bounds, entry.bounds)
		}
		return bounds
	} else if !node.isLeaf && len(node.children) > 0 {
		bounds := node.children[0].bounds
		for _, child := range node.children[1:] {
			bounds = expandBounds(bounds, child.bounds)
		}
		return bounds
	}
	return geom.Rect[Scalar]{}
}

// Region represents a collection of non-overlapping rectangles.
// It's used for efficient region operations and clipping.
type Region struct {
	rects  []geom.Rect[Scalar]
	bounds geom.Rect[Scalar]
}

// NewRegion creates a new empty region.
func NewRegion() *Region {
	return &Region{}
}

// NewRegionFromRect creates a new region from a single rectangle.
func NewRegionFromRect(rect geom.Rect[Scalar]) *Region {
	return &Region{
		rects:  []geom.Rect[Scalar]{rect},
		bounds: rect,
	}
}

// AddRect adds a rectangle to the region.
func (r *Region) AddRect(rect geom.Rect[Scalar]) {
	// TODO: Implement proper region union algorithm
	// This would merge overlapping rectangles efficiently
	r.rects = append(r.rects, rect)
	if len(r.rects) == 1 {
		r.bounds = rect
	} else {
		r.bounds = expandBounds(r.bounds, rect)
	}
}

// Intersect computes the intersection of this region with another region.
func (r *Region) Intersect(other *Region) *Region {
	// TODO: Implement region intersection
	result := NewRegion()

	for _, rect1 := range r.rects {
		for _, rect2 := range other.rects {
			if intersects(rect1, rect2) {
				// Calculate intersection rectangle
				intersection := geom.NewRect(
					max(rect1.Left, rect2.Left),
					max(rect1.Top, rect2.Top),
					min(rect1.Right, rect2.Right),
					min(rect1.Bottom, rect2.Bottom),
				)
				result.AddRect(intersection)
			}
		}
	}

	return result
}

// Union computes the union of this region with another region.
func (r *Region) Union(other *Region) *Region {
	// TODO: Implement efficient region union
	result := NewRegion()

	// Simple implementation - add all rectangles
	for _, rect := range r.rects {
		result.AddRect(rect)
	}
	for _, rect := range other.rects {
		result.AddRect(rect)
	}

	return result
}

// Contains checks if a point is contained within the region.
func (r *Region) Contains(point geom.Point[Scalar]) bool {
	for _, rect := range r.rects {
		if point.X >= rect.Left && point.X <= rect.Right &&
			point.Y >= rect.Top && point.Y <= rect.Bottom {
			return true
		}
	}
	return false
}

// Bounds returns the overall bounds of the region.
func (r *Region) Bounds() geom.Rect[Scalar] {
	return r.bounds
}

// Rects returns the rectangles that make up this region.
func (r *Region) Rects() []geom.Rect[Scalar] {
	return r.rects
}

// IsEmpty returns true if the region contains no rectangles.
func (r *Region) IsEmpty() bool {
	return len(r.rects) == 0
}

// Area returns the total area covered by the region.
func (r *Region) Area() Scalar {
	var total Scalar
	for _, rect := range r.rects {
		total += (rect.Right - rect.Left) * (rect.Bottom - rect.Top)
	}
	return total
}

// String returns a string representation of the region.
func (r *Region) String() string {
	return fmt.Sprintf("Region{rects: %d, bounds: %v}", len(r.rects), r.bounds)
}

// SpatialIndex provides a generic interface for spatial indexing structures.
type SpatialIndex interface {
	// Insert adds a rectangle with the given ID to the index.
	Insert(bounds geom.Rect[Scalar], id int)
	// Search finds all IDs whose rectangles intersect with the query.
	Search(query geom.Rect[Scalar]) []int
	// Size returns the number of entries in the index.
	Size() int
	// Clear removes all entries from the index.
	Clear()
}

// QuadTree implements a quadtree spatial index as an alternative to R-tree.
type QuadTree struct {
	root     *quadNode
	bounds   geom.Rect[Scalar]
	maxDepth int
	maxItems int
	size     int
}

// quadNode represents a node in the quadtree.
type quadNode struct {
	bounds   geom.Rect[Scalar]
	children [4]*quadNode // NW, NE, SW, SE
	entries  []rtreeEntry
	level    int
}

// NewQuadTree creates a new quadtree with the specified bounds and limits.
func NewQuadTree(bounds geom.Rect[Scalar], maxDepth, maxItems int) *QuadTree {
	return &QuadTree{
		bounds:   bounds,
		maxDepth: maxDepth,
		maxItems: maxItems,
	}
}

// Insert implements SpatialIndex interface.
func (qt *QuadTree) Insert(bounds geom.Rect[Scalar], id int) {
	entry := rtreeEntry{bounds: bounds, id: id}

	if qt.root == nil {
		qt.root = &quadNode{
			bounds:  qt.bounds,
			entries: []rtreeEntry{entry},
			level:   0,
		}
	} else {
		qt.insertEntry(qt.root, entry)
	}

	qt.size++
}

// Search implements SpatialIndex interface.
func (qt *QuadTree) Search(query geom.Rect[Scalar]) []int {
	if qt.root == nil {
		return nil
	}

	var results []int
	qt.searchNode(qt.root, query, &results)
	return results
}

// Size implements SpatialIndex interface.
func (qt *QuadTree) Size() int {
	return qt.size
}

// Clear implements SpatialIndex interface.
func (qt *QuadTree) Clear() {
	qt.root = nil
	qt.size = 0
}

// insertEntry inserts an entry into the quadtree.
func (qt *QuadTree) insertEntry(node *quadNode, entry rtreeEntry) {
	// If node has children, try to insert into appropriate child
	if node.children[0] != nil {
		for i := 0; i < 4; i++ {
			if intersects(node.children[i].bounds, entry.bounds) {
				qt.insertEntry(node.children[i], entry)
				return
			}
		}
	}

	// Add to this node
	node.entries = append(node.entries, entry)

	// Check if we need to subdivide
	if len(node.entries) > qt.maxItems && node.level < qt.maxDepth {
		qt.subdivide(node)
	}
}

// subdivide divides a quadtree node into four children.
func (qt *QuadTree) subdivide(node *quadNode) {
	centerX := (node.bounds.Left + node.bounds.Right) / 2
	centerY := (node.bounds.Top + node.bounds.Bottom) / 2

	// Create four child nodes
	node.children[0] = &quadNode{ // NW
		bounds: geom.NewRect(node.bounds.Left, node.bounds.Top, centerX, centerY),
		level:  node.level + 1,
	}
	node.children[1] = &quadNode{ // NE
		bounds: geom.NewRect(centerX, node.bounds.Top, node.bounds.Right, centerY),
		level:  node.level + 1,
	}
	node.children[2] = &quadNode{ // SW
		bounds: geom.NewRect(node.bounds.Left, centerY, centerX, node.bounds.Bottom),
		level:  node.level + 1,
	}
	node.children[3] = &quadNode{ // SE
		bounds: geom.NewRect(centerX, centerY, node.bounds.Right, node.bounds.Bottom),
		level:  node.level + 1,
	}

	// Redistribute entries to children
	entries := node.entries
	node.entries = nil

	for _, entry := range entries {
		qt.insertEntry(node, entry)
	}
}

// searchNode searches for intersecting entries in a quadtree node.
func (qt *QuadTree) searchNode(node *quadNode, query geom.Rect[Scalar], results *[]int) {
	if !intersects(node.bounds, query) {
		return
	}

	// Check entries in this node
	for _, entry := range node.entries {
		if intersects(entry.bounds, query) {
			*results = append(*results, entry.id)
		}
	}

	// Check children
	for i := 0; i < 4; i++ {
		if node.children[i] != nil {
			qt.searchNode(node.children[i], query, results)
		}
	}
}
