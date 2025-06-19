// Package shader provides shader bundling and asset management.
package shader

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Bundle represents a collection of compiled shaders.
type Bundle interface {
	// Name returns the name of this shader bundle.
	Name() string

	// Version returns the version of this shader bundle.
	Version() string

	// GetShader retrieves a shader from the bundle by name.
	GetShader(name string) (CompiledShader, error)

	// ListShaders returns a list of all shader names in the bundle.
	ListShaders() []string

	// HasShader checks if a shader exists in the bundle.
	HasShader(name string) bool

	// Size returns the total size of the bundle in bytes.
	Size() int64

	// Metadata returns bundle metadata.
	Metadata() BundleMetadata
}

// BundleMetadata contains information about a shader bundle.
type BundleMetadata struct {
	Name        string
	Version     string
	Description string
	Author      string
	BuildTime   string
	TargetAPI   string
	ShaderCount int
	Checksum    string
}

// BundleBuilder helps build shader bundles from multiple sources.
type BundleBuilder interface {
	// SetName sets the bundle name.
	SetName(name string)

	// SetVersion sets the bundle version.
	SetVersion(version string)

	// SetDescription sets the bundle description.
	SetDescription(description string)

	// AddShader adds a compiled shader to the bundle.
	AddShader(name string, shader CompiledShader) error

	// AddShaderFromFile adds a shader from a file.
	AddShaderFromFile(name, filename string, shaderType ShaderType) error

	// AddShadersFromDirectory adds all shaders from a directory.
	AddShadersFromDirectory(directory string) error

	// Build creates the final shader bundle.
	Build() (Bundle, error)

	// SaveToFile saves the bundle to a file.
	SaveToFile(filename string) error
}

// DefaultBundle provides the default implementation of Bundle.
type DefaultBundle struct {
	metadata BundleMetadata
	shaders  map[string]CompiledShader
	size     int64
}

// NewDefaultBundle creates a new default bundle.
func NewDefaultBundle(metadata BundleMetadata) *DefaultBundle {
	return &DefaultBundle{
		metadata: metadata,
		shaders:  make(map[string]CompiledShader),
		size:     0,
	}
}

// Name implements Bundle.
func (b *DefaultBundle) Name() string {
	return b.metadata.Name
}

// Version implements Bundle.
func (b *DefaultBundle) Version() string {
	return b.metadata.Version
}

// GetShader implements Bundle.
func (b *DefaultBundle) GetShader(name string) (CompiledShader, error) {
	shader, exists := b.shaders[name]
	if !exists {
		return CompiledShader{}, fmt.Errorf("shader '%s' not found in bundle", name)
	}
	return shader, nil
}

// ListShaders implements Bundle.
func (b *DefaultBundle) ListShaders() []string {
	names := make([]string, 0, len(b.shaders))
	for name := range b.shaders {
		names = append(names, name)
	}
	return names
}

// HasShader implements Bundle.
func (b *DefaultBundle) HasShader(name string) bool {
	_, exists := b.shaders[name]
	return exists
}

// Size implements Bundle.
func (b *DefaultBundle) Size() int64 {
	return b.size
}

// Metadata implements Bundle.
func (b *DefaultBundle) Metadata() BundleMetadata {
	return b.metadata
}

// addShader adds a shader to the bundle (internal method).
func (b *DefaultBundle) addShader(name string, shader CompiledShader) {
	b.shaders[name] = shader
	b.size += int64(len(shader.Bytecode))
	b.metadata.ShaderCount = len(b.shaders)
}

// DefaultBundleBuilder provides the default implementation of BundleBuilder.
type DefaultBundleBuilder struct {
	metadata BundleMetadata
	shaders  map[string]CompiledShader
	compiler Compiler
}

// NewDefaultBundleBuilder creates a new default bundle builder.
func NewDefaultBundleBuilder(compiler Compiler) *DefaultBundleBuilder {
	return &DefaultBundleBuilder{
		metadata: BundleMetadata{
			TargetAPI: "WebGPU",
			Version:   "1.0.0",
		},
		shaders:  make(map[string]CompiledShader),
		compiler: compiler,
	}
}

// SetName implements BundleBuilder.
func (b *DefaultBundleBuilder) SetName(name string) {
	b.metadata.Name = name
}

// SetVersion implements BundleBuilder.
func (b *DefaultBundleBuilder) SetVersion(version string) {
	b.metadata.Version = version
}

// SetDescription implements BundleBuilder.
func (b *DefaultBundleBuilder) SetDescription(description string) {
	b.metadata.Description = description
}

// AddShader implements BundleBuilder.
func (b *DefaultBundleBuilder) AddShader(name string, shader CompiledShader) error {
	if _, exists := b.shaders[name]; exists {
		return fmt.Errorf("shader '%s' already exists in bundle", name)
	}
	b.shaders[name] = shader
	return nil
}

// AddShaderFromFile implements BundleBuilder.
func (b *DefaultBundleBuilder) AddShaderFromFile(name, filename string, shaderType ShaderType) error {
	// Read shader source from file
	source, err := b.readFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read shader file '%s': %v", filename, err)
	}

	shaderSource := ShaderSource{
		Name:    name,
		Source:  source,
		Type:    shaderType,
		Defines: make(map[string]string),
	}

	compiled, err := b.compiler.Compile(shaderSource)
	if err != nil {
		return fmt.Errorf("failed to compile shader '%s': %v", name, err)
	}

	return b.AddShader(name, compiled)
}

