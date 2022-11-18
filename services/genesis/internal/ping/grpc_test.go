package ping

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type timing struct {
	startDelay   time.Duration
	healthyDelay time.Duration
	timeout      time.Duration
}

func TestImmediateRunning(t *testing.T) {
	err := runTest(t, timing{
		startDelay:   time.Duration(0),
		healthyDelay: time.Duration(0),
		timeout:      5 * time.Second,
	})
	assert.Nil(t, err)
}

func TestStartupDelay(t *testing.T) {
	err := runTest(t, timing{
		startDelay:   3 * time.Second,
		healthyDelay: time.Duration(0),
		timeout:      10 * time.Second,
	})
	assert.Nil(t, err)
}

func TestHealthyDelay(t *testing.T) {
	err := runTest(t, timing{
		startDelay:   3 * time.Second,
		healthyDelay: 5 * time.Second,
		timeout:      15 * time.Second,
	})
	assert.Nil(t, err)
}

//-------------------------------------------------------------------------------------------------

type mockServer struct {
	*grpc_health_v1.UnimplementedHealthServer
	deadline time.Time
}

func (s *mockServer) Check(
	ctx context.Context, request *grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	if time.Now().Before(s.deadline) {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

//-------------------------------------------------------------------------------------------------

func runTest(t *testing.T, deadlines timing) error {
	done := make(chan struct{}, 1)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, deadlines.timeout)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return runHealthServer(t, deadlines, done)
	})
	eg.Go(func() error {
		pinger := NewGrpc(nil)
		err := pinger.AwaitHealthy(ctx, "localhost", 5404)
		// Stop grpc server
		done <- struct{}{}
		time.Sleep(time.Second)
		return err
	})
	return eg.Wait()
}

func runHealthServer(t *testing.T, deadlines timing, done chan struct{}) error {
	time.Sleep(deadlines.startDelay)

	sock, err := net.Listen("tcp", "127.0.0.1:5404")
	if err != nil {
		return fmt.Errorf("failed to listen on port 5404: %s", err)
	}

	server := grpc.NewServer()
	service := &mockServer{deadline: time.Now().Add(deadlines.startDelay + deadlines.healthyDelay)}
	grpc_health_v1.RegisterHealthServer(server, service)

	go func() {
		select {
		case <-time.After(deadlines.timeout - deadlines.startDelay):
		case <-done:
		}
		server.GracefulStop()
	}()

	if err := server.Serve(sock); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %s", err)
	}

	return nil
}
