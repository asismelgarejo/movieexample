package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Registry defines a service registry-
type Registry interface {
	// Register creates a server instance record in the registry.
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error

	// Deregister removes a server instance record from the registry.
	Deregister(ctx context.Context, instanceID string, serviceName string) error

	// ServiceAddress returns the list of addresses of the active instances of a given service.
	ServiceAddress(ctx context.Context, serviceID string) ([]string, error)

	// ReportHealthyState is a push mechanism to report healthy state to the registry.
	ReportHealthyState(instanceID string, serviceName string) error
}

// ErrNotFound is returned when no server addresses are found
var ErrNotFound = errors.New("no service addresses found")

// GenerateInstanceID generates a pseudo-random service
// instance identifier, using a service name
// suffixed by dash and a random number.
func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
