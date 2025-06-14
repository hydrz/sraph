package wgpu

//go:generate go run gen.go -go wgpu.go

import (
	_ "github.com/goccy/go-yaml"
	_ "github.com/santhosh-tekuri/jsonschema/v5"
)

type Backend interface {
	Type() BackendType

	Init() error
	Release() error
}

var registeredBackends = make(map[BackendType]Backend)

func RegisterBackend(backend Backend) {
	if backend == nil {
		panic("gpu: RegisterBackend called with nil backend")
	}
	registeredBackends[backend.Type()] = backend
}

func GetBackend(backendType BackendType) (Backend, bool) {
	backend, exists := registeredBackends[backendType]
	if !exists {
		return nil, false
	}
	return backend, true
}

func GetAllBackends() []Backend {
	backends := make([]Backend, 0, len(registeredBackends))
	for _, backend := range registeredBackends {
		backends = append(backends, backend)
	}
	return backends
}
