package mtglib

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

/*
Count:
ALL: select count(*)
DISTINCT: select count(distinct oracle_id)

All Distinct cards:
select * from helvault_library group by oracle_id;
*/

var db *sql.DB

func init() {
	var err error
	log.Println("loading db")

	db, err = sql.Open("sqlite3", "data/library.db")
	if err != nil {
		log.Fatal(err)
	}
}

func truncateTable(name string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM " + name)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	log.Printf("Deleted %d rows", rows)
	return err
}
