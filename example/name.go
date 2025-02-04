package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/djthorpe/go-pg"
)

type Name struct {
	Id        uint64 `json:"id,omitempty"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Frequency uint64 `json:"frequency"`
}

type NameList struct {
	Count uint64  `json:"count"`
	Names []*Name `json:"names"`
}

// Create a new name from the CSV file which consists of three cells
// name, gender, frequency
func NewName(cells ...string) *Name {
	if len(cells) != 3 {
		return nil
	}
	name := cells[0]
	gender := cells[1]
	frequency, err := strconv.ParseUint(cells[2], 10, 64)
	if err != nil {
		return nil
	}
	return &Name{
		Name:      name,
		Gender:    gender,
		Frequency: frequency,
	}
}

// Stringify
func (obj Name) String() string {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// Reader - bind to object
func (obj *Name) Scan(row pg.Row) error {
	return row.Scan(&obj.Id, &obj.Name, &obj.Gender, &obj.Frequency)
}

func (list *NameList) Scan(row pg.Row) error {
	var obj Name
	if err := row.Scan(&obj.Id, &obj.Name, &obj.Gender, &obj.Frequency); err != nil {
		return err
	} else {
		list.Names = append(list.Names, &obj)
	}
	return nil
}

// Bind count
func (list *NameList) ScanCount(row pg.Row) error {
	return row.Scan(&list.Count)
}

// Selector - select rows from database
func (list NameList) Select(bind *pg.Bind, op pg.Op) (string, error) {
	switch op {
	case pg.List:
		return `SELECT id, name, gender, frequency FROM names`, nil
	default:
		return "", fmt.Errorf("Unsupported operation: %q", op)
	}
}

// Selector - select rows from database
func (obj Name) Select(bind *pg.Bind, op pg.Op) (string, error) {
	bind.Set("id", obj.Id)
	switch op {
	case pg.Get:
		return `SELECT id, name, gender, frequency FROM names WHERE id=@id`, nil
	case pg.Patch:
		return `UPDATE names SET @patch WHERE id=@id RETURNING id, name, gender, frequency`, nil
	default:
		return "", fmt.Errorf("Unsupported operation: %q", op)
	}
}

// Writer - insert object
func (obj Name) Insert(bind *pg.Bind) (string, error) {
	bind.Set("name", obj.Name)
	bind.Set("gender", obj.Gender)
	bind.Set("frequency", obj.Frequency)
	return `INSERT INTO names (name, gender, frequency) VALUES (@name, @gender, @frequency) RETURNING id, name, gender, frequency`, nil
}

// Writer - patch object
func (obj Name) Patch(bind *pg.Bind) error {
	// Reset the patch parameters
	bind.Del("patch")

	// Append patches
	if obj.Name != "" {
		bind.Append("patch", `name=`+bind.Set("name", obj.Name))
	}
	if obj.Gender != "" {
		bind.Append("patch", `gender=`+bind.Set("gender", obj.Gender))
	}
	if obj.Frequency != 0 {
		bind.Append("patch", `frequency=`+bind.Set("frequency", obj.Frequency))
	}

	// If nothing was patched, then return an error
	if patch := bind.Join("patch", ", "); patch == "" {
		return fmt.Errorf("No patch parameters")
	} else {
		bind.Set("patch", patch)
	}

	// Return success
	return nil
}
