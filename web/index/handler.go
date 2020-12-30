package index

import (
	"math"
	"net/http"
	"text/template"

	"github.com/mbndr/mtglib"
	"github.com/mbndr/mtglib/web/index/filter"
	"github.com/mbndr/mtglib/web/index/sorting"
)

var templates = []string{
	"html/index.html",
	"html/snippets/pagination.html",
	"html/snippets/sorting.html",
	"html/snippets/advancedSearch.html",
}

// Handler wraps the data for the default http handler (only global data which does not change)
type Handler struct {
	cards             *mtglib.CardCollection
	DistinctCardCount int
	TotalCardCount    int // ALL cards (not only distinct)
	symbols           mtglib.SymbolCollection
	tpl               *template.Template
}

// NewHandler returns a new IndexHandler
func NewHandler(cards *mtglib.CardCollection, symbols mtglib.SymbolCollection, totalCardCount int) *Handler {
	return &Handler{
		cards:             cards,
		DistinctCardCount: cards.Count(),
		TotalCardCount:    totalCardCount,
		symbols:           symbols,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == "GET" {
		var err error
		h.tpl, err = template.ParseFiles(templates...)
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
		request:          r,
		Handler:          h,
		ShownOracleIDs:   resultOracleIDs[offset:idRangeTo],
		TotalResults:     len(resultOracleIDs),
		SearchTerm:       r.FormValue("search"),
		Pagination:       p,
		Sorting:          s,
		IsAdvancedSearch: len(f) != 0,
	})
}

// uses the request to build filters
func buildFilters(r *http.Request) []filter.Filter {
	var filters []filter.Filter

	// Quick name search
	quickSearch := r.URL.Query().Get("search")

	if quickSearch != "" {
		filters = append(filters, &filter.NameFilter{
			Name: quickSearch,
		})
		return filters
	}

	// Advances search filters
	ruleText := r.URL.Query().Get("rule")
	name := r.URL.Query().Get("name")
	typ := r.URL.Query().Get("type")
	rarity := r.URL.Query().Get("rarity")
	colors := r.URL.Query()["colors"]
	monocolorOnly := r.URL.Query().Get("monocolorOnly")

	if ruleText != "" {
		filters = append(filters, &filter.RuleFilter{
			Text: ruleText,
		})
	}

	if name != "" {
		filters = append(filters, &filter.NameFilter{
			Name: name,
		})
	}

	if typ != "" {
		filters = append(filters, &filter.TypeFilter{
			Type: typ,
		})
	}

	if rarity != "" {
		filters = append(filters, &filter.RarityFilter{
			Rarity: rarity,
		})
	}

	if len(colors) != 0 || monocolorOnly == "1" {
		filters = append(filters, &filter.ColorFilter{
			Colors:        colors,
			MonocolorOnly: monocolorOnly == "1",
		})
	}

	return filters
}

// get from url and validate it
func buildSorting(r *http.Request) *sorting.Sorting {
	s := sortingFromURL(r)
	if s == nil || !s.Validate() {
		return &sorting.Sorting{
			SortBy:    "name",
			SortOrder: "asc",
		}
	}
	return s
}

func sortingFromURL(r *http.Request) *sorting.Sorting {
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		return nil
	}

	sortOrder := r.URL.Query().Get("order")
	if sortOrder == "" {
		sortOrder = "asc"
	}

	return &sorting.Sorting{
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}
