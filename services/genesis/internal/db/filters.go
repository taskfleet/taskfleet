package db

import (
	"fmt"
	"time"
)

// Statement implements the Filter interface.
func (s InstanceStatus) Statement() string {
	switch s {
	case InstanceStatusRequested:
		return "deleted_at IS NULL AND booted_at IS NULL"
	case InstanceStatusBooting:
		return "deleted_at IS NULL AND booted_at IS NOT NULL AND started_at IS NULL"
	case InstanceStatusRunning:
		return "deleted_at IS NULL AND started_at IS NOT NULL"
	case InstanceStatusDeleted:
		return "deleted_at IS NOT NULL"
	default:
		panic("invalid instance status")
	}
}

// Values implements the Filter interface.
func (s InstanceStatus) Values() []interface{} {
	return []interface{}{}
}

//-------------------------------------------------------------------------------------------------

type filterSince struct {
	field     string
	timestamp time.Time
}

func (f filterSince) Statement() string {
	return fmt.Sprintf("%s < %%s", f.field)
}

func (f filterSince) Values() []interface{} {
	return []interface{}{f.timestamp}
}

// FilterStatusSince returns a filter that excludes instances that entered the specified status
// less than the specified duration ago.
func FilterStatusSince(status InstanceStatus, duration time.Duration) Filter {
	field := func() string {
		switch status {
		case InstanceStatusRequested:
			return "created_at"
		case InstanceStatusBooting:
			return "booted_at"
		case InstanceStatusRunning:
			return "started_at"
		case InstanceStatusDeleted:
			return "deleted_at"
		default:
			panic("invalid instance status")
		}
	}()
	ts := time.Now().Add(-duration)
	return filterSince{field, ts}
}

//-------------------------------------------------------------------------------------------------

type filterOwner struct {
	owner string
}

func (f filterOwner) Statement() string {
	return "owner = %s"
}

func (f filterOwner) Values() []interface{} {
	return []interface{}{f.owner}
}

// FilterOwner returns a filter that includes all instances from a specified owner.
func FilterOwner(owner string) Filter {
	return filterOwner{owner}
}

//-------------------------------------------------------------------------------------------------

type filterTriaged struct {
	triaged bool
}

func (f filterTriaged) Statement() string {
	return "is_deletion_triaged = %s"
}

func (f filterTriaged) Values() []interface{} {
	return []interface{}{f.triaged}
}

// FilterDeletionTriaged filters instances based on whether their deletion has been acknowledged.
// Whenever this flag is set, it is implied that the `deleted_at` timestamp is set.
func FilterDeletionTriaged(triaged bool) Filter {
	return filterTriaged{triaged}
}
