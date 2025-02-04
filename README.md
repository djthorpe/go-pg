# go-pg

Postgresql Support for Go. This module provides:

* Binding SQL statements to named arguments;
* Support for mapping go structures to SQL tables, and vice versa;
* Easy semantics for Insert, Delete, Update, Get and List operations;
* Support for tracing and observability.

Documentation: <https://pkg.go.dev/github.com/djthorpe/go-pg>

## Motivation

The package provides a simple way to interact with a Postgresql database from Go. 

* **Insert** - Insert a row into a table, and return the inserted row;
* **Delete** - Delete one or more rows from a table, and optionally return the deleted rows;
* **Update** - Update one or more rows in a table, and optionally return the updated rows;
* **Get** - Get a single row from a table;
* **List** - Get a list of rows from a table.

In order to support database operations on go types, those types need to implement one or
more of the following interfaces:

### Selector

```go
type Selector interface {
  // Bind row selection variables, returning the SQL statement required for the operation
  // The operation can be Get, Insert, Patch, Delete or List
  Select(*Bind, Op) (string, error)
}
```

### Reader

```go
type Reader interface {
  // Scan a row into the receiver
  Scan(Row) error
}
```

```go
type ListReader interface {
  // Scan a count of returned rows into the receiver
  ScanCount(Row) error
}
```

### Writer

```go
type Writer interface {
  // Bind insert variables, returning the SQL statement required for the insert
  Insert(*Bind) (string, error)

  // Bind patch variables
  Patch(*Bind) error
}
```
