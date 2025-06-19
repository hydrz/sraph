// Package app provides the application framework and runtime for the Sraph UI toolkit.
//
// This package defines the core application structure, including the engine,
// scheduler, state management, and lifecycle management.
package app

import (
	"context"
	"time"

	"github.com/opensraph/sraph/gio"
	"github.com/opensraph/sraph/render"
	"github.com/opensraph/sraph/widget"
)

// App represents the main application instance.
type App interface {
	// Run starts the application with the given root widget.
	Run(ctx context.Context, root widget.Widget) error

	// Stop stops the application gracefully.
	Stop()

	// Engine returns the application engine.
	Engine() Engine

	// Scheduler returns the application scheduler.
	Scheduler() Scheduler

	// State returns the application state manager.
	State() StateManager
}

// Engine represents the core application engine.
type Engine interface {
	// Initialize initializes the engine with the given configuration.
	Initialize(config *EngineConfig) error

	// Shutdown shuts down the engine.
	Shutdown()

	// Renderer returns the renderer instance.
	Renderer() render.Renderer

	// WindowManager returns the window manager.
	WindowManager() gio.WindowManager

	// Frame performs a single frame update.
	Frame(deltaTime time.Duration) error

	// SetRootWidget sets the root widget for the application.
	SetRootWidget(root widget.Widget)

	// GetRootWidget returns the current root widget.
	GetRootWidget() widget.Widget
}

// EngineConfig contains configuration options for the engine.
type EngineConfig struct {
	// WindowTitle specifies the initial window title.
	WindowTitle string

	// WindowWidth specifies the initial window width.
	WindowWidth int

	// WindowHeight specifies the initial window height.
	WindowHeight int

	// Resizable indicates if the window should be resizable.
	Resizable bool

	// VSync indicates if vertical synchronization should be enabled.
	VSync bool

	// MaxFPS specifies the maximum frames per second (0 for unlimited).
	MaxFPS int

	// DebugMode indicates if debug features should be enabled.
	DebugMode bool
}

// Scheduler manages the execution of tasks and frame updates.
type Scheduler interface {
	// Schedule schedules a task to run on the next frame.
	Schedule(task Task)

	// ScheduleWithDelay schedules a task to run after the specified delay.
	ScheduleWithDelay(task Task, delay time.Duration)

	// ScheduleRepeating schedules a task to run repeatedly with the given interval.
	ScheduleRepeating(task Task, interval time.Duration) ScheduledTask

	// RequestFrame requests a frame update.
	RequestFrame()

	// SetFrameCallback sets the callback to be called on each frame.
	SetFrameCallback(callback FrameCallback)

	// Start starts the scheduler.
	Start()

	// Stop stops the scheduler.
	Stop()
}

// Task represents a task that can be scheduled for execution.
type Task interface {
	// Execute executes the task.
	Execute()

	// Name returns the name of the task for debugging purposes.
	Name() string
}

// ScheduledTask represents a scheduled task that can be cancelled.
type ScheduledTask interface {
	// Cancel cancels the scheduled task.
	Cancel()

	// IsActive returns true if the task is still active.
	IsActive() bool
}

// FrameCallback is called on each frame update.
type FrameCallback func(deltaTime time.Duration)

// StateManager manages application state and state changes.
type StateManager interface {
	// SetState sets a state value with the given key.
	SetState(key string, value interface{})

	// GetState gets a state value by key.
	GetState(key string) (interface{}, bool)

	// RemoveState removes a state value by key.
	RemoveState(key string)

	// Subscribe subscribes to state changes for the given key.
	Subscribe(key string, callback StateChangeCallback) Subscription

	// Clear clears all state values.
	Clear()
}

// StateChangeCallback is called when a state value changes.
type StateChangeCallback func(key string, oldValue, newValue interface{})

// Subscription represents a subscription to state changes.
type Subscription interface {
	// Unsubscribe unsubscribes from state changes.
	Unsubscribe()

	// IsActive returns true if the subscription is still active.
	IsActive() bool
}

// LifecycleManager manages application lifecycle events.
type LifecycleManager interface {
	// AddListener adds a lifecycle event listener.
	AddListener(listener LifecycleListener)

	// RemoveListener removes a lifecycle event listener.
	RemoveListener(listener LifecycleListener)

	// NotifyStarted notifies all listeners that the application has started.
	NotifyStarted()

	// NotifyPaused notifies all listeners that the application has been paused.
	NotifyPaused()

	// NotifyResumed notifies all listeners that the application has been resumed.
	NotifyResumed()

	// NotifyWillTerminate notifies all listeners that the application will terminate.
	NotifyWillTerminate()
}

// LifecycleListener listens to application lifecycle events.
type LifecycleListener interface {
	// OnStarted is called when the application has started.
	OnStarted()

	// OnPaused is called when the application has been paused.
	OnPaused()

	// OnResumed is called when the application has been resumed.
	OnResumed()

	// OnWillTerminate is called when the application will terminate.
	OnWillTerminate()
}

// ResourceManager manages application resources such as textures, fonts, etc.
type ResourceManager interface {
	// LoadTexture loads a texture resource.
	LoadTexture(path string) (interface{}, error)

	// LoadFont loads a font resource.
	LoadFont(path string) (interface{}, error)

	// LoadShader loads a shader resource.
	LoadShader(path string) (interface{}, error)

	// UnloadResource unloads a resource by path.
	UnloadResource(path string)

	// GetResource gets a loaded resource by path.
	GetResource(path string) (interface{}, bool)

	// Clear clears all loaded resources.
	Clear()
}

// SimpleTask is a simple implementation of Task using a function.
type SimpleTask struct {
	name string
	fn   func()
}

// Execute implements Task.
func (t *SimpleTask) Execute() {
	if t.fn != nil {
		t.fn()
	}
}

// Name implements Task.
func (t *SimpleTask) Name() string {
	return t.name
}

// NewTask creates a new SimpleTask with the given name and function.
func NewTask(name string, fn func()) Task {
	return &SimpleTask{
		name: name,
		fn:   fn,
	}
}

// NewApp creates a new application instance.
func NewApp() App {
	// TODO: Implement app creation
	return nil
}

// NewEngine creates a new engine instance.
func NewEngine() Engine {
	// TODO: Implement engine creation
	return nil
}

// NewScheduler creates a new scheduler instance.
func NewScheduler() Scheduler {
	// TODO: Implement scheduler creation
	return nil
}

// NewStateManager creates a new state manager instance.
func NewStateManager() StateManager {
	// TODO: Implement state manager creation
	return nil
}

// NewLifecycleManager creates a new lifecycle manager instance.
func NewLifecycleManager() LifecycleManager {
	// TODO: Implement lifecycle manager creation
	return nil
}

// NewResourceManager creates a new resource manager instance.
func NewResourceManager() ResourceManager {
	// TODO: Implement resource manager creation
	return nil
}
