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
	indexes  map[string]int
}

func New() *Registry {
	return &Registry{
		services: make(map[string][]ServiceInstance),
		indexes:  make(map[string]int),
	}
}

func (r *Registry) Register(service string, instance []ServiceInstance) {
	//lock
	r.mu.Lock()
	//Defer unlock
	defer r.mu.Unlock()

	r.services[service] = append(r.services[service], instance...)
}
func (r *Registry) Next(service string) (ServiceInstance, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	instances := r.services[service]

	if len(instances) == 0 {
		return ServiceInstance{}, false
	}
	index := r.indexes[service]
	instance := instances[index]

	r.indexes[service] = (index + 1) % len(instances)

	return instance, true
}
func (r *Registry) Get(service string) []ServiceInstance {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.services[service]
}
