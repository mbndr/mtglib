package web

import (
	"html/template"
	"net/http"

	"github.com/mbndr/mtglib"
)

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
	shownOracleIDs := []string{
		"1d8bd391-b133-4c8b-95ad-7f7ab99137c7",
		"4f2d4538-dc1d-4c09-964b-b0d7c240fb7d",
		"9645c1cb-2305-4bc0-89d0-a50815e91573",
	}

	h.tpl.Execute(w, &ViewVars{
		Handler:        h,
		ShownOracleIDs: shownOracleIDs,
	})
}
