package mtglib

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
)

// ScryfallCard is a card object parsed from the Scryfall API
// Only necessary fields are fetched from json
type ScryfallCard struct {
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

func readScryfallRecords(reader io.Reader) ([]*ScryfallCard, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	records := []*ScryfallCard{}

	err = json.Unmarshal(bytes, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func insertScryfallRecords(records []*ScryfallCard) error {
	sql := `INSERT INTO scryfall_cards(
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
		set_identifier,
		set_name
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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
		imageURI := rec.ImageURIs["normal"]
		colors := strings.Join(rec.Colors, "|")
		colorIdentity := strings.Join(rec.ColorIdentity, "|")

		_, err = stmt.Exec(
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

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ImportScryfall imports the Scryfall bulk data into the database
func ImportScryfall(r io.Reader) error {
	err := truncateTable("scryfall_cards")
	if err != nil {
		return err
	}

	records, err := readScryfallRecords(r)
	if err != nil {
		return err
	}

	return insertScryfallRecords(records)
}
