package impl

import (
	"fmt"
	"sync"
)

// ResourceTracker tracks WebGPU resources for debugging and cleanup
type ResourceTracker struct {
	mu        sync.RWMutex
	resources map[uintptr]*ResourceInfo
	enabled   bool
}

// ResourceInfo contains information about a tracked resource
type ResourceInfo struct {
	Type      string
	Label     string
	CreatedAt string // stack trace or timestamp
}

var globalResourceTracker = &ResourceTracker{
	resources: make(map[uintptr]*ResourceInfo),
	enabled:   false, // Enable for debugging
}

// TrackResource registers a resource for tracking
func TrackResource(ptr uintptr, resourceType, label string) {
	if !globalResourceTracker.enabled {
		return
	}

	globalResourceTracker.mu.Lock()
	defer globalResourceTracker.mu.Unlock()

	globalResourceTracker.resources[ptr] = &ResourceInfo{
		Type:      resourceType,
		Label:     label,
		CreatedAt: "TODO: stack trace",
	}
}

// UntrackResource removes a resource from tracking
func UntrackResource(ptr uintptr) {
	if !globalResourceTracker.enabled {
		return
	}

	globalResourceTracker.mu.Lock()
	defer globalResourceTracker.mu.Unlock()

	delete(globalResourceTracker.resources, ptr)
}

// GetResourceReport returns a report of all tracked resources
func GetResourceReport() string {
	if !globalResourceTracker.enabled {
		return "Resource tracking is disabled"
	}

	globalResourceTracker.mu.RLock()
	defer globalResourceTracker.mu.RUnlock()

	report := fmt.Sprintf("Active WebGPU Resources: %d\n", len(globalResourceTracker.resources))
	for ptr, info := range globalResourceTracker.resources {
		report += fmt.Sprintf("  %x: %s '%s' \n", ptr, info.Type, info.Label)
	}

	return report
}

// EnableResourceTracking enables resource tracking for debugging
func EnableResourceTracking() {
	globalResourceTracker.enabled = true
}

// DisableResourceTracking disables resource tracking
func DisableResourceTracking() {
	globalResourceTracker.enabled = false
}
