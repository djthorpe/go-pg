package pg

import (
	"context"
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// SQL

const (
	// Check for existence of schema
	schemaExists = `	
		SELECT EXISTS (
  			SELECT 1 FROM pg_catalog.pg_namespace WHERE	nspname = ${'schema'}
		)
	`

	// Create schema
	schemaCreate = `
		CREATE SCHEMA IF NOT EXISTS ${"schema"}
	`

	// Drop schema
	schemaDrop = `
		DROP SCHEMA IF EXISTS ${"schema"} CASCADE
	`
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type exists struct {
	Exists bool
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func SchemaExists(ctx context.Context, conn Conn, name string) (bool, error) {
	var exists exists
	if err := conn.With("schema", name).Get(ctx, &exists, exists); err != nil {
		return false, err
	} else {
		return exists.Exists, nil
	}
}

func SchemaCreate(ctx context.Context, conn Conn, name string) error {
	return conn.With("schema", name).Exec(ctx, schemaCreate)
}

func SchemaDrop(ctx context.Context, conn Conn, name string) error {
	return conn.With("schema", name).Exec(ctx, schemaDrop)
}

////////////////////////////////////////////////////////////////////////////////
// SCAN

func (r exists) Select(bind *Bind, op Op) (string, error) {
	switch op {
	case Get:
		return schemaExists, nil
	default:
		return "", fmt.Errorf("invalid operation %q", op)
	}
}

func (r *exists) Scan(row Row) error {
	return row.Scan(&r.Exists)
}
