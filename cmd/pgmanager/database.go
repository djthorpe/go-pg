package main

import (
	"fmt"

	// Packages
	schema "github.com/djthorpe/go-pg/pkg/manager/schema"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type DatabaseCommands struct {
	ListDatabase   ListDatabaseCommand   `cmd:"" name:"databases" help:"List databases."`
	GetDatabase    GetDatabaseCommand    `cmd:"" name:"database" help:"Get database."`
	CreateDatabase CreateDatabaseCommand `cmd:"" name:"create-database" help:"Create database."`
	DeleteDatabase DeleteDatabaseCommand `cmd:"" name:"delete-database" help:"Delete database."`
}

type ListDatabaseCommand struct{}

type GetDatabaseCommand struct {
	Name string `arg:"" name:"name" help:"Database name"`
}

type DeleteDatabaseCommand struct {
	GetDatabaseCommand
}

type CreateDatabaseCommand struct {
	GetDatabaseCommand
	Owner string `name:"owner" help:"Database owner"`
}

///////////////////////////////////////////////////////////////////////////////
// COMMANDS

func (cmd *ListDatabaseCommand) Run(ctx *Globals) error {
	client, err := ctx.Client()
	if err != nil {
		return err
	}

	// List databases
	databases, err := client.ListDatabases(ctx.ctx)
	if err != nil {
		return err
	}

	// Print
	fmt.Println(databases)
	return nil
}

func (cmd *GetDatabaseCommand) Run(ctx *Globals) error {
	client, err := ctx.Client()
	if err != nil {
		return err
	}

	// Get one database
	database, err := client.GetDatabase(ctx.ctx, cmd.Name)
	if err != nil {
		return err
	}

	// Print
	fmt.Println(database)
	return nil
}

func (cmd *CreateDatabaseCommand) Run(ctx *Globals) error {
	client, err := ctx.Client()
	if err != nil {
		return err
	}

	// Get one database
	database, err := client.CreateDatabase(ctx.ctx, schema.DatabaseMeta{
		Name:  cmd.Name,
		Owner: cmd.Owner,
	})
	if err != nil {
		return err
	}

	// Print
	fmt.Println(database)
	return nil
}

func (cmd *DeleteDatabaseCommand) Run(ctx *Globals) error {
	client, err := ctx.Client()
	if err != nil {
		return err
	}

	// Get one database
	if err := client.DeleteDatabase(ctx.ctx, cmd.Name); err != nil {
		return err
	}

	// Return success
	return nil
}
