package impl

import (
	"fmt"
	"sync"
	"sync/atomic"
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
	RefCount  int32
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
		RefCount:  1,
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

// AddResourceRef increments reference count
func AddResourceRef(ptr uintptr) {
	if !globalResourceTracker.enabled {
		return
	}

	globalResourceTracker.mu.RLock()
	resource, exists := globalResourceTracker.resources[ptr]
	globalResourceTracker.mu.RUnlock()

	if exists {
		atomic.AddInt32(&resource.RefCount, 1)
	}
}

// ReleaseResourceRef decrements reference count
func ReleaseResourceRef(ptr uintptr) {
	if !globalResourceTracker.enabled {
		return
	}

	globalResourceTracker.mu.RLock()
	resource, exists := globalResourceTracker.resources[ptr]
	globalResourceTracker.mu.RUnlock()

	if exists {
		if atomic.AddInt32(&resource.RefCount, -1) <= 0 {
			UntrackResource(ptr)
		}
	}
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
		report += fmt.Sprintf("  %x: %s '%s' (refs: %d)\n", ptr, info.Type, info.Label, info.RefCount)
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
