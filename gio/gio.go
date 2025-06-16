// Package gio provides a set of types and interfaces for handling events in a graphical input/output context.
package gio

import (
	"fmt"
	"runtime"
)

// EnumerateDrivers returns all available driver types on the current platform
func EnumerateDrivers() []DriverType {
	availableDrivers := make([]DriverType, 0)

	// Check each registered driver type
	for driverType := range registeredDrivers {
		availableDrivers = append(availableDrivers, driverType)
	}

	return availableDrivers
}

// GetBestDriver returns the best available driver for the current platform
func GetBestDriver() (Driver, error) {
	// Platform-specific driver priority order
	var preferredOrder []DriverType

	switch runtime.GOOS {
	case "linux":
		// On Linux, prefer Wayland if available, then X11
		preferredOrder = []DriverType{
			DriverTypeWayland,
			DriverTypeX11,
		}
	case "darwin":
		// On macOS, prefer Cocoa
		preferredOrder = []DriverType{
			DriverTypeCocoa,
		}
	case "windows":
		// On Windows, prefer Win32
		preferredOrder = []DriverType{
			DriverTypeWin32,
		}
	case "android", "ios":
		// On mobile platforms
		preferredOrder = []DriverType{
			DriverTypeMobile,
		}
	default:
		// For unknown platforms, try all available drivers
		preferredOrder = GetAllDriverTypes()
	}

	// Try to get drivers in preferred order
	for _, driverType := range preferredOrder {
		if driver, err := GetDriver(driverType); err == nil {
			return driver, nil
		}
	}

	// If no preferred driver is available, try any available driver
	availableDrivers := GetAllDriverTypes()
	if len(availableDrivers) == 0 {
		return nil, NewError(ErrorCodeDriverNotFound,
			"no graphics drivers available on this platform", nil)
	}

	for _, driverType := range availableDrivers {
		if driver, err := GetDriver(driverType); err == nil {
			return driver, nil
		}
	}

	return nil, NewError(ErrorCodeDriverInitFailed,
		"failed to initialize any available graphics driver", nil)
}

// CreateWindow creates a window using the best available driver
func CreateWindow(o ...NewWindowOptions) (Window, error) {
	// Merge options with defaults
	attr := DefaultNewWindowOptions()
	if len(o) > 0 {
		attr = attr.Apply(o...)
	}

	// Get the best available driver
	driver, err := GetBestDriver()
	if err != nil {
		return nil, fmt.Errorf("gio: failed to get graphics driver: %w", err)
	}

	// Create window using the driver
	window, err := driver.CreateWindow(attr)
	if err != nil {
		return nil, fmt.Errorf("gio: failed to create window: %w", err)
	}

	return window, nil
}

// CreateWindowWithDriver creates a window using a specific driver type
func CreateWindowWithDriver(driverType DriverType, o ...NewWindowOptions) (Window, error) {
	// Merge options with defaults
	attr := DefaultNewWindowOptions()
	if len(o) > 0 {
		attr = attr.Apply(o...)
	}

	// Get the specific driver
	driver, err := GetDriver(driverType)
	if err != nil {
		return nil, fmt.Errorf("gio: failed to get %d driver: %w", driverType, err)
	}

	// Create window using the driver
	window, err := driver.CreateWindow(attr)
	if err != nil {
		return nil, fmt.Errorf("gio: failed to create window with driver %d: %w", driverType, err)
	}

	return window, nil
}

// GetPlatformInfo returns information about the current platform and available drivers
func GetPlatformInfo() PlatformInfo {
	availableDrivers := EnumerateDrivers()

	var bestDriver DriverType = DriverTypeUnknown
	if driver, err := GetBestDriver(); err == nil {
		bestDriver = driver.Type()
	}

	return PlatformInfo{
		OS:               runtime.GOOS,
		Arch:             runtime.GOARCH,
		AvailableDrivers: availableDrivers,
		BestDriver:       bestDriver,
	}
}

// PlatformInfo contains information about the current platform
type PlatformInfo struct {
	OS               string       // Operating system (linux, windows, darwin, etc.)
	Arch             string       // Architecture (amd64, arm64, etc.)
	AvailableDrivers []DriverType // List of available graphics drivers
	BestDriver       DriverType   // The best driver for this platform
}

// String returns a human-readable description of the platform info
func (pi PlatformInfo) String() string {
	driverNames := make([]string, len(pi.AvailableDrivers))
	for i, dt := range pi.AvailableDrivers {
		driverNames[i] = dt.String()
	}

	return fmt.Sprintf("Platform: %s/%s, Available drivers: %v, Best driver: %s",
		pi.OS, pi.Arch, driverNames, pi.BestDriver.String())
}
