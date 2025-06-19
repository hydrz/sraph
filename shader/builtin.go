// Package shader - Built-in shader sources
package shader

import _ "embed"

// Built-in shader sources embedded at compile time

//go:embed builtin/solid.wgsl
var solidShaderSource string

//go:embed builtin/texture.wgsl
var textureShaderSource string

//go:embed builtin/gradient.wgsl
var gradientShaderSource string

// BuiltinShaderNames contains the names of all built-in shaders.
var BuiltinShaderNames = []string{
	"solid",
	"texture",
	"gradient",
}

// GetBuiltinShaderSource returns the source code for a built-in shader.
func GetBuiltinShaderSource(name string) (string, bool) {
	switch name {
	case "solid":
		return solidShaderSource, true
	case "texture":
		return textureShaderSource, true
	case "gradient":
		return gradientShaderSource, true
	default:
		return "", false
	}
}

// BuiltinShaderInfo contains metadata about built-in shaders.
type BuiltinShaderInfo struct {
	Name        string
	Description string
	HasVertex   bool
	HasFragment bool
	Uniforms    []string
}

// GetBuiltinShaderInfo returns information about a built-in shader.
func GetBuiltinShaderInfo(name string) (BuiltinShaderInfo, bool) {
	switch name {
	case "solid":
		return BuiltinShaderInfo{
			Name:        "solid",
			Description: "Renders solid colors with transform support",
			HasVertex:   true,
			HasFragment: true,
			Uniforms:    []string{"transform", "color"},
		}, true
	case "texture":
		return BuiltinShaderInfo{
			Name:        "texture",
			Description: "Renders textured quads with color modulation",
			HasVertex:   true,
			HasFragment: true,
			Uniforms:    []string{"transform", "color_modulation", "texture", "texture_sampler"},
		}, true
	case "gradient":
		return BuiltinShaderInfo{
			Name:        "gradient",
			Description: "Renders linear and radial gradients",
			HasVertex:   true,
			HasFragment: true,
			Uniforms:    []string{"transform", "start_color", "end_color", "start_point", "end_point", "gradient_type"},
		}, true
	default:
		return BuiltinShaderInfo{}, false
	}
}
