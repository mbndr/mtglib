package mtglib

import (
	"database/sql"
	"strings"

	"github.com/mbndr/mtglib/db"
)

const sqlStmtDistinctCards = `SELECT s.scryfall_id,
								s.oracle_id,
								s.name,
								s.image_uri,
								s.mana_cost,
								s.cmc,
								s.type_line,
								s.oracle_text,
								s.colors,
								s.color_identity,
								s.set_code,
								s.set_name,
								SUM(h.quantity)
							FROM helvault_library h
								INNER JOIN scryfall_cards s
								ON s.scryfall_id = h.scryfall_id
							GROUP BY oracle_id`

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
	// When a card has multiple faces
	CardFaces []CardFace
}

// CardFace represents a face of a multiface card
type CardFace struct {
	Colors   []string
	ImageURI string
	ManaCost string
	Name     string
	TypeLine string
}

// TotalLibraryCardCount returns the total amount of cards in collection
func TotalLibraryCardCount() int {
	var count int
	err := db.Select("SELECT SUM(quantity) FROM helvault_library", func(rows *sql.Rows) error {
		return rows.Scan(&count)
	})
	if err != nil {
		count = -1
	}

	return count
}

// LoadCards returns an map with the form oracle_id -> card and a slice with all oracle IDs
func LoadCards() (map[string]Card, []string, error) {
	cards := make(map[string]Card)
	oracleIDs := []string{}

	err := db.Select(sqlStmtDistinctCards, func(rows *sql.Rows) error {
		var c Card

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
		if err == nil {
			c.Colors = strings.Split(colorsStr, "|")
			c.ColorIdentity = strings.Split(colorIdentityStr, "|")

			cardFaces, err := loadCardFaces(c.ScryfallID)
			if err != nil {
				return err
			}

			c.CardFaces = cardFaces
			cards[c.OracleID] = c
			oracleIDs = append(oracleIDs, c.OracleID)
		}

		return err
	})
	if err != nil {
		return nil, nil, err
	}

	return cards, oracleIDs, nil
}

func loadCardFaces(scryfallID string) ([]CardFace, error) {
	var cardFaces []CardFace

	err := db.Select("SELECT colors, image_uri, mana_cost, name, type_line FROM scryfall_card_faces WHERE card_id = ?", func(rows *sql.Rows) error {
		var cf CardFace
		var colorsStr string
		err := rows.Scan(
			&colorsStr, &cf.ImageURI, &cf.ManaCost, &cf.Name, &cf.TypeLine,
		)
		cf.Colors = strings.Split(colorsStr, "|")
		cardFaces = append(cardFaces, cf)
		return err
	}, scryfallID)

	if err != nil {
		return nil, err
	}

	return cardFaces, nil
}
