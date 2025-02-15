package pg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	// Packages
	types "github.com/djthorpe/go-pg/pkg/types"
	pgx "github.com/jackc/pgx/v5"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Bind represents a set of variables and arguments to be used in a query.
// The vars are substituted in the query string itself, while the args are
// passed as arguments to the query.
type Bind struct {
	sync.RWMutex
	vars pgx.NamedArgs
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new Bind object with the given name/value pairs
// Returns nil if the number of arguments is not even
func NewBind(pairs ...any) *Bind {
	if len(pairs)%2 != 0 {
		return nil
	}

	// Populate the vars map
	vars := make(pgx.NamedArgs, len(pairs)>>1)
	for i := 0; i < len(pairs); i += 2 {
		if key, ok := pairs[i].(string); !ok || key == "" {
			return nil
		} else {
			vars[key] = pairs[i+1]
		}
	}

	// Return the Bind object
	return &Bind{vars: vars}
}

// Make a copy of the bind object
func (bind *Bind) Copy(pairs ...any) *Bind {
	if len(pairs)%2 != 0 {
		return nil
	}

	varsCopy := make(pgx.NamedArgs, len(bind.vars)+len(pairs)>>1)
	for key, value := range bind.vars {
		varsCopy[key] = value
	}
	for i := 0; i < len(pairs); i += 2 {
		if key, ok := pairs[i].(string); !ok || key == "" {
			return nil
		} else {
			varsCopy[key] = pairs[i+1]
		}
	}

	// Return the copied Bind object
	return &Bind{vars: varsCopy}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (bind *Bind) MarshalJSON() ([]byte, error) {
	return json.Marshal(bind.vars)
}

func (bind *Bind) String() string {
	data, err := json.MarshalIndent(bind.vars, "", "  ")
	if err != nil {
		return err.Error()
	} else {
		return string(data)
	}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Set a bind var, and return the parameter name
func (bind *Bind) Set(key string, value any) string {
	bind.Lock()
	defer bind.Unlock()

	if key == "" {
		return ""
	}
	bind.vars[key] = value
	return "@" + key
}

// Get a bind var
func (bind *Bind) Get(key string) any {
	bind.RLock()
	defer bind.RUnlock()
	return bind.vars[key]
}

// Return true if there is a bind var with the given key
func (bind *Bind) Has(key string) bool {
	bind.RLock()
	defer bind.RUnlock()

	_, ok := bind.vars[key]
	return ok
}

// Delete a bind var
func (bind *Bind) Del(key string) {
	bind.Lock()
	defer bind.Unlock()
	delete(bind.vars, key)
}

// Join a bind var with a separator, when it is a
// []any and return the result as a string. Returns
// an empty string if the key does not exist.
func (bind *Bind) Join(key, sep string) string {
	bind.RLock()
	defer bind.RUnlock()

	if value, ok := bind.vars[key]; !ok {
		return ""
	} else if v, ok := value.([]any); ok {
		str := make([]string, len(v))
		for i, value := range v {
			str[i] = fmt.Sprint(value)
		}
		return strings.Join(str, sep)
	} else {
		return fmt.Sprint(value)
	}
}

// Append a bind var to a list. Returns false if the key
// is not a list, or the value is not a list.
func (bind *Bind) Append(key string, value any) bool {
	bind.RLock()
	defer bind.RUnlock()

	// Create a new list if it doesn't exist
	if _, ok := bind.vars[key]; !ok {
		bind.vars[key] = make([]any, 0, 5)
	}

	// Check type
	if _, ok := bind.vars[key].([]any); !ok {
		return false
	}

	// Append value
	bind.vars[key] = append(bind.vars[key].([]any), value)

	// Return success
	return true
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - QUERY

// Query a row and return the result
func (bind *Bind) QueryRow(ctx context.Context, conn pgx.Tx, query string) pgx.Row {
	bind.RLock()
	defer bind.RUnlock()
	return conn.QueryRow(ctx, bind.Replace(query), bind.vars)
}

// Query a set of rows and return the result
func (bind *Bind) Query(ctx context.Context, conn pgx.Tx, query string) (pgx.Rows, error) {
	bind.RLock()
	defer bind.RUnlock()
	return conn.Query(ctx, bind.Replace(query), bind.vars)
}

// Execute a query
func (bind *Bind) Exec(ctx context.Context, conn pgx.Tx, query string) error {
	bind.RLock()
	defer bind.RUnlock()
	_, err := conn.Exec(ctx, bind.Replace(query), bind.vars)
	return err
}

// Queue a query
func (bind *Bind) queuerow(batch *pgx.Batch, query string, reader Reader) {
	bind.RLock()
	defer bind.RUnlock()
	queuedquery := batch.Queue(bind.Replace(query), bind.vars)
	queuedquery.QueryRow(func(row pgx.Row) error {
		return reader.Scan(row)
	})
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Return a query string with ${subtitution} replaced by the values:
//   - ${key} => value
//   - ${'key'} => 'value'
//   - ${"key"} => "value"
//   - $1 => $1
//   - $$ => $$
func (bind *Bind) Replace(query string) string {
	fetch := func(key string) string {
		return fmt.Sprint(bind.vars[key])
	}
	return os.Expand(query, func(key string) string {
		if key == "$" { // $$ => $$
			return "$$"
		}
		if types.IsNumeric(key) {
			return "$" + key // $1 => $1
		}
		if types.IsSingleQuoted(key) { // ${'key'} => 'value'
			// Special case where value is []string and single quote for IN (${key})
			key := strings.Trim(key, "'")
			value := bind.vars[key]
			switch v := value.(type) {
			case []string:
				result := make([]string, len(v))
				for i, s := range v {
					result[i] = types.Quote(s)
				}
				return strings.Join(result, ",")
			default:
				return types.Quote(fetch(key))
			}
		}
		if types.IsDoubleQuoted(key) { // ${"key"} => "value"
			return types.DoubleQuote(fetch(strings.Trim(key, "\"")))
		}
		return fetch(key) // ${key} => value
	})
}
