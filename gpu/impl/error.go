package impl

import (
	"log"
	"sync"

	. "github.com/opensraph/sraph/gpu/webgpu"
)

// ErrorScope represents a pushed error scope
type ErrorScope struct {
	Filter ErrorFilter
	Errors []CapturedError
}

// CapturedError represents a captured error
type CapturedError struct {
	Type    ErrorType
	Message string
}

// ErrorManager manages error scopes and uncaptured errors
type ErrorManager struct {
	mu                    sync.RWMutex
	errorScopes           []ErrorScope
	uncapturedCallback    UncapturedErrorCallbackInfo
	hasUncapturedCallback bool
}

// NewErrorManager creates a new error manager
func NewErrorManager() *ErrorManager {
	return &ErrorManager{
		errorScopes: make([]ErrorScope, 0),
	}
}

// PushErrorScope pushes an error scope
func (em *ErrorManager) PushErrorScope(filter ErrorFilter) {
	em.mu.Lock()
	defer em.mu.Unlock()

	scope := ErrorScope{
		Filter: filter,
		Errors: make([]CapturedError, 0),
	}
	em.errorScopes = append(em.errorScopes, scope)
}

// PopErrorScope pops an error scope and returns captured errors
func (em *ErrorManager) PopErrorScope() (ErrorType, string, bool) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if len(em.errorScopes) == 0 {
		return ErrorTypeNoError, "", false
	}

	// Pop the top scope
	scope := em.errorScopes[len(em.errorScopes)-1]
	em.errorScopes = em.errorScopes[:len(em.errorScopes)-1]

	// Return the first error if any
	if len(scope.Errors) > 0 {
		err := scope.Errors[0]
		return err.Type, err.Message, true
	}

	return ErrorTypeNoError, "", true
}

// CaptureError captures an error in the current scope
func (em *ErrorManager) CaptureError(errorType ErrorType, message string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	// Check if error should be captured by any scope
	for i := len(em.errorScopes) - 1; i >= 0; i-- {
		scope := &em.errorScopes[i]
		if em.shouldCaptureError(scope.Filter, errorType) {
			scope.Errors = append(scope.Errors, CapturedError{
				Type:    errorType,
				Message: message,
			})
			return
		}
	}

	// No scope captured the error, call uncaptured error callback
	if em.hasUncapturedCallback && em.uncapturedCallback.Callback != nil {
		em.uncapturedCallback.Callback(nil, errorType, message)
	} else {
		// Default logging if no callback is set
		log.Printf("Uncaptured WebGPU error [%v]: %s", errorType, message)
	}
}

// SetUncapturedErrorCallback sets the uncaptured error callback
func (em *ErrorManager) SetUncapturedErrorCallback(callback UncapturedErrorCallbackInfo) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.uncapturedCallback = callback
	em.hasUncapturedCallback = true
}

// shouldCaptureError checks if an error should be captured by a filter
func (em *ErrorManager) shouldCaptureError(filter ErrorFilter, errorType ErrorType) bool {
	switch filter {
	case ErrorFilterValidation:
		return errorType == ErrorTypeValidation
	case ErrorFilterOutOfMemory:
		return errorType == ErrorTypeOutOfMemory
	case ErrorFilterInternal:
		return errorType == ErrorTypeInternal
	default:
		return false
	}
}

// Global error manager
var globalErrorManager = NewErrorManager()

// GetGlobalErrorManager returns the global error manager
func GetGlobalErrorManager() *ErrorManager {
	return globalErrorManager
}
