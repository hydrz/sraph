// Package shader provides shader compilation and management for the Sraph UI toolkit.
//
// This package handles loading, compiling, and managing GPU shaders,
// including built-in shaders for common rendering operations.
package shader

import (
	"fmt"

	"github.com/opensraph/sraph/gpu"
)

// ShaderType defines the type of shader.
type ShaderType uint8

const (
	// ShaderTypeVertex represents a vertex shader.
	ShaderTypeVertex ShaderType = iota
	// ShaderTypeFragment represents a fragment/pixel shader.
	ShaderTypeFragment
	// ShaderTypeCompute represents a compute shader.
	ShaderTypeCompute
)

// Shader represents a compiled GPU shader.
type Shader interface {
	// ID returns the unique identifier for this shader.
	ID() ShaderID

	// Type returns the type of this shader.
	Type() ShaderType

	// Source returns the original shader source code.
	Source() string

	// Module returns the compiled GPU shader module.
	Module() gpu.ShaderModule

	// Uniforms returns the list of uniform variables in this shader.
	Uniforms() []UniformInfo

	// Destroy releases GPU resources associated with this shader.
	Destroy()
}

// ShaderID represents a unique identifier for a shader.
type ShaderID uint32

// UniformInfo describes a uniform variable in a shader.
type UniformInfo struct {
	Name     string
	Type     UniformType
	Size     int
	Offset   int
	Binding  int
	Location int
}

// UniformType defines the type of a uniform variable.
type UniformType uint8

const (
	UniformTypeFloat UniformType = iota
	UniformTypeVec2
	UniformTypeVec3
	UniformTypeVec4
	UniformTypeMat2
	UniformTypeMat3
	UniformTypeMat4
	UniformTypeInt
	UniformTypeIVec2
	UniformTypeIVec3
	UniformTypeIVec4
	UniformTypeSampler2D
	UniformTypeSamplerCube
)

// ShaderManager manages shader resources and compilation.
type ShaderManager interface {
	// LoadShader loads and compiles a shader from source.
	LoadShader(source string, shaderType ShaderType) (Shader, error)

	// LoadShaderFromFile loads and compiles a shader from a file.
	LoadShaderFromFile(filename string, shaderType ShaderType) (Shader, error)

	// GetShader retrieves a shader by ID.
	GetShader(id ShaderID) Shader

	// CreateProgram creates a shader program from vertex and fragment shaders.
	CreateProgram(vertexShader, fragmentShader Shader) (ShaderProgram, error)

	// GetBuiltinShader retrieves a built-in shader by name.
	GetBuiltinShader(name string) Shader

	// DestroyShader destroys a shader and releases its resources.
	DestroyShader(id ShaderID)

	// Shutdown releases all shader resources.
	Shutdown()
}

// ShaderProgram represents a complete shader program with vertex and fragment shaders.
type ShaderProgram interface {
	// ID returns the unique identifier for this program.
	ID() ProgramID

	// VertexShader returns the vertex shader.
	VertexShader() Shader

	// FragmentShader returns the fragment shader.
	FragmentShader() Shader

	// Uniforms returns all uniform variables in this program.
	Uniforms() []UniformInfo

	// GetUniformLocation returns the location of a uniform variable.
	GetUniformLocation(name string) int

	// RenderPipeline returns the GPU render pipeline for this program.
	RenderPipeline() gpu.RenderPipeline

	// Destroy releases GPU resources associated with this program.
	Destroy()
}

// ProgramID represents a unique identifier for a shader program.
type ProgramID uint32

// ShaderSource represents shader source code with metadata.
type ShaderSource struct {
	Name     string
	Source   string
	Type     ShaderType
	Includes []string
	Defines  map[string]string
}

// DefaultShaderManager provides the default implementation of ShaderManager.
type DefaultShaderManager struct {
	device         gpu.Device
	shaders        map[ShaderID]Shader
	programs       map[ProgramID]ShaderProgram
	nextShaderID   ShaderID
	nextProgramID  ProgramID
	builtinShaders map[string]Shader
}

// NewDefaultShaderManager creates a new default shader manager.
func NewDefaultShaderManager(device gpu.Device) *DefaultShaderManager {
	manager := &DefaultShaderManager{
		device:         device,
		shaders:        make(map[ShaderID]Shader),
		programs:       make(map[ProgramID]ShaderProgram),
		nextShaderID:   1,
		nextProgramID:  1,
		builtinShaders: make(map[string]Shader),
	}

	// Load built-in shaders
	manager.loadBuiltinShaders()

	return manager
}

// LoadShader implements ShaderManager.
func (m *DefaultShaderManager) LoadShader(source string, shaderType ShaderType) (Shader, error) {
	// Create shader module
	module := m.device.CreateShaderModule(gpu.ShaderModuleDescriptor{
		Label: fmt.Sprintf("Shader_%d", m.nextShaderID),
		Code:  source,
	})

	shader := &DefaultShader{
		id:         m.nextShaderID,
		shaderType: shaderType,
		source:     source,
		module:     module,
	}

	m.shaders[m.nextShaderID] = shader
	m.nextShaderID++

	return shader, nil
}

