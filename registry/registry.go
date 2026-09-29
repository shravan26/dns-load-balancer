package registry

import (
	"sync"
)

type ServiceInstance struct {
	Host string
	Port string
}

type Registry struct {
	mu       sync.RWMutex
	services map[string][]ServiceInstance
}

func New() *Registry {
	return &Registry{
		services: make(map[string][]ServiceInstance),
	}
}

func (r *Registry) Register(service string, instance ServiceInstance) {
	//lock
	r.mu.Lock()
	//Defer unlock
	defer r.mu.Unlock()

	r.services[service] = append(r.services[service], instance)
}

func (r *Registry) Get(service string) []ServiceInstance {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.services[service]
}
