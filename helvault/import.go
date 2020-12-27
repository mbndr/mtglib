package helvault

import (
	"database/sql"
	"encoding/csv"
	"os"
	"strconv"

	"github.com/mbndr/mtglib/db"
	"github.com/pkg/errors"
)

type helvaultRecord struct {
	scryfallID string
	quantity   int
}

// Import data from helvault
func Import(filename string) error {
	if _, err := db.ExecSingle("DELETE FROM helvault_library"); err != nil {
		return errors.Wrap(err, "Cannot delete table")
	}

	records, err := readCSV(filename)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO helvault_library(scryfall_id, quantity) VALUES(?, ?)", len(records), func(stmt *sql.Stmt, i int) (sql.Result, error) {
		return stmt.Exec(records[i].scryfallID, records[i].quantity)
	})

	return err
}

func readCSV(filename string) ([]helvaultRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ','

	lines, err := r.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot read CSV")
	}

	records := make([]helvaultRecord, len(lines)-1)

	for i, line := range lines[1:] {
		quantity, err := strconv.Atoi(line[5])
		if err != nil {
			return nil, errors.Wrap(err, "Cannot convert quantity")
		}

		records[i] = helvaultRecord{
			quantity:   quantity,
			scryfallID: line[6],
		}
	}

	return records, nil
}
