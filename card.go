package mtglib

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
								SUM(h.quantity),
								s.rarity
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
	Quantity      int
	Rarity        string
	// When a card has multiple faces
	CardFaces []CardFace
}

// CardFace represents a face of a multiface card
type CardFace struct {
	Colors     []string
	ImageURI   string
	ManaCost   string
	Name       string
	TypeLine   string
	OracleText string
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

// LoadCards returns all cards in the library loaded from db
func LoadCards() (*CardCollection, error) {
	var cc CardCollection
	cc.cards = make(map[string]Card)

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
			&c.Quantity,
			&c.Rarity,
		)
		if err == nil {
			c.Colors = strToColorArr(colorsStr)
			c.ColorIdentity = strToColorArr(colorIdentityStr)

			cardFaces, err := loadCardFaces(c.ScryfallID)
			if err != nil {
				return err
			}

			c.CardFaces = cardFaces
			cc.cards[c.OracleID] = c
			oracleIDs = append(oracleIDs, c.OracleID)
		}

		return err
	})
	if err != nil {
		return nil, err
	}

	return &cc, nil
}

func loadCardFaces(scryfallID string) ([]CardFace, error) {
	var cardFaces []CardFace

	err := db.Select("SELECT colors, image_uri, mana_cost, name, type_line, oracle_text FROM scryfall_card_faces WHERE card_id = ?", func(rows *sql.Rows) error {
		var cf CardFace
		var colorsStr string
		err := rows.Scan(
			&colorsStr, &cf.ImageURI, &cf.ManaCost, &cf.Name, &cf.TypeLine, &cf.OracleText,
		)

		cf.Colors = strToColorArr(colorsStr)

		cardFaces = append(cardFaces, cf)
		return err
	}, scryfallID)

	if err != nil {
		return nil, err
	}

	return cardFaces, nil
}

func strToColorArr(s string) []string {
	if s != "" {
		return strings.Split(s, "|")
	}
	return []string{}
}

// CardCollection wraps cards and helper methods
type CardCollection struct {
	cards map[string]Card // key: oracleID
}

// Get returns a card or nil
func (cc *CardCollection) Get(oracleID string) *Card {
	card, ok := cc.cards[oracleID]
	if !ok {
		return nil
	}
	return &card
}

// Count returns the count of cards
func (cc *CardCollection) Count() int {
	return len(cc.cards)
}

const card404 = "/static/img/card_404.jpg"

// CardImageURL returns the local url of an card image.
// Downloads the image if it doesn't exists.
func (cc *CardCollection) CardImageURL(oracleID string) string {
	card := cc.Get(oracleID)
	if card == nil {
		return card404
	}

	imgPath := fmt.Sprintf("./resources/%s.jpg", oracleID)

	err := downloadImageIfNotExist(imgPath, card.ImageURI)
	if err != nil {
		return card404
	}

	return fmt.Sprintf("/resources/%s.jpg", oracleID)
}

// FaceImageURLs returns a list of card face images.
// If the card itself has an image, use this instead.
func (cc *CardCollection) FaceImageURLs(oracleID string) []string {
	card := cc.Get(oracleID)
	if card == nil {
		return []string{card404}
	}

	uris := make([]string, len(card.CardFaces))

	// if the parent card has an image, use this instead for each face
	if card.ImageURI != "" {
		for i := range uris {
			uris[i] = cc.CardImageURL(card.OracleID)
		}
		return uris
	}

	// parent doesn't have an image, download own
	for i, f := range card.CardFaces {
		err := downloadImageIfNotExist(fmt.Sprintf("./resources/%s_%d.jpg", oracleID, i), f.ImageURI)
		if err != nil {
			uris[i] = card404
			continue
		}
		uris[i] = fmt.Sprintf("/resources/%s_%d.jpg", oracleID, i)
	}

	return uris
}

func downloadImageIfNotExist(dstPath, srcURI string) error {
	// Download file if it does not exist
	if _, err := os.Stat(dstPath); err != nil {
		res, err := http.Get(srcURI)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		err = ioutil.WriteFile(dstPath, buf.Bytes(), 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
