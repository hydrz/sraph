package geom

import (
	"image"
	"image/color"
)

// ImagePathReceiver is a PathReceiver that renders paths to an image.RGBA.
// It implements both PathReceiver and image.Image interfaces.
type ImagePathReceiver[T Scalar] struct {
	img         *image.RGBA
	bounds      image.Rectangle
	strokeColor color.RGBA
	fillColor   color.RGBA
	currentPos  Point[T]
	pathStarted bool
	subpaths    [][]Point[T]
	currentPath []Point[T]
}

// NewImagePathReceiver creates a new ImagePathReceiver with the given dimensions.
func NewImagePathReceiver[T Scalar](width, height int) *ImagePathReceiver[T] {
	bounds := image.Rect(0, 0, width, height)
	img := image.NewRGBA(bounds)

	return &ImagePathReceiver[T]{
		img:         img,
		bounds:      bounds,
		strokeColor: color.RGBA{0, 0, 0, 255},       // Black stroke
		fillColor:   color.RGBA{128, 128, 128, 255}, // Gray fill
		subpaths:    make([][]Point[T], 0),
	}
}

// SetStrokeColor sets the stroke color for drawing lines.
func (r *ImagePathReceiver[T]) SetStrokeColor(c color.RGBA) {
	r.strokeColor = c
}

// SetFillColor sets the fill color for filling paths.
func (r *ImagePathReceiver[T]) SetFillColor(c color.RGBA) {
	r.fillColor = c
}

// Image returns the underlying image.RGBA.
func (r *ImagePathReceiver[T]) Image() *image.RGBA {
	return r.img
}

// Implement image.Image interface
func (r *ImagePathReceiver[T]) ColorModel() color.Model {
	return r.img.ColorModel()
}

func (r *ImagePathReceiver[T]) Bounds() image.Rectangle {
	return r.img.Bounds()
}

func (r *ImagePathReceiver[T]) At(x, y int) color.Color {
	return r.img.At(x, y)
}

// PathReceiver interface implementation

// MoveTo implements PathReceiver.
func (r *ImagePathReceiver[T]) MoveTo(p2 Point[T], willBeClosed bool) {
	r.currentPos = p2
	r.pathStarted = true

	// Start a new subpath
	if len(r.currentPath) > 0 {
		r.subpaths = append(r.subpaths, r.currentPath)
	}
	r.currentPath = []Point[T]{p2}
}

// LineTo implements PathReceiver.
func (r *ImagePathReceiver[T]) LineTo(p2 Point[T]) {
	if !r.pathStarted {
		r.MoveTo(p2, false)
		return
	}

	r.drawLine(r.currentPos, p2)
	r.currentPos = p2
	r.currentPath = append(r.currentPath, p2)
}

// QuadTo implements PathReceiver.
func (r *ImagePathReceiver[T]) QuadTo(cp, p2 Point[T]) {
	if !r.pathStarted {
		r.MoveTo(cp, false)
	}

	r.drawQuadraticBezier(r.currentPos, cp, p2)
	r.currentPos = p2
	r.currentPath = append(r.currentPath, p2)
}

// ConicTo implements PathReceiver.
func (r *ImagePathReceiver[T]) ConicTo(cp, p2 Point[T], weight float64) bool {
	// Convert conic to quadratic approximation
	// For simplicity, treat as quadratic curve
	r.QuadTo(cp, p2)
	return true
}

// CubicTo implements PathReceiver.
func (r *ImagePathReceiver[T]) CubicTo(cp1, cp2, p2 Point[T]) {
	if !r.pathStarted {
		r.MoveTo(cp1, false)
	}

	r.drawCubicBezier(r.currentPos, cp1, cp2, p2)
	r.currentPos = p2
	r.currentPath = append(r.currentPath, p2)
}

// Close implements PathReceiver.
func (r *ImagePathReceiver[T]) Close() {
	if len(r.currentPath) > 1 {
		// Draw line back to start
		startPoint := r.currentPath[0]
		r.drawLine(r.currentPos, startPoint)
		r.currentPos = startPoint
	}
}

