package gc

import (
	"context"
	"fmt"
	"time"

	"github.com/borchero/zeus/pkg/zeus"
	"github.com/google/uuid"
	db "go.taskfleet.io/services/genesis/db/gen"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.uber.org/zap"
)

type purger struct {
	client    providers.Provider
	database  db.Connection
	instances map[uuid.UUID]providers.Instance
	err       error
}

func newPurger(
	client providers.Provider, database db.Connection, instances []providers.Instance,
) *purger {
	instanceMap := make(map[uuid.UUID]providers.Instance)
	for _, i := range instances {
		instanceMap[i.Meta.ID] = i
	}
	return &purger{client: client, database: database, instances: instanceMap}
}

func (p *purger) purgeRequested(ctx context.Context) {
	p.purgeStatus(ctx, db.InstanceStatusRequested)
}

func (p *purger) purgeCreated(ctx context.Context) {
	p.purgeStatus(ctx, db.InstanceStatusBooting)
}

func (p *purger) purgeRunning(ctx context.Context) {
	if p.err != nil {
		return
	}

	// First, we get all instances which are currently running
	iterator := p.database.ListInstances(ctx, db.InstanceStatusRunning)
	count := 0
	_, err := iterator.ForEach(func(instance *db.Instance) error {
		// If the instance cannot be found on the cloud provider, we delete it locally and triage
		// the deletion
		if _, ok := p.instances[instance.ID]; !ok {
			if err := instance.SetDeleted(ctx); err != nil {
				return fmt.Errorf("failed to flag instance %q as deleted: %s", instance.ID, err)
			}

			// If deletion was successful, we can triage the deletion
			if err := instance.TriageDeletion(ctx); err != nil {
				return fmt.Errorf("failed to triage deletion of instance %q: %s", instance.ID, err)
			}

			count++
		}
		return nil
	})
	p.err = err

	if err == nil {
		zeus.Logger(ctx).Info("successfully flagged non-existing instances as deleted",
			zap.Int("count", count),
		)
	} else {
		zeus.Logger(ctx).Error("purging non-existing instances failed, skipping further steps")
	}
}

func (p *purger) purgeDeleted(ctx context.Context) {
	if p.err != nil {
		return
	}

	status := db.InstanceStatusDeleted
	iterator := p.database.ListInstances(
		ctx, status, db.FilterStatusSince(status, 15*time.Minute), db.FilterDeletionTriaged(false),
	)
	count, err := iterator.ForEach(func(instance *db.Instance) error {
		return p.deleteIfRunningAndTriageDeletion(ctx, instance)
	})
	p.err = err

	if err == nil {
		zeus.Logger(ctx).Info("successfully completed purge of deleted instances",
			zap.Int("count", count))
	} else {
		zeus.Logger(ctx).Error("purging deleted instances failed")
	}
}

//-------------------------------------------------------------------------------------------------

func (p *purger) purgeStatus(ctx context.Context, status db.InstanceStatus) {
	if p.err != nil {
		return
	}

	iterator := p.database.ListInstances(ctx, status, db.FilterStatusSince(status, 15*time.Minute))
	count, err := iterator.ForEach(func(instance *db.Instance) error {
		if err := instance.SetDeleted(ctx); err != nil {
			return fmt.Errorf("failed to flag instance %q as deleted: %s", instance.ID, err)
		}
		return p.deleteIfRunningAndTriageDeletion(ctx, instance)
	})
	p.err = err

	if err == nil {
		zeus.Logger(ctx).Info("successfully completed purge of non-started instances",
			zap.Int("count", count),
		)
	} else {
		zeus.Logger(ctx).Error("purging non-started instances failed, skipping further steps")
	}
}

func (p *purger) deleteIfRunningAndTriageDeletion(
	ctx context.Context, dbInstance *db.Instance,
) error {
	// Check if we can find the instance on the cloud provider - if so, we delete it. If
	// deletion was successful or it cannot be found, we triage its deletion.
	if instance, ok := p.instances[dbInstance.ID]; ok {
		if err := p.client.Instances().Delete(ctx, instance.Meta); err != nil {
			return fmt.Errorf(
				"failed to delete instance %q from GCP: %s", instance.Meta.ID, err,
			)
		}
	}

	// At this point, we can set the deletion to be triaged
	if err := dbInstance.TriageDeletion(ctx); err != nil {
		return fmt.Errorf(
			"failed to triage deletion of instance %q: %s", dbInstance.ID, err,
		)
	}

	return nil
}
