package mtglib

import (
	"encoding/csv"
	"io"
	"strconv"
)

// HelvaultRecord is a line in the csv export
type HelvaultRecord struct {
	ScryfallID string
	Quantity   int
}

func readHelvaultCsv(reader io.Reader) ([]*HelvaultRecord, error) {
	r := csv.NewReader(reader)
	r.Comma = ','

	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	records := []*HelvaultRecord{}

	for _, line := range lines[1:] {
		quantity, err := strconv.Atoi(line[5])
		if err != nil {
			return nil, err
		}

		records = append(records, &HelvaultRecord{
			Quantity:   quantity,
			ScryfallID: line[6],
		})
	}

	return records, nil
}

func insertHelvaultRecords(records []*HelvaultRecord) error {
	sql := `INSERT INTO helvault_library(
			scryfall_id,
			quantity
		) VALUES (?, ?)`

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, rec := range records {
		_, err = stmt.Exec(
			rec.ScryfallID,
			rec.Quantity,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// LoadHelvault imports a csv into the database
func LoadHelvault(r io.Reader) error {
	err := truncateTable("helvault_library")
	if err != nil {
		return err
	}

	records, err := readHelvaultCsv(r)
	if err != nil {
		return err
	}

	return insertHelvaultRecords(records)
}
