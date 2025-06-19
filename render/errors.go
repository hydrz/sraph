package render

import (
	"errors"
	"fmt"
)

// Common render errors
var (
	ErrInvalidPipelineConfiguration = errors.New("invalid pipeline configuration")
	ErrNotImplemented               = errors.New("not implemented")
	ErrMissingComputeShader         = errors.New("missing compute shader")
	ErrMissingVertexShader          = errors.New("missing vertex shader")
	ErrMissingFragmentShader        = errors.New("missing fragment shader")
	ErrInvalidShaderStage           = errors.New("invalid shader stage")
	ErrInvalidVertexDescriptor      = errors.New("invalid vertex descriptor")
	ErrInvalidRenderTarget          = errors.New("invalid render target")
	ErrInvalidTexture               = errors.New("invalid texture")
	ErrInvalidBuffer                = errors.New("invalid buffer")
	ErrDeviceNotFound               = errors.New("device not found")
	ErrContextNotInitialized        = errors.New("context not initialized")
	ErrCommandBufferInvalidState    = errors.New("command buffer in invalid state")
	ErrRenderPassNotActive          = errors.New("render pass not active")
	ErrComputePassNotActive         = errors.New("compute pass not active")
)

// ValidationError represents a validation error with detailed information
type ValidationError struct {
	Message string
	Field   string
	Value   interface{}
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string, value interface{}) *ValidationError {
	return &ValidationError{
		Message: message,
		Field:   field,
		Value:   value,
	}
}

// ShaderCompilationError represents a shader compilation error
type ShaderCompilationError struct {
	ShaderSource string
	Message      string
	Line         int
	Column       int
}

// Error implements the error interface
func (e *ShaderCompilationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("shader compilation error at line %d:%d: %s", e.Line, e.Column, e.Message)
	}
	return fmt.Sprintf("shader compilation error: %s", e.Message)
}

// ResourceError represents a resource-related error
type ResourceError struct {
	ResourceType string
	ResourceID   string
	Operation    string
	Cause        error
}

// Error implements the error interface
func (e *ResourceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("resource error: %s operation on %s '%s' failed: %v",
			e.Operation, e.ResourceType, e.ResourceID, e.Cause)
	}
	return fmt.Sprintf("resource error: %s operation on %s '%s' failed",
		e.Operation, e.ResourceType, e.ResourceID)
}

// Unwrap returns the underlying cause
func (e *ResourceError) Unwrap() error {
	return e.Cause
}

// ContextError represents a context-related error
type ContextError struct {
	Operation string
	Message   string
}

// Error implements the error interface
func (e *ContextError) Error() string {
	return fmt.Sprintf("context error during %s: %s", e.Operation, e.Message)
}

// DeviceError represents a device-related error
type DeviceError struct {
	DeviceType string
	Operation  string
	Message    string
}

// Error implements the error interface
func (e *DeviceError) Error() string {
	return fmt.Sprintf("device error (%s) during %s: %s", e.DeviceType, e.Operation, e.Message)
}

// PipelineError represents a pipeline-related error
type PipelineError struct {
	PipelineType string
	Stage        string
	Message      string
}

// Error implements the error interface
func (e *PipelineError) Error() string {
	if e.Stage != "" {
		return fmt.Sprintf("pipeline error (%s) in stage %s: %s", e.PipelineType, e.Stage, e.Message)
	}
	return fmt.Sprintf("pipeline error (%s): %s", e.PipelineType, e.Message)
}

// ErrorWithContext wraps an error with additional context
type ErrorWithContext struct {
	Err     error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *ErrorWithContext) Error() string {
	return e.Err.Error()
}

// Unwrap returns the underlying error
func (e *ErrorWithContext) Unwrap() error {
	return e.Err
}

// GetContext returns the error context
func (e *ErrorWithContext) GetContext() map[string]interface{} {
	return e.Context
}

// WithContext adds context to an error
func WithContext(err error, context map[string]interface{}) *ErrorWithContext {
	return &ErrorWithContext{
		Err:     err,
		Context: context,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsShaderError checks if an error is a shader compilation error
func IsShaderError(err error) bool {
	var shaderErr *ShaderCompilationError
	return errors.As(err, &shaderErr)
}

// IsResourceError checks if an error is a resource error
func IsResourceError(err error) bool {
	var resourceErr *ResourceError
	return errors.As(err, &resourceErr)
}
