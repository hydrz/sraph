package gio

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type DriverType uint8

const (
	DriverTypeUnknown DriverType = iota
	DriverTypeX11
	DriverTypeWayland
	DriverTypeCocoa
	DriverTypeWin32
	DriverTypeMobile
)

// Driver interface with context support
type Driver interface {
	Type() DriverType
	Name() string
	Version() string
	CreateWindow(ctx Context, o NewWindowOptions) (Window, error)
	Initialize(ctx Context) error
	Shutdown(ctx Context) error
	IsInitialized() bool
}

// DriverFactory creates driver instances
type DriverFactory func(config *Config) (Driver, error)

// DriverManager manages driver registration and lifecycle
type DriverManager struct {
	factories map[DriverType]DriverFactory
	drivers   map[DriverType]Driver
	mu        sync.RWMutex
}

var (
	globalDriverManager = &DriverManager{
		factories: make(map[DriverType]DriverFactory),
		drivers:   make(map[DriverType]Driver),
	}
)

// RegisterDriver registers a driver factory
func RegisterDriver(driverType DriverType, factory DriverFactory) {
	globalDriverManager.RegisterDriver(driverType, factory)
}

// GetDriver returns an initialized driver instance
func GetDriver(ctx context.Context, driverType DriverType) (Driver, error) {
	return globalDriverManager.GetDriver(ctx, driverType)
}

// GetAllDriverTypes returns all registered driver types
func GetAllDriverTypes() []DriverType {
	return globalDriverManager.GetAllDriverTypes()
}

// ShutdownAllDrivers shuts down all initialized drivers
func ShutdownAllDrivers(ctx context.Context) error {
	return globalDriverManager.ShutdownAll(ctx)
}

func (dm *DriverManager) RegisterDriver(driverType DriverType, factory DriverFactory) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.factories[driverType] = factory
}

func (dm *DriverManager) GetDriver(ctx context.Context, driverType DriverType) (Driver, error) {
	// Fast path: check if driver is already initialized
	dm.mu.RLock()
	if driver, exists := dm.drivers[driverType]; exists && driver.IsInitialized() {
		dm.mu.RUnlock()
		return driver, nil
	}
	dm.mu.RUnlock()

	// Slow path: need to create and initialize driver
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Double-check after acquiring write lock
	if driver, exists := dm.drivers[driverType]; exists && driver.IsInitialized() {
		return driver, nil
	}

	// Get factory
	factory, exists := dm.factories[driverType]
	if !exists {
		return nil, NewError(ErrorCodeDriverNotFound,
			fmt.Sprintf("driver type %d not registered", driverType), nil)
	}

	// Create driver with timeout context
	config := GetConfig()
	driver, err := factory(config)
	if err != nil {
		return nil, NewError(ErrorCodeDriverInitFailed, "failed to create driver", err)
	}

	// Initialize driver with timeout
	gioCtx := NewContext(ctx)
	initCtx, cancel := NewContextWithTimeout(gioCtx, time.Second*30)
	defer cancel()

	if err := driver.Initialize(initCtx); err != nil {
		return nil, NewError(ErrorCodeDriverInitFailed, "failed to initialize driver", err)
	}

	dm.drivers[driverType] = driver
	return driver, nil
}

func (dm *DriverManager) GetAllDriverTypes() []DriverType {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	types := make([]DriverType, 0, len(dm.factories))
	for driverType := range dm.factories {
		types = append(types, driverType)
	}
	return types
}

func (dm *DriverManager) ShutdownAll(ctx context.Context) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	gioCtx := NewContext(ctx)
	var errors []error

	// Shutdown all drivers and collect errors
	for driverType, driver := range dm.drivers {
		if driver.IsInitialized() {
			if err := driver.Shutdown(gioCtx); err != nil {
				errors = append(errors, fmt.Errorf("failed to shutdown driver %d: %w", driverType, err))
			}
		}
		delete(dm.drivers, driverType)
	}

	// Return combined error if any occurred
	if len(errors) > 0 {
		return fmt.Errorf("shutdown errors: %v", errors)
	}

	return nil
}
