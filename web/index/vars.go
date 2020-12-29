package index

import "github.com/mbndr/mtglib"

// indexVars changes for each request and is given to the template.
type indexVars struct {
	Handler        *Handler
	ShownOracleIDs []string
	SearchTerm     string
}

// GetCard returns a Card object for an oracleID.
func (v *indexVars) GetCard(oracleID string) *mtglib.Card {
	return v.Handler.cards.Get(oracleID)
}

// CardURL returns the link to the cards image.
func (v *indexVars) CardURL(oracleID string) string {
	return v.Handler.cards.GetImageURL(oracleID)
}
