package detail

import (
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/mbndr/mtglib"
)

// Handler handles requests for a card
type Handler struct {
	cards   *mtglib.CardCollection
	symbols mtglib.SymbolCollection
	tpl     *template.Template
}

// NewHandler returns a new DetailHandler
func NewHandler(cards *mtglib.CardCollection, symbols mtglib.SymbolCollection) *Handler {
	return &Handler{
		cards:   cards,
		symbols: symbols,
	}
}

// ServeHTTP reloads the template and renders the HTML
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/detail/") && r.Method == "GET" {
		var err error
		h.tpl, err = template.ParseFiles("html/detail.html")
		if err != nil {
			http.Error(w, "500 internal server error\n"+err.Error(), http.StatusInternalServerError)
			return
		}
		h.serveHTML(w, r)
		return
	}

	http.Error(w, "404 page not found", http.StatusNotFound)
}

func (h *Handler) serveHTML(w http.ResponseWriter, r *http.Request) {
	oracleID := path.Base(r.URL.Path)
	card := h.cards.Get(oracleID)
	if card == nil {
		http.Error(w, "Oracle ID not in library", http.StatusNotFound)
		return
	}

	h.tpl.Execute(w, &detailVars{
		Handler: h,
		Card:    *card,
	})
}
