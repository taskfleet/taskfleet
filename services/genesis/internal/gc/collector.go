package gc

import (
	"context"
	"fmt"

	"github.com/borchero/zeus/pkg/zeus"
	"go.taskfleet.io/services/genesis/internal/db"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.uber.org/zap"
)

// GarbageCollector implements the logic for purging orphaned instances.
type GarbageCollector interface {

	// Collect attempts purging orphaned instances. If this attempt fails, an error is returned and
	// this method should be called again after a reasonable amount of time. An error is also
	// returned if the collection is not finished until the given context is cancelled.
	Collect(ctx context.Context) error
}

type garbageCollector struct {
	database db.Connection
	client   providers.Provider
}

// Options wraps the services that are required by the garbage collector.
type Options struct {
	Database db.Connection
	Client   providers.Provider
}

// NewGarbageCollector initializes a new garbage collector to collect orphanes instances.
func NewGarbageCollector(options Options) GarbageCollector {
	return &garbageCollector{
		database: options.Database,
		client:   options.Client,
	}
}

//-------------------------------------------------------------------------------------------------

func (d *garbageCollector) Collect(ctx context.Context) error {
	// First, we retrieve all instances that are actually running
	instances, err := d.client.Instances().List(ctx)
	if err != nil {
		return fmt.Errorf("failed listing instances running on GCP: %s", err)
	}
	zeus.Logger(ctx).Debug("successfully fetched instances from GCP",
		zap.Int("count", len(instances)))

	// Now we got the ground truth data. A couple of different failures could have happened and we
	// want to get rid of them. All these failure cases assume that the current state has not been
	// changed for 15 minutes. This equals the timeout for operations performed by genesis. Also,
	// all instances in the database are ignored if their `is_deletion_triaged` flag is set to
	// true.
	p := newPurger(d.client, d.database, instances)

	// 1. We requested an instance that was never flagged as booted. We delete it in the database
	//    and potentially delete it on the cloud provider if it is running.
	// 2. An instance was successfully created (i.e. is in booting state), but it was never flagged
	//    to be running. We delete it in the database and on the cloud provider.
	// 3. An instance is flagged to be running but cannot be found on the cloud provider. It is
	//    flagged as deleted in the database and the deletion is set to be triaged.
	// 4. An instance is set to be deleted but it is still shown as running on the cloud provider.
	//    It is deleted from the cloud provider and the deletion is set to be triaged on successful
	//    deletion.
	// 5. [Not actually a failure case] An instance was set to be deleted and cannot be found on
	//    the cloud provider. We triage its deletion.
	p.purgeRequested(ctx)
	p.purgeCreated(ctx)
	p.purgeRunning(ctx)
	p.purgeDeleted(ctx)

	return p.err
}
