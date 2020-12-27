package scryfall

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/mbndr/mtglib/db"
	"github.com/pkg/errors"
)

const (
	APISymbols = "https://api.scryfall.com/symbology"
)

type symbol struct {
	Symbol string `json:"symbol"`
	SvgURI string `json:"svg_uri"`
	Title  string `json:"english"`
}

type symbolList struct {
	// Not checking if the list is multi page!
	Data []symbol `json:"data"`
}

// ImportSymbols imports all card symbols
func ImportSymbols() error {
	if _, err := db.ExecSingle("DELETE FROM scryfall_symbols"); err != nil {
		return errors.Wrap(err, "Cannot delete table")
	}

	// Getting data
	res, err := http.Get(APISymbols)
	if err != nil {
		return errors.Wrap(err, "Cannot download symbols")
	}
	defer res.Body.Close()

	var list symbolList

	if err = json.NewDecoder(res.Body).Decode(&list); err != nil {
		return errors.Wrap(err, "Cannot parse symbols")
	}

	// Inserting data
	_, err = db.Exec("INSERT INTO scryfall_symbols(symbol, svg_uri, title) VALUES(?, ?, ?)", len(list.Data), func(stmt *sql.Stmt, i int) (sql.Result, error) {
		return stmt.Exec(list.Data[i].Symbol, list.Data[i].SvgURI, list.Data[i].Title)
	})

	return err
}
