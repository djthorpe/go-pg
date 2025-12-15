package test

import (
	"context"
	"errors"

	// Packages
	pg "github.com/djthorpe/go-pg"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	pgxContainer = "ghcr.io/mutablelogic/docker-postgres:17-bookworm"
	//pgxContainer = "postgis/postgis:16-master" // Postgresql container
	pgxPort = "5432/tcp"
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewPgxContainer creates a new PostgreSQL container and connection pool.
func NewPgxContainer(ctx context.Context, name string, verbose bool, tracer pg.TraceFn) (*Container, pg.PoolConn, error) {
	// Create a new container with postgresql package
	container, err := NewContainer(ctx, name, pgxContainer,
		OptEnv("POSTGRES_REPLICATION_PASSWORD", "password"),
		OptPostgres("postgres", "password", name),                            // User, Password, Database
		OptPostgresSetting("shared_preload_libraries", "pg_stat_statements"), // Enable pg_stat_statements
		OptPostgresSetting("wal_level", "logical"),                           // Enable logical replication
	)
	if err != nil {
		return nil, nil, err
	}

	host, _ := container.GetEnv("POSTGRES_HOST")
	port, err := container.GetPort(pgxPort)
	if err != nil {
		return nil, nil, err
	}

	// Create a connection pool
	pool, err := pg.NewPool(ctx,
		pg.WithCredentials("postgres", "password"),
		pg.WithDatabase(name),
		pg.WithHostPort(host, port),
		pg.WithTrace(tracer),
	)
	if err != nil {
		return nil, nil, errors.Join(err, container.Close(ctx))
	} else if err := pool.Ping(ctx); err != nil {
		return nil, nil, errors.Join(err, container.Close(ctx))
	}

	// Return success
	return container, pool, nil
}
