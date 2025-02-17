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
	"path/filepath"
	"regexp"
	"strconv"

	// Packages
	pg "github.com/djthorpe/go-pg"
)

var (
	reFilename = regexp.MustCompile(`^yob(\d{4})\.txt$`)
)

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func ingest(ctx context.Context, url string, conn pg.Conn) (int, error) {
	n := 0

	// Download the ZIP file
	resp, err := http.Get(url)
	if err != nil {
		return n, err
	}
	defer resp.Body.Close()

	// If not found, then return the error
	if resp.StatusCode != http.StatusOK {
		return n, errors.New(resp.Status)
	}

	// Read the response body into a buffer
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return n, fmt.Errorf("Failed to read response body: %w", err)
	}

	// Create a zip reader from the buffer
	zReader, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return n, fmt.Errorf("Failed to create zip reader: %w", err)
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
			return n, fmt.Errorf("Failed to open zip file: %w", err)
		}
		defer zf.Close()

		// Determine the year from the filename
		var year uint64
		if year_ := reFilename.FindStringSubmatch(filepath.Base(file.Name)); len(year_) != 2 {
			log.Print("Skipping file:", file.Name)
			continue
		} else if year__, err := strconv.ParseUint(year_[1], 10, 64); err != nil {
			log.Print("Skipping file:", file.Name)
			continue
		} else {
			year = year__
		}

		// Process the file inside a transaction
		if err := conn.Tx(ctx, func(conn pg.Conn) error {
			// Do bulk insert
			return conn.Bulk(ctx, func(conn pg.Conn) error {
				m, err := ingestFile(ctx, zf, conn, year)
				if err != nil {
					return err
				}
				n += m
				return nil
			})
		}); err != nil {
			return n, err
		}

		log.Print("  ", n, " records ingested")
	}

	// Return success
	return n, nil
}

func ingestFile(ctx context.Context, r io.Reader, conn pg.Conn, year uint64) (int, error) {
	n := 0
	decoder := csv.NewReader(r)
	for {
		record, err := decoder.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return n, err
		}

		// Insert the record into the database
		var name Name
		if err := conn.Insert(ctx, &name, NewName(year, record...)); err != nil {
			return n, err
		} else {
			n++
		}
	}
	return n, nil
}
