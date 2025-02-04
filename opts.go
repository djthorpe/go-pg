package pg

import (
	"fmt"
	"net"
	"net/url"
	"slices"
	"sort"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type opt struct {
	TraceFn
	Verbose bool
	url.Values
	bind *Bind
}

// Opt is a function which applies options for a connection pool
type Opt func(*opt) error

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultPort     = "5432"
	defaultMaxConns = "10"
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Apply options to the opt struct
func apply(opts ...Opt) (*opt, error) {
	var o opt

	// Set defaults
	o.Values = make(url.Values)
	o.Set("host", "localhost")
	o.Set("port", DefaultPort)
	o.Set("pool_max_conns", defaultMaxConns)
	o.bind = NewBind()

	// Apply options
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}

	// Return success
	return &o, nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Set connection pool username and password. If the database name is not set,
// then the username will be used as the default database name.
func WithCredentials(user, password string) Opt {
	return func(o *opt) error {
		if user != "" {
			o.Set("user", user)
		}
		if password != "" {
			o.Set("password", password)
		}
		// TODO: Possible the db name is not being correctly set
		if !o.Has("dbname") && user != "" {
			o.Set("dbname", user)
		}

		// Return success
		return nil
	}
}

// Set the database name for the connection. If the user name is not set,
// then the database name will be used as the user name.
func WithDatabase(name string) Opt {
	return func(o *opt) error {
		if name == "" {
			o.Del("dbname")
		} else {
			o.Set("dbname", name)
		}
		if !o.Has("user") && name != "" {
			o.Set("user", name)
		}
		return nil
	}
}

// Set the address (host) or (host:port) for the connection
func WithAddr(addr string) Opt {
	return func(o *opt) error {
		if !strings.Contains(addr, ":") {
			return WithHostPort(addr, DefaultPort)(o)
		} else if host, port, err := net.SplitHostPort(addr); err != nil {
			return err
		} else {
			return WithHostPort(host, port)(o)
		}
	}
}

// Set the hostname and port for the connection. If the port is not set, then
// the default port 5432 will be used.
func WithHostPort(host, port string) Opt {
	return func(o *opt) error {
		if host != "" {
			o.Set("host", host)
		}
		if port != "" {
			o.Set("port", port)
		}
		return nil
	}
}

// Set the postgresql SSL mode. Valid values are "disable", "allow", "prefer",
// "require",  "verify-ca", "verify-full". See
// https://www.postgresql.org/docs/current/libpq-ssl.html for more information.
func WithSSLMode(mode string) Opt {
	return func(o *opt) error {
		if mode != "" {
			o.Set("sslmode", mode)
		}
		return nil
	}
}

// Set the trace function for the connection pool. The signature of the trace
// function is func(query string, args any, err error) and is called for every
// query executed by the connection pool.
func WithTrace(fn TraceFn) Opt {
	return func(o *opt) error {
		o.TraceFn = fn
		return nil
	}
}

// Set the bind variables for the connection pool
func WithBind(k string, v any) Opt {
	return func(o *opt) error {
		o.bind.Set(k, v)
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (o *opt) encode(skip ...string) []string {
	// We sort the keys to ensure that the URL is deterministic
	keys := make([]string, 0, len(o.Values))
	for key := range o.Values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Encode the values
	var parts []string
	for _, key := range keys {
		if slices.Contains(skip, key) {
			continue
		}
		if value := o.Values.Get(key); value != "" {
			parts = append(parts, fmt.Sprintf("%v=%v", key, o.Values.Get(key)))
		}
	}

	return parts
}

// Encode the options as a connection string
func (o *opt) Encode() string {
	return strings.Join(o.encode(), " ")
}
