package index

import (
	"math"
	"net/http"
	"text/template"

	"github.com/mbndr/mtglib"
	"github.com/mbndr/mtglib/web/index/filter"
	"github.com/mbndr/mtglib/web/index/sorting"
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
	p := newPagination(r)
	s := buildSorting(r)
	f := buildFilters(r)

	offset := (p.CurrentPage - 1) * p.ItemsPerPage

	resultOracleIDs, err := filter.GetFilterResult(f, s)
	if err != nil {
		http.Error(w, "500 internal server error\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	p.TotalPages = int(math.Ceil(float64(len(resultOracleIDs)) / float64(p.ItemsPerPage)))

	// normalize slice range
	if offset > len(resultOracleIDs) {
		offset = len(resultOracleIDs)
	}

	idRangeTo := offset + p.ItemsPerPage
	if idRangeTo > len(resultOracleIDs) {
		idRangeTo = len(resultOracleIDs)
	}

	h.tpl.Execute(w, &indexVars{
		Handler:        h,
		ShownOracleIDs: resultOracleIDs[offset:idRangeTo],
		SearchTerm:     r.FormValue("search"),
		Pagination:     p,
		Sorting:        s,
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

// get from url and validate it
func buildSorting(r *http.Request) *sorting.Sorting {
	s := sortingFromURL(r)
	if s == nil || !s.Validate() {
		return nil
	}
	return s
}

func sortingFromURL(r *http.Request) *sorting.Sorting {
	sortBy, ok := r.URL.Query()["sort"]
	if !ok {
		return nil
	}

	sortOrder, ok := r.URL.Query()["order"]
	if !ok {
		sortOrder = []string{"asc"}
	}

	return &sorting.Sorting{
		SortBy:    sortBy[0],
		SortOrder: sortOrder[0],
	}
}
