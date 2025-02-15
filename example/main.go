package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	// Packages
	"github.com/djthorpe/go-pg/pkg/test"
)

func main() {
	// Cancel on CTRL+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Create a postgres database
	log.Print("Starting postgresql")
	container, conn, err := test.NewPgxContainer(ctx, "postgres", true)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	defer container.Close(context.Background())

	// Create the 'names' table
	if err := conn.Exec(ctx, "CREATE TABLE names (id SERIAL PRIMARY KEY, name TEXT, gender TEXT, frequency INT, year INT)"); err != nil {
		panic(err)
	}

	// Import data
	if _, err := ingest(ctx, "https://www.ssa.gov/oact/babynames/names.zip", conn); err != nil {
		panic(err)
	}

	// Start API server
	log.Print("Starting web server http://localhost:8080")
	server := http.Server{
		Addr:    ":8080",
		Handler: RegisterHandlers(http.NewServeMux(), conn),
	}
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
