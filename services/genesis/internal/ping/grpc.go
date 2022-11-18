package ping

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/borchero/zeus/pkg/zeus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Grpc can be used to check if a gRPC server running a health server is ready.
type Grpc struct {
	dialOption grpc.DialOption
}

// NewGrpc initializes a new type able to check if a gRPC server is running. For all checks, the
// given TLS configuration is reused. If the given TLS configuration yields insecure credentials,
// no TLS is used.
func NewGrpc(tls *tls.Config) *Grpc {
	if tls == nil {
		return &Grpc{
			dialOption: grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	}
	return &Grpc{
		dialOption: grpc.WithTransportCredentials(credentials.NewTLS(tls)),
	}
}

// AwaitHealthy waits for the gRPC service at the specified hostname to be running. The error that
// this method returns is typically raised by the context being cancelled. The method tries to
// connect to the given IP indefinitely otherwise, ignoring errors encountered along the way.
func (g *Grpc) AwaitHealthy(ctx context.Context, host string, port uint16) error {
	address := fmt.Sprintf("%s:%d", host, port)
	var err error

	zeus.Logger(ctx).Debug("starting health check", zap.String("address", address))

	// First, we wait until we get a valid connection
	conn, err := grpc.DialContext(ctx, address, g.dialOption, grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed dialing target: %s", err)
	}
	defer conn.Close() // nolint:errcheck

	// Using the connection, we can initialize our client
	client := grpc_health_v1.NewHealthClient(conn)
	request := &grpc_health_v1.HealthCheckRequest{Service: ""}

	// Then, we run the health check until it returns a positive result (using exponential backoff)
	backoff := time.Second
	for {
		response, err := client.Check(ctx, request)
		if err == nil && response.Status == grpc_health_v1.HealthCheckResponse_SERVING {
			// In this case we're done!
			return nil
		} else if err == nil {
			zeus.Logger(ctx).Debug("grpc health check unknown response status",
				zap.Int32("response", int32(response.Status)))
		} else {
			zeus.Logger(ctx).Debug("grpc health check failed", zap.Error(err))
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			backoff = backoff * 2
		}
	}
}
