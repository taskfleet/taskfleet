package mercury

import (
	"context"

	"github.com/borchero/zeus/pkg/zeus"
	"golang.org/x/sync/errgroup"
)

// Runnable is a type which can be run and exists as a context is cancelled.
// It may return an error.
type Runnable interface {
	Run(ctx context.Context) error
}

// Runtime is a type which enables running a set of runnables concurrently and await their
// completion
type Runtime struct {
	eg  *errgroup.Group
	ctx context.Context
}

// NewRuntime initializes a new runtime which exits once the given context expires.
func NewRuntime(ctx context.Context) *Runtime {
	r := Runtime{}
	r.eg, r.ctx = errgroup.WithContext(ctx)
	return &r
}

// Schedule adds the given runnable to the runtime. The name is used to assign a logger name to the
// context such that each runnable logs in its own "namespace".
func (r *Runtime) Schedule(name string, runnable Runnable) *Runtime {
	r.eg.Go(func() error {
		return runnable.Run(zeus.WithName(r.ctx, name))
	})
	return r
}

// Await waits for all processes of the runtime to have finished.
func (r *Runtime) Await() error {
	return r.eg.Wait()
}
