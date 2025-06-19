// Package font - Rectangle packer functionality
package font

import (
	"github.com/opensraph/sraph/geom"
)

// skylineRectanglePacker implements a skyline-based rectangle packing algorithm.
// Based on Jukka Jylanki's work and adapted from Skia's implementation.
type skylineRectanglePacker struct {
	width     int
	height    int
	skyline   []skylineSegment
	areaSoFar int32
}

// skylineSegment represents a segment of the skyline.
type skylineSegment struct {
	x     int
	y     int
	width int
}

// NewSkylineRectanglePacker creates a new skyline rectangle packer.
func NewSkylineRectanglePacker(width, height int) RectanglePacker {
	return &skylineRectanglePacker{
		width:     width,
		height:    height,
		skyline:   []skylineSegment{{x: 0, y: 0, width: width}},
		areaSoFar: 0,
	}
}

// AddRect implements RectanglePacker.
func (p *skylineRectanglePacker) AddRect(width, height int) (geom.Point[geom.I32], bool) {
	if width <= 0 || height <= 0 {
		return geom.Point[geom.I32]{}, false
	}

	bestIndex := -1
	bestY := p.height + 1
	bestWastedHorizontalSpace := p.width + 1

	// Find the best position
	for i := 0; i < len(p.skyline); i++ {
		if p.rectangleFits(i, width, height) {
			y := p.computeSegmentY(i, width, height)
			if y < bestY || (y == bestY && p.wastedHorizontalSpace(i, width) < bestWastedHorizontalSpace) {
				bestIndex = i
				bestY = y
				bestWastedHorizontalSpace = p.wastedHorizontalSpace(i, width)
			}
		}
	}

	if bestIndex == -1 {
		return geom.Point[geom.I32]{}, false
	}

	position := geom.Point[geom.I32]{X: geom.I32(p.skyline[bestIndex].x), Y: geom.I32(bestY)}
	p.addSkylineLevel(bestIndex, p.skyline[bestIndex].x, bestY, width, height)
	p.areaSoFar += int32(width * height)

	return position, true
}

// GetPercentFull implements RectanglePacker.
func (p *skylineRectanglePacker) GetPercentFull() float32 {
	totalArea := p.width * p.height
	if totalArea == 0 {
		return 0
	}
	return float32(p.areaSoFar) / float32(totalArea)
}

// Reset implements RectanglePacker.
func (p *skylineRectanglePacker) Reset() {
	p.skyline = []skylineSegment{{x: 0, y: 0, width: p.width}}
	p.areaSoFar = 0
}

// rectangleFits checks if a rectangle can fit starting at the given skyline index.
func (p *skylineRectanglePacker) rectangleFits(skylineIndex, width, height int) bool {
	x := p.skyline[skylineIndex].x
	if x+width > p.width {
		return false
	}

	y := p.computeSegmentY(skylineIndex, width, height)
	return y+height <= p.height
}

// computeSegmentY computes the Y coordinate for a rectangle at the given skyline index.
func (p *skylineRectanglePacker) computeSegmentY(skylineIndex, width, height int) int {
	y := p.skyline[skylineIndex].y
	x := p.skyline[skylineIndex].x

	for i := skylineIndex; i < len(p.skyline) && p.skyline[i].x < x+width; i++ {
		if p.skyline[i].y > y {
			y = p.skyline[i].y
		}
	}

	return y
}

// wastedHorizontalSpace computes the wasted horizontal space for a rectangle.
func (p *skylineRectanglePacker) wastedHorizontalSpace(skylineIndex, width int) int {
	x := p.skyline[skylineIndex].x
	y := p.computeSegmentY(skylineIndex, width, 0)
	wastedArea := 0

	for i := skylineIndex; i < len(p.skyline) && p.skyline[i].x < x+width; i++ {
		if p.skyline[i].y < y {
			wastedArea += (y - p.skyline[i].y) * p.skyline[i].width
		}
	}

	return wastedArea
}

// addSkylineLevel adds a new level to the skyline after placing a rectangle.
func (p *skylineRectanglePacker) addSkylineLevel(skylineIndex, x, y, width, height int) {
	newNode := skylineSegment{x: x, y: y + height, width: width}

	// Insert the new node and remove overlapped nodes
	p.skyline = append(p.skyline[:skylineIndex], append([]skylineSegment{newNode}, p.skyline[skylineIndex:]...)...)

	// Remove nodes that are now covered
	i := skylineIndex + 1
	for i < len(p.skyline) && p.skyline[i].x < x+width {
		if p.skyline[i].x+p.skyline[i].width <= x+width {
			// This node is completely covered
			p.skyline = append(p.skyline[:i], p.skyline[i+1:]...)
		} else {
			// This node is partially covered
			p.skyline[i].width -= x + width - p.skyline[i].x
			p.skyline[i].x = x + width
			break
		}
	}

	// Merge nodes with the same height
	p.mergeSkylineNodes()
}

// mergeSkylineNodes merges adjacent skyline nodes with the same height.
func (p *skylineRectanglePacker) mergeSkylineNodes() {
	for i := 0; i < len(p.skyline)-1; i++ {
		if p.skyline[i].y == p.skyline[i+1].y {
			p.skyline[i].width += p.skyline[i+1].width
			p.skyline = append(p.skyline[:i+1], p.skyline[i+2:]...)
			i--
		}
	}
}
