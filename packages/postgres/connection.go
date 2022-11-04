package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // loads the Postgres driver
)

// Connection represents an initialized client connection to Postgres. It extends `sqlx.DB`.
type Connection struct {
	*sqlx.DB
}

// NewConnection initializes a new connection by using the given connection string. The context is
// used to ping the database to ensure connectivity.
func NewConnection(config ConnectionConfig, options ...ConnectionOption) (*Connection, error) {
	db, err := sqlx.Open("postgres", config.sqlxString())
	if err != nil {
		return nil, fmt.Errorf("failed creating connection: %s", err)
	}

	for _, option := range options {
		option.apply(db)
	}

	return &Connection{db}, nil
}

// MustNewConnection is the same as `NewConnection` but panics on error.
func MustNewConnection(config ConnectionConfig, options ...ConnectionOption) *Connection {
	conn, err := NewConnection(config, options...)
	if err != nil {
		panic(err)
	}
	return conn
}
