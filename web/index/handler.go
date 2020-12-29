package index

import (
	"net/http"
	"text/template"

	"github.com/mbndr/mtglib"
	"github.com/mbndr/mtglib/web/index/filter"
)

// Handler wraps the data for the default http handler (only global data which does not change)
type Handler struct {
	cards             *mtglib.CardCollection
	DistinctCardCount int
	TotalCardCount    int // ALL cards (not only distinct)
	tpl               *template.Template
}

// NewHandler returns a new IndexHandler
func NewHandler(cards *mtglib.CardCollection, symbols mtglib.SymbolCollection, totalCardCount int) *Handler {
	return &Handler{
		cards:             cards,
		DistinctCardCount: cards.Count(),
		TotalCardCount:    totalCardCount,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == "GET" {
		var err error
		h.tpl, err = template.ParseFiles("html/index.html")
		if err != nil {
			http.Error(w, "500 internal server error\n"+err.Error(), http.StatusInternalServerError)
			return
		}
		h.serveHTML(w, r)
		return
	}

	http.Error(w, "404 page not found", http.StatusNotFound)
}

// TODO: proper limit and offset
func (h *Handler) serveHTML(w http.ResponseWriter, r *http.Request) {
	shownOracleIDs, err := filter.GetFilterResult(buildFilters(r), 0, 30)
	if err != nil {
		http.Error(w, "500 internal server error\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	h.tpl.Execute(w, &indexVars{
		Handler:        h,
		ShownOracleIDs: shownOracleIDs,
		SearchTerm:     r.FormValue("search"),
	})
}

func buildFilters(r *http.Request) []filter.Filter {
	var filters []filter.Filter

	if r.FormValue("search") != "" {
		filters = append(filters, &filter.NameFilter{
			Name: r.FormValue("search"),
		})
	}

	return filters
}
