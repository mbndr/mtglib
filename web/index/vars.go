package index

import (
	"github.com/mbndr/mtglib"
	"github.com/mbndr/mtglib/web/index/sorting"
)

// indexVars changes for each request and is given to the template.
type indexVars struct {
	Handler        *Handler
	ShownOracleIDs []string
	SearchTerm     string
	Pagination     pagination
	Sorting        *sorting.Sorting
	// TODO: filters (and get rid of SearchTerm then?)
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
