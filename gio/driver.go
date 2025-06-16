package gio

import (
	"fmt"
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

type Driver interface {
	Type() DriverType
	CreateWindow(o ...NewWindowOptions) (Window, error)
}

type DirverFactory func() (Driver, error)

var registersDriver = make(map[DriverType]DirverFactory)

func RegisterDriver(driverType DriverType, factory DirverFactory) {
	registersDriver[driverType] = factory
}

func GetDriver(driverType DriverType) (Driver, error) {
	factory, ok := registersDriver[driverType]
	if !ok {
		return nil, fmt.Errorf("driver type %d not registered", driverType)
	}
	return factory()
}

func GetAllDrivers() []DriverType {
	var drivers []DriverType
	for driverType := range registersDriver {
		drivers = append(drivers, driverType)
	}
	return drivers
}