// PathEnd implements PathReceiver.
func (r *ImagePathReceiver[T]) PathEnd() {
	// Add current path to subpaths if not empty
	if len(r.currentPath) > 0 {
		r.subpaths = append(r.subpaths, r.currentPath)
		r.currentPath = nil
	}
	r.pathStarted = false
}

// Clear clears the image with the given background color.
func (r *ImagePathReceiver[T]) Clear(bgColor color.RGBA) {
	for y := r.bounds.Min.Y; y < r.bounds.Max.Y; y++ {
		for x := r.bounds.Min.X; x < r.bounds.Max.X; x++ {
			r.img.SetRGBA(x, y, bgColor)
		}
	}
}

// Fill fills all closed paths with the current fill color.
func (r *ImagePathReceiver[T]) Fill() {
	for _, path := range r.subpaths {
		if len(path) >= 3 {
			r.fillPolygon(path)
		}
	}
}

// Drawing helper methods

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// drawLine draws a line between two points using Bresenham's algorithm.
func (r *ImagePathReceiver[T]) drawLine(p1, p2 Point[T]) {
	x0 := int(p1.X.Float64())
	y0 := int(p1.Y.Float64())
	x1 := int(p2.X.Float64())
	y1 := int(p2.Y.Float64())

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}

	err := dx - dy
	x, y := x0, y0

	for {
		r.setPixelSafe(x, y, r.strokeColor)

		if x == x1 && y == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// drawQuadraticBezier draws a quadratic Bézier curve.
func (r *ImagePathReceiver[T]) drawQuadraticBezier(p0, p1, p2 Point[T]) {
	steps := 50 // Number of line segments to approximate the curve
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		point := r.evaluateQuadraticBezier(p0, p1, p2, T(t))

		if i == 0 {
			continue
		}

		prevT := float64(i-1) / float64(steps)
		prevPoint := r.evaluateQuadraticBezier(p0, p1, p2, T(prevT))
		r.drawLine(prevPoint, point)
	}
}

// drawCubicBezier draws a cubic Bézier curve.
func (r *ImagePathReceiver[T]) drawCubicBezier(p0, p1, p2, p3 Point[T]) {
	steps := 50 // Number of line segments to approximate the curve
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		point := r.evaluateCubicBezier(p0, p1, p2, p3, T(t))

		if i == 0 {
			continue
		}

		prevT := float64(i-1) / float64(steps)
		prevPoint := r.evaluateCubicBezier(p0, p1, p2, p3, T(prevT))
		r.drawLine(prevPoint, point)
	}
}

// evaluateQuadraticBezier evaluates a quadratic Bézier curve at parameter t.
func (r *ImagePathReceiver[T]) evaluateQuadraticBezier(p0, p1, p2 Point[T], t T) Point[T] {
	oneMinusT := T(1) - t
	oneMinusTSq := oneMinusT * oneMinusT
	tSq := t * t

	x := oneMinusTSq*p0.X + T(2)*oneMinusT*t*p1.X + tSq*p2.X
	y := oneMinusTSq*p0.Y + T(2)*oneMinusT*t*p1.Y + tSq*p2.Y

	return Point[T]{x, y}
}

// evaluateCubicBezier evaluates a cubic Bézier curve at parameter t.
func (r *ImagePathReceiver[T]) evaluateCubicBezier(p0, p1, p2, p3 Point[T], t T) Point[T] {
	oneMinusT := T(1) - t
	oneMinusTCub := oneMinusT * oneMinusT * oneMinusT
	oneMinusTSq := oneMinusT * oneMinusT
	tSq := t * t
	tCub := t * t * t

	x := oneMinusTCub*p0.X + T(3)*oneMinusTSq*t*p1.X + T(3)*oneMinusT*tSq*p2.X + tCub*p3.X
	y := oneMinusTCub*p0.Y + T(3)*oneMinusTSq*t*p1.Y + T(3)*oneMinusT*tSq*p2.Y + tCub*p3.Y

	return Point[T]{x, y}
}

