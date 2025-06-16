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
