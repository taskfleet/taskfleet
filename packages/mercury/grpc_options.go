package mercury

import (
	"crypto/tls"
	"strings"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// GrpcOption allows to customize the gRPC server.
type GrpcOption interface {
	apply(grpc *Grpc)
	serverOptions() []grpc.ServerOption
}

//-------------------------------------------------------------------------------------------------
// PROMETHEUS METRICS
//-------------------------------------------------------------------------------------------------

type grpcOptionPrometheus struct{}

// WithPrometheusMetrics enables Prometheus metrics for streaming and unary RPC calls.
func WithPrometheusMetrics() GrpcOption {
	return grpcOptionPrometheus{}
}

func (grpcOptionPrometheus) apply(grpc *Grpc) {
	// noop
}

func (grpcOptionPrometheus) serverOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.ChainStreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	}
}

//-------------------------------------------------------------------------------------------------
// REQUEST VALIDATION
//-------------------------------------------------------------------------------------------------

type grpcOptionRequestValidation struct{}

// WithRequestValidation enables automatic validation of incoming requests for streaming and unary
// API calls. Invalid argument errors will be returned automatically.
func WithRequestValidation() GrpcOption {
	return grpcOptionRequestValidation{}
}

func (grpcOptionRequestValidation) apply(grpc *Grpc) {
	// noop
}

func (grpcOptionRequestValidation) serverOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(grpc_validator.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(grpc_validator.StreamServerInterceptor()),
	}
}

//-------------------------------------------------------------------------------------------------
// LOGGING
//-------------------------------------------------------------------------------------------------

type grpcOptionLogging struct {
	logger           *zap.Logger
	filterHealthLogs bool
}

// WithLogger sets the specified logger to log incoming requests and their error codes. If health
// logs are filtered, all requests to the `/grpc.health` service will not be logged.
func WithLogger(logger *zap.Logger, filterHealthLogs bool) GrpcOption {
	return grpcOptionLogging{logger, filterHealthLogs}
}

func (grpcOptionLogging) apply(grpc *Grpc) {
	// noop
}

func (o grpcOptionLogging) serverOptions() []grpc.ServerOption {
	options := []grpc_zap.Option{}
	if o.filterHealthLogs {
		options = append(options, grpc_zap.WithDecider(excludeGrpcHealthLogs))
	}
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(grpc_zap.UnaryServerInterceptor(o.logger, options...)),
		grpc.ChainStreamInterceptor(grpc_zap.StreamServerInterceptor(o.logger, options...)),
	}
}

func excludeGrpcHealthLogs(path string, err error) bool {
	if err != nil {
		return true
	}
	return !strings.HasPrefix(path, "/grpc.health")
}

//-------------------------------------------------------------------------------------------------
// HEALTH
//-------------------------------------------------------------------------------------------------

type grpcOptionHealth struct{}

// WithHealthService attaches the gRPC health service to the gRPC server.
func WithHealthService() GrpcOption {
	return grpcOptionHealth{}
}

func (grpcOptionHealth) apply(grpc *Grpc) {
	grpc.health = health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpc.Server, grpc.health)
}

func (grpcOptionHealth) serverOptions() []grpc.ServerOption {
	return []grpc.ServerOption{}
}

//-------------------------------------------------------------------------------------------------
// TLS
//-------------------------------------------------------------------------------------------------

type grpcOptionTLS struct {
	config *tls.Config
}

// WithTLS uses the specified TLS configuration for the gRPC server. If a value of `nil` is
// provided, the server will not use TLS.
func WithTLS(config *tls.Config) GrpcOption {
	return grpcOptionTLS{config}
}

func (o grpcOptionTLS) apply(grpc *Grpc) {
	if o.config == nil {
		grpc.tlsSetting = "disabled"
	} else if o.config.ClientCAs != nil && o.config.ClientAuth == tls.RequireAndVerifyClientCert {
		grpc.tlsSetting = "mtls"
	} else {
		grpc.tlsSetting = "enabled"
	}
}

func (o grpcOptionTLS) serverOptions() []grpc.ServerOption {
	if o.config == nil {
		return nil
	}
	return []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(o.config)),
	}
}
