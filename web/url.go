package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mbndr/mtglib"
)

// TODO: to scryfall package

// Get the url of an card image.
// Downloads the image if not already done.
func cardURL(oracleID string, cards map[string]mtglib.Card) string {
	card, ok := cards[oracleID]
	if !ok {
		return card404
	}

	imgPath := fmt.Sprintf("./resources/%s.jpg", oracleID)

	// Download file if it does not exist
	if _, err := os.Stat(imgPath); err != nil {
		res, err := http.Get(card.ImageURI)
		if err != nil {
			log.Println(err)
			return card404
		}
		defer res.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		err = ioutil.WriteFile(imgPath, buf.Bytes(), 0755)
		if err != nil {
			return card404
		}
	}

	return fmt.Sprintf("/resources/%s.jpg", oracleID)
}