// readFile reads the content of a file.
func (b *DefaultBundleBuilder) readFile(filename string) (string, error) {
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

// AddShadersFromDirectory implements BundleBuilder.
func (b *DefaultBundleBuilder) AddShadersFromDirectory(directory string) error {
	return filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Determine shader type from file extension
		ext := filepath.Ext(path)
		var shaderType ShaderType

		switch ext {
		case ".vert", ".vs":
			shaderType = ShaderTypeVertex
		case ".frag", ".fs":
			shaderType = ShaderTypeFragment
		case ".comp", ".cs":
			shaderType = ShaderTypeCompute
		case ".wgsl":
			// For WGSL files, we need to parse the content to determine type
			// For now, default to vertex
			shaderType = ShaderTypeVertex
		default:
			// Skip unknown file types
			return nil
		}

		// Use filename without extension as shader name
		name := filepath.Base(path)
		name = name[:len(name)-len(ext)]

		return b.AddShaderFromFile(name, path, shaderType)
	})
}

// Build implements BundleBuilder.
func (b *DefaultBundleBuilder) Build() (Bundle, error) {
	bundle := NewDefaultBundle(b.metadata)

	for name, shader := range b.shaders {
		bundle.addShader(name, shader)
	}

	return bundle, nil
}

// SaveToFile implements BundleBuilder.
func (b *DefaultBundleBuilder) SaveToFile(filename string) error {
	// TODO: Implement bundle serialization
	// This would serialize the bundle to a binary format
	// that can be loaded efficiently at runtime
	return fmt.Errorf("bundle serialization not implemented")
}

// BundleLoader loads shader bundles from various sources.
type BundleLoader interface {
	// LoadFromFile loads a bundle from a file.
	LoadFromFile(filename string) (Bundle, error)

	// LoadFromBytes loads a bundle from byte data.
	LoadFromBytes(data []byte) (Bundle, error)

	// LoadEmbedded loads an embedded bundle.
	LoadEmbedded(name string) (Bundle, error)
}

// DefaultBundleLoader provides the default implementation of BundleLoader.
type DefaultBundleLoader struct {
	embeddedBundles map[string]Bundle
}

// NewDefaultBundleLoader creates a new default bundle loader.
func NewDefaultBundleLoader() *DefaultBundleLoader {
	return &DefaultBundleLoader{
		embeddedBundles: make(map[string]Bundle),
	}
}

// LoadFromFile implements BundleLoader.
func (l *DefaultBundleLoader) LoadFromFile(filename string) (Bundle, error) {
	// TODO: Implement bundle deserialization
	return nil, fmt.Errorf("bundle loading not implemented")
}

// LoadFromBytes implements BundleLoader.
func (l *DefaultBundleLoader) LoadFromBytes(data []byte) (Bundle, error) {
	// TODO: Implement bundle deserialization from bytes
	return nil, fmt.Errorf("bundle loading not implemented")
}

// LoadEmbedded implements BundleLoader.
func (l *DefaultBundleLoader) LoadEmbedded(name string) (Bundle, error) {
	bundle, exists := l.embeddedBundles[name]
	if !exists {
		return nil, fmt.Errorf("embedded bundle '%s' not found", name)
	}
	return bundle, nil
}

// RegisterEmbeddedBundle registers an embedded bundle.
func (l *DefaultBundleLoader) RegisterEmbeddedBundle(name string, bundle Bundle) {
	l.embeddedBundles[name] = bundle
}

// BundleRegistry manages multiple shader bundles.
type BundleRegistry struct {
	bundles map[string]Bundle
	loader  BundleLoader
}

// NewBundleRegistry creates a new bundle registry.
func NewBundleRegistry(loader BundleLoader) *BundleRegistry {
	return &BundleRegistry{
		bundles: make(map[string]Bundle),
		loader:  loader,
	}
}

// RegisterBundle registers a bundle in the registry.
func (r *BundleRegistry) RegisterBundle(bundle Bundle) {
	r.bundles[bundle.Name()] = bundle
}

// GetBundle retrieves a bundle by name.
func (r *BundleRegistry) GetBundle(name string) (Bundle, error) {
	bundle, exists := r.bundles[name]
	if !exists {
		// Try to load the bundle
		loadedBundle, err := r.loader.LoadEmbedded(name)
		if err != nil {
			return nil, fmt.Errorf("bundle '%s' not found", name)
		}
		r.bundles[name] = loadedBundle
		return loadedBundle, nil
	}
	return bundle, nil
}

// GetShader retrieves a shader from any registered bundle.
func (r *BundleRegistry) GetShader(bundleName, shaderName string) (CompiledShader, error) {
	bundle, err := r.GetBundle(bundleName)
	if err != nil {
		return CompiledShader{}, err
	}
	return bundle.GetShader(shaderName)
}

// ListBundles returns a list of all registered bundle names.
func (r *BundleRegistry) ListBundles() []string {
	names := make([]string, 0, len(r.bundles))
	for name := range r.bundles {
		names = append(names, name)
	}
	return names
}
