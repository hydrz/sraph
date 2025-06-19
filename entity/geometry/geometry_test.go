package geometry

import (
	"testing"

	"github.com/opensraph/sraph/geom"
)

// TestRectGeometry tests basic rectangle geometry functionality
func TestRectGeometry(t *testing.T) {
	rect := geom.Rect{
		Origin: geom.Point{X: 10, Y: 20},
		Size:   geom.Size{Width: 100, Height: 50},
	}

	geom := NewRectGeometry(rect)
	if geom == nil {
		t.Fatal("NewRectGeometry returned nil")
	}

	bounds := geom.GetBounds()
	if bounds != rect {
		t.Errorf("Expected bounds %v, got %v", rect, bounds)
	}

	// TODO: Add more comprehensive tests for position buffer, UVs, etc.
}

// TestCircleGeometry tests basic circle geometry functionality
func TestCircleGeometry(t *testing.T) {
	center := geom.Point{X: 50, Y: 50}
	radius := float32(25)

	geom := NewCircleGeometry(center, radius)
	if geom == nil {
		t.Fatal("NewCircleGeometry returned nil")
	}

	bounds := geom.GetBounds()
	expectedBounds := geom.Rect{
		Origin: geom.Point{X: center.X - radius, Y: center.Y - radius},
		Size:   geom.Size{Width: 2 * radius, Height: 2 * radius},
	}

	if bounds != expectedBounds {
		t.Errorf("Expected bounds %v, got %v", expectedBounds, bounds)
	}

	// TODO: Add more comprehensive tests
}

// TestLineGeometry tests basic line geometry functionality
func TestLineGeometry(t *testing.T) {
	start := geom.Point{X: 0, Y: 0}
	end := geom.Point{X: 100, Y: 100}
	width := float32(5)

	geom := NewLineGeometry(start, end, width)
	if geom == nil {
		t.Fatal("NewLineGeometry returned nil")
	}

	bounds := geom.GetBounds()

	// Line bounds should encompass both points plus stroke width
	expectedMinX := min(start.X, end.X) - width/2
	expectedMaxX := max(start.X, end.X) + width/2
	expectedMinY := min(start.Y, end.Y) - width/2
	expectedMaxY := max(start.Y, end.Y) + width/2

	expectedBounds := geom.Rect{
		Origin: geom.Point{X: expectedMinX, Y: expectedMinY},
		Size: geom.Size{
			Width:  expectedMaxX - expectedMinX,
			Height: expectedMaxY - expectedMinY,
		},
	}

	if bounds != expectedBounds {
		t.Errorf("Expected bounds %v, got %v", expectedBounds, bounds)
	}

	// TODO: Add more comprehensive tests
}

// TestRoundRectGeometry tests basic rounded rectangle geometry functionality
func TestRoundRectGeometry(t *testing.T) {
	rect := geom.Rect{
		Origin: geom.Point{X: 10, Y: 20},
		Size:   geom.Size{Width: 100, Height: 50},
	}
	radius := float32(10)

	geom := NewRoundRectGeometry(rect, radius)
	if geom == nil {
		t.Fatal("NewRoundRectGeometry returned nil")
	}

	bounds := geom.GetBounds()
	if bounds != rect {
		t.Errorf("Expected bounds %v, got %v", rect, bounds)
	}

	// TODO: Add tests for corner radius validation, tessellation, etc.
}

// TestEllipseGeometry tests basic ellipse geometry functionality
func TestEllipseGeometry(t *testing.T) {
	rect := geom.Rect{
		Origin: geom.Point{X: 25, Y: 25},
		Size:   geom.Size{Width: 50, Height: 30},
	}

	geom := NewEllipseGeometry(rect)
	if geom == nil {
		t.Fatal("NewEllipseGeometry returned nil")
	}

	bounds := geom.GetBounds()
	if bounds != rect {
		t.Errorf("Expected bounds %v, got %v", rect, bounds)
	}

	// TODO: Add more comprehensive tests
}

// TestSuperellipseGeometry tests basic superellipse geometry functionality
func TestSuperellipseGeometry(t *testing.T) {
	rect := geom.Rect{
		Origin: geom.Point{X: 10, Y: 10},
		Size:   geom.Size{Width: 80, Height: 60},
	}
	exponent := float32(2.5)

	geom := NewSuperellipseGeometry(rect, exponent)
	if geom == nil {
		t.Fatal("NewSuperellipseGeometry returned nil")
	}

	bounds := geom.GetBounds()
	if bounds != rect {
		t.Errorf("Expected bounds %v, got %v", rect, bounds)
	}

	if geom.GetExponent() != exponent {
		t.Errorf("Expected exponent %v, got %v", exponent, geom.GetExponent())
	}

	// TODO: Add more comprehensive tests
}

// TestRoundSuperellipseGeometry tests basic rounded superellipse geometry functionality
func TestRoundSuperellipseGeometry(t *testing.T) {
	rect := geom.Rect{
		Origin: geom.Point{X: 5, Y: 5},
		Size:   geom.Size{Width: 90, Height: 70},
	}
	exponent := float32(3.0)
	cornerRadius := float32(8.0)

	geom := NewRoundSuperellipseGeometry(rect, exponent, cornerRadius)
	if geom == nil {
		t.Fatal("NewRoundSuperellipseGeometry returned nil")
	}

	bounds := geom.GetBounds()
	if bounds != rect {
		t.Errorf("Expected bounds %v, got %v", rect, bounds)
	}

	if geom.GetExponent() != exponent {
		t.Errorf("Expected exponent %v, got %v", exponent, geom.GetExponent())
	}

	if geom.GetCornerRadius() != cornerRadius {
		t.Errorf("Expected corner radius %v, got %v", cornerRadius, geom.GetCornerRadius())
	}

	// TODO: Add more comprehensive tests
}

// Helper functions for min/max since we don't want to use built-in functions
func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
