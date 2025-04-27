package pg

import (
	"context"
	"fmt"
	"sync"

	// Packages
	types "github.com/djthorpe/go-pg/pkg/types"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Listener is an interface for listening to notifications
type Listener interface {
	// Listen to a topic
	Listen(context.Context, string) error

	// Unlisten from a topic
	Unlisten(context.Context, string) error

	// Wait for a notification and return it
	WaitForNotification(context.Context) (*Notification, error)

	// Free resources
	Close(context.Context) error
}

type listener struct {
	sync.Mutex
	pool *pgxpool.Pool
	conn *pgxpool.Conn
}

var _ Listener = (*listener)(nil)

type Notification struct {
	Channel string
	Payload []byte
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewListener return a Listener for the given pool. If pool is nil then
// return nil
func (pg *poolconn) Listener() Listener {
	l := new(listener)
	l.pool = pg.conn.Pool
	return l
}

// Close the connection to the database
func (l *listener) Close(ctx context.Context) error {
	l.Lock()
	defer l.Unlock()

	if l.conn == nil {
		return nil
	}

	// Release below would take care of cleanup and potentially put the
	// connection back into rotation, but in case a Listen was invoked without a
	// subsequent Unlisten on the same topic, close the connection explicitly to
	// guarantee no other caller will receive a partially tainted connection.
	err := l.conn.Conn().Close(ctx)

	// Release the connection
	l.conn.Release()
	l.conn = nil

	// Return any errors
	return err
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Connect to the database, and listen to a topic
func (l *listener) Listen(ctx context.Context, topic string) error {
	l.Lock()
	defer l.Unlock()

	// Acquire a connection
	if l.conn == nil {
		conn, err := l.pool.Acquire(ctx)
		if err != nil {
			return err
		} else {
			l.conn = conn
		}
	}

	// Listen to the topic
	_, err := l.conn.Exec(ctx, "LISTEN "+types.DoubleQuote(topic))
	return err
}

// Unlisten issues an UNLISTEN from the supplied topic.
func (l *listener) Unlisten(ctx context.Context, topic string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the connection is nil
	if l.conn == nil {
		return fmt.Errorf("connection is nil")
	}

	// Unlisten from a topic
	_, err := l.conn.Exec(ctx, "UNLISTEN "+types.DoubleQuote(topic))
	return err
}

// WaitForNotification blocks until receiving a notification and returns it.
func (l *listener) WaitForNotification(ctx context.Context) (*Notification, error) {
	l.Lock()
	defer l.Unlock()

	// Wait for a notification
	if l.conn == nil || l.conn.Conn() == nil {
		return nil, fmt.Errorf("connection is nil")
	}
	n, err := l.conn.Conn().WaitForNotification(ctx)
	if err != nil {
		return nil, err
	}

	// Return the notification
	return &Notification{
		Channel: n.Channel,
		Payload: []byte(n.Payload),
	}, nil
}
