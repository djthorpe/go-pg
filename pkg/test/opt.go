package test

import (
	"fmt"

	// Packages
	nat "github.com/docker/go-connections/nat"
	testcontainers "github.com/testcontainers/testcontainers-go"
	wait "github.com/testcontainers/testcontainers-go/wait"

	// Anonymous imports
	_ "github.com/jackc/pgx/v5/stdlib"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Opt func(*opts) error

type opts struct {
	req testcontainers.ContainerRequest
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func OptCommand(cmd []string) Opt {
	return func(o *opts) error {
		o.req.Cmd = cmd
		return nil
	}
}

func OptEnv(name, value string) Opt {
	return func(o *opts) error {
		if o.req.Env == nil {
			o.req.Env = make(map[string]string)
		}
		o.req.Env[name] = value
		return nil
	}
}

func OptPorts(ports ...string) Opt {
	return func(o *opts) error {
		o.req.ExposedPorts = ports
		o.appendWaitStrategy(wait.ForExposedPort())
		return nil
	}
}

func OptPostgres(user, password, database string) Opt {
	return func(o *opts) error {
		if err := OptEnv("POSTGRES_USER", user)(o); err != nil {
			return err
		}
		if err := OptEnv("POSTGRES_PASSWORD", password)(o); err != nil {
			return err
		}
		if err := OptEnv("POSTGRES_DB", database)(o); err != nil {
			return err
		}
		if err := OptPorts(pgxPort)(o); err != nil {
			return err
		}
		o.appendWaitStrategy(wait.ForSQL(nat.Port(pgxPort), "pgx", func(host string, port nat.Port) string {
			return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port.Port(), database)
		}))
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (o *opts) appendWaitStrategy(strategy wait.Strategy) {
	if o.req.WaitingFor == nil {
		o.req.WaitingFor = strategy
	} else if multi, ok := o.req.WaitingFor.(*wait.MultiStrategy); ok {
		multi.Strategies = append(multi.Strategies, strategy)
	} else {
		o.req.WaitingFor = wait.ForAll(o.req.WaitingFor, strategy)
	}
}
