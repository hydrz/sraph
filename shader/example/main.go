// Package main demonstrates usage of the Sraph shader package
package main

import (
	"fmt"
	"log"

	"github.com/opensraph/sraph/gpu"
	"github.com/opensraph/sraph/shader"
)

func main() {
	fmt.Println("Sraph Shader Package Demo")

	// Initialize GPU subsystem (placeholder)
	if err := gpu.Initialize(); err != nil {
		log.Fatal("Failed to initialize GPU:", err)
	}
	defer gpu.Shutdown()

	// Create a mock device for demo purposes
	// In real usage, you would get this from the GPU package
	device := createMockDevice()

	// Demo 1: Shader Manager Usage
	demoShaderManager(device)

	// Demo 2: Shader Compilation
	demoShaderCompilation()

	// Demo 3: Shader Bundling
	demoShaderBundling()

	// Demo 4: Built-in Shaders
	demoBuiltinShaders()

	// Demo 5: Shader Utilities
	demoShaderUtils()
}

func demoShaderManager(device gpu.Device) {
	fmt.Println("\n=== Shader Manager Demo ===")

	// Create shader manager
	manager := shader.NewDefaultShaderManager(device)
	defer manager.Shutdown()

	// Get built-in shaders
	solidShader := manager.GetBuiltinShader("solid")
	if solidShader != nil {
		fmt.Printf("Loaded built-in solid shader: ID=%d\n", solidShader.ID())
		fmt.Printf("Uniforms: %v\n", solidShader.Uniforms())
	}

	textureShader := manager.GetBuiltinShader("texture")
	if textureShader != nil {
		fmt.Printf("Loaded built-in texture shader: ID=%d\n", textureShader.ID())
	}

	gradientShader := manager.GetBuiltinShader("gradient")
	if gradientShader != nil {
		fmt.Printf("Loaded built-in gradient shader: ID=%d\n", gradientShader.ID())
	}

	// Create a shader program
	if solidShader != nil && textureShader != nil {
		program, err := manager.CreateProgram(solidShader, textureShader)
		if err != nil {
			fmt.Printf("Error creating program: %v\n", err)
		} else {
			fmt.Printf("Created shader program: ID=%d\n", program.ID())
		}
	}
}

func demoShaderCompilation() {
	fmt.Println("\n=== Shader Compilation Demo ===")

	// Create WGSL compiler
	compiler := shader.NewWGSLCompiler()
	compiler.SetOptimizationLevel(shader.OptimizationLevelBasic)
	compiler.AddDefine("USE_TEXTURE", "1")
	compiler.AddDefine("MAX_LIGHTS", "8")

	// Example shader source
	shaderSource := shader.ShaderSource{
		Name: "custom_vertex",
		Source: `
			struct VertexInput {
				@location(0) position: vec3<f32>,
				@location(1) normal: vec3<f32>,
			}

			struct VertexOutput {
				@builtin(position) clip_position: vec4<f32>,
				@location(0) world_normal: vec3<f32>,
			}

			@vertex
			fn vs_main(in: VertexInput) -> VertexOutput {
				var out: VertexOutput;
				out.clip_position = vec4<f32>(in.position, 1.0);
				out.world_normal = in.normal;
				return out;
			}
		`,
		Type:    shader.ShaderTypeVertex,
		Defines: map[string]string{"VERTEX_LIGHTING": "1"},
	}

	// Compile shader
	compiled, err := compiler.Compile(shaderSource)
	if err != nil {
		fmt.Printf("Compilation error: %v\n", err)
	} else {
		fmt.Printf("Successfully compiled shader '%s'\n", compiled.Source.Name)
		fmt.Printf("Bytecode size: %d bytes\n", len(compiled.Bytecode))
		fmt.Printf("Target API: %s\n", compiled.Metadata.TargetAPI)
	}

	// Validate shader
	if err := compiler.ValidateShader(shaderSource); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	} else {
		fmt.Println("Shader validation passed")
	}
}

