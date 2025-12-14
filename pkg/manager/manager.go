package manager

import (
	"context"

	// Packages
	pg "github.com/djthorpe/go-pg"
	schema "github.com/djthorpe/go-pg/pkg/manager/schema"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	conn pg.PoolConn
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new database manager
func New(ctx context.Context, conn pg.PoolConn) (*Manager, error) {
	self := new(Manager)
	self.conn = conn.With("schema", schema.CatalogSchema).(pg.PoolConn)

	// Bootstrap dblink
	//if err := schema.Bootstrap(ctx, self.conn); err != nil {
	//	return nil, err
	//}

	// Return success
	return self, nil
}
