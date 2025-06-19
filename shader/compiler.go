// Package shader provides shader compilation utilities.
package shader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Compiler provides shader compilation and preprocessing capabilities.
type Compiler interface {
	// Compile compiles shader source code to the target format.
	Compile(source ShaderSource) (CompiledShader, error)

	// PreprocessShader preprocesses shader source with includes and defines.
	PreprocessShader(source ShaderSource) (string, error)

	// ValidateShader validates shader syntax and semantics.
	ValidateShader(source ShaderSource) error

	// OptimizeShader optimizes compiled shader code.
	OptimizeShader(shader CompiledShader) (CompiledShader, error)
}

// CompiledShader represents compiled shader bytecode.
type CompiledShader struct {
	Source   ShaderSource
	Bytecode []byte
	Metadata CompilationMetadata
}

// CompilationMetadata contains information about the compilation process.
type CompilationMetadata struct {
	TargetAPI    string
	Version      string
	Optimization OptimizationLevel
	Warnings     []string
	Errors       []string
}

// OptimizationLevel defines shader optimization levels.
type OptimizationLevel uint8

const (
	// OptimizationLevelNone performs no optimization.
	OptimizationLevelNone OptimizationLevel = iota
	// OptimizationLevelBasic performs basic optimizations.
	OptimizationLevelBasic
	// OptimizationLevelFull performs aggressive optimizations.
	OptimizationLevelFull
)

// WGSLCompiler compiles WGSL (WebGPU Shading Language) shaders.
type WGSLCompiler struct {
	optimization OptimizationLevel
	defines      map[string]string
	includePaths []string
}

// NewWGSLCompiler creates a new WGSL compiler.
func NewWGSLCompiler() *WGSLCompiler {
	return &WGSLCompiler{
		optimization: OptimizationLevelBasic,
		defines:      make(map[string]string),
		includePaths: make([]string, 0),
	}
}

// SetOptimizationLevel sets the optimization level for compilation.
func (c *WGSLCompiler) SetOptimizationLevel(level OptimizationLevel) {
	c.optimization = level
}

// AddDefine adds a preprocessor define.
func (c *WGSLCompiler) AddDefine(name, value string) {
	c.defines[name] = value
}

// AddIncludePath adds a path for shader includes.
func (c *WGSLCompiler) AddIncludePath(path string) {
	c.includePaths = append(c.includePaths, path)
}

// Compile implements Compiler.
func (c *WGSLCompiler) Compile(source ShaderSource) (CompiledShader, error) {
	// Preprocess the shader
	preprocessed, err := c.PreprocessShader(source)
	if err != nil {
		return CompiledShader{}, err
	}

	// Validate the shader
	validatedSource := source
	validatedSource.Source = preprocessed
	if err := c.ValidateShader(validatedSource); err != nil {
		return CompiledShader{}, err
	}

	// For WGSL, the "bytecode" is actually the preprocessed source
	// since WebGPU accepts WGSL source directly
	compiled := CompiledShader{
		Source:   validatedSource,
		Bytecode: []byte(preprocessed),
		Metadata: CompilationMetadata{
			TargetAPI:    "WebGPU",
			Version:      "1.0",
			Optimization: c.optimization,
			Warnings:     make([]string, 0),
			Errors:       make([]string, 0),
		},
	}

	return compiled, nil
}

// PreprocessShader implements Compiler.
func (c *WGSLCompiler) PreprocessShader(source ShaderSource) (string, error) {
	result := source.Source

	// Apply defines
	for name, value := range c.defines {
		define := fmt.Sprintf("#define %s %s", name, value)
		result = define + "\n" + result
	}

	// Apply shader-specific defines
	for name, value := range source.Defines {
		define := fmt.Sprintf("#define %s %s", name, value)
		result = define + "\n" + result
	}

	// Process includes
	processedResult, err := c.processIncludes(result, source.Includes)
	if err != nil {
		return "", err
	}

	return processedResult, nil
}

