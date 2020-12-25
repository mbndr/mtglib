package mtglib

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// ViewVars changes for each request
type ViewVars struct {
	Handler        *indexHandler
	ShownOracleIDs []string
}

// GetCard returns a Card object for an oracleID
func (v *ViewVars) GetCard(oracleID string) *Card {
	return v.Handler.cards[oracleID]
}

// CardURL returns the link to the cards image
// Downloads the image if not already done
func (v *ViewVars) CardURL(oracleID string) string {
	card := v.Handler.cards[oracleID]
	if card == nil {
		return "/static/img/card_404.jpg"
	}

	imgPath := fmt.Sprintf("./resources/%s.jpg", oracleID)

	// Download file if it does not exist
	if _, err := os.Stat(imgPath); err != nil {
		log.Println("Downloading image: " + card.ImageURI)

		res, err := http.Get(card.ImageURI)
		if err != nil {
			log.Println(err)
			return "/static/img/card_404.jpg"
		}
		defer res.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		err = ioutil.WriteFile(imgPath, buf.Bytes(), 0755)
		if err != nil {
			log.Println(err)
		}
	}

	return fmt.Sprintf("/resources/%s.jpg", oracleID)
}

// IndexHandler wraps the data for the default http handler (only global data which does not change)
type indexHandler struct {
	cards             map[string]*Card
	oracleIDs         []string // ALL OracleIDs
	DistinctCardCount int
	TotalCardCount    int // ALL cards (not only distinct)
	tpl               *template.Template
}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == "GET" {
		h.serveHTML(w, r)
		return
	}

	http.Error(w, "404 page not found", http.StatusNotFound)
}

func (h *indexHandler) serveHTML(w http.ResponseWriter, r *http.Request) {
	shownOracleIDs := []string{
		"1d8bd391-b133-4c8b-95ad-7f7ab99137c7",
		"4f2d4538-dc1d-4c09-964b-b0d7c240fb7d",
	}
	// TODO: filters and pretty cards

	h.tpl.Execute(w, &ViewVars{
		Handler:        h,
		ShownOracleIDs: shownOracleIDs,
	})
}

// StartServer starts the http server
func StartServer(port string) error {
	// load all cards
	cards, oracleIDs, err := LoadCards()
	if err != nil {
		return err
	}

	log.Printf("Loaded %d cards\n", len(cards))

	tpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		return err
	}

	handler := &indexHandler{
		DistinctCardCount: len(cards),
		TotalCardCount:    getTotalCardCount(),
		cards:             cards,
		oracleIDs:         oracleIDs,
		tpl:               tpl,
	}

	// setting up webserver
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", handler)

	return http.ListenAndServe(port, nil)
}

func httpError(w http.ResponseWriter, err error) {
	http.Error(w, "500 internal server error", http.StatusInternalServerError)
}
