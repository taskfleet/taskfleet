package db

import (
	"context"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"       // file source
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/postgres"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()

	// First, run migration to obtain a clean slate
	var dbConfig postgres.ConnectionConfig
	envconfig.MustProcess("DB", &dbConfig)

	migration, err := migrate.New("file://../../migrations", dbConfig.String())
	require.Nil(t, err)
	if err := migration.Down(); err != nil {
		if err != migrate.ErrNoChange {
			t.Fatalf("failed to apply down migrations: %s", err)
		}
	}
	require.Nil(t, migration.Up())

	// Then, create some instances
	c, err := postgres.NewConnection(dbConfig)
	require.Nil(t, err)
	conn := NewConnection(c)

	instance := Instance{
		Provider:        typedefs.ProviderGoogleCloudPlatform,
		AccountName:     "my-project",
		Zone:            "my-zone",
		Owner:           "my-owner",
		MachineType:     "n1-standard-1",
		IsSpot:          false,
		CPUCount:        1,
		MemoryMB:        3840,
		BootImage:       "projects/debian-cloud/global/images/family/debian-10",
		BootDiskSizeGiB: 10,
	}

	// Initialize some without GPU
	instances := make([]*Instance, 10)
	for i := 0; i < 5; i++ {
		instances[i] = newInstance(ctx, t, conn, instance)
	}

	// Initialize some with GPU
	gpuKind := typedefs.GPUNvidiaTeslaK80
	var gpuCount int32 = 1
	instance.GPUKind = &gpuKind
	instance.GPUCount = &gpuCount
	for i := 5; i < 10; i++ {
		instances[i] = newInstance(ctx, t, conn, instance)
	}

	// Change the state for some of them
	// Delete instance #1 without booting
	if err := instances[0].SetDeleted(ctx); err != nil {
		t.Fatalf("failed to set instance deleted: %s", err)
	}

	// Boot instances 4-10 (incl)
	for i := 3; i < 10; i++ {
		if err := instances[i].SetBooting(
			ctx, typedefs.CPUIntelXeonBroadwell, "test.dns.internal",
		); err != nil {
			t.Fatalf("failed to set instance booted: %s", err)
		}
	}

	// Start instances 7-10 (incl)
	for i := 6; i < 10; i++ {
		if err := instances[i].SetRunning(ctx); err != nil {
			t.Fatalf("failed to set instance started: %s", err)
		}
	}

	// List instances and check number of instances makes sense
	it := conn.ListInstances(ctx, InstanceStatusRequested)
	requestedInstances, err := it.Collect()
	require.Nil(t, err)
	assert.Len(t, requestedInstances, 2)

	it = conn.ListInstances(ctx, InstanceStatusBooting)
	bootingInstances, err := it.Collect()
	require.Nil(t, err)
	assert.Len(t, bootingInstances, 3)

	it = conn.ListInstances(ctx, InstanceStatusRunning)
	runningInstances, err := it.Collect()
	require.Nil(t, err)
	assert.Len(t, runningInstances, 4)

	it = conn.ListInstances(ctx, InstanceStatusDeleted)
	deletedInstances, err := it.Collect()
	require.Nil(t, err)
	assert.Len(t, deletedInstances, 1)

	// Run migrations to reset the database to a clean slate
	require.Nil(t, migration.Down())
	require.Nil(t, migration.Up())
}

//-------------------------------------------------------------------------------------------------

func newInstance(
	ctx context.Context, t *testing.T, conn Connection, instance Instance,
) *Instance {
	uuid, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("failed creating UUID: %s", err)
	}
	instance.ID = uuid

	result, err := conn.CreateInstance(ctx, instance)
	if err != nil {
		t.Fatalf("failed creating instance: %s", err)
	}
	return result
}
