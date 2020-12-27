package scryfall

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mbndr/mtglib/db"
	"github.com/pkg/errors"
)

var (
	// APIBulkCards returns info of bulk card data
	APIBulkCards = "https://api.scryfall.com/bulk-data/all-cards"
)

// ScryfallCard is a card object parsed from the Scryfall API
// Only necessary fields are fetched from json
type card struct {
	ScryfallID    string            `json:"id"`
	OracleID      string            `json:"oracle_id"`
	Name          string            `json:"name"`
	Lang          string            `json:"lang"`
	ImageURIs     map[string]string `json:"image_uris"`
	ManaCost      string            `json:"mana_cost"`
	Cmc           float32           `json:"cmc"`
	TypeLine      string            `json:"type_line"`
	OracleText    string            `json:"oracle_text"`
	Colors        []string          `json:"colors"`
	ColorIdentity []string          `json:"color_identity"`
	Set           string            `json:"set"`
	SetName       string            `json:"set_name"`
}

type bulkObject struct {
	DownloadURI string `json:"download_uri"`
}

// ImportCards imports all cards (bulk data).
func ImportCards() error {
	if _, err := db.ExecSingle("DELETE FROM scryfall_cards"); err != nil {
		return errors.Wrap(err, "Cannot delete table")
	}

	res, err := http.Get(APIBulkCards)
	if err != nil {
		return errors.Wrap(err, "Cannot download bulk object")
	}
	defer res.Body.Close()

	var bulkData bulkObject

	if err = json.NewDecoder(res.Body).Decode(&bulkData); err != nil {
		return errors.Wrap(err, "Cannot parse bulk object")
	}

	return importFromBulk(bulkData.DownloadURI)
}

func importFromBulk(bulkURI string) error {
	res, err := http.Get(bulkURI)
	if err != nil {
		return errors.Wrap(err, "Cannot download cards")
	}
	defer res.Body.Close()

	var records []card

	if err = json.NewDecoder(res.Body).Decode(&records); err != nil {
		return errors.Wrap(err, "Cannot parse cards")
	}

	_, err = db.Exec(`INSERT INTO scryfall_cards(
		scryfall_id,
		oracle_id,
		name,
		lang,
		image_uri,
		mana_cost,
		cmc,
		type_line,
		oracle_text,
		colors,
		color_identity,
		set_code,
		set_name
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, len(records), func(stmt *sql.Stmt, i int) (sql.Result, error) {
		rec := records[i]

		imageURI := rec.ImageURIs["normal"]
		colors := strings.Join(rec.Colors, "|")
		colorIdentity := strings.Join(rec.ColorIdentity, "|")

		return stmt.Exec(
			rec.ScryfallID,
			rec.OracleID,
			rec.Name,
			rec.Lang,
			imageURI,
			rec.ManaCost,
			rec.Cmc,
			rec.TypeLine,
			rec.OracleText,
			colors,
			colorIdentity,
			rec.Set,
			rec.SetName,
		)
	})

	return err
}
