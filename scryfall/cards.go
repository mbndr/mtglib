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
	ImageURIs     map[string]string `json:"image_uris"`
	ManaCost      string            `json:"mana_cost"`
	Cmc           float32           `json:"cmc"`
	TypeLine      string            `json:"type_line"`
	OracleText    string            `json:"oracle_text"`
	Colors        []string          `json:"colors"`
	ColorIdentity []string          `json:"color_identity"`
	Set           string            `json:"set"`
	SetName       string            `json:"set_name"`
	CardFaces     []cardFace        `json:"card_faces"'`
}

type cardFace struct {
	CardID    string            `json:"-"` // must be set manually
	Colors    []string          `json:"colors"`
	ImageURIs map[string]string `json:"image_uris"`
	ManaCost  string            `json:"mana_cost"`
	Name      string            `json:"name"`
	TypeLine  string            `json:"type_line"`
}

type bulkObject struct {
	DownloadURI string `json:"download_uri"`
}

// ImportCards imports all cards (bulk data).
func ImportCards() error {
	if _, err := db.ExecSingle("DELETE FROM scryfall_cards"); err != nil {
		return errors.Wrap(err, "Cannot delete table")
	}

	if _, err := db.ExecSingle("DELETE FROM scryfall_card_faces"); err != nil {
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

	// collect card faces while inserting cards
	var cardFaces []cardFace

	_, err = db.Exec(`INSERT INTO scryfall_cards(
		scryfall_id,
		oracle_id,
		name,
		image_uri,
		mana_cost,
		cmc,
		type_line,
		oracle_text,
		colors,
		color_identity,
		set_code,
		set_name
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, len(records), func(stmt *sql.Stmt, i int) (sql.Result, error) {
		rec := records[i]

		for j := range rec.CardFaces {
			rec.CardFaces[j].CardID = rec.ScryfallID
		}
		cardFaces = append(cardFaces, rec.CardFaces...)

		return stmt.Exec(
			rec.ScryfallID,
			rec.OracleID,
			rec.Name,
			rec.ImageURIs["normal"],
			rec.ManaCost,
			rec.Cmc,
			rec.TypeLine,
			rec.OracleText,
			colorsToString(rec.Colors),
			colorsToString(rec.ColorIdentity),
			rec.Set,
			rec.SetName,
		)
	})

	return importCardFaces(cardFaces)
}

func importCardFaces(cardFaces []cardFace) error {
	_, err := db.Exec(`INSERT INTO scryfall_card_faces(
		card_id,
		colors,
		image_uri,
		mana_cost,
		name,
		type_line
	) VALUES (?, ?, ?, ?, ?, ?)`, len(cardFaces), func(stmt *sql.Stmt, i int) (sql.Result, error) {
		rec := cardFaces[i]

		return stmt.Exec(
			rec.CardID,
			colorsToString(rec.Colors),
			rec.ImageURIs["normal"],
			rec.ManaCost,
			rec.Name,
			rec.TypeLine,
		)
	})

	return err
}

func colorsToString(colors []string) string {
	return strings.Join(colors, "|")
}
