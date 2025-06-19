package inputs

import (
	"testing"

	"github.com/opensraph/sraph/geom"
)

// TestTextureFilterInput tests basic texture filter input functionality
func TestTextureFilterInput_Basic(t *testing.T) {
	input := NewTextureFilterInput(nil)

	if input.IsValid() {
		t.Error("Expected invalid input with nil texture")
	}

	if input.GetInputType() != FilterInputTypeTexture {
		t.Errorf("Expected input type %v, got %v", FilterInputTypeTexture, input.GetInputType())
	}
}

// TestPlaceholderFilterInput tests basic placeholder filter input functionality
func TestPlaceholderFilterInput_Basic(t *testing.T) {
	size := geom.Size{Width: 100, Height: 100}
	color := geom.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0}

	input := NewPlaceholderFilterInput(size, color)

	if !input.IsValid() {
		t.Error("Expected valid placeholder input")
	}

	if input.GetInputType() != FilterInputTypePlaceholder {
		t.Errorf("Expected input type %v, got %v", FilterInputTypePlaceholder, input.GetInputType())
	}

	if input.GetPlaceholderSize() != size {
		t.Errorf("Expected size %v, got %v", size, input.GetPlaceholderSize())
	}

	if input.GetPlaceholderColor() != color {
		t.Errorf("Expected color %v, got %v", color, input.GetPlaceholderColor())
	}
}

// TestContentsFilterInput tests basic contents filter input functionality
func TestContentsFilterInput_Basic(t *testing.T) {
	input := NewContentsFilterInput(nil)

	if input.IsValid() {
		t.Error("Expected invalid input with nil contents")
	}

	if input.GetInputType() != FilterInputTypeContents {
		t.Errorf("Expected input type %v, got %v", FilterInputTypeContents, input.GetInputType())
	}
}

// TestFilterContentsFilterInput tests basic filter contents filter input functionality
func TestFilterContentsFilterInput_Basic(t *testing.T) {
	input := NewFilterContentsFilterInput(nil)

	if input.IsValid() {
		t.Error("Expected invalid input with nil filter contents")
	}

	if input.GetInputType() != FilterInputTypeFilter {
		t.Errorf("Expected input type %v, got %v", FilterInputTypeFilter, input.GetInputType())
	}
}

// TestFilterInputCloning tests cloning functionality for all input types
func TestFilterInputCloning(t *testing.T) {
	// Test texture input cloning
	textureInput := NewTextureFilterInput(nil)
	clonedTexture := textureInput.Clone()

	if clonedTexture.GetInputType() != textureInput.GetInputType() {
		t.Error("Cloned texture input type mismatch")
	}

	// Test placeholder input cloning
	size := geom.Size{Width: 50, Height: 50}
	color := geom.Color{R: 0.5, G: 0.5, B: 0.5, A: 1.0}
	placeholderInput := NewPlaceholderFilterInput(size, color)
	clonedPlaceholder := placeholderInput.Clone()

	if clonedPlaceholder.GetInputType() != placeholderInput.GetInputType() {
		t.Error("Cloned placeholder input type mismatch")
	}

	// Test contents input cloning
	contentsInput := NewContentsFilterInput(nil)
	clonedContents := contentsInput.Clone()

	if clonedContents.GetInputType() != contentsInput.GetInputType() {
		t.Error("Cloned contents input type mismatch")
	}

	// Test filter contents input cloning
	filterInput := NewFilterContentsFilterInput(nil)
	clonedFilter := filterInput.Clone()

	if clonedFilter.GetInputType() != filterInput.GetInputType() {
		t.Error("Cloned filter input type mismatch")
	}
}

// TestFilterInputBounds tests bounds functionality
func TestFilterInputBounds(t *testing.T) {
	size := geom.Size{Width: 200, Height: 150}
	color := geom.Color{R: 0.0, G: 1.0, B: 0.0, A: 1.0}

	input := NewPlaceholderFilterInput(size, color)
	bounds := input.GetBounds()

	expectedBounds := geom.Rect{
		Origin: geom.Point{X: 0, Y: 0},
		Size:   size,
	}

	if bounds != expectedBounds {
		t.Errorf("Expected bounds %v, got %v", expectedBounds, bounds)
	}

	// Test bounds modification
	newBounds := geom.Rect{
		Origin: geom.Point{X: 10, Y: 10},
		Size:   geom.Size{Width: 300, Height: 200},
	}

	input.SetBounds(newBounds)
	updatedBounds := input.GetBounds()

	if updatedBounds != newBounds {
		t.Errorf("Expected updated bounds %v, got %v", newBounds, updatedBounds)
	}
}
