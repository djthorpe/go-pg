package pg

import (
	"context"
	"strings"

	// Packages
	pgx "github.com/jackc/pgx/v5"
)

//////////////////////////////////////////////////////////////////////////////
// TYPES

// Tracer is a postgresql query tracer
type tracer struct {
	TraceFn
	SQL  string
	Args []any
}

// TraceFn is a function which is called when a query is executed,
// with the execution context, the SQL and arguments, and the error
// if any was generated
type TraceFn func(context.Context, string, any, error)

//////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new query tracer
func NewTracer(fn TraceFn) *tracer {
	if fn == nil {
		return nil
	}
	return &tracer{
		TraceFn: fn,
	}
}

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (tracer *tracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	tracer.SQL = data.SQL
	tracer.Args = data.Args
	return ctx
}

func (tracer *tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	tracer.TraceFn(ctx, strings.TrimSpace(tracer.SQL), args(tracer.Args), data.Err)
}

//////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func args(args []any) any {
	if len(args) == 0 {
		return nil
	}
	if len(args) == 1 {
		return args[0]
	}
	return args
}
