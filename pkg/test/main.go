package test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	// Packages
	pg "github.com/djthorpe/go-pg"
)

/////////////////////////////////////////////////////////////////////
// TYPES

// Conn is a wrapper around pg.PoolConn which provides a test connection
type Conn struct {
	pg.PoolConn
	t *testing.T
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	timeout = 2 * time.Minute
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func Main(m *testing.M, conn *Conn) {
	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Name of executable
	name, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Start the container
	verbose := slices.Contains(os.Args, "-test.v=true")
	container, pool, err := NewPgxContainer(ctx, filepath.Base(name), verbose, func(ctx context.Context, sql string, args any, err error) {
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
		if verbose || err != nil {
			if args == nil {
				log.Printf("SQL: %v", sql)
			} else {
				log.Printf("SQL: %v, ARGS: %v", sql, args)
			}
		}
	})
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	defer container.Close(ctx)

	// Set the connection
	*conn = Conn{pool, nil}

	// Run tests
	os.Exit(m.Run())
}

// Begin a test
func (c *Conn) Begin(t *testing.T) *Conn {
	t.Log("Begin", t.Name())
	return &Conn{c.PoolConn, t}
}

// End a test
func (c *Conn) Close() {
	if c.t != nil {
		c.t.Log("End", c.t.Name())
	}
}