// LoadShaderFromFile implements ShaderManager.
func (m *DefaultShaderManager) LoadShaderFromFile(filename string, shaderType ShaderType) (Shader, error) {
	// TODO: Implement file loading
	return nil, fmt.Errorf("file loading not implemented")
}

// GetShader implements ShaderManager.
func (m *DefaultShaderManager) GetShader(id ShaderID) Shader {
	return m.shaders[id]
}

// CreateProgram implements ShaderManager.
func (m *DefaultShaderManager) CreateProgram(vertexShader, fragmentShader Shader) (ShaderProgram, error) {
	program := &DefaultShaderProgram{
		id:             m.nextProgramID,
		vertexShader:   vertexShader,
		fragmentShader: fragmentShader,
		device:         m.device,
	}

	// Create render pipeline
	err := program.createRenderPipeline()
	if err != nil {
		return nil, err
	}

	m.programs[m.nextProgramID] = program
	m.nextProgramID++

	return program, nil
}

// GetBuiltinShader implements ShaderManager.
func (m *DefaultShaderManager) GetBuiltinShader(name string) Shader {
	return m.builtinShaders[name]
}

// DestroyShader implements ShaderManager.
func (m *DefaultShaderManager) DestroyShader(id ShaderID) {
	if shader, exists := m.shaders[id]; exists {
		shader.Destroy()
		delete(m.shaders, id)
	}
}

// Shutdown implements ShaderManager.
func (m *DefaultShaderManager) Shutdown() {
	// Destroy all shaders
	for _, shader := range m.shaders {
		shader.Destroy()
	}
	m.shaders = make(map[ShaderID]Shader)

	// Destroy all programs
	for _, program := range m.programs {
		program.Destroy()
	}
	m.programs = make(map[ProgramID]ShaderProgram)
}

// loadBuiltinShaders loads the built-in shaders.
func (m *DefaultShaderManager) loadBuiltinShaders() {
	// TODO: Load built-in shaders from embedded source code
	// This would load shaders like solid color, texture, gradient, etc.
}

// DefaultShader provides the default implementation of Shader.
type DefaultShader struct {
	id         ShaderID
	shaderType ShaderType
	source     string
	module     gpu.ShaderModule
	uniforms   []UniformInfo
}

// ID implements Shader.
func (s *DefaultShader) ID() ShaderID {
	return s.id
}

// Type implements Shader.
func (s *DefaultShader) Type() ShaderType {
	return s.shaderType
}

// Source implements Shader.
func (s *DefaultShader) Source() string {
	return s.source
}

// Module implements Shader.
func (s *DefaultShader) Module() gpu.ShaderModule {
	return s.module
}

// Uniforms implements Shader.
func (s *DefaultShader) Uniforms() []UniformInfo {
	return s.uniforms
}

// Destroy implements Shader.
func (s *DefaultShader) Destroy() {
	s.module.Release()
}

// DefaultShaderProgram provides the default implementation of ShaderProgram.
type DefaultShaderProgram struct {
	id             ProgramID
	vertexShader   Shader
	fragmentShader Shader
	device         gpu.Device
	renderPipeline gpu.RenderPipeline
	uniforms       []UniformInfo
}

// ID implements ShaderProgram.
func (p *DefaultShaderProgram) ID() ProgramID {
	return p.id
}

// VertexShader implements ShaderProgram.
func (p *DefaultShaderProgram) VertexShader() Shader {
	return p.vertexShader
}

// FragmentShader implements ShaderProgram.
func (p *DefaultShaderProgram) FragmentShader() Shader {
	return p.fragmentShader
}

// Uniforms implements ShaderProgram.
func (p *DefaultShaderProgram) Uniforms() []UniformInfo {
	return p.uniforms
}

// GetUniformLocation implements ShaderProgram.
func (p *DefaultShaderProgram) GetUniformLocation(name string) int {
	for _, uniform := range p.uniforms {
		if uniform.Name == name {
			return uniform.Location
		}
	}
	return -1
}

// RenderPipeline implements ShaderProgram.
func (p *DefaultShaderProgram) RenderPipeline() gpu.RenderPipeline {
	return p.renderPipeline
}

// Destroy implements ShaderProgram.
func (p *DefaultShaderProgram) Destroy() {
	p.renderPipeline.Release()
}

// createRenderPipeline creates the GPU render pipeline for this program.
func (p *DefaultShaderProgram) createRenderPipeline() error {
	// TODO: Create render pipeline from vertex and fragment shaders
	// This would involve setting up the pipeline descriptor with
	// vertex layout, blend state, etc.
	return nil
}
