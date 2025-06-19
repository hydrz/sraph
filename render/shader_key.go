package render

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// ShaderKey represents a unique identifier for shader functions and pipelines
// Used for caching and identifying compiled shaders
type ShaderKey struct {
	hash [32]byte
}

// NewShaderKey creates a new shader key from the given data
func NewShaderKey(data []byte) ShaderKey {
	return ShaderKey{
		hash: sha256.Sum256(data),
	}
}

// NewShaderKeyFromString creates a shader key from a string
func NewShaderKeyFromString(str string) ShaderKey {
	return NewShaderKey([]byte(str))
}

// NewShaderKeyFromComponents creates a shader key from multiple components
func NewShaderKeyFromComponents(components ...interface{}) ShaderKey {
	hasher := sha256.New()

	for _, component := range components {
		switch v := component.(type) {
		case string:
			hasher.Write([]byte(v))
		case []byte:
			hasher.Write(v)
		case int:
			binary.Write(hasher, binary.LittleEndian, int64(v))
		case int32:
			binary.Write(hasher, binary.LittleEndian, v)
		case int64:
			binary.Write(hasher, binary.LittleEndian, v)
		case uint32:
			binary.Write(hasher, binary.LittleEndian, v)
		case uint64:
			binary.Write(hasher, binary.LittleEndian, v)
		case float32:
			binary.Write(hasher, binary.LittleEndian, v)
		case float64:
			binary.Write(hasher, binary.LittleEndian, v)
		case bool:
			if v {
				hasher.Write([]byte{1})
			} else {
				hasher.Write([]byte{0})
			}
		default:
			// For other types, use string representation
			hasher.Write([]byte(fmt.Sprintf("%v", v)))
		}
	}

	var key ShaderKey
	copy(key.hash[:], hasher.Sum(nil))
	return key
}

// String returns the string representation of the shader key
func (k ShaderKey) String() string {
	return fmt.Sprintf("%x", k.hash)
}

// Bytes returns the raw bytes of the shader key
func (k ShaderKey) Bytes() []byte {
	return k.hash[:]
}

// Equal returns true if two shader keys are equal
func (k ShaderKey) Equal(other ShaderKey) bool {
	return k.hash == other.hash
}

// IsValid returns true if the shader key is valid (non-zero)
func (k ShaderKey) IsValid() bool {
	for _, b := range k.hash {
		if b != 0 {
			return true
		}
	}
	return false
}

// ShaderKeyBuilder helps build shader keys incrementally
type ShaderKeyBuilder struct {
	components []interface{}
}

// NewShaderKeyBuilder creates a new shader key builder
func NewShaderKeyBuilder() *ShaderKeyBuilder {
	return &ShaderKeyBuilder{
		components: make([]interface{}, 0, 8),
	}
}

// Add adds a component to the shader key
func (b *ShaderKeyBuilder) Add(component interface{}) *ShaderKeyBuilder {
	b.components = append(b.components, component)
	return b
}

// AddString adds a string component
func (b *ShaderKeyBuilder) AddString(str string) *ShaderKeyBuilder {
	return b.Add(str)
}

// AddInt adds an integer component
func (b *ShaderKeyBuilder) AddInt(value int) *ShaderKeyBuilder {
	return b.Add(value)
}

// AddBool adds a boolean component
func (b *ShaderKeyBuilder) AddBool(value bool) *ShaderKeyBuilder {
	return b.Add(value)
}

// AddBytes adds a byte slice component
func (b *ShaderKeyBuilder) AddBytes(data []byte) *ShaderKeyBuilder {
	return b.Add(data)
}

// Build creates the final shader key
func (b *ShaderKeyBuilder) Build() ShaderKey {
	return NewShaderKeyFromComponents(b.components...)
}

// Reset clears the builder for reuse
func (b *ShaderKeyBuilder) Reset() {
	b.components = b.components[:0]
}
