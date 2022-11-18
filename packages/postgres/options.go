package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// ConnectionOption is an interface that allows customization of created database connections.
type ConnectionOption interface {
	apply(db *sqlx.DB)
}

//-------------------------------------------------------------------------------------------------

// WithConnectionPoolSize limits the number of concurrent database connections to the specified
// number.
func WithConnectionPoolSize(size int) ConnectionOption {
	return poolSizeOption{size}
}

type poolSizeOption struct {
	size int
}

func (o poolSizeOption) apply(db *sqlx.DB) {
	db.SetMaxOpenConns(o.size)
}

//-------------------------------------------------------------------------------------------------

// WithIdlePoolSize limits the number of idle database connections to the specified number.
func WithIdlePoolSize(size int) ConnectionOption {
	return idlePoolSizeOption{size}
}

type idlePoolSizeOption struct {
	size int
}

func (o idlePoolSizeOption) apply(db *sqlx.DB) {
	db.SetMaxIdleConns(o.size)
}

//-------------------------------------------------------------------------------------------------

// WithConnectionTimeout limits the total age of database connections to the given duration.
func WithConnectionTimeout(timeout time.Duration) ConnectionOption {
	return connectionTimeoutOption{timeout}
}

type connectionTimeoutOption struct {
	timeout time.Duration
}

func (o connectionTimeoutOption) apply(db *sqlx.DB) {
	db.SetConnMaxLifetime(o.timeout)
}

//-------------------------------------------------------------------------------------------------

// WithIdleTimeout limits the duration for which idle connections may exist.
func WithIdleTimeout(timeout time.Duration) ConnectionOption {
	return idleTimeoutOption{timeout}
}

type idleTimeoutOption struct {
	timeout time.Duration
}

func (o idleTimeoutOption) apply(db *sqlx.DB) {
	db.SetConnMaxIdleTime(o.timeout)
}
