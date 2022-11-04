package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.taskfleet.io/packages/postgres"
)

type conn struct {
	db *postgres.Connection
}

// NewConnection opens a new connection to the database with the specified connection string.
func NewConnection(client *postgres.Connection) Connection {
	return &conn{client}
}

//-------------------------------------------------------------------------------------------------

func (c *conn) CreateInstance(ctx context.Context, instance Instance) (*Instance, error) {
	// Validate
	if err := instance.validateInput(); err != nil {
		return nil, err
	}

	// Create instance
	instance.CreatedAt = time.Now()

	// Execute database statement
	if err := c.db.Insert(ctx, "instances", instance); err != nil {
		return nil, err
	}
	instance.db = c.db
	return &instance, nil
}

func (c *conn) GetInstance(ctx context.Context, id uuid.UUID) (*Instance, error) {
	query := "SELECT * FROM instances WHERE id = $1 AND deleted_at IS NULL"
	row := c.db.QueryRowxContext(ctx, query, id)
	var instance Instance
	if err := row.StructScan(&instance); err != nil {
		if err == sql.ErrNoRows {
			return &instance, ErrNotExist
		}
		return nil, err
	}
	instance.db = c.db
	return &instance, nil
}

func (c *conn) ListInstances(ctx context.Context, filters ...Filter) InstanceIterator {
	// First, build the correct query
	query := "SELECT * FROM instances"
	idx := 1
	values := []interface{}{}
	if len(filters) > 0 {
		conditions := make([]string, len(filters))
		for i, filter := range filters {
			values = append(values, filter.Values()...)
			placeholders := make([]interface{}, len(filter.Values()))
			for i := 0; i < len(filter.Values()); i++ {
				placeholders[i] = fmt.Sprintf("$%d", idx)
				idx++
			}
			conditions[i] = fmt.Sprintf(filter.Statement(), placeholders...)
		}
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Then execute it
	rows, err := c.db.QueryxContext(ctx, query, values...)
	return InstanceIterator{db: c.db, rows: rows, err: err}
}
