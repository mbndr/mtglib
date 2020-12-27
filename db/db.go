package db

import (
	"database/sql"

	// database adapter
	_ "github.com/mattn/go-sqlite3"
)

// SelectRowCallback is called on each row in a select statement.
type SelectRowCallback func(*sql.Rows) error

// ExecRowCallback is called x times (count in Exec() method).
type ExecRowCallback func(*sql.Stmt, int) (sql.Result, error)

var db *sql.DB

// Open opens the database.
func Open(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	return err
}

// Select runs a select statement and calls a callback on each row.
func Select(sqlStatement string, callback SelectRowCallback, args ...interface{}) error {
	rows, err := db.Query(sqlStatement, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		callback(rows)
	}

	return rows.Err()
}

// ExecSingle executes a statement a single time.
func ExecSingle(sqlStatement string, args ...interface{}) (int64, error) {
	return Exec(sqlStatement, 1, func(stmt *sql.Stmt, i int) (sql.Result, error) {
		return stmt.Exec(args...)
	})
}

// Exec prepares a statement calls a callback count times.
// The callback is responsible for executing the statement.
// A list can help to execute the statement for multiple records.
func Exec(sqlStatement string, count int, callback ExecRowCallback) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(sqlStatement)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	var rowsAffected int64 = 0

	for i := 0; i < count; i++ {
		res, err := callback(stmt, i)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		ra, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		rowsAffected += ra
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
