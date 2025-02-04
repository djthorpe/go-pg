package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/djthorpe/go-pg"
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
	if err := conn.Exec(ctx, "CREATE TABLE names (id SERIAL PRIMARY KEY, name TEXT, gender TEXT, frequency INT)"); err != nil {
		panic(err)
	}

	// Import data
	if err := ingest(ctx, "https://www.ssa.gov/oact/babynames/names.zip", conn); err != nil {
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

func ingest(ctx context.Context, url string, conn pg.Conn) error {
	// Download the ZIP file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// If not found, then return the error
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	// Read the response body into a buffer
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %w", err)
	}

	// Create a zip reader from the buffer
	zReader, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return fmt.Errorf("Failed to create zip reader: %w", err)
	}

	// Iterate through each file in the zip archive
	for _, file := range zReader.File {
		if file.FileInfo().IsDir() {
			continue
		} else if filepath.Ext(file.Name) != ".txt" && filepath.Ext(file.Name) != ".csv" {
			continue
		}

		// Open the file inside the zip
		log.Print("Processing file:", file.Name)
		zf, err := file.Open()
		if err != nil {
			return fmt.Errorf("Failed to open zip file: %w", err)
		}
		defer zf.Close()

		// Process the file inside a transaction
		if err := conn.Tx(ctx, func(conn pg.Conn) error {
			return ingestFile(zf, conn)
		}); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

func ingestFile(r io.Reader, conn pg.Conn) error {
	decoder := csv.NewReader(r)
	for {
		record, err := decoder.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// Insert the record into the database
		var name Name
		if err := conn.Insert(context.Background(), &name, NewName(record...)); err != nil {
			return err
		} else {
			log.Print("Inserted:", name)
		}
	}
	return nil
}
