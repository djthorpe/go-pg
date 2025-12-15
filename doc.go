// Package pg provides PostgreSQL support for Go, built on top of pgx.
//
// It provides binding SQL statements to named arguments, mapping Go structures
// to SQL tables, and easy semantics for Insert, Delete, Update, Get and List
// operations. The package also supports bulk insert operations, transactions,
// and tracing for observability.
//
// # Connection Pool
//
// Create a connection pool using NewPool:
//
//	pool, err := pg.NewPool(ctx,
//	    pg.WithURL("postgres://user:pass@localhost:5432/dbname"),
//	)
//	if err != nil {
//	    panic(err)
//	}
//	defer pool.Close()
//
// # Executing Queries
//
// Use bind variables with the With method:
//
//	err := pool.With("table", "users").Exec(ctx, `SELECT * FROM ${"table"}`)
//
// # CRUD Operations
//
// Implement the Reader, Writer, and Selector interfaces on your types to
// enable Insert, Update, Delete, Get, and List operations:
//
//	err := pool.Insert(ctx, &obj, obj)    // Insert
//	err := pool.Get(ctx, &obj, selector)  // Get
//	err := pool.List(ctx, &list, request) // List
//	err := pool.Update(ctx, &obj, selector, writer) // Update
//	err := pool.Delete(ctx, &obj, selector) // Delete
package pg
