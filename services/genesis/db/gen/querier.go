// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateInstance(ctx context.Context, arg CreateInstanceParams) error
	CreateInstanceType(ctx context.Context, arg CreateInstanceTypeParams) error
	GetInstance(ctx context.Context, id uuid.UUID) (Instance, error)
	GetInstanceType(ctx context.Context, arg GetInstanceTypeParams) (InstanceType, error)
	ListPastBootedInstances(ctx context.Context, minAge pgtype.Interval) ([]Instance, error)
	ListPastDeletedUntriagedInstances(ctx context.Context, minAge pgtype.Interval) ([]Instance, error)
	ListPastRequestedInstances(ctx context.Context, minAge pgtype.Interval) ([]Instance, error)
	ListRunningInstances(ctx context.Context, owner InstanceOwner) ([]Instance, error)
	SetInstanceBooting(ctx context.Context, arg SetInstanceBootingParams) error
	SetInstanceDeleted(ctx context.Context, arg SetInstanceDeletedParams) error
	SetInstanceDeletionTriaged(ctx context.Context, id uuid.UUID) error
	SetInstanceRunning(ctx context.Context, arg SetInstanceRunningParams) error
}

var _ Querier = (*Queries)(nil)
