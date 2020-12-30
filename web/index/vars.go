package index

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/mbndr/mtglib"
	"github.com/mbndr/mtglib/web/index/sorting"
)

// indexVars changes for each request and is given to the template.
type indexVars struct {
	request          *http.Request
	Handler          *Handler
	ShownOracleIDs   []string
	TotalResults     int
	SearchTerm       string
	Pagination       pagination
	Sorting          *sorting.Sorting
	IsAdvancedSearch bool
}

// GetCard returns a Card object for an oracleID.
func (v *indexVars) GetCard(oracleID string) *mtglib.Card {
	return v.Handler.cards.Get(oracleID)
}

// CardURL returns the link to the cards image.
func (v *indexVars) CardURL(oracleID string) string {
	return v.Handler.cards.GetImageURL(oracleID)
}

// CardURL returns the link to the cards image.
func (v *indexVars) SortingIs(s string) bool {
	return v.Sorting != nil && v.Sorting.SortBy == s
}

// check if a url value is set and return it (string only, not list)
func (v *indexVars) QueryValue(s string) string {
	val, ok := v.request.URL.Query()[s]
	if !ok {
		return ""
	}
	return val[0]
}

func (v *indexVars) QueryValueIs(key, val string) bool {
	return v.QueryValue(key) == val
}

func (v *indexVars) QueryValueContains(key, val string) bool {
	urlValues, ok := v.request.URL.Query()[key]
	if !ok {
		return false
	}

	for _, x := range urlValues {
		if x == val {
			return true
		}
	}

	return false
}

func (v *indexVars) ManaSymbol(identifier rune) template.HTML {
	return v.Handler.symbols.HTMLImages(fmt.Sprintf("{%c}", identifier))
}
