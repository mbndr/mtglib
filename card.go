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

func getTotalCardCount() int {
	rows, err := db.Query(SQLCountAllCards)
	if err != nil {
		return -1
	}

	var count int = -1

	if rows.Next() {
		rows.Scan(&count)
	}

	return count
}

// LoadCards returns an map with the form oracle_id -> card and a slice with all oracle IDs
func LoadCards() (map[string]*Card, []string, error) {
	cards := make(map[string]*Card)
	oracleIDs := []string{}

	rows, err := db.Query(SQLDistinctCards)
	if err != nil {
		return nil, nil, err
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
			return nil, nil, err
		}

		c.Colors = strings.Split(colorsStr, "|")
		c.ColorIdentity = strings.Split(colorIdentityStr, "|")

		cards[c.OracleID] = c
		oracleIDs = append(oracleIDs, c.OracleID)
	}

	return cards, oracleIDs, nil
}
