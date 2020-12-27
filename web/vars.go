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

const card404 = "/static/img/card_404.jpg"

// ViewVars changes for each request and is given to the template.
type ViewVars struct {
	Handler        *indexHandler
	ShownOracleIDs []string
}

// GetCard returns a Card object for an oracleID.
func (v *ViewVars) GetCard(oracleID string) mtglib.Card {
	return v.Handler.cards[oracleID]
}

// CardURL returns the link to the cards image.
// Downloads the image if not already done.
func (v *ViewVars) CardURL(oracleID string) string {
	card, ok := v.Handler.cards[oracleID]
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

	return imgPath
}
