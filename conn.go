package pg

import (
	"context"
	"errors"

	// Packages
	pgx "github.com/jackc/pgx/v5"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Conn interface {
	// Return a new connection with bound parameters
	With(...any) Conn

	// Perform a transaction within a function
	Tx(context.Context, func(conn Conn) error) error

	// Execute a query
	Exec(context.Context, string) error

	// Perform an insert
	Insert(context.Context, Reader, Writer) error

	// Perform a patch
	Patch(context.Context, Reader, Selector, Writer) error

	// Perform a delete
	Delete(context.Context, Reader, Selector) error

	// Perform a get
	Get(context.Context, Reader, Selector) error

	// Perform a list. If the reader is a ListReader, then the
	// count of items is also calculated
	List(context.Context, Reader, Selector) error
}

// Operation type
type Op uint

// Row scanner
type Row pgx.Row

// Bind a row to an object
type Reader interface {
	// Scan row into a result
	Scan(Row) error
}

// Bind a row to an object, and also count the number of rows
type ListReader interface {
	Reader

	// Scan count into the result
	ScanCount(Row) error
}

// Bind an object to bind parameters for inserting or patching
type Writer interface {
	// Set bind parameters for an insert
	Insert(*Bind) (string, error)

	// Set bind parameters for a patch (update)
	Patch(*Bind) error
}

// Bind selection parameters for getting, patching or deleting
type Selector interface {
	// Set bind parameters for getting, patching or deleting
	Select(*Bind, Op) (string, error)
}

// Concrete connection
type conn struct {
	conn pgx.Tx
	bind *Bind
}

// Ensure interfaces are satisfied
var _ Conn = (*conn)(nil)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

// Operations
const (
	None Op = iota
	Get
	Insert
	Patch
	Delete
	List
)

func (o Op) String() string {
	switch o {
	case Get:
		return "GET"
	case Insert:
		return "INSERT"
	case Patch:
		return "PATCH"
	case Delete:
		return "DELETE"
	case List:
		return "LIST"
	}
	return "UNKNOWN"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - CONN

// Return a new connection with new bound parameters
func (p *conn) With(params ...any) Conn {
	return &conn{p.conn, p.bind.Copy(params...)}
}

// Perform a transaction, then commit or rollback
func (p *conn) Tx(ctx context.Context, fn func(conn Conn) error) error {
	return tx(ctx, p.conn, p.bind, fn)
}

// Execute a query
func (p *conn) Exec(ctx context.Context, query string) error {
	return p.bind.Exec(ctx, p.conn, query)
}

// Perform an insert, binding parameters from
// the writer, and scanning the result into the reader
func (p *conn) Insert(ctx context.Context, reader Reader, writer Writer) error {
	return insert(ctx, p.conn, p.bind, reader, writer)
}

// Perform a patch, selecting using the selector, binding parameters from
// the writer, and scanning the result into the reader
func (p *conn) Patch(ctx context.Context, reader Reader, sel Selector, writer Writer) error {
	return patch(ctx, p.conn, p.bind, reader, sel, writer)
}

// Perform a delete, binding parameters with the selector and scanning the
// deleted data into the reader
func (p *conn) Delete(ctx context.Context, reader Reader, sel Selector) error {
	return del(ctx, p.conn, p.bind, reader, sel)
}

// Perform a get, binding parameters with the selector and scanning a single
// row into the reader
func (p *conn) Get(ctx context.Context, reader Reader, sel Selector) error {
	return get(ctx, p.conn, p.bind, reader, sel)
}

// Perform a list, binding parameters with the selector and scanning rows
// into the reader
func (p *conn) List(ctx context.Context, reader Reader, sel Selector) error {
	return list(ctx, p.conn, p.bind, reader, sel)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func tx(ctx context.Context, tx pgx.Tx, bind *Bind, fn func(conn Conn) error) error {
	tx, err := tx.Begin(ctx)
	if err != nil {
		return err
	}

	tx_ := &conn{tx, bind.Copy()}
	if err := fn(tx_); err != nil {
		return errors.Join(notfound(err), tx.Rollback(ctx))
	} else {
		return errors.Join(notfound(err), tx.Commit(ctx))
	}
}

func insert(ctx context.Context, conn pgx.Tx, bind *Bind, reader Reader, writer Writer) error {
	if query, err := writer.Insert(bind); err != nil {
		return err
	} else {
		return notfound(reader.Scan(bind.QueryRow(ctx, conn, query)))
	}
}

func patch(ctx context.Context, conn pgx.Tx, bind *Bind, reader Reader, sel Selector, writer Writer) error {
	if query, err := sel.Select(bind, Patch); err != nil {
		return err
	} else if err := writer.Patch(bind); err != nil {
		return err
	} else {
		return notfound(reader.Scan(bind.QueryRow(ctx, conn, query)))
	}
}

func del(ctx context.Context, conn pgx.Tx, bind *Bind, reader Reader, sel Selector) error {
	if query, err := sel.Select(bind, Delete); err != nil {
		return err
	} else {
		return notfound(reader.Scan(bind.QueryRow(ctx, conn, query)))
	}
}

func get(ctx context.Context, conn pgx.Tx, bind *Bind, reader Reader, sel Selector) error {
	if query, err := sel.Select(bind, Get); err != nil {
		return err
	} else {
		return notfound(reader.Scan(bind.QueryRow(ctx, conn, query)))
	}
}

func list(ctx context.Context, conn pgx.Tx, bind *Bind, reader Reader, sel Selector) error {
	// Set groupby, orderby and offsetlimit
	bind.Set("groupby", "")
	bind.Set("orderby", "")
	bind.Set("offsetlimit", "")

	// Bind
	query, err := sel.Select(bind, List)
	if err != nil {
		return notfound(err)
	}

	// Count the number of rows if the reader is a ListReader
	if counter, ok := reader.(ListReader); ok {
		if err := count(ctx, conn, query, bind, counter); err != nil {
			return notfound(err)
		}
	}

	// Execute
	rows, err := bind.Query(ctx, conn, query+` ${groupby} ${orderby} ${offsetlimit}`)
	if err != nil {
		return notfound(err)
	}
	defer rows.Close()

	// Read rows
	for rows.Next() {
		if err := reader.Scan(rows); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// Return success
	return nil
}

func count(ctx context.Context, conn pgx.Tx, query string, bind *Bind, reader ListReader) error {
	// Make a subquery
	return reader.ScanCount(bind.QueryRow(ctx, conn, `WITH sq AS (`+query+` ${groupby}) SELECT COUNT(*) AS "count" FROM sq`))
}

func notfound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	} else {
		return err
	}
}
