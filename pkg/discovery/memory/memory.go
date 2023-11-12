package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"movieexample.com/pkg/discovery"
)

type servicename string
type instanceid string

// Registry define a in-memory service registry.
type Registry struct {
	sync.RWMutex
	serviceAddrs map[servicename]map[instanceid]*serviceInstance
}
type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// NewRegistry creates a new in-memory service registry instance.
func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[servicename]map[instanceid]*serviceInstance{}}
}

func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[servicename(serviceName)]; !ok {
		r.serviceAddrs[servicename(serviceName)] = map[instanceid]*serviceInstance{}
	}
	r.serviceAddrs[servicename(serviceName)][instanceid(instanceID)] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[servicename(serviceName)]; !ok {
		return nil
	}
	delete(r.serviceAddrs[servicename(serviceName)], instanceid(instanceID))
	return nil
}

// ReportHealthyState is a push mechanism for
// reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[servicename(serviceName)]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := r.serviceAddrs[servicename(serviceName)][instanceid(instanceID)]; !ok {
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddrs[servicename(serviceName)][instanceid(instanceID)].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of
// active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddrs[servicename(serviceName)]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, i := range r.serviceAddrs[servicename(serviceName)] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
