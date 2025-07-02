package postgre_db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB represents a PostgreSQL database connection pool.
type DB struct {
	connPool *pgxpool.Pool
}

// NewDB creates a new DB instance with a connection pool.
// It takes a context and a connection string as parameters.
// It returns a pointer to a DB instance and an error.
// The error is returned if there was an issue parsing the connection string or creating the connection pool.
func NewDB(ctx context.Context, connString string) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return &DB{connPool: pool}, nil
}
