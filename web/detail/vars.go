package detail

import (
	"html/template"

	"github.com/mbndr/mtglib"
)

type detailVars struct {
	Handler *Handler
	Card    mtglib.Card
}

// CardURL returns the link to the cards image.
func (v *detailVars) CardImageURL(oracleID string) string {
	return v.Handler.cards.GetImageURL(oracleID)
}

func (v *detailVars) ManaSymbols(symbols string) template.HTML {
	return v.Handler.symbols.HTMLImages(symbols)
}

// Parse symbols in text and structure it with paragraphs
func (v *detailVars) ParseOracleText(oracleText string) template.HTML {
	return v.Handler.symbols.ParseInText(oracleText)
}
