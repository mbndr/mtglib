package mtglib

import (
	"strings"
)

// Card get from db, for general use (joined data)
type Card struct {
	ScryfallID    string
	OracleID      string
	Name          string
	ImageURI      string
	ManaCost      string // {2}{B}{W}
	Cmc           float32
	TypeLine      string
	OracleText    string
	Colors        []string
	ColorIdentity []string
	SetCode       string
	SetName       string
	Quantity      int
}

// LoadCards returns an map with the form oracle_id -> card
func LoadCards() (map[string]*Card, error) {
	cards := make(map[string]*Card)

	rows, err := db.Query(SQLDistinctCards)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := &Card{}

		// special treatment
		var colorsStr string
		var colorIdentityStr string

		err := rows.Scan(
			&c.ScryfallID,
			&c.OracleID,
			&c.Name,
			&c.ImageURI,
			&c.ManaCost,
			&c.Cmc,
			&c.TypeLine,
			&c.OracleText,
			&colorsStr,
			&colorIdentityStr,
			&c.SetCode,
			&c.SetName,
			&c.Quantity,
		)
		if err != nil {
			return nil, err
		}

		c.Colors = strings.Split(colorsStr, "|")
		c.ColorIdentity = strings.Split(colorIdentityStr, "|")

		cards[c.OracleID] = c
	}

	return cards, nil
}
