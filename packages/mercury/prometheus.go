package mercury

import (
	"context"
	"fmt"
	"net/http"

	"github.com/borchero/zeus/pkg/zeus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Prometheus wraps a set of properties associated with a Prometheus server.
type Prometheus struct {
	port uint16
}

// NewPrometheus creates a new Prometheus HTTP server on the provided port. The metrics will be
// provided on the path `/metrics` on the specified port.
func NewPrometheus(port uint16) *Prometheus {
	return &Prometheus{port}
}

// Run runs the HTTP server serving Prometheus metrics. It listens to the /metrics endpoint.
func (p *Prometheus) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", p.port),
		Handler: mux,
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		zeus.Logger(ctx).Info("http server started", zap.Uint16("port", p.port))
		err := s.ListenAndServe()
		switch err {
		case http.ErrServerClosed:
			// This is not actually an error, but rather a "graceful shutdown"
			return nil
		default:
			return err
		}
	})
	eg.Go(func() error {
		<-ctx.Done()
		if err := s.Shutdown(ctx); err != nil {
			zeus.Logger(ctx).Debug("failed to shut down http server", zap.Error(err))
		}
		return ctx.Err()
	})
	err := eg.Wait()
	zeus.Logger(ctx).Debug("http server exited")
	return err
}
