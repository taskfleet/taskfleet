package mercury

import (
	"context"
	"fmt"
	"net"

	"github.com/borchero/zeus/pkg/zeus"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

// Grpc wraps a set of properties associated with a gRPC server.
type Grpc struct {
	Server     *grpc.Server
	health     *health.Server
	port       uint16
	tlsSetting string
}

// NewGrpc creates a new gRPC server. The caller must call `Run` on the returned instance to
// actually start the server.
func NewGrpc(port uint16, options ...GrpcOption) (*Grpc, error) {
	// Setup interceptors
	serverOptions := []grpc.ServerOption{}
	for _, option := range options {
		serverOptions = append(serverOptions, option.serverOptions()...)
	}

	// Setup server
	grpc := Grpc{
		Server:     grpc.NewServer(serverOptions...),
		port:       port,
		tlsSetting: "disabled",
	}
	for _, option := range options {
		option.apply(&grpc)
	}
	return &grpc, nil
}

// Run runs the gRPC server and terminates once the context expires.
func (g *Grpc) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		sock, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
		if err != nil {
			return fmt.Errorf("grpc server failed to create socket: %s", err)
		}
		zeus.Logger(ctx).Info(
			"grpc server started",
			zap.Uint16("port", g.port),
			zap.String("tls", g.tlsSetting),
		)
		return g.Server.Serve(sock)
	})
	eg.Go(func() error {
		<-ctx.Done()
		g.Server.Stop()
		return ctx.Err()
	})
	err := eg.Wait()
	zeus.Logger(ctx).Debug("grpc server exited")
	return err
}
