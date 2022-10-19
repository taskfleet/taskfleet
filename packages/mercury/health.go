package mercury

import (
	"sync"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Health represents a type which ought to be used to capture the health status of the application.
type Health interface {
	// SetHealthy sets the health status of the application.
	SetHealthy(healthy bool)
}

type grpcHealth struct {
	health *health.Server
	mutex  sync.Mutex
}

// Health returns an instance which allows setting the health status of the gRPC server. If the
// gRPC server does not have the health service enabled, this function panics.
func (i *Grpc) Health() Health {
	if i.health == nil {
		panic("grpc server does not have a health server set")
	}
	return &grpcHealth{health: i.health}
}

// SetHealthy implements the Health interface.
func (h *grpcHealth) SetHealthy(healthy bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if healthy {
		h.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	} else {
		h.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}
}
