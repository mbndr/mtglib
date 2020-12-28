package web

import (
	"html/template"
	"net/http"

	"github.com/mbndr/mtglib"
)

const card404 = "/static/img/card_404.jpg"

// indexVars changes for each request and is given to the template.
type indexVars struct {
	Handler        *indexHandler
	ShownOracleIDs []string
}

// GetCard returns a Card object for an oracleID.
func (v *indexVars) GetCard(oracleID string) mtglib.Card {
	return v.Handler.cards[oracleID]
}

// CardURL returns the link to the cards image.
func (v *indexVars) CardURL(oracleID string) string {
	return cardURL(oracleID, v.Handler.cards)
}

// IndexHandler wraps the data for the default http handler (only global data which does not change)
type indexHandler struct {
	cards             map[string]mtglib.Card
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
	shownOracleIDs := h.oracleIDs[80:100]

	h.tpl.Execute(w, &indexVars{
		Handler:        h,
		ShownOracleIDs: shownOracleIDs,
	})
}
