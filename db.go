package mtglib

import (
	"database/sql"
	"log"

	// databse adapter
	_ "github.com/mattn/go-sqlite3"
)

const (
	// DBHelvaultLibrary name
	DBHelvaultLibrary = "helvault_library"
	// DBScryfallCards name
	DBScryfallCards = "scryfall_cards"

	// DefaultCardSelectFields are the field fetched for getting a card
	DefaultCardSelectFields = "s.scryfall_id, s.oracle_id, s.name, s.image_uri, s.mana_cost, s.cmc, s.type_line, s.oracle_text, s.colors, s.color_identity, s.set_code, s.set_name"

	// SQLCountAllCards selects the count of ALL cards in the library
	SQLCountAllCards = "SELECT SUM(quantity) AS count FROM " + DBHelvaultLibrary
	// SQLDistinctCards selects the count of the distinct cards (set doesn't matter)
	SQLDistinctCards = "SELECT " + DefaultCardSelectFields + ", SUM(h.quantity) as quantity FROM " + DBHelvaultLibrary + " h INNER JOIN " + DBScryfallCards + " s ON s.scryfall_id = h.scryfall_id GROUP BY oracle_id"
)

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
