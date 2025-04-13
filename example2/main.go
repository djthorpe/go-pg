package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	// Packages
	pg "github.com/djthorpe/go-pg"
	test "github.com/djthorpe/go-pg/pkg/test"
)

func main() {
	// Cancel on CTRL+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Create a postgres database
	log.Print("Starting postgresql")
	container, conn, err := test.NewPgxContainer(ctx, "postgres", true, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	defer container.Close(context.Background())

	// Create the 'record' table
	if err := conn.Exec(ctx, "CREATE TABLE record (id SERIAL PRIMARY KEY, value INT NOT NULL DEFAULT 0)"); err != nil {
		panic(err)
	}

	// Create 10 records
	log.Print("Creating 10 records")
	for i := 0; i < 10; i++ {
		var result Record
		if err := conn.Insert(ctx, &result, Record{}); err != nil {
			panic(err)
		}
		log.Print("Insert=", result)
	}

	// Update 10 records
	log.Print("Update 10 records")
	var result RecordList
	if err := conn.Update(ctx, &result, Record{}, nil); err != nil {
		panic(err)
	}
	for _, record := range result {
		log.Print("Update=", record)
	}

	// Select 10 records
	log.Print("Select 10 records")
	var result2 RecordList
	if err := conn.List(ctx, &result2, Record{}); err != nil {
		panic(err)
	}
	for _, record := range result2 {
		log.Print("Select=", record)
	}

	// Delete 10 records
	log.Print("Delete 10 records")
	var result3 RecordList
	if err := conn.Delete(ctx, &result3, Record{}); err != nil {
		panic(err)
	}
	for _, record := range result3 {
		log.Print("Delete=", record)
	}

	log.Print("Done")
}

///////////////////////////////////////////////////////////////////////////////

type Record struct {
	Id    int
	Value int
}

type RecordList []Record

func (r Record) Insert(_ *pg.Bind) (string, error) {
	return "INSERT INTO record DEFAULT VALUES RETURNING id, value", nil
}

func (r Record) Update(*pg.Bind) error {
	return pg.ErrNotImplemented
}

func (r Record) Select(_ *pg.Bind, op pg.Op) (string, error) {
	switch op {
	case pg.Update:
		return "UPDATE record SET value=id RETURNING id, value", nil
	case pg.List:
		return "SELECT id, value FROM record", nil
	case pg.Delete:
		return "DELETE FROM record RETURNING id, value", nil
	}
	return "", pg.ErrNotImplemented
}

func (r *Record) Scan(row pg.Row) error {
	return row.Scan(&r.Id, &r.Value)
}

func (r *RecordList) Scan(row pg.Row) error {
	var record Record
	if err := record.Scan(row); err != nil {
		return err
	}
	*r = append(*r, record)
	return nil
}