// fillPolygon fills a polygon using an improved scanline algorithm with even-odd rule.
func (r *ImagePathReceiver[T]) fillPolygon(vertices []Point[T]) {
	if len(vertices) < 3 {
		return
	}

	// Find bounding box
	minY := int(vertices[0].Y.Float64())
	maxY := minY
	minX := int(vertices[0].X.Float64())
	maxX := minX

	for _, v := range vertices {
		x := int(v.X.Float64())
		y := int(v.Y.Float64())
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
	}

	// Clamp to image bounds
	if minY < r.bounds.Min.Y {
		minY = r.bounds.Min.Y
	}
	if maxY >= r.bounds.Max.Y {
		maxY = r.bounds.Max.Y - 1
	}
	if minX < r.bounds.Min.X {
		minX = r.bounds.Min.X
	}
	if maxX >= r.bounds.Max.X {
		maxX = r.bounds.Max.X - 1
	}

	// Scanline fill using even-odd rule
	for y := minY; y <= maxY; y++ {
		intersections := r.findIntersections(vertices, float64(y)+0.5) // Use middle of scanline

		// Sort intersections
		for i := 0; i < len(intersections)-1; i++ {
			for j := i + 1; j < len(intersections); j++ {
				if intersections[i] > intersections[j] {
					intersections[i], intersections[j] = intersections[j], intersections[i]
				}
			}
		}

		// Fill between pairs of intersections (even-odd rule)
		for i := 0; i < len(intersections)-1; i += 2 {
			if i+1 < len(intersections) {
				x1 := intersections[i]
				x2 := intersections[i+1]

				// Ensure x1 <= x2
				if x1 > x2 {
					x1, x2 = x2, x1
				}

				// Clamp to bounds
				if x1 < minX {
					x1 = minX
				}
				if x2 > maxX {
					x2 = maxX
				}

				for x := x1; x <= x2; x++ {
					r.setPixelSafe(x, y, r.fillColor)
				}
			}
		}
	}
}

// findIntersections finds x-coordinates where the polygon edges intersect with a horizontal scanline.
func (r *ImagePathReceiver[T]) findIntersections(vertices []Point[T], y float64) []int {
	var intersections []int
	n := len(vertices)

	for i := 0; i < n; i++ {
		j := (i + 1) % n

		y1 := vertices[i].Y.Float64()
		y2 := vertices[j].Y.Float64()

		// Skip horizontal edges
		if abs(int(y1-y2)) < 1 {
			continue
		}

		// Check if scanline intersects this edge
		// Use strict inEqity to avoid double-counting vertices
		if (y1 < y && y <= y2) || (y2 < y && y <= y1) {
			x1 := vertices[i].X.Float64()
			x2 := vertices[j].X.Float64()

			// Calculate intersection x-coordinate using linear interpolation
			t := (y - y1) / (y2 - y1)
			x := x1 + t*(x2-x1)

			intersections = append(intersections, int(x+0.5)) // Round to nearest integer
		}
	}

	return intersections
}

// setPixelSafe sets a pixel with bounds checking.
func (r *ImagePathReceiver[T]) setPixelSafe(x, y int, c color.RGBA) {
	if x >= r.bounds.Min.X && x < r.bounds.Max.X &&
		y >= r.bounds.Min.Y && y < r.bounds.Max.Y {
		r.img.SetRGBA(x, y, c)
	}
}

// RasterizePathSource is a convenience function to rasterize a PathSource to an image.
func RasterizePathSource[T Scalar](source PathSource[T], width, height int) *image.RGBA {
	receiver := NewImagePathReceiver[T](width, height)

	// Clear with white background
	receiver.Clear(color.RGBA{255, 255, 255, 255})

	// Dispatch the path
	source.Dispatch(receiver)

	// Fill the path
	receiver.Fill()

	return receiver.Image()
}

// RasterizePathSourceWithColors is a convenience function with custom colors.
func RasterizePathSourceWithColors[T Scalar](
	source PathSource[T],
	width, height int,
	strokeColor, fillColor, bgColor color.RGBA,
) *image.RGBA {
	receiver := NewImagePathReceiver[T](width, height)

	receiver.SetStrokeColor(strokeColor)
	receiver.SetFillColor(fillColor)
	receiver.Clear(bgColor)

	source.Dispatch(receiver)
	receiver.Fill()

	return receiver.Image()
}
