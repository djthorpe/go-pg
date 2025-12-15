package schema

import (
	"fmt"
	"strings"
	"time"

	// Packages
	pg "github.com/djthorpe/go-pg"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Statement represents a row from pg_stat_statements
type Statement struct {
	UserID            uint64  `json:"userid"`                        // OID of user who executed the statement
	UserName          string  `json:"username,omitempty"`            // Name of the user (joined from pg_roles)
	DatabaseID        uint64  `json:"dbid"`                          // OID of database in which the statement was executed
	DatabaseName      string  `json:"database,omitempty"`            // Name of the database (joined from pg_database)
	QueryID           int64   `json:"queryid"`                       // Hash code to identify identical normalized queries
	Query             string  `json:"query"`                         // Text of a representative statement
	Calls             int64   `json:"calls"`                         // Number of times the statement was executed
	TotalExecTime     float64 `json:"total_exec_time"`               // Total time spent executing the statement, in milliseconds
	MinExecTime       float64 `json:"min_exec_time"`                 // Minimum time spent executing the statement, in milliseconds
	MaxExecTime       float64 `json:"max_exec_time"`                 // Maximum time spent executing the statement, in milliseconds
	MeanExecTime      float64 `json:"mean_exec_time"`                // Mean time spent executing the statement, in milliseconds
	StddevExecTime    float64 `json:"stddev_exec_time"`              // Population standard deviation of time spent executing the statement
	Rows              int64   `json:"rows"`                          // Total number of rows retrieved or affected by the statement
	SharedBlksHit     int64   `json:"shared_blks_hit,omitempty"`     // Total number of shared block cache hits by the statement
	SharedBlksRead    int64   `json:"shared_blks_read,omitempty"`    // Total number of shared blocks read by the statement
	SharedBlksDirtied int64   `json:"shared_blks_dirtied,omitempty"` // Total number of shared blocks dirtied by the statement
	SharedBlksWritten int64   `json:"shared_blks_written,omitempty"` // Total number of shared blocks written by the statement
	LocalBlksHit      int64   `json:"local_blks_hit,omitempty"`      // Total number of local block cache hits by the statement
	LocalBlksRead     int64   `json:"local_blks_read,omitempty"`     // Total number of local blocks read by the statement
	LocalBlksDirtied  int64   `json:"local_blks_dirtied,omitempty"`  // Total number of local blocks dirtied by the statement
	LocalBlksWritten  int64   `json:"local_blks_written,omitempty"`  // Total number of local blocks written by the statement
	TempBlksRead      int64   `json:"temp_blks_read,omitempty"`      // Total number of temp blocks read by the statement
	TempBlksWritten   int64   `json:"temp_blks_written,omitempty"`   // Total number of temp blocks written by the statement
}

// StatementList is a list of statements with a total count
type StatementList struct {
	Count uint64      `json:"count"`
	Body  []Statement `json:"body"`
}

// StatementListRequest contains parameters for listing statements
type StatementListRequest struct {
	pg.OffsetLimit

	// Filter by database name
	Database *string `json:"database,omitempty"`

	// Filter by user name
	User *string `json:"user,omitempty"`

	// Order by field (calls, total_exec_time, mean_exec_time, rows)
	OrderBy string `json:"order_by,omitempty"`

	// Order direction (asc, desc) - default desc
	OrderDir string `json:"order_dir,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	StatementListLimit = 100
)

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s Statement) String() string {
	// Truncate query for display
	query := s.Query
	if len(query) > 60 {
		query = query[:57] + "..."
	}
	return fmt.Sprintf("Statement{query=%q calls=%d total_time=%.2fms mean_time=%.2fms rows=%d}",
		query, s.Calls, s.TotalExecTime, s.MeanExecTime, s.Rows)
}

func (l StatementList) String() string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("StatementList{count=%d body=[", l.Count))
	for i, s := range l.Body {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(s.String())
	}
	result.WriteString("]}")
	return result.String()
}

///////////////////////////////////////////////////////////////////////////////
// SELECTOR

func (r *StatementListRequest) Select(bind *pg.Bind, op pg.Op) (string, error) {
	if op != pg.List {
		return "", pg.ErrNotImplemented.Withf("unsupported operation %q", op)
	}

	// Build WHERE clause
	var where []string
	if r.Database != nil && *r.Database != "" {
		bind.Set("database", *r.Database)
		where = append(where, "d.datname = @database")
	}
	if r.User != nil && *r.User != "" {
		bind.Set("user", *r.User)
		where = append(where, "u.rolname = @user")
	}

	if len(where) > 0 {
		bind.Set("where", "WHERE "+strings.Join(where, " AND "))
	} else {
		bind.Set("where", "")
	}

	// Build ORDER BY clause
	orderBy := "total_exec_time"
	if r.OrderBy != "" {
		switch strings.ToLower(r.OrderBy) {
		case "calls":
			orderBy = "calls"
		case "total_exec_time", "total_time":
			orderBy = "total_exec_time"
		case "mean_exec_time", "mean_time":
			orderBy = "mean_exec_time"
		case "rows":
			orderBy = "rows"
		case "min_exec_time", "min_time":
			orderBy = "min_exec_time"
		case "max_exec_time", "max_time":
			orderBy = "max_exec_time"
		}
	}

	orderDir := "DESC"
	if r.OrderDir != "" && strings.ToUpper(r.OrderDir) == "ASC" {
		orderDir = "ASC"
	}
	bind.Set("orderby", fmt.Sprintf("ORDER BY %s %s", orderBy, orderDir))

	// Set offset/limit
	r.OffsetLimit.Bind(bind, StatementListLimit)

	return statementList, nil
}

///////////////////////////////////////////////////////////////////////////////
// READER

func (s *Statement) Scan(row pg.Row) error {
	return row.Scan(
		&s.UserID, &s.UserName, &s.DatabaseID, &s.DatabaseName, &s.QueryID, &s.Query,
		&s.Calls, &s.TotalExecTime, &s.MinExecTime, &s.MaxExecTime, &s.MeanExecTime, &s.StddevExecTime,
		&s.Rows, &s.SharedBlksHit, &s.SharedBlksRead, &s.SharedBlksDirtied, &s.SharedBlksWritten,
		&s.LocalBlksHit, &s.LocalBlksRead, &s.LocalBlksDirtied, &s.LocalBlksWritten,
		&s.TempBlksRead, &s.TempBlksWritten,
	)
}

func (l *StatementList) Scan(row pg.Row) error {
	var statement Statement
	if err := statement.Scan(row); err != nil {
		return err
	}
	l.Body = append(l.Body, statement)
	return nil
}

func (l *StatementList) ScanCount(row pg.Row) error {
	return row.Scan(&l.Count)
}

///////////////////////////////////////////////////////////////////////////////
// QUERIES

// pg_stat_statements query - joins with pg_roles and pg_database for names
const statementSelect = `
	SELECT
		s.userid::BIGINT,
		COALESCE(u.rolname, '') AS username,
		s.dbid::BIGINT,
		COALESCE(d.datname, '') AS database,
		s.queryid,
		s.query,
		s.calls,
		s.total_exec_time,
		s.min_exec_time,
		s.max_exec_time,
		s.mean_exec_time,
		s.stddev_exec_time,
		s.rows,
		s.shared_blks_hit,
		s.shared_blks_read,
		s.shared_blks_dirtied,
		s.shared_blks_written,
		s.local_blks_hit,
		s.local_blks_read,
		s.local_blks_dirtied,
		s.local_blks_written,
		s.temp_blks_read,
		s.temp_blks_written
	FROM
		public.pg_stat_statements s
	LEFT JOIN
		pg_catalog.pg_roles u ON s.userid = u.oid
	LEFT JOIN
		pg_catalog.pg_database d ON s.dbid = d.oid
`

const statementList = `WITH q AS (` + statementSelect + `) SELECT * FROM q ${where} ${orderby}`

// Ensure time is imported (used in potential future extensions)
var _ = time.Now
