package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	// Packages
	kong "github.com/alecthomas/kong"
	httpclient "github.com/djthorpe/go-pg/pkg/manager/httpclient"
	client "github.com/mutablelogic/go-client"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Globals struct {
	// Debug option
	Debug bool `name:"debug" help:"Enable debug logging"`

	// HTTP server options
	HTTP struct {
		Prefix string `name:"prefix" help:"HTTP path prefix" default:"/api/v1"`
		Listen string `name:"http" help:"HTTP Listen address" default:":8080"`
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
	RoleCommands
	SchemaCommands
	ObjectCommands
	ServerCommands
	TablespaceCommands
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
	host, port, err := net.SplitHostPort(g.HTTP.Listen)
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
