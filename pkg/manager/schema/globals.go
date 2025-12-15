package schema

import (
	"context"

	// Packages
	pg "github.com/djthorpe/go-pg"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	CatalogSchema  = "pg_catalog"
	APIPrefix      = "/pg/v1"
	DefaultAclRole = "PUBLIC"
)

const (
	// Maximum number of items to return in a list query, for each type
	RoleListLimit       = 100
	DatabaseListLimit   = 100
	SchemaListLimit     = 100
	ObjectListLimit     = 100
	ConnectionListLimit = 100
	TablespaceListLimit = 100
	ExtensionListLimit  = 100
	SettingListLimit    = 500
)

const (
	pgTimestampFormat    = "2006-01-02 15:04:05"
	pgObfuscatedPassword = "********"
	defaultSchema        = "public"
	reservedPrefix       = "pg_"
)

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

// Bootstrap creates the dblink extension required for remote database queries.
// This should be called once when initializing the manager.
func Bootstrap(ctx context.Context, conn pg.PoolConn) error {
	return conn.Exec(ctx, dblinkCreateExtension)
}

const (
	dblinkCreateExtension = `CREATE EXTENSION IF NOT EXISTS dblink WITH SCHEMA ` + defaultSchema
)
