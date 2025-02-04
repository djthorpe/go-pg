# go-pg

Postgresql Support for Go. This module provides:

* Binding SQL statements to named arguments;
* Support for mapping go structures to SQL tables, and vice versa;
* Easy semantics for Insert, Delete, Update, Get and List operations;
* Support for tracing and observability.

Documentation: <https://pkg.go.dev/github.com/djthorpe/go-pg>

## Motivation

The package provides a simple way to interact with a Postgresql database from Go, to reduce
the amount of boilerplate code required to interact with the database. The supported operations
align with API calls _POST_, _PUT_, _GET_, _DELETE_ and _PATCH_.

* **Insert** - Insert a row into a table, and return the inserted row (_POST_ or _PUT_);
* **Delete** - Delete one or more rows from a table, and optionally return the deleted rows (_DELETE_);
* **Update** - Update one or more rows in a table, and optionally return the updated rows (_PATCH_);
* **Get** - Get a single row from a table (_GET_);
* **List** - Get a list of rows from a table (_GET_).

In order to support database operations on go types, those types need to implement one or
more of the following interfaces:

### Selector

```go
type Selector interface {
  // Bind row selection variables, returning the SQL statement required for the operation
  // The operation can be Get, Patch, Delete or List
  Select(*Bind, Op) (string, error)
}
```

A type which implements a `Selector` interface can be used to select rows from a table, for get, list, update
and deleting operations.

### Reader

```go
type Reader interface {
  // Scan a row into the receiver
  Scan(Row) error
}
```

A type which implements a `Reader` interface can be used to translate SQL types to the type instance. If multiple
rows are returned, then the `Scan` method is called repeatly until no more rows are returned.

```go
type ListReader interface {
  // Scan a count of returned rows into the receiver
  ScanCount(Row) error
}
```

A type which implements a `ListReader` interface can be used to scan the count of rows returned.

### Writer

```go
type Writer interface {
  // Bind insert variables, returning the SQL statement required for the insert
  Insert(*Bind) (string, error)

  // Bind patch variables
  Patch(*Bind) error
}
```

A type which implements a `Writer` interface can be used to bind object instance variables to SQL parameters.
An example of how to implement an API gateway using this package is shown below.

## Database Server Connection Pool

You can create a connection pool to a database server using the `pg.NewPool` function:

```go
import (
  pg "github.com/djthorpe/go-pg"
)

func main() {
  pool, err := pg.NewPool(ctx,
    pg.WithHostPort(host, port),
    pg.WithCredentials("postgres", "password"),
    pg.WithDatabase(name),
  )
  if err != nil {
      panic(err)
  }
  defer pool.Close()

  // ...
}
```

The options that can be passed to `pg.NewPool` are:

TODO

## Executing Statements and Transactions

* Executing Statements
* Binding Named Arguments
* Replacing Named Arguments
* Executing Transactions

## Implementing Get

TODO

## Implementing List

TODO

## Implementing Insert

TODO

## Implementing Patch

TODO

## Implementing Delete

TODO

## Notify and Listen

TODO

## Schema Support

* Checking if a schema exists
* Creating a schema
* Dropping a schema

## Error Handing and Tracing

TODO

## Testing Support

TODO
