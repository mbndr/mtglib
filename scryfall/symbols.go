package scryfall

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/mbndr/mtglib"
	"github.com/mbndr/mtglib/db"
	"github.com/pkg/errors"
)

const (
	APISymbols = "https://api.scryfall.com/symbology"
)

type symbolList struct {
	// Not checking if the list is multi page!
	Data []mtglib.Symbol `json:"data"`
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

	return downloadSymbols(list)
}

func downloadSymbols(list symbolList) error {
	// Download
	for _, sym := range list.Data {
		res, err := http.Get(sym.SvgURI)
		if err != nil {
			return errors.Wrap(err, "Cannot download symbol image")
		}
		defer res.Body.Close()

		file, err := os.OpenFile(path.Join("resources", path.Base(sym.SvgURI)), os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			return errors.Wrap(err, "Cannot create symbol file")
		}
		defer file.Close()

		_, err = io.Copy(file, res.Body)
		if err != nil {
			return errors.Wrap(err, "Cannot copy symbol")
		}
	}

	return nil
}
