// Package shader - Utility functions and helpers
package shader

import (
	"fmt"
	"strings"
)

// ShaderUtils provides utility functions for shader operations.
type ShaderUtils struct{}

// NewShaderUtils creates a new shader utilities instance.
func NewShaderUtils() *ShaderUtils {
	return &ShaderUtils{}
}

// ParseShaderType attempts to determine shader type from source content.
func (su *ShaderUtils) ParseShaderType(source string) (ShaderType, error) {
	source = strings.ToLower(source)

	// Check for WGSL entry point attributes
	if strings.Contains(source, "@vertex") {
		return ShaderTypeVertex, nil
	}
	if strings.Contains(source, "@fragment") {
		return ShaderTypeFragment, nil
	}
	if strings.Contains(source, "@compute") {
		return ShaderTypeCompute, nil
	}

	// Check for GLSL version and shader type
	if strings.Contains(source, "#version") {
		if strings.Contains(source, "gl_position") {
			return ShaderTypeVertex, nil
		}
		if strings.Contains(source, "gl_fragcolor") || strings.Contains(source, "fragcolor") {
			return ShaderTypeFragment, nil
		}
	}

	return ShaderTypeVertex, fmt.Errorf("unable to determine shader type from source")
}

// ExtractUniforms extracts uniform declarations from WGSL shader source.
func (su *ShaderUtils) ExtractUniforms(source string) []UniformInfo {
	var uniforms []UniformInfo
	lines := strings.Split(source, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Look for uniform declarations
		if strings.HasPrefix(trimmed, "@group(") && strings.Contains(trimmed, "@binding(") {
			// Parse WebGPU uniform declaration
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if uniform := su.parseWGSLUniform(trimmed, nextLine); uniform.Name != "" {
					uniforms = append(uniforms, uniform)
				}
			}
		}
	}

	return uniforms
}

// parseWGSLUniform parses a WGSL uniform declaration.
func (su *ShaderUtils) parseWGSLUniform(groupBinding, declaration string) UniformInfo {
	uniform := UniformInfo{}

	// Extract binding information
	if strings.Contains(groupBinding, "@binding(") {
		parts := strings.Split(groupBinding, "@binding(")
		if len(parts) > 1 {
			bindingPart := strings.Split(parts[1], ")")[0]
			fmt.Sscanf(bindingPart, "%d", &uniform.Binding)
		}
	}

	// Parse declaration
	if strings.HasPrefix(declaration, "var<uniform>") {
		// var<uniform> uniformName: UniformType;
		parts := strings.Fields(declaration)
		if len(parts) >= 3 {
			nameType := parts[1]
			if strings.Contains(nameType, ":") {
				nameParts := strings.Split(nameType, ":")
				uniform.Name = strings.TrimSpace(nameParts[0])
				typeName := strings.TrimSpace(strings.TrimSuffix(nameParts[1], ";"))
				uniform.Type = su.parseUniformType(typeName)
			}
		}
	} else if strings.HasPrefix(declaration, "var ") {
		// var textureName: texture_2d<f32>;
		parts := strings.Fields(declaration)
		if len(parts) >= 3 {
			nameType := parts[1]
			if strings.Contains(nameType, ":") {
				nameParts := strings.Split(nameType, ":")
				uniform.Name = strings.TrimSpace(nameParts[0])
				typeName := strings.TrimSpace(strings.TrimSuffix(nameParts[1], ";"))
				uniform.Type = su.parseUniformType(typeName)
			}
		}
	}

	return uniform
}

// parseUniformType converts a type string to UniformType.
func (su *ShaderUtils) parseUniformType(typeName string) UniformType {
	switch typeName {
	case "f32":
		return UniformTypeFloat
	case "vec2<f32>":
		return UniformTypeVec2
	case "vec3<f32>":
		return UniformTypeVec3
	case "vec4<f32>":
		return UniformTypeVec4
	case "mat2x2<f32>":
		return UniformTypeMat2
	case "mat3x3<f32>":
		return UniformTypeMat3
	case "mat4x4<f32>":
		return UniformTypeMat4
	case "i32":
		return UniformTypeInt
	case "vec2<i32>":
		return UniformTypeIVec2
	case "vec3<i32>":
		return UniformTypeIVec3
	case "vec4<i32>":
		return UniformTypeIVec4
	case "texture_2d<f32>":
		return UniformTypeSampler2D
	case "texture_cube<f32>":
		return UniformTypeSamplerCube
	default:
		return UniformTypeFloat
	}
}

// ValidateShaderSource performs basic validation on shader source.
func (su *ShaderUtils) ValidateShaderSource(source ShaderSource) []string {
	var issues []string

	if source.Name == "" {
		issues = append(issues, "shader name is empty")
	}

	if source.Source == "" {
		issues = append(issues, "shader source is empty")
	}

	// Check for common WGSL issues
	if strings.Contains(source.Source, "@vertex") && !strings.Contains(source.Source, "-> @builtin(position)") {
		issues = append(issues, "vertex shader should return position")
	}

	if strings.Contains(source.Source, "@fragment") && !strings.Contains(source.Source, "-> @location(0)") {
		issues = append(issues, "fragment shader should return color at location 0")
	}

	return issues
}

// CreateShaderVariant creates a shader variant with different defines.
func (su *ShaderUtils) CreateShaderVariant(base ShaderSource, defines map[string]string) ShaderSource {
	variant := base
	variant.Name = base.Name + "_variant"

	if variant.Defines == nil {
		variant.Defines = make(map[string]string)
	}

	// Merge defines
	for key, value := range defines {
		variant.Defines[key] = value
	}

	return variant
}

// ShaderTemplateProcessor processes shader templates with replacements.
type ShaderTemplateProcessor struct {
	replacements map[string]string
}

// NewShaderTemplateProcessor creates a new template processor.
func NewShaderTemplateProcessor() *ShaderTemplateProcessor {
	return &ShaderTemplateProcessor{
		replacements: make(map[string]string),
	}
}

// SetReplacement sets a template replacement.
func (stp *ShaderTemplateProcessor) SetReplacement(key, value string) {
	stp.replacements[key] = value
}

// ProcessTemplate processes a shader template with replacements.
func (stp *ShaderTemplateProcessor) ProcessTemplate(template string) string {
	result := template

	for key, value := range stp.replacements {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}