// processIncludes processes #include directives in shader source.
func (c *WGSLCompiler) processIncludes(source string, includes []string) (string, error) {
	lines := strings.Split(source, "\n")
	var result strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#include") {
			// Extract include filename
			parts := strings.Fields(trimmed)
			if len(parts) < 2 {
				return "", fmt.Errorf("invalid #include directive: %s", line)
			}

			filename := strings.Trim(parts[1], "\"<>")

			// Try to find and read the include file
			includeContent, err := c.readIncludeFile(filename)
			if err != nil {
				return "", fmt.Errorf("failed to read include file '%s': %v", filename, err)
			}

			result.WriteString(includeContent)
			result.WriteString("\n")
		} else {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

// readIncludeFile reads an include file from the include paths.
func (c *WGSLCompiler) readIncludeFile(filename string) (string, error) {
	// Try each include path
	for _, includePath := range c.includePaths {
		fullPath := filepath.Join(includePath, filename)
		if content, err := c.readFile(fullPath); err == nil {
			return content, nil
		}
	}

	// Try current directory
	if content, err := c.readFile(filename); err == nil {
		return content, nil
	}

	return "", fmt.Errorf("include file not found: %s", filename)
}

// readFile reads the content of a file.
func (c *WGSLCompiler) readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ValidateShader implements Compiler.
func (c *WGSLCompiler) ValidateShader(source ShaderSource) error {
	// Basic validation - check for required entry points
	switch source.Type {
	case ShaderTypeVertex:
		if !strings.Contains(source.Source, "@vertex") {
			return fmt.Errorf("vertex shader missing @vertex entry point")
		}
	case ShaderTypeFragment:
		if !strings.Contains(source.Source, "@fragment") {
			return fmt.Errorf("fragment shader missing @fragment entry point")
		}
	case ShaderTypeCompute:
		if !strings.Contains(source.Source, "@compute") {
			return fmt.Errorf("compute shader missing @compute entry point")
		}
	}

	// TODO: Implement more comprehensive validation
	// - Check for syntax errors
	// - Validate uniform declarations
	// - Check for undefined variables
	// - Validate function signatures

	return nil
}

// OptimizeShader implements Compiler.
func (c *WGSLCompiler) OptimizeShader(shader CompiledShader) (CompiledShader, error) {
	// For WGSL, optimization is typically handled by the GPU driver
	// But we could implement some basic optimizations here:
	// - Dead code elimination
	// - Constant folding
	// - Function inlining

	optimized := shader

	switch c.optimization {
	case OptimizationLevelNone:
		// No optimization
	case OptimizationLevelBasic:
		// Basic optimizations
		optimized = c.applyBasicOptimizations(shader)
	case OptimizationLevelFull:
		// Aggressive optimizations
		optimized = c.applyFullOptimizations(shader)
	}

	return optimized, nil
}

// applyBasicOptimizations applies basic shader optimizations.
func (c *WGSLCompiler) applyBasicOptimizations(shader CompiledShader) CompiledShader {
	// TODO: Implement basic optimizations
	// - Remove unused variables
	// - Simplify constant expressions
	return shader
}

// applyFullOptimizations applies aggressive shader optimizations.
func (c *WGSLCompiler) applyFullOptimizations(shader CompiledShader) CompiledShader {
	// TODO: Implement full optimizations
	// - Function inlining
	// - Loop unrolling
	// - Vectorization
	return shader
}

// CompilerError represents a shader compilation error.
type CompilerError struct {
	Line    int
	Column  int
	Message string
	Type    ErrorType
}

// ErrorType defines the type of compilation error.
type ErrorType uint8

const (
	ErrorTypeSyntax ErrorType = iota
	ErrorTypeSemantic
	ErrorTypeWarning
)

// Error implements error interface.
func (e CompilerError) Error() string {
	return fmt.Sprintf("line %d:%d: %s", e.Line, e.Column, e.Message)
}

// CompilerDiagnostics contains compilation diagnostics.
type CompilerDiagnostics struct {
	Errors   []CompilerError
	Warnings []CompilerError
}

// HasErrors returns whether there are compilation errors.
func (d CompilerDiagnostics) HasErrors() bool {
	for _, err := range d.Errors {
		if err.Type == ErrorTypeSyntax || err.Type == ErrorTypeSemantic {
			return true
		}
	}
	return false
}

// HasWarnings returns whether there are compilation warnings.
func (d CompilerDiagnostics) HasWarnings() bool {
	return len(d.Warnings) > 0
}
