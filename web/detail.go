package web

import (
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/mbndr/mtglib"
)

type detailVars struct {
	Handler *detailHandler
	Card    mtglib.Card
}

// CardURL returns the link to the cards image.
func (v *detailVars) CardURL(oracleID string) string {
	return cardURL(oracleID, v.Handler.cards)
}

func (v *detailVars) ManaSymbols(symbols string) template.HTML {
	return v.Handler.symbols.HTMLImages(symbols)
}

// Parse symbols in text and structure it with paragraphs
func (v *detailVars) ParseOracleText(oracleText string) template.HTML {
	return v.Handler.symbols.ParseInText(oracleText)
}

// IndexHandler wraps the data for the default http handler (only global data which does not change)
type detailHandler struct {
	cards   map[string]mtglib.Card
	symbols mtglib.SymbolCollection
	tpl     *template.Template
}

func (h *detailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/detail/") && r.Method == "GET" {
		h.serveHTML(w, r)
		return
	}

	http.Error(w, "404 page not found", http.StatusNotFound)
}

func (h *detailHandler) serveHTML(w http.ResponseWriter, r *http.Request) {
	oracleID := path.Base(r.URL.Path)
	card, ok := h.cards[oracleID]
	if !ok {
		http.Error(w, "Oracle ID not in library", http.StatusNotFound)
		return
	}

	h.tpl.Execute(w, &detailVars{
		Handler: h,
		Card:    card,
	})
}
