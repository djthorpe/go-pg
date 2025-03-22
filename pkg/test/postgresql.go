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

// Create a new postgresql database and pool, using a unique container name
func NewPgxContainer(ctx context.Context, name string, verbose bool, tracer pg.TraceFn) (*Container, pg.PoolConn, error) {
	// Create a new container with postgresql package
	container, err := NewContainer(ctx, name, pgxContainer,
		OptEnv("POSTGRES_REPLICATION_PASSWORD", "password"),
		OptPostgres("postgres", "password", name), // User, Password, Database
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
