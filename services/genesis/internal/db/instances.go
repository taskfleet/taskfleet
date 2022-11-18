package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.taskfleet.io/packages/postgres"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

var (
	errNotRetrievedFromDatabase = errors.New("instance has not been retrieved from database")
)

// Instance represents an instance stored in the database.
type Instance struct {
	db *postgres.Connection

	ID          uuid.UUID              `db:"id"`
	Provider    typedefs.CloudProvider `db:"provider"`
	AccountName string                 `db:"account_name"`
	Zone        string                 `db:"zone"`
	Owner       string                 `db:"owner"`

	InstanceType      string `db:"instance_type"`
	IsSpot            bool   `db:"is_spot"`
	CPUCountRequested int32  `db:"cpu_count_requested"`
	MemoryMBRequested int32  `db:"memory_mb_requested"`
	MemoryMBReserved  int32  `db:"memory_mb_reserved"`
	BootImage         string `db:"boot_image"`

	Hostname *string `db:"hostname"`

	CreatedAt         time.Time  `db:"created_at"`
	BootedAt          *time.Time `db:"booted_at"`
	StartedAt         *time.Time `db:"started_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
	IsDeletionTriaged bool       `db:"is_deletion_triaged"`
}

// InstanceType represents an instance type stored in the database.
type InstanceType struct {
	db *postgres.Connection

	Provider        typedefs.CloudProvider   `db:"provider"`
	Name            string                   `db:"name"`
	CPUCount        int32                    `db:"cpu_count"`
	CPUArchitecture typedefs.CPUArchitecture `db:"cpu_architecture"`
	MemoryMiB       int32                    `db:"memory_mib"`
	GPUKind         *typedefs.GPUKind        `db:"gpu_kind"`
	GPUCount        *int32                   `db:"gpu_count"`
}

//-------------------------------------------------------------------------------------------------

// SetBooting sets a timestamp indicating that the instance booted at the current point in time.
// Additional information that is only available upon boot is passed here. An error is returned if
// the instance this method is called on has not been retrieved from the database.
func (i *Instance) SetBooting(
	ctx context.Context, cpuKind typedefs.CPUKind, hostname string,
) error {
	t, err := i.setNow(ctx, "booted_at", map[string]interface{}{
		"cpu_kind": cpuKind,
		"hostname": hostname,
	})
	if err != nil {
		return err
	}
	i.BootedAt = t
	return nil
}

// SetRunning sets a timestamp indicating that the instance became responsive (i.e. its primary
// process started up) at the current point in time. An error is returned if the instance this
// method is called on has not been retrieved from the database.
func (i *Instance) SetRunning(ctx context.Context) error {
	t, err := i.setNow(ctx, "started_at", map[string]interface{}{})
	if err != nil {
		return err
	}
	i.StartedAt = t
	return nil
}

// SetDeleted flags the instance as deleted at the current point in time. An error is returned if
// the instance this method is called on has not been retrieved from the database.
func (i *Instance) SetDeleted(ctx context.Context) error {
	t, err := i.setNow(ctx, "deleted_at", map[string]interface{}{})
	if err != nil {
		return err
	}
	i.DeletedAt = t
	return nil
}

// TriageDeletion sets a flag which indicates that the instance has been successfully deleted.
func (i *Instance) TriageDeletion(ctx context.Context) error {
	_, err := i.db.ExecContext(
		ctx, "UPDATE instances SET is_deletion_triaged = TRUE WHERE id = $1", i.ID,
	)
	return err
}

func (i *Instance) setNow(
	ctx context.Context, field string, additionalUpdates map[string]interface{},
) (*time.Time, error) {
	if i.db == nil {
		return nil, errNotRetrievedFromDatabase
	}

	now := time.Now()

	// Collect all update statements and values
	updateStatements := []string{fmt.Sprintf("%s = $2", field)}
	values := []interface{}{i.ID, now}
	j := 3
	for column, value := range additionalUpdates {
		updateStatements = append(updateStatements, fmt.Sprintf("%s = $%d", column, j))
		values = append(values, value)
		j++
	}

	// Then execute the query
	query := fmt.Sprintf(
		"UPDATE instances SET %s WHERE id = $1", strings.Join(updateStatements, ", "),
	)
	if _, err := i.db.ExecContext(ctx, query, values...); err != nil {
		return nil, err
	}

	// And set time
	return &now, nil
}

//-------------------------------------------------------------------------------------------------

func (i Instance) validateInput() error {
	if i.Provider == "" {
		return ValidationError{Message: "instance must be associated with cloud provider"}
	}
	if i.Zone == "" {
		return ValidationError{Message: "zone must not be empty"}
	}
	if i.AccountName == "" {
		return ValidationError{Message: "account name must not be empty"}
	}
	if i.Owner == "" {
		return ValidationError{Message: "owner must not be empty"}
	}
	if i.MachineType == "" {
		return ValidationError{Message: "machine type must not be empty"}
	}
	if i.CPUCount < 1 {
		return ValidationError{Message: "cpu count must be at least 1"}
	}
	if i.MemoryMB < 1024 {
		return ValidationError{Message: "memory must be at least 1024 MB"}
	}
	if (i.GPUKind == nil) != (i.GPUCount == nil) {
		return ValidationError{Message: "gpu configuration is incosistent"}
	}
	if i.GPUKind != nil {
		if *i.GPUCount < 1 {
			return ValidationError{Message: "gpu count must be at least 1 if set"}
		}
	}
	if i.BootImage == "" {
		return ValidationError{Message: "boot image must not be empty"}
	}
	if i.BootDiskSizeGiB < 10 {
		return ValidationError{Message: "boot disk size must be at least 10 GiB"}
	}
	if i.DiskSizeHDDGiB < 0 {
		return ValidationError{Message: "hdd disk size must not be negative"}
	}
	if i.DiskSizeSSDStandardGiB < 0 {
		return ValidationError{Message: "standard ssd disk size must not be negative"}
	}
	if i.DiskSizeSSDHPGiB < 0 {
		return ValidationError{Message: "high-performance ssd disk size must not be negative"}
	}
	if i.BootedAt != nil || i.StartedAt != nil || i.DeletedAt != nil {
		return ValidationError{Message: "timestamps are output-only"}
	}
	return nil
}
