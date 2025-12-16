package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	// Packages
	kong "github.com/alecthomas/kong"
	client "github.com/mutablelogic/go-client"
	httpclient "github.com/mutablelogic/go-pg/pkg/manager/httpclient"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Globals struct {
	// Debug option
	Debug bool `name:"debug" help:"Enable debug logging"`

	// HTTP server options
	HTTP struct {
		Prefix string `name:"prefix" help:"HTTP path prefix" default:"/api/v1"`
		Addr   string `name:"addr" env:"PG_ADDR" help:"HTTP Listen address" default:":8080"`
	} `embed:"" prefix:"http."`

	// Private fields
	ctx    context.Context
	cancel context.CancelFunc
}

type CLI struct {
	Globals
	ConnectionCommands
	DatabaseCommands
	ExtensionCommands
	ReplicationSlotCommands
	RoleCommands
	SchemaCommands
	ObjectCommands
	ServerCommands
	SettingCommands
	StatementCommands
	TablespaceCommands
	VersionCommands
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func main() {
	cli := new(CLI)
	ctx := kong.Parse(cli)

	// Create the context and cancel function
	cli.Globals.ctx, cli.Globals.cancel = signal.NotifyContext(context.Background(), os.Interrupt)
	defer cli.Globals.cancel()

	// Call the Run() method of the selected parsed command.
	if err := ctx.Run(&cli.Globals); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (g *Globals) Client() (*httpclient.Client, error) {
	host, port, err := net.SplitHostPort(g.HTTP.Addr)
	if err != nil {
		return nil, err
	}

	// Client options
	opts := []client.ClientOpt{}
	if g.Debug {
		opts = append(opts, client.OptTrace(os.Stderr, true))
	}

	// Create a client
	url := fmt.Sprintf("http://%s:%s%s", host, port, g.HTTP.Prefix)
	return httpclient.New(url, opts...)
}
