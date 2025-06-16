package gio

import (
	"errors"
	"fmt"
)

// Error codes for better error handling
type ErrorCode int

const (
	ErrorCodeUnknown ErrorCode = iota
	ErrorCodeWindowNotInitialized
	ErrorCodeInvalidWindowState
	ErrorCodeDriverNotFound
	ErrorCodeDriverInitFailed
	ErrorCodeEventHandlerFull
	ErrorCodeInvalidConfiguration
	ErrorCodeResourceExhausted
	ErrorCodeOperationTimeout
)

// GioError represents a structured error in the gio package
type GioError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *GioError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("gio error %d: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("gio error %d: %s", e.Code, e.Message)
}

func (e *GioError) Unwrap() error {
	return e.Cause
}

// NewError creates a new GioError
func NewError(code ErrorCode, message string, cause error) *GioError {
	return &GioError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Predefined errors
var (
	ErrWindowNotInitialized = NewError(ErrorCodeWindowNotInitialized, "window not initialized", nil)
	ErrInvalidWindowState   = NewError(ErrorCodeInvalidWindowState, "invalid window state", nil)
	ErrDriverNotFound       = NewError(ErrorCodeDriverNotFound, "driver not found", nil)
	ErrDriverInitFailed     = NewError(ErrorCodeDriverInitFailed, "driver initialization failed", nil)
	ErrEventHandlerFull     = NewError(ErrorCodeEventHandlerFull, "event handler queue is full", nil)
	ErrInvalidConfiguration = NewError(ErrorCodeInvalidConfiguration, "invalid configuration", nil)
	ErrResourceExhausted    = NewError(ErrorCodeResourceExhausted, "system resources exhausted", nil)
	ErrOperationTimeout     = NewError(ErrorCodeOperationTimeout, "operation timeout", nil)
)

// IsError checks if an error is a specific GioError
func IsError(err error, code ErrorCode) bool {
	var gioErr *GioError
	if errors.As(err, &gioErr) {
		return gioErr.Code == code
	}
	return false
}
