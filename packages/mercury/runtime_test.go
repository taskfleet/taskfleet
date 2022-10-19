package mercury

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunning(t *testing.T) {
	ctx := context.Background()

	checkRunning(ctx, t, []GrpcOption{})

	checkRunning(ctx, t, []GrpcOption{
		WithPrometheusMetrics(),
	})

	checkRunning(ctx, t, []GrpcOption{
		WithRequestValidation(),
	})

	checkRunning(ctx, t, []GrpcOption{
		WithHealthService(),
	})

	checkRunning(ctx, t, []GrpcOption{
		WithPrometheusMetrics(),
		WithRequestValidation(),
	})

	checkRunning(ctx, t, []GrpcOption{
		WithPrometheusMetrics(),
		WithHealthService(),
	})

	checkRunning(ctx, t, []GrpcOption{
		WithRequestValidation(),
		WithHealthService(),
	})

	checkRunning(ctx, t, []GrpcOption{
		WithPrometheusMetrics(),
		WithRequestValidation(),
		WithHealthService(),
	})
}

func checkRunning(
	ctx context.Context,
	t *testing.T,
	grpcOpts []GrpcOption,
) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	grpc, err := NewGrpc(5404, grpcOpts...)
	require.Nil(t, err)
	prometheus := NewPrometheus(9090)
	rt := NewRuntime(ctx).
		Schedule("grpc", grpc).
		Schedule("prometheus", prometheus)

	// Test whether ports are open
	time.Sleep(250 * time.Millisecond)
	checkPort(ctx, t, 5404)
	checkPort(ctx, t, 9090)

	assert.EqualValues(t, context.DeadlineExceeded, rt.Await())
}

func checkPort(ctx context.Context, t *testing.T, port uint16) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 500*time.Millisecond)
	assert.Nil(t, err)
	if err == nil {
		assert.Nil(t, conn.Close())
	}
}