func demoShaderBundling() {
	fmt.Println("\n=== Shader Bundling Demo ===")

	// Create compiler and bundle builder
	compiler := shader.NewWGSLCompiler()
	builder := shader.NewDefaultBundleBuilder(compiler)

	// Configure bundle
	builder.SetName("MyCustomShaders")
	builder.SetVersion("1.0.0")
	builder.SetDescription("Custom shaders for my application")

	// Add built-in shaders to demonstrate
	solidSource, _ := shader.GetBuiltinShaderSource("solid")
	solidShaderSource := shader.ShaderSource{
		Name:   "solid",
		Source: solidSource,
		Type:   shader.ShaderTypeVertex,
	}

	compiled, err := compiler.Compile(solidShaderSource)
	if err == nil {
		builder.AddShader("custom_solid", compiled)
	}

	// Build bundle
	bundle, err := builder.Build()
	if err != nil {
		fmt.Printf("Bundle creation error: %v\n", err)
		return
	}

	fmt.Printf("Created bundle: %s v%s\n", bundle.Name(), bundle.Version())
	fmt.Printf("Bundle size: %d bytes\n", bundle.Size())
	fmt.Printf("Shaders in bundle: %v\n", bundle.ListShaders())

	// Test shader retrieval
	if bundle.HasShader("custom_solid") {
		shader, err := bundle.GetShader("custom_solid")
		if err != nil {
			fmt.Printf("Error retrieving shader: %v\n", err)
		} else {
			fmt.Printf("Retrieved shader: %s\n", shader.Source.Name)
		}
	}
}

func demoBuiltinShaders() {
	fmt.Println("\n=== Built-in Shaders Demo ===")

	// List all built-in shaders
	fmt.Println("Available built-in shaders:")
	for _, name := range shader.BuiltinShaderNames {
		info, exists := shader.GetBuiltinShaderInfo(name)
		if exists {
			fmt.Printf("  %s: %s\n", info.Name, info.Description)
			fmt.Printf("    Uniforms: %v\n", info.Uniforms)
		}
	}

	// Get shader source
	solidSource, exists := shader.GetBuiltinShaderSource("solid")
	if exists {
		fmt.Printf("\nSolid shader source (first 200 chars):\n%s...\n",
			truncateString(solidSource, 200))
	}
}

func demoShaderUtils() {
	fmt.Println("\n=== Shader Utilities Demo ===")

	utils := shader.NewShaderUtils()

	// Test shader type detection
	vertexSource := `
		@vertex
		fn vs_main() -> @builtin(position) vec4<f32> {
			return vec4<f32>(0.0, 0.0, 0.0, 1.0);
		}
	`

	shaderType, err := utils.ParseShaderType(vertexSource)
	if err != nil {
		fmt.Printf("Type detection error: %v\n", err)
	} else {
		fmt.Printf("Detected shader type: %d (vertex)\n", shaderType)
	}

	// Test uniform extraction
	gradientSource, _ := shader.GetBuiltinShaderSource("gradient")
	uniforms := utils.ExtractUniforms(gradientSource)
	fmt.Printf("Extracted %d uniforms from gradient shader:\n", len(uniforms))
	for _, uniform := range uniforms {
		fmt.Printf("  %s (binding: %d)\n", uniform.Name, uniform.Binding)
	}

	// Test shader validation
	testSource := shader.ShaderSource{
		Name:   "test",
		Source: vertexSource,
		Type:   shader.ShaderTypeVertex,
	}

	issues := utils.ValidateShaderSource(testSource)
	if len(issues) > 0 {
		fmt.Printf("Validation issues: %v\n", issues)
	} else {
		fmt.Println("Shader validation passed")
	}

	// Test template processing
	template := `
		struct Uniforms {
			color: vec4<f32>,
			{{EXTRA_FIELDS}}
		}

		@fragment
		fn fs_main() -> @location(0) vec4<f32> {
			return {{DEFAULT_COLOR}};
		}
	`

	processor := shader.NewShaderTemplateProcessor()
	processor.SetReplacement("EXTRA_FIELDS", "transform: mat4x4<f32>,")
	processor.SetReplacement("DEFAULT_COLOR", "vec4<f32>(1.0, 0.0, 0.0, 1.0)")

	processed := processor.ProcessTemplate(template)
	fmt.Printf("\nTemplate processing result:\n%s\n", processed)
}

// Mock device for demo purposes
func createMockDevice() gpu.Device {
	// In a real implementation, this would return an actual GPU device
	// For now, return nil as we're just demonstrating the API
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
