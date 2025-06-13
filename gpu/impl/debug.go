package impl

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// DebugInfo contains debugging information
type DebugInfo struct {
	ObjectType  string
	ObjectID    uint64
	Label       string
	CreatedAt   time.Time
	IsDestroyed bool
}

// DebugManager manages debugging and profiling information
type DebugManager struct {
	mu           sync.RWMutex
	objects      map[uint64]*DebugInfo
	enabled      bool
	nextObjectID uint64
}

var globalDebugManager = &DebugManager{
	objects: make(map[uint64]*DebugInfo),
	enabled: false,
}

// EnableDebug enables debug mode
func EnableDebug() {
	globalDebugManager.mu.Lock()
	defer globalDebugManager.mu.Unlock()
	globalDebugManager.enabled = true
}

// DisableDebug disables debug mode
func DisableDebug() {
	globalDebugManager.mu.Lock()
	defer globalDebugManager.mu.Unlock()
	globalDebugManager.enabled = false
}

// RegisterObject registers an object for debugging
func RegisterDebugObject(objectType, label string) uint64 {
	if !globalDebugManager.enabled {
		return 0
	}

	globalDebugManager.mu.Lock()
	defer globalDebugManager.mu.Unlock()

	globalDebugManager.nextObjectID++
	objectID := globalDebugManager.nextObjectID

	globalDebugManager.objects[objectID] = &DebugInfo{
		ObjectType:  objectType,
		ObjectID:    objectID,
		Label:       label,
		CreatedAt:   time.Now(),
		IsDestroyed: false,
	}

	return objectID
}

// GetDebugReport returns a debug report of all objects
func GetDebugReport() string {
	if !globalDebugManager.enabled {
		return "Debug mode is disabled"
	}

	globalDebugManager.mu.RLock()
	defer globalDebugManager.mu.RUnlock()

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("WebGPU Debug Report (%d objects)\n", len(globalDebugManager.objects)))
	builder.WriteString("===============================================\n")

	for _, info := range globalDebugManager.objects {
		status := "Active"
		if info.IsDestroyed {
			status = "Destroyed"
		}

		builder.WriteString(fmt.Sprintf(
			"[%s] %s '%s' , Created: %s\n",
			status, info.ObjectType, info.Label, info.CreatedAt.Format(time.RFC3339),
		))
	}

	return builder.String()
}

// Performance profiler for debugging
type PerformanceProfiler struct {
	mu      sync.RWMutex
	metrics map[string]*PerformanceMetric
}

type PerformanceMetric struct {
	TotalCalls      uint64
	TotalDuration   time.Duration
	AverageDuration time.Duration
	LastCall        time.Time
}

var globalProfiler = &PerformanceProfiler{
	metrics: make(map[string]*PerformanceMetric),
}

// StartProfiling starts profiling for a function
func StartProfiling(name string) func() {
	if !globalDebugManager.enabled {
		return func() {}
	}

	start := time.Now()
	return func() {
		duration := time.Since(start)

		globalProfiler.mu.Lock()
		defer globalProfiler.mu.Unlock()

		metric, exists := globalProfiler.metrics[name]
		if !exists {
			metric = &PerformanceMetric{}
			globalProfiler.metrics[name] = metric
		}

		metric.TotalCalls++
		metric.TotalDuration += duration
		metric.AverageDuration = metric.TotalDuration / time.Duration(metric.TotalCalls)
		metric.LastCall = time.Now()
	}
}

// GetPerformanceReport returns a performance report
func GetPerformanceReport() string {
	if !globalDebugManager.enabled {
		return "Debug mode is disabled"
	}

	globalProfiler.mu.RLock()
	defer globalProfiler.mu.RUnlock()

	var builder strings.Builder
	builder.WriteString("WebGPU Performance Report\n")
	builder.WriteString("=========================\n")

	for name, metric := range globalProfiler.metrics {
		builder.WriteString(fmt.Sprintf(
			"%s: %d calls, Total: %v, Avg: %v, Last: %s\n",
			name, metric.TotalCalls, metric.TotalDuration,
			metric.AverageDuration, metric.LastCall.Format(time.RFC3339),
		))
	}

	return builder.String()
}
