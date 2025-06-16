package gio

import "fmt"

type DriverType uint8

const (
	DriverTypeUnknown DriverType = iota
	DriverTypeX11
	DriverTypeWayland
	DriverTypeCocoa
	DriverTypeWin32
	DriverTypeMobile
)

// String returns the string representation of the driver type
func (dt DriverType) String() string {
	switch dt {
	case DriverTypeX11:
		return "X11"
	case DriverTypeWayland:
		return "Wayland"
	case DriverTypeCocoa:
		return "Cocoa"
	case DriverTypeWin32:
		return "Win32"
	case DriverTypeMobile:
		return "Mobile"
	default:
		return "Unknown"
	}
}

// IsAvailable checks if the driver type is available on the current platform
func (dt DriverType) IsAvailable() bool {
	_, ok := registeredDrivers[dt]
	return ok
}

// Driver interface with context support
type Driver interface {
	Type() DriverType
	CreateWindow(o NewWindowOptions) (Window, error)
}

// DriverFactory creates driver instances
type DriverFactory func() (Driver, error)

var registeredDrivers = make(map[DriverType]DriverFactory)

// RegisterDriver registers a driver factory
func RegisterDriver(driverType DriverType, factory DriverFactory) {
	if factory == nil {
		panic("gio: RegisterDriver factory cannot be nil")
	}
	registeredDrivers[driverType] = factory
}

// GetDriver returns an initialized driver instance
func GetDriver(driverType DriverType) (Driver, error) {
	if factory, ok := registeredDrivers[driverType]; ok {
		return factory()
	}
	return nil, fmt.Errorf("gio: driver type %d not registered", driverType)
}

// GetAllDriverTypes returns all registered driver types
func GetAllDriverTypes() []DriverType {
	types := make([]DriverType, 0, len(registeredDrivers))
	for dt := range registeredDrivers {
		types = append(types, dt)
	}
	return types
}
